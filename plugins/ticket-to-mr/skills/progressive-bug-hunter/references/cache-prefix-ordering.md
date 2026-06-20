# Prompt-cache-preserving retrieval ordering

Distilled from research report *Prompt cache breakpoint architecture* (+ the
cache-layout finding in *Selective context retrieval*). A bug hunt is
**iterative** — each escalation re-sends the prompt. If volatile retrieval
output sits in the cached prefix, every iteration pays a full cache *write* and
gets no *read* hit, so iteration cost compounds. Correct ordering keeps each
hunt iteration cheap (cache read ≈ 0.1× input; up to 90% cost / 85% latency on
the stable prefix; independent measurement 41–80% cost, 13–31% TTFT).

## The rule: static-first, dynamic-last
Every byte upstream of the last cache breakpoint must be byte-identical across
iterations. Layout:
```
stable prefix  (cache-hit territory — set ONCE, reuse every iteration)
┌───────────────────────────────────────────────┐
│ system prompt + tool defs + stable repo-map    │
└───────────────────────────────────────────────┘
            ── cache_control breakpoint ──
┌───────────────────────────────────────────────┐
│ this iteration's grep hits / read spans /      │
│ escalation scratchpad   (volatile, after BP)   │
└───────────────────────────────────────────────┘
dynamic suffix (re-sent / re-priced each iteration)
```

## Operating constraints (from the mechanism)
- Stable prefix = system + tool definitions + a stable repo-map/file-tree.
  Place all dynamic content (grep results, read spans, the evolving suspect
  list) **after** the breakpoint.
- Never put volatile tool/grep results in the prefix — caching changing text
  forces a rewrite every turn and can *increase* latency (the documented
  "Don't Break the Cache" anti-pattern; system-prompt-only / exclude-tool-
  results placement gave 28–31% TTFT).
- Strict prefix match: one character of drift (a timestamp, a templated path,
  a `tool_choice` toggle) invalidates the cache. Keep the prefix deterministic.
- Order retrieved spans **deterministically** (same query → same ordered list).
  A nondeterministic rerank (float ties, timestamps) silently breaks the cache.
- Lookback is bounded (~20 content blocks) and breakpoints are few (≤4): for a
  long hunt, re-anchor with a rolling-window breakpoint rather than relying on
  a turn-1 write still being reachable late in the session.
- Provider-agnostic intent: Anthropic/Bedrock use explicit `cache_control` /
  `cachePoint`; OpenAI/Gemini infer from the prefix hash. The portable
  behavior the agent controls is **prefix stability + deterministic ordering**,
  not a specific marker. Caching complements retrieval; it does not replace it.
