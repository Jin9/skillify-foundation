---
name: search-sources
description: >
  Stage 2 of the researcher workflow. Given a structured research_plan
  (output of plan-research) and a depth knob, gather candidate sources for
  every sub-question and emit a deduplicated, diversity-checked list of
  normalised source records that extract-findings can consume.
  Backend-agnostic: wraps whatever search tool the host configures (Tavily,
  Exa, Brave, Serper, Bing, Google CSE, embeddings store, etc.) behind a
  fixed output schema. Use when invoked as the search-sources stage of
  workflows/researcher.yaml, or when a caller says "gather sources for this
  research plan", "find sources for these sub-questions", "scout the web for
  this plan". Do NOT use to plan the research (use plan-research), to
  extract claims from sources (use extract-findings), or to synthesize a
  report (use synthesize-report); do NOT edit workflows/researcher.yaml or
  any sibling skill.
compatibility: claude-code, codex, copilot, gemini, antigravity
allowed-tools: WebSearch  # plus any MCP search backend the host configures (Tavily, Exa, Brave, Serper, Bing, Google CSE, etc.); skill is backend-agnostic by design
version: 0.1.1
stage: search-sources
workflow: researcher
inputs:
  - { name: research_plan, type: object, required: true, source: stages.plan-research.research_plan, description: "from plan-research; expects sub_questions[] with stable ids" }
  - { name: depth,         type: string, required: true, source: workflow.inputs.depth,              description: "quick | standard | deep" }
outputs:
  - { name: sources, type: array, description: "list of source records: {source_id, url, title, snippet, retrieved_at, sub_question_ids, ...}" }
---

# search-sources

## Purpose

Turn a `research_plan` into a vetted `sources` list in one LLM call. The job is
**finding and normalising candidate sources** — not reading them in depth, not
extracting claims, not writing prose. Downstream stages assume every emitted
source has a real retrievable URL (no fabrications), is tagged with the
sub-question(s) it addresses, carries a stable `source_id` (so later stages cite
without re-resolving URLs), and that the list is deduplicated and
diversity-checked.

## Workflow contract

Reads from `workflows/researcher.yaml`. Second of six stages
(plan-research → **search-sources** → extract-findings → synthesize-report →
review-report → reporting-research-run).

### Inputs

| Name            | Type   | Notes                                                                                                  |
|-----------------|--------|--------------------------------------------------------------------------------------------------------|
| `research_plan` | object | Output of `plan-research`. Must expose an iterable `sub_questions` collection (each with an id + text). |
| `depth`         | string | One of `quick`, `standard`, `deep`. Controls queries-per-sub-question and target source count.          |

The skill assumes no shape from `research_plan` beyond `sub_questions`;
pass-through fields (topic, audience, framing) are allowed and may bias query
phrasing but must not be required.

### Outputs

| Name      | Type  | Notes                                                                                  |
|-----------|-------|----------------------------------------------------------------------------------------|
| `sources` | array | List of source records (see **Source records** below). Order: most-relevant first within each sub-question, then sub-question order from the plan. |

## Source records

Each `sources[i]` is an object whose **required** fields are `source_id`, `url`,
`title`, `snippet`, `retrieved_at`, `query`, `sub_question_ids`, and
`source_type`; optional fields (`published_date`, `author`, `publisher`,
`relevance_score`, `notes`) are emitted only when derivable, never fabricated.
Full field types, rules, and the optional `coverage` side channel are in
`references/source-record-schema.md`.

## Depth → query budget

Shapes how many queries are issued per sub-question and how many sources are
retained. Treat the numbers as soft targets; the diversity and dedup gates always
trump quotas.

| `depth`    | queries per sub-question | total source target | per-source-type variants                            |
|------------|--------------------------|---------------------|------------------------------------------------------|
| `quick`    | 1                        | 5–8                 | One literal/keyword query.                           |
| `standard` | 2                        | 10–20               | One keyword + one semantic/intent.                   |
| `deep`     | 3–4                      | 25–40               | Keyword + semantic + recency-filtered + contrarian.  |

`deep` may iterate: after the first pass, if any sub-question is below half the
per-question target or coverage is flagged, formulate one extra query naming the
gap.

## Procedure (the LLM call)

This is the single LLM call the stage makes. Execute the steps in order.

1. **Parse the plan.** Extract `sub_questions`. If the plan has zero
   sub-questions, emit `sources: []` with a `coverage` array noting the empty
   plan, and stop. Do not invent sub-questions.
2. **Formulate queries.** For each sub-question, write queries per the Depth →
   query budget table. Rules:
   - Each query must be self-contained — no pronouns like "this"/"that" that
     depend on the plan text.
   - **Keyword variant** uses the sub-question's core entities + a constraining
     term (year, jurisdiction, method). Good for proper nouns, spec lookups,
     statutes, exact quotes.
   - **Semantic/intent variant** rephrases the question for an embeddings backend.
     Good for conceptual coverage.
   - **Recency-filtered variant** (deep only) adds a recency clause (e.g.
     `after:2025-01-01`) when the topic is time-sensitive.
   - **Contrarian variant** (deep only) inverts the implicit hypothesis (e.g.
     "evidence against X", "limitations of Y"). Counters confirmation bias.
   - Avoid leading phrasing ("why X is good"); prefer neutral nouns.
   - Drop any query that duplicates another after lowercasing and trimming
     stopwords.
3. **Call the configured search backend** for each query. The skill is
   backend-agnostic: the host supplies the search tool (Tavily, Exa, Brave,
   Serper, Bing, Google CSE, a local index, or an MCP search server). The LLM
   emits only sources a tool actually returned. If no search tool is available,
   stop and emit `sources: []` with a per-sub-question `coverage` entry noting
   `gap_note: "no search backend available"`.
4. **Normalise each result** into the Source-record schema
   (`references/source-record-schema.md`). Map the backend's field names to the
   required fields; fill `retrieved_at` with the current UTC timestamp if omitted;
   assign a fresh `source_id` per unique URL (`s1`, `s2`, … in kept order).
5. **Deduplicate.** Apply in order:
   1. Exact URL match (after stripping `utm_*`, `ref`, and trailing slash).
   2. Same `(canonical_domain, normalized_title)` near-match.
   3. Snippet 5-gram Jaccard overlap ≥ 0.7.
   When two records collide, keep the higher `relevance_score` (or the earlier
   `published_date` if scores tie) and merge their `sub_question_ids`.
6. **Apply diversity caps.**
   - No single canonical domain may exceed 30% of the final list (or more than 2
     sources when the list has fewer than 7).
   - If the topic admits primary/academic sources, ensure at least one
     `source_type ∈ {paper, docs, dataset}` per sub-question when the backend
     returned any. If only blog/news exist, leave it and log it in `coverage`.
7. **Score relevance.** Use the backend score when present. Otherwise judge each
   source's `relevance_score` on 0–1 by whether the snippet substantively
   addresses the sub-question (not the topic in general). Drop sources below 0.3.
8. **Build coverage summary.** Per sub-question, count surviving sources and write
   a one-line `gap_note` when coverage is below `floor(target / 2)` or sources are
   all one type or all paywalled.
9. **Validate before emitting** (see Validation gate).
10. **Emit** the `sources` array (and `coverage` if produced). Order: sub-question
    order from the plan, then `relevance_score` descending within each.

## Validation gate

Before emitting, the LLM call must self-check:

1. Every `sources[i].url` was produced by an actual search-tool result in this
   call. No URL was constructed, completed, or guessed.
2. Every required field is present and non-empty on every record.
3. Every `sources[i].sub_question_ids[*]` references a real id from
   `research_plan.sub_questions`. No orphan ids.
4. No duplicate `source_id` and no duplicate canonical URL.
5. No single canonical domain exceeds the diversity cap from step 6.
6. The total source count is within the depth target band (above the floor) OR
   the `coverage` array explains the shortfall.
7. `retrieved_at` is a valid ISO 8601 UTC timestamp.

If a check fails, fix and re-emit. If a fix is impossible (e.g. backend returned
nothing for a sub-question), keep the source list valid and record the gap in
`coverage`; do not pad with low-quality sources to hit the target.

## Constraints

- DO NOT fabricate URLs, titles, snippets, authors, dates, or any other field.
  Only emit what a tool returned or what is directly visible in a fetched page.
- DO NOT read sources in depth in this stage — that is `extract-findings`.
  Snippets are the unit of work here.
- DO NOT make claims about source content beyond restating the snippet.
- DO NOT prefer a single search backend; use whatever tool the host wired in.
- DO NOT invent sub-questions or rewrite the plan. The plan is read-only.
- DO NOT edit `workflows/researcher.yaml` or any sibling skill folder.
- DO NOT emit paywalled-only coverage silently. Surface it via `coverage`.
- DO NOT exceed the diversity cap to pad the list. Empty seats beat homogeneous
  noise.

## Failure modes

Degenerate-search handling (backend down, zero results, paywalls, hallucinated
URLs, etc.): see `references/failure-modes.md`.

## References

| Need | File |
|------|------|
| Full source-record field spec + `coverage` side channel | `references/source-record-schema.md` |
| Symptom → action table for degenerate searches | `references/failure-modes.md` |
| A complete worked input → output | `examples/heat-pump-retrofit.md` |

Prior-art (GPT Researcher, Perplexity / Gemini Deep Research, Tavily / Exa /
academic query-expansion literature) informed the source-record schema, the
depth → query-budget table, and the diversity cap but is not loaded into model
context.
