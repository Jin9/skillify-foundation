---
name: synthesize-report
description: >
  Compose a markdown research report from structured findings, tuned to a named
  audience, with inline numeric citations and a sources list. Stage 4 of 6 in
  the squad-researcher workflow: receives `topic`, `research_plan`, `findings`,
  and `audience`; emits `draft_report` (markdown) and `cited_sources` (the
  subset of source records actually cited in the prose). Use when the user
  says "draft the report", "synthesize this research", "write the markdown
  report", "compose findings into a report", or as the synthesize-report
  stage of `workflows/researcher.yaml`. Do NOT use to plan research
  (plan-research), fetch sources (search-sources), extract claims
  (extract-findings), or critique the draft (review-report). Do NOT invent
  findings or sources not present in the input.
compatibility: claude-code, codex, copilot, gemini, antigravity
allowed-tools: ""  # no external tools — composes prose from in-context findings
version: 0.2.2
stage: synthesize-report
workflow: researcher
inputs:
  - { name: topic,         type: string, required: true,  source: workflow.inputs.topic,              description: "the subject" }
  - { name: research_plan, type: object, required: true,  source: stages.plan-research.research_plan, description: "from plan-research" }
  - { name: findings,      type: array,  required: true,  source: stages.extract-findings.findings,   description: "from extract-findings" }
  - { name: audience,      type: string, required: false, source: workflow.inputs.audience,           description: "who the report is for" }
outputs:
  - { name: draft_report,  type: string, description: "markdown body with inline [n] citations and a Sources section" }
  - { name: cited_sources, type: array,  description: "ordered list of source records actually cited; subset of input sources" }
---

# Skill: synthesize-report

## Purpose

Turn structured `findings` into a coherent markdown research report on `topic`,
written for `audience`, structured along the research plan's sub-questions.
The skill is the **author**: it picks which findings to elevate, organizes
them into sections, writes prose, places inline citations, and emits the
canonical list of sources actually cited so the downstream `review-report`
stage can audit prose against citations.

**Atomicity:** one stage, one LLM call. The skill is responsible for the
final markdown body and the cited-sources list — nothing else. It does not
search, plan, extract, or self-review.

## Inputs

The `inputs:` frontmatter block is the source of truth for shape and source
of each input (`topic`, `research_plan`, `findings`, `audience`). Notes below
cover the non-obvious wiring.

**Implicit input via `research_plan`:** the workflow does NOT pass `depth`
directly to this stage — depth lives on `research_plan.depth` (one of
`quick` / `standard` / `deep`). Read it from there. If missing, default to
`standard`.

**Source records:** `findings[i].source_id` references the source list
produced by `search-sources`. The full source list is reachable via the
`findings` array (every distinct `source_id` referenced in `findings` is a
candidate source). Do NOT cite a source that no `finding` references.

## Outputs

Both outputs are required.

### `draft_report` — markdown string

Sectioned markdown body. The section set depends on `research_plan.depth`
(`quick` 3 sections / `standard` 5–6 / `deep` 7+) and a body section may carry
an inline ASCII diagram (0–3 per report, `standard`/`deep` only). The canonical
per-depth section set and the full diagram catalog (When / Count / Where /
Format / allowed types / Grounding / Scan test) live in
[`references/output-contract.md`](references/output-contract.md). That file is
the single source of truth; `augment-diagrams` also grounds its backfill rules
against its Diagrams section.

### `cited_sources` — array

Ordered list of source records actually referenced in `draft_report`. Each
record has at least:

```yaml
- id: <integer, matches the [n] marker in prose>
  title: <string, human-readable>
  url: <string, verifiable>
  findings_supported: [<finding_id>, ...]   # which input findings this source backs
```

Additional fields from the upstream source records (author, accessed date,
publisher) carry through verbatim. The list is a **subset** of the input
sources (dereferenced via `findings`) — never invent a source not already
linked to a finding.

`cited_sources[k].id` MUST match every `[id]` marker that appears in
`draft_report`. The mapping is bijective: every `[n]` in prose ↔ exactly one
entry in `cited_sources`, and vice versa. The `## Sources` section in
`draft_report` is a markdown rendering of this same list and MUST agree.

## Procedure

1. **Resolve depth.** Read `research_plan.depth`. If absent, default to
   `standard`. Pick the section set from
   [`references/output-contract.md`](references/output-contract.md).

2. **Resolve audience.** Map `audience` to the closest profile per
   [`references/audience-profiles.md`](references/audience-profiles.md). If
   `audience` is a free-form string that doesn't match cleanly, map to
   `general` and proceed. State the chosen profile internally (do not
   surface in the report).

3. **Select findings.** Compute a score per finding and select up to the
   depth cap (quick ≤ 6, standard ≤ 15, deep ≤ 30):

   - **+3** if the finding answers a `research_plan.sub_questions` entry
     directly.
   - **Confidence bias** from the qualitative `findings[i].confidence`
     bucket emitted by `extract-findings`: `high` → +1.0, `medium` → +0.5,
     `low` → 0. (Never expect a numeric `0–1` here — upstream bans it.)
   - **+1** if backed by ≥2 distinct `source_id`s (i.e. `supporting_source_ids`
     non-empty).
   - **+1** if the finding contradicts the apparent default narrative
     (signal value); a non-empty `disputed_by` is a strong hint.
   - **+1** if the finding is high-relevance to the chosen audience
     profile (impact/risk for executive, mechanism for expert,
     "what it means for me" for general).

   Drop unselected findings from the prose (they stay in the input;
   `review-report` may recall them). If two findings collide (same claim,
   different sources), keep the higher-scored one and merge their
   `source_id`s.

4. **Outline the report.** Before drafting any prose, produce a private
   outline mapping each section to the findings it will cite. Verify
   every selected finding lands in at least one section; verify every
   sub-question gets at least one section (or an explicit "no evidence
   found" sentence in Limitations).

5. **Assign citation ids.** Walk the outline section-by-section. The
   first time a source appears, assign it the next integer `id` starting
   at 1. Reuse the same `id` for subsequent references to that source.
   Build `cited_sources` as you go.

6. **Draft section-by-section.** Write each section against its assigned
   findings only. Apply audience adaptation per
   [`references/audience-profiles.md`](references/audience-profiles.md).
   Place `[n]` cites adjacent to the claim, not at paragraph end. Length
   per section per the depth target — if a section runs short on findings,
   prefer cutting it (or merging with another) over padding.

   **Diagram decision (standard + deep only).** After the first prose
   paragraph of a body section, decide whether one of the diagram types
   in [`references/output-contract.md`](references/output-contract.md)
   → Diagrams clarifies the section. Apply the scan test —
   "would a reader who reads only the diagram still get the load-bearing
   structure?" If yes, emit it after that first paragraph. If no, skip.
   Cap the report at 3 diagrams total; on `quick` depth, always skip.

7. **Render the Sources section.** Append `## Sources` to the report.
   Render `cited_sources` as a numbered list. Format:
   `1. <title> — <url>` (one line per source; add author and date inline
   when present, separated by ` · `).

8. **Self-check before emit.** Run the Validation gate below. Adjust if
   any gate fails.

9. **Emit.** Return `draft_report` (markdown string) and `cited_sources`
   (array).

## Audience adaptation

Four profiles — `general`, `executive`, `expert`, `academic` — each setting
tone, vocabulary, depth bias, and examples/framing. Full profile table,
free-form-string mapping rules, and the validation-gate tone checks live in
[`references/audience-profiles.md`](references/audience-profiles.md).
Audience does NOT change which findings are selected (Step 3 already weighs
relevance) — it changes how they're framed when drafted. When `audience` is
absent or unrecognizable, default to `general` (informed adult reader, not
absolute novice).

## Citation style — inline numeric

- Marker: bracketed integer, no superscript, no space before the bracket:
  `… increases throughput by 23% [4].`
- Multiple sources for one claim: `… [2][5]` (no comma, two brackets).
- Range never: write `[3][4][5]`, not `[3-5]` — review-report's audit is
  per-id.
- Cite **at the claim**, not at paragraph end. The fluency-vs-precision
  trade-off in long-form LLM writing favors precision: citation recall
  drops sharply when cites cluster at the end.
- The `## Sources` section is the canonical render of `cited_sources` and
  uses the same integer ids.

## Validation gate

See [`references/validation-gate.md`](references/validation-gate.md) for the
17-check pre-emit gate (title verbatim, depth/section match, the bijective
`[n]` ↔ `cited_sources[].id` checks, Sources agreement, finding/sub-question
coverage, audience tone, length band, no invented sources, the diagram
count/fence/placement/grounding checks, and every-sentence-traces-to-a-finding).
Any failure → fix and re-run the relevant step; do not emit a known-broken
draft.

## Anti-patterns

Sweep [`references/anti-patterns.md`](references/anti-patterns.md) before
emitting `draft_report` — 11 failure modes covering hallucinated synthesis,
ghost citations, citation drift, audience default-drift, padding,
contradiction laundering, topic restatement, coverage decay, source-list
disagreement, decorative diagrams, and diagrams-as-new-claims. If any
applies to a section, fix the section before emitting.

## Worked example

See `examples/deep-research-systems.md` for a complete worked input → output
illustration at `standard` depth, `general` audience. It shows: depth-appropriate
section set, citations adjacent to claims, an inline v0.2.0 ASCII diagram after
the introducing paragraph of a body section, and the `cited_sources` ↔
`## Sources` bijection.

## Troubleshooting

See [`references/troubleshooting.md`](references/troubleshooting.md) for the
edge-case lookup table — empty findings, missing depth, unrecognized
audience, contradicting findings, short sections, uncovered sub-questions,
upstream-source defects, and length-overruns. Default for any unlisted
case: refuse to invent and surface the gap in `## Limitations & Open
Questions`.

## Notes for downstream stages

- **review-report** consumes `draft_report` and `cited_sources` together. It
  may re-order, drop, or add citations — and on doing so re-emits its own
  `cited_sources`. Keep the bijection `[n] ↔ cited_sources[].id` intact in
  this stage so review-report has a clean baseline to diff against.
- **Workflow output binding:** `cited_sources` from this stage feeds
  review-report; the workflow's final `sources` output comes from
  review-report's `cited_sources`, not this stage's.
