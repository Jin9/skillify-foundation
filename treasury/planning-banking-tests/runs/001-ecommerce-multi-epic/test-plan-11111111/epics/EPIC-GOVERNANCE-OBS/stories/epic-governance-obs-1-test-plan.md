---
artifact_type: qa-story-plan
blocking_oqs:
  - OQ-20
coverage_status: partial
epic_id: EPIC-GOVERNANCE-OBS
story_id: EPIC-GOVERNANCE-OBS-1
test_case_count: 7
test_plan_id: TPLAN-EPIC-SHOPPILOT-MVP
---

# Story Test Plan — EPIC-GOVERNANCE-OBS-1

## Test Case Roster

| ID | Scenario | Type | Pyramid | Owner | Risk |
|---|---|---|---|---|---|
| TC-EPIC-GOVERNANCE-OBS-1-001 | Stock increase with reason is recorded | integration | integration | SDET | high |
| TC-EPIC-GOVERNANCE-OBS-1-002 | Adjustment without reason is rejected | functional | integration | SDET | medium |
| TC-EPIC-GOVERNANCE-OBS-1-003 | Cannot reduce stock below currently-reserved | integration | integration | SDET | high |
| TC-EPIC-GOVERNANCE-OBS-1-004 | Stock can never become negative under any path | integration | integration | SDET | high |
| TC-EPIC-GOVERNANCE-OBS-1-005 | Banking-grade audit assertion: 7-field audit row on inventory adjust | integration | integration | SDET | critical |
| TC-EPIC-GOVERNANCE-OBS-1-006 | Banking-grade authz on admin-only inventory-adjust endpoint | security | integration | SDET | high |
| TC-EPIC-GOVERNANCE-OBS-1-007 | Idempotency: double-click inventory adjust does not double-apply | integration | integration | SDET | high |

## Test Case Details

### TC-EPIC-GOVERNANCE-OBS-1-001 — Stock increase with reason is recorded

- Scenario type: happy
- Test type: integration
- Pyramid level: integration
- Environment: env-gov-integration
- Owner: SDET
- Reviewer role: QA-Lead
- Risk: high
- Smoke subset: True
- Tags: audit-emission, happy, inventory-adjust
- Compliance tags: CCA-TH-LOG-90D

#### Expected Assertions
- [state_change] P1 stock_available transitions from 5 to 15 after +10 adjustment with reason='restocked_from_supplier'.
- [audit_event_emitted] Audit event 'inventory.adjusted' enumerates all 7 fields: event='inventory.adjusted', actor='admin:aoy', ts (ISO-8601), before={stock_available:5}, after={stock_available:15}, reason='restocked_from_supplier', idem_key (non-empty). (AP-Q12)

### TC-EPIC-GOVERNANCE-OBS-1-002 — Adjustment without reason is rejected

- Scenario type: error
- Test type: functional
- Pyramid level: integration
- Environment: env-gov-integration
- Owner: SDET
- Reviewer role: QA-Lead
- Risk: medium
- Smoke subset: False
- Tags: error, inventory-adjust, validation

#### Expected Assertions
- [error_response] Request rejected with field-level error code 'reason_required'.
- [state_invariant] P1 stock_available is unchanged after rejected adjustment.
- [no_audit_event] No 'inventory.adjusted' audit event is emitted for a rejected adjustment.

### TC-EPIC-GOVERNANCE-OBS-1-003 — Cannot reduce stock below currently-reserved

- Scenario type: edge_case
- Test type: integration
- Pyramid level: integration
- Environment: env-gov-integration
- Owner: SDET
- Reviewer role: QA-Lead
- Risk: high
- Smoke subset: False
- Tags: edge, invariants, inventory-adjust

#### Expected Assertions
- [error_response] Adjustment rejected with reason='would_break_reservations' when invariant violated.
- [state_invariant] stock_available + stock_reserved >= existing reservations remains true.

### TC-EPIC-GOVERNANCE-OBS-1-004 — Stock can never become negative under any path

- Scenario type: edge_case
- Test type: integration
- Pyramid level: integration
- Environment: env-gov-integration
- Owner: SDET
- Reviewer role: QA-Lead
- Risk: high
- Smoke subset: False
- Tags: edge, invariants, inventory-adjust

#### Expected Assertions
- [error_response] Adjustment rejected with reason='stock_cannot_go_negative'.
- [state_invariant] P1 stock_available stays at 1 after rejected -5 adjustment.

### TC-EPIC-GOVERNANCE-OBS-1-005 — Banking-grade audit assertion: 7-field audit row on inventory adjust

- Scenario type: banking_grade_audit
- Test type: integration
- Pyramid level: integration
- Environment: env-gov-integration
- Owner: SDET
- Reviewer role: Compliance-Lead
- Risk: critical
- Smoke subset: True
- Tags: AP-Q12, banking_grade_audit, cross-cutting
- Compliance tags: CCA-TH-LOG-90D

#### Expected Assertions
- [audit_field_event] Field 'event' equals 'inventory.adjusted' (AP-Q12 field 1).
- [audit_field_actor] Field 'actor' equals 'admin:aoy' (AP-Q12 field 2).
- [audit_field_ts] Field 'ts' is ISO-8601 timestamp (AP-Q12 field 3).
- [audit_field_before] Field 'before' captures pre-state snapshot (AP-Q12 field 4).
- [audit_field_after] Field 'after' captures post-state snapshot (AP-Q12 field 5).
- [audit_field_reason] Field 'reason' is non-empty mandatory string (AP-Q12 field 6).
- [audit_field_idem_key] Field 'idem_key' present and non-empty (AP-Q12 field 7).

### TC-EPIC-GOVERNANCE-OBS-1-006 — Banking-grade authz on admin-only inventory-adjust endpoint

- Scenario type: banking_grade_authz
- Test type: security
- Pyramid level: integration
- Environment: env-gov-integration
- Owner: SDET
- Reviewer role: Compliance-Lead
- Risk: high
- Smoke subset: False
- Tags: admin-only, banking_grade_authz
- Compliance tags: CCA-TH-LOG-90D

#### Expected Assertions
- [http_status] Non-admin caller receives HTTP 403 (or 404) on POST /admin/inventory/adjust.
- [audit_event_emitted] authz_denied audit event emitted with 7 standard fields (event='authz_denied', actor, ts, before={authz_state}, after={authz_state}, reason='not_admin', idem_key).

### TC-EPIC-GOVERNANCE-OBS-1-007 — Idempotency: double-click inventory adjust does not double-apply

- Scenario type: banking_grade_idempotency
- Test type: integration
- Pyramid level: integration
- Environment: env-gov-integration
- Owner: SDET
- Reviewer role: QA-Lead
- Risk: high
- Smoke subset: False
- Tags: banking_grade_idempotency, inventory-adjust

#### Expected Assertions
- [state_invariant] Repeated POST with same idem_key produces stock delta of +10 once, not twice.
- [audit_event_count] Exactly one 'inventory.adjusted' audit event with that idem_key exists.

