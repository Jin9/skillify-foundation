---
template_version: 0.1.0
story_id: STORY_CATALOG_ADMIN_CATEGORY
epic_id: EPIC_CATALOG
title: Admin creates, edits, deactivates a category
as_a: Admin
i_want: to create, edit, and deactivate categories
so_that: the catalog organisation reflects the merchant's structure
acceptance_criteria:
  - given: An ADMIN token and an existing category at the same level as another with the same proposed name
    when: Admin POSTs category create
    then: |
      Response is 409 CONFLICT and no category is created (CATE-001).
  - given: An ADMIN token and an existing category C with status=ACTIVE
    when: Admin PATCHes category C to status=INACTIVE
    then: |
      C remains in the database, products in C continue to exist, but C no longer
      appears in category-filter listings (CATE-003).
priority: Must
size: S
edge_cases:
  - case: Admin attempts to create a category with a slug that already exists globally
    expected: 409 CONFLICT on the slug; no category created (CATE-002)
  - case: Customer token attempts the admin-category endpoint
    expected: 403 FORBIDDEN (AUTH-003)
requirement_refs: [CATE-001, CATE-002, CATE-003, AUTH-003]
definition_of_ready: [default]
definition_of_done: [default]
change_log:
  - date: 2026-05-08
    author: BA (Claude Opus 4.7 1M, dry-run #2)
    change: Created from requirement §8.4 CATE-001..003.
---

# STORY_CATALOG_ADMIN_CATEGORY — Admin creates, edits, deactivates a category

## User narrative

As an **Admin**, I want **to create, edit, and deactivate categories** so that **the catalog organisation reflects the merchant's structure**.

## Why this story exists in EPIC_CATALOG

Categories drive the search-tab filter chips on the customer mobile webview. Backend-only this run; no admin UI.

## Edge cases

| Case | Expected |
|---|---|
| Duplicate slug globally | 409 CONFLICT |
| CUSTOMER token | 403 FORBIDDEN |

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | BA (Claude Opus 4.7 1M, dry-run #2) | Created from requirement §8.4 CATE-001..003. |
