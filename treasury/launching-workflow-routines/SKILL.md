---
name: launching-workflow-routines
description: >-
  Trigger a pre-designed multi-node workflow routine BY NAME from a routines/
  registry: resolve, script-validate, show an overview run plan, then execute
  node-by-node with human gates, script-verified artifacts under
  tmp/runs/routines/, and a final run report. Light modes: list routines;
  dry-run shows the plan only. Use when the user says "run the X routine",
  "start the X routine", "trigger the X workflow routine", "list my routines",
  "what workflow routines do I have", "dry-run the X routine", or "show the
  run plan for X". Nodes dispatch existing skills as black boxes at
  small/mid/frontier tiers; a node never launches another routine. Do NOT use
  to design workflows or gates (agentic-workflow-design), the fixed in-session
  7-phase pipeline (composing-agent-pipelines), the platform-bound scaffold
  launcher (orchestrating-agent-scaffold), a single external-model consult
  (delegating-to-cli-models), or single hardcoded runners like
  jira-fix-mr-workflow: those are dispatch targets, not this skill.
compatibility: claude-code, codex, copilot, gemini, antigravity
version: 0.1.0
---

# Launching Workflow Routines

## Purpose

Run one pre-designed multi-node workflow routine, selected BY NAME from a registry, end to end: the agent resolves the definition, validates it deterministically, shows the plan, then executes each node by dispatching its executor skill as a black box — with human gates, script-verified artifacts, and a final run report. The launcher is a thin dispatcher: routines are data, executors do the work, the host agent is the runtime.

## Contract

- **Preconditions:** a resolvable registry (see Registry resolution) and a routine file that exits 0 from `scripts/validate_routine.py`.
- **Postconditions:** a new run directory `tmp/runs/routines/ROUTINE-YYYYMMDD-HHMMSS/` containing every declared node artifact verified by `scripts/check_run.py`, plus `INDEX.md` and `run-report.md`. List and dry-run modes write nothing.
- **Failure output:** a partial `run-report.md` naming the failed node, per-node honesty statuses, the verification evidence, and one lesson per failure.

## When to use / when not

Use for: "run the X routine", "start the X routine", "trigger the X workflow routine", "list my routines", "dry-run the X routine", "show the run plan for X".

Not this skill:

- Designing a workflow, its stages, or its gates — use agentic-workflow-design.
- The fixed in-session 7-phase multi-agent pipeline — use composing-agent-pipelines.
- Operating the agent-scaffold platform workflow — use orchestrating-agent-scaffold.
- One advisory prompt to an external model CLI — use delegating-to-cli-models.
- A hardcoded runner invoked directly (jira-fix-mr-workflow, running-ba-pipeline, and similar) — use that runner. Such runners ARE valid executors inside a routine.
- Authoring or repairing a routine definition — copy `templates/routine-template.md`, edit by hand, validate. This skill never edits routines.

## Registry resolution

Search in order; first hit wins:

1. A registry path or routine file the user names explicitly.
2. `routines/` at the current project root.
3. The skill's own `examples/` directory (ships `examples/demo-inventory-digest.md`).

Routine files are `*.md` excluding `INDEX.md`. The registry's `INDEX.md` table is the human catalog; glob the directory when it is absent or stale. List mode merges all three locations and labels each routine's source.

## Modes

| Mode | Trigger wording | Effect |
|------|-----------------|--------|
| list | "list my routines", "what workflow routines do I have" | show the catalog, write nothing |
| dry-run | "dry-run the X routine", "show the run plan for X" | Steps 1-3 only, write nothing |
| run | "run the X routine", "start the X routine", "trigger the X workflow routine" | full workflow |

Detection rules: dry-run wording wins over run wording when both appear. A run or dry-run request whose name matches no routine falls back to list, then asks which routine is meant. Never guess between near-matches — show them and ask.

## Workflow

Follow the numbered steps. Gate wording, dispatch idiom, tier meanings, honesty statuses, and failure detail are normative in `references/run-protocol.md`; the routine schema is normative in `references/routine-format.md`. Both scripts live in this skill's `scripts/` directory — invoke them by absolute path with `python3`.

### Step 0 — List (list mode only)

- Entry: list mode. Print one table: routine, summary, node count, gate count, version, source file.
- Exit: table shown. END.

### Step 1 — Resolve

- Entry: run or dry-run mode with a name candidate.
- Match the candidate against `routine:` frontmatter values and filename stems, in registry search order.
- Exit: exactly one routine file. On zero or multiple matches: show the candidates or the list, ask which is meant, and STOP.

### Step 2 — Validate

- Run `scripts/validate_routine.py` on the resolved file.
- Exit: exit code 0 (warnings are allowed and shown).
- On errors: print the script output verbatim, do NOT execute, do NOT edit the routine file. END.

### Step 3 — Run plan (overview first)

- Render the plan table: order, node id, purpose, executor with mode, tier, gate, on_fail, declared outputs.
- Below it: run-level inputs with values — ask for any required input still missing, apply stated defaults for optional ones — plus `requires:` preconditions, gate count, and failure policy.
- Exit: plan shown, all required inputs collected. Dry-run ENDS here, having written nothing.

### Step 4 — Gate 0: launch confirmation

- Ask the Gate 0 question exactly as worded in the Gates section of `references/run-protocol.md`.
- Exit ONLY on explicit approve. Abort ends with nothing written. Inspecting a node loops back to this gate.

### Step 5 — Create the run directory

- Create `tmp/runs/routines/ROUTINE-YYYYMMDD-HHMMSS/` (UTC timestamp, never reuse a directory).
- Write its `INDEX.md` seeded with the run plan and a status table, every node `pending`.
- Exit: run dir and `INDEX.md` exist.

### Step 6 — Node loop (strictly sequential, file order)

For each node:

1. Mark it `running` in the run dir `INDEX.md`.
2. Resolve declared inputs: earlier artifacts must exist in the run dir, `user.*` values come from Step 3, `external:` paths must exist. A missing input is a node failure — go to substep 9.
3. If `gate: before` or `both`: show exactly what will be dispatched (executor, tier, resolved inputs, expected outputs); wait for approve, skip, or abort.
4. Runtime recursion check: refuse dispatch if the executor is this skill or names a routine file.
5. Dispatch the executor as a black box per `references/run-protocol.md`: sub-agent when the host supports one, inline otherwise; pass only the node's input contract and exact output paths. `agent-inline` nodes the agent performs itself at the node's tier.
6. Verify with `scripts/check_run.py` using `--node` for this node id. Never accept the executor's own success claim.
7. Record the honesty status (`built` or `planned`; upgrade earlier nodes to `wired` when this node consumed their artifacts) in `INDEX.md`.
8. If `gate: after` or `both`: show the artifact summary plus the verification RESULT line; wait for approve, redo (one human-corrected re-dispatch maximum), or abort.
9. On failure: apply `on_fail` — the node's value, else `default_on_fail`, else stop. Stop marks all remaining nodes `skipped` and jumps to Step 7; continue records the failure and proceeds. NO automatic retries.

### Step 7 — Final sweep

- Run `scripts/check_run.py` over the whole run dir (no `--node`) and capture its RESULT line.

### Step 8 — Run report

- Fill `templates/run-report-template.md` as `run-report.md` in the run dir; finalize every `INDEX.md` status.
- Re-verify with `scripts/check_run.py` adding `--final`; exit 0 is required for a `complete` verdict.
- Present: verdict, the RESULT line quoted verbatim, per-node honesty statuses, artifact paths, failures with lessons.
- Exit: `run-report.md` exists. END.

## Output contract

```
tmp/runs/routines/ROUTINE-YYYYMMDD-HHMMSS/
├── INDEX.md          # seeded at Step 5, status-updated through the run
├── run-report.md     # written at Step 8 from the template
├── 01-NODE.md        # one artifact per node, exactly as declared
└── ...               # any further declared outputs
```

The launcher writes nothing outside this directory and never modifies the routine definition or the registry.

## Guards

- Nodes run strictly sequentially, in file order. No parallel groups in v1.
- NO automatic retries; the only re-execution is one human-approved redo at an after-gate. Re-launching a routine always means a fresh run dir.
- A routine node MUST NOT launch another routine, and an executor must never be this skill — enforced by the validator and re-checked at dispatch.
- Never skip Gate 0. Only explicit approval counts — silence is not approval.
- Never fabricate or paraphrase verification output; quote the script's RESULT line.
- Never edit a routine definition, in any mode, even to fix a validation error — report the error instead.
- Stay vendor-neutral in plans and reports: tiers (small, mid, frontier), never model names.
- v1 non-goals: parallel node execution, automatic retries, nested routines, routine authoring.

## Examples

Full run:

> User: run the demo-inventory-digest routine
> Agent: resolves `examples/demo-inventory-digest.md`; validation ok; shows the 2-node plan and asks for target_dir; Gate 0 — user replies approve; creates the run dir; node inventory runs inline at small tier, verification RESULT: pass, after-gate approved; node digest likewise; final sweep "RESULT: pass nodes=2 artifacts=2 missing=0"; writes run-report.md — verdict complete, both nodes built, 01-inventory.md wired.

List, then dry-run:

> User: what workflow routines do I have
> Agent: one table — three registry routines plus the packaged demo, with node and gate counts and sources.
> User: show the run plan for research-squad-chain
> Agent: 6-node plan (plan, search, extract, synthesize, review, report-run) with tiers, gates, inputs topic/depth/audience, and the web-search precondition. Nothing written.

## References

| Need | File |
|------|------|
| Routine schema (authoritative) | `references/routine-format.md` |
| Dispatch, gates, honesty, failure detail | `references/run-protocol.md` |
| New routine skeleton | `templates/routine-template.md` |
| Run report skeleton | `templates/run-report-template.md` |
| Registry INDEX skeleton | `templates/registry-index-template.md` |
| Definition validator | `scripts/validate_routine.py` |
| Artifact verifier | `scripts/check_run.py` |
| Shared parser | `scripts/routine_parser.py` |
| Smoke-test routine | `examples/demo-inventory-digest.md` |
