#!/usr/bin/env python3
"""
validate_story_trace.py — story↔AC trace + EPIC→Story integrity.

Verifies:
  1. Every story file under requirement/<EPIC_SLUG>/STORY_*.md has a parseable
     YAML frontmatter with story_id, epic_id, acceptance_criteria.
  2. Every story's epic_id resolves to an existing requirement/<EPIC_SLUG>/<EPIC_SLUG>.md.
  3. Every AC in ba.json is referenced by ≥1 story (heuristic: matches by AC index
     since ba.json's AC items don't carry IDs; we check that the count of stories'
     ACs is ≥ the count in ba.json).

Usage:
  python3 validate_story_trace.py <ba.json> <requirement-dir>

Exit 0 = pass; non-zero = trace failure.
"""
from __future__ import annotations

import json
import re
import sys
from pathlib import Path


def fail(msg: str) -> None:
    print(f"validate_story_trace: FAIL {msg}", file=sys.stderr)
    sys.exit(2)


def parse_frontmatter(text: str) -> dict | None:
    """Tiny YAML-frontmatter parser — only enough to extract story_id, epic_id,
    and count acceptance_criteria items. Avoids requiring PyYAML."""
    m = re.match(r"^---\s*\n(.*?)\n---\s*\n", text, re.DOTALL)
    if not m:
        return None
    body = m.group(1)
    out: dict = {"acceptance_criteria_count": 0}
    in_ac = False
    for raw in body.splitlines():
        stripped = raw.strip()
        if not stripped or stripped.startswith("#"):
            continue
        if not raw.startswith(" ") and ":" in raw:
            key, _, val = raw.partition(":")
            key = key.strip()
            val = val.strip()
            in_ac = (key == "acceptance_criteria")
            if val and not in_ac:
                out[key] = val.strip().strip('"').strip("'")
            elif val == "" and not in_ac:
                out[key] = None
        elif in_ac and stripped.startswith("- "):
            out["acceptance_criteria_count"] += 1
    return out


def main(ba_path: str, req_dir: str) -> None:
    ba = json.loads(Path(ba_path).read_text())
    ac_count = len(ba.get("acceptance_criteria", []))
    req = Path(req_dir)
    if not req.is_dir():
        fail(f"requirement dir not found: {req_dir}")

    # Only directories whose name matches the EPIC slug pattern (EPIC_<DOMAIN>...)
    # are treated as epic folders; everything else (.git, design, scratch dirs) is ignored.
    epic_dirs = [
        d for d in req.iterdir()
        if d.is_dir() and not d.name.startswith(".") and d.name.startswith("EPIC_")
    ]
    if not epic_dirs:
        fail(
            f"no EPIC subfolders under {req_dir} — BA must emit "
            f"requirement/<EPIC_SLUG>/<EPIC_SLUG>.md per docs/templates/epic.md"
        )

    epic_ids: set[str] = set()
    story_count = 0
    story_ac_total = 0

    for ed in epic_dirs:
        epic_md = ed / f"{ed.name}.md"
        if not epic_md.exists():
            fail(f"missing EPIC description: {epic_md}")
        fm = parse_frontmatter(epic_md.read_text())
        if not fm or "epic_id" not in fm:
            fail(f"EPIC {epic_md} missing frontmatter or epic_id")
        epic_ids.add(fm["epic_id"])

        for story in ed.glob("STORY_*.md"):
            story_count += 1
            sfm = parse_frontmatter(story.read_text())
            if not sfm:
                fail(f"story {story} has no parseable frontmatter")
            if "story_id" not in sfm:
                fail(f"story {story} missing story_id")
            if "epic_id" not in sfm:
                fail(f"story {story} missing epic_id")
            if sfm["epic_id"] not in epic_ids and sfm["epic_id"] != fm["epic_id"]:
                fail(
                    f"story {sfm['story_id']} references unknown epic_id "
                    f"{sfm['epic_id']!r}"
                )
            story_ac_total += sfm.get("acceptance_criteria_count", 0)

    if story_count == 0:
        fail("no STORY_*.md files found under any epic dir")
    if story_ac_total < ac_count:
        fail(
            f"AC coverage gap: ba.json declares {ac_count} ACs but stories "
            f"collectively define only {story_ac_total} ACs — every AC needs a "
            f"story owner"
        )

    print(
        f"validate_story_trace: PASS "
        f"({len(epic_ids)} EPICs, {story_count} stories, "
        f"{story_ac_total} story-ACs covering {ac_count} ba ACs)"
    )


if __name__ == "__main__":
    if len(sys.argv) != 3:
        print(
            "usage: validate_story_trace.py <path-to-ba.json> <path-to-requirement-dir>",
            file=sys.stderr,
        )
        sys.exit(64)
    main(sys.argv[1], sys.argv[2])
