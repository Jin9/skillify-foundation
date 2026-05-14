#!/usr/bin/env python3
"""Diff STORY-3 hand-authored baseline vs current fan-out output on 4 quality metrics.

Metrics per design brief / Scenario 4 in tests/assertions/coverage-completeness.md:
  1. AC coverage ratio       (≥ 95% required, 100% hard for banking-grade)
  2. Banking-grade coverage  (≥ 100% hard)
  3. OQ-to-test linkage      (≥ 95% required)
  4. PII / regulator coverage(≥ 95% required)

Inputs:
  - baseline/STORY-3-reference.json  (hand-authored)
  - baseline/STORY-3-current.json    (extracted from fan-out output-normalized.json)
  - The BA brief                     (ground truth)
"""

from __future__ import annotations

import json
import re
import sys
from pathlib import Path

RUN = Path(__file__).resolve().parent.parent
BRIEF = Path("/Users/IF640063/Desktop/example/qa-bootstrap-kit/e-commerce-v5/output-e5f8b9c2/output.json")
BASELINE = RUN / "baseline" / "STORY-3-reference.json"
CURRENT = RUN / "baseline" / "STORY-3-current.json"


def norm(s: str) -> str:
    """Normalize a scenario name for fuzzy comparison: lowercase, collapse whitespace and dashes."""
    s = s.lower().strip()
    s = re.sub(r"[-–—]+", " ", s)
    s = re.sub(r"\s+", " ", s)
    return s


def scenario_match(brief_name: str, plan_refs: set[str]) -> bool:
    """Fuzzy match: brief_name is covered if any plan_ref shares a normalized prefix."""
    bn = norm(brief_name)
    for p in plan_refs:
        pn = norm(p)
        if bn == pn:
            return True
        # Truncation tolerance: one is a prefix of the other after norm
        if bn.startswith(pn) or pn.startswith(bn):
            return True
        # Major-word overlap
        bw, pw = set(bn.split()), set(pn.split())
        if bw and pw and len(bw & pw) >= max(2, min(len(bw), len(pw)) // 2):
            return True
    return False


def get_story3_brief() -> dict:
    brief = json.loads(BRIEF.read_text())
    for s in brief["stories"]:
        if s["id"] == "EPIC-CART-CHECKOUT-3":
            return s
    raise RuntimeError("STORY-3 not found in brief")


def coverage_metric(brief_story: dict, plan_tcs: list[dict]) -> tuple[int, int, list[str]]:
    """AC coverage: how many BA Gherkin scenarios are covered by plan TCs."""
    brief_scenarios = [ac["scenario_name"] for ac in brief_story.get("acceptance_criteria", [])]
    plan_refs = {tc["scenario_ref"] for tc in plan_tcs}
    covered = sum(1 for s in brief_scenarios if scenario_match(s, plan_refs))
    uncovered = [s for s in brief_scenarios if not scenario_match(s, plan_refs)]
    return covered, len(brief_scenarios), uncovered


def banking_grade_metric(brief_story: dict, plan_tcs: list[dict], plan_root: dict | None = None) -> tuple[int, int, list[str]]:
    """Banking-grade coverage: every banking_grade_concerns.<concern>.status='applies' row
    has >=1 TC with matching tag. Tag conventions vary across agents; this matcher accepts
    canonical AND common short-form aliases (e.g. `banking_grade_audit` for `audit_events`).
    pii_fields + regulatory are satisfied by the presence of any test_data_specs.pii_field_refs
    / any compliance_tests, since those are the artifacts the concern represents."""
    aliases = {
        "audit_events": {"banking_grade_audit_events", "banking_grade_audit"},
        "authn_authz": {"banking_grade_authn_authz", "banking_grade_authz"},
        "idempotency": {"banking_grade_idempotency"},
        "reversibility": {"banking_grade_reversibility"},
        "pii_fields": {"banking_grade_pii_fields", "banking_grade_pii"},
        "regulatory": {"banking_grade_regulatory"},
        "tipping_off": {"banking_grade_tipping_off"},
        "notification": {"banking_grade_notification"},
    }
    concerns = brief_story.get("banking_grade_concerns", {})
    applies = [k for k, v in concerns.items() if isinstance(v, dict) and v.get("status") == "applies"]
    plan_tags = {t for tc in plan_tcs for t in tc.get("tags", [])}

    plan_has_pii_specs = bool(plan_root and any(s.get("pii_field_refs") for s in plan_root.get("test_data_specs", [])))
    plan_has_compliance = bool(plan_root and plan_root.get("compliance_tests"))

    covered_set = set()
    for c in applies:
        # Tag-based coverage
        if any(t in plan_tags for t in aliases.get(c, {f"banking_grade_{c}"})):
            covered_set.add(c)
            continue
        # Artifact-based coverage for cross-cutting concerns
        if c == "pii_fields" and plan_has_pii_specs:
            covered_set.add(c)
            continue
        if c == "regulatory" and plan_has_compliance:
            covered_set.add(c)
            continue
    uncovered = [c for c in applies if c not in covered_set]
    return len(covered_set), len(applies), uncovered


def oq_linkage_metric(brief: dict, plan_tcs: list[dict], plan_nfrs: list[dict], plan_root: dict | None = None) -> tuple[int, int]:
    """OQ linkage: count of distinct OQ references across STORY-3 TCs + NFRs + coverage_gaps + tl_design_dependencies.
    Different agents put OQ refs in different fields; this rolls up all of them."""
    oq_refs = set()
    for tc in plan_tcs:
        oq_refs.update(tc.get("blocking_oqs", []) or [])
    for n in plan_nfrs:
        oq_refs.update(n.get("blocking_oqs", []) or [])
    if plan_root:
        for g in plan_root.get("coverage_gaps", []) or []:
            # Sweep description for OQ-NN patterns
            for token in (g.get("description", "") + " " + g.get("required_action", "")).split():
                if token.startswith("OQ-") and len(token) >= 5:
                    oq_refs.add(token.strip(".,;:()"))
        for d in plan_root.get("tl_design_dependencies", []) or []:
            for token in d.get("description", "").split():
                if token.startswith("OQ-") and len(token) >= 5:
                    oq_refs.add(token.strip(".,;:()"))
        # Top-level "blocking_oqs" list (some agents emit this — may be strings or objects)
        for item in plan_root.get("blocking_oqs", []) or []:
            if isinstance(item, str):
                oq_refs.add(item)
            elif isinstance(item, dict):
                for key in ("oq_id", "id", "ref"):
                    if key in item:
                        oq_refs.add(item[key])
                        break
    return len(oq_refs), len(brief.get("open_questions", []))


def pii_regulator_metric(brief: dict, plan: dict) -> tuple[int, int]:
    """PII + regulator coverage on STORY-3's scope.
    STORY-3 touches checkout PII (email, phone, name, shipping_address, payment-related PII).
    Compliance: PDPA, TRC, CPA all apply.
    We measure: of BA's 4 regulators, how many are referenced in STORY-3 fragment's
    compliance_tests or compliance_tags? Of BA's 10 PII fields, how many are referenced
    in fixtures or assertions for STORY-3?"""
    brief_regs = {r["code"].replace("-", "_") for r in brief.get("regulatory_dependencies", [])}
    brief_pii = {p["field"] for p in brief.get("pii_inventory", [])}

    plan_regs = set()
    for c in plan.get("compliance_tests", []):
        plan_regs.add(c.get("regulator_code", ""))
    for tc in plan.get("test_cases", []):
        plan_regs.update(tc.get("compliance_tags", []) or [])

    plan_pii = set()
    for spec in plan.get("test_data_specs", []):
        plan_pii.update(spec.get("pii_field_refs", []) or [])

    regs_covered = len(plan_regs & brief_regs)
    pii_covered = len(plan_pii & brief_pii)
    return (regs_covered + pii_covered), (len(brief_regs) + len(brief_pii))


def evaluate(label: str, brief: dict, brief_story: dict, plan_root: dict) -> dict:
    plan_tcs = plan_root.get("test_cases", [])
    plan_nfrs = plan_root.get("nfr_tests", []) or plan_root.get("nfr_tests_for_epic", [])

    ac_cov, ac_total, ac_miss = coverage_metric(brief_story, plan_tcs)
    bg_cov, bg_total, bg_miss = banking_grade_metric(brief_story, plan_tcs, plan_root)
    oq_count, oq_total = oq_linkage_metric(brief, plan_tcs, plan_nfrs, plan_root)
    pii_reg, pii_reg_total = pii_regulator_metric(brief, plan_root)

    result = {
        "label": label,
        "tc_count": len(plan_tcs),
        "nfr_count": len(plan_nfrs),
        "ac_coverage": f"{ac_cov}/{ac_total} = {ac_cov/ac_total*100:.0f}%",
        "ac_uncovered": ac_miss,
        "banking_grade_coverage": f"{bg_cov}/{bg_total} = {bg_cov/bg_total*100 if bg_total else 100:.0f}%",
        "banking_grade_uncovered": bg_miss,
        "oq_linkage_count": oq_count,
        "pii_regulator_coverage": f"{pii_reg}/{pii_reg_total} = {pii_reg/pii_reg_total*100:.0f}%",
    }
    return result


def main() -> int:
    if not BASELINE.exists():
        print(f"error: baseline not found yet: {BASELINE}", file=sys.stderr)
        return 2
    if not CURRENT.exists():
        print(f"error: current extract not found: {CURRENT}", file=sys.stderr)
        return 2

    brief = json.loads(BRIEF.read_text())
    brief_story = get_story3_brief()
    baseline = json.loads(BASELINE.read_text())
    current = json.loads(CURRENT.read_text())

    print("=" * 70)
    print("STORY-3 QUALITY DIFF — hand-authored baseline vs fan-out current")
    print("=" * 70)
    print()
    print(f"Brief STORY-3: id={brief_story['id']}, ACs={len(brief_story.get('acceptance_criteria', []))}")
    bg_applies = [k for k, v in brief_story.get("banking_grade_concerns", {}).items() if isinstance(v, dict) and v.get("status") == "applies"]
    print(f"               banking_grade_applies: {bg_applies}")
    print()

    base = evaluate("BASELINE (hand-authored)", brief, brief_story, baseline)
    curr = evaluate("CURRENT  (fan-out + normalized)", brief, brief_story, current)

    print(f"{'METRIC':<32}  {'BASELINE':>22}  {'CURRENT':>22}")
    print("-" * 80)
    for k in ["tc_count", "nfr_count", "ac_coverage", "banking_grade_coverage", "oq_linkage_count", "pii_regulator_coverage"]:
        b_val = str(base.get(k, ""))[:22]
        c_val = str(curr.get(k, ""))[:22]
        print(f"  {k:<30}  {b_val:>22}  {c_val:>22}")

    print()
    if base["ac_uncovered"]:
        print(f"BASELINE ac_uncovered: {base['ac_uncovered']}")
    if curr["ac_uncovered"]:
        print(f"CURRENT  ac_uncovered: {curr['ac_uncovered']}")
    if base["banking_grade_uncovered"]:
        print(f"BASELINE banking_grade_uncovered: {base['banking_grade_uncovered']}")
    if curr["banking_grade_uncovered"]:
        print(f"CURRENT  banking_grade_uncovered: {curr['banking_grade_uncovered']}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
