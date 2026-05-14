---
artifact_type: qa-story-plan
blocking_oqs: []
coverage_status: partial
epic_id: EPIC-IDENTITY
story_id: EPIC-IDENTITY-3
test_case_count: 5
test_plan_id: TPLAN-EPIC-SHOPPILOT-MVP
---

# Story Test Plan — EPIC-IDENTITY-3

## Test Case Roster

| ID | Scenario | Type | Pyramid | Owner | Risk |
|---|---|---|---|---|---|
| TC-EPIC-IDENTITY-3-001 | Customer creates a first address and it becomes the default | e2e | e2e | SDET | high |
| TC-EPIC-IDENTITY-3-002 | Customer changes the default address | functional | integration | SDET | high |
| TC-EPIC-IDENTITY-3-003 | Delete refused when address referenced by an active order | integration | integration | SDET | high |
| TC-EPIC-IDENTITY-3-004 | Customer can only edit their own addresses | security | integration | Security-tester | critical |
| TC-EPIC-IDENTITY-3-005 | Profile-edit audit event (banking_grade_audit derived from address mutation) | compliance | integration | Compliance-mapper | high |

## Test Case Details

### TC-EPIC-IDENTITY-3-001 — Customer creates a first address and it becomes the default

- Scenario type: happy
- Test type: e2e
- Pyramid level: e2e
- Environment: ENV-IDENTITY-e2e
- Owner: SDET
- Reviewer role: qa-lead
- Risk: high
- Smoke subset: True
- Tags: address-book, default-invariant, happy, must
- Compliance tags: PDPA-TH-2019
- Depends on: TC-EPIC-IDENTITY-2-001

#### Expected Assertions
- [db_state] new address row exists with is_default=true and matches fixture Thai postal fields
- [downstream_visibility] address is selectable at checkout

### TC-EPIC-IDENTITY-3-002 — Customer changes the default address

- Scenario type: happy
- Test type: functional
- Pyramid level: integration
- Environment: ENV-IDENTITY-integration
- Owner: SDET
- Reviewer role: qa-lead
- Risk: high
- Smoke subset: False
- Tags: address-book, default-invariant
- Compliance tags: PDPA-TH-2019
- Depends on: TC-EPIC-IDENTITY-3-001

#### Expected Assertions
- [db_state] address B.is_default=true after operation
- [db_state] address A.is_default=false after operation
- [invariant] exactly one address has is_default=true at all times (DB count predicate)

### TC-EPIC-IDENTITY-3-003 — Delete refused when address referenced by an active order

- Scenario type: error
- Test type: integration
- Pyramid level: integration
- Environment: ENV-IDENTITY-integration
- Owner: SDET
- Reviewer role: qa-lead
- Risk: high
- Smoke subset: False
- Tags: address-book, error, referential-integrity
- Compliance tags: PDPA-TH-2019
- Depends on: TC-EPIC-IDENTITY-3-001

#### Expected Assertions
- [response_body] delete refused with human-readable message naming the blocking order number
- [db_state] address row remains in customer's address book

### TC-EPIC-IDENTITY-3-004 — Customer can only edit their own addresses

- Scenario type: banking_grade_authz
- Test type: security
- Pyramid level: integration
- Environment: ENV-IDENTITY-integration
- Owner: Security-tester
- Reviewer role: security-reviewer
- Risk: critical
- Smoke subset: False
- Tags: banking_grade_authz, cross-customer, ownership
- Compliance tags: PDPA-TH-2019, banking_grade_authz
- Depends on: TC-EPIC-IDENTITY-3-001

#### Expected Assertions
- [response_status] C2's edit attempt on C1's address returns HTTP 403 (or domain equivalent)
- [db_state] no state transition on address row
- [audit_log_field] audit event 'authz_denied' emitted with actor=C2, attempted='address.edit', subject='address:A', required='owner'

### TC-EPIC-IDENTITY-3-005 — Profile-edit audit event (banking_grade_audit derived from address mutation)

- Scenario type: banking_grade_audit
- Test type: compliance
- Pyramid level: integration
- Environment: ENV-IDENTITY-integration
- Owner: Compliance-mapper
- Reviewer role: compliance-officer
- Risk: high
- Smoke subset: False
- Tags: address-mutation, banking_grade_audit, profile-edit
- Compliance tags: PDPA-TH-2019, banking_grade_audit
- Depends on: TC-EPIC-IDENTITY-3-002

#### Expected Assertions
- [audit_log_field] audit event has event='customer.address.update'
- [audit_log_field] audit event has actor=customer_id
- [audit_log_field] audit event has ts ISO-8601
- [audit_log_field] audit event has before={is_default:<prior>} (no street/recipient PII)
- [audit_log_field] audit event has after={is_default:<new>}
- [audit_log_field] audit event has reason='customer_self_edit'
- [audit_log_field] audit event has idem_key from request
- [audit_log_negative] audit payload does NOT contain street, recipient name, phone, or postcode

