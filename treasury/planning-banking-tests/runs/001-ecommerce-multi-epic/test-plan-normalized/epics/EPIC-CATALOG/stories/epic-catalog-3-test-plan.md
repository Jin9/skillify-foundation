---
artifact_type: qa-story-plan
blocking_oqs:
  - OQ-34
  - OQ-35
coverage_status: partial
epic_id: EPIC-CATALOG
story_id: EPIC-CATALOG-3
test_case_count: 7
test_plan_id: TPLAN-EPIC-SHOPPILOT-MVP
---

# Story Test Plan — EPIC-CATALOG-3

## Test Case Roster

| ID | Scenario | Type | Pyramid | Owner | Risk |
|---|---|---|---|---|---|
| TC-EPIC-CATALOG-3-001 | Admin creates a product with all required fields | integration | integration | SDET | high |
| TC-EPIC-CATALOG-3-002 | Duplicate SKU is rejected | integration | integration | SDET | medium |
| TC-EPIC-CATALOG-3-003 | Price change emits an audit event | integration | integration | SDET | critical |
| TC-EPIC-CATALOG-3-004 | Edit attempt to mutate SKU is rejected | integration | integration | SDET | high |
| TC-EPIC-CATALOG-3-005 | Idempotent product-edit replay produces no duplicate audit | integration | integration | SDET | critical |
| TC-EPIC-CATALOG-3-006 | Non-admin caller is rejected from admin product CRUD surface (banking_grade authn_authz) | security | integration | Security-Tester | critical |
| TC-EPIC-CATALOG-3-007 | Soft-delete and price edit are reversible by inverse admin action with full audit trail (banking_grade reversibility) | integration | integration | SDET | high |

## Test Case Details

### TC-EPIC-CATALOG-3-001 — Admin creates a product with all required fields

- Scenario type: happy
- Test type: integration
- Pyramid level: integration
- Environment: env-integration
- Owner: SDET
- Reviewer role: QA-Peer
- Risk: high
- Smoke subset: True
- Tags: admin, happy, product-create

#### Expected Assertions
- [state] a new product row is created with the submitted values
- [state] the response includes the product id
- [state] the product appears in the admin catalog list with status=active

### TC-EPIC-CATALOG-3-002 — Duplicate SKU is rejected

- Scenario type: error
- Test type: integration
- Pyramid level: integration
- Environment: env-integration
- Owner: SDET
- Reviewer role: QA-Peer
- Risk: medium
- Smoke subset: False
- Tags: admin, error, sku, uniqueness
- Depends on: TC-EPIC-CATALOG-3-001

#### Expected Assertions
- [state] the request is rejected with a field-level error on SKU
- [state] no new product is created

### TC-EPIC-CATALOG-3-003 — Price change emits an audit event

- Scenario type: banking_grade_audit
- Test type: integration
- Pyramid level: integration
- Environment: env-integration
- Owner: SDET
- Reviewer role: QA-Lead
- Risk: critical
- Smoke subset: False
- Tags: admin, audit, banking_grade_audit_events, price-change
- Depends on: TC-EPIC-CATALOG-3-001

#### Expected Assertions
- [state] product P1 has price=999 after the edit is saved
- [audit] an audit event is emitted with payload containing event='product.price_changed', actor, ts (ISO-8601), before.price=1290, after.price=999, reason, and idem_key
- [audit] the audit event payload does NOT contain any password, token, or card data

### TC-EPIC-CATALOG-3-004 — Edit attempt to mutate SKU is rejected

- Scenario type: error
- Test type: integration
- Pyramid level: integration
- Environment: env-integration
- Owner: SDET
- Reviewer role: QA-Peer
- Risk: high
- Smoke subset: False
- Tags: admin, error, immutable, sku
- Depends on: TC-EPIC-CATALOG-3-001

#### Expected Assertions
- [state] the request is rejected with a field-level error stating SKU is immutable
- [state] no change is persisted
- [audit] the audit log records no successful SKU change

### TC-EPIC-CATALOG-3-005 — Idempotent product-edit replay produces no duplicate audit

- Scenario type: banking_grade_idempotency
- Test type: integration
- Pyramid level: integration
- Environment: env-integration
- Owner: SDET
- Reviewer role: QA-Lead
- Risk: critical
- Smoke subset: False
- Tags: admin, banking_grade_idempotency, replay
- Depends on: TC-EPIC-CATALOG-3-003

#### Expected Assertions
- [state] no additional state change occurs when the same idempotency_key='ik-prod-001' request is replayed
- [audit] no duplicate audit event is emitted for the replayed request
- [state] the response returns the original result

### TC-EPIC-CATALOG-3-006 — Non-admin caller is rejected from admin product CRUD surface (banking_grade authn_authz)

- Scenario type: banking_grade_authz
- Test type: security
- Pyramid level: integration
- Environment: env-integration
- Owner: Security-Tester
- Reviewer role: Security-Reviewer
- Risk: critical
- Smoke subset: False
- Tags: admin, authz, banking_grade_authn_authz
- Depends on: TC-EPIC-CATALOG-3-001

#### Expected Assertions
- [authz] guest caller invoking the admin create-product endpoint receives 401/403 and no product row is created
- [authz] customer-role caller invoking the admin edit-product endpoint receives 403 and no state change is persisted
- [audit] rejected admin-endpoint calls do not emit a successful-write audit event

### TC-EPIC-CATALOG-3-007 — Soft-delete and price edit are reversible by inverse admin action with full audit trail (banking_grade reversibility)

- Scenario type: banking_grade_reversibility
- Test type: integration
- Pyramid level: integration
- Environment: env-integration
- Owner: SDET
- Reviewer role: QA-Lead
- Risk: high
- Smoke subset: False
- Tags: admin, audit, banking_grade_reversibility, soft-delete
- Depends on: TC-EPIC-CATALOG-3-003

#### Expected Assertions
- [state] flipping status from active to paused and back to active leaves the product row in the original active state
- [state] re-editing price from 999 back to 1290 restores the original price value
- [audit] each reversing action emits its own audit event with before/after deltas; no audit row is mutated or removed

