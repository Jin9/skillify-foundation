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

Turn a `research_plan` into a vetted `sources` list. One LLM call per
invocation. The stage's job is **finding and normalising candidate
sources**, not reading them in depth, not extracting claims, not writing
prose. Downstream stages assume:

- every emitted source has a real, retrievable URL (no fabrications),
- every source is tagged with which sub-question(s) it addresses,
- the list is deduplicated and diversity-checked,
- sources have a stable `source_id` so later stages can cite without
  re-resolving URLs.

## Workflow contract

Reads from `workflows/researcher.yaml`. This stage is the second of six
(plan-research → **search-sources** → extract-findings → synthesize-report
→ review-report → reporting-research-run).

### Inputs

| Name            | Type   | Notes                                                                                                  |
|-----------------|--------|--------------------------------------------------------------------------------------------------------|
| `research_plan` | object | Output of `plan-research`. Must expose an iterable `sub_questions` collection (each with an id + text). |
| `depth`         | string | One of `quick`, `standard`, `deep`. Controls queries-per-sub-question and target source count.          |

The skill does not assume any further shape from `research_plan` beyond
`sub_questions`; pass-through fields (topic, audience, framing) are
allowed and may be used to bias query phrasing but must not be required.

### Outputs

| Name      | Type  | Notes                                                                                  |
|-----------|-------|----------------------------------------------------------------------------------------|
| `sources` | array | List of source records using the schema in **Source-record schema** below. Order: most-relevant first within each sub-question, then sub-question order from the plan. |

## Source-record schema

Each element of `sources` is an object with these fields. Required fields
MUST be present and non-empty. Optional fields are emitted when the search
backend supplies them or they can be inferred from the result page; never
fabricated.

| Field              | Required | Type                                | Description                                                                                       |
|--------------------|----------|-------------------------------------|---------------------------------------------------------------------------------------------------|
| `source_id`        | yes      | string                              | Short stable id, kebab/numbered (e.g. `s1`, `s2`, `s12`). Unique within the call.                  |
| `url`              | yes      | string (URL)                        | Canonical URL returned by the search backend. Must be one the tool actually returned.              |
| `title`            | yes      | string                              | Page or document title as returned (trimmed). If backend gave nothing, use the URL's last path segment and mark `title_inferred: true`. |
| `snippet`          | yes      | string                              | 1–3 sentence summary from the backend or the page's meta description. Verbatim or near-verbatim.   |
| `retrieved_at`     | yes      | string (ISO 8601 UTC)               | Timestamp the search backend returned this result. Use the call time if the backend omits it.      |
| `query`            | yes      | string                              | Exact query string that surfaced this source. Used for dedup and debug.                            |
| `sub_question_ids` | yes      | array of string                     | Ids of the plan sub-questions this source addresses. Length ≥ 1.                                   |
| `source_type`      | yes      | enum                                | One of `web`, `paper`, `docs`, `news`, `forum`, `video`, `book`, `dataset`, `other`.               |
| `published_date`   | no       | string (ISO 8601 date) or null      | Publication date if extractable. Omit or `null` otherwise. Never guess.                            |
| `author`           | no       | string or null                      | Author or byline if present.                                                                       |
| `publisher`        | no       | string or null                      | Site / journal / org name.                                                                         |
| `relevance_score`  | no       | float in [0, 1] or null             | Backend-supplied score (Tavily `score`, Exa similarity, etc.) or LLM-judged. Omit if not derivable. |
| `notes`            | no       | string or null                      | One short note from the agent — e.g. "paywalled", "primary source", "opinion piece".               |

A `coverage` array MAY be emitted alongside `sources` on the same output
envelope for the convenience of `review-report`; the workflow contract
only consumes `sources`, so `coverage` is a non-binding side channel.
Shape: `[{sub_question_id, source_count, gap_note}]`.

## Depth → query budget

The `depth` knob shapes both how many queries are issued per sub-question
and how many sources are retained. Treat the numbers as soft targets; the
diversity and dedup gates always trump quotas.

| `depth`    | queries per sub-question | total source target | per-source-type variants                            |
|------------|--------------------------|---------------------|------------------------------------------------------|
| `quick`    | 1                        | 5–8                 | One literal/keyword query.                           |
| `standard` | 2                        | 10–20               | One keyword + one semantic/intent.                   |
| `deep`     | 3–4                      | 25–40               | Keyword + semantic + recency-filtered + contrarian.  |

`deep` may iterate: after the first pass, if any sub-question is below
half the per-question target or coverage is flagged, formulate one extra
query that names the gap.

## Procedure (the LLM call)

This is the single LLM call the stage makes. Execute the steps in order.

1. **Parse the plan.** Extract `sub_questions`. If the plan has zero
   sub-questions, emit `sources: []` with a `coverage` array noting the
   empty plan, and stop. Do not invent sub-questions.
2. **Formulate queries.** For each sub-question, write queries per the
   Depth → query budget table. Rules:
   - Each query must be self-contained — no pronouns like "this" or "that"
     that depend on the plan text.
   - **Keyword variant** uses the sub-question's core entities + a
     constraining term (year, jurisdiction, method, etc.). Good for
     proper nouns, spec lookups, statutes, exact quotes.
   - **Semantic/intent variant** rephrases the underlying question for an
     embeddings-based backend. Good for conceptual coverage.
   - **Recency-filtered variant** (deep only) adds a recency clause
     (e.g. `after:2025-01-01`) when the topic is time-sensitive.
   - **Contrarian variant** (deep only) inverts the implicit hypothesis
     (e.g. "evidence against X", "limitations of Y"). Counters confirmation
     bias.
   - Avoid leading phrasing ("why X is good"); prefer neutral nouns.
   - Drop any query that duplicates another after lowercasing and trimming
     stopwords.
3. **Call the configured search backend** for each query. The skill is
   backend-agnostic: the host environment supplies the search tool (Tavily,
   Exa, Brave, Serper, Bing, Google CSE, a local index, or an MCP search
   server). The LLM does not invent results — it only emits sources that
   came back from a tool call. If no search tool is available in the
   environment, stop and emit `sources: []` with a `coverage` entry per
   sub-question noting `gap_note: "no search backend available"`.
4. **Normalise each result** into the Source-record schema. Map the
   backend's field names to the schema's required fields. Fill
   `retrieved_at` with the current UTC timestamp if the backend omits it.
   Assign a fresh `source_id` per unique URL (e.g. `s1`, `s2`, … in the
   order kept).
5. **Deduplicate.** Apply in order:
   1. Exact URL match (after stripping `utm_*`, `ref`, and trailing slash).
   2. Same `(canonical_domain, normalized_title)` near-match.
   3. Snippet 5-gram Jaccard overlap ≥ 0.7.
   When two records collide, keep the one with the higher
   `relevance_score`, or the earlier `published_date`-aware one if scores
   tie, and merge their `sub_question_ids`.
6. **Apply diversity caps.**
   - No single canonical domain may account for more than 30% of the final
     list (or more than 2 sources when the list has fewer than 7).
   - If the topic admits primary / academic sources, ensure at least one
     `source_type ∈ {paper, docs, dataset}` per sub-question when the
     backend returned any. If only blog/news sources exist, leave it and
     log it in `coverage`.
7. **Score relevance.** Use the backend score when present. If the
   backend gives none and the sub-question is well-formed, judge each
   source's `relevance_score` on a 0–1 scale based on whether the snippet
   substantively addresses the sub-question (not the topic in general).
   Drop sources scoring below 0.3.
8. **Build coverage summary.** For each sub-question, count surviving
   sources and write a one-line `gap_note` when coverage is below
   `floor(target / 2)` or when sources are all of one type or all
   paywalled.
9. **Validate before emitting** (see Validation gate).
10. **Emit** the `sources` array (and `coverage` if produced). Order:
    by sub-question order from the plan, then by `relevance_score`
    descending within each sub-question.

## Validation gate

Before emitting, the LLM call must self-check:

1. Every `sources[i].url` was produced by an actual search-tool result in
   this call. No URL was constructed, completed, or guessed.
2. Every required field in the Source-record schema is present and
   non-empty on every record.
3. Every `sources[i].sub_question_ids[*]` references a real id from
   `research_plan.sub_questions`. No orphan ids.
4. No duplicate `source_id` and no duplicate canonical URL.
5. No single canonical domain exceeds the diversity cap from step 6.
6. The total source count is within the depth target band (above the
   floor) OR the `coverage` array explains the shortfall.
7. `retrieved_at` is a valid ISO 8601 UTC timestamp.

If any check fails, fix and re-emit. If a fix is impossible (e.g. backend
returned nothing for a sub-question), keep the source list valid and
record the gap in `coverage`; do not pad with low-quality sources just to
hit the target.

## Constraints

- DO NOT fabricate URLs, titles, snippets, authors, dates, or any other
  field. Only emit what a tool returned or what is directly visible in a
  fetched page.
- DO NOT read sources in depth in this stage — that is `extract-findings`.
  Snippets are the unit of work here.
- DO NOT make claims about source content beyond restating the snippet.
- DO NOT prefer a single search backend; the skill is backend-agnostic.
  Use whatever tool the host has wired in.
- DO NOT invent sub-questions or rewrite the plan. The plan is read-only.
- DO NOT edit `workflows/researcher.yaml` or any sibling skill folder.
- DO NOT emit paywalled-only coverage silently. Surface it via `coverage`.
- DO NOT exceed the diversity cap to pad the list. Empty seats beat
  homogeneous noise.

## Failure modes and how to handle them

| Failure                                  | Action                                                                                                       |
|------------------------------------------|--------------------------------------------------------------------------------------------------------------|
| Search backend unavailable               | Emit `sources: []`; per-sub-question `coverage` entries naming the missing backend.                          |
| Sub-question returns zero results        | Log gap in `coverage`; do not invent sources or paraphrase the LLM's own knowledge as a source.              |
| All top results paywalled                | Keep them with `notes: "paywalled"`; log in `coverage` so the report can hedge or so review can flag it.     |
| Backend ranks one domain into every slot | Apply the 30% cap; drop overflow even if relevance is high.                                                  |
| Backend score missing                    | Self-score per step 7; mark `relevance_score` as the agent's estimate by leaving it numeric with no fanfare. |
| Suspected hallucinated URL               | Drop the record. Do not attempt to "fix" the URL; do not surround it with disclaimers.                       |

## Worked example

See `examples/heat-pump-retrofit.md` for a complete worked input → output
illustration. It shows: query formulation under `standard` depth (2 queries
per sub-question), the source-record schema in use, and how `coverage[]`
surfaces thin sub-questions rather than padding the source list.

## References

This skill is self-contained. Background prior-art (GPT Researcher,
Perplexity Deep Research, Gemini Deep Research, and the Tavily / Exa /
academic query-expansion literature) informed the source-record schema,
the depth → query-budget table, and the diversity cap, but is not loaded
into the model context.
