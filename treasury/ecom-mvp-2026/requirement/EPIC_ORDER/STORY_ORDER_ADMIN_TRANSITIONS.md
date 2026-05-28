---
template_version: 0.1.0
story_id: STORY_ORDER_ADMIN_TRANSITIONS
epic_id: EPIC_ORDER
title: Admin moves PAID → PACKING → SHIPPED → DELIVERED
as_a: Admin
i_want: to advance an order through fulfillment statuses
so_that: customers see accurate status and tracking
acceptance_criteria:
  - given: An order in PACKING and an ADMIN token
    when: Admin POSTs order ship without a tracking number
    then: |
      Response is VALIDATION_ERROR and the order remains PACKING (ORD-006).
  - given: An order in DELIVERED state
    when: Admin POSTs order cancel
    then: |
      Response is 409 INVALID_ORDER_STATE and the order remains DELIVERED (ORD-003).
priority: Must
size: M
edge_cases:
  - case: Admin tries to transition order from SHIPPED back to PAID
    expected: 409 INVALID_ORDER_STATE; order stays SHIPPED (ORD-003, §9.3)
  - case: Customer token attempts the admin-transition endpoint
    expected: 403 FORBIDDEN (AUTH-003)
  - case: Admin POSTs SHIP with a tracking number on an order in PAID (skipping PACKING)
    expected: 409 INVALID_ORDER_STATE per §9.2 — only PACKING → SHIPPED is allowed
requirement_refs: [ORD-003, ORD-006, ORD-009, AUTH-003]
definition_of_ready: [default]
definition_of_done: [default]
change_log:
  - date: 2026-05-08
    author: BA (Claude Opus 4.7 1M, dry-run #2)
    change: Created from requirement §8.9 ORD-006 + §9.2 + §9.3.
---

# STORY_ORDER_ADMIN_TRANSITIONS — Admin moves PAID → PACKING → SHIPPED → DELIVERED

## User narrative

As an **Admin**, I want **to advance an order through fulfillment statuses** so that **customers see accurate status and tracking**.

## Why this story exists in EPIC_ORDER

Backend-only this run (no admin UI). The state machine in §9.2 / §9.3 is enforced at this surface. Tracking number on SHIPPED is a hard requirement (ORD-006).

## Edge cases

| Case | Expected |
|---|---|
| SHIPPED → PAID rollback | 409 INVALID_ORDER_STATE |
| CUSTOMER token | 403 FORBIDDEN |
| PAID → SHIPPED skip | 409 INVALID_ORDER_STATE |

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | BA (Claude Opus 4.7 1M, dry-run #2) | Created from requirement §8.9 ORD-006. |
