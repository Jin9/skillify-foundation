# POST /api/v1/identity/auth/logout

## Summary

Revokes the caller's refresh token by setting `revoked_at` on the matching `refresh_tokens` row. Requires a valid Bearer access token (to authenticate the user) plus the refresh token in the request body (to identify the specific session). Cross-user logout is prevented: the jti's `user_id` must match the authenticated user. Re-revoking an already-revoked jti returns 200 (idempotent per `identity.logout` contract v0.1.0 — see AMBIG-IDENTITY-001).

## Story refs

- `STORY_AUTH_LOGOUT` (AC-09)

## Contract ref

[`identity.logout`](../../architecture/contracts.json) — `contract_version: 0.1.0`

## Auth

- Tier: **Bearer JWT** (access token required in `Authorization: Bearer <token>` header).
- `common/middleware.Auth` validates the token and injects `user_id` into request context.

## Request

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "required": ["refreshToken"],
  "properties": {
    "refreshToken": { "type": "string", "minLength": 1 }
  }
}
```

## Response (success)

HTTP 200:

```json
{
  "code": "SUCCESS",
  "message": "Logged out",
  "data": null,
  "traceId": "0af7651916cd43dd8448eb211c80319c"
}
```

## Response (errors)

| Code | HTTP | Trigger |
|---|--:|---|
| `AUTH_MISSING`     | 401 | Access token absent from `Authorization` header |
| `AUTH_INVALID`     | 401 | Access token signature invalid, expired, iss/aud mismatch |
| `VALIDATION_ERROR` | 400 | `refreshToken` missing or not a parseable JWT |
| `AUTH_REVOKED`     | 401 | jti absent from `refresh_tokens` or `user_id` in jti does not match authenticated user |

## Business logic steps

1. `common/middleware.Auth` validates the Bearer access token; returns `AUTH_MISSING` (401) or `AUTH_INVALID` (401) on failure; injects `user_id` into context.
2. Bind and validate request (`common/validator`); return `VALIDATION_ERROR` if `refreshToken` is missing or not parseable as JWT.
3. Parse `refreshToken` via `common/token.ParseES256` to extract `jti` and `sub` claims. If parse fails, return `VALIDATION_ERROR`.
4. Verify `token.Claims.Sub == authenticated user_id` to prevent cross-user session revocation. On mismatch return `AUTH_REVOKED`.
5. `UPDATE refresh_tokens SET revoked_at = NOW() WHERE jti = $jti AND user_id = $userId AND revoked_at IS NULL`.
6. If `RowsAffected == 1`: return 200 SUCCESS with null data.
7. If `RowsAffected == 0`: check whether `SELECT EXISTS(... WHERE jti = $jti AND user_id = $userId AND revoked_at IS NOT NULL)`. If already revoked: return 200 SUCCESS (idempotent per AMBIG-IDENTITY-001 resolution). If jti absent or user mismatch: return `AUTH_REVOKED`.
8. Return `common/wrapper.Success(ctx, nil)`.

## Side effects

- `UPDATE` one row in `identity.refresh_tokens` (sets `revoked_at`).
- No outbox event.
- Caller (Next.js route handler) clears `access` and `refresh` HttpOnly cookies after receiving 200.

## Idempotency

Idempotent per `identity.logout` contract v0.1.0: re-revoking an already-revoked jti for the same user returns 200 SUCCESS. Client treats any 200 as "logged out". Ambiguity AMBIG-IDENTITY-001 flagged for TL — resolution implemented here follows `contracts.json` (idempotent 200 on repeat).

## Performance

- p95 < 100ms (single UPDATE, no bcrypt).
- Expected QPS at MVP: 2–10.

## Test cases

- `logout_revokes_refresh_token`
- `logout_requires_valid_access_token`
- `logout_with_missing_access_token_returns_AUTH_MISSING`
- `logout_with_invalid_access_token_returns_AUTH_INVALID`
- `logout_with_missing_refreshToken_returns_VALIDATION_ERROR`
- `logout_already_revoked_jti_returns_200_idempotent`
- `logout_cross_user_jti_returns_AUTH_REVOKED`
- `logout_unknown_jti_returns_AUTH_REVOKED`

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | Tech-Designer (Sonnet 4.6, dry-run #2) | Created |
