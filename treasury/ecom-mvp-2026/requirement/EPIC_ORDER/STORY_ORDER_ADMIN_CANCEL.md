---
template_version: 0.1.0
story_id: STORY_ORDER_ADMIN_CANCEL
epic_id: EPIC_ORDER
title: Admin cancels PENDING_PAYMENT or PAID with reason; refused after PACKING
as_a: Admin
i_want: to cancel an order with a reason from PENDING_PAYMENT or PAID
so_that: I can refund stuck or problematic orders before fulfillment begins
acceptance_criteria:
  - given: An order in PAID state and an ADMIN token
    when: Admin POSTs order cancel with reason='out of stock'
    then: |
      Order moves to CANCELLED, the reason is persisted on the cancellation record,
      and stock previously sold is released back to availableQty per MVP policy
      (ORD-005).
  - given: An ADMIN token and an order whose current status is PACKING (or SHIPPED, or DELIVERED)
    when: Admin POSTs order cancel with a reason
    then: |
      admin-cancel from PENDING_PAYMENT or PAID succeeds and emits order.cancelled;
      admin-cancel from PACKING or later returns 409 INVALID_ORDER_STATE and the
      order is unchanged (ORD-005, §9.2, §9.3).
  - given: An ADMIN token and an order in PENDING_PAYMENT
    when: Admin POSTs order cancel with reason='admin override' on the PENDING_PAYMENT order
    then: |
      Order moves to CANCELLED, reserved stock is released back to availableQty, the
      reason is persisted, and a status-history row is written (ORD-005, INV-006).
priority: Should
size: M
edge_cases:
  - case: Admin POSTs cancel without supplying a reason
    expected: 400 VALIDATION_FAILED on the reason field; no state change (ORD-005 requires a reason)
  - case: Customer token attempts admin-cancel
    expected: 403 FORBIDDEN (AUTH-003)
  - case: Two admins attempt to cancel the same PAID order concurrently
    expected: At most one cancellation succeeds; the second returns 409 INVALID_ORDER_STATE
requirement_refs: [ORD-003, ORD-005, ORD-009, AUTH-003, INV-006]
definition_of_ready: [default]
definition_of_done: [default]
change_log:
  - date: 2026-05-08
    author: BA (Claude Opus 4.7 1M, dry-run #2)
    change: Created from requirement §8.9 ORD-005 + §9.2 + §9.3; explicit auto-restock policy on PAID cancel.
---

# STORY_ORDER_ADMIN_CANCEL — Admin cancels PENDING_PAYMENT or PAID with reason; refused after PACKING

## User narrative

As an **Admin**, I want **to cancel an order with a reason from PENDING_PAYMENT or PAID** so that **I can refund stuck or problematic orders before fulfillment begins**.

## Why this story exists in EPIC_ORDER

The MVP policy: admin can cancel until packing starts. After PACKING, the package is being prepared; cancelling there is a fulfillment problem the MVP does not solve. Auto-restock on PAID cancel is locked from runs/01 to keep inventory honest without a refund flow.

## Edge cases

| Case | Expected |
|---|---|
| Missing reason | 400 VALIDATION_FAILED |
| CUSTOMER token | 403 FORBIDDEN |
| Concurrent admin cancels | At most one succeeds |

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | BA (Claude Opus 4.7 1M, dry-run #2) | Created from requirement §8.9 ORD-005. |
