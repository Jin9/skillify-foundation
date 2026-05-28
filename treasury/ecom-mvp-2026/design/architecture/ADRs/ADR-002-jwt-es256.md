# ADR-002 — JWT ES256 for end-user authentication

- **Status:** Accepted (orchestrator-locked)
- **Date:** 2026-05-08
- **Deciders:** Tech-Lead (Claude Opus 4.7 1M, dry-run #2)

## Context

End-user (customer + admin) authentication needs a stateless transport that every backend service can validate locally without hitting Identity. The orchestrator locked JWT signed with ES256 via `common/token`. The remaining design surface is the claim shape, the access/refresh split, the rotation rule, and the cookie posture for the browser.

## Decision

- **Signing:** ES256 via `common/token`. Identity holds the signing key; every other backend holds only the public key (verifier-only) and validates signature + iss + aud + exp + tokenType on every request via `common/middleware` JWT middleware.
- **Tokens:** access JWT (short-lived, 15min suggested) + refresh JWT (long-lived, 14d suggested; `tokenType=REFRESH`). Access carries `sub`, `role` (CUSTOMER|ADMIN), `iat`, `exp`, `iss="shoppilot-identity"`, `aud="shoppilot-api"`, `jti`, `tokenType="ACCESS"`. Refresh carries the same shape with `tokenType="REFRESH"`.
- **Rotation:** every successful `identity.refresh` mints a NEW access AND a NEW (rotated) refresh token; the old refresh-jti is revoked atomically in the same tx. Reuse of an old refresh-jti outside the 5-second soft-idempotency window returns `AUTH_REVOKED`.
- **Logout:** writes the refresh-jti into a revoked store (Postgres in MVP — Redis in v2) with TTL = remaining exp.
- **Browser posture:** access cookie + refresh cookie are both `HttpOnly`, `Secure`, `SameSite=Lax`, `Path=/`. The frontend SPA NEVER sees the JWT bytes; auth state is inferred from a non-sensitive `auth-hint` cookie or a `/me` probe. Detailed 6-step `withAuth` flow lives in `cross-cutting.auth.token_acquisition_flow`.
- **Email NOT in claims:** Identity does not put `email` into the access JWT extras. Reasoning: minimizes blast radius of token leakage; preserves AUTH-007 generic-error semantics; the buyerEmail dataflow (REV-L2-001) goes via `identity.profile.read` instead. See ADR-005.

## Consequences

- **Positive:** every backend validates locally; Identity is not a single point of failure on the request path.
- **Positive:** access expiry is short enough that a stolen token has limited usefulness; rotation revokes refresh-jti so a stolen refresh is one-shot.
- **Negative:** key rotation requires a coordinated public-key push to all 7 backends (operationally non-trivial; deferred to a v2 ADR).
- **Negative:** revoked-jti store is a hot read on every refresh; in PRD this should move to Redis with PG fallback.
