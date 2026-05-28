# POST /api/v1/cart/cart/remove-item

## Summary

Removes a specific cart line item from the authenticated customer's cart. The remove is idempotent: if the `cartItemId` was already removed (or never existed in this cart), the endpoint returns 200 SUCCESS with the current cart state — no error. Crucially, the product's current catalog status is NOT checked on this path; INACTIVE and DELETED product lines can always be removed so the customer can clear them and proceed to checkout (STORY_CART_REMOVE AC #2, CART-005). Returns the enriched cart after the delete.

## Story refs

- `STORY_CART_REMOVE` (CART-003, CART-004, CART-005)

## Contract ref

[`cart.remove-item`](../../architecture/contracts.json#cart.remove-item)

## Auth

- Tier: **customer_jwt**
- `common/middleware.JWTMiddleware` validates ES256 signature, `iss`, `aud`, `exp`, `role`.
- Ownership enforced implicitly: `DELETE WHERE id=$1 AND cart_id=$2` — `cart_id` derived from `carts WHERE user_id=claims.Sub`. A `cartItemId` belonging to another user's cart will simply match 0 rows (idempotent no-op), not a 403.

## Request

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "required": ["cartItemId"],
  "properties": {
    "cartItemId": {
      "type": "string",
      "format": "uuid",
      "description": "UUID of the cart_items row to remove"
    }
  }
}
```

## Response (success)

HTTP 200 — envelope:

```json
{
  "code": "SUCCESS",
  "message": "Cart item removed",
  "traceId": "0af7651916cd43dd8448eb211c80319c",
  "data": {
    "cart": {
      "customerUserId": "01935b9c-0001-7000-8000-000000000001",
      "items": [],
      "subtotal":  0,
      "updatedAt": "2026-05-08T10:40:00Z"
    }
  }
}
```

**Idempotent no-op:** if the item was already removed, `items[]` reflects the current state of the cart (which may be empty or contain other lines). `subtotal` is recomputed from remaining checkoutable lines.

## Response (errors)

| Code              | HTTP | Trigger                                              |
|---|--:|---|
| `AUTH_MISSING`    | 401  | No `Authorization: Bearer` header                    |
| `AUTH_INVALID`    | 401  | JWT validation failed                                |
| `VALIDATION_ERROR`| 400  | `cartItemId` is empty or not a valid UUID format     |
| `INTERNAL_ERROR`  | 500  | Unexpected DB error                                  |

**Note:** NOT_FOUND is NOT returned for an already-removed line — the idempotent contract returns 200. No UPSTREAM_TIMEOUT is possible on this path — the remove path does NOT call catalog or inventory.

## Business logic steps

1. Bind and validate request body; return VALIDATION_ERROR if `cartItemId` is empty or not a valid UUID.
2. Extract `userID = claims.Sub`.
3. Fetch `cart_id` from `cart.carts WHERE user_id=$1`. If no cart exists → return empty cart shape with 200 SUCCESS (idempotent; nothing to remove).
4. `DELETE FROM cart.cart_items WHERE id=$1 AND cart_id=$2`. Capture `RowsAffected()`:
   - `RowsAffected() == 1`: item deleted.
   - `RowsAffected() == 0`: item not in this cart (already removed, never existed, or belongs to another user's cart) — treat as idempotent no-op.
5. `UPDATE cart.carts SET updated_at=now() WHERE cart_id=$1` — always update even on no-op to advance the staleness marker.
6. Return enriched cart via `CartService.ReadCart(ctx, userID)` (recomputes subtotal over remaining checkoutable lines — CART-004, CART-003).

**No catalog call on this path:** INACTIVE/DELETED product lines can always be removed. Calling catalog here would break the remove-as-escape-hatch for non-checkoutable items.

## Side effects

- DELETE 0 or 1 row from `cart.cart_items`.
- UPDATE `carts.updated_at` (always — even on no-op).
- No outbox event.

## Idempotency

**Idempotent by construction.** `DELETE WHERE id=$1 AND cart_id=$2` with 0 rows affected is a safe no-op. The 200 SUCCESS response is returned in both the deletion and no-op cases, with the current cart state as the response body.

## Performance

- p95 target: < 300ms.
- Expected QPS at MVP: 5–20.
- No upstream HTTP calls on this path — DB-only. The response time is dominated by the ReadCart fan-out (catalog + inventory) called after the delete.

## Test cases

- `remove_item_deletes_line_and_recalculates_subtotal`
- `remove_item_inactive_product_line_removes_successfully_no_catalog_call`
- `remove_item_already_removed_returns_200_idempotent`
- `remove_item_cartItemId_belonging_to_other_user_returns_200_no_change`
- `remove_item_no_cart_returns_200_empty_cart`
- `remove_item_invalid_cartItemId_uuid_returns_VALIDATION_ERROR`
- `remove_item_unauthenticated_returns_401`

## Change log

| Date       | Author                                        | Change          |
|---|---|---|
| 2026-05-08 | Tech-Designer (Claude Sonnet 4.6, dry-run #2) | Initial version |
