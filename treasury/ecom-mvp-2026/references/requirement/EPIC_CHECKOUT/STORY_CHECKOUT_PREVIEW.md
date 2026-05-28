---
template_version: 0.1.0
story_id: STORY_CHECKOUT_PREVIEW
epic_id: EPIC_CHECKOUT
title: Customer previews a checkout (server-priced totals before commit)
as_a: Customer
i_want: to preview my checkout totals before committing
so_that: I can review the server-computed subtotal, shipping fee, and grand total before placing the order
acceptance_criteria:
  - given: An authenticated CUSTOMER with a non-empty cart and a valid owned address
    when: Customer POSTs checkout preview with that addressId
    then: |
      Response is 200 with code=SUCCESS and data contains items (with snapshot price),
      subtotal, shippingFee, grandTotal — all computed server-side from authoritative
      product prices; no order is created and no stock is reserved (CHK-001, CHK-002,
      CHK-005).
  - given: A cart whose only product transitioned to INACTIVE after add
    when: Customer POSTs checkout preview
    then: |
      Response rejects with a per-item error referencing that product, no order is
      created, and stock is not reserved (CHK-003, CART-005).
priority: Must
size: M
edge_cases:
  - case: Customer POSTs preview with an addressId belonging to a different customer
    expected: 403 FORBIDDEN or VALIDATION_ERROR; no order created (CHK-002)
  - case: Customer POSTs preview with an empty cart
    expected: 400 VALIDATION_FAILED citing the empty cart; no order created (CHK-001)
requirement_refs: [CHK-001, CHK-002, CHK-003, CHK-005]
definition_of_ready: [default]
definition_of_done: [default]
change_log:
  - date: 2026-05-08
    author: BA (Claude Opus 4.7 1M, dry-run #2)
    change: Created from requirement §8.8; preview is non-mutating per §13.4.
---

# STORY_CHECKOUT_PREVIEW — Customer previews a checkout (server-priced totals before commit)

## User narrative

As a **Customer**, I want **to preview my checkout totals before committing** so that **I can review the server-computed subtotal, shipping fee, and grand total before placing the order**.

## Why this story exists in EPIC_CHECKOUT

The mobile webview's checkout route renders totals from this endpoint. CHK-005 — server-priced everything — is the single most important rule.

## Edge cases

| Case | Expected |
|---|---|
| Address belongs to another customer | 403 FORBIDDEN or VALIDATION_ERROR |
| Empty cart | 400 VALIDATION_FAILED |

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | BA (Claude Opus 4.7 1M, dry-run #2) | Created from requirement §8.8. |
