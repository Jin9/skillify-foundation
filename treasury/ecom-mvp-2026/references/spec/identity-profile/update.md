# POST /api/v1/identity/profile/update

## Summary

Partially updates the authenticated user's mutable profile fields (`name` and/or `phone`). At least one field must be provided. Email and `role` are not updatable through this endpoint. A dynamic `SET` clause is built from only the provided fields to avoid overwriting untouched columns. `updated_at` is always bumped. Returns the updated profile fields on success.

## Story refs

- `STORY_AUTH_PROFILE` (AC-12)

## Contract ref

[`identity.profile.update`](../../architecture/contracts.json) — `contract_version: 0.1.0`

## Auth

- Tier: **Bearer JWT** (access token required in `Authorization: Bearer <token>` header).
- `common/middleware.Auth` validates the token and injects `user_id` into request context.

## Request

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "minProperties": 1,
  "properties": {
    "name":  { "type": "string", "minLength": 1, "maxLength": 120 },
    "phone": { "type": "string", "pattern": "^[0-9+\\- ]{7,20}$" }
  },
  "description": "At least one field required. email and role are NOT updatable here."
}
```

## Response (success)

HTTP 200:

```json
{
  "code": "SUCCESS",
  "message": "Profile updated",
  "data": {
    "userId":    "01935b9c-1234-7000-abcd-000000000001",
    "email":     "alice@example.com",
    "name":      "Alice Updated",
    "phone":     "+66812345678",
    "updatedAt": "2026-05-08T10:00:00Z"
  },
  "traceId": "0af7651916cd43dd8448eb211c80319c"
}
```

## Response (errors)

| Code | HTTP | Trigger |
|---|--:|---|
| `AUTH_MISSING`     | 401 | Access token absent from `Authorization` header |
| `AUTH_INVALID`     | 401 | Access token signature invalid, expired, iss/aud mismatch |
| `VALIDATION_ERROR` | 400 | Empty request body, neither `name` nor `phone` provided, or any field violates its constraint |

## Business logic steps

1. `common/middleware.Auth` validates Bearer access token; returns `AUTH_MISSING` or `AUTH_INVALID` on failure; injects `user_id` into context.
2. Bind and validate request (`common/validator`); return `VALIDATION_ERROR` if body is empty or fields violate constraints.
3. Verify at least one of `name` or `phone` is present in the parsed body; return `VALIDATION_ERROR` if neither is provided.
4. Build a dynamic `SET` clause including only the fields that were explicitly provided. Always append `updated_at = NOW()`.
5. `UPDATE users SET <dynamic_fields>, updated_at = NOW() WHERE id = $userId`.
6. Return the updated row fields (`userId`, `email`, `name`, `phone`, `updatedAt`) via `common/wrapper.Success(ctx, data)`.

## Side effects

- `UPDATE` one row in `identity.users` (bumps `updated_at`; sets only provided fields).
- No outbox event.

## Idempotency

Idempotent for the same payload: repeating the same update produces the same final values. `updated_at` is still bumped on each call, but the observable data fields are unchanged.

## Performance

- p95 < 50ms (single UPDATE; no bcrypt).
- Expected QPS at MVP: 1–5.

## Test cases

- `profile_update_name_and_phone_succeeds`
- `profile_update_name_only_does_not_overwrite_phone`
- `profile_update_phone_only_does_not_overwrite_name`
- `profile_update_empty_body_returns_VALIDATION_ERROR`
- `profile_update_invalid_phone_pattern_returns_VALIDATION_ERROR`
- `profile_update_name_exceeds_max_length_returns_VALIDATION_ERROR`
- `profile_update_with_missing_access_token_returns_AUTH_MISSING`
- `profile_update_with_invalid_access_token_returns_AUTH_INVALID`

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | Tech-Designer (Sonnet 4.6, dry-run #2) | Created |
