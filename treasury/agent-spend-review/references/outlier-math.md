# Outlier arithmetic

The squad-spend report uses three outlier classes. Each maps to a
specific calculation against the LiteLLM JSON.

## User outlier

> A user spent > median(by_user) + 2 × stddev(by_user).

```bash
just llm-spend 168 | jq -r '
  .by_user
  | sort_by(.spend)
  | (length / 2 | floor) as $mid
  | .[$mid].spend as $median
  | (map(.spend) | add) as $sum
  | length as $n
  | (map(.spend) | map(. - ($sum / $n)) | map(. * .) | add / $n | sqrt) as $sigma
  | map(select(.spend > ($median + 2 * $sigma)) | "\(.user): $\(.spend) (median=$\($median), σ=$\($sigma))")
  | .[]
'
```

When the squad has < 5 active users, σ is noisy. Disable user-outlier
detection and only surface "highest spender" for context.

## Workflow outlier

> Any single workflow's `total_cost_usd` > $10.

Walk every `state.json` in the squad's worktrees:

```bash
find . -name state.json -path '*/.agent/*' \
  -exec jq -r 'select(.total_cost_usd > 10) | "\(.workflow_id): $\(.total_cost_usd)"' {} \;
```

> $10 is the PLAYBOOK ceiling without team-lead acknowledgement. Anything
above that is a postmortem trigger.

## Model outlier

> One model's per-request average > 3× the squad median per-request.

```bash
just llm-spend 168 | jq -r '
  .by_model
  | map(. + {avg: (.spend / .requests)})
  | sort_by(.avg)
  | (length / 2 | floor) as $mid
  | .[$mid].avg as $median
  | map(select(.avg > 3 * $median) | "\(.model): avg $\(.avg) per request (squad median $\($median))")
  | .[]
'
```

When the squad has only 1–2 models in active use, model-outlier
detection is noise. Skip the section if `by_model | length < 3`.

## Median-vs-mean

The skill uses **median** as the central tendency, not mean. Means are
distorted by exactly the outliers you're trying to detect.

## Reporting precision

- Round dollar values to two decimal places.
- Round σ to two decimal places.
- Show ratios (3×, 2.4×) to one decimal place.
- Never report a per-request cost below $0.001 — that level of precision
  is noise from LiteLLM rounding.
