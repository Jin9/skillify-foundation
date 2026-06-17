#!/usr/bin/env python3
"""Deterministic ENVELOPE structural check for an extract-anything output.

Validates the loss-bounding envelope { contract, _meta{...} } only. It does NOT
validate the inner `contract` against a downstream target_schema — full schema
conformance is a different job (see universal-spec-validator). This is the
post-condition integrity gate the skill runs before emitting.

Usage:
    validate_extraction.py <extraction.json>

Exit: 0 = envelope valid · 1 = invalid or usage/parse error.
Reads stdlib only.
"""

from __future__ import annotations

import json
import sys

COVERAGE = {"high", "partial", "low"}
CONFIDENCE = {"high", "medium", "low"}
MODE = {"conform", "infer"}
META_REQUIRED = ("mode", "source_ref", "coverage", "confidence", "dropped")
META_ALLOWED = set(META_REQUIRED) | {"provenance", "idempotency_key", "notes"}


def check(envelope: object) -> list[str]:
    errors: list[str] = []

    if not isinstance(envelope, dict):
        return ["top level must be a JSON object"]

    for key in ("contract", "_meta"):
        if key not in envelope:
            errors.append(f"missing top-level key: {key}")
    for key in envelope:
        if key not in ("contract", "_meta"):
            errors.append(f"unexpected top-level key: {key}")

    if "contract" in envelope and not isinstance(envelope["contract"], (dict, list)):
        errors.append("contract must be an object or array")

    meta = envelope.get("_meta")
    if not isinstance(meta, dict):
        errors.append("_meta must be an object")
        return errors

    for key in META_REQUIRED:
        if key not in meta:
            errors.append(f"_meta missing required key: {key}")
    for key in meta:
        if key not in META_ALLOWED:
            errors.append(f"_meta has unexpected key: {key}")

    if meta.get("mode") not in MODE:
        errors.append(f"_meta.mode must be one of {sorted(MODE)}")
    if meta.get("coverage") not in COVERAGE:
        errors.append(f"_meta.coverage must be one of {sorted(COVERAGE)}")
    if meta.get("confidence") not in CONFIDENCE:
        errors.append(f"_meta.confidence must be one of {sorted(CONFIDENCE)}")

    ref = meta.get("source_ref")
    if not isinstance(ref, str) or not ref.strip():
        errors.append("_meta.source_ref must be a non-empty string")

    dropped = meta.get("dropped")
    if not isinstance(dropped, list):
        errors.append("_meta.dropped must be an array (empty is allowed)")
    else:
        for i, entry in enumerate(dropped):
            if not isinstance(entry, dict) or set(entry) != {"item", "reason"}:
                errors.append(f"_meta.dropped[{i}] must be an object with exactly item, reason")
                continue
            for field in ("item", "reason"):
                if not isinstance(entry[field], str) or not entry[field].strip():
                    errors.append(f"_meta.dropped[{i}].{field} must be a non-empty string")

    prov = meta.get("provenance")
    if prov is not None:
        if not isinstance(prov, list):
            errors.append("_meta.provenance must be an array when present")
        else:
            for i, entry in enumerate(prov):
                if not isinstance(entry, dict) or set(entry) != {"field", "from"}:
                    errors.append(f"_meta.provenance[{i}] must be an object with exactly field, from")
                    continue
                for field in ("field", "from"):
                    if not isinstance(entry[field], str) or not entry[field].strip():
                        errors.append(f"_meta.provenance[{i}].{field} must be a non-empty string")

    for key in ("idempotency_key", "notes"):
        if key in meta and not isinstance(meta[key], str):
            errors.append(f"_meta.{key} must be a string when present")

    return errors


def main(argv: list[str]) -> int:
    if len(argv) != 2:
        print("usage: validate_extraction.py <extraction.json>", file=sys.stderr)
        return 1
    try:
        with open(argv[1], encoding="utf-8") as handle:
            envelope = json.load(handle)
    except OSError as exc:
        print(f"error: cannot read {argv[1]}: {exc}", file=sys.stderr)
        return 1
    except json.JSONDecodeError as exc:
        print(f"error: invalid JSON in {argv[1]}: {exc}", file=sys.stderr)
        return 1

    errors = check(envelope)
    if errors:
        print(f"INVALID envelope: {argv[1]}", file=sys.stderr)
        for error in errors:
            print(f"  - {error}", file=sys.stderr)
        return 1
    print(f"PASS envelope: {argv[1]}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main(sys.argv))
