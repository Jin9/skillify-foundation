# Review Conditions (detail)

Apply all four to each person, per milestone. "More than 95% similar" means the task
descriptions are effectively the same activity, allowing trivial wording differences.

## Condition 1 — Same task over more than 3 consecutive days

A task that is more than 95% the same, repeated on more than 3 consecutive **working days** →
REJECT.

- "Consecutive" is counted by billed working days (mandays claimed): skip weekends/holidays and
  leave days, but if the next *working* day is still the same task, it counts as consecutive.
- Weekends/holidays do NOT reset the counter — judge by the next working day.

## Condition 2 — 2-day repeat pattern over the monthly limit

The same task on 2 consecutive days is allowed, but the number of such 2-day patterns in a month
must not exceed:

```
max_allowed = round(3 * M / 20)
```

`M` = number of mandays billed that month (from Summary_Manday, or counted from the timesheet),
including Saturday if the person genuinely worked Saturday.

| M (days) | max_allowed |
|----------|-------------|
| ~20      | 3           |
| ~14      | 2           |
| ~10      | 2           |

If the count of 2-day repeat patterns exceeds `max_allowed`, REJECT the days that push it over.

## Condition 3 — Work outside the JO scope

Using the JO scope the user provided at intake, flag any day whose task description contains a
keyword or project code outside that scope — for example a ticket prefix from a different
project, or a feature name belonging to another product.

## Condition 4 — Duplicate work between people

Two people with the **same role**, on the same day (or 1 day apart), with more than 95% identical
task → REJECT.

**Who is the "copy"**: REJECT the person whose description is broader/more generic or carries less
context (no specific defect ID, no unique scenario) compared with the other. If it is genuinely
unclear, mark **both** and let the user decide.
