# OTel span / metric definition — <agent / service name>

Fill in and emit alongside a JSON instance conforming to
`schemas/otel-genai-spans.schema.json`.

## semconv opt-in

- `OTEL_SEMCONV_STABILITY_OPT_IN`: `gen_ai_latest_experimental`  (required)
- Status: Development (semantic-conventions v1.36.0 origin; breaking by v1.37.0)

## Spans

### invoke_agent (root)

- Span name: `invoke_agent {gen_ai.agent.name}`
- Span kind: **CLIENT**
- Role: root (agent-invocation boundary)
- Required attributes:
  - `gen_ai.operation.name`  = `invoke_agent`
  - `gen_ai.provider.name`   = `<anthropic | openai | gcp.gemini | gcp.vertex_ai>`
  - `gen_ai.usage.input_tokens`   (billable; includes cached tokens)
  - `gen_ai.usage.output_tokens`  (billable; from provider `usage` block)
- Conditional: `gen_ai.agent.name`, `gen_ai.agent.id`, `gen_ai.conversation.id`
- High-cardinality attrs (span only, never metric labels): `user_id`,
  `request_id`, `session_id`

### execute_tool (child)

- Span name: `execute_tool`
- Span kind: **INTERNAL**
- Parent: `invoke_agent`  (skip if a more specific instrumentation, e.g. MCP,
  already covers the tool)
- Carry: latency, outcome; `gen_ai.tool.name` on the parent CLIENT span

## Metric: gen_ai.client.token.usage

- Instrument: **histogram**
- Unit: `{token}`
- Explicit bucket boundaries (exact, do not alter):
  `[1, 4, 16, 64, 256, 1024, 4096, 16384, 65536, 262144, 1048576, 4194304, 16777216, 67108864]`
- Split by: `gen_ai.token.type` in {input, output}
- Reporting rule: report **billable** tokens, not used tokens
- Allowed labels: `{ model, operation, provider, status }`  — and nothing else
- Forbidden labels: `user_id`, `request_id`  (→ span attrs + exemplars)
- Exemplars: attach `(trace_id, span_id)` so a hot bucket resolves to a trace

## Raw payload capture

- `optIn`: `false`  (default OFF; never automatic)
- `maxChars`: `<= 1000`  (cap 500–1000)
- Storage: span *events*; full payload → object storage; pointer on span
- Secrets/PII redacted before telemetry

## Sampling

- Tail sampling: keep all error traces + ~10% successful; tune over time
- Guard against never-completing traces (unbounded agent loops)
