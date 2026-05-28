# POST /api/v1/cart/cart/update-item

## Summary

Updates the quantity of a specific cart line item. If `qty` is set to 0, the line is removed (same effect as `remove-item` — no zero-qty lines are permitted in the cart per CART-003). If `qty >= 1`, a best-effort stock check is performed via `inventory.stock.bulk-read`; if the requested qty exceeds `availableQty`, the update is rejected with a VALIDATION_ERROR that includes the current `availableQty` (CART-002). This stock check is advisory — no reservation is created. Returns the enriched cart after the write.

## Story refs

- `STORY_CART_UPDATE_QTY` (CART-002, CART-004)
- `STORY_CART_ADD` (CART-003 — qty=0 remove path)

## Contract ref

[`cart.update-item`](../../architecture/contracts.json#cart.update-item)

## Auth

- Tier: **customer_jwt**
- `common/middleware.JWTMiddleware` validates ES256 signature, `iss`, `aud`, `exp`, `role`.
- `claims.Sub` is used as `user_id`; the handler resolves `cart_id` from `carts WHERE user_id=$1` and verifies the `cartItemId` belongs to that cart. Cross-user access returns NOT_FOUND.

## Request

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "required": ["cartItemId", "qty"],
  "properties": {
    "cartItemId": {
      "type": "string",
      "format": "uuid",
      "description": "UUID of the cart_items row to update"
    },
    "qty": {
      "type": "integer",
      "minimum": 0,
      "description": "New quantity. 0 = remove the line. Negative = VALIDATION_ERROR."
    }
  }
}
```

## Response (success)

HTTP 200 — envelope:

```json
{
  "code": "SUCCESS",
  "message": "Cart item updated",
  "traceId": "0af7651916cd43dd8448eb211c80319c",
  "data": {
    "cart": {
      "customerUserId": "01935b9c-0001-7000-8000-000000000001",
      "items": [
        {
          "cartItemId":    "01935b9c-0001-7000-8000-000000000010",
          "productId":     "01935b9c-0001-7000-8000-000000000002",
          "sku":           "SHIRT-RED-M",
          "name":          "Red Cotton Shirt (M)",
          "image":         "https://cdn.example.com/shirt-red-m.jpg",
          "currentPrice":  590,
          "qty":           3,
          "lineSubtotal":  1770,
          "availableQty":  5,
          "productStatus": "ACTIVE",
          "checkoutable":  true,
          "available":     true,
          "unavailReason": ""
        }
      ],
      "subtotal":  1770,
      "updatedAt": "2026-05-08T10:35:00Z"
    }
  }
}
```

**qty=0 success response:** same shape but the updated line is absent from `items[]`; `subtotal` is recomputed from remaining lines.

## Response (errors)

| Code              | HTTP | Trigger                                                                                    |
|---|--:|---|
| `AUTH_MISSING`    | 401  | No `Authorization: Bearer` header                                                          |
| `AUTH_INVALID`    | 401  | JWT validation failed                                                                      |
| `VALIDATION_ERROR`| 400  | `cartItemId` empty/invalid UUID; `qty < 0`; OR `qty > availableQty` (CART-002 — response `data` includes `{availableQty: N}`) |
| `NOT_FOUND`       | 404  | No cart for this user, or `cartItemId` not in this user's cart                             |
| `UPSTREAM_TIMEOUT`| 504  | `inventory.stock.bulk-read` call exceeded deadline (qty >= 1 path only)                   |
| `INTERNAL_ERROR`  | 500  | Unexpected DB error                                                                        |

## Business logic steps

1. Bind and validate request body; return VALIDATION_ERROR if `cartItemId` is not a valid UUID or `qty < 0`.
2. Extract `userID = claims.Sub`.
3. Fetch `cart_id` from `cart.carts WHERE user_id=$1`. If no row → NOT_FOUND.
4. Verify `cartItemId` belongs to this cart: `SELECT product_id FROM cart.cart_items WHERE id=$1 AND cart_id=$2`. If 0 rows → NOT_FOUND.
5. **If qty == 0:** delegate to internal remove logic — `DELETE FROM cart.cart_items WHERE id=$1 AND cart_id=$2`; `UPDATE carts SET updated_at=now()`; return enriched cart (CART-003). Skip steps 6–8.
6. **If qty >= 1:** call `inventory.stock.bulk-read([productId])` via `common/httpclient.Post`. On network error → UPSTREAM_TIMEOUT. Extract `availableQty` from response.
7. If `qty > availableQty` → return VALIDATION_ERROR with detail `{availableQty: N}` (CART-002). No DB write.
8. `UPDATE cart.cart_items SET qty=$1 WHERE id=$2 AND cart_id=$3`; `UPDATE carts SET updated_at=now()`. Single pgx transaction.
9. Return enriched cart via `CartService.ReadCart(ctx, userID)`.

## Side effects

- DELETE one row from `cart.cart_items` (qty=0 path).
- UPDATE one row in `cart.cart_items` (qty>=1 path).
- UPDATE `carts.updated_at` in all paths.
- No outbox event.

## Idempotency

Last-writer-wins on `qty`. No idempotency key required. Concurrent updates for the same `cartItemId` serialize on the pgx row lock; the last committed write wins.

## Performance

- p95 target: < 500ms.
- Expected QPS at MVP: 5–20.
- For qty>=1 path: one inventory HTTP call adds ~20–50ms round-trip inside the cluster.

## Test cases

- `update_item_qty_increases_within_stock_succeeds`
- `update_item_qty_zero_removes_line_and_recalculates_subtotal`
- `update_item_qty_exceeds_available_returns_VALIDATION_ERROR_with_availableQty`
- `update_item_unknown_cartItemId_returns_NOT_FOUND`
- `update_item_cross_user_cartItemId_returns_NOT_FOUND`
- `update_item_negative_qty_returns_VALIDATION_ERROR`
- `update_item_unauthenticated_returns_401`
- `update_item_inventory_timeout_returns_UPSTREAM_TIMEOUT`

## Change log

| Date       | Author                                        | Change          |
|---|---|---|
| 2026-05-08 | Tech-Designer (Claude Sonnet 4.6, dry-run #2) | Initial version |
