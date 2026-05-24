#!/usr/bin/env python3
"""business-logic-extractor: provenance & ledger self-check (deterministic).

Validates the skill's OWN output before it is delivered:
  1. every `path:line` / `path:line-line` citation in the spec resolves to an
     existing file with that line in range (relative to --root, default cwd);
  2. a loss ledger exists, is non-empty, and contains a coverage section;
  3. each rule section (### R...) carries a decision-logic flowchart
     (```ascii), provenance-annotated pseudo-code (```text), and at least one
     resolvable file:line citation; code fences are balanced. (Skipped when the
     spec uses no rule headings.)

No network, no LLM, no code execution — stdlib only. Exit 0 = ok,
1 = unresolved provenance / missing-or-empty ledger / rule render gap,
2 = usage error.

Usage:
  python3 check_provenance.py business-logic-spec.md [--ledger loss-ledger.md] [--root .]
"""
import argparse, re, sys
from pathlib import Path

# `path/to/file.ext:120` or `:120-138`, inside backticks or bare.
CITE = re.compile(r'`?([\w./\-]+\.[A-Za-z0-9_]+):(\d+)(?:-(\d+))?`?')

RULE_HEAD = re.compile(r'(?m)^###\s+R\S.*$')   # "### R1 — ...", "### R12 - ..."
SEC_HEAD = re.compile(r'(?m)^##\s+\S.*$')        # top-level "## " section breaks


def _is_self_ref(fp: str) -> bool:
    return fp.endswith(("spec.md", "ledger.md"))


def check_rules(text: str) -> list[str]:
    """Per-rule render contract: each `### R` section must carry an ASCII
    flowchart, a pseudo-code block, and >=1 non-self file:line citation; code
    fences must be balanced. Citation *resolution* is left to the global pass.
    Returns [] when the spec uses no rule headings (different output shape)."""
    problems: list[str] = []
    if "```" in text and text.count("```") % 2 != 0:
        problems.append("unbalanced code fence (``` count is odd) — an "
                        "ascii/pseudo-code block is not closed")
    heads = [(m.start(), m.group().strip()) for m in RULE_HEAD.finditer(text)]
    if not heads:
        return problems
    bounds = sorted({s for s, _ in heads}
                    | {m.start() for m in SEC_HEAD.finditer(text)}
                    | {len(text)})
    for start, head in heads:
        end = min(b for b in bounds if b > start)
        section = text[start:end]
        rid = head.lstrip("# ").split()[0]  # "R9" from "### R9 - name"
        if "```ascii" not in section:
            problems.append(f"{rid}: missing decision-logic flowchart (```ascii block)")
        if "```text" not in section:
            problems.append(f"{rid}: missing pseudo-code block (```text block)")
        if not [fp for fp, _1, _2 in CITE.findall(section) if not _is_self_ref(fp)]:
            problems.append(f"{rid}: no file:line citation in the rule section")
    return problems


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("spec")
    ap.add_argument("--ledger", default=None,
                    help="loss ledger path (default: sibling loss-ledger.md)")
    ap.add_argument("--root", default=".", help="repo root for resolving citations")
    a = ap.parse_args()

    spec_path = Path(a.spec)
    if not spec_path.is_file():
        print(f"error: spec not found: {a.spec}", file=sys.stderr)
        sys.exit(2)
    root = Path(a.root)
    text = spec_path.read_text(encoding="utf-8", errors="replace")

    cites = CITE.findall(text)
    if not cites:
        print("FAIL: no `file:line` provenance citations found in the spec — "
              "every rule must cite real code.")
        sys.exit(1)

    unresolved = []
    checked = 0
    for fp, l1, l2 in cites:
        # skip self-references to the spec's own companion files
        if fp.endswith(("spec.md", "ledger.md")):
            continue
        checked += 1
        target = (root / fp)
        if not target.is_file():
            unresolved.append(f"{fp}:{l1} — file not found under {a.root}")
            continue
        try:
            n = sum(1 for _ in target.open("rb"))
        except Exception as e:
            unresolved.append(f"{fp}:{l1} — unreadable ({e})")
            continue
        lo, hi = int(l1), int(l2 or l1)
        if lo < 1 or hi < lo or hi > n:
            unresolved.append(f"{fp}:{l1}{'-'+l2 if l2 else ''} "
                              f"— line out of range (file has {n} lines)")

    # ledger presence + non-empty + coverage section
    ledger = Path(a.ledger) if a.ledger else spec_path.with_name("loss-ledger.md")
    ledger_problems = []
    if not ledger.is_file():
        ledger_problems.append(f"loss ledger missing: {ledger}")
    else:
        lt = ledger.read_text(encoding="utf-8", errors="replace").strip()
        if len(lt) < 80:
            ledger_problems.append("loss ledger present but effectively empty")
        if not re.search(r"coverage", lt, re.I):
            ledger_problems.append("loss ledger has no coverage section")

    rule_problems = check_rules(text)

    ok = not unresolved and not ledger_problems and not rule_problems
    print(f"provenance citations checked: {checked}")
    if unresolved:
        print(f"UNRESOLVED ({len(unresolved)}):")
        for u in unresolved:
            print(f"  - {u}")
    if ledger_problems:
        print("LEDGER:")
        for p in ledger_problems:
            print(f"  - {p}")
    if rule_problems:
        print(f"RULES ({len(rule_problems)}):")
        for p in rule_problems:
            print(f"  - {p}")
    print("OK — every cited rule resolves, renders an ASCII flowchart + pseudo-code, "
          "and the loss ledger is present."
          if ok else "FAIL — fix the above before delivering the spec.")
    sys.exit(0 if ok else 1)


if __name__ == "__main__":
    main()
