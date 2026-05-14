# Review Checklist (Rust)

Applied at steps 4 + 5 of `SKILL.md` — a deterministic YES/NO scan over
every file in `code_under_review` and `tests_under_review`. Every NO becomes
one finding. The check itself does NOT assign severity — `severity-guide.md`
does that at step 7.

## Convention overrides (read first)

This checklist names specific Rust crates as defaults: `axum` for HTTP,
`sqlx` for persistence, `tracing` for observability, `thiserror` / `anyhow`
for errors, `tokio` for async runtime, `chrono` for time, `wiremock` for
HTTP test doubles, `validator` for input validation, `rust_decimal` for
money. These are STAND-INS. If `convention_overrides` declares a different
crate playing the same role (e.g. `sql_layer: diesel`, `logging_crate:
slog`, `time_crate: jiff`, `http_framework: actix-web`,
`money_type: integer_minor_units`), evaluate the equivalent rule against the
declared crate. The checklist asks "is the role fulfilled," not "is the
literal crate used." When `convention_overrides` is absent, treat the named
crate as the expected default.

Deep references with the code shapes each item is asking about:
- `rust-project-structure.md` (project layout, modules, Cargo, lints)
- `rust-idioms-and-async.md` (errors, async, unsafe, ownership, traits)
- `rust-banking-patterns.md` (money, time, sqlx, tracing, serde, tests, supply chain)
- `coverage-and-unsafe-shapes.md` (cross-checking `decision_metadata.rust_specific` against the code)
- `glob-policy-idioms.md` (glob syntax and precedence for `critical_paths` / per-module `expect_invariant_policy`)
- `implement-stage-output-contract.md` (forward-compat extension keys the implement skill is expected to emit)

**`critical_paths` and severity promotion.** When the input declares
`critical_paths: [{globs, reason, promote_to?}]`, any finding in an
eligible category whose `file` matches a `globs` entry is promoted to
the declared `promote_to` (default P1) and must carry the marker
`[critical: <reason>]` in `evidence`. Cosmetic categories
(`observability`, `errors`, `scope`, `project_structure`,
`cargo_hygiene`, `lint_baseline`, `claims_vs_reality`) are NOT promoted.
See `severity-guide.md` § Critical-paths severity promotion.

**`expect_invariant_policy` per-module.** Object form:
`{default, forbid_under: [globs], allow_under: [globs]}`. Resolution is
restrictive-wins (`forbid_under` takes precedence). See
`severity-guide.md` § Documented-invariant `expect` exception § Per-module
override.

Source: the same generate-side discipline checklist used by the Rust
implement stage, re-organized so generate-side and review-side cover
identical ground from opposite ends.

## A. Data correctness (code)

- [ ] Every command handler enforces an invariant from the design's L1 (cite the line that holds the invariant).
- [ ] Every query handler is side-effect-free (no inserts / updates / deletes / event publishes — `sqlx::query!` only against `SELECT`, no `query::execute` on a mutation statement).
- [ ] Every persistence call sits inside an explicit `sqlx::Transaction` (or repo-equivalent) scope OR is a single-statement write whose atomicity the comment / commit message justifies.
- [ ] No cross-aggregate write inside a single transaction without outbox / saga / compensation.
- [ ] Concurrency-sensitive paths use `SELECT ... FOR UPDATE`, an advisory lock, or a documented compare-and-swap (cite the lock acquisition).
- [ ] No edit to a committed migration file. New migrations are additive (new `migrations/NNN_*.sql`, never an edit to a numbered existing one).
- [ ] Every monetary value uses `rust_decimal::Decimal`, an integer minor-unit type (`Money<C>` newtype over `i64`/`i128`), or `bigdecimal::BigDecimal`. No `f32` / `f64` for money. See `rust-banking-patterns.md` § 1.
- [ ] Every monetary arithmetic call uses `checked_*` / `saturating_*` (per design); no bare `+` / `-` / `*` on a money type unless wrapped by a newtype whose `Add` impl is `checked_*`.
- [ ] Every `Decimal` rounding call passes an explicit `RoundingStrategy::*` argument.
- [ ] Every monetary function signature carries a currency newtype or an explicit `currency` parameter.
- [ ] Only one time crate is used across the repo (`chrono` and `time` and `jiff` are not mixed).
- [ ] `DateTime<Local>` does not cross any persistence or network boundary.
- [ ] `Instant::elapsed()` (monotonic) is used for elapsed-time measurement; `SystemTime` / `DateTime<Utc>` subtraction is not used as a duration.

## B. Idempotency (code)

- [ ] Every external side-effect handler reads `idempotency_key` at the boundary (cite the extractor / request struct field / header read).
- [ ] A persistence call to the dedup store is visible (cite the `sqlx::query!` against the idempotency table or the repo trait method).
- [ ] A replay branch returns the stored response without re-running the side-effect (cite the match arm or `if let Some(prev) = ...` branch).
- [ ] A same-key + different-`request_hash` branch returns a 409 / domain-equivalent error (cite the branch).
- [ ] Kafka / message consumers dedup on `(topic, partition, offset)` AND on business idempotency key (e.g. `rdkafka` consumer storing offsets manually plus a business-key check).
- [ ] If `idempotency_strategy` declares "naturally idempotent," the code path is in fact read-only (no inserts / updates / deletes / publishes / external calls).

## C. Errors (code)

- [ ] Every `?`-propagated error has either a `thiserror` `#[from]` / `#[source]` link OR an `anyhow::Context::context("...")` call (cite a `?` that drops the cause).
- [ ] No `.map_err(|e| Custom::Msg(e.to_string()))` lossy chain when `#[from]` would carry the source.
- [ ] Library public API does NOT return `anyhow::Result` (use `thiserror` enum for stable, matchable variants).
- [ ] Every error variant returned from a handler has a class (`client | server | dependency`) at the type level — an enum discriminant, a `fn class(&self) -> ErrorClass`, or a `From<MyError> for HttpStatusClass` impl.
- [ ] No `unwrap()` / `expect(...)` / `panic!()` / `unreachable!()` / `todo!()` / `unimplemented!()` in a request or worker path (cite the call site). Documented-invariant `expect("invariant: ...")` exception per `severity-guide.md`.
- [ ] No `_ = result` on a `Result` whose error variant matters. No `let _ = ...?;` swallow. No `if let Ok(_) = ...` that discards the error case silently.
- [ ] No string-sniffing on errors (`err.to_string().contains("not found")`). Branch on typed variants instead.
- [ ] No SQL fragment / internal ID / stack trace / env value / PII in any error response shape (cite the leak).

## D. Cancellation / async hygiene (code)

- [ ] Every public async function in a handler / use-case path is reachable from a request with a timeout or cancellation token (cite the timeout source).
- [ ] No `tokio::spawn(...)` that outlives the parent request unless the spawn is declared by the design (background task / outbox publisher). Use `JoinSet` to collect detached results.
- [ ] No `std::thread::sleep`, blocking `std::fs::*`, blocking `std::net::*`, or other blocking call inside an `async fn` (use `tokio::fs`, `tokio::time::sleep`, `tokio::task::spawn_blocking`).
- [ ] No CPU-bound work >10 µs inside an `async fn` without `tokio::task::spawn_blocking` (or a yield point).
- [ ] No `std::sync::Mutex` / `RwLock` guard held across an `.await` point (use `tokio::sync::Mutex` only when truly needed, otherwise restructure the scope).
- [ ] No `tokio::spawn` capturing `Rc` / `RefCell` or a borrow that would compile only under the `current_thread` runtime.
- [ ] No non-cancel-safe future (`AsyncReadExt::read_exact`, `AsyncWriteExt::write_all`, `Mutex::lock`, in-flight transaction) inside a `tokio::select!` branch.
- [ ] No resource needing async cleanup hidden in `Drop` (use explicit `async fn close(self)` / `finish(self)`).
- [ ] Cancellation token (or deadline) is passed into every `sqlx`, `reqwest`, `tonic`, and message-producer call on the path.

## E. Observability (code)

- [ ] Every declared failure mode in the design has a counter increment in the code (`metrics::counter!`, `prometheus::Counter`, repo equivalent).
- [ ] Every handler is wrapped by a `tracing` span (`#[tracing::instrument(skip(...), fields(...))]` attribute OR an explicit `span!(...).entered()`); span name follows the repo's convention.
- [ ] All `tracing::*!` calls use structured fields (`field = %value` for Display, `field = ?value` for Debug) not format-string interpolation (`"…{}…"`).
- [ ] Every log macro call includes `trace_id` (typically inherited from the span; cite if the span is missing).
- [ ] Span fields and `metrics::*!` labels are bounded — no user IDs, account IDs, request IDs, free-form strings, or PII as labels.
- [ ] No `info!` / `warn!` inside a hot loop body.
- [ ] `traceparent` header is extracted at the HTTP boundary and made the parent of the root handler span.

## F. Security (code)

- [ ] AuthN check is present at the transport boundary (cite the middleware / extractor / tower layer).
- [ ] AuthZ check is against the use case's required permission, not the transport route.
- [ ] Every input validated using the repo's validator (`validator` crate, `serde` `try_from`, hand-rolled `TryFrom<RawDto, Error = ValidationError>`) before reaching the use case.
- [ ] Every externally-deserialized DTO carries `#[serde(deny_unknown_fields)]`.
- [ ] `rename_all` is consistent across DTOs in the same API surface.
- [ ] Validation happens at the deserialization boundary (`TryFrom<RawDto>` or `validator::Validate`), not deep in the handler.
- [ ] All SQL parameterized — `sqlx::query!` / `query_as!` macros, or `query(...).bind(...)`. No `format!`-built SQL. No `+` concatenation into SQL.
- [ ] No hardcoded secrets / tokens / connection strings / private keys.
- [ ] No `reqwest` / `hyper` / `tonic` client calls to external hosts inside tests.
- [ ] No hand-rolled crypto where a repo-approved helper exists. No direct use of `ring` / `RustCrypto` low-level primitives (`BlockCipher`, `StreamCipher`, raw `Mac`) where the repo's auth crate wraps them.
- [ ] No `rand::random()` for keys / nonces (use `OsRng` + `RngCore::try_fill_bytes`).
- [ ] PII fields named in the design are masked in logs / error responses (custom `Debug` / `Display` impl, or `#[serde(skip)]` / `#[serde(skip_serializing)]` on the response shape).
- [ ] No `#[derive(Debug)]` on a struct that carries a PAN, email, full name, SSN, account number, or other named-PII field.
- [ ] No `unsafe { ... }` block or `unsafe fn` without a `// SAFETY:` comment naming the invariants AND a declaration in `implement_stage_output.decision_metadata`.
- [ ] Every `unsafe fn` carries a `# Safety` rustdoc section.

## G. Audit (code)

- [ ] Every state-changing path emits exactly one audit event (cite the emit — typically a method on an `AuditPublisher` trait or an outbox table insert).
- [ ] Every emitted `event_type` appears in `implement_stage_output.audit_events_emitted`.
- [ ] Audit payload includes actor, action, target, timestamp, trace_id, decision_metadata.

## H. Compensation (code)

- [ ] Every irreversible external side-effect call has a corresponding entry in `implement_stage_output.compensating_actions`.
- [ ] Each `compensating_actions[].action_skill_ref` is plausible (named in a way the workflow can resolve).
- [ ] If compensation is impossible, the workflow stage is declared `human-queue` (read this from the design / workflow definition, not the code).

## I. Tests (tests_under_review)

- [ ] A test file or `#[cfg(test)] mod tests` block exists for every production file in `code_under_review`.
- [ ] Per-file coverage claim `>= test_coverage_target` from input is plausible (the test exercises the branches required to reach the number — success path, each error variant, replay path, validation rejection).
- [ ] Tests are parameterized where appropriate (`rstest` cases, a `for case in vec![...]` loop, or `#[test_case]`) for handlers, services, and repo unit tests.
- [ ] No test depends on wall-clock time, network, or random ordering. `tokio::time::pause()` / mock clocks used instead of real sleeps.
- [ ] No test calls `reqwest` / `hyper` / `tonic` against a real host — `wiremock::MockServer` or an in-process fake repo only.
- [ ] No test depends on `Utc::now()` / `Instant::now()` directly without an injected `Clock` trait.
- [ ] No `#[ignore]` without a named unblock condition. No `t.Skip`-style runtime skip without justification.
- [ ] No real PII in fixtures. Synthetic data only — no Luhn-valid PAN, no valid-TLD email matching real customers, no real-shape SSN.
- [ ] No secret read from environment in the test (no `std::env::var("PROD_DB_URL")`). Test config comes from fixtures.
- [ ] `#[tokio::test]` flavor (`current_thread` vs `multi_thread`) matches the production runtime when the test exercises runtime-specific behavior.

## J. Scope discipline

- [ ] No crate `use`d that is not in the relevant `Cargo.toml` (cross-reference each `use` statement). Workspace-inherited deps are fine if declared at workspace level.
- [ ] No edit to a generated artifact (`OUT_DIR`-emitted modules, `tonic-build` / `prost` outputs, vendored bindings, `src/generated/` files).
- [ ] No reformatting / cleanup outside the changed files.
- [ ] Public contract (exported `pub fn`, `pub struct`, `pub enum`, `pub trait`) has not widened beyond what the design declares.

## K. Claims-vs-reality (used at step 6)

- [ ] Every `audit_events_emitted` entry has a matching emit call in the code.
- [ ] Every `compensating_actions[].trigger` has a matching call site in the code.
- [ ] `idempotency_strategy` text matches the code's actual behavior (not aspirational text).
- [ ] Every `decision_metadata.pattern_choices` entry is consistent with the code.
- [ ] Every `uncertainty_flag` raised by Generate is triaged: resolved-by-code / still-open / needs-escalation.
- [ ] When present, every `decision_metadata.rust_specific.unsafe_declarations[]` entry points at a real `unsafe` block in `code_under_review` at the cited `file_path:line_start`. Conversely, every `unsafe` block in `code_under_review` appears in this list. See `coverage-and-unsafe-shapes.md` § Unsafe tracking and `implement-stage-output-contract.md`.
- [ ] When present, `tests_generated[].coverage_pct` claims are within ±2 pp of the sum computed from `decision_metadata.rust_specific.per_function_coverage[]` (when both are emitted by the implement stage).
- [ ] When present, `decision_metadata.rust_specific.clippy_lint_floor` matches the actual `[lints]` table / `#![deny(...)]` in `code_under_review`.
- [ ] `audit_metadata.rule_walk` MUST contain exactly one entry for each rule in {B1..B11, A1..A7, R1..R12, C1..C4} (34 entries total). Each `not_applicable` outcome carries a non-empty `rationale`. Each `finding` outcome corresponds to at least one entry in `findings[].rule_violated`.

## L. Project structure (Rust)

- [ ] No edits to a generated artifact under `src/generated/` or any file embedded via `include!(concat!(env!("OUT_DIR"), ...))`.
- [ ] No member crate redeclares `[profile.*]` or `[patch]` (silently ignored at workspace).
- [ ] `[[bin]]` does not reuse the `[lib]` name; target ambiguity avoided.
- [ ] Any `optional = true` dependency is gated by a feature.
- [ ] Features are additive (enabling one does not remove an API).
- [ ] `pub` widening of a previously `pub(crate)` item is justified in the design (counts as API surface change).
- [ ] `#[non_exhaustive]` is preserved on public enums and structs the design declares as forward-compat.
- [ ] No glob `pub use foo::*;` that re-exports items not enumerated in the design.
- [ ] No `#[path = "..."]` loading code from outside the crate's `src/` directory.

## M. Cargo & toolchain hygiene (Rust)

- [ ] `Cargo.toml` has `rust-version` (MSRV) pinned.
- [ ] `rust-toolchain.toml` channel is a specific version (e.g. `1.82.0`), not a moving label (`stable` / `beta` / `nightly`).
- [ ] `Cargo.lock` edits in the diff are explained by a `Cargo.toml` change or a documented `cargo update` invocation.
- [ ] No yanked crate in `Cargo.lock` (cross-check `cargo audit`).
- [ ] `[lints]` table or `#![deny(...)]` / `#![forbid(...)]` is consistent with the repo's documented lint floor.
- [ ] `unsafe_code = "forbid"` is not weakened to `deny` or `allow` in this diff.
- [ ] No blanket `#[allow(clippy::all)]` / `#[allow(warnings)]` introduced near banking-grade logic.
- [ ] No `clippy.toml` edits that raise complexity thresholds rather than refactor.

## N. Supply chain (Rust)

- [ ] A `deny.toml` is committed at the repo root for `cargo-deny`.
- [ ] No `build.rs` (project's or a `[build-dependencies]` crate's) performs network I/O, shells out to `curl` / `wget`, or writes outside `OUT_DIR`.
- [ ] `build.rs` carries appropriate `cargo::rerun-if-changed=...` / `rerun-if-env-changed=...` directives.
- [ ] No newly-added dependency is a typo-squat or low-reputation publisher; proc-macro additions especially scrutinized.
- [ ] No direct or transitive dep listed in RUSTSEC with status `vulnerability` or `unmaintained`.

## Routing of NO answers

A NO becomes a finding. Severity is assigned at step 7 using
`severity-guide.md`. As a quick mental map (the matrix is authoritative):

| Section | Default severity |
|---------|------------------|
| A (data — including money/time), B (idempotency), F (security — including serde / crypto / PII / unsafe), G (audit), H (compensation), N (supply chain — RUSTSEC, build.rs, crypto) | `P1` |
| C (errors), D (cancellation / async / mutex-across-await), E (observability), I (tests), K (claims-vs-reality), L (project structure — non-cosmetic), M (cargo hygiene — non-cosmetic) | `P2` |
| J (scope discipline — cosmetic), L (style-only structural items), M (missing MSRV, lint-baseline weakening on non-sensitive logic), N (`deny.toml` absence when other gates exist) | `P3` |
