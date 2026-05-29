# Output contract — section sets and diagram catalog

Canonical home for `synthesize-report`'s `draft_report` shape: the per-depth
section set, and the full inline-diagram catalog. `SKILL.md` → Outputs points
here; this file is the single source of truth. `augment-diagrams` also treats
the Diagrams section below as canonical (it backfills against these rules).

Both outputs (`draft_report`, `cited_sources`) are required. The
`cited_sources` shape and the `[n] ↔ cited_sources[].id` bijection are defined
in `SKILL.md` → Outputs → `cited_sources`.

## `draft_report` section set by depth

The section set depends on `research_plan.depth`.

### `quick` (≈300–600 words, 3 sections)

1. `# <topic>` — title is the input `topic` verbatim.
2. `## TL;DR` — 2–4 sentences answering the topic question, with inline `[n]` cites.
3. `## Findings` — bullet list, each bullet 1–2 sentences with at least one `[n]` cite.
4. `## Sources` — numbered list.

### `standard` (≈800–1500 words, 5–6 sections)

1. `# <topic>`
2. `## Executive Summary` — 3–5 sentences. Answers the topic question first; busy readers stop here.
3. `## Background` — 1 short paragraph framing why the topic matters and what was asked.
4. `## Key Findings` — bullet list, each bullet 1–2 sentences with `[n]` cites.
5. `## <one section per research-plan sub-question>` — 2–5 themed body sections; titles taken from `research_plan.sub_questions`. 2–4 paragraphs each with inline `[n]` cites and optional bullet sub-points.
6. `## Limitations & Open Questions` — what wasn't answered, where confidence is low, where sources disagree.
7. `## Sources` — numbered list.

### `deep` (≈2000–4000 words, 7+ sections)

Same as `standard` plus:
- A `## Methodology` section (1 paragraph) between Background and Key Findings, naming the research plan's angle and any scope limits.
- 3–5 themed body sections instead of 2–3, with deeper exposition (contrasts, mechanisms, comparative claims, parameter values where present in findings).
- Optionally a `## Synthesis` section after themed bodies for cross-cutting claims that combine multiple findings — clearly labeled as synthesis, not as a primary finding.

## Diagrams (canonical rules; in v0.2.0)

> **Cross-skill dependency:** `augment-diagrams` (the out-of-band backfill
> skill) treats the rules below as canonical and references this file
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
