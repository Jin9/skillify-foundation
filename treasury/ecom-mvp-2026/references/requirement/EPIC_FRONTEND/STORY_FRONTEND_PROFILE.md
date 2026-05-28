---
template_version: 0.1.0
story_id: STORY_FRONTEND_PROFILE
epic_id: EPIC_FRONTEND
title: profile tab shows name/email/default address and exposes logout
as_a: Customer
i_want: a Thai-language profile tab showing my name, email, default address, and a logout action
so_that: I can review my profile and log out
acceptance_criteria:
  - given: An authenticated customer on the profile tab
    when: The tab renders and the customer taps the logout action
    then: |
      The tab shows the customer's name, email, and default address (or an empty
      state for default address if none); the logout action calls the identity
      logout endpoint, clears the locally cached access token, and routes the
      customer to the home tab as an unauthenticated visitor (CUST-001, AUTH-004).
priority: Must
size: S
edge_cases:
  - case: Customer's logout call fails due to network error
    expected: The frontend still clears the local token cache and routes to home; the customer is effectively logged out client-side, and the next refresh attempt with the stale token will be rejected when connectivity returns
requirement_refs: [CUST-001, AUTH-004, "Frontend Spec §3"]
definition_of_ready: [default]
definition_of_done: [default]
change_log:
  - date: 2026-05-08
    author: BA (Claude Opus 4.7 1M, dry-run #2)
    change: Created NEW for dry-run #2 from Frontend Spec §3 profile row.
---

# STORY_FRONTEND_PROFILE — profile tab shows name/email/default address and exposes logout

## User narrative

As a **Customer**, I want **a Thai-language profile tab showing my name, email, default address, and a logout action** so that **I can review my profile and log out**.

## Why this story exists in EPIC_FRONTEND

The profile tab is the fifth bottom-tab route. It's the customer's session-control surface.

## Edge cases

| Case | Expected |
|---|---|
| Logout network failure | Local token cleared anyway; effective logout client-side |

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | BA (Claude Opus 4.7 1M, dry-run #2) | Created NEW for dry-run #2 from Frontend Spec §3. |
