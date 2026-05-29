# Validation gate

Before emitting `draft_report` and `cited_sources`, confirm all of the
following. Any failure → fix and re-run the relevant Procedure step; do not
emit a known-broken draft.

1. **Title is the input `topic` verbatim** in an `# H1`.
2. **Section set matches depth** (per
   [`output-contract.md`](output-contract.md)).
3. **Every `[n]` marker** in `draft_report` has a matching entry in
   `cited_sources` with that `id`.
4. **Every `cited_sources` entry** has at least one `[n]` reference in
   `draft_report`.
5. **`## Sources` section** in `draft_report` agrees with `cited_sources`
   (same ids, same order, same titles/urls).
6. **Every `cited_sources[i].findings_supported`** is non-empty and every
   referenced finding id exists in the input `findings`.
7. **Every selected finding** is cited at least once OR explicitly
   accounted for in `## Limitations & Open Questions`.
8. **Every `research_plan.sub_questions` entry** is addressed by at least
   one body section OR an explicit "no evidence found" note.
9. **Audience tone check** — no `general` report contains undefined jargon
   in body prose; no `executive` report defers the recommendation past
   the executive summary; no `expert` report omits mechanism / parameter
   detail when the findings supply it.
10. **Length is within ±50%** of the depth target (300–600 / 800–1500 /
    2000–4000 words). Over-budget → cut. Severely under → only OK if
    findings genuinely don't support more, and Limitations must say so.
11. **No invented sources** — every `cited_sources[i].url` and `title`
    appears verbatim in the input source records (dereferenced via
    `findings`).
12. **Diagram count and depth gate.** Zero fenced code blocks on `quick`
    depth. ≤ 3 fenced blocks total on `standard` / `deep`.
13. **Diagram fences are valid markdown.** Every triple-backtick opens
    and closes; no nested fences; no language tag on diagram blocks.
14. **No `[n]` markers inside any fenced block.** Citations live in the
    surrounding prose, not in the diagram.
15. **Diagram placement.** No fenced blocks in `## Executive Summary`,
    `## Background`, `## Methodology`, `## Key Findings`, `## Limitations
    & Open Questions`, or `## Sources`. Diagrams live only in themed body
    sections or `## Synthesis`.
16. **Diagram grounding.** Every entity / arrow / layer in each diagram
    traces to a finding selected for the same section.
17. **Every declarative sentence in `draft_report`** traces to at least one
    finding in the input `findings`. Cross-finding synthesis appears only in
    a labeled `## Synthesis` subsection (deep depth only), with each
    synthesis sentence citing its underlying findings.
