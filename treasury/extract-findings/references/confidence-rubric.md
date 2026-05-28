# Confidence rubric (qualitative, evidence-strength-based)

`extract-findings` tags every finding with `confidence ∈ {high, medium, low}`.
The bucket is chosen by the **strength of the evidence relative to the
claim**, not by the writer's vibe.

| Bucket | Evidence pattern |
|---|---|
| `high` | Direct verbatim quote from a primary source. OR same claim corroborated by ≥2 independent sources (`supporting_source_ids` populated) without contradiction. |
| `medium` | Close paraphrase from a single source. OR direct quote from a single secondary source. OR direct quote with mild hedging in the source ("may", "in some cases"). |
| `low` | Inferred from indirect evidence (the source implies but does not state). OR single source with significant hedging. OR contradicted by another source in `disputed_by`. |

When uncertain between two buckets, pick the lower one.

## Why ordinal, not numeric

Never emit a numeric confidence score. Numeric `0.0–1.0` scores from LLMs are
poorly calibrated — published evals report Expected Calibration Error > 0.37
on LLM-generated confidences — and they collapse into overconfidence. The
ordinal `high | medium | low` scale is both more useful (downstream stages
can act on bucket membership) and more honest (it doesn't pretend at
precision the model can't actually deliver).

## How downstream stages consume it

- `synthesize-report` weights `high`-confidence findings more in the prose
  and tends to surface `low`-confidence ones with hedging language.
- `review-report` scans for `low`-confidence findings during the consistency
  pass and may flag the corresponding sentences for softening.
- `reporting-research-run` includes a confidence histogram in the per-stage stats row
  for `extract-findings`.
