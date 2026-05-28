# POST /api/v1/order/order/list-admin

## Summary

Admin-only paged listing of all orders with optional filters: `status`, `dateFrom`/`dateTo` (createdAt), and `search` (matches `order_number` prefix OR `buyer_email_snapshot` LOWER-exact). Backed by indexes on `(status, created_at DESC)`, `order_number` UNIQUE, and `LOWER(buyer_email_snapshot)`. Writes one `admin_action_log` row per call (audit-deferred fallback per the brief).

## Story refs

- `STORY_ORDER_LIST_DETAIL` (admin path is implicit in the EPIC; admin endpoint covers ADM-003/ADM-004 — backend-only, no admin UI per BA out_of_scope)

## Contract ref

[`order.list-admin`](../../architecture/contracts.json#order.list-admin)

## Auth

- Tier: **admin_jwt**.
- `transport/http` admin-role decorator checks `claims.role = ADMIN`; CUSTOMER tokens reject 403 `AUTH_FORBIDDEN`.

## Request

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "page":     { "type": "integer", "minimum": 1, "default": 1 },
    "limit":    { "type": "integer", "minimum": 1, "maximum": 100, "default": 20 },
    "status":   {
      "type": "string",
      "enum": ["PENDING_PAYMENT","PAID","PAYMENT_FAILED","PAYMENT_EXPIRED","PACKING","SHIPPED","DELIVERED","CANCELLED"]
    },
    "dateFrom": { "type": "string", "format": "date" },
    "dateTo":   { "type": "string", "format": "date" },
    "search":   { "type": "string", "maxLength": 254, "description": "matches order_number prefix OR buyer_email_snapshot LOWER-exact" }
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
        "orderId": "01935b9c-...",
        "orderNumber": "ORD-20260508-000123",
        "userId": "0192d-customer",
        "buyerEmail": "alice@example.com",
        "status": "PAID",
        "subtotal": 980,
        "shippingFee": 60,
        "total": 1040,
        "itemCount": 3,
        "createdAt": "2026-05-08T07:12:33Z"
      }
    ],
    "total": 421
  },
  "traceId": "0af7651916cd43dd8448eb211c80319c"
}
```

## Response (errors)

| Code | HTTP | Trigger |
|---|--:|---|
| `VALIDATION_ERROR` | 400 | `dateFrom > dateTo`; `status` not in enum; `page`/`limit` out of bounds |
| `AUTH_MISSING` | 401 | no `Authorization` header |
| `AUTH_INVALID` | 401 | JWT signature/exp/iss/aud invalid |
| `AUTH_FORBIDDEN` | 403 | `claims.role != ADMIN` (AUTH-003) |

## Business logic steps

1. `transport/http` admin-decorator: reject 403 `AUTH_FORBIDDEN` if `claims.role != ADMIN`.
2. Bind + validate body. If `dateFrom` and `dateTo` both present: `dateFrom <= dateTo` else 400 `VALIDATION_ERROR`.
3. If `search` is non-empty:
   - If it begins with `ORD-` → treat as `order_number` prefix match (backed by `idx_orders_order_number`).
   - Else → treat as a customer email; lowercase it and query `WHERE LOWER(buyer_email_snapshot) = LOWER($1)` (backed by `idx_orders_buyer_email_lower`).
4. `BEGIN;` (small audit tx)
5. `repo.OrdersListAdmin(filters, page, limit)` — composes the filter + ORDER BY `created_at DESC` + LIMIT/OFFSET.
6. `repo.OrdersCountAdmin(filters)` — for `total`.
7. `repo.AdminActionLogInsert(tx, {actor_user_id=claims.sub, endpoint='order.list-admin', order_id=NULL, request_summary={filters: ...}, outcome_code='SUCCESS'});`
8. `COMMIT;`
9. Map results to `OrderSummary[]` with `buyerEmail` included (admin convenience).
10. Wrap with `common/wrapper.Success`.

## Side effects

- `INSERT admin_action_log` — one row per call, even on rejection paths (the audit log records the rejection with `outcome_code='AUTH_FORBIDDEN'` etc.).

## Idempotency

- None — operation is naturally idempotent (read).
- Audit row is appended on every call; that's not "idempotent" semantically but is ledger-safe (every read is recorded).

## Performance

- p95 < 500ms with reasonable filter cardinality.
- Worst-case full-status scan with no date filter is bounded by LIMIT; the `(status, created_at DESC)` index covers the typical admin operator flow ("show me everything PAID today").

## Test cases

- `list_admin_returns_all_orders_when_no_filters`
- `list_admin_filters_by_status`
- `list_admin_filters_by_date_range`
- `list_admin_inverted_date_range_returns_400_VALIDATION_ERROR`
- `list_admin_search_by_order_number_prefix`
- `list_admin_search_by_email_uses_buyer_email_snapshot`
- `list_admin_search_by_email_is_case_insensitive`
- `list_admin_writes_audit_row`
- `list_admin_customer_token_returns_403_AUTH_FORBIDDEN`
- `list_admin_no_jwt_returns_401_AUTH_MISSING`

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | Tech-Designer (Opus 4.7 1M, dry-run #2) | Created |
