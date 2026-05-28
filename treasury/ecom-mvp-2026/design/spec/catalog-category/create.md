# POST /api/v1/catalog/category/create

## Summary

Admin-only endpoint that creates a new product category in the catalog taxonomy. Categories support a self-referencing tree structure (via `parentCategoryId`). Name uniqueness is enforced within the same level/parent: two siblings (same parent) cannot share a name, and two root-level categories cannot share a name. Slug is globally unique across all categories. Implements requirement CATE-001, CATE-002.

## Story refs

- `STORY_CATALOG_ADMIN_CATEGORY`

## Contract ref

[`catalog.category.create`](../../architecture/contracts.json#catalog.category.create)

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
  "required": ["name", "slug", "level"],
  "properties": {
    "name": {
      "type": "string",
      "minLength": 1,
      "maxLength": 255
    },
    "slug": {
      "type": "string",
      "minLength": 1,
      "maxLength": 255,
      "description": "URL-safe identifier; globally unique; lowercase-alphanumeric-hyphen recommended"
    },
    "parentCategoryId": {
      "type": "string",
      "format": "uuid",
      "description": "Omit or null for root-level category"
    },
    "level": {
      "type": "integer",
      "minimum": 1,
      "description": "Depth hint; informational in MVP; no strict depth-vs-parent enforcement"
    }
  }
}
```

**Notes:**
- `parentCategoryId` absent/null = root-level category.
- `level` is informational; the server does not validate level against the parent's actual depth in MVP.
- `slug` should be globally unique and URL-safe. Recommended format: `lowercase-alphanumeric-hyphen`.

## Response (success)

HTTP 201, envelope:

```json
{
  "code": "CREATED",
  "message": "Category created",
  "data": {
    "categoryId": "01935b9c-0000-7000-0000-000000000010"
  },
  "traceId": "0af7651916cd43dd8448eb211c80319c"
}
```

## Response (errors)

| Code | HTTP | Trigger |
|---|--:|---|
| `AUTH_MISSING` | 401 | No `Authorization` header |
| `AUTH_FORBIDDEN` | 403 | `claims.role = CUSTOMER` |
| `VALIDATION_ERROR` | 400 | Missing required fields; empty `name` or `slug`; `level < 1`; malformed UUID in `parentCategoryId` |
| `DUPLICATE_CATEGORY_NAME` | 409 | A category with the same `name` already exists at the same level (same `parentCategoryId` or both root-level) — DB unique-constraint violation (CATE-001) |
| `CONFLICT` | 409 | `slug` already exists globally across all categories — DB unique-constraint violation (CATE-002) |
| `NOT_FOUND` | 404 | `parentCategoryId` provided but does not exist in `categories` table |

## Business logic steps

1. JWT middleware validates token; injects `claims`. Returns `AUTH_MISSING` / `AUTH_FORBIDDEN` before handler.
2. Bind and validate body (`common/validator`). Return `VALIDATION_ERROR` on schema violations.
3. If `parentCategoryId` provided: verify it exists: `SELECT id FROM categories WHERE id = $1`. If not found: return `NOT_FOUND (404)`. (Note: parent may be active or inactive — no active requirement for parent in MVP.)
4. Generate `categoryId = uuid_v7()` in Go.
5. `INSERT INTO categories (id, name, slug, parent_category_id, level, active, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, TRUE, NOW(), NOW())`.
6. On DB unique-constraint violation:
   - Constraint `categories_parent_category_id_name_key` OR `categories_name_parent_null_key` → return `DUPLICATE_CATEGORY_NAME (409)`.
   - Constraint `categories_slug_key` → return `CONFLICT (409)`.
7. Return `CREATED` envelope with `{categoryId}`.

## Side effects

- INSERT one row into `catalog.categories`.
- No outbox event emitted for category creation.

## Idempotency

None — category create is not idempotent. Retrying with the same name at the same level returns `DUPLICATE_CATEGORY_NAME (409)`. Retrying with the same slug returns `CONFLICT (409)`.

## Performance

- **p95 target:** < 300ms (admin write).
- **Expected QPS at MVP:** < 1.
- Simple INSERT with constraint check; sub-millisecond.

## Test cases

- `category_create_valid_root_level_returns_201`
- `category_create_valid_child_category_returns_201`
- `category_create_duplicate_name_same_parent_returns_409_DUPLICATE_CATEGORY_NAME`
- `category_create_duplicate_name_root_level_returns_409_DUPLICATE_CATEGORY_NAME`
- `category_create_duplicate_name_different_parent_succeeds`
- `category_create_duplicate_slug_returns_409_CONFLICT`
- `category_create_absent_parentCategoryId_returns_404`
- `category_create_CUSTOMER_token_returns_403`
- `category_create_no_token_returns_401`
- `category_create_missing_name_returns_400_VALIDATION_ERROR`
- `category_create_missing_slug_returns_400_VALIDATION_ERROR`
- `category_create_level_zero_returns_400_VALIDATION_ERROR`

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | Tech-Designer (Claude Sonnet 4.6, dispatch 2) | Created |
