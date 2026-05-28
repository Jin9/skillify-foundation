# POST /api/v1/checkout/checkout/preview

## Summary

Customer-facing read-only endpoint that returns the **server-priced** checkout breakdown for the caller's current cart against a chosen owned address. Implements §11.2 (BA-PR-002) shipping-fee tier from the SINGLE source `pricing.go` (subtotal < 1500 THB → 60; subtotal ≥ 1500 THB → 0; boundary inclusive at 1500). NEVER mutates state — does not reserve stock, does not write `idempotency_keys`, does not call `identity.profile.read` (buyerEmail is only sourced at commit-time). Closes the CHK-005 "server-priced everything" rule on the read path: the client renders verbatim what the server returns.

## Story refs

- `STORY_CHECKOUT_PREVIEW`
- `STORY_CHECKOUT_SHIPPING_FEE` (boundary cases live here on the read path; commit re-applies the rule on the write path)

## Contract ref

[`checkout.preview`](../../architecture/contracts.json#checkout.preview)

## Auth

- Tier: **customer JWT (Bearer)**.
- Header: `Authorization: Bearer <access JWT>`. Validated via `common/middleware/jwt_middleware`; `claims.role` MUST be `CUSTOMER`.
- `AUTH_MISSING` (401) if header absent; `AUTH_FORBIDDEN` (403) if `role != CUSTOMER` (e.g., admin token must not reach this route).
- Forwarded customer JWT is propagated to `cart.read` (step 1) and `identity.address.list` (step 2). `inventory.stock.bulk-read` (step 4) is hit with `X-Internal-Secret` only — the customer JWT is NOT forwarded into Inventory.

## Request

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "required": ["shippingAddressId"],
  "properties": {
    "shippingAddressId": {
      "type": "string",
      "format": "uuid",
      "minLength": 1,
      "maxLength": 64,
      "description": "UUID v7; MUST belong to the caller (caller-scoped lookup via identity.address.list)."
    },
    "couponCode": {
      "type": "string",
      "description": "REJECTED with VALIDATION_ERROR; coupon engine deferred (CHK-006 / CHK-AMBIG-002). Field presence — even empty string — is the trigger."
    }
  }
}
```

Validation rules beyond the schema:
- If `couponCode` field is PRESENT in the request body (even empty string) → 400 `VALIDATION_ERROR` with message `"couponCode is not supported in MVP"` (CHK-AMBIG-002).
- If any field other than `{shippingAddressId, couponCode}` is present and non-null → 400 `VALIDATION_ERROR` (strict body shape; the schema's `additionalProperties: false` enforces this at bind time).

Required headers:
- `Authorization: Bearer <access JWT>`.
- `Content-Type: application/json`.
- `Idempotency-Key`: OPTIONAL; ignored on preview (preview is read-only; no row is written).

## Response (success)

HTTP 200, envelope:

```json
{
  "code": "SUCCESS",
  "message": "preview computed",
  "data": {
    "items": [
      {
        "productId": "01935b9c-1a17-7c3d-9e8f-aaaaaaaaaaaa",
        "sku": "SKU-A",
        "qty": 2,
        "currentPrice": 599,
        "lineSubtotal": 1198,
        "availableQty": 12,
        "checkoutable": true
      },
      {
        "productId": "01935b9c-1a17-7c3d-9e8f-bbbbbbbbbbbb",
        "sku": "SKU-B",
        "qty": 1,
        "currentPrice": 301,
        "lineSubtotal": 301,
        "availableQty": 0,
        "checkoutable": false
      }
    ],
    "subtotal": 1198,
    "shippingFee": 60,
    "couponDiscount": 0,
    "total": 1258,
    "blockers": ["INSUFFICIENT_STOCK:SKU-B"],
    "addressSnapshot": {
      "addressId": "01935b9c-2a17-7c3d-9e8f-cccccccccccc",
      "receiverName": "Alice Example",
      "phone": "0812345678",
      "addressLine": "12/3 Sukhumvit Soi 24",
      "province": "Bangkok",
      "district": "Khlong Toei",
      "postalCode": "10110"
    }
  },
  "traceId": "0af7651916cd43dd8448eb211c80319c"
}
```

Notes on the success body:
- `items[].lineSubtotal` and `subtotal` SUM only over `checkoutable=true` lines (§11.2 + cross-cutting.pricing.formula.subtotal). Non-checkoutable lines are returned with `checkoutable=false` for client UX (so the cart screen can flag them) but excluded from the priced totals.
- `shippingFee` is computed by `pricing.ComputeShippingFee(subtotalMinor)` — the SINGLE source. Boundary at exactly 1500 THB → 0; subtotal=1499 → 60 (§11.2 / BA-PR-002).
- `total = subtotal + shippingFee - couponDiscount`. `couponDiscount` is always 0 in MVP.
- `addressSnapshot` is the chosen address joined from `identity.address.list` for this preview only (not persisted; commit re-snapshots).

## Response (errors)

| Code | HTTP | Trigger |
|---|--:|---|
| `VALIDATION_ERROR` | 400 | bad body, empty `shippingAddressId`, `couponCode` field present (CHK-AMBIG-002), or extra unknown fields |
| `AUTH_MISSING` | 401 | no `Authorization` header on this protected route |
| `AUTH_FORBIDDEN` | 403 | `claims.role != CUSTOMER` (admin token; service-account token; etc.) |
| `ADDRESS_NOT_OWNED` | 403 | `shippingAddressId` not in caller's `identity.address.list` (CHK-002) |
| `ADDRESS_INCOMPLETE` | 400 | address row exists but is missing required fields (defensive; identity.address.create rejects this earlier) |
| `DATABASE_UNAVAILABLE` | 503 | propagated from cart/identity/catalog/inventory's DB |
| `UPSTREAM_TIMEOUT` | 504 | a downstream call exceeded its budget after 1 retry |

Notes on the empty-cart edge case:
- `STORY_CHECKOUT_PREVIEW.edge: empty cart`. Per the existing handler (`handler_preview.go`), preview returns 200 with `items=[]`, `subtotal=0`, `shippingFee=0` (since `0 < 1500` would yield 60, the empty-items branch shortcuts before pricing applies — see implementation), `total=0`, `blockers=[]`. CHK-001's hard-block on empty cart is enforced ONLY at commit (preview is advisory-only per §13.4 of the requirement). Reviewer-L1 should test that this empty-cart 200 response is not surprising to the frontend — it is the documented behaviour.

## Business logic steps

1. **Auth + role gate:** `common/middleware/jwt_middleware` validates the Bearer JWT (ES256, iss/aud/exp); enforce `claims.role == CUSTOMER`.
2. **Body validation:** bind `rawBody` (strict; reject unknown fields). If `couponCode` field is present → 400 `VALIDATION_ERROR` ("couponCode is not supported in MVP"). Require `shippingAddressId` (max 64 chars).
3. **Step 1 — `cart.read`:** POST `cart:8084/api/v1/cart/cart/read` with the forwarded Bearer JWT. Timeout 800ms; 1 retry on 5xx. Returns the caller's cart items (`{cartItemId, productId, sku, qty}`-shaped). If empty, short-circuit and return the empty-cart success body (see "Notes on the empty-cart edge case" above).
4. **Step 2 — `identity.address.list`:** POST `identity:8081/api/v1/identity/address/list` with the forwarded Bearer JWT. Timeout 500ms; 1 retry. Find the entry matching `req.shippingAddressId`; absent → 403 `ADDRESS_NOT_OWNED` (CHK-002). Incomplete → 400 `ADDRESS_INCOMPLETE`.
5. **Step 3 — `catalog.product.detail` (parallel fan-out, cap 50):** for each cart item, POST `catalog:8082/api/v1/catalog/product/detail` via `errgroup`. Timeout 500ms per call; 1 retry per call. Collect per-line `{ProductID, Status, PriceMinor, Name, ImageURL}`. Lines with `status IN {INACTIVE, DELETED}` or `error == NOT_FOUND` are flagged non-checkoutable and a `PRODUCT_INACTIVE:<productId>` / `PRODUCT_DELETED:<productId>` blocker is added.
6. **Step 4 — `inventory.stock.bulk-read`:** POST `inventory:8083/api/v1/inventory/stock/bulk-read` with `X-Internal-Secret`. Timeout 800ms; 1 retry. Map `{sku → availableQty}`. Lines where `qty > availableQty` are flagged non-checkoutable; a `INSUFFICIENT_STOCK:<sku>` blocker is added.
7. **Server-side pricing (using `pricing.go` — SINGLE source):** for each line where `checkoutable=true`, compute `lineSubtotalMinor = priceSnapshot * qty` (int64 minor units); `subtotalMinor = sum(lineSubtotalMinor)`. Then `shippingFeeMinor = ComputeShippingFee(subtotalMinor)` (60 if `<1500`, else 0; boundary inclusive at 1500 — BA-PR-002). `couponDiscountMinor = 0`. `grandTotalMinor = ComputeGrandTotal(subtotalMinor, shippingFeeMinor, 0)`.
8. **Convert to whole-THB numbers** via `minorToTHB` and assemble the response: `{items, subtotal, shippingFee, couponDiscount, total, blockers, addressSnapshot}`. Wrap in `wrapper.Respond` with `code=SUCCESS`, `message='preview computed'`, traceId from the request context.

## Side effects

- **None.** Read-only by construction (per `checkout.preview.idempotency_rules: "Read; no commit."`).
- No `idempotency_keys` row is created (`Idempotency-Key` is ignored).
- No `saga_log` row is written (saga_log is exclusively a `commit` artifact).
- No outbox event.
- No cache invalidation.

## Idempotency

None required at the API level — this endpoint is naturally idempotent (read-only). Repeated calls with the same body produce the same response within the cart/catalog/inventory drift window. The `Idempotency-Key` header, if supplied by the client, is silently ignored — the body alone determines the response.

## Performance

- p95 < 250 ms (cross-component PERF-002 gives checkout-as-a-whole < 500 ms; preview is the cheaper of the two endpoints).
- Expected QPS at MVP: 2–10.
- Latency budget breakdown (success path, no retries):
  - Step 1 cart.read ≤ 800ms (typical 30–80ms)
  - Step 2 identity.address.list ≤ 500ms (typical 20–50ms)
  - Step 3 catalog.product.detail × N parallel ≤ 500ms wall-clock (errgroup; typical 50–100ms wall-clock for N≤10)
  - Step 4 inventory.stock.bulk-read ≤ 800ms (typical 30–80ms)
  - Local pricing arithmetic < 1 ms.
  - Total typical: 130–310 ms p95.

## Sequence diagram

(Optional per the api-spec.md template — preview is non-orchestrating and the steps are linear. Provided here for reviewer convenience.)

```mermaid
sequenceDiagram
  autonumber
  participant FE as Frontend (Next.js route handler)
  participant CO as Checkout
  participant CART as Cart
  participant ID as Identity
  participant CAT as Catalog
  participant INV as Inventory
  participant PG as Postgres (pricing.go local)

  FE->>+CO: POST /api/v1/checkout/checkout/preview<br/>Bearer + {shippingAddressId}
  CO->>CO: jwt_middleware + role==CUSTOMER<br/>reject couponCode/extra fields
  CO->>+CART: POST /cart/cart/read (Bearer)
  CART-->>-CO: items[]
  alt items == []
    CO-->>FE: 200 {items:[], total:0, blockers:[]}<br/>(advisory; CHK-001 enforced at commit)
  else items > 0
    CO->>+ID: POST /identity/address/list (Bearer)
    ID-->>-CO: addresses[]
    alt shippingAddressId not in addresses
      CO-->>FE: 403 ADDRESS_NOT_OWNED (CHK-002)
    else found
      par fan-out per productId (errgroup; cap 50)
        CO->>+CAT: POST /catalog/product/detail (productId)
        CAT-->>-CO: {status, priceMinor, name, imageUrl}
      end
      CO->>+INV: POST /inventory/stock/bulk-read<br/>X-Internal-Secret + {skus}
      INV-->>-CO: stocks[{sku, availableQty}]
      CO->>PG: pricing.ComputeShippingFee(subtotalMinor)<br/>(SINGLE SOURCE; 1500 inclusive — BA-PR-002)
      CO->>PG: pricing.ComputeGrandTotal(...)
      CO-->>-FE: 200 {items, subtotal, shippingFee, total, blockers, addressSnapshot}
    end
  end
```

## Test cases

- `preview_with_authenticated_customer_and_valid_address_returns_priced_breakdown`
- `preview_subtotal_1499_THB_returns_shippingFee_60_total_1559_BA_PR_002`
- `preview_subtotal_1500_THB_returns_shippingFee_0_total_1500_BA_PR_002_inclusive_boundary`
- `preview_subtotal_9999_THB_returns_shippingFee_0_total_9999_no_progressive_discount`
- `preview_address_not_owned_returns_403_ADDRESS_NOT_OWNED`
- `preview_couponCode_field_present_returns_400_VALIDATION_ERROR`
- `preview_extra_unknown_field_returns_400_VALIDATION_ERROR`
- `preview_inactive_product_in_cart_marks_line_non_checkoutable_and_adds_blocker`
- `preview_deleted_product_in_cart_marks_line_non_checkoutable_with_PRODUCT_DELETED_blocker`
- `preview_qty_greater_than_availableQty_marks_line_non_checkoutable_with_INSUFFICIENT_STOCK_blocker`
- `preview_empty_cart_returns_200_with_items_empty_advisory_only`
- `preview_admin_token_returns_403_AUTH_FORBIDDEN`
- `preview_no_authorization_header_returns_401_AUTH_MISSING`
- `preview_idempotency_key_header_is_ignored_on_repeated_call`
- `preview_does_NOT_create_idempotency_keys_row` (negative side-effect smoke)
- `preview_does_NOT_create_inventory_reservation` (negative side-effect smoke)
- `preview_does_NOT_call_identity_profile_read` (negative call-graph smoke; profile.read is commit-only per REV-L2-001)
- `preview_inventory_503_returns_503_DATABASE_UNAVAILABLE`
- `preview_catalog_timeout_marks_affected_lines_non_checkoutable_but_returns_200`
- `preview_uses_pricing_go_single_source` (grep-style: shippingFee computation MUST go through `pricing.ComputeShippingFee`)

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | Tech-Designer (Opus 4.7 1M, dry-run #2) | Created. Pins §11.2 single-source via `pricing.go`; documents that `identity.profile.read` is NOT called on preview (commit-only per REV-L2-001); strict body shape with explicit `couponCode` reject. |
