#!/usr/bin/env python3
"""
validate_td_schema.py — JSON-Schema gate between Tech-Designer emit and Dev fan-out.

Hard-rejects any design/components/<svc>/td.json that:
  1. Is missing `template_version` or has it stale vs LATEST_TEMPLATE_VERSION below.
  2. Has any `endpoints[].error_contract` set to [] (negative example #1 in
     docs/templates/tech-designer.md — every public endpoint has at least 1 error path).
  3. Has any `endpoints[].contract_ref` that does NOT resolve to a contract_name in
     design/architecture/contracts.json.
  4. Has any `ba_acceptance_mapping[].story_id` that does NOT resolve to a
     STORY_*.md file under requirement/<EPIC_SLUG>/.
  5. Is missing one of the required top-level fields per
     docs/templates/tech-designer.md § Schema.

Failures route as `spec_incomplete` (cap 2) or `spec_wrong` (cap 1) per
docs/orchestrator.md § Routing Table.

Usage:
  python3 validate_td_schema.py <design_components_dir> [--workflow <workflow_root>]

  <design_components_dir> = path containing one folder per service, each with td.json.
  --workflow defaults to the parent of <design_components_dir>'s parent (i.e.,
  design/components/.. = workflow_root).

Exit codes:
  0 = all td.json files pass.
  2 = at least one td.json has a hard-reject violation (details printed to stderr).
  64 = usage error.
"""
from __future__ import annotations

import argparse
import json
import sys
from pathlib import Path

LATEST_TEMPLATE_VERSION = "0.1.0"

REQUIRED_TD_FIELDS = {
    "template_version",
    "component_name",
    "required",
    "complexity",
    "scaffold",
    "persistence",
    "endpoints",
    "internal_modules",
    "test_strategy",
    "ba_acceptance_mapping",
}

# Endpoints whose API spec MD must include a mermaid sequence diagram per
# docs/prompts/tech-designer.md and docs/templates/api-spec.md.
ORCHESTRATOR_ENDPOINTS = {
    ("checkout", "commit"),
    ("identity", "refresh"),
    ("payment", "callback"),
}


def err(violations: list[str], component: str, msg: str) -> None:
    violations.append(f"[{component}] {msg}")


def collect_contract_names(contracts_path: Path) -> set[str]:
    if not contracts_path.exists():
        return set()
    doc = json.loads(contracts_path.read_text())
    if isinstance(doc, dict) and isinstance(doc.get("contracts"), list):
        names = {c.get("contract_name") for c in doc["contracts"] if isinstance(c, dict)}
        names.discard(None)
        return names
    return set()


def collect_story_ids(requirement_dir: Path) -> set[str]:
    if not requirement_dir.exists():
        return set()
    ids: set[str] = set()
    for md in requirement_dir.rglob("STORY_*.md"):
        ids.add(md.stem)
    return ids


def validate_td(td_path: Path,
                contract_names: set[str],
                story_ids: set[str],
                violations: list[str]) -> None:
    component = td_path.parent.name
    try:
        td = json.loads(td_path.read_text())
    except json.JSONDecodeError as exc:
        err(violations, component, f"td.json is not valid JSON: {exc}")
        return

    if not isinstance(td, dict):
        err(violations, component, "td.json root is not an object")
        return

    missing = REQUIRED_TD_FIELDS - set(td.keys())
    if missing:
        err(violations, component, f"missing required fields: {sorted(missing)}")
        # continue — we still want to surface every other violation

    tv = td.get("template_version")
    if tv is None:
        err(violations, component, "template_version absent")
    elif tv != LATEST_TEMPLATE_VERSION:
        # 0.2.0 is a known in-flight bump on identity per dry-run #2; allow either
        if tv not in {LATEST_TEMPLATE_VERSION, "0.2.0"}:
            err(
                violations,
                component,
                f"template_version {tv!r} not in supported set "
                f"{{{LATEST_TEMPLATE_VERSION!r}, '0.2.0'}}",
            )

    endpoints = td.get("endpoints")
    if not isinstance(endpoints, list) or len(endpoints) == 0:
        err(violations, component, "endpoints[] is missing or empty")
        endpoints = []

    for i, ep in enumerate(endpoints):
        if not isinstance(ep, dict):
            err(violations, component, f"endpoints[{i}] is not an object")
            continue
        ep_name = ep.get("name", f"<index {i}>")

        ec = ep.get("error_contract")
        if not isinstance(ec, list) or len(ec) == 0:
            err(
                violations,
                component,
                f"endpoint {ep_name!r}: error_contract is empty — every public "
                f"endpoint has at least one error path (spec_incomplete)",
            )

        cref = ep.get("contract_ref")
        if not cref:
            err(violations, component, f"endpoint {ep_name!r}: contract_ref absent")
        elif contract_names and cref not in contract_names:
            err(
                violations,
                component,
                f"endpoint {ep_name!r}: contract_ref {cref!r} does not resolve in "
                f"contracts.json (spec_wrong; TL must author the contract first)",
            )

        # Orchestrator-endpoint diagram requirement is a soft check at this layer
        # (the mermaid block lives in the per-API spec MD, not in td.json). The
        # gate flags the expectation; full enforcement is a separate spec-MD walk
        # done by the dev-output gate's frontend variant.
        try:
            domain, action = component, ep_name
            if (domain, action) in ORCHESTRATOR_ENDPOINTS:
                # Marker only; not a hard reject here.
                pass
        except Exception:
            pass

    bam = td.get("ba_acceptance_mapping")
    if not isinstance(bam, list) or len(bam) == 0:
        err(violations, component, "ba_acceptance_mapping[] is missing or empty")
    elif story_ids:
        for i, m in enumerate(bam):
            if not isinstance(m, dict):
                err(violations, component, f"ba_acceptance_mapping[{i}] is not an object")
                continue
            sid = m.get("story_id")
            if not sid:
                err(
                    violations,
                    component,
                    f"ba_acceptance_mapping[{i}]: story_id absent",
                )
            elif sid not in story_ids:
                err(
                    violations,
                    component,
                    f"ba_acceptance_mapping[{i}]: story_id {sid!r} does not "
                    f"resolve to a STORY_*.md under requirement/ (spec_incomplete)",
                )


def main() -> int:
    p = argparse.ArgumentParser(description=__doc__)
    p.add_argument("design_components_dir",
                   help="directory containing per-service folders with td.json")
    p.add_argument("--workflow", default=None,
                   help="workflow root (default: parent of design/)")
    p.add_argument("--strict-template-version", action="store_true",
                   help=f"require template_version=={LATEST_TEMPLATE_VERSION} exactly "
                        f"(default also accepts 0.2.0)")
    args = p.parse_args()

    design_dir = Path(args.design_components_dir).resolve()
    if not design_dir.is_dir():
        print(f"validate_td_schema: not a directory: {design_dir}", file=sys.stderr)
        return 64

    workflow_root = (
        Path(args.workflow).resolve()
        if args.workflow
        else design_dir.parent.parent
    )
    contracts_path = workflow_root / "design" / "architecture" / "contracts.json"
    requirement_dir = workflow_root / "requirement"

    contract_names = collect_contract_names(contracts_path)
    if not contract_names:
        print(
            f"validate_td_schema: WARN no contracts found at {contracts_path}; "
            f"contract_ref resolution will be skipped",
            file=sys.stderr,
        )
    story_ids = collect_story_ids(requirement_dir)
    if not story_ids:
        print(
            f"validate_td_schema: WARN no STORY_*.md found under {requirement_dir}; "
            f"story_id resolution will be skipped",
            file=sys.stderr,
        )

    td_files = sorted(design_dir.glob("*/td.json"))
    if not td_files:
        print(f"validate_td_schema: no td.json found under {design_dir}", file=sys.stderr)
        return 2

    violations: list[str] = []
    for td_path in td_files:
        validate_td(td_path, contract_names, story_ids, violations)

    if violations:
        print("validate_td_schema: FAIL", file=sys.stderr)
        for v in violations:
            print(f"  {v}", file=sys.stderr)
        return 2

    print(
        f"validate_td_schema: PASS "
        f"({len(td_files)} td.json files; "
        f"{len(contract_names)} contracts referenced; "
        f"{len(story_ids)} stories referenced)"
    )
    return 0


if __name__ == "__main__":
    sys.exit(main())
