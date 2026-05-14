---
artifact_type: qa-story-plan
blocking_oqs: []
coverage_status: complete
epic_id: EPIC-ORDER-FULFILL
story_id: EPIC-ORDER-FULFILL-3
test_case_count: 6
test_plan_id: TPLAN-EPIC-SHOPPILOT-MVP
---

# Story Test Plan — EPIC-ORDER-FULFILL-3

## Test Case Roster

| ID | Scenario | Type | Pyramid | Owner | Risk |
|---|---|---|---|---|---|
| TC-EPIC-ORDER-FULFILL-3-001 | EPIC-ORDER-FULFILL-3::Customer views own order with full snapshot | e2e | e2e | sdet-fulfilment | high |
| TC-EPIC-ORDER-FULFILL-3-002 | EPIC-ORDER-FULFILL-3::Snapshot price unaffected by catalog change | integration | integration | sdet-fulfilment | high |
| TC-EPIC-ORDER-FULFILL-3-003 | EPIC-ORDER-FULFILL-3::Cross-customer access is rejected | security | integration | security-tester | critical |
| TC-EPIC-ORDER-FULFILL-3-004 | EPIC-ORDER-FULFILL-3::Customer cannot advance state (read-only) | security | integration | security-tester | critical |
| TC-EPIC-ORDER-FULFILL-3-005 | EPIC-ORDER-FULFILL-3::Order-detail PII fields rendered to owner only | compliance | integration | compliance-mapper | high |
| TC-EPIC-ORDER-FULFILL-3-006 | EPIC-ORDER-FULFILL-3::Status timeline truthful display | e2e | e2e | sdet-fulfilment | medium |

## Test Case Details

### TC-EPIC-ORDER-FULFILL-3-001 — EPIC-ORDER-FULFILL-3::Customer views own order with full snapshot

- Scenario type: happy
- Test type: e2e
- Pyramid level: e2e
- Environment: env-e2e-fulfil
- Owner: sdet-fulfilment
- Reviewer role: qa-lead
- Risk: high
- Smoke subset: True
- Tags: customer-self-read, snapshot
- Compliance tags: PDPA-subject-access

#### Expected Assertions
- [ui-render] order-detail page renders order_number ORD-YYYYMMDD-NNNNNN, line snapshot (price/name/SKU), address snapshot, totals (subtotal, discount, shipping, grand_total), current status, status timeline, tracking_number
- [data-isolation] page returns only data owned by authenticated C1

### TC-EPIC-ORDER-FULFILL-3-002 — EPIC-ORDER-FULFILL-3::Snapshot price unaffected by catalog change

- Scenario type: edge_case
- Test type: integration
- Pyramid level: integration
- Environment: env-integration-fulfil
- Owner: sdet-fulfilment
- Reviewer role: qa-lead
- Risk: high
- Smoke subset: False
- Tags: section-6.8, snapshot-immutability
- Compliance tags: CPA-truthful-display

#### Expected Assertions
- [snapshot] order detail still shows 1290 for P1 line after catalog price changed to 999
- [catalog-divergence] catalog page shows new 999 price; snapshot vs catalog divergence is by design

### TC-EPIC-ORDER-FULFILL-3-003 — EPIC-ORDER-FULFILL-3::Cross-customer access is rejected

- Scenario type: banking_grade_authz
- Test type: security
- Pyramid level: integration
- Environment: env-integration-fulfil
- Owner: security-tester
- Reviewer role: security-tester
- Risk: critical
- Smoke subset: True
- Tags: banking_grade_authz, cross-tenant, iaa
- Compliance tags: PDPA-art-37

#### Expected Assertions
- [authz] C2 GET /orders/O1 (owned by C1) returns HTTP 403 (or 404 anti-enumeration)
- [no-leak] response body contains no order detail, no order number echo, no line items
- [audit] authz_denied audit row {actor:C2, attempted:'order.read', subject:O1}

### TC-EPIC-ORDER-FULFILL-3-004 — EPIC-ORDER-FULFILL-3::Customer cannot advance state (read-only)

- Scenario type: banking_grade_authz
- Test type: security
- Pyramid level: integration
- Environment: env-integration-fulfil
- Owner: security-tester
- Reviewer role: security-tester
- Risk: critical
- Smoke subset: False
- Tags: banking_grade_authz, customer-read-only, rbac

#### Expected Assertions
- [authz] customer C1 POST to state-advance endpoint on own order returns 403
- [state] order state unchanged

### TC-EPIC-ORDER-FULFILL-3-005 — EPIC-ORDER-FULFILL-3::Order-detail PII fields rendered to owner only

- Scenario type: regulatory
- Test type: compliance
- Pyramid level: integration
- Environment: env-integration-fulfil
- Owner: compliance-mapper
- Reviewer role: compliance-mapper
- Risk: high
- Smoke subset: False
- Tags: owner-only, pdpa, pii-render
- Compliance tags: PDPA-subject-access
- Depends on: TC-EPIC-ORDER-FULFILL-3-001

#### Expected Assertions
- [pii] shipping address, recipient_name, phone, email rendered only when authenticated user is the order owner

### TC-EPIC-ORDER-FULFILL-3-006 — EPIC-ORDER-FULFILL-3::Status timeline truthful display

- Scenario type: happy
- Test type: e2e
- Pyramid level: e2e
- Environment: env-e2e-fulfil
- Owner: sdet-fulfilment
- Reviewer role: qa-lead
- Risk: medium
- Smoke subset: False
- Tags: status-timeline, truthful-display
- Compliance tags: CPA-truthful-display
- Depends on: TC-EPIC-ORDER-FULFILL-3-001

#### Expected Assertions
- [ui] status timeline lists each transition with ts in customer-local TZ; matches audit rows 1:1

