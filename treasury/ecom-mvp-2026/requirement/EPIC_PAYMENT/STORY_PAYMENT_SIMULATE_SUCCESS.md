---
template_version: 0.1.0
story_id: STORY_PAYMENT_SIMULATE_SUCCESS
epic_id: EPIC_PAYMENT
title: Customer simulates payment success; order → PAID, stock reserved → sold; idempotent
as_a: Customer
i_want: to simulate a successful payment on the mock payment route
so_that: my order moves to PAID and the reserved stock is converted to sold
acceptance_criteria:
  - given: A payment intent in REQUIRES_PAYMENT for order O
    when: The same payment-success callback is delivered twice
    then: |
      Order O ends in PAID exactly once, soldQty increased exactly once, and the
      second callback returns 200 with no additional state change (PAY-002, PAY-005).
  - given: Any payment intent that resolves
    when: It transitions to a terminal status
    then: |
      The persisted payment row contains mockPaymentRef, providerStatus, and paidAt
      (PAY-007).
priority: Must
size: M
edge_cases:
  - case: Payment-success callback for an order in CANCELLED state
    expected: 409 INVALID_ORDER_STATE per ORD-003; order remains CANCELLED; no stock conversion
  - case: Payment-success callback for an unknown paymentIntentId
    expected: NOT_FOUND; no state change
requirement_refs: [PAY-002, PAY-005, PAY-007, INV-005]
definition_of_ready: [default]
definition_of_done: [default]
change_log:
  - date: 2026-05-08
    author: BA (Claude Opus 4.7 1M, dry-run #2)
    change: Created from requirement §8.10 PAY-002 + PAY-005 + PAY-007.
---

# STORY_PAYMENT_SIMULATE_SUCCESS — Customer simulates payment success; order → PAID, stock reserved → sold; idempotent

## User narrative

As a **Customer**, I want **to simulate a successful payment on the mock payment route** so that **my order moves to PAID and the reserved stock is converted to sold**.

## Why this story exists in EPIC_PAYMENT

The payment route's "success" simulate button drives this endpoint. Idempotency is the non-negotiable: a duplicate callback (e.g. retried by the frontend) must not double-apply.

## Edge cases

| Case | Expected |
|---|---|
| Callback for CANCELLED order | 409 INVALID_ORDER_STATE; no state change |
| Callback for unknown intent | NOT_FOUND |

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | BA (Claude Opus 4.7 1M, dry-run #2) | Created from requirement §8.10 PAY-002. |
