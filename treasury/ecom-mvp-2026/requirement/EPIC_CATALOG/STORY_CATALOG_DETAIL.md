---
template_version: 0.1.0
story_id: STORY_CATALOG_DETAIL
epic_id: EPIC_CATALOG
title: Guest/Customer reads product detail
as_a: Guest
i_want: to read a product's detail (name, description, images, price, stock status, category, review summary)
so_that: I can decide whether to add it to cart
acceptance_criteria:
  - given: An ACTIVE and visible product P
    when: Client GETs product detail by P's id
    then: |
      Response is 200 with code=SUCCESS and data contains name, description, images,
      price, stock status, category, and a review summary block (placeholder this run
      since reviews are deferred) (CAT-006).
  - given: A product P that has been soft-deleted by admin
    when: Client GETs product detail by P's id
    then: |
      The listing endpoint no longer returns P, and the detail endpoint either returns
      NOT_FOUND or returns the product but flagged so it cannot be added to cart;
      historical orders that include P still resolve P's snapshot fields from the order
      item (CAT-009).
priority: Must
size: S
edge_cases:
  - case: Detail requested for a product whose status is DRAFT
    expected: NOT_FOUND for guests/customers; only admins can see DRAFT
  - case: Detail requested for a malformed product id (non-UUID)
    expected: 400 VALIDATION_FAILED
requirement_refs: [CAT-006, CAT-009, CAT-010]
definition_of_ready: [default]
definition_of_done: [default]
change_log:
  - date: 2026-05-08
    author: BA (Claude Opus 4.7 1M, dry-run #2)
    change: Created from requirement §8.3 CAT-006 + CAT-009 + CAT-010.
---

# STORY_CATALOG_DETAIL — Guest/Customer reads product detail

## User narrative

As a **Guest**, I want **to read a product's detail** so that **I can decide whether to add it to cart**.

## Why this story exists in EPIC_CATALOG

The detail route on the customer mobile webview (Frontend Spec §3) is built on this endpoint. CAT-009's "still resolves in historical orders" is a strong consistency requirement that lives here AND in the order epic.

## Edge cases

| Case | Expected |
|---|---|
| Detail for DRAFT product | NOT_FOUND for guests/customers |
| Malformed product id | 400 VALIDATION_FAILED |

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | BA (Claude Opus 4.7 1M, dry-run #2) | Created from requirement §8.3 CAT-006 + CAT-009. |
