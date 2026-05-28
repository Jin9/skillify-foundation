---
name: technical-debt-management
description: >
  Classify technical debt, prioritize it by business risk and cost-of-delay, make the dangerous debt
  visible, and negotiate a fix budget with PM/business. Use when the user asks "how do we deal with
  this tech debt", "is this debt worth paying down now", "make our tech debt visible", or "negotiate
  debt budget with the PM". Produces a Tech-Debt Ledger artifact plus a PM-negotiation summary and a
  checklist. Do NOT use to decide or bound a single refactor's scope (that is a separate refactor-decision step).
---

# technical-debt-management

## Purpose
Take an area of debt, classify each item, rank it by business risk and cost-of-delay, surface the dangerous debt that blocks customers, and turn it into a budget conversation the PM/business can decide on — so debt is **tracked and negotiated, never silently accumulated.**

## When to use
- Triggers: *"how do we deal with this tech debt"*, *"is this debt worth paying down now"*, *"make our tech debt visible"*, *"negotiate debt budget with the PM"*. This **classifies and prioritizes the whole debt portfolio**; deciding or bounding one specific refactor's scope is a separate refactor-decision step.
- **Not this skill:** a single bug's root cause → `progressive-bug-hunter`.

## Input
- An area, module, or service with suspected debt (plus any incident/PR history).
- (Optional) current sprint allocation, deadline, PM constraints.

## Output
A **Tech-Debt Ledger** + checklist + negotiation summary. Skeleton:
```
# Tech-Debt Ledger — <area>
| Item | Class (design/infra/service/pattern/code-smell) | Business risk | Triage (pay-now/tolerate/live-with) | Fix man-days | Cost of delay | Explosion signals present? | Recommendation |
|------|--------------------------------------------------|---------------|-------------------------------------|--------------|---------------|----------------------------|----------------|
| …    | …                                                | …             | …                                   | …            | …             | …                          | …              |
Allocation target:  tech-debt ≤5% vs ~95% feature this sprint   # push debt lower if the estimate allows
PM-negotiation summary: <risk to business> → <fix time / benefits / downsides> → <ask>
Verdict:  PAY NOW  |  TOLERATE  |  LIVE WITH IT  |  NEEDS PM DECISION
```

## Decision rules
1. **Tech debt = anything with risk of causing a problem** — system design, infra, microservices, coding patterns. **Bad code is only a code smell — a subset.** Classify before triaging.
2. Triage by customer impact: **pay-now** = risk of a bug that blocks customers · **tolerate** = below standard but still functions · **live-with-it** = general nits.
3. Hold the line on allocation: **tech debt ≤5% vs ~95% feature** — push debt lower when the estimate allows.
4. Watch the **explosion signals** — rising build/test time, same bug in new shapes, slower features at same size, "don't-touch-that" zones, rotting docs, declining coverage, CI flakiness, bus-factor onboarding, repeating incidents, a refactor that never finishes. Their presence promotes priority.
5. Negotiate, don't unilaterally schedule: analyze the debt's **business risk level** → analyze fix time/benefits/downsides → **present and agree with PM/business**.
6. Make debt **visible and automatically gated** — back the ledger with architecture fitness-function conformance checks in CI so drift is caught before it compounds.

## Checklist
- [ ] Each item classified (design/infra/service/pattern/code-smell)
- [ ] Each item triaged pay-now / tolerate / live-with by customer impact
- [ ] Business risk + fix man-days + cost-of-delay recorded
- [ ] Explosion signals checked per item
- [ ] Allocation target stated (≤5% debt)
- [ ] PM-negotiation summary written (risk → fix/benefit/downside → ask)

## Anti-patterns (never do)
- Treat "bad code" as the whole of tech debt — it is only a code-smell subset.
- Let debt accumulate without tracking, or hide a customer-blocking risk from the PM.
- Schedule a debt-fix budget unilaterally instead of agreeing it with PM/business.
- Ignore explosion signals until the refactor never finishes.
- Recommend a large refactor without an impact assessment.

## Example
**Input:** "Our fintech wallet ledger service keeps getting slower to change and a balance-mismatch bug keeps coming back in new shapes."
**Output (excerpt):** *Item:* non-idempotent ledger writes — *class:* coding pattern → *business risk:* duplicated debits block customers — *triage:* **pay-now** — *fix:* ~4 man-days — *cost of delay:* rising support tickets + repeat incidents — *explosion signals:* same bug new shapes + features slower at same size = **yes**. *Item:* inconsistent error logging — *triage:* **tolerate**. *PM summary:* "Ledger write idempotency is a customer-blocking risk; 4 man-days now vs growing incident load — request a 1-day-this-sprint slot." → **NEEDS PM DECISION**.

## Human approval gate
**Stop.** This skill **recommends and produces the ledger; it never commits a fix budget.** The PM/business owns the spend decision after seeing risk vs fix-time vs downside; the lead presents but does not answer business on this alone. The lead signs off on the final ledger artifact.
