# Review Rubric (Rust)

The canonical scan list applied at step 4 of `SKILL.md`. Each item is the
same banking-grade rule used by the upstream Rust implement stage — but
re-cast as an *adversarial question* the reviewer asks of the code. The
mapping is 1:1: every rule a Generate stage must follow, the Review stage
must verify.

The framing differs deliberately:
- Generate's self-review checklist: "Did I do X?" (first-person, self-discipline)
- This file: "Show me where X is. Cite the file:line." (third-person, adversarial)

When a question's answer is "no" or "not visible in the code," emit a
finding. Severity is assigned at step 7 using `severity-guide.md`.

## Convention overrides

Crate names appearing in any question below (`axum`, `sqlx`, `tracing`,
`thiserror`, `anyhow`, `tokio`, `chrono`, `wiremock`, `validator`,
`rust_decimal`, etc.) are DEFAULTS. If `convention_overrides` declares a
different crate playing the same role (e.g. `sql_layer: diesel`,
`logging_crate: slog`, `time_crate: jiff`, `http_framework: actix-web`),
the rule applies to the declared crate. The rule asks "is the role
fulfilled," not "is the literal crate used." When `convention_overrides`
is absent, treat the named crate as the expected default.

## Deep references for Rust-specific code shape comparison

| Need | File |
|------|------|
| Project structure, module system, Cargo.toml schema, lint configuration, generated-code rules, RUSTSEC categories | `rust-project-structure.md` |
| Error handling, async/tokio pitfalls, unsafe discipline, ownership traps, trait/generic discipline | `rust-idioms-and-async.md` |
| Money/decimal, time, `sqlx` persistence, `tracing`/`metrics` observability, `serde` validation, test discipline, supply chain | `rust-banking-patterns.md` |
| Rust coverage shapes (`cargo-llvm-cov`, `cargo-tarpaulin`, LCOV, `grcov`); unsafe-tracking (`cargo-geiger`, `cargo-vet`, MIRI, clippy lints) — basis for `decision_metadata.rust_specific` cross-checks | `coverage-and-unsafe-shapes.md` |
| Glob syntax (`globset` / gitignore semantics), precedence model (restrictive-wins), severity-zone naming — basis for `critical_paths` and per-module `expect_invariant_policy` | `glob-policy-idioms.md` |
| `implement_stage_output` shape, forward-compat extension keys, co-evolution checklist for the Rust implement skill | `implement-stage-output-contract.md` |

These references contain the deep code-shape patterns each adversarial
question is asking about. Consult them when a scan question is ambiguous
on what to look for.

## Severity promotion via `critical_paths`

When the input declares `critical_paths: [{globs, reason, promote_to?}]`,
the eligible category list in `severity-guide.md` § Critical-paths
severity promotion governs which findings get promoted. The Rust-specific
rules most commonly affected: R1 (unsafe), R3 (blocking-in-async), R7
(money), R8 (time), R9 (persistence), R12 (supply chain), plus all
banking-grade B/A rules. Findings in promoted categories MUST carry the
marker `[critical: <reason>]` in `evidence`. Cosmetic categories
(`observability`, `errors`, `scope`, `project_structure`,
`cargo_hygiene`, `lint_baseline`, `claims_vs_reality`) are ineligible.

## Base rule questions (11)

Source for the rules: the banking-grade decision-rules set carried over
from `treasury/crafting-backend-code/references/decision-rules.md`, with
Rust idioms substituted for Go idioms.

| # | Rule | Adversarial scan question |
|---|------|---------------------------|
| B1 | Pattern first for templates | Has the code copied an example crate's module path, binary shape, or business capability as if it were mandatory? Cite the over-borrowed identifier. |
| B2 | Repo first for existing services | Does the new code mirror the existing crate's error type, `tracing` subscriber wiring, validator, DI shape (typically constructor injection through traits), and test style? Cite a sibling file that the code diverges from. |
| B3 | Contracts before code | Does every public function / handler / message exposed by the new code match the contract declared in the design (request / response struct shape, versioning, error envelope, auth)? Cite the design line and the code line. |
| B4 | One owner per piece of state | Does the new code mutate a table / aggregate / cache key / event stream that another crate already owns? Cite the cross-crate write. |
| B5 | Transactions are explicit | For every multi-statement write path: is a `sqlx::Transaction` (or repo-equivalent `begin()` / `commit()`) visible? For cross-aggregate work: is outbox / saga / compensation present? Cite the missing boundary. |
| B6 | Idempotency is required for retries | For every command handler / webhook / consumer / outbox dispatcher: where is the dedup key persisted and checked? Cite the missing dedup. |
| B7 | Cancellation propagates | Does every public async function accept either an explicit `CancellationToken` (or per-repo equivalent) OR honour caller cancellation by not spawning detached `tokio::spawn` tasks that outlive the request? Does the code propagate the timeout / deadline into every `sqlx`, `reqwest`, `tonic`, and message-producer call? Cite a path that drops cancellation. |
| B8 | Errors are operational signals | Is every `?`-propagated error wrapped with a typed cause (`#[from]` / `#[source]` on a `thiserror` enum, or `anyhow::Context::context(...)`)? Is every public error response free of SQL fragments, internal IDs, stack traces, env values, PII? Cite the leak. |
| B9 | Observability follows failure modes | For every declared failure mode in the design: is there a metric increment (`metrics::counter!`, Prometheus handle, etc.), a `tracing::error!` / `warn!` log line, and a span? Cite the missing instrumentation. |
| B10 | Security is not a cleanup task | Are authN / authZ / input validation / SSRF guards / parameterized `sqlx::query!` (NOT `query(format!(...))`) / secret handling via the repo's secrets crate / PII masking / approved crypto all present at the boundary they belong? Cite the missing control. |
| B11 | Generated and migrated artifacts are special | Has the code edited a generated file (`OUT_DIR` files, `tonic-build` / `prost` generated modules, committed migration SQL files), or a `Cargo.lock` line under `[[package]]` not driven by `cargo update`? Cite the edit. |

## v2 augmentation questions (7)

Source for the augmentations: the v2 implementation-rules set carried over
from the implement stage and re-cast for Rust.

| # | Augmentation | Adversarial scan question |
|---|--------------|---------------------------|
| A1 | Canonical audit event shape | For every state-changing path: is exactly one audit event emitted with `event_type / actor / action / target / timestamp / trace_id / decision_metadata`? Does the event_type appear in `implement_stage_output.audit_events_emitted`? Cite the missing emit or the mismatch. |
| A2 | Idempotency key format + replay behavior | For every external side-effect path: is the UUID-v4 key read at the boundary (extractor / message header)? Is `(key, request_hash, response_payload, status)` persisted? Does a same-key replay return the stored response without re-running the side-effect? Does a same-key + different-`request_hash` return a 409 / domain-equivalent? Cite the missing piece. |
| A3 | Compensating-action discipline | For every irreversible external side-effect in the code: is the corresponding `compensating_actions[].trigger` declared by Generate and is its referenced action_skill_ref plausible? Cite an irreversible call site that has no declared compensation. |
| A4 | Error classification (`client \| server \| dependency`) | Does every `Error` variant returned from a handler carry a class at the type level (an enum discriminant or a trait method like `fn class(&self) -> ErrorClass`), not inferred at the edge? Does the class drive the HTTP status, retry decision, and `tracing` error attribute correctly? Cite the unclassified error. |
| A5 | Test fixtures discipline | Do any test files call `reqwest`, `hyper`, `tonic` clients, etc. against a real host, read secrets from environment, mutate shared state between cases (static `Mutex`, lazily-initialized global database), `#[ignore]` without a named unblock condition, or rely on `tokio::time::sleep` for ordering? Cite the offender. |
| A6 | Convention discovery overrides templates | If the code follows a template default but the target crate's existing code uses a different convention (e.g. custom `Error` enum vs. `anyhow::Error`, or `slog` vs. `tracing`): did Generate emit a `convention_conflict` uncertainty flag? Cite the divergence the flag should have named. |
| A7 | No silent dependency additions | Does the code `use` a crate not declared in `Cargo.toml` (workspace-level or local)? Did Generate emit a `dependency_addition` uncertainty flag for any added `[dependencies]` line? Cite the import / dependency entry. |

## Rust-specific adversarial questions (12)

These augment the base rules with Rust-language-specific code shapes.
Severity assignment is in `severity-guide.md`. The "Deep ref" column
points to the file with examples and code shapes.

| # | Item | Adversarial scan question | Deep ref |
|---|------|---------------------------|----------|
| R1 | Unsafe Rust | Does any file contain an `unsafe { ... }` block or `unsafe fn`? For each: (a) is there a `// SAFETY:` comment naming the invariants (not "this is safe")? (b) does each `unsafe fn` carry a `# Safety` rustdoc section? (c) are `unsafe impl Send` / `unsafe impl Sync` blocks justified for every field including future ones? (d) is `mem::transmute`, `Pin::new_unchecked`, `Box::leak`, or `std::hint::unreachable_unchecked` used without a written rationale? (e) should this crate be `#![forbid(unsafe_code)]` but is not? | `rust-idioms-and-async.md` § 3 |
| R2 | Panicking APIs on production paths | Does any function reachable from the handler / worker / request path call `.unwrap()`, `.expect(...)`, `panic!()`, `unreachable!()`, `todo!()`, or `unimplemented!()`? Cite the call site. Documented-invariant `expect("invariant: <reason>")` may be downgraded per `severity-guide.md` § Documented-invariant `expect` exception. | `rust-idioms-and-async.md` § 1 |
| R3 | Blocking-in-async + mutex-across-await + send/static | In any `async fn` or `.await`-reachable path: (a) does the code call `std::thread::sleep`, blocking `std::fs::*`, blocking `std::net::*`, or run >10 µs of CPU work without `tokio::task::spawn_blocking`? (b) is any `std::sync::Mutex` / `RwLock` guard held across `.await` (producing a `!Send` future)? (c) does any `tokio::spawn` capture `Rc` / `RefCell` or a borrow that compiles only under the `current_thread` runtime? (d) is `tokio::spawn` used in detach mode on the request path where `JoinSet` would propagate failures? (e) does any `tokio::select!` branch contain a non-cancel-safe future (`AsyncReadExt::read_exact`, `AsyncWriteExt::write_all`, `Mutex::lock`, in-flight transaction)? (f) does any I/O-issuing function lack a deadline or `CancellationToken`? | `rust-idioms-and-async.md` § 2 |
| R4 | Lint baseline | Does the crate or workspace declare a `[lints]` table (Rust 1.74+) or crate-level `#![deny(...)]` / `#![forbid(...)]` consistent with the design or repo convention? Is `unsafe_code = "forbid"` weakened to `deny` or `allow` in this diff? Are blanket `#[allow(clippy::all)]` / `#[allow(warnings)]` introduced near sensitive logic? Does `clippy.toml` raise complexity thresholds rather than refactor? | `rust-project-structure.md` § Lint Configuration |
| R5 | Project structure & module visibility | Does the diff widen visibility (`pub(crate)` → `pub`) without a SemVer-relevant signal? Is `#[non_exhaustive]` removed from a public enum or struct? Does a glob `pub use foo::*;` re-export items not enumerated in the design? Is `name.rs` present alongside `name/mod.rs` (half-rename)? Is `#[path = "..."]` used to load code from outside the crate's `src/` (privacy escape hatch / vendoring without provenance)? Are integration tests in `tests/*.rs` using `#[path]` to import private items? Has a sealed-trait `Sealed` supertrait become accidentally public? | `rust-project-structure.md` § Module System, § Visibility and API Surface |
| R6 | Cargo manifest & toolchain hygiene | Does any member redeclare `[profile.*]` or `[patch]` (silently ignored at workspace)? Is `resolver` absent from a virtual workspace (defaults to "1")? Does `[[bin]]` reuse the `[lib]` name (target ambiguity)? Is an `optional = true` dependency present without a feature gating its `use` statement? Are features non-additive (mutually exclusive without `compile_error!`)? Is `Cargo.lock` edited without a `Cargo.toml` change or `cargo update` invocation? Is `rust-toolchain.toml` `channel` a moving label (`stable` / `beta` / `nightly`) instead of a pinned version? Is the pinned MSRV in `package.rust-version` absent or lower than the toolchain channel? | `rust-project-structure.md` § Cargo.toml Schema, § Cargo.lock Discipline, § Toolchain Pinning |
| R7 | Money arithmetic discipline | Does any monetary value hold `f32` / `f64`? Does any monetary `+` / `-` / `*` use the bare operator instead of `checked_*` (panic-in-debug, wrap-in-release)? Does a `rust_decimal::Decimal` rounding call (`round_dp`, `round`) omit `RoundingStrategy::*`? Does any monetary function signature accept a bare `Decimal` or `i64` without a currency newtype or explicit `currency: Currency` parameter? | `rust-banking-patterns.md` § 1 |
| R8 | Time discipline | Does any handler call `Utc::now()` / `Local::now()` / `SystemTime::now()` / `Instant::now()` directly instead of through an injected `Clock` trait? Is any `DateTime<Local>` used across a persistence or network boundary? Is `SystemTime` subtraction used to measure elapsed time (instead of `Instant::elapsed()`)? Are two time crates (`chrono` and `time` and `jiff`) mixed in the same repo? | `rust-banking-patterns.md` § 2 |
| R9 | Persistence (Rust-flavor) | Does any SQL string come from `format!` or `+` concatenation instead of `sqlx::query!` / `query_as!` / `query(...).bind(...)`? Is a runtime `sqlx::query(...)` used where a compile-time-checked `query!` macro would compile equivalently? Is the `.sqlx/` offline cache absent from version control when `query!` is used? Is `PgPoolOptions::max_connections(<literal>)` hardcoded instead of read from config? Does any read-then-write of a money column skip `FOR UPDATE` or serializable-isolation discipline? Does a function that calls `pool.begin()` ever fail to call `tx.commit()` on the success path (silent rollback on success)? Has any committed migration file been edited (instead of a new additive migration)? | `rust-banking-patterns.md` § 3 |
| R10 | Observability (Rust-flavor) | Does every handler carry `#[tracing::instrument(skip(...), fields(...))]`? Are log macros using structured fields (`field = %value` / `field = ?value`) rather than format strings (`"...{}..."`)? Are any unbounded values (user IDs, request IDs, account IDs, IPs, emails, free-form strings) used as `metrics` labels or recurring span fields? Does `info!` / `warn!` appear inside a tight loop? Is the `traceparent` header extracted at the HTTP boundary and made the parent of the root handler span? | `rust-banking-patterns.md` § 4 |
| R11 | Input validation & serde discipline | Does every externally-deserialized DTO carry `#[serde(deny_unknown_fields)]`? Is `rename_all` consistent across DTOs in the same API surface? Is validation performed at the deserialization boundary (`TryFrom<RawDto>` or `validator::Validate`) rather than deep in the handler? Does any PII-bearing struct use `#[derive(Debug)]` or the default `Display` impl without redaction? Does any PII field appear in a response body without `#[serde(skip_serializing)]`? | `rust-banking-patterns.md` § 5 |
| R12 | Supply chain hygiene | Is a `deny.toml` committed at the repo root for `cargo-deny`? Is `rust-version` (MSRV) pinned in `Cargo.toml`? Does any `build.rs` (project's or a `[build-dependencies]` crate's) perform network I/O, shell out to `curl` / `wget`, or write outside `OUT_DIR`? Are any dependencies listed in RUSTSEC with status `vulnerability`, `yanked`, or `unmaintained`? Does this PR add a proc-macro dep from an unfamiliar publisher? Is any low-level crypto primitive (`BlockCipher`, `StreamCipher`, raw `Mac`, `ring` internals) used directly in business code instead of through a project-local wrapper? Is `rand::random()` used for keys / nonces (instead of `OsRng` + `RngCore::try_fill_bytes`)? | `rust-banking-patterns.md` § 7 |

## Workflow-contract questions (4)

In addition to the rules, the Review stage verifies the Generate stage's
contract-level claims. These do not map to a base rule — they exist
because the workflow exists.

| # | Contract item | Adversarial scan question |
|---|---------------|---------------------------|
| C1 | Companion tests exist | Is there a test (in-module `#[cfg(test)] mod tests`, integration `tests/*.rs`, or `#[tokio::test]` async test) for every production file in `code_under_review`? Cite the unaccompanied production file. |
| C2 | Per-file coverage claim plausibility | For each `tests_generated[].coverage_pct` claim: does the test file appear to exercise the branches required to reach that number (success path, error variants, replay path, validation rejection)? Cite a coverage claim with an obviously thin test. |
| C3 | `decision_metadata.pattern_choices` consistency | For each pattern choice declared by Generate (e.g. "use `thiserror` not `anyhow` at the boundary", "no audit on read-only path"): is the code consistent with the choice? Cite the inconsistency. |
| C4 | `uncertainty_flags` triage | For each flag Generate raised: is it still unresolved, was it resolved by the code in a way the design accepts, or does it need escalation? Output a triage line per flag. |

## Out-of-scope by design

This skill does NOT cover:

- Full OWASP Top 10 walk per artifact — that is
  `validating-banking-implementation` or `reviewing-software-security`
  scope. The Review stage spot-checks security via B10 + A4 + R11 + R12
  only.
- Chaos test planning — separate stage.
- Performance / load testing, micro-benchmark review (`criterion` harness
  shape) — separate stage.
- Memory-leak / borrow-checker findings the compiler catches — only
  `unsafe` discipline (R1) is reviewed.
- Architectural / boundary critique — that was the design stage's job.
  The Review stage takes the design as given.

If a reviewer spots an issue outside scope, surface it as an
`uncertainty_flag` of kind `other`, not as a finding.
