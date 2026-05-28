---
template_version: 0.1.0
story_id: STORY_CART_REMOVE
epic_id: EPIC_CART
title: Customer removes a cart line
as_a: Customer
i_want: to remove a cart line
so_that: items I no longer want do not affect my checkout total
acceptance_criteria:
  - given: An authenticated CUSTOMER with a cart line for product P
    when: Customer POSTs cart remove on that line
    then: |
      Response is 200 with code=SUCCESS, the cart no longer contains a line for P, and
      the cart subtotal is recalculated from the remaining lines (CART-003, CART-004).
  - given: An authenticated CUSTOMER with a cart line for product P that has become INACTIVE
    when: Customer POSTs cart remove on that line
    then: |
      The line is removed regardless of product status; remove must work even for
      INACTIVE/DELETED items so the customer can clear them and proceed to checkout
      (CART-003, CART-005).
priority: Must
size: S
edge_cases:
  - case: Remove a cart line that does not belong to the requesting customer
    expected: NOT_FOUND or FORBIDDEN; no state change
  - case: Remove a cart line that has already been removed
    expected: NOT_FOUND or 200 idempotent (Tech-Designer chooses); no state change either way
requirement_refs: [CART-003, CART-004, CART-005]
definition_of_ready: [default]
definition_of_done: [default]
change_log:
  - date: 2026-05-08
    author: BA (Claude Opus 4.7 1M, dry-run #2)
    change: Created from requirement §8.6 CART-003 + CART-005; "remove must work for INACTIVE items" added to unblock CART-005's "non-checkoutable" requirement.
---

# STORY_CART_REMOVE — Customer removes a cart line

## User narrative

As a **Customer**, I want **to remove a cart line** so that **items I no longer want do not affect my checkout total**.

## Why this story exists in EPIC_CART

The cart tab's "remove" action (Frontend Spec §3) is built on this endpoint. CART-005 says inactive items must be flagged non-checkoutable; remove is the customer's escape hatch for those items.

## Edge cases

| Case | Expected |
|---|---|
| Remove not owned | NOT_FOUND or FORBIDDEN |
| Already-removed line | NOT_FOUND or idempotent 200 |

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | BA (Claude Opus 4.7 1M, dry-run #2) | Created from requirement §8.6 CART-003 + CART-005. |
