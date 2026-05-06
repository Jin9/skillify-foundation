#!/usr/bin/env python3
"""Print the phase status of a pipeline run.

Reads manifest.json from .agent-pipelines/<task-id>/ and prints a status
table. Read-only.
"""

from __future__ import annotations

import argparse
import json
import sys
from pathlib import Path


PHASE_KEYS = ("plan", "gather", "analyze", "review", "validate", "decide", "compact")

STATUS_GLYPHS = {
    "pending":  "·",
    "running":  "▶",
    "done":     "✓",
    "failed":   "✗",
    "skipped":  "—",
}


def fail(message: str) -> None:
    print(f"error: {message}", file=sys.stderr)
    raise SystemExit(1)


def find_task_dir(task_id: str) -> Path:
    cwd = Path.cwd()
    task_dir = cwd / ".agent-pipelines" / task_id
    if not task_dir.is_dir():
        fail(f"task directory does not exist: {task_dir}")
    return task_dir


def load_manifest(task_dir: Path) -> dict:
    manifest_path = task_dir / "manifest.json"
    if not manifest_path.is_file():
        fail(f"manifest.json missing in {task_dir}")
    try:
        return json.loads(manifest_path.read_text(encoding="utf-8"))
    except json.JSONDecodeError as exc:
        fail(f"manifest.json is not valid JSON: {exc}")


def render(manifest: dict) -> str:
    lines: list[str] = []
    lines.append(f"task_id: {manifest['task_id']}")
    lines.append(f"domain:  {manifest['domain']}")
    lines.append(f"shape:   {' -> '.join(manifest['phase_shape'])}")
    lines.append("")
    lines.append(f"{'phase':<10} {'status':<10} {'artifact':<22} workers")
    lines.append("-" * 60)
    for key in PHASE_KEYS:
        phase = manifest["phases"][key]
        status = phase["status"]
        glyph = STATUS_GLYPHS.get(status, "?")
        artifact = phase["artifact"]
        delegation_count = len(phase.get("delegations", []))
        lines.append(f"{key:<10} {glyph} {status:<8} {artifact:<22} {delegation_count}")
    halted = manifest.get("halted_reason")
    if halted:
        lines.append("")
        lines.append(f"HALTED: {halted}")
    return "\n".join(lines)


def main() -> int:
    parser = argparse.ArgumentParser(description="Print pipeline phase status.")
    parser.add_argument("task_id", help="Task id, e.g. audit-auth-a3f2c1")
    parser.add_argument("--json", action="store_true", help="Emit raw manifest as JSON instead of the table")
    args = parser.parse_args()

    task_dir = find_task_dir(args.task_id)
    manifest = load_manifest(task_dir)

    if args.json:
        print(json.dumps(manifest, indent=2))
    else:
        print(render(manifest))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
