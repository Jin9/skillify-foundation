# Banking-Grade Rust Patterns Reference

A deterministic, checklist-friendly reference for adversarial code review of one Rust feature in a banking workflow. Every section ends with an **adversarial scan question** the reviewer must ask. Findings are cited as `file:line`. The final table maps each pattern to its default severity and review-skill category.

---

## Section 1: Money and Arithmetic

### 1.1 Never use `f32` / `f64` for money

IEEE-754 binary floats cannot represent `0.10` exactly: the closest `f64` value is `0.1000000000000000055511151231257827021181583404541015625`. Summing `0.1 + 0.2` yields `0.30000000000000004`. Any ledger that uses `f64` will silently drift.

**Accepted alternatives**, in preference order for transactional ledger code:

| Type | When to use |
|---|---|
| **Integer minor units** (`i64` cents / `i128` satoshis) with a currency tag | Default for ledger entries, balances, postings. Deterministic, no rounding ambiguity, trivial to serialize. |
| **`rust_decimal::Decimal`** | When fractional precision beyond minor units matters (FX, interest accrual). 96-bit mantissa, max ~28 significant digits. |
| **`bigdecimal::BigDecimal`** | Arbitrary precision; only when `rust_decimal`'s 28-digit ceiling is provably insufficient. Slower and heap-allocating. |

See `rust_decimal` precision constraints at https://docs.rs/rust_decimal/latest/rust_decimal/struct.Decimal.html.

### 1.2 Integer overflow discipline

`i64 + i64` panics in debug builds (because `overflow-checks = true`) but **wraps silently in release**. Tests passing in debug do not prove production safety. Per https://doc.rust-lang.org/std/primitive.i64.html:

| Method | Semantics | Use when |
|---|---|---|
| `checked_add` | `Option<i64>`; `None` on overflow | Default for money. Force the caller to handle overflow. |
| `saturating_add` | Clamps at `MIN`/`MAX` | Quota counters, never balances. |
| `wrapping_add` | Two's complement wrap | Cryptography, hashing. **Never** for money. |
| `overflowing_add` | `(i64, bool)` | When the overflow bit must be inspected. |

Bare `+` / `-` / `*` on a money type is a finding unless the surrounding type is a newtype whose `Add` impl is itself `checked_*`.

### 1.3 Rounding rules

`rust_decimal::RoundingStrategy` (https://docs.rs/rust_decimal/latest/rust_decimal/enum.RoundingStrategy.html) variants:

- `MidpointNearestEven` (banker's rounding) — usually the legal default in EU/IFRS contexts.
- `MidpointAwayFromZero` — "round half up" in plain language.
- `MidpointTowardZero`
- `ToZero`, `AwayFromZero`, `ToNegativeInfinity`, `ToPositiveInfinity`

Every rounding call site MUST specify a strategy explicitly:

```rust
use rust_decimal::{Decimal, RoundingStrategy};
let fee = gross.round_dp_with_strategy(2, RoundingStrategy::MidpointNearestEven);
```

A call to `.round_dp(2)` or `.round()` without an explicit strategy is a finding — the default may change across versions and is opaque at the call site.

### 1.4 Currency tagging

A bare `Decimal` or `i64` traveling through a function signature loses currency information. Two `Decimal` values of unknown currency cannot be safely added. Require a newtype:

```rust
pub struct Money<C: Currency>(Decimal, PhantomData<C>);
```

or an explicit `currency: Currency` parameter alongside every amount. Mixed-currency arithmetic must be a compile error or an explicit `convert()` call against a `FxRate`.

### Adversarial scan questions

- Cite any `f32` / `f64` declaration that holds a monetary value at `file:line`.
- Cite any `Decimal` or integer arithmetic on a money value at `file:line` that uses `+` / `-` / `*` instead of `checked_*`.
- Cite a `Decimal` rounding call at `file:line` that omits the `RoundingStrategy` argument.
- Cite a function signature accepting `amount: Decimal` (or `i64`) without a currency parameter or newtype wrapper.

---

## Section 2: Time

### 2.1 Pick one time crate per repo

| Crate | Notes |
|---|---|
| **`chrono`** | Most mature; absorbed `time 0.1`'s platform functionality in 0.4.30. https://docs.rs/chrono/latest/chrono/ |
| **`time`** | Smaller surface, no implicit locale. |
| **`jiff`** | Newer, tzdb-aware, designed against `chrono`'s pitfalls. |

Mixing two time crates in one repo is a finding — round-trip conversions through `SystemTime` lose precision and zone metadata.

### 2.2 UTC discipline

`DateTime<Utc>` only at boundaries; `DateTime<Local>` is forbidden in persistence layers and across service boundaries. `Local::now()` reads the OS timezone, which is non-deterministic across deployment regions.

### 2.3 Monotonic vs wall clock

Per https://doc.rust-lang.org/std/time/struct.Instant.html:

- `std::time::Instant` is **monotonically nondecreasing**. Use for durations, timeouts, retry backoffs, latency histograms.
- `std::time::SystemTime` and `chrono::DateTime<Utc>` are **wall-clock**, subject to NTP step and admin adjustment. Use for timestamps written to a row or emitted in an event.

Mixing the two — e.g. subtracting two `SystemTime`s to measure elapsed work — produces negative durations on NTP step. Always prefer `Instant::elapsed()` for measurement.

### 2.4 Test-time injection

Banking handlers MUST NOT call `Utc::now()` or `Instant::now()` directly. Inject a clock:

```rust
pub trait Clock: Send + Sync {
    fn now(&self) -> DateTime<Utc>;
}
```

In tests, either substitute a fake `Clock` or use `tokio::time::pause()` (https://docs.rs/tokio/latest/tokio/time/fn.pause.html) — note this only freezes Tokio's `Instant`, not `std::time::Instant` or `chrono::Utc::now()`.

### Adversarial scan questions

- Cite any direct call to `Utc::now()` / `Local::now()` / `SystemTime::now()` / `Instant::now()` inside a handler or domain function at `file:line`.
- Cite any subtraction of two `SystemTime` or `DateTime<Utc>` values used as an elapsed measurement.
- Cite any use of `DateTime<Local>` in a struct that crosses the persistence or network boundary.

---

## Section 3: Persistence with `sqlx`

### 3.1 Macros over runtime queries

`sqlx::query!`, `sqlx::query_as!`, `sqlx::query_scalar!` are compile-time checked against the live database or the offline cache. Runtime `sqlx::query("...").bind(x)` is acceptable only when the SQL is dynamic in a way the macros cannot express (e.g. dynamic column lists from an allowlisted enum). https://docs.rs/sqlx/latest/sqlx/

```rust
let row = sqlx::query_as!(Account, "SELECT id, balance_cents FROM accounts WHERE id = $1", id)
    .fetch_one(&pool).await?;
```

### 3.2 `SQLX_OFFLINE` and `.sqlx/`

Reviewers should expect a `.sqlx/` directory in version control when the project uses `query!`. The macros consult `.sqlx/query-<hash>.json` when `SQLX_OFFLINE=true`. A repo with `query!` calls but no `.sqlx/` directory fails reproducible builds.

### 3.3 No string-built SQL

`format!("SELECT ... WHERE id = {}", id)` is a SQL injection finding regardless of the call site. Every variable must be a bind parameter (`$1`, `?`, `$N`).

### 3.4 Transactions

```rust
let mut tx = pool.begin().await?;
sqlx::query!(...).execute(&mut *tx).await?;
sqlx::query!(...).execute(&mut *tx).await?;
tx.commit().await?;
```

If `tx` is dropped without `commit()`, the destructor rolls back. This means a `?` short-circuit between `begin()` and `commit()` is **safe by default** — the rollback is the right behavior. However:

- A `tx.rollback().await?` made explicit at every error path is preferred for auditability.
- A function that `begin()`s but never `commit()`s on the success path is a finding (silent rollback on success).

### 3.5 Pool sizing

```rust
PgPoolOptions::new()
    .max_connections(config.db.max_connections)
    .connect(&config.db.url).await?
```

Hard-coded `.max_connections(5)` is a finding — pool size must come from config so production can size against `max_connections` on the Postgres side.

### 3.6 Migrations

`sqlx::migrate!("./migrations")` (https://docs.rs/sqlx/latest/sqlx/macro.migrate.html) embeds files matching `YYYYMMDDHHMMSS_name.sql`. Discipline:

- Additive only. No destructive `DROP COLUMN` in the same migration as code that reads it.
- Idempotent where possible (`CREATE TABLE IF NOT EXISTS`, `ADD COLUMN IF NOT EXISTS`).
- A `build.rs` containing `println!("cargo:rerun-if-changed=migrations");` to force rebuild on migration edits.

### 3.7 Row-level locking

- `SELECT ... FOR UPDATE` legitimately appears when reading-then-writing a balance row inside a transaction.
- `pg_advisory_xact_lock(key)` legitimately appears when coordinating cross-row work (e.g. one-and-only-one outbox processor).
- `SELECT ... FOR UPDATE SKIP LOCKED` is the canonical pattern for a work-queue dequeue.

Reads of a balance followed by writes **without** `FOR UPDATE` and **without** a serialization-level transaction is a P1 correctness finding.

### Adversarial scan questions

- Cite any `format!` or `+` string concatenation that feeds a SQL string.
- Cite any `sqlx::query` (runtime) call where a `sqlx::query!` macro would compile equivalently.
- Cite any `PgPoolOptions::max_connections(<literal>)` not sourced from config.
- Cite any read-then-write of a money column without `FOR UPDATE` or equivalent isolation.

---

## Section 4: Observability

### 4.1 `tracing` idioms

Every handler MUST carry `#[tracing::instrument(skip(...), fields(...))]`:

```rust
#[tracing::instrument(skip(pool, body), fields(account_id = %req.account_id))]
async fn transfer(pool: &PgPool, body: TransferDto, req: AuthCtx) -> Result<...> { ... }
```

Structured fields — `field = %value` (Display) or `field = ?value` (Debug) — are mandatory. Per https://docs.rs/tracing/latest/tracing/, the macros recognize these sigils. Format-string logging (`info!("user = {}", id)`) is a finding: it produces an unstructured message field that backends cannot index.

### 4.2 Log levels

| Level | Use |
|---|---|
| `error!` | Faults the operator must investigate. |
| `warn!` | Degraded conditions, retries, fallbacks. |
| `info!` | One-per-request boundary events. Never inside a hot loop. |
| `debug!` | Development; gated off in prod. |
| `trace!` | Verbose; gated off in prod. |

`info!` inside a `for` loop over a result set is a finding — it produces unbounded log volume.

### 4.3 Label cardinality

`tracing` span fields and `metrics` labels must be **bounded**. Unbounded values forbidden as labels:

- User IDs, account IDs, request IDs, IPs, emails, free-form strings.
- High-cardinality identifiers belong in **span fields for traces** (sampled), not in **metric labels** (every value gets a time series). https://docs.rs/metrics/latest/metrics/ explicitly notes this cost.

Acceptable label values: enum-like bounded sets — `method`, `endpoint_template`, `status_class`, `currency`, `region`.

### 4.4 `metrics` crate

```rust
metrics::counter!("transfers_total", "currency" => currency.as_str()).increment(1);
metrics::histogram!("transfer_duration_seconds").record(elapsed.as_secs_f64());
```

https://docs.rs/metrics/latest/metrics/ — `counter!`, `gauge!`, `histogram!`. Same cardinality rules as tracing.

### 4.5 Trace propagation

`tracing-opentelemetry` (https://docs.rs/tracing-opentelemetry/latest/tracing_opentelemetry/) wires `tracing` spans into the OpenTelemetry `Context`, then an `opentelemetry-otlp` exporter ships to an OTLP collector. Trace ID must be extracted from inbound request headers (`traceparent`) at the HTTP layer and made the parent of the root handler span.

### Adversarial scan questions

- Cite any handler function without `#[tracing::instrument]` at `file:line`.
- Cite any `info!("…{}…", …)` format-string log instead of structured fields.
- Cite any span field or metric label whose value is a user ID, free-form string, or otherwise unbounded.
- Cite any `info!` / `warn!` inside a loop body.

---

## Section 5: Input Validation and `serde`

### 5.1 `deny_unknown_fields`

Per https://serde.rs/container-attrs.html, `#[serde(deny_unknown_fields)]` errors on unknown keys during deserialization. Every externally-sourced DTO must carry it. Without it, a client sending `{"amount": 100, "amunt": 50}` is silently accepted with `amount = 100`.

```rust
#[derive(Deserialize)]
#[serde(deny_unknown_fields, rename_all = "snake_case")]
pub struct TransferDto { pub from: AccountId, pub to: AccountId, pub amount: Money<Usd> }
```

### 5.2 Case-convention consistency

`#[serde(rename_all = "snake_case")]` (or `camelCase` — pick one project-wide). Mixed conventions across DTOs in the same API surface is a finding.

### 5.3 Validation as deserialization

Two patterns, both acceptable:

**Pattern A — `TryFrom<RawDto>`.** Raw deserialization yields a `RawDto`; `TryFrom<RawDto>` produces a validated `Dto` or `ValidationError`. The raw type must be `pub(crate)` and documented as accepted-but-unvalidated.

```rust
#[derive(Deserialize)]
pub(crate) struct RawTransfer { amount: String, currency: String }

pub struct Transfer { amount: Money<Usd> }

impl TryFrom<RawTransfer> for Transfer { type Error = ValidationError; /* ... */ }
```

**Pattern B — `validator` crate.** https://docs.rs/validator/latest/validator/ — `#[derive(Validate)]` plus per-field `#[validate(range(min=0))]`, `#[validate(email)]`, etc. Caller invokes `dto.validate()?` after deserialization.

Mixing both patterns in the same project is a smell. Validating deep inside a handler (rather than at the deserialization boundary) is a P2 finding — it lets unvalidated values propagate.

### 5.4 PII redaction

- `#[serde(skip_serializing)]` on fields that must never appear in a response body.
- `#[serde(skip)]` on fields that must never appear in either direction.
- Custom `Debug` / `Display` impls that print `***` instead of the raw value:

```rust
impl std::fmt::Debug for Pan {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "Pan(****{})", &self.0[self.0.len() - 4..])
    }
}
```

A PII-bearing struct that uses `#[derive(Debug)]` is a finding — `Debug` output reaches logs.

### Adversarial scan questions

- Cite any externally-deserialized DTO missing `#[serde(deny_unknown_fields)]`.
- Cite any handler that calls `serde_json::from_str` (or equivalent) without an immediate `try_into()` / `.validate()`.
- Cite any struct with PII fields that uses `#[derive(Debug)]` or default `Display`.

---

## Section 6: Tests

### 6.1 `#[tokio::test]` flavor

```rust
#[tokio::test(flavor = "current_thread")]
async fn unit() { ... }

#[tokio::test(flavor = "multi_thread", worker_threads = 2)]
async fn integration() { ... }
```

`current_thread` is the default and avoids requiring `Send` on captured locals. Choose `multi_thread` only when the test specifically exercises cross-thread behavior; that choice imposes `Send` bounds on every `.await`ed future. Mismatched flavor against the production runtime is a smell.

### 6.2 `rstest` parameterized cases

https://docs.rs/rstest/latest/rstest/

```rust
#[rstest]
#[case(Decimal::ZERO, Decimal::ZERO, Decimal::ZERO)]
#[case(dec!(0.1), dec!(0.2), dec!(0.3))]
fn add(#[case] a: Decimal, #[case] b: Decimal, #[case] expected: Decimal) {
    assert_eq!(a + b, expected);
}
```

### 6.3 HTTP test doubles via `wiremock`

https://docs.rs/wiremock/latest/wiremock/ — `MockServer::start().await`, then `Mock::given(method("POST")).and(path("/charge")).respond_with(ResponseTemplate::new(200)).mount(&server).await`. A test that calls `reqwest::get("https://real.host.example/...")` is a P1 finding — non-deterministic, network-dependent.

### 6.4 Property-based tests

`proptest` (https://docs.rs/proptest/latest/proptest/) for invariants:

```rust
proptest! {
    #[test]
    fn debit_credit_preserves_total(a in 0i64..1_000_000, b in 0i64..1_000_000) {
        let mut acct = Account::with(a);
        let mut other = Account::with(0);
        transfer(&mut acct, &mut other, b).ok();
        prop_assert_eq!(acct.balance() + other.balance(), a);
    }
}
```

`quickcheck` is an older alternative; `proptest` is preferred for new code because of integrated shrinking and `prop_assert!` (which returns a failure instead of panicking from inside the harness).

### 6.5 `cargo nextest`

https://nexte.st/ — process-per-test isolation, ~3x faster than `cargo test`. Assertions that rely on shared global state (`static mut`, `OnceCell`) may pass under `cargo test` and fail under `nextest`; that exposes a real bug, not a nextest bug.

### 6.6 Doctests

Legitimate for documenting public crate APIs with copy-pasteable examples. **Not** legitimate as a substitute for unit tests — they cannot use `mockall`, are slow to compile, and run sequentially.

### 6.7 Fixtures, fakes, mocks

| Pattern | When |
|---|---|
| Hand-rolled trait stub | Small surface, one or two methods. |
| `mockall::automock` | Trait with several methods, expectation matching needed. https://docs.rs/mockall/latest/mockall/ |
| `mockall::mock!` | Complex trait the `automock` attribute cannot handle (generics, multiple impls). |
| Fakes (in-memory `Repo` impl) | Multi-test reuse of stateful behavior. |

### 6.8 No real PII in fixtures

Synthetic data only — `fake` crate generators, or hand-crafted obviously-fake strings (`"alice@example.test"`, `"4242 4242 4242 4242"`). Real PANs, real emails, real names from production datasets in fixtures is a P1 finding.

### Adversarial scan questions

- Cite any test that issues a real network request (no `MockServer` between it and the host).
- Cite any test relying on `Utc::now()` or wall-clock without `tokio::time::pause()` or a fake clock.
- Cite any fixture or test data file containing a real-looking PII value (PAN passing Luhn, valid TLD email).

---

## Section 7: Supply Chain and Security

### 7.1 `cargo audit`

https://rustsec.org/ — `cargo audit` cross-references `Cargo.lock` against the RUSTSEC advisory database. Advisory IDs follow `RUSTSEC-YYYY-NNNN`. Also flags yanked crates (a crate version withdrawn from crates.io). CI must fail on any unresolved advisory.

### 7.2 `cargo deny`

https://embarkstudios.github.io/cargo-deny/ — broader scope than `cargo-audit`. A `deny.toml` configures four checks:

| Check | Enforces |
|---|---|
| `advisories` | RUSTSEC IDs, severity floor. |
| `licenses` | Allowlist of SPDX identifiers. |
| `bans` | Banned crates, banned multiple-versions, banned wildcards. |
| `sources` | Allowed registries / git origins. |

Banking-grade repos pin a committed `deny.toml`. Absence of `deny.toml` is a finding.

### 7.3 `cargo vet`

https://github.com/mozilla/cargo-vet — supply-chain audit *attestations*. Each transitive dependency must be marked as audited by a trusted party (your team, Mozilla, Google, etc.). Attestations are stored in `supply-chain/audits.toml`.

### 7.4 `cargo-semver-checks`

https://github.com/obi1kenobi/cargo-semver-checks — detects removed public items, visibility narrowing, attribute changes that constitute SemVer breakage. Run pre-publish on internal published crates.

### 7.5 `cargo-msrv`

https://github.com/foresterre/cargo-msrv — pins the Minimum Supported Rust Version. Compliance demands an MSRV declared in `Cargo.toml` (`rust-version = "1.XX"`) and verified in CI. Unpinned MSRV means a `cargo update` can silently raise the toolchain requirement.

### 7.6 Yanked crates

`Cargo.lock` records exact versions; if a maintainer yanks a version (e.g. it published a regression), `cargo audit` reports `yanked` against that line. Yanked dependencies block CI.

### 7.7 `build.rs` review

Per https://doc.rust-lang.org/cargo/reference/build-scripts.html, build scripts execute with full user privileges and no sandbox before any source code compiles. Any `build.rs` (project's own or in a `[build-dependencies]`) that opens a socket, downloads a tarball, or shells out to `curl` / `wget` is a P1 supply-chain finding. Legitimate `build.rs` content: `cargo:rerun-if-changed=…`, `bindgen`, `cc::Build`, protobuf compilation.

### 7.8 Hand-rolled crypto

Acceptable: `ring`, `RustCrypto` primitives (`aes-gcm`, `chacha20poly1305`, `ed25519-dalek`, `argon2`) wrapped behind a project-local crate that fixes the algorithm, key size, IV/nonce strategy, and AAD policy.

Findings:
- Direct use of low-level RustCrypto traits (`BlockCipher`, `StreamCipher`) inline in business code.
- Custom KDFs, custom MAC constructions, custom padding.
- `rand::random()` for keys or nonces (use `rand::rngs::OsRng` and `RngCore::try_fill_bytes`).
- Roll-your-own JWT / JWE / JWS.

### Adversarial scan questions

- Is there a `deny.toml` committed at the repo root? If not, cite its absence.
- Is the `rust-version` key set in `Cargo.toml`? If not, cite its absence.
- Cite any `build.rs` (or `[build-dependencies]` crate's `build.rs`) that performs network I/O or shells out.
- Cite any direct use of a low-level crypto primitive outside a wrapping crate.
- Cite any `unsafe` block introduced by this change that is not paired with a `// SAFETY:` comment justifying invariants.

---

## Severity and Category Table

| Pattern | Default tier | Category |
|---|---|---|
| `f32` / `f64` holding money | **P1** | `data_correctness` |
| Bare `+` / `-` on money type (no `checked_*`) | **P1** | `data_correctness` |
| Rounding without explicit `RoundingStrategy` | **P2** | `data_correctness` |
| Untagged `Decimal` money in signature | **P2** | `data_correctness` |
| Direct `Utc::now()` in handler | **P2** | `tests` |
| `DateTime<Local>` crossing boundary | **P2** | `data_correctness` |
| `SystemTime` subtraction for elapsed | **P2** | `data_correctness` |
| `format!`-built SQL | **P1** | `security` |
| Runtime `query` where `query!` fits | **P3** | `data_correctness` |
| Hardcoded pool `max_connections` | **P3** | `scope` |
| Read-then-write balance without `FOR UPDATE` | **P1** | `data_correctness` |
| Missing `#[serde(deny_unknown_fields)]` on external DTO | **P2** | `security` |
| Validation deep in handler (not at boundary) | **P2** | `data_correctness` |
| `#[derive(Debug)]` on PII-bearing struct | **P2** | `security` |
| Missing `#[tracing::instrument]` on handler | **P2** | `observability` |
| Format-string `info!` instead of structured fields | **P3** | `observability` |
| Unbounded label / span field cardinality | **P2** | `observability` |
| `info!` inside hot loop | **P3** | `observability` |
| Real network call in test | **P1** | `tests` |
| Real PII in fixtures | **P1** | `security` |
| Wall-clock-dependent test without `pause()` / fake clock | **P2** | `tests` |
| Missing `deny.toml` | **P2** | `supply_chain` |
| Missing `rust-version` (MSRV) pin | **P3** | `cargo_hygiene` |
| Yanked crate in `Cargo.lock` | **P2** | `supply_chain` |
| Unresolved RUSTSEC advisory | **P1** | `supply_chain` |
| Network I/O in `build.rs` | **P1** | `supply_chain` |
| Inline low-level crypto in business code | **P1** | `security` |
| Claim in PR description not backed by code | **P2** | `claims_vs_reality` |
| `unsafe` block without `// SAFETY:` comment | **P1** | `unsafe_rust` |

Severity may be raised one tier when the affected code path is on the money-movement critical path, and lowered one tier when the affected code path is internal-only and behind a feature flag.
