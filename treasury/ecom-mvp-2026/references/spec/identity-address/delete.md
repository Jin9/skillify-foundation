# POST /api/v1/identity/address/delete

## Summary

Soft-deletes a delivery address owned by the authenticated user by setting `deleted_at = NOW()` and clearing `is_default = false` on the target row (ID-DEC-3). No hard delete is ever performed — the row is preserved for audit. The operation is idempotent: re-deleting an already soft-deleted address returns 200. Historical orders are unaffected because they snapshot the address at `checkout.commit` time. The `is_default` flag is explicitly cleared to preserve the partial UNIQUE index invariant.

## Story refs

- `STORY_AUTH_ADDRESS_ADD` (AC-16, AC-17)

## Contract ref

[`identity.address.delete`](../../architecture/contracts.json) — `contract_version: 0.1.0`

## Auth

- Tier: **Bearer JWT** (access token required in `Authorization: Bearer <token>` header).
- `common/middleware.Auth` validates the token and injects `user_id` into request context.

## Request

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "required": ["addressId"],
  "properties": {
    "addressId": { "type": "string", "format": "uuid" }
  }
}
```

## Response (success)

HTTP 200:

```json
{
  "code": "SUCCESS",
  "message": "Address deleted",
  "data": null,
  "traceId": "0af7651916cd43dd8448eb211c80319c"
}
```

## Response (errors)

| Code | HTTP | Trigger |
|---|--:|---|
| `AUTH_MISSING`     | 401 | Access token absent from `Authorization` header |
| `AUTH_INVALID`     | 401 | Access token signature invalid, expired, iss/aud mismatch |
| `VALIDATION_ERROR` | 400 | `addressId` missing or not a valid UUID |
| `NOT_FOUND`        | 404 | `addressId` truly absent from `addresses` table (never existed or belongs to different user — the `never-existed` case only) |

## Business logic steps

1. `common/middleware.Auth` validates Bearer access token; injects `user_id` into context.
2. Bind and validate request (`common/validator`); return `VALIDATION_ERROR` if `addressId` is missing or not a valid UUID.
3. `UPDATE addresses SET deleted_at = NOW(), is_default = false, updated_at = NOW() WHERE id = $addressId AND user_id = $userId AND deleted_at IS NULL`.
4. If `RowsAffected == 1`: return 200 SUCCESS with null data.
5. If `RowsAffected == 0`: run `SELECT EXISTS(SELECT 1 FROM addresses WHERE id = $addressId AND user_id = $userId)`.
   - If EXISTS is true: the row exists but is already soft-deleted → return 200 SUCCESS (idempotent, ID-DEC-3).
   - If EXISTS is false: the row never existed or belongs to another user → return `NOT_FOUND` (404).
6. Return `common/wrapper.Success(ctx, nil)`.

## Side effects

- `UPDATE` one row in `identity.addresses`: sets `deleted_at = NOW()`, `is_default = false`, `updated_at = NOW()`.
- Clears `is_default` explicitly so the partial UNIQUE index `UNIQUE(user_id) WHERE is_default = true AND deleted_at IS NULL` is not violated by future `set-default` calls on remaining live addresses.
- Historical orders retain their snapshotted address data — no cascade effect.
- No outbox event.

## Idempotency

Idempotent (ID-DEC-3). Re-deleting an already soft-deleted address owned by the same user returns 200 SUCCESS. Simplifies client retry logic on network failures.

## Performance

- p95 < 50ms (single UPDATE + conditional SELECT EXISTS; no bcrypt).
- Expected QPS at MVP: 1–5.

## Test cases

- `address_delete_sets_deleted_at`
- `address_delete_clears_is_default`
- `delete_already_deleted_address_returns_200_idempotent`
- `address_delete_never_existed_returns_NOT_FOUND`
- `address_delete_other_users_address_returns_NOT_FOUND`
- `address_delete_invalid_addressId_uuid_returns_VALIDATION_ERROR`
- `address_delete_with_missing_access_token_returns_AUTH_MISSING`
- `address_delete_with_invalid_access_token_returns_AUTH_INVALID`

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | Tech-Designer (Sonnet 4.6, dry-run #2) | Created; ID-DEC-3 soft-delete pattern with idempotent re-delete |
