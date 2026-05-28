#!/usr/bin/env python3
"""
validate_ba.py — structural check on a BA spec (ba.json).

Refuses fan-out on a malformed BA artifact. Catches:
  - missing required top-level fields per docs/templates/ba.md schema
  - non-semver template_version
  - empty in_scope / acceptance_criteria
  - acceptance_criteria items missing given/when/then
  - unfalsifiable `then` clauses (heuristic: "is good", "works", "looks correct")
  - empty edge_cases (the documented smell from ba.md negative example #1)

Usage:
  python3 validate_ba.py <path-to-ba.json>

Exit 0 = pass; non-zero = structural failure with diagnostic on stderr.
"""
from __future__ import annotations

import json
import re
import sys
from pathlib import Path

REQUIRED_TOP = {
    "template_version", "workflow_name", "purpose",
    "in_scope", "out_of_scope", "acceptance_criteria",
    "edge_cases", "success_conditions", "compliance_sensitive",
}
ALLOWED_TOP = REQUIRED_TOP | {"non_functional"}
SEMVER = re.compile(r"^\d+\.\d+\.\d+$")
UNFALSIFIABLE = re.compile(
    r"\b(is good|works correctly|looks correct|is fine|is right|properly handled)\b",
    re.IGNORECASE,
)


def fail(msg: str) -> None:
    print(f"validate_ba: FAIL {msg}", file=sys.stderr)
    sys.exit(2)


def main(path: str) -> None:
    p = Path(path)
    if not p.exists():
        fail(f"file not found: {path}")
    try:
        d = json.loads(p.read_text())
    except json.JSONDecodeError as e:
        fail(f"not valid JSON: {e}")

    extras = set(d.keys()) - ALLOWED_TOP
    if extras:
        fail(f"unexpected top-level keys: {sorted(extras)}")
    missing = REQUIRED_TOP - set(d.keys())
    if missing:
        fail(f"missing required top-level keys: {sorted(missing)}")

    if not SEMVER.match(d["template_version"]):
        fail(f"template_version not semver: {d['template_version']!r}")
    if not isinstance(d["compliance_sensitive"], bool):
        fail("compliance_sensitive must be a boolean")
    if len(d["purpose"]) > 200:
        fail(f"purpose exceeds 200 chars: {len(d['purpose'])}")

    if not isinstance(d["in_scope"], list) or not d["in_scope"]:
        fail("in_scope must be a non-empty array")
    if not isinstance(d["acceptance_criteria"], list) or not d["acceptance_criteria"]:
        fail("acceptance_criteria must be a non-empty array")

    for i, ac in enumerate(d["acceptance_criteria"]):
        if not isinstance(ac, dict) or set(ac.keys()) != {"given", "when", "then"}:
            fail(f"acceptance_criteria[{i}] must be {{given, when, then}}")
        for k, v in ac.items():
            if not isinstance(v, str) or not v.strip():
                fail(f"acceptance_criteria[{i}].{k} must be non-empty string")
        if UNFALSIFIABLE.search(ac["then"]):
            fail(
                f"acceptance_criteria[{i}].then is unfalsifiable: {ac['then']!r}"
                " — replace with a concrete checkable assertion"
            )

    if not isinstance(d["edge_cases"], list):
        fail("edge_cases must be an array")
    if not d["edge_cases"]:
        fail(
            "edge_cases is empty — this is the documented smell from "
            "docs/templates/ba.md negative example #1"
        )
    for i, ec in enumerate(d["edge_cases"]):
        if not isinstance(ec, dict) or set(ec.keys()) != {"case", "expected"}:
            fail(f"edge_cases[{i}] must be {{case, expected}}")

    print(f"validate_ba: PASS {path} "
          f"(in_scope={len(d['in_scope'])}, "
          f"AC={len(d['acceptance_criteria'])}, "
          f"edges={len(d['edge_cases'])})")


if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("usage: validate_ba.py <path-to-ba.json>", file=sys.stderr)
        sys.exit(64)
    main(sys.argv[1])
