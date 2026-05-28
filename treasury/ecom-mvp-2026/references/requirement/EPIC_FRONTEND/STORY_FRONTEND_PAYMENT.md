---
template_version: 0.1.0
story_id: STORY_FRONTEND_PAYMENT
epic_id: EPIC_FRONTEND
title: payment route exposes mock-QR + 3 simulate buttons; navigates by service-returned status
as_a: Customer
i_want: a Thai-language payment route with a mock QR display and three simulate buttons (success/failed/timeout)
so_that: I can complete the demo payment loop end-to-end
acceptance_criteria:
  - given: An authenticated customer on the payment route after a successful checkout
    when: The customer taps simulate-success, simulate-failed, or simulate-timeout
    then: |
      The app calls the payment service for the chosen outcome, then routes to
      orderSuccess on PAID, or to orderResult on PAYMENT_FAILED / PAYMENT_EXPIRED,
      surfacing the order number and a deep link into orderDetail (Frontend Spec §5,
      PAY-002, PAY-003, PAY-004).
  - given: An authenticated customer on the payment route
    when: The customer taps a simulate button
    then: |
      A short processing indicator (Frontend Spec §5: ~1.4s) is shown while the
      payment service resolves; navigation depends on the resulting status returned
      by the service, NOT on a local optimistic guess.
priority: Must
size: M
edge_cases:
  - case: Customer closes the browser on the payment route before tapping any simulate button
    expected: The order remains PENDING_PAYMENT until the reservation TTL elapses, then transitions to PAYMENT_EXPIRED with stock released; the orders tab still surfaces the order with the resulting status (PAY-004, INV-004)
  - case: Customer taps simulate-success twice quickly
    expected: The payment service is idempotent on (paymentIntentId, providerStatus); the order ends in PAID exactly once and the frontend routes to orderSuccess once (PAY-005)
requirement_refs: [PAY-002, PAY-003, PAY-004, PAY-005, "Frontend Spec §3", "Frontend Spec §5"]
definition_of_ready: [default]
definition_of_done: [default]
change_log:
  - date: 2026-05-08
    author: BA (Claude Opus 4.7 1M, dry-run #2)
    change: Created NEW for dry-run #2 from Frontend Spec §3 payment row + §5 mock payment.
---

# STORY_FRONTEND_PAYMENT — payment route exposes mock-QR + 3 simulate buttons; navigates by service-returned status

## User narrative

As a **Customer**, I want **a Thai-language payment route with a mock QR display and three simulate buttons (success/failed/timeout)** so that **I can complete the demo payment loop end-to-end**.

## Why this story exists in EPIC_FRONTEND

This is the only route where the customer drives a state-machine transition that affects inventory. The non-negotiable contract: navigation is driven by the service's returned status, never by a frontend guess.

## Edge cases

| Case | Expected |
|---|---|
| Browser closed before tap | Order eventually PAYMENT_EXPIRED via sweeper |
| Double-tap simulate-success | Idempotent; one transition; one navigation |

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | BA (Claude Opus 4.7 1M, dry-run #2) | Created NEW for dry-run #2 from Frontend Spec §3 + §5. |
