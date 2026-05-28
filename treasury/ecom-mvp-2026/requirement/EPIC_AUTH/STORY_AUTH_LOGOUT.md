---
template_version: 0.1.0
story_id: STORY_AUTH_LOGOUT
epic_id: EPIC_AUTH
title: User logs out and refresh token is revoked
as_a: Customer
i_want: to log out
so_that: my refresh token cannot be used to mint new access tokens after I leave
acceptance_criteria:
  - given: A logged-in user holding a valid refresh token
    when: The client POSTs identity logout with that refresh token, then re-POSTs identity refresh with the same refresh token
    then: |
      The first call (logout) succeeds with code=SUCCESS; the second call (refresh) is
      rejected with UNAUTHORIZED/INVALID_TOKEN and no new access token is issued
      (AUTH-004).
priority: Should
size: S
edge_cases:
  - case: Logout called with an already-revoked refresh token
    expected: Idempotent — second logout returns success without raising an error; no state change
  - case: Logout called without a refresh token in the body
    expected: 400 VALIDATION_FAILED; no state change
requirement_refs: [AUTH-004]
definition_of_ready: [default]
definition_of_done: [default]
change_log:
  - date: 2026-05-08
    author: BA (Claude Opus 4.7 1M, dry-run #2)
    change: Created from requirement §8.1 AUTH-004.
---

# STORY_AUTH_LOGOUT — User logs out and refresh token is revoked

## User narrative

As a **Customer**, I want **to log out** so that **my refresh token cannot be used to mint new access tokens after I leave**.

## Why this story exists in EPIC_AUTH

Logout is the only mechanism in MVP for invalidating a refresh token. Without it, a stolen refresh token would mint access tokens for its full TTL.

## Edge cases

| Case | Expected |
|---|---|
| Logout called with an already-revoked refresh token | Idempotent; second logout returns success with no state change |
| Logout called without a refresh token | 400 VALIDATION_FAILED; no state change |

## Implementation notes (advisory)

- Tech-Designer chooses the revocation mechanism (deny-list, rotation, etc.).
- The frontend profile tab's logout action calls this endpoint and then clears the locally cached access token.

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | BA (Claude Opus 4.7 1M, dry-run #2) | Created from requirement §8.1 AUTH-004. |
