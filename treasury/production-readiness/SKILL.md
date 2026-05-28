---
name: production-readiness
description: >
  Verify a feature is safe to ship — observability, rollback, runbook, migration safety, and a passing smoke test — and return a clear go / no-go / conditional verdict. Use when the user asks "is this safe to ship", "production-readiness check", "pre-release review", or "do we have rollback and a runbook". Confirms the whole ops surface, scopes alerts to High+ severity, and requires a rollback plan before a human owns the final go decision. Do NOT use for the live, firing P0 incident itself.
---

# production-readiness

## Purpose
Verify a feature is safe to ship — observability, rollback, runbook, migration safety, and a passing smoke test — and
return a clear go / no-go / conditional verdict. This is the **go/no-go gate** that confirms the whole ops surface.

## When to use
- Triggers: *"is this safe to ship"*, *"production-readiness check"*, *"pre-release review"*, *"do we have rollback and a runbook"*.
- **Not this skill:** the live, firing P0 incident itself (that is incident response, not a pre-ship gate); security sign-off → `reviewing-software-security`.

## Input
- The feature + its operational surface: logs, metrics, config, migrations, alerts.

## Output
A **Production-Readiness** artifact + checklist + verdict. Skeleton:

```
# Production Readiness — <feature>
Observability:  logs (correlation id?) · metrics · traces · dashboard
Alerts:         <High+ only — what fires, to whom>
Rollback:       <plan — feature flag? revert? migration reverse?>
Runbook:        <link — "what to do when X breaks">
Migration:      <safe + reversible?>
Smoke test:     <critical-journey check — pass/fail>
Severity map:   P0 unusable · P1 core degraded · P2 edge · P3 cosmetic
Verdict:  GO  |  NO-GO  |  CONDITIONAL (conditions: …)
```

Second artifact — the **runbook** skeleton (so on-call has it before ship):

```
# Runbook — <feature> / <likely failure>
Symptom / alert:
Blast radius (who's affected):
Immediate mitigation (short-term fix):
Diagnosis steps:
Rollback:
Escalation / who to call:
Long-term fix (ticket):
```

## Decision rules
1. Production behavior must be **debuggable**: logs (with correlation/request/business IDs), metrics, and a dashboard; document the **runbook** before ship.
2. **Alert for High+ severity only** — don't page people constantly.
3. **Define severity up front**: P0 = customers can't use it · P1 = core degraded (~30–40% slower) · P2 = edge case · P3 = wording/pixel.
4. A **rollback plan is mandatory** — never ship something you can't reverse, and never fix prod by bypassing the pipeline.
5. Treat **tech-debt-explosion signals** (rising CI flakiness, build/test time, repeating incidents, silently rising error rate) as no-go signals.
6. A **green smoke test on the critical journey** is required before GO.

## Checklist
- [ ] Logs/metrics/traces present (correlation id)
- [ ] Alerts scoped to High+ severity
- [ ] Rollback plan defined and tested
- [ ] Runbook written for likely failures
- [ ] Migration safe + reversible
- [ ] Smoke test on the critical journey passes
- [ ] Severity + on-call defined

## Anti-patterns (never do)
- Ship with no rollback or runbook.
- Alert on everything → alert fatigue.
- Fix production directly, bypassing the PR/pipeline.
- Hide a known problem from management; ignore a silently rising error rate — an incident is incoming.
- Expose system vulnerabilities — defer security sign-off to `reviewing-software-security`.

## Example
**Input:** wallet top-up. **Output (excerpt):** logs w/ correlation id ✅, alert on error-rate >1% ✅, rollback =
feature-flag off ✅, runbook "provider down" ✅, migration reversible ✅, smoke test ✅ — but flag config not yet set →
**CONDITIONAL GO** (condition: set the production flag).

## Human approval gate
**Stop.** A human makes the final go / no-go before release; the skill produces the verdict + checklist and never
deploys or bypasses the pipeline.
