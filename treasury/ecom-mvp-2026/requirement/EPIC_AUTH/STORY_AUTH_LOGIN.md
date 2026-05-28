---
template_version: 0.1.0
story_id: STORY_AUTH_LOGIN
epic_id: EPIC_AUTH
title: User logs in with email and password
as_a: Customer
i_want: to log in with my email and password
so_that: I can access cart, checkout, orders, and profile
acceptance_criteria:
  - given: A registered user with the correct password
    when: The client POSTs identity login with email and the correct password
    then: |
      Response is 200 with code=SUCCESS, body contains an access token and a refresh
      token, the access token decodes to ES256 with shorter expiry than the refresh
      token, and both expiries are present (AUTH-002, AUTH-005).
  - given: A registered user but the request carries the wrong password
    when: The client POSTs identity login
    then: |
      Response is 401 with code=AUTH_INVALID and the SAME human-readable message as
      the unknown-email path (AUTH-007). No tokens are returned.
  - given: No user is registered with the supplied email
    when: The client POSTs identity login with an unknown email
    then: |
      Response is 401 with the EXACT same code=AUTH_INVALID and the same human-readable
      message as the wrong-password path. The response shape MUST NOT differ between the
      unknown-email and wrong-password cases (AUTH-007).
priority: Must
size: M
edge_cases:
  - case: Empty email or empty password in the request body
    expected: 400 VALIDATION_FAILED with field-level error; no token issued
  - case: Email with leading/trailing whitespace ("  alice@example.com  ")
    expected: Email is normalised (trimmed + lowercased) before lookup; login succeeds if otherwise valid
  - case: User account exists but status is SUSPENDED
    expected: 401 with the documented code; no tokens issued
requirement_refs: [AUTH-002, AUTH-005, AUTH-006, AUTH-007]
definition_of_ready: [default]
definition_of_done: [default]
change_log:
  - date: 2026-05-08
    author: BA (Claude Opus 4.7 1M, dry-run #2)
    change: Created from requirement §8.1 AUTH-002 + AUTH-005 + AUTH-007.
---

# STORY_AUTH_LOGIN — User logs in with email and password

## User narrative

As a **Customer**, I want **to log in with my email and password** so that **I can access cart, checkout, orders, and profile**.

## Why this story exists in EPIC_AUTH

Login is the gateway to every authenticated route in the customer mobile webview. Frontend Spec §3 marks orders / profile / checkout / payment / orderDetail / addresses as `Auth?: required`, all of which depend on this story producing a valid access + refresh token pair.

## Edge cases

| Case | Expected |
|---|---|
| Empty email or password | 400 VALIDATION_FAILED with field-level error |
| Email with whitespace padding | Normalised before lookup; login succeeds if otherwise valid |
| User exists but status=SUSPENDED | 401 with the documented code; no tokens issued |

## Implementation notes (advisory)

- Reviewer-L1 should diff the wrong-password and unknown-email responses byte-for-byte at the envelope `code` and `message` level. AUTH-007 is enumeration-resistance, not just "return some error".
- Tech-Designer chooses the access/refresh TTLs; BA only requires that access < refresh and both are non-zero.

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | BA (Claude Opus 4.7 1M, dry-run #2) | Created from requirement §8.1 AUTH-002 + AUTH-005 + AUTH-007. |
