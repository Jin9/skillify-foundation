# Downstream contract

Loaded on demand when a downstream stage breaks on a field this skill produced.
Each later stage consumes these fields from `research_plan`:

| Stage                | Consumes from `research_plan`                                                                |
|----------------------|----------------------------------------------------------------------------------------------|
| `search-sources`     | `sub_questions[].question`, `expected_source_types`, `out_of_scope`, `notes_for_downstream`. |
| `extract-findings`   | `sub_questions[].id` (grouping key), `sub_questions[].question`, `thesis_question`.          |
| `synthesize-report`  | `thesis_question`, `coverage_dimensions`, `sub_questions[]`, `audience`.                     |
| `review-report`      | `thesis_question`, `success_criteria`, `out_of_scope`, `audience`.                           |

If a downstream stage breaks on a field this skill produced, fix it here (and
bump the version), not there. The plan is the contract.
