# Backend Thinking Model (L1 → L4)

Apply this layered sequence before writing or recommending code. Do not jump to implementation without passing through contract and data questions.

## L1 — Business Invariant

What rule must remain true after this change? What is the success state and who owns it?

Surface the invariant in plain language before naming any technology. Examples:
- "A loan disbursement occurs at most once per approved application."
- "A consumer that crashes mid-batch must not double-process the same Kafka offset."
- "An underwriter override is recorded with actor identity and reason, even if the system is partitioned."

If the invariant is unclear, stop and ask before proceeding. Do not invent invariants.

## L2 — Service Boundary

Name the boundary the change lives inside.

- **Service / template ownership**: which deployable owns the change.
- **API / message contract**: the shape promised to callers (REST, gRPC, Kafka topic, event schema).
- **Data ownership**: which service writes the source of truth; everyone else reads or subscribes.
- **Transactional boundary**: the unit of atomicity. Cross-aggregate writes in one transaction are a smell.
- **Authorization boundary**: who is allowed to invoke the contract; how identity is asserted.

## L3 — Technical Strategy

Choose the implementation strategy without writing code yet.

- **Persistence model**: row schema, indexes, normalization, denormalization for read paths.
- **Query / command split**: read model vs write model; CQRS only when the read shape diverges materially from the write shape.
- **Transaction scope**: ACID inside one aggregate; outbox / saga / compensating actions across aggregates.
- **Retries and timeouts**: per-call budget, jitter, retry classes (transient vs permanent), circuit breakers where the dependency justifies them.
- **Idempotency**: command id keys, deduplication windows, consumer-side dedup tables.
- **Observability**: log lines, metrics, traces, alerts, RED/USE per service-level objective.
- **Testing approach**: unit boundary, integration with real dependencies, contract tests against the API/message shape, migration dry-runs.
- **Migration approach**: backward-compatible schema, deprecation path, rollback plan.

## L4 — Implementation

Now write the smallest safe change.

- **Handlers**: thin transport, validate input, call use case.
- **Services / use cases**: orchestrate domain logic, hold the invariant, call repositories.
- **Repositories / adapters**: persistence and external IO; one place per dependency.
- **Schemas / migrations**: additive when possible; expand-then-contract for breaking changes.
- **Event consumers**: idempotent, dead-letter on poison, log correlation id.
- **Tests**: regression test for the bug or invariant; contract test if the boundary changed.
- **Docs**: update OpenAPI / proto / event-schema and any developer-facing readme.

## Evaluation axes

For every L2/L3 decision, name where it lands on at least two of these four axes:

- correctness ↔ latency
- consistency ↔ availability
- coupling ↔ autonomy
- simplicity ↔ operability

Stating the trade-off keeps the design honest. A choice that claims to win on every axis is usually under-analyzed.

## Fast path

For isolated compile errors, narrow test fixes, or one-line query fixes, compress L1–L3 into one sentence then jump to L4. State `L1–L3 skipped: isolated fix` so traceability stays explicit.
