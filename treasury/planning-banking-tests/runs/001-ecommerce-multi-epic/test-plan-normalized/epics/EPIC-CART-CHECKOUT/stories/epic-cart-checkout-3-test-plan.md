---
artifact_type: qa-story-plan
blocking_oqs:
  - OQ-CART-checkout-p95-target
  - OQ-CART-reservation-ttl
coverage_status: partial
epic_id: EPIC-CART-CHECKOUT
story_id: EPIC-CART-CHECKOUT-3
test_case_count: 10
test_plan_id: TPLAN-EPIC-SHOPPILOT-MVP
---

# Story Test Plan — EPIC-CART-CHECKOUT-3

## Test Case Roster

| ID | Scenario | Type | Pyramid | Owner | Risk |
|---|---|---|---|---|---|
| TC-EPIC-CART-CHECKOUT-3-001 | Successful checkout reserves stock and creates order in awaiting-payment | e2e | e2e | SDET | critical |
| TC-EPIC-CART-CHECKOUT-3-002 | Race-resolution on last unit — only one customer wins | integration | integration | SDET | critical |
| TC-EPIC-CART-CHECKOUT-3-003 | Idempotent submit against double-click produces one order | integration | integration | SDET | critical |
| TC-EPIC-CART-CHECKOUT-3-004 | Order creation emits audit event with required payload | functional | integration | SDET | critical |
| TC-EPIC-CART-CHECKOUT-3-005 | Address from another customer is rejected | security | integration | SDET | critical |
| TC-EPIC-CART-CHECKOUT-3-006 | Payment failure releases reservation (reversibility on checkout-side) | integration | integration | SDET | critical |
| TC-EPIC-CART-CHECKOUT-3-007 | Expired reservation TTL sweep releases stock | integration | integration | SDET | high |
| TC-EPIC-CART-CHECKOUT-3-008 | Concurrent reservation on multi-unit (5 customers vs stock=3) | integration | integration | SDET | high |
| TC-EPIC-CART-CHECKOUT-3-009 | Server re-validates coupon at submit (consistency with story-2) | integration | integration | SDET | high |
| TC-EPIC-CART-CHECKOUT-3-010 | Order number uniqueness and unguessability | functional | integration | SDET | medium |

## Test Case Details

### TC-EPIC-CART-CHECKOUT-3-001 — Successful checkout reserves stock and creates order in awaiting-payment

- Scenario type: happy
- Test type: e2e
- Pyramid level: e2e
- Environment: ENV-CHK-e2e
- Owner: SDET
- Reviewer role: Payment-Lead
- Risk: critical
- Smoke subset: True
- Tags: checkout-submit, happy

#### Expected Assertions
- [state] P1.stock_available transitions 5 -> 3 (delta -2)
- [state] P1.stock_reserved increases by 2
- [state] order row created with status='awaiting-payment' and grand_total matching server recompute
- [format] order.number matches /^ORD-\d{8}-\d{6}$/ and is unique across DB
- [snapshot] order_line snapshot captures price, name, sku at submit time (frozen vs current product fields)
- [audit_event] audit_sink contains 'order.created' event with the 7 audit fields (event,actor,ts,before,after,reason,idem_key)

### TC-EPIC-CART-CHECKOUT-3-002 — Race-resolution on last unit — only one customer wins

- Scenario type: edge_case
- Test type: integration
- Pyramid level: integration
- Environment: ENV-CHK-race-deterministic
- Owner: SDET
- Reviewer role: Payment-Lead
- Risk: critical
- Smoke subset: False
- Tags: concurrency, edge_case, last-unit, race
- Depends on: TC-EPIC-CART-CHECKOUT-3-001

#### Expected Assertions
- [state] exactly one of {O_C1, O_C2} exists with status='awaiting-payment'; the other does not exist
- [error_response] the losing customer receives a structured error referencing P1 by id; no PII of the winning customer is leaked
- [invariant] P1.stock_available >= 0 at every observed step (no negative reading)
- [invariant] P1.stock_reserved + P1.stock_available + P1.stock_sold == initial total across the race window

### TC-EPIC-CART-CHECKOUT-3-003 — Idempotent submit against double-click produces one order

- Scenario type: banking_grade_idempotency
- Test type: integration
- Pyramid level: integration
- Environment: ENV-CHK-e2e
- Owner: SDET
- Reviewer role: Payment-Lead
- Risk: critical
- Smoke subset: True
- Tags: banking_grade_idempotency, checkout, replay-safety

#### Expected Assertions
- [state] exactly one order row exists for idempotency_key='ik-chk-001'; second submit returns existing order_id without DB insert
- [state] P1.stock_reserved equals the single-reservation value (not doubled) after replay
- [response_equality] second submit response body equals first response (order_id, grand_total, status all identical)
- [audit_event] only one 'order.created' audit event is emitted across the two calls (no duplicate)

### TC-EPIC-CART-CHECKOUT-3-004 — Order creation emits audit event with required payload

- Scenario type: banking_grade_audit
- Test type: functional
- Pyramid level: integration
- Environment: ENV-CHK-e2e
- Owner: SDET
- Reviewer role: Compliance
- Risk: critical
- Smoke subset: False
- Tags: banking_grade_audit, order-created

#### Expected Assertions
- [audit_event] audit_sink contains exactly one event with 7 fields enumerated: (1) event=='order.created', (2) actor=='customer:C1', (3) ts is valid ISO-8601 within +/-2s of submit time, (4) before==null, (5) after=={order_id:<uuid>, grand_total:<int>}, (6) reason=='checkout_submit', (7) idem_key matches request header (AP-Q12)
- [negative] audit payload does NOT contain any of: password, token, card_number, pan, cvv, full_address_line, phone_number (allow-list assertion)
- [ordering] audit event is durably written before the API response is returned (write-before-respond invariant)

### TC-EPIC-CART-CHECKOUT-3-005 — Address from another customer is rejected

- Scenario type: banking_grade_authz
- Test type: security
- Pyramid level: integration
- Environment: ENV-CHK-e2e
- Owner: SDET
- Reviewer role: Payment-Lead
- Risk: critical
- Smoke subset: True
- Tags: banking_grade_authz, checkout, ownership

#### Expected Assertions
- [authz] HTTP 403 (or domain error 'address_not_owned') returned when C1 references address A owned by C2
- [negative] no order row created; no stock reservation taken; no notification sent
- [audit_event] audit event 'authz_denied' emitted with actor='customer:C1', resource='address:A', reason='cross_owner_access'

### TC-EPIC-CART-CHECKOUT-3-006 — Payment failure releases reservation (reversibility on checkout-side)

- Scenario type: banking_grade_reversibility
- Test type: integration
- Pyramid level: integration
- Environment: ENV-CHK-e2e
- Owner: SDET
- Reviewer role: Payment-Lead
- Risk: critical
- Smoke subset: False
- Tags: banking_grade_reversibility, release-on-fail
- Depends on: TC-EPIC-CART-CHECKOUT-3-001

#### Expected Assertions
- [state] on payment failure signal, stock_reserved decremented and stock_available restored to pre-checkout value (no leak)
- [state] order transitions awaiting-payment -> payment-failed; no order row leaks in awaiting-payment indefinitely
- [audit_event] audit event 'order.payment_failed' emitted with before.status='awaiting-payment', after.status='payment-failed'
- [invariant] stock_available + stock_reserved + stock_sold conservation holds across the reversal

### TC-EPIC-CART-CHECKOUT-3-007 — Expired reservation TTL sweep releases stock

- Scenario type: edge_case
- Test type: integration
- Pyramid level: integration
- Environment: ENV-CHK-race-deterministic
- Owner: SDET
- Reviewer role: Payment-Lead
- Risk: high
- Smoke subset: False
- Tags: edge_case, reservation, ttl-expiry
- Depends on: TC-EPIC-CART-CHECKOUT-3-001

#### Expected Assertions
- [state] after advancing frozen clock past reservation TTL and running sweep, order moves to 'payment-timeout' and stock_reserved decrements by 2
- [state] stock_available restored to original value
- [audit_event] audit event 'order.payment_timeout' emitted with reason='reservation_ttl_expired'

### TC-EPIC-CART-CHECKOUT-3-008 — Concurrent reservation on multi-unit (5 customers vs stock=3)

- Scenario type: edge_case
- Test type: integration
- Pyramid level: integration
- Environment: ENV-CHK-race-deterministic
- Owner: SDET
- Reviewer role: Payment-Lead
- Risk: high
- Smoke subset: False
- Tags: edge_case, multi-unit, race

#### Expected Assertions
- [invariant] exactly 3 orders created in awaiting-payment; remaining 2 customers receive structured 'out_of_stock_at_confirm' error
- [invariant] P1.stock_available reaches 0 (not negative) at end of race
- [invariant] no customer-PII fields leak across error responses

### TC-EPIC-CART-CHECKOUT-3-009 — Server re-validates coupon at submit (consistency with story-2)

- Scenario type: error
- Test type: integration
- Pyramid level: integration
- Environment: ENV-CHK-e2e
- Owner: SDET
- Reviewer role: QA-Lead
- Risk: high
- Smoke subset: False
- Tags: coupon, error, server-revalidate

#### Expected Assertions
- [error_response] submit fails with reason='coupon_expired' or 'coupon_invalid_at_confirm'
- [negative] no order row created; no stock reservation taken

### TC-EPIC-CART-CHECKOUT-3-010 — Order number uniqueness and unguessability

- Scenario type: boundary
- Test type: functional
- Pyramid level: integration
- Environment: ENV-CHK-e2e
- Owner: SDET
- Reviewer role: QA-Lead
- Risk: medium
- Smoke subset: False
- Tags: boundary, id-generation

#### Expected Assertions
- [format] 100 sequential order_numbers all match /^ORD-\d{8}-\d{6}$/ and are pairwise distinct
- [negative] order_numbers are NOT a tight sequential counter visible to client (gaps allowed; not trivially guessable - assert that consecutive order_numbers differ by >1 in at least 5% of samples, or are random within the 6-digit suffix space)

