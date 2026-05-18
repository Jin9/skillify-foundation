# Cross-provider accounting traps

Three providers, three usage shapes, three divergence points. Anthropic
surfaces cache reads as a first-class field; OpenAI buries reasoning tokens
*inside* `completion_tokens`; Gemini's `candidatesTokenCount` flips meaning
between the Gemini API and Vertex AI. Homogenizing those into a single
internal row is lossy unless the schema is designed for it.

## Anthropic: summarized vs billed thinking

Every `messages.create` response carries
`usage { input_tokens, output_tokens, cache_creation_input_tokens,
cache_read_input_tokens }`, with `cache_creation` decomposed into
`ephemeral_5m_input_tokens` / `ephemeral_1h_input_tokens`.

The accounting trap is severe: on Claude 4 you are billed for the **full**
thinking process but the API returns only a **summary**. `budget_tokens:
10000` may show ~500 visible tokens and bill ~10,000 output — telemetry that
re-counts the response prose understates the bill by 5–20x. Per-step
telemetry MUST use `usage.output_tokens` from the response, not a re-counted
visible-text length.

## The March 2026 default-TTL change

Anthropic silently changed the Claude Code default prompt-cache TTL from **1h
to 5m** on March 6 2026 with no blog post or changelog entry. Without explicit
`cache_control`, cache-hit rates collapsed to near zero. Recovery:
`"anthropic-beta": "prompt-caching-2024-07-31"` header or
`cache_control { type: "ephemeral", ttl: 3600 }` per message. The cache-hit
ratio `cache_read_input_tokens / (cache_read_input_tokens + input_tokens)` is
the canonical regression alert for upstream TTL or eviction changes.

## OpenAI: reasoning folded inside completion

`usage` is the parent of `prompt_tokens_details { cached_tokens, audio_tokens }`
and `completion_tokens_details { reasoning_tokens, accepted_prediction_tokens,
rejected_prediction_tokens, audio_tokens }`. The frequently-misread fact:
`reasoning_tokens` is **not** a separate counter — it is a subcomponent of
`completion_tokens`, "counted in the total completion tokens for purposes of
billing, output, and context window limits." A naive `completion_tokens` sum
already includes reasoning; the details object only gives the breakdown.

## Gemini API vs Vertex AI divergence

`usageMetadata` is `{ promptTokenCount, candidatesTokenCount,
cachedContentTokenCount, thoughtsTokenCount }`. The single largest
cross-provider gotcha: on the Gemini API
(`generativelanguage.googleapis.com`), `candidatesTokenCount` already
**includes** `thoughtsTokenCount`; on Vertex AI (`vertex.googleapis.com`), it
**excludes** it. Any homogenization layer that doesn't tag the platform before
summing will silently overcount on the Gemini API path or undercount on the
Vertex path.

## Experimental-status gating

The OpenTelemetry GenAI conventions are still officially **Development**
status — not Stable. They began at semantic-conventions **v1.36.0** and were
new enough that **breaking changes already occurred by v1.37.0** (cache-token
attributes, an Evaluation Event, reasoning-content message parts added). The
project manages breakage through the environment variable
`OTEL_SEMCONV_STABILITY_OPT_IN=gen_ai_latest_experimental`. The default emits
whichever version the instrumentation already emits, so instrumentation
pinned to a SemConv version needs an explicit opt-in to roll forward, or
dashboards keyed off `gen_ai.*` attribute names will break. Always set the
opt-in explicitly.

## The homogenized 5–8-column usage table

"Homogenized usage" is a five-to-eight-column thing, not two — at minimum
`{ input_tokens_billable, output_tokens_billable, cache_read_tokens,
cache_write_tokens_short_ttl, cache_write_tokens_long_ttl,
reasoning_or_thinking_tokens, tool_use_prompt_tokens, audio_tokens }`.
Lossiness in the homogenization is acceptable **only if the raw `usage` JSON
is also stored** for later reconciliation. Keep the reasoning/thinking series
separately addressable per provider: OpenAI's reasoning is hidden and
non-round-tripping, Anthropic's is summarized-and-shorter-than-billed,
Gemini's is non-billable-separately but rate-tier-changing.
