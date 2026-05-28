# POST /api/v1/cart/cart/clear-on-checkout

## Summary

Internal endpoint called by the Checkout service during `checkout.commit` (step 9 of the orchestration). Removes only the specified cart line items (those moved into the order) from the customer's cart. Items not listed in `cartItemIds` remain in the cart (CHK-010). The operation is idempotent on `orderId`: a replay returns the cached `removed` count without re-executing the DELETE. Cross-user leakage is prevented by validating `customerUserId == claims.Sub` before any DB operation.

## Story refs

- `STORY_CART_REMOVE` (CHK-010 — checkout integration)

## Contract ref

[`cart.clear-on-checkout`](../../architecture/contracts.json#cart.clear-on-checkout)

## Auth

- Tier: **customer_jwt** (internal use — Checkout forwards the customer's own Bearer JWT)
- `common/middleware.JWTMiddleware` validates ES256 signature, `iss`, `aud`, `exp`.
- **Additional cross-user guard:** handler validates `request.customerUserId == claims.Sub` — VALIDATION_ERROR on mismatch. This is the primary leakage-prevention mechanism in the absence of a SERVICE role token in MVP (AMB-002 resolution).
- This endpoint is NOT protected by `X-Internal-Secret`. Checkout calls it using the customer's Bearer token forwarded from `checkout.commit`'s inbound `Authorization` header.

## Request

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "required": ["customerUserId", "cartItemIds", "orderId"],
  "properties": {
    "customerUserId": {
      "type": "string",
      "format": "uuid",
      "description": "Must equal JWT claims.Sub. Validated before any DB operation."
    },
    "cartItemIds": {
      "type": "array",
      "items": { "type": "string", "format": "uuid" },
      "minItems": 1,
      "description": "UUIDs of cart_items rows to remove. Only these rows are deleted; all others remain."
    },
    "orderId": {
      "type": "string",
      "format": "uuid",
      "description": "Idempotency anchor. Replay with the same orderId returns the cached removed count."
    }
  }
}
```

## Response (success)

HTTP 200 — envelope:

```json
{
  "code": "SUCCESS",
  "message": "Cart cleared for checkout",
  "traceId": "0af7651916cd43dd8448eb211c80319c",
  "data": {
    "removed": 3
  }
}
```

**Replay response (idempotent, same `orderId`):**

```json
{
  "code": "SUCCESS",
  "message": "Cart cleared for checkout",
  "traceId": "0af7651916cd43dd8448eb211c80319c",
  "data": {
    "removed": 3
  }
}
```

Replay returns the cached `removed` count from `cart_clear_idempotency`. The value is identical to the first call regardless of whether the cart items still exist.

**`removed: 0` is valid:** occurs on replay after all items have already been individually removed, or when `cartItemIds` contained IDs not present in this cart.

## Response (errors)

| Code              | HTTP | Trigger                                                                           |
|---|--:|---|
| `AUTH_MISSING`    | 401  | No `Authorization: Bearer` header                                                 |
| `AUTH_INVALID`    | 401  | JWT validation failed                                                             |
| `VALIDATION_ERROR`| 400  | `customerUserId != claims.Sub`; or `cartItemIds` is empty; or `orderId` empty/invalid UUID; or `customerUserId` invalid UUID |
| `INTERNAL_ERROR`  | 500  | Unexpected DB error (tx failure, pgx pool exhausted)                              |

## Business logic steps

1. Bind and validate request body.
2. Extract `userID = claims.Sub`.
3. **Cross-user guard (first, before any DB operation):** if `request.customerUserId != userID` → return VALIDATION_ERROR immediately. Do not log the mismatch as a warn — log as `slog.Error` with `msg="clear_on_checkout_user_mismatch"` and `remoteAddr`, `traceId`, `requestedUserId` (never log `claims.Sub` value in full to avoid leaking PII to log stores; log only a prefix).
4. **Idempotency check:** `SELECT removed FROM cart.cart_clear_idempotency WHERE order_id=$1`. If row found → return `{removed: N}` immediately (200 SUCCESS). No DB write.
5. Fetch `cart_id` from `cart.carts WHERE user_id=$1`. If no cart → return `{removed: 0}` (nothing to clear; still persist idempotency record with removed=0).
6. Open pgx transaction.
7. `DELETE FROM cart.cart_items WHERE cart_id=$1 AND id=ANY($2::uuid[])`. Use pgx array parameter — **never** build the `IN` clause via string formatting. Capture `RowsAffected()` as `removed`.
8. `UPDATE cart.carts SET updated_at=now() WHERE cart_id=$1`.
9. `INSERT INTO cart.cart_clear_idempotency (order_id, removed, created_at) VALUES ($1,$2,now()) ON CONFLICT (order_id) DO NOTHING` — commited atomically with the DELETE in the same transaction.
10. Commit transaction.
11. Return `{removed: N}`.

## Side effects

- DELETE N rows from `cart.cart_items` (N = `len(cartItemIds)` that existed; may be 0).
- UPDATE `carts.updated_at`.
- INSERT one row into `cart.cart_clear_idempotency` (or no-op if `ON CONFLICT`).
- No outbox event.
- Leftover cart items not in `cartItemIds` are untouched (CHK-010).

## Idempotency

Idempotent on `orderId`. Implementation:

1. Pre-check: `SELECT removed FROM cart_clear_idempotency WHERE order_id=$1`. If found → return cached `removed` immediately.
2. Post-write: `INSERT INTO cart_clear_idempotency ON CONFLICT (order_id) DO NOTHING` inside the same pgx transaction as the DELETE.

This two-step pattern (check-before + atomic insert) ensures:
- Single execution: the DELETE runs at most once per `orderId`.
- Replay safety: concurrent replays after the first call race on the `ON CONFLICT`; only one wins; all others read the cached value on their next SELECT.
- The `order_id` PK is globally unique (UUID v7 from the Order service); no namespace collision is possible.

## Performance

- p95 target: < 200ms.
- Expected QPS at MVP: matches checkout.commit rate (1–5 QPS).
- DB-only path on replay (no DELETE). Single-transaction path on first call.

## Test cases

- `clear_on_checkout_removes_specified_items_only_leaves_others`
- `clear_on_checkout_replay_same_orderId_returns_cached_removed_count_no_second_delete`
- `clear_on_checkout_removed_zero_when_items_already_gone`
- `clear_on_checkout_customerUserId_mismatch_returns_VALIDATION_ERROR`
- `clear_on_checkout_empty_cartItemIds_returns_VALIDATION_ERROR`
- `clear_on_checkout_invalid_orderId_returns_VALIDATION_ERROR`
- `clear_on_checkout_no_cart_returns_removed_zero`
- `clear_on_checkout_unauthenticated_returns_401`
- `clear_on_checkout_concurrent_replay_both_return_same_removed_count` (integration — race test)

## Change log

| Date       | Author                                        | Change          |
|---|---|---|
| 2026-05-08 | Tech-Designer (Claude Sonnet 4.6, dry-run #2) | Initial version |
