# Run protocol

Normative detail for executing a validated routine. The numbered workflow lives in `SKILL.md`; this page defines the dispatch idiom, gates, honesty statuses, and failure handling it uses.

## Dispatch idiom (black box)

For each node, dispatch the executor without coupling to its internals:

- If the host supports sub-agents, dispatch to a sub-agent instructed to apply the executor skill (with `executor_mode` when set) at the node's tier. Otherwise the agent applies the executor skill inline, in sequence.
- Pass ONLY: routine name, node id, `purpose`, executor skill + mode, the tier meaning, the resolved input values/paths, and the exact output paths to write inside the run dir.
- Never pass other nodes' internals, the rest of the routine, or earlier gate discussions. If the executor needs something not in the input contract, the routine definition is wrong — stop and report, do not improvise.
- `agent-inline` nodes are performed directly by the launching agent at the declared tier's effort level.
- Runtime no-recursion check: refuse dispatch if the executor equals `launching-workflow-routines` or names a routine file — even if validation was somehow skipped.

## Tier meanings

| Tier | Use | Effort |
|---|---|---|
| `small` | mechanical transforms, listings, compaction, bookkeeping | cheapest capable model class, minimal reasoning |
| `mid` | structured drafting, digestion, search, coordination maps | balanced class |
| `frontier` | synthesis, judgment, adversarial review, high-stakes drafting | highest-capability class available |

When the host cannot vary model class, tiers degrade to effort guidance for the single available model.

## Gates

Explicit human approval — silence, timeouts, or the agent's own judgment are NEVER approval. Wording:

- **Gate 0 (launch)** — after showing the run plan: "Reply approve to launch, abort to cancel, or name a node to inspect first."
- **`gate: before`** — before dispatch: show exactly what will be dispatched (executor, tier, resolved inputs, expected outputs) and ask "approve to dispatch this node, skip to mark it skipped, or abort the run."
- **`gate: after`** — after verification: show the artifact summary plus the `check_run.py` result and ask "approve to continue, redo to re-dispatch this node once with your correction, or abort."
- A `redo` at an after-gate is the only permitted re-execution, it requires the human's correction in hand, and it overwrites the node's artifacts before re-verification.

Record every gate decision (gate, node, decision, one-line human rationale if given) in the run report's gates log.

## Honesty statuses

Reported per node in `INDEX.md` and `run-report.md`; derived from `check_run.py` evidence, never from executor claims:

| Status | Meaning |
|---|---|
| `planned` | declared but not (yet) real: artifact missing, empty, or a stub under 64 bytes |
| `built` | artifact exists, is substantive, and passed script verification |
| `wired` | built AND actually consumed as an input by a later executed node |
| `failed` | node dispatch or verification failed |
| `skipped` | not executed (human skip at a gate, or an earlier `stop` failure) |

`INDEX.md` progress vocabulary during the run: `pending`, `running`, then one of the five statuses above.

## Failure handling

- Missing declared input at dispatch time = node failure before any dispatch happens.
- Node failure = executor error, refusal, or `check_run.py --node` exit 1 after execution.
- Apply `on_fail` (node value, else `default_on_fail`, else `stop`): `stop` marks all remaining nodes `skipped` and jumps to the final sweep and run report; `continue` records the failure and proceeds — the final RESULT is still `fail` if any declared artifact is missing.
- NO automatic retries. One human-approved `redo` per after-gate is the only exception. Otherwise re-launching the routine is a fresh run in a fresh run dir.
- Every failure gets one line in the run report: what failed, the evidence line from the script, and one lesson.

## Run directory contract

- Path: `tmp/runs/routines/<routine-name>-<YYYYMMDD-HHMMSS>/` (UTC timestamp), created at Step 5, never reused.
- Contents: one artifact per node as declared, `INDEX.md` (seeded with the run plan, updated at every status change), `run-report.md` (written last from `templates/run-report-template.md`).
- The launcher never writes outside the run dir, never edits the routine definition, and never deletes anything.
