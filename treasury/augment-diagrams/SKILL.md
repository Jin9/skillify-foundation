---
name: augment-diagrams
description: >
  OPTIONAL out-of-band stage (NOT part of the standard 1–6 pipeline). Backfill
  0–3 inline ASCII diagrams into an existing final research report,
  insertion-only, without rewriting any prose. Receives `final_report`,
  `findings`, `cited_sources`, `research_plan`, `topic`, and (optionally)
  `audience`; emits `augmented_report` (same prose plus fenced diagrams in
  body sections) and `augment_notes` (per-section decision log). Use when the
  user says "add diagrams to this report", "backfill ASCII diagrams", "augment
  with diagrams", or to bring pre-v0.7.1 runs up to the in-report-diagram
  contract introduced in `synthesize-report@0.2.0`. Do NOT use to draft a
  fresh report (synthesize-report), audit prose (review-report), invent
  findings or sources, or modify any existing sentence in the input report.
compatibility: claude-code, codex, copilot, gemini, antigravity
allowed-tools: ""  # no external tools — text-in / text-out (the driver script handles file I/O around the skill)
version: 0.1.2
stage: augment-diagrams
workflow: researcher
inputs:
  - { name: final_report,   type: string, required: true,  source: artifact.05-final_report.md,           description: "existing markdown report; treated as immutable prose" }
  - { name: findings,       type: array,  required: true,  source: artifact.03-findings.json,             description: "from extract-findings; objects with id, claim, evidence, source_id, confidence" }
  - { name: cited_sources,  type: array,  required: true,  source: artifact.05-cited_sources.json,        description: "from review-report; id ↔ {title, url, findings_supported[]}" }
  - { name: research_plan,  type: object, required: true,  source: artifact.01-research_plan.json,        description: "from plan-research; this skill reads `.depth` and `.sub_questions`" }
  - { name: topic,          type: string, required: true,  source: workflow.inputs.topic,                 description: "the subject; used only for tie-breaking and notes" }
  - { name: audience,       type: string, required: false, source: workflow.inputs.audience,              description: "audience used for the run; for log only, does not change diagram rules" }
outputs:
  - { name: augmented_report, type: string, description: "input report verbatim with 0–3 fenced ASCII diagrams inserted into body sections; insertion-only diff against input" }
  - { name: augment_notes,    type: string, description: "markdown decision log: per candidate section, the diagram type chosen or skip reason, grounding finding ids, scan-test rationale" }
---

# Skill: augment-diagrams

## Purpose

Take an already-final research report and add 0–3 inline ASCII diagrams to
its body sections **without rewriting a single sentence**. The skill is an
**inserter**, not an author: it reads the existing prose, finds body
sections where a diagram would carry structural signal the prose does not
already make scannable, and emits the same report with diagrams interleaved
after the first prose paragraph of each chosen section.

**Atomicity:** one stage, one LLM call. The skill is responsible for diagram
generation + placement, nothing else. It does not re-draft, re-cite, or
re-review.

**Why it exists separately from `synthesize-report`:** the diagram rules
were added to `synthesize-report@0.2.0` (squad-researcher v0.7.1). Runs
produced before that point have no diagrams. Re-running
`synthesize-report` on those runs would regenerate the prose from
`findings` — losing the audited final wording that `review-report` already
shipped. This skill backfills diagrams against the audited prose instead.

## When to use this skill

- Backfilling diagrams into a pre-v0.7.1 run whose `05-final_report.md` you
  want to preserve verbatim.
- Standalone invocation against any run directory whose `03-findings.json`,
  `05-cited_sources.json`, and `01-research_plan.json` are present.
- Trigger phrases: "add diagrams to this report", "backfill ASCII
  diagrams", "augment with diagrams", "insert diagrams into the final
  report".

## Do NOT use this skill to

- Draft a fresh report — use `synthesize-report`. New runs get diagrams
  inline during stage 4; this skill is for runs that already passed
  stage 5.
- Audit, edit, or otherwise modify the report's prose — use
  `review-report`. This skill's output preserves prose byte-for-byte.
- Invent findings, sources, or `[n]` markers. All diagram entities must
  trace to findings already cited by the input report.
- Edit the `## Sources` numbered list. The bijection
  `[n] ↔ cited_sources[].id` is fixed at input time.
- Add diagrams to `quick`-depth reports. Same gate as `synthesize-report`.

## Inputs

The six inputs and their run-directory sources are declared in the
frontmatter `inputs:` block. The skill reads `research_plan.depth` to gate
output and uses `cited_sources[i].findings_supported[]` joined to
`findings[]` to resolve each section's grounding set (see Procedure Step 3).

## Outputs

### `augmented_report` — markdown string

The input `final_report` **byte-for-byte**, with 0–3 fenced code blocks
inserted. Each inserted block is preceded and followed by exactly one
blank line so the surrounding prose's paragraph structure is unaffected.
Every other line of the input — every sentence, every `[n]` marker,
every heading, every blank line, the `## Sources` list — appears in the
output unchanged and in the same order.

#### Diagrams

**Canonical rules:** the When, Count, Where, Format, Allowed diagram types,
and Scan test rules live in
[`../synthesize-report/references/output-contract.md`](../synthesize-report/references/output-contract.md)
under **Diagrams**. Apply them verbatim. If those rules change, this skill
follows — no need to edit here.

**Backfill-mode delta — Grounding:** every entity, arrow, or layer MUST
trace to a finding **already cited in this section's existing prose**.
Resolve the section's grounding set by walking each `[n]` marker in the
section, looking up `cited_sources[i]` where `cited_sources[i].id == n`,
and joining `cited_sources[i].findings_supported[]` against `findings[*].id`.
That set is the only material the diagram may reference. (This differs
from `synthesize-report`'s grounding rule, which uses "findings selected
for the section" because `synthesize-report` controls finding selection
in its Step 3; backfill has no selection step, only an audit of what the
audited prose already cites.) Procedure Step 3 below operationalises this.

Citations therefore live in the surrounding prose, NOT inside the fenced
block; no `[n]` markers inside any fence.

### `augment_notes` — markdown string

A short markdown document recording the decision for each candidate body
section. Suggested shape:

```
## augment-diagrams notes — <run_id>

depth: <quick|standard|deep>
audience: <string>
diagrams_emitted: <integer 0–3>

### Section: <heading>
- decision: <emit | skip>
- diagram_type: <tree | flow | stack | fan | matrix | n/a>
- grounding_finding_ids: [<id>, …]
- scan_test: <passed | failed — reason>
- skip_reason: <only if decision == skip>

…
```

The notes are written to `06-augment_notes.md` alongside the run by the
driver; this skill returns them as a string.

## Procedure

1. **Resolve depth.** Read `research_plan.depth`. Default to `standard` if
   absent. If `quick`, short-circuit: emit `augmented_report` equal to
   `final_report` and `augment_notes` containing the single line
   `skipped — quick depth`. Do not run any subsequent steps.

2. **Parse sections.** Split `final_report` by H2 headings. Bucket each
   section as **candidate** (themed body section, or `## Synthesis` on
   `deep` depth) or **ineligible** (`## Executive Summary`,
   `## Background`, `## Methodology`, `## Key Findings`,
   `## Limitations & Open Questions`, `## Sources`). Preserve the exact
   byte offsets of section starts and the first prose paragraph end —
   the insertion point is the newline immediately after the first prose
   paragraph of the chosen section.

3. **Resolve each candidate's grounding set.** Walk the section's prose
   and collect every `[n]` marker that appears in it. For each marker,
   look up `cited_sources[i]` where `cited_sources[i].id == n`, then
   join its `findings_supported[]` against `findings[*].id` to produce
   the set of findings the section is already grounded in. That set is
   the **only** material the diagram may reference. If a section has
   zero `[n]` markers, mark it `skip — no grounded findings` and move
   on.

4. **Choose a diagram type per candidate (at most one).** Inspect the
   section's prose and grounding set. Decide whether one of the allowed
   types (tree / flow / stack / fan / matrix) clarifies structure that
   is **already present** in the prose. Apply the scan test: "would a
   reader who scans only the diagram still get the section's
   load-bearing structure?" If no, mark `skip — failed scan test`.

5. **Rank and select winners.** Rank surviving candidates by signal
   strength (higher grounding-set size + clearer structural shape =
   stronger). Select at most 3. The "cut don't stretch" rule binds: if
   only 1 candidate clearly passes the scan test, emit 1.

6. **Compose each diagram against its grounding set only.** Use entities
   that name findings or finding-clusters the section's prose already
   cites. Respect the format constraints (≤ 15 lines, ≤ 70 cols, allowed
   character set, no language tag, no nested fences, no `[n]` markers).

7. **Build `augmented_report` by line-walking the input.** For each
   winning section, insert at the position recorded in Step 2 (after the
   first prose paragraph): one blank line, the fenced block, one blank
   line. Every other byte of `final_report` MUST appear in
   `augmented_report` in the same order. The `## Sources` section is
   byte-identical.

8. **Self-check before emit.** Run the Validation gate below. Adjust if
   any gate fails, then emit `augmented_report` and `augment_notes`.

## Constraints

The hard "DO NOT" rules (insertion-only diff, no `[n]` in fences, no
boilerplate-section diagrams, no padding to 3, no novel entities, no
`## Sources` edits, no tabs/nested-fence/language-tag, no quick-depth
diagrams) are the same set the anti-pattern sweep enforces. Canonical list:
[`references/anti-patterns.md`](references/anti-patterns.md).

## Validation gate

See [`references/validation-gate.md`](references/validation-gate.md) for
the 10-check pre-emit gate (insertion-only diff, depth/count caps, fence
validity, character set, line/column caps, placement, grounding, byte-
identical Sources section, untouched title). Any failure → fix and re-run
the relevant Procedure step.

## Anti-patterns

Sweep [`references/anti-patterns.md`](references/anti-patterns.md) before
emitting — 10 failure modes covering prose touch-up, cited diagrams,
boilerplate-section diagrams, padding to 3, novel-entity diagrams,
Sources-list edits, decorated headings, language-tagged fences, quick-
depth overrides, and section-reorder-via-insertion. If any applies, fix
the issue rather than relying on review-report to catch it (backfill is
by definition past the stage-5 gate).

## Worked example

See `examples/capability-cost-router-stack.md` for a complete worked
input → output illustration on a `deep`-depth report at `expert` audience.
It shows: one body-section paragraph from the input, the stack/layer
diagram inserted immediately after it, a one-line note on why the scan
test passes, and a snippet of `augment_notes` recording the decision.

## Troubleshooting

See [`references/troubleshooting.md`](references/troubleshooting.md) for
the edge-case lookup table — quick depth, missing depth, ungrounded
sections, all-fail scan tests, duplicate-structure sections, bullet-list
first paragraphs, body-less reports, and oversized diagrams. Default for
any unlisted case: emit `augmented_report` equal to `final_report` (zero
diagrams) and log the reason in `augment_notes`.

## Notes for downstream stages

- This skill is **out of band**: it is not part of the standard 1 → 6
  pipeline. New runs get diagrams from `synthesize-report` directly.
- The driver typically writes `augment_notes` to
  `<run_dir>/06-augment_notes.md` and overwrites `<run_dir>/05-final_report.md`
  in place with `augmented_report`. Existence of `06-augment_notes.md`
  is the idempotency marker — re-running the driver over the same run
  directory should be a no-op.
- The driver is responsible for the pre-flight `~/.Trash` snapshot of
  `tmp/runs/` and `ResearchVault/` and for the post-run sync of the
  augmented `05-final_report.md` to `ResearchVault/reports/<slug>.md`.
  Neither concern is the skill's.
