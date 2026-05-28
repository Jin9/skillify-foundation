---
template_version: 0.1.0
story_id: STORY_FRONTEND_CART
epic_id: EPIC_FRONTEND
title: cart tab lists lines with quantity stepper, remove, server-computed subtotal
as_a: Customer
i_want: a Thai-language cart tab with each line's image, name, ฿ price, quantity stepper, remove, and a server-computed subtotal
so_that: I can review and adjust my purchases before checkout
acceptance_criteria:
  - given: An authenticated CUSTOMER with three cart lines on the cart tab
    when: The cart tab renders
    then: |
      Each line shows the placeholder image, Thai product name, unit price with ฿
      prefix, line subtotal, a quantity stepper (+/−), and a remove action; the cart
      subtotal at the bottom equals the sum of line subtotals returned by the cart
      service (CART-003, CART-004).
  - given: An authenticated CUSTOMER on the cart tab with a line for product P qty=2
    when: The customer taps the + on the quantity stepper
    then: |
      The frontend calls the cart service edit endpoint with qty=3; on the success
      response the cart re-renders with qty=3 and the updated subtotal (CART-002,
      CART-004).
priority: Must
size: M
edge_cases:
  - case: Customer attempts to step quantity above availableQty
    expected: The cart service returns the documented validation error; the cart tab surfaces a Thai message and the qty does not increase
  - case: Customer's cart contains an INACTIVE/DELETED product
    expected: The line is visually flagged non-checkoutable and the 'go to checkout' CTA does not navigate until the line is removed (CART-005, CHK-003)
requirement_refs: [CART-002, CART-003, CART-004, CART-005, "Frontend Spec §3"]
definition_of_ready: [default]
definition_of_done: [default]
change_log:
  - date: 2026-05-08
    author: BA (Claude Opus 4.7 1M, dry-run #2)
    change: Created NEW for dry-run #2 from Frontend Spec §3 cart row.
---

# STORY_FRONTEND_CART — cart tab lists lines with quantity stepper, remove, server-computed subtotal

## User narrative

As a **Customer**, I want **a Thai-language cart tab with each line's image, name, ฿ price, quantity stepper, remove, and a server-computed subtotal** so that **I can review and adjust my purchases before checkout**.

## Why this story exists in EPIC_FRONTEND

The cart tab is one of the 5 bottom-tab routes and the staging ground for checkout. Its critical responsibility is to surface the server-computed subtotal — the frontend must NEVER compute totals locally.

## Edge cases

| Case | Expected |
|---|---|
| Step over availableQty | Server validation error; qty unchanged |
| INACTIVE line in cart | Flagged; 'go to checkout' blocked until removed |

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | BA (Claude Opus 4.7 1M, dry-run #2) | Created NEW for dry-run #2 from Frontend Spec §3. |
