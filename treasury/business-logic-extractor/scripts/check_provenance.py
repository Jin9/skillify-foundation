#!/usr/bin/env python3
"""business-logic-extractor: provenance & ledger self-check (deterministic).

Validates the skill's OWN output before it is delivered:
  1. every `path:line` / `path:line-line` citation in the spec resolves to an
     existing file with that line in range (relative to --root, default cwd);
  2. a loss ledger exists, is non-empty, and contains a coverage section.

No network, no LLM, no code execution — stdlib only. Exit 0 = ok,
1 = unresolved provenance / missing-or-empty ledger, 2 = usage error.

Usage:
  python3 check_provenance.py business-logic-spec.md [--ledger loss-ledger.md] [--root .]
"""
import argparse, re, sys
from pathlib import Path

# `path/to/file.ext:120` or `:120-138`, inside backticks or bare.
CITE = re.compile(r'`?([\w./\-]+\.[A-Za-z0-9_]+):(\d+)(?:-(\d+))?`?')


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

    ok = not unresolved and not ledger_problems
    print(f"provenance citations checked: {checked}")
    if unresolved:
        print(f"UNRESOLVED ({len(unresolved)}):")
        for u in unresolved:
            print(f"  - {u}")
    if ledger_problems:
        print("LEDGER:")
        for p in ledger_problems:
            print(f"  - {p}")
    print("OK — every cited rule resolves and the loss ledger is present."
          if ok else "FAIL — fix the above before delivering the spec.")
    sys.exit(0 if ok else 1)


if __name__ == "__main__":
    main()
