---
name: openclaw-orchestrator
description: Orchestrates the OpenClaw multi-agent squad (BA, Architect, Developer, QA) via Discord for enterprise banking workflows. Enforces strict human-in-the-loop validation, sandbox constraints, and routing protocol. Do NOT use for direct production deployment or single-agent trivial tasks.
---

# OpenClaw Multi-Agent Orchestrator

## Purpose

To manage and securely route tasks across the OpenClaw squad of specialized AI agents (BA, Architect, Developer, QA) and 1 Human DevOps Engineer via a Discord hub. This skill enforces enterprise banking-grade security, KYC/AML compliance checks, transactional integrity, adversarial testing, and explicit human-in-the-loop approvals.

## When to use this skill

- Orchestrating multi-stage enterprise banking task execution (e.g., refactoring a legacy Loan Origination System).
- Routing deliverables safely between the BA, Architect, Developer, and QA agents.
- Managing human-in-the-loop code review cycles in the `#output-review` Discord channel.
- Operating the Mission Control Kanban board to track multi-agent progress.
- Trigger terms: "OpenClaw squad", "orchestrate agents", "route to QA", "human approval", "god-pm".

## Modes

### `route`
Direct tasks or intermediate deliverables to the appropriate specialized agent based on the phase (Requirements → Architecture → Implementation → Security).

### `review`
Enforce the human tech-lead approval gate in the Discord `#output-review` channel before permitting any artifact to proceed to the next lifecycle phase or deployment.

## Core workflow

1. **Intake & Discovery (BA Agent)**: Receive feature request and route to the Business Analyst agent to clarify requirements with strict KYC/AML/PCI-DSS compliance mapping.
2. **Architecture (Architect Agent)**: Route the BA specs to the Architect agent for deep codebase reasoning, transaction integrity design, and concurrency safety mapping.
3. **Implementation (Developer Agent)**: Route the approved design to the Developer agent inside its isolated Docker sandbox to execute code changes, write unit tests (≥80% coverage), and implement circuit breakers.
4. **Validation (QA Agent)**: Pass the developer's artifacts to the QA Engineer for adversarial security tests (OWASP), concurrency checks, and chaos testing.
5. **Human Approval & Deployment**: Route the final, QA-passed artifact to the `#output-review` Discord channel for explicit Human DevOps deployment.

## Output format

Produce structured routing commands for the Discord communication layer, status updates for the Mission Control Kanban board, and clear, structured approval requests in the `#output-review` channel.

## Constraints

- DO NOT grant any agent CI/CD deployment access; only the Human DevOps Engineer deploys to production.
- DO NOT allow the BA or Architect agents to execute writing operations (enforce Read-Only boundary).
- DO NOT allow the Developer agent to execute code outside its isolated Docker Sandbox.
- DO NOT bypass or automatically approve the `#output-review` explicit human step.
- DO NOT duplicate identity definitions; reference the canonical identity files in `identities/`.

## Troubleshooting

| Signal | Action |
|--------|--------|
| Agent stalled | Verify the respective agent's process or sandbox container is running via Mission Control. |
| Missing approval | Ping the designated human tech lead in the Discord `#output-review` channel. |
| Sandbox breach | Immediately halt routing and revoke Developer sandbox tokens; alert DevOps. |
| Scope creep | Re-route back to the BA agent to force requirements clarification and constraint checking. |

## Validation gate

Before sending the response or executing a routing command, re-check:

1. The deliverable has successfully passed through the required prior agent (e.g., cannot route to QA if the Developer phase is incomplete).
2. The current routing mode (`route` vs `review`) matches the task lifecycle state.
3. No production deployment or out-of-sandbox execution commands are included in the generated output.
4. If proceeding to final deployment, explicitly verify that Human DevOps approval has been registered in the `#output-review` channel.
