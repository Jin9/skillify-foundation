---
template_version: 0.1.0
story_id: STORY_PAYMENT_INTENT_CREATE
epic_id: EPIC_PAYMENT
title: Checkout completion creates a REQUIRES_PAYMENT intent linked to the order
as_a: System
i_want: to create a payment intent on checkout completion
so_that: the customer's payment route has a target to simulate against
acceptance_criteria:
  - given: A successful checkout that produced order O with total T
    when: Server-side payment intent creation runs
    then: |
      Payment intent exists with status REQUIRES_PAYMENT, amount=T, and is linked to
      O (PAY-001).
priority: Must
size: S
edge_cases:
  - case: Checkout fails (e.g. INSUFFICIENT_STOCK)
    expected: No payment intent is created
  - case: Two payment intents created for the same order
    expected: Forbidden by design — one order, one intent in MVP; Tech-Designer enforces uniqueness
requirement_refs: [PAY-001]
definition_of_ready: [default]
definition_of_done: [default]
change_log:
  - date: 2026-05-08
    author: BA (Claude Opus 4.7 1M, dry-run #2)
    change: Created from requirement §8.10 PAY-001.
---

# STORY_PAYMENT_INTENT_CREATE — Checkout completion creates a REQUIRES_PAYMENT intent linked to the order

## User narrative

As a **System**, I want **to create a payment intent on checkout completion** so that **the customer's payment route has a target to simulate against**.

## Why this story exists in EPIC_PAYMENT

The payment intent is what the simulate endpoints act on. PAY-001 mandates REQUIRES_PAYMENT as the initial status; this is the gate to PAID / PAYMENT_FAILED / PAYMENT_EXPIRED.

## Edge cases

| Case | Expected |
|---|---|
| Checkout fails | No payment intent created |
| Two intents per order | Disallowed by uniqueness |

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | BA (Claude Opus 4.7 1M, dry-run #2) | Created from requirement §8.10 PAY-001. |
