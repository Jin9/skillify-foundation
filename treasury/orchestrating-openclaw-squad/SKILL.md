---
name: orchestrating-openclaw-squad
description: Orchestrates a multi-agent squad (Business Analyst, Architect, Developer, QA Engineer) with strict human-in-the-loop validation, sandbox isolation, and explicit approval gates for enterprise and regulated workflows. Use when the user says "OpenClaw squad", "orchestrate agents", "route to QA", "human approval gate", "multi-agent enterprise workflow", "god-pm", or asks to manage a Requirements → Architecture → Implementation → Validation lifecycle with a human tech-lead gate. Do NOT use for generic multi-agent decomposition without a strict human approval gate (use `composing-agent-pipelines`), single-agent trivial tasks, or direct production deployment commands.
---

# Orchestrating the OpenClaw Squad

## Purpose

Manage and securely route tasks across a squad of specialized AI agents (BA, Architect, Developer, QA) plus a human DevOps approval gate. Enforces enterprise-grade security posture, compliance mapping, transactional integrity, adversarial testing, and explicit human approvals for any artifact bound for production.

## When to use this skill

- Orchestrating multi-stage enterprise task execution (e.g., refactoring a regulated banking system, modernizing a legacy origination flow).
- Routing deliverables safely between BA → Architect → Developer → QA agents.
- Managing human tech-lead approval cycles before any artifact proceeds to deploy.
- Tracking multi-agent progress on a shared coordination / kanban board.
- Trigger terms: "OpenClaw squad", "orchestrate agents", "route to QA", "human approval gate", "god-pm", "multi-agent enterprise workflow".

## When NOT to use

- Generic multi-agent decomposition without a strict human approval gate — use `composing-agent-pipelines`.
- Single-agent tasks with no multi-stage routing — invoke the appropriate specialist agent directly.
- Direct production deployment commands — only the human DevOps gate may approve deployment.

## Modes

### `route`

Direct a task or intermediate deliverable to the appropriate specialized agent based on lifecycle phase (Requirements → Architecture → Implementation → Validation). Emit a routing command using `templates/routing-command.json`.

### `review`

Enforce the human tech-lead approval gate in the configured review channel before permitting any artifact to proceed to the next phase or to deployment. Emit an approval request using `templates/approval-request.md`.

## Core workflow

1. **Intake & discovery (BA agent)**: receive the feature request and route to the Business Analyst agent for requirements clarification with strict compliance mapping (KYC / AML / PCI-DSS or domain-appropriate regulatory framework).
2. **Architecture (Architect agent)**: route the BA specs to the Architect agent for codebase reasoning, transaction-integrity design, and concurrency-safety mapping.
3. **Implementation (Developer agent)**: route the approved design to the Developer agent inside an isolated sandbox to execute code changes, write unit tests (≥ 80% coverage target), and implement defensive controls (circuit breakers, idempotency, retries).
4. **Validation (QA agent)**: pass the Developer's artifacts to the QA agent for adversarial security tests (OWASP Top 10), concurrency checks, and chaos test plans.
5. **Human approval & deployment**: route the final, QA-passed artifact to the configured review channel for explicit human approval. Only the human DevOps role triggers production deployment.

Routing rules, sandbox boundaries, and lifecycle gates are codified in `references/routing-protocol.md`. The OpenClaw-specific deployment binding (review-channel name, tracking-board schema, squad identity files) is in `references/openclaw-deployment.md`. Other deployments may map "review channel" and "tracking board" to different infrastructure.

## Output format

- **Routing command** — `templates/routing-command.json`. JSON message containing target agent, source artifact, lifecycle phase, and routing reason.
- **Approval request** — `templates/approval-request.md`. Markdown request containing artifact summary, QA verdict, residual risks, and the explicit approval ask.
- **Status update** — short markdown suitable for the tracking board (one bullet per phase status change).

## Constraints

- DO NOT grant any agent CI/CD or production-deployment access; only the human DevOps role may deploy.
- DO NOT allow the BA or Architect agents to execute write operations on source files (read-only boundary).
- DO NOT allow the Developer agent to execute code outside the isolated sandbox.
- DO NOT bypass or auto-approve the human review gate.
- DO NOT duplicate identity definitions; reference canonical identity files at the deployment level (see `references/openclaw-deployment.md`).

## Validation gate

Before sending the response or executing a routing command, re-check:

1. The deliverable has successfully passed through the required prior agent (cannot route to QA before Developer phase is complete).
2. The current routing mode (`route` vs `review`) matches the task lifecycle state.
3. No production-deployment or out-of-sandbox execution commands are included in the generated output.
4. If proceeding to final deployment, explicit human approval has been registered in the configured review channel.

## References

| Need | File |
|------|------|
| Generic routing constraints, sandbox boundaries, lifecycle gates, troubleshooting | `references/routing-protocol.md` |
| OpenClaw deployment specifics — channel names, tracking board, squad identity | `references/openclaw-deployment.md` |

## Templates

- `templates/routing-command.json` — routing-message schema.
- `templates/approval-request.md` — human-approval request format.
