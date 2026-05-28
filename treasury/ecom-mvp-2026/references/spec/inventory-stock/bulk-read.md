# POST /api/v1/inventory/stock/bulk-read

## Summary

Multi-SKU stock read. Returns `availableQty / reservedQty / soldQty` for up to 200 SKUs in a single call, preserving input order; missing SKUs are omitted (not zeroed). This is the read path used by Catalog (`product.list` with `inStock=true` filter), Cart (`subtotal / availableQty` checks), and Checkout preview.

## Story refs

- `STORY_INVENTORY_RESERVE_AND_CONVERT` (visibility into the live state)
- `STORY_CART_VALIDATE_AVAILABILITY` (cart-side caller)
- `STORY_CHECKOUT_PREVIEW` (checkout preview pre-commit availability check)

## Contract ref

[`inventory.stock.bulk-read`](../../architecture/contracts.json#inventory.stock.bulk-read)

## Auth

- Tier: **internal_secret**.
- Header: `X-Internal-Secret: <32-byte secret>`.
- Constant-time compare per [ADR-008](../../architecture/ADRs/ADR-008-internal-shared-secret.md).
- Callers per contract: catalog, cart, checkout (all backend services on the private network plane).

## Request

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "required": ["skus"],
  "properties": {
    "skus": {
      "type": "array",
      "minItems": 1,
      "maxItems": 200,
      "items": { "type": "string", "minLength": 1, "maxLength": 64 }
    }
  }
}
```

## Response (success)

HTTP 200, envelope. Result `stocks[]` is in INPUT ORDER (caller-supplied); missing SKUs are omitted (no zero-row placeholder).

```json
{
  "code": "SUCCESS",
  "message": "ok",
  "data": {
    "stocks": [
      { "sku": "SKU-A", "availableQty": 8, "reservedQty": 2, "soldQty": 0 },
      { "sku": "SKU-B", "availableQty": 0, "reservedQty": 0, "soldQty": 5 }
    ]
  },
  "traceId": "0af7651916cd43dd8448eb211c80319c"
}
```

## Response (errors)

| Code | HTTP | Trigger |
|---|--:|---|
| `VALIDATION_ERROR` | 400 | `skus` empty, `skus` length > 200, or any element invalid (empty / > 64 chars) |
| `AUTH_INVALID` | 401 | `X-Internal-Secret` missing or mismatched |
| `DATABASE_UNAVAILABLE` | 503 | Postgres unreachable |

Note: the contract intentionally does NOT include a `NOT_FOUND` per-SKU error — missing SKUs are simply omitted from the response. Callers compute "missing" by comparing input `skus` to `data.stocks[].sku`.

## Business logic steps

1. `common/middleware.InternalAuth(secret)` validates `X-Internal-Secret`.
2. Bind + validate via `common/validator`: `len(skus) ∈ [1, 200]`; per-element non-empty, ≤ 64 chars.
3. Single SQL: `SELECT sku, available_qty, reserved_qty, sold_qty FROM stock_levels WHERE sku = ANY($1)` — Postgres scans efficiently with the PK index.
4. Build a map `sku → row`; iterate the input `skus` slice in order; emit only the present SKUs into `stocks[]`.
5. `wrapper.Respond(c, SUCCESS, {stocks})`.

## Side effects

None — read-only.

## Idempotency

Naturally idempotent — read.

## Performance

- p95 < 50ms for 200 SKUs (single `WHERE sku = ANY($1)` query) — comfortably within the checkout preview p95 < 500ms budget (BA non-functional latency PERF-002 indirect).
- Expected QPS at MVP: 50–200 (catalog list + cart open + checkout preview all hit this).

## Test cases

- `stock_bulk_read_preserves_input_order`
- `stock_bulk_read_drops_missing_skus_no_placeholder`
- `stock_bulk_read_with_201_skus_returns_VALIDATION_ERROR`
- `stock_bulk_read_with_empty_array_returns_VALIDATION_ERROR`
- `stock_bulk_read_returns_AUTH_INVALID_when_secret_missing`
- `stock_bulk_read_p95_under_50ms_for_200_skus_smoke`

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | Tech-Designer (Opus 4.7 1M) | Created |
