---
template_version: 0.1.0
story_id: STORY_FRONTEND_ORDERS
epic_id: EPIC_FRONTEND
title: orders tab lists customer's orders with all/active/done filter and Thai status badges
as_a: Customer
i_want: a Thai-language orders tab listing all my orders with status badges and an all/active/done filter
so_that: I can quickly find an active order or a past order
acceptance_criteria:
  - given: An authenticated customer with three orders in different statuses
    when: The customer opens the orders tab and toggles the all / active / done filter
    then: |
      The list renders one row per order with a coloured status badge carrying a
      Thai label, and tapping a row opens orderDetail for that order (Frontend Spec
      §3 + ORD-001).
  - given: An authenticated customer with no orders
    when: The orders tab renders
    then: |
      An empty-state placeholder in Thai is shown with a CTA back to the home tab
      (Frontend Spec §3).
priority: Must
size: M
edge_cases:
  - case: Customer's orders include one in PAYMENT_EXPIRED
    expected: The expired order shows under "done" (or the documented filter bucket) with a distinct Thai label and the appropriate badge colour
  - case: Customer scrolls past the first page of orders
    expected: The list loads the next page using documented pagination; existing rows do not re-render
requirement_refs: [ORD-001, ORD-002, "Frontend Spec §3"]
definition_of_ready: [default]
definition_of_done: [default]
change_log:
  - date: 2026-05-08
    author: BA (Claude Opus 4.7 1M, dry-run #2)
    change: Created NEW for dry-run #2 from Frontend Spec §3 orders row.
---

# STORY_FRONTEND_ORDERS — orders tab lists customer's orders with all/active/done filter and Thai status badges

## User narrative

As a **Customer**, I want **a Thai-language orders tab listing all my orders with status badges and an all/active/done filter** so that **I can quickly find an active order or a past order**.

## Why this story exists in EPIC_FRONTEND

The orders tab is the customer's hub for order tracking. Status badges with Thai labels are explicitly mentioned in Frontend Spec §6 (`<StatusBadge>`).

## Edge cases

| Case | Expected |
|---|---|
| PAYMENT_EXPIRED in list | Distinct badge + Thai label |
| Pagination | Next page grows the list |

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | BA (Claude Opus 4.7 1M, dry-run #2) | Created NEW for dry-run #2 from Frontend Spec §3. |
