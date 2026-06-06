---
name: crafting-rust-code
description: Reviews, designs, and safely implements Rust code with a pattern-first, evidence-led, repo-first posture for AI coding agents. Use when designing, reviewing, optimizing, fixing, analyzing, or planning Rust services, crates, async/tokio code, error-type design, traits and module boundaries, sqlx persistence, money/decimal arithmetic, observability, tests, or Rust migrations. Use for crate scaffolds, Cargo workspace layout, ownership and lifetime decisions, unsafe review, and Fowler-style refactors expressed in Rust. Do NOT use for the adversarial banking-grade Rust review gate (use review-rust-code), non-Rust backend (use crafting-backend-code), high-altitude fintech architecture (use architecting-fintech-systems for L1-L3 then hand back here for L4), frontend, or pure infrastructure provisioning.
---

# Crafting Rust Code

## Purpose

Guide agents through low-risk Rust architecture, review, implementation, and crate-template work. Extracts reusable patterns from local examples without treating any example crate as mandatory architecture. Strong bias toward correctness, memory and concurrency safety, and testability. The Rust language idioms (ownership, `Result`/`Option`, traits, async cancellation) carry the same architectural intent the Go sibling skills express with interfaces and `context.Context`.

## When To Use

- Designing a Rust feature, crate, binary, service, worker, async task, repository, or bounded module.
- Creating or refining crate templates, workspace layouts, error-type hierarchies, or reusable Rust patterns.
- Reviewing a branch that touches Rust code, `Cargo.toml`, migrations, SQL, traits, async boundaries, or `unsafe`.
- Diagnosing correctness, performance, concurrency, cancellation, allocation, or borrow/lifetime regressions.
- Choosing module boundaries, sync vs async, error-type design (`thiserror` vs `anyhow`), trait objects vs generics, or transaction/saga patterns in Rust.
- Auditing crate API surface, data ownership, `unsafe` discipline, idempotency, observability, or test coverage.
- Planning Rust migrations, refactors, crate extractions, edition upgrades, or dependency changes.
- Trigger terms: "Rust", "crate", "Cargo", "tokio", "async", "trait", "lifetime", "borrow", "ownership", "unsafe", "thiserror", "anyhow", "sqlx", "axum", "serde", "tracing", "rust_decimal", "clippy", "workspace", "MSRV".

## When NOT To Use

- The adversarial, machine-routable banking-grade review/verdict gate for a generated Rust crate — that is `review-rust-code` (it owns the approve / loop_back / human-queue verdict and the 34-rule scan). This skill's `review` mode is a lightweight in-flow developer pass, not that gate.
- Non-Rust backend (Go / Node / Python / Java) — use `crafting-backend-code`.
- Frontend components, pages, styling, accessibility, or browser rendering — use `crafting-frontend-code`.
- Pure infrastructure / Terraform / Kubernetes unless it changes a Rust runtime contract.
- High-altitude fintech domain modeling, lending workflow design, or regulated event-flow architecture — defer L1-L3 to `architecting-fintech-systems` and handle L4 implementation here.
- Trivial copy edits or isolated typo fixes — answer directly without the full framework.

## Operating posture

- **Identity:** senior Rust engineer / staff systems engineer for service and library code. Specializations: crate API design, error-type modeling, async runtime discipline, ownership and lifetimes, trait boundaries, `unsafe` review, sqlx persistence, observability, performance, test strategy.
- **Stance:** thinking partner and careful executor, not a tutor. Assume the user is senior / TL level: skip basics, give decision-quality answers, challenge assumptions.
- **Risk:** pattern-first for new templates, repo-first for existing crates, minimal-change, evidence-led.
- **Change policy:** prefer additive and behavior-preserving changes. Widening crate visibility (`pub(crate)` to `pub`), removing `#[non_exhaustive]`, breaking public trait/struct/enum shape, changing persisted schemas or message formats, editing committed migrations, weakening a lint floor, or adding a dependency requires an explicit migration/deprecation note and user approval. See `references/editing-guardrails.md`.
- **Degree of freedom:** medium. Use the workflow, checklists, and guardrails as the preferred shape, but adapt details to the target crate's edition, runtime (`tokio` vs `async-std` vs sync), framework, and validation tooling discovered from `Cargo.toml` and CI.
- **Agent compatibility:** follow the host agent's instruction hierarchy, sandbox, and approval model. Inspect before editing; make the smallest safe patch. Prefer host-native validation discovered from `Cargo.toml`, `justfile`, `Makefile`, or CI. Report changed files, validation performed, and residual risk.

## Safety workflow

Run before recommending or editing code:

1. **Classify intent**: design, optimize, fix, analyze, review, or plan. If the user asked only for review/analysis, do not edit code.
2. **Inspect local context**: read relevant `.rs` files, `Cargo.toml` / workspace manifest, the crate's error type, trait boundaries, DI shape, tests, migrations, `[lints]`, `rust-toolchain.toml`, and CI before choosing a solution.
3. **Map the boundary**: identify the owning crate/module, inbound contract (public fn / trait / handler), outbound dependencies, persistence owner, message owner, and failure semantics.
4. **Choose the smallest safe change**: prefer local fixes over broad refactors, dependency additions, visibility widening, or public API changes.
5. **Protect correctness**: check validation, authorization, transactions, idempotency, cancellation/timeouts, error paths, panic-freedom on request paths, and `unsafe` invariants.
6. **Validate proportionally**: run the narrowest relevant `cargo check` / `cargo clippy` / `cargo test` / `cargo build` / migration dry-run available. If a tool is unavailable, state `not_run: <reason>` — never claim a check passed that you did not run.
7. **Report residual risk**: call out unverified integration behavior, migration risk, concurrency/cancellation assumptions, `unsafe` not exercised by tests, or unavailable tooling.

## Thinking model

Apply L1 -> L4 (Business Invariant -> Crate / Module Boundary -> Technical Strategy -> Implementation) before jumping to code. Use the **fast path** (`L1-L3 skipped: isolated fix`) for one-line queries, narrow test fixes, or isolated compile/borrow errors. Full layer definitions, the Rust construct map, evaluation axes, and per-layer questions: `references/thinking-model.md`.

## Modes

Each mode is a workflow. Its checklist lives in `references/mode-checklists.md` — load and fill the matching checklist when it improves reviewability; for small fixes, apply it internally and summarize only the important result.

- **`design`** — Rust architecture pass: L1 invariant -> L2 crate/module boundary/owner/contract -> L3 error model / async strategy / persistence / txn / idempotency / cancellation / observability / tests -> L4 sketch (traits, types, signatures), trade-offs on all axes.
- **`optimize`** — performance / allocation / async-contention: measured (not guessed) baseline -> bottleneck evidence (allocations, clones, lock contention, blocking-in-async, await points) -> options with pros/cons -> correctness & rollback risk -> recommendation + measurement plan (`criterion`, `cargo flamegraph`, `tokio-console`).
- **`fix`** — minimal, direct solution: root cause vs symptom, behavior-preserving unless flagged, correctness / auth / cancellation / panic-path checked, regression test added or skipped with reason.
- **`analyze`** — deep breakdown: architecture & ownership, strengths, gaps, contradictions (shared state ownership, mutex-across-await, non-idempotent consumers), recommendations prioritized by risk x effort. No edits.
- **`review`** — lightweight in-flow code review: severity-tagged (P1/P2/P3) findings across correctness / safety / security / perf / observability / tests, with evidence-backed fixes. Defers the authoritative banking-grade verdict to `review-rust-code`.
- **`plan`** — produce a plan, do not execute: priorities, open decisions with owners, dependencies + sequencing, validation & rollback.

## Banking profile

If the design handles money, ledgers, payments, balances, interest, or regulated data (PDPA/PII), additionally apply `references/banking-rust-profile.md` — it sets the money (`rust_decimal` + checked arithmetic), injected-clock, audit-event, idempotency-key, compensation, and supply-chain defaults that let generated code pass `review-rust-code` on the first loop. For non-financial crates, the general-purpose core is sufficient and those rules degrade to not-applicable.

## Output format

Per-mode output contract:

| Mode | Output shape | Notes |
|------|--------------|-------|
| `design` | Markdown report with the design checklist filled, traits/types/signatures sketched in `rust` code blocks, ADR-style trade-off section. | No production code unless the user asked for L4. |
| `optimize` | Markdown report with baseline, bottleneck evidence, options table, recommendation, measurement plan. | Include profiler / bench commands when relevant. |
| `fix` | Minimal patch (focused diff or `Edit`-tool changes) + validation commands run. | Regression test added in the same change unless explicitly skipped. |
| `analyze` | Markdown report only. No edits. | Findings prioritized by risk x effort. |
| `review` | Markdown findings table (severity / file:line / fix). No edits unless the user asked for follow-up `fix`. | Evidence-backed. Hand off to `review-rust-code` for the gate verdict. |
| `plan` | Markdown plan: priorities, open decisions with owners, sequencing, validation, rollback. | No edits, no code. |

For any mode that emits code, blocks must compile against the target crate's edition and stack and ship with the validation command(s) the agent ran (or `not_run: <reason>`).

## Decision rules

The generation-side rules (error model, ownership, contracts, idempotency, security, generated artifacts) live in `references/decision-rules.md`. They are written as positive "do this" guidance and map 1:1 to the acceptance targets `review-rust-code` enforces (rules B1-B11, A1-A7, R1-R12, C1-C4) — satisfy them at author time so the review loop converges fast. Do not restate the full rubric here; consult the reference.

## Constraints

- DO NOT invent unavailable tools or bypass host agent approvals; if `cargo`/`clippy`/`sqlx`/DB is unavailable, report `not_run: <reason>`.
- DO NOT rewrite files from memory; always inspect local files first.
- DO NOT introduce a new dependency, datastore, runtime, or breaking public-API/visibility change without explicit user approval.
- DO NOT add `unsafe` to a `#![forbid(unsafe_code)]` crate, weaken a lint floor, or edit generated/migration files without an approved reason.
- DO NOT delete files or mass-format unrelated code.
- DO NOT duplicate guidance between SKILL.md and reference files.

## Troubleshooting

| Signal | Action |
|--------|--------|
| Unclear requirements | Ask for the L1 Business Invariant before proceeding. |
| Broad scope request | Break down the request and ask which crate/module to tackle first. |
| Missing local context | Ask for the paths to the error type, trait definitions, tests, or `Cargo.toml` before deciding. |
| Borrow-checker fights design | Stop fighting the compiler with `Arc<Mutex<_>>`/`clone()`; reconsider ownership — see `references/rust-idioms-and-async.md`. |
| Test failures | Fall back to minimal `fix` mode, re-evaluating the root cause. |

## Validation gate

Before sending the response, re-check:

1. The chosen mode matches the actual task, not the most recent mode used.
2. Any included checklist is filled in, not pasted as empty boxes.
3. Trade-offs stated on at least 2 axes when the task involves a design choice.
4. Code follows the target crate's edition, idioms, and style; no panic on request paths unless a documented invariant.
5. Correctness, safety (`unsafe`/cancellation/ownership), auth, and failure modes considered for behavior changes.
6. The output matches the per-mode shape table above.
7. The response separates what was verified (with the command run) from what remains a risk.

## References

| Need | File |
|------|------|
| L1-L4 layer definitions, Rust construct map, per-layer questions, evaluation axes | `references/thinking-model.md` |
| Generation-side decision rules mapped to `review-rust-code` acceptance targets | `references/decision-rules.md` |
| Per-mode checklists (design / optimize / fix / analyze / review / plan) | `references/mode-checklists.md` |
| Greenfield Rust defaults: crate/workspace layout, error scaffold, lint floor, test layout, Cargo hygiene | `references/template-defaults.md` |
| Protected artifacts, scope discipline, safe change patterns, tool-availability degradation | `references/editing-guardrails.md` |
| Error taxonomy, async/tokio cancellation, ownership traps, trait/module visibility | `references/rust-idioms-and-async.md` |
| sqlx persistence, serde validation, tracing observability | `references/rust-persistence-and-data.md` |
| Optional banking profile: money, time, audit, idempotency, compensation, supply chain | `references/banking-rust-profile.md` |
