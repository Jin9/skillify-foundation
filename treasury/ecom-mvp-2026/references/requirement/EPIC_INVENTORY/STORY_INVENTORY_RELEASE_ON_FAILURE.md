---
template_version: 0.1.0
story_id: STORY_INVENTORY_RELEASE_ON_FAILURE
epic_id: EPIC_INVENTORY
title: Payment failed/timeout/customer-cancel releases reserved → available
as_a: System
i_want: to release a reservation back to availableQty whenever the order does not settle
so_that: abandoned or failed orders do not lock inventory permanently
acceptance_criteria:
  - given: A SKU with availableQty=7 reservedQty=3 soldQty=0 from an active reservation
    when: Payment-failed is simulated for the owning order
    then: |
      availableQty=10 reservedQty=0 soldQty=0 and the order moves to PAYMENT_FAILED
      (INV-006, PAY-003).
  - given: A SKU with availableQty=5 reservedQty=2 soldQty=0 and a customer-cancel on the owning PENDING_PAYMENT order
    when: Customer cancels the order
    then: |
      availableQty=7 reservedQty=0 soldQty=0; reserved stock is released and the
      order is CANCELLED (INV-006, ORD-004).
priority: Must
size: M
edge_cases:
  - case: Two release operations attempted on the same reservation (e.g. failed callback delivered twice)
    expected: Idempotent on (paymentIntentId, providerStatus); reservation released exactly once (PAY-005, INV-006)
  - case: Release on a reservation that has already been converted to sold (e.g. payment succeeded first)
    expected: 409 INVALID_ORDER_STATE on the order; no stock change
requirement_refs: [INV-006, PAY-003, PAY-004, PAY-005, ORD-004]
definition_of_ready: [default]
definition_of_done: [default]
change_log:
  - date: 2026-05-08
    author: BA (Claude Opus 4.7 1M, dry-run #2)
    change: Created from requirement §8.5 INV-006 + §8.10 PAY-003/004/005 + §8.9 ORD-004.
---

# STORY_INVENTORY_RELEASE_ON_FAILURE — Payment failed/timeout/customer-cancel releases reserved → available

## User narrative

As a **System**, I want **to release a reservation back to availableQty whenever the order does not settle** so that **abandoned or failed orders do not lock inventory permanently**.

## Why this story exists in EPIC_INVENTORY

This is the failure-path coverage of the inventory state machine. Without it, every abandoned cart or failed payment would silently consume stock until manual cleanup.

## Edge cases

| Case | Expected |
|---|---|
| Duplicate release | Idempotent; one release |
| Release on already-sold | 409 INVALID_ORDER_STATE; no stock change |

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | BA (Claude Opus 4.7 1M, dry-run #2) | Created from requirement §8.5 + §8.10 + §8.9. |
