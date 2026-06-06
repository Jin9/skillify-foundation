# Mode Checklists

Per-mode checklists for `crafting-rust-code`. Load and fill the checklist for the mode classified in the Safety workflow (Step 1). Include the filled checklist in the output when it improves reviewability; for small fixes, apply it internally and summarize only the important result.

## `design` — Rust architecture pass

```
- [ ] L1 business invariant + success state stated
- [ ] L2 crate/module boundary + data owner + public contract (fn/trait/type) named
- [ ] L3 error model, async strategy, ownership/sharing, persistence, txn, idempotency, cancellation, observability, tests chosen
- [ ] L4 traits, types, signatures, modules, migrations, and tests sketched
- [ ] Banking profile applied if design touches money/ledger/PII
- [ ] Trade-offs stated on all axes (incl. safety/ownership vs ergonomics)
```

## `optimize` — performance / allocation / async-contention

```
- [ ] Current baseline measured or explicitly unavailable (criterion / flamegraph / tokio-console)
- [ ] Bottleneck identified, not guessed (allocations, clones, lock contention, blocking-in-async, await fan-out, DB query, serialization, monomorphization bloat)
- [ ] Options enumerated with pros/cons (borrow vs clone, Arc vs message-pass, dyn vs generic, buffer reuse)
- [ ] Correctness and rollback risks stated
- [ ] Recommendation + measurement plan
```

## `fix` — minimal, direct solution

```
- [ ] Root cause vs symptom called out
- [ ] Behavior-preserving unless flagged
- [ ] Correctness, auth, cancellation, and panic-path impact checked
- [ ] No new panic on a request path; no Mutex-across-await introduced
- [ ] Existing tests checked before adding new deps
- [ ] Regression test added or skipped with reason
```

## `analyze` — deep breakdown

```
- [ ] Current architecture and ownership summarized
- [ ] Strengths
- [ ] Gaps (correctness, safety/unsafe, security, performance, observability, tests, operability)
- [ ] Contradictions (shared state ownership, mutex-across-await, non-idempotent consumers, panic on request path)
- [ ] Recommendations prioritized by risk and effort
```

## `review` — lightweight in-flow code review

```
- [ ] Risks flagged by severity (P1/P2/P3)
- [ ] Correctness, safety (unsafe/cancellation/ownership), security, perf, observability, tests checked
- [ ] Concrete fix suggested for each finding
- [ ] Findings evidence-backed with file:line when local code is available
- [ ] Authoritative banking-grade verdict deferred to review-rust-code (do not emit approve/loop_back here)
```

## `plan` — produce a plan, do not execute

```
- [ ] Priorities set
- [ ] Open decisions listed with owners (product, backend, data, infra, security)
- [ ] Dependencies + sequencing stated (trait contract, schema, migration, feature flag, MSRV/edition, rollout)
- [ ] Validation and rollback plan included
```
