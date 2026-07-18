#!/usr/bin/env python3
"""Deterministic schema validation for routine definition files.

Usage: python3 validate_routine.py ROUTINE.md [ROUTINE.md ...]

Prints `error:` / `warning:` lines to stderr and `ok: PATH` per passing
file. Exit 0 when every file passes; exit 1 when any file has at least
one error. Warnings never affect the exit code.
"""

from __future__ import annotations

import argparse
import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent))

from routine_parser import (
    FRONTMATTER_OPTIONAL,
    FRONTMATTER_REQUIRED,
    GATES,
    KEBAB_RE,
    NO_INPUT,
    NODE_OPTIONAL,
    NODE_REQUIRED,
    ON_FAIL,
    TIERS,
    parse_routine,
)

LAUNCHER_NAME = "launching-workflow-routines"
MAX_NODES = 12
MAX_SUMMARY = 200
MAX_NAME = 64


def validate_file(path: Path, sibling_stems: set) -> tuple:
    errors = []
    warnings = []
    routine = parse_routine(path)
    errors.extend(routine.errors)

    fm = routine.frontmatter
    for key in FRONTMATTER_REQUIRED:
        if key not in fm:
            errors.append(f"frontmatter missing required field: {key}")
    allowed_fm = set(FRONTMATTER_REQUIRED) | set(FRONTMATTER_OPTIONAL)
    for key in fm:
        if key not in allowed_fm:
            errors.append(f"frontmatter unknown field: {key}")

    name = fm.get("routine", "")
    if name:
        if not KEBAB_RE.fullmatch(name):
            errors.append(f"routine must be lowercase kebab-case: {name}")
        elif len(name) > MAX_NAME:
            errors.append(f"routine name too long: {len(name)} > {MAX_NAME}")
        if name != path.stem:
            errors.append(f"routine must match the filename stem: {name} != {path.stem}")

    summary = fm.get("summary", "")
    if summary:
        if len(summary) > MAX_SUMMARY:
            errors.append(f"summary too long: {len(summary)} > {MAX_SUMMARY}")
        if "<" in summary or ">" in summary:
            errors.append("summary must not contain angle brackets")

    default_on_fail = fm.get("default_on_fail")
    if default_on_fail is not None and default_on_fail not in ON_FAIL:
        errors.append(f"default_on_fail must be one of {'|'.join(ON_FAIL)}: {default_on_fail}")

    if routine.nodes_heading_count != 1:
        errors.append(f"file must contain exactly one '## Nodes' section (found {routine.nodes_heading_count})")
    if not routine.nodes:
        errors.append("routine must declare at least one '### Node:' block")
    elif len(routine.nodes) > MAX_NODES:
        errors.append(f"too many nodes: {len(routine.nodes)} > {MAX_NODES}")

    seen_ids = set()
    declared_outputs = {}
    allowed_node = set(NODE_REQUIRED) | set(NODE_OPTIONAL)
    for position, node in enumerate(routine.nodes, start=1):
        label = f"node '{node.node_id}'"
        if not KEBAB_RE.fullmatch(node.node_id):
            errors.append(f"{label}: id must be lowercase kebab-case")
        if node.node_id in seen_ids:
            errors.append(f"{label}: duplicate node id")
        seen_ids.add(node.node_id)

        for key in NODE_REQUIRED:
            if key not in node.fields:
                errors.append(f"{label}: missing required key: {key}")
        for key in node.fields:
            if key not in allowed_node:
                errors.append(f"{label}: unknown key: {key}")

        tier = node.fields.get("tier")
        if tier is not None and tier not in TIERS:
            errors.append(f"{label}: tier must be one of {'|'.join(TIERS)}: {tier}")
        gate = node.fields.get("gate")
        if gate is not None and gate not in GATES:
            errors.append(f"{label}: gate must be one of {'|'.join(GATES)}: {gate}")
        on_fail = node.fields.get("on_fail")
        if on_fail is not None and on_fail not in ON_FAIL:
            errors.append(f"{label}: on_fail must be one of {'|'.join(ON_FAIL)}: {on_fail}")

        executor = node.fields.get("executor", "")
        if executor:
            if executor == LAUNCHER_NAME:
                errors.append(f"{label}: recursion - executor must never be {LAUNCHER_NAME}")
            elif executor.endswith(".md"):
                errors.append(f"{label}: executor must be a skill name, not a file: {executor}")
            elif executor in sibling_stems:
                errors.append(f"{label}: recursion - executor matches sibling routine '{executor}'")

        for token in node.inputs():
            if token == NO_INPUT or token.startswith("external:"):
                continue
            if token.startswith("user."):
                input_key = token[len("user.") :]
                if input_key not in routine.inputs:
                    errors.append(f"{label}: input '{token}' has no entry in '## Inputs'")
            elif token not in declared_outputs:
                errors.append(f"{label}: input artifact '{token}' is not an output of an earlier node")

        outputs = node.outputs()
        if "outputs" in node.fields and not outputs:
            errors.append(f"{label}: outputs must declare at least one artifact path")
        for out in outputs:
            if out.startswith("/") or ".." in out.split("/"):
                errors.append(f"{label}: output must be run-dir-relative without '..': {out}")
            if out in declared_outputs:
                errors.append(f"{label}: output duplicated (also declared by node '{declared_outputs[out]}'): {out}")
            else:
                declared_outputs[out] = node.node_id
        if outputs:
            prefix = f"{position:02d}-"
            if not outputs[0].split("/")[-1].startswith(prefix):
                warnings.append(f"{label}: first output does not carry the '{prefix}' sequence prefix: {outputs[0]}")

    return errors, warnings


def main() -> int:
    parser = argparse.ArgumentParser(description="Validate routine definition files.")
    parser.add_argument("files", nargs="+", help="Routine .md files to validate")
    args = parser.parse_args()

    failed = False
    for raw in args.files:
        path = Path(raw).expanduser().resolve()
        if not path.is_file():
            print(f"error: {raw}: not a file", file=sys.stderr)
            failed = True
            continue
        sibling_stems = {
            candidate.stem
            for candidate in path.parent.glob("*.md")
            if candidate != path and candidate.name != "INDEX.md"
        }
        errors, warnings = validate_file(path, sibling_stems)
        for warning in warnings:
            print(f"warning: {path.name}: {warning}", file=sys.stderr)
        if errors:
            failed = True
            for error in errors:
                print(f"error: {path.name}: {error}", file=sys.stderr)
        else:
            print(f"ok: {path}")
    return 1 if failed else 0


if __name__ == "__main__":
    raise SystemExit(main())
