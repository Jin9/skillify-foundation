# Rust Thinking Model (L1 -> L4)

Apply this layered sequence before writing or recommending Rust code. Do not jump to implementation without passing through contract and ownership questions.

## L1 — Business Invariant

What rule must remain true after this change? What is the success state and who owns it?

Surface the invariant in plain language before naming any crate or type. Examples:
- "A loan disbursement occurs at most once per approved application."
- "A consumer that crashes mid-batch must not double-process the same message."
- "An override is recorded with actor identity and reason even under partition."

If the invariant is unclear, stop and ask. Do not invent invariants.

## L2 — Crate / Module Boundary

Name the boundary the change lives inside.

- **Crate / module ownership**: which crate or `mod` owns the change; which binary or library exposes it.
- **API contract**: the public shape promised to callers — public `fn` signatures, `trait` definitions, request/response `struct`s, error `enum`, message schema. Public API is forever until a SemVer bump; treat it as a contract.
- **Data ownership**: which crate writes the source of truth; everyone else reads or subscribes.
- **Transactional boundary**: the unit of atomicity (one `sqlx::Transaction` per aggregate). Cross-aggregate writes in one transaction are a smell — reach for outbox / saga / compensation.
- **Authorization boundary**: who may invoke the contract; how identity is asserted and where it is checked.

## L3 — Technical Strategy

Choose the implementation strategy without writing code yet.

- **Error model**: domain/library errors as a `thiserror` enum with `#[from]`/`#[source]`; `anyhow` only in binaries / top-level glue. Each variant carries a class (client / server / dependency) that drives HTTP status, retry, and the `tracing` attribute. See `rust-idioms-and-async.md` § 1.
- **Async strategy**: sync vs `async`; runtime (`tokio` multi-thread vs current-thread); where `.await` points sit; cancellation propagation; whether any task needs `spawn`/`JoinSet`. See `rust-idioms-and-async.md` § 2.
- **Ownership & sharing**: who owns each piece of state; borrow vs move vs `Arc`; avoid `Arc<Mutex<_>>` as a default — prefer message-passing (a manager task owning state, requests over a channel) when state is shared across tasks.
- **Persistence model**: table schema, indexes, `sqlx` compile-time macros, transaction scope, `FOR UPDATE` on read-then-write. See `rust-persistence-and-data.md`.
- **Trait vs generic**: trait objects (`dyn`, dynamic dispatch, object-safe) vs generics (`impl Trait` / `<T: Trait>`, monomorphized) — pick by call-site count, binary-size, and whether the boundary needs to be stored heterogeneously.
- **Idempotency**: command-id keys, dedup window, consumer-side dedup table; replay returns the stored response.
- **Observability**: `#[tracing::instrument]` on handlers, structured fields, metrics per failure mode, bounded label cardinality.
- **Testing approach**: `#[cfg(test)]` unit boundary, `tests/*.rs` integration, `#[tokio::test]` async, contract tests against the public API, migration dry-runs, injected clock for determinism.

## L4 — Implementation

Now write the smallest safe change. Map the Fowler/DDD intent onto Rust constructs:

| Architectural role | Rust construct |
|--------------------|----------------|
| Handler / transport | thin `async fn` handler; validate input, call use case, format response |
| Service / use case | module `fn` or method holding the invariant, calling repositories |
| Repository / adapter | a `trait` + impl struct; one per dependency; injected via constructor |
| Domain model | owned `struct` with private fields + intention-revealing methods (Tell, Don't Ask) |
| Value Object | newtype (`struct Money(Decimal)`, `struct MerchantId(Uuid)`) guarding invalid states |
| DTO | `#[derive(Serialize, Deserialize)]` struct at the boundary; `#[serde(deny_unknown_fields)]` |
| Sentinel / special case | a `thiserror` enum variant (`#[error("...")] NotFound`) |
| Strategy / polymorphism | a `trait` with bounded impls, replacing a conditional explosion |

- **Handlers**: thin transport; no business decisions.
- **Repositories**: persistence and external IO behind a trait; one place per dependency.
- **Schemas / migrations**: additive when possible; expand-then-contract for breaking changes; never edit a committed migration.
- **Consumers**: idempotent, dead-letter on poison, correlation id in the span.
- **Tests**: regression test for the bug or invariant; contract test if the boundary changed.

## Evaluation axes

For every L2/L3 decision, name where it lands on at least two of these axes:

- correctness ↔ latency
- safety / ownership clarity ↔ ergonomics (e.g. `Arc<Mutex>` vs message-passing; `clone()` vs borrow)
- coupling ↔ autonomy
- simplicity ↔ operability

Stating the trade-off keeps the design honest. A choice that claims to win on every axis is usually under-analyzed.

## Fast path

For isolated compile/borrow errors, narrow test fixes, or one-line query fixes, compress L1-L3 into one sentence then jump to L4. State `L1-L3 skipped: isolated fix` so traceability stays explicit.
