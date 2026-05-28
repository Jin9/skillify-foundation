---
template_version: 0.1.0
story_id: STORY_FRONTEND_AUTH_GUARD
epic_id: EPIC_FRONTEND
title: Auth-required routes redirect to login and resume original route on success
as_a: Customer
i_want: deep links to auth-required routes to prompt me to log in and then take me back where I was going
so_that: I do not lose context when the app prompts me for credentials
acceptance_criteria:
  - given: An unauthenticated visitor on the cart, orders, profile, checkout, payment, orderDetail, or addresses route
    when: The route is opened without a stored access token
    then: |
      The app redirects to the login route and, after successful login, returns the
      customer to the route they originally requested (Frontend Spec §3 'Auth?'
      column).
  - given: An authenticated customer on any route
    when: The customer's access token expires and a subsequent API call returns 401
    then: |
      The frontend handles the 401, attempts to refresh via the identity refresh
      endpoint; on refresh success the customer continues on the same route; on
      refresh failure the customer is redirected to the login route (AUTH-005).
priority: Must
size: M
edge_cases:
  - case: Customer with an expired refresh token attempts a deep link into orderDetail
    expected: Refresh fails; customer is redirected to login; after re-login the original orderDetail route is resumed
  - case: Customer navigates to a public route (home/search/detail/login) without a token
    expected: No redirect; the page renders normally as a guest
  - case: Customer rapidly navigates between auth-required routes
    expected: Only one login redirect is enqueued; the final destination on login success is the most recently requested route
requirement_refs: [AUTH-005, "Frontend Spec §3"]
definition_of_ready: [default]
definition_of_done: [default]
change_log:
  - date: 2026-05-08
    author: BA (Claude Opus 4.7 1M, dry-run #2)
    change: Created NEW for dry-run #2; the Auth? column from Frontend Spec §3 needed an explicit guard story so Tech-Designer doesn't invent the policy.
---

# STORY_FRONTEND_AUTH_GUARD — Auth-required routes redirect to login and resume original route on success

## User narrative

As a **Customer**, I want **deep links to auth-required routes to prompt me to log in and then take me back where I was going** so that **I do not lose context when the app prompts me for credentials**.

## Why this story exists in EPIC_FRONTEND

The Frontend Spec §3 marks 8 routes as `Auth?: required` (orders, profile, checkout, payment, orderSuccess, orderResult, orderDetail, addresses). Without a guard story, every Tech-Designer would invent their own redirect policy.

## Edge cases

| Case | Expected |
|---|---|
| Expired refresh token | Redirect to login; resume after re-login |
| Public route without token | No redirect; renders as guest |
| Rapid navigation | One redirect; final destination is most-recent |

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | BA (Claude Opus 4.7 1M, dry-run #2) | Created NEW for dry-run #2 to anchor the Auth? guard policy. |
