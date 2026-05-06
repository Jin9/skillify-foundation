# Launch confirmation

Before kicking off `just workflow`, surface this block. Wait for explicit
user confirmation.

```
🚀 Launch plan

Goal:     <one sentence, ≤500 chars>
ID:       <kebab-case slug, ^[a-z][a-z0-9-]{0,31}$>

Profile:  <code|literature|dataset|custom-name>
Stages:   <stage1> → <stage2> → ... → done
Gated:    <stage names that will pause for approval>

Cap:      $<n>.<nn>     (tier: <one-shot|plan-only|full|ceiling>)
Sandbox:  <on|off|n/a>  (rationale: <why>)

Models (per profile):
  <stage1>: <agent>/<model>
  ...

Pre-flight findings (if any):
  <bullet>
  <bullet>

Command (you type this):
  IMPLEMENT_SANDBOXED=<0|1> WORKFLOW_PROFILE=<name> \
    just workflow <id> "<goal>" <cap>

Reply "go" to proceed, or correct any field above.
```

## Notes

- `Sandbox` reads `n/a` when the profile skips `implement`.
- The user must type the `just workflow` command themselves only if they
  prefer to. The orchestrator is allowed to run `just workflow` once the
  user replies "go" — it is not a gated stage. Approval gates apply only to
  stages within a workflow.
- Pre-flight findings are advisory. They do not seed `state.json`; if they
  matter, save a prompt via `drafting-stage-prompt`.
