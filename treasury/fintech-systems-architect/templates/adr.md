# ADR-NNNN: <decision title>

- **Status:** Proposed | Accepted | Deprecated | Superseded by ADR-XXXX
- **Date:** YYYY-MM-DD
- **Owners:** <product, domain, data, infra, security, compliance — one named owner per relevant track>
- **Stakeholders:** <who must approve, who must be informed>

## Context

What problem are we solving? What regulatory or financial invariant must hold? What forced the decision now (incident, growth, compliance, deprecation)? Cite the L1 Business Intent and the L2 bounded context affected. Keep this section evidence-led: link to incidents, metrics, regulatory clauses, or upstream ADRs rather than asserting "this is bad."

## Decision

State the chosen direction in one paragraph. Name the bounded contexts, the contracts (REST / gRPC / Kafka topic / event schema), the data ownership, the transactional and authorization boundaries, and the failure-handling posture.

## Alternatives considered

For each rejected option, name it and state the deciding trade-off in one or two sentences. At least two alternatives are required for any non-trivial decision.

- **Option A — <name>:** <one-line trade-off and reason rejected>
- **Option B — <name>:** <one-line trade-off and reason rejected>
- **Option C — <name>:** <one-line trade-off and reason rejected>

## Consequences

### Positive

- What does this enable, simplify, or de-risk?

### Negative

- What new complexity, cost, or operational burden does this add?

### Trade-off axes

State where the decision lands on at least 2 of:

- complexity ↔ maintainability
- cost ↔ performance
- speed ↔ safety
- coupling ↔ flexibility

### Regulatory and compliance impact

- PDPA / GDPR / PCI-DSS / BOT / OJK / MAS implications. Mark `n/a` only if no regulated data class is touched.
- Audit log changes (what new event must be logged, retained how long, in which sink).
- Data residency / cross-border transfer impact.

## Hand-off to `crafting-backend-code`

The implementation skill picks up from this section. Provide:

- **Bounded contexts and aggregates** (names, owners, invariants).
- **Contracts** (REST endpoints, gRPC services, Kafka topics with key + schema + compatibility).
- **Data ownership** (who writes the source of truth; subscribers).
- **Transactional boundary** (what is atomic, what is eventually consistent, what compensation handles partial failure).
- **Authorization boundary** (which actor can invoke; how identity is asserted and verified).
- **Idempotency keys and retry posture** (per command / consumer).
- **Observability requirements** (log lines, metrics, traces, alerts that must exist).
- **Migration strategy** (additive vs expand-then-contract; rollback plan).
- **Validation strategy** (integration test boundary, contract test, migration dry-run).
- **Open questions for L4** (anything the architect deliberately deferred).

## Validation and rollback

- **Validation:** how will we know the design works in PRD? Name the metric, dashboard, or test that confirms the invariant.
- **Rollback:** if the design fails, what is the safest reversal path? Name the feature flag, deployment, or dual-write window that enables it.

## References

- <upstream ADRs, incidents, RFCs, regulatory clauses, design docs>
