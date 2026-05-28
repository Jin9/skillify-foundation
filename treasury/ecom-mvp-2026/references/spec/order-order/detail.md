# POST /api/v1/order/order/detail

## Summary

Returns the full order detail including the immutable item snapshots (price/name/image — ORD-007), address snapshot, status, tracking number, and the chronological status history (ORD-009). Honors customer-scope: a CUSTOMER caller can only read their OWN orders; a foreign-id miss returns 404 `ORDER_NOT_OWNED` (no existence leak — ORD-001). ADMIN callers may read any order.

## Story refs

- `STORY_ORDER_LIST_DETAIL`

## Contract ref

[`order.detail`](../../architecture/contracts.json#order.detail)

## Auth

- Tier: **customer_jwt** (single handler; `claims.role` branches between CUSTOMER-scoped and ADMIN-unscoped).
- CUSTOMER: `WHERE id=$1 AND user_id=claims.sub`; on miss → 404 `ORDER_NOT_OWNED` (NEVER 403 — preserves no-existence-leak).
- ADMIN: `WHERE id=$1`; on miss → 404 `NOT_FOUND`.

## Request

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "required": ["orderId"],
  "properties": {
    "orderId": { "type": "string", "format": "uuid" }
  }
}
```

## Response (success)

HTTP 200, envelope:

```json
{
  "code": "SUCCESS",
  "message": "ok",
  "data": {
    "orderId": "01935b9c-...-uuidv7",
    "orderNumber": "ORD-20260508-000123",
    "status": "PACKING",
    "items": [
      {
        "productId": "0192f...",
        "sku": "BLU-T-001-M",
        "productNameSnapshot": "Blue Tee — M",
        "productImageUrlSnapshot": "https://cdn/.../blue-tee-m.webp",
        "priceSnapshot": 490,
        "qty": 2,
        "lineSubtotal": 980
      }
    ],
    "addressSnapshot": {
      "addressId": "0192e...",
      "receiverName": "Somchai Phakdi",
      "phone": "+66 2 123 4567",
      "addressLine": "123/4 Soi Aree",
      "province": "Bangkok",
      "district": "Phaya Thai",
      "postalCode": "10400"
    },
    "subtotal": 980,
    "total": 1040,
    "trackingNumber": null,
    "statusHistory": [
      { "fromStatus": null, "toStatus": "PENDING_PAYMENT", "actorUserId": null, "actorRole": "SYSTEM", "reason": "checkout.commit step 7", "at": "2026-05-08T07:12:33Z" },
      { "fromStatus": "PENDING_PAYMENT", "toStatus": "PAID", "actorUserId": null, "actorRole": "SYSTEM", "reason": "payment.completed eventId=01935...", "at": "2026-05-08T07:14:01Z" },
      { "fromStatus": "PAID", "toStatus": "PACKING", "actorUserId": "0192d-admin", "actorRole": "ADMIN", "reason": null, "at": "2026-05-08T09:00:00Z" }
    ],
    "createdAt": "2026-05-08T07:12:33Z",
    "updatedAt": "2026-05-08T09:00:00Z"
  },
  "traceId": "0af7651916cd43dd8448eb211c80319c"
}
```

## Response (errors)

| Code | HTTP | Trigger |
|---|--:|---|
| `VALIDATION_ERROR` | 400 | `orderId` missing or not UUID-shaped |
| `AUTH_MISSING` | 401 | no `Authorization` header |
| `AUTH_INVALID` | 401 | JWT signature/exp/iss/aud invalid |
| `ORDER_NOT_OWNED` | 404 | CUSTOMER role and (id not found OR id not owned by `claims.sub`) — single error code on both branches to prevent enumeration (ORD-001 no-existence-leak) |
| `NOT_FOUND` | 404 | ADMIN role and `id` does not exist |

## Business logic steps

1. Bind + validate body. UUID format check.
2. Read `claims.sub` and `claims.role` from context.
3. Branch by role:
   - CUSTOMER: `repo.OrdersFindByIdScopedToUser(orderId, claims.sub)`. Miss → 404 `ORDER_NOT_OWNED`.
   - ADMIN: `repo.OrdersFindById(orderId)`. Miss → 404 `NOT_FOUND`.
4. `repo.OrderItemsListByOrderId(orderId)` — returns rows IN INSERTION ORDER (no ORDER BY because `id` is uuidv7 = creation-time-sortable; the API spec does not mandate any particular item order but production sets ORDER BY `created_at ASC`).
5. `repo.StatusHistoryListByOrderIdAsc(orderId)` — backed by `idx_order_status_history_order_occurred`; ORDER BY `occurred_at ASC`.
6. Assemble `OrderDetailResponse` per the schema above. CAT-009: `items[].productNameSnapshot` resolves from the snapshot column even if `catalog.products` was soft-deleted.
7. Wrap with `common/wrapper.Success`.

## Side effects

- None. Read-only.

## Idempotency

- None — naturally idempotent (read).

## Performance

- p95 < 300ms (PERF-003: order detail < 300ms).
- Three indexed queries (PK + FK + indexed FK) → typical < 30ms cold, < 10ms warm.

## Test cases

- `detail_customer_own_order_returns_full_payload`
- `detail_customer_other_users_order_returns_404_ORDER_NOT_OWNED`
- `detail_unknown_id_for_customer_returns_404_ORDER_NOT_OWNED`
- `detail_unknown_id_for_admin_returns_404_NOT_FOUND`
- `detail_admin_any_order_returns_full_payload`
- `detail_status_history_sorted_ascending`
- `detail_resolves_snapshots_after_product_soft_delete` (CAT-009)
- `detail_no_jwt_returns_401`

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | Tech-Designer (Opus 4.7 1M, dry-run #2) | Created |
