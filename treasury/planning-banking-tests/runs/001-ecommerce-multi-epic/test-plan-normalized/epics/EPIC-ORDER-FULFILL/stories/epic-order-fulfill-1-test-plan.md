---
artifact_type: qa-story-plan
blocking_oqs:
  - BA-OQ-citation-status-pending
  - TL-state-machine-model-pending
coverage_status: partial
epic_id: EPIC-ORDER-FULFILL
story_id: EPIC-ORDER-FULFILL-1
test_case_count: 12
test_plan_id: TPLAN-EPIC-SHOPPILOT-MVP
---

# Story Test Plan — EPIC-ORDER-FULFILL-1

## Test Case Roster

| ID | Scenario | Type | Pyramid | Owner | Risk |
|---|---|---|---|---|---|
| TC-EPIC-ORDER-FULFILL-1-001 | EPIC-ORDER-FULFILL-1::Advance paid order to packing | integration | integration | sdet-fulfilment | high |
| TC-EPIC-ORDER-FULFILL-1-002 | EPIC-ORDER-FULFILL-1::Advance paid order to packing | integration | integration | sdet-fulfilment | critical |
| TC-EPIC-ORDER-FULFILL-1-003 | EPIC-ORDER-FULFILL-1::Shipped transition requires tracking number | integration | integration | sdet-fulfilment | high |
| TC-EPIC-ORDER-FULFILL-1-004 | EPIC-ORDER-FULFILL-1::Shipped transition without tracking number is rejected | functional | unit | sdet-fulfilment | medium |
| TC-EPIC-ORDER-FULFILL-1-005 | EPIC-ORDER-FULFILL-1::Shipped transition happy::packing->shipped | integration | integration | sdet-fulfilment | high |
| TC-EPIC-ORDER-FULFILL-1-006 | EPIC-ORDER-FULFILL-1::Backward transition is rejected | functional | unit | sdet-fulfilment | high |
| TC-EPIC-ORDER-FULFILL-1-007 | EPIC-ORDER-FULFILL-1::Irreversibility - paid->shipping cannot be reversed | integration | integration | sdet-fulfilment | critical |
| TC-EPIC-ORDER-FULFILL-1-008 | EPIC-ORDER-FULFILL-1::Delivered terminal - re-cancel rejected | functional | unit | sdet-fulfilment | high |
| TC-EPIC-ORDER-FULFILL-1-009 | EPIC-ORDER-FULFILL-1::Idempotent transition replay produces no duplicate audit | integration | integration | sdet-fulfilment | critical |
| TC-EPIC-ORDER-FULFILL-1-010 | EPIC-ORDER-FULFILL-1::Admin-only state advancement | security | integration | security-tester | critical |
| TC-EPIC-ORDER-FULFILL-1-011 | EPIC-ORDER-FULFILL-1::AP-Q11 no tipping-off in transition errors | security | integration | security-tester | medium |
| TC-EPIC-ORDER-FULFILL-1-012 | EPIC-ORDER-FULFILL-1::Audit field contract on every flip | contract | contract | sdet-fulfilment | critical |

## Test Case Details

### TC-EPIC-ORDER-FULFILL-1-001 — EPIC-ORDER-FULFILL-1::Advance paid order to packing

- Scenario type: happy
- Test type: integration
- Pyramid level: integration
- Environment: env-integration-fulfil
- Owner: sdet-fulfilment
- Reviewer role: qa-lead
- Risk: high
- Smoke subset: True
- Tags: state-machine, transition-paid-to-packing

#### Expected Assertions
- [state] order O1.status transitions from 'paid' to 'packing' atomically
- [audit] audit row 'order.status_changed' written with before='paid', after='packing', actor='admin:aoy'

### TC-EPIC-ORDER-FULFILL-1-002 — EPIC-ORDER-FULFILL-1::Advance paid order to packing

- Scenario type: banking_grade_audit
- Test type: integration
- Pyramid level: integration
- Environment: env-integration-fulfil
- Owner: sdet-fulfilment
- Reviewer role: compliance-mapper
- Risk: critical
- Smoke subset: False
- Tags: AP-Q12, banking_grade_audit
- Compliance tags: PDPA-art-39, TRC-record-keeping
- Depends on: TC-EPIC-ORDER-FULFILL-1-001

#### Expected Assertions
- [audit-fields] audit event contains all 7 mandatory fields: actor, action, subject_id, before_state, after_state, ts (ISO-8601 UTC), idem_key
- [audit-immutability] audit row insert-only; no UPDATE/DELETE permitted on audit table

### TC-EPIC-ORDER-FULFILL-1-003 — EPIC-ORDER-FULFILL-1::Shipped transition requires tracking number

- Scenario type: happy
- Test type: integration
- Pyramid level: integration
- Environment: env-integration-fulfil
- Owner: sdet-fulfilment
- Reviewer role: qa-lead
- Risk: high
- Smoke subset: True
- Tags: state-machine, transition-packing-to-shipped
- Depends on: TC-EPIC-ORDER-FULFILL-1-001

#### Expected Assertions
- [state] order O1.status='shipped' and O1.tracking_number='TH123456789' persisted
- [audit] audit row 'order.status_changed' before='packing', after='shipped' with tracking_number captured
- [ui] customer order-detail page renders tracking_number after transition

### TC-EPIC-ORDER-FULFILL-1-004 — EPIC-ORDER-FULFILL-1::Shipped transition without tracking number is rejected

- Scenario type: error
- Test type: functional
- Pyramid level: unit
- Environment: env-integration-fulfil
- Owner: sdet-fulfilment
- Reviewer role: qa-lead
- Risk: medium
- Smoke subset: False
- Tags: field-required, state-machine, validation

#### Expected Assertions
- [error] response returns field-level error code 'tracking_number_required'
- [state] O1 remains in 'packing' (no transition persisted)
- [audit] no 'order.status_changed' audit row created; only optional authz_denied/validation log

### TC-EPIC-ORDER-FULFILL-1-005 — EPIC-ORDER-FULFILL-1::Shipped transition happy::packing->shipped

- Scenario type: happy
- Test type: integration
- Pyramid level: integration
- Environment: env-integration-fulfil
- Owner: sdet-fulfilment
- Reviewer role: qa-lead
- Risk: high
- Smoke subset: True
- Tags: state-machine, transition-shipped-to-delivered
- Depends on: TC-EPIC-ORDER-FULFILL-1-003

#### Expected Assertions
- [state] order O1 transitions from 'shipped' to 'delivered'; delivered is terminal
- [audit] audit row before='shipped', after='delivered' with all 7 fields

### TC-EPIC-ORDER-FULFILL-1-006 — EPIC-ORDER-FULFILL-1::Backward transition is rejected

- Scenario type: error
- Test type: functional
- Pyramid level: unit
- Environment: env-integration-fulfil
- Owner: sdet-fulfilment
- Reviewer role: qa-lead
- Risk: high
- Smoke subset: False
- Tags: invalid-transition, irreversible, state-machine

#### Expected Assertions
- [error] attempt to move O1 from 'shipped' back to 'packing' rejected with reason='invalid_transition'
- [state] O1 remains in 'shipped'

### TC-EPIC-ORDER-FULFILL-1-007 — EPIC-ORDER-FULFILL-1::Irreversibility - paid->shipping cannot be reversed

- Scenario type: banking_grade_reversibility
- Test type: integration
- Pyramid level: integration
- Environment: env-integration-fulfil
- Owner: sdet-fulfilment
- Reviewer role: compliance-mapper
- Risk: critical
- Smoke subset: False
- Tags: banking_grade_reversibility, irreversible, no-api-reverse
- Depends on: TC-EPIC-ORDER-FULFILL-1-003

#### Expected Assertions
- [no-api-reverse] no API endpoint exists to move status from 'shipped' back to 'packing' or 'paid'; attempt returns 'invalid_transition'
- [compensating-action] system documents that post-shipped corrections are OOB manual (no auto-compensation)

### TC-EPIC-ORDER-FULFILL-1-008 — EPIC-ORDER-FULFILL-1::Delivered terminal - re-cancel rejected

- Scenario type: banking_grade_reversibility
- Test type: functional
- Pyramid level: unit
- Environment: env-integration-fulfil
- Owner: sdet-fulfilment
- Reviewer role: qa-lead
- Risk: high
- Smoke subset: False
- Tags: banking_grade_reversibility, no-api-reverse, terminal-state
- Depends on: TC-EPIC-ORDER-FULFILL-1-005

#### Expected Assertions
- [error] cancel request on delivered order rejected with reason='order_terminal'
- [state] O1 remains in 'delivered' (terminal)
- [explicit-doc] this transition cannot be reversed via API - documented and tested

### TC-EPIC-ORDER-FULFILL-1-009 — EPIC-ORDER-FULFILL-1::Idempotent transition replay produces no duplicate audit

- Scenario type: banking_grade_idempotency
- Test type: integration
- Pyramid level: integration
- Environment: env-integration-fulfil
- Owner: sdet-fulfilment
- Reviewer role: qa-lead
- Risk: critical
- Smoke subset: False
- Tags: banking_grade_idempotency, replay-safe
- Depends on: TC-EPIC-ORDER-FULFILL-1-001

#### Expected Assertions
- [state] after replay with same idem_key='ik-fulfil-001', O1 remains in 'packing'
- [audit-dedup] exactly one audit row exists for the transition; replay does NOT emit duplicate
- [response] replay returns the original stored response payload (same body, same idempotent outcome)

### TC-EPIC-ORDER-FULFILL-1-010 — EPIC-ORDER-FULFILL-1::Admin-only state advancement

- Scenario type: banking_grade_authz
- Test type: security
- Pyramid level: integration
- Environment: env-integration-fulfil
- Owner: security-tester
- Reviewer role: security-tester
- Risk: critical
- Smoke subset: True
- Tags: admin-only, banking_grade_authz, rbac
- Compliance tags: PDPA-art-37

#### Expected Assertions
- [authz] customer-role attempt to call state-advance endpoint returns HTTP 403
- [state] O1 state unchanged after customer attempt
- [audit] authz_denied audit row emitted with actor (customer), attempted action, subject

### TC-EPIC-ORDER-FULFILL-1-011 — EPIC-ORDER-FULFILL-1::AP-Q11 no tipping-off in transition errors

- Scenario type: error
- Test type: security
- Pyramid level: integration
- Environment: env-integration-fulfil
- Owner: security-tester
- Reviewer role: compliance-mapper
- Risk: medium
- Smoke subset: False
- Tags: AP-Q11, error-copy-review, no-tipping-off

#### Expected Assertions
- [copy-review] invalid_transition and order_terminal error copy contains no regulatory state hints (status copy is operational only per BA tipping_off=not_applicable)

### TC-EPIC-ORDER-FULFILL-1-012 — EPIC-ORDER-FULFILL-1::Audit field contract on every flip

- Scenario type: banking_grade_audit
- Test type: contract
- Pyramid level: contract
- Environment: env-integration-fulfil
- Owner: sdet-fulfilment
- Reviewer role: compliance-mapper
- Risk: critical
- Smoke subset: False
- Tags: AP-Q12, banking_grade_audit, contract
- Compliance tags: TRC-record-keeping
- Depends on: TC-EPIC-ORDER-FULFILL-1-001, TC-EPIC-ORDER-FULFILL-1-003, TC-EPIC-ORDER-FULFILL-1-005

#### Expected Assertions
- [schema] for each transition (paid->packing, packing->shipped, shipped->delivered) audit JSON validates against the 7-field schema (actor, action, subject_id, before_state, after_state, ts, idem_key)

