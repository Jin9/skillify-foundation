---
name: daily-planner
description: >
  Turn a raw list of tasks into a prioritized, trackable daily plan and keep it
  current across days. Use when the user says "here are my tasks, help me
  prioritize", "plan my day", "what should I work on first", "make my daily
  plan", "update my plan / track my tasks", or "carry over what I didn't
  finish". Ranks each task by impact x urgency (P1 do-now / P2 schedule /
  P3 delegate / P4 drop-defer) with an S/M/L effort tag, then maintains a
  status-tracked plan file (todo / doing / done / blocked) that rolls unfinished
  work forward on the next run. Do NOT use for team delivery timelines, critical
  path, or man-day estimates (delivery-planning / risk-estimation), or for
  decomposing a research topic (plan-research).
---

# Daily Planner

## Purpose

Turn a raw, unordered list of tasks into a prioritized daily plan and maintain it as a
status-tracked file across days. This skill prioritizes and tracks *personal* work; it does not
estimate team delivery effort, build a critical path, or scope a research topic.

## When to use this skill

- Use when: the user drops a raw task list and asks to "prioritize", "plan my day", or "what should I work on first".
- Use when: the user asks to "make/update my daily plan", "track my tasks", or "carry over what I didn't finish".
- Do NOT use when: the user wants a team delivery timeline, critical path, or man-day estimate — that is `delivery-planning` / `risk-estimation`.
- Do NOT use when: the user wants to decompose a research topic into sub-questions — that is `plan-research`.

## Input

- Required: a list of raw tasks (free text, any order).
- Optional context: deadlines, focus-hour capacity for the day, dependencies, and a directory to store plans in.
- Default plan location: `~/daily-plans/`. The skill creates this folder if it is missing. Use a user-supplied directory instead when one is given.

## Choose the flow first

- **Flow A — Plan a new day**: no plan exists for the target date, or the user asks for a fresh plan.
- **Flow B — Update and track**: a plan for the target date already exists, or the user asks to change task status / show what is left.

If a same-date plan file already exists and the user asked for a fresh plan, switch to Flow B and update it — do not regenerate from scratch.

## Flow A — Plan a new day

1. **Parse**: split the raw input into discrete, single-outcome tasks. Keep the user's wording; do not merge or drop any task.
2. **Score**: rate each task on impact (high/low) and urgency (high/low) per the Decision rules, mapping to a P-band (P1–P4). Add an effort tag S / M / L.
3. **Carry over**: glob the plan directory for the most recent prior `daily-plan-*.md`. Pull every `todo`, `doing`, and `blocked` item forward (skip `done`), re-score them alongside the new tasks, and apply the aging bump. If none is found, note "no prior plan".
4. **Sort and capacity-check**: order P1 → P4, then larger-impact / smaller-effort first within a band. If the user gave focus hours, sum the planned effort and flag any overflow rather than silently cutting tasks.
5. **Write**: render the plan from `templates/daily-plan-template.md` and save to `~/daily-plans/daily-plan-<YYYY-MM-DD>.md` (target date, or today). Create the folder if missing.

## Flow B — Update and track

1. **Read** the existing dated plan file. Never delete it.
2. **Apply** the status changes the user names (todo → doing → done, or → blocked with a reason).
3. **Recompute** the remaining view: keep finished items checked, surface what is still open, and re-flag capacity if the day changed.
4. **List blocked** items with a reason and a next action.
5. **Rewrite** the same dated file in place (read-then-write). Do not create a second file for the same date.

## Output contract

- File: `~/daily-plans/daily-plan-<YYYY-MM-DD>.md` (one file per day; folder auto-created).
- Structure follows `templates/daily-plan-template.md`:
  - A `Capacity` line (focus hours vs. planned effort) when capacity is known.
  - `Priorities` grouped by band: `P1 — Do now`, `P2 — Schedule`, `P3 — Delegate`, `P4 — Drop / Defer`, each task carrying an effort tag and status.
  - A `Tracking` checklist using the status legend.
  - A `Blocked` section with reason + next action.
  - A `Carried over` section naming the source date.
- Status legend (used throughout): `[ ]` todo · `[~]` doing · `[x]` done · `[!]` blocked.

## Decision rules

1. **Impact x urgency mapping**: high impact + high urgency = P1 Do now; high impact + low urgency = P2 Schedule; low impact + high urgency = P3 Delegate; low impact + low urgency = P4 Drop / Defer.
2. **Impact** = real consequence if left undone or clear progress toward a stated goal. **Urgency** = a genuine deadline or time-sensitivity, not merely "feels pressing".
3. **Effort tag**: S = quick (well under an hour), M = a few hours, L = most of a day or more. Use the user's own sizing when given; otherwise estimate and mark it as an assumption.
4. **Carry-over aging**: a task carried forward for 2+ days gets its urgency bumped one level (so neglected work surfaces), capped at P1.
5. **Blocked handling**: a blocked task keeps its P-band, records the blocker and the next action to unblock, and does not count against today's capacity.
6. **Capacity honesty**: when focus hours are given, do not plan P1 + P2 effort beyond them — flag the overflow and recommend deferring the lowest-priority overflow tasks to P4.

## Checklist

- [ ] Every input task has a P-band, an effort tag, and a status — none dropped.
- [ ] Carry-over from the latest prior plan applied, or "no prior plan" noted.
- [ ] Over-capacity flagged when planned effort exceeds the stated focus hours.
- [ ] Blocked items list a reason and a next action.
- [ ] File written to `~/daily-plans/daily-plan-<date>.md` (folder created if missing); a same-date file is updated, not clobbered.
- [ ] No invented deadlines, capacity, or importance — assumptions are marked.

## Anti-patterns

- Never silently drop a task; low-value work goes to P4 and stays visible.
- Never overwrite an existing same-date plan from scratch — switch to Flow B and update it.
- Never invent deadlines, capacity, or task importance the user did not state; ask once or mark an assumption.
- Never plan P1 + P2 effort beyond stated capacity without flagging the overflow.
- Never delete the plan file or any prior-day file.

## Example

**User says**: "Plan my day: finish the auth bug, reply to Sam, read the new RFC, renew my passport (deadline Friday). I have ~5 focus hours."

**Result** (excerpt of `~/daily-plans/daily-plan-2026-06-01.md`):

```markdown
# Daily Plan — 2026-06-01
Capacity: 5h · Planned: ~5h

## Priorities
### P1 — Do now
- [ ] Finish the auth bug  · effort: M · status: todo
- [ ] Renew passport (deadline Fri)  · effort: S · status: todo
### P2 — Schedule
- [ ] Read the new RFC  · effort: M · status: todo
### P3 — Delegate
- [ ] Reply to Sam  · effort: S · status: todo
```

## Human approval gate

Present the plan (Flow A) or the status diff (Flow B) for confirmation before writing when you made
assumptions, when planned effort exceeds capacity, or when a same-date file would change. Write the
file only after the user confirms or when the request is unambiguous.
