# POST /api/v1/order/order/list-mine

## Summary

Customer-facing list of the requesting customer's own orders, sorted by `createdAt DESC`. Optionally filtered by `status`. Returns paged `OrderSummary[]` with `total` and `page`. Server enforces `customerUserId = claims.sub` — listing of someone else's orders is impossible by construction (ORD-001).

## Story refs

- `STORY_ORDER_LIST_DETAIL`

## Contract ref

[`order.list-mine`](../../architecture/contracts.json#order.list-mine)

## Auth

- Tier: **customer_jwt** (CUSTOMER role; ADMIN tokens are accepted for listing their OWN orders if any but the use case is admins reading customers' orders via `order.list-admin` instead — Reviewer-L1 may tighten to CUSTOMER-only if desired).
- Header: `Authorization: Bearer <jwt>`. Backend validates via `common/middleware` JWT (ES256, iss/aud/exp) and injects `token.Claims` into context.
- Server enforces `WHERE user_id = claims.sub`. NO client-supplied user filter is honored.

## Request

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "page":   { "type": "integer", "minimum": 1, "default": 1 },
    "limit":  { "type": "integer", "minimum": 1, "maximum": 100, "default": 20 },
    "status": {
      "type": "string",
      "enum": ["PENDING_PAYMENT","PAID","PAYMENT_FAILED","PAYMENT_EXPIRED","PACKING","SHIPPED","DELIVERED","CANCELLED"]
    }
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
    "orders": [
      {
        "orderId": "01935b9c-...-uuidv7",
        "orderNumber": "ORD-20260508-000123",
        "status": "PAID",
        "subtotal": 980,
        "shippingFee": 60,
        "total": 1040,
        "itemCount": 3,
        "createdAt": "2026-05-08T07:12:33Z"
      }
    ],
    "total": 17,
    "page": 1
  },
  "traceId": "0af7651916cd43dd8448eb211c80319c"
}
```

## Response (errors)

| Code | HTTP | Trigger |
|---|--:|---|
| `VALIDATION_ERROR` | 400 | `page`/`limit` out of bounds; `status` not in the 8-value enum |
| `AUTH_MISSING` | 401 | no `Authorization` header on a protected endpoint |
| `AUTH_INVALID` | 401 | JWT signature/exp/iss/aud invalid |

## Business logic steps

1. Bind + validate request body (`common/validator`); coerce `page` default 1, `limit` default 20.
2. Read `claims.sub` from context (set by `common/middleware` JWT middleware).
3. `repo.OrdersListByUser(ctx, userId=claims.sub, status=req.status, page, limit)` — backed by index `idx_orders_user_id_created_at_desc`. ORDER BY `created_at DESC`.
4. `repo.OrdersCountByUser(ctx, userId, status)` — for `total`. (Single COUNT(*) on the same predicate.)
5. Map `[]Order` -> `[]OrderSummary` (drop snapshot/internal columns).
6. Wrap with `common/wrapper.Success(ctx, data)`.

## Side effects

- None. Read-only.

## Idempotency

- None — operation is naturally idempotent (read).

## Performance

- p95 < 300ms (BA `non_functional.latency`; PERF-003 order detail < 300ms; this list query is a covering-index scan on `(user_id, created_at)`).
- Expected QPS at MVP: 5–10 (mobile webview "Orders" tab).

## Test cases

- `list_mine_returns_only_own_orders_sorted_desc`
- `list_mine_status_filter_active_excludes_DELIVERED`
- `list_mine_pagination_page2_offset_correct`
- `list_mine_other_customers_orders_excluded`
- `list_mine_no_jwt_returns_401_AUTH_MISSING`
- `list_mine_invalid_status_returns_400_VALIDATION_ERROR`

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | Tech-Designer (Opus 4.7 1M, dry-run #2) | Created |
