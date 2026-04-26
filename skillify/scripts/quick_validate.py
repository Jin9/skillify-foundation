#!/usr/bin/env python3
"""Deterministic SKILL.md structure and frontmatter checks."""

from __future__ import annotations

import argparse
import re
import sys
from pathlib import Path


NAME_RE = re.compile(r"^[a-z0-9]+(?:-[a-z0-9]+)*$")
BANNED_DOCS = {
    "README.md",
    "CHANGELOG.md",
    "INSTALLATION_GUIDE.md",
    "QUICK_REFERENCE.md",
    "CONTRIBUTING.md",
}


def parse_frontmatter(text: str) -> tuple[dict[str, str], list[str], list[str]]:
    lines = text.splitlines()
    if not lines or lines[0].strip() != "---":
        return {}, [], ["SKILL.md must start with YAML frontmatter delimiter ---"]

    close_index = None
    for index, line in enumerate(lines[1:], start=1):
        if line.strip() == "---":
            close_index = index
            break
    if close_index is None:
        return {}, [], ["SKILL.md frontmatter must close with ---"]

    raw = lines[1:close_index]
    fields: dict[str, str] = {}
    errors: list[str] = []
    index = 0
    while index < len(raw):
        line = raw[index]
        stripped = line.strip()
        if not stripped or stripped.startswith("#"):
            index += 1
            continue
        if line.startswith((" ", "\t")):
            errors.append(f"unexpected indented frontmatter line: {line}")
            index += 1
            continue
        if ":" not in line:
            errors.append(f"invalid frontmatter line: {line}")
            index += 1
            continue

        key, value = line.split(":", 1)
        key = key.strip()
        value = value.strip()
        if value in {">", "|", ">-", "|-", ">+", "|+"}:
            block: list[str] = []
            index += 1
            while index < len(raw):
                next_line = raw[index]
                if next_line and not next_line.startswith((" ", "\t")):
                    break
                block.append(next_line.strip())
                index += 1
            if value.startswith(">"):
                fields[key] = " ".join(part for part in block if part)
            else:
                fields[key] = "\n".join(block)
            continue

        fields[key] = value.strip("\"'")
        index += 1

    return fields, lines[close_index + 1 :], errors


def validate(skill_dir: Path) -> list[str]:
    errors: list[str] = []
    skill_file = skill_dir / "SKILL.md"
    if not skill_file.exists():
        return [f"missing required file: {skill_file}"]
    if not skill_file.is_file():
        return [f"SKILL.md is not a file: {skill_file}"]

    text = skill_file.read_text(encoding="utf-8")
    fields, body, parse_errors = parse_frontmatter(text)
    errors.extend(parse_errors)
    if parse_errors:
        return errors

    name = fields.get("name", "")
    description = fields.get("description", "")
    if not name:
        errors.append("frontmatter missing required field: name")
    elif not NAME_RE.fullmatch(name):
        errors.append("frontmatter name must be lowercase kebab-case")
    elif len(name) > 64:
        errors.append("frontmatter name must be 64 characters or fewer")
    elif "claude" in name or "anthropic" in name:
        errors.append("frontmatter name must not contain reserved vendor names")
    elif skill_dir.name != name:
        errors.append(f"frontmatter name must match folder name: {name} != {skill_dir.name}")

    if not description:
        errors.append("frontmatter missing required field: description")
    else:
        if len(description) > 1024:
            errors.append(f"frontmatter description is too long: {len(description)} > 1024")
        if "<" in description or ">" in description:
            errors.append("frontmatter description must not contain XML angle brackets")
        trigger_markers = (
            "use when",
            "use for",
            "use after",
            "use specifically when",
            "triggers on",
            "activate when",
        )
        if not any(marker in description.lower() for marker in trigger_markers):
            errors.append(
                "frontmatter description should include trigger language "
                '(one of: "Use when", "Use for", "Use after", '
                '"Triggers on", "Activate when")'
            )

    for key, value in fields.items():
        if "<" in value or ">" in value:
            errors.append(f"frontmatter field must not contain XML angle brackets: {key}")

    total_lines = len(text.splitlines())
    if total_lines > 500:
        errors.append(f"SKILL.md is too long: {total_lines} lines > 500")

    if not any(line.startswith("## ") for line in body):
        errors.append("SKILL.md body should contain section headings")

    for path in skill_dir.rglob("*"):
        if path.is_file() and path.name in BANNED_DOCS:
            errors.append(f"banned human-facing doc inside skill folder: {path.relative_to(skill_dir)}")

    references_dir = skill_dir / "references"
    if references_dir.exists():
        for path in references_dir.rglob("*"):
            if path.is_file() and path.parent != references_dir:
                errors.append(f"references must be one level deep: {path.relative_to(skill_dir)}")

    return errors


def main() -> int:
    parser = argparse.ArgumentParser(description="Validate SKILL.md frontmatter and structure.")
    parser.add_argument("skill_dir", help="Path to a skill folder")
    args = parser.parse_args()

    skill_dir = Path(args.skill_dir).expanduser().resolve()
    if not skill_dir.exists():
        print(f"error: skill folder does not exist: {skill_dir}", file=sys.stderr)
        return 1
    if not skill_dir.is_dir():
        print(f"error: path is not a directory: {skill_dir}", file=sys.stderr)
        return 1

    errors = validate(skill_dir)
    if errors:
        for error in errors:
            print(f"error: {error}", file=sys.stderr)
        return 1

    print(f"ok: {skill_dir}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
