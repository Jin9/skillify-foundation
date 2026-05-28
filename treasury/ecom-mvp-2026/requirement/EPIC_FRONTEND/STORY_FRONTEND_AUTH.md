---
template_version: 0.1.0
story_id: STORY_FRONTEND_AUTH
epic_id: EPIC_FRONTEND
title: login route supports register and login toggle; on success token is stored client-side
as_a: Guest
i_want: a single Thai-language login route that lets me toggle between login and register
so_that: I can authenticate quickly without leaving the app
acceptance_criteria:
  - given: An unauthenticated visitor on the login route in login mode
    when: The visitor submits valid email and password
    then: |
      The frontend calls the identity login endpoint, stores the returned access
      token client-side, and navigates back to the originally requested route or
      the home tab if no route was pending (US-AUTH-002).
  - given: An unauthenticated visitor on the login route in register mode
    when: The visitor submits a valid registration form (email, password, name)
    then: |
      The frontend calls the identity register endpoint and, on success, transitions
      the same route into login mode (or auto-logs-in per design); the customer is
      then routed home or to the originally requested route (US-AUTH-001).
priority: Must
size: M
edge_cases:
  - case: Visitor submits login with wrong password
    expected: The form surfaces the standard envelope error message in Thai (the backend's generic AUTH-007 message); no token is stored
  - case: Visitor submits register with an already-registered email
    expected: The form surfaces a duplicate-email error message in Thai per the standard envelope; no second account is created
  - case: Visitor toggles between login and register repeatedly
    expected: Each toggle clears the form so credentials never carry across modes
requirement_refs: [AUTH-001, AUTH-002, AUTH-007, "Frontend Spec §3"]
definition_of_ready: [default]
definition_of_done: [default]
change_log:
  - date: 2026-05-08
    author: BA (Claude Opus 4.7 1M, dry-run #2)
    change: Created NEW for dry-run #2 from Frontend Spec §3 login row.
---

# STORY_FRONTEND_AUTH — login route supports register and login toggle; on success token is stored client-side

## User narrative

As a **Guest**, I want **a single Thai-language login route that lets me toggle between login and register** so that **I can authenticate quickly without leaving the app**.

## Why this story exists in EPIC_FRONTEND

The login route is the gateway to every authenticated route. The toggle pattern (login ↔ register on one route) is from the Frontend Spec §3.

## Edge cases

| Case | Expected |
|---|---|
| Wrong password | Generic AUTH-007 error in Thai; no token stored |
| Duplicate email register | Duplicate-email error in Thai |
| Toggle clears form | Credentials don't bleed across modes |

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | BA (Claude Opus 4.7 1M, dry-run #2) | Created NEW for dry-run #2 from Frontend Spec §3. |
