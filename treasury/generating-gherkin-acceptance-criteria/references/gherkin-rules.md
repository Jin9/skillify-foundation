# Gherkin Rules

## Scenario Shape

Each scenario must include:

- `scenario_name`: concise, test-intent name.
- `scenario_type`: one of the schema enum values.
- `given[]`: concrete starting state and actor.
- `when`: exactly one trigger or user/system action.
- `then[]`: observable outcomes, state changes, emitted events, visible messages, or payload fields.

Avoid subjective assertions such as "is happy", "works", "fast", or "appropriate".

## Required Scenario Set

Every story needs:

- at least one `happy` scenario
- at least one `error` or `edge_case` scenario
- at least one banking-grade scenario when banking-grade concerns apply

State-change or notification stories require:

- `banking_grade_idempotency`: replay with same idempotency key creates no duplicate effect.
- `banking_grade_audit`: audit event includes event, actor, timestamp, before, after, reason, and idempotency key.

Customer-facing messages require:

- `banking_grade_tipping_off`: customer copy avoids forbidden terms and uses approved neutral phrasing.

## Unresolved Values

Do not invent thresholds, dates, owners, regulator citations, or policy values. Use a named TBD token plus the related open question, for example `TBD_by_Legal_OQ-3`.

## Banking-Grade Matrix

Each story must retain all seven rows: `pii_fields`, `audit_events`, `idempotency`, `reversibility`, `authn_authz`, `regulatory`, and `tipping_off`. Each row needs `status` and a justification of at least 10 characters.
