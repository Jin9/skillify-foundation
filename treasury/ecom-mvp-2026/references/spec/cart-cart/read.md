# POST /api/v1/cart/cart/read

## Summary

Returns the authenticated customer's cart enriched with live product data and availability flags. The cart service fans out in parallel to `catalog.product.detail` (per item, semaphore-bounded at 10 concurrent calls) and `inventory.stock.bulk-read` (one batched call for all product IDs) within a 3000ms timeout. Each line item carries an `available: bool` flag and an `unavailReason` so the frontend can render non-checkoutable items distinctly — no line is silently dropped from the response. Subtotal is computed server-side from `checkoutable=true` lines only (CART-004). Informational only — `checkout.commit` recomputes the authoritative total (CHK-005).

## Story refs

- `STORY_CART_ADD` (CART-004, CART-005, CART-007)
- `STORY_CART_REMOVE` (CART-005)

## Contract ref

[`cart.read`](../../architecture/contracts.json#cart.read)

## Auth

- Tier: **customer_jwt**
- `common/middleware.JWTMiddleware` validates ES256 signature, `iss`, `aud`, `exp`, `role`.
- `claims.Sub` is the only userId source — no body parameter for user scoping.

## Request

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {},
  "description": "Empty body. userId is derived exclusively from JWT claims.Sub."
}
```

## Response (success)

HTTP 200 — envelope:

```json
{
  "code": "SUCCESS",
  "message": "Cart retrieved",
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
        },
        {
          "cartItemId":    "01935b9c-0001-7000-8000-000000000011",
          "productId":     "01935b9c-0001-7000-8000-000000000003",
          "sku":           "PANTS-BLUE-L",
          "name":          "Blue Linen Pants (L)",
          "image":         "https://cdn.example.com/pants-blue-l.jpg",
          "currentPrice":  890,
          "qty":           1,
          "lineSubtotal":  890,
          "availableQty":  0,
          "productStatus": "ACTIVE",
          "checkoutable":  false,
          "available":     false,
          "unavailReason": "OUT_OF_STOCK"
        }
      ],
      "subtotal":  1180,
      "updatedAt": "2026-05-08T10:30:00Z"
    }
  }
}
```

**Empty cart response (no rows or no cart):**

```json
{
  "code": "SUCCESS",
  "message": "Cart retrieved",
  "traceId": "...",
  "data": {
    "cart": {
      "customerUserId": "01935b9c-0001-7000-8000-000000000001",
      "items":     [],
      "subtotal":  0,
      "updatedAt": "2026-05-08T10:00:00Z"
    }
  }
}
```

**Field semantics:**

| Field            | Description                                                                 |
|---|---|
| `currentPrice`   | Live price from `catalog.product.detail`. Integer THB. Never stored in cart.|
| `lineSubtotal`   | `currentPrice * qty`. Integer THB.                                          |
| `availableQty`   | From `inventory.stock.bulk-read`. 0 if product_id not in inventory (uninitialized SKU). |
| `checkoutable`   | `true` only if `productStatus == ACTIVE AND availableQty > 0`               |
| `available`      | Same as `checkoutable` in current implementation. Explicit field per calibration requirement. |
| `unavailReason`  | One of `""`, `"PRODUCT_INACTIVE"`, `"PRODUCT_DELETED"`, `"OUT_OF_STOCK"`. Empty when `available=true`. |
| `subtotal`       | Sum of `lineSubtotal` for `checkoutable=true` lines only (CART-004). Informational — not authoritative. |
| `updatedAt`      | From `carts.updated_at`; advanced on every cart mutation (CART-007).        |

## Response (errors)

| Code              | HTTP | Trigger                                                                                      |
|---|--:|---|
| `AUTH_MISSING`    | 401  | No `Authorization: Bearer` header                                                            |
| `AUTH_INVALID`    | 401  | JWT validation failed                                                                        |
| `UPSTREAM_TIMEOUT`| 504  | Any catalog.product.detail call OR the inventory.stock.bulk-read call exceeds the 3000ms fan-out deadline or returns a network error. Strict: no partial success — the entire read fails. |
| `INTERNAL_ERROR`  | 500  | Unexpected DB error fetching cart / cart_items                                               |

## Business logic steps

1. Extract `userID = claims.Sub`.
2. `SELECT * FROM cart.carts WHERE user_id=$1`. If no row → return empty cart shape (200 SUCCESS; no fan-out).
3. `SELECT * FROM cart.cart_items WHERE cart_id=$1 ORDER BY added_at ASC`. If 0 rows → return empty cart shape (200 SUCCESS; no fan-out).
4. Enforce `MAX_CART_ITEMS=50`: if `len(items) > 50`, trim to first 50 and log a `slog.Warn` with `cart_id` and original count.
5. Collect `productIDs = [item.productId for item in items]`.
6. Create fan-out context: `fanoutCtx, cancel = context.WithTimeout(ctx, 3*time.Second); defer cancel()`.
7. Pre-allocate `catalogResults []CatalogDetail` of length `len(items)` (indexed, not map — no concurrent write race).
8. Start `errgroup`: `g, gCtx = errgroup.WithContext(fanoutCtx)`.
9. **Inventory goroutine** (single, inside `g`): call `inventory.stock.bulk-read(productIDs)` via `common/httpclient.Post` with `gCtx`. On error → return error (cancels all goroutines). Store results in `stockByProductID map[string]int` (written by this single goroutine only — no race).
10. **Catalog goroutines** (one per item, inside `g`, semaphore=10): for each item `i`:
    - Acquire semaphore slot (`sem <- struct{}{}`); respect `gCtx` cancellation.
    - Call `catalog.product.detail(item.productId)` via `common/httpclient.Post` with `gCtx`.
    - On error → release semaphore, return error (cancels all remaining goroutines).
    - Store result at `catalogResults[i]` (write-by-index; safe without mutex).
    - Release semaphore (`defer <-sem`).
11. `err = g.Wait()`. If `err != nil` → return UPSTREAM_TIMEOUT (strict; no partial result).
12. For each item, compute:
    - `currentPrice = catalogResults[i].price`
    - `productStatus = catalogResults[i].status`
    - `availableQty = stockByProductID[item.productId]` (default 0 if key absent — uninitialized SKU → OUT_OF_STOCK)
    - `checkoutable = (productStatus == "ACTIVE" && availableQty > 0)`
    - `available = checkoutable`
    - `unavailReason` based on the failing condition (PRODUCT_INACTIVE, PRODUCT_DELETED, OUT_OF_STOCK)
    - `lineSubtotal = currentPrice * qty` (integer arithmetic)
13. `subtotal = sum(lineSubtotal for items where checkoutable=true)`.
14. Return `Cart{customerUserId, items, subtotal, updatedAt: carts.updated_at}`.

## Side effects

None — read-only endpoint. No DB writes, no outbox events.

## Idempotency

Read — trivially idempotent. Repeated calls return the current state.

## Performance

- p95 target: < 1000ms (dominated by fan-out; semaphore=10 means worst-case 5 batches of 10 catalog calls at ~50ms each = ~250ms + inventory call).
- Expected QPS at MVP: 10–50.
- `MAX_CART_ITEMS=50` bounds worst-case fan-out to 50 catalog calls + 1 inventory call.
- 3000ms timeout provides sufficient headroom; p99 catalog call is expected < 200ms within the cluster.

## Sequence diagram

```mermaid
sequenceDiagram
  participant FE   as Frontend
  participant CART as Cart Service
  participant DB   as Postgres (cart schema)
  participant CAT  as Catalog Service
  participant INV  as Inventory Service

  FE->>+CART: POST /api/v1/cart/cart/read (Bearer JWT)
  CART->>+DB: SELECT carts + cart_items WHERE user_id=$1
  DB-->>-CART: cart row + N item rows

  note over CART: MAX_CART_ITEMS=50 enforced; fan-out context timeout=3000ms

  par errgroup fan-out (all within 3000ms timeout)
    CART->>+INV: POST /api/v1/inventory/stock/bulk-read {skus: [p1,p2,...pN]}
    INV-->>-CART: {stocks: [{productId, availableQty}, ...]}
  and
    loop per item (semaphore=10 concurrent)
      CART->>+CAT: POST /api/v1/catalog/product/detail {productId: pI}
      CAT-->>-CART: {sku, name, image, price, status}
    end
  end

  note over CART: any upstream error → UPSTREAM_TIMEOUT (strict)
  note over CART: enrich items; compute checkoutable + unavailReason; sum subtotal

  CART-->>-FE: 200 {cart: {items[], subtotal, updatedAt}}
```

## Test cases

- `read_cart_empty_cart_returns_empty_shape_no_fanout`
- `read_cart_no_cart_row_returns_empty_shape`
- `read_cart_returns_enriched_items_with_live_price`
- `read_cart_inactive_product_checkoutable_false_excluded_from_subtotal`
- `read_cart_deleted_product_checkoutable_false_excluded_from_subtotal`
- `read_cart_out_of_stock_product_checkoutable_false_excluded_from_subtotal`
- `read_cart_subtotal_only_includes_checkoutable_lines`
- `read_cart_updated_at_reflects_last_mutation`
- `read_cart_max_50_items_enforced_trims_to_50`
- `read_cart_catalog_timeout_any_item_returns_UPSTREAM_TIMEOUT`
- `read_cart_inventory_timeout_returns_UPSTREAM_TIMEOUT`
- `read_cart_unauthenticated_returns_401`
- `read_cart_concurrent_10_reads_all_succeed_no_race` (integration)

## Change log

| Date       | Author                                        | Change          |
|---|---|---|
| 2026-05-08 | Tech-Designer (Claude Sonnet 4.6, dry-run #2) | Initial version; fan-out sequence diagram included per optional calibration |
