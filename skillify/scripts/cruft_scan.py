#!/usr/bin/env python3
"""Advisory scanner for dated prompt patterns ("cruft") in a skill folder.

Cruft is relative to a model generation: instructions that helped an older
model (pressure language, "think step by step" scaffolds, narration
suppressors, pinned model names, step choreography for judgment work) lower
output quality on current frontier models. This scanner finds the greppable
forms and reports them with a confidence level. It is ADVISORY by default
(exit 0): every hit is a hypothesis to classify against the keep list in
references/model-generation-fit.md, not an error. --strict exits 1 on
High-confidence findings (or Medium with --fail-on medium) and is meant only as
the exit check of the Re-baseline sub-flow.

Keep-list exemptions built in:
  - fenced code blocks are skipped (--include-fenced to scan them)
  - frontmatter is scanned only for pinned model names
  - text inside `backticks` or "double quotes" is a mention, not an instruction
  - sections headed "When (not) to use", "Scope", "Trigger", "Do not use" are
    exempt from pressure and prohibition signals (trigger text may carry urgency)
  - table rows downgrade fossil/brittle signals to Low and skip scaffold signals
    entirely (a detection table names the pattern as data, not as an instruction)
  - a prohibition that carries a reason (";", "because", "so that", "otherwise",
    "would", "since", "to avoid") does not count toward a prohibition wall
  - inline suppression: <!-- cruft-scan: allow <signal.id> --> on the line or
    the line above

Usage:
  python3 cruft_scan.py <skill_dir> [<skill_dir> ...] [--strict]
      [--fail-on {high,medium}] [--format {text,json}] [--allow ID ...]
      [--include-examples] [--include-fenced] [--self-test]
Exit codes: 0 advisory run or strict pass; 1 strict fail or self-test fail;
2 usage error or missing path. stdlib only.
"""

from __future__ import annotations

import argparse
import json
import re
import sys
import tempfile
from pathlib import Path

LEVELS = {"high": 3, "medium": 2, "low": 1}
FLAGS = re.I

# signal id -> (group, confidence, action, why)
SIGNALS = {
    "pressure.density": ("1a", "medium", "rewrite",
        "Capitalised emphasis over-applies on current models; state the one or two real constraints plainly, with the reason."),
    "pressure.no_reason": ("1a", "low", "flag",
        "An emphasised rule with no adjacent reason reads as register, not information; add the reason or drop the emphasis."),
    "pressure.hedge": ("1a", "medium", "rewrite",
        "Hedges on real requirements are read literally as permission to under-deliver; state the requirement."),
    "pressure.trait_claim": ("1a", "medium", "rewrite",
        "Trait claims about the model describe a past generation; state the desired behavior instead."),
    "scaffold.think": ("1b", "high", "remove",
        "Current models plan and reason natively; thinking scaffolds cause over-planning and are redundant at best."),
    "scaffold.show_reasoning": ("1b", "high", "remove",
        "Asking the model to reproduce its reasoning can trigger a refusal on current models; ask for the verified result and its evidence."),
    "scaffold.cadence": ("1b", "high", "remove",
        "Fixed narration cadences were tuned against older verbosity; ask for an opening line and a standalone recap instead."),
    "scaffold.numeric_cap": ("1b", "medium", "rewrite",
        "Numeric output caps starve reasoning on hard problems; prefer qualitative length guidance."),
    "overspec.prohibition_run": ("1c", "medium", "rewrite",
        "A wall of unreasoned prohibitions anchors the model toward the failures it names; keep the ones that guard a current failure or real policy and say why."),
    "overspec.reinforcement": ("1c", "medium", "remove",
        "Repetition as reinforcement makes the model reconcile wordings; say it once, in the right place."),
    "overspec.grader": ("1c", "high", "remove",
        "Describing the grader pushes effort toward being watched; state every requirement the grader checks instead."),
    "overspec.virtue": ("1c", "low", "flag",
        "Generic virtues restate trained defaults; delete unless the line carries a specific standard."),
    "fossil.model_name": ("1d", "medium", "rewrite",
        "Pinned model names rot at the next release; route by model-cost tier and effort hint, and keep names only as dated examples in data tables."),
    "fossil.suppressor": ("1d", "high", "rewrite",
        "Narration suppressors were written against chatty models; current models under-narrate with these present."),
    "fossil.anti_format": ("1d", "high", "remove",
        "Anti-formatting rules were written against over-formatting models; current models under-format, so the rule strips structure the reader wanted."),
    "fossil.reminder_cadence": ("1d", "medium", "remove",
        "Instruction re-insertion is a retention crutch; current models retain a once-stated instruction."),
    "fossil.date_conditional": ("1d", "medium", "rewrite",
        "Date-conditional guidance rots; state the current rule as the only rule."),
    "fossil.relative_phrasing": ("1d", "low", "flag",
        "Migration-relative phrasing describes a diff the model never saw; write the current rule directly."),
    "fossil.identity_stub": ("1d", "low", "flag",
        "An identity stub is fine as a focus-setter but not as a substitute for audience, product, and quality bar."),
    "brittle.hardcoded_path": ("2", "medium", "rewrite",
        "Absolute machine paths rot and break portability; use repo-relative paths or a named location the host provides."),
    "brittle.version_pin": ("2", "low", "flag",
        "A version pin in a rule line rots as tooling ships; re-verify or move it to a dated data table."),
    "brittle.incident": ("2", "low", "flag",
        "History narratives carry the incident, not the behavior; state the current rule."),
    "contract.missing": ("add", "medium", "add",
        "Current models pause or diverge when a skill and the request disagree; add the operating contract (templates/operating-contract.md) with its instruction-priority, autonomy, stop, verification, delegation, and progress lines."),
    "tier.missing": ("add", "low", "add",
        "No model-cost tier is declared; add a 'model-cost tier:' line or inline [tier/effort] tags so hosts can route by tier, never by model name."),
    "vendor.body_mention": ("house", "medium", "rewrite",
        "A vendor or model name in body prose breaks portability across hosts; say 'the agent' and route by tier."),
}

CAPS_RE = re.compile(r"\b(MUST|NEVER|ALWAYS|CRITICAL|IMPORTANT)\b")
REASON_RE = re.compile(r"(;|\bbecause\b|\bso that\b|\botherwise\b|\bwould\b|\bsince\b|\bto avoid\b|\bthe reason\b)", FLAGS)
PROHIBIT_RE = re.compile(r"^\s*(?:[-*]|\d+\.)\s*(?:\*\*)?(?:do not|don'?t|never|avoid|must not)\b", FLAGS)
LIST_RE = re.compile(r"^\s*(?:[-*]|\d+\.)\s+")
EXEMPT_HEADING_RE = re.compile(r"^#{1,6}\s*(when (not )?to use|scope|trigger|do not use|when not to use)", FLAGS)
HEADING_RE = re.compile(r"^#{1,6}\s")
TABLE_RE = re.compile(r"^\s*\|")
FENCE_RE = re.compile(r"^\s*(```|~~~)")
ALLOW_RE = re.compile(r"<!--\s*cruft-scan:\s*allow\s+([a-z_.]+)[^>]*-->", FLAGS)
CONTRACT_KEYS = ["Instruction priority", "Autonomy", "Stop conditions", "Verification", "Progress"]
OPTIONAL_CONTRACT_KEYS = ["Delegation"]  # the template allows deleting this line when the skill never delegates

LINE_PATTERNS = {
    "pressure.hedge": (re.compile(r"\b(try to|if possible|ideally|where possible)\b", FLAGS), "list"),
    "pressure.trait_claim": (re.compile(r"\byou (tend to|often|sometimes)\b|\bdon'?t be too \w+", FLAGS), "any"),
    "scaffold.think": (re.compile(r"think step[- ]by[- ]step|take a deep breath|<scratchpad>|<thinking>|use the think tool|\bplan before acting\b", FLAGS), "any"),
    "scaffold.show_reasoning": (re.compile(r"(show|explain|reveal|include|output)\s+(all\s+of\s+)?your\s+(full\s+)?(reasoning|thinking|chain of thought)|\bchain[- ]of[- ]thought\b.*\b(before|first)\b", FLAGS), "any"),
    "scaffold.cadence": (re.compile(r"\bevery \d+ (tool calls?|messages?|steps?|turns?)\b", FLAGS), "any"),
    "scaffold.numeric_cap": (re.compile(r"\b(at most|no more than|under|maximum of|fewer than|limit(?:ed)? to)\s+\d+\s+(words|sentences|bullets|bullet points)\b", FLAGS), "any"),
    "overspec.reinforcement": (re.compile(r"\b(Remember,|Again,|As stated above|As mentioned (above|earlier))", FLAGS), "any"),
    "overspec.grader": (re.compile(r"\bhidden tests?\b|\byou will be (graded|scored|evaluated)\b", FLAGS), "any"),
    "overspec.virtue": (re.compile(r"^\s*[-*]?\s*(Always )?[Bb]e (accurate|thorough|helpful|careful|clear|concise)\b.{0,40}$", FLAGS), "any"),
    "fossil.model_name": (re.compile(r"\b(claude-2|claude-3(?:[.-]\d)?|claude-instant|(?:claude|sonnet|opus|haiku)[ -]?3\.[57]|gpt-4(?:[.-]\w+)?|gpt-5(?:\.\d+)?(?:-\w+)?|opus 4(?:\.\d)?|sonnet 4(?:\.\d)?|haiku 4(?:\.\d)?|gemini[ -]?2(?:\.\d)?|gemini[ -]?3(?:\.\d)?)\b", FLAGS), "any"),
    "fossil.suppressor": (re.compile(r"\bhold (all )?(findings|results)\b|\bdon'?t narrate\b|\bno interim (updates?|progress)\b|\bdo not (narrate|report progress)\b", FLAGS), "any"),
    "fossil.anti_format": (re.compile(r"\bnever use (bullets|bullet points|headers|headings|bold|markdown)\b|\bno (bullets|bullet points|headers|headings|bold)\b|\bdo not use (bullets|bullet points|headers|headings|bold|markdown)\b", FLAGS), "any"),
    "fossil.date_conditional": (re.compile(r"\bif\b.{0,30}\b(before|after) \d{4}\b", FLAGS), "any"),
    "fossil.relative_phrasing": (re.compile(r"\b(no longer|used to|previously|as of this version)\b.{0,40}\b(need|needs|needed|require|required|support|supported|use|used|accept|accepted|applies|work|works)\b", FLAGS), "list"),
    "fossil.identity_stub": (re.compile(r"^You are (a|an) (helpful|expert)\b"), "any"),
    "brittle.hardcoded_path": (re.compile(r"(/Users/[A-Za-z0-9_.-]+/|/home/[A-Za-z0-9_.-]+/|\b[A-Z]:\\)"), "any"),
    "brittle.version_pin": (re.compile(r"\bv?\d+\.\d+\.\d+\b"), "list"),
    "brittle.incident": (re.compile(r"\b(PR|MR|issue) #?\d+\b|\b(we|this) (once|previously) (saw|hit|had)\b", FLAGS), "any"),
    "vendor.body_mention": (re.compile(r"\bClaude\b(?! Code)(?!\.md)|\bAnthropic\b|\bOpenAI\b(?! Codex)|\bChatGPT\b"), "any"),
}
# signals where a mention inside quotes/backticks is talk-about, not an instruction
MENTION_SENSITIVE = {"scaffold.think", "scaffold.show_reasoning", "scaffold.cadence", "scaffold.numeric_cap",
                     "fossil.suppressor", "fossil.anti_format", "fossil.model_name", "vendor.body_mention",
                     "overspec.reinforcement", "overspec.grader", "pressure.trait_claim", "pressure.hedge",
                     "fossil.relative_phrasing", "overspec.virtue"}
MODEL_FILES = {"fossil.model_name"}


def strip_mentions(line: str) -> str:
    line = re.sub(r"`[^`]*`", "`x`", line)
    line = re.sub(r'"[^"]*"', '"x"', line)
    line = re.sub(r"“[^”]*”", "“x”", line)
    return line


def split_frontmatter(text: str):
    if text.startswith("---"):
        m = re.match(r"^---\n(.*?)\n---\n?", text, re.S)
        if m:
            fm_lines = m.group(1).count("\n") + 3
            return m.group(1), text[m.end():], fm_lines
    return "", text, 0


class Finding:
    def __init__(self, file, line, signal, evidence, confidence=None):
        self.file, self.line, self.signal, self.evidence = file, line, signal, evidence
        group, conf, action, why = SIGNALS[signal]
        self.group, self.action, self.why = group, action, why
        self.confidence = confidence or conf

    def as_dict(self):
        return {"file": self.file, "line": self.line, "signal": self.signal, "group": self.group,
                "confidence": self.confidence, "action": self.action, "evidence": self.evidence, "why": self.why}


def scan_file(path: Path, rel: str, include_fenced: bool, is_skill_md: bool):
    text = path.read_text(encoding="utf-8", errors="replace")
    fm, body, fm_offset = split_frontmatter(text)
    findings = []
    allowed_next = set()

    # frontmatter: pinned model names only
    for i, ln in enumerate(fm.splitlines(), start=2):
        pat, _ = LINE_PATTERNS["fossil.model_name"]
        if pat.search(strip_mentions(ln)):
            findings.append(Finding(rel, i, "fossil.model_name", ln.strip()[:120], "low"))

    lines = body.splitlines()
    in_fence = False
    exempt = False
    caps_total = 0
    reminder_count = 0
    body_lines = 0
    run = []  # consecutive unreasoned prohibition lines (line numbers)
    heading_hits = 0
    numbered = 0
    has_contract = False
    contract_keys_found = set()
    has_tier = False
    tier_re = re.compile(r"model-cost tier:\s*(small|mid|frontier)|\[(small|mid|frontier)(/(low|medium|high))?\]", FLAGS)

    def flush_run():
        nonlocal run
        if len(run) >= 3:
            findings.append(Finding(rel, run[0], "overspec.prohibition_run",
                                    "%d consecutive prohibition lines without a stated reason (lines %d-%d)" % (len(run), run[0], run[-1])))
        run = []

    for idx, raw in enumerate(lines):
        lineno = idx + 1 + fm_offset
        if FENCE_RE.match(raw):
            in_fence = not in_fence
            flush_run()
            continue
        if in_fence and not include_fenced:
            continue
        line_allows = set(ALLOW_RE.findall(raw)) | allowed_next
        allowed_next = set(ALLOW_RE.findall(raw)) if ALLOW_RE.search(raw) and not raw.strip().startswith("<!--") is False else set()
        if raw.strip().startswith("<!--") and ALLOW_RE.search(raw):
            allowed_next = set(a.lower() for a in ALLOW_RE.findall(raw))
            continue
        if HEADING_RE.match(raw):
            if tier_re.search(raw):
                has_tier = True
            exempt = bool(EXEMPT_HEADING_RE.match(raw))
            heading_hits += 1 if re.search(r"\b(step|phase|workflow)\b", raw, FLAGS) else 0
            if re.match(r"^##\s+Operating contract\s*$", raw.strip()):
                has_contract = True
            flush_run()
            continue
        if not raw.strip():
            flush_run()
            continue
        body_lines += 1
        is_table = bool(TABLE_RE.match(raw))
        is_list = bool(LIST_RE.match(raw))
        if is_list and re.match(r"^\s*\d+\.", raw):
            numbered += 1
        if tier_re.search(raw):
            has_tier = True
        for key in CONTRACT_KEYS + OPTIONAL_CONTRACT_KEYS:
            if re.match(r"^\s*-\s*%s:" % re.escape(key), raw):
                contract_keys_found.add(key)
        if re.match(r"^\s*reminder:", raw, FLAGS):
            reminder_count += 1

        # prohibition wall (skip exempt sections and tables)
        if not exempt and not is_table and PROHIBIT_RE.match(raw) and not REASON_RE.search(raw):
            run.append(lineno)
        else:
            flush_run()

        # pressure counters (skip exempt sections)
        if not exempt and not is_table:
            caps = CAPS_RE.findall(strip_mentions(raw))
            caps_total += len(caps)
            if "!!" in raw:
                caps_total += 1
            if caps and not REASON_RE.search(raw) and "pressure.no_reason" not in line_allows:
                findings.append(Finding(rel, lineno, "pressure.no_reason", raw.strip()[:120]))

        for sig, (pat, where) in LINE_PATTERNS.items():
            if sig in line_allows:
                continue
            if exempt and sig.startswith(("pressure.", "overspec.")):
                continue
            if where == "list" and not is_list:
                continue
            target = strip_mentions(raw) if sig in MENTION_SENSITIVE else raw
            m = pat.search(target)
            if not m:
                continue
            conf = None
            if is_table and sig.startswith("scaffold."):
                continue  # a detection or reference table names the pattern as data, not as an instruction
            if sig == "fossil.model_name":
                conf = "low" if is_table or not re.search(r"\b(use|run|route|prefer|only|on|with|via)\b", raw, FLAGS) else "medium"
            elif is_table and sig.startswith(("fossil.", "brittle.")):
                conf = "low"
            findings.append(Finding(rel, lineno, sig, raw.strip()[:120], conf))
    flush_run()

    if caps_total >= 3 and body_lines and (caps_total * 100.0 / body_lines) >= 2.0:
        findings.append(Finding(rel, 0, "pressure.density",
                                "%d capitalised emphasis markers in %d body lines" % (caps_total, body_lines)))
    if reminder_count >= 2:
        findings.append(Finding(rel, 0, "fossil.reminder_cadence", "%d 'reminder:' lines" % reminder_count))
    if is_skill_md:
        multi_step = numbered >= 2 or heading_hits >= 2
        missing = [k for k in CONTRACT_KEYS if k not in contract_keys_found]
        if not has_contract or missing:
            ev = "no '## Operating contract' section" if not has_contract else "contract missing keys: " + ", ".join(missing)
            findings.append(Finding(rel, 0, "contract.missing", ev, None if multi_step else "low"))
        if not has_tier:
            findings.append(Finding(rel, 0, "tier.missing", "no 'model-cost tier:' line or [tier/effort] tag"))
    return findings


def scan_skill(skill_dir: Path, include_examples: bool, include_fenced: bool, allow: set):
    files = []
    skill_md = skill_dir / "SKILL.md"
    if skill_md.is_file():
        files.append(skill_md)
    for sub in ("references", "templates") + (("examples",) if include_examples else ()):
        d = skill_dir / sub
        if d.is_dir():
            files.extend(sorted(p for p in d.glob("*.md") if p.is_file()))
    findings = []
    for f in files:
        rel = f.relative_to(skill_dir).as_posix()
        findings.extend(scan_file(f, rel, include_fenced, f.name == "SKILL.md" and f.parent == skill_dir))
    findings = [x for x in findings if x.signal not in allow]
    findings.sort(key=lambda x: (-LEVELS[x.confidence], x.file, x.line))
    return files, findings


def summarize(findings):
    s = {"high": 0, "medium": 0, "low": 0}
    groups = {}
    for f in findings:
        s[f.confidence] += 1
        g = f.signal.split(".")[0]
        groups[g] = groups.get(g, 0) + 1
    return s, groups


def render_text(skill_dir, files, findings, strict, fail_on):
    out = ["cruft_scan: %s" % skill_dir]
    for f in findings:
        loc = "%s:%d" % (f.file, f.line) if f.line else f.file
        out.append("%-28s [%-6s] %-26s %s" % (loc, f.confidence.capitalize(), f.signal, f.action))
        out.append("    evidence: %s" % f.evidence)
        out.append("    why: %s" % f.why)
    s, groups = summarize(findings)
    out.append("summary: high=%d medium=%d low=%d files=%d (%s)" % (
        s["high"], s["medium"], s["low"], len(files),
        " ".join("%s=%d" % kv for kv in sorted(groups.items())) or "no findings"))
    failed = strict and any(LEVELS[f.confidence] >= LEVELS[fail_on] for f in findings)
    if strict:
        n = sum(1 for f in findings if LEVELS[f.confidence] >= LEVELS[fail_on])
        out.append("result: strict %s, %d finding(s) at %s or above (exit %d)" % ("FAIL" if failed else "pass", n, fail_on, 1 if failed else 0))
    else:
        out.append("result: advisory (exit 0)")
    return "\n".join(out), failed


# ---------------------------------------------------------------- self-test
SELF_TEST = {
    "pressure.density": ("## Rules\n- You MUST do A\n- You MUST do B\n- NEVER do C\n- ALWAYS do D\n", "## Rules\n- Do A, because B.\n- Do C so that D.\n"),
    "pressure.hedge": ("## Rules\n- Try to include a summary\n", "## Rules\n- Include a summary.\n"),
    "pressure.trait_claim": ("## Rules\nYou tend to over-explain, so keep it short.\n", "## Rules\nKeep responses to the length the question needs.\n"),
    "scaffold.think": ("## Rules\nThink step by step before answering.\n", "## Rules\nDo not add \"think step by step\" scaffolds; the agent plans natively.\n"),
    "scaffold.show_reasoning": ("## Rules\nShow your reasoning before acting.\n", "## Rules\nReport the verified result and quote its evidence.\n"),
    "scaffold.cadence": ("## Rules\nSummarize progress every 3 tool calls.\n", "## Rules\nOpen with one line and close with a standalone recap.\n"),
    "scaffold.numeric_cap": ("## Rules\nRespond in at most 120 words.\n", "## Rules\nBe selective about what you include.\n"),
    "overspec.prohibition_run": ("## Rules\n- Do not use tabs\n- Never use spaces\n- Avoid comments\n", "## Rules\n- Do not use tabs; the linter rejects them.\n- Never force-push; history would be lost.\n- Avoid comments, because the docs cover it.\n"),
    "overspec.reinforcement": ("## Rules\nRemember, always validate input.\n", "## Rules\nValidate input at system boundaries.\n"),
    "overspec.grader": ("## Rules\nYou will be graded on completeness.\n", "## Rules\nCover every stated requirement.\n"),
    "overspec.virtue": ("## Rules\n- Be accurate and thorough.\n", "## Rules\n- Cite the line number for every finding.\n"),
    "fossil.model_name": ("## Rules\nUse claude-3-5-sonnet for extraction.\n", "## Rules\nUse a small-tier model for extraction; `gpt-5.5` is only a dated example.\n"),
    "fossil.suppressor": ("## Rules\nHold all findings for the final response.\n", "## Rules\nRemove \"hold all findings\" lines; current models under-narrate.\n"),
    "fossil.anti_format": ("## Rules\nNever use bullets in replies.\n", "## Rules\nUse lists when the content is genuinely parallel.\n"),
    "fossil.reminder_cadence": ("## Rules\nreminder: check the log\nreminder: check the log again\n", "## Rules\nCheck the log once before reporting.\n"),
    "fossil.date_conditional": ("## Rules\nIf the run is after 2025, use the new API.\n", "## Rules\nUse the current API.\n"),
    "fossil.relative_phrasing": ("## Rules\n- The tool no longer needs a token.\n", "## Rules\n- The tool needs no token.\n"),
    "fossil.identity_stub": ("You are a helpful assistant.\n", "You are the release-notes writer for this repository; the audience is on-call engineers.\n"),
    "brittle.hardcoded_path": ("## Rules\nRead /Users/alice/notes.md first.\n", "## Rules\nRead `notes.md` in the working directory first.\n"),
    "brittle.version_pin": ("## Rules\n- Requires tool v1.2.3 exactly.\n", "## Rules\n- Requires the tool version listed in `references/versions.md`.\n"),
    "brittle.incident": ("## Rules\nWe once saw PR #123 break the build, so run tests.\n", "## Rules\nRun the tests before reporting.\n"),
    "vendor.body_mention": ("## Rules\nAsk Claude to run the command.\n", "## Rules\nRun the command. Say \"Run the command\" instead of \"Ask Claude to run the command\". Claude Code is a supported host; see `CLAUDE.md`.\n"),
    "contract.missing": ("## Workflow\n1. Read\n2. Write\n3. Verify\n", "## Workflow\n1. Read\n2. Write\n\n## Operating contract\n\n- Instruction priority: the user wins.\n- Autonomy: act.\n- Stop conditions: none.\n- Verification: quote the PASS line.\n- Delegation: none.\n- Progress: opening line and recap.\n- Model-cost tier: small.\n"),
    "tier.missing": ("## Workflow\n1. Read\n", "## Workflow\nEvery step is model-cost tier: small.\n1. Read\n"),
}
FM = "---\nname: probe-skill\ndescription: Probe. Use when testing the scanner.\n---\n\n# Probe\n\n"


def self_test() -> int:
    failures = []
    with tempfile.TemporaryDirectory() as td:
        for sig, (hit, control) in SELF_TEST.items():
            for kind, body in (("hit", hit), ("control", control)):
                d = Path(td) / ("%s-%s" % (sig.replace(".", "-"), kind))
                d.mkdir()
                (d / "SKILL.md").write_text(FM + body, encoding="utf-8")
                _, findings = scan_skill(d, False, False, set())
                got = {f.signal for f in findings}
                if kind == "hit" and sig not in got:
                    failures.append("%s: expected a hit, got %s" % (sig, sorted(got)))
                if kind == "control" and sig in got:
                    failures.append("%s: control should not hit (got %s)" % (sig, sorted(got)))
    if failures:
        print("self-test FAIL:\n  " + "\n  ".join(failures))
        return 1
    print("self-test ok: %d signals, hit and control each" % len(SELF_TEST))
    return 0


def main() -> int:
    ap = argparse.ArgumentParser(description="Advisory dated-pattern scanner for skill folders.")
    ap.add_argument("skill_dirs", nargs="*")
    ap.add_argument("--strict", action="store_true", help="exit 1 when findings at --fail-on or above exist")
    ap.add_argument("--fail-on", choices=("high", "medium"), default="high")
    ap.add_argument("--format", choices=("text", "json"), default="text")
    ap.add_argument("--allow", nargs="*", default=[], help="signal ids to suppress")
    ap.add_argument("--include-examples", action="store_true")
    ap.add_argument("--include-fenced", action="store_true")
    ap.add_argument("--self-test", action="store_true")
    args = ap.parse_args()
    if args.self_test:
        return self_test()
    if not args.skill_dirs:
        ap.print_usage(sys.stderr)
        return 2
    allow = set(a.lower() for a in args.allow)
    any_failed = False
    reports = []
    for raw in args.skill_dirs:
        d = Path(raw).expanduser().resolve()
        if not d.is_dir() or not (d / "SKILL.md").is_file():
            print("error: not a skill folder (no SKILL.md): %s" % d, file=sys.stderr)
            return 2
        files, findings = scan_skill(d, args.include_examples, args.include_fenced, allow)
        s, groups = summarize(findings)
        failed = args.strict and any(LEVELS[f.confidence] >= LEVELS[args.fail_on] for f in findings)
        any_failed = any_failed or failed
        if args.format == "json":
            reports.append({"skill_dir": str(d), "findings": [f.as_dict() for f in findings],
                            "summary": {**s, "files": len(files), "groups": groups},
                            "strict": args.strict, "failed": failed})
        else:
            text, _ = render_text(d, files, findings, args.strict, args.fail_on)
            print(text)
            if len(args.skill_dirs) > 1:
                print()
    if args.format == "json":
        print(json.dumps(reports if len(reports) > 1 else reports[0], indent=2))
    return 1 if any_failed else 0


if __name__ == "__main__":
    sys.exit(main())
