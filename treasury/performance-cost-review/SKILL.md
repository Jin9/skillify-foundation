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
Apply all six; full rationale for each is in `references/budget-method.md`.
1. **Derive the budget from NFRs and volume, not vibes** — targets per Critical User Journey, not per service; budget only the dimensions that matter.
2. **Baseline before you budget** — data-driven once data exists, NFR estimate before it; state the mode; an un-baselined budget is **PROVISIONAL**.
3. **Budget cost in cost-per-unit, not the total bill** — per request / *verified* successful task / tenant; count instance + dev-maintenance cost.
4. **Exhaust cheap levers before buying capacity** — cache/batch -> query/index -> cardinality/retention -> rightsize/autoscale -> only then more instances; cite each saving.
5. **Guardrails are a hierarchy** — soft alert -> hard cap -> circuit breaker; attribute by tenant/feature/env; page High+ only.
6. **Headroom is the deliverable** — report % to ceiling and time-to-wall; re-baseline on cadence and budget-consequence triggers, not one raw reading.

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

## Anti-patterns and example
The "never do" list (set a budget with no baseline, budget the total bill, page on every breach, divide cost by attempts not verified successes, etc.) is in `references/anti-patterns.md`. A complete worked wallet-top-up budget is in `examples/wallet-top-up.md`.

## Human approval gate
**Stop.** A human owns the budget numbers and any spend or scaling decision — committing money or capacity is a business call. A one-way-door spend (reserved instances, a tier upgrade, a bigger cluster) follows decide-then-explain: set a clear recommendation with rationale, but it still needs sign-off. The lead recommends the targets, the headroom, and the lever order; the eng owner and finance approve the budget and the scaling trigger. The skill produces the budget; it does not provision, scale, or commit spend.
