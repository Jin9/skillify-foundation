# Instrumentation plan — <agent / service name>

Ordered checklist. The reported good path: initialize the OpenTelemetry SDK
as early as possible, auto-instrument the LLM client first, then wrap the
top-level agent call and each tool, then make the metric / sampling / privacy
decisions deliberately.

## 1. SDK init first

- [ ] Initialize the OpenTelemetry SDK as early as possible in the process
      (before the LLM client is constructed).
- [ ] Set `OTEL_SEMCONV_STABILITY_OPT_IN=gen_ai_latest_experimental`
      explicitly (conventions are Development status; v1.36.0 origin,
      breaking by v1.37.0).
- [ ] Configure the OTLP exporter / Collector endpoint.
- [ ] Enable W3C trace-context propagation across service hops and sub-agents.

## 2. Auto-instrument the LLM client

- [ ] Add auto-instrumentation for the LLM client first so token and model
      spans appear immediately.
- [ ] Confirm `gen_ai.usage.input_tokens` / `gen_ai.usage.output_tokens`
      populate from the provider `usage` block (billable, not re-counted).

## 3. Wrap invoke_agent and each execute_tool step

- [ ] Wrap the top-level agent call in an `invoke_agent {agent.name}` span,
      kind **CLIENT**, root.
- [ ] Wrap each tool in an `execute_tool` span, kind **INTERNAL**, child of
      `invoke_agent` (skip tools already traced by a more specific
      instrumentation, e.g. MCP).
- [ ] Verify nesting: per-call under per-turn under per-agent.
- [ ] Required attributes present: `gen_ai.operation.name`,
      `gen_ai.provider.name`, `gen_ai.usage.input_tokens`,
      `gen_ai.usage.output_tokens`.

## 4. Metric + sampling + privacy policy

- [ ] Emit the `gen_ai.client.token.usage` histogram, unit `{token}`, with
      the exact 14 bucket boundaries; split by `gen_ai.token.type`
      in {input, output}.
- [ ] Metric labels restricted to {model, operation, provider, status}.
- [ ] `user_id` / `request_id` on span attributes + histogram exemplars only,
      NEVER on metric labels.
- [ ] Raw prompt/completion capture opt-in (default OFF); store as span
      events; cap at 500–1000 chars; full payload to object storage with a
      span pointer; redact secrets first.
- [ ] Tail sampling: keep all error traces + ~10% of successful ones; tune.
- [ ] Run `scripts/span_attr_check.py <spec-file>` and confirm PASS.
