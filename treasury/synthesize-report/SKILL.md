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

| Input | Shape | Source |
|---|---|---|
| `topic` | string | workflow input, verbatim |
| `research_plan` | object with at least `{ sub_questions[], depth }` plus optional `{ angle, scope_notes, audience_notes }` from `plan-research` | upstream stage `plan-research` |
| `findings` | array of objects with at least `{ id, claim, evidence, source_id, confidence }` from `extract-findings` | upstream stage `extract-findings` |
| `audience` | string — one of `general` / `executive` / `expert` / `academic` (free-form strings tolerated; map to closest profile) | workflow input |

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

Sectioned markdown body. Section set depends on `research_plan.depth`:

**`quick` (≈300–600 words, 3 sections):**

1. `# <topic>` — title is the input `topic` verbatim.
2. `## TL;DR` — 2–4 sentences answering the topic question, with inline `[n]` cites.
3. `## Findings` — bullet list, each bullet 1–2 sentences with at least one `[n]` cite.
4. `## Sources` — numbered list.

**`standard` (≈800–1500 words, 5–6 sections):**

1. `# <topic>`
2. `## Executive Summary` — 3–5 sentences. Answers the topic question first; busy readers stop here.
3. `## Background` — 1 short paragraph framing why the topic matters and what was asked.
4. `## Key Findings` — bullet list, each bullet 1–2 sentences with `[n]` cites.
5. `## <one section per research-plan sub-question>` — 2–5 themed body sections; titles taken from `research_plan.sub_questions`. 2–4 paragraphs each with inline `[n]` cites and optional bullet sub-points.
6. `## Limitations & Open Questions` — what wasn't answered, where confidence is low, where sources disagree.
7. `## Sources` — numbered list.

**`deep` (≈2000–4000 words, 7+ sections):**

Same as `standard` plus:
- A `## Methodology` section (1 paragraph) between Background and Key Findings, naming the research plan's angle and any scope limits.
- 3–5 themed body sections instead of 2–3, with deeper exposition (contrasts, mechanisms, comparative claims, parameter values where present in findings).
- Optionally a `## Synthesis` section after themed bodies for cross-cutting claims that combine multiple findings — clearly labeled as synthesis, not as a primary finding.

#### Diagrams (new in v0.2.0)

> **Cross-skill dependency:** `augment-diagrams` (the out-of-band backfill
> skill) treats the rules below as canonical and references this section
> directly. When editing these rules, also re-read
> `skills/augment-diagrams/SKILL.md` → Outputs → Diagrams to confirm the
> backfill grounding-rule delta still applies cleanly.

A body section may carry an inline ASCII diagram when structure does more work
than another paragraph. Rules:

- **When:** depth ∈ `standard` or `deep`. On `quick`, do not emit diagrams.
- **Count:** 0–3 across the whole report. Zero is a valid answer. Never pad —
  same "cut don't stretch" rule as prose.
- **Where:** inside a body section (the themed sub-question sections, or a
  deep-depth `## Synthesis`), placed AFTER the first prose paragraph that
  introduces the concept the diagram visualizes. Never in `## Executive
  Summary`, `## Background`, `## Methodology`, `## Key Findings`,
  `## Limitations & Open Questions`, or `## Sources`.
- **Format:** a fenced code block with no language tag. ≤ 15 lines.
  ≤ 70 columns per line. ASCII plus Unicode box-drawing characters only
  (`─ │ ┌ ┐ └ ┘ ├ ┤ ┬ ┴ ┼ → ← ↔`). No tabs. No nested fences.
- **Allowed diagram types** (pick the one the section needs):
  - **Tree** — taxonomies, type hierarchies.
  - **Flow** — `A → B → C` pipelines or state transitions.
  - **Stack/layer** — boxed layers stacked top-to-bottom.
  - **Fan-in / fan-out** — multiple inputs into a fuser or one source into many sinks.
  - **Comparison matrix** — 2–3 column ASCII contrast (only when a markdown
    table can't carry the same visual contrast).
- **Grounding:** every entity, arrow, or layer in a diagram MUST trace to
  findings already selected for the section the diagram sits in. A diagram
  is a structural rendering of claims the prose already makes — not a new
  claim. Citations therefore live in the surrounding prose, NOT inside the
  fenced block. No `[n]` markers inside a fenced block.
- **Scan test:** before emitting a diagram, ask "would a reader who scans
  ONLY the diagram still get the section's load-bearing structure?" If no,
  skip it — that diagram is decoration, not communication.

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
   `standard`. Pick the section set from the Outputs table above.

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
   in Outputs → Diagrams clarifies the section. Apply the scan test —
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

Before emitting, confirm all of the following. Any failure → fix and re-run
the relevant step; do not emit a known-broken draft.

1. **Title is the input `topic` verbatim** in an `# H1`.
2. **Section set matches depth** (per Outputs table).
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
