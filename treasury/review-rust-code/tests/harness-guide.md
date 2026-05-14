# Test Harness — review-rust-code

> **File-naming note.** Originally `tests/README.md` per the project's
> standard naming, renamed to `harness-guide.md` because
> `scripts/quick_validate.py` (BANNED_DOCS) rejects any `README.md`
> inside a skill folder. Same content; the validator wins.

## How the harness invokes the skill

The shared test harness (`COGNITIVE_OS.md` Section 9
`internal/test_harness/runner.go`) drives this skill the same way the
workflow engine will:

1. Load `SKILL.md` + every file in `references/` into LLM context.
2. Read `tests/cases/NNN-<name>.input.json`; validate against
   `schemas/input.json`. Case-level validation failure = harness bug,
   fix the case.
3. Send the input to the LLM as the stage payload; capture the raw response.
4. Parse the response as JSON; validate against `schemas/output.json`.
5. Compare against `tests/cases/NNN-<name>.expected.json` using the
   assertion set below.
6. Emit a per-case `TestResult` with per-assertion pass/fail.

Run the full suite for this skill:

```bash
go run ./cmd/harness -skill treasury/review-rust-code -case-glob 'tests/cases/*.input.json'
```

## How assertions check schema validity

| Layer | Checks | On failure |
|-------|--------|-----------|
| `SchemaValidityAssertion` | Output JSON validates against `schemas/output.json` | Fail case — schema break is always a blocker |
| Structural-diff against `expected.json` | `verdict`, `loop_back_target_stage`, finding *categories* + *severities* (not full evidence text), `claims_unverified` set | Fail with diff |
| `audit_metadata` envelope check | Numeric fields present and non-negative | Fail with diff |

Evidence snippets and exact line numbers are NOT diffed verbatim — the LLM
may pick a different but equally valid citation. The harness checks the
*shape* of the verdict and findings; the human reviews exact wording.

## Banking-grade assertions

These run on every case in addition to schema validity.

| Assertion | What it verifies | Maps to |
|-----------|------------------|---------|
| `VerdictRoutingAssertion` | If `verdict == "loop_back"`, `loop_back_target_stage` is non-null. If `verdict != "loop_back"`, it is null. | severity-guide.md verdict matrix |
| `P1RoutingAssertion` | Any `P1` finding produces `verdict == "human-queue"`. No exceptions. | severity-guide.md P1 rule |
| `UnverifiedClaimsAssertion` | If `claims_unverified` is non-empty AND no `P1` is present, verdict is `loop_back` to `implement`. | severity-guide.md claims rule |
| `DesignAmbiguityRoutingAssertion` | If any `uncertainty_flag.kind == "design_ambiguity"`, verdict is `loop_back` and target is `design`. | severity-guide.md design override |
| `EvidenceCitationAssertion` | Every finding has either a file:line pair OR an explicit `[needs verification]` tag in `evidence`. | SKILL.md Anti-Patterns rule against fabrication |
| `NoCodeEmissionAssertion` | The output JSON contains no field that would deliver code (this skill is read-only). | SKILL.md Anti-Patterns DO NOT emit code |
| `RulesEvaluatedFloorAssertion` | `audit_metadata.rules_evaluated >= 34` (11 base + 7 augmentation + 12 Rust-specific + 4 contract). Lower means the reviewer skipped rules. Fast-path check; `rule_walk` below is the authoritative depth check. | review-rubric.md scope |
| `RuleWalkCompletenessAssertion` | `audit_metadata.rule_walk` MUST contain exactly one entry for each rule in {B1..B11, A1..A7, R1..R12, C1..C4} (34 distinct rule IDs). Each entry's `outcome` is one of `finding`/`no_finding`/`not_applicable`. `not_applicable` entries MUST carry a non-empty `rationale`. Missing rule IDs OR extra IDs fail this assertion — depth check beyond the count floor. | schemas/output.json `audit_metadata.rule_walk` |
| `CriticalPathPromotionAssertion` | When any input `critical_paths[].globs` matches a finding's `file` and the finding's `category` is eligible for promotion (see severity-guide.md § Critical-paths severity promotion), the finding's `severity` MUST equal or exceed the matching `critical_paths[].promote_to` (default `P1`). Evidence MUST contain the marker `[critical: <reason>]`. | severity-guide.md § Critical-paths severity promotion |
| `RuleWalkFindingConsistencyAssertion` | Every `findings[].rule_violated` value MUST correspond to a `rule_walk` entry with `outcome == "finding"`. Conversely, every `rule_walk` entry with `outcome == "finding"` MUST be cited by at least one entry in `findings[]`. | schemas/output.json `rule_walk` |

## Adding a new case

1. Pick the next prefix: `002-`, `003-`, …
2. Common cases worth adding:
   - `002-missing-idempotency.input.json` — handler with side-effect but no key. Expected: `human-queue`, `P1`, category `idempotency`.
   - `003-design-ambiguity.input.json` — design says "TBD" on auth. Expected: `loop_back` to `design`.
   - `004-unverified-audit-claim.input.json` — Generate claims event_type X, code emits Y. Expected: `loop_back` to `implement`, `claims_unverified` non-empty.
   - `005-unwrap-in-handler.input.json` — `.unwrap()` on a request path. Expected: `human-queue`, `P1`, category `panicking_api`.
   - `006-unsafe-block.input.json` — `unsafe { ... }` block with no `// SAFETY:` comment and not declared in `decision_metadata`. Expected: `human-queue`, `P1`, category `unsafe_rust`.
   - `007-blocking-in-async.input.json` — `std::thread::sleep` inside an `async fn`. Expected: `loop_back` to `implement`, `P2`, category `async_blocking`.
3. Write `NNN-<name>.input.json` — must validate against `schemas/input.json`.
4. Write `NNN-<name>.expected.json` — must validate against `schemas/output.json`, including a complete `audit_metadata.rule_walk` array (all 34 rule IDs).
5. Re-run the suite.

### Cases shipped today

Every Rust-specific rule (R1–R12) now has at least one fixture. The
banking-grade rule families (B6, A1, A2 idempotency; A5 test discipline)
also have coverage.

| Case | Exercises | Verdict |
|------|-----------|---------|
| `001-clean-handler` | Happy path; backward-compat for `rust_specific` block | `approve` |
| `002-missing-idempotency` | B6 + A1 + A2 (idempotency + audit claim) | `human-queue`, claims_unverified non-empty |
| `005-unwrap-in-handler` | R2 (panicking API on production path) | `human-queue` |
| `006-unsafe-block-undocumented` | R1 + R4 + C3 (unsafe undeclared + lint floor weakened + claim inconsistency) | `human-queue`, claims_unverified non-empty |
| `007-blocking-in-async` (with `critical_paths`) | R3 + critical-path promotion P2 → P1 | `human-queue` (promoted) |
| `008-pub-use-glob` | R5 + C3 (glob re-export widens API surface; pattern_choice inconsistency) | `loop_back` to `implement` |
| `009-cargo-lock-edit-no-update` | R6 (Cargo.lock edited without manifest change or `cargo update`) | `loop_back` to `implement` |
| `010-utc-now-in-handler` | R8 + C3 + A5 (`Utc::now()` in handler instead of injected Clock; non-reproducible test) | `loop_back` to `implement` |
| `011-format-string-info` | R10 × 2 + C3 (format-string log + unbounded metric label cardinality) | `loop_back` to `implement` |
| `012-missing-deny-unknown-fields` | R11 × 2 + C3 + C2 (missing `#[serde(deny_unknown_fields)]` + PII in `Debug`; thin coverage claim) | `loop_back` to `implement` |
| `013-rustsec-yanked` | R12 + C3 (yanked tokio 1.35.0 in `Cargo.lock`; pinned in `Cargo.toml`) | `loop_back` to `implement` |

Future iterations could add fixtures for: explicit `unsafe` *with* a
`// SAFETY:` comment + declaration (the happy path of R1, to avoid
false positives); `tokio::select!` with a non-cancel-safe future (R3
sub-case); SQL `format!` injection (R9); compile-time-checked
`sqlx::query!` happy path (R9 negative); per-module
`expect_invariant_policy: {default: forbid, allow_under: [...]}`
demonstration; `critical_paths.promote_to: P2` (downgrading the
promotion target).

## What a passing run looks like

```
ok   001-clean-handler                schema=PASS verdict=approve     routing=PASS evidence=PASS claims=PASS rules=34 rule_walk=34
ok   002-missing-idempotency          schema=PASS verdict=human-queue routing=PASS rules=34 rule_walk=34
ok   005-unwrap-in-handler            schema=PASS verdict=human-queue routing=PASS rules=34 rule_walk=34
ok   006-unsafe-block-undocumented    schema=PASS verdict=human-queue routing=PASS rules=34 rule_walk=34
ok   007-blocking-in-async            schema=PASS verdict=human-queue routing=PASS critical_path_promoted=PASS rules=34 rule_walk=34
ok   008-pub-use-glob                 schema=PASS verdict=loop_back   routing=PASS rules=34 rule_walk=34
ok   009-cargo-lock-edit-no-update    schema=PASS verdict=loop_back   routing=PASS rules=34 rule_walk=34
ok   010-utc-now-in-handler           schema=PASS verdict=loop_back   routing=PASS rules=34 rule_walk=34
ok   011-format-string-info           schema=PASS verdict=loop_back   routing=PASS rules=34 rule_walk=34
ok   012-missing-deny-unknown-fields  schema=PASS verdict=loop_back   routing=PASS rules=34 rule_walk=34
ok   013-rustsec-yanked               schema=PASS verdict=loop_back   routing=PASS rules=34 rule_walk=34
```

`verdict=loop_back` or `verdict=human-queue` are also valid pass states for
cases written to exercise those routes — the harness checks shape, not the
verdict value.
