---
name: scoping-technical-requirements
description: >
  Turns a messy business requirement into a clear, bounded technical scope a team can safely act on, surfacing the open questions and non-functional requirements that decide everything downstream. Use when the user asks "help me scope this requirement", "is this ticket ready to build", "what's missing from this ask", or "turn this business request into technical scope". Produces a Requirement Analysis artifact with in-scope/non-scope, quantified NFRs, BA/PM open questions, risk flags, and a ready/needs-clarification verdict. Do NOT use for extracting rules from existing code (that belongs to business-logic-extractor) or for task-breakdown and estimation mechanics.
---

# scoping-technical-requirements

## Purpose
Turn a messy business requirement into a clear, bounded technical scope a team can safely act on — surfacing the
questions and non-functional requirements that decide everything downstream.

## When to use
- Triggers: *"help me scope this requirement"*, *"is this ticket ready to build"*, *"what's missing from this ask"*, *"turn this business request into technical scope"*.
- **Not this skill:** extracting rules from *existing code* → `business-logic-extractor`; task-breakdown and estimation mechanics are a separate planning step.

## Input
- A raw requirement: ticket / feature request / business ask / meeting note.
- (Optional) expected scale, current system context, deadline.

## Output
A **Requirement Analysis** artifact + a filled checklist + a verdict. Skeleton:

```
# Requirement Analysis — <feature>
Business goal (1 sentence):
In scope:
  - …
Out of scope (non-scope):
  - …
Assumptions:
  - …
Open questions for BA/PM:        # ambiguity becomes a question, never a silent assumption
  - …
Non-functional requirements:     # derived from expected volume
  - Volume/throughput:  Latency:  Data growth:  Cost ceiling:
Technical impact (systems/domains touched):
  - …
Risk flags:  [ ] unclear  [ ] likely to change  [ ] coupled to uncertain legacy
Verdict:  READY for feasibility  |  NEEDS CLARIFICATION
```

## Decision rules
1. Derive **non-functional requirements from expected user volume before anything else**. No NFRs → not ready.
2. If the requirement is **unclear, likely to change, or tied to uncertain legacy**, flag it and **do not estimate yet** — these are the patterns that guarantee a wrong estimate.
3. Make ambiguity an **open question for the BA**, not a buried assumption; poor business requirement analysis is the #1 cause of timeline blowups.
4. **Simple first** — capture the minimal requirement; defer nice-to-haves (YAGNI).
5. If clarifying the ask needs **>30 min of discussion**, it needs a written design doc downstream.

## Checklist
- [ ] Business goal restated in one sentence
- [ ] In-scope and non-scope both listed explicitly
- [ ] NFRs quantified (volume, latency, data growth, cost)
- [ ] Open questions for BA/PM captured
- [ ] Risk flags set (unclear / volatile / legacy-coupled)
- [ ] Assumptions made explicit
- [ ] Verdict recorded (ready / needs-clarification)

## Anti-patterns (never do)
- Convert ambiguity into silent assumptions instead of open questions.
- Estimate or commit before NFRs exist.
- Let the PM "just push for output" without understanding technical effort.
- Gold-plate the scope — adding what isn't needed yet violates simple-first/YAGNI.
- Answer the business on technical matters, or finalize scope, without the team knowing.

## Example
**Input:** "Let users top up their wallet with a new payment provider."
**Output (excerpt):** *Goal:* enable wallet top-up via Provider X. *In scope:* top-up flow, provider callback, ledger
entry. *Non-scope:* refunds, payouts. *Open questions:* expected TPS? settlement SLA? idempotency key from provider?
*NFRs:* ~20 top-ups/s peak, <2s p95, idempotent. *Risk flags:* ✅ likely-to-change (provider contract not final) →
**NEEDS CLARIFICATION** before feasibility.

## Human approval gate
**Stop.** The BA/PM must confirm in-scope vs non-scope and answer the open questions before the requirement advances to
feasibility. The lead does not finalize scope unilaterally when it affects the team.
