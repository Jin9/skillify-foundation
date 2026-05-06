# Suggested team-lead actions

The skill never acts. It surfaces options for the team-lead. Use this
table to populate the "What the team-lead might consider" section of
the readout.

## By outlier class

| Outlier | Possible action |
|---|---|
| User > 2σ above median | 5-minute walkthrough with the user at the Monday meeting. Likely cause: tight cap thrash (multiple aborts) or one expensive workflow. Not necessarily a problem. |
| Workflow > $10 | Postmortem within 48h (PLAYBOOK §4 trigger). Hand off to `agent-workflow-postmortem`. |
| Model 3× squad median per-request | Spot-check whether the model assignment in profiles is still appropriate. Hand off to `authoring-scaffold-profile` if a different model is warranted. |
| Drift > 25% | Spot-check `cost_source: estimate` stages. Possibly rotate `LITELLM_MASTER_KEY` if the proxy was unreachable repeatedly. |
| LiteLLM total far above expected | Audit recent prompt-library entries for prompts that explode token use under longer goals. Hand off to `drafting-stage-prompt` to update or anti-pattern. |

## Action shape

Each action in the readout is a one-liner:

```markdown
- <description> — owner: @<handle> — by: <date>
```

Owners are GitHub handles, not first names. Dates are ≤ 14 days.

## Things NOT to recommend

- Freezing a user's API access.
- Lowering everyone's cap to "punish" a single overrun.
- Rotating the master key as a routine action (it is a quarterly /
  off-boarding action, not a weekly one).
- Switching models mid-week without prompt-library re-validation.

## Escalation actions

- Approval-log integrity check failed → security incident channel,
  same-day.
- Real PII detected in archived logs from postmortems → compliance
  contact, same-day.
- Sandbox bypass confirmed → security incident, pull pinned proxy
  digest for audit.

These appear as a separate "Escalations" block above the actions list,
with explicit team-lead ack.
