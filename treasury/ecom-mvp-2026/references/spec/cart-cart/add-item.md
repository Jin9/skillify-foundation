# POST /api/v1/cart/cart/add-item

## Summary

Adds a product to the authenticated customer's cart. If the same SKU already exists in the cart, the quantities are merged into a single line (qty += requested qty) — a duplicate add never creates a second row (CART-006). The endpoint verifies the product is ACTIVE via `catalog.product.detail` before writing; adding an INACTIVE or DELETED product is rejected. No stock reservation is created at add-time (INV-002). Returns the full enriched cart (live prices, availability flags) after the write.

## Story refs

- `STORY_CART_ADD` (CART-001, CART-006, CART-007)

## Contract ref

[`cart.add-item`](../../architecture/contracts.json#cart.add-item)

## Auth

- Tier: **customer_jwt**
- `common/middleware.JWTMiddleware` validates ES256 signature, `iss=shoppilot-identity`, `aud=shoppilot-api`, `exp`, and `role`.
- `claims.role` must be `CUSTOMER`; ADMIN callers receive `AUTH_FORBIDDEN` (not expected in normal flow).
- `claims.Sub` is used as `user_id` for all DB operations — never from the request body.

## Request

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "required": ["productId", "qty"],
  "properties": {
    "productId": {
      "type": "string",
      "format": "uuid",
      "description": "UUID of the product to add to the cart"
    },
    "qty": {
      "type": "integer",
      "minimum": 1,
      "description": "Quantity to add. Must be >= 1. qty=0 or negative is a VALIDATION_ERROR."
    }
  }
}
```

## Response (success)

HTTP 200 — envelope:

```json
{
  "code": "SUCCESS",
  "message": "Item added to cart",
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
          "qty":           2,
          "lineSubtotal":  1180,
          "availableQty":  8,
          "productStatus": "ACTIVE",
          "checkoutable":  true,
          "available":     true,
          "unavailReason": ""
        }
      ],
      "subtotal":  1180,
      "updatedAt": "2026-05-08T10:30:00Z"
    }
  }
}
```

**Note:** `subtotal` is the sum of `lineSubtotal` for `checkoutable=true` lines only. All prices are integer THB (no fractional THB in MVP).

## Response (errors)

| Code              | HTTP | Trigger                                                                 |
|---|--:|---|
| `AUTH_MISSING`    | 401  | No `Authorization: Bearer` header on the request                       |
| `AUTH_INVALID`    | 401  | JWT signature/exp/iss/aud/role validation failed                        |
| `VALIDATION_ERROR`| 400  | `productId` is empty or not a valid UUID; or `qty` < 1                 |
| `NOT_FOUND`       | 404  | `catalog.product.detail` returned NOT_FOUND for the given `productId`  |
| `PRODUCT_INACTIVE`| 409  | Product exists but `status == INACTIVE`; cannot add to cart             |
| `PRODUCT_DELETED` | 409  | Product exists but `status == DELETED`; cannot add to cart              |
| `UPSTREAM_TIMEOUT`| 504  | `catalog.product.detail` call exceeded deadline or returned network error|
| `INTERNAL_ERROR`  | 500  | Unexpected DB error (pgx pool exhausted, tx failure)                   |

## Business logic steps

1. Bind and validate request body via `common/wrapper.BindJSON[AddItemRequest]`; return 400 VALIDATION_ERROR on malformed body or constraint violation.
2. Extract `userID = claims.Sub` from context via `common/token.ClaimsFromContext(ctx)`.
3. Call `catalog.product.detail(productId)` via `common/httpclient.Post`; on network error or deadline exceeded → UPSTREAM_TIMEOUT. Map response:
   - `NOT_FOUND` → return NOT_FOUND
   - `status == INACTIVE` → return PRODUCT_INACTIVE
   - `status == DELETED` → return PRODUCT_DELETED
   - `status != ACTIVE` (e.g. DRAFT) → return PRODUCT_INACTIVE
4. Open a pgx transaction.
5. UPSERT `carts`: `INSERT INTO cart.carts (cart_id, user_id, updated_at) VALUES ($1,$2,now()) ON CONFLICT (user_id) DO UPDATE SET updated_at=now() RETURNING cart_id` — get or create the cart for this user. `cart_id = common/generator.NewUUIDv7()` for new rows.
6. Merge the item: `INSERT INTO cart.cart_items (id, cart_id, product_id, qty, added_at) VALUES ($1,$2,$3,$4,now()) ON CONFLICT (cart_id, product_id) DO UPDATE SET qty = cart_items.qty + EXCLUDED.qty` — single parameterized statement; no read-then-write race. `id = common/generator.NewUUIDv7()` for new rows.
7. Commit the transaction.
8. Call `CartService.ReadCart(ctx, userID)` to build and return the enriched cart (live prices, availability flags, recomputed subtotal).

## Side effects

- UPSERT one row in `cart.carts` (creates cart on first add; updates `updated_at` otherwise).
- INSERT or UPDATE one row in `cart.cart_items` (merge on same SKU).
- No outbox event — cart mutations are not published to Kafka.
- No inventory reservation — INV-002 prohibits reservation at add-time.

## Idempotency

**Same-SKU merge is the idempotency mechanism (CART-006).** Two identical add-item calls for the same product serialize on the UNIQUE row lock in Postgres; the final qty equals the sum of both requests. No separate idempotency table or Idempotency-Key header required. A completely duplicate call (same productId, same qty) results in qty doubling — this is the correct and documented behavior per the merge contract.

## Performance

- p95 target: < 500ms (dominated by the synchronous `catalog.product.detail` call).
- Expected QPS at MVP: 5–30.
- DB operations are sub-10ms on pgx pool; the catalog call is the bottleneck.

## Test cases

- `add_item_new_sku_creates_line`
- `add_item_same_sku_merges_qty_sums_correctly`
- `add_item_concurrent_same_sku_10_goroutines_single_row_qty_10` (integration)
- `add_item_unauthenticated_returns_401`
- `add_item_invalid_product_id_returns_VALIDATION_ERROR`
- `add_item_qty_zero_returns_VALIDATION_ERROR`
- `add_item_qty_negative_returns_VALIDATION_ERROR`
- `add_item_inactive_product_returns_PRODUCT_INACTIVE`
- `add_item_deleted_product_returns_PRODUCT_DELETED`
- `add_item_catalog_timeout_returns_UPSTREAM_TIMEOUT`

## Change log

| Date       | Author                                        | Change          |
|---|---|---|
| 2026-05-08 | Tech-Designer (Claude Sonnet 4.6, dry-run #2) | Initial version |
