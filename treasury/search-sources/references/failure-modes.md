# Failure modes and how to handle them

Loaded on demand by `search-sources` when a search comes back degenerate. None of
these warrant fabricating data — degrade and record the gap in `coverage`.

| Failure                                  | Action                                                                                                       |
|------------------------------------------|--------------------------------------------------------------------------------------------------------------|
| Search backend unavailable               | Emit `sources: []`; per-sub-question `coverage` entries naming the missing backend.                          |
| Sub-question returns zero results        | Log gap in `coverage`; do not invent sources or paraphrase the LLM's own knowledge as a source.              |
| All top results paywalled                | Keep them with `notes: "paywalled"`; log in `coverage` so the report can hedge or so review can flag it.     |
| Backend ranks one domain into every slot | Apply the 30% cap; drop overflow even if relevance is high.                                                  |
| Backend score missing                    | Self-score per step 7; mark `relevance_score` as the agent's estimate by leaving it numeric with no fanfare. |
| Suspected hallucinated URL               | Drop the record. Do not attempt to "fix" the URL; do not surround it with disclaimers.                       |
