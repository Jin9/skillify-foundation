# Routing Protocol — Generic Rules

This file is the authoritative source for routing rules, sandbox boundaries, lifecycle gates, and troubleshooting. The orchestrator skill points here for operational detail. Deployment-specific bindings (which channel, which board, which identity files) live in `openclaw-deployment.md`.

## Routing rules

- **Phase order is enforced.** Requirements → Architecture → Implementation → Validation → Human approval. A routing command that skips a phase MUST be rejected and surfaced to the human tech lead.
- **One routing command per artifact transition.** Do not bundle multiple agents into one command. Each `templates/routing-command.json` instance routes exactly one artifact to exactly one agent for one phase.
- **Routing commands are read-only events.** They do not mutate source files. Any file mutation happens inside the target agent's sandbox.
- **The Developer agent receives at most one approved design at a time.** Concurrent implementation against the same bounded context is a coordination failure; queue or merge before routing.
- **The QA agent is never bypassed.** Every artifact bound for human review MUST first carry a QA verdict (pass / conditional / reject).

## Sandbox boundaries

- **BA and Architect agents are read-only.** They may read source, references, and prior artifacts. They MUST NOT write to source files. Their outputs are markdown specifications and architecture documents persisted to the tracking board.
- **Developer agent runs in an isolated sandbox** (e.g., a container, ephemeral worktree, or restricted-scope environment). The sandbox MUST forbid network access to production systems and MUST forbid CI/CD trigger access.
- **QA agent runs in a read-only inspection environment** with the Developer's artifacts mounted. It MUST NOT mutate source files and MUST NOT execute deployment commands.
- **No agent has production credentials.** Production secrets are accessible only via the human DevOps role after explicit approval.

## Lifecycle gates

- **Soft gate after BA, Architect, Developer, QA**: surface a one-line phase summary to the tracking board and proceed unless the human halts.
- **Hard gate before deployment**: pause and route to the configured review channel. A natural-language acknowledgement from the configured human tech lead role is required. Auto-approval rules (timeouts, presence checks, "✓" reactions without explicit role match) MUST NOT satisfy this gate.

## Coverage and quality thresholds

- Unit-test coverage target on Developer outputs: **≥ 80%** of changed lines.
- QA verdict states: `pass`, `conditional` (must list mitigations), `reject` (must list blocking findings). Only `pass` and `conditional with mitigations applied` may proceed to human review.
- A `reject` verdict re-routes to the Developer agent with the QA findings; do not iterate the QA agent more than twice on the same artifact without escalating to the human tech lead.

## Failure modes

| Signal | Action |
|--------|--------|
| Agent process or sandbox container is unresponsive | Mark phase failed on the tracking board; alert the human DevOps role. Do not silently retry. |
| Missing approval after the configured timeout | Surface a reminder in the review channel and pause routing. Do not auto-approve. |
| Sandbox-isolation breach detected | Immediately halt routing, revoke the Developer sandbox tokens, alert DevOps, and refuse to emit further routing commands until cleared. |
| Scope creep detected (artifact exceeds the BA's stated requirements) | Re-route to the BA agent with the diff highlighted; force re-clarification before resuming. |
| QA verdict missing or empty | Treat as `reject` and route back to the Developer agent for QA-evidence regeneration. |
| Two consecutive QA rejects on the same artifact | Escalate to the human tech lead; do not auto-iterate. |
| Routing command attempts to skip a phase | Reject the command, write a routing-violation entry to the tracking board, surface to the human tech lead. |

## Audit trail

Every routing command and approval request MUST be persisted to the tracking board with: artifact id, source phase, target phase, target agent, timestamp, actor (orchestrator or human role), and outcome. The audit trail is append-only.
