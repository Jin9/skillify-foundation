# Observability — ShopPilot MVP (B2C E-Commerce Platform)

The observability contract every service must satisfy. Authored by Tech-Lead alongside the architecture; validated by Reviewer-L1 (per service) and Reviewer-L2 (cross-component).

Local-stack docker-compose snippet shipped at `B2C E-Commerce Platform/observability/docker-compose.yml` referenced from this spec.

---

## Structured logging

All services use `slog` (Go `log/slog`) with the JSON handler. Every log event MUST include the following fields. Field plumbing comes from `common/middleware` (`AccessLogMiddleware`, `RequestIDMiddleware`, `TraceContextTraceIDMiddleware`):

| Field | Source | Example |
|---|---|---|
| `time` | slog default | `2026-05-08T17:48:01.234Z` (RFC 3339) |
| `level` | slog default | `INFO` / `WARN` / `ERROR` |
| `msg` | call site | `"login_attempted"` (kebab-case verb) |
| `service` | env var `SERVICE_NAME` | `identity` / `catalog` / ... |
| `version` | build-time `-ldflags -X` | git short SHA |
| `requestId` | `common/middleware.RequestIDMiddleware` | UUID v7 |
| `userId` | from JWT claims if present | empty for unauthenticated |
| `path` | request path | `/api/v1/identity/user/login` |
| `method` | HTTP method | `POST` |
| `status` | response status code | `200` / `401` / `500` |
| `latencyMs` | wall-clock from middleware entry to write | float |
| `traceId` | `wrapper.CtxTraceID` (W3C `traceparent`) | hex 32 |
| `spanId` | from W3C `traceparent` | hex 16 |

**Redaction list** (NEVER logged):

- Passwords, password hashes
- Raw access JWTs, raw refresh JWTs, JWT signing keys
- The internal shared secret (env `INTERNAL_SHARED_SECRET`) — neither the configured value nor the supplied `X-Internal-Secret` header value may appear in logs. The `internal_auth_reject` log line includes `remoteAddr`, `x_forwarded_for`, `path`, `requestId`; never the supplied secret.
- HMAC keys (n/a in MVP — payment callback is in-process)
- Customer payment card data (n/a — mock payment only)
- Full request bodies on `identity/user/login` + `identity/user/register` + `payment/intent/callback` — log only field NAMES, never values

Per `common/serror`, errors wrap with source location; the wrap chain is logged via `slog.Error(msg, "err", err)` — slog renders the chain in the `err` field.

**Internal-auth audit posture** (REV-L2-002 + REV-L2-003 closure):

- On startup, every recipient service (checkout, order, inventory, payment) logs:
  ```
  slog.Info("internal_auth_configured",
      "service", os.Getenv("SERVICE_NAME"),
      "sha256_prefix", hex.EncodeToString(sha256.Sum256([]byte(secret))[:4]))
  ```
  Ops can confirm at a glance that all four services hash the same secret without ever logging the secret itself.
- On every reject, recipient logs:
  ```
  slog.Error("internal_auth_reject",
      "remoteAddr", c.Request.RemoteAddr,
      "x_forwarded_for", c.Request.Header.Get("X-Forwarded-For"),
      "path", c.Request.URL.Path,
      "requestId", c.GetString("requestId"))
  ```

## Metrics

Each service exposes `/metrics` on its main port. Prometheus scrapes every 15s.

Standard metrics (every service):

| Metric | Type | Labels |
|---|---|---|
| `http_requests_total` | counter | `service, path, method, status` |
| `http_request_duration_seconds` | histogram (le buckets: 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10) | `service, path, method, status` |
| `http_in_flight` | gauge | `service` |
| `db_query_duration_seconds` | histogram | `service, query_kind` |
| `db_connections_in_use` | gauge | `service, pool` |
| `kafka_messages_published_total` | counter | `service, topic, status` |
| `kafka_messages_consumed_total` | counter | `service, topic, status` |
| `outbox_pending_count` | gauge | `service` |
| `internal_auth_rejects_total` | counter | `service, path` |

Per-domain metrics:

| Service | Metric | Type |
|---|---|---|
| inventory | `inventory_reservations_total{outcome=created|expired|committed|released}` | counter |
| inventory | `inventory_reservation_sweep_batch_size` | histogram |
| inventory | `inventory_reservation_pending_expiry` | gauge |
| inventory | `inventory_stock_negative_blocked_total` (CHECK constraint hit count) | counter |
| payment | `payment_callback_dedup_hits_total{kind=replay|terminal_mismatch}` | counter |
| payment | `payment_intent_state_total{state=requires_payment|succeeded|failed|expired|cancelled}` | counter |
| checkout | `checkout_orchestration_step_duration_seconds{step=load_cart|fetch_email|validate_address|preview_stock|reserve|create_order|create_intent|clear_cart}` | histogram |
| checkout | `checkout_saga_status_total{status=committed|degraded|failed_no_compensation}` | counter |
| checkout | `checkout_idempotency_replay_total` | counter |
| order | `order_state_transitions_total{from_status, to_status, actor_role}` | counter |
| order | `order_consumer_late_event_total{event_type, current_status}` (logged when state-driven consumer no-ops) | counter |
| catalog | `catalog_inventory_call_failures_total` (when listing with inStock=true) | counter |
| identity | `identity_login_failures_total` (per IP/min surfaced via the alert rule) | counter |
| identity | `identity_refresh_rotation_total{outcome=ok|reused|revoked}` | counter |

## Distributed tracing

- **Propagation:** W3C `traceparent` header on every HTTP request and Kafka message header. All services use `common/middleware.TraceContextTraceIDMiddleware` which wires the trace id into `wrapper.CtxTraceID` so it surfaces in the response envelope's `traceId` field automatically.
- **Exporter:** OpenTelemetry SDK with OTLP gRPC exporter to Tempo (default endpoint `tempo:4317`).
- **Span boundaries:** every public HTTP handler creates a server span via `otelgin`; every downstream HTTP call and Kafka publish is a child span. Database queries are spans only when query duration > 50ms (sampling).
- **Sampling:** parent-based, with 10% probabilistic sampling for traces that originate in the frontend; 100% sampling for traces that originate in error paths.
- **Baggage:** never propagate user PII via baggage.

## Audit log

Audit Service is OUT_OF_SCOPE for the MVP run per BA. When re-enabled, every admin-initiated mutation emits an audit event with this taxonomy:

| Field | Required | Example |
|---|---|---|
| `event_id` | Yes | UUID v7 |
| `actor_user_id` | Yes | claims.sub of admin |
| `actor_role` | Yes | `ADMIN` |
| `action` | Yes (enum) | `CREATE_PRODUCT` / `UPDATE_PRODUCT` / `SOFT_DELETE_PRODUCT` / `CREATE_CATEGORY` / `DEACTIVATE_CATEGORY` / `ADJUST_STOCK` / `UPDATE_ORDER_STATUS` / `CANCEL_ORDER` |
| `entity_type` | Yes | `product` / `category` / `stock` / `order` |
| `entity_id` | Yes | UUID |
| `before_summary` | Optional | redacted snapshot |
| `after_summary` | Optional | redacted snapshot |
| `correlation_id` | Yes | requestId of the originating request |
| `occurred_at` | Yes | RFC 3339 UTC |

Sensitive fields (passwords, payment data, secret material) MUST be redacted in `before_summary` / `after_summary`.

In MVP the equivalent posture is achieved per-service:
- `inventory.stock_adjustments` row carries `actor_user_id`, `delta`, `reason`, `created_at` — local audit (INV-007).
- `order.order_status_history` row carries `actor_user_id`, `actor_role`, `from_status`, `to_status`, `reason`, `at` — local audit (ORD-009).
- `order.admin_action_log` row captures admin-cancel reason on PAID->CANCELLED.

## Dashboards

| Dashboard | Intent | Key panels |
|---|---|---|
| Service Health | per-service overview | request rate, error rate, p95/p99 latency, in-flight, db connections in use, outbox pending |
| Customer Journey funnel | spine end-to-end | browse -> cart -> checkout commit -> payment success rate; per-step error rate; checkout step duration histogram |
| Saga Compensation | catch stuck flows | reservations RESERVED > 5min count; orders PENDING_PAYMENT > 15min count; checkout_saga_status_total{status='degraded'} rate; outbox_pending_count |
| Auth Anomaly | spot abuse | identity_login_failures_total per IP/min; identity_refresh_rotation_total{outcome='reused'}; AUTH_REVOKED rate; internal_auth_rejects_total |
| State Machine | order transitions | order_state_transitions_total heatmap (from_status x to_status); order_consumer_late_event_total (PR-003 race signal) |
| Schema Performance | query hotness | top 10 slow queries by p95 (from db_query_duration_seconds); index hit ratio; lock wait events |

## Alert rules

```yaml
groups:
  - name: service-health
    rules:
      - alert: HighErrorRate
        expr: |
          sum(rate(http_requests_total{status=~"5.."}[5m])) by (service) /
          sum(rate(http_requests_total[5m])) by (service) > 0.05
        for: 5m
        labels: { severity: page }
        annotations:
          summary: "{{ $labels.service }} error rate > 5% over 5min"

      - alert: P95LatencyBreach
        expr: |
          histogram_quantile(0.95,
            sum(rate(http_request_duration_seconds_bucket[5m])) by (le, service, path)
          ) > 1
        for: 10m
        labels: { severity: warn }

      - alert: SagaInFlightStuck
        expr: |
          (inventory_reservations_total{outcome="created"} -
           (inventory_reservations_total{outcome="expired"} +
            inventory_reservations_total{outcome="committed"} +
            inventory_reservations_total{outcome="released"})) > 0
          and inventory_reservation_pending_expiry > 0
        for: 15m
        labels: { severity: page }
        annotations:
          summary: "Inventory has > 0 reservations stuck past TTL — sweeper is degraded"

      - alert: SagaDegradedRate
        expr: |
          sum(rate(checkout_saga_status_total{status="degraded"}[10m])) > 0.05
        for: 10m
        labels: { severity: page }
        annotations:
          summary: "Checkout DEGRADED saga rate > 5% — investigate compensation path"

      - alert: OutboxBacklog
        expr: outbox_pending_count > 1000
        for: 10m
        labels: { severity: warn }

      - alert: InternalAuthRejectSpike
        expr: sum(rate(internal_auth_rejects_total[5m])) by (service) > 10
        for: 5m
        labels: { severity: page }
        annotations:
          summary: "{{ $labels.service }} internal-auth rejects > 10/sec — possible secret rotation gap or attack"

      - alert: PaymentDedupReplaySpike
        expr: rate(payment_callback_dedup_hits_total{kind="replay"}[5m]) > 1
        for: 10m
        labels: { severity: warn }
        annotations:
          summary: "Payment callback replay rate elevated — provider redelivery storm"
```

## Local stack

Reference: `B2C E-Commerce Platform/observability/docker-compose.yml`. Services:

- `loki` (3100) — log aggregation; `promtail` sidecar tails stdout from each service container.
- `tempo` (3200, 4317 OTLP) — trace storage.
- `prometheus` (9090) — metrics scrape; reads `prometheus.yml` listing all backend service `/metrics` endpoints.
- `grafana` (3000) — dashboards; auto-provisioned from `dashboards/`.

Run: `docker-compose -f observability/docker-compose.yml up -d`.

## Per-service obligations

Each service's Dev output MUST satisfy:

- [ ] `slog` JSON handler wired in `main.go` with required field set.
- [ ] `common/middleware.RequestIDMiddleware` registered (sets `requestId` UUID v7 in context).
- [ ] `common/middleware.TraceContextTraceIDMiddleware` registered (sets `wrapper.CtxTraceID` in context; `wrapper.Respond` auto-injects into `Response.TraceID`).
- [ ] `common/middleware.AccessLogMiddleware` registered (logs the per-request envelope with the field set above).
- [ ] `/metrics` endpoint exposed via `promhttp.Handler()` or equivalent — UNAUTHENTICATED (Prometheus has no JWT).
- [ ] OTLP exporter initialized in `main.go`; spans created via `otelgin` middleware.
- [ ] Per-domain counters declared in `app/<domain>/metrics.go` per the per-service table above.
- [ ] Redaction list enforced in `common/logger`'s replacer.
- [ ] For checkout/order/inventory/payment: `internal_auth_configured` startup log line + `internal_auth_reject` log line on every reject (NEVER the secret value); `internal_auth_rejects_total` counter incremented on reject.
- [ ] For Order: `order_state_transitions_total` incremented on every status mutation (sole-writer rule means one site per transition); `order_consumer_late_event_total` incremented when the state-driven consumer no-ops (PR-003 visibility).
- [ ] For Inventory: `inventory_reservations_total{outcome=...}` incremented at every reservation lifecycle change; `inventory_stock_negative_blocked_total` incremented every time the CHECK constraint catches an attempted negative.
- [ ] For Payment: `payment_callback_dedup_hits_total{kind=...}` incremented on every dedup unique-violation catch.
- [ ] For Checkout: `checkout_orchestration_step_duration_seconds` recorded around every step 1-9; `checkout_saga_status_total` incremented per saga outcome; `checkout_idempotency_replay_total` incremented on duplicate Idempotency-Key replay.
