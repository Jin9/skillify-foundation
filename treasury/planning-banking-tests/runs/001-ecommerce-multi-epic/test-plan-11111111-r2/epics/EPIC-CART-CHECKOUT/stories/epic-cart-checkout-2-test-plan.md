---
artifact_type: qa-story-plan
blocking_oqs: []
coverage_status: complete
epic_id: EPIC-CART-CHECKOUT
story_id: EPIC-CART-CHECKOUT-2
test_case_count: 7
test_plan_id: TPLAN-EPIC-SHOPPILOT-MVP
---

# Story Test Plan — EPIC-CART-CHECKOUT-2

## Test Case Roster

| ID | Scenario | Type | Pyramid | Owner | Risk |
|---|---|---|---|---|---|
| TC-EPIC-CART-CHECKOUT-2-001 | Valid fixed-THB coupon applies to subtotal | functional | integration | SDET | high |
| TC-EPIC-CART-CHECKOUT-2-002 | Percent coupon with cap applied correctly | functional | integration | SDET | high |
| TC-EPIC-CART-CHECKOUT-2-003 | Expired coupon rejected at submit | functional | integration | SDET | high |
| TC-EPIC-CART-CHECKOUT-2-004 | Discount cannot make total negative | functional | integration | SDET | high |
| TC-EPIC-CART-CHECKOUT-2-005 | Idempotent coupon-apply against double-click | functional | integration | SDET | high |
| TC-EPIC-CART-CHECKOUT-2-006 | Coupon redemption audit event | functional | integration | SDET | high |
| TC-EPIC-CART-CHECKOUT-2-007 | Coupon quota reversibility on order cancel pre-payment | functional | integration | SDET | high |

## Test Case Details

### TC-EPIC-CART-CHECKOUT-2-001 — Valid fixed-THB coupon applies to subtotal

- Scenario type: happy
- Test type: functional
- Pyramid level: integration
- Environment: ENV-CART-integration
- Owner: SDET
- Reviewer role: QA-Lead
- Risk: high
- Smoke subset: True
- Tags: coupon-fixed, happy

#### Expected Assertions
- [computation] discount line == -100 THB exactly
- [computation] grand_total == 560 (subtotal 600 - discount 100 + shipping 60 per shipping rule at subtotal<1500)
- [server_recompute] server-side total computation is the authoritative value; client-supplied total is ignored

### TC-EPIC-CART-CHECKOUT-2-002 — Percent coupon with cap applied correctly

- Scenario type: happy
- Test type: functional
- Pyramid level: integration
- Environment: ENV-CART-integration
- Owner: SDET
- Reviewer role: QA-Lead
- Risk: high
- Smoke subset: False
- Tags: cap, coupon-percent, happy

#### Expected Assertions
- [computation] raw 10% would be 400 but cap clamps discount to exactly 300
- [computation] grand_total == 3700 with shipping=0 (subtotal>=1500 free-shipping threshold)

### TC-EPIC-CART-CHECKOUT-2-003 — Expired coupon rejected at submit

- Scenario type: error
- Test type: functional
- Pyramid level: integration
- Environment: ENV-CART-integration
- Owner: SDET
- Reviewer role: QA-Lead
- Risk: high
- Smoke subset: False
- Tags: coupon-revalidation, error

#### Expected Assertions
- [error_response] response.reason == 'coupon_expired'
- [negative] no order row is created in DB
- [ui_state] cart re-renders with coupon row removed

### TC-EPIC-CART-CHECKOUT-2-004 — Discount cannot make total negative

- Scenario type: edge_case
- Test type: functional
- Pyramid level: integration
- Environment: ENV-CART-integration
- Owner: SDET
- Reviewer role: QA-Lead
- Risk: high
- Smoke subset: False
- Tags: clamp, edge_case, non-negative-total

#### Expected Assertions
- [computation] discount is clamped to subtotal value (80), not the coupon face value (100)
- [computation] grand_total == 0 + shipping; never negative
- [negative] no refund/credit obligation row is created

### TC-EPIC-CART-CHECKOUT-2-005 — Idempotent coupon-apply against double-click

- Scenario type: banking_grade_idempotency
- Test type: functional
- Pyramid level: integration
- Environment: ENV-CART-integration
- Owner: SDET
- Reviewer role: Payment-Lead
- Risk: high
- Smoke subset: True
- Tags: banking_grade_idempotency, coupon, replay-safety

#### Expected Assertions
- [state] cart has exactly one coupon row after replay (not two)
- [state] coupon.quota_remaining decremented by 1 (not 2) across the two requests
- [response_equality] replay returns the original cart-state body

### TC-EPIC-CART-CHECKOUT-2-006 — Coupon redemption audit event

- Scenario type: banking_grade_audit
- Test type: functional
- Pyramid level: integration
- Environment: ENV-CART-integration
- Owner: SDET
- Reviewer role: Compliance
- Risk: high
- Smoke subset: False
- Tags: banking_grade_audit, coupon

#### Expected Assertions
- [audit_event] audit_sink contains exactly one event with event='coupon.applied', actor='customer:C1', ts present in ISO-8601, coupon_code='WELCOME100', amount_applied=100, before=null, after.coupon_code='WELCOME100', reason='coupon_apply'
- [negative] audit payload does NOT include any password, token, card_number, or cvv field (asserted by allow-list)

### TC-EPIC-CART-CHECKOUT-2-007 — Coupon quota reversibility on order cancel pre-payment

- Scenario type: banking_grade_reversibility
- Test type: functional
- Pyramid level: integration
- Environment: ENV-CART-integration
- Owner: SDET
- Reviewer role: Payment-Lead
- Risk: high
- Smoke subset: False
- Tags: banking_grade_reversibility, coupon-quota

#### Expected Assertions
- [state] after order cancel pre-payment, coupon.quota_remaining is restored to pre-apply value
- [audit_event] audit event 'coupon.quota_reversed' emitted with reason='order_cancelled_pre_payment'

