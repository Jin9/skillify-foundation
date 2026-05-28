# Validation gate

Before emitting `augmented_report` and `augment_notes`, confirm every check
below. Any failure → fix and re-run the relevant Procedure step; do not
emit a known-broken output.

1. **Insertion-only diff.** Every line of `final_report` appears in
   `augmented_report` in the same relative order. The only new content
   is fenced code blocks and the single blank lines on either side of
   each block.
2. **Diagram count gate.** Zero fenced code blocks on `quick` depth.
   ≤ 3 newly inserted fenced blocks on `standard` / `deep`.
3. **Fence validity.** Every triple-backtick opens and closes; no
   nested fences; no language tag on any newly inserted block.
4. **No `[n]` inside any newly inserted fence.** Citations live in
   surrounding prose.
5. **Character set.** Every newly inserted fenced block contains only
   the allowed character set: ASCII printable plus
   `─ │ ┌ ┐ └ ┘ ├ ┤ ┬ ┴ ┼ → ← ↔` (and spaces / newlines).
6. **Line and column caps.** Every newly inserted fenced block is ≤ 15
   lines and ≤ 70 columns per line.
7. **Placement.** No newly inserted fenced block sits in
   `## Executive Summary`, `## Background`, `## Methodology`,
   `## Key Findings`, `## Limitations & Open Questions`, or
   `## Sources`. On `standard` depth, no block sits in `## Synthesis`
   (that section is `deep`-only).
8. **Grounding.** Every entity / arrow / layer in each newly inserted
   block names a finding the section already cites (transitively via
   `cited_sources[i].findings_supported`). No diagram entity is novel
   to the diagram.
9. **Sources section byte-identical.** The `## Sources` heading and
   everything after it appears in `augmented_report` byte-for-byte
   identical to `final_report`.
10. **Topic / title untouched.** The `# H1` line at the top of the
    report is byte-identical between input and output.

## Why these ten

Items 1 and 9–10 enforce the insertion-only contract — the single most
important invariant this skill maintains. Items 2 and 7 enforce the depth
gate and section gate. Items 3, 5, and 6 enforce the fence shape. Items 4
and 8 ensure diagrams do not smuggle in new claims via the structural
side-door.
