---
artifact_type: qa-story-plan
blocking_oqs: []
coverage_status: complete
epic_id: EPIC-CART-CHECKOUT
story_id: EPIC-CART-CHECKOUT-1
test_case_count: 6
test_plan_id: TPLAN-EPIC-SHOPPILOT-MVP
---

# Story Test Plan — EPIC-CART-CHECKOUT-1

## Test Case Roster

| ID | Scenario | Type | Pyramid | Owner | Risk |
|---|---|---|---|---|---|
| TC-EPIC-CART-CHECKOUT-1-001 | Add active product to empty cart | functional | integration | SDET | medium |
| TC-EPIC-CART-CHECKOUT-1-002 | Re-adding same SKU merges quantities | functional | integration | SDET | medium |
| TC-EPIC-CART-CHECKOUT-1-003 | Reject quantity that exceeds available stock | functional | integration | SDET | high |
| TC-EPIC-CART-CHECKOUT-1-004 | Stale-cart warning when product becomes paused | functional | integration | SDET | high |
| TC-EPIC-CART-CHECKOUT-1-005 | Idempotent add-to-cart against double-click | functional | integration | SDET | high |
| TC-EPIC-CART-CHECKOUT-1-006 | Cart owner authorization | security | integration | SDET | high |

## Test Case Details

### TC-EPIC-CART-CHECKOUT-1-001 — Add active product to empty cart

- Scenario type: happy
- Test type: functional
- Pyramid level: integration
- Environment: ENV-CART-integration
- Owner: SDET
- Reviewer role: QA-Lead
- Risk: medium
- Smoke subset: True
- Tags: cart-add, happy

#### Expected Assertions
- [state] cart contains exactly one line: product=P1, quantity=2
- [computation] cart.subtotal == 2 * P1.price (exact, no rounding drift)

### TC-EPIC-CART-CHECKOUT-1-002 — Re-adding same SKU merges quantities

- Scenario type: happy
- Test type: functional
- Pyramid level: integration
- Environment: ENV-CART-integration
- Owner: SDET
- Reviewer role: QA-Lead
- Risk: medium
- Smoke subset: False
- Tags: cart-merge, happy

#### Expected Assertions
- [state] cart has exactly one line for P1 with quantity=5
- [negative] no second cart line for P1 is created (cart.lines.length == 1)

### TC-EPIC-CART-CHECKOUT-1-003 — Reject quantity that exceeds available stock

- Scenario type: error
- Test type: functional
- Pyramid level: integration
- Environment: ENV-CART-integration
- Owner: SDET
- Reviewer role: QA-Lead
- Risk: high
- Smoke subset: True
- Tags: error, stock-validation

#### Expected Assertions
- [error_response] response.reason == 'quantity_exceeds_stock'
- [state] cart line for P1 is unchanged from pre-call snapshot

### TC-EPIC-CART-CHECKOUT-1-004 — Stale-cart warning when product becomes paused

- Scenario type: edge_case
- Test type: functional
- Pyramid level: integration
- Environment: ENV-CART-integration
- Owner: SDET
- Reviewer role: QA-Lead
- Risk: high
- Smoke subset: False
- Tags: edge_case, stale-cart

#### Expected Assertions
- [ui_state] cart page renders a warning row referencing P1 with availability='unavailable'
- [negative] proceed-to-checkout button is disabled or returns error when P1 is in cart

### TC-EPIC-CART-CHECKOUT-1-005 — Idempotent add-to-cart against double-click

- Scenario type: banking_grade_idempotency
- Test type: functional
- Pyramid level: integration
- Environment: ENV-CART-integration
- Owner: SDET
- Reviewer role: Payment-Lead
- Risk: high
- Smoke subset: True
- Tags: banking_grade_idempotency, replay-safety

#### Expected Assertions
- [state] P1 cart quantity remains 2 after replay (no increment to 4)
- [response_equality] replay response body byte-equals original add response (modulo timestamps)

### TC-EPIC-CART-CHECKOUT-1-006 — Cart owner authorization

- Scenario type: banking_grade_authz
- Test type: security
- Pyramid level: integration
- Environment: ENV-CART-integration
- Owner: SDET
- Reviewer role: QA-Lead
- Risk: high
- Smoke subset: False
- Tags: banking_grade_authz, ownership

#### Expected Assertions
- [authz] C2 GET /cart returns C2's cart only; C1's cart contents not leaked in any field
- [authz] C2 attempting to read C1's cart by id returns 403 (no 404 vs 403 oracle distinction)

