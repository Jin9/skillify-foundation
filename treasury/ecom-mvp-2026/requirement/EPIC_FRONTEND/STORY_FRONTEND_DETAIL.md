---
template_version: 0.1.0
story_id: STORY_FRONTEND_DETAIL
epic_id: EPIC_FRONTEND
title: detail route shows product info, ฿ price, stock status, add-to-cart CTA
as_a: Guest
i_want: a Thai-language product detail page with name, image, ฿ price, stock status, and an add-to-cart CTA
so_that: I can read the product's details and add it to my cart
acceptance_criteria:
  - given: An ACTIVE and visible product P
    when: The visitor opens the detail route for P
    then: |
      The detail page shows P's name, image placeholder, price in ฿ Thai-Baht prefix
      format, stock status, description, and an emerald 'add to cart' CTA;
      tapping the CTA when authenticated calls the cart service add endpoint and
      gives feedback when complete (CAT-006, CART-001).
priority: Must
size: M
edge_cases:
  - case: Visitor taps the add-to-cart CTA while unauthenticated
    expected: The frontend either prompts login or redirects to the login route, then resumes the add-to-cart action on success (Frontend Spec §3 'Auth?' for cart actions)
  - case: Product P has stock status=out-of-stock
    expected: The add-to-cart CTA is visually disabled and tapping it does not call the cart service
requirement_refs: [CAT-006, CART-001, "Frontend Spec §3"]
definition_of_ready: [default]
definition_of_done: [default]
change_log:
  - date: 2026-05-08
    author: BA (Claude Opus 4.7 1M, dry-run #2)
    change: Created NEW for dry-run #2 from Frontend Spec §3 detail row.
---

# STORY_FRONTEND_DETAIL — detail route shows product info, ฿ price, stock status, add-to-cart CTA

## User narrative

As a **Guest**, I want **a Thai-language product detail page with name, image, ฿ price, stock status, and an add-to-cart CTA** so that **I can read the product's details and add it to my cart**.

## Why this story exists in EPIC_FRONTEND

Detail is the conversion point: this is where the customer commits to add-to-cart.

## Edge cases

| Case | Expected |
|---|---|
| Add-to-cart while unauthenticated | Login prompt/redirect, then resume |
| Out-of-stock product | CTA disabled |

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | BA (Claude Opus 4.7 1M, dry-run #2) | Created NEW for dry-run #2 from Frontend Spec §3. |
