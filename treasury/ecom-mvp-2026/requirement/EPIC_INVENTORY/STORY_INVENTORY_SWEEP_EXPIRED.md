---
template_version: 0.1.0
story_id: STORY_INVENTORY_SWEEP_EXPIRED
epic_id: EPIC_INVENTORY
title: Sweeper expires reservations past TTL and releases stock
as_a: System
i_want: a background sweeper that releases reservations whose TTL has elapsed
so_that: customers who walk away from checkout do not permanently lock stock
acceptance_criteria:
  - given: A reservation created at checkout that has aged past its TTL with the order still in PENDING_PAYMENT
    when: The sweeper job runs
    then: |
      The order moves PENDING_PAYMENT → PAYMENT_EXPIRED, the reserved stock is
      released back to availableQty, and a status_history row records the transition
      (INV-004, PAY-004, ORD-009).
  - given: A reservation created at checkout that is still within its TTL
    when: The sweeper job runs
    then: |
      The reservation is NOT touched and the order remains PENDING_PAYMENT (INV-004).
priority: Must
size: M
edge_cases:
  - case: Sweeper crashes mid-batch and is restarted
    expected: No reservation is double-released; the resume picks up where the previous run left off (idempotent within the batch)
  - case: Sweeper runs on a reservation whose order is already in CANCELLED state
    expected: Sweeper does not touch the reservation or order (already terminal)
requirement_refs: [INV-004, INV-006, INV-008, PAY-004, ORD-009]
definition_of_ready: [default]
definition_of_done: [default]
change_log:
  - date: 2026-05-08
    author: BA (Claude Opus 4.7 1M, dry-run #2)
    change: Created NEW for dry-run #2 to give the PAYMENT_EXPIRED edge case (browser walks away) a clear owner; runs/01 implied this but did not own it.
---

# STORY_INVENTORY_SWEEP_EXPIRED — Sweeper expires reservations past TTL and releases stock

## User narrative

As a **System**, I want **a background sweeper that releases reservations whose TTL has elapsed** so that **customers who walk away from checkout do not permanently lock stock**.

## Why this story exists in EPIC_INVENTORY

The Frontend Spec payment route may never receive a simulate action — the customer might close the browser. This story is the safety net: the sweeper guarantees that PENDING_PAYMENT cannot stay locked forever.

## Edge cases

| Case | Expected |
|---|---|
| Sweeper crash + restart | Idempotent; no double-release |
| Sweeper hits CANCELLED order | No-op |

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | BA (Claude Opus 4.7 1M, dry-run #2) | Created NEW for dry-run #2; sweeper has its own owner. |
