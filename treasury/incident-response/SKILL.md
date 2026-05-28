---
name: incident-response
description: >
  Stabilizes a firing production incident — classify severity, identify blast radius, drive a short-term mitigation, communicate, then run a blameless postmortem with prevention follow-ups. Use when the user asks "we have a P0/incident — what now", "triage this outage", "write the postmortem", or "how should on-call/rotation work". Produces an Incident artifact (severity, blast radius, timeline, mitigation, comms cadence) and a Postmortem artifact (root cause, fixes, owned prevention items). Do NOT use for localizing the offending code (progressive-bug-hunter) or a security breach's disclosure specifics (reviewing-software-security).
---

# incident-response

## Purpose
Stabilize a firing incident — classify, scope blast radius, mitigate short-term, communicate on cadence — then run a **blameless** postmortem that ships owned prevention follow-ups.

## When to use
- Triggers: *"we have a P0/incident — what now"*, *"triage this outage"*, *"write the postmortem"*, *"how should on-call/rotation work"*.
- **Not this skill:** localizing the offending code → `progressive-bug-hunter`; a security breach's disclosure specifics → `reviewing-software-security`.

## Input
- A live alert / outage report, or a resolved incident needing a postmortem.
- (Optional) severity, affected CUJs, dashboards, on-call roster.

## Output
Two short artifacts + checklist. Skeleton:
```
# Incident — <id>
Severity:        P0 | P1 | P2 | P3
Blast radius:    <users/services affected>
Timeline:        detect → mitigate → recover (timestamps)
Short-term fix:  <mitigation plan>
Comms log:       cadence + audience (business/leadership)
Monitor:         what was changed, now being watched

# Postmortem — <incident>
Summary (third-party readable):
Impact (system + users):
Root-cause summary:
Short-term fix process:
Long-term fix process:
Prevention action items:  owner · due · status
```

## Decision rules
1. **First 5 minutes, in order:** assess problem **scope** → analyze initial cause + draft a **short-term fix plan** → **contact SRE/Operations** → discuss and begin solving → **monitor** what was fixed.
2. Classify severity to size the response: **P0** customers can't use it at all; **P1** core ~30–40% degraded; **P2** edge, core journey intact; **P3** wording/pixel.
3. On-call: **High/critical-urgent → tech lead is primary**; medium/low → **rotate every 4–5 days**; dashboard-first.
4. Postmortem captures: summary · impact · **root-cause summary** · short-term fix · long-term fix.
5. Stay **blameless** — analyze problem/cause/fix/prevention, name the owner to coach and improve future assignment; **rarely blame** unless truly unacceptable.
6. Comms to business/leadership: analyze problem → cause + remediation **options** → explain step-by-step (cause + probability, short- and long-term fix, effort).
7. Lesson: thoroughly vet tech stack, community maintenance, and prevention **before committing** — a non-prod pass is not production proof.
8. Run the incident by **role**, not by crowd: Incident Commander commands (never debugs) · Comms Lead · Ops Lead · Scribe; hold a **severity-bound public-comms cadence** (e.g. SEV-1 first message within minutes even if incomplete, then every 30 min).

## Checklist
- [ ] Severity assigned (P0–P3) and blast radius scoped
- [ ] Short-term mitigation drafted; SRE/Operations contacted
- [ ] Incident Commander / Comms / Ops / Scribe roles assigned
- [ ] Leadership comms sent on cadence
- [ ] Fix monitored after deploy
- [ ] Blameless postmortem written (summary, impact, root cause, short + long-term fix)
- [ ] Prevention action items have owner + due date

## Anti-patterns (never do)
- Skip **contacting SRE/Operations** in the first 5 minutes.
- Blame the engineer instead of analyzing cause/prevention.
- Fix production directly, bypassing the PR/pipeline.
- Hide technical problems from management.
- Close the incident with prevention items that have no owner or due date.

## Example
**Input:** "Wallet top-ups are failing for everyone — Provider X callbacks timing out."
**Output (excerpt):** *Severity:* **P0** (customers can't top up). *Blast radius:* all top-up users; ledger writes paused. *First 5 min:* scope confirmed via dashboard → short-term plan = switch to polling fallback → SRE paged → mitigating → monitoring. *Roles:* tech lead as Incident Commander, Comms Lead to leadership every 30 min. *Postmortem (later):* root cause = provider TLS renegotiation; long-term = circuit breaker + provider SLA review; prevention item owned by payments squad, due next sprint.

## Human approval gate
**Stop.** Severity declaration, the leadership comms, and postmortem signoff are **human-owned**; the skill drafts them. It **never** fixes production directly or bypasses the PR/pipeline, and prevention action items are committed by their owners — not the skill.
