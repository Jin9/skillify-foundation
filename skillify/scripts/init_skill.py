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

Every step is model-cost tier: [small | mid | frontier] unless tagged.

Goal: [what a finished result looks like, in one sentence]
Constraints: [the one or two real constraints, each with its reason]
Done when: [the evidence that the job is complete]

1. Read the user's request and identify the concrete input.
2. Apply the reusable workflow this skill owns; number further steps only where order matters.
3. Produce the output contract below.
4. Verify the result against the evidence named below before responding.

## Output Format

[State exact files, formats, paths, or response sections.]

## Operating contract

- Instruction priority: the user's request in this session takes precedence over this skill; repo-policy files (AGENTS.md, CLAUDE.md, or the host equivalent) take precedence over this skill's defaults. If following a line here would make you pause, ask for permission, leave requested work unfinished, or diverge from what the user asked, follow the user, say which line you set aside, and quote it.
- Autonomy: "can you", "help me", and "please" are instructions. Once the inputs above are present, act; do not ask for confirmation of work the user already authorized. Ask at most one question per run, and only when a required input is missing and cannot be taken from the request, pasted text, or the workspace; otherwise state the assumption in your opening line and proceed.
- Stop conditions: stop and ask only before [the skill's irreversible action, if any], or when finishing would change the scope the user set. When the user asks a question rather than for a change, the assessment is the deliverable. Before ending your turn, check your last paragraph: if it is a plan or a promise, do that work now.
- Verification: before claiming success, check [the specific evidence: exit code, PASS line, re-read of the written file] and quote it in the recap. Do not describe a check you did not run. Do not add tests for reversible, low-impact changes.
- Delegation: [none, and delete this line | which parts may run as parallel sub-agents and what each returns].
- Progress: open with one line saying what you are about to do and which files you will touch; close with a recap that stands on its own (what changed, what was verified, what remains).
- Model-cost tier: [small | mid | frontier] by default; steps that differ are tagged inline as [tier/effort].

## Constraints

- Stay inside this skill's one responsibility; hand adjacent tasks to their own skill.
- Keep long reference material in `references/`; `SKILL.md` loads on every trigger.

## Verification

- [ ] Frontmatter name matches the folder.
- [ ] Description includes concrete trigger phrases.
- [ ] Workflow, operating contract, and output contract are explicit.

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
