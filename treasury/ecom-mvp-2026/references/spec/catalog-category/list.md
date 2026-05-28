# POST /api/v1/catalog/category/list

## Summary

Public endpoint that returns a flat list of categories, ordered by `level ASC, name ASC`. Drives the category filter chips on the customer mobile webview. By default (`activeOnly=true`) returns only active categories; callers may pass `activeOnly=false` to include inactive categories. No authentication is required or enforced — the endpoint is registered outside the JWT middleware group. See AMB-CAT-002 for the access-control ambiguity on `activeOnly=false` (MVP: no gate; post-MVP: optional ADMIN check). Implements requirement CATE-003.

## Story refs

- `STORY_CATALOG_ADMIN_CATEGORY`
- `STORY_CATALOG_LIST` (category filter chips depend on this endpoint)

## Contract ref

[`catalog.category.list`](../../architecture/contracts.json#catalog.category.list)

## Auth

- Tier: **none** (public endpoint — registered outside JWT middleware group).
- No `Authorization` header required or enforced.
- Known limitation: `activeOnly=false` is not gated to ADMIN in MVP. Any caller may pass it to see inactive categories. Documented in AMB-CAT-002; post-MVP mitigation is an optional-auth middleware with ADMIN check for `activeOnly=false`.

## Request

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "activeOnly": {
      "type": "boolean",
      "default": true,
      "description": "true = return only active categories (storefront default); false = return all including inactive"
    }
  }
}
```

**Notes:**
- Omitting `activeOnly` is equivalent to `activeOnly=true`.
- No pagination parameters — MVP returns the full list (category counts at MVP scale are expected to be small; < 200 rows).

## Response (success)

HTTP 200, envelope:

```json
{
  "code": "SUCCESS",
  "message": "OK",
  "data": {
    "categories": [
      {
        "categoryId":       "01935b9c-0000-7000-0000-000000000001",
        "name":             "Footwear",
        "slug":             "footwear",
        "parentCategoryId": null,
        "level":            1,
        "active":           true
      },
      {
        "categoryId":       "01935b9c-0000-7000-0000-000000000010",
        "name":             "Running Shoes",
        "slug":             "running-shoes",
        "parentCategoryId": "01935b9c-0000-7000-0000-000000000001",
        "level":            2,
        "active":           true
      }
    ]
  },
  "traceId": "0af7651916cd43dd8448eb211c80319c"
}
```

**Each category item fields:**
- `categoryId` — UUID
- `name` — display name
- `slug` — URL-safe identifier
- `parentCategoryId` — UUID or `null` for root-level categories
- `level` — depth hint (1 = root)
- `active` — boolean visibility flag

**Ordering:** `level ASC, name ASC` — root categories first, then children, alphabetically within each level.

**Tree reconstruction:** Returned as a flat list. Client reconstructs the tree using `parentCategoryId`. No recursive CTE required server-side in MVP.

## Response (errors)

| Code | HTTP | Trigger |
|---|--:|---|
| `DATABASE_UNAVAILABLE` | 503 | Postgres unreachable |

No `VALIDATION_ERROR` — the request schema has no required fields and `activeOnly` is a simple boolean with a default.

## Business logic steps

1. Bind body (`common/validator`). Apply default: `activeOnly=true` if omitted.
2. Execute query: `SELECT id, name, slug, parent_category_id, level, active FROM categories WHERE ($1::boolean IS NULL OR active = $1) ORDER BY level ASC, name ASC`.
   - `activeOnly=true` → `$1 = TRUE` → `active = TRUE`
   - `activeOnly=false` → `$1 = FALSE` → no active filter (all categories)
3. Map rows to response items.
4. Return `SUCCESS` envelope with `categories` array (empty array if no categories found — NOT an error).

## Side effects

- None (read-only).

## Idempotency

Read — trivially idempotent.

## Performance

- **p95 target:** < 100ms (lightweight read; no joins; small table).
- **Expected QPS at MVP:** 50–200 (storefront filter chip population on every browse session load).
- Single table scan with `INDEX (active, level)` covering the default `activeOnly=true` path.

## Test cases

- `category_list_activeOnly_true_returns_only_active_categories`
- `category_list_activeOnly_false_returns_all_including_inactive`
- `category_list_default_activeOnly_true_when_omitted`
- `category_list_ordered_by_level_asc_name_asc`
- `category_list_empty_returns_200_empty_array`
- `category_list_deactivated_category_absent_from_activeOnly_true`
- `category_list_parentCategoryId_null_for_root_categories`

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | Tech-Designer (Claude Sonnet 4.6, dispatch 2) | Created |
