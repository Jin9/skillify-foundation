---
artifact_type: qa-story-plan
blocking_oqs:
  - OQ-12
  - OQ-16
  - OQ-37
  - OQ-49
coverage_status: partial
epic_id: EPIC-CATALOG
story_id: EPIC-CATALOG-1
test_case_count: 5
test_plan_id: TPLAN-EPIC-SHOPPILOT-MVP
---

# Story Test Plan — EPIC-CATALOG-1

## Test Case Roster

| ID | Scenario | Type | Pyramid | Owner | Risk |
|---|---|---|---|---|---|
| TC-EPIC-CATALOG-1-001 | Guest sees only active products in listing | integration | integration | SDET | high |
| TC-EPIC-CATALOG-1-002 | Filter by category narrows results | integration | integration | SDET | medium |
| TC-EPIC-CATALOG-1-003 | Filter by in-stock-only excludes out-of-stock items | integration | integration | SDET | medium |
| TC-EPIC-CATALOG-1-004 | Price range filter rejects invalid bounds | integration | integration | SDET | medium |
| TC-EPIC-CATALOG-1-005 | Public listing endpoint never surfaces admin-only fields or non-active products (banking_grade authn_authz) | security | integration | Security-Tester | high |

## Test Case Details

### TC-EPIC-CATALOG-1-001 — Guest sees only active products in listing

- Scenario type: happy
- Test type: integration
- Pyramid level: integration
- Environment: env-ci
- Owner: SDET
- Reviewer role: QA-Peer
- Risk: high
- Smoke subset: True
- Tags: catalog, happy, listing

#### Expected Assertions
- [state] the listing returns exactly the 12 active products
- [state] no product with status in {draft, paused, deleted} appears in the response
- [state] pagination metadata is present in the response

### TC-EPIC-CATALOG-1-002 — Filter by category narrows results

- Scenario type: happy
- Test type: integration
- Pyramid level: integration
- Environment: env-ci
- Owner: SDET
- Reviewer role: QA-Peer
- Risk: medium
- Smoke subset: False
- Tags: category, filter, happy, listing
- Depends on: TC-EPIC-CATALOG-1-001

#### Expected Assertions
- [state] the listing returns only products whose category_id = Electronics
- [state] the count badge reflects the filtered count
- [state] the URL contains the filter parameter for shareability

### TC-EPIC-CATALOG-1-003 — Filter by in-stock-only excludes out-of-stock items

- Scenario type: happy
- Test type: integration
- Pyramid level: integration
- Environment: env-ci
- Owner: SDET
- Reviewer role: QA-Peer
- Risk: medium
- Smoke subset: False
- Tags: filter, happy, listing, stock
- Depends on: TC-EPIC-CATALOG-1-001

#### Expected Assertions
- [state] the response returns only the 10 products where stock_available > 0
- [state] no product with stock_available = 0 appears

### TC-EPIC-CATALOG-1-004 — Price range filter rejects invalid bounds

- Scenario type: error
- Test type: integration
- Pyramid level: integration
- Environment: env-ci
- Owner: SDET
- Reviewer role: QA-Peer
- Risk: medium
- Smoke subset: False
- Tags: error, filter, listing, validation
- Depends on: TC-EPIC-CATALOG-1-001

#### Expected Assertions
- [state] the listing returns an empty result with an explanation that min must be <= max
- [state] no server error is raised
- [state] the previous filter state is recoverable

### TC-EPIC-CATALOG-1-005 — Public listing endpoint never surfaces admin-only fields or non-active products (banking_grade authn_authz)

- Scenario type: banking_grade_authz
- Test type: security
- Pyramid level: integration
- Environment: env-ci
- Owner: Security-Tester
- Reviewer role: Security-Reviewer
- Risk: high
- Smoke subset: False
- Tags: authz, banking_grade_authn_authz, listing
- Depends on: TC-EPIC-CATALOG-1-001

#### Expected Assertions
- [authz] anonymous guest caller receives only status=active products
- [state] response payload omits admin-only fields (e.g. cost price, supplier notes) regardless of caller
- [state] no product with status in {draft, paused, deleted} appears even when caller crafts arbitrary query parameters

