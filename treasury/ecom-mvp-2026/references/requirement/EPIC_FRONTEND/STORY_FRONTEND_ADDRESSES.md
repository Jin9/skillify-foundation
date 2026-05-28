---
template_version: 0.1.0
story_id: STORY_FRONTEND_ADDRESSES
epic_id: EPIC_FRONTEND
title: addresses route lists/add addresses; supports picker mode for checkout
as_a: Customer
i_want: a Thai-language addresses route to list and add addresses, plus a picker mode used by checkout
so_that: I can manage my shipping addresses and pick one during checkout
acceptance_criteria:
  - given: An authenticated customer on the addresses route with at least one saved address
    when: The route renders
    then: |
      Each address card shows the receiver name, phone, address line, district,
      province, and postal code; a default-address indicator marks the customer's
      current default; the page exposes an 'add new address' CTA in emerald
      (Frontend Spec §3 + CUST-003, CUST-004).
  - given: An authenticated customer entering the addresses route from the checkout route in picker mode
    when: The customer taps an address card to select it
    then: |
      The frontend returns the customer to the checkout route with that address
      preselected; checkout's totals and items recalculate against the picked
      address (Frontend Spec §3 + CHK-002).
priority: Must
size: M
edge_cases:
  - case: Customer submits the add-address form with a missing required field
    expected: The form surfaces a field-level Thai error per the standard envelope; no address row created
  - case: Customer enters picker mode from checkout with no saved addresses
    expected: The page presents the add-new-address form first, then returns to checkout with the new address selected
requirement_refs: [CUST-003, CUST-004, CHK-002, "Frontend Spec §3"]
definition_of_ready: [default]
definition_of_done: [default]
change_log:
  - date: 2026-05-08
    author: BA (Claude Opus 4.7 1M, dry-run #2)
    change: Created NEW for dry-run #2 from Frontend Spec §3 addresses row.
---

# STORY_FRONTEND_ADDRESSES — addresses route lists/add addresses; supports picker mode for checkout

## User narrative

As a **Customer**, I want **a Thai-language addresses route to list and add addresses, plus a picker mode used by checkout** so that **I can manage my shipping addresses and pick one during checkout**.

## Why this story exists in EPIC_FRONTEND

The addresses route serves two purposes: a standalone management surface and a picker invoked from checkout. CUST-003 + CUST-004 + CHK-002 are the binding rules.

## Edge cases

| Case | Expected |
|---|---|
| Add-address with missing field | Field-level Thai error |
| Picker mode with no saved addresses | Add form first, then return to checkout |

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | BA (Claude Opus 4.7 1M, dry-run #2) | Created NEW for dry-run #2 from Frontend Spec §3. |
