---
template_version: 0.1.0
story_id: STORY_CHECKOUT_COMMIT
epic_id: EPIC_CHECKOUT
title: Customer commits checkout (idempotent, reserves stock, creates order)
as_a: Customer
i_want: to commit my checkout safely even if my network retries
so_that: I create exactly one order and exactly one stock reservation regardless of retries
acceptance_criteria:
  - given: An authenticated CUSTOMER with a non-empty cart, a valid owned address, and an Idempotency-Key K1
    when: Customer POSTs checkout commit twice with the same Idempotency-Key K1
    then: |
      Exactly one order is created, both calls return that same order id, and stock is
      reserved exactly once (CHK-001, CHK-002, CHK-007, CHK-008, CHK-009).
  - given: A cart with one line of qty=5 against a SKU with availableQty=2
    when: Customer POSTs checkout commit
    then: |
      Response is 409 INSUFFICIENT_STOCK with details listing that productId,
      requestedQty=5, availableQty=2, and no order is created (CHK-004).
  - given: A valid checkout request whose body claims total=1 baht
    when: Customer POSTs checkout commit
    then: |
      The created order's total equals the server-computed sum of (current price *
      quantity) per item plus the §11.2 shipping fee, ignoring the client-supplied
      total (CHK-005).
  - given: A successful checkout that ordered the cart's only line
    when: Customer GETs cart afterward
    then: |
      Cart is empty (CHK-010).
priority: Must
size: L
edge_cases:
  - case: Client replays a checkout with a previously used Idempotency-Key but with a different cart payload
    expected: Server returns the originally created order for that key without creating a new order, and does NOT silently honor the new payload (CHK-009)
  - case: Customer POSTs commit with an addressId belonging to a different customer
    expected: 403 FORBIDDEN or VALIDATION_ERROR; no order created (CHK-002)
  - case: Two concurrent commits race on the last unit of a SKU
    expected: Exactly one succeeds with reservation; the other returns 409 INSUFFICIENT_STOCK; availableQty never goes negative (INV-008, CHK-004)
  - case: Cart contains a mix of in-stock and out-of-stock items
    expected: 409 INSUFFICIENT_STOCK with per-item details for every short item; no partial order is created (CHK-004)
  - case: Cart contains an INACTIVE/DELETED product
    expected: Per-item error referencing that product; no order created; stock not reserved (CHK-003, CART-005)
requirement_refs: [CHK-001, CHK-002, CHK-003, CHK-004, CHK-005, CHK-007, CHK-008, CHK-009, CHK-010, INV-008]
definition_of_ready: [default]
definition_of_done: [default]
change_log:
  - date: 2026-05-08
    author: BA (Claude Opus 4.7 1M, dry-run #2)
    change: Created from requirement §8.8 CHK-001..010; coupon (CHK-006) explicitly NOT in this story per scope.
---

# STORY_CHECKOUT_COMMIT — Customer commits checkout (idempotent, reserves stock, creates order)

## User narrative

As a **Customer**, I want **to commit my checkout safely even if my network retries** so that **I create exactly one order and exactly one stock reservation regardless of retries**.

## Why this story exists in EPIC_CHECKOUT

This is the consistency point of the entire purchase. Idempotency-Key + server-priced totals + per-item INSUFFICIENT_STOCK + cart-clear = a transaction that does not leak orders, stock, or money.

## Edge cases

| Case | Expected |
|---|---|
| Idempotency-Key replay with mismatched payload | Original order returned; new payload ignored |
| Address not owned | 403 FORBIDDEN or VALIDATION_ERROR |
| Concurrent commits on last unit | One wins; one gets INSUFFICIENT_STOCK |
| Mix of in-stock and out-of-stock | All short items listed; no partial order |
| INACTIVE/DELETED in cart | Per-item error; no order created |

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | BA (Claude Opus 4.7 1M, dry-run #2) | Created from requirement §8.8 CHK-001..010. |
