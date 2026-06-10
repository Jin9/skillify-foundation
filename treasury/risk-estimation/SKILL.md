---
name: risk-estimation
description: >
  Size complexity and separate known work from unknown risk, producing a man-day estimate with an explicit confidence and uncertainty multiplier plus the assumptions behind it, to feed delivery planning rather than replace it. Use when the user asks "how complex is this", "estimate the risk/effort", "how confident is this timeline", or "what's the unknown here". Produces a known-work estimate at a mid-level baseline, an unknown-risk list with spike time-boxes, an applied complexity multiplier, a confidence rating, and do-not-estimate flags. Do NOT use for turning the sized work into a sequenced task plan and critical path (use delivery-planning), or for deciding whether the thing can be built at all and what the options are (use technical-feasibility).
---

# risk-estimation

## Purpose
Size complexity and **separate known work from unknown risk** — producing a man-day estimate with an explicit **uncertainty multiplier**, a confidence level, and the assumptions behind it. This *feeds* delivery planning; it does not sequence or schedule.

## When to use
- Triggers: *"how complex is this"*, *"estimate the risk/effort"*, *"how confident is this timeline"*, *"what's the unknown here"*.
- **Not this skill:** negotiating the timeline politically is stakeholder-communication work; this skill sizes, it does not negotiate.

## Input
- A bounded requirement / chosen solution (ideally post scoping-technical-requirements + feasibility).
- (Optional) infra/maintenance scope, legacy touchpoints, known team velocity.

## Output
A **Risk & Estimate** artifact + a filled checklist + a recommended next step. Skeleton:
```
# Risk & Estimate — <feature>
Known work (man-days @ mid-level baseline):
  - <task> — <man-days>
Unknown risk (item — likelihood — spike? — time-box — question it must answer):
  - <risk> — <H/M/L> — spike: <y/n> — <Nd> — "<one clear question>"
Complexity multiplier applied:  ×<n>  — why: <reason>
Confidence:  <H / M / L>
Do-not-estimate flags:  [ ] unclear  [ ] likely to change  [ ] coupled to uncertain legacy
Assumptions:
  - …
Recommended next step:  ESTIMATE → delivery-planning  |  SPIKE FIRST  |  CLARIFY (no firm estimate)
```

## Decision rules
1. **Estimation process**: analyze the business requirement → analyze the problem → assume the solving process → estimate man-days. No skipped steps.
2. **Apply an uncertainty multiplier** — increase man-day sizing by the **complexity level** of the requirement, and state the multiplier explicitly.
3. **Do NOT give a firm estimate** when the requirement is unclear, likely to change, or coupled to uncertain legacy — these guarantee a wrong number; flag and route to clarify or spike.
4. Estimate against a **mid-level developer** baseline; usually the lead estimates for the team.
5. **Round-number estimates → the team hasn't thought it through;** and a **refactor running 2× its estimate → stop and reassess** rather than push on.
6. **Time-box a spike** by infra/maintenance sizing and the nature of the task, and give it **one clear question** it must answer.

## Checklist
- [ ] Business requirement + problem analyzed before any number (process followed)
- [ ] Known work itemized in man-days at the mid-level baseline
- [ ] Unknown risks listed with likelihood
- [ ] Spikes time-boxed, each with a single answerable question
- [ ] Complexity multiplier applied and justified
- [ ] Confidence rated (H/M/L)
- [ ] Do-not-estimate flags checked (unclear / volatile / legacy-coupled)
- [ ] Assumptions explicit
- [ ] Recommended next step recorded

## Anti-patterns (never do)
- Give a **firm estimate** on unclear, volatile, or uncertain-legacy work.
- Hand over **round-number** estimates that signal no thinking.
- Estimate to a senior's speed instead of the **mid-level baseline**.
- Drop the **complexity multiplier** and quote raw known-work man-days as if certain.
- Let a runaway spike or 2×-over refactor keep burning instead of reassessing.
- Promise a timeline without consulting the team.
- Hide technical problems from management.

## Example
**Input:** "Wallet top-up via Provider X." **Output (excerpt):** *Known work* @ mid-level: API contract 1d, ledger write 1d, idempotency 1d = 3d. *Unknown risk:* Provider X callback semantics — likelihood H — spike: yes — 1d — "Does the provider resend callbacks, and with what idempotency key?" *Multiplier:* ×1.5 (external dependency + unsettled contract). *Confidence:* M. *Flags:* ✅ likely-to-change (provider contract not final) → **SPIKE FIRST**, no firm date yet.

## Human approval gate
**Stop.** The team validates the assumptions and the multiplier; the business owns any date that follows. The skill produces the estimate + risk; it never commits a timeline. The lead signs off on the artifact and does not promise a date without the team.
