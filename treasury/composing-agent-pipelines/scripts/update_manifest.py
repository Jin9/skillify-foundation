#!/usr/bin/env python3
"""Update phase status, timestamps, and delegated calls in a pipeline manifest.

Called by the orchestrator at phase boundaries:

  python3 update_manifest.py --task-id audit-X --phase plan --status running \
      --add-delegation "writer|Plan: decompose audit task"
  python3 update_manifest.py --task-id audit-X --phase plan --status done

Read-only on the host: writes only to <cwd>/.agent-pipelines/<task-id>/manifest.json.
"""

from __future__ import annotations

import argparse
import json
import sys
from datetime import datetime, timezone
from pathlib import Path


PHASE_KEYS = {"plan", "gather", "analyze", "review", "validate", "decide", "compact"}
STATUS_VALUES = {"pending", "running", "done", "failed", "skipped"}


def fail(message: str) -> None:
    print(f"error: {message}", file=sys.stderr)
    raise SystemExit(1)


def load(task_dir: Path) -> dict:
    manifest_path = task_dir / "manifest.json"
    if not manifest_path.is_file():
        fail(f"manifest.json missing in {task_dir}")
    try:
        return json.loads(manifest_path.read_text(encoding="utf-8"))
    except json.JSONDecodeError as exc:
        fail(f"manifest.json is not valid JSON: {exc}")


def save(task_dir: Path, manifest: dict) -> None:
    manifest_path = task_dir / "manifest.json"
    manifest_path.write_text(json.dumps(manifest, indent=2) + "\n", encoding="utf-8")


def parse_delegation(spec: str) -> dict:
    if "|" not in spec:
        fail(f"--add-delegation must be 'worker_type|description', got: {spec!r}")
    worker_type, description = spec.split("|", 1)
    worker_type = worker_type.strip()
    description = description.strip()
    if not worker_type or not description:
        fail(f"--add-delegation requires non-empty worker_type and description, got: {spec!r}")
    return {
        "worker_type": worker_type,
        "description": description,
        "background_id": None,
    }


def main() -> int:
    parser = argparse.ArgumentParser(description="Update a phase in a pipeline manifest.")
    parser.add_argument("--task-id", required=True)
    parser.add_argument("--phase", required=True, choices=sorted(PHASE_KEYS))
    parser.add_argument("--status", required=True, choices=sorted(STATUS_VALUES))
    parser.add_argument(
        "--add-delegation",
        action="append",
        default=[],
        help="Format: 'worker_type|description'. May repeat.",
    )
    parser.add_argument("--note", help="Set phases.<phase>.notes to this string.")
    parser.add_argument("--halt", help="Set top-level halted_reason and mark phase failed.")
    args = parser.parse_args()

    cwd = Path.cwd()
    task_dir = cwd / ".agent-pipelines" / args.task_id
    if not task_dir.is_dir():
        fail(f"task directory does not exist: {task_dir}")

    manifest = load(task_dir)
    if args.phase not in manifest.get("phases", {}):
        fail(f"phase {args.phase!r} not in manifest")

    phase = manifest["phases"][args.phase]
    now = datetime.now(timezone.utc).isoformat(timespec="seconds")

    if args.status == "running":
        if phase.get("started_at") is None:
            phase["started_at"] = now
        phase["status"] = "running"
    elif args.status in {"done", "failed"}:
        if phase.get("started_at") is None:
            phase["started_at"] = now
        phase["ended_at"] = now
        phase["status"] = args.status
    elif args.status == "skipped":
        phase["status"] = "skipped"
    elif args.status == "pending":
        phase["status"] = "pending"

    for spec in args.add_delegation:
        phase.setdefault("delegations", []).append(parse_delegation(spec))

    if args.note:
        phase["notes"] = args.note

    if args.halt:
        manifest["halted_reason"] = args.halt
        phase["status"] = "failed"
        if phase.get("ended_at") is None:
            phase["ended_at"] = now

    manifest["updated_at"] = now
    save(task_dir, manifest)
    print(f"ok: {task_dir}/manifest.json (phase {args.phase} -> {phase['status']})")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
