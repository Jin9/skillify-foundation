# Critique-to-plan cross-walk

The mechanism behind Question 2 of the five-question check.

## Mechanics

For every issue in `critique.md` at severity P1 or P2:

1. Extract the issue title and a one-sentence description.
2. Search `plan.md` for explicit mitigation text — file/line references,
   acceptance criteria, or test names that close the gap.
3. Tag the issue:

| Tag | Meaning |
|---|---|
| `addressed` | Plan has a numbered step that closes the gap with explicit acceptance criteria. |
| `partial` | Plan acknowledges the issue but the mitigation is vague ("we'll be careful with X"). |
| `unaddressed` | No mention in plan. |
| `accepted` | Plan explicitly states "we accept this risk because <reason>". |

4. Compute the verdict:

| Worst tag | Verdict |
|---|---|
| All `addressed` or `accepted` | Q2 = yes |
| Any `partial` or `accepted` without reason | Q2 = unsure (recommend reject; let user upgrade plan) |
| Any `unaddressed` P1 | Q2 = no (force reject) |
| Any `unaddressed` P2 | Q2 = unsure (recommend reject) |

## What counts as "explicit mitigation"

Concrete enough that an implementing agent could not skip it:

- A numbered plan step with file/line target and acceptance criteria.
- A test name that, if it passes, proves the gap is closed.
- A flag/config change with the new value spelled out.

What does NOT count:

- "We'll add validation" — too vague.
- "Following best practices" — empty.
- "TBD during implementation" — the gate is your last chance.

## Output shape per issue

```
[<severity>] <issue title>
  Source:    critique.md:<line>
  Tag:       <addressed|partial|unaddressed|accepted>
  Mitigation: <plan.md:<line>> "<short quote>"  (or "—" if unaddressed)
```

The reviewer block in `templates/gate-review.md` includes a section to
list every P1/P2 issue with these fields. Surface the full list — do not
summarize or drop entries.
