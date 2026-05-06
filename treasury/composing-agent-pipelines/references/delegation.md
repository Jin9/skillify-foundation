# Delegation Guide

Use this file to map pipeline phases onto the host agent's delegation
capabilities. The skill is portable: it does not require a specific agent tool
name or vendor-specific worker type.

## Worker Classes

| Class | Required capability | Best for |
|---|---|---|
| `writer` | Read relevant inputs and write exactly one declared artifact path. | Plan, Analyze, Review, Validate, Decide, and Gather evidence files. |
| `reader` | Read/search only and return findings to the orchestrator. | Optional evidence collection when the orchestrator will write the artifact. |
| `inline` | No delegation; the orchestrator performs the phase directly. | Fallback when delegation is unavailable or the user did not ask for multi-agent work. |

If the host exposes named worker roles, choose the narrowest role that satisfies
the required capability. For Codex-style hosts, `worker` maps to `writer` and
`explorer` maps to `reader`. For hosts with only one general worker, use it as
`writer` only when the worker can write the declared artifact.

## Default Mapping

| Phase | Default | Rationale |
|---|---|---|
| Plan | `writer` or `inline` | Produces `01-plan.md`; hard gate follows. |
| Gather | parallel `writer` workers, one per sub-question | Sub-questions are independent; each worker writes one `02-evidence/qN.md`. |
| Analyze | `writer` or `inline` | Synthesis is sequential and writes `03-analysis.md`. |
| Review | `writer` or `inline` | One adversarial pass writes `04-review.md`. |
| Validate | `writer` or `inline` | One sequential pass writes `05-validation.md`; parallel validators would race. |
| Decide | `writer` or `inline` | One decision artifact writes `06-decision.md`. |
| Compact | `inline` | The orchestrator reads artifacts and writes `07-final.md` plus `summary.md`. |

## Parallelism Rules

Only Gather parallelizes by default. Launch all Gather delegations in one batch
when the host supports parallel calls. If the host cannot parallelize, run the
Gather prompts sequentially but still keep one output file per sub-question.

All other phases are single-worker or inline. Review and Validate each run at
most once per pipeline.

## Recording Rules

Record every delegated or inline phase boundary with
`scripts/update_manifest.py`:

```bash
python3 scripts/update_manifest.py --task-id <id> --phase gather \
  --status running \
  --add-delegation "writer|Gather q1: auth middleware evidence"
```

Use `inline|<description>` when no worker is delegated.

## Worker Prompt Rules

- The worker may write only the artifact path named in its prompt.
- The worker must not edit `manifest.json`; only the orchestrator updates it.
- The worker must not invoke `composing-agent-pipelines` recursively.
- The worker must not broaden its scope beyond its phase input and declared artifact.
- A `reader` worker must return findings to the orchestrator instead of writing source files.

## Background And Isolation

Background execution is allowed only for Gather and only when the user explicitly
asks. Record the background id in the manifest when the host exposes one.

Worktree or sandbox isolation is disabled by default. Allow it only for
read-heavy Gather work when the user asks, and close the isolated environment
before declaring the phase done.
