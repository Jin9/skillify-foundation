# POST /api/v1/identity/address/create

## Summary

Creates a new delivery address for the authenticated user. When `isDefault = true`, the operation runs in a Postgres transaction that clears `is_default` on all existing live addresses before inserting the new one, ensuring the partial UNIQUE index invariant `UNIQUE(user_id) WHERE is_default = true AND deleted_at IS NULL` is never violated. When `isDefault = false` (default), a plain `INSERT` is performed without a transaction.

## Story refs

- `STORY_AUTH_ADDRESS_ADD` (AC-13)

## Contract ref

[`identity.address.create`](../../architecture/contracts.json) — `contract_version: 0.1.0`

## Auth

- Tier: **Bearer JWT** (access token required in `Authorization: Bearer <token>` header).
- `common/middleware.Auth` validates the token and injects `user_id` into request context.

## Request

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "required": ["receiverName", "phone", "addressLine1", "province", "district", "postalCode"],
  "properties": {
    "receiverName": { "type": "string", "minLength": 1, "maxLength": 100 },
    "phone":        { "type": "string", "pattern": "^[0-9+\\- ]{7,20}$" },
    "addressLine1": { "type": "string", "minLength": 1, "maxLength": 255 },
    "addressLine2": { "type": "string", "maxLength": 255 },
    "province":     { "type": "string", "minLength": 1, "maxLength": 100 },
    "district":     { "type": "string", "minLength": 1, "maxLength": 100 },
    "subDistrict":  { "type": "string", "maxLength": 100 },
    "postalCode":   { "type": "string", "pattern": "^[0-9]{5}$" },
    "isDefault":    { "type": "boolean", "default": false }
  }
}
```

## Response (success)

HTTP 200:

```json
{
  "code": "SUCCESS",
  "message": "Address created",
  "data": {
    "addressId":    "01935b9c-abcd-7000-0000-000000000099",
    "receiverName": "Alice",
    "addressLine1": "123 Main St",
    "province":     "Bangkok",
    "isDefault":    true
  },
  "traceId": "0af7651916cd43dd8448eb211c80319c"
}
```

## Response (errors)

| Code | HTTP | Trigger |
|---|--:|---|
| `AUTH_MISSING`     | 401 | Access token absent from `Authorization` header |
| `AUTH_INVALID`     | 401 | Access token signature invalid, expired, iss/aud mismatch |
| `VALIDATION_ERROR` | 400 | Any required field missing, `phone` pattern mismatch, `postalCode` not 5 digits, or field exceeds max length |

## Business logic steps

1. `common/middleware.Auth` validates Bearer access token; injects `user_id` into context.
2. Bind and validate request (`common/validator`); return `VALIDATION_ERROR` on any constraint violation.
3. Generate new `address_id = uuid.NewV7()`.
4. **If `isDefault == true`:**
   - Open `pgx` transaction.
   - `UPDATE addresses SET is_default = false, updated_at = NOW() WHERE user_id = $userId AND deleted_at IS NULL`.
   - `INSERT INTO addresses(id, user_id, receiver_name, phone, address_line1, address_line2, province, district, sub_district, postal_code, is_default = true, created_at = NOW(), updated_at = NOW())`.
   - `COMMIT`.
5. **If `isDefault == false`:**
   - Single `INSERT INTO addresses(...)` with `is_default = false` (no transaction needed).
6. Return `common/wrapper.Success(ctx, data)` with new `addressId`, `receiverName`, `addressLine1`, `province`, `isDefault`.

## Side effects

- `INSERT` one row into `identity.addresses`.
- If `isDefault = true`: `UPDATE` all existing live addresses for the user to clear `is_default` (within same transaction).
- Partial UNIQUE index `UNIQUE(user_id) WHERE is_default = true AND deleted_at IS NULL` acts as DB-level safety net.
- No outbox event.

## Idempotency

Not idempotent — each call creates a new address row with a new UUID. Duplicate submissions produce multiple address rows.

## Performance

- p95 < 50ms (single INSERT or short transaction; no bcrypt).
- Expected QPS at MVP: 1–5.

## Test cases

- `address_create_with_valid_payload_succeeds`
- `address_create_with_is_default_true_clears_previous_default`
- `address_create_with_is_default_false_does_not_affect_existing_defaults`
- `address_create_missing_required_field_returns_VALIDATION_ERROR`
- `address_create_invalid_phone_returns_VALIDATION_ERROR`
- `address_create_invalid_postal_code_returns_VALIDATION_ERROR`
- `address_create_with_missing_access_token_returns_AUTH_MISSING`
- `address_create_with_invalid_access_token_returns_AUTH_INVALID`

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | Tech-Designer (Sonnet 4.6, dry-run #2) | Created |
