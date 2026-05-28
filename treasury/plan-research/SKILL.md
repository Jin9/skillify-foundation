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
remaining five stages of the researcher workflow can execute against
without re-asking the user anything. The plan locks in scope, names the
sub-questions to investigate, declares what is intentionally out of scope,
and states the success criteria the final report will be judged against.

This is the **planning** stage. It produces a plan; it does not run
searches, read sources, or draft prose.

## When to use this skill

- Stage 1 dispatch from `workflows/researcher.yaml`.
- User prompts: "plan research on X", "decompose this topic", "draft a
  research outline for X", "give me sub-questions for researching X".
- Any caller that has a `topic` plus optional `depth` and `audience` and
  needs a downstream-consumable plan.

Do NOT use this skill to:
- Run web searches or fetch sources — that is `search-sources` (stage 2).
- Extract claims from sources — that is `extract-findings` (stage 3).
- Write or revise the report — `synthesize-report` / `review-report`.
- Plan non-research workflows (code refactors, multi-step ops, etc.).
- Loop with the user for clarifications — atomic stages cannot pause.

## Inputs

Bound from the workflow per `workflows/researcher.yaml`:

| Name       | Type   | Required | Default     | Notes |
|------------|--------|----------|-------------|-------|
| `topic`    | string | yes      | —           | The subject to research. May be a question, noun phrase, or short brief. |
| `depth`    | string | no       | `standard`  | One of `quick` \| `standard` \| `deep`. Governs sub-question count and granularity. |
| `audience` | string | no       | `general`   | Free-text descriptor of the reader (e.g., `general`, `technical`, `executives`, `policy`, `researchers`). Governs vocabulary and emphasis. |

If `topic` is empty or whitespace only, emit a `plan_skipped` failure shape
(see Output contract) — do not invent a topic.

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
  "plan_skipped": false               // true only on the failure shape below
}
```

Failure shape (return this when `topic` is empty/whitespace, or when the
topic is unintelligible after one parsing attempt):

```jsonc
{
  "topic": "<echo or empty>",
  "depth": "<echo>",
  "audience": "<echo>",
  "thesis_question": "",
  "coverage_dimensions": [],
  "sub_questions": [],
  "out_of_scope": [],
  "success_criteria": [],
  "notes_for_downstream": "plan_skipped: <reason>",
  "plan_skipped": true
}
```

Downstream stages (`search-sources`, `extract-findings`,
`synthesize-report`, `review-report`) treat `plan_skipped: true` as a
hard stop; the workflow returns an empty `report_markdown`.

## Procedure

Execute all of the following in a single LLM call. Do not split into
multiple calls and do not ask the user clarifying questions.

1. **Normalize the topic.** Strip leading/trailing whitespace. If the
   topic is empty after normalization, emit the failure shape with reason
   `empty_topic` and stop. If the topic is ambiguous (e.g., a single
   acronym with multiple meanings), pick the most general / most likely
   interpretation, record the picked interpretation in
   `notes_for_downstream` (`assumed_interpretation: ...`), and add the
   alternative readings to `out_of_scope`.

2. **Set the dial values from `depth`.** Apply the Depth dial table
   below. The count is a hard target; the granularity rule shapes how
   each sub-question is phrased.

3. **Draft the thesis question.** Convert `topic` into one declarative
   question whose answer is exactly the report. Examples:
   - topic `"vector databases for RAG"` → thesis `"Which vector databases
     are best suited to retrieval-augmented generation workloads, and on
     what dimensions do they differ?"`
   - topic `"GLP-1 drugs and cardiovascular risk"` → thesis `"What does
     current evidence say about cardiovascular outcomes for patients on
     GLP-1 agonists?"`
   The thesis question must be answerable; reject thesis phrasings
   without a verb or without a clear answer-shape.

4. **Choose 3–6 coverage dimensions (MECE pillars).** These are the
   conceptual axes the report will be organized along. Standard pillars
   to consider — pick those that fit the topic; do not force-fit:
   - `definitions_and_scope` — what the topic *is* and isn't.
   - `mechanisms_or_how_it_works` — internal structure / dynamics.
   - `landscape_or_options` — comparators, players, variants.
   - `evidence_and_outcomes` — empirical results, benchmarks, studies.
   - `trade_offs_or_constraints` — costs, risks, limitations.
   - `temporal_or_trajectory` — history, recent shifts, trends.
   - `counterpoints_or_controversies` — disagreements, edge cases.
   - `practical_implications` — what the audience should do with this.
   The pillars must be **mutually exclusive** (no pair is substitutable)
   and **collectively exhaustive** for the topic at the chosen depth.

5. **Generate sub-questions, one per pillar (mostly).** Each sub-question
   must:
   - Be a **complete natural-language question** ending in `?` — not a
     keyword string. Bad: `"vector DB latency"`. Good: `"How do leading
     vector databases compare on ingestion and query latency for ~10M
     documents?"`
   - Be **single-search-tractable** — a competent searcher should be
     able to answer it from ~1–3 sources at the chosen depth.
   - Be **non-overlapping** with every other sub-question. If two
     candidates overlap, merge them or move one to `out_of_scope`.
   - Be **audience-conditioned** — see the Audience dial.
   - Carry a one-sentence `rationale` and a list of 1–3
     `expected_source_types` (`academic`, `industry`, `news`, `docs`,
     `primary`, `gov`, `community`, `vendor`).
   At `deep` depth, two sub-questions may share a pillar if the second
   drills into an explicit edge case or contrasting view; the second
   must reference the first in its rationale.

6. **Write the `out_of_scope` list.** At least one entry. Explicitly
   exclude:
   - Adjacent topics the user did NOT ask about (e.g., for "vector
     databases for RAG", exclude "training embedding models from
     scratch").
   - Alternative interpretations from step 1 (if any).
   - Time/region/language scopes the topic implicitly excludes.
   This list is load-bearing — downstream `search-sources` uses it to
   suppress drifting queries.

7. **Write 3–6 success criteria.** Each is a one-line testable rubric
   item the final report must satisfy. Mix coverage criteria
   ("Every coverage_dimension is addressed with at least one cited
   source") with quality criteria ("Trade-offs section names at least
   two trade-offs and a concrete example each") and audience criteria
   ("Vocabulary is appropriate for a non-specialist reader; technical
   terms are defined inline on first use" — when audience is `general`).

8. **Write `notes_for_downstream`.** One short paragraph (≤ 80 words)
   handing off planner intuition to the search / synthesis stages:
   any topic-specific search tips, must-include sources or authors,
   known controversies that should not be flattened, and time-sensitivity
   warnings (e.g., "this field changed materially in the last 6 months").

9. **Self-check before emitting.** Apply every item in the "Validation
   checklist" section below. Fix in place; do not make a second LLM call.

10. **Emit the JSON object.** Output the `research_plan` object only;
    no prose preamble, no fences in workflow context (the workflow
    parses the JSON).

## Depth dial

| `depth`    | Sub-question count | Granularity rule                                              | Coverage style                                                                 |
|------------|--------------------|---------------------------------------------------------------|--------------------------------------------------------------------------------|
| `quick`    | 3                  | Each Q answerable from 1 strong source.                       | Pick the 3 most load-bearing pillars only. Drop temporal and controversies.    |
| `standard` | 4–6                | Each Q answerable from 2–3 sources.                           | One pillar per Q; include trade-offs and at least one counterpoint pillar.     |
| `deep`     | 7–10               | Sub-Qs may drill into mechanisms, edge cases, temporal shifts. | Include counterpoints, temporal, and practical-implications pillars by default. |

If `depth` is missing or unrecognized, treat it as `standard`.

## Audience dial

Adjust sub-question phrasing and pillar mix, not count:

- `general` (default) — include a `definitions_and_scope` pillar; prefer
  plain language; one sub-question should ask about real-world
  implications. Avoid acronyms without inline expansion in the
  sub-question text.
- `technical` / `researchers` / `engineers` — drop the definitions pillar
  if the topic is mainstream in the field; allow jargon; emphasize
  `mechanisms`, `evidence`, `trade_offs`. Add benchmarks/metrics to
  `expected_source_types` where relevant.
- `executives` / `decision_makers` / `policy` — emphasize
  `trade_offs`, `practical_implications`, `counterpoints`. At least one
  sub-question should ask "what should the audience do with this?" or
  "what is the decision rule?"
- For any custom audience string, infer the closest of the above and
  record the inference in `notes_for_downstream`.

## Failure modes and how to avoid them

| Failure                                          | Symptom                                                              | Fix during step                                                                   |
|--------------------------------------------------|----------------------------------------------------------------------|-----------------------------------------------------------------------------------|
| Sub-question too broad                           | Answer would require an entire literature review.                    | Step 5: narrow to one mechanism / one region / one timeframe / one comparator.    |
| Sub-question too granular                        | Answer is a single fact already in the topic statement.              | Step 5: merge into a sibling sub-question or drop.                                |
| Redundant / overlapping sub-questions            | Two Qs share answer-space; second adds no new evidence.              | Step 5: merge, or move the weaker to `out_of_scope`.                              |
| Sub-questions phrased as keywords                | No `?`; reads like a search query.                                   | Step 5: rewrite as a complete question. Reject keyword strings in the self-check. |
| Missing `out_of_scope`                           | Downstream search drifts into adjacent topics.                       | Step 6: must produce ≥ 1 exclusion; infer from topic edges if none obvious.       |
| No `success_criteria`                            | `review-report` has nothing to check against.                        | Step 7: always produce 3–6.                                                        |
| Audience-blind plan                              | Plan reads the same regardless of `audience` input.                  | Step 5 + Audience dial: shift pillar mix and vocabulary.                          |
| Depth-blind plan                                 | Plan reads the same regardless of `depth` input.                     | Step 2 + Depth dial: enforce hard sub-question count.                             |
| Thesis is a noun phrase, not a question          | `thesis_question` lacks a verb or `?`.                               | Step 3: convert to declarative question.                                          |
| Hidden clarifying question                       | Plan contains "TBD pending user input" or similar.                   | Step 1: assume an interpretation; record it; do not defer.                        |

## Worked example

See `examples/vector-db-rag.md` for a complete worked input → output illustration.
It shows: a `standard`/`technical` plan with 5 MECE sub-questions, jargon kept
audience-appropriate, an explicit `out_of_scope` list, and 5 testable success
criteria the downstream `review-report` stage can audit against.

## Constraints

- DO NOT make a second LLM call. Atomic stage = one call.
- DO NOT ask the user clarifying questions. Assume; record the
  assumption in `notes_for_downstream`.
- DO NOT emit prose around the JSON in workflow context.
- DO NOT exceed the Depth dial's sub-question count.
- DO NOT phrase sub-questions as keyword strings — always full questions.
- DO NOT leave `out_of_scope` or `success_criteria` empty.
- DO NOT modify `workflows/researcher.yaml` or files outside this skill
  folder.
- DO NOT invent sources, URLs, or citations — that is `search-sources`'s
  job. The plan has no citations.

## Validation checklist

Before returning, verify all of:

- [ ] `topic`, `depth`, `audience` echoed verbatim from input.
- [ ] `thesis_question` is one well-formed question ending in `?`.
- [ ] `coverage_dimensions` length is 3–6 and matches the sub-questions' pillars.
- [ ] `sub_questions` length matches Depth dial (3 / 4–6 / 7–10).
- [ ] Every `sub_question.question` ends in `?` and is a complete sentence.
- [ ] Every `sub_question.coverage_dimension` appears in `coverage_dimensions`.
- [ ] No two `sub_question.question` values are substitutable.
- [ ] `out_of_scope` has ≥ 1 entry.
- [ ] `success_criteria` has 3–6 entries, all testable.
- [ ] `notes_for_downstream` is ≤ 80 words.
- [ ] `plan_skipped` is `false` unless the failure shape was emitted.

## Downstream contract

| Stage                | Consumes from `research_plan`                                                                |
|----------------------|----------------------------------------------------------------------------------------------|
| `search-sources`     | `sub_questions[].question`, `expected_source_types`, `out_of_scope`, `notes_for_downstream`. |
| `extract-findings`   | `sub_questions[].id` (grouping key), `sub_questions[].question`, `thesis_question`.          |
| `synthesize-report`  | `thesis_question`, `coverage_dimensions`, `sub_questions[]`, `audience`.                     |
| `review-report`      | `thesis_question`, `success_criteria`, `out_of_scope`, `audience`.                           |

If a downstream stage breaks on a field this skill produced, fix it
here (and bump the version), not there. The plan is the contract.

## References

This skill is self-contained; no sibling reference files are required at
runtime. Background prior-art (GPT Researcher, OpenAI Deep Research,
Perplexity Deep Research, STORM, Elicit, Static-DRA, MAST failure
taxonomy, audience-writing literature) informed the Depth dial,
coverage-dimension list, and self-check rules above but is not loaded
into the model context.
