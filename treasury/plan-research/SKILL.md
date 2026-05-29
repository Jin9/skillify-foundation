---
name: plan-research
description: >
  Stage 1 of the squad-researcher workflow. Decompose a research topic into
  a structured plan — thesis question, 3–10 MECE sub-questions (count gated
  by depth), coverage dimensions, out-of-scope list, success criteria — that
  downstream stages (search-sources, extract-findings, synthesize-report,
  review-report) consume as a single JSON object. Use when invoked as the
  `plan-research` stage in `workflows/researcher.yaml`, or when the user
  asks to "plan research on X", "decompose this topic", "draft a research
  outline", or "give me sub-questions for X". One stage = one LLM call;
  no clarifying-question loops, no recursion. Do NOT use to actually run
  searches (that is search-sources), to write the report (synthesize-report),
  or to plan non-research workflows.
compatibility: claude-code, codex, copilot, gemini, antigravity
allowed-tools: ""  # no external tools — pure text-in / JSON-out LLM call
version: 0.1.1
stage: plan-research
workflow: researcher
inputs:
  - { name: topic,    type: string, required: true,  source: workflow.inputs.topic,    description: "the subject to research" }
  - { name: depth,    type: string, required: false, source: workflow.inputs.depth,    description: "quick | standard | deep; default standard" }
  - { name: audience, type: string, required: false, source: workflow.inputs.audience, description: "who the report is written for; default general" }
outputs:
  - { name: research_plan, type: object, description: "structured plan: thesis question, sub-questions, scope, success criteria" }
---

# Skill: plan-research

## Purpose

Convert a raw research `topic` into a structured `research_plan` that the
remaining five stages of the researcher workflow execute against without
re-asking the user anything. The plan locks scope, names the sub-questions to
investigate, declares what is out of scope, and states the success criteria the
final report is judged against. This is the **planning** stage: it produces a
plan; it does not run searches, read sources, or draft prose.

## When to use this skill

- Stage 1 dispatch from `workflows/researcher.yaml`.
- User prompts: "plan research on X", "decompose this topic", "draft a research
  outline for X", "give me sub-questions for researching X".
- Any caller with a `topic` (plus optional `depth`, `audience`) needing a
  downstream-consumable plan.

Do NOT use this skill to run searches (`search-sources`, stage 2), extract claims
(`extract-findings`, stage 3), write/revise the report (`synthesize-report` /
`review-report`), plan non-research workflows, or loop with the user for
clarifications — atomic stages cannot pause.

## Inputs

Bound from the workflow per `workflows/researcher.yaml`:

| Name       | Type   | Required | Default     | Notes |
|------------|--------|----------|-------------|-------|
| `topic`    | string | yes      | —           | The subject to research. May be a question, noun phrase, or short brief. |
| `depth`    | string | no       | `standard`  | One of `quick` \| `standard` \| `deep`. Governs sub-question count and granularity. |
| `audience` | string | no       | `general`   | Free-text reader descriptor (e.g., `general`, `technical`, `executives`, `policy`). Governs vocabulary and emphasis. |

If `topic` is empty or whitespace only, emit the failure shape (see Output
contract) — do not invent a topic.

## Output contract

Emit a single JSON object named `research_plan` matching this schema:

```jsonc
{
  "topic": "string",                  // echo of the input topic
  "depth": "quick|standard|deep",     // echo
  "audience": "string",               // echo
  "thesis_question": "string",        // the one question the report must answer
  "coverage_dimensions": ["string"],  // 3-6 MECE pillars; each sub_question maps to one
  "sub_questions": [                  // count gated by depth (see Depth dial)
    {
      "id": "q1",                     // q1..qN, contiguous
      "question": "string",           // natural-language, searchable, testable
      "rationale": "string",          // one sentence on why this matters to the topic
      "coverage_dimension": "string", // must appear in coverage_dimensions
      "expected_source_types": ["string"]  // e.g., academic, news, docs, primary, industry, gov
    }
  ],
  "out_of_scope": ["string"],         // explicit exclusions, ≥ 1
  "success_criteria": ["string"],     // rubric items the final report must hit, 3-6
  "notes_for_downstream": "string",   // hints for search-sources / synthesize-report
  "plan_skipped": false               // true only on the failure shape
}
```

**Failure shape.** When `topic` is empty/whitespace, or unintelligible after one
parsing attempt, emit the same object with every content field empty (`""` /
`[]`), `notes_for_downstream` set to `"plan_skipped: <reason>"`, and
`plan_skipped: true`. Downstream stages treat `plan_skipped: true` as a hard stop
and the workflow returns an empty `report_markdown`.

## Procedure

Execute all of the following in a single LLM call. Do not split into multiple
calls and do not ask the user clarifying questions.

1. **Normalize the topic.** Strip whitespace. If empty after normalization, emit
   the failure shape with reason `empty_topic` and stop. If ambiguous (e.g., an
   acronym with multiple meanings), pick the most likely interpretation, record
   it in `notes_for_downstream` (`assumed_interpretation: ...`), and add the
   alternative readings to `out_of_scope`.
2. **Set the dial values from `depth`.** Apply the Depth dial below. The count is
   a hard target; the granularity rule shapes how each sub-question is phrased.
3. **Draft the thesis question.** Convert `topic` into one declarative question
   whose answer is exactly the report (e.g., topic `"vector databases for RAG"` →
   `"Which vector databases are best suited to RAG workloads, and on what
   dimensions do they differ?"`). Reject phrasings without a verb or a clear
   answer-shape.
4. **Choose 3–6 coverage dimensions (MECE pillars).** Pick from the pillar
   catalog in `references/coverage-dimensions.md`; do not force-fit. Pillars must
   be mutually exclusive and collectively exhaustive for the topic at this depth.
5. **Generate sub-questions, one per pillar (mostly).** Each must:
   - Be a complete natural-language question ending in `?` — not a keyword
     string. Bad: `"vector DB latency"`. Good: `"How do leading vector databases
     compare on ingestion and query latency for ~10M documents?"`
   - Be single-search-tractable (answerable from ~1–3 sources at this depth).
   - Be non-overlapping with every other sub-question; merge or move overlaps to
     `out_of_scope`.
   - Be audience-conditioned — see the Audience dial.
   - Carry a one-sentence `rationale` and 1–3 `expected_source_types`
     (`academic`, `industry`, `news`, `docs`, `primary`, `gov`, `community`,
     `vendor`).
   At `deep` depth, two sub-questions may share a pillar if the second drills into
   an explicit edge case or contrasting view and references the first in its
   rationale.
6. **Write `out_of_scope` (≥ 1 entry).** Exclude adjacent topics not asked about,
   alternative interpretations from step 1, and implicit time/region/language
   scopes. Load-bearing: `search-sources` uses it to suppress drifting queries.
7. **Write 3–6 success criteria.** Each a one-line testable rubric item the final
   report must satisfy — mix coverage, quality, and audience criteria.
8. **Write `notes_for_downstream`** (≤ 80 words): topic-specific search tips,
   must-include sources/authors, known controversies not to flatten, and
   time-sensitivity warnings.
9. **Self-check before emitting.** Apply every item in
   `references/validation-checklist.md`. Fix in place; do not make a second call.
10. **Emit the JSON object.** Output the `research_plan` object only — no prose
    preamble, no fences in workflow context (the workflow parses the JSON).

## Depth dial

| `depth`    | Sub-question count | Granularity rule                                              | Coverage style                                                                 |
|------------|--------------------|---------------------------------------------------------------|--------------------------------------------------------------------------------|
| `quick`    | 3                  | Each Q answerable from 1 strong source.                       | Pick the 3 most load-bearing pillars only. Drop temporal and controversies.    |
| `standard` | 4–6                | Each Q answerable from 2–3 sources.                           | One pillar per Q; include trade-offs and at least one counterpoint pillar.     |
| `deep`     | 7–10               | Sub-Qs may drill into mechanisms, edge cases, temporal shifts. | Include counterpoints, temporal, and practical-implications pillars by default. |

If `depth` is missing or unrecognized, treat it as `standard`.

## Audience dial

Adjust sub-question phrasing and pillar mix, not count:

- `general` (default) — include a `definitions_and_scope` pillar; prefer plain
  language; one sub-question asks about real-world implications; expand acronyms
  inline.
- `technical` / `researchers` / `engineers` — drop the definitions pillar if the
  topic is mainstream in the field; allow jargon; emphasize `mechanisms`,
  `evidence`, `trade_offs`; add benchmarks/metrics to `expected_source_types`.
- `executives` / `decision_makers` / `policy` — emphasize `trade_offs`,
  `practical_implications`, `counterpoints`; at least one sub-question asks "what
  should the audience do?" or "what is the decision rule?"
- For any custom audience string, infer the closest of the above and record the
  inference in `notes_for_downstream`.

## Constraints

- DO NOT make a second LLM call. Atomic stage = one call.
- DO NOT ask the user clarifying questions. Assume; record the assumption in
  `notes_for_downstream`.
- DO NOT emit prose around the JSON in workflow context.
- DO NOT exceed the Depth dial's sub-question count.
- DO NOT phrase sub-questions as keyword strings — always full questions.
- DO NOT leave `out_of_scope` or `success_criteria` empty.
- DO NOT modify `workflows/researcher.yaml` or files outside this skill folder.
- DO NOT invent sources, URLs, or citations — that is `search-sources`'s job.

## References

| Need | File |
|------|------|
| The MECE pillar catalog for step 4 | `references/coverage-dimensions.md` |
| Self-check list for step 9 | `references/validation-checklist.md` |
| Symptom → fix table for bad plans | `references/failure-modes.md` |
| Which fields each later stage consumes | `references/downstream-contract.md` |
| A complete worked input → output | `examples/vector-db-rag.md` |

Prior-art (GPT Researcher, OpenAI / Perplexity Deep Research, STORM, Elicit,
Static-DRA, MAST failure taxonomy, audience-writing literature) informed the
dials, pillar catalog, and self-check rules but is not loaded into model context.
