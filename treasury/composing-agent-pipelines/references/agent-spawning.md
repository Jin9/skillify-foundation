# Agent Spawning Guide

The Agent tool exposes several `subagent_type` values. Each phase uses a specific type for a specific reason. This file is the rationale; the Modes table in `SKILL.md` is the authoritative mapping.

## Available subagent_types and their write capability

| Type | Documented tools | Empirical write behavior | Best for |
|------|------------------|--------------------------|----------|
| `Explore` | All tools except Agent, ExitPlanMode, Edit, Write, NotebookEdit | **Inconsistent.** In end-to-end testing, some Explore agents wrote artifacts and others returned content in their reply without writing. Cannot rely on Explore for write-required phases. | Pure read-only research where the orchestrator transcribes returns. |
| `Plan` | All tools except Agent, ExitPlanMode, Edit, Write, NotebookEdit | **Empirically appears to write** when prompted to produce a planning artifact at a specific path. The harness may grant special write capability for plan files. Behavior is undocumented and should not be relied on. | Reserved for explicit planning artifacts; not used by this skill due to documentation/behavior mismatch. |
| `general-purpose` | Full tool access including Edit and Write. Slowest and most expensive. | **Reliable.** Always writes when instructed. | All write-required phases of this pipeline. |

## Why this skill uses `general-purpose` for every spawned phase

Earlier versions of this skill prescribed `Plan` for the Plan/Review/Decide phases and `Explore` for Gather/Validate. End-to-end testing surfaced two issues:

1. **`Explore` writes inconsistently.** Out of 5 parallel Gather agents in the e2e test, only 2 wrote their evidence files; the other 3 returned content in chat replies. The orchestrator had to transcribe the missing 3 to disk. This is a reliability bug that compounds with N parallelism.
2. **`Plan`'s write capability is undocumented.** The Plan agent in the e2e test produced `01-plan.md` correctly, but the documented tool list excludes Write. Relying on undocumented behavior is fragile.

`general-purpose` is the only subagent type that reliably writes. The cost premium (typically a few hundred extra tokens per spawn vs. `Explore`) is worth the reliability gain in a 7-phase pipeline where partial-write failures cascade.

## Per-phase rationale

### Plan phase → `subagent_type: general-purpose`
Plan needs to produce `01-plan.md` with structured sub-questions. `general-purpose` writes reliably; `Plan` writes inconsistently per the documentation.

### Gather phase → `subagent_type: general-purpose` (parallel, N agents)
Each sub-question becomes one agent. The agent reads source files (Read/Grep) and writes a single evidence file (`02-evidence/qN.md`). Parallel because sub-questions are independent and the read work is the bottleneck.

### Analyze phase → `subagent_type: general-purpose`
Reads all evidence files and writes `03-analysis.md`. Single agent — synthesis is sequential by nature.

### Review phase → `subagent_type: general-purpose` (adversarial prompt)
Reads `03-analysis.md` and writes `04-review.md`. The agent is given a devil's-advocate framing in the prompt. The agent must NOT edit the analysis itself — only write the review file.

### Validate phase → `subagent_type: general-purpose` (1 agent, all claims)
Reads analysis + review and writes `05-validation.md` containing one section per claim. Earlier versions specified N parallel `Explore` agents (one per claim), but two issues forced the change: `Explore` does not reliably write, and parallel writes to the same file would race. One sequential agent is simpler and demonstrably works.

### Decide phase → `subagent_type: general-purpose`
Reads analysis (+ review + validation if present) and writes `06-decision.md`. Single agent.

### Compact phase → no spawn (inline)
The orchestrator (main agent) reads everything and writes `07-final.md` + `summary.md` directly. No subagent needed; the main agent already has the context.

### Compose mode → dispatcher only
Compose itself spawns nothing. It reads the manifest, calls the appropriate phase mode, and runs `update_manifest.py` at phase boundaries.

## Parallelism rules

Only the **Gather** phase parallelizes. Issue all N `Agent` tool calls in a single message — the runtime parallelizes them. Sequential issuance defeats the design.

All other phases spawn a single agent, including Validate (one agent handles every claim sequentially).

## Background mode

Background (`run_in_background: true`) is only acceptable for the Gather phase, and only when the user explicitly requests it (e.g., "this will take a while, run it in the background"). Record the returned background id in the manifest via `update_manifest.py --add-agent-call "general-purpose|Gather qN: ..."` (the script accepts the call before the agent completes; the orchestrator must reconcile completion later). All other phases run synchronously.

## Worktree isolation

`isolation: worktree` is forbidden by default. Allow only for Gather, and only when the user explicitly asks (e.g., they're worried about read-side-effects). The phase must `ExitWorktree` before declaring done; leftover worktrees are a leak.

## Anti-patterns

- Never call `Agent` with the same `subagent_type: <this-skill-name>` (no recursive Compose).
- Never spawn N agents sequentially when N > 1 — single message, multiple calls.
- Never give an agent permission to edit `manifest.json` — only the orchestrator (via `scripts/update_manifest.py`) writes the manifest.
- Never use `Explore` for a phase that must produce a written artifact. Use `general-purpose`.
- Never assume `Plan` will write; if a phase needs a written artifact, use `general-purpose`.
