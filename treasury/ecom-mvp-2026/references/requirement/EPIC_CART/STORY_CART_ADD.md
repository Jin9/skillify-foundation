---
template_version: 0.1.0
story_id: STORY_CART_ADD
epic_id: EPIC_CART
title: Customer adds an active product to cart
as_a: Customer
i_want: to add an active product to my cart with a quantity greater than zero
so_that: I can collect items before checkout
acceptance_criteria:
  - given: An authenticated CUSTOMER and an ACTIVE product P with availableQty=10
    when: Customer POSTs cart add for P qty=2
    then: |
      Response is 200 with code=SUCCESS, the cart contains a single line for P with
      quantity=2, the cart's subtotal reflects 2*P.price, and updatedAt is advanced
      (CART-001, CART-004, CART-007).
  - given: An authenticated CUSTOMER and an ACTIVE product P
    when: Customer POSTs cart add for P qty=2 twice in succession
    then: |
      Cart contains a single line for P with quantity=4 and updatedAt advanced; the
      cart NEVER contains two separate lines for the same SKU (CART-006, CART-007).
priority: Must
size: S
edge_cases:
  - case: Cart-add for a product whose status is INACTIVE or DELETED
    expected: 400 VALIDATION_FAILED on the productId; cart unchanged (CART-001, CART-005)
  - case: Cart-add with qty=0 or qty=-1
    expected: 400 VALIDATION_FAILED on the qty field; cart unchanged
  - case: Cart-add by an unauthenticated guest
    expected: 401 UNAUTHORIZED; cart unchanged
requirement_refs: [CART-001, CART-004, CART-005, CART-006, CART-007]
definition_of_ready: [default]
definition_of_done: [default]
change_log:
  - date: 2026-05-08
    author: BA (Claude Opus 4.7 1M, dry-run #2)
    change: Created from requirement §8.6 CART-001 + CART-006 + CART-007.
---

# STORY_CART_ADD — Customer adds an active product to cart

## User narrative

As a **Customer**, I want **to add an active product to my cart with a quantity greater than zero** so that **I can collect items before checkout**.

## Why this story exists in EPIC_CART

This is the entry point into the cart. The merge-on-duplicate rule (CART-006) is what keeps the cart UI clean — the customer never sees two separate lines for the same SKU.

## Edge cases

| Case | Expected |
|---|---|
| Product status INACTIVE/DELETED | 400 VALIDATION_FAILED |
| qty <= 0 | 400 VALIDATION_FAILED |
| Unauthenticated | 401 UNAUTHORIZED |

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | BA (Claude Opus 4.7 1M, dry-run #2) | Created from requirement §8.6 CART-001 + CART-006 + CART-007. |
