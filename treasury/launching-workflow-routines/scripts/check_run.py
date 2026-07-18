#!/usr/bin/env python3
"""Verify run-directory artifacts against a routine's declared outputs.

Usage:
  python3 check_run.py ROUTINE.md RUN_DIR              # check all nodes
  python3 check_run.py ROUTINE.md RUN_DIR --node ID    # check one node (repeatable)
  python3 check_run.py ROUTINE.md RUN_DIR --final      # also require run-report.md

Trusts only the filesystem, never executor self-reports. Prints a
per-artifact table plus a machine line:
  RESULT: pass|fail nodes=N artifacts=M missing=K
Exit 0 only when every checked declared artifact exists and is non-empty,
INDEX.md is present and non-empty, and (with --final) run-report.md too.
"""

from __future__ import annotations

import argparse
import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent))

from routine_parser import parse_routine

STUB_BYTES = 64
RESERVED = {"INDEX.md", "run-report.md"}


def main() -> int:
    parser = argparse.ArgumentParser(description="Verify run artifacts for a routine.")
    parser.add_argument("routine", help="Routine definition .md file")
    parser.add_argument("run_dir", help="Run directory to verify")
    parser.add_argument("--node", action="append", help="Check only this node id (repeatable)")
    parser.add_argument("--final", action="store_true", help="Also require a non-empty run-report.md")
    args = parser.parse_args()

    routine_path = Path(args.routine).expanduser().resolve()
    run_dir = Path(args.run_dir).expanduser().resolve()

    routine = parse_routine(routine_path)
    if routine.errors:
        for error in routine.errors:
            print(f"error: {routine_path.name}: {error}", file=sys.stderr)
        print("RESULT: fail nodes=0 artifacts=0 missing=0")
        return 1
    if not run_dir.is_dir():
        print(f"error: run dir does not exist: {run_dir}", file=sys.stderr)
        print("RESULT: fail nodes=0 artifacts=0 missing=0")
        return 1

    nodes = routine.nodes
    if args.node:
        known = {node.node_id for node in nodes}
        unknown = [node_id for node_id in args.node if node_id not in known]
        if unknown:
            print(f"error: unknown node id(s): {', '.join(unknown)}", file=sys.stderr)
            print("RESULT: fail nodes=0 artifacts=0 missing=0")
            return 1
        selected = set(args.node)
        nodes = [node for node in nodes if node.node_id in selected]

    artifact_failures = 0
    structural_failures = 0
    rows = []
    for node in nodes:
        for out in node.outputs():
            artifact = run_dir / out
            if not artifact.is_file():
                status, size = "MISSING", "-"
                artifact_failures += 1
            else:
                byte_count = artifact.stat().st_size
                size = str(byte_count)
                if byte_count == 0:
                    status = "EMPTY"
                    artifact_failures += 1
                elif byte_count < STUB_BYTES:
                    status = "ok (suspect-stub)"
                else:
                    status = "ok"
            rows.append((node.node_id, out, status, size))

    index_file = run_dir / "INDEX.md"
    if not index_file.is_file() or index_file.stat().st_size == 0:
        print(f"error: run dir INDEX.md missing or empty: {index_file}", file=sys.stderr)
        structural_failures += 1
    if args.final:
        report = run_dir / "run-report.md"
        if not report.is_file() or report.stat().st_size == 0:
            print(f"error: run-report.md missing or empty (required by --final): {report}", file=sys.stderr)
            structural_failures += 1

    declared = {out for node in routine.nodes for out in node.outputs()}
    for candidate in sorted(run_dir.rglob("*")):
        if not candidate.is_file():
            continue
        rel = candidate.relative_to(run_dir).as_posix()
        if rel in declared or rel in RESERVED:
            continue
        print(f"warning: undeclared file in run dir: {rel}", file=sys.stderr)

    width = max([len(row[1]) for row in rows] + [len("artifact")])
    print(f"{'node':<20} {'artifact':<{width}}  {'status':<18} bytes")
    for node_id, out, status, size in rows:
        print(f"{node_id:<20} {out:<{width}}  {status:<18} {size}")

    verdict = "pass" if artifact_failures + structural_failures == 0 else "fail"
    print(f"RESULT: {verdict} nodes={len(nodes)} artifacts={len(rows)} missing={artifact_failures}")
    return 0 if verdict == "pass" else 1


if __name__ == "__main__":
    raise SystemExit(main())
