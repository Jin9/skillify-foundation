# POST /api/v1/inventory/stock/read

## Summary

Single-SKU stock read. Returns the live `availableQty / reservedQty / soldQty` for one SKU. Used internally by other backend services (catalog detail, checkout preview fallback, admin tooling) over the private network plane. Customer-facing reads go through `inventory.stock.bulk-read` for efficiency.

## Story refs

- `STORY_INVENTORY_RESERVE_AND_CONVERT` (read view of the state machine)
- `STORY_INVENTORY_ADMIN_ADJUST` (admin needs to verify post-adjust qty)

## Contract ref

[`inventory.stock.read`](../../architecture/contracts.json#inventory.stock.read)

## Auth

- Tier: **internal_secret**.
- Header: `X-Internal-Secret: <32-byte secret>`.
- Recipient validates with `crypto/subtle.ConstantTimeCompare` per [ADR-008](../../architecture/ADRs/ADR-008-internal-shared-secret.md). Plain `==` / `!=` is FORBIDDEN.
- Frontend MUST NOT call this endpoint and MUST NOT load `INTERNAL_SHARED_SECRET`.

## Request

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "required": ["sku"],
  "properties": {
    "sku": { "type": "string", "minLength": 1, "maxLength": 64 }
  }
}
```

## Response (success)

HTTP 200, envelope:

```json
{
  "code": "SUCCESS",
  "message": "ok",
  "data": {
    "sku": "SKU-WIDGET-001",
    "availableQty": 8,
    "reservedQty": 2,
    "soldQty": 0
  },
  "traceId": "0af7651916cd43dd8448eb211c80319c"
}
```

## Response (errors)

| Code | HTTP | Trigger |
|---|--:|---|
| `VALIDATION_ERROR` | 400 | `sku` missing, empty, or > 64 chars |
| `AUTH_INVALID` | 401 | `X-Internal-Secret` missing or mismatched (constant-time compare returned 0) |
| `NOT_FOUND` | 404 | No `stock_levels` row for the supplied `sku` |
| `DATABASE_UNAVAILABLE` | 503 | Postgres unreachable |

## Business logic steps

1. `common/middleware.InternalAuth(secret)` validates `X-Internal-Secret`. On mismatch logs `slog.Error("internal_auth_reject", ...)` (no secret value in log) and returns 401 `AUTH_INVALID`.
2. Bind + validate request body via `common/validator` (`sku` non-empty, ≤ 64).
3. `SELECT sku, available_qty, reserved_qty, sold_qty FROM stock_levels WHERE sku = $1` — point lookup on PK.
4. If `pgx.ErrNoRows` → `wrapper.Respond(c, NOT_FOUND, ...)`.
5. Wrap success envelope via `common/wrapper.Respond(c, SUCCESS, data)`.

## Side effects

None — read-only.

## Idempotency

Naturally idempotent — read endpoint, no state mutation.

## Performance

- p95 < 50ms (point lookup on PK).
- Expected QPS at MVP: 5–20 (most reads are batched via `bulk-read`).

## Test cases

- `stock_read_returns_envelope_with_seeded_zero_qtys`
- `stock_read_returns_NOT_FOUND_for_unknown_sku`
- `stock_read_returns_VALIDATION_ERROR_for_empty_sku`
- `stock_read_returns_AUTH_INVALID_when_X_Internal_Secret_missing`
- `stock_read_envelope_carries_traceId`

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | Tech-Designer (Opus 4.7 1M) | Created |
