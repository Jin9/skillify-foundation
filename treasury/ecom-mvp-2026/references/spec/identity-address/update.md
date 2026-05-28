# POST /api/v1/identity/address/update

## Summary

Partially updates a specific delivery address owned by the authenticated user. `addressId` identifies the target row; at least one additional mutable field must be provided (PATCH semantics). A dynamic `SET` clause is built from only the supplied fields to avoid overwriting untouched columns. Ownership and soft-delete status are enforced via `WHERE user_id = $userId AND deleted_at IS NULL`; a zero-row update returns 404.

## Story refs

- `STORY_AUTH_ADDRESS_ADD` (AC-15)

## Contract ref

[`identity.address.update`](../../architecture/contracts.json) — `contract_version: 0.1.0`

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
  "minProperties": 2,
  "properties": {
    "addressId":    { "type": "string", "format": "uuid" },
    "receiverName": { "type": "string", "minLength": 1, "maxLength": 100 },
    "phone":        { "type": "string", "pattern": "^[0-9+\\- ]{7,20}$" },
    "addressLine1": { "type": "string", "minLength": 1, "maxLength": 255 },
    "addressLine2": { "type": "string", "maxLength": 255 },
    "province":     { "type": "string", "minLength": 1, "maxLength": 100 },
    "district":     { "type": "string", "minLength": 1, "maxLength": 100 },
    "subDistrict":  { "type": "string", "maxLength": 100 },
    "postalCode":   { "type": "string", "pattern": "^[0-9]{5}$" }
  },
  "description": "addressId required. At least one other field must also be present (minProperties: 2)."
}
```

## Response (success)

HTTP 200:

```json
{
  "code": "SUCCESS",
  "message": "Address updated",
  "data": {
    "addressId":    "01935b9c-abcd-7000-0000-000000000099",
    "receiverName": "Alice Updated",
    "addressLine1": "456 New St",
    "province":     "Chiang Mai",
    "isDefault":    false,
    "updatedAt":    "2026-05-08T10:30:00Z"
  },
  "traceId": "0af7651916cd43dd8448eb211c80319c"
}
```

## Response (errors)

| Code | HTTP | Trigger |
|---|--:|---|
| `AUTH_MISSING`     | 401 | Access token absent from `Authorization` header |
| `AUTH_INVALID`     | 401 | Access token signature invalid, expired, iss/aud mismatch |
| `VALIDATION_ERROR` | 400 | `addressId` missing or invalid UUID; only `addressId` provided (no mutable fields); or any field violates its constraint |
| `NOT_FOUND`        | 404 | `addressId` not found, soft-deleted, or belongs to a different user (`RowsAffected == 0`) |

## Business logic steps

1. `common/middleware.Auth` validates Bearer access token; injects `user_id` into context.
2. Bind and validate request (`common/validator`); return `VALIDATION_ERROR` if `addressId` is missing/invalid or no mutable fields are provided.
3. Build a dynamic `SET` clause from only the fields explicitly present in the request body. Always append `updated_at = NOW()`.
4. `UPDATE addresses SET <dynamic_fields>, updated_at = NOW() WHERE id = $addressId AND user_id = $userId AND deleted_at IS NULL`. No pre-flight `SELECT` — rely on `RowsAffected`.
5. If `RowsAffected == 0`: return `NOT_FOUND` (404). This covers not-found, soft-deleted, and wrong-user cases uniformly.
6. If `RowsAffected == 1`: return updated fields via `common/wrapper.Success(ctx, data)`.

## Side effects

- `UPDATE` one row in `identity.addresses` (only provided fields + `updated_at`).
- No outbox event.

## Idempotency

Idempotent for the same payload: repeating the update with identical values produces the same result. `updated_at` is still bumped on each call.

## Performance

- p95 < 50ms (single UPDATE with PK lookup; no bcrypt).
- Expected QPS at MVP: 1–5.

## Test cases

- `address_update_with_valid_fields_succeeds`
- `address_update_only_addressId_returns_VALIDATION_ERROR`
- `address_update_invalid_addressId_uuid_returns_VALIDATION_ERROR`
- `address_update_not_found_returns_NOT_FOUND`
- `address_update_soft_deleted_address_returns_NOT_FOUND`
- `address_update_other_users_address_returns_NOT_FOUND`
- `address_update_invalid_phone_returns_VALIDATION_ERROR`
- `address_update_with_missing_access_token_returns_AUTH_MISSING`
- `address_update_with_invalid_access_token_returns_AUTH_INVALID`

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | Tech-Designer (Sonnet 4.6, dry-run #2) | Created |
