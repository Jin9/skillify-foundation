---
name: observability-telemetry-instrumenter
description: >
  Instrument agent code with OpenTelemetry GenAI semantic conventions: an
  invoke_agent CLIENT root span, execute_tool INTERNAL child spans, the
  gen_ai.client.token.usage histogram with the exact 14 bucket boundaries,
  low-cardinality metric labels, opt-in raw prompt/completion capture, payload
  caps, and experimental-status gating. Use when the user asks to "add
  OpenTelemetry to this agent", "instrument per-step token usage", "emit
  gen_ai spans and the token histogram", "wire up invoke_agent and
  execute_tool spans", or "make agent spend observable". Output: an
  instrumentation spec object plus an instrumentation plan and a per-provider
  token-accounting reconciliation skeleton. Do NOT use for the inter-agent
  handoff contract (multi-agent-handoff-architect), credential scrubbing in
  the pipeline (devops-infrastructure-hardener), governance policy
  (governance-policy-generator), generic non-AI APM, feature/service
  observability design (observability-design), or after-the-fact spend
  review (reviewing-agent-spend).
---

# Observability Telemetry Instrumenter

## Purpose

Produce the instrumentation contract that makes an agent's execution and
per-step token spend visible under the OpenTelemetry GenAI semantic
conventions: an `invoke_agent` CLIENT root span, `execute_tool` INTERNAL child
spans, and the `gen_ai.client.token.usage` histogram with its exact 14 bucket
boundaries. The deliverable hard-enforces low-cardinality metric labels,
opt-in raw payload capture, payload caps with object-storage pointers, and
experimental-status gating. This instruments the telemetry; it does not review
spend, scrub credentials, or define governance.

## When to use this skill

- Use when: "add OpenTelemetry to this agent" / "instrument per-step token usage".
- Use when: "emit gen_ai spans and the token histogram" / "wire up invoke_agent and execute_tool spans".
- Use when: "make agent spend observable" / "set up the GenAI semantic conventions telemetry".
- Do NOT use when: the ask is the inter-agent handoff contract, credential scrubbing in the pipeline, governance policy, generic non-AI APM, or after-the-fact spend review. Hand those to the appropriate skill.

## Inputs

Read access to the agent's code or framework (the LLM client call site, the
agent loop, each tool entry point) and the target provider(s) — Anthropic,
OpenAI, Gemini/Vertex. Optionally an existing OTel SDK init and collector
config. None are generated here; only the instrumentation spec and the
planning skeletons are produced.

## Workflow

1. **Scope and map the hierarchy.** Identify the agent-invocation boundary, the per-turn exchange, and each per-call and per-tool-call site. Fix the nesting: `execute_tool` (per-call) under the per-turn exchange under the per-agent `invoke_agent` root. Entry: code/framework read access. Exit: a span map. See `references/span-hierarchy.md`.
2. **Define spans.** Require an `invoke_agent {agent.name}` root span of kind CLIENT and `execute_tool` child spans of kind INTERNAL. Require attributes `gen_ai.operation.name`, `gen_ai.provider.name`, `gen_ai.usage.input_tokens`, `gen_ai.usage.output_tokens`. Propagate W3C trace context across services. See `references/span-hierarchy.md` and `templates/otel-span-metric-spec.md`.
3. **Define the token metric.** Require the `gen_ai.client.token.usage` histogram in units `{token}` with the EXACT 14 explicit bucket boundaries `[1, 4, 16, 64, 256, 1024, 4096, 16384, 65536, 262144, 1048576, 4194304, 16777216, 67108864]`, split by `gen_ai.token.type` in {input, output}. Report billable tokens, never used tokens; input includes cached tokens. See `references/token-metric-spec.md`.
4. **Set cardinality and privacy gates.** Restrict metric labels to model / operation / provider / status. Put `user_id`, `request_id` on span attributes and histogram exemplars, NEVER on metric labels. Make raw prompt/completion capture opt-in only (default off); cap captured content at 500–1000 chars with an object-storage pointer; redact secrets before telemetry. See `references/cardinality-and-privacy.md`.
5. **Gate experimental conventions and reconcile providers.** Require `OTEL_SEMCONV_STABILITY_OPT_IN=gen_ai_latest_experimental`; record that conventions began at semantic-conventions v1.36.0 with breaking changes by v1.37.0 (Development status). Populate the homogenized 5–8-column per-provider usage table and its trap callouts. See `references/cross-provider-traps.md` and `templates/token-accounting-reconciliation.md`.
6. **Self-check, then emit.** Run `scripts/span_attr_check.py <spec-file>`: it PASSES only if the histogram `bucketBoundaries` exactly equals the 14-element array, no metric label set contains `user_id`/`request_id`, and raw-payload `optIn` defaults to false. Emit the spec and the two planning skeletons only after PASS. The spec conforms to `schemas/otel-genai-spans.schema.json`.

Do not silently emit whichever SemConv version the instrumentation already
emits; the opt-in is explicit. Do not re-count visible response text as
output tokens; use the provider `usage` block (billable).

## Output contract

Three artifacts shaped by `templates/`:

- **Instrumentation spec** — a JSON object conforming to `schemas/otel-genai-spans.schema.json`: the `invoke_agent` (CLIENT, root) and `execute_tool` (INTERNAL, child) span definitions, required `gen_ai.*` attributes, the `tokenUsageHistogram` with the 14 bucket boundaries, the restricted metric label set, and `rawPayloadCapture` with `optIn: false` default and `maxChars` ≤ 1000.
- **Instrumentation plan** — the ordered SDK-init-first → auto-instrument → wrap spans → metric/sampling/privacy checklist (`templates/instrumentation-plan.md`, with `templates/otel-span-metric-spec.md`).
- **Token-accounting reconciliation** — the homogenized 5–8-column per-provider usage table with the raw-usage-JSON-retained note and per-provider trap callouts (`templates/token-accounting-reconciliation.md`).

No credential scrubbing, governance policy, handoff contract, or spend review is produced.

## Constraints

- DO NOT emit a metric label outside {model, operation, provider, status}; high-cardinality `user_id`/`request_id` go on span attributes and histogram exemplars only.
- MUST use the EXACT 14 bucket boundaries `[1, 4, 16, 64, 256, 1024, 4096, 16384, 65536, 262144, 1048576, 4194304, 16777216, 67108864]`; any deviation fails the self-check.
- MUST report billable tokens, not used tokens; `gen_ai.usage.input_tokens` includes cached tokens; `count_tokens` is a pre-flight estimate only.
- NEVER capture raw prompts/completions automatically — opt-in only, default off, capped at 500–1000 chars with an object-storage pointer, secrets redacted first.
- MUST gate experimental conventions behind `OTEL_SEMCONV_STABILITY_OPT_IN=gen_ai_latest_experimental` (Development status; v1.36.0 origin, breaking by v1.37.0).
- DO NOT duplicate the methodology here; it lives one level deep in `references/`.

## Validation

- [ ] Frontmatter `name` equals folder, kebab-case, no XML, description under 1024 chars with triggers + negatives.
- [ ] Workflow requires the `invoke_agent` CLIENT root / `execute_tool` INTERNAL child hierarchy, the exact 14-bucket histogram, opt-in capped capture, and the experimental opt-in, and ends with a `scripts/span_attr_check.py` self-check.
- [ ] Output contract names the spec, plan, and reconciliation artifacts and the no-scrubbing / no-governance / no-handoff / no-spend-review boundary.

## References

- Span hierarchy & W3C context: `references/span-hierarchy.md`
- Token metric & billable-tokens rule: `references/token-metric-spec.md`
- Cardinality & privacy gates: `references/cardinality-and-privacy.md`
- Cross-provider accounting traps & experimental gating: `references/cross-provider-traps.md`
- Schema: `schemas/otel-genai-spans.schema.json`
- Skeletons: `templates/instrumentation-plan.md`, `templates/otel-span-metric-spec.md`, `templates/token-accounting-reconciliation.md`; self-check: `scripts/span_attr_check.py`
