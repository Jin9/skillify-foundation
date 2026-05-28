# POST /api/v1/catalog/category/update

## Summary

Admin-only endpoint that updates the `name` and/or `slug` of an existing category. Applies a partial patch: at least one of `name` or `slug` must be provided. Name uniqueness within the same level/parent and global slug uniqueness are both re-enforced on update. Inactive categories can still be updated. Implements requirement CATE-002.

## Story refs

- `STORY_CATALOG_ADMIN_CATEGORY`

## Contract ref

[`catalog.category.update`](../../architecture/contracts.json#catalog.category.update)

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
    "categoryId": {
      "type": "string",
      "format": "uuid"
    },
    "name": {
      "type": "string",
      "minLength": 1,
      "maxLength": 255
    },
    "slug": {
      "type": "string",
      "minLength": 1,
      "maxLength": 255
    }
  }
}
```

**Notes:**
- At least one of `name` or `slug` must be provided. Empty patch (only `categoryId`) returns `VALIDATION_ERROR`.
- `name` must remain unique within the category's current parent level.
- `slug` must remain globally unique.
- `parentCategoryId` and `level` are not editable via this endpoint in MVP.

## Response (success)

HTTP 200, envelope:

```json
{
  "code": "SUCCESS",
  "message": "Category updated",
  "data": {
    "categoryId": "01935b9c-0000-7000-0000-000000000010",
    "name":       "Running & Trail Shoes",
    "slug":       "running-trail-shoes",
    "active":     true,
    "updatedAt":  "2026-05-08T11:00:00Z"
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
| `VALIDATION_ERROR` | 400 | Empty patch (no `name` or `slug` provided); malformed UUID in `categoryId` |
| `DUPLICATE_CATEGORY_NAME` | 409 | New `name` conflicts with a sibling category at the same level (CATE-001) |
| `CONFLICT` | 409 | New `slug` already taken globally (CATE-002) |

## Business logic steps

1. JWT middleware validates token; injects `claims`. Returns `AUTH_MISSING` / `AUTH_FORBIDDEN` before handler.
2. Bind and validate body. Return `VALIDATION_ERROR` if neither `name` nor `slug` provided (empty patch).
3. Verify category exists: `SELECT id, name, slug, parent_category_id, level, active FROM categories WHERE id = $1`. If not found: return `NOT_FOUND`.
4. Apply update: `UPDATE categories SET name = COALESCE($newName, name), slug = COALESCE($newSlug, slug), updated_at = NOW() WHERE id = $1 RETURNING id, name, slug, active, updated_at`.
5. On DB unique-constraint violation:
   - `categories_parent_category_id_name_key` OR `categories_name_parent_null_key` → return `DUPLICATE_CATEGORY_NAME (409)`.
   - `categories_slug_key` → return `CONFLICT (409)`.
6. Return `SUCCESS` envelope with updated category fields.

## Side effects

- UPDATE one row in `catalog.categories`.
- No outbox event emitted for category update.

## Idempotency

Last-writer-wins on `name` and `slug`. Re-submitting the same values returns 200 with no visible change.

## Performance

- **p95 target:** < 300ms (admin write).
- **Expected QPS at MVP:** < 1.
- Simple PK lookup + COALESCE UPDATE; sub-millisecond.

## Test cases

- `category_update_name_returns_200_with_new_name`
- `category_update_slug_returns_200_with_new_slug`
- `category_update_both_name_and_slug_returns_200`
- `category_update_duplicate_name_same_parent_returns_409_DUPLICATE_CATEGORY_NAME`
- `category_update_duplicate_slug_returns_409_CONFLICT`
- `category_update_absent_categoryId_returns_404`
- `category_update_empty_patch_returns_400_VALIDATION_ERROR`
- `category_update_inactive_category_can_be_updated`
- `category_update_CUSTOMER_token_returns_403`
- `category_update_no_token_returns_401`

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | Tech-Designer (Claude Sonnet 4.6, dispatch 2) | Created |
