# POST /api/v1/catalog/product/detail

## Summary

Returns the full detail record for a single product. Public for ACTIVE and INACTIVE products; DELETED and DRAFT products return `NOT_FOUND` to guest/customer callers but are visible to ADMIN callers. Includes the full image array (ordered by `sort_order`), a live stock-status indicator computed from a best-effort call to `inventory.stock.read`, and a review summary placeholder (`{averageRating: null, reviewCount: 0}`) while the Reviews module is deferred. This endpoint is also called by `cart.add-item` and `checkout.commit` step 4 to enforce `PRODUCT_INACTIVE`/`PRODUCT_DELETED` checks. Implements requirement CAT-006, CAT-009, CAT-010.

## Story refs

- `STORY_CATALOG_DETAIL`

## Contract ref

[`catalog.product.detail`](../../architecture/contracts.json#catalog.product.detail)

## Auth

- Tier: **optional** — no Authorization header required.
- If present and valid (ES256, ADMIN role): DELETED/DRAFT products are returned.
- If absent, invalid, or CUSTOMER role: DELETED/DRAFT products return `NOT_FOUND`.
- Handler uses `common/token` to parse the Authorization header if present; does NOT reject missing auth.

## Request

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "required": ["productId"],
  "properties": {
    "productId": { "type": "string", "format": "uuid" }
  }
}
```

## Response (success)

HTTP 200, envelope:

```json
{
  "code": "SUCCESS",
  "message": "OK",
  "data": {
    "productId":   "01935b9c-0000-7000-0000-000000000001",
    "sku":         "SHOE-RED-42",
    "name":        "Running Shoe Red 42",
    "description": "Lightweight running shoe for road surfaces.",
    "images": [
      "https://cdn.example.com/shoes/red-42/01.jpg",
      "https://cdn.example.com/shoes/red-42/02.jpg"
    ],
    "price": 1290.00,
    "category": {
      "categoryId": "01935b9c-0000-7000-0000-000000000010",
      "name":       "Running Shoes",
      "slug":       "running-shoes"
    },
    "stockStatus":   "IN_STOCK",
    "reviewSummary": { "averageRating": null, "reviewCount": 0 },
    "status":        "ACTIVE"
  },
  "traceId": "0af7651916cd43dd8448eb211c80319c"
}
```

**`stockStatus` enum:**
- `IN_STOCK` — `availableQty > 10`
- `LOW` — `availableQty` in 1..10
- `OUT_OF_STOCK` — `availableQty == 0` or sku absent from inventory response or inventory unavailable

**`reviewSummary`:** Always `{averageRating: null, reviewCount: 0}` in MVP (Reviews module deferred per EPIC_CATALOG out-of-scope).

## Response (errors)

| Code | HTTP | Trigger |
|---|--:|---|
| `VALIDATION_ERROR` | 400 | `productId` is not a valid UUID |
| `NOT_FOUND` | 404 | Product row absent in DB; OR product status is DELETED/DRAFT and caller is not ADMIN |
| `UPSTREAM_TIMEOUT` | 504 | `inventory.stock.read` times out (> 2s) — logged as WARN; `stockStatus` returned as `OUT_OF_STOCK`; endpoint does NOT return 504 to caller for this path (safe fallback; detail remains functional) |

## Business logic steps

1. Bind and validate `productId` — must be a valid UUID. Return `VALIDATION_ERROR` if malformed.
2. Optionally parse `Authorization` header if present: `common/token.VerifyClaims` (or `ParseUnverifiedClaims` for optional-auth pattern). Extract `claims.role`. If header absent or token invalid, treat caller as guest (no role).
3. Query product + category join: `SELECT p.*, c.id AS cat_id, c.name AS cat_name, c.slug AS cat_slug FROM products p JOIN categories c ON c.id = p.category_id WHERE p.id = $1`.
4. If row not found: return `NOT_FOUND`.
5. Apply status visibility rules:
   - `status IN ('DELETED', 'DRAFT')` AND `claims.role != 'ADMIN'` → return `NOT_FOUND`.
   - `status = 'INACTIVE'` → return product with `status = 'INACTIVE'` (allows cart/checkout to flag `PRODUCT_INACTIVE`).
   - `status = 'ACTIVE'` → return normally.
6. Fetch image array: `SELECT url FROM product_images WHERE product_id = $1 ORDER BY sort_order ASC`.
7. Fetch stock status (best-effort, 2s timeout): `POST http://inventory-svc/api/v1/inventory/stock/read {sku: product.sku}` via `common/httpclient`.
   - `availableQty > 10` → `IN_STOCK`
   - `1..10` → `LOW`
   - `0` or absent → `OUT_OF_STOCK`
   - Timeout / 503: log `WARN`; set `stockStatus = "OUT_OF_STOCK"`; continue.
8. Build `reviewSummary = {averageRating: null, reviewCount: 0}` (fixed placeholder).
9. Return envelope via `common/wrapper.Success`.

## Side effects

- No state mutations (read-only for catalog data).
- One best-effort HTTP call to inventory service for `stockStatus`.

## Idempotency

None — read operation; trivially idempotent.

## Performance

- **p95 target:** < 300ms (BA non_functional latency).
- **Expected QPS at MVP:** 50–200 (product detail views + cart/checkout upstream calls).
- Main query is a PK lookup on `products.id` + category join — sub-millisecond on indexed PK.
- Inventory call adds ~20–50ms best-case; failure path (timeout) capped at 2s but that path returns `OUT_OF_STOCK` without blocking the response with a 504.

## Test cases

- `product_detail_ACTIVE_product_returns_full_record`
- `product_detail_includes_images_ordered_by_sort_order`
- `product_detail_stockStatus_IN_STOCK_when_availableQty_gt_10`
- `product_detail_stockStatus_LOW_when_availableQty_1_to_10`
- `product_detail_stockStatus_OUT_OF_STOCK_when_availableQty_0`
- `product_detail_inventory_timeout_returns_OUT_OF_STOCK_not_504`
- `product_detail_DELETED_product_returns_NOT_FOUND_for_guest`
- `product_detail_DELETED_product_returns_record_for_ADMIN`
- `product_detail_DRAFT_product_returns_NOT_FOUND_for_guest`
- `product_detail_DRAFT_product_returns_record_for_ADMIN`
- `product_detail_INACTIVE_product_returns_record_with_INACTIVE_status`
- `product_detail_absent_productId_returns_NOT_FOUND`
- `product_detail_malformed_UUID_returns_VALIDATION_ERROR`
- `product_detail_reviewSummary_is_always_placeholder`

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | Tech-Designer (Claude Sonnet 4.6, dispatch 2) | Created |
