# Architecture

## Contents
- Reference architecture (layered structure)
- Layer responsibilities
- Event-driven rules
- Sync vs async heuristic
- Event schema evolution
- Aggregate & data boundaries
- Anti-patterns

## Reference architecture

```
API Gateway
    ↓
Orchestrator (Process Manager)
    ↓
Domain Processor (Aggregate)
    ↓
Adapter (External Systems)
```

## Layer responsibilities

- **Orchestrator**: controls flow only. Contains no business logic.
- **Domain Processor**: owns business rules, emits events, owns state transitions.
- **Adapter**: protocol translation only. Carries no domain meaning.

## Event-driven rules

- Prefer event chaining over direct coupling.
- Sync + async mix is allowed — the choice MUST be explicit per integration.
- Every event MUST have: clear ownership, clear meaning, idempotency strategy.

## Sync vs async heuristic

Prefer **sync** when:
- Strong consistency is required (caller must see the result before the next action).
- Fan-out is low (1–2 downstream consumers).
- End-to-end latency budget is under 100ms.
- Failure must surface to the caller immediately.

Otherwise prefer **async** (events).

## Event schema evolution

- Additive by default: new fields are optional, never required.
- NEVER rename or remove a field in place — deprecate first, remove in a later major version.
- Version events explicitly (`loan.application.submitted.v1`, `.v2`).
- Use a schema registry (Confluent, AWS Glue) for any event crossing team boundaries.
- Consumers MUST be tolerant readers: ignore unknown fields, assume no field order.

## Aggregate & data boundaries

- Strict data ownership per domain — no shared mutable tables across domains.
- Cross-domain interaction goes through events (preferred) or a controlled API.
- One aggregate equals one transactional boundary.

## Anti-patterns (never)

- Over-engineering without proven scale requirement.
- Forcing async everywhere.
- Shared database across domain boundaries.
- Business logic hidden inside an orchestrator.
- Blind retries without state tracking.
