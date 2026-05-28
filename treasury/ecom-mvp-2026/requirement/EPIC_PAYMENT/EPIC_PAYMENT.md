---
template_version: 0.1.0
epic_id: EPIC_PAYMENT
title: Mock payment intent, simulate-success/failed/timeout, idempotent callbacks
summary: |
  Checkout completion creates a payment intent with status REQUIRES_PAYMENT. The customer
  drives one of three simulate actions (success / failed / timeout) on the payment route
  of the mobile webview. The payment service transitions the owning order to PAID,
  PAYMENT_FAILED, or PAYMENT_EXPIRED and applies the appropriate stock conversion. All
  callbacks are idempotent; amount mismatches return PAYMENT_AMOUNT_MISMATCH without
  transitioning the order.
business_value: Mock payment unblocks the entire commerce loop locally without a real PSP. Idempotent + amount-validated callbacks are what make the mock a credible substitute that real Omise/Stripe integrations can drop into later.
in_scope_stories:
  - STORY_PAYMENT_INTENT_CREATE
  - STORY_PAYMENT_SIMULATE_SUCCESS
  - STORY_PAYMENT_SIMULATE_FAILED_OR_TIMEOUT
  - STORY_PAYMENT_AMOUNT_VALIDATION
out_of_scope:
  - Real payment gateway integration (§4.2 #1) deferred
  - Saved cards / tokenisation — no PSP, nothing to tokenise
  - Refunds (§4.2 #6) deferred
  - Multi-currency (§4.2 #11) deferred — THB only
success_metrics:
  - Duplicate success callback for same paymentIntentId is idempotent: order PAID once, soldQty incremented once (PAY-005)
  - amount mismatch returns 409 PAYMENT_AMOUNT_MISMATCH and does NOT transition the order (PAY-006)
  - Failed/timeout always release reserved stock back to available (PAY-003, PAY-004, INV-006)
  - Every resolved intent persists mockPaymentRef, providerStatus, paidAt (PAY-007)
dependencies:
  - EPIC_CHECKOUT
  - EPIC_ORDER
  - EPIC_INVENTORY
compliance_sensitive: false
change_log:
  - date: 2026-05-08
    author: BA (Claude Opus 4.7 1M, dry-run #2)
    change: Created from requirement §8.10; simulate triggers map to the customer-driven payment route from Frontend Spec §3.
---

# EPIC_PAYMENT — Mock payment intent, simulate-success/failed/timeout, idempotent callbacks

## Why

Real PSPs are out of scope (§4.2 #1). The mock payment service is what makes the rest of the system testable locally: checkout creates a payment intent, the customer's mobile webview drives a simulate-success / simulate-failed / simulate-timeout action on the payment route, and the payment service applies the corresponding side-effects on the order and inventory.

Three behavioral guarantees keep the mock honest enough to swap for a real PSP later without ripping out logic:

1. **Idempotency on (paymentIntentId, providerStatus).** PAY-005: a duplicate success callback must not double-transition the order or double-mutate stock. The Frontend Spec payment route's 1.4-second processing animation does NOT decide state — only the payment service does.

2. **Amount validation.** PAY-006: a callback whose amount does not match the order total returns 409 PAYMENT_AMOUNT_MISMATCH without transitioning the order. Otherwise an attacker could pay 1 baht for a 2480-baht order.

3. **Stock release on non-success.** PAY-003 + PAY-004 + INV-006: failed and timeout always release the reservation. The frontend's orderResult route consumes the resulting status (PAYMENT_FAILED / PAYMENT_EXPIRED) and surfaces the order number plus a deep link into orderDetail.

## Stories in this epic

| Story ID | Title | Priority |
|---|---|---|
| `STORY_PAYMENT_INTENT_CREATE` | Checkout completion creates a REQUIRES_PAYMENT intent linked to the order | Must |
| `STORY_PAYMENT_SIMULATE_SUCCESS` | Customer simulates payment success; order → PAID, stock reserved → sold; idempotent | Must |
| `STORY_PAYMENT_SIMULATE_FAILED_OR_TIMEOUT` | Customer simulates payment failed or timeout; order → PAYMENT_FAILED / PAYMENT_EXPIRED, stock released | Must |
| `STORY_PAYMENT_AMOUNT_VALIDATION` | Callback with mismatched amount returns PAYMENT_AMOUNT_MISMATCH; order unchanged | Must |

## Out-of-scope rationale

- **Real PSP** — §4.2 #1 explicitly defers; the mock service is sized to be drop-in-replaceable, not an actual gateway.
- **Tokenisation / saved cards** — without a PSP there is nothing to tokenise.
- **Refunds** — §4.2 #6 deferred; cancel + restock substitutes for refund in MVP.
- **Multi-currency** — §4.2 #11 deferred; THB only, no FX in MVP.

## Risks

- **Idempotency key choice** — PAY-005 specifies idempotency on (paymentIntentId, providerStatus). Tech-Designer must persist this combination atomically; otherwise a redelivered callback can double-apply.
- **Cross-service write** — payment success transitions the order AND mutates inventory. Failure-tolerance non-functional requires no half-state. Tech-Designer chooses the consistency mechanism (saga, two-phase commit, outbox).
- **Frontend race on simulate** — the Frontend Spec payment route shows a 1.4-second processing animation. It must not navigate to orderSuccess before the payment service confirms the transition. The frontend consumes the *resulting* order status, not a local optimistic guess.
- **Amount-mismatch leak** — PAY-006 should not leak the actual expected amount in the error message; the code is enough.

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | BA (Claude Opus 4.7 1M, dry-run #2) | Created from requirement §8.10; simulate triggers mapped to the customer-driven payment route. |
