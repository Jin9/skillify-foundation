# The four principles of an effective handoff protocol

These four principles distinguish effective from fragile handoff
protocols. The skill enforces all four. They are cited verbatim from the
synthesis in the agent-handoff-protocols research report.

## 1. Ownership is binary

> **Ownership is binary.** At any moment exactly one agent may mutate the
> shared state. Frameworks that allow concurrent mutation without explicit
> ownership transfer — or that lack clear termination criteria — accumulate
> subtle corruption that manifests as unexplained contradictions downstream.

Engineering consequences:

- Mark every shared-state element. Assign exactly one writing agent at
  any moment; all other agents are read-only consumers.
- Ownership transfer is an explicit, point-in-time event — the handoff
  itself. The originating agent ceases to be responsible once it hands
  off; the receiving agent takes over with the accumulated context.
- Define explicit termination criteria. Role ambiguity with no agent
  owning a terminal condition produces infinite wait loops.
- Concurrent writers are forbidden: "allowing two agents to mutate the
  same context concurrently produces state corruption that cascades
  forward through the pipeline." Blackboard / shared-graph designs must
  implement synchronization or they suffer race conditions and stale
  reads (a documented state-synchronization anti-pattern).
- Cognition's narrowed position lands here: writes stay single-threaded;
  additional agents may contribute intelligence (review, critique) but
  not act.

## 2. Handoff payloads are contracts, not notes

> **Handoff payloads are contracts, not notes.** Effective teams treat
> inter-agent transfers like public APIs: typed fields, versioned schemas
> (include `schemaVersion` in the payload), JSON Schema validation at
> generation time, and repair-on-failure semantics (validate with Pydantic
> or Guardrails; on failure, run a repair prompt and retry). Free-text
> handoff notes are the dominant source of context loss in practice.

Engineering consequences:

- Required typed fields: task ID, intent, current state, confidence
  signal, provenance/trace ID — plus `schemaVersion` in the payload.
- Validate the produced payload at generation time against the JSON
  Schema. On failure: run a repair prompt and retry before the handoff
  is allowed to proceed. Never let an unvalidated payload cross the
  boundary.
- A handoff "should be structured data, not a long narrative." Free-text
  notes are debt that accumulates with every additional agent.
- This is structurally aligned with MCP's JSON-RPC 2.0 typed messages
  and OpenAI Agents SDK explicit schema validation for the handoff
  payload.

## 3. Observability must span agent boundaries

> **Observability must span agent boundaries.** Single-agent tracing is
> insufficient. Effective handoff observability captures nested spans
> preserving parent-child relationships across agent transitions, enables
> correlation IDs to reconstruct the full execution graph, and feeds
> production traces back into the eval suite so that every regression
> becomes a test.

Engineering consequences:

- Emit nested spans whose parent-child relationships survive every agent
  transition. The `provenance` object carries `traceId` (constant across
  the whole graph), `spanId` (this handoff), and optionally
  `parentSpanId`.
- A single correlation ID must let you reconstruct the full execution
  graph. In multi-agent debugging you read N parallel transcripts plus
  the orchestrator's reasoning — without correlation IDs the
  inter-agent-misalignment failures are essentially unreachable.
- This skill defines the hooks and the correlation-ID field; the actual
  telemetry pipeline (OTEL exporters, backends) is the responsibility of
  `observability-telemetry-instrumenter`.

## 4. Autonomy is a spectrum calibrated to reversibility

> **Autonomy is a spectrum, not a switch.** Effective protocols assign a
> governance level per handoff based on task reversibility and risk: fully
> supervised (human approves every action) for schema changes or data
> deletion; conditional autonomy within defined bounds for code
> deployments; higher autonomy for reversible low-risk tasks.

The four-level model (assign one level per handoff):

- **Level 1 — human approves every action.** Irreversible changes:
  schema changes, data deletion, production disbursement.
- **Level 2 — conditional autonomy with escalation.** Bounded risk:
  code deployments within defined bounds; escalate on threshold breach.
- **Level 3 — supervised with audit.** Routine tasks; actions logged
  and auditable after the fact.
- **Level 4 — full autonomy.** Reversible low-stakes tasks only.

Teams that map their SDLC control-transfer points before scaling agent
volume adopt successfully 3x more often than those that attempt ad hoc
deployment. The autonomy level travels with the contract in the optional
`autonomyLevel` field; the *policy* deciding what each level may do is
governed by `governance-policy-generator`, not this skill.
