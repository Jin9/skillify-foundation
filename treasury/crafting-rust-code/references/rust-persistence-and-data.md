# Rust Persistence, Serde & Observability

Generation-side guidance for data access, input validation, and instrumentation. Constructive companion to `review-rust-code`'s R9/R10/R11 scans.

## 1. Persistence (sqlx)

- **Compile-time-checked SQL**: prefer `sqlx::query!` / `query_as!` / `query_scalar!` over runtime `sqlx::query(...)`. The macro checks the SQL against the database schema at build time and infers result types.
- **Never build SQL by `format!` or `+`**: always parameterize (`$1`, `?`, `.bind(...)`). String-built SQL is an injection vector and defeats compile-time checking (B10/R9).
- **Offline cache**: when using the `query!` macros, commit the `.sqlx/` offline cache (`cargo sqlx prepare`) so builds and CI succeed without a live database. In a sandbox without a DB, rely on this cache and note it (`unavailable: no live DB; using .sqlx/ offline cache`).
- **Transactions**: wrap multi-statement writes in a `sqlx::Transaction` (`let mut tx = pool.begin().await?; ...; tx.commit().await?;`). Always `commit()` on the success path — relying on rollback-on-drop silently discards a "successful" write. Cross-aggregate work uses outbox/saga/compensation, not one wide transaction.
- **Read-then-write of a contended row** (balances, counters): use `SELECT ... FOR UPDATE` or serializable isolation to avoid lost updates.
- **Pool config from config**, not a hardcoded `max_connections(<literal>)`.
- **Migrations are immutable once committed**: add a new additive migration; never edit a shipped one.

## 2. Input validation & serde

- **Harden inbound DTOs**: every externally-deserialized struct carries `#[serde(deny_unknown_fields)]` so unexpected payloads are rejected, and a consistent `rename_all` across the API surface.
- **Validate at the boundary**, not deep in the handler: implement `TryFrom<RawDto>` for the validated domain type, or use `validator::Validate`. The handler should receive an already-valid value.
- **PII discipline**: a struct holding PII must not derive a naive `Debug`/`Display` that prints it — implement a redacting `Debug` or wrap the field in a redacting newtype. PII fields in responses carry `#[serde(skip_serializing)]` unless the contract requires them.

## 3. Observability (tracing)

- **Instrument handlers**: `#[tracing::instrument(skip(...), fields(...))]` on each handler / worker entry. `skip` large or sensitive args; add the correlation fields you want on every child span.
- **Structured fields, not format strings**: `tracing::info!(order_id = %id, amount = ?amt, "order placed")`, not `info!("order {} placed for {}", id, amt)`. Structured fields are queryable; format strings are not.
- **Bounded cardinality**: never use unbounded values (user IDs, request IDs, account IDs, IPs, emails, free-form strings) as `metrics` labels or recurring span fields — they explode the metrics backend. Use them as span fields on a single span, not as metric labels.
- **A metric + log + span per declared failure mode**: each failure the design names gets a `metrics::counter!` increment, a `tracing::error!`/`warn!` line, and lives inside a span. No `info!`/`warn!` inside a tight loop.
- **Trace propagation**: extract the `traceparent` header at the HTTP boundary and make it the parent of the root handler span so traces stitch across services.
