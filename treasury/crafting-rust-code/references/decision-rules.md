# Rust Decision Rules (generation side)

These are the positive, author-time rules for producing idiomatic, safe, review-ready Rust. They are the constructive mirror of the adversarial scan `review-rust-code` runs: each rule below is phrased "do this," and the reviewer asks "show me where you did." Satisfy them at author time so the review loop converges in one pass.

Mapping column references `review-rust-code` rule IDs (`treasury/review-rust-code/references/review-rubric.md`). Do not re-derive that rubric here — these are generation directives, not a duplicate scan list.

## Architecture & contracts

- **Pattern first for templates** (B1): extract reusable patterns from example crates without copying their module paths, binary shape, or business capabilities as mandatory rules.
- **Repo first for existing crates** (B2): match the existing crate's error type, `tracing` wiring, validator, DI shape (constructor injection through traits), and test style before importing a preferred pattern.
- **Contracts before code** (B3): decide public `fn`/`trait`/`struct` shape, request/response types, versioning, error envelope, and authorization before implementing. Public API is a SemVer contract.
- **One owner per piece of state** (B4): one crate/module writes a given table, aggregate, cache key, or event stream. Others read or subscribe.

## Errors & control flow

- **Errors are operational signals** (B8, A4): model domain/library errors as a `thiserror` enum; preserve cause with `#[from]`/`#[source]`; reserve `anyhow` for binaries/glue. Give each variant a class (client / server / dependency) at the type level that drives HTTP status, retry, and the `tracing` attribute. Never leak SQL, internal IDs, stack traces, env values, or PII in a public error.
- **No panics on production paths** (R2): no `.unwrap()` / `.expect()` / `panic!` / `unreachable!` / `todo!` / `unimplemented!` on any path reachable from a handler/worker/request. The only allowance is `expect("invariant: <reason>")` for a genuinely-unreachable documented invariant.
- **Unsafe is exceptional** (R1): prefer `#![forbid(unsafe_code)]`. Any `unsafe` block needs a `// SAFETY:` comment naming the invariants; any `unsafe fn` needs a `# Safety` rustdoc section; `unsafe impl Send/Sync` is justified per field. In banking code, absence of `unsafe` is the norm.

## Concurrency & async

- **Cancellation propagates** (B7, R3): every public `async fn` honours caller cancellation — no detached `tokio::spawn` outliving the request (use `JoinSet`), no `std::sync::Mutex` guard held across `.await`, no blocking `std::{thread,fs,net}` in async (use `spawn_blocking` or the async equivalent), no non-cancel-safe future in a `tokio::select!` branch. Propagate the deadline into every `sqlx` / `reqwest` / `tonic` call.
- **Ownership over locking**: prefer message-passing (manager task owns state, callers send requests over a channel) to `Arc<Mutex<_>>` shared mutable state when state crosses tasks. See `rust-idioms-and-async.md` § 2.

## Persistence, security, observability

- **Transactions are explicit** (B5, R9): wrap multi-statement writes in a visible `sqlx::Transaction`; `commit()` on the success path (never rely on rollback-on-drop for success). Cross-aggregate work uses outbox/saga/compensation.
- **Idempotency for retries** (B6, A2): every command handler / webhook / consumer / outbox dispatcher reads a dedup key at the boundary, persists `(key, request_hash, response, status)`, replays the stored response on same-key, and returns 409 on same-key + different hash.
- **Security is not a cleanup task** (B10, R11): authN/authZ, boundary validation (`TryFrom<RawDto>` / `validator`), parameterized `sqlx::query!` (never `format!`-built SQL), secrets via the repo's secret path, PII masking, approved crypto only.
- **Observability follows failure modes** (B9, R10): `#[tracing::instrument(skip(...), fields(...))]` on handlers; structured fields not format strings; a metric + log + span per declared failure mode; bounded label cardinality (no user IDs / IPs / emails as labels).

## Hygiene & contract

- **Generated and migrated artifacts are special** (B11): do not hand-edit `OUT_DIR` / `prost` / `tonic-build` output, committed migrations, or `Cargo.lock` `[[package]]` lines outside `cargo update`.
- **No silent dependency additions** (A7): every `use` resolves to a declared `Cargo.toml` dependency; flag and justify any added `[dependencies]` line.
- **Lint floor holds** (R4, R6): keep the crate/workspace `[lints]` table and `#![deny(...)]`/`#![forbid(...)]` intact; do not weaken `unsafe_code = "forbid"` or add blanket `#[allow(...)]` near sensitive logic; keep `resolver`, MSRV (`rust-version`), and `rust-toolchain.toml` pinning consistent.
- **API surface stays narrow** (R5): do not widen `pub(crate)` to `pub`, remove `#[non_exhaustive]`, or glob-re-export undesigned items without a SemVer signal.
- **Companion tests exist** (C1, C2, C3, C4): every production file ships a test; coverage claims are plausible from the branches exercised; the code is consistent with any declared pattern choice; any open uncertainty is surfaced, not hidden.

## Banking profile (apply only for financial / regulated designs)

R7 (money), R8 (time), A1 (audit events), A3 (compensation), R12 (supply chain) are set in `banking-rust-profile.md`. For non-financial crates they degrade to not-applicable.
