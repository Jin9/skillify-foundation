# Token metric spec

## The metric

The headline is the metric `gen_ai.client.token.usage`: a **histogram**
instrument measured in units of `{token}` that records the number of input and
output tokens. It requires the attributes `gen_ai.operation.name`,
`gen_ai.provider.name`, and `gen_ai.token.type`.

## The exact 14 bucket boundaries

`gen_ai.client.token.usage` is defined as a histogram with explicit bucket
boundaries:

```
[1, 4, 16, 64, 256, 1024, 4096, 16384, 65536, 262144, 1048576, 4194304, 16777216, 67108864]
```

This is a 14-bucket roughly-geometric ladder covering single-token outputs to
multi-million-token long-context calls — 1 token to ~67M tokens. Use exactly
these boundaries; do not round, extend, or trim. `scripts/span_attr_check.py`
fails any spec whose `bucketBoundaries` is not exactly this array.

## input vs output split

`gen_ai.token.type` is how input and output are kept apart. It has two
well-known values:

- `input` — "prompt, input, etc."
- `output` — "completion, response, etc."

Alongside the metric, individual spans carry `gen_ai.usage.input_tokens` and
`gen_ai.usage.output_tokens`.

## Normative rules that carry production weight

- **Billable, not used.** When a system reports both raw used tokens and
  billable tokens, instrumentation **MUST** report billable tokens, so cost
  numbers stay consistent across tools. Per-step telemetry MUST use
  `usage.output_tokens` from the response, not a re-counted visible-text
  length (Anthropic Claude 4 bills full thinking but returns a summary —
  re-counting prose understates the bill by 5–20x).
- **Input includes cached tokens.** `gen_ai.usage.input_tokens` SHOULD include
  all input token types, *including cached tokens*. A dashboard built on
  OpenInference `llm.token_count.prompt_details.cache_read` does not
  automatically transfer to one built on `gen_ai.usage.input_tokens`, where
  cached reads are folded in.
- **Pre-flight is not billing.** A pre-flight `count_tokens` call (e.g.
  Anthropic `POST /v1/messages/count_tokens`) is a useful budgeting input but
  is **not** a substitute for the post-response usage block; the billed count
  incorporates provider-side template overhead the estimate may not mirror.
  Tool definitions are a silent fixed input-token cost (~200–500 tokens per
  provider per request; ~346 cross-provider average) — they land in
  `gen_ai.usage.input_tokens`, not a separate dimension.

## Low overhead

OpenTelemetry overhead is reported at under 1ms per call while LLM API latency
of 100ms to 30s dominates entirely, and traces flow asynchronously, so this
instrumentation can safely run in production. Without per-span token
accounting, an agent burning 50,000 tokens on a task that should take 3,000 is
invisible.
