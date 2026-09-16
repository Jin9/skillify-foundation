#!/usr/bin/env python3
"""Line-rule linter for plain, direct, analogy-free explanations.

Not a language model, not a rewriter; reports violations, never rewrites.

Usage:
    python3 check_plain_style.py [--lang auto|en|th] [--cap N] FILE [FILE ...]
    python3 check_plain_style.py - <<'EOF'      (draft on stdin)
    ...
    EOF

Output: one "file:line: RULE_ID message" per violation on stderr;
"ok: N input(s) clean (en, 97 words)" on stdout when everything passes. The
count in parentheses is the budget used above the tail; a multi-pass input
prints one count per pass ("103+94 words"); Thai inputs count characters.

Exit codes:
    0  all input clean
    1  one or more violations
    2  usage or I/O error

Preprocessing (before every rule):
    - fenced code blocks (``` or ~~~) are ignored; inline `code` spans are
      masked to one letterless token
    - list markers, heading marks, and blockquote marks are stripped; bold
      and italic asterisks are removed
    - a line of three or more dashes splits the input into passes (a first
      pass and its "more" layers); the cap is applied per pass
    - the last prose line of a pass is the tail when it starts with Say,
      Reply, Type, or Ask for and contains "more", or starts with พิมพ์ and
      contains เพิ่ม; the tail is skipped by every rule
    - language is "th" when at least 30% of the letters are Thai, else "en";
      it changes the unit of P3 and P4 and switches on the Thai P5 list

Rules:
    P1   analogy or stand-in scenario: imagine / picture this / think of X
         as / pretend / metaphor / analogy / kind of like, or a comparison
         marker followed by "a" or "an" ("works like a queue", "similar to a
         waiter"); Thai figurative markers (นึกภาพ, เปรียบเสมือน, ราวกับ, ...).
         Comparisons to a named real thing ("like the `handler` you wrote",
         "similar to `parseToken`") pass.
    P2   filler: openers (Sure, Certainly, Of course, Happy to, Let me
         explain, แน่นอน, ยินดี) and closers (great question, in summary,
         hope this helps, let me know if, หวังว่าจะช่วย, โดยสรุป, ...)
    P3   sentence length: more than 20 words in one English sentence, or
         more than 90 Thai characters without a space
    P4   cap per pass, tail excluded: 120 English words or 500 Thai
         characters by default (--cap overrides)
    P5   complex word with a plain replacement (utilize -> use, prior to ->
         before, via -> through, e.g. -> for example, ...); in Thai, also a
         textbook coinage or essay marker with its chat replacement
         (ฐานข้อมูล -> database, ดำเนินการ -> ทำ, ทั้งนี้ -> delete, ...)
    P6   third-person reader: "the user should/must/needs to", "the reader",
         ผู้อ่าน, ผู้ใช้ควร/ต้อง; address the reader as "you"
    P7   clause joiners in prose: semicolon, em dash, spaced dash or en dash
         (digit ranges such as 2019–2020 pass)
    P8   question mark in prose (rhetorical questions); phrase a question to
         the user as an imperative instead
    P9   skeleton of the first pass: a bold first line, 2 to 5 numbered
         lines, and a "For you:" (สำหรับคุณ:) line

Known limits (deliberate):
    - No language model: a clean run is necessary, not sufficient. The manual
      pass owns personification, "suppose/say you have a ..." scenarios, a
      bare "as if", soft Thai markers used figuratively (เหมือน, เหมือนกับ,
      คล้ายกับ, สมมติว่า, เปรียบเทียบ), one idea per Thai sentence, the limit
      of three glossary terms, and whether every term is defined at first use.
    - Thai has no word boundaries, so Thai markers match as plain substrings;
      the P5 Thai list therefore also trips ทำการตลาด and ทำการบ้าน.
    - Only a line of dashes splits passes; "***", "___", and setext headings
      do not.
    - A closing line that is not recognised as the tail counts toward the cap
      (and "Let me know if" trips P2 by design).
    - The English sentence splitter over-splits abbreviations, which can only
      shorten sentences; an unterminated sentence ends at the line end.
    - Single-backtick spans only; an unbalanced backtick leaves the rest of
      the line unmasked.
    - Language is decided once per input; mixed text is counted in the
      dominant unit (Thai runs count as one word each in "en", English words
      count by character in "th").
    - P6 flags only obligation phrases; "the user sees a blank page" passes.
    - One report per rule per line; the rerun loop catches the rest.
"""

import argparse
import re
import sys
from pathlib import Path

DEFAULT_CAP = {"en": 120, "th": 500}
UNIT = {"en": "words", "th": "chars"}
MAX_SENTENCE_WORDS = 20
MAX_THAI_RUN = 90
THAI_RATIO = 0.30

_FENCE_RE = re.compile(r"^\s*(?:```|~~~)")
_CODE_SPAN_RE = re.compile(r"`[^`]*`")
_MARKER_RE = re.compile(r"^\s*(?:[-*+]|\d+[.)])\s+|^\s*#+\s+|^\s*>\s*")
_SPLIT_RE = re.compile(r"^\s*-{3,}\s*$")
_TAIL_RE = re.compile(r"^(?:say|reply|type|ask for)\b.*\bmore\b|^พิมพ์.*เพิ่ม", re.IGNORECASE)
_WORD_RE = re.compile(r"[^\W_]")
_NUMBERED_RE = re.compile(r"^\s*\d+[.)]\s")
_FOR_YOU_RE = re.compile(r"^(?:for you|สำหรับคุณ)\s*:", re.IGNORECASE)
_SENTENCE_RE = re.compile(r'(?<=[.!?])\s+|(?<=[.!?][)"”])\s+')

P1_MARK = re.compile(
    r"\bimagine\b|\bpicture (?:this|a|an|yourself)\b"
    r"|\bthink of (?:[^\s.,;:]+\s+){1,4}(?:as|like)\b|\bpretend\b"
    r"|\bmetaphor(?:ical(?:ly)?)?\b|\banalog(?:y|ies|ous)\b"
    r"|\bkind of like\b|\bsort of like\b|\ba bit like\b",
    re.IGNORECASE,
)
P1_CMP = re.compile(
    r"\b(?:is|are|was|were|be|it's|that's|works?|acts?|behaves?|just|much|"
    r"a lot|something|exactly)\s+like\s+(?:a|an)\b"
    r"|\bsimilar to\s+(?:a|an)\b|\bin the same way(?: that)?\s+(?:a|an)\b"
    r"|\blike when (?:you|a|an|someone|people)\b"
    r"|\bas if (?:it|this|that|you|the \w+) (?:were|was|had) (?:a|an)\b",
    re.IGNORECASE,
)
P1_TH = [
    "นึกภาพ", "จินตนาการ", "เปรียบเสมือน", "เปรียบเหมือน", "เปรียบได้กับ",
    "เปรียบดัง", "อุปมา", "เหมือนกับว่า", "เสมือนว่า", "ราวกับ", "เหมือนเวลา",
    "เหมือนตอนที่", "คล้ายๆ กับ", "คล้ายๆกับ", "ประมาณว่า", "ทำนองเดียวกับ",
]
P2_START = re.compile(
    r"^(?:sure|certainly|of course|absolutely|definitely)\b"
    r"|^happy to\b|^let me (?:explain|break)|^แน่นอน|^ยินดี|^ได้เลย",
    re.IGNORECASE,
)
P2_ANY = re.compile(
    r"\b(?:great|good|excellent) question\b|\bin summary\b"
    r"|\bto summari[sz]e\b|\bin conclusion\b|\b(?:i )?hope (?:this|that) helps\b"
    r"|\blet me know if\b|หวังว่าจะช่วย|หวังว่าจะเป็นประโยชน์|คำถามดี"
    r"|ขอบคุณสำหรับคำถาม|โดยสรุป|ถ้ามีคำถามเพิ่มเติม|หากมีคำถาม",
    re.IGNORECASE,
)
P5 = [
    (re.compile(pattern, re.IGNORECASE), fix)
    for pattern, fix in [
        (r"\butili[sz](?:e[sd]?|ing)\b", "use"),
        (r"\bleverag(?:e[sd]?|ing)\b", "use"),
        (r"\bfacilitat(?:e[sd]?|ing)\b", "help"),
        (r"\bin order to\b", "to"),
        (r"\bprior to\b", "before"),
        (r"\bsubsequent(?:ly)?\b", "later"),
        (r"\bit should be noted\b|\bit is important to note\b|\bit's worth noting\b",
         "just state the point"),
        (r"\baforementioned\b", "name the thing again"),
        (r"\b(?:thus|hence)\b", "so"),
        (r"\bwhereby\b", "where"),
        (r"\bvia\b", "through"),
        (r"\be\.g\.", "for example"),
        (r"\bi\.e\.", "that is"),
        (r"\b(?:additionally|furthermore|moreover)\b", "also"),
        (r"\bin the event that\b", "if"),
        (r"\bwith regards? to\b", "about"),
        (r"\bcommence[sd]?\b", "start"),
    ]
]
# ponytail: fixed high-precision list, substring matched; extend from real slips
P5_TH = [
    (re.compile(pattern), fix)
    for pattern, fix in [
        (r"ฐานข้อมูล", "database"),
        (r"แม่ข่าย", "server"),
        (r"ลูกข่าย", "client"),
        (r"คิวรีย่อย", "subquery"),
        (r"ส่วนต่อประสาน", "interface"),
        (r"(?:โปรแกรม|ซอฟต์แวร์)ประยุกต์", "app"),
        (r"ระเบียน", "record"),
        (r"ดำเนินการ", "ทำ"),
        (r"ทำการ", "the verb after it, alone"),
        (r"ในกรณีที่", "ถ้า"),
        (r"เพื่อที่จะ", "เพื่อ"),
        (r"กล่าวคือ", "คือ"),
        (r"ได้แก่", "คือ"),
        (r"อย่างไรก็ตาม", "แต่"),
        (r"ดังกล่าว", "the thing's name again"),
        (r"ทั้งนี้|อนึ่ง", "nothing, delete it"),
    ]
]
P6 = re.compile(
    r"\bthe readers?\b|\bthe user (?:should|must|needs? to|has to|will need to)\b"
    r"|ผู้อ่าน|ผู้ใช้(?:ควร|ต้อง|จะต้อง)",
    re.IGNORECASE,
)
P7 = re.compile(r";|—|(?<!\d)–(?!\d)|\s-\s|\s--\s")


def is_thai(ch):
    return "฀" <= ch <= "๿"


def prose_of(raw):
    """Mask inline code, strip list/heading/quote markers, drop bold marks."""
    text = _CODE_SPAN_RE.sub("0", raw)
    text = _MARKER_RE.sub("", text, count=1)
    return text.replace("*", "").strip()


def words(text):
    return sum(1 for token in text.split() if _WORD_RE.search(token))


def chars(text):
    return sum(1 for ch in text if not ch.isspace())


def check_line(prose, lineno, lang, report):
    m = P1_MARK.search(prose) or P1_CMP.search(prose)
    if m:
        report(lineno, "P1", f'analogy marker "{m.group(0)}"')
    else:
        for marker in P1_TH:
            if marker in prose:
                report(lineno, "P1", f'analogy marker "{marker}"')
                break
    m = P2_START.search(prose) or P2_ANY.search(prose)
    if m:
        report(lineno, "P2", f'filler "{m.group(0)}"')
    if lang == "en":
        for sentence in _SENTENCE_RE.split(prose):
            n = words(sentence)
            if n > MAX_SENTENCE_WORDS:
                report(lineno, "P3", f"sentence of {n} words (max {MAX_SENTENCE_WORDS})")
                break
    else:
        for run in prose.split():
            n = sum(1 for ch in run if is_thai(ch))
            if n > MAX_THAI_RUN:
                report(lineno, "P3",
                       f"Thai run of {n} characters without a space (max {MAX_THAI_RUN})")
                break
    for pattern, fix in P5 + (P5_TH if lang == "th" else []):
        m = pattern.search(prose)
        if m:
            report(lineno, "P5", f'use "{fix}" instead of "{m.group(0)}"')
            break
    m = P6.search(prose)
    if m:
        report(lineno, "P6", f'"{m.group(0)}": address the reader as you')
    m = P7.search(prose)
    if m:
        what = "semicolon" if m.group(0) == ";" else "dash"
        report(lineno, "P7", f"{what} in prose: split into two sentences")
    if "?" in prose:
        report(lineno, "P8", "question mark in prose: state it, or phrase the ask as an imperative")


def split_passes(text):
    """Return passes as lists of (lineno, raw, prose); fenced lines get empty prose."""
    passes = [[]]
    in_fence = False
    for lineno, raw in enumerate(text.splitlines(), 1):
        if _FENCE_RE.match(raw):
            in_fence = not in_fence
            passes[-1].append((lineno, raw, ""))
        elif in_fence:
            passes[-1].append((lineno, raw, ""))
        elif _SPLIT_RE.match(raw):
            passes.append([])
        else:
            passes[-1].append((lineno, raw, prose_of(raw)))
    return [p for p in passes if any(prose for _, _, prose in p)]


def detect_lang(passes):
    prose = " ".join(p for entries in passes for _, _, p in entries)
    thai = sum(1 for ch in prose if is_thai(ch))
    latin = sum(1 for ch in prose if ch.isascii() and ch.isalpha())
    return "th" if thai and thai / (thai + latin) >= THAI_RATIO else "en"


def check_text(text, label, lang_opt, cap_opt, out):
    count = 0

    def report(lineno, rule, msg):
        nonlocal count
        count += 1
        print(f"{label}:{lineno}: {rule} {msg}", file=out)

    passes = split_passes(text)
    if not passes:
        report(1, "P9", "no prose found")
        return count, "empty"
    lang = detect_lang(passes) if lang_opt == "auto" else lang_opt
    cap = cap_opt or DEFAULT_CAP[lang]
    measure = words if lang == "en" else chars
    used = []
    for index, entries in enumerate(passes, 1):
        prose_lines = [e for e in entries if e[2]]
        first = prose_lines[0]
        body = prose_lines[:-1] if _TAIL_RE.match(prose_lines[-1][2]) else prose_lines
        for lineno, _, prose in body:
            check_line(prose, lineno, lang, report)
        total = sum(measure(prose) for _, _, prose in body)
        used.append(total)
        if total > cap:
            report(first[0], "P4", f"pass {index}: {total} {UNIT[lang]} (cap {cap}, tail excluded)")
        if index == 1:
            if not first[1].strip().startswith("**"):
                report(first[0], "P9", "first line must be the bold one-line answer")
            numbered = sum(1 for _, raw, _ in entries if _NUMBERED_RE.match(raw))
            if not 2 <= numbered <= 5:
                report(first[0], "P9", f"{numbered} numbered step lines (need 2 to 5)")
            if not any(_FOR_YOU_RE.match(prose) for _, _, prose in body):
                report(first[0], "P9", 'missing "For you:" line')
    summary = f"{lang}, {'+'.join(str(n) for n in used)} {UNIT[lang]}"
    return count, summary


def main(argv=None):
    parser = argparse.ArgumentParser(
        description="Line-rule linter for plain, direct, analogy-free explanations "
                    "(reports violations, never rewrites)."
    )
    parser.add_argument("files", nargs="+", metavar="FILE",
                        help='markdown draft to check, or "-" for stdin')
    parser.add_argument("--lang", choices=("auto", "en", "th"), default="auto",
                        help="force the counting unit; default detects Thai")
    parser.add_argument("--cap", type=int, default=None,
                        help="budget per pass above the tail (default 120 words or 500 Thai chars)")
    args = parser.parse_args(argv)
    if args.cap is not None and args.cap <= 0:
        parser.error("--cap must be a positive integer")

    violations = 0
    summaries = []
    for name in args.files:
        if name == "-":
            text, label = sys.stdin.read(), "<stdin>"
        else:
            try:
                text = Path(name).read_text(encoding="utf-8")
            except OSError as exc:
                print(f"error: cannot read {name}: {exc}", file=sys.stderr)
                return 2
            label = name
        count, summary = check_text(text, label, args.lang, args.cap, sys.stderr)
        violations += count
        summaries.append(summary)

    if violations:
        return 1
    print(f"ok: {len(args.files)} input(s) clean ({'; '.join(summaries)})")
    return 0


if __name__ == "__main__":
    sys.exit(main())
