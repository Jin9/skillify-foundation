# Troubleshooting

Lookup table for common signals encountered while drafting. The default
guidance for any unlisted edge case: refuse to invent — surface the gap in
`## Limitations & Open Questions` and let `review-report` decide what to do.

| Signal | Action |
|---|---|
| Input `findings` is empty | Refuse to draft. Emit a one-paragraph report under `## Limitations` stating that no findings were available, with an empty `cited_sources: []`. Do not invent. |
| `research_plan.depth` is missing | Default to `standard`. Note nothing in the report; review-report can flag. |
| `audience` is an unrecognized free-form string | Map to `general` (per [`audience-profiles.md`](audience-profiles.md)) and proceed. |
| Two findings make contradicting claims about the same fact | Cite both, name the contradiction in prose, flag in Limitations. Do not pick a winner. |
| Section runs short on findings | Cut the section or merge with an adjacent one. Do not pad. |
| One sub-question has no supporting findings | Add a one-line note in `## Limitations & Open Questions` ("No evidence found on <sub-question>."). Do NOT fabricate. |
| A finding cites a source not in the input source list | Refuse to cite that finding. (This is an upstream defect — review-report or extract-findings should catch it.) |
| Length would exceed the depth's upper bound | Cut the lowest-scored body section first; if still over, tighten paragraphs. Do not silently drop citations. |
