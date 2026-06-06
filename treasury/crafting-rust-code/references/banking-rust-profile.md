# Banking Rust Profile (optional)

Apply this profile **only** when the design handles money, ledgers, payments, balances, interest, or regulated data (PDPA/PII, AMLA-relevant flows). It layers banking-grade defaults onto the general-purpose core so generated Rust passes `review-rust-code`'s financial rules (R7 money, R8 time, A1 audit, A2 idempotency, A3 compensation, R12 supply chain) on the first review loop. For non-financial crates, skip this file — those rules degrade to not-applicable.

## 1. Money arithmetic (R7)

- **Never `f32`/`f64` for money.** Use `rust_decimal::Decimal` (fixed-point, ~28 significant digits, fast) or integer minor units (`i64` cents) with a documented scale. Floats lose pennies and are non-deterministic across platforms.
- **Checked arithmetic only**: `amount.checked_add(other)` / `checked_sub` / `checked_mul`, returning a domain error on overflow — never the bare `+`/`-`/`*` (panic in debug, wrap in release).
- **Explicit rounding**: every `Decimal::round_dp` / `round` call names a `RoundingStrategy` (e.g. `MidpointNearestEven` for banker's rounding). Never rely on a default.
- **Currency-tagged signatures**: wrap money in a newtype that carries currency (`struct Money { amount: Decimal, currency: Currency }`) or pass an explicit `currency: Currency` param. A bare `Decimal`/`i64` in a monetary signature is a defect — it lets you add THB to USD.

```rust
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub struct Money { amount: Decimal, currency: Currency }

impl Money {
    pub fn checked_add(self, rhs: Money) -> Result<Money, MoneyError> {
        if self.currency != rhs.currency { return Err(MoneyError::CurrencyMismatch); }
        let amount = self.amount.checked_add(rhs.amount).ok_or(MoneyError::Overflow)?;
        Ok(Money { amount, currency: self.currency })
    }
}
```

## 2. Time discipline (R8)

- **Inject a `Clock` trait** — never call `Utc::now()` / `Local::now()` / `SystemTime::now()` / `Instant::now()` directly in a handler or domain method. Injection makes time testable and auditable.
- **UTC at boundaries**: persist and transmit `DateTime<Utc>`; `DateTime<Local>` must not cross a persistence or network boundary.
- **Monotonic for durations**: measure elapsed time with `Instant::elapsed()`, not `SystemTime` subtraction (wall clock can jump).
- **One time crate**: pick `chrono` **or** `time` **or** `jiff` — mixing them in one repo is a defect.
- Tests freeze time via the injected clock or `tokio::time::pause()`.

## 3. Canonical audit events (A1)

- Every state-changing path emits **exactly one** audit event with the canonical shape: `event_type`, `actor`, `action`, `target`, `timestamp`, `trace_id`, `decision_metadata`.
- Emit via the repo's audit publisher / outbox insert / dedicated `tracing` audit target — not an ad-hoc log line.
- Keep an enumerated list of emitted `event_type`s so it can be cross-checked against the design (and, in a pipeline, against `implement_stage_output.audit_events_emitted`).

## 4. Idempotency key + replay (A2)

- Read a **UUID-v4** idempotency key at the boundary (request header or message envelope) for every external side-effect.
- Persist `(key, request_hash, response_payload, status)`.
- **Replay**: same key + same `request_hash` returns the stored response without re-running the side-effect. Same key + **different** `request_hash` returns `409 Conflict` (or the domain equivalent). Read-only paths are naturally idempotent and need no key.

## 5. Compensating actions (A3)

- For every **irreversible** external side-effect (fund transfer, third-party charge, message publish), declare the compensating action up front: what reverses it and under what trigger.
- Pair the forward action with a named compensation (e.g. `transfer` -> `transfer_reversal`) so a downstream failure can unwind cleanly. Saga/outbox coordinates the sequence; do not bury an irreversible call inside an uncompensated path.

## 6. Supply-chain hygiene (R12)

- Commit a **`deny.toml`** at the repo root and run `cargo-deny` (licenses, advisories, bans, sources) in CI; commit a `cargo-audit` / RUSTSEC check. Zero unresolved `vulnerability`/`yanked`/`unmaintained` advisories of high/critical severity.
- Pin MSRV (`package.rust-version`).
- **No network I/O or shelling out in `build.rs`** (project's or any `[build-dependencies]` crate's) and no writes outside `OUT_DIR`.
- Crypto only through a vetted, project-local wrapper — no direct `BlockCipher`/`StreamCipher`/raw `Mac`/low-level primitives in business code, no custom cryptography.
- **Randomness for keys/nonces**: `OsRng` + `RngCore::try_fill_bytes`, never `rand::random()`.

## Degradation note

When `cargo-deny` / `cargo-audit` / a live DB / network is unavailable in the sandbox, report the exact status (`not_run: cargo-deny unavailable`, `unavailable: no DB for sqlx prepare`) per `editing-guardrails.md` — never claim a banking control passed that you could not verify.
