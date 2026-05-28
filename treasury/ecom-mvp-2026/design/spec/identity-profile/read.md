# POST /api/v1/identity/profile/read

## Summary

Returns the authenticated user's own profile, including `defaultAddressId` for checkout pre-fill. `user_id` is extracted from the Bearer access token — no request body is needed. This is a first-class endpoint (not a stub), closed per REV-L2-001. Email is intentionally kept out of the JWT claims (AUTH-007 blast-radius protection); callers such as `checkout.commit` use this endpoint to retrieve `email` for the `buyerEmail` field on order creation (ID-DEC-5).

## Story refs

- `STORY_AUTH_PROFILE` (AC-10, AC-11)

## Contract ref

[`identity.profile.read`](../../architecture/contracts.json) — `contract_version: 0.1.0`

## Auth

- Tier: **Bearer JWT** (access token required in `Authorization: Bearer <token>` header).
- `common/middleware.Auth` validates the token and injects `user_id` (JWT `sub`) into request context.

## Request

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {},
  "description": "No request body. user_id is extracted from the Bearer access token claims (claims.sub)."
}
```

## Response (success)

HTTP 200:

```json
{
  "code": "SUCCESS",
  "message": "Profile fetched",
  "data": {
    "userId":           "01935b9c-1234-7000-abcd-000000000001",
    "email":            "alice@example.com",
    "name":             "Alice",
    "phone":            "+66812345678",
    "defaultAddressId": "01935b9c-5678-7000-abcd-000000000002",
    "role":             "CUSTOMER",
    "status":           "ACTIVE",
    "createdAt":        "2026-01-01T08:00:00Z",
    "updatedAt":        "2026-05-01T12:00:00Z"
  },
  "traceId": "0af7651916cd43dd8448eb211c80319c"
}
```

`phone` and `defaultAddressId` are `null` when not set.

## Response (errors)

| Code | HTTP | Trigger |
|---|--:|---|
| `AUTH_MISSING` | 401 | Access token absent from `Authorization` header |
| `AUTH_INVALID` | 401 | Access token signature invalid, expired, iss/aud mismatch |
| `NOT_FOUND`    | 404 | User record not found (edge case: account deleted after token issuance) |

## Business logic steps

1. `common/middleware.Auth` validates the Bearer access token; returns `AUTH_MISSING` or `AUTH_INVALID` on failure; injects `user_id` into context.
2. Query: `SELECT u.id, u.email, u.name, u.phone, u.role, u.status, u.created_at, u.updated_at, a.id AS default_address_id FROM users u LEFT JOIN addresses a ON a.user_id = u.id AND a.is_default = true AND a.deleted_at IS NULL WHERE u.id = $userId`.
3. If no user row returned: return `NOT_FOUND` (404).
4. Map row to response DTO; `password_hash` is NEVER included in the response.
5. Return `common/wrapper.Success(ctx, data)`.

## Side effects

None — read-only query; no state mutation.

## Idempotency

Fully idempotent — repeated calls with the same valid token return identical data (aside from `updatedAt` if a concurrent profile update occurs).

## Performance

- p95 < 50ms (single SELECT + optional LEFT JOIN; no bcrypt).
- Expected QPS at MVP: 10–30 (called by profile page and by `checkout.commit` step 2).

## Test cases

- `profile_read_returns_self_with_default_address_id`
- `profile_read_by_checkout_returns_email`
- `profile_read_phone_is_null_when_not_set`
- `profile_read_default_address_id_is_null_when_no_default`
- `profile_read_with_missing_access_token_returns_AUTH_MISSING`
- `profile_read_with_invalid_access_token_returns_AUTH_INVALID`
- `profile_read_for_deleted_user_returns_NOT_FOUND`
- `profile_read_never_returns_password_hash`

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | Tech-Designer (Sonnet 4.6, dry-run #2) | Created; REV-L2-001 closure — promoted from stub to first-class endpoint |
