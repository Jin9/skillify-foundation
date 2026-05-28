---
template_version: 0.1.0
story_id: STORY_AUTH_PROFILE
epic_id: EPIC_AUTH
title: Customer reads and edits own profile
as_a: Customer
i_want: to read my profile and edit my name and phone
so_that: I can keep my contact details current without changing my immutable email
acceptance_criteria:
  - given: An authenticated customer
    when: The client GETs identity profile
    then: |
      Response is 200 with code=SUCCESS and data contains name, email, phone, and
      defaultAddress (CUST-001).
  - given: An authenticated customer
    when: The client PATCHes identity profile with new name and phone, and separately PATCHes identity profile with a new email
    then: |
      The name+phone PATCH succeeds with code=SUCCESS and the persisted user reflects
      the new values; the email PATCH is rejected with VALIDATION_ERROR and the
      persisted email is unchanged (CUST-002).
priority: Must
size: S
edge_cases:
  - case: PATCH with a phone number containing invalid characters
    expected: 400 VALIDATION_FAILED citing the phone format rule
  - case: PATCH with a name that is empty after trimming
    expected: 400 VALIDATION_FAILED on the name field
requirement_refs: [CUST-001, CUST-002]
definition_of_ready: [default]
definition_of_done: [default]
change_log:
  - date: 2026-05-08
    author: BA (Claude Opus 4.7 1M, dry-run #2)
    change: Created from requirement §8.2 CUST-001 + CUST-002.
---

# STORY_AUTH_PROFILE — Customer reads and edits own profile

## User narrative

As a **Customer**, I want **to read my profile and edit my name and phone** so that **I can keep my contact details current without changing my immutable email**.

## Why this story exists in EPIC_AUTH

The mobile webview's profile tab (Frontend Spec §3) is built on this surface. CUST-002 locks email immutability in MVP — this is what protects uniqueness without a separate identity-merge flow.

## Edge cases

| Case | Expected |
|---|---|
| PATCH with malformed phone | 400 VALIDATION_FAILED on the phone field |
| PATCH with empty name | 400 VALIDATION_FAILED on the name field |

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | BA (Claude Opus 4.7 1M, dry-run #2) | Created from requirement §8.2 CUST-001 + CUST-002. |
