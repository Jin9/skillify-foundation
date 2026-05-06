# Squad spend readout — week of YYYY-Www

**Window:** last <N> hours (<start ISO> → <end ISO>)
**Total spend:** $<X.XX> (LiteLLM `/spend/logs`, source of truth)
**Total requests:** <N>
**Active worktrees scanned:** <N>

## Top spenders (by user)

| User | Spend | Requests | Notes |
|------|-------|----------|-------|
| @<handle> | $<X.XX> | <N> | <one-liner if context warrants> |
| @<handle> | $<X.XX> | <N> | |
| @<handle> | $<X.XX> | <N> | |

## Top models (by spend)

| Model | Spend | Requests | Avg / req |
|-------|-------|----------|-----------|
| <model> | $<X.XX> | <N> | $<0.0X> |
| <model> | $<X.XX> | <N> | $<0.0X> |
| <model> | $<X.XX> | <N> | $<0.0X> |

## Outliers

### User outliers (> median + 2σ)

- @<handle> — $<X.XX> spent, vs. median $<X.XX>, σ $<X.XX>. Context:
  <one-line: "ran 4 implement attempts on the JWT refactor", or "n/a">.

(Or "no user outliers this window.")

### Workflow outliers (single workflow > $10)

- `<workflow_id>` — $<X.XX>. State: `<path/to/.agent/state.json>`.
  Postmortem trigger: <yes — hand off to `agent-workflow-postmortem`>.

(Or "no workflow outliers this window.")

### Model outliers (avg per-request > 3× squad median)

- <model> — avg $<0.0X>/req, squad median $<0.0X>. Context:
  <one-line>.

(Or "no model outliers this window.")

## Drift findings

- LiteLLM total: $<X.XX>
- Sum of state.json: $<X.XX>
- Drift: <±X.X%>
- `cost_source` mix: <X stages real, Y stages estimate>

<one paragraph: which stages drifted, why, and what to spot-check>

## What the team-lead might consider

- [ ] <action> — owner: @<handle> — by: <YYYY-MM-DD>
- [ ] <action> — owner: @<handle> — by: <YYYY-MM-DD>

## Escalations

(Skip this section if none.)

- [ ] <escalation> — channel: <#name> — owner: @<handle> — same-day

## Source

- LiteLLM: `/spend/logs?<window>`
- State snapshots: `find . -name state.json -path '*/.agent/*'`
- Generated: <ISO timestamp>
