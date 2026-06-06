# Editing Guardrails

## Read Before Write

- Read before writing. Do not rewrite files from memory or apply broad mechanical changes without inspecting current code, the crate's `Cargo.toml`, its error type, and its trait boundaries.

## Do Not Introduce Without Request

- Do not introduce a new dependency, a new async runtime, a new ORM/SQL layer, a new logging/tracing stack, a new error crate, or a new CI gate unless requested. If a dependency is genuinely required, explain why the existing stack cannot solve the problem and flag it explicitly (no silent `[dependencies]` additions).

## Protected Artifacts — Require Explicit Approval

Require explicit approval or an accepted migration plan before any of:

- **Public API / SemVer surface**: widening visibility (`pub(crate)` -> `pub`), removing `#[non_exhaustive]`, changing a public `fn`/`trait`/`struct`/`enum` shape, or glob-re-exporting (`pub use foo::*`) undesigned items.
- **Safety floor**: adding `unsafe` to a `#![forbid(unsafe_code)]` crate, or weakening `unsafe_code = "forbid"` to `deny`/`allow`.
- **Lint floor**: removing/loosening the `[lints]` table or adding blanket `#[allow(warnings)]` / `#[allow(clippy::all)]`.
- **Data & contracts**: database schema, persisted formats, message schemas, generated clients, auth/authz/session behavior.
- **Generated & migrated files**: `OUT_DIR` / `prost` / `tonic-build` output, committed migration SQL (add a new additive migration instead), or `Cargo.lock` `[[package]]` lines outside `cargo update`.
- **Toolchain**: `rust-toolchain.toml` channel, `package.rust-version` (MSRV), edition.

## Scope Discipline

- Do not delete files, run `cargo fmt` across unrelated files, or clean up code outside the requested scope.
- Keep dependency changes separate from feature fixes when possible.
- Do not chase borrow-checker errors by sprinkling `clone()` / `Arc<Mutex<_>>` — reconsider ownership (see `rust-idioms-and-async.md` § 2).

## Safe Change Patterns

- Prefer feature flags, additive fields/variants (behind `#[non_exhaustive]`), compatibility shims, and expand-then-contract migrations for risky changes.
- Preserve public APIs and trait interfaces unless the task is explicitly a breaking refactor with an approved deprecation path.

## Tool-availability degradation

The host sandbox may lack `cargo`, a compiler toolchain, a database, or network access. When a validation step cannot run:

- State the exact status: `not_run: cargo unavailable in sandbox`, `unavailable: no database for sqlx offline prepare`, or `blocked: network egress denied for cargo-audit`.
- **Never** claim a check passed that you did not run, and never silently omit it. An unrun check is a residual risk to report, not a pass.
- Prefer the narrowest check that *can* run (`cargo check` over a full `cargo test`) and say which you ran. For `sqlx::query!` macros without a live DB, rely on the committed `.sqlx/` offline cache and note it.
