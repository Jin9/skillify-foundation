# Worked example: wallet top-up budget

Loaded on demand by `performance-cost-review`.

**Input:** "Wallet top-up is live at ~100 RPS / 80k users-day — set its
performance and cost budget."

**Output (excerpt):** *Workload:* 100 RPS peak, 80k top-ups/day. *Baseline:* p95
420 ms, $0.0021 / successful top-up (load test + 2 weeks prod — **data-driven**).
*Perf budget (CUJ = top-up success):* p95 ≤ 500 ms, throughput ≥ 150 RPS, DB
connections ≤ 60% of pool. *Cost budget:* ≤ $0.003 / successful top-up, ≤ $2.4k/mo
at current volume; instance + on-call maintenance counted. *Headroom:* 16%
latency, 33% throughput; connections hit the wall at ~2.2× volume (~7 weeks at
current growth). *Cheap levers:* provider-response cache (−30% calls) and a
partial index on the hot lookup (budget +3× write) before adding read replicas.
*Guardrails:* alert at 80% of any ceiling -> hard-cap unauthenticated retries ->
breaker at provider p99 > 2 s. *Re-review:* at 2× volume or if cost / top-up
rises 20%. → **Verdict: WITHIN BUDGET** — levers queued before spend.
