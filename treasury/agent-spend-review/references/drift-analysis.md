# Drift analysis (LiteLLM vs. state.json)

LiteLLM `/spend/logs` is the source of truth — it counts what was
actually spent at the model. `state.json.total_cost_usd` is a snapshot
the dispatcher fills from LiteLLM at stage completion, falling back to
hardcoded estimates when LiteLLM is unreachable.

Drift between the two is itself a finding worth surfacing.

## Common drift causes

| Drift | Likely cause |
|---|---|
| `state.json` < LiteLLM | Stages completed before LiteLLM had logged the spend; cron lag in LiteLLM's accounting. |
| `state.json` > LiteLLM | Stages used `cost_source: estimate` (LiteLLM was unreachable at completion). |
| `state.json` ≪ LiteLLM | Multiple worktrees ran in parallel against the same LiteLLM but only one was sampled. |
| `state.json` ≫ LiteLLM | Estimate is wildly out of date for the model's current pricing. |

## Computing drift

```bash
LITELLM_TOTAL=$(just llm-spend 168 | jq -r '.total_usd')
STATE_TOTAL=$(
  find . -name state.json -path '*/.agent/*' \
    -exec jq -r '.total_cost_usd // 0' {} \; \
  | awk '{s+=$1} END {printf "%.2f", s}'
)
DRIFT=$(python3 -c "
ll=$LITELLM_TOTAL
st=$STATE_TOTAL
if ll == 0:
    print('n/a')
else:
    print(f'{(st-ll)/ll*100:+.1f}%')
")
echo "LiteLLM: \$$LITELLM_TOTAL  state.json sum: \$$STATE_TOTAL  drift: $DRIFT"
```

## Drift thresholds

| Drift | Action |
|---|---|
| < 10% | Note as "within tolerance"; no action. |
| 10–25% | Surface as "drift; spot-check stages with `cost_source: estimate`". |
| > 25% | Surface as a finding; offer to draft a postmortem stub. |

## What "estimate" means

The dispatcher only marks `cost_source: estimate` when the LiteLLM call
to `/spend/logs` failed — usually because the proxy was down or the
master key was unset. If a stage's `cost_source` is `estimate`, the
dispatcher used the runner's hardcoded fallback, which is a static
guess from the runner script. These guesses drift fast.

## Actions on drift findings

- Spot-check the stage with the worst estimate by reading the runner.
- If the runner's estimate is materially wrong, surface a profile-author
  hand-off via `authoring-scaffold-profile` (the runner usually reads
  the estimate from a profile or a default).
- Consider rotating `LITELLM_MASTER_KEY` if drift correlates with a
  permission failure on `/spend/logs`.

These are recommendations only; the team-lead decides.
