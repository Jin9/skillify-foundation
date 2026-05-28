---
name: performance-cost-review
description: >
  Define and review a feature or service's performance budget and cost budget — latency/throughput/resource targets per critical user journey, cost-per-unit, baseline vs target, headroom, and the cheap-levers-first plan to stay within them — so resource efficiency is a measured budget, not a vibe. Use when the user asks "set a performance budget for this", "what's our cost budget / are we overspending", "how much headroom before we must scale", or "define cost-per-request / unit economics". Produces a Performance & Cost Budget artifact with targets, baseline, headroom, ordered cheap levers, and guardrails. Do NOT use to instrument or alert on production (a monitoring step), or to run the load test or read the cloud bill itself (use k6/JMeter/Locust and AWS Cost Explorer / GCP Billing).
---

# performance-cost-review

## Purpose
Turn a feature's NFRs and baselines into a concrete **performance + cost budget** — targets per critical user journey, cost-per-unit, **headroom**, and an ordered cheap-levers-first plan to stay inside them — so resource efficiency is a measured budget rather than a qualitative aside. This skill **recommends** the budget and the lever order; a human approves the numbers and any spend or scaling decision.

## When to use
- Triggers: *"set a performance budget for this"*, *"what's our cost budget / are we overspending"*, *"how much headroom before we must scale"*, *"define cost-per-request / unit economics"*.
- **Not this skill:** instrumenting / measuring or alerting on production behavior is a monitoring step — do that to *produce the baselines* this skill consumes; running the load test or reading the cloud bill itself uses your load-testing tool (e.g. k6 / JMeter / Locust) and cloud cost dashboard (e.g. AWS Cost Explorer / GCP Billing).

## Input
- The feature / service + its **Critical User Journeys (CUJs)**, expected volume, and the NFRs derived from business logic + volume.
- The **cost model**: the per-unit cost driver plus instance cost *and* dev-maintenance cost.
- Baseline metrics if they exist (production telemetry or a load test); otherwise an NFR-derived estimate.

## Output
A **Performance & Cost Budget** artifact + checklist + verdict. Skeleton:
```
# Performance & Cost Budget — <feature/service>
Workload:     volume / RPS / data size; NFRs derived from business logic + volume
Baseline:     measured metrics (data-driven)  OR  "no data yet — NFR estimate" (intuition mode)
Perf budget:  per CUJ — latency p50/p95/p99 · throughput · resource ceiling (CPU/mem/IO/connections)
Cost budget:  cost-per-unit (per request | per successful task | per tenant) · monthly target · instance + maintenance cost
Headroom:     % to ceiling per dimension · time-to-wall at current growth
Cheap levers: ordered — cache/batch · query/index (budget write-amplification) · cardinality/retention · rightsize/autoscale — exhaust before scaling spend; cite each lever's expected saving
Guardrails:   soft alert (trend/anomaly) -> hard cap (runaway/unauth) -> circuit breaker (degrade); attribution dims (tenant/feature/env)
Re-review:    re-baseline cadence + triggers — volume Nx · cost-per-unit drift >X% vs baseline · perf trend crosses budget midline · new cost driver
Verdict:      WITHIN BUDGET | AT RISK — levers before spend | OVER — needs decision | PROVISIONAL — re-baseline at <trigger>
```

## Decision rules
1. **Derive the budget from NFRs and volume, not vibes.** Compute the targets from business logic + expected volume, and set them **per Critical User Journey, not per service** — a service-wide average hides the one journey that is failing. Budget only the few dimensions that matter; do not target everything.
2. **Baseline before you budget — data-driven once data exists, estimate before it.** Use measured metrics once the project has enough signal (load test or active users); use an NFR-derived estimate for early or urgent cases. **State which mode you are in**; a budget with no baseline is labeled **PROVISIONAL** with a re-baseline trigger, never presented as firm.
3. **Budget cost in cost-per-unit, not the total bill.** Pick the unit that scales with the business — per request, per *verified* successful task, or per tenant. The monthly total cannot separate efficiency from growth. Count **instance cost and dev-maintenance cost**, not just the invoice. For a per-outcome unit, divide by *verified* successes — failures inflate the numerator only, so attempts-as-denominator makes a cheap-but-failing path look efficient.
4. **Exhaust cheap levers before buying capacity.** Apply in order: caching / batching -> query + index tuning (budget write-amplification: each secondary index ≈ 3–4× insert cost) -> metric/log cardinality + retention (cardinality and retention fixes alone routinely cut 30–50% of an observability bill) -> rightsizing / autoscaling -> only then bigger or more instances. Cite each lever's expected saving so the ordering is justified, not asserted.
5. **Guardrails are a hierarchy, not a single threshold.** Stack soft alert (trend / anomaly %) -> hard cap (runaway loop, unauthenticated spend) -> circuit breaker (shed load / degrade). Attribute cost and performance by tenant / feature / environment so the driver is findable. Page for **High+ severity only** — guardrails escalate through alert and cap before anyone is paged.
6. **Headroom is the deliverable, not just the current number.** Report **% to the ceiling per dimension** and the **time-to-wall** at current growth — that lead time is what lets leadership plan a scaling or spend decision. Re-baseline on a cadence and on **budget-consequence triggers** (volume Nx · cost-per-unit drift past a threshold · a perf trend crossing the budget midline · a new cost driver), not on a single raw reading.

## Checklist
- [ ] NFRs derived from business logic + expected volume; targets set **per CUJ**, not per service
- [ ] Baseline labeled **data-driven** (measured) or **PROVISIONAL** (NFR estimate) — mode stated
- [ ] Perf budget: latency p50/p95/p99, throughput, and resource ceilings per CUJ
- [ ] Cost budget in **cost-per-unit** (request / successful-task / tenant); instance + maintenance cost both counted
- [ ] Success denominator defined and independently audited for any per-outcome cost unit
- [ ] Cheap levers **ordered and exhausted** before any capacity / spend increase; each lever's expected saving cited
- [ ] Guardrail hierarchy (soft alert -> hard cap -> circuit breaker) + cost/perf attribution dimensions
- [ ] Headroom (% to ceiling + time-to-wall) reported per dimension
- [ ] Re-review cadence + budget-consequence triggers set
- [ ] Verdict recorded

## Anti-patterns (never do)
- Set a budget with no baseline and no NFR derivation — that is a wish, not a budget.
- Budget the total monthly bill instead of cost-per-unit — it cannot tell efficiency from growth.
- Average a target across the whole service instead of per user journey — it hides the journey that is failing.
- Scale capacity or jump to a bigger instance before exhausting caching, query/index, cardinality, and rightsizing levers.
- Add an index or a high-cardinality metric dimension without budgeting its write / cost amplification.
- Divide a per-outcome cost by attempts instead of *verified* successes — a cheap-but-failing path then looks efficient.
- Report a single current number with no headroom or time-to-wall — leadership cannot plan a scaling decision from it.
- Page on every threshold breach or cosmetic trend — guardrails are tiered (alert -> cap -> breaker); pages are High+ only.

## Example
**Input:** "Wallet top-up is live at ~100 RPS / 80k users-day — set its performance and cost budget."
**Output (excerpt):** *Workload:* 100 RPS peak, 80k top-ups/day. *Baseline:* p95 420 ms, $0.0021 / successful top-up (load test + 2 weeks prod — **data-driven**). *Perf budget (CUJ = top-up success):* p95 ≤ 500 ms, throughput ≥ 150 RPS, DB connections ≤ 60% of pool. *Cost budget:* ≤ $0.003 / successful top-up, ≤ $2.4k/mo at current volume; instance + on-call maintenance counted. *Headroom:* 16% latency, 33% throughput; connections hit the wall at ~2.2× volume (~7 weeks at current growth). *Cheap levers:* provider-response cache (−30% calls) and a partial index on the hot lookup (budget +3× write) before adding read replicas. *Guardrails:* alert at 80% of any ceiling -> hard-cap unauthenticated retries -> breaker at provider p99 > 2 s. *Re-review:* at 2× volume or if cost / top-up rises 20%. → **Verdict: WITHIN BUDGET** — levers queued before spend.

## Human approval gate
**Stop.** A human owns the budget numbers and any spend or scaling decision — committing money or capacity is a business call. A one-way-door spend (reserved instances, a tier upgrade, a bigger cluster) follows decide-then-explain: set a clear recommendation with rationale, but it still needs sign-off. The lead recommends the targets, the headroom, and the lever order; the eng owner and finance approve the budget and the scaling trigger. The skill produces the budget; it does not provision, scale, or commit spend.
