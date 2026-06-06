# Rust Template Defaults

Use these defaults only when they match the target crate or the task is greenfield. Treat any existing crates as examples of reusable patterns, not as a required layout. When the repo already has a convention, the repo wins (decision-rules B2).

## Crate & workspace layout

- Prefer a Cargo **workspace** for multi-crate projects: a virtual root `Cargo.toml` with `[workspace] members = [...]`, **`resolver = "2"`** declared explicitly (a virtual workspace defaults to "1"), and shared settings in `[workspace.package]`, `[workspace.dependencies]`, and `[workspace.lints]`.
- One crate per bounded capability. Library crate (`src/lib.rs`) for reusable logic; binary crate (`src/main.rs`) for the deployable. Do not reuse the `[lib]` name for a `[[bin]]` (target ambiguity).
- Module tree: prefer `foo/mod.rs` **or** `foo.rs` + `foo/` (2018+ style) — never both for the same module (half-rename). Keep modules small and intention-named.
- Keep the public API surface minimal: default to `pub(crate)`; promote to `pub` only what the contract requires. Mark public enums/structs that may grow `#[non_exhaustive]`.

## Error type scaffold

```rust
// Domain / library errors: a typed enum, one variant per failure mode the caller may branch on.
#[derive(Debug, thiserror::Error)]
pub enum OrderError {
    #[error("order {0} not found")]
    NotFound(OrderId),
    #[error("insufficient balance")]
    InsufficientBalance,
    #[error(transparent)]
    Db(#[from] sqlx::Error),
}

impl OrderError {
    pub fn class(&self) -> ErrorClass { /* client | server | dependency */ }
}
```

- Libraries return the `thiserror` enum, not `anyhow::Result`. Binaries / top-level glue may use `anyhow` with `.context(...)`.
- Do not collapse a source error into a string (`map_err(|e| Custom(e.to_string()))`) — it discards the chain.

## Async runtime defaults

- `#[tokio::main]` for the binary; `#[tokio::test]` for async tests. Choose `flavor = "multi_thread"` for servers, `current_thread` only for deliberately single-threaded contexts.
- Every public `async fn` that issues IO carries a deadline or `CancellationToken`. No detached `tokio::spawn` on a request path — use `JoinSet`.

## Lint floor

Declare a workspace `[lints]` table (Rust 1.74+) and keep it intact:

```toml
[workspace.lints.rust]
unsafe_code = "forbid"          # weaken only with an approved, documented reason

[workspace.lints.clippy]
all = "warn"
pedantic = "warn"
```

- Do not introduce blanket `#[allow(clippy::all)]` / `#[allow(warnings)]` near sensitive logic. Prefer a narrow, commented `#[allow(...)]` at the smallest scope when genuinely needed.

## Cargo & toolchain hygiene

- Pin the toolchain in `rust-toolchain.toml` with a concrete `channel = "1.NN.N"`, not a moving label (`stable`/`beta`).
- Pin MSRV in `package.rust-version` (>= the toolchain floor you test against).
- Declare `[profile.*]` and `[patch]` only at the workspace root (member redeclarations are silently ignored).
- Gate every `optional = true` dependency behind a feature and gate its `use` with `#[cfg(feature = "...")]`. Keep features additive.
- Edit `Cargo.lock` only via `cargo update` / a `Cargo.toml` change.

## Persistence-backed crates

- Prefer migrations (e.g. `sqlx migrate` / `refinery`), `sqlx::query!`/`query_as!` compile-time macros, a committed `.sqlx/` offline cache, explicit `Transaction` scopes, indexes tied to query plans, and tests for concurrency-sensitive paths. See `rust-persistence-and-data.md`.

## HTTP / message crates

- Prefer stable request/response types, boundary validation, a consistent error envelope, `#[serde(deny_unknown_fields)]` on inbound DTOs, and backward-compatible additions.

## Test layout

```rust
#[cfg(test)]
mod tests {
    use super::*;

    #[tokio::test]
    async fn rejects_duplicate_command() { /* ... */ }
}
```

- Unit tests in-module under `#[cfg(test)]`; cross-crate/integration tests in `tests/*.rs`. Cover the success path, each error variant, the replay path, and each validation rejection. Use an injected clock and `tokio::time::pause()` for determinism — never `sleep` for ordering, never a real network host (see `editing-guardrails.md` and the banking profile).

## Security-sensitive code

- Prefer vetted crates and repo-approved crypto/auth wrappers, centralized secret handling, `OsRng` for keys/nonces, and no custom cryptography.
