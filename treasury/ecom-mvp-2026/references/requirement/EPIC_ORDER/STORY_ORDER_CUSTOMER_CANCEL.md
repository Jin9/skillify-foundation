---
template_version: 0.1.0
story_id: STORY_ORDER_CUSTOMER_CANCEL
epic_id: EPIC_ORDER
title: Customer cancels own order while PENDING_PAYMENT
as_a: Customer
i_want: to cancel my own order while it is still PENDING_PAYMENT
so_that: I can change my mind without paying, and the reserved stock is released
acceptance_criteria:
  - given: An order in PENDING_PAYMENT belonging to the requesting customer
    when: Customer POSTs order cancel
    then: |
      Order moves to CANCELLED, reserved stock is released back to availableQty, and
      a status-history row is written (ORD-004, INV-006).
  - given: An order that has already moved to PAID belonging to the requesting customer
    when: Customer POSTs order cancel
    then: |
      Customer-cancel returns 409 INVALID_ORDER_STATE; the order remains PAID; only
      ADMIN can cancel a PAID order with a reason (ORD-004, ORD-005).
priority: Must
size: M
edge_cases:
  - case: Customer attempts to cancel an order belonging to a different customer
    expected: NOT_FOUND or FORBIDDEN; no state change
  - case: Customer attempts to cancel an order in DELIVERED state
    expected: 409 INVALID_ORDER_STATE; order stays DELIVERED
requirement_refs: [ORD-003, ORD-004, ORD-005, ORD-009, INV-006]
definition_of_ready: [default]
definition_of_done: [default]
change_log:
  - date: 2026-05-08
    author: BA (Claude Opus 4.7 1M, dry-run #2)
    change: Created from requirement §8.9 ORD-004 + §9.2 + §9.3.
---

# STORY_ORDER_CUSTOMER_CANCEL — Customer cancels own order while PENDING_PAYMENT

## User narrative

As a **Customer**, I want **to cancel my own order while it is still PENDING_PAYMENT** so that **I can change my mind without paying, and the reserved stock is released**.

## Why this story exists in EPIC_ORDER

This is the only customer-controlled state transition. After PAID, customer self-service cancel ends; only admin can cancel further.

## Edge cases

| Case | Expected |
|---|---|
| Cancel another customer's order | NOT_FOUND or FORBIDDEN |
| Cancel a DELIVERED order | 409 INVALID_ORDER_STATE |

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | BA (Claude Opus 4.7 1M, dry-run #2) | Created from requirement §8.9 ORD-004. |
