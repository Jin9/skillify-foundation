---
artifact_type: qa-story-plan
blocking_oqs: []
coverage_status: complete
epic_id: EPIC-CART-CHECKOUT
story_id: EPIC-CART-CHECKOUT-4
test_case_count: 8
test_plan_id: TPLAN-EPIC-SHOPPILOT-MVP
---

# Story Test Plan — EPIC-CART-CHECKOUT-4

## Test Case Roster

| ID | Scenario | Type | Pyramid | Owner | Risk |
|---|---|---|---|---|---|
| TC-EPIC-CART-CHECKOUT-4-001 | Successful payment moves order to paid and converts stock to sold | e2e | e2e | SDET | critical |
| TC-EPIC-CART-CHECKOUT-4-002 | Failed payment releases stock and moves order to payment-failed | integration | integration | SDET | critical |
| TC-EPIC-CART-CHECKOUT-4-003 | Timed-out payment releases stock and moves order to payment-timeout | integration | integration | SDET | critical |
| TC-EPIC-CART-CHECKOUT-4-004 | Duplicate webhook callback does NOT double-effect | integration | integration | SDET | critical |
| TC-EPIC-CART-CHECKOUT-4-005 | Price-mismatch payment is rejected | integration | integration | SDET | critical |
| TC-EPIC-CART-CHECKOUT-4-006 | Payment audit event field enumeration (cross-outcome) | functional | integration | SDET | critical |
| TC-EPIC-CART-CHECKOUT-4-007 | Payment reversibility — failure reverses reservation fully | integration | integration | SDET | critical |
| TC-EPIC-CART-CHECKOUT-4-008 | Webhook origin authorization (mock provider signature) | security | integration | SDET | critical |

## Test Case Details

### TC-EPIC-CART-CHECKOUT-4-001 — Successful payment moves order to paid and converts stock to sold

- Scenario type: happy
- Test type: e2e
- Pyramid level: e2e
- Environment: ENV-CHK-e2e
- Owner: SDET
- Reviewer role: Payment-Lead
- Risk: critical
- Smoke subset: True
- Tags: happy, payment-success

#### Expected Assertions
- [state] O1.status transitions awaiting-payment -> paid
- [state] P1.stock_reserved -2, P1.stock_sold +2
- [state] shipment_record created with order_id=O1, status='to-pack'
- [state] notification_log row 'payment.succeeded' inserted
- [audit_event] audit event 'order.payment_succeeded' emitted with 7 fields including before.status='awaiting-payment', after.status='paid'

### TC-EPIC-CART-CHECKOUT-4-002 — Failed payment releases stock and moves order to payment-failed

- Scenario type: error
- Test type: integration
- Pyramid level: integration
- Environment: ENV-CHK-e2e
- Owner: SDET
- Reviewer role: Payment-Lead
- Risk: critical
- Smoke subset: True
- Tags: error, payment-failed

#### Expected Assertions
- [state] O1.status transitions awaiting-payment -> payment-failed
- [state] P1.stock_reserved -2, P1.stock_available +2
- [negative] no shipment_record exists for O1
- [audit_event] audit event 'order.payment_failed' emitted

### TC-EPIC-CART-CHECKOUT-4-003 — Timed-out payment releases stock and moves order to payment-timeout

- Scenario type: edge_case
- Test type: integration
- Pyramid level: integration
- Environment: ENV-CHK-race-deterministic
- Owner: SDET
- Reviewer role: Payment-Lead
- Risk: critical
- Smoke subset: False
- Tags: edge_case, payment-timeout

#### Expected Assertions
- [state] O1.status transitions awaiting-payment -> payment-timeout after sweep runs past TTL
- [state] P1.stock_reserved -2, P1.stock_available +2
- [audit_event] audit event 'order.payment_timeout' emitted

### TC-EPIC-CART-CHECKOUT-4-004 — Duplicate webhook callback does NOT double-effect

- Scenario type: banking_grade_idempotency
- Test type: integration
- Pyramid level: integration
- Environment: ENV-CHK-e2e
- Owner: SDET
- Reviewer role: Payment-Lead
- Risk: critical
- Smoke subset: True
- Tags: banking_grade_idempotency, replay-safety, webhook

#### Expected Assertions
- [state] after replay of provider_event_id='evt-001', O1.status remains 'paid' (no extra transition)
- [state] P1.stock_sold unchanged by replay (no second +2)
- [audit_event] audit_sink contains exactly one 'order.payment_succeeded' event for O1 (not two)
- [response_equality] second webhook response equals first (same outcome body)

### TC-EPIC-CART-CHECKOUT-4-005 — Price-mismatch payment is rejected

- Scenario type: error
- Test type: integration
- Pyramid level: integration
- Environment: ENV-CHK-e2e
- Owner: SDET
- Reviewer role: Payment-Lead
- Risk: critical
- Smoke subset: False
- Tags: error, price-mismatch

#### Expected Assertions
- [state] O1 does NOT transition to 'paid'; remains in 'awaiting-payment' or moves to 'mismatch' sub-state
- [audit_event] audit event 'order.payment_mismatch_rejected' emitted with payload {expected:560, received:500}
- [negative] stock_sold not incremented; shipment_record not created

### TC-EPIC-CART-CHECKOUT-4-006 — Payment audit event field enumeration (cross-outcome)

- Scenario type: banking_grade_audit
- Test type: functional
- Pyramid level: integration
- Environment: ENV-CHK-e2e
- Owner: SDET
- Reviewer role: Compliance
- Risk: critical
- Smoke subset: False
- Tags: banking_grade_audit, payment

#### Expected Assertions
- [audit_event] for each outcome in {success, failure, timeout}, audit_sink emits one event whose payload contains exactly the 7 fields per AP-Q12: (1) event, (2) actor, (3) ts (ISO-8601), (4) before (prior status), (5) after (new status), (6) reason ('payment_signal'|'reservation_ttl_expired'), (7) idem_key (provider_event_id)
- [negative] no audit payload contains PAN, CVV, full card number, or password (PCI-out-of-scope allow-list)

### TC-EPIC-CART-CHECKOUT-4-007 — Payment reversibility — failure reverses reservation fully

- Scenario type: banking_grade_reversibility
- Test type: integration
- Pyramid level: integration
- Environment: ENV-CHK-e2e
- Owner: SDET
- Reviewer role: Payment-Lead
- Risk: critical
- Smoke subset: False
- Tags: banking_grade_reversibility, payment-failure
- Depends on: TC-EPIC-CART-CHECKOUT-4-002

#### Expected Assertions
- [invariant] stock_available + stock_reserved + stock_sold equals pre-checkout total after reversal
- [state] no order leaks: every order whose latest payment outcome is 'failure' is in terminal 'payment-failed' (not stuck in awaiting-payment)

### TC-EPIC-CART-CHECKOUT-4-008 — Webhook origin authorization (mock provider signature)

- Scenario type: banking_grade_authz
- Test type: security
- Pyramid level: integration
- Environment: ENV-CHK-e2e
- Owner: SDET
- Reviewer role: Payment-Lead
- Risk: critical
- Smoke subset: False
- Tags: banking_grade_authz, webhook-origin

#### Expected Assertions
- [authz] callback from an origin other than the configured mock PSP is rejected with HTTP 401 or 403
- [negative] no order state change occurs on rejected callback
- [audit_event] audit event 'webhook_authz_denied' emitted with caller_origin field

