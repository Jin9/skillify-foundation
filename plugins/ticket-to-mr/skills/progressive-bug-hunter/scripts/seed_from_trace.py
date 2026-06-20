#!/usr/bin/env python3
"""progressive-bug-hunter: Tier-0 anchor helper (optional, deterministic).

Parse a stack trace / error text and emit a prioritized lexical seed plan:
file:line frames (deepest first) + candidate identifiers + ready-to-run
ripgrep seed queries. No network, no code execution, no LLM — stdlib only.
It only PROPOSES searches; the agent runs them per the workflow.

Usage:
  python3 seed_from_trace.py <trace-file>
  cat trace.txt | python3 seed_from_trace.py -
"""
import json, re, sys

# Frame patterns across common ecosystems.
FRAME_RES = [
    re.compile(r'File "(?P<file>[^"]+)", line (?P<line>\d+), in (?P<sym>\S+)'),   # Python
    re.compile(r'at (?P<sym>[\w$.<>]+)\s*\((?P<file>[^):]+):(?P<line>\d+)'),       # JS/Java
    re.compile(r'(?P<file>[\w./\-]+\.\w+):(?P<line>\d+):(?:\d+)?'),                # generic file:line[:col]
    re.compile(r'(?P<file>[\w./\-]+\.go):(?P<line>\d+)'),                          # Go
]
EXC_RE = re.compile(r'\b([A-Z][A-Za-z0-9_]*(?:Error|Exception|Panic|Failure|Fault))\b')
IDENT_RE = re.compile(r'\b([A-Za-z_][A-Za-z0-9_]{2,})\b')
STOP = {"the", "and", "for", "line", "file", "error", "exception", "traceback",
        "most", "recent", "call", "last", "self", "none", "null", "true",
        "false", "return", "import", "from", "test", "assert", "object"}


def main():
    arg = sys.argv[1] if len(sys.argv) > 1 else "-"
    text = sys.stdin.read() if arg == "-" else open(arg, encoding="utf-8",
                                                     errors="replace").read()
    frames, seen = [], set()
    for rx in FRAME_RES:
        for m in rx.finditer(text):
            d = m.groupdict()
            key = (d.get("file"), d.get("line"))
            if d.get("file") and key not in seen:
                seen.add(key)
                frames.append({"file": d["file"], "line": int(d["line"]),
                               "symbol": d.get("sym", "")})
    # deepest frame first = most-recent call last in Python; keep source order
    # but expose reversed too so the agent can pick.
    exceptions = sorted(set(EXC_RE.findall(text)))
    quoted = re.findall(r'"([^"]{3,60})"|\'([^\']{3,60})\'', text)
    msgs = sorted({a or b for a, b in quoted})
    idents = [w for w in dict.fromkeys(IDENT_RE.findall(text))
              if w.lower() not in STOP and not w.isdigit()]
    # rank identifiers: those appearing in frame symbols first
    framesyms = " ".join(f["symbol"] for f in frames)
    idents.sort(key=lambda w: (w not in framesyms, -len(w)))
    top_idents = idents[:12]

    seeds = []
    for e in exceptions:
        seeds.append(f'rg -n -- {e!r}')
    for f in frames[:8]:
        seeds.append(f'sed -n {max(1, f["line"]-8)},{f["line"]+8}p '
                     f'{f["file"]}   # frame: {f["symbol"] or "?"}')
    for w in top_idents[:8]:
        seeds.append(f'rg -n -w -- {w!r}')

    plan = {
        "frames": frames,
        "exceptions": exceptions,
        "messages": msgs[:8],
        "candidate_identifiers": top_idents,
        "seed_queries_priority_order": seeds,
        "note": "Tier-0 anchors only. Run the queries per the workflow; "
                "start Tier 1 (parallel grep), escalate per references/retrieval-ladder.md.",
    }
    print(json.dumps(plan, indent=2))


if __name__ == "__main__":
    main()
