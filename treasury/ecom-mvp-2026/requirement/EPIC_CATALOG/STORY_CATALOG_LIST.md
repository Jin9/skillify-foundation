---
template_version: 0.1.0
story_id: STORY_CATALOG_LIST
epic_id: EPIC_CATALOG
title: Guest/Customer lists products with pagination, filter, sort
as_a: Guest
i_want: to list products with pagination, filter by category and price, and sort
so_that: I can find products that match what I'm looking for
acceptance_criteria:
  - given: 1000 ACTIVE products seeded in catalog
    when: Client GETs product listing with page=1 limit=20 sort=price_asc
    then: |
      Response returns exactly 20 items sorted ascending by price, all with
      status=ACTIVE and visible=true, and includes total count and page metadata
      (CAT-001, CAT-002).
  - given: Products belonging to category C and others outside it, where some have availableQty=0
    when: Client GETs product listing with categoryId=C and minPrice=100 and maxPrice=500 and inStock=true
    then: |
      Every returned item belongs to category C, has price in [100,500], and has
      availableQty greater than zero (CAT-003, CAT-004, CAT-005).
priority: Must
size: M
edge_cases:
  - case: Listing requested with limit greater than the documented maximum (e.g. limit=10000)
    expected: Server caps at the documented maximum and returns that many; total reflects the unfiltered match count
  - case: Listing requested with sort=invalid_value
    expected: 400 VALIDATION_FAILED on the sort param; no listing returned
  - case: Listing on an empty catalog
    expected: Response is 200 with items=[] and total=0
requirement_refs: [CAT-001, CAT-002, CAT-003, CAT-004, CAT-005, CAT-010]
definition_of_ready: [default]
definition_of_done: [default]
change_log:
  - date: 2026-05-08
    author: BA (Claude Opus 4.7 1M, dry-run #2)
    change: Created from requirement §8.3 CAT-001..005 + CAT-010.
---

# STORY_CATALOG_LIST — Guest/Customer lists products with pagination, filter, sort

## User narrative

As a **Guest**, I want **to list products with pagination, filter by category and price, and sort** so that **I can find products that match what I'm looking for**.

## Why this story exists in EPIC_CATALOG

This is the funnel mouth: the home tab and the search tab of the customer mobile webview both consume this endpoint. PERF-001 caps p95 listing latency at 300 ms on a 1,000-SKU dataset.

## Edge cases

| Case | Expected |
|---|---|
| limit greater than documented maximum | Capped at documented max; response notes total |
| invalid sort value | 400 VALIDATION_FAILED |
| empty catalog | 200 with items=[] and total=0 |

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | BA (Claude Opus 4.7 1M, dry-run #2) | Created from requirement §8.3 CAT-001..005. |
