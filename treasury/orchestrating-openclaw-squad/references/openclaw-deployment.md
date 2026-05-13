# OpenClaw Deployment Binding

This file documents the OpenClaw squad's concrete deployment of the multi-agent orchestrator. It is one supported configuration; other deployments may map the abstract "review channel", "tracking board", and "squad identities" to different infrastructure.

## Squad roster

The OpenClaw squad consists of four AI agents and one human role:

| Role | Identity file | Sandbox | Permissions |
|------|---------------|---------|-------------|
| Business Analyst (BA) | `identities/ba.md` | Read-only | Reads source and references; writes specifications to the tracking board only. |
| Architect | `identities/architect.md` | Read-only | Reads source, references, and BA specs; writes architecture documents to the tracking board. |
| Developer | `identities/developer.md` | Isolated Docker sandbox; no production network access | Mutates source inside the sandbox; commits to a feature branch. |
| QA Engineer | `identities/qa.md` | Read-only inspection environment with Developer artifacts mounted | Writes QA verdicts and findings to the tracking board. |
| Human DevOps Engineer (Tech Lead) | n/a | n/a | Sole authority to approve deployment to production. |

`identities/` is maintained at the deployment level (project-root or squad-config repo) and is not duplicated inside this skill.

## Communication infrastructure

- **Hub:** Discord workspace dedicated to the OpenClaw squad.
- **Review channel (hard approval gate):** `#output-review` — the only channel where human deployment approval is registered.
- **Routing channels (soft phase transitions, audit trail):**
  - `#ba-intake` — feature requests routed to the BA agent.
  - `#architecture` — BA specs routed to the Architect agent.
  - `#implementation` — approved designs routed to the Developer agent.
  - `#qa` — Developer artifacts routed to the QA agent.
- **Status broadcasts:** `#status` — orchestrator emits one bullet per phase status change.

## Tracking board

- **Mission Control Kanban** (Notion / Trello / Jira project, deployment-specific) with columns mirroring lifecycle phases: `Backlog`, `BA Intake`, `Architecture`, `Implementation`, `QA`, `Awaiting Approval`, `Deployed`.
- Each card carries: artifact id, current phase, agent on duty, timestamps, attached routing commands, attached approval requests, audit trail.
- Cards are append-only. Phase reverts (e.g., QA reject → back to Implementation) create a new entry, never overwrite.

## Approval policy

- The human tech-lead role explicitly named in the squad config is the only role whose approval satisfies the hard gate.
- Approval is registered by the tech lead posting an explicit acknowledgement in `#output-review` referencing the artifact id. Reactions, presence checks, or timeouts do not satisfy the gate.
- An approval covers exactly the artifact id named; it does not implicitly extend to subsequent artifacts.

## Deployment authority

- Production deployment is triggered manually by the human DevOps engineer after `#output-review` approval.
- The orchestrator MUST NOT emit deployment commands and MUST NOT route deployment artifacts to any AI agent.
- Sandbox-token rotation, secret rotation, and CI/CD pipeline configuration are owned by the human DevOps role; agents have no access.

## Mapping to the generic protocol

When `routing-protocol.md` references abstract terms, the OpenClaw deployment maps them as follows:

| Abstract term | OpenClaw binding |
|---------------|------------------|
| Review channel | Discord `#output-review` |
| Tracking board | Mission Control Kanban |
| Squad identity | `identities/{ba,architect,developer,qa}.md` |
| Configured human tech-lead role | the named role in the squad config (single named human) |
| Routing channel | the per-phase Discord channels listed above |

A different deployment can keep this skill unchanged and supply its own deployment file.
