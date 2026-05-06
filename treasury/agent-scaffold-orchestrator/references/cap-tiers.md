# Cap tier selection

Pick the smallest cap that fits the workflow shape. Lower caps fail loudly
when an agent goes off the rails; higher caps mask runaway spend until the
LiteLLM report on Monday.

## Tiers (from PLAYBOOK §2)

| Cap | When to use |
|---|---|
| `$0.50` | One-shot lookup, single research stage, literature ping. |
| `$2.00` | `plan` + `critique` only. No `implement`. |
| `$5.00` | Full pipeline including `implement` and `test`. |
| `$10.00` | Hard ceiling without team-lead acknowledgement. |
| `> $10.00` | Postmortem trigger. Surface the team-lead-acknowledgement requirement before launching. |

## Heuristics

- **Doubling rule.** If the goal mentions a stack you do not know well, pick
  the next tier up. Wasted spend on under-capped failures usually exceeds the
  delta.
- **Sandboxed implement** rarely changes the cap — the bottleneck is model
  cost, not the proxy.
- **Caps are pre-set, not adjustable** post-launch without `just raise-cap`.
  Raising the cap requires an explicit user instruction; the orchestrator
  must not raise unilaterally.

## What "real" cost means

The cap is compared against `state.json.total_cost_usd`, which is filled
from LiteLLM `/spend/logs` (`cost_source: real`) when available, otherwise
from the runner's hardcoded estimate (`cost_source: estimate`). Estimates
drift; surface the source when reporting spend.

## Reporting cap exhaustion

When the dispatcher exits with rc=2 and ntfy says spend cap reached:

1. Run `just llm-spend 1` to confirm against LiteLLM.
2. Tell the user the gap (`$<spent>/$<cap>`) and the `cost_source` mix.
3. Offer two paths: raise the cap (user must request explicitly) or accept
   the partial result.
4. If `> $1` was spent on a failed stage, surface the postmortem trigger
   from PLAYBOOK §4 and hand off to `agent-workflow-postmortem`.
