# Validation checklist

Loaded on demand by `plan-research` step 9 (self-check before emitting). Verify
all of the following; fix in place — do not make a second LLM call.

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
