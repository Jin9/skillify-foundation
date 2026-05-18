# Handoff design — &lt;task / system name&gt;

> Skeleton. Replace every bracketed slot. Delete guidance lines before delivery.

## 1. Topology

- Chosen topology: [orchestrator-worker (default) | hierarchical/supervisor | other — justify]
- Agents / roles: [list each role]
- Decision-record reference: see the companion topology-decision-record.

## 2. Binary write-ownership map

Exactly one writer per shared-state element at any moment. All others
are read-only.

| Shared state element | Sole writer (per phase) | Readers | Ownership-transfer point | Terminal owner |
|---|---|---|---|---|
| [e.g. plan] | [planner] | [coder, reviewer] | [planner -> coder on intent=implement] | [reviewer] |
| [e.g. codebase diff] | [coder] | [reviewer, tester] | [coder -> reviewer on intent=review] | [reviewer] |

Concurrent-writer check: [confirm no row allows two simultaneous writers]

## 3. The contract

- Canonical schema: `schemas/handoff-contract.schema.json` (JSON Schema
  draft 2020-12, `additionalProperties: false`).
- Required fields present: taskId, intent, state, confidence,
  provenance, schemaVersion. No free-text note/narrative/prose field.
- Per-edge payload notes:

| Handoff edge | intent enum used | state.summary content (fielded) | confidence gate |
|---|---|---|---|
| [planner -> coder] | implement | [scope, constraints, decisions] | [escalate if &lt; 0.6] |

## 4. Generation-time validation and repair strategy

- Validator: [Pydantic | Guardrails | JSON Schema validator]
- On validation failure: [run repair prompt -> retry -> re-validate; max
  N attempts then escalate]
- Versioning: `schemaVersion` = [e.g. 1.0.0]; receivers reject
  incompatible major.

## 5. Cross-boundary observability hooks

- Correlation ID: `provenance.traceId` constant across the whole graph.
- Span nesting: `provenance.spanId` child of producing agent's span;
  `provenance.parentSpanId` preserves parent-child across transitions.
- Hook points: [emit span on ownership transfer; on validation failure;
  on confidence-gate breach]
- Note: actual telemetry export/backends are out of scope here — defer
  to `observability-telemetry-instrumenter`.

## 6. Autonomy per handoff

| Handoff edge | Reversibility | Autonomy level (1–4) | Rationale |
|---|---|---|---|
| [coder -> deploy] | [hard to reverse] | [1 or 2] | [human approves irreversible] |
