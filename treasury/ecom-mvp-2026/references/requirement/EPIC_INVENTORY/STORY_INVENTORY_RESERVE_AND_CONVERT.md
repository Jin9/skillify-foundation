---
template_version: 0.1.0
story_id: STORY_INVENTORY_RESERVE_AND_CONVERT
epic_id: EPIC_INVENTORY
title: Checkout reserves; payment success converts reserved → sold
as_a: System
i_want: to reserve stock at checkout and convert reserved → sold on payment success
so_that: stock movements are transactionally safe and visible
acceptance_criteria:
  - given: A SKU with availableQty=10, reservedQty=0, soldQty=0
    when: Checkout reserves 3 units and then payment-success is simulated
    then: |
      After reservation availableQty=7 reservedQty=3 soldQty=0; after payment success
      availableQty=7 reservedQty=0 soldQty=3 (INV-002, INV-003, INV-005).
  - given: A SKU with availableQty=2 and a checkout requesting qty=5
    when: Checkout attempts reservation
    then: |
      Reservation fails with 409 INSUFFICIENT_STOCK and the SKU's availableQty stays
      at 2 (INV-008, CHK-004).
priority: Must
size: M
edge_cases:
  - case: Two concurrent reservations on the last unit
    expected: One succeeds with reservation; the other returns 409 INSUFFICIENT_STOCK; availableQty never goes negative (INV-008)
  - case: Reservation for qty greater than availableQty by exactly 1
    expected: 409 INSUFFICIENT_STOCK referencing the requestedQty and the available
requirement_refs: [INV-002, INV-003, INV-005, INV-008, CHK-004]
definition_of_ready: [default]
definition_of_done: [default]
change_log:
  - date: 2026-05-08
    author: BA (Claude Opus 4.7 1M, dry-run #2)
    change: Created from requirement §8.5 INV-002 + INV-003 + INV-005 + INV-008.
---

# STORY_INVENTORY_RESERVE_AND_CONVERT — Checkout reserves; payment success converts reserved → sold

## User narrative

As a **System**, I want **to reserve stock at checkout and convert reserved → sold on payment success** so that **stock movements are transactionally safe and visible**.

## Why this story exists in EPIC_INVENTORY

This is the happy path of the inventory state machine. The §10.3 worked example is the spec; this story tests both halves of it.

## Edge cases

| Case | Expected |
|---|---|
| Concurrent reservation on last unit | One wins; one INSUFFICIENT_STOCK; never negative |
| Reservation off-by-1 over available | 409 INSUFFICIENT_STOCK with details |

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | BA (Claude Opus 4.7 1M, dry-run #2) | Created from requirement §8.5. |
