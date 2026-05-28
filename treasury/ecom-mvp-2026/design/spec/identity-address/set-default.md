# POST /api/v1/identity/address/set-default

## Summary

Atomically promotes one of the authenticated user's live delivery addresses to the default. Runs in a single Postgres transaction: first clears `is_default` on all other non-deleted addresses for the user, then sets `is_default = true` on the target. If the target row is not found, soft-deleted, or belongs to a different user, the transaction is rolled back and `NOT_FOUND` is returned. The partial UNIQUE index `UNIQUE(user_id) WHERE is_default = true AND deleted_at IS NULL` provides a DB-level safety net. Setting an already-default address as default is idempotent (returns 200).

## Story refs

- `STORY_AUTH_ADDRESS_ADD` (AC-18, AC-19)

## Contract ref

[`identity.address.set-default`](../../architecture/contracts.json) — `contract_version: 0.1.0`

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
  "message": "Default address set",
  "data": {
    "addressId": "01935b9c-abcd-7000-0000-000000000099",
    "isDefault": true
  },
  "traceId": "0af7651916cd43dd8448eb211c80319c"
}
```

## Response (errors)

| Code | HTTP | Trigger |
|---|--:|---|
| `AUTH_MISSING`     | 401 | Access token absent from `Authorization` header |
| `AUTH_INVALID`     | 401 | Access token signature invalid, expired, iss/aud mismatch |
| `VALIDATION_ERROR` | 400 | `addressId` missing or not a valid UUID |
| `NOT_FOUND`        | 404 | `addressId` not found, soft-deleted, or belongs to a different user (detected when step-2 UPDATE `RowsAffected == 0`, then ROLLBACK) |

## Business logic steps

1. `common/middleware.Auth` validates Bearer access token; injects `user_id` into context.
2. Bind and validate request (`common/validator`); return `VALIDATION_ERROR` if `addressId` is missing or not a valid UUID.
3. Open `pgx` transaction.
4. **Step 1 (clear others):** `UPDATE addresses SET is_default = false, updated_at = NOW() WHERE user_id = $userId AND deleted_at IS NULL AND is_default = true AND id <> $addressId`. Safe to run even when no other defaults exist (zero rows affected is fine here).
5. **Step 2 (set target):** `UPDATE addresses SET is_default = true, updated_at = NOW() WHERE id = $addressId AND user_id = $userId AND deleted_at IS NULL`. Capture `RowsAffected`.
6. If `RowsAffected == 0`: `ROLLBACK`; return `NOT_FOUND` (404). This covers not-found, soft-deleted, and wrong-user cases.
7. If `RowsAffected == 1`: `COMMIT`.
8. Return `common/wrapper.Success(ctx, data)` with `addressId` and `isDefault: true`.

## Side effects

- Single `pgx` transaction updates up to two sets of rows in `identity.addresses`:
  - Clears `is_default` on previously-default address(es) for the user (step 1).
  - Sets `is_default = true` on the target address (step 2).
- Partial UNIQUE index `UNIQUE(user_id) WHERE is_default = true AND deleted_at IS NULL` enforces single-default invariant at DB level as a safety net.
- No outbox event.

## Idempotency

Idempotent — setting an address that is already the default returns 200 SUCCESS. Step 1 becomes a no-op (no other defaults to clear) and step 2 updates with the same value.

## Performance

- p95 < 50ms (two UPDATEs within a short Postgres transaction; no bcrypt).
- Expected QPS at MVP: 1–5.

## Test cases

- `address_create_then_set_default`
- `set_default_clears_previous_default_atomically`
- `set_default_on_already_default_address_returns_200_idempotent`
- `set_default_not_found_address_returns_NOT_FOUND`
- `set_default_soft_deleted_address_returns_NOT_FOUND`
- `set_default_other_users_address_returns_NOT_FOUND`
- `set_default_invalid_addressId_uuid_returns_VALIDATION_ERROR`
- `set_default_with_missing_access_token_returns_AUTH_MISSING`
- `set_default_with_invalid_access_token_returns_AUTH_INVALID`

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | Tech-Designer (Sonnet 4.6, dry-run #2) | Created |
