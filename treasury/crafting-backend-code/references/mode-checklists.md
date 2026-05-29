# Mode Checklists

Per-mode checklists for `crafting-backend-code`. Load and fill the checklist for the mode classified in the Safety workflow (Step 1). Include the filled checklist in the output when it improves reviewability; for small fixes, apply it internally and summarize only the important result.

## `design` — backend architecture pass

```
- [ ] L1 business invariant + success state stated
- [ ] L2 service/template boundary + data owner + contract named
- [ ] L3 persistence, transaction, auth, idempotency, failure handling, observability, tests chosen
- [ ] L4 files/packages, schemas, handlers, services, adapters, and tests sketched
- [ ] Trade-offs stated on all 4 axes
```

## `optimize` — performance / scalability / cost trade-off

```
- [ ] Current baseline measured or explicitly unavailable
- [ ] Bottleneck identified, not guessed (CPU, allocation, DB query, lock contention, network, queue lag, serialization, cold start)
- [ ] Options enumerated with pros/cons
- [ ] Correctness and rollback risks stated
- [ ] Recommendation + measurement plan
```

## `fix` — minimal, direct solution

```
- [ ] Root cause vs symptom called out
- [ ] Behavior-preserving unless flagged
- [ ] Data correctness, auth, and failure-path impact checked
- [ ] Existing tests checked before adding new tools
- [ ] Regression test added or skipped with reason
```

## `analyze` — deep breakdown

```
- [ ] Current architecture and ownership summarized
- [ ] Strengths
- [ ] Gaps (correctness, security, performance, observability, tests, operability)
- [ ] Contradictions (shared data ownership, hidden cross-service transactions, non-idempotent consumers)
- [ ] Recommendations prioritized by risk and effort
```

## `review` — code / API / architecture review

```
- [ ] Risks flagged by severity (P1/P2/P3)
- [ ] Correctness, security, data integrity, performance, observability, and tests checked
- [ ] Concrete fix suggested for each finding
- [ ] Findings evidence-backed with file/line references when local code is available
```

## `plan` — produce a plan, do not execute

```
- [ ] Priorities set
- [ ] Open decisions listed with owners (product, backend, data, infra, security)
- [ ] Dependencies + sequencing stated (schema, API contract, feature flag, migration, rollout)
- [ ] Validation and rollback plan included
```
