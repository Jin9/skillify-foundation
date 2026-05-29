# Source-record schema

Loaded on demand by `search-sources` step 4 (normalisation). Each element of the
emitted `sources` array is an object with these fields. Required fields MUST be
present and non-empty. Optional fields are emitted only when the search backend
supplies them or they can be inferred from the result page; never fabricated.

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

## `coverage` side channel

A `coverage` array MAY be emitted alongside `sources` on the same output envelope
for the convenience of `review-report`; the workflow contract only consumes
`sources`, so `coverage` is a non-binding side channel.
Shape: `[{sub_question_id, source_count, gap_note}]`.
