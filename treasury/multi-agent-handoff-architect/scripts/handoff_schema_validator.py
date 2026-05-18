#!/usr/bin/env python3
"""multi-agent-handoff-architect: handoff-contract schema self-check.

Validates the skill's OWN output (the JSON Schema contract) before it is
delivered. A handoff payload is a contract, not a note: the schema PASSES
only if it

  1. declares all six required fields in its `required` array:
     taskId, intent, state, confidence, provenance, schemaVersion;
  2. sets `additionalProperties: false` at the top level;
  3. contains NO property whose name matches /note|narrative|freetext|
     free_text|prose/i anywhere in the schema (free-text handoffs are
     the dominant source of context loss and are forbidden by design).

No network, no LLM, no code execution — stdlib only (json). Exit 0 = ok,
1 = schema fails one or more rules, 2 = usage / parse error.

Usage:
  python3 handoff_schema_validator.py schemas/handoff-contract.schema.json
"""
import argparse
import json
import re
import sys
from pathlib import Path

REQUIRED_FIELDS = (
    "taskId",
    "intent",
    "state",
    "confidence",
    "provenance",
    "schemaVersion",
)
FREETEXT_RE = re.compile(r"note|narrative|freetext|free_text|prose", re.I)


def collect_property_names(node):
    """Recursively collect every property name declared anywhere."""
    names = []
    if isinstance(node, dict):
        props = node.get("properties")
        if isinstance(props, dict):
            names.extend(props.keys())
        for value in node.values():
            names.extend(collect_property_names(value))
    elif isinstance(node, list):
        for item in node:
            names.extend(collect_property_names(item))
    return names


def main():
    ap = argparse.ArgumentParser(
        description="Self-check a handoff-contract JSON Schema."
    )
    ap.add_argument("schema", help="path to the JSON Schema file")
    a = ap.parse_args()

    schema_path = Path(a.schema)
    if not schema_path.is_file():
        print(f"error: schema not found: {a.schema}", file=sys.stderr)
        sys.exit(2)

    try:
        schema = json.loads(schema_path.read_text(encoding="utf-8"))
    except json.JSONDecodeError as e:
        print(f"error: schema is not valid JSON: {e}", file=sys.stderr)
        sys.exit(2)

    if not isinstance(schema, dict):
        print("error: schema root must be a JSON object", file=sys.stderr)
        sys.exit(2)

    problems = []

    # 1. all six required fields declared in top-level `required`
    required = schema.get("required")
    if not isinstance(required, list):
        problems.append("top-level `required` array is missing or not a list")
        required = []
    missing = [f for f in REQUIRED_FIELDS if f not in required]
    if missing:
        problems.append(
            "required array missing fields: " + ", ".join(missing)
        )

    # 2. additionalProperties: false at top level
    if schema.get("additionalProperties") is not False:
        problems.append(
            "top-level `additionalProperties` must be false "
            "(open contracts admit free-text drift)"
        )

    # 3. no free-text-style property anywhere
    bad = sorted(
        {n for n in collect_property_names(schema) if FREETEXT_RE.search(n)}
    )
    if bad:
        problems.append(
            "forbidden free-text property name(s): " + ", ".join(bad)
        )

    ok = not problems
    print(f"required fields expected: {', '.join(REQUIRED_FIELDS)}")
    if problems:
        print(f"FAIL ({len(problems)}):")
        for p in problems:
            print(f"  - {p}")
        print("FAIL — fix the above before delivering the schema.")
        sys.exit(1)

    print(
        "PASS — six required fields declared, additionalProperties:false, "
        "no free-text property. Handoff payload is a contract, not a note."
    )
    sys.exit(0)


if __name__ == "__main__":
    main()
