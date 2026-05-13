---
name: composing-agent-pipelines
description: >
  Composes a portable multi-agent pipeline (Plan, Gather, Analyze, Review,
  Validate, Decide, Compact) for research, code review, implementation planning,
  or trade-off analysis. Use when the user asks to "compose agent pipeline",
  "run multi-agent pipeline", "spawn agents for this task", "audit with multiple
  agents", "plan and review with agents", "decide with agent help", "research
  this with agents", or wants modular phase invocation. Each phase emits a
  versioned artifact under .agent-pipelines/[task-id]/. Do NOT use for
  single-agent tasks, editing repo-policy files, or strict human-in-the-loop
  approval workflows (use orchestrating-openclaw-squad instead).
---

# Composing Agent Pipelines

## Purpose

Drive a structured, multi-phase agent pipeline that turns one task into
evidence-backed deliverables. The skill decomposes a task into sub-questions
(Plan), delegates bounded work to host-native worker agents when available
(Gather and other write-required phases), synthesizes findings (Analyze), runs
an adversarial critique (Review), fact-checks claims against source (Validate),
produces recommendations (Decide), and compacts the result into a final report
(Compact). A Compose mode runs the full pipeline end-to-end, adapting the phase
set to the task domain. Each mode is also invocable standalone.

## When to use this skill

- Use when a task is too large for a single agent and needs structured decomposition with citations.
- Use when the user wants reproducible artifacts (plan, evidence, analysis, review, decisions) on disk for later inspection or hand-off.
- Use when the answer must survive an adversarial critique pass (Review) or be fact-checked against source (Validate).
- Use after the user names a target domain — research, code review/audit, implementation planning, or decision/trade-off analysis.
- Triggers on: "run the pipeline", "compose phases", "spawn pipeline", "gather then analyze", "decompose this task", "agent pipeline", "research with agents".
- Do NOT use for trivial single-step questions, single-file edits, or one-shot lookups — answer directly or use a single host-native helper.
- Do NOT use to invoke another instance of this skill (no recursive Compose).
- Do NOT use to edit `AGENTS.md`, `CLAUDE.md`, or repo-policy files.

## Modes

The pipeline has eight modes, dispatched from user phrasing or chained by Compose:

- **Plan** — "decompose this", "break into sub-questions". Hard approval gate. Writes `01-plan.md`.
- **Gather** — "gather sources", "scout codebase". The only parallel phase (N workers when available). Writes `02-evidence/qN.md` per sub-question.
- **Analyze** — "synthesize", "find patterns". Writes `03-analysis.md`.
- **Review** — "critique this", "devil's advocate". One pass, no loop. Writes `04-review.md`.
- **Validate** — "fact-check claims", "verify against source". One pass, no loop. Writes `05-validation.md`.
- **Decide** — "recommend a path", "compare options". Hard approval gate. Writes `06-decision.md`.
- **Compact** — "produce final report". Inline. Writes `07-final.md` + `summary.md`.
- **Compose** — "run the full pipeline", "end-to-end". Dispatches per domain shape.

Full per-mode prompt skeletons, input/output paths, exit gates, and delegation
choices are in `references/mode-playbooks.md` — read that before running any
mode beyond the universal preamble. Domain → phase shapes (which modes Compose
runs for each domain) are in `references/domain-shapes.md`.

## Universal preamble

Run before every mode.

1. **Detect mode** from the user's phrasing using the Modes table. If the user only described a task with no mode hint, default to Compose.
2. **Detect domain** if invoking Compose. Map phrasing to one of:
   - General research / investigation
   - Code review / audit
   - Implementation planning
   - Decision / trade-off analysis
   If ambiguous, ask one disambiguating question rather than guess. Print the chosen domain shape from `references/domain-shapes.md` so the user can override before phases start.
3. **Establish input contract**:
   - Compose / Plan: collect the user's task prompt and target domain.
   - Standalone phase: ask for the path to the upstream artifact if not stated. Do not fabricate a path.
4. **Initialize pipeline state**: run `scripts/init_pipeline.py --slug <short-slug>` to create `<cwd>/.agent-pipelines/<task-id>/` with a `manifest.json`. The script prints the absolute task directory; record it for subsequent phase calls. To resume an existing pipeline, pass `--resume <task-id>`.
5. **Announce output contract**: list the exact artifact paths the chosen mode will write. Wait for user confirmation when the mode includes a hard gate (Plan, Decide).
6. **Mark phase running**: run `scripts/update_manifest.py --task-id <id> --phase <name> --status running` before delegation or inline work. Record each delegated call with `--add-delegation "<worker_type>|<description>"`.
7. **Run the mode** following its playbook. Use the host's delegation mechanism only when the user explicitly asked for multi-agent work and the host supports it; otherwise run the phase inline and record `inline|<description>`. Only Gather is parallel by default; all other phases are single-worker or inline.
8. **Exit through the validation gate** below before proceeding. On success, run `scripts/update_manifest.py --task-id <id> --phase <name> --status done` to mark the phase complete; on gate failure, mark `--status failed` and surface the blocker.

## Output contract

Every phase writes artifacts under `<cwd>/.agent-pipelines/<task-id>/`:

```
.agent-pipelines/<task-id>/
├── manifest.json        # phase status, timestamps, delegated calls
├── 01-plan.md
├── 02-evidence/
│   ├── q1.md
│   └── q2.md
├── 03-analysis.md
├── 04-review.md
├── 05-validation.md
├── 06-decision.md
├── 07-final.md
└── summary.md
```

Format details, manifest schema, and resumption rules: see `references/state-convention.md`.

The final user-facing return for Compose is `07-final.md` (full report) plus `summary.md` (TL;DR with citations). For standalone phase invocation, the return is the single artifact that phase produced.

## Quality gates

Every mode exits through these gates.

1. **Artifact written** to the path declared in the Modes table; manifest updated to `done`.
2. **Citations preserved**: every finding/claim/recommendation references at least one upstream artifact path or source file. Phases that produced no evidence must explicitly say so.
3. **Hard approval gates** at Plan and Decide: pause and surface the artifact, do not auto-proceed.
4. **Soft gates** at Gather, Analyze, Review, Validate: surface a one-line summary and proceed unless the user halts. Cap Review and Validate at one pass each — do not loop.
5. **Anti-pattern sweep**: see `references/anti-patterns.md` before declaring the mode done.

If a gate fails after one retry, stop and surface the blocker rather than guessing further.

## Constraints

- DO NOT invoke Compose from inside any phase (no recursion).
- DO NOT loop Review or Validate beyond one pass — escalate remaining issues to the user.
- DO NOT use `isolation: worktree` for any phase by default; only Gather or Validate may use it when explicitly requested by the user, and the worktree must be exited after the phase.
- DO NOT use background execution for synchronous phases. Background mode is reserved for explicitly long-running Gather only, and the manifest must record the background worker id when the host exposes one.
- DO NOT reuse a task-id without `--resume`. The init script refuses to overwrite existing state.
- DO NOT include destructive shell commands in any artifact or example.
- DO NOT duplicate phase mechanics here that live in `references/mode-playbooks.md`.

## Troubleshooting

| Signal | Action |
|--------|--------|
| User invoked a phase without an upstream artifact | Ask for the artifact path; do not fabricate one. |
| Sub-questions exceed 8 in Plan output | Push back to the user — pipeline is over-decomposing. Suggest narrowing scope. |
| Review surfaces > 0 P1 issues | Do not skip Validate. Wait for user to address P1 or explicitly accept. |
| Validate marks a P1 claim as refuted | Loop back to Analyze with the validation report; do not proceed to Decide. |
| Phase worker returned an empty artifact | Surface the failure, do not write a placeholder. Ask user how to proceed. |
| Pipeline state directory already exists | Refuse and ask the user to either pass `--resume <task-id>` or pick a new slug. |

## References

| Need | Reference |
|------|-----------|
| Per-mode prompt skeletons, parallelism, exit gates | `references/mode-playbooks.md` |
| Domain → phase subset mapping | `references/domain-shapes.md` |
| Pipeline directory layout and manifest schema | `references/state-convention.md` |
| Delegation and inline fallback rules | `references/delegation.md` |
| Pipeline-specific failure modes | `references/anti-patterns.md` |

## Templates and scripts

- `templates/01-plan.md` – Plan artifact skeleton.
- `templates/02-evidence.md` – Evidence-bundle skeleton (per sub-question).
- `templates/03-analysis.md` – Analysis artifact skeleton.
- `templates/04-review.md` – Review artifact skeleton (P1/P2/P3).
- `templates/05-validation.md` – Validation artifact skeleton (claim verdicts).
- `templates/06-decision.md` – Decision artifact skeleton (options table).
- `templates/07-final.md` – Final compacted report skeleton.
- `scripts/init_pipeline.py` – Create `.agent-pipelines/[task-id]/` with manifest.
- `scripts/pipeline_status.py` – Read manifest.json and print phase status.
- `scripts/update_manifest.py` – Mark phases running/done/failed and append delegation records.

## Examples

- `examples/research-walkthrough.md` – End-to-end research pipeline.
- `examples/code-review-walkthrough.md` – Code-review pipeline with all phases.
- `examples/decision-walkthrough.md` – Decision pipeline (skip Review, keep Validate).
