#!/usr/bin/env python3
"""Create a minimal skill folder skeleton."""

from __future__ import annotations

import argparse
import re
import sys
from pathlib import Path


NAME_RE = re.compile(r"^[a-z0-9]+(?:-[a-z0-9]+)*$")


def fail(message: str) -> None:
    print(f"error: {message}", file=sys.stderr)
    raise SystemExit(1)


def validate_name(name: str) -> None:
    if len(name) > 64:
        fail("skill name must be 64 characters or fewer")
    if not NAME_RE.fullmatch(name):
        fail("skill name must be lowercase kebab-case")
    if "claude" in name or "anthropic" in name:
        fail("skill name must not contain reserved vendor names")


def skill_markdown(name: str, title: str, description: str) -> str:
    return f"""---
name: {name}
description: >
  {description}
---

# {title}

## Purpose

[One sentence describing the reusable job this skill performs.]

## When to use this skill

- Use when: [specific trigger phrase or user intent]
- Use when: [another specific trigger phrase or user intent]
- Use when: [a third specific trigger phrase or user intent]
- Do NOT use when: [out-of-scope scenario]

## Workflow

1. Read the user's request and identify the concrete input.
2. Apply the reusable workflow this skill owns.
3. Produce the output contract below.
4. Validate the result before responding.

## Output Format

[State exact files, formats, paths, or response sections.]

## Constraints

- DO NOT handle adjacent tasks outside this skill's responsibility.
- DO NOT duplicate long reference material in `SKILL.md`.

## Validation

- [ ] Frontmatter name matches the folder.
- [ ] Description includes concrete trigger phrases.
- [ ] Workflow and output contract are explicit.

## References

- For [topic]: see `references/[filename].md` (create as needed)
"""


def main() -> int:
    parser = argparse.ArgumentParser(description="Create a minimal SKILL.md folder skeleton.")
    parser.add_argument("name", help="Skill folder/name in lowercase kebab-case")
    parser.add_argument("--root", default=".", help="Parent directory to create the skill in")
    parser.add_argument("--title", help="Human-readable title")
    parser.add_argument(
        "--description",
        default='[What this skill does]. Use when the user asks to "[trigger phrase]". Do NOT use for [negative trigger].',
        help="Frontmatter description text",
    )
    parser.add_argument("--references", action="store_true", help="Create references/ directory")
    parser.add_argument("--templates", action="store_true", help="Create templates/ directory")
    parser.add_argument("--scripts", action="store_true", help="Create scripts/ directory")
    parser.add_argument("--assets", action="store_true", help="Create assets/ directory")
    parser.add_argument("--examples", action="store_true", help="Create examples/ directory")
    args = parser.parse_args()

    validate_name(args.name)
    root = Path(args.root).expanduser().resolve()
    if not root.exists():
        fail(f"root does not exist: {root}")
    if not root.is_dir():
        fail(f"root is not a directory: {root}")

    skill_dir = root / args.name
    if skill_dir.exists():
        fail(f"target already exists: {skill_dir}")

    title = args.title or args.name.replace("-", " ").title()
    skill_dir.mkdir()
    (skill_dir / "SKILL.md").write_text(
        skill_markdown(args.name, title, args.description),
        encoding="utf-8",
    )

    for dirname, enabled in {
        "references": args.references,
        "templates": args.templates,
        "scripts": args.scripts,
        "assets": args.assets,
        "examples": args.examples,
    }.items():
        if enabled:
            (skill_dir / dirname).mkdir()

    print(skill_dir)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
