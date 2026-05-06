# Postmortem: <workflow_id> — <one-line summary>

**Date:** YYYY-MM-DD (UTC)
**Workflow:** `<workflow_id>`
**Researcher:** `@<github-handle>`
**Severity:** <wasted-spend | rejected-gate | failed-stage | cap-overrun | sandbox-failure>
**Cost impact:** $X.XX (real $A, estimate $B)
**Wall-clock impact:** <minutes>
**Archive:** `docs/postmortems/YYYY-MM-DD-<workflow_id>/`

## What happened

Plain-English narrative, two to four paragraphs. Cover:
- The goal the squad was pursuing.
- The stage that failed or surprised.
- How and when the researcher noticed.
- The recovery action.

Avoid jargon; this is read by the whole squad.

## Timeline (UTC)

Pull from `.agent/runlog.jsonl`:

```
- HH:MM:SS — workflow kicked off, cap $X
- HH:MM:SS — research done ($X, real)
- HH:MM:SS — plan done ($X, real)
- HH:MM:SS — critique done ($X, real) — verdict summary
- HH:MM:SS — implement gate approved by @<handle>
- HH:MM:SS — implement failed: <one-line cause>
- HH:MM:SS — researcher noticed via <ntfy | dashboard | tail>
- HH:MM:SS — `just abort` run by @<handle>
```

## Root cause

One sentence naming the **mechanism** and the **system gap** that allowed
the failure. Not "the model hallucinated"; rather, what allowed the
hallucination to slip past critique without being caught.

## Contributing factors

Conditions that enabled the failure but did not cause it. Examples:

- prompt phrasing was looser than usual
- spend cap default in the profile is too generous
- approver was on a phone notification (Question 1 violated)
- repo had uncommitted changes that confused diffs
- LiteLLM unreachable at runtime; cost source fell back to estimate

## Escalations (if any)

- [ ] Approval-log integrity failure — escalated to <channel>, <owner>, <date>
- [ ] Sandbox bypass confirmed — escalated to <channel>, <owner>, <date>
- [ ] Real PII in logs — redacted, rotated, escalated to <compliance contact>

Skip this section if no escalation triggered.

## What we'll change

Each action lands as a PR — not a good intention. Owners are GitHub
handles. Dates are within 14 days of the postmortem.

- [ ] <action> — owner: `@<handle>` — by: YYYY-MM-DD — hand off: `<sibling-skill>`
- [ ] <action> — owner: `@<handle>` — by: YYYY-MM-DD
- [ ] <action> — owner: `@<handle>` — by: YYYY-MM-DD

## Audit links

- `docs/postmortems/<date>-<id>/runlog.jsonl`
- `docs/postmortems/<date>-<id>/stages/`
- `docs/postmortems/<date>-<id>/approvals.md` and `approvals-verify.txt`
- `docs/postmortems/<date>-<id>/spend.txt`
- (if applicable) `docs/postmortems/<date>-<id>/proxy.log`

## Sign-off

- Researcher: `@<handle>` — date
- Team-lead ack: `@<handle>` — date
