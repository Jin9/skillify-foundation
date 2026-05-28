---
template_version: 0.1.0
story_id: STORY_CART_UPDATE_QTY
epic_id: EPIC_CART
title: Customer updates quantity of a cart line
as_a: Customer
i_want: to update the quantity of a cart line
so_that: I can buy more or less of an item without removing it
acceptance_criteria:
  - given: An authenticated CUSTOMER with a cart line for product P qty=2 and P availableQty=5
    when: Customer POSTs cart edit setting qty=3
    then: |
      Response is 200 with code=SUCCESS, the cart line for P now has quantity=3, and
      cart subtotal reflects 3*P.price (CART-002, CART-004).
  - given: An authenticated CUSTOMER with a cart line for product P qty=2 and P availableQty=1
    when: Customer POSTs cart edit setting qty=3
    then: |
      Response is a validation error referencing availableQty=1 and the cart line is
      unchanged (still qty=2) (CART-002).
priority: Must
size: S
edge_cases:
  - case: Cart-edit with qty=0
    expected: Either 400 VALIDATION_FAILED on the qty field, OR the line is removed (Tech-Designer chooses); BA only requires that qty=0 cannot end with a zero-qty line in the cart
  - case: Cart-edit on a cart line that does not belong to the requesting customer
    expected: NOT_FOUND or FORBIDDEN; the line is unchanged
  - case: Cart-edit referencing a cartItemId that does not exist
    expected: NOT_FOUND
requirement_refs: [CART-002, CART-004]
definition_of_ready: [default]
definition_of_done: [default]
change_log:
  - date: 2026-05-08
    author: BA (Claude Opus 4.7 1M, dry-run #2)
    change: Created from requirement §8.6 CART-002 + CART-004.
---

# STORY_CART_UPDATE_QTY — Customer updates quantity of a cart line

## User narrative

As a **Customer**, I want **to update the quantity of a cart line** so that **I can buy more or less of an item without removing it**.

## Why this story exists in EPIC_CART

The cart tab's quantity stepper (Frontend Spec §3) is built on this endpoint. CART-002's availableQty cap is a best-effort UX guard; the authoritative cap lives in checkout (CHK-004).

## Edge cases

| Case | Expected |
|---|---|
| qty=0 | Either 400 or remove-the-line; BA only requires no zero-qty lines persist |
| Cart line not owned by requester | NOT_FOUND or FORBIDDEN |
| Unknown cartItemId | NOT_FOUND |

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | BA (Claude Opus 4.7 1M, dry-run #2) | Created from requirement §8.6 CART-002 + CART-004. |
