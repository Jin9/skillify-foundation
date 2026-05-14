---
artifact_type: qa-story-plan
blocking_oqs: []
coverage_status: complete
epic_id: EPIC-CATALOG
story_id: EPIC-CATALOG-2
test_case_count: 4
test_plan_id: TPLAN-EPIC-SHOPPILOT-MVP
---

# Story Test Plan — EPIC-CATALOG-2

## Test Case Roster

| ID | Scenario | Type | Pyramid | Owner | Risk |
|---|---|---|---|---|---|
| TC-EPIC-CATALOG-2-001 | Active product PDP shows add-to-cart affordance | integration | integration | SDET | high |
| TC-EPIC-CATALOG-2-002 | Paused product PDP shows unavailable message instead of add-to-cart | integration | integration | SDET | high |
| TC-EPIC-CATALOG-2-003 | Soft-deleted product reached via direct URL returns 404 or unavailable view | integration | integration | SDET | medium |
| TC-EPIC-CATALOG-2-004 | PDP response never surfaces admin-only product attributes regardless of caller (banking_grade authn_authz) | security | integration | Security-Tester | high |

## Test Case Details

### TC-EPIC-CATALOG-2-001 — Active product PDP shows add-to-cart affordance

- Scenario type: happy
- Test type: integration
- Pyramid level: integration
- Environment: env-ci
- Owner: SDET
- Reviewer role: QA-Peer
- Risk: high
- Smoke subset: True
- Tags: happy, pdp

#### Expected Assertions
- [state] the page renders name, SKU, description, all images, current price, compare-at price, in-stock indicator, average rating 4.2, and review count 11
- [state] the add-to-cart button is enabled

### TC-EPIC-CATALOG-2-002 — Paused product PDP shows unavailable message instead of add-to-cart

- Scenario type: error
- Test type: integration
- Pyramid level: integration
- Environment: env-ci
- Owner: SDET
- Reviewer role: QA-Peer
- Risk: high
- Smoke subset: False
- Tags: error, i18n-th, pdp, unavailable
- Depends on: TC-EPIC-CATALOG-2-001

#### Expected Assertions
- [state] the page renders product metadata for reference
- [state] the add-to-cart button is replaced by an unavailable-message block with the Thai unavailable copy
- [state] no add-to-cart request can be sent from this view

### TC-EPIC-CATALOG-2-003 — Soft-deleted product reached via direct URL returns 404 or unavailable view

- Scenario type: edge_case
- Test type: integration
- Pyramid level: integration
- Environment: env-ci
- Owner: SDET
- Reviewer role: QA-Peer
- Risk: medium
- Smoke subset: False
- Tags: edge_case, pdp, soft-delete
- Depends on: TC-EPIC-CATALOG-2-001

#### Expected Assertions
- [state] the response is HTTP 404 OR an unavailable view with the same Thai unavailable copy
- [state] no add-to-cart affordance is rendered

### TC-EPIC-CATALOG-2-004 — PDP response never surfaces admin-only product attributes regardless of caller (banking_grade authn_authz)

- Scenario type: banking_grade_authz
- Test type: security
- Pyramid level: integration
- Environment: env-ci
- Owner: Security-Tester
- Reviewer role: Security-Reviewer
- Risk: high
- Smoke subset: False
- Tags: authz, banking_grade_authn_authz, pdp
- Depends on: TC-EPIC-CATALOG-2-001

#### Expected Assertions
- [state] PDP response omits admin-only attributes such as cost price and supplier notes for anonymous callers
- [authz] PDP response payload is identical (no extra fields) when called with arbitrary or forged caller headers

