# Pre-flight `subagent_type` mapping

The scaffold's stage runners (`scripts/stage-runners/<stage>.sh`) talk to
their configured models directly (Gemini for research, Codex for plan and
implement, Claude for critique and review). The orchestrator's job is to
spawn *Claude Code* sub-agents *before* `just workflow` to reduce wasted
spend or *between* stages when the user asks for a richer summary.

## Spawn rules

1. Spawn pre-flight only when the goal is ambiguous, the repo is unfamiliar,
   or the cap is tight enough that a single misfire wastes >25% of the
   budget. For routine work, kick off the scaffold directly.
2. Spawn between stages only on explicit user request ("summarize the plan
   for me", "what are the top three risks in the critique?").
3. Do not spawn sub-agents *during* a running stage to help with that stage —
   the runner is already executing the agent that owns that work.
4. One sub-agent per goal unless the work is genuinely parallelizable; only
   `Explore` is safe to fan out (4–6 in one message at most).

## subagent_type → job

| Job | `subagent_type` | Use when |
|---|---|---|
| Find a file, symbol, or wiring point | `Explore` | "Where is auth wired?", "Which configs reference Kong?" |
| Draft a high-level plan to seed the scaffold | `Plan` | The goal is a refactor or feature spanning multiple modules. |
| Open-ended investigation (>3 tool calls) | `general-purpose` | "What's the simplest way to migrate from X to Y here?" |
| Hand-off to a security-focused review | invoke `expert-software-security-reviewer` skill | The goal touches auth, KYC, disbursement, secrets, or PII. |
| Hand-off to backend code review | invoke `crafting-backend-code` skill | The goal is a backend refactor or API design choice. |
| Hand-off to a Claude-native multi-phase pipeline | invoke `composing-agent-pipelines` skill | The user explicitly asked for the Claude pipeline, not the scaffold. |
| Read changed-code review | (no scaffold equivalent) — use Claude's `/review` slash if installed | Post-implement readout of the diff. |

## What pre-flight outputs feed into

Sub-agent output goes to the orchestrator's reply, never directly into
`.agent/`. If the user wants to seed the scaffold's `plan` stage, route the
output through a saved prompt (`prompts/library/plan/<topic>.md`) via the
`drafting-stage-prompt` skill — do not hand-edit `.agent/`.

## When NOT to spawn

- The user has already typed the goal precisely and the cap is generous.
- The previous stage's artifact is fresh and the user just wants to know
  status.
- Spawning would duplicate work the running stage is already doing.

## Reporting sub-agent results

After every pre-flight call, surface to the user:

- One-paragraph summary.
- Decision: launch as planned, adjust profile/cap, or abort.
- The agent's findings are advisory — they do not bind the scaffold runner.
