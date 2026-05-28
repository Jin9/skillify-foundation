---
template_version: 0.1.0
epic_id: EPIC_AUTH
title: Customer identity, session, and address management
summary: |
  Customers must register, log in, refresh sessions, log out, manage profile fields,
  and maintain shipping addresses used by checkout. This is the gateway to every
  authenticated path in the customer journey (cart, checkout, orders, profile).
business_value: Unblocks every authenticated path in the platform. Without a stable identity surface, the customer mobile webview cannot reach the cart, checkout, payment, or order routes.
in_scope_stories:
  - STORY_AUTH_REGISTER
  - STORY_AUTH_LOGIN
  - STORY_AUTH_LOGOUT
  - STORY_AUTH_PROFILE
  - STORY_AUTH_ADDRESS_ADD
out_of_scope:
  - Social login (Google/Facebook/Apple) — deferred per requirement §4.2
  - Multi-factor authentication — deferred per requirement §4.2
  - Password reset via email — deferred (no notification infrastructure in MVP, §8.14)
  - Account suspension self-service — admin-only field in §12.1
success_metrics:
  - Login success rate >= 99% on valid credentials with the correct password
  - Generic login error message identical for unknown-email and wrong-password paths (AUTH-007)
  - Refresh-token revocation effective on the very next refresh attempt (AUTH-004)
  - p95 login latency < 300 ms (PERF-001 family budget)
dependencies: []
compliance_sensitive: false
change_log:
  - date: 2026-05-08
    author: BA (Claude Opus 4.7 1M, dry-run #2)
    change: Created from requirement §8.1 + §8.2; carries over from runs/01 with refactored success metrics and explicit address scope (CUST-003, CUST-004).
---

# EPIC_AUTH — Customer identity, session, and address management

## Why

Every authenticated route in the ShopPilot mobile webview — cart, orders, profile, checkout, payment, addresses, orderDetail — depends on a working identity surface. The Frontend Spec §3 marks these routes as `Auth?: required`, which means the frontend must redirect to the login route when no valid token is present.

The original requirement (§8.1, §8.2) lays down the full backbone: register, login, logout, profile read/edit, and shipping addresses with a default selection used by checkout. We are intentionally NOT shipping social login, MFA, or email-based password reset in MVP — those are explicitly deferred (§4.2 and the absence of any email/notification infrastructure per §8.14).

The single most important behavioural rule here is AUTH-007: login failure must return the same generic error code and message for unknown-email and wrong-password cases so attackers cannot enumerate registered emails through the response shape.

## Stories in this epic

| Story ID | Title | Priority |
|---|---|---|
| `STORY_AUTH_REGISTER` | Guest registers a CUSTOMER account | Must |
| `STORY_AUTH_LOGIN` | User logs in and receives access + refresh tokens | Must |
| `STORY_AUTH_LOGOUT` | User logs out and refresh token is revoked | Should |
| `STORY_AUTH_PROFILE` | Customer reads and edits own profile (name, phone) | Must |
| `STORY_AUTH_ADDRESS_ADD` | Customer adds a shipping address and selects default | Must |

## Out-of-scope rationale

- **Social login / MFA / email-based password reset** — the original requirement §4.2 explicitly defers these and §8.14 confirms there is no email infrastructure in MVP, so we cannot mail a reset link.
- **Account suspension self-service** — the User entity has a status enum (ACTIVE, SUSPENDED) per §12.1 but no story exposes a customer-facing path to mutate it; admin-only.

## Risks

- **AUTH-007 enumeration risk** — a wrong-password and unknown-email response must be byte-for-byte identical at the envelope `code` and `message` level. Reviewer-L1 should diff both responses.
- **Refresh-token rotation race** — logout must be atomic against an in-flight refresh; otherwise a revoked token could still mint a fresh access token within a few-millisecond window.
- **Email immutability** — CUST-002 forbids email edits in MVP. The profile PATCH endpoint must reject email changes with VALIDATION_ERROR; otherwise a customer could "rename" themselves around uniqueness.
- **Address ownership** — CHK-002 requires checkout to verify the chosen addressId belongs to the requesting userId. The address create/list endpoints must enforce userId scoping at the query level.

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | BA (Claude Opus 4.7 1M, dry-run #2) | Created from requirement §8.1 + §8.2; carries forward from runs/01 with refactored success metrics and explicit address scope. |
