# Failure modes and how to avoid them

Loaded on demand when a draft plan looks wrong. Each row maps a symptom to the
procedure step that fixes it (steps are in `SKILL.md`).

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
