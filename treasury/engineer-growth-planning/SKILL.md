---
name: engineer-growth-planning
description: >
  Grows an engineer — runs useful 1-on-1s, chooses teach-by-doing vs teach-by-telling, spots underperformance early, delegates to build capability, and reviews thinking (not just output) without becoming the bottleneck. Use when the user asks "how do I grow this engineer", "structure a 1-on-1", "should I delegate this or do it myself", or "this junior's code works but the thinking is wrong". Produces a Growth Plan / 1-on-1 notes with ladder rung, missing skills, teach-mode and delegation decisions, underperformance signals and action, the 6-question agenda, and a promotion-signal check. Do NOT use to review the PR's code itself, or for formal HR performance management.
---

# engineer-growth-planning

## Purpose
Grow an engineer deliberately — place them on the growth ladder, pick teach-by-doing vs telling, delegate to build capability, catch underperformance early, and review **thinking, not just output** — without becoming the bottleneck.

## When to use
- Triggers: *"how do I grow this engineer"*, *"structure a 1-on-1"*, *"should I delegate this or do it myself"*, *"this junior's code works but the thinking is wrong"*.
- **Not this skill:** formal HR performance management; pure stakeholder / up-management is a different concern.

## Input
- The engineer + situation: a 1-on-1 due, a delegation call, a growth question, or code that works but reasons wrong.
- (Optional) recent work, current ladder rung, timeline pressure.

## Output
A **Growth Plan / 1-on-1 notes** artifact + checklist. Skeleton:
```
# Growth Plan — <engineer>
Current rung → next rung:   # onboard → asks → solo → reviews others → supports squads → supports legacy →
                            # discusses w/ legacy → designs APIs → designs L3 (partial) → designs L4 in own scope
Missing skills:
  - …
Teach mode:        teach-by-doing (deep experience fast) | teach-by-telling (overview/foundation, then direction)
Delegation:        delegate (~90%) | guide-don't-solve (blocked) | do-it-myself (~10%, critical + hard deadline)  — why:
Underperformance:  signals seen [quality↓ / delivery↓ / less support / less opinion] → action:
1-on-1 agenda:     1 work lately?  2 working w/ team?  3 anything to discuss (career/blockers)?
                   4 their feedback on the lead  5 lead's feedback on them  6 goals they want
Promotion signals: ownership · supports team & legacy · handles complex problems · shares & mentors
```

## Decision rules
1. Place the engineer on the **growth ladder** and target the **next rung**: onboard → takes tasks but asks → takes tasks alone → reviews others → supports other squads → supports legacy → discusses with legacy team → designs APIs → designs L3 (partial) → designs L4 in own scope.
2. Run the **1-on-1 with the 6 questions** — work lately · working with the team · anything to discuss (career/blockers) · their feedback on the lead · lead's feedback on them · goals they want.
3. Choose the **teach mode**: **teach-by-doing** to build deep experience fast; **teach-by-telling** to lay the overview/foundation first, then set direction for doing.
4. Watch **underperformance signals** — quality drops, delivery slows, less support work, less opinion-sharing — then check which criteria + severity, recommend an example & teach problem-solving, warn on work style, and if critical → 1-on-1.
5. **Delegate ~90%**; do it yourself only if critical with a hard deadline (~10%); when the junior is blocked, **guide — don't solve for them**.
6. Judge promotion **beyond output**: ownership · supporting team & legacy · handling complex problems · sharing & mentoring.
7. When the **code works but the thinking is wrong**, classify must-change / consider / nit; if there's time, explain and **teach the new way of thinking**, ask them to fix and resubmit — share the thinking so an agent can help check at scale.

## Checklist
- [ ] Current rung and next rung identified
- [ ] Missing skills for the next rung listed
- [ ] Teach mode chosen (doing vs telling) with reason
- [ ] Delegation decision recorded (delegate / guide / do-it-myself + why)
- [ ] Underperformance signals checked; action set if triggered
- [ ] 1-on-1 agenda covers all 6 questions
- [ ] Promotion signals assessed

## Anti-patterns (never do)
- Talk down to or condescend to juniors — it undermines the highest goal: growing the team's ownership and productivity.
- Solve the junior's blocked problem for them instead of guiding.
- Fix the wrong thinking silently — skip the teaching when there is time to teach it.
- Become the bottleneck by hoarding work instead of delegating ~90%.

## Example
**Input:** "My junior's wallet-ledger code passes tests, but he double-writes the ledger then compensates instead of making the top-up idempotent."
**Output (excerpt):** *Rung:* takes tasks but asks → target *takes tasks alone*. *Classify:* must-change (correctness). *Teach mode:* teach-by-telling — lay the idempotency/ledger-invariant foundation, then have him redo it. *Delegation:* guide-don't-solve — he reworks it. *Action:* explain the new thinking, ask to fix and resubmit; note pattern for the next 1-on-1.

## Human approval gate
**Stop.** The skill recommends growth moves and drafts 1-on-1 notes — it never grades, promotes, or manages anyone. The human lead owns every people decision and every conversation; the mentor's highest duty is growing ownership, never talking down.
