# Treasury — Production Skill Library

93 top-level skills, grouped by purpose. Each skill name links to its folder. The **Extra assets** column lists files or folders beside `SKILL.md`.

| Purpose Group | Count |
|---------------|-------|
| [Banking, BA Delivery & Requirements](#banking-ba-delivery-requirements) | 14 |
| [Architecture, Engineering Decisions & Planning](#architecture-engineering-decisions-planning) | 19 |
| [Implementation, Platform Templates & Code Review](#implementation-platform-templates-code-review) | 12 |
| [Testing, QA & Validation](#testing-qa-validation) | 4 |
| [Agent Orchestration & Workflow Infrastructure](#agent-orchestration-workflow-infrastructure) | 10 |
| [Research, Debate & Knowledge Synthesis](#research-debate-knowledge-synthesis) | 14 |
| [Security, Governance & Compliance](#security-governance-compliance) | 5 |
| [Observability, Cost & Incident Operations](#observability-cost-incident-operations) | 6 |
| [Code Analysis, Productivity & Publishing](#code-analysis-productivity-publishing) | 9 |
| **Total** | **92** |

---

## Banking, BA Delivery & Requirements (14)

Business-analysis workflows for regulated banking and delivery handoffs.

| Skill | Purpose | Source family | Extra assets |
|-------|---------|---------------|--------------|
| [`analyzing-banking-requirements`](./analyzing-banking-requirements/) | Business Analyst persona for enterprise banking and lending workflows. Extracts raw business requests, runs strict KYC / AML / PCI-DSS compliance mapping, drafts... | business-analyse | references |
| [`assembling-tl-handoff`](./assembling-tl-handoff/) | Stage 5 of the BA pipeline: assemble the Handoff Bundle the Tech Lead signs for — the raw requirement plus the Scope Sheet, Story Set, clear Governance Check, and... | business-analyse | — |
| [`assessing-ba-feasibility`](./assessing-ba-feasibility/) | Stage 4 of the BA pipeline: consume the Story Set, a clear Governance Check, and the raw requirement, and emit a typed Feasibility Note contract (at least two... | business-analyse | — |
| [`checking-ba-governance`](./checking-ba-governance/) | Stage 3 of the BA pipeline: sweep a Story Set for governance and privacy gaps and emit a typed Governance Check contract (PII flags, compliance flags, and blockers... | business-analyse | — |
| [`drafting-ba-stories`](./drafting-ba-stories/) | Stage 2 of the BA pipeline: consume a confirmed Scope Sheet and emit a typed Story Set — one epic plus INVEST user stories, each a full user-story-template instance... | business-analyse | references · templates |
| [`eliciting-banking-brief`](./eliciting-banking-brief/) | Convert raw BA input (Jira, Slack, meeting notes, email, mixed prose) into a structured epic-plus-stories brief with banking-grade fields force-evaluated,... | business-analyse | audit · references · schemas · scripts · tests |
| [`evaluate-banking-compliance`](./evaluate-banking-compliance/) | Scan structured banking Epics and Stories for PII, regulatory dependencies, tipping-off risk, missing Legal or Privacy roles, and TL-handoff governance blockers.... | craft | references · templates · examples |
| [`extract-brief-structure`](./extract-brief-structure/) | Convert redacted raw BA input from Jira, Slack, meeting notes, emails, or mixed prose into a strict Epic and Story JSON skeleton. Use when a user asks to "parse... | craft | references · templates · examples |
| [`orchestrate-banking-brief-pipeline`](./orchestrate-banking-brief-pipeline/) | Run the decomposed banking BA brief pipeline across extraction, compliance, ambiguity, and Gherkin stages with validated JSON handoffs and pluggable model commands.... | craft | references · scripts · templates · examples |
| [`running-ba-pipeline`](./running-ba-pipeline/) | Run the five-stage agentic Business-Analysis pipeline that turns a raw requirement into a Tech-Lead-ready handoff: it carries one task_id and a typed envelope... | business-analyse | references |
| [`running-business-analysis-workflow`](./running-business-analysis-workflow/) | Guide an AI agent through a structured business-analysis workflow for software and IT projects: define the problem and goals, analyze the current (as-is) state,... | craft | references · templates |
| [`scoping-ba-intake`](./scoping-ba-intake/) | Stage 1 of the BA pipeline: ingest a raw requirement (Jira, email, meeting notes, document) and emit a typed Scope Sheet contract wrapped in the pipeline's shared... | business-analyse | — |
| [`scoping-technical-requirements`](./scoping-technical-requirements/) | Turns a messy business requirement into a clear, bounded technical scope a team can safely act on, surfacing the open questions and non-functional requirements that... | business-analyse | — |
| [`sweep-ambiguities`](./sweep-ambiguities/) | Detect linguistic ambiguities and hidden requirements in a structured banking brief using the eight ambiguity detectors and ten elicitation frames. Use when a user... | craft | references · templates · examples |

## Architecture, Engineering Decisions & Planning (19)

Design, modeling, planning, and engineering decision support.

| Skill | Purpose | Source family | Extra assets |
|-------|---------|---------------|--------------|
| [`api-contract-design`](./api-contract-design/) | Design a stable, clear API contract (action, request/response, error cases, validation, idempotency, versioning) for an API you expose, so consumers can rely on it... | business-analyse | — |
| [`architecting-fintech-systems`](./architecting-fintech-systems/) | Senior fintech domain architect for lending, loan origination, KYC, credit decisioning, disbursement, and regulated financial systems. Owns L1–L3 design (business... | squad-delivery | references · templates |
| [`architecture-decision`](./architecture-decision/) | Chooses a design among options and records it as an ADR with explicit trade-offs, protecting long-term maintainability and setting clear direction the team can... | business-analyse | — |
| [`data-modeling`](./data-modeling/) | Design the persistence model safely (schema, source of truth, indexes, constraints) with an expand/contract migration and rollback plan that preserves backward... | business-analyse | — |
| [`defining-engineering-standards`](./defining-engineering-standards/) | Define and own the team's reusable standard — tech-stack baseline, DDD domain boundaries, a C4 L3/L4 design skeleton, API naming and folder/layer conventions, and... | business-analyse | — |
| [`delivery-planning`](./delivery-planning/) | Break a chosen solution into a sequenced, estimated execution plan — tasks, critical path, blockers, and a phased timeline honest enough to tell the business. Use... | business-analyse | — |
| [`designing-tech-lead-handoff`](./designing-tech-lead-handoff/) | Convert an approved BA epic-and-stories brief plus a UX design pack into the full Tech-Lead architecture handoff: integration contracts, component map, infra spec,... | squad-delivery | audit · references · schemas · templates · tests |
| [`domain-modeling`](./domain-modeling/) | Defines domain boundaries and DDD building blocks so ownership is clear, the model matches the business, and domain logic stays separate from infrastructure. Use... | business-analyse | — |
| [`engineer-growth-planning`](./engineer-growth-planning/) | Grows an engineer — runs useful 1-on-1s, chooses teach-by-doing vs teach-by-telling, spots underperformance early, delegates to build capability, and reviews... | business-analyse | — |
| [`engineering-doc-planning`](./engineering-doc-planning/) | Decides what to document and what to deliberately skip, for which audience, keeping docs minimal, findable, and actionable across ADRs, C4 design docs, API... | business-analyse | — |
| [`generate-ux-pack`](./generate-ux-pack/) | Produce a v1.1 UX-design intake pack from a UX team's drop (bundled prototype HTML, Frontend Spec markdown, BA brief directory). Emits a structured... | squad-delivery | references · schemas · RATIONALE.md |
| [`integration-design`](./integration-design/) | Design resilient integration with an external/3rd-party system (timeouts, retries, fallback, error mapping, idempotency, and a clear ownership/support model) so... | business-analyse | — |
| [`model-selection`](./model-selection/) | Assign each agent or stage in an agentic workflow a model and reasoning effort by criteria — a privacy/data-class gate, capability-to-role match,... | business-analyse | — |
| [`pr-design-review`](./pr-design-review/) | Review a PR's design and maintainability — business-logic completeness, test coverage, over-engineering, and template conformance — returning tagged, teachable... | business-analyse | — |
| [`production-readiness`](./production-readiness/) | Verify a feature is safe to ship — observability, rollback, runbook, migration safety, and a passing smoke test — and return a clear go / no-go / conditional... | business-analyse | — |
| [`refactor-decision`](./refactor-decision/) | Decide whether a refactor is necessary, bound its scope, protect existing behavior with tests first, and avoid cosmetic refactors that have no business reason. Use... | business-analyse | — |
| [`risk-estimation`](./risk-estimation/) | Size complexity and separate known work from unknown risk, producing a man-day estimate with an explicit confidence and uncertainty multiplier plus the assumptions... | business-analyse | — |
| [`technical-debt-management`](./technical-debt-management/) | Classify technical debt, prioritize it by business risk and cost-of-delay, make the dangerous debt visible, and negotiate a fix budget with PM/business. Use when... | business-analyse | — |
| [`technical-feasibility`](./technical-feasibility/) | Decides whether and how a requirement is buildable on the current stack, surfacing the options, dependencies, and risks before anyone commits to a design or a date.... | business-analyse | — |

## Implementation, Platform Templates & Code Review (12)

Implementation, platform template, refactoring, and code-review skills.

| Skill | Purpose | Source family | Extra assets |
|-------|---------|---------------|--------------|
| [`crafting-backend-code`](./crafting-backend-code/) | Reviews, designs, and safely implements backend code and microservice templates with a pattern-first, evidence-led posture. Use when designing, reviewing,... | existing treasury | references |
| [`crafting-frontend-code`](./crafting-frontend-code/) | Reviews, designs, and safely implements frontend code with a conservative, repo-first posture for AI coding agents. Use when designing, reviewing, optimizing,... | existing treasury | references |
| [`crafting-rust-code`](./crafting-rust-code/) | Reviews, designs, and safely implements Rust code with a pattern-first, evidence-led, repo-first posture for AI coding agents. Use when designing, reviewing, optimizing, fixing, analyzing, or planning Rust services, crates, async/tokio code, error-type design, sqlx persistence, money/decimal... | existing treasury | references |
| [`implement-backend-feature`](./implement-backend-feature/) | Generate production-grade Go backend code for one microservice feature from an approved design document. Use when implementing a Go HTTP handler from a design spec.... | squad-delivery | references · schemas · tests · RATIONALE.md |
| [`implement-frontend-feature`](./implement-frontend-feature/) | Generate production-grade React/TypeScript code for one frontend feature from an approved UI design, with banking-grade discipline: WCAG 2.1 AA a11y, no any outside... | squad-delivery | references · schemas · tests · RATIONALE.md |
| [`implementing-go-template-requirements`](./implementing-go-template-requirements/) | Apply a single requirement (spec line, ticket, bug report, user story) to a Go service that follows the go-template scaffold by editing ONLY business logic under... | agentic | references · templates · examples |
| [`platform-common`](./platform-common/) | Shared Go infrastructure library for A-Team Krungthai DGL microservices — Gin middleware, Kafka producer/consumer, JWT, structured slog logging, response envelopes,... | agentic | references |
| [`platform-go-service`](./platform-go-service/) | Scaffold or extend a Go microservice in this repository (go-template) using A-Team platform conventions: DDD aggregate-per-package, CQRS handler/consumer split,... | agentic | references |
| [`refactoring-go-services`](./refactoring-go-services/) | Incrementally refactor messy Go microservices toward clean DDD/CQRS architecture while preserving behavior. Use when the user asks "clean up this Go service"... | agentic | examples · references |
| [`review-backend-code`](./review-backend-code/) | Adversarially verify Go backend code emitted by a Generate stage against the approved design and the 11 banking-grade decision rules + v2 augmentations, then issue... | agentic | references · schemas · tests · RATIONALE.md |
| [`review-frontend-code`](./review-frontend-code/) | Adversarially verify React/TypeScript code emitted by a Generate stage against the approved UI design and the 12 banking-grade frontend non-negotiables + 9 v2... | squad-delivery | references · schemas · tests · RATIONALE.md |
| [`review-rust-code`](./review-rust-code/) | Adversarially verify Rust backend code emitted by a Generate stage against the approved design and the 11 banking-grade decision rules + v2 augmentations, then... | existing treasury | references · schemas · tests · RATIONALE.md |

## Testing, QA & Validation (4)

Test planning, acceptance criteria, QA, and implementation validation.

| Skill | Purpose | Source family | Extra assets |
|-------|---------|---------------|--------------|
| [`generating-gherkin-acceptance-criteria`](./generating-gherkin-acceptance-criteria/) | Translate resolved banking stories into strict BDD Gherkin acceptance criteria with happy, error, audit, idempotency, and tipping-off-safe scenarios. Use when a... | craft | references · templates · examples |
| [`planning-banking-tests`](./planning-banking-tests/) | Convert a BA brief (output of eliciting-banking-brief v1.2+) into a structured QA test plan covering every Gherkin scenario, banking-grade concern, NFR target,... | squad-delivery | references · schemas · scripts · tests |
| [`testing-strategy`](./testing-strategy/) | Decide what to test and at which level (unit / integration / contract / e2e / regression) and map a requirement's business logic into concrete test scenarios, so... | business-analyse | — |
| [`validating-banking-implementation`](./validating-banking-implementation/) | QA Engineer persona for enterprise banking implementation artifacts. Adversarial OWASP Top 10 testing, transactional-integrity audit, race / deadlock analysis, and... | squad-delivery | references |

## Agent Orchestration & Workflow Infrastructure (10)

Multi-agent orchestration, scaffold configuration, prompts, gates, and workflow infrastructure.

| Skill | Purpose | Source family | Extra assets |
|-------|---------|---------------|--------------|
| [`agent-context-initializer`](./agent-context-initializer/) | Generate a minimal, human-curatable AGENTS.md (under 200 lines, target 100-150) for an agentic-squad repo: exactly six sections (commands, testing, project... | existing treasury | references · scripts · templates |
| [`agentic-workflow-design`](./agentic-workflow-design/) | Design how an AI-agent pipeline is supervised — its stages and owners, the human-in-the-loop approval gates, the never-do guardrails and command-safety policy,... | business-analyse | references |
| [`authoring-scaffold-profile`](./authoring-scaffold-profile/) | Author or modify an agent-scaffold profile file at profiles/NAME.sh. Sets STAGES, GATED_STAGES, and per-stage AGENT / MODEL / PROMPT_PREFIX env vars. Validates that... | existing treasury | references · templates |
| [`composing-agent-pipelines`](./composing-agent-pipelines/) | Composes a portable multi-agent pipeline (Plan, Gather, Analyze, Review, Validate, Decide, Compact) for research, code review, implementation planning, or trade-off... | existing treasury | references · scripts · templates · examples |
| [`developing-langgraph-workflows`](./developing-langgraph-workflows/) | Guide LangGraph v1.x implementation, refactoring, and review workflows. Use when implementing, refactoring, or reviewing LangGraph StateGraph orchestration... | langgraph-claude-agent | references · templates |
| [`drafting-stage-prompt`](./drafting-stage-prompt/) | Curate stage prompts for the agent-scaffold prompt library at prompts/library/STAGE/TOPIC.md. Captures the Friday prompt-review ritual: a stage prompt with a... | existing treasury | references · templates |
| [`multi-agent-handoff-architect`](./multi-agent-handoff-architect/) | Design inter-agent handoff APIs as versioned JSON Schema contracts (required taskId, intent, state, confidence, provenance/trace, schemaVersion) with binary... | existing treasury | references · schemas · scripts · templates |
| [`orchestrating-agent-scaffold`](./orchestrating-agent-scaffold/) | Orchestrates a multi-stage research-squad workflow on top of the agent-scaffold (just / state.json / zellij / LiteLLM / ntfy / sandbox). Plans the run, picks... | existing treasury | references · scripts · templates |
| [`orchestrating-openclaw-squad`](./orchestrating-openclaw-squad/) | Orchestrates a multi-agent squad (Business Analyst, Architect, Developer, QA Engineer) with strict human-in-the-loop validation, sandbox isolation, and explicit... | existing treasury | references · templates |
| [`reviewing-implement-gate`](./reviewing-implement-gate/) | Walk a researcher through the five-question check before approving the implement gate in an agent-scaffold workflow. Reads .agent/stages/plan.md and... | existing treasury | references · templates |

## Research, Debate & Knowledge Synthesis (14)

Research pipeline, debate, source synthesis, and knowledge-vault workflows.

| Skill | Purpose | Source family | Extra assets |
|-------|---------|---------------|--------------|
| [`augment-diagrams`](./augment-diagrams/) | OPTIONAL out-of-band stage (NOT part of the standard 1–6 pipeline). Backfill 0–3 inline ASCII diagrams into an existing final research report, insertion-only,... | squad-researcher | references · examples |
| [`cross-examine`](./cross-examine/) | Stage 3 of the squad-brainstorm workflow, run once per panel CLI per debate round (skipped entirely when rounds=quick). The invoked panelist reads the OTHER two... | squad-brainstorm | references · examples |
| [`extract-findings`](./extract-findings/) | Stage 3 of the researcher workflow. Given a research_plan (from plan-research) and sources (from search-sources), pull out structured, source-grounded claims... | squad-researcher | references · examples |
| [`frame-debate`](./frame-debate/) | Stage 1 of the squad-brainstorm workflow. Read a finished squad-researcher run (05-final_report.md + 03-findings.json + 01-research_plan.json) and distill it into a... | squad-brainstorm | references · examples |
| [`opening-debate-panel`](./opening-debate-panel/) | Stage 2 of the squad-brainstorm workflow, run once per panel CLI (P-codex via `codex exec`, P-gemini via `gemini -p`, P-claude via `claude -p`). Given the debate... | squad-brainstorm | references · examples |
| [`plan-research`](./plan-research/) | Stage 1 of the squad-researcher workflow. Decompose a research topic into a structured plan — thesis question, 3–10 MECE sub-questions (count gated by depth),... | squad-researcher | examples · references |
| [`report-debate`](./report-debate/) | Stage 6 (meta) of the squad-brainstorm workflow. Reads prior-stage artifacts from the debate directory plus an optional metrics.json and produces a single... | squad-brainstorm | examples · references |
| [`reporting-research-run`](./reporting-research-run/) | Stage 6 (meta) of the squad-researcher workflow. Reads prior-stage artifacts from the run directory plus an optional metrics.json sidecar, and produces a single... | squad-researcher | examples · references |
| [`research-vault-librarian`](./research-vault-librarian/) | Read-only librarian and context-only workflow guide for a CLAUDE.md-governed Obsidian research-report vault (the ResearchVault: deep-research reports in domain... | existing treasury | references · scripts |
| [`review-report`](./review-report/) | Self-review pass over a draft research report. Verifies every cited claim is grounded in the input findings, checks audience and topic fit, surfaces... | squad-researcher | references · examples |
| [`revise-positions`](./revise-positions/) | Stage 4 of the squad-brainstorm workflow, run once per panel CLI per debate round (skipped when rounds=quick). The invoked panelist reads ONLY the critiques aimed... | squad-brainstorm | references · examples |
| [`search-sources`](./search-sources/) | Stage 2 of the researcher workflow. Given a structured research_plan (output of plan-research) and a depth knob, gather candidate sources for every sub-question and... | squad-researcher | examples · references |
| [`synthesize-consensus`](./synthesize-consensus/) | Stage 5 of the squad-brainstorm workflow. A single Claude moderator reads the attributed final panelist positions plus the cross-examination exchanges and produces... | squad-brainstorm | references · examples |
| [`synthesize-report`](./synthesize-report/) | Compose a markdown research report from structured findings, tuned to a named audience, with inline numeric citations and a sources list. Stage 4 of 6 in the... | squad-researcher | references · examples |

## Security, Governance & Compliance (5)

Security review, policy, sandboxing, spec validation, and governance controls.

| Skill | Purpose | Source family | Extra assets |
|-------|---------|---------------|--------------|
| [`configuring-sandbox-allowlist`](./configuring-sandbox-allowlist/) | Edit the agent-scaffold sandbox hostname allowlist at docker/sandbox-proxy/filter. Adds, reviews, or removes hostname patterns, refusing wildcards, host-internal... | existing treasury | references · templates |
| [`devops-infrastructure-hardener`](./devops-infrastructure-hardener/) | Audit an agentic CI/CD or agent-skill codebase for exposed static credentials and emit a remediation plan plus a runtime-secrets architecture that replaces them... | existing treasury | references · scripts · templates |
| [`governance-policy-generator`](./governance-policy-generator/) | Emit policy-as-code agent-governance artifacts architecturally separated from the agent they govern: an Open Policy Agent (OPA) Rego allowlist with default-deny... | existing treasury | references · scripts · templates |
| [`reviewing-software-security`](./reviewing-software-security/) | Defensive software security review for Go/Gin services, Kafka event flows, MySQL/RDS, Kubernetes workloads, Kong/APISIX gateways, and DDD/CQRS event-driven lending... | squad-delivery | references · templates |
| [`universal-spec-validator`](./universal-spec-validator/) | Validate an agent spec (tool/function JSON Schema, MCP tool manifest, OpenAPI tool def, SKILL.md frontmatter, or data contract) as a CI / pre-commit ENFORCEMENT... | existing treasury | references · scripts · templates |

## Observability, Cost & Incident Operations (6)

Telemetry, cost governance, incident handling, readiness, and postmortems.

| Skill | Purpose | Source family | Extra assets |
|-------|---------|---------------|--------------|
| [`authoring-workflow-postmortem`](./authoring-workflow-postmortem/) | Author a postmortem for an agent-scaffold workflow within 48h of a trigger event: failed stage with more than $1 spent, rejected implement gate, cap overrun, or... | existing treasury | references · templates |
| [`incident-response`](./incident-response/) | Stabilizes a firing production incident — classify severity, identify blast radius, drive a short-term mitigation, communicate, then run a blameless postmortem with... | business-analyse | — |
| [`observability-design`](./observability-design/) | Designs how a feature is observed in production — correlation/request/business IDs, logs, metrics, traces, dashboards, and alerts — so production behavior is... | business-analyse | — |
| [`observability-telemetry-instrumenter`](./observability-telemetry-instrumenter/) | Instrument agent code with OpenTelemetry GenAI semantic conventions: an invoke_agent CLIENT root span, execute_tool INTERNAL child spans, the... | existing treasury | references · schemas · scripts · templates |
| [`performance-cost-review`](./performance-cost-review/) | Define and review a feature or service's performance budget and cost budget — latency/throughput/resource targets per critical user journey, cost-per-unit, baseline... | business-analyse | examples · references |
| [`reviewing-agent-spend`](./reviewing-agent-spend/) | Run the Monday cost-review ritual for an agent-scaffold squad. Wraps just llm-spend, segments by user and model, flags greater-than-2-sigma outliers and any single... | existing treasury | references · templates |

## Code Analysis, Productivity & Publishing (9)

Code analysis, productivity, translation, publishing, and human-readable report generation.

| Skill | Purpose | Source family | Extra assets |
|-------|---------|---------------|--------------|
| [`business-logic-extractor`](./business-logic-extractor/) | Extract the implemented business logic and rules FROM a codebase — cross-referenced with requirements and agent/execution traces — into a faithful, traceable... | existing treasury | references · scripts · templates |
| [`daily-planner`](./daily-planner/) | Turn a raw task list into a prioritized, trackable daily plan and keep it current across days. Ranks by impact x urgency (P1–P4) with S/M/L effort and maintains a status-tracked... | skillify | templates |
| [`generating-pseudocode`](./generating-pseudocode/) | Analyzes requirements or existing code to generate clean, language-agnostic pseudocode that bridges high-level intent and implementation, readable by a Python, Go,... | existing treasury | references |
| [`organizing-local-files`](./organizing-local-files/) | Plans and applies safe, offline organization of local files and folders. Use when the user says "organize this folder", "clean up my files", "classify my... | existing treasury | references · templates |
| [`progressive-bug-hunter`](./progressive-bug-hunter/) | Localize and diagnose a bug by progressively retrieving the MINIMAL sufficient code context: start with cheap agentic grep / structured search, escalate to... | existing treasury | references · scripts · templates |
| [`publishing-git-review-requests`](./publishing-git-review-requests/) | Publishes one local change to GitHub or GitLab for hosted review: prepare an intended commit, create or attach origin, push the branch, and open a PR or MR. Use... | existing treasury | references · templates |
| [`rendering-readable-html`](./rendering-readable-html/) | Render structured content the agent already has - tabular data, a markdown or plain-text report, or findings produced this session - into one clean, self-contained,... | existing treasury | references · scripts · templates · examples |
| [`thai-translator`](./thai-translator/) | Translates English documents into structurally identical bilingual (Thai/English) documents. Use when the user asks to 'translate to thai', 'mirror this document in... | business-analyse | — |
| [`transcribing-media-to-minutes`](./transcribing-media-to-minutes/) | Reads a Thai/English bilingual meeting recording (audio or video, mostly Thai) and produces faithful, plain-language Minutes of Meeting plus a timestamped transcript sidecar. Routes by length/complexity to fast vs frontier model variants (Gemini examples; portable across hosts). Use when asked to "minute this recording" / "ทำรายงานการประชุม". | skillify | references · templates · examples |

## Alphabetical Index

| Skill | Purpose Group | Source family |
|-------|---------------|---------------|
| [`agent-context-initializer`](./agent-context-initializer/) | Agent Orchestration & Workflow Infrastructure | existing treasury |
| [`agentic-workflow-design`](./agentic-workflow-design/) | Agent Orchestration & Workflow Infrastructure | business-analyse |
| [`analyzing-banking-requirements`](./analyzing-banking-requirements/) | Banking, BA Delivery & Requirements | business-analyse |
| [`api-contract-design`](./api-contract-design/) | Architecture, Engineering Decisions & Planning | business-analyse |
| [`architecting-fintech-systems`](./architecting-fintech-systems/) | Architecture, Engineering Decisions & Planning | squad-delivery |
| [`architecture-decision`](./architecture-decision/) | Architecture, Engineering Decisions & Planning | business-analyse |
| [`assembling-tl-handoff`](./assembling-tl-handoff/) | Banking, BA Delivery & Requirements | business-analyse |
| [`assessing-ba-feasibility`](./assessing-ba-feasibility/) | Banking, BA Delivery & Requirements | business-analyse |
| [`augment-diagrams`](./augment-diagrams/) | Research, Debate & Knowledge Synthesis | squad-researcher |
| [`authoring-scaffold-profile`](./authoring-scaffold-profile/) | Agent Orchestration & Workflow Infrastructure | existing treasury |
| [`authoring-workflow-postmortem`](./authoring-workflow-postmortem/) | Observability, Cost & Incident Operations | existing treasury |
| [`business-logic-extractor`](./business-logic-extractor/) | Code Analysis, Productivity & Publishing | existing treasury |
| [`checking-ba-governance`](./checking-ba-governance/) | Banking, BA Delivery & Requirements | business-analyse |
| [`composing-agent-pipelines`](./composing-agent-pipelines/) | Agent Orchestration & Workflow Infrastructure | existing treasury |
| [`configuring-sandbox-allowlist`](./configuring-sandbox-allowlist/) | Security, Governance & Compliance | existing treasury |
| [`crafting-backend-code`](./crafting-backend-code/) | Implementation, Platform Templates & Code Review | existing treasury |
| [`crafting-frontend-code`](./crafting-frontend-code/) | Implementation, Platform Templates & Code Review | existing treasury |
| [`crafting-rust-code`](./crafting-rust-code/) | Implementation, Platform Templates & Code Review | existing treasury |
| [`cross-examine`](./cross-examine/) | Research, Debate & Knowledge Synthesis | squad-brainstorm |
| [`daily-planner`](./daily-planner/) | Code Analysis, Productivity & Publishing | skillify |
| [`data-modeling`](./data-modeling/) | Architecture, Engineering Decisions & Planning | business-analyse |
| [`defining-engineering-standards`](./defining-engineering-standards/) | Architecture, Engineering Decisions & Planning | business-analyse |
| [`delivery-planning`](./delivery-planning/) | Architecture, Engineering Decisions & Planning | business-analyse |
| [`designing-tech-lead-handoff`](./designing-tech-lead-handoff/) | Architecture, Engineering Decisions & Planning | squad-delivery |
| [`developing-langgraph-workflows`](./developing-langgraph-workflows/) | Agent Orchestration & Workflow Infrastructure | langgraph-claude-agent |
| [`devops-infrastructure-hardener`](./devops-infrastructure-hardener/) | Security, Governance & Compliance | existing treasury |
| [`domain-modeling`](./domain-modeling/) | Architecture, Engineering Decisions & Planning | business-analyse |
| [`drafting-ba-stories`](./drafting-ba-stories/) | Banking, BA Delivery & Requirements | business-analyse |
| [`drafting-stage-prompt`](./drafting-stage-prompt/) | Agent Orchestration & Workflow Infrastructure | existing treasury |
| [`eliciting-banking-brief`](./eliciting-banking-brief/) | Banking, BA Delivery & Requirements | business-analyse |
| [`engineer-growth-planning`](./engineer-growth-planning/) | Architecture, Engineering Decisions & Planning | business-analyse |
| [`engineering-doc-planning`](./engineering-doc-planning/) | Architecture, Engineering Decisions & Planning | business-analyse |
| [`evaluate-banking-compliance`](./evaluate-banking-compliance/) | Banking, BA Delivery & Requirements | craft |
| [`extract-brief-structure`](./extract-brief-structure/) | Banking, BA Delivery & Requirements | craft |
| [`extract-findings`](./extract-findings/) | Research, Debate & Knowledge Synthesis | squad-researcher |
| [`frame-debate`](./frame-debate/) | Research, Debate & Knowledge Synthesis | squad-brainstorm |
| [`generate-ux-pack`](./generate-ux-pack/) | Architecture, Engineering Decisions & Planning | squad-delivery |
| [`generating-gherkin-acceptance-criteria`](./generating-gherkin-acceptance-criteria/) | Testing, QA & Validation | craft |
| [`generating-pseudocode`](./generating-pseudocode/) | Code Analysis, Productivity & Publishing | existing treasury |
| [`governance-policy-generator`](./governance-policy-generator/) | Security, Governance & Compliance | existing treasury |
| [`implement-backend-feature`](./implement-backend-feature/) | Implementation, Platform Templates & Code Review | squad-delivery |
| [`implement-frontend-feature`](./implement-frontend-feature/) | Implementation, Platform Templates & Code Review | squad-delivery |
| [`implementing-go-template-requirements`](./implementing-go-template-requirements/) | Implementation, Platform Templates & Code Review | agentic |
| [`incident-response`](./incident-response/) | Observability, Cost & Incident Operations | business-analyse |
| [`integration-design`](./integration-design/) | Architecture, Engineering Decisions & Planning | business-analyse |
| [`model-selection`](./model-selection/) | Architecture, Engineering Decisions & Planning | business-analyse |
| [`multi-agent-handoff-architect`](./multi-agent-handoff-architect/) | Agent Orchestration & Workflow Infrastructure | existing treasury |
| [`observability-design`](./observability-design/) | Observability, Cost & Incident Operations | business-analyse |
| [`observability-telemetry-instrumenter`](./observability-telemetry-instrumenter/) | Observability, Cost & Incident Operations | existing treasury |
| [`opening-debate-panel`](./opening-debate-panel/) | Research, Debate & Knowledge Synthesis | squad-brainstorm |
| [`orchestrate-banking-brief-pipeline`](./orchestrate-banking-brief-pipeline/) | Banking, BA Delivery & Requirements | craft |
| [`orchestrating-agent-scaffold`](./orchestrating-agent-scaffold/) | Agent Orchestration & Workflow Infrastructure | existing treasury |
| [`orchestrating-openclaw-squad`](./orchestrating-openclaw-squad/) | Agent Orchestration & Workflow Infrastructure | existing treasury |
| [`organizing-local-files`](./organizing-local-files/) | Code Analysis, Productivity & Publishing | existing treasury |
| [`performance-cost-review`](./performance-cost-review/) | Observability, Cost & Incident Operations | business-analyse |
| [`plan-research`](./plan-research/) | Research, Debate & Knowledge Synthesis | squad-researcher |
| [`planning-banking-tests`](./planning-banking-tests/) | Testing, QA & Validation | squad-delivery |
| [`platform-common`](./platform-common/) | Implementation, Platform Templates & Code Review | agentic |
| [`platform-go-service`](./platform-go-service/) | Implementation, Platform Templates & Code Review | agentic |
| [`pr-design-review`](./pr-design-review/) | Architecture, Engineering Decisions & Planning | business-analyse |
| [`production-readiness`](./production-readiness/) | Architecture, Engineering Decisions & Planning | business-analyse |
| [`progressive-bug-hunter`](./progressive-bug-hunter/) | Code Analysis, Productivity & Publishing | existing treasury |
| [`publishing-git-review-requests`](./publishing-git-review-requests/) | Code Analysis, Productivity & Publishing | existing treasury |
| [`refactor-decision`](./refactor-decision/) | Architecture, Engineering Decisions & Planning | business-analyse |
| [`refactoring-go-services`](./refactoring-go-services/) | Implementation, Platform Templates & Code Review | agentic |
| [`rendering-readable-html`](./rendering-readable-html/) | Code Analysis, Productivity & Publishing | existing treasury |
| [`report-debate`](./report-debate/) | Research, Debate & Knowledge Synthesis | squad-brainstorm |
| [`reporting-research-run`](./reporting-research-run/) | Research, Debate & Knowledge Synthesis | squad-researcher |
| [`research-vault-librarian`](./research-vault-librarian/) | Research, Debate & Knowledge Synthesis | existing treasury |
| [`review-backend-code`](./review-backend-code/) | Implementation, Platform Templates & Code Review | agentic |
| [`review-frontend-code`](./review-frontend-code/) | Implementation, Platform Templates & Code Review | squad-delivery |
| [`review-report`](./review-report/) | Research, Debate & Knowledge Synthesis | squad-researcher |
| [`review-rust-code`](./review-rust-code/) | Implementation, Platform Templates & Code Review | existing treasury |
| [`reviewing-agent-spend`](./reviewing-agent-spend/) | Observability, Cost & Incident Operations | existing treasury |
| [`reviewing-implement-gate`](./reviewing-implement-gate/) | Agent Orchestration & Workflow Infrastructure | existing treasury |
| [`reviewing-software-security`](./reviewing-software-security/) | Security, Governance & Compliance | squad-delivery |
| [`revise-positions`](./revise-positions/) | Research, Debate & Knowledge Synthesis | squad-brainstorm |
| [`risk-estimation`](./risk-estimation/) | Architecture, Engineering Decisions & Planning | business-analyse |
| [`running-ba-pipeline`](./running-ba-pipeline/) | Banking, BA Delivery & Requirements | business-analyse |
| [`running-business-analysis-workflow`](./running-business-analysis-workflow/) | Banking, BA Delivery & Requirements | craft |
| [`scoping-ba-intake`](./scoping-ba-intake/) | Banking, BA Delivery & Requirements | business-analyse |
| [`scoping-technical-requirements`](./scoping-technical-requirements/) | Banking, BA Delivery & Requirements | business-analyse |
| [`search-sources`](./search-sources/) | Research, Debate & Knowledge Synthesis | squad-researcher |
| [`sweep-ambiguities`](./sweep-ambiguities/) | Banking, BA Delivery & Requirements | craft |
| [`synthesize-consensus`](./synthesize-consensus/) | Research, Debate & Knowledge Synthesis | squad-brainstorm |
| [`synthesize-report`](./synthesize-report/) | Research, Debate & Knowledge Synthesis | squad-researcher |
| [`technical-debt-management`](./technical-debt-management/) | Architecture, Engineering Decisions & Planning | business-analyse |
| [`technical-feasibility`](./technical-feasibility/) | Architecture, Engineering Decisions & Planning | business-analyse |
| [`testing-strategy`](./testing-strategy/) | Testing, QA & Validation | business-analyse |
| [`thai-translator`](./thai-translator/) | Code Analysis, Productivity & Publishing | business-analyse |
| [`transcribing-media-to-minutes`](./transcribing-media-to-minutes/) | Code Analysis, Productivity & Publishing | skillify |
| [`universal-spec-validator`](./universal-spec-validator/) | Security, Governance & Compliance | existing treasury |
| [`validating-banking-implementation`](./validating-banking-implementation/) | Testing, QA & Validation | squad-delivery |

**Import policy** — duplicate source skills are resolved by destination name. Active skill folders outrank generated runs, archive/demo folders, and test fixtures; ties inside a class use newest local `SKILL.md`, then source path. Imported skill content is preserved as-is.

**Validation note** — some imported skills use richer YAML frontmatter than `skillify/scripts/quick_validate.py` accepts. The catalog therefore treats `SKILL.md` presence and README link integrity as the blocking checks for this import.
