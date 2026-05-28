# POST /api/v1/catalog/product/create

## Summary

Admin-only endpoint that creates a new product in the catalog. The product INSERT, its initial images, the initial status history row, and an `outbox_events` row for `events.product.created` are all written in a single database transaction — guaranteeing that if the product is persisted, the event will be published to Kafka (at-least-once via the outbox relay goroutine). Inventory consumes `product.created` and seeds `stock_levels(sku, 0, 0, 0)`. Implements requirement CAT-007, CAT-010.

## Story refs

- `STORY_CATALOG_ADMIN_PRODUCT_CRUD`

## Contract ref

[`catalog.product.create`](../../architecture/contracts.json#catalog.product.create)

## Auth

- Tier: **required — ADMIN role**.
- `Authorization: Bearer <jwt>` header required.
- `common/middleware` JWT middleware validates ES256 signature, `iss`, `aud`, `exp`.
- If `claims.role != 'ADMIN'`: return `AUTH_FORBIDDEN (403)`.
- If header absent: return `AUTH_MISSING (401)`.

## Request

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "required": ["sku", "name", "price", "categoryId"],
  "properties": {
    "sku": {
      "type": "string",
      "minLength": 1,
      "maxLength": 100
    },
    "name": {
      "type": "string",
      "minLength": 1,
      "maxLength": 255
    },
    "description": {
      "type": "string",
      "default": ""
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
    "categoryId": {
      "type": "string",
      "format": "uuid"
    },
    "status": {
      "type": "string",
      "enum": ["DRAFT", "ACTIVE"],
      "default": "DRAFT"
    }
  }
}
```

**Notes:**
- `sku` is case-insensitively unique; stored as-is but uniqueness checked via `LOWER(sku)`.
- `price` must be > 0 and may have up to 2 decimal places; stored as `NUMERIC(12,2)`.
- `status` defaults to `DRAFT` if omitted. `INACTIVE` and `DELETED` are forbidden on create.
- `images` is optional; if omitted the product is created with no images.

## Response (success)

HTTP 201, envelope:

```json
{
  "code": "CREATED",
  "message": "Product created",
  "data": {
    "productId": "01935b9c-0000-7000-0000-000000000001"
  },
  "traceId": "0af7651916cd43dd8448eb211c80319c"
}
```

## Response (errors)

| Code | HTTP | Trigger |
|---|--:|---|
| `AUTH_MISSING` | 401 | No `Authorization` header |
| `AUTH_FORBIDDEN` | 403 | Token valid but `claims.role = CUSTOMER` |
| `VALIDATION_ERROR` | 400 | Missing required fields; `price <= 0`; invalid `status` enum (`INACTIVE`/`DELETED`); malformed UUID in `categoryId`; `images` array > 20; any image URL empty or > 2048 chars; `sku` empty or > 100 chars |
| `DUPLICATE_SKU` | 409 | `LOWER(sku)` already exists in `products` table (DB unique-constraint violation) |
| `NOT_FOUND` | 404 | `categoryId` does not exist in `categories` OR category exists but `active = false` |

## Business logic steps

1. JWT middleware validates token and injects `claims` into Gin context. Returns `AUTH_MISSING` or `AUTH_FORBIDDEN` before handler runs.
2. Bind and validate request body (`common/validator`). Return `VALIDATION_ERROR` on schema violations.
3. Verify `categoryId` exists and is active: `SELECT id FROM categories WHERE id = $1 AND active = TRUE`. If not found: return `NOT_FOUND`.
4. Generate `productId = uuid_v7()` in Go (UUID v7 — needed in outbox payload in the same tx before COMMIT).
5. Begin transaction (`pgx`):
   a. `INSERT INTO products (id, sku, name, description, price, category_id, status, visible, created_at, updated_at)`.
   b. If `images` non-empty: batch `INSERT INTO product_images (id, product_id, url, sort_order)` — one row per URL with `sort_order = array index`.
   c. `INSERT INTO product_status_history (id, product_id, from_status=NULL, to_status=<status>, actor_user_id=claims.sub, created_at)` — `from_status=NULL` marks initial creation.
   d. Generate `outboxId = uuid_v7()`. `INSERT INTO outbox_events (id=outboxId, aggregate_id=productId, event_type='product.created', payload_json={eventId: uuid_v7(), occurredAt: RFC3339, productId, sku}, created_at, published_at=NULL)`.
   e. `COMMIT`.
6. If DB returns unique-constraint violation on `LOWER(sku)` unique index: return `DUPLICATE_SKU (409)` — no outbox row written (tx rolled back).
7. Return `CREATED` envelope with `{productId}`.

**Outbox relay (async, started at boot):** `StartOutboxRelay(ctx, db, kafkaProducer)` — goroutine polls `outbox_events WHERE published_at IS NULL ORDER BY id LIMIT 100 FOR UPDATE SKIP LOCKED`, publishes each to `ecom.catalog.events` (partition key = `payload_json->>'sku'`), then `UPDATE SET published_at = NOW()`. Inventory consumer seeds `stock_levels(sku, 0, 0, 0)` on receipt; dedup on `eventId`.

## Side effects

- INSERT one row into `catalog.products`.
- INSERT 0..20 rows into `catalog.product_images`.
- INSERT one row into `catalog.product_status_history` (`from_status=NULL`).
- INSERT one row into `catalog.outbox_events` (`event_type='product.created'`, `published_at=NULL`).
- All four INSERTs are in a single transaction — atomic; no partial writes.
- Async: relay goroutine publishes `product.created` to `ecom.catalog.events`. Inventory seeds `stock_levels`.

## Idempotency

None at the API level — product create is not idempotent. Retrying with the same SKU returns `DUPLICATE_SKU (409)`. No `Idempotency-Key` header required.

## Performance

- **p95 target:** < 300ms (admin write; no storefront SLA).
- **Expected QPS at MVP:** < 5 (admin-only; infrequent catalog updates).
- Single-transaction write; no external HTTP calls during create.

## Sequence diagram

```mermaid
sequenceDiagram
    participant Admin as Admin Client
    participant CAT as Catalog Service
    participant DB as Postgres (catalog schema)
    participant RELAY as Outbox Relay (goroutine)
    participant KAFKA as Kafka (ecom.catalog.events)
    participant INV as Inventory Service

    Admin->>+CAT: POST /api/v1/catalog/product/create {sku, name, price, categoryId, ...}
    CAT->>CAT: JWT middleware (ADMIN check)
    CAT->>DB: SELECT categories WHERE id=$1 AND active=TRUE
    DB-->>CAT: category row
    CAT->>DB: BEGIN
    CAT->>DB: INSERT products
    CAT->>DB: INSERT product_images (batch)
    CAT->>DB: INSERT product_status_history (from_status=NULL)
    CAT->>DB: INSERT outbox_events (product.created, published_at=NULL)
    CAT->>DB: COMMIT
    CAT-->>-Admin: 201 CREATED {productId}

    Note over RELAY: Async — runs independently on 5s tick
    RELAY->>DB: SELECT outbox_events WHERE published_at IS NULL ... FOR UPDATE SKIP LOCKED
    DB-->>RELAY: pending rows
    RELAY->>KAFKA: Produce (topic=ecom.catalog.events, key=sku, payload)
    RELAY->>DB: UPDATE outbox_events SET published_at=NOW()

    KAFKA-->>INV: product.created event
    INV->>INV: dedup on eventId; seed stock_levels(sku, 0, 0, 0)
```

## Test cases

- `product_create_ADMIN_with_valid_input_returns_201_productId`
- `product_create_default_status_is_DRAFT`
- `product_create_status_ACTIVE_product_appears_in_listing`
- `product_create_CUSTOMER_token_returns_403_AUTH_FORBIDDEN`
- `product_create_no_token_returns_401_AUTH_MISSING`
- `product_create_duplicate_SKU_case_insensitive_returns_409_DUPLICATE_SKU`
- `product_create_inactive_category_returns_404`
- `product_create_absent_category_returns_404`
- `product_create_price_zero_returns_400_VALIDATION_ERROR`
- `product_create_price_negative_returns_400_VALIDATION_ERROR`
- `product_create_images_over_20_returns_400_VALIDATION_ERROR`
- `product_create_status_INACTIVE_returns_400_VALIDATION_ERROR`
- `product_create_status_DELETED_returns_400_VALIDATION_ERROR`
- `product_create_outbox_row_written_with_published_at_NULL`
- `product_create_relay_publishes_event_and_marks_published_at`
- `product_create_all_four_inserts_in_single_tx_rollback_on_error`

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | Tech-Designer (Claude Sonnet 4.6, dispatch 2) | Created |
