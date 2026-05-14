# Rust Project Structure Reference

Authoritative reference for adversarial code review of Rust backend crates. Every section ends with an **Adversarial scan question** the reviewer must answer when reading the feature under review.

## Workspace vs Single Crate

Use a **single crate** when one publishable unit owns all code. Use a **workspace** when multiple crates share a lockfile, target dir, and dependency versions (services + shared libs, multi-binary deployments, plugin trees).

```toml
# Root Cargo.toml (virtual workspace)
[workspace]
members  = ["crates/*", "services/api"]
exclude  = ["crates/legacy-*"]
resolver = "2"            # required for edition 2021+; "3" on 1.84+

[workspace.package]
edition      = "2021"
rust-version = "1.74"
license      = "Apache-2.0"

[workspace.dependencies]
serde  = { version = "1", features = ["derive"] }
tokio  = { version = "1", features = ["rt-multi-thread"] }

[workspace.lints.rust]
unsafe_code = "forbid"
```

```toml
# Member crate inherits
[package]
name           = "billing-core"
version.workspace      = true
edition.workspace      = true
rust-version.workspace = true

[dependencies]
serde = { workspace = true, features = ["rc"] }   # additive features OK
tokio.workspace = true

[lints]
workspace = true
```

`[patch]`, `[replace]`, and `[profile.*]` are only honored in the **root** manifest ([workspaces.html](https://doc.rust-lang.org/cargo/reference/workspaces.html)). `workspace.lints` requires Rust 1.74+.

**Adversarial scan questions**
- Does a member redeclare `[profile.*]` or `[patch]` (silently ignored)?
- Does any member pin a different version of a workspace dep instead of `.workspace = true`?
- Is `resolver` absent from a virtual workspace (defaults to "1" silently)?

## Standard Directory Layout

| Path | Role | Notes |
|---|---|---|
| `src/lib.rs` | Library root | Module tree root; public API surface |
| `src/main.rs` | Default binary | Implicit `[[bin]]` named after the package |
| `src/bin/<name>.rs` | Extra binaries | Name = file stem (or dir name for multi-file) |
| `tests/*.rs` | Integration tests | Each file = separate crate; **public API only** |
| `examples/*.rs` | Example binaries | Built with `cargo build --examples` |
| `benches/*.rs` | Benchmarks | Built with `cargo bench` |
| `build.rs` | Build script | Runs before crate compile |

Unit tests live in-source under `#[cfg(test)] mod tests { use super::*; ... }` and **can see private items**. Integration tests in `tests/*.rs` are external crates and can only call `pub` items ([ch11-03](https://doc.rust-lang.org/book/ch11-03-test-organization.html)). Shared test helpers must live in `tests/common/mod.rs` (subdirectory `mod.rs` form), not `tests/common.rs`, or Cargo treats them as a standalone test target.

**Adversarial scan questions**
- Is logic in `src/main.rs` that should be in `src/lib.rs` to be reachable by `tests/`?
- Does a `tests/foo.rs` import a private item (impossible — confirms it's a unit test misplaced)?
- Are integration tests reaching into the crate via `#[path]` to bypass `pub` (privacy escape hatch)?

## Module System

Two file forms for `mod foo;` ([reference/items/modules](https://doc.rust-lang.org/reference/items/modules.html)):

```
src/foo.rs              <-- modern preferred form
src/foo/bar.rs          <-- child modules live in foo/

src/foo/mod.rs          <-- legacy form; do not mix with foo.rs (E0761)
```

```rust
// lib.rs
mod internal;            // private; descendants only
pub mod api;             // public surface
pub(crate) mod util;     // visible to whole crate, not consumers
pub(super) fn helper(){} // visible to parent module
pub(in crate::api) fn x(){} // visible inside crate::api subtree

pub use api::v1::Client;          // re-export flattens path
pub use internal::types::*;       // glob re-export: avoid; pollutes API

#[path = "../shared/codegen.rs"]  // override file location
mod codegen;

mod inline { pub fn k() {} }      // inline form (rare outside tests/macros)
```

`pub use foo::*;` is a contract hazard: adding a new `pub` item to `foo` silently widens the parent's API ([visibility-and-privacy](https://doc.rust-lang.org/reference/visibility-and-privacy.html)).

**Adversarial scan questions**
- Is `name.rs` AND `name/mod.rs` both present (compile error, but check for half-renames)?
- Is `#[path]` used to load code from outside the crate's `src/` (vendoring without provenance)?
- Does a glob `pub use foo::*;` re-export items the author did not intend to publish?
- Did a previously `pub(crate)` item become `pub` in this diff (API-surface widening — SemVer)?

## Cargo.toml Schema

```toml
[package]
name         = "billing-core"
version      = "0.4.2"
edition      = "2021"          # 2015/2021/2024
rust-version = "1.74"          # MSRV; cargo enforces

[lib]
name       = "billing_core"
crate-type = ["rlib"]          # add "cdylib" only if FFI

[[bin]]
name = "billd"
path = "src/bin/billd.rs"

[dependencies]
serde  = { version = "1", features = ["derive"] }
tonic  = { version = "0.12", optional = true }

[dev-dependencies]
proptest = "1"

[build-dependencies]
prost-build = "0.13"

[features]
default = ["grpc"]
grpc    = ["dep:tonic"]        # dep:foo = no implicit feature, namespaced
sqlite  = ["dep:rusqlite", "rusqlite?/bundled"]   # 1.60+ conditional

[profile.release]
lto           = "thin"
codegen-units = 1
panic         = "abort"
strip         = "symbols"

[patch.crates-io]
ring = { git = "https://github.com/example/ring", rev = "..." }

[lints.rust]        # Rust 1.74+
unsafe_code        = "forbid"
unreachable_pub    = "warn"
missing_docs       = "warn"

[lints.clippy]
pedantic           = { level = "warn", priority = -1 }
expect_used        = "deny"
unwrap_used        = "deny"
```

Feature-gate code: `#[cfg(feature = "grpc")] pub mod grpc;`. Features must be **additive** (enabling one must not remove API) ([features](https://doc.rust-lang.org/cargo/reference/features.html)).

**Adversarial scan questions**
- Does `[[bin]]` reuse the same `name` as `[lib]` causing target-name ambiguity?
- Is an `optional = true` dependency present without a feature gating its `extern crate` / `use` (always compiled)?
- Are features non-additive (mutually exclusive without `compile_error!`)?
- Is `crate-type = ["cdylib"]` set without an FFI boundary (forces dynamic linking on consumers)?
- Are `unwrap_used`/`expect_used` allowed in a payments path?

## Cargo.lock Discipline

`cargo new` defaults to committing `Cargo.lock` ([FAQ](https://doc.rust-lang.org/cargo/faq.html#why-have-cargolock-in-version-control)). Practical rule for backend code:

| Crate kind | Commit `Cargo.lock`? |
|---|---|
| Binary / service / application | **Yes** |
| Library published to crates.io | Optional — does not bind consumers |
| Workspace containing any binary | **Yes** (single root lockfile) |

`Cargo.lock` is machine-managed: legitimate diffs come only from `cargo update`, `cargo add`, `cargo remove`, or `cargo build` resolving a new manifest. Hand edits are a red flag.

**Adversarial scan questions**
- Does the diff edit `Cargo.lock` without a corresponding `Cargo.toml` change or `cargo update` invocation in CI?
- Was a `Cargo.lock` deleted in this PR for a service crate?
- Does the lockfile pin a yanked version (cross-check RUSTSEC)?

## Toolchain Pinning

```toml
# rust-toolchain.toml at repo root
[toolchain]
channel    = "1.82.0"          # never "stable" in regulated builds
components = ["clippy", "rustfmt", "rust-src"]
targets    = ["x86_64-unknown-linux-musl"]
profile    = "minimal"
```

`channel = "stable"` makes the build non-reproducible across rustup updates ([rustup overrides](https://rust-lang.github.io/rustup/overrides.html)).

**Adversarial scan questions**
- Is `channel` a moving label (`stable`/`beta`/`nightly`) instead of a pinned version?
- Is the pinned MSRV in `package.rust-version` lower than the toolchain channel (silent regression risk)?

## Lint Configuration

Prefer the `[lints]` table (Rust 1.74+) over crate attributes — it composes through workspace inheritance and survives `cargo fix --edition`. Use crate attributes only for what `[lints]` cannot express (e.g., `#![forbid(unsafe_code)]` enforced at file load).

```rust
// lib.rs — banking-grade floor
#![forbid(unsafe_code)]                          // unforbiddable downstream
#![deny(rust_2018_idioms, future_incompatible)]
#![warn(missing_docs, unreachable_pub)]
#![allow(clippy::module_name_repetitions)]       // surgical opt-out
```

Companion files: `clippy.toml` (e.g., `cognitive-complexity-threshold = 20`), `rustfmt.toml` (e.g., `edition = "2021"`, `max_width = 100`). `forbid` cannot be downgraded by inner `#[allow]`; `deny` can ([attributes/diagnostics](https://doc.rust-lang.org/reference/attributes/diagnostics.html)).

**Adversarial scan questions**
- Is `unsafe_code = "forbid"` weakened to `deny` or `allow` in this diff?
- Are blanket `#[allow(clippy::all)]` / `#[allow(warnings)]` introduced near sensitive logic?
- Does `clippy.toml` raise complexity thresholds rather than refactor?

## Visibility and API Surface

The **public contract** is the transitive closure of items reachable through `pub`/`pub use` from `lib.rs`. Widening visibility is a SemVer-relevant change.

```rust
#[non_exhaustive]
pub enum LedgerError { Insufficient, Frozen }  // consumers must use `_ =>`

#[non_exhaustive]
pub struct Config { pub host: String }         // forces ..Default::default()

// Sealed trait: only this crate may implement
pub trait Currency: sealed::Sealed { fn iso(&self) -> &str; }
mod sealed { pub trait Sealed {} }
impl sealed::Sealed for crate::Usd {}
impl Currency for crate::Usd { fn iso(&self) -> &str { "USD" } }
```

`#[non_exhaustive]` is the canonical knob for future-proof enums/structs ([api-guidelines/future-proofing](https://rust-lang.github.io/api-guidelines/future-proofing.html)).

**Adversarial scan questions**
- Did a previously `pub(crate)` item become `pub` without a SemVer-minor bump?
- Was `#[non_exhaustive]` removed from a public enum/struct?
- Did a sealed-trait `Sealed` supertrait become accidentally public?
- Are struct fields `pub` directly (pinning representation) instead of behind accessors?

## Build Script (`build.rs`)

Legitimate uses ([build-scripts](https://doc.rust-lang.org/cargo/reference/build-scripts.html)):
- Compile bundled C/C++ via `cc`
- Probe linker / target with `CARGO_CFG_TARGET_*`
- Generate Rust from `.proto`/`.fbs`/IDL into `OUT_DIR`
- Embed compile-time constants (`cargo::rustc-env=...`)

Red flags: network calls, reading `$HOME`, writing outside `OUT_DIR`, executing downloaded binaries, behavior gated only on environment variables without `rerun-if-env-changed`.

```rust
// build.rs — minimum hygiene
fn main() {
    println!("cargo::rerun-if-changed=proto/billing.proto");
    println!("cargo::rerun-if-env-changed=PROTOC");
    tonic_build::compile_protos("proto/billing.proto").unwrap();
}
```

**Adversarial scan questions**
- Does `build.rs` open sockets, spawn `curl`/`wget`, or read paths outside the crate?
- Are there no `rerun-if-changed` directives (full-package scan, caches over-invalidate)?
- Does the script write into `src/` instead of `OUT_DIR` (commits regenerate-on-build artifacts)?
- Does it `panic!()` to "skip" platforms instead of using `cfg`?

## Generated Code

```rust
// src/proto.rs
include!(concat!(env!("OUT_DIR"), "/billing.v1.rs"));
```

Generated output belongs in `OUT_DIR` only; vendored bindings (e.g., `bindgen --no-rebuild`) belong under `src/generated/` with a header banner. **Never hand-edit** files marked generated — the next build rewrites them.

**Adversarial scan questions**
- Is a file under `src/generated/` modified in this diff without regenerator change?
- Is `OUT_DIR` output checked into git?
- Does `include!` reference a path outside `OUT_DIR` (escapes Cargo's clean model)?

## Relevant RUSTSEC Advisory Categories

Cross-check the advisory database ([rustsec.org/advisories](https://rustsec.org/advisories/)) for structural risks:

| Category | What to look for |
|---|---|
| **Yanked / removed crate** | Dep version no longer on crates.io; `Cargo.lock` still references it |
| **Malicious code (supply chain)** | Newly added dep with low download count, typo-squat name, or recent `cargo owner` change |
| **Unmaintained** | Advisories tagged `unmaintained` (e.g., `async-std`, pre-0.17 `ring`); plan migration |
| **Build-time RCE / corruption** | `build.rs` of a dep executing untrusted input; e.g., `pyo3` PYO3_CONFIG_FILE class |
| **Unsoundness** | `pub` API exposes UB without `unsafe` — escalate if crate is in the trust boundary |
| **Proc-macro compromise** | Proc-macro crates can run arbitrary code at compile time; pin tightly, audit updates |
| **Notice** | Behavioral or licensing notice; not a vuln but may force version floor |

**Adversarial scan questions**
- Is any direct or transitive dep listed in RUSTSEC with status `vulnerability` or `unmaintained`?
- Does this PR add a proc-macro dep from an unfamiliar publisher?
- Is `cargo deny check advisories` wired into CI for this workspace?
