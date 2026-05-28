---
template_version: 0.1.0
story_id: STORY_INVENTORY_ADMIN_ADJUST
epic_id: EPIC_INVENTORY
title: Admin adjusts stock; record persists delta, reason, actorUserId
as_a: Admin
i_want: to adjust stock with a delta and a reason
so_that: physical inventory matches recorded inventory and the change is auditable
acceptance_criteria:
  - given: An ADMIN token and a SKU
    when: Admin POSTs stock-adjust with delta=-2 and reason='damaged'
    then: |
      availableQty decreases by 2 and a stock-adjustment record exists with the
      delta, reason, and admin actorUserId (INV-007).
  - given: An ADMIN token and a SKU with availableQty=1
    when: Admin POSTs stock-adjust with delta=-2 and reason='correction'
    then: |
      Adjustment is rejected because the resulting availableQty would be -1; no stock
      change and no record persisted (INV-008).
priority: Must
size: S
edge_cases:
  - case: Admin POSTs stock-adjust with delta=0
    expected: 400 VALIDATION_FAILED on delta; no record persisted
  - case: Admin POSTs stock-adjust without a reason
    expected: 400 VALIDATION_FAILED on reason; no record persisted
  - case: CUSTOMER token attempts stock-adjust
    expected: 403 FORBIDDEN (AUTH-003)
requirement_refs: [INV-007, INV-008, AUTH-003]
definition_of_ready: [default]
definition_of_done: [default]
change_log:
  - date: 2026-05-08
    author: BA (Claude Opus 4.7 1M, dry-run #2)
    change: Created from requirement §8.5 INV-007 + INV-008.
---

# STORY_INVENTORY_ADMIN_ADJUST — Admin adjusts stock; record persists delta, reason, actorUserId

## User narrative

As an **Admin**, I want **to adjust stock with a delta and a reason** so that **physical inventory matches recorded inventory and the change is auditable**.

## Why this story exists in EPIC_INVENTORY

INV-007 is the only customer-of-system path for stock changes outside the order lifecycle. The audit module is deferred (§8.15), but the adjustment record itself is in scope.

## Edge cases

| Case | Expected |
|---|---|
| delta=0 | 400 VALIDATION_FAILED |
| missing reason | 400 VALIDATION_FAILED |
| CUSTOMER token | 403 FORBIDDEN |

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | BA (Claude Opus 4.7 1M, dry-run #2) | Created from requirement §8.5 INV-007. |
