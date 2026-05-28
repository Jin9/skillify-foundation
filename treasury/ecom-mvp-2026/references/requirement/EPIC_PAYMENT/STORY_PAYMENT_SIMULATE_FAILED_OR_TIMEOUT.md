---
template_version: 0.1.0
story_id: STORY_PAYMENT_SIMULATE_FAILED_OR_TIMEOUT
epic_id: EPIC_PAYMENT
title: Customer simulates payment failed or timeout; stock released
as_a: Customer
i_want: to simulate a failed or timed-out payment on the mock payment route
so_that: my order ends in PAYMENT_FAILED or PAYMENT_EXPIRED and the reserved stock is released
acceptance_criteria:
  - given: A SKU with availableQty=7 reservedQty=3 soldQty=0 from an active reservation
    when: Payment-failed is simulated for the owning order
    then: |
      availableQty=10 reservedQty=0 soldQty=0 and the order moves to PAYMENT_FAILED
      (PAY-003, INV-006).
  - given: An order in PENDING_PAYMENT and an active reservation
    when: Payment-timeout is simulated by the customer
    then: |
      Order moves to PAYMENT_EXPIRED, reservedQty is released back to availableQty,
      and a status-history row is written (PAY-004, INV-006, ORD-009).
priority: Must
size: M
edge_cases:
  - case: Payment-failed callback delivered twice for the same paymentIntentId
    expected: Idempotent; order ends PAYMENT_FAILED exactly once, stock released exactly once (PAY-005)
  - case: Payment-timeout simulated on an order already in PAID
    expected: 409 INVALID_ORDER_STATE; no state change (ORD-003)
requirement_refs: [PAY-003, PAY-004, PAY-005, INV-006, ORD-009]
definition_of_ready: [default]
definition_of_done: [default]
change_log:
  - date: 2026-05-08
    author: BA (Claude Opus 4.7 1M, dry-run #2)
    change: Created from requirement §8.10 PAY-003 + PAY-004.
---

# STORY_PAYMENT_SIMULATE_FAILED_OR_TIMEOUT — Customer simulates payment failed or timeout; stock released

## User narrative

As a **Customer**, I want **to simulate a failed or timed-out payment on the mock payment route** so that **my order ends in PAYMENT_FAILED or PAYMENT_EXPIRED and the reserved stock is released**.

## Why this story exists in EPIC_PAYMENT

The payment route's "failed" and "timeout" simulate buttons drive these endpoints. Stock release is the most important guarantee — otherwise abandoned checkouts permanently lock inventory.

## Edge cases

| Case | Expected |
|---|---|
| Failed callback delivered twice | Idempotent; one transition, one stock release |
| Timeout on already-PAID order | 409 INVALID_ORDER_STATE |

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | BA (Claude Opus 4.7 1M, dry-run #2) | Created from requirement §8.10 PAY-003 + PAY-004. |
