---
template_version: 0.1.0
story_id: STORY_FRONTEND_SEARCH
epic_id: EPIC_FRONTEND
title: search tab supports query input, category filter chips, and product grid with pagination
as_a: Guest
i_want: to search and filter products from the search tab
so_that: I can find the products I'm looking for using a query, category, and price filters
acceptance_criteria:
  - given: A visitor on the search tab with a populated catalog
    when: The visitor types a query and selects category and price filters
    then: |
      The product grid re-renders with results matching the catalog service's
      pagination + filter contract (page, limit, categoryId, minPrice, maxPrice,
      inStock); each card carries name, ฿ price, category, and tapping it opens the
      detail route (Frontend Spec §3 search row + CAT-001..005).
  - given: A visitor on the search tab whose filters yield zero results
    when: The grid renders the response
    then: |
      An empty-state placeholder in Thai is shown; the bottom tab bar remains
      anchored; resetting filters re-renders the unfiltered result set.
priority: Must
size: M
edge_cases:
  - case: Visitor scrolls past the first page of results
    expected: The grid loads the next page (page=2) using the documented pagination contract; the visible items grow without re-rendering the existing items
  - case: Visitor toggles inStock=true filter
    expected: Only products with availableQty greater than zero render in the grid
requirement_refs: [CAT-001, CAT-002, CAT-003, CAT-004, CAT-005, "Frontend Spec §3"]
definition_of_ready: [default]
definition_of_done: [default]
change_log:
  - date: 2026-05-08
    author: BA (Claude Opus 4.7 1M, dry-run #2)
    change: Created NEW for dry-run #2 from Frontend Spec §3 search row.
---

# STORY_FRONTEND_SEARCH — search tab supports query input, category filter chips, and product grid with pagination

## User narrative

As a **Guest**, I want **to search and filter products from the search tab** so that **I can find the products I'm looking for using a query, category, and price filters**.

## Why this story exists in EPIC_FRONTEND

The search tab is the second of the five bottom-tab routes and the primary discovery surface beyond the home tab.

## Edge cases

| Case | Expected |
|---|---|
| Pagination | Next-page load grows the visible list |
| inStock=true | Only stocked products render |

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | BA (Claude Opus 4.7 1M, dry-run #2) | Created NEW for dry-run #2 from Frontend Spec §3. |
