---
artifact_type: qa-story-plan
blocking_oqs:
  - BA-OQ-citation-status-pending
  - TL-compensation-hook-spec-pending
coverage_status: partial
epic_id: EPIC-ORDER-FULFILL
story_id: EPIC-ORDER-FULFILL-2
test_case_count: 8
test_plan_id: TPLAN-EPIC-SHOPPILOT-MVP
---

# Story Test Plan — EPIC-ORDER-FULFILL-2

## Test Case Roster

| ID | Scenario | Type | Pyramid | Owner | Risk |
|---|---|---|---|---|---|
| TC-EPIC-ORDER-FULFILL-2-001 | EPIC-ORDER-FULFILL-2::Admin cancels paid order with reason and stock restores | integration | integration | sdet-fulfilment | high |
| TC-EPIC-ORDER-FULFILL-2-002 | EPIC-ORDER-FULFILL-2::Admin cancel without reason is rejected | functional | unit | sdet-fulfilment | medium |
| TC-EPIC-ORDER-FULFILL-2-003 | EPIC-ORDER-FULFILL-2::Cannot cancel after packing has started | functional | unit | sdet-fulfilment | high |
| TC-EPIC-ORDER-FULFILL-2-004 | EPIC-ORDER-FULFILL-2::Manual refund expectation captured but NOT auto-disbursed | integration | integration | sdet-fulfilment | critical |
| TC-EPIC-ORDER-FULFILL-2-005 | EPIC-ORDER-FULFILL-2::OOB manual refund operator workflow | e2e | e2e | sdet-fulfilment | high |
| TC-EPIC-ORDER-FULFILL-2-006 | EPIC-ORDER-FULFILL-2::Audit on admin-cancel | contract | contract | sdet-fulfilment | critical |
| TC-EPIC-ORDER-FULFILL-2-007 | EPIC-ORDER-FULFILL-2::Admin-only cancel authz | security | integration | security-tester | critical |
| TC-EPIC-ORDER-FULFILL-2-008 | EPIC-ORDER-FULFILL-2::Idempotent admin-cancel replay | integration | integration | sdet-fulfilment | high |

## Test Case Details

### TC-EPIC-ORDER-FULFILL-2-001 — EPIC-ORDER-FULFILL-2::Admin cancels paid order with reason and stock restores

- Scenario type: happy
- Test type: integration
- Pyramid level: integration
- Environment: env-integration-fulfil
- Owner: sdet-fulfilment
- Reviewer role: qa-lead
- Risk: high
- Smoke subset: True
- Tags: admin-cancel, state-machine, stock-restore

#### Expected Assertions
- [state] O1 transitions from 'paid' to 'cancelled'
- [stock] P1.stock_sold decreased by 2; P1.stock_available increased by 2 (compensation hook executed)
- [audit] audit row 'order.admin_cancelled' before='paid', after='cancelled', reason='customer_phone_request' captured

### TC-EPIC-ORDER-FULFILL-2-002 — EPIC-ORDER-FULFILL-2::Admin cancel without reason is rejected

- Scenario type: error
- Test type: functional
- Pyramid level: unit
- Environment: env-integration-fulfil
- Owner: sdet-fulfilment
- Reviewer role: qa-lead
- Risk: medium
- Smoke subset: False
- Tags: field-required, validation

#### Expected Assertions
- [error] request rejected with reason='cancel_reason_required'
- [state] O1 remains in 'paid'

### TC-EPIC-ORDER-FULFILL-2-003 — EPIC-ORDER-FULFILL-2::Cannot cancel after packing has started

- Scenario type: edge_case
- Test type: functional
- Pyramid level: unit
- Environment: env-integration-fulfil
- Owner: sdet-fulfilment
- Reviewer role: qa-lead
- Risk: high
- Smoke subset: False
- Tags: boundary, invalid-transition, state-machine

#### Expected Assertions
- [error] request rejected with reason='cannot_cancel_after_packing'
- [state] O1 remains in 'packing'

### TC-EPIC-ORDER-FULFILL-2-004 — EPIC-ORDER-FULFILL-2::Manual refund expectation captured but NOT auto-disbursed

- Scenario type: banking_grade_reversibility
- Test type: integration
- Pyramid level: integration
- Environment: env-integration-fulfil
- Owner: sdet-fulfilment
- Reviewer role: compliance-mapper
- Risk: critical
- Smoke subset: False
- Tags: banking_grade_reversibility, compensating-action, out-of-band-refund
- Compliance tags: CPA-refund-disclosure
- Depends on: TC-EPIC-ORDER-FULFILL-2-001

#### Expected Assertions
- [audit] audit row records prior payment captured (money taken) - compensation required OOB
- [no-auto-refund] no payment gateway refund API is invoked by the system on cancellation
- [notification] notification log row 'order.cancelled' created for Finance team to disburse manually

### TC-EPIC-ORDER-FULFILL-2-005 — EPIC-ORDER-FULFILL-2::OOB manual refund operator workflow

- Scenario type: banking_grade_reversibility
- Test type: e2e
- Pyramid level: e2e
- Environment: env-e2e-fulfil
- Owner: sdet-fulfilment
- Reviewer role: compliance-mapper
- Risk: high
- Smoke subset: False
- Tags: banking_grade_reversibility, oob-workflow, operator-triggered
- Compliance tags: CPA-refund-disclosure, TRC-record-keeping
- Depends on: TC-EPIC-ORDER-FULFILL-2-004

#### Expected Assertions
- [workflow] operator marks notification log row 'refund_disbursed_at' with timestamp + external_ref; system records compensating-action audit row
- [audit] compensating-action audit row contains actor (operator), original_order_id, refund_external_ref, ts, idem_key
- [no-state-change] order remains 'cancelled' (OOB does not change order state, only annotates compensation)

### TC-EPIC-ORDER-FULFILL-2-006 — EPIC-ORDER-FULFILL-2::Audit on admin-cancel

- Scenario type: banking_grade_audit
- Test type: contract
- Pyramid level: contract
- Environment: env-integration-fulfil
- Owner: sdet-fulfilment
- Reviewer role: compliance-mapper
- Risk: critical
- Smoke subset: False
- Tags: AP-Q12, banking_grade_audit
- Compliance tags: TRC-record-keeping
- Depends on: TC-EPIC-ORDER-FULFILL-2-001

#### Expected Assertions
- [audit-fields] admin_cancelled audit row contains 7 mandatory fields plus reason text

### TC-EPIC-ORDER-FULFILL-2-007 — EPIC-ORDER-FULFILL-2::Admin-only cancel authz

- Scenario type: banking_grade_authz
- Test type: security
- Pyramid level: integration
- Environment: env-integration-fulfil
- Owner: security-tester
- Reviewer role: security-tester
- Risk: critical
- Smoke subset: False
- Tags: admin-only, banking_grade_authz, rbac

#### Expected Assertions
- [authz] customer or other-tenant admin attempts to cancel are rejected with HTTP 403
- [audit] authz_denied audit emitted with actor, attempted action, subject

### TC-EPIC-ORDER-FULFILL-2-008 — EPIC-ORDER-FULFILL-2::Idempotent admin-cancel replay

- Scenario type: banking_grade_idempotency
- Test type: integration
- Pyramid level: integration
- Environment: env-integration-fulfil
- Owner: sdet-fulfilment
- Reviewer role: qa-lead
- Risk: high
- Smoke subset: False
- Tags: banking_grade_idempotency, replay-safe
- Depends on: TC-EPIC-ORDER-FULFILL-2-001

#### Expected Assertions
- [state] replay with same idem_key does not double-cancel; stock not double-restored
- [audit-dedup] exactly one admin_cancelled audit row exists

