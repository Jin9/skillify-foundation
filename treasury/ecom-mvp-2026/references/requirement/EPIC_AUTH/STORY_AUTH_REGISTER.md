---
template_version: 0.1.0
story_id: STORY_AUTH_REGISTER
epic_id: EPIC_AUTH
title: Guest registers a CUSTOMER account
as_a: Guest
i_want: to register a new account with my email, password, and name
so_that: I can log in and use authenticated features (cart, checkout, orders)
acceptance_criteria:
  - given: An email not yet registered, a password meeting policy, and a non-empty name
    when: The client POSTs identity register with those fields
    then: |
      Response is 200/201 with code=SUCCESS, no password material in the body, the
      created user has role=CUSTOMER, and the persisted password is a one-way hash
      distinct from the plaintext (AUTH-001, AUTH-006).
  - given: An email already registered to an existing user
    when: The client POSTs identity register with that same email
    then: |
      Response is 409 with a code that signals duplicate email; no second user row is
      created and the existing user is unchanged (AUTH-001).
priority: Must
size: M
edge_cases:
  - case: Email with leading/trailing whitespace ("  alice@example.com  ")
    expected: Email is normalised (trimmed + lowercased) before uniqueness check; the registration succeeds if otherwise valid
  - case: Empty name or empty password
    expected: 400 VALIDATION_FAILED with field-level error; no user row created
  - case: Password shorter than the policy minimum
    expected: 400 VALIDATION_FAILED citing the password rule that failed; no user row created
requirement_refs: [AUTH-001, AUTH-006]
definition_of_ready: [default]
definition_of_done: [default]
change_log:
  - date: 2026-05-08
    author: BA (Claude Opus 4.7 1M, dry-run #2)
    change: Created from requirement §8.1 AUTH-001 + AUTH-006.
---

# STORY_AUTH_REGISTER — Guest registers a CUSTOMER account

## User narrative

As a **Guest**, I want **to register a new account with my email, password, and name** so that **I can log in and use authenticated features (cart, checkout, orders)**.

## Why this story exists in EPIC_AUTH

Registration is the entry point for every customer in the system. Without it, the only path into the authenticated routes (cart-checkout-payment-orders) is admin-seeded user accounts, which is not a customer-acquisition story.

## Edge cases

| Case | Expected |
|---|---|
| Email with whitespace padding | Normalised (trim + lowercase) before uniqueness check; registration succeeds |
| Empty name or password | 400 VALIDATION_FAILED with field-level error |
| Weak password | 400 VALIDATION_FAILED citing the policy rule that failed |

## Implementation notes (advisory; BA does NOT prescribe)

- Tech-Designer chooses the password-hashing algorithm and cost factor.
- Tech-Designer chooses the password policy (min length, character classes); BA only requires that violations return VALIDATION_FAILED with a field-level reason.

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | BA (Claude Opus 4.7 1M, dry-run #2) | Created from requirement §8.1 AUTH-001 + AUTH-006. |
