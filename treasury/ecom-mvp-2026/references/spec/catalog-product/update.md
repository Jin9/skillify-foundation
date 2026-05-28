# POST /api/v1/catalog/product/update

## Summary

Admin-only endpoint that applies a partial update to an existing product. Supports updating name, description, images (full replacement), price, status (DRAFT, ACTIVE, INACTIVE — not DELETED), and category. SKU is not editable in MVP. Every update writes an audit row to `product_status_history` and emits a `product.updated` outbox event (no consumer in MVP — written for future extensibility). A `SELECT ... FOR UPDATE` read-then-write pattern prevents lost updates in concurrent admin sessions. Implements requirement CAT-008, CAT-010.

## Story refs

- `STORY_CATALOG_ADMIN_PRODUCT_CRUD`

## Contract ref

[`catalog.product.update`](../../architecture/contracts.json#catalog.product.update)

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
  "required": ["productId"],
  "properties": {
    "productId": {
      "type": "string",
      "format": "uuid"
    },
    "name": {
      "type": "string",
      "minLength": 1,
      "maxLength": 255
    },
    "description": {
      "type": "string"
    },
    "images": {
      "type": "array",
      "items": { "type": "string", "minLength": 1, "maxLength": 2048 },
      "maxItems": 20
    },
    "price": {
      "type": "number",
      "exclusiveMinimum": 0
    },
    "status": {
      "type": "string",
      "enum": ["DRAFT", "ACTIVE", "INACTIVE"]
    },
    "categoryId": {
      "type": "string",
      "format": "uuid"
    }
  }
}
```

**Notes:**
- At least one of `name`, `description`, `images`, `price`, `status`, `categoryId` must be present. Empty patch returns `VALIDATION_ERROR`.
- `status = 'DELETED'` is forbidden via update — use `product.soft-delete` instead.
- `sku` is NOT in this schema — SKU is locked/immutable in MVP.
- If `images` is present, it **replaces** the entire image set (DELETE + re-INSERT inside same tx).
- `price` must be > 0 if provided.
- `categoryId` if provided must reference an existing active category.

## Response (success)

HTTP 200, envelope:

```json
{
  "code": "SUCCESS",
  "message": "Product updated",
  "data": {
    "productId":   "01935b9c-0000-7000-0000-000000000001",
    "sku":         "SHOE-RED-42",
    "name":        "Running Shoe Red 42 (v2)",
    "description": "Updated description.",
    "images": [
      "https://cdn.example.com/shoes/red-42/01-v2.jpg"
    ],
    "price": 1390.00,
    "category": {
      "categoryId": "01935b9c-0000-7000-0000-000000000010",
      "name":       "Running Shoes",
      "slug":       "running-shoes"
    },
    "status":    "ACTIVE",
    "updatedAt": "2026-05-08T10:30:00Z"
  },
  "traceId": "0af7651916cd43dd8448eb211c80319c"
}
```

## Response (errors)

| Code | HTTP | Trigger |
|---|--:|---|
| `AUTH_MISSING` | 401 | No `Authorization` header |
| `AUTH_FORBIDDEN` | 403 | `claims.role = CUSTOMER` |
| `NOT_FOUND` | 404 | `productId` not in DB; OR product `status = 'DELETED'` (treat DELETED as non-existent for update) |
| `VALIDATION_ERROR` | 400 | Empty patch (no optional fields provided); `price <= 0`; `status = 'DELETED'` in request body; malformed UUID; `images` > 20 items; `categoryId` inactive or absent |

## Business logic steps

1. JWT middleware validates token; injects `claims`. Returns `AUTH_MISSING` / `AUTH_FORBIDDEN` before handler.
2. Bind and validate body (`common/validator`). Return `VALIDATION_ERROR` if schema violated or patch is empty.
3. Begin transaction (`pgx`):
   a. `SELECT id, sku, name, description, price, category_id, status, visible FROM products WHERE id = $1 FOR UPDATE`.
   b. If row not found OR `status = 'DELETED'`: ROLLBACK; return `NOT_FOUND`.
   c. If `categoryId` provided: `SELECT id FROM categories WHERE id = $newCategoryId AND active = TRUE`. If not found: ROLLBACK; return `VALIDATION_ERROR` (category inactive or absent).
   d. Apply patch: compute new values for each provided field; keep old value for audit `old_value` capture.
   e. `UPDATE products SET name=..., description=..., price=..., category_id=..., status=..., updated_at=NOW() WHERE id=$1`.
   f. If `images` provided: `DELETE FROM product_images WHERE product_id = $1`; batch `INSERT INTO product_images (id, product_id, url, sort_order)`.
   g. `INSERT INTO product_status_history (product_id, from_status=oldStatus, to_status=newStatus, changed_field=NULL, old_value=<JSON of old fields>, new_value=<JSON of new fields>, actor_user_id=claims.sub)`.
   h. `INSERT INTO outbox_events (id=uuid_v7(), aggregate_id=productId, event_type='product.updated', payload_json={eventId, occurredAt, productId, sku, changedFields:[list]}, published_at=NULL)`.
   i. `COMMIT`.
4. Fetch updated product + images + category for response (or use RETURNING clause).
5. Return `SUCCESS` envelope.

## Side effects

- UPDATE one row in `catalog.products`.
- If `images` provided: DELETE all existing `product_images` for this product + INSERT new set.
- INSERT one row into `catalog.product_status_history`.
- INSERT one row into `catalog.outbox_events` (`event_type='product.updated'`; no consumer in MVP).
- All writes in a single transaction.
- If `status` changed to `INACTIVE`: product is immediately hidden from storefront listing; `cart.read` will flag affected cart items as non-checkoutable on next call.

## Idempotency

Last-writer-wins. No `Idempotency-Key` header. Concurrent admin updates to the same product are serialized by the `SELECT ... FOR UPDATE` row lock.

## Performance

- **p95 target:** < 300ms (admin write).
- **Expected QPS at MVP:** < 5.
- Read-then-write pattern: one SELECT FOR UPDATE + conditional INSERTs/UPDATE — all sub-millisecond on PK lookups.

## Test cases

- `product_update_name_and_price_returns_200_with_updated_fields`
- `product_update_status_to_INACTIVE_hides_product_from_storefront`
- `product_update_status_DELETED_in_request_returns_400_VALIDATION_ERROR`
- `product_update_empty_patch_returns_400_VALIDATION_ERROR`
- `product_update_absent_productId_returns_404`
- `product_update_DELETED_product_returns_404`
- `product_update_images_replaces_entire_image_set`
- `product_update_categoryId_inactive_returns_400_VALIDATION_ERROR`
- `product_update_writes_product_status_history_row`
- `product_update_outbox_event_product_updated_written`
- `product_update_CUSTOMER_token_returns_403`
- `product_update_no_token_returns_401`
- `product_update_price_zero_returns_400_VALIDATION_ERROR`
- `product_update_concurrent_writes_last_writer_wins`

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | Tech-Designer (Claude Sonnet 4.6, dispatch 2) | Created |
