# Token-accounting reconciliation — <agent / service name>

Homogenized usage is a 5-to-8-column thing, not two. Fill one row per
per-call usage block. Lossiness is acceptable ONLY because the raw `usage`
JSON is retained for later reconciliation.

## Homogenized per-provider usage table

| step / call | provider / platform | input (billable) | output (billable) | cache-read | cache-write short-TTL (5m) | cache-write long-TTL (1h) | reasoning / thinking | tool-use prompt | audio |
|---|---|---|---|---|---|---|---|---|---|
| <call 1> | anthropic | | | | | | | | |
| <call 2> | openai | | | | n/a | n/a | | | |
| <call 3> | gemini-api / vertex | | | | | | | | |

- input billable = Anthropic `input_tokens` / OpenAI `prompt_tokens` /
  Gemini `promptTokenCount`  (includes cached tokens per the OTel rule)
- output billable = Anthropic `output_tokens` / OpenAI `completion_tokens` /
  Gemini `candidatesTokenCount`  (platform-conditional — see Gemini trap)
- cache-read = `cache_read_input_tokens` / `cached_tokens` / implicit
- cache-write short/long-TTL = Anthropic
  `ephemeral_5m_input_tokens` / `ephemeral_1h_input_tokens`

## Raw usage JSON retained

- [ ] The verbatim provider `usage` / `usageMetadata` JSON is stored per call
      (homogenization is lossy; raw is the durable reconciliation primitive,
      with a price-table version per step).

## Per-provider trap callouts

- **Anthropic — summarized vs billed thinking.** `budget_tokens: 10000` may
  show ~500 visible tokens while billing ~10,000 output. Use
  `usage.output_tokens`, never re-counted visible text (5–20x understatement).
- **Anthropic — March 2026 TTL change.** Default prompt-cache TTL silently
  moved 1h → 5m (Mar 6 2026). Recover with `anthropic-beta:
  prompt-caching-2024-07-31` header or `cache_control.ttl=3600`. Alert on
  cache-hit ratio `cache_read / (cache_read + input)`.
- **OpenAI — reasoning folded in.** `reasoning_tokens` is a subcomponent of
  `completion_tokens` (counted in total for billing/output/context); not a
  separate counter.
- **Gemini — API vs Vertex divergence.** `candidatesTokenCount` *includes*
  `thoughtsTokenCount` on the Gemini API but *excludes* it on Vertex AI. Tag
  the platform before summing or you over/undercount.
