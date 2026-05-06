#!/usr/bin/env python3
"""Initialize a pipeline state directory under .agent-pipelines/<task-id>/.

Creates the directory tree, writes a skeleton manifest.json, and prints the
absolute task directory so the orchestrator can record it. Refuses to overwrite
an existing task-id; use --resume to continue an existing pipeline.

Read-only on the host: only writes inside the new task directory.
"""

from __future__ import annotations

import argparse
import json
import re
import secrets
import sys
from datetime import datetime, timezone
from pathlib import Path


SLUG_RE = re.compile(r"^[a-z0-9]+(?:-[a-z0-9]+)*$")
TASK_ID_RE = re.compile(r"^[a-z0-9]+(?:-[a-z0-9]+)*-[0-9a-f]{6}$")
DOMAINS = {"research", "code-review", "implementation", "decision"}
PHASE_KEYS = ("plan", "gather", "analyze", "review", "validate", "decide", "compact")

DOMAIN_SHAPES: dict[str, list[str]] = {
    "research":       ["plan", "gather", "analyze", "review", "validate", "compact"],
    "code-review":    ["plan", "gather", "analyze", "review", "validate", "decide", "compact"],
    "implementation": ["plan", "gather", "analyze", "decide", "compact"],
    "decision":       ["plan", "gather", "analyze", "decide", "validate", "compact"],
}


def fail(message: str) -> None:
    print(f"error: {message}", file=sys.stderr)
    raise SystemExit(1)


def scrub(prompt: str) -> str:
    """Replace token-like and credential-like substrings with [redacted].

    Conservative: errs on the side of redacting too much rather than leaking
    secrets into the manifest.
    """
    redacted = prompt
    redacted = re.sub(r"\b[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,}\b", "[redacted-email]", redacted)
    redacted = re.sub(r"(?i)\b(password|passwd|secret|token|api[_-]?key|bearer)\s*[:=]\s*\S+", r"\1=[redacted]", redacted)
    redacted = re.sub(r"\b[A-Za-z0-9+/=_\-]{24,}\b", "[redacted-token]", redacted)
    return redacted


def validate_slug(slug: str) -> None:
    if not (3 <= len(slug) <= 24):
        fail("slug must be 3–24 characters")
    if not SLUG_RE.fullmatch(slug):
        fail("slug must be lowercase kebab-case (letters, digits, hyphens)")


def make_task_id(slug: str) -> str:
    return f"{slug}-{secrets.token_hex(3)}"


def empty_phase(artifact: str) -> dict:
    return {
        "status": "pending",
        "started_at": None,
        "ended_at": None,
        "artifact": artifact,
        "delegations": [],
    }


def build_manifest(task_id: str, slug: str, domain: str, prompt: str) -> dict:
    now = datetime.now(timezone.utc).isoformat(timespec="seconds")
    shape = DOMAIN_SHAPES[domain]
    phase_artifacts = {
        "plan": "01-plan.md",
        "gather": "02-evidence/",
        "analyze": "03-analysis.md",
        "review": "04-review.md",
        "validate": "05-validation.md",
        "decide": "06-decision.md",
        "compact": "07-final.md",
    }
    phases = {}
    for key in PHASE_KEYS:
        phases[key] = empty_phase(phase_artifacts[key])
        if key not in shape:
            phases[key]["status"] = "skipped"
    return {
        "task_id": task_id,
        "slug": slug,
        "domain": domain,
        "created_at": now,
        "updated_at": now,
        "phase_shape": shape,
        "phases": phases,
        "user_prompt": scrub(prompt),
        "halted_reason": None,
    }


def init_new(args: argparse.Namespace) -> Path:
    if args.domain not in DOMAINS:
        fail(f"domain must be one of: {', '.join(sorted(DOMAINS))}")
    validate_slug(args.slug)
    task_id = make_task_id(args.slug)

    cwd = Path.cwd()
    pipelines_root = cwd / ".agent-pipelines"
    pipelines_root.mkdir(parents=True, exist_ok=True)

    task_dir = pipelines_root / task_id
    if task_dir.exists():
        fail(f"task directory already exists: {task_dir}")
    task_dir.mkdir()
    (task_dir / "02-evidence").mkdir()

    manifest = build_manifest(task_id, args.slug, args.domain, args.prompt)
    (task_dir / "manifest.json").write_text(
        json.dumps(manifest, indent=2) + "\n",
        encoding="utf-8",
    )
    return task_dir


def resume_existing(args: argparse.Namespace) -> Path:
    if not TASK_ID_RE.fullmatch(args.resume):
        fail("--resume task-id must match <slug>-<6 hex chars>")
    cwd = Path.cwd()
    task_dir = cwd / ".agent-pipelines" / args.resume
    if not task_dir.is_dir():
        fail(f"task directory does not exist: {task_dir}")
    manifest_path = task_dir / "manifest.json"
    if not manifest_path.is_file():
        fail(f"manifest.json missing in {task_dir}")
    try:
        manifest = json.loads(manifest_path.read_text(encoding="utf-8"))
    except json.JSONDecodeError as exc:
        fail(f"manifest.json is not valid JSON: {exc}")
    next_phase = next(
        (p for p in manifest["phase_shape"] if manifest["phases"][p]["status"] == "pending"),
        None,
    )
    if next_phase is None:
        print(f"resumed: {task_dir}")
        print("status: all phases complete or skipped")
    else:
        print(f"resumed: {task_dir}")
        print(f"next phase: {next_phase}")
    return task_dir


def main() -> int:
    parser = argparse.ArgumentParser(description="Initialize or resume a pipeline state directory.")
    parser.add_argument("--slug", help="Short slug (3–24 chars, kebab-case). Required when not resuming.")
    parser.add_argument("--domain", choices=sorted(DOMAINS), help="Task domain. Required when not resuming.")
    parser.add_argument("--prompt", default="", help="User task prompt (scrubbed before persistence).")
    parser.add_argument("--resume", help="Existing task-id to resume.")
    args = parser.parse_args()

    if args.resume:
        if args.slug or args.domain:
            fail("--slug and --domain are not allowed with --resume")
        task_dir = resume_existing(args)
    else:
        if not args.slug or not args.domain:
            fail("--slug and --domain are required when not resuming")
        task_dir = init_new(args)
        print(f"created: {task_dir}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
