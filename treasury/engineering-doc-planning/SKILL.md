---
name: engineering-doc-planning
description: >
  Decides what to document and what to deliberately skip, for which audience, keeping docs minimal, findable, and actionable across ADRs, C4 design docs, API contracts, system context, ubiquitous language, runbooks, and onboarding. Use when the user asks "what should we document here", "write this up as an ADR/design doc", "what goes in the runbook vs not", or "who is this doc for". Produces a Doc Plan with doc type, owner (lead vs team), audience, write-or-skip calls with reasons, where-to-find-it, plus the chosen doc's skeleton. Do NOT use to make the architecture decision an ADR records, or to design the API contract itself.
---

# engineering-doc-planning

## Purpose
Decide **what to write down, what to skip, and for whom** — so docs stay minimal, findable, and actionable instead of stale clutter. Pick the doc type (ADR · C4 L2–L4 · API spec · runbook · ubiquitous language · onboarding), name its owner and audience, then scaffold it. This decides **what to document** — it does not make the decision an ADR records, nor write the API contract itself.

## When to use
- Triggers: *"what should we document here"*, *"write this up as an ADR/design doc"*, *"what goes in the runbook vs not"*, *"who is this doc for"*.
- **Not this skill:** there is nothing here to extract from existing code.

## Input
- The thing to be documented: a decision, design, API, incident procedure, domain term, or onboarding gap.
- (Optional) audience, existing docs, where readers look first.

## Output
A **Doc Plan** artifact + checklist + the chosen doc's skeleton. Skeleton:
```
# Doc Plan — <feature/decision>
Doc type:   ADR | C4 (L2 network / L3 components / L4 entity+flow) | API spec | runbook | ubiquitous language | onboarding
Owner:      lead (high-level L3) | team (detail + API spec, owner-driven)
Audience:   future self | team | SRE monitors | stakeholders (data analysts, 3rd parties) | new joiners
Write or SKIP:  WRITE — <reason>   |   SKIP — <reason: PoC notes / unresolved-incident notes>
Findable where:  <repo path / wiki space / runbook index>
— Chosen skeleton (ADR) —
Context:        # forces / problem
Decision:       # what we chose
Consequences:   # trade-offs accepted
Status:         Proposed → Accepted → Superseded   # never edit an Accepted ADR — supersede it
```

## Decision rules
1. Use the **C4 doc shape** — **L2 network view, L3 high-level components, L4 entity + business-flow logic**, plus a components sheet (microservices, DB, API mapping) and an API spec (request, response, business logic, change-log).
2. **Ownership splits the work:** the **lead writes the high-level L3 doc**; the **team writes the detail + API spec** (owner-driven); API spec is reviewed by squad-lead + team; the **lead signs off**.
3. **Write** system context & boundary · ubiquitous language · decision records · manual incident runbook · new-joiner onboarding. **Do NOT write** PoC notes or notes on unresolved incidents.
4. **Always name the audience** before writing — future self, team, SRE monitors, stakeholders (data analysts, 3rd parties), new joiners — and write only for them.
5. A **junior asking the same thing 3×** means **documentation is missing** — write it.
6. A **TODO comment older than 6 months** → **delete it or do it**; never let it rot.
7. Record an architecture decision as an ADR (Context / Decision / Consequences / Status); **never edit an Accepted ADR — supersede it**.

## Checklist
- [ ] Doc type chosen and justified
- [ ] Owner assigned (lead vs team) per the L3/detail split
- [ ] Audience named explicitly
- [ ] Write-or-skip decided with a reason (skip PoC / unresolved-incident notes)
- [ ] "Findable where" recorded — readers can locate it
- [ ] Chosen doc's skeleton filled

## Anti-patterns (never do)
- Document PoC notes or unresolved-incident notes — they go stale and mislead.
- Write without naming an audience — audienceless docs serve no one.
- Edit an Accepted ADR in place instead of superseding it.
- Leave a junior asking the same question 3× without writing the missing doc.

## Example
**Input:** "We picked Provider X for wallet top-ups — write it up."
**Output (excerpt):** *Doc type:* ADR. *Owner:* lead (decision record). *Audience:* team + future self. *Write/skip:* WRITE the ADR; SKIP the spike notes from the bake-off (PoC). *Findable:* `/docs/adr/0007-payment-provider-x.md`. *Skeleton — Context:* need idempotent top-up + settlement callback; *Decision:* Provider X over Y; *Consequences:* vendor lock-in on payouts; *Status:* Proposed → (sign-off) Accepted.

## Human approval gate
**Stop.** A doc that affects the team ships only after review: API spec is reviewed by squad-lead + team, and the **lead signs off** on the high-level doc — the lead does not finalize team-facing docs unilaterally.
