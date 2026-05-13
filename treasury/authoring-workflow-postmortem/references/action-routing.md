# Action-item routing

Postmortem actions land as PRs, not good intentions. Map each action to
the right artifact and skill.

## Routing table

| Mechanism | PR target | Hand-off skill |
|---|---|---|
| Bad library prompt | `prompts/library/<stage>/<topic>.md` (update or new entry) | `drafting-stage-prompt` |
| Anti-pattern to record | `prompts/library/_anti-patterns.md` (one-line append) | `drafting-stage-prompt` (loss mode) |
| Wrong cap default | `profiles/<name>.sh` (`SPEND_CAP_USD` or env var) | `authoring-scaffold-profile` |
| Wrong stage list | `profiles/<name>.sh` (`STAGES`) | `authoring-scaffold-profile` |
| Wrong gated stages | `profiles/<name>.sh` (`GATED_STAGES`) | `authoring-scaffold-profile` |
| Sandbox allowlist gap | `docker/sandbox-proxy/filter` + proxy rebuild | `configuring-sandbox-allowlist` |
| Gate-discipline failure | `docs/PLAYBOOK.md` §10 (clarify a question) | direct PR |
| Approval-log integrity bug | scaffold `dispatch.sh` / `approve.sh` | direct PR |
| Stale-PID handling | scaffold `common.sh` | direct PR |
| LiteLLM/state cost drift | docs note + maybe a runner cost-source flag | direct PR |
| Doctor missed a check | scaffold `doctor.sh` | direct PR |

## Action shape

Every action item is one row:

```markdown
- [ ] <action> — owner: @<github-handle> — by: YYYY-MM-DD
```

### Validation

- Owner is a GitHub handle, not a name. ("@chinnawat" beats "Chin".)
- Due date ≤ 14 days from postmortem date.
- Action verb is concrete: "update", "add", "remove", "rebuild", "rotate".
  Not "investigate", "consider", "look into".
- Each action maps to exactly one PR target (link to file path).

## Cross-skill hand-off

When the action requires a sibling skill:

1. Note the skill in the action row, e.g.:
   `— hand off to drafting-stage-prompt`
2. Do not invoke the skill from this one. The owner runs the sibling
   skill themselves so the PR has the right author.

## Escalation actions

Some mechanisms require escalation, not a PR:

- Approval-log integrity failure → security incident channel before PR.
- Sandbox bypass confirmed → security incident; pull the proxy image
  pinned digest for audit.
- Real PII in archived logs → redact immediately, rotate any exposed
  secrets, escalate to compliance/PII contact per the squad's playbook.

These appear as a separate "Escalations" section above the actions list,
with an explicit ack from the team-lead before the postmortem ships.
