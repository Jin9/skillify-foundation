#!/usr/bin/env python3
"""Verify local resource pointers from SKILL.md resolve."""

from __future__ import annotations

import argparse
import re
import sys
from pathlib import Path


# Match resource paths only when backtick-delimited (e.g. `references/foo.md`).
# A bare path token in prose (e.g. "examples/framing" meaning "examples and
# framing") is not a link and must not be treated as a missing resource.
RESOURCE_RE = re.compile(
    r"`(?P<path>(?:references|templates|scripts|assets|examples)/[A-Za-z0-9._/\-]+(?:#[A-Za-z0-9._\-]+)?)`"
)
MARKDOWN_LINK_RE = re.compile(r"\[[^\]]+\]\((?P<target>[^)]+)\)")


def clean_target(target: str) -> str:
    target = target.strip().strip("`'\"")
    target = target.rstrip(".,;:")
    return target


def local_targets(text: str) -> set[str]:
    targets: set[str] = set()
    for match in RESOURCE_RE.finditer(text):
        targets.add(clean_target(match.group("path")))
    for match in MARKDOWN_LINK_RE.finditer(text):
        target = clean_target(match.group("target"))
        if not target or target.startswith(("#", "http://", "https://", "mailto:")):
            continue
        if target.startswith(("../", "/")):
            continue
        targets.add(target)
    return targets


def check(skill_dir: Path) -> list[str]:
    skill_file = skill_dir / "SKILL.md"
    if not skill_file.exists():
        return [f"missing required file: {skill_file}"]

    text = skill_file.read_text(encoding="utf-8")
    missing: list[str] = []
    for target in sorted(local_targets(text)):
        path_text = target.split("#", 1)[0]
        if not path_text:
            continue
        candidate = skill_dir / path_text
        if not candidate.exists():
            missing.append(path_text)
    return missing


def main() -> int:
    parser = argparse.ArgumentParser(description="Check local links and resource pointers in SKILL.md.")
    parser.add_argument("skill_dir", help="Path to a skill folder")
    args = parser.parse_args()

    skill_dir = Path(args.skill_dir).expanduser().resolve()
    if not skill_dir.exists():
        print(f"error: skill folder does not exist: {skill_dir}", file=sys.stderr)
        return 1
    if not skill_dir.is_dir():
        print(f"error: path is not a directory: {skill_dir}", file=sys.stderr)
        return 1

    missing = check(skill_dir)
    if missing:
        for path in missing:
            print(f"error: missing local resource: {path}", file=sys.stderr)
        return 1

    print(f"ok: {skill_dir}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
