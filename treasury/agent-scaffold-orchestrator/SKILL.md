---
name: agent-scaffold-orchestrator
description: >
  Orchestrates a multi-stage research-squad workflow on top of the agent-scaffold
  (just / state.json / zellij / LiteLLM / ntfy / sandbox). Plans the run, picks
  profile and cap, spawns fit-for-job sub-agents for pre-flight, kicks off the
  workflow via just, then monitors state.json and stage logs. Use when the user
  says "kick off the workflow", "start an agent-scaffold run", "what's the
  workflow status", "research squad workflow", "stage failed", "it's stuck",
  "show me the plan / critique / review", "abort the workflow", "resume the
  workflow", or pastes a goal that maps onto research → plan → critique → 🔒
  implement → review → test. Honors AGENTS.md: never auto-approves, never edits
  .agent/, never runs production commands. Do NOT use for the Claude-native
  pipeline (use composing-agent-pipelines), the BA/Architect/Dev/QA squad (use
  openclaw-orchestrator), or for deciding to approve implement (use
  reviewing-implement-gate).
---

# agent-scaffold orchestrator

## Purpose

Drive a goal through the agent-scaffold pipeline by planning the run, spawning
fit-for-job sub-agents where richer reasoning helps, dispatching the workflow
through `just`, and monitoring state without ever editing it. The skill owns
*how Claude operates the scaffold*; the scaffold owns the actual stage
execution and the model API calls.

## When to use this skill

- A goal should run through the scaffold's stages (research → plan → critique → 🔒 implement → review → test).
- The user asks for status, logs, or stage output of an active workflow.
- A workflow has stalled, failed, or hit an approval gate.
- The user wants help choosing a profile, cap tier, or sandbox flag.

Do NOT use this skill to:
- Decide *whether to approve* an implement gate → `reviewing-implement-gate`.
- Edit `.agent/`, `prompts/library/`, profiles, or sandbox allowlists → relevant sibling skill.
- Run prod-affecting commands (forbidden by `AGENTS.md`).

## Hard constraints (from AGENTS.md)

Non-negotiable. Refuse and explain if pushed.

1. NEVER `just approve <stage>` on the user's behalf — humans type approvals.
2. NEVER edit `.agent/` files; state mutations go through the dispatcher.
3. NEVER `just raise-cap` without explicit user instruction.
4. NEVER start a parallel stage while one is `running` or `blocked`.
5. NEVER modify source code directly — that is the `implement` stage's job.
6. NEVER run `git push`, `terraform apply`, `kubectl apply`, DB migrations, or anything else that touches production.

## Universal preamble

Run before every operation.

1. Confirm the working directory is an agent-scaffold checkout: `.agent/`,
   `justfile`, and `dispatch.sh` exist. If not, stop and tell the user
   the skill needs the scaffold present.
2. Read `.agent/state.json` and `.agent/status.md` to ground the next response
   in current state. Never assume yesterday's state.
3. Detect the operation from the user's phrasing using the Modes table below.
   If multiple match, prefer the most recently spoken intent and ask one
   disambiguating question rather than guess.

## Modes

| Mode | Trigger phrasing | What this skill does |
|---|---|---|
| Plan-and-launch | "kick off", "start a workflow", "run X through the scaffold" | Plan profile/cap/sandbox, optionally spawn pre-flight Agent calls, confirm with user, run `just workflow`. |
| Status | "what's the status", "where are we", "show me progress" | Read `state.json` + tail current stage log; one-paragraph summary. |
| Stage-readout | "show me the plan / critique / review / research" | Cat `.agent/stages/<stage>.md` and summarize. |
| Failure-triage | "stage failed", "it's stuck", "it errored" | Read failed-stage log, surface verbatim error, suggest fixes (do not auto-retry). |
| Gate-surface | "approval pending", "what's blocked" | Surface the gate; for `implement`, hand off to `reviewing-implement-gate`. |
| Resume | "resume", "continue", "it's been hanging" | `just resume` (recovers stale PIDs) and re-summarize. |
| Abort | "stop everything", "abort", "kill it" | `just reject <stage>` if a gate is pending; otherwise `just abort`. Confirm before invoking. |

## Plan-and-launch workflow

1. **Restate the goal** in one sentence, ≤ 500 chars, no newlines.
2. **Choose profile** — `code` (default), `literature` (no implement/test), `dataset` (no test). See `references/profile-decision.md`.
3. **Choose cap** — smallest tier that fits, per `references/cap-tiers.md`. Anything above $10 requires explicit team-lead acknowledgment.
4. **Decide sandbox** — `IMPLEMENT_SANDBOXED=1` for credential / customer-data / fintech repos. Default to sandboxed when in doubt.
5. **Pre-flight (optional)** — spawn `Explore` / `Plan` / `general-purpose` sub-agents *before* `just workflow` when the goal is ambiguous or repo-unfamiliar. Pick `subagent_type` per `references/agent-spawning.md`. Skip for routine, well-scoped goals.
6. **Confirm the plan** with the user using `templates/launch-confirmation.md`. Wait for explicit "go".
7. **Kick off** with a kebab-case `workflow_id` (`^[a-z][a-z0-9-]{0,31}$`):
   ```bash
   IMPLEMENT_SANDBOXED=<0|1> WORKFLOW_PROFILE=<name> \
     just workflow <id> "<goal>" <cap>
   ```
8. **Switch to monitoring.**

## Monitoring loop

While a stage is `running` or `blocked`, the orchestrator watches; it does
not "help". Read `.agent/state.json` plus the redacted tail of the current
stage log on demand — never on a fixed timer. Suggest `just goto <stage>`
when `state.json` exposes a `zellij_tab` field. When `approval_pending`
becomes non-null, surface the gate; for `implement`, hand off to
`reviewing-implement-gate`. Cadence rules and summary shapes:
`references/monitoring-loops.md`.

## Failure triage

When a stage flips to `failed`: read the redacted log tail, surface the
failing command verbatim (no paraphrased stack traces), suggest **one** fix
path, and stop. Do not auto-retry. Hand off to `drafting-stage-prompt` for
prompt fixes or `agent-workflow-postmortem` when the trigger conditions
fire (see `references/failure-triage.md` for the classification table).

## Output format

For every response, lead with a one-line state summary in the shape:

```
<workflow_id>: <current_stage> <status> (<elapsed>m, $<spent>/$<cap>)
```

Then one of: a status paragraph, a launch-confirmation block (template), a
stage-readout summary, or a failure-triage block. Always cite the file path
the information came from (e.g., `state.json`, `stages/critique.log:42`).

## Tool inventory

| Need | Tool | Notes |
|---|---|---|
| Read state | `Read` on `.agent/state.json`, `status.md`, `stages/*.md`, `stages/*.log` | jq via Bash for filtered reads. |
| Run scaffold verbs | `Bash` running `just <verb>` | Never bypass `just`; never call the dispatcher directly. |
| Pre-flight reasoning | `Agent` tool with `subagent_type` from `references/agent-spawning.md` | Pick by job, not by name familiarity. |
| Cost truth | `Bash` running `just llm-spend <hours>` | LiteLLM `/spend/logs` is authoritative; `state.json` cost is a snapshot. |
| Tail live logs | `Bash` running `tail -n 50 .agent/stages/<stage>.log` | Logs are already redacted at write time. |

## Validation gate

Before responding, verify:

1. The user's request matches one of the Modes; if not, ask a disambiguating
   question rather than guess.
2. No hard-constraint violation in the planned action (no auto-approve, no
   `.agent/` writes, no prod commands).
3. Every claim about state cites the file it was read from.
4. If a stage is `running` or `blocked`, no parallel stage is being started.
5. If kicking off a new workflow, the user explicitly confirmed the
   launch-confirmation block.

## Constraints

- DO NOT auto-approve gates, edit `.agent/`, or run prod-affecting commands.
- DO NOT spawn sub-agents *during* a `running` stage to "help" — wait.
- DO NOT fabricate `state.json` field values; read the file.
- DO NOT loop status checks unprompted.
- DO NOT raise the spend cap without an explicit user request.
- DO NOT use `run_in_background` for pre-flight Agent calls unless the user
  asked; pre-flight findings should land before launch confirmation.

## Troubleshooting

| Signal | Action |
|---|---|
| `state.json` missing or schema mismatch | Run `just doctor`; surface output. Do not write a fresh state file. |
| Stage `running` for hours, pid alive | Suggest `just goto <stage>` for live view. Do not abort without user say-so. |
| Stage `running` for hours, pid dead | `just resume` (auto-recovers stale PIDs). |
| `approval_pending` set but user not notified | Check `NTFY_TOPIC` env, surface the gate in the response. Do not retry ntfy. |
| Spend cap tripped | Run `just llm-spend 1`, confirm real cost; raise-cap requires explicit user request. |
| Two researchers on the same checkout | Surface `lsof .agent/state.json.lock`; suggest a git worktree. |

Full failure runbook: `references/failure-triage.md`.

## References

| Need | Reference |
|---|---|
| Profile selection (code/literature/dataset/custom) | `references/profile-decision.md` |
| Cap tier selection | `references/cap-tiers.md` |
| Pre-flight `subagent_type` mapping | `references/agent-spawning.md` |
| Monitoring cadence and summary shapes | `references/monitoring-loops.md` |
| Failure triage flowchart | `references/failure-triage.md` |
| State.json field cheatsheet | `references/state-cheatsheet.md` |

## Templates and scripts

- `templates/launch-confirmation.md` — pre-launch confirmation block.
- `templates/status-summary.md` — status-paragraph skeleton.
- `templates/failure-triage.md` — failure-response skeleton.
- `scripts/state_summary.sh` — jq one-liner for the canonical status line.
