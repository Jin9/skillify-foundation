---
name: generating-gherkin-acceptance-criteria
description: >
  Translate resolved banking stories into strict BDD Gherkin acceptance criteria
  with happy, error, audit, idempotency, and tipping-off-safe scenarios. Use when
  a user asks to "write Gherkin ACs for these banking stories", "add
  idempotency scenarios to this epic", or "format the acceptance criteria". Do
  NOT use to extract notes, run compliance checks, or sweep hidden requirements.
---

# Skill: Generating Gherkin Acceptance Criteria

## Purpose

Produce testable Given/When/Then acceptance criteria for extracted, compliance-checked, ambiguity-swept banking stories.

## Input

Accept a pipeline state JSON containing completed extraction, compliance, and ambiguity results.

Before writing criteria, read `references/gherkin-rules.md`, then inspect `templates/stage4_gherkin.schema.json`.

## Procedure

1. Review each story, its business rules, compliance rows, open questions, and assumptions. Do not generate final AC for unresolved P1 blockers unless the scenario explicitly verifies the blocker or mitigation.
2. For each story, write at least one `happy` scenario and at least one `error` or `edge_case` scenario.
3. For every state-change or notification story, add `banking_grade_idempotency` and `banking_grade_audit` scenarios.
4. For every customer-facing message or rejection flow, add a `banking_grade_tipping_off` scenario that verifies forbidden vocabulary is absent and safe phrasing is used.
5. Keep each `When` to one trigger. Use concrete actors, state names, thresholds, timestamps, payload fields, and observable outcomes. If a value is unresolved, use a named TBD owner and related open question.
6. Populate or update each story's `banking_grade_concerns` seven-row matrix.
7. Return strict JSON matching `templates/stage4_gherkin.schema.json`.

## Output Contract

Return only JSON with:

- `stage: "gherkin"`
- `stage_status: "complete"` or `"blocked"`
- `stories[]` with `acceptance_criteria[]`
- `gherkin_quality_findings[]`
- `processing_metadata.ac_generation`

Any story with `banking_grade_concerns.idempotency.status: "applies"` must include a `banking_grade_idempotency` scenario.

## References

- `references/gherkin-rules.md` - scenario structure, banking-grade templates, and quality rules.
- `templates/stage4_gherkin.schema.json` - exact stage output contract.
- `examples/stage4_example.json` - compact example of valid stage output.
