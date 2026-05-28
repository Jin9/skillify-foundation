---
template_version: 0.1.0
story_id: STORY_CATALOG_ADMIN_PRODUCT_CRUD
epic_id: EPIC_CATALOG
title: Admin creates, edits, soft-deletes a product
as_a: Admin
i_want: to create, edit, and soft-delete products
so_that: the catalog reflects the merchant's current inventory without breaking historical orders
acceptance_criteria:
  - given: An ADMIN token and a category that exists
    when: Admin POSTs product create with unique SKU, price > 0, that categoryId
    then: |
      Product is created with status DRAFT or ACTIVE per request, the returned id is a
      UUID, and a duplicate-SKU retry returns 409 CONFLICT (CAT-007, CAT-010).
  - given: An ADMIN token and an existing product
    when: Admin PUTs product edit with new price and new status
    then: |
      The product reflects the new fields; the edit is recorded server-side for the
      future audit hookup (CAT-008).
  - given: An ADMIN token and an existing product that appears in a historical order
    when: Admin POSTs product soft-delete on that product
    then: |
      Listing endpoints stop returning the product, but GETting the prior order's
      detail still resolves the product's snapshot fields (CAT-009).
priority: Must
size: M
edge_cases:
  - case: Admin creates a product with price=0 or negative price
    expected: 400 VALIDATION_FAILED on the price field; no product created
  - case: Admin creates a product with a categoryId that does not exist
    expected: 400 VALIDATION_FAILED citing the unknown category; no product created
  - case: Customer token attempts the admin-create endpoint
    expected: 403 FORBIDDEN; no product created (AUTH-003)
requirement_refs: [CAT-007, CAT-008, CAT-009, CAT-010, AUTH-003]
definition_of_ready: [default]
definition_of_done: [default]
change_log:
  - date: 2026-05-08
    author: BA (Claude Opus 4.7 1M, dry-run #2)
    change: Created from requirement §8.3 CAT-007..010; admin UI deferred but endpoints kept in MVP.
---

# STORY_CATALOG_ADMIN_PRODUCT_CRUD — Admin creates, edits, soft-deletes a product

## User narrative

As an **Admin**, I want **to create, edit, and soft-delete products** so that **the catalog reflects the merchant's current inventory without breaking historical orders**.

## Why this story exists in EPIC_CATALOG

Admin product management is in scope at the *backend* layer this run (no admin UI). The endpoints are exercised via Postman/curl. CAT-009's soft-delete-vs-history rule is the single highest-stakes correctness rule in this epic.

## Edge cases

| Case | Expected |
|---|---|
| price <= 0 | 400 VALIDATION_FAILED |
| Unknown categoryId | 400 VALIDATION_FAILED |
| CUSTOMER token on admin endpoint | 403 FORBIDDEN |

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | BA (Claude Opus 4.7 1M, dry-run #2) | Created from requirement §8.3 CAT-007..010. |
