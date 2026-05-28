# POST /api/v1/catalog/product/soft-delete

## Summary

Admin-only endpoint that soft-deletes a product by setting its `status` to `'DELETED'`. Does NOT physically remove any rows — the product row, its images, and its status history are all preserved. Once soft-deleted, the product is excluded from storefront listing and returns `NOT_FOUND` to non-admin callers on `product.detail`. Historical orders that reference this product are unaffected because the Order service stores `productNameSnapshot`, `productImageUrlSnapshot`, and `priceSnapshot` on `order_items` at checkout time and never calls Catalog post-snapshot (CAT-009). The operation is idempotent: re-deleting an already-DELETED product returns 200 with `status=DELETED`. No outbox event emitted for soft-delete. Implements requirement CAT-009.

## Story refs

- `STORY_CATALOG_ADMIN_PRODUCT_CRUD`

## Contract ref

[`catalog.product.soft-delete`](../../architecture/contracts.json#catalog.product.soft-delete)

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
    "productId": { "type": "string", "format": "uuid" }
  }
}
```

## Response (success)

HTTP 200, envelope:

```json
{
  "code": "SUCCESS",
  "message": "Product deleted",
  "data": {
    "productId": "01935b9c-0000-7000-0000-000000000001",
    "status":    "DELETED"
  },
  "traceId": "0af7651916cd43dd8448eb211c80319c"
}
```

## Response (errors)

| Code | HTTP | Trigger |
|---|--:|---|
| `AUTH_MISSING` | 401 | No `Authorization` header |
| `AUTH_FORBIDDEN` | 403 | `claims.role = CUSTOMER` |
| `NOT_FOUND` | 404 | `productId` row does not exist in `products` table (never-existed product) |

**Note:** If the product exists but is already `status='DELETED'`, this returns HTTP 200 with `status=DELETED` (idempotent — NOT a 404 or 409).

## Business logic steps

1. JWT middleware validates token; injects `claims`. Returns `AUTH_MISSING` / `AUTH_FORBIDDEN` before handler.
2. Bind and validate `productId` — must be valid UUID. Return `VALIDATION_ERROR` if malformed.
3. Idempotency check: `SELECT status FROM products WHERE id = $1`.
   - Row not found: return `NOT_FOUND (404)`.
   - `status = 'DELETED'`: return `200 SUCCESS {productId, status: "DELETED"}` immediately — no writes.
4. Begin transaction (`pgx`):
   a. `UPDATE products SET status='DELETED', updated_at=NOW() WHERE id=$1 AND status != 'DELETED'`.
   b. Check `rowsAffected`. If 0 (concurrent soft-delete won the race): treat as idempotent success.
   c. `INSERT INTO product_status_history (product_id, from_status=<previously-read status>, to_status='DELETED', changed_field=NULL, old_value=NULL, new_value=NULL, actor_user_id=claims.sub)`.
   d. `COMMIT`.
5. Return `200 SUCCESS {productId, status: "DELETED"}`.

**No outbox event** is emitted for soft-delete — no consumer in the defined contracts requires it. Inventory does not need to react; Order uses snapshot fields.

## Side effects

- UPDATE one row in `catalog.products` (`status='DELETED'`, `updated_at`).
- INSERT one row in `catalog.product_status_history` (from_status=prior, to_status='DELETED').
- Both in a single transaction.
- Storefront effect: `product.list` excludes this product immediately (partial indexes filter `status != 'ACTIVE'`). `product.detail` returns `NOT_FOUND` to non-admin callers.
- Historical orders: unaffected — Order snapshots are self-contained.

## Idempotency

**Idempotent by design.** Re-calling with an already-DELETED `productId`:
- Returns `200 SUCCESS {productId, status: "DELETED"}`.
- Writes NO additional `product_status_history` row.
- Writes NO outbox event.

## Performance

- **p95 target:** < 300ms (admin write).
- **Expected QPS at MVP:** < 1 (infrequent catalog management operation).
- Simple PK lookup + conditional UPDATE; sub-millisecond.

## Test cases

- `product_soft_delete_changes_status_to_DELETED`
- `product_soft_delete_idempotent_returns_200_no_extra_history_row`
- `product_soft_delete_excludes_product_from_listing`
- `product_soft_delete_product_detail_returns_NOT_FOUND_for_guest`
- `product_soft_delete_product_detail_returns_record_for_ADMIN`
- `product_soft_delete_absent_productId_returns_404`
- `product_soft_delete_CUSTOMER_token_returns_403`
- `product_soft_delete_no_token_returns_401`
- `product_soft_delete_writes_status_history_row_from_prior_status`
- `product_soft_delete_no_outbox_event_written`
- `product_soft_delete_historical_order_snapshot_unaffected`

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | Tech-Designer (Claude Sonnet 4.6, dispatch 2) | Created |
