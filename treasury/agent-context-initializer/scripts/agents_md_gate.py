#!/usr/bin/env python3
"""agent-context-initializer: candidate AGENTS.md gate (deterministic).

Validates a candidate AGENTS.md before it is delivered for the human pass:
  1. line count must not exceed the hard cap of 200 (WARN, not fail, > 150);
  2. all six required sections must be present (case-insensitive heading
     match on: commands, testing, project structure, code style,
     git workflow, boundaries);
  3. every prohibition line (containing don't / do not / never / must not /
     prohibited) must be accompanied — same line or the adjacent
     bullet/line — by an allowed-alternative cue (do / instead / use /
     prefer / allowed).

No network, no LLM, no code execution — stdlib only. Exit 0 = ok (may WARN),
1 = a blocking failure, 2 = usage error. Mirrors check_provenance.py.

Usage:
  python3 agents_md_gate.py path/to/AGENTS.md [--self-test]
"""
import argparse
import re
import sys
from pathlib import Path

HARD_CAP = 200
WARN_OVER = 150

REQUIRED_SECTIONS = (
    "commands",
    "testing",
    "project structure",
    "code style",
    "git workflow",
    "boundaries",
)

PROHIBITION = re.compile(r"\b(don't|do not|never|must not|prohibited)\b", re.I)
ALTERNATIVE = re.compile(r"\b(do |instead|use |prefer|allowed)\b", re.I)
HEADING = re.compile(r"^\s{0,3}#{1,6}\s+(.*\S)\s*$")


def strip_comments(text: str) -> str:
    """Remove HTML comment blocks so skeleton guidance is not scored."""
    return re.sub(r"<!--.*?-->", "", text, flags=re.S)


def find_missing_sections(scored: str) -> list:
    headings = []
    for line in scored.splitlines():
        m = HEADING.match(line)
        if m:
            headings.append(m.group(1).lower())
    blob = " ".join(headings)
    return [s for s in REQUIRED_SECTIONS if s not in blob]


def find_unpaired_prohibitions(scored: str) -> list:
    lines = scored.splitlines()
    unpaired = []
    for i, raw in enumerate(lines):
        line = raw.strip()
        if not line or not PROHIBITION.search(line):
            continue
        window = line
        # adjacent non-blank lines (previous and next) count as accompanying.
        for j in (i - 1, i + 1):
            if 0 <= j < len(lines):
                window += "\n" + lines[j]
        if not ALTERNATIVE.search(window):
            unpaired.append((i + 1, line))
    return unpaired


def gate(path: Path) -> int:
    text = path.read_text(encoding="utf-8", errors="replace")
    total = len(text.splitlines())
    scored = strip_comments(text)

    failures = []
    warnings = []

    if total > HARD_CAP:
        failures.append(
            f"line count {total} exceeds hard cap of {HARD_CAP}"
        )
    elif total > WARN_OVER:
        warnings.append(
            f"line count {total} exceeds target of {WARN_OVER} "
            f"(under {HARD_CAP} so not fatal — trim toward 100-150)"
        )

    missing = find_missing_sections(scored)
    if missing:
        failures.append(
            "missing required section heading(s): " + ", ".join(missing)
        )

    unpaired = find_unpaired_prohibitions(scored)
    if unpaired:
        failures.append(
            f"{len(unpaired)} prohibition(s) with no allowed alternative "
            f"(pair every don't with a do):"
        )

    print(f"lines: {total}  (cap {HARD_CAP}, target <= {WARN_OVER})")
    print(f"required sections present: {6 - len(missing)}/6")
    for w in warnings:
        print(f"WARN: {w}")
    for f in failures:
        print(f"FAIL: {f}")
    for ln, txt in unpaired:
        print(f"  - line {ln}: {txt}")

    if failures:
        print("FAIL — fix the above before the human curation pass.")
        return 1
    print("OK — gate passed; still requires a human editorial pass before commit.")
    return 0


def self_test() -> int:
    """Template must pass; a bad fixture must fail."""
    here = Path(__file__).resolve().parent
    template = here.parent / "templates" / "AGENTS.md"
    print("[self-test] templates/AGENTS.md (expect exit 0)")
    rc_good = gate(template)
    print()

    fixture = Path("/tmp/agents_md_gate_bad_fixture.md")
    fixture.write_text(
        "# AGENTS.md\n\n"
        "## Commands\n- Build: `make`\n\n"
        "## Testing\n- `pytest -q`\n\n"
        "## Project structure\n- `src/` core\n\n"
        "## Code style\nSee config.\n\n"
        "## Git workflow\nOpen a PR.\n\n"
        # NOTE: no Boundaries section (missing required section)
        "Never force-push.\n",  # bare prohibition, no alternative
        encoding="utf-8",
    )
    print("[self-test] /tmp bad fixture (expect exit 1: missing section + bare prohibition)")
    rc_bad = gate(fixture)
    print()

    ok = (rc_good == 0) and (rc_bad == 1)
    print(
        f"[self-test] template rc={rc_good} (want 0), "
        f"fixture rc={rc_bad} (want 1) -> "
        + ("PASS" if ok else "FAIL")
    )
    return 0 if ok else 1


def main() -> int:
    ap = argparse.ArgumentParser(
        description="Deterministic gate for a candidate AGENTS.md."
    )
    ap.add_argument(
        "agents_md",
        nargs="?",
        help="path to the candidate AGENTS.md",
    )
    ap.add_argument(
        "--self-test",
        action="store_true",
        help="run the built-in self-test (template passes, bad fixture fails)",
    )
    a = ap.parse_args()

    if a.self_test:
        return self_test()

    if not a.agents_md:
        print("error: agents_md path required (or use --self-test)", file=sys.stderr)
        return 2
    path = Path(a.agents_md)
    if not path.is_file():
        print(f"error: file not found: {a.agents_md}", file=sys.stderr)
        return 2
    return gate(path)


if __name__ == "__main__":
    raise SystemExit(main())
