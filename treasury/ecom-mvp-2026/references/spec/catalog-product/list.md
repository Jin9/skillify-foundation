# POST /api/v1/catalog/product/list

## Summary

Public storefront endpoint that returns a paginated, filterable, sortable list of ACTIVE+visible products. Drives the home tab and search tab of the customer mobile webview. Supports filters for category, price range, in-stock status, and a name keyword search. When `inStock=true` is requested, catalog makes a synchronous call to `inventory.stock.bulk-read` to filter out out-of-stock products server-side. Implements requirement CAT-001 through CAT-005.

## Story refs

- `STORY_CATALOG_LIST`

## Contract ref

[`catalog.product.list`](../../architecture/contracts.json#catalog.product.list)

## Auth

- Tier: **none** (public endpoint — registered outside JWT middleware group).
- No `Authorization` header required or validated.

## Request

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "page":  { "type": "integer", "minimum": 1, "default": 1 },
    "limit": { "type": "integer", "minimum": 1, "maximum": 100, "default": 20 },
    "sort":  {
      "type": "string",
      "enum": ["price_asc", "price_desc", "name_asc", "name_desc", "created_desc"],
      "default": "created_desc"
    },
    "filter": {
      "type": "object",
      "additionalProperties": false,
      "properties": {
        "categoryId": { "type": "string", "format": "uuid" },
        "minPrice":   { "type": "number", "minimum": 0 },
        "maxPrice":   { "type": "number", "minimum": 0 },
        "inStock":    { "type": "boolean" },
        "search":     { "type": "string", "maxLength": 200 }
      }
    }
  }
}
```

**Notes:**
- `maxPrice` must be >= `minPrice` when both are supplied.
- `search` is trimmed; ignored if empty after trim.
- Omitting `page`/`limit`/`sort` uses the stated defaults.

## Response (success)

HTTP 200, envelope:

```json
{
  "code": "SUCCESS",
  "message": "OK",
  "data": {
    "items": [
      {
        "productId":    "01935b9c-0000-7000-0000-000000000001",
        "sku":          "SHOE-RED-42",
        "name":         "Running Shoe Red 42",
        "price":        1290.00,
        "thumbnailUrl": "https://cdn.example.com/shoes/red-42/01.jpg",
        "category": {
          "categoryId": "01935b9c-0000-7000-0000-000000000010",
          "name":       "Running Shoes",
          "slug":       "running-shoes"
        },
        "status": "ACTIVE"
      }
    ],
    "total": 1000,
    "page":  1,
    "limit": 20
  },
  "traceId": "0af7651916cd43dd8448eb211c80319c"
}
```

**`items` element fields:**
- `productId` — UUID v7
- `sku` — stock-keeping unit string
- `name` — product display name
- `price` — NUMERIC as decimal; THB
- `thumbnailUrl` — first image URL ordered by `sort_order ASC`; `null` if no images
- `category` — `{categoryId, name, slug}`
- `status` — always `"ACTIVE"` in storefront listing

## Response (errors)

| Code | HTTP | Trigger |
|---|--:|---|
| `VALIDATION_ERROR` | 400 | `page < 1`, `limit` not in 1..100, unknown `sort` enum value, malformed UUID in `categoryId`, `maxPrice < minPrice`, `search` > 200 chars |
| `UPSTREAM_TIMEOUT` | 504 | `inStock=true` AND `inventory.stock.bulk-read` returns 503 or times out (> 2s); catalog does NOT silently return an unfiltered list |

## Business logic steps

1. Bind and validate request body against the schema above (`common/validator`). Return `VALIDATION_ERROR` on failure.
2. Apply defaults: `page=1`, `limit=20`, `sort=created_desc` if omitted.
3. Build a single parameterized query: `WHERE status='ACTIVE' AND visible=TRUE` plus conditional clauses for `categoryId`, `minPrice`, `maxPrice`, and `search` (ILIKE `'%' || $n || '%'`). No dynamic SQL string concatenation.
4. Run COUNT query with the same WHERE clause to compute `total` (separate query or window function).
5. Run paged query with `ORDER BY (sort_col ASC|DESC, id ASC)` for stable pagination, `LIMIT $limit OFFSET ($page-1)*$limit`.
6. For each product row, fetch its first image via a LEFT JOIN or subquery on `product_images` ordered by `sort_order ASC`.
7. If `inStock=true`:
   a. Extract `skus` from the candidate product rows.
   b. Call `POST http://inventory-svc/api/v1/inventory/stock/bulk-read {skus:[...]}` via `common/httpclient` with a 2s timeout.
   c. On success: build `map[sku]availableQty`; filter out products where `availableQty == 0` or sku is absent from the inventory response.
   d. On timeout or 503: return `UPSTREAM_TIMEOUT (504)` — do not silently return the unfiltered list.
8. Wrap filtered items in the response envelope via `common/wrapper.Success`.

## Side effects

- No state mutations (read-only).
- If `inStock=true`: one synchronous HTTP call to inventory service.

## Idempotency

None — read operation; trivially idempotent.

## Performance

- **p95 target:** < 300ms on a 1,000-SKU dataset in local Docker Compose (PERF-001).
- **Expected QPS at MVP:** 20–100 (storefront browse + search).
- Partial indexes `idx_products_listing` (category filter) and `idx_products_price` (price range) reduce scan to ACTIVE+visible rows only. Validate with `EXPLAIN ANALYZE` after seeding 1k rows.
- `inStock=true` adds one inventory HTTP call (~20–50ms on same cluster); still within the 300ms budget.

## Test cases

- `product_list_returns_20_items_page1_sorted_price_asc`
- `product_list_filters_by_categoryId`
- `product_list_filters_by_price_range`
- `product_list_inStock_true_filters_out_zero_qty_items`
- `product_list_inStock_true_inventory_timeout_returns_504`
- `product_list_excludes_DRAFT_INACTIVE_DELETED_products`
- `product_list_empty_catalog_returns_200_items_empty_total_0`
- `product_list_limit_greater_than_100_returns_VALIDATION_ERROR`
- `product_list_invalid_sort_returns_VALIDATION_ERROR`
- `product_list_maxPrice_less_than_minPrice_returns_VALIDATION_ERROR`
- `product_list_categoryId_malformed_UUID_returns_VALIDATION_ERROR`
- `product_list_search_keyword_matches_name_case_insensitive`
- `product_list_stable_pagination_no_duplicate_items_across_pages`

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | Tech-Designer (Claude Sonnet 4.6, dispatch 2) | Created |
