#!/usr/bin/env python3
"""
validate_components.py — orphan-dep check on Tech-Lead's component list.

Every `dependencies` ref in components.json must resolve to a `contract_name`
in contracts.json. This is the structural check that would have caught
CHK-AMBIG-001 in dry-run #1 (Checkout TD referenced
`order.create-from-checkout` and `order.cancel-on-checkout-failure` before TL
authored them as contracts).

Also verifies:
  - Each component has the required fields per docs/templates/tech-lead-components.md.
  - complexity is in {simple, standard, complex}.
  - complexity_rationale is present when complexity != "standard".

Usage:
  python3 validate_components.py <contracts.json> <components.json>

Exit 0 = pass; non-zero = orphan dep or schema break.
"""
from __future__ import annotations

import json
import re
import sys
from pathlib import Path


REQUIRED_COMP_FIELDS = {
    "component_name", "required", "complexity", "dependencies", "parallel_safe",
}
ALLOWED_COMPLEXITY = {"simple", "standard", "complex"}
COMP_NAME = re.compile(r"^[a-z][a-z0-9-]*$")


def fail(msg: str) -> None:
    print(f"validate_components: FAIL {msg}", file=sys.stderr)
    sys.exit(2)


def main(contracts_path: str, components_path: str) -> None:
    cdoc = json.loads(Path(contracts_path).read_text())
    codoc = json.loads(Path(components_path).read_text())

    if "contracts" not in cdoc or not isinstance(cdoc["contracts"], list):
        fail(f"contracts.json missing 'contracts' array")
    if "components" not in codoc or not isinstance(codoc["components"], list):
        fail(f"components.json missing 'components' array")

    contract_names = {c.get("contract_name") for c in cdoc["contracts"]}
    contract_names.discard(None)
    if not contract_names:
        fail("contracts.json has no contract_name entries")

    seen: set[str] = set()
    orphan_log: list[tuple[str, str]] = []

    for i, comp in enumerate(codoc["components"]):
        name = comp.get("component_name", f"<index {i}>")

        missing = REQUIRED_COMP_FIELDS - set(comp.keys())
        if missing:
            fail(f"component {name!r} missing required fields: {sorted(missing)}")

        if not COMP_NAME.match(comp["component_name"]):
            fail(f"component_name not lowercase-kebab: {comp['component_name']!r}")
        if comp["component_name"] in seen:
            fail(f"duplicate component_name: {comp['component_name']!r}")
        seen.add(comp["component_name"])

        if comp["complexity"] not in ALLOWED_COMPLEXITY:
            fail(
                f"component {name!r}: complexity {comp['complexity']!r} not in "
                f"{sorted(ALLOWED_COMPLEXITY)}"
            )
        if comp["complexity"] != "standard" and not comp.get("complexity_rationale"):
            fail(
                f"component {name!r}: complexity={comp['complexity']!r} "
                f"requires complexity_rationale"
            )
        if not isinstance(comp["required"], bool):
            fail(f"component {name!r}: required must be boolean")
        if not isinstance(comp["parallel_safe"], bool):
            fail(f"component {name!r}: parallel_safe must be boolean")
        if not isinstance(comp["dependencies"], list):
            fail(f"component {name!r}: dependencies must be array")

        for dep in comp["dependencies"]:
            if dep not in contract_names:
                orphan_log.append((comp["component_name"], dep))

    if orphan_log:
        details = "\n  ".join(f"{c} -> {d}" for c, d in orphan_log)
        fail(
            f"orphan dependencies (component → missing contract_name):\n  "
            f"{details}\n"
            f"Tech-Lead must author these contracts in contracts.json before fan-out."
        )

    print(
        f"validate_components: PASS "
        f"({len(codoc['components'])} components, "
        f"{len(cdoc['contracts'])} contracts, 0 orphan deps)"
    )


if __name__ == "__main__":
    if len(sys.argv) != 3:
        print(
            "usage: validate_components.py <contracts.json> <components.json>",
            file=sys.stderr,
        )
        sys.exit(64)
    main(sys.argv[1], sys.argv[2])
