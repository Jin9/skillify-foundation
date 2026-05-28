# POST /api/v1/catalog/category/deactivate

## Summary

Admin-only endpoint that deactivates a category by setting `active = FALSE`. Products belonging to this category are NOT automatically hidden or soft-deleted — they retain their own `status` field. The deactivated category will no longer appear in `category.list` when `activeOnly=true` (the storefront default), and it will no longer be accepted as a valid `categoryId` when creating new products. The operation is idempotent: deactivating an already-inactive category returns 200 with `active = false`. Implements requirement CATE-003.

## Story refs

- `STORY_CATALOG_ADMIN_CATEGORY`

## Contract ref

[`catalog.category.deactivate`](../../architecture/contracts.json#catalog.category.deactivate)

## Auth

- Tier: **required — ADMIN role**.
- `Authorization: Bearer <jwt>` required.
- `common/middleware` JWT middleware validates before handler runs.
- `claims.role != 'ADMIN'`: return `AUTH_FORBIDDEN (403)`.
- Header absent: return `AUTH_MISSING (401)`.

## Request

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "required": ["categoryId"],
  "properties": {
    "categoryId": { "type": "string", "format": "uuid" }
  }
}
```

## Response (success)

HTTP 200, envelope:

```json
{
  "code": "SUCCESS",
  "message": "Category deactivated",
  "data": {
    "categoryId": "01935b9c-0000-7000-0000-000000000010",
    "active":     false
  },
  "traceId": "0af7651916cd43dd8448eb211c80319c"
}
```

## Response (errors)

| Code | HTTP | Trigger |
|---|--:|---|
| `AUTH_MISSING` | 401 | No `Authorization` header |
| `AUTH_FORBIDDEN` | 403 | `claims.role = CUSTOMER` |
| `NOT_FOUND` | 404 | `categoryId` does not exist in `categories` table |

**Note:** If the category already has `active = false`, this returns HTTP 200 with `active = false` (idempotent — NOT a 404 or 409).

## Business logic steps

1. JWT middleware validates token; injects `claims`. Returns `AUTH_MISSING` / `AUTH_FORBIDDEN` before handler.
2. Bind and validate `categoryId` — must be valid UUID. Return `VALIDATION_ERROR` if malformed.
3. Verify category exists: `SELECT id, active FROM categories WHERE id = $1`. If not found: return `NOT_FOUND`.
4. If `active = false` already: return `200 SUCCESS {categoryId, active: false}` immediately — no write.
5. `UPDATE categories SET active = FALSE, updated_at = NOW() WHERE id = $1`.
6. Return `200 SUCCESS {categoryId, active: false}`.

## Side effects

- UPDATE one row in `catalog.categories` (`active=FALSE`, `updated_at`).
- No cascade to products — products in this category retain their own status.
- Storefront effect: `category.list?activeOnly=true` no longer returns this category. Products in this category still appear in `product.list` if `status=ACTIVE` AND `visible=TRUE`.
- New product creates with this `categoryId` will return `NOT_FOUND` (category no longer `active=true`).

## Idempotency

**Idempotent by design.** Re-calling with an already-inactive `categoryId` returns 200 with `active=false` and writes nothing.

## Performance

- **p95 target:** < 300ms (admin write).
- **Expected QPS at MVP:** < 1.
- Simple PK lookup + conditional UPDATE; sub-millisecond.

## Test cases

- `category_deactivate_active_category_sets_active_false`
- `category_deactivate_idempotent_returns_200_no_write`
- `category_deactivate_absent_categoryId_returns_404`
- `category_deactivate_products_in_category_remain_visible`
- `category_deactivate_category_absent_from_activeOnly_list`
- `category_deactivate_category_still_in_all_categories_list`
- `category_deactivate_CUSTOMER_token_returns_403`
- `category_deactivate_no_token_returns_401`

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | Tech-Designer (Claude Sonnet 4.6, dispatch 2) | Created |
