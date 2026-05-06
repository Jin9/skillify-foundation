# Failure triage block

```
❌ Stage <stage> failed (<elapsed>s, $<spent>)

Last log lines (.agent/stages/<stage>.log):
  <line N-4>
  <line N-3>
  <line N-2>
  <line N-1>
  <line N>

Failure class: <stale-pid|cap-tripped|rate-limit|sandbox-egress|bad-prompt|real-bug>

Suggested next step:
  <one path — do not chain>

Postmortem trigger: <yes|no — based on cost, gate, cap>

Source: stages/<stage>.log, runlog.md (last 30 lines)
```

## Filling rules

- Quote log lines verbatim. Do not paraphrase stack traces.
- Pick one failure class. If unsure, surface "ambiguous — surfacing both"
  and list two — the user picks.
- Suggest exactly one fix. Chained suggestions ("try resume, then …, then
  …") trick the user into auto-piloting through the failure.
- Postmortem trigger is `yes` if any of: cost > $1, gate rejected, cap
  exceeded, sandbox egress test failed. Hand off to
  `agent-workflow-postmortem` when triggered.
