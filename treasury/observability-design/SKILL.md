---
name: observability-design
description: >
  Designs how a feature is observed in production — correlation/request/business IDs, logs, metrics, traces, dashboards, and alerts — so production behavior is debuggable and only meaningful problems page anyone. Use when the user asks "what should we log/trace here", "design the dashboards and alerts", "make this debuggable in prod", or "what should alert on-call". Produces an Observability Design artifact with an ID strategy, log/metric/trace plan, dashboard panels, SLO burn-rate alert recipe, and quiet signals to watch. Do NOT use for defining cost/performance budgets.
---

# observability-design

## Purpose
Design the telemetry a feature needs so production behavior is **debuggable** and only **meaningful** problems page a human — IDs, logs, metrics, traces, dashboards, and alerts wired to the user journeys that matter.

## When to use
- Triggers: *"what should we log/trace here"*, *"design the dashboards and alerts"*, *"make this debuggable in prod"*, *"what should alert on-call"*.
- **Not this skill:** defining cost/performance budgets is a separate budgeting step.

## Input
- The feature + its **Critical User Journeys (CUJs)**, expected volume, and current telemetry stack.
- (Optional) existing dashboards, SLO targets, log conventions.

## Output
An **Observability Design** artifact + checklist + verdict. Skeleton:
```
# Observability Design — <feature>
IDs:            correlation_id · request_id · business_id (e.g. wallet_id)
Logs:           events + levels (structured, IDs on every line)
Metrics:        RED (rate/errors/duration) + USE; SLIs per CUJ
Traces:         spans across <service hops> with the IDs propagated
Dashboard:      panels (RED per CUJ, saturation, dependency health)
Alerts:         SLO burn-rate windows → page | ticket; High+ severity only
Quiet signals:  [ ] error rate silently rising  [ ] latency creep  [ ] saturation trend
Verdict:        OBSERVABLE  |  GAPS — not ship-ready
```

## Decision rules
1. Make production behavior **debuggable first**: structured logs carrying **correlation / request / business IDs**, plus metrics and traces — core monitoring lives on a **dashboard**.
2. **Page for High+ severity only.** Automation alerts at High+; everything else is a dashboard/ticket — don't page people constantly.
3. Tie page-worthiness to **severity** — P0 customers can't use it; P1 core ~30–40% degraded; P2 edge; P3 cosmetic.
4. Instrument the **quiet signals** — a silently rising error rate means an **incident is incoming**; watch it before it pages.
5. After any fix ships, **monitor exactly what changed**.
6. Define SLIs as `good/valid` events per **CUJ (not per service)**; alert on **multi-window burn rate** (fast 14.4× / 1h confirmed at 5m → page; 6× / 6h; 3× / 24h; two-window AND gate to stop flapping); cap 4–6 SLOs/service; bucket histograms at the SLO threshold.

## Checklist
- [ ] Correlation / request / business IDs defined and propagated through every hop
- [ ] Logs structured, leveled, and ID-stamped
- [ ] RED/USE metrics + one SLI per Critical User Journey
- [ ] Dashboard panels cover the CUJs and dependency health
- [ ] Alerts are SLO burn-rate based and **High+ only**; each routes page vs ticket
- [ ] Quiet signals listed (silent error-rate rise, latency creep, saturation)
- [ ] Verdict recorded

## Anti-patterns (never do)
- Page on P2/P3 or on raw thresholds that flap — pages are for High+ only.
- Ship logs/traces with no correlation/request/business ID — undebuggable in prod.
- Ignore a silently rising error rate until it becomes a P0.
- Define SLOs per service instead of per user journey.
- Expose secrets or credentials in plaintext in logs or traces.

## Example
**Input:** "We just launched wallet top-up via Provider X — make it observable."
**Output (excerpt):** *IDs:* `correlation_id` end-to-end, `request_id` per call, `wallet_id` + `provider_ref` as business IDs. *SLI (CUJ = top-up success):* `successful_topups / valid_topups` over 28d, SLO 99.5%. *Alert:* fast burn 14.4× over 1h (verified at 5m) → **page**; 3× over 24h → ticket. *Quiet signal:* provider-callback error rate creeping up. → **OBSERVABLE**.

## Human approval gate
**Stop.** A human owns **which signals page on-call** — alert rules are a one-way-door choice that commits other people's sleep. The lead recommends the alert recipe and SLO targets; the on-call owners and SRE sign off before it goes live, and the lead does not change paging policy unilaterally.
