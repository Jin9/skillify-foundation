# Contract field specification

The handoff payload is a public API. It carries six required fields plus
`schemaVersion`. Free-text `notes`/`narrative`/`prose` fields are
forbidden. The canonical schema is
`schemas/handoff-contract.schema.json` (JSON Schema draft 2020-12,
`additionalProperties: false`).

## Required fields and their semantics

- **`taskId`** (string) — stable identifier for the unit of work.
  Constant across every handoff for the same task, so the full chain is
  joinable. Maps to the "task ID" required field for handoff payloads
  designed as APIs.
- **`intent`** (enum string) — the receiving agent's next required
  action as a discrete, enumerable directive (plan, implement, review,
  test, deploy, investigate, escalate, finalize). The "detected intent"
  field of a structured context object. Never a free-text instruction.
- **`state`** (object) — the structured *current state*: `status`,
  a bounded fielded `summary`, lightweight `artifacts` references, and
  explicit `decisions` already made (so downstream agents do not
  silently re-decide them — the Flappy-Bird implicit-decision failure).
  This is a typed structured context object, not a transcript.
- **`confidence`** (number, 0–1) — the confidence signal. Low-confidence
  outputs must be surfaced before the handoff proceeds; cascading-error
  mitigation depends on agents flagging low confidence before passing
  work forward.
- **`provenance`** (object) — `traceId` (correlation ID constant across
  the whole execution graph), `spanId` (this handoff), `sourceAgent`,
  and optional `parentSpanId` to preserve parent-child span nesting
  across agent transitions. This is the provenance/trace ID required
  field, and the substrate for boundary-spanning observability.
- **`schemaVersion`** (semver string, e.g. `1.0.0`) — version of the
  contract carried *in the payload*.

## Versioning and repair-on-failure

- Version the payload with `schemaVersion`. A breaking field change
  requires a major bump; receivers reject an incompatible major rather
  than guess.
- Validate the produced payload at generation time against the JSON
  Schema (Pydantic / Guardrails / a JSON Schema validator). On failure:
  run a repair prompt and retry, then re-validate. An unvalidated
  payload never crosses the boundary. This repair-on-failure loop is
  mandatory, not optional.
- Validation is generation-time, not consumption-time only: catching the
  malformed payload before it is handed off prevents the receiving agent
  from reconstructing context from a broken contract.

## Structured-context-object sizing

- Target the `state` object at **200–500 tokens**. Full-context
  forwarding runs **5,000–20,000 tokens** and scales quadratically with
  handoff depth — explicitly rejected as a default.
- Subagents return the minimum needed: a ~200-token summary from a
  5,000-token read, with the read-tokens discarded when the subagent
  finishes. Persist large work products externally and pass lightweight
  `artifacts` references, not inlined content. (Token *instrumentation*
  and budget telemetry are owned by
  `observability-telemetry-instrumenter`; this skill only sizes the
  contract.)

## MCP / JSON-RPC typed-message alignment

The contract is deliberately aligned with MCP's JSON-RPC 2.0 typed
messages and the OpenAI Agents SDK's explicit schema validation for the
handoff payload: schema-validated communication (typed messages plus
explicit termination criteria and ownership boundaries) is the
structural mitigation for role ambiguity. Designing the handoff as a
JSON Schema contract keeps it interoperable across frameworks rather
than locked to one vendor's handoff idiom.
