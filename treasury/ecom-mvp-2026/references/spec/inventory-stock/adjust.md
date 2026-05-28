# POST /api/v1/inventory/stock/adjust

## Summary

Admin-only stock adjustment. The admin issues a signed `delta` against an SKU's `available_qty` along with a free-form `reason`. The endpoint persists a `stock_adjustments` audit row in the same transaction as the `stock_levels` UPDATE — so an audit row never exists without its corresponding stock change, and vice-versa. Implements requirement INV-007 + INV-008 + AUTH-003.

## Story refs

- `STORY_INVENTORY_ADMIN_ADJUST`

## Contract ref

[`inventory.stock.adjust`](../../architecture/contracts.json#inventory.stock.adjust)

## Auth

- Tier: **admin_jwt**.
- Header: `Authorization: Bearer <jwt>` with `claims.role == 'ADMIN'`.
- CUSTOMER token → 403 `AUTH_FORBIDDEN` per `cross-cutting.auth.role_rejection_rules`.
- No admin UI in MVP — admin uses curl/Postman per BA `out_of_scope`.

## Request

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "required": ["sku", "delta", "reason"],
  "properties": {
    "sku":    { "type": "string",  "minLength": 1, "maxLength": 64 },
    "delta":  { "type": "integer", "not": { "const": 0 } },
    "reason": { "type": "string",  "minLength": 1, "maxLength": 255 }
  }
}
```

## Response (success)

HTTP 200, envelope:

```json
{
  "code": "SUCCESS",
  "message": "stock adjusted",
  "data": {
    "sku": "SKU-WIDGET-001",
    "availableQty": 12,
    "adjustmentId": "01935b9c-8a9e-7c3f-9d11-22aa33bb44cc"
  },
  "traceId": "0af7651916cd43dd8448eb211c80319c"
}
```

## Response (errors)

| Code | HTTP | Trigger |
|---|--:|---|
| `VALIDATION_ERROR` | 400 | `delta == 0`, `reason` missing/empty/> 255 chars, or `sku` invalid |
| `AUTH_INVALID` | 401 | JWT missing / expired / signature invalid |
| `AUTH_FORBIDDEN` | 403 | JWT valid but `claims.role != 'ADMIN'` |
| `NOT_FOUND` | 404 | No `stock_levels` row for `sku` |
| `STOCK_NEGATIVE_INVARIANT` | 409 | `current.available_qty + delta < 0` (Layer-1 guard) |
| `IDEMPOTENCY_KEY_REUSED` | 409 | Same client-supplied `Idempotency-Key` + different request body (only if header present) |
| `DATABASE_UNAVAILABLE` | 503 | Postgres unreachable |

## Business logic steps

1. `common/middleware.JWT()` validates ES256 signature, `iss`, `aud`, `exp`. Injects `token.Claims` into context.
2. Role guard: `claims.role == 'ADMIN'` else `AUTH_FORBIDDEN`.
3. Bind + validate body (`common/validator`): `delta != 0`, `reason ∈ [1, 255]`, `sku ∈ [1, 64]`.
4. (Optional) If `Idempotency-Key` header present: `SELECT idempotency_keys WHERE key=$1 AND endpoint='inventory.stock.adjust' AND ...` — replay branch on hash match; 409 on hash mismatch. (Per `cross-cutting.idempotency`; scope per-(endpoint, actorUserId).)
5. `BEGIN TX`.
6. `SELECT sku, available_qty, version FROM stock_levels WHERE sku = $1 FOR UPDATE` — row lock per `cross-cutting.persistence.concurrency`.
7. If row missing → `ROLLBACK`; return `NOT_FOUND`.
8. **Layer-1 no-negative guard:** if `row.available_qty + delta < 0` → `ROLLBACK`; return `STOCK_NEGATIVE_INVARIANT`. (Layer-2 `CHECK (available_qty >= 0)` is the defense-in-depth backstop.)
9. `UPDATE stock_levels SET available_qty = available_qty + $delta, version = version + 1, updated_at = NOW() WHERE sku = $1`.
10. `adjustmentId = uuidv7()`; `INSERT INTO stock_adjustments (id, sku, delta, reason, actor_user_id) VALUES ($1, $2, $3, $4, $claimsSub)` — actor pulled from `claims.sub`.
11. (Optional idempotency persist) If `Idempotency-Key` header present, `INSERT idempotency_keys (...)` inside the same tx.
12. `COMMIT`.
13. `wrapper.Respond(c, SUCCESS, {sku, availableQty: row.available_qty + delta, adjustmentId})`.

## Side effects

- `UPDATE stock_levels` (one row) — `available_qty`, `version`, `updated_at`.
- `INSERT stock_adjustments` (one row) — `delta`, `reason`, `actor_user_id`.
- (Optional) `INSERT idempotency_keys` if header present.
- No outbox event. (Stock adjustments are an admin-driven audit signal; downstream consumers don't need an event in MVP.)

## Idempotency

- API-level: optional client-supplied `Idempotency-Key` header per `cross-cutting.idempotency`. Scope: per-(endpoint=`inventory.stock.adjust`, actorUserId).
- Same key + same body → return persisted envelope verbatim.
- Same key + different body → 409 `IDEMPOTENCY_KEY_REUSED`.
- Without `Idempotency-Key` header: NOT idempotent (each call applies a fresh delta).

## Performance

- p95 < 100ms (single tx, FOR UPDATE on one PK row, two INSERTs/UPDATEs).
- Expected QPS at MVP: < 1 (admin manual operation).

## Test cases

- `stock_adjust_admin_increments_and_writes_audit_row`
- `stock_adjust_admin_decrements_and_writes_audit_row`
- `stock_adjust_negative_invariant_rejects_no_change_no_audit`
- `stock_adjust_delta_zero_returns_VALIDATION_ERROR`
- `stock_adjust_missing_reason_returns_VALIDATION_ERROR`
- `stock_adjust_customer_token_returns_AUTH_FORBIDDEN`
- `stock_adjust_unknown_sku_returns_NOT_FOUND`
- `stock_adjust_idempotency_key_replay_returns_persisted_envelope`
- `stock_adjust_idempotency_key_with_different_body_returns_409`
- `stock_adjust_envelope_carries_traceId`

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | Tech-Designer (Opus 4.7 1M) | Created |
