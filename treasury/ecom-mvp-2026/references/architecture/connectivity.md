# Connectivity — ShopPilot MVP (B2C E-Commerce Platform)

End-to-end request and event flows. Operator's tracing aid. Companion to [`./infra-summary.md`](./infra-summary.md) and [`./infra-topology.md`](./infra-topology.md). The contract surface is in [`./contracts.json`](./contracts.json).

There is no API gateway in MVP — Next.js route handlers play that role.

---

## Synchronous request flows

### Customer registration / login

```
Browser
   |  POST /api/proxy/identity/user/register   (or /login)
   v
Next.js route handler  (no withAuth — public route)
   |  POST identity:8081/api/v1/identity/user/register   (or /login)
   v
Identity service
   - register: bcrypt(password); INSERT users; respond {userId, email, name, role=CUSTOMER}
   - login:    bcrypt.CompareHashAndPassword; mint access JWT (ES256, 15min) + refresh JWT (ES256, 14d)
   |  Response envelope {code, message, data, traceId}
   v
Next.js route handler
   - on /login success: Set-Cookie access (HttpOnly, Secure, SameSite=Lax, Max-Age=15min) + refresh (HttpOnly, ..., Max-Age=14d)
   v
Browser  (cookies set; SPA reads only the non-sensitive auth-hint cookie for UI gating)
```

### Customer reads profile (and any other Bearer-JWT route)

```
Browser
   |  POST /api/proxy/identity/profile/read
   v
Next.js route handler -> withAuth (6-step flow per cross-cutting.auth.token_acquisition_flow):
   1. read 'access' cookie
   2. if present and exp > now+30s skew  -> Authorization: Bearer <access>; forward
   3. else read 'refresh' cookie; if absent -> 401 AUTH_MISSING
   4. POST identity:8081/api/v1/identity/user/refresh {refreshToken}
   5. on 200: Set-Cookie BOTH cookies with rotated values (single-use rotation; old jti revoked)
   6. on AUTH_INVALID/AUTH_REVOKED: clear both cookies (Max-Age=0); 401 AUTH_REVOKED -> SPA redirects to /login
   |  POST identity:8081/api/v1/identity/profile/read   (Authorization: Bearer ...)
   v
Identity service -> SELECT user; respond {userId, email, name, phone, defaultAddressId, ...}
```

### Customer browse (public — no auth)

```
Browser -> Next.js page (Server Component) reads via /api/proxy/catalog/product/list
       -> POST catalog:8082/api/v1/catalog/product/list   (no Authorization header)
       -> Catalog SELECT products WHERE status='ACTIVE'; if filter.inStock=true, sync call to inventory:8083/api/v1/inventory/stock/bulk-read (X-Internal-Secret on this internal hop)
       -> respond {items, total, page, limit}
```

### Cart mutations (customer JWT)

```
Browser -> Next.js route handler -> withAuth -> Bearer
   POST cart:8084/api/v1/cart/cart/add-item  {productId, qty}
   - cart server validates productId via catalog.product.detail (PRODUCT_INACTIVE/PRODUCT_DELETED reject)
   - cart_items UPSERT (cart_id, sku) MERGE quantity per CART-006
```

### Checkout commit (the orchestrated 10-step flow)

```
Browser -> Next.js route handler -> withAuth + generates Idempotency-Key (UUID v4; persisted in localStorage per pending checkout)
   POST checkout:8085/api/v1/checkout/checkout/commit
        Headers: Authorization: Bearer <jwt>; Idempotency-Key: <client UUID>
        Body:    {shippingAddressId}

Step 1: Checkout -> cart:8084/api/v1/cart/cart/read                                                  (Bearer; customer-context)
Step 2: Checkout -> identity:8081/api/v1/identity/profile/read                                        (Bearer; pulls buyerEmail; REV-L2-001 fix)
Step 3: Checkout -> identity:8081/api/v1/identity/address/list                                        (Bearer; validates shippingAddressId ownership)
Step 4: Checkout -> catalog:8082/api/v1/catalog/product/detail (per item, parallel via errgroup)      (no auth needed for read; fan-in)
        Checkout -> inventory:8083/api/v1/inventory/stock/bulk-read                                   (X-Internal-Secret)
Step 5: Checkout BEGIN tx; INSERT idempotency_keys (status=INFLIGHT)
Step 6: Checkout -> inventory:8083/api/v1/inventory/reservation/create                                (X-Internal-Secret; Idempotency-Key=newOrderId)
Step 7: Checkout -> order:8086/api/v1/order/internal/create-from-checkout                             (X-Internal-Secret; Idempotency-Key=newOrderId; includes buyerEmail)
Step 8: Checkout -> payment:8087/api/v1/payment/intent/create                                         (X-Internal-Secret; Idempotency-Key=newOrderId)
Step 9: Checkout -> cart:8084/api/v1/cart/cart/clear-on-checkout                                      (best-effort; Bearer)
Step 10: Checkout UPDATE idempotency_keys SET response_envelope = ...; COMMIT

Returns to browser: {orderId, orderNumber, status: PENDING_PAYMENT, subtotal, shippingFee, total, paymentIntent: {paymentIntentId, amount, status: REQUIRES_PAYMENT}}
```

### Mock payment simulate (customer-driven)

```
Browser (payment route) -> withAuth -> POST /api/proxy/payment/intent/simulate {paymentIntentId, outcome: SUCCESS|FAILED|TIMEOUT}
                       -> POST payment:8087/api/v1/payment/intent/simulate (Bearer; ownership check: claims.sub == intent.customerUserId)
                       -> Payment internally enqueues self-callback to its own /api/v1/payment/intent/callback (sync; same-process)
                       -> Callback handler: BEGIN tx; INSERT payment_callback_dedup (paymentIntentId, providerStatus); UPDATE payment_records (mockPaymentRef, providerStatus, paidAt); INSERT outbox_events (payment.completed | payment.failed | payment.expired); COMMIT
                       -> outbox-relay (every 5s) drains to Kafka topic ecom.payment.events
```

### Admin endpoints (CUSTOMER-token-rejected)

```
Admin (curl/Postman, no UI) -> POST identity:8081/api/v1/identity/user/login (admin credential pre-seeded)
                            -> Bearer JWT with role=ADMIN
                            -> POST catalog:8082/api/v1/catalog/product/create     (Bearer; role=ADMIN check; CUSTOMER -> 403 AUTH_FORBIDDEN)
                            -> POST inventory:8083/api/v1/inventory/stock/adjust   (same)
                            -> POST order:8086/api/v1/order/order/update-status-admin (transitions table per ORD-003 / §9.2)
```

## Asynchronous event flows

### `ecom.payment.events` (per-orderId partitioned)

| Event | Producer | Consumers | Consumer effect (state-driven via SELECT FOR UPDATE) |
|---|---|---|---|
| `events.payment.completed` | payment (outbox) | order, inventory | order: PENDING_PAYMENT -> PAID; inventory: RESERVED -> COMMITTED (or no-op if RELEASED — PR-003 race) |
| `events.payment.failed` | payment (outbox) | order, inventory | order: PENDING_PAYMENT -> PAYMENT_FAILED; inventory: RESERVED -> RELEASED (or no-op) |
| `events.payment.expired` | payment (outbox) | order, inventory | order: PENDING_PAYMENT -> PAYMENT_EXPIRED (no-op if already moved by reservation.expired path); inventory: RESERVED -> RELEASED |

### `ecom.order.events` (per-orderId partitioned)

| Event | Producer | Consumers | Consumer effect |
|---|---|---|---|
| `events.order.cancelled` | order (outbox) | inventory, payment | inventory: state-driven release per current reservation status (RESERVED -> available, or COMMITTED -> sold->available per BA's PR-006 MVP restock policy); payment: REQUIRES_PAYMENT -> CANCELLED, else no-op |

### `ecom.inventory.events` (per-orderId partitioned)

| Event | Producer | Consumers | Consumer effect |
|---|---|---|---|
| `events.reservation.expired` | inventory.sweeper (outbox) | order | order: PENDING_PAYMENT -> PAYMENT_EXPIRED; status_history actorRole=SYSTEM, reason='reservation expired (TTL)' |

### `ecom.catalog.events` (per-sku partitioned)

| Event | Producer | Consumers | Consumer effect |
|---|---|---|---|
| `events.product.created` | catalog (outbox) | inventory | inventory UPSERT stock_levels(sku) with all-zero quantities so subsequent inStock filters resolve correctly |

## System-internal flows (no HTTP route, no caller)

| Job | Cadence | Owner | Effect |
|---|---|---|---|
| `inventory.reservation.sweep-expired` | every 30s | inventory | SELECT FROM reservations WHERE status='RESERVED' AND expires_at < NOW() FOR UPDATE SKIP LOCKED LIMIT 200; per row: SELECT stock_levels FOR UPDATE; reserved_qty -= qty; available_qty += qty; status='RELEASED' release_reason='EXPIRED'; INSERT outbox events.reservation.expired; COMMIT. Closes the abandoned-checkout stock-leak hole (PR-001 / BA edge case 'browser walks away'). |
| `payment.intent.expiry-sweeper` | every 60s | payment | SELECT payment_intents WHERE status='REQUIRES_PAYMENT' AND expires_at < NOW() FOR UPDATE SKIP LOCKED; transition to EXPIRED; emits events.payment.expired via outbox. |
| `outbox-relay` | every 5s per service | catalog, inventory, order, payment | SELECT FROM outbox_events WHERE published_at IS NULL ORDER BY id LIMIT 100; publish to Kafka via common/kafka; UPDATE published_at = NOW(). At-least-once; consumers dedup on eventId. |
| `checkout.saga-reconciler` | every 60s | checkout | SELECT saga_log WHERE status='DEGRADED' AND created_at < NOW() - INTERVAL '1 hour'; emit metric `checkout_saga_degraded_unreconciled_total`; do NOT auto-remediate (manual ops in MVP). REV-L2-008 followup. |

## Auth boundaries

| Path | Auth tier | Notes |
|---|---|---|
| `/api/v1/identity/user/register|login` | none | public |
| `/api/v1/identity/user/refresh|logout` | (refresh token in body) | rotate jti in tx; revoked-jti store in identity_db |
| `/api/v1/identity/profile/read|update` | customer JWT | `claims.sub == users.id` |
| `/api/v1/identity/address/*` | customer JWT | scoped by claims.sub; address ownership check on update/delete/set-default |
| `/api/v1/catalog/product/list|detail` | none | public; storefront filters out non-ACTIVE |
| `/api/v1/catalog/product/create|update|soft-delete` | admin JWT | role=ADMIN |
| `/api/v1/catalog/category/create|update|deactivate` | admin JWT | role=ADMIN |
| `/api/v1/catalog/category/list` | none | public |
| `/api/v1/inventory/stock/read|bulk-read` | X-Internal-Secret | callers: catalog, cart, checkout, order |
| `/api/v1/inventory/stock/adjust` | admin JWT | role=ADMIN; writes stock_adjustments(actor_user_id) |
| `/api/v1/inventory/reservation/create|commit|release` | X-Internal-Secret | constant-time compare; idempotencyKey == orderId per cross-cutting.idempotency.internal_orderId_rule |
| `/api/v1/cart/*` | customer JWT | scoped by claims.sub; cart row keyed on customer_user_id |
| `/api/v1/checkout/checkout/preview|commit` | customer JWT + Idempotency-Key on commit | per-(endpoint, customerUserId) idempotency in checkout_db |
| `/api/v1/order/order/list-mine|detail|cancel-mine` | customer JWT | server enforces customer_user_id == claims.sub (ORD-001) |
| `/api/v1/order/order/list-admin|update-status-admin` | admin JWT | role=ADMIN |
| `/api/v1/order/internal/create-from-checkout|cancel-on-checkout-failure` | X-Internal-Secret | callers: checkout only |
| `/api/v1/payment/intent/create` | X-Internal-Secret | caller: checkout |
| `/api/v1/payment/intent/simulate` | customer JWT + ownership check | mock-only |
| `/api/v1/payment/intent/callback` | (in-process) | mock self-callback; in production this would be HMAC-signed webhook |

## Failure modes

| Spine flow | Failure | Compensation |
|---|---|---|
| Step 2 (identity.profile.read) fails | no state mutated | return 502 UPSTREAM_TIMEOUT envelope |
| Step 6 (inventory.reservation.create) fails | no state mutated | return error to client (e.g. INSUFFICIENT_STOCK with per-item availableQty/requestedQty) |
| Step 7 (order.create-from-checkout) fails after step 6 OK | inventory has a RESERVED reservation but no Order | Checkout calls inventory.reservation.release reason=ORDER_CANCELLED; on retry-exhaustion log CRITICAL + saga DEGRADED; reservation TTL (15min) is the safety net; saga reconciler surfaces the row at age > 1h |
| Step 8 (payment.intent.create) fails after step 7 OK | order is PENDING_PAYMENT; reservation RESERVED; intent absent | Checkout calls order.cancel-on-checkout-failure reason='payment intent creation failed: <upstream code>' THEN inventory.reservation.release reason=ORDER_CANCELLED; events.order.cancelled flows to Inventory consumer (state-driven; COMMITTED -> sold->available, RESERVED -> available) |
| Step 9 (cart.clear-on-checkout) fails after success | order created and paid; cart still holds the items | log WARN; do NOT roll back order/payment; cart inconsistency is cosmetic (customer can manually remove) |
| `events.payment.completed` consumer crashes mid-tx | tx rolls back; offset NOT committed | Kafka rebalance redelivers; consumed_events PK + state-driven SELECT FOR UPDATE prevents double-apply |
| `events.payment.completed` arrives AFTER `events.order.cancelled` (PR-003 race) | inventory is COMMITTED then order.cancelled lands | Inventory state-driven consumer reads CURRENT reservation status (COMMITTED) and applies sold->available branch — event.fromStatus is audit/debug only; current-row-read wins. Order's events.payment.completed handler reads CURRENT order status (CANCELLED) and no-ops. ADR-009. |
| `payment.callback` arrives twice for same (paymentIntentId, providerStatus) | dedup unique-constraint violation | Handler catches the PG unique violation; returns the cached envelope; NEVER re-emits a domain event |
| Internal call rejected (X-Internal-Secret mismatch) | recipient logs `internal_auth_reject` (no secret in log); returns 401 AUTH_INVALID | Caller treats as upstream error and runs the standard compensation matrix |
| Postgres unavailable on a stateful service | service returns 503 DATABASE_UNAVAILABLE with traceId envelope | reads do not fall back to stale or fabricated data; orchestrator-level retry policy is per-call (3 retries with jitter); checkout's saga DEGRADED captures the case |
