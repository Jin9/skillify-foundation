# Implement-Stage Output Contract

Documents the expected shape of `implement_stage_output` consumed by
`review-rust-code`. The Review stage trust-but-verifies these fields
against the code emitted by the upstream Generate stage. The forward-
compat extension keys close the Schema-co-evolution gap so the Rust
implement skill (when it ships) can declare Rust-specific facts in a
documented, schema-validated shape rather than free-form blobs.

## Backward-compatibility note

Today's `schemas/input.json` accepts payloads shaped after
`treasury/implement-backend-feature` (the Go implement skill), with
`additionalProperties: true` on every nested object. The Rust implement
skill is forward-compat: when it ships, it MAY emit additional fields
under `decision_metadata.rust_specific` and per-test under
`tests_generated[]`. The Review skill consumes these when present and
ignores them when absent — no payload that validates today will stop
validating tomorrow.

## Required top-level keys (today)

| Key | Type | Notes |
|---|---|---|
| `files_generated` | array of `{path, lines}` | Production files emitted by Generate. |
| `tests_generated` | array of `{path, coverage_pct, kind?}` | Companion tests. |
| `idempotency_strategy` | string | e.g. `"header_uuid_v4"`, `"naturally_idempotent_read_only"`. |
| `compensating_actions` | array of `{trigger, action_skill_ref}` | Empty array if no irreversible side-effects. |
| `audit_events_emitted` | array of string `event_type` values | Empty for read-only paths. |
| `uncertainty_flags` | array of `{kind, location, note}` | Flagged ambiguities. |
| `decision_metadata` | object with `pattern_choices` | Plus optional `rust_specific` (below). |

## Rust-flavor extension keys (forward-compat)

These keys SHOULD appear in the Rust implement skill's output. The
Review skill consumes them when present and treats absence as
backward-compat.

### `decision_metadata.rust_specific`

| Key | Type | Source |
|---|---|---|
| `schema_version` | string | `"1.1-rust"` or successor; signals which extension shape applies. |
| `unsafe_declarations` | array | Per-block records mirroring clippy `--message-format=json` shape. See `references/coverage-and-unsafe-shapes.md` § Unsafe tracking for the canonical fields (`file_path`, `line_start`, `line_end`, `kind`, `safety_comment_present`, `safety_justification`, `miri_status`, `miri_commit`, `audit_ref`, `allow_lint_attribute`). |
| `async_runtime_features` | array of strings | e.g. `["tokio_rt_multi_thread", "tokio_macros"]`. Reviewer cross-checks against `#[tokio::test]` flavors and `tokio::spawn` usage. |
| `per_function_coverage` | array | Per-function coverage entries — keys converge across `cargo-llvm-cov` / `cargo-tarpaulin` / LCOV; see `coverage-and-unsafe-shapes.md` § Coverage (`function_name`, `file_path`, `entry_line`, `execution_count`, `lines_covered`, `lines_total`, `regions_covered`, `regions_total`, `branches_covered`, `branches_total`). |
| `clippy_lint_floor` | string | Mirrors `convention_overrides.lint_floor` so reviewer can cross-check declared vs actual `[lints]` table. |
| `cargo_audit_summary` | object | Optional `{advisories_unresolved: int, yanked_count: int, last_run_commit: string}`. |
| `cargo_deny_summary` | object | Optional `{advisories_pass: bool, licenses_pass: bool, bans_pass: bool, sources_pass: bool, last_run_commit: string}`. |

### `tests_generated[]`

| Key | Type | Notes |
|---|---|---|
| `coverage_kind` | enum `line\|branch\|region\|function` | Optional. Defaults to `line`. Drives the plausibility check at C2. |
| `test_runtime_flavor` | enum `current_thread\|multi_thread` | Optional. Defaults to `current_thread`. Mismatch with production runtime is a smell flagged under A5. |
| `per_function_coverage_ref` | string | Optional. Index into `decision_metadata.rust_specific.per_function_coverage[]` for the function this test exercises. |

## How the Review skill consumes the extension

| Extension key | Cross-checked against | Rule (when violation) |
|---|---|---|
| `unsafe_declarations[]` | Every `unsafe { ... }` block found in `code_under_review` must appear here; every entry here must point at a real `unsafe` block at the cited `file_path:line_start` | R1 — missing declaration or unannotated block = P1 (`unsafe_rust`) |
| `async_runtime_features` | `#[tokio::test]` flavors, `tokio::spawn` usage, `tokio_*` feature gates in `Cargo.toml` | R3 — feature mismatch or runtime/flavor mismatch = P2 (`async_blocking`) |
| `per_function_coverage[]` | `tests_generated[].coverage_pct` claim plausibility — sum of covered lines / total lines must be within ±2 pp of the claim | C2 — implausible coverage = P2 (`tests`) |
| `clippy_lint_floor` | `[lints]` table or crate-level `#![deny(...)]` / `#![forbid(...)]` in code | R4 — floor weakened or absent = P2 (`lint_baseline`) |
| `cargo_audit_summary.advisories_unresolved > 0` | RUSTSEC walk (manual or via `cargo audit` output) | R12 — unresolved advisory = P1 (`supply_chain`) |
| `cargo_audit_summary.yanked_count > 0` | `Cargo.lock` cross-check | R12 — yanked dep = P2 (`supply_chain`) |
| `cargo_deny_summary.advisories_pass == false` | `deny.toml` presence and freshness | R12 — `deny.toml` absent or gating failed = P2 (`supply_chain`) |

## Versioning

When the Rust implement skill ships its v1 output, both skills tag the
co-versioned contract via
`decision_metadata.rust_specific.schema_version`.

- Bump **major** when a required field is removed or changes type.
- Bump **minor** for additive fields (the default).
- Bump **patch** for description-only clarifications.

The Review skill MUST tolerate unknown fields (forward-compat) and MUST
log an `uncertainty_flag` of kind `convention_conflict` when a
`schema_version` it does not understand appears.

## Co-evolution checklist for the Rust implement skill author

When the implement skill ships, the schemas can stay aligned by walking:

1. Compare implement output keys against this contract's "Required top-
   level keys" table. Any new required key needs a `schemas/input.json`
   update here AND a corresponding rule walk entry in
   `references/review-rubric.md`.
2. Compare implement output keys against the "Rust-flavor extension"
   table. Any new optional key gets added to
   `decision_metadata.rust_specific` with `additionalProperties: true`
   preserved at parent level so unknown keys don't fail validation.
3. Bump `schema_version`. Add a "what changed" entry to this file.
4. Add a test case to `tests/cases/` exercising the new field.

## Today's verified gap

Until the Rust implement skill exists:

- `unsafe_declarations` is always empty in real payloads.
- `per_function_coverage` is not emitted.
- `cargo_audit_summary` and `cargo_deny_summary` are not emitted.
- `clippy_lint_floor` is set to a free-text value if at all.

The Review skill therefore falls back to walking the code directly for
R1 (unsafe), R4 (lint baseline), R12 (supply chain) without the
implement-side declaration. Findings emitted this way carry confidence
`Medium` by default since the reviewer is inferring rather than
cross-checking.
