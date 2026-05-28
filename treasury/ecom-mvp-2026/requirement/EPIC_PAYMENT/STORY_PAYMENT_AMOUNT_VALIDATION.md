---
template_version: 0.1.0
story_id: STORY_PAYMENT_AMOUNT_VALIDATION
epic_id: EPIC_PAYMENT
title: Callback with mismatched amount returns PAYMENT_AMOUNT_MISMATCH; order unchanged
as_a: System
i_want: to reject payment callbacks whose amount does not equal the intent's amount
so_that: customers cannot pay less than the order total
acceptance_criteria:
  - given: A payment intent for order O with amount T
    when: A payment-success callback arrives with amount != T
    then: |
      Response is 409 PAYMENT_AMOUNT_MISMATCH, order O remains PENDING_PAYMENT, and
      stock reservation is unchanged (PAY-006).
  - given: A payment intent for order O with amount T
    when: A payment-success callback arrives with amount == T (exactly)
    then: |
      Response is 200 with code=SUCCESS; order O moves to PAID and stock converts
      reserved → sold (PAY-002, PAY-006).
priority: Must
size: S
edge_cases:
  - case: Callback amount differs by 1 satang due to rounding
    expected: 409 PAYMENT_AMOUNT_MISMATCH; the comparison is exact in MVP, no tolerance
  - case: Callback amount is zero
    expected: 409 PAYMENT_AMOUNT_MISMATCH; order remains PENDING_PAYMENT
requirement_refs: [PAY-002, PAY-006]
definition_of_ready: [default]
definition_of_done: [default]
change_log:
  - date: 2026-05-08
    author: BA (Claude Opus 4.7 1M, dry-run #2)
    change: Created from requirement §8.10 PAY-006.
---

# STORY_PAYMENT_AMOUNT_VALIDATION — Callback with mismatched amount returns PAYMENT_AMOUNT_MISMATCH; order unchanged

## User narrative

As a **System**, I want **to reject payment callbacks whose amount does not equal the intent's amount** so that **customers cannot pay less than the order total**.

## Why this story exists in EPIC_PAYMENT

PAY-006 is the only mechanism in MVP that prevents a 1-baht payment from settling a 2,480-baht order. There is no other guard — the payment service is the integrity boundary.

## Edge cases

| Case | Expected |
|---|---|
| Off-by-1-satang | 409 MISMATCH; exact comparison only |
| Zero amount | 409 MISMATCH; order unchanged |

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | BA (Claude Opus 4.7 1M, dry-run #2) | Created from requirement §8.10 PAY-006. |
