# Approval Request — <artifact id>

- **Artifact id:** <id from the tracking board>
- **Issued at:** <ISO 8601 timestamp>
- **Tracking board card:** <link>
- **Source artifact:** <branch / PR / commit / specification path>

## Summary

<One paragraph: what changed, why it changed, which bounded context or system area it touches.>

## Lifecycle status

| Phase | Agent | Status | Notes |
|-------|-------|--------|-------|
| BA Intake | <BA agent identity> | done | <link to spec> |
| Architecture | <Architect agent identity> | done | <link to design doc> |
| Implementation | <Developer agent identity> | done | <link to PR or feature branch> |
| QA | <QA agent identity> | done | <verdict: pass / conditional> |

## QA verdict

- **Verdict:** <pass | conditional>
- **Coverage on changed lines:** <e.g., 86%>
- **Adversarial findings:** <count by severity, link to QA report>
- **Concurrency / chaos checks:** <pass | conditional with mitigations | n/a — explain>
- **Mitigations applied (for conditional verdicts):** <list, with link to commit>

## Residual risks

<Bulleted list of risks the human approver should weigh — what remains uncertain, what the rollback path is, what to watch in the first hours after deploy.>

## Compliance and regulatory check

- **Regulatory frameworks touched:** <e.g., PDPA, GDPR, PCI-DSS, BOT, OJK, MAS, or n/a>
- **PII handling:** <masked / encrypted / not in scope>
- **Audit-log changes:** <new event types, retention, sink>
- **Dual-control checkpoints:** <if the change touches disbursement / override / refund / write-off>

## Rollback plan

- **Trigger:** <what signal indicates rollback>
- **Mechanism:** <feature flag / revert commit / dual-write window>
- **Owner during rollback window:** <human role>

## Approval ask

> Approval needed from: <named human tech-lead role>
> Approve in this channel by replying with: `APPROVE <artifact id>`
> Reactions, presence checks, and timeouts do NOT satisfy this gate.
