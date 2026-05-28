# POST /api/v1/identity/address/list

## Summary

Returns all non-deleted delivery addresses for the authenticated user, ordered with the default address first and then by `created_at` ascending. Returns an empty array (not 404) when no addresses exist. `user_id` is extracted from the Bearer access token — no request body is needed.

## Story refs

- `STORY_AUTH_ADDRESS_ADD` (AC-14)

## Contract ref

[`identity.address.list`](../../architecture/contracts.json) — `contract_version: 0.1.0`

## Auth

- Tier: **Bearer JWT** (access token required in `Authorization: Bearer <token>` header).
- `common/middleware.Auth` validates the token and injects `user_id` into request context.

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
  "message": "Addresses fetched",
  "data": [
    {
      "addressId":    "01935b9c-abcd-7000-0000-000000000099",
      "receiverName": "Alice",
      "phone":        "+66812345678",
      "addressLine1": "123 Main St",
      "addressLine2": null,
      "province":     "Bangkok",
      "district":     "Pathum Wan",
      "subDistrict":  null,
      "postalCode":   "10330",
      "isDefault":    true,
      "createdAt":    "2026-01-15T09:00:00Z"
    }
  ],
  "traceId": "0af7651916cd43dd8448eb211c80319c"
}
```

`data` is an empty array `[]` when the user has no addresses. `addressLine2` and `subDistrict` are `null` when not set.

## Response (errors)

| Code | HTTP | Trigger |
|---|--:|---|
| `AUTH_MISSING` | 401 | Access token absent from `Authorization` header |
| `AUTH_INVALID` | 401 | Access token signature invalid, expired, iss/aud mismatch |

## Business logic steps

1. `common/middleware.Auth` validates Bearer access token; returns `AUTH_MISSING` or `AUTH_INVALID` on failure; injects `user_id` into context.
2. `SELECT id, receiver_name, phone, address_line1, address_line2, province, district, sub_district, postal_code, is_default, created_at FROM addresses WHERE user_id = $userId AND deleted_at IS NULL ORDER BY is_default DESC, created_at ASC`.
3. Map rows to response DTO array.
4. Return `common/wrapper.Success(ctx, data)` with the array (empty array if no rows — never 404).

## Side effects

None — read-only query; no state mutation.

## Idempotency

Fully idempotent — repeated calls with the same valid token return identical data (aside from real-time changes from concurrent add/update/delete operations).

## Performance

- p95 < 50ms (single SELECT with filter index `INDEX(user_id) WHERE deleted_at IS NULL`; no bcrypt).
- Expected QPS at MVP: 5–15.

## Test cases

- `address_list_returns_only_non_deleted`
- `address_list_returns_empty_array_when_no_addresses`
- `address_list_default_address_appears_first`
- `address_list_does_not_return_soft_deleted_addresses`
- `address_list_with_missing_access_token_returns_AUTH_MISSING`
- `address_list_with_invalid_access_token_returns_AUTH_INVALID`

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | Tech-Designer (Sonnet 4.6, dry-run #2) | Created |
