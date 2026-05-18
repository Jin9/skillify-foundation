# Cardinality and privacy

## The cardinality wall

Agent systems generate near-unique values constantly — every thread, tool
result, or retrieved document looks distinct — and high cardinality in span
names or attributes causes storage bloat and slow queries. Prometheus's hard
constraint is that high-cardinality labels (`user_id`, `request_id`,
`session_id`) destroy series count.

The canonical footgun: adding a `user_id` label to a single request counter
produced **over 5,000,000 time series across 1,000,000 users in 90 days**,
causing severe performance degradation. Production guidance explicitly forbids
`user_id` on metric labels.

## Low-cardinality metric labels vs high-cardinality span attributes

- **Metric labels (low cardinality, allowed):** a `gen_ai.client.token.usage`
  histogram labeled only by `gen_ai.request.model` /
  `gen_ai.operation.name` / `gen_ai.provider.name` / `error.type` (status) is
  the production-safe shape — `{model, operation, provider, status}`.
- **Span attributes + histogram exemplars (high cardinality):** `user_id`,
  `request_id`, and everything more specific belong on spans. OTel exemplars
  attach a `(trace_id, span_id)` to a metric data point — for histograms, the
  bucket the value lands in. When a token-usage histogram bucket lights up
  unexpectedly, the operator clicks an exemplar and lands on a representative
  trace. Without exemplars the histogram says *something* spiked but not
  *which step*; with exemplars the path is one click.

This skill REQUIRES that no metric label set contains `user_id` or
`request_id`; the self-check fails otherwise.

## Opt-in raw payload capture

OpenTelemetry explicitly recommends that sensitive GenAI payloads — full
instructions, inputs, and outputs — be **opt-in** rather than captured
automatically. Capturing prompts and completions must be opt-in (for example
an `ENABLE_SENSITIVE_DATA`-style flag), typically enabled in development for
debugging and disabled in production. Default is OFF; capture is never
automatic.

## Payload caps and object-storage pointers

Prompts and completions can be huge, so cap captured content at **500–1000
characters** and, in production, store the full payload in external object
storage with only a pointer on the span. Store prompt/completion content as
span *events* rather than attributes, because events can be filtered or
dropped at the Collector without touching application code. Redact secrets and
personal data before they ever reach telemetry.

## The never-completing-trace / tail-sampling hazard

An agent loop can keep generating spans indefinitely, which breaks tail
sampling — the Collector normally waits for a trace to "complete" before
deciding whether to keep it, and an unbounded trace forces it to buffer every
span in memory. Manage volume with tail sampling: a sensible starting policy
keeps all error traces and a 10% sample of successful ones, then is tuned as
the team learns what is useful. Every powerful default (capture everything,
keep every span) has a cost the team should choose, not inherit.
