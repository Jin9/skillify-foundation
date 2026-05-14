# Severity Guide & Verdict Matrix (Rust)

Applied at step 7 of `SKILL.md` to classify findings and compute the verdict.
Banking-grade flavored: every routing decision is deterministic, and the
verdict is a function of the findings, not the other way around.

Source: the same banking-grade severity discipline applied by the
language-neutral review stage, with Rust-specific categories added (unsafe
code, panicking APIs in production paths, async-blocking violations,
project structure, Cargo hygiene, supply chain, lint baseline).

## Severity definitions

### `P1` — block

Banking-grade non-negotiables. A `P1` finding ALWAYS routes to
`human-queue` — never `loop_back`, never `approve`. The Generate stage
cannot fix `P1` issues in another iteration without a person looking.

Categories that are inherently `P1`:

- Data correctness — missing transaction boundary on a state-mutating flow,
  cross-aggregate write without outbox / saga, lost update on a
  concurrency-sensitive path.
- Idempotency — missing key on an external side-effect path, dedup store
  not consulted on replay, silent overwrite on same-key conflict.
- AuthN / AuthZ — handler reachable without authentication, ownership
  predicate missing, tenant-scoping absent from a query.
- Secret exposure — credentials / tokens / connection strings in code,
  fixtures, test data, logs, or error messages.
- Injection — `format!`-built SQL (anything but parameterized
  `sqlx::query!` / `query_as!` / `query(...).bind(...)`), command
  interpolation via `Command::new("sh").arg(format!(...))`, deserialization
  of untrusted input into a type with `#[serde(deny_unknown_fields)]`
  removed, SSRF on a user-controlled URL.
- Audit — state change with no audit event emitted, or with an event_type
  not declared in Generate output.
- Compensation — irreversible external call with no declared compensating
  action and no `human-queue` escalation in the workflow.
- PII handling — PII written to a log line, error response, or test fixture
  without masking.
- Real (or real-looking — Luhn-valid PAN, valid-TLD email, real-shape SSN)
  PII in test fixtures.
- Hand-rolled crypto where a repo-approved helper exists; direct use of
  low-level crypto primitives (`BlockCipher`, `StreamCipher`, raw `Mac`,
  `ring` internals) outside a wrapping crate; `rand::random()` for keys
  or nonces.
- Hand-edited generated file (`OUT_DIR`, `tonic-build` / `prost` outputs,
  vendored bindings) or committed migration.
- Unsafe Rust without a `// SAFETY:` comment AND a declaration in
  `implement_stage_output.decision_metadata`; `unsafe fn` without a
  `# Safety` rustdoc section; `unsafe impl Send` / `unsafe impl Sync`
  without per-field justification; `mem::transmute`, `Pin::new_unchecked`,
  `Box::leak`, or `std::hint::unreachable_unchecked` without rationale;
  `'static` widened via `transmute` or `Box::leak`.
- `unwrap()` / `expect(...)` / `panic!()` / `unreachable!()` / `todo!()` /
  `unimplemented!()` on a request or worker path. (Documented-invariant
  exception below.)
- Money — `f32` / `f64` holding a monetary value; bare `+` / `-` / `*` on
  a money type without `checked_*` wrapping.
- Persistence — read-then-write of a money column without `FOR UPDATE`
  or serializable-isolation discipline.
- Async — blocking std I/O (`std::thread::sleep`, blocking `std::fs::*`,
  `std::net::*`) inside an `async fn`; CPU-bound work in async without
  `spawn_blocking`; `std::sync::MutexGuard` held across `.await`;
  `tokio::spawn` capturing `Rc` / `RefCell` on a multi-thread runtime;
  non-cancel-safe future inside a `tokio::select!` branch on a
  money-critical path.
- Errors — branching on `err.to_string()` instead of typed variants.
- Supply chain — unresolved RUSTSEC advisory of severity `high`/`critical`
  in the dependency graph; network I/O in `build.rs` or any
  `[build-dependencies]` crate's `build.rs`.
- Tests — real network call (not via `wiremock::MockServer`) in a test.

### `P2` — fix before merge

Correctness or discipline issues the next Generate iteration can address.
`P2` routes to `loop_back` (target = `implement`).

Categories typically `P2`:

- `?`-propagated error with no `#[from]` / `#[source]` / `.context(...)`
  wrapper — cause lost in the chain.
- `.map_err(|e| Custom::Msg(e.to_string()))` lossy chain when `#[from]`
  would carry the source.
- Library public API returning `anyhow::Result` instead of a `thiserror`
  enum.
- Error variant not classified (`client | server | dependency`) at the
  type level; retry policy will be guessed at the edge.
- Cancellation token / deadline not propagated into a downstream call.
- `tokio::spawn` detached on a request path (no `JoinSet`).
- `Rc<RefCell<_>>` reachable from an `async fn`.
- Non-cancel-safe future inside `tokio::select!` on a non-money path.
- Async resource cleanup hidden in `Drop` (use explicit `async fn close`).
- Blocking call inside an `async fn` (`std::thread::sleep`, blocking
  `std::fs::*`, `std::sync::Mutex` across `.await`) on a non-money path.
- Missing observability on a declared failure mode (counter / log / span).
- Missing `#[tracing::instrument]` on a handler.
- Unbounded label / span field cardinality.
- `unreachable!()` used as control flow (not a true impossibility).
- `expect("...")` without a documented invariant on a non-request path.
- `#[derive(Debug)]` on a PII-bearing struct.
- Missing `#[serde(deny_unknown_fields)]` on an externally-deserialized DTO.
- Validation deep in handler instead of at deserialization boundary.
- Money — `Decimal` rounding without explicit `RoundingStrategy`; monetary
  signature without currency tagging.
- Time — `DateTime<Local>` across boundary; `SystemTime` /
  `DateTime<Utc>` subtraction as elapsed measurement; direct `Utc::now()`
  in handler instead of injected `Clock`.
- Test missing for a declared failure mode.
- Per-file coverage claim implausible (test obviously thin).
- Pattern choice declared by Generate not consistent with the code.
- Convention divergence not flagged by Generate.
- Wall-clock-dependent test without `tokio::time::pause()` or injected clock.
- `unwrap()` / `expect(...)` in setup / `main` / one-shot init code.
- Project structure — glob `pub use foo::*;` re-exporting items not
  declared in the design; `pub` widening of a previously `pub(crate)` item
  without SemVer-relevant signal; `#[non_exhaustive]` removed from a
  public enum/struct.
- Cargo hygiene — `Cargo.lock` edited without `Cargo.toml` change or
  `cargo update`; `[[bin]]` reusing `[lib]` name; non-additive features.
- Lint baseline — `unsafe_code = "forbid"` weakened in this diff.
- Supply chain — missing `deny.toml`; yanked crate in `Cargo.lock`; new
  proc-macro dep from unfamiliar publisher.

### `P3` — note, don't block

Style, scope, or future-improvement items. `P3` does NOT block `approve`.
The finding is carried forward as a note for the next sprint / refactor.

Categories typically `P3`:

- Code duplication that could be extracted (within reason — banking-grade
  prefers boring over clever).
- Doc-comment / rustdoc style.
- Test could be parameterized (`rstest`, `#[test_case]`) but isn't (when
  correctness is unaffected).
- Naming inconsistency that doesn't break discovery.
- Adjacent issue spotted outside `code_under_review` (file the reviewer
  noticed in passing).
- Crate could be `#![forbid(unsafe_code)]` but is not.
- Runtime `sqlx::query` used where `sqlx::query!` would compile.
- Hardcoded `PgPoolOptions::max_connections` (not config-driven).
- Missing `rust-version` (MSRV) pin in `Cargo.toml`.
- Format-string `info!` log instead of structured fields.
- `info!` / `warn!` inside a tight loop.
- Missing `#[from]` / `#[source]` on a wrapping error variant.
- `context(format!(...))` on a hot path (closure form unused).
- Public unsealed trait the design declares as closed.
- Clippy lint that is `warn`-level in the repo's lint floor and not on a
  banking-grade path.

## Documented-invariant `expect` exception

The default rule (R2) classifies `.unwrap()` / `.expect(...)` / `panic!()`
on a request or worker path as `P1`. **Exception**:
`.expect("invariant: <reason>")` where the invariant is established by a
verifiable precondition earlier in the same function (or by the function's
documented contract / type-system guarantees) may be downgraded as follows.

| Reviewer confidence the invariant is real and unviolated | Severity |
|----------------------------------------------------------|----------|
| **High** — the invariant is documented, the precondition is visible at file:line, and the reviewer can reason from the code alone | `P3` |
| **Medium** — invariant is plausible but the precondition spans multiple files or requires caller context | `P2` |
| **Low** — invariant is asserted but not visibly established | `P1` (no downgrade) |

The `expect` string MUST begin with the literal `invariant:` prefix to be
eligible for downgrade. `expect("ok")`, `expect("should work")`,
`expect("never None")`, etc. remain `P1`. `.unwrap()` is never eligible
for downgrade under this exception — it carries no invariant string.

If the input's `convention_overrides.expect_invariant_policy` is set to
`forbid` (string form), this exception does not apply globally — every
`expect` on a request path remains `P1` regardless of invariant comment.

### Per-module override (object form)

`expect_invariant_policy` may also be an object:

```json
{"default": "allow", "forbid_under": ["src/payments/**"], "allow_under": ["src/cli/**"]}
```

Resolution:

1. Compute the default policy from `default`.
2. If the finding's file matches any `forbid_under` glob, policy = `forbid`.
3. Else if the finding's file matches any `allow_under` glob, policy = `allow`.
4. Else policy = `default`.

`forbid_under` wins over `allow_under` when both match (restrictive
precedence — matches Semgrep `paths.include`/`paths.exclude` and
SonarQube `sonar.inclusions`/`sonar.exclusions`). Globs are matched
with gitignore semantics via the `globset` crate; see
`references/glob-policy-idioms.md` for syntax. `!` negation is not
allowed inside the glob lists — use `allow_under` to express re-includes.

## Critical-paths severity promotion

Top-level input field `critical_paths: [{globs, reason, promote_to?}]`
declares per-path severity zones. When a finding's `file` matches one of
the `globs`, severity is promoted to the declared `promote_to` (default
`P1`) — but **only** for categories where promotion makes sense.

| Category | Eligible for promotion? |
|----------|-------------------------|
| `data_correctness` (incl. money / time) | yes |
| `idempotency` | yes |
| `security` (auth, injection, PII, crypto, unsafe wrap) | yes |
| `audit` | yes |
| `compensation` | yes |
| `async_blocking` (R3 — canonical case) | yes |
| `cancellation` | yes |
| `unsafe_rust` | yes |
| `panicking_api` | already `P1` by default on production paths; promotion is a no-op |
| `supply_chain` | yes |
| `observability` | no — cosmetic; promotion would create noise |
| `errors` | no — discipline, not correctness |
| `tests` | no — exception: real PII in fixtures is already `P1` |
| `scope`, `claims_vs_reality` | no |
| `project_structure`, `cargo_hygiene`, `lint_baseline` | no |

When promotion applies, the finding's `evidence` MUST include the marker
`[critical: <reason>]` (the `reason` value from the matching
`critical_paths` entry) so the audit trail makes the promotion visible.

Glob matching uses the `globset` crate (gitignore semantics — see
`references/glob-policy-idioms.md`). If multiple `critical_paths`
entries match, the first match in array order wins for `reason` —
purely cosmetic, since `promote_to` is conventionally identical across
overlapping entries.

The harness enforces this via `CriticalPathPromotionAssertion`: any
finding in an eligible category whose `file` matches any
`critical_paths.globs` MUST have severity `>= promote_to`.

## Confidence

Borrow the discipline from `reviewing-software-security`:

- **High** — file:line cited, behavior reproducible from the code alone.
- **Medium** — pattern is present but the exploit / failure path requires
  context the reviewer doesn't have (e.g. upstream caller behavior,
  runtime feature flag).
- **Low** — suspicion only. The reviewer cannot cite file:line confidently.

**Hard rule** (verbatim from the source): never publish a `P1` / `P2` at
`Low` confidence without an explicit `[needs verification]` tag in the
finding's `evidence` field. Do not fabricate file:line references,
standards identifiers, CWE numbers, clippy lint names, or RUSTSEC
advisory IDs — withhold instead.

When confidence is `Low`, drop the severity by one tier (`P1` → `P2`,
`P2` → `P3`) before applying the verdict matrix, unless the design
itself flagged the area as known-fragile.

## Verdict matrix

Applied after every finding has a severity. The verdict is the highest
escalation any finding produces, with one extra check for unsubstantiated
claims.

| Condition | Verdict | `loop_back_target_stage` |
|-----------|---------|--------------------------|
| Any `P1` finding | `human-queue` | null |
| `claims_unverified` non-empty (and no `P1`) | `loop_back` | `implement` |
| Any `P2` finding (and no `P1` / no unverified claims) | `loop_back` | `implement` |
| `uncertainty_flag` of kind `design_ambiguity` raised in step 2 | `loop_back` | `design` (overrides `implement` routing) |
| Only `P3` findings (and no `P1` / no `P2` / no unverified claims / no design ambiguity) | `approve` | null |
| No findings at all | `approve` | null |

Notes:

- `design_ambiguity` always wins routing — fixing it at `implement` is
  treating the symptom. Send it back to `design`.
- `human-queue` from a `P1` is final for this stage. The workflow policy
  decides whether the human can override.
- `approve` with `P3` findings still emits them — the next stage receives
  them as a notes list, not as blockers.

## Standards identifiers

When a finding maps to a published standard, cite it in
`finding.standards_ref`. Withhold (omit the field) rather than guess.

Acceptable identifier shapes:

- `CWE-###`
- `OWASP A0#:2021` or `OWASP API#:2023`
- `ASVS V#.#.#`
- `NIST SSDF PW.#.#` / `PS.#.#`
- `CIS Kubernetes #.#.#`
- `SLSA L#`
- `RUSTSEC-YYYY-NNNN` (advisory database identifier; do not invent one)
- Clippy lint name in the shape `clippy::<lint_name>` (only if the lint
  actually exists — withhold if uncertain)

If unsure of the exact identifier, omit the field. Do NOT invent.
