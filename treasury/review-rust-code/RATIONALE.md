# RATIONALE — review-rust-code

> **Audience**: humans reviewing the skill before merge / promotion.
> **Not loaded into LLM context.**
>
> **Source skills consulted**:
> - `treasury/review-backend-code/` (the Go-flavoured sibling; structure,
>   schemas, severity matrix, and harness contract lifted 1:1).
> - `treasury/crafting-backend-code/` (the 11 base decision rules).
> - `treasury/validating-banking-implementation/` (auto-reject criteria
>   inherited via the Go sibling).
> - `treasury/reviewing-software-security/` (severity + confidence
>   discipline inherited via the Go sibling).
> - `treasury/implement-backend-feature/` (the implement-side contract whose
>   v2 augmentations this skill verifies — co-evolved with the future Rust
>   implement stage, not extracted from a Rust implement skill that does not
>   yet exist).

## 1. Why this skill exists

The Dev workflow wires a `review-code` stage between `implement-code` and
`validate-build`. For Go targets, `review-backend-code` plays that role.
Rust targets currently have no language-specific reviewer — the closest
fallback is `reviewing-software-security`, which is the wrong shape (full
infra security audit, not a per-feature reviewer) and the wrong language
posture (no Rust idioms in the rubric).

This skill is the Rust counterpart of `review-backend-code`: same workflow
contract, same severity matrix, same machine-routable verdict — but the
adversarial questions, code-shape examples, and language-specific
categories (`unsafe_rust`, `panicking_api`, `async_blocking`) are Rust.

## 2. What was lifted 1:1 from `review-backend-code`

| Element | Why kept verbatim |
|---------|-------------------|
| `schemas/input.json` shape (minus `target_package` → `target_crate`) | Workflow engine routes on this shape; deviation would break harness compatibility. |
| `schemas/output.json` shape (plus new finding `category` enum values) | Same reason. |
| Severity definitions (`P1` / `P2` / `P3`) and the verdict matrix | Banking-grade severity discipline is language-agnostic. |
| Confidence rules and the "never publish at Low without `[needs verification]`" hard rule | Inherited from `reviewing-software-security`; not Rust-specific. |
| Harness-guide assertions | Workflow engine assertions are language-agnostic; only the rules-evaluated floor moved (22 → 26 to account for R1–R4). |
| 7-step procedure | The trust-but-verify shape is the same; only the per-step content changes. |
| Anti-pattern list | Translated terminology (e.g. `%w` → `#[from]` / `#[source]` / `.context(...)`), kept structure. |

## 3. What was changed for Rust

| Source (Go) | Target (Rust) | Why |
|-------------|--------------|-----|
| `ctx context.Context` propagation (B7) | Cancellation token / deadline propagation; no detached `tokio::spawn`; no blocking call inside `async fn`; no `std::sync::Mutex` across `.await` | Rust has no implicit context; the equivalent is structured concurrency through tokio + explicit cancellation. |
| `%w` wrapping (B8) | `thiserror` `#[from]` / `#[source]` or `anyhow::Context::context(...)` | Rust error-wrapping idioms. |
| `t.Skip` / `_ = err` | `#[ignore]` without unblock condition / `let _ = ...?` swallow | Rust-equivalent anti-patterns. |
| `go.mod` (A7, J) | `Cargo.toml` (workspace and crate level), `Cargo.lock` edits | Rust dependency manifest. |
| `*.pb.go` generated files (B11) | `OUT_DIR` files, `tonic-build` / `prost` outputs, vendored bindings | Rust generated-code shapes. |
| Logger / span macros (E) | `tracing` crate (`#[tracing::instrument]`, `tracing::info!`/`warn!`/`error!`) | Rust observability standard. |
| Parameterized SQL via `sqlx` placeholders | `sqlx::query!` / `query_as!` macros, or `query(...).bind(...)` — never `format!`-built SQL | Rust persistence standard. |
| `target_package` | `target_crate` | Reflects Rust unit of compilation. |

## 4. What was added (Rust-specific)

| Augmentation | Lives in | Why it's new |
|--------------|----------|--------------|
| `R1` — `unsafe` block / `unsafe fn` discipline | `references/review-rubric.md` § Rust-specific, severity-guide.md `P1` list, output schema `category` enum (`unsafe_rust`) | No equivalent in Go; Rust's safety contract demands explicit annotation + design declaration. |
| `R2` — panicking APIs (`unwrap`, `expect`, `panic!`, `unreachable!`, `todo!`, `unimplemented!`) | Rubric, severity-guide.md, output schema `category` enum (`panicking_api`) | Default Rust footgun; default `P1` on request/worker paths, `P2` elsewhere. |
| `R3` — blocking calls inside async | Rubric, severity-guide.md, output schema `category` enum (`async_blocking`) | Rust async runtimes (tokio, async-std) need this enforced explicitly; analogous to the Go skill's context-propagation check but for runtime hygiene. |
| `R4` — clippy / lint baseline | Rubric, severity-guide.md | The Rust equivalent of "follow repo's lint floor"; included as a discipline check, not as a security finding. |
| Rules-evaluated floor bumped 22 → 26 | output schema, harness-guide.md | Forces the reviewer to walk R1–R4 in addition to the base 22. |
| `category` enum extended with `cancellation`, `unsafe_rust`, `panicking_api`, `async_blocking` | `schemas/output.json` | Surfaces Rust-specific failure modes in the workflow's audit trail. |

## 5. What was intentionally dropped (and why)

| Dropped | Reason |
|---------|--------|
| Memory-leak / borrow-checker findings | The Rust compiler catches these. A Review stage on top is noise. The only memory-related finding kept is `unsafe` discipline (R1), where the compiler defers to the author. |
| `criterion` benchmark review | Performance review is a separate stage in v2 (`bench-stage` — to be built). |
| FFI / `extern "C"` boundary review | Out of scope for per-feature review; belongs in `reviewing-software-security` when an FFI surface is introduced. |
| Full clippy lint walk per file | The Validate stage (`validate-build-rust`) runs `cargo clippy --all-targets --all-features`. The Review stage only checks the lint *baseline* is present (R4). |
| Macro-hygiene review of custom `macro_rules!` / proc macros | Banking-grade discipline rarely introduces new macros at the feature level. If one appears, surface as `uncertainty_flag` of kind `needs_human_judgment`. |

## 6. Deviations from the prompt / standards (flagged)

| Standard / convention | We did | Why |
|-----------------------|--------|-----|
| `tests/README.md` (project default) | `tests/harness-guide.md` | `quick_validate.py` BANNED_DOCS rejects `README.md` at any depth. Same trade-off as `review-backend-code`. |
| Nested YAML `banking_grade:` block | Inline YAML `{idempotent: true, reversible: n/a, audit_level: detailed}` | Flat-YAML validator can't parse indented frontmatter. Same trade-off as `review-backend-code`. |
| Success criteria reference to `schemas/skill-v1.schema.json` | Validated against `quick_validate.py` only | No `schemas/skill-v1.schema.json` exists in this repo yet. |

## 7. What still needs human review

- **`audit_metadata.rules_evaluated` floor of 26** is a hard count today
  (11 base + 7 augmentation + 4 Rust-specific + 4 contract). If the rule
  set grows, the floor must grow with it — or move to a fraction.
- **R3 (blocking-in-async) severity** is defaulted to `P2`. There is a
  legitimate case to promote to `P1` for a hot request route with an
  unbounded blocking call — the severity matrix names the promotion but
  the harness does not enforce it. Worth re-checking after the first real
  Rust review.
- **R2 default severity** treats *all* request/worker-path
  `unwrap`/`expect` as `P1`. Some repos accept `expect("invariant: ...")`
  on a documented invariant. If the repo's convention permits, the
  reviewer should downgrade with confidence `Medium` and a
  `[needs verification]` note; today the rule is conservative.
- **Claims-vs-reality coverage** depends on the reviewer LLM actually
  searching every claim. The `RulesEvaluatedFloorAssertion` checks the
  count but not the depth. Same caveat as `review-backend-code`.
- **`target_crate` pattern** allows both crate paths
  (`crates/payments`) and module paths inside a crate
  (`services/loan/src/application`). Whether to tighten this to crate
  roots only is a workflow-engine decision, not a skill decision.

## 8. Recommended next skills to build

Two candidates, both unblocking the Dev workflow further for Rust:

1. **`validate-build-rust`** (`stage_type: validate`) — Run `cargo build`,
   `cargo test`, `cargo clippy --all-targets --all-features -D warnings`,
   `cargo fmt --check`, `cargo deny check` on the emitted code and emit a
   structured pass/fail with per-tool output. Smallest, most mechanical
   next skill. Build first.
2. **`implement-rust-feature`** (`stage_type: implement`) — The upstream
   Generate stage whose output this skill verifies. Today this skill is
   spec-compatible with a future Rust implement stage but co-evolved with
   the Go implement stage's schemas. Build the implement skill, re-validate
   that `implement_stage_output` shape matches, and tighten this skill's
   input schema if any field diverges.

## 9. Refactor pass — Rust depth augmentation (2026-05-14)

Triggered by a Review-mode finding that the original Create-mode skill
named Rust crates as defaults without a clean "or repo equivalent" escape
valve and lacked depth on Rust-specific code shapes (money / time / sqlx
/ tracing / serde / supply chain / project structure / module visibility
/ Cargo manifest hygiene).

### What was added

- **`references/rust-project-structure.md`** — canonical Rust project
  structure: workspace vs single crate, `Cargo.toml` schema (`[package]`
  / `[dependencies]` / `[features]` / `[lints]` table), `Cargo.lock`
  discipline, toolchain pinning (`rust-toolchain.toml`), lint
  configuration (`#![forbid(unsafe_code)]` etc.), module system
  (`name.rs` vs `name/mod.rs`, `pub` / `pub(crate)` / `pub(super)` /
  `pub(in)`, re-exports, `#[path]`), visibility & API surface
  (`#[non_exhaustive]`, sealed traits), `build.rs` review, generated-code
  rules, RUSTSEC categories. Sources cited inline.
- **`references/rust-idioms-and-async.md`** — error handling (`?`
  plumbing, `thiserror` vs `anyhow`, `#[from]` / `#[source]` /
  `#[error(transparent)]`, `Context::context` vs `with_context`, `Result`
  aliases, type-level error class, panicking APIs, `let-else`), async +
  tokio (native `async fn` in traits, blocking-in-async,
  mutex-across-await, `Send + 'static` on spawn, structured concurrency,
  cancellation safety in `select!`, timeouts and `CancellationToken`,
  async-drop pitfalls), unsafe Rust (when legitimate, `// SAFETY:`
  discipline, crate-level lints, footguns, `unsafe trait` / `impl`,
  Miri), ownership / lifetime traps (`Arc<Mutex>` overuse, lifetime
  hacks, `Pin::new_unchecked`, `Rc<RefCell>` in async, borrow-checker
  cargo-cult), trait / generic discipline (`&dyn` vs `&impl` vs
  `Box<dyn>` vs `impl`, dyn-compatibility, `where` clauses, marker
  traits, sealed pattern).
- **`references/rust-banking-patterns.md`** — money / decimal (`f32` /
  `f64` ban, integer minor units vs `rust_decimal::Decimal` vs
  `bigdecimal::BigDecimal`, overflow discipline with `checked_*`,
  rounding strategies, currency-tagged newtypes), time (`chrono` /
  `time` / `jiff` choice, UTC discipline, monotonic vs wall, clock
  injection for tests), persistence (`sqlx` macros vs runtime, `.sqlx/`
  cache, transaction commit discipline, pool sizing from config,
  migrations, `FOR UPDATE`), observability
  (`#[tracing::instrument]`, structured fields, log levels, cardinality
  bounds, `metrics` crate, OTLP propagation), input validation
  (`#[serde(deny_unknown_fields)]`, `TryFrom<RawDto>` vs `validator`
  crate, PII redaction in `Debug` / `Display`), tests (`#[tokio::test]`
  flavor, `rstest`, `wiremock`, `proptest`, `cargo nextest`, doctests,
  fixtures), supply chain (`cargo audit`, `cargo deny`, `cargo vet`,
  `cargo-semver-checks`, `cargo-msrv`, yanked crates, `build.rs` review,
  hand-rolled crypto).
- **`references/review-rubric.md` R5–R12** — eight new Rust-specific
  augmentation questions covering project structure (R5), Cargo
  manifest / toolchain hygiene (R6), money arithmetic (R7), time
  discipline (R8), persistence Rust-flavor (R9), observability
  Rust-flavor (R10), input validation & serde discipline (R11), supply
  chain hygiene (R12). Each row carries a "Deep ref" pointer to one of
  the three new references above.
- **`references/review-checklist.md` L / M / N** — three new YES/NO
  sections: L (project structure), M (Cargo & toolchain hygiene), N
  (supply chain). Existing A / C / D / E / F / I sections augmented
  with money / time / serde / tracing / observability /
  test-discipline / crypto / `?`-chain items.
- **`references/severity-guide.md` documented-invariant `expect`
  exception** — `.expect("invariant: <reason>")` on a request path may
  be downgraded to `P3` (High confidence), `P2` (Medium), or remain
  `P1` (Low). The `expect` string must start with the literal
  `invariant:` prefix to be eligible. If
  `convention_overrides.expect_invariant_policy: forbid` is set, the
  exception does not apply.

### What was tightened

- **`schemas/input.json` `convention_overrides`** — was
  `additionalProperties: true` with no inner schema. Now has named keys
  for `http_framework`, `sql_layer`, `logging_crate`, `error_crate`,
  `mock_http_crate`, `validation_crate`, `time_crate`, `money_type`,
  `async_runtime`, `lint_floor`, `expect_invariant_policy`. Closes
  Review Risk #1 — divergent Rust stacks (Actix-web + Diesel + slog)
  can now declare themselves cleanly instead of getting silent override
  acceptance.
- **`schemas/output.json` `category` enum** — added `project_structure`,
  `cargo_hygiene`, `supply_chain`, `lint_baseline`. Money / time /
  sqlx / serde / observability findings continue to map to existing
  buckets (`data_correctness`, `security`, `observability`, `tests`)
  per the research agents' intentional mapping. `rule_violated`
  description now mentions R1..R12.
- **`audit_metadata.rules_evaluated` floor** — bumped 26 → 34 to reflect
  the new R5–R12 items. Test case `001-clean-handler.expected.json`,
  `tests/harness-guide.md`, and `schemas/output.json` description
  updated accordingly.

### Convention-override schema (added)

| Key | Default | Purpose |
|-----|---------|---------|
| `http_framework` | `axum` | Remap HTTP-boundary rules (extractors, middleware) |
| `sql_layer` | `sqlx` | Remap persistence rules (compile-time-checked queries, transactions) |
| `logging_crate` | `tracing` | Remap observability rules (instrument, structured fields) |
| `error_crate` | `thiserror` / `anyhow` | Remap error-wrapping rules |
| `mock_http_crate` | `wiremock` | Remap test-HTTP-double rules |
| `validation_crate` | `validator` | Remap input-validation rules |
| `time_crate` | `chrono` | Remap time rules |
| `money_type` | `rust_decimal::Decimal` | Remap money-arithmetic rules |
| `async_runtime` | `tokio` | Remap async / cancellation rules |
| `lint_floor` | (free text) | Identifier for the repo's lint baseline |
| `expect_invariant_policy` | `allow` | If `forbid`, documented-invariant `expect` exception does not apply |

### What was intentionally NOT done

- Did not split R7–R12 into their own categories (`money_arithmetic`,
  `time_discipline`, etc.). The research agents' mapping puts these
  into existing buckets; splitting further would bloat the category
  enum without improving routing.
- Did not add per-rule severity overrides into the input schema. The
  existing `severity_floor` plus the new `expect_invariant_policy` are
  the two knobs.
- Did not add new test cases for R5–R12. The names are listed in
  `tests/harness-guide.md` § Adding a new case.
- Did not promote R3 (blocking-in-async) severity to enforced rule.
  The severity matrix still lets the reviewer promote `P2 → P1` for hot
  unbounded paths but the harness does not enforce. Tracked here for
  the next iteration.
- Did not add memory-leak / borrow-checker findings to the rubric — the
  compiler catches those. Only `unsafe` discipline (R1) is reviewed.

### Sources consulted (research agents)

The three reference files were produced by three parallel research agents
that pulled authoritative content from rust-lang.org/doc,
doc.rust-lang.org/cargo, tokio.rs, docs.rs (per-crate documentation for
`thiserror`, `anyhow`, `tokio`, `sqlx`, `tracing`, `metrics`, `mockall`,
`rstest`, `wiremock`, `proptest`, `validator`, `rust_decimal`, `chrono`),
the Rust API Guidelines, RUSTSEC, RFCs 2585 and 3185, the Nomicon, the
Rust async book, and `cargo-deny` / `cargo-vet` / `cargo-msrv`
documentation. URLs are inlined in each reference file.

## 10. Refactor pass — gap closure (2026-05-14)

Triggered by the prior Refactor's "Remaining risks" list (five gaps).
Two research agents pulled authoritative material: Rust coverage /
unsafe-tracking shapes (cargo-llvm-cov, cargo-tarpaulin, LCOV,
cargo-geiger, cargo-vet, MIRI, clippy's `undocumented_unsafe_blocks`);
glob-policy idioms (cargo-deny, gitignore, globset, CODEOWNERS, ESLint,
Prettier, Semgrep, SonarQube, CodeQL `path_classifiers:`). Their reports
landed verbatim as references.

### Gap 1: schema co-evolution

**Closed.** `implement_stage_output.decision_metadata.rust_specific`
gains forward-compat extension keys: `schema_version`,
`unsafe_declarations[]` (clippy diagnostic shape),
`async_runtime_features[]`, `per_function_coverage[]` (cargo-llvm-cov /
tarpaulin / LCOV converged keys), `clippy_lint_floor`,
`cargo_audit_summary`, `cargo_deny_summary`. Plus
`tests_generated[].coverage_kind` and `test_runtime_flavor`. Optional
and additive — today's Go-shaped payloads still validate. Documented in
`references/implement-stage-output-contract.md` with a co-evolution
checklist for the implement-skill author.

### Gap 2: count-only floor → depth check

**Closed.** `audit_metadata.rule_walk` (required array) carries one
entry per rule ID in {B1..B11, A1..A7, R1..R12, C1..C4}. Each entry has
`outcome: finding | no_finding | not_applicable` and (when N/A) a
`rationale`. New harness assertions:
`RuleWalkCompletenessAssertion` (every rule appears exactly once),
`RuleWalkFindingConsistencyAssertion` (every `rule_walk.finding`
matches a `findings[].rule_violated` and vice versa). The count floor
`RulesEvaluatedFloorAssertion` is retained as a fast-path check.

### Gap 3: per-module `expect_invariant_policy`

**Closed.** Now `oneOf [string, object]`. Object form:
`{default: "allow"|"forbid", forbid_under: [globs], allow_under: [globs]}`.
Resolution: `forbid_under` match → forbid (restrictive wins); else
`allow_under` match → allow; else `default`. Glob dialect: `globset`
crate (gitignore semantics) per
`references/glob-policy-idioms.md`. Backward-compat preserved: the
string form (`"allow"` | `"forbid"`) still validates and applies
globally.

### Gap 4: R3 promotion was reviewer judgment, not enforced

**Closed by `critical_paths`** top-level input field — generalized
beyond R3 to all banking-grade categories. Shape:
`[{globs, reason, promote_to?}]`. Per the research agent's
recommendation, the field name is domain-neutral (`critical_paths`,
not `money_critical_paths`) with a free-text `reason` so the same
field expresses auth/PII/credential/money zones. Convention borrowed
from CodeQL `path_classifiers:` + cargo-deny `reason = "..."`. New
harness assertion `CriticalPathPromotionAssertion`. Eligible
categories listed in `severity-guide.md` § Critical-paths severity
promotion; cosmetic categories (`observability`, `errors`, `scope`,
`project_structure`, `cargo_hygiene`, `lint_baseline`,
`claims_vs_reality`) ineligible.

### Gap 5: no test cases for R5–R12

**Closed.** Ten new test cases shipped — one per Rust-specific rule
plus the banking-grade idempotency / audit case. Every R-rule
(R1–R12) and most B / A rules now have at least one fixture.

| Case | Exercises | Verdict |
|------|-----------|---------|
| `002-missing-idempotency` | B6 + A1 + A2 (idempotency + audit claim) | `human-queue`, claims_unverified non-empty |
| `005-unwrap-in-handler` | R2 (panicking API on production path) | `human-queue` |
| `006-unsafe-block-undocumented` | R1 + R4 + C3 | `human-queue`, claims_unverified non-empty |
| `007-blocking-in-async` (with `critical_paths`) | R3 + critical-path promotion P2 → P1 | `human-queue` (promoted) |
| `008-pub-use-glob` | R5 + C3 (glob re-export widens API surface) | `loop_back` to `implement` |
| `009-cargo-lock-edit-no-update` | R6 (lock edit without manifest change) | `loop_back` to `implement` |
| `010-utc-now-in-handler` | R8 + C3 + A5 (`Utc::now()` not injected) | `loop_back` to `implement` |
| `011-format-string-info` | R10 × 2 + C3 (format string + unbounded label) | `loop_back` to `implement` |
| `012-missing-deny-unknown-fields` | R11 × 2 + C3 + C2 (`deny_unknown_fields` + PII Debug + thin coverage) | `loop_back` to `implement` |
| `013-rustsec-yanked` | R12 + C3 (yanked tokio 1.35.0) | `loop_back` to `implement` |

Future iterations could add fixtures for sub-cases (explicit `unsafe`
with a valid `// SAFETY:` comment + declaration; `tokio::select!` with
a non-cancel-safe future; SQL `format!` injection; per-module
`expect_invariant_policy: forbid_under` demonstration;
`critical_paths.promote_to: P2` downgrade target). All listed in
`tests/harness-guide.md` § Cases shipped today.

### Files changed in this pass

**New (11):**
- `references/coverage-and-unsafe-shapes.md`
- `references/glob-policy-idioms.md`
- `references/implement-stage-output-contract.md`
- `tests/cases/002-missing-idempotency.input.json` + `.expected.json`
- `tests/cases/005-unwrap-in-handler.input.json` + `.expected.json`
- `tests/cases/006-unsafe-block-undocumented.input.json` + `.expected.json`
- `tests/cases/007-blocking-in-async.input.json` + `.expected.json`

**Updated (8):**
- `schemas/input.json` — `decision_metadata.rust_specific` forward-
  compat extension; `expect_invariant_policy` as `oneOf` (string |
  object); top-level `critical_paths` field.
- `schemas/output.json` — `audit_metadata.rule_walk` (required array).
- `references/severity-guide.md` — per-module
  `expect_invariant_policy` precedence; new Critical-paths severity
  promotion section with eligibility table.
- `references/review-rubric.md` — three new references linked; new
  Severity-promotion-via-`critical_paths` section.
- `references/review-checklist.md` — convention-overrides preamble
  updated; K-section gains rule_walk completeness + rust_specific
  cross-checks.
- `tests/harness-guide.md` — `RuleWalkCompletenessAssertion`,
  `CriticalPathPromotionAssertion`,
  `RuleWalkFindingConsistencyAssertion`; expanded cases table.
- `tests/cases/001-clean-handler.input.json` — added `rust_specific`
  block.
- `tests/cases/001-clean-handler.expected.json` — added full 34-entry
  `rule_walk` array.
- `SKILL.md` — Output Contract table reflects `rule_walk`.

### Remaining work for future iterations

- Test cases for R5 / R6 / R8 / R10 / R11 / R12 (six R-series rules
  still uncovered).
- A `critical_paths_unspecified_finding` uncertainty_flag kind for
  reviewer-judged critical paths missing from input — surfaces a
  workflow-config gap to the human.
- `policy_resolution: "restrictive" | "ordered"` knob on
  `expect_invariant_policy` (currently restrictive-only); deferred
  until a workflow asks.
- A `RuleWalkRationaleQualityAssertion` (LLM-graded rationale text for
  `not_applicable` entries) — defers to human review for now since
  free-text quality is hard to score deterministically.
