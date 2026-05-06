# Stage inventory

Every name in a profile's `STAGES` must match a runner script. The
default scaffold ships these:

| Stage | Runner path | Notes |
|---|---|---|
| `research` | `scripts/stage-runners/research.sh` | Long-context model recommended (Gemini default). |
| `plan` | `scripts/stage-runners/plan.sh` | Structured planning (Codex default). |
| `critique` | `scripts/stage-runners/critique.sh` | Architectural / security review (Claude default). |
| `implement` | `scripts/stage-runners/implement.sh` | Sandbox variant: `implement-sandboxed.sh`, selected when `IMPLEMENT_SANDBOXED=1`. |
| `review` | `scripts/stage-runners/review.sh` | Post-implementation diff review. |
| `test` | `scripts/stage-runners/test.sh` | Auto-detects go / npm / pytest / cargo / bats. |

## How runners are resolved

The dispatcher computes `runner = "scripts/stage-runners/${stage}.sh"`,
with the special case for sandboxed implement. If a stage has no runner,
the dispatcher fails the stage with `no runner script` and the workflow
halts.

## Custom stages

If the user names a stage that has no runner, **stop**. This skill does
not author runner scripts — they require state-mutation logic that
belongs in the scaffold itself.

Tell the user:

> Custom stage `<name>` has no runner at
> `scripts/stage-runners/<name>.sh`. You'll need to author the runner
> first (model the default runners). Once the runner exists, re-run this
> skill to write the profile.

## How to verify a runner exists

```bash
test -x scripts/stage-runners/<stage>.sh && echo "ok" || echo "missing"
```

Run this for each stage in the proposed `STAGES` before writing the
profile. Surface any missing stage to the user.

## Stage-skipping rules

Skipping a stage is supported by leaving it out of `STAGES`. Common
shapes:

- Full code pipeline: `research plan critique implement review test`.
- Literature: `research plan critique review` (no implement, no test).
- Dataset / EDA: `research plan critique implement review` (no test).

A profile with no `implement` should set `GATED_STAGES=""` since
`implement` is the only built-in gate.
