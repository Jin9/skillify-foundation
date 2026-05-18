---
name: multi-agent-handoff-architect
description: >
  Design inter-agent handoff APIs as versioned JSON Schema contracts
  (required taskId, intent, state, confidence, provenance/trace,
  schemaVersion) with binary single-writer ownership, observability
  spanning agent boundaries via correlation IDs, and autonomy calibrated
  to task reversibility, preventing the 17.2x error amplification seen in
  unstructured agent loops; forbids free-text narrative handoffs. Use
  when the user asks to "design an agent handoff contract / API", "define
  the payload between planner and coder agents", "make multi-agent
  handoffs reliable", or "choose single-agent vs multi-agent and the
  handoff topology". Output: a
  handoff-design doc, a topology decision record, and a JSON Schema. Do
  NOT use for repo-level AGENTS.md squad-role prose (use
  agent-context-initializer), handoff tracing/token instrumentation (use
  observability-telemetry-instrumenter), governing agent permissions (use
  governance-policy-generator), or single-pipeline stage composition (use
  composing-agent-pipelines).
---

# Multi Agent Handoff Architect

## Purpose

Design the transfer of control, context, and responsibility between AI
agents as a versioned API, not a note. Every handoff is a JSON Schema
contract carrying typed fields (task ID, intent, current state,
confidence signal, provenance/trace ID, `schemaVersion`), validated at
generation time with repair-on-failure; exactly one agent may mutate
shared state at any moment; observability spans agent boundaries via
correlation IDs; and autonomy is calibrated to task reversibility. This
prevents context degradation and the 17.2x error amplification seen in
unstructured agent loops, where multi-agent systems fail at 41–87% in
production when handoffs are ad hoc and under-specified. Free-text
narrative handoffs are forbidden — they are the dominant source of
context loss in practice.

## When to use this skill

- Use when: "design an agent handoff contract / API" / "define the payload between planner and coder agents".
- Use when: "make our multi-agent handoffs reliable" / "schema for inter-agent state transfer".
- Use when: "choose single-agent vs multi-agent and the handoff topology" / "stop free-text handoffs between agents".
- Do NOT use when: the ask is repo-level AGENTS.md squad-role prose (hand to `agent-context-initializer`), tracing/token instrumentation of the handoffs (hand to `observability-telemetry-instrumenter`), governing what agents are permitted to do (hand to `governance-policy-generator`), or composing one linear pipeline's stages (hand to `composing-agent-pipelines`). Hand those to the appropriate skill.

## Inputs

A description of the task and the candidate agents/roles, the shared
state surface, the tool surface (shared or partitioned), and any
reversibility/risk constraints on the actions involved. Read access to
existing handoff code or schemas if revising. No agents, runtime, or
telemetry pipeline are built here — only the contract, topology, and
ownership design are produced.

## Workflow

1. **Decide single vs. multi-agent first (entry: a task; exit: a topology decision).** Apply the four-step rule: if the task is a deterministic 90%+ flowchart, recommend a workflow, not agents; otherwise recommend ONE agent by default (tightly coupled, mostly sequential, global context, under 20 tools). Escalate to a squad only when the task is parallelizable, read-heavy, with independent sub-problems AND a hard verifier that rejects rather than blends outputs. Default the topology to orchestrator-worker. Record the rationale in `templates/topology-decision-record.md`. See `references/topology-decision-rule.md`.
2. **Map binary write-ownership.** For the chosen topology, mark every shared-state element and assign exactly one writing agent at any moment; all others are read-only consumers. Define explicit ownership-transfer points and terminal conditions. Concurrent writers are forbidden — they cause state corruption that cascades forward. See `references/four-principles.md`.
3. **Specify the contract fields.** For each handoff edge define the payload using the six required fields plus `schemaVersion`: `taskId`, `intent`, `state`, `confidence` (numeric 0–1), `provenance` (`traceId`/`spanId`/`sourceAgent`), `schemaVersion` (semver). No free-text `notes`/`narrative`/`prose` field is permitted. Size structured context objects at 200–500 tokens, not full forwarding (5,000–20,000). See `references/contract-field-spec.md`.
4. **Emit the JSON Schema and validation strategy.** Produce or reference `schemas/handoff-contract.schema.json` (draft 2020-12, `additionalProperties:false`, the six required fields). Specify generation-time validation with repair-on-failure: validate the produced payload, and on failure run a repair prompt and retry before the handoff proceeds. See `references/contract-field-spec.md`.
5. **Wire cross-boundary observability and autonomy.** Define nested spans preserving parent-child relationships across agent transitions and a correlation ID that reconstructs the full execution graph. Assign an autonomy level per handoff from the four-level reversibility model (Level 1 human-approves irreversible … Level 4 full autonomy for reversible low-stakes). Flag anti-patterns to avoid. See `references/handoff-anti-patterns.md` and `references/handoff-evidence.md`.
6. **Self-check, then emit.** Run `scripts/handoff_schema_validator.py` against the schema: it must PASS only if all six required fields are declared, `additionalProperties:false` is set, and no `note`/`narrative`/`freetext`/`prose` property exists. Confirm exactly one writer per state element and that no payload carries free text. Emit the design doc, the decision record, and the schema; recommend a human review of any Level 1–2 handoff.

## Output contract

Three artifacts (shaped by `templates/` and `schemas/`):

- **Handoff design doc** (`templates/handoff-design.md`) — chosen topology, binary write-ownership map, the contract (pointer to `schemas/handoff-contract.schema.json`), generation-time validation/repair strategy, and cross-boundary observability hooks.
- **Topology decision record** (`templates/topology-decision-record.md`) — single-vs-multi rationale, the four-step gate checklist, chosen topology, and the justified boundary.
- **JSON Schema** (`schemas/handoff-contract.schema.json`) — draft 2020-12, required `taskId`/`intent`/`state`/`confidence`/`provenance`/`schemaVersion`, `additionalProperties:false`, no free-text field.

No agents, runtime, telemetry pipeline, governance policy, or AGENTS.md prose are produced.

## Constraints

- DO NOT permit free-text or narrative handoff payloads — handoff payloads are contracts, not notes; free-text notes are the dominant source of context loss.
- MUST emit a JSON Schema contract with the six required fields plus `schemaVersion`, validated at generation time with repair-on-failure.
- MUST enforce binary ownership — at any moment exactly one agent may mutate shared state; all others are read-only.
- MUST recommend a single agent unless the task is parallelizable, read-heavy, with independent sub-problems AND a hard verifier that rejects rather than blends outputs; default multi-agent topology to orchestrator-worker.
- NEVER use full-context forwarding where a structured context object (200–500 tokens) suffices.
- DO NOT take on AGENTS.md squad-role prose, telemetry instrumentation, agent-permission governance, or single-pipeline stage composition — defer to the sibling skills named in the description.
- DO NOT duplicate the methodology here; it lives one level deep in `references/`.

## Validation

- [ ] Frontmatter `name` equals folder, kebab-case, no XML, description under 1024 chars with triggers + negatives and sibling cross-refs.
- [ ] Workflow runs single-vs-multi first, maps binary ownership, specifies the six-field versioned contract, and ends with a schema self-check via `scripts/handoff_schema_validator.py`.
- [ ] Output contract names all three artifacts, forbids free-text payloads, and states the no-runtime/no-governance boundary.

## References

- The four principles (ownership, contracts, observability, autonomy): `references/four-principles.md`
- Single-agent-first and topology decision rule: `references/topology-decision-rule.md`
- Contract field semantics, versioning, repair-on-failure, context sizing: `references/contract-field-spec.md`
- Handoff failure modes and anti-patterns: `references/handoff-anti-patterns.md`
- Carried-from-brief 17.2x note + verified corroborating statistics: `references/handoff-evidence.md`
- Skeletons: `templates/handoff-design.md`, `templates/topology-decision-record.md`; schema: `schemas/handoff-contract.schema.json`; schema self-check: `scripts/handoff_schema_validator.py`
