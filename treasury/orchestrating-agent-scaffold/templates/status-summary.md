# Status summary

Use after a status / where-are-we / show-progress request.

```
<workflow_id>: <current_stage> <status> (<elapsed>m, $<spent>/$<cap>, src=<real|estimate>)

Stage status:
  ✅ <done stages>
  ⏳ <running stage> — last log: "<last 1 line>"
  🔒 <blocked stage> — gate pending
  ⏸ <pending stages>

Cost: $<spent> / $<cap>  (mix: <X real, Y estimate>)
Started: <iso>            Elapsed: <Xh Ym>

Source: state.json (mtime <iso>); stages/<stage>.log (last <n> lines)
```

When `approval_pending` is non-null, append the approval block from
`references/monitoring-loops.md`.

When the most recent stage is `failed`, append the failure-triage block
from `templates/failure-triage.md`.
