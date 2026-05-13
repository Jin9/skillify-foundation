---
name: reviewing-agent-spend
description: >
  Run the Monday cost-review ritual for an agent-scaffold squad. Wraps
  just llm-spend, segments by user and model, flags greater-than-2-sigma
  outliers and any single workflow over $10, and drafts a one-page
  team-lead readout. Treats LiteLLM /spend/logs as the source of truth
  and surfaces drift against state.json estimates. Use when the user
  says "weekly cost review", "Monday cost review", "agent spend report",
  "cost outliers this week", "who spent the most", "spend by model",
  "anomaly check on agent spend", or after a workflow finishes and the
  squad wants a snapshot. Does NOT raise any spend caps, freeze any
  users, or rotate any keys — those are team-lead actions taken outside
  this skill. Do NOT use for application-level cost analysis (cloud
  bills, RDS, K8s) or for predicting next-week spend.
---

# Reviewing agent spend

## Purpose

The PLAYBOOK §4 Monday ritual: a 15-minute team-lead-led review of the
last 7 days of agent spend. This skill produces the readout — pulls
LiteLLM `/spend/logs`, segments by user and model, flags outliers, and
drafts a paste-into-standup summary. The team-lead acts on the readout;
the skill does not act.

## When to use this skill

- The user says "weekly cost review", "Monday cost review", "agent spend
  report", "spend outliers", "who spent the most".
- A postmortem references cost drift and the user wants the squad-level
  picture.
- A team-lead asks for the readout before the Monday meeting.

Do NOT use this skill to:
- Raise or lower a spend cap. Caps are per-workflow and require explicit
  user instruction (handled by `orchestrating-agent-scaffold`).
- Rotate `LITELLM_MASTER_KEY`. That is a quarterly or off-boarding
  action, taken by the team-lead outside this skill.
- Predict next-week spend. Forecasting requires modeling beyond this
  skill's scope.
- Replace `state.json` cost values. State stays as-is; the skill reports
  the drift, the team-lead decides whether to file a postmortem.

## Universal preamble

1. Confirm the working directory is an agent-scaffold checkout
   (`docker/litellm/`, `justfile` exist).
2. Confirm LiteLLM is reachable:
   ```bash
   curl -fsS http://127.0.0.1:4000/health/readiness
   ```
   If unreachable, run `just llm-up` first (or surface the failure to the
   user — do not produce a phantom report from `state.json` alone).
3. Ask the user for the window if not stated. Default to 168 hours (7
   days) for the Monday ritual; 24 hours for ad-hoc daily checks.

## Core workflow

### 1 — Pull the spend report

```bash
just llm-spend <hours>
```

This wraps `spend_report` (the scaffold's `cost.sh` helper) and prints a
JSON report shaped:

```jsonc
{
  "hours": 168,
  "requests": 412,
  "total_usd": 23.47,
  "by_user":  [{"user": "alice", "spend": 9.12}, ...],
  "by_model": [{"model": "claude-critique", "spend": 6.30, "requests": 41}, ...]
}
```

Capture the JSON for the readout.

### 2 — Cross-check against state.json snapshots

Walk every active worktree's `.agent/state.json` (one per worktree).
For each, compute:

```bash
jq -r '.workflow_id, .total_cost_usd' .agent/state.json
```

Sum and compare against `total_usd` from LiteLLM. Drift > 10% is
worth flagging — usually means a stage's `cost_source` was `estimate`.

### 3 — Detect outliers

Two outlier classes:

| Outlier | Threshold |
|---|---|
| User outlier | spend > median(by_user) + 2 × stddev(by_user) |
| Workflow outlier | any single workflow's `total_cost_usd` > $10 |
| Model outlier | one model's per-request average > 3× the squad median per-request |

Compute via `jq` or hand off to a shell calculator. See
`references/outlier-math.md` for the exact arithmetic.

### 4 — Draft the readout

Use `templates/spend-readout.md`. Sections:

- Window and total
- Top three users by spend (no shaming — surface for context)
- Top three models by spend
- Outliers list with one-line context per entry
- Drift findings (LiteLLM vs. state.json)
- Action items the team-lead might consider (recommendations only)

### 5 — Surface to the user

Print the readout in the response. Suggest the team-lead paste it into
the squad standup channel, and offer to draft postmortem stubs for any
workflows over $10.

## Output format

A markdown readout matching `templates/spend-readout.md`. The skill
does not write a file by default; the team-lead pastes it where they
want. If the user asks for a file, write at
`docs/spend-reports/YYYY-Www-spend.md` (ISO-week format).

## Constraints

- DO NOT raise/lower caps, rotate keys, or freeze users.
- DO NOT shame individuals. Outliers are surfaced for context, not
  judgement; a researcher 2σ above median may have legitimate context.
- DO NOT fabricate cost numbers. If LiteLLM is unreachable, stop and
  surface the failure rather than estimate.
- DO NOT include real PII or customer names from log metadata.
- DO NOT produce a "by-key" segmentation if `LITELLM_MASTER_KEY` is
  shared across the squad — that breaks the per-user attribution and
  produces misleading data.
- DO NOT predict next-week spend.

## Validation gate

Before producing the readout:

1. LiteLLM `/health/readiness` returned OK.
2. The window matches what the user requested (default 168 if unstated).
3. Every outlier has a one-line context.
4. Every workflow > $10 is named and linked to its `state.json` (path).
5. Drift findings note `cost_source` mix (`real` vs. `estimate`).

## References

| Need | Reference |
|---|---|
| Outlier arithmetic (median, σ, per-request) | `references/outlier-math.md` |
| Drift interpretation (LiteLLM vs. state.json) | `references/drift-analysis.md` |
| Suggested team-lead actions per outlier class | `references/lead-actions.md` |

## Templates

- `templates/spend-readout.md` — one-page squad readout.
