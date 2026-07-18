# Treasury — Production Skill Library

131 top-level skills, grouped by purpose. Each skill name links to its folder. The **Extra assets** column lists files or folders beside `SKILL.md`.

| Purpose Group | Count |
|---------------|-------|
| [Banking, BA Delivery & Requirements](#banking-ba-delivery-requirements) | 15 |
| [Architecture, Engineering Decisions & Planning](#architecture-engineering-decisions-planning) | 22 |
| [Implementation, Platform Templates & Code Review](#implementation-platform-templates-code-review) | 15 |
| [Testing, QA & Validation](#testing-qa-validation) | 12 |
| [Agent Orchestration & Workflow Infrastructure](#agent-orchestration-workflow-infrastructure) | 17 |
| [Research, Debate & Knowledge Synthesis](#research-debate-knowledge-synthesis) | 15 |
| [Security, Governance & Compliance](#security-governance-compliance) | 7 |
| [Observability, Cost & Incident Operations](#observability-cost-incident-operations) | 9 |
| [Code Analysis, Productivity & Publishing](#code-analysis-productivity-publishing) | 19 |
| **Total** | **131** |

---

## Banking, BA Delivery & Requirements (15)

Business-analysis workflows for regulated banking and delivery handoffs.

| Skill | Purpose | Source family | Extra assets |
|-------|---------|---------------|--------------|
| [`analyzing-banking-requirements`](./analyzing-banking-requirements/) | Business Analyst persona for enterprise banking and lending workflows. Produces a compliance-validated PRD with API contracts from a raw business request — extracting the ask, running strict KYC / AML / PCI-DSS compli... | business-analyse | references |
| [`assembling-tl-handoff`](./assembling-tl-handoff/) | Stage 5 of the BA pipeline: assemble the Handoff Bundle the Tech Lead signs for — the raw requirement plus the Scope Sheet, Story Set, clear Governance Check, and buildable-or-phased Feasibility Note linked by path, p... | business-analyse | — |
| [`assessing-ba-feasibility`](./assessing-ba-feasibility/) | Stage 4 of the BA pipeline: consume the Story Set, a clear Governance Check, and the raw requirement, and emit a typed Feasibility Note contract (at least two options with pros and cons, dependencies, risks with mitig... | business-analyse | — |
| [`checking-ba-governance`](./checking-ba-governance/) | Stage 3 of the BA pipeline: sweep a Story Set for governance and privacy gaps and emit a typed Governance Check contract (PII flags, compliance flags, and blockers each with a named resolve owner, plus a verdict of cl... | business-analyse | — |
| [`drafting-ba-stories`](./drafting-ba-stories/) | Stage 2 of the BA pipeline: consume a confirmed Scope Sheet and emit a typed Story Set — one epic plus INVEST user stories, each a full user-story-template instance (Description, Business Logic with decision tables, G... | business-analyse | references · templates |
| [`eliciting-banking-brief`](./eliciting-banking-brief/) | Convert raw BA input (Jira, Slack, meeting notes, email, mixed prose) into a structured epic-plus-stories brief with banking-grade fields force-evaluated, ambiguities surfaced as open questions, and governance gaps (L... | business-analyse | audit · references · schemas · scripts · tests |
| [`evaluate-banking-compliance`](./evaluate-banking-compliance/) | Scan structured banking Epics and Stories for PII, regulatory dependencies, tipping-off risk, missing Legal or Privacy roles, and TL-handoff governance blockers. Use when a user asks to "check this brief for banking c... | craft | examples · references · templates |
| [`extract-brief-structure`](./extract-brief-structure/) | Convert redacted raw BA input from Jira, Slack, meeting notes, emails, or mixed prose into a strict Epic and Story JSON skeleton. Use when a user asks to "parse this raw Jira ticket into a brief", "extract the epics a... | craft | examples · references · templates |
| [`orchestrate-banking-brief-pipeline`](./orchestrate-banking-brief-pipeline/) | Run the decomposed banking BA brief pipeline across extraction, compliance, ambiguity, and Gherkin stages with validated JSON handoffs and pluggable model commands. Use when a user asks to "run the decomposed banking... | craft | examples · references · scripts · templates |
| [`researching-ba-problem-space`](./researching-ba-problem-space/) | Run upstream BA problem-space discovery BEFORE intake — investigate the problem, frame opportunities, surface assumptions and the four product risks (value, usability, feasibility, viability), and map the banking regu... | workflow-pack | references · schemas |
| [`running-ba-pipeline`](./running-ba-pipeline/) | Run the five-stage agentic Business-Analysis pipeline that turns a raw requirement into a Tech-Lead-ready handoff: it carries one task_id and a typed envelope across intake, story drafting, governance, feasibility, an... | business-analyse | references |
| [`running-business-analysis-workflow`](./running-business-analysis-workflow/) | Guide an AI agent through a structured business-analysis workflow for software and IT projects: define the problem and goals, analyze the current (as-is) state, elicit functional and non-functional requirements, desig... | craft | references · templates |
| [`scoping-ba-intake`](./scoping-ba-intake/) | Stage 1 of the BA pipeline: ingest a raw requirement (Jira, email, meeting notes, document) and emit a typed Scope Sheet contract wrapped in the pipeline's shared envelope, ready for the BA/PM scope-confirm gate (G1)... | business-analyse | — |
| [`scoping-technical-requirements`](./scoping-technical-requirements/) | Turns a messy business requirement into a clear, bounded technical scope a team can safely act on, surfacing the open questions and non-functional requirements that decide everything downstream. Use when the user asks... | business-analyse | — |
| [`sweep-ambiguities`](./sweep-ambiguities/) | Detect linguistic ambiguities and hidden requirements in a structured banking brief using the eight ambiguity detectors and ten elicitation frames. Use when a user asks to "find the hidden requirements in this brief",... | craft | examples · references · templates |

## Architecture, Engineering Decisions & Planning (22)

Design, modeling, planning, and engineering decision support.

| Skill | Purpose | Source family | Extra assets |
|-------|---------|---------------|--------------|
| [`api-contract-design`](./api-contract-design/) | Design a stable, clear API contract (action, request/response, error cases, validation, idempotency, versioning) for an API you expose, so consumers can rely on it and it stays backward-compatible. Use when the user a... | business-analyse | — |
| [`architecting-fintech-systems`](./architecting-fintech-systems/) | Senior fintech domain architect for lending, loan origination, KYC, credit decisioning, disbursement, and regulated financial systems. Owns L1–L3 design (business intent → bounded contexts → technical strategy) using... | squad-delivery | references · templates |
| [`architecture-decision`](./architecture-decision/) | Chooses a design among options and records it as an ADR with explicit trade-offs, protecting long-term maintainability and setting clear direction the team can follow. Use when the user asks "should we split this into... | business-analyse | — |
| [`befe-contract-design`](./befe-contract-design/) | Design the two-sided backend/frontend shared contract for parallel development — a design-first OpenAPI/AsyncAPI single source of truth, generated client types, a consumer mock/stub so the frontend is unblocked before... | workflow-pack | references · schemas |
| [`data-modeling`](./data-modeling/) | Design the persistence model safely (schema, source of truth, indexes, constraints) with an expand/contract migration and rollback plan that preserves backward compatibility. Use when the user asks "design the table/s... | business-analyse | — |
| [`defining-engineering-standards`](./defining-engineering-standards/) | Define and own the team's reusable standard — tech-stack baseline, DDD domain boundaries, a C4 L3/L4 design skeleton, API naming and folder/layer conventions, and the project scaffold — so the team builds consistently... | business-analyse | — |
| [`delivery-planning`](./delivery-planning/) | Break a chosen solution into a sequenced, estimated execution plan — tasks, critical path, blockers, and a phased timeline honest enough to tell the business. Use when the user asks to "break this into tasks", "estima... | business-analyse | — |
| [`designing-tech-lead-handoff`](./designing-tech-lead-handoff/) | Convert an approved BA epic-and-stories brief plus a UX design pack into the full Tech-Lead architecture handoff: integration contracts, component map, infra spec, ADRs, per-service API specs, per-story L4 specs, and... | squad-delivery | audit · references · schemas · templates · tests |
| [`domain-modeling`](./domain-modeling/) | Defines domain boundaries and DDD building blocks so ownership is clear, the model matches the business, and domain logic stays separate from infrastructure. Use when the user asks "define the domain model", "where ar... | business-analyse | — |
| [`engineer-growth-planning`](./engineer-growth-planning/) | Grows an engineer — runs useful 1-on-1s, chooses teach-by-doing vs teach-by-telling, spots underperformance early, delegates to build capability, and reviews thinking (not just output) without becoming the bottleneck.... | business-analyse | — |
| [`engineering-doc-planning`](./engineering-doc-planning/) | Decides what to document and what to deliberately skip, for which audience, keeping docs minimal, findable, and actionable across ADRs, C4 design docs, API contracts, system context, ubiquitous language, runbooks, and... | business-analyse | — |
| [`generate-ux-pack`](./generate-ux-pack/) | Produce a UX-design intake pack (vendored v1.1 contract) from a UX team's drop (bundled prototype HTML, Frontend Spec markdown, BA brief directory). Emits a structured `ux-design-{idem8}/` tree with tokens.json (W3C d... | squad-delivery | references · schemas |
| [`integration-design`](./integration-design/) | Design resilient integration with an external/3rd-party system (timeouts, retries, fallback, error mapping, idempotency, and a clear ownership/support model) so their failures do not become our outages. Use when the u... | business-analyse | — |
| [`model-selection`](./model-selection/) | Assign each agent or stage in an agentic workflow a model and reasoning effort by criteria — a privacy/data-class gate, capability-to-role match, cost-per-successful-task, reasoning effort, structured decoding, and a... | business-analyse | — |
| [`pr-design-review`](./pr-design-review/) | Review a PR's design and maintainability — business-logic completeness, test coverage, over-engineering, and template conformance — returning tagged, teachable comments plus a merge-or-iterate verdict. Use when the us... | business-analyse | — |
| [`principal-advisor`](./principal-advisor/) | Principal/staff-level engineering sparring partner for explicit technology decisions and idea exploration — compares stack and architecture options, tests feasibility, steelmans then challenges assumpt... | skillify | references · templates |
| [`production-readiness`](./production-readiness/) | Verify a feature is safe to ship — observability, rollback, runbook, migration safety, and a passing smoke test — and return a clear go / no-go / conditional verdict. Use when the user asks "is this safe to ship", "pr... | business-analyse | — |
| [`red-teaming-implementation-plan`](./red-teaming-implementation-plan/) | Adversarially red-team an implementation plan or Tech-Lead design BEFORE any code is written, then issue a machine-readable PROCEED, REVISE, or BLOCK verdict with severity-ranked findings. Use when asked to red-team a... | workflow-pack | references · schemas |
| [`refactor-decision`](./refactor-decision/) | Decide whether a refactor is necessary, bound its scope, protect existing behavior with tests first, and avoid cosmetic refactors that have no business reason. Use when the user asks "should we refactor this", "how bi... | business-analyse | — |
| [`risk-estimation`](./risk-estimation/) | Size complexity and separate known work from unknown risk, producing a man-day estimate with an explicit confidence and uncertainty multiplier plus the assumptions behind it, to feed delivery planning rather than repl... | business-analyse | — |
| [`technical-debt-management`](./technical-debt-management/) | Classify technical debt, prioritize it by business risk and cost-of-delay, make the dangerous debt visible, and negotiate a fix budget with PM/business. Use when the user asks "how do we deal with this tech debt", "is... | business-analyse | — |
| [`technical-feasibility`](./technical-feasibility/) | Decides whether and how a requirement is buildable on the current stack, surfacing the options, dependencies, and risks before anyone commits to a design or a date. Use when the user asks "can we build this", "is this... | business-analyse | — |

## Implementation, Platform Templates & Code Review (15)

Implementation, platform template, refactoring, and code-review skills.

| Skill | Purpose | Source family | Extra assets |
|-------|---------|---------------|--------------|
| [`clean-go-service`](./clean-go-service/) | Incrementally refactor messy Go microservices toward clean DDD and CQRS architecture while preserving behavior. Use when asked to clean up, refactor, restructure, or improve Go service code without changing external b... | rust-research | — |
| [`crafting-backend-code`](./crafting-backend-code/) | Reviews, designs, and safely implements backend code and microservice templates with a pattern-first, evidence-led posture. Use when designing, reviewing, optimizing, fixing, analyzing, or planning backend services, m... | existing treasury | references |
| [`crafting-frontend-code`](./crafting-frontend-code/) | Reviews, designs, and safely implements frontend code with a conservative, repo-first posture for AI coding agents. Use when designing, reviewing, optimizing, fixing, analyzing, or planning React/TypeScript frontend f... | existing treasury | references |
| [`crafting-rust-code`](./crafting-rust-code/) | Reviews, designs, and safely implements Rust code with a pattern-first, evidence-led, repo-first posture for AI coding agents. Use when designing, reviewing, optimizing, fixing, analyzing, or planning Rust services, c... | existing treasury | references |
| [`golite`](./golite/) | Implement, improve, and refactor Go backend code as a senior engineer: simple, high-impact, readable code with lean comments, idiomatic Go conventions, Martin Fowler refactoring idioms, a handler/service/access layering kernel, and table-driven testify... | platform-mgmt | references |
| [`implement-backend-feature`](./implement-backend-feature/) | Generate production-grade Go backend code for one microservice feature from an approved design document. Use when implementing a Go HTTP handler from a design spec. Use when generating a CQRS command or query handler... | squad-delivery | RATIONALE.md · references · schemas · tests |
| [`implement-frontend-feature`](./implement-frontend-feature/) | Generate production-grade React/TypeScript code for one frontend feature from an approved UI design, with banking-grade discipline: WCAG 2.1 AA a11y, no any outside parsers, no localStorage auth tokens, no unsanitized... | squad-delivery | RATIONALE.md · references · schemas · tests |
| [`implementing-go-template-requirements`](./implementing-go-template-requirements/) | Apply a single requirement (spec line, ticket, bug report, user story) to a Go service that follows the go-template scaffold by editing ONLY business logic under `app/[domain]/` plus narrow `register*` wiring in `rout... | agentic | examples · references · templates |
| [`langgraph-professional`](./langgraph-professional/) | Guide professional LangGraph v1.x implementation, refactoring, and review workflows. Use when the user asks "Implement a LangGraph v1.x workflow in this repo using professional StateGraph patterns", "Refactor this Lan... | langgraph-claude-agent | references · templates |
| [`platform-common`](./platform-common/) | Shared Go infrastructure library for A-Team Krungthai DGL microservices — Gin middleware, Kafka producer/consumer, JWT, structured slog logging, response envelopes, database/Redis/Firestore/GCS/S3 connectors, AES/RSA... | agentic | references |
| [`platform-go-service`](./platform-go-service/) | Scaffold or extend a Go microservice in this repository (go-template) using A-Team platform conventions: DDD aggregate-per-package, CQRS handler/consumer split, Fowler-style Repository / Cache / Gateway in `access/`,... | agentic | references |
| [`refactoring-go-services`](./refactoring-go-services/) | Incrementally refactor messy Go microservices toward clean DDD/CQRS architecture while preserving behavior — one code smell and one Fowler-style refactoring action per iteration, verified by tests or build after each... | agentic | examples · references |
| [`review-backend-code`](./review-backend-code/) | Adversarially verify Go backend code emitted by a Generate stage against the approved design and the 11 banking-grade decision rules + v2 augmentations, then issue a machine-readable verdict (approve, loop_back, or hu... | agentic | RATIONALE.md · references · schemas · tests |
| [`review-frontend-code`](./review-frontend-code/) | Adversarially verify React/TypeScript code emitted by a Generate stage against the approved UI design and the 12 banking-grade frontend non-negotiables + 9 v2 augmentations, then issue a machine-readable verdict (appr... | squad-delivery | RATIONALE.md · references · schemas · tests |
| [`review-rust-code`](./review-rust-code/) | Adversarially verify Rust backend code emitted by a Generate stage against the approved design and the 11 banking-grade decision rules + v2 augmentations, then issue a machine-readable verdict (approve, loop_back, or... | existing treasury | RATIONALE.md · references · schemas · tests |

## Testing, QA & Validation (12)

Test planning, acceptance criteria, QA, and implementation validation.

| Skill | Purpose | Source family | Extra assets |
|-------|---------|---------------|--------------|
| [`authoring-e2e-test-suite`](./authoring-e2e-test-suite/) | Derive end-to-end scenarios from user journeys or stories AND execute them against a running system in a sandbox via the CI/test runner, surface flaky scenarios by bounded re-run, and emit a runner-backed PASS, FAIL,... | workflow-pack | references · schemas |
| [`contract-testing-pact`](./contract-testing-pact/) | Run consumer-driven contract tests with Pact in a sandbox via the CI/test runner, verify the provider against each consumer contract, query the Pact Broker can-i-deploy gate, and emit a runner-backed PASS, FAIL, or ER... | workflow-pack | references · schemas |
| [`executing-backend-unit-tests`](./executing-backend-unit-tests/) | Execute backend unit tests against built backend artifacts in a sandbox via the CI/test runner, measure real line coverage, surface flaky tests, catch unit-level defects, and emit a runner-backed PASS, FAIL, or ERROR... | workflow-pack | references · schemas |
| [`executing-frontend-unit-tests`](./executing-frontend-unit-tests/) | Execute frontend unit tests against built frontend artifacts in a sandbox via the CI/test runner, measure real per-file coverage, surface flaky tests, catch unit-level defects, and emit a runner-backed PASS, FAIL, or... | workflow-pack | references · schemas |
| [`executing-integration-tests`](./executing-integration-tests/) | Execute integration tests (post-implementation SIT) against a wired-up system in a sandbox via the CI/test runner, verify clean teardown, catch integration defects, and emit a runner-backed PASS, FAIL, or ERROR gate t... | workflow-pack | references · schemas |
| [`executing-qa-test-suite`](./executing-qa-test-suite/) | Execute a planned QA test roster against a running system, measure run-system coverage, surface flaky tests, catch live defects, and emit a results-backed PASS, FAIL, or ERROR QA gate that feeds a human verification l... | workflow-pack | references · schemas |
| [`generating-gherkin-acceptance-criteria`](./generating-gherkin-acceptance-criteria/) | Translate resolved banking stories into strict BDD Gherkin acceptance criteria with happy, error, audit, idempotency, and tipping-off-safe scenarios. Use when a user asks to "write Gherkin ACs for these banking storie... | craft | examples · references · templates |
| [`planning-banking-tests`](./planning-banking-tests/) | Convert a BA brief (output of eliciting-banking-brief v1.2+) into a structured QA test plan covering every Gherkin scenario, banking-grade concern, NFR target, regulatory dependency, and compliance requirement. Emits... | squad-delivery | references · schemas · scripts · tests |
| [`running-accessibility-tests`](./running-accessibility-tests/) | Run automated accessibility tests against a built front end during FE review or SIT and emit a PASS, FAIL, or ERROR gate for WCAG 2.1 AA — surfacing automated violations AND the needs-review items the engine cannot de... | workflow-pack | references · schemas |
| [`running-smoke-tests`](./running-smoke-tests/) | Run a small set of critical-path smoke probes against a freshly deployed live release and emit a PASS, FAIL, or ERROR sanity gate from real probe results — each probe green or red with its measured latency, never an a... | workflow-pack | references · schemas |
| [`testing-strategy`](./testing-strategy/) | Decide what to test and at which level (unit / integration / contract / e2e / regression) and map a requirement's business logic into concrete test scenarios, so coverage guarantees behavior rather than a number. Use... | business-analyse | — |
| [`validating-banking-implementation`](./validating-banking-implementation/) | QA Engineer persona for enterprise banking implementation artifacts. Adversarial OWASP Top 10 testing, transactional-integrity audit, race / deadlock analysis, and chaos-test plan drafting on completed Developer-stage... | squad-delivery | references |

## Agent Orchestration & Workflow Infrastructure (17)

Multi-agent orchestration, scaffold configuration, prompts, gates, and workflow infrastructure.

| Skill | Purpose | Source family | Extra assets |
|-------|---------|---------------|--------------|
| [`agent-context-initializer`](./agent-context-initializer/) | Generate a minimal, human-curatable AGENTS.md (under 200 lines, target 100-150) for an agentic-squad repo: exactly six sections (commands, testing, project structure, code style, git workflow, boundaries), every prohi... | existing treasury | references · scripts · templates |
| [`agentic-workflow-design`](./agentic-workflow-design/) | Design how an AI-agent pipeline is supervised — its stages and owners, the human-in-the-loop approval gates, the never-do guardrails and command-safety policy, accountability for AI-generated output, and observability... | business-analyse | references |
| [`authoring-scaffold-profile`](./authoring-scaffold-profile/) | Author or modify an agent-scaffold profile file at profiles/NAME.sh. Sets STAGES, GATED_STAGES, and per-stage AGENT / MODEL / PROMPT_PREFIX env vars. Validates that every STAGE has a runner script, that GATED_STAGES i... | existing treasury | references · templates |
| [`composing-agent-pipelines`](./composing-agent-pipelines/) | Composes a portable multi-agent pipeline (Plan, Gather, Analyze, Review, Validate, Decide, Compact) for research, code review, implementation planning, or trade-off analysis. Use when the user asks to "compose agent p... | existing treasury | examples · references · scripts · templates |
| [`delegating-to-cli-models`](./delegating-to-cli-models/) | Drive external model CLIs as advisory sub-agents from one orchestrating agent: codex (GPT, gpt-5.5) and agy (Antigravity/Gemini), plus optional headless self-delegation. Dispatch headlessly with a watchdog and the exact model label, ground-truth the backend, then adjudic... | skillify | examples · references · scripts · templates |
| [`developing-langgraph-workflows`](./developing-langgraph-workflows/) | Guide professional LangGraph v1.x implementation, refactoring, and review workflows. Use when the user asks "Implement a LangGraph v1.x workflow in this repo using professional StateGraph patterns", "Refactor this Lan... | langgraph-claude-agent | references · templates |
| [`discord-harness-extend`](./discord-harness-extend/) | Add or modify one tool in a Discord LLM harness while preserving the gate-lives-outside-the-model invariant: write the @tool handler, re-validate side effects inside the handler, classify ALLOW/CONFIRM/DENY in config... | discord-llm-harness | references · scripts · templates |
| [`discord-harness-operate`](./discord-harness-operate/) | Run, configure, monitor, and deploy the Discord MySQL-operator bot without changing its tool code: the .env-secrets vs config.yaml-policy split, model swap via models.work, the Grafana/Prometheus stack, /healthz, and... | discord-llm-harness | references · scripts |
| [`drafting-stage-prompt`](./drafting-stage-prompt/) | Curate stage prompts for the agent-scaffold prompt library at prompts/library/STAGE/TOPIC.md. Captures the Friday prompt-review ritual: a stage prompt with a why-it-works paragraph, success metric, failure mode, model... | existing treasury | references · templates |
| [`handoff-revoke`](./handoff-revoke/) | Reverse a previously emitted deploy handoff as a SAGA compensating action — revoke the issued short-lived deploy credentials, signal the release control plane to halt or roll back the promotion tied to a handoff recei... | workflow-pack | references · schemas |
| [`handoff-to-deploy`](./handoff-to-deploy/) | Hand a QA-signed-off release off to the deploy/release control plane behind a mandatory synchronous named-human approval, mint short-lived OIDC deploy credentials, and emit an immutable handoff receipt that ties the l... | workflow-pack | references · schemas |
| [`launching-workflow-routines`](./launching-workflow-routines/) | Trigger a pre-designed multi-node workflow routine BY NAME from a routines/ registry: resolve, script-validate, show an overview run plan, then execute node-by-node with human gates, script-verified artifacts under... | skillify | examples · references · scripts · templates |
| [`multi-agent-handoff-architect`](./multi-agent-handoff-architect/) | Design inter-agent handoff APIs as versioned JSON Schema contracts (required taskId, intent, state, confidence, provenance/trace, schemaVersion) with binary single-writer ownership, observability spanning agent bounda... | existing treasury | references · schemas · scripts · templates |
| [`orchestrating-agent-scaffold`](./orchestrating-agent-scaffold/) | Orchestrates a multi-stage research-squad workflow on top of the agent-scaffold (just / state.json / zellij / LiteLLM / ntfy / sandbox). Plans the run, picks profile and cap, spawns fit-for-job sub-agents for pre-flig... | existing treasury | references · scripts · templates |
| [`orchestrating-openclaw-squad`](./orchestrating-openclaw-squad/) | Orchestrates a multi-agent squad (Business Analyst, Architect, Developer, QA Engineer) with strict human-in-the-loop validation, sandbox isolation, and explicit approval gates for enterprise and regulated workflows. U... | existing treasury | references · templates |
| [`reviewing-implement-gate`](./reviewing-implement-gate/) | Walk a researcher through the five-question check before approving the implement gate in an agent-scaffold workflow. Reads .agent/stages/plan.md and .agent/stages/critique.md, confirms IMPLEMENT_SANDBOXED for credenti... | existing treasury | references · templates |
| [`roundtripping-dashboard-data-contract`](./roundtripping-dashboard-data-contract/) | Edit the squad-delivery dashboard's data model through a validated JSON contract instead of hand-editing its gzip+base64 bundle. Use when the user asks to "change a dashboard stage, gate, owner, or tier", "update the... | workflow-pack | references · schemas · scripts |

## Research, Debate & Knowledge Synthesis (15)

Research pipeline, debate, source synthesis, and knowledge-vault workflows.

| Skill | Purpose | Source family | Extra assets |
|-------|---------|---------------|--------------|
| [`augment-diagrams`](./augment-diagrams/) | OPTIONAL out-of-band stage (NOT part of the standard 1–6 pipeline). Backfill 0–3 inline ASCII diagrams into an existing final research report, insertion-only, without rewriting any prose. Receives `final_report`, `fin... | squad-researcher | examples · references |
| [`cross-examine`](./cross-examine/) | Stage 3 of the squad-brainstorm workflow, run once per panel CLI per debate round (skipped entirely when rounds=quick). The invoked panelist reads the OTHER two panelists' latest positions and produces a structured, r... | squad-brainstorm | examples · references |
| [`extract-findings`](./extract-findings/) | Stage 3 of the researcher workflow. Given a research_plan (from plan-research) and sources (from search-sources), pull out structured, source-grounded claims organized against the plan's sub-questions. Every claim car... | squad-researcher | examples · references |
| [`frame-debate`](./frame-debate/) | Stage 1 of the squad-brainstorm workflow. Read a finished squad-researcher run (05-final_report.md + 03-findings.json + 01-research_plan.json) and distill it into a neutral debate brief — a contestable proposition, th... | squad-brainstorm | examples · references |
| [`opening-debate-panel`](./opening-debate-panel/) | Stage 2 of the squad-brainstorm workflow, run once per panel CLI (P-codex via `codex exec`, P-gemini via `gemini -p`, P-claude via `claude -p`). Given the debate brief and the shared grounding pack, the invoked paneli... | squad-brainstorm | examples · references |
| [`panel-open`](./panel-open/) | Stage 2 of the squad-brainstorm workflow, run once per panel CLI (P-codex via `codex exec`, P-gemini via `gemini -p`, P-claude via `claude -p`). Given the debate brief and the shared grounding pack, the invoked paneli... | squad-brainstorm | examples · references |
| [`plan-research`](./plan-research/) | Stage 1 of the squad-researcher workflow. Decompose a research topic into a structured plan — thesis question, 3–10 MECE sub-questions (count gated by depth), coverage dimensions, out-of-scope list, success criteria —... | squad-researcher | examples · references |
| [`report-debate`](./report-debate/) | Stage 6 (meta) of the squad-brainstorm workflow. Reads prior-stage artifacts from the debate directory plus an optional metrics.json and produces a single 00-debate_report.md summarising the run: panel roster with too... | squad-brainstorm | examples · references |
| [`reporting-research-run`](./reporting-research-run/) | Stage 6 (meta) of the squad-researcher workflow. Reads prior-stage artifacts from the run directory plus an optional metrics.json sidecar, and produces a single 00-run_report.md summarising the run: per-stage stats ta... | squad-researcher | examples · references |
| [`research-vault-librarian`](./research-vault-librarian/) | Read-only librarian and context-only workflow guide for a CLAUDE.md-governed Obsidian research-report vault (the ResearchVault: deep-research reports in domain folders, an index.md catalog, per-domain MOC maps, fed by... | existing treasury | references · scripts |
| [`review-report`](./review-report/) | Self-review pass over a draft research report. Verifies every cited claim is grounded in the input findings, checks audience and topic fit, surfaces internal-consistency issues, and may perform bounded citation surger... | squad-researcher | examples · references |
| [`revise-positions`](./revise-positions/) | Stage 4 of the squad-brainstorm workflow, run once per panel CLI per debate round (skipped when rounds=quick). The invoked panelist reads ONLY the critiques aimed at itself, then produces a revised position: concede w... | squad-brainstorm | examples · references |
| [`search-sources`](./search-sources/) | Stage 2 of the researcher workflow. Given a structured research_plan (output of plan-research) and a depth knob, gather candidate sources for every sub-question and emit a deduplicated, diversity-checked list of norma... | squad-researcher | examples · references |
| [`synthesize-consensus`](./synthesize-consensus/) | Stage 5 of the squad-brainstorm workflow. A single Claude moderator reads the attributed final panelist positions plus the cross-examination exchanges and produces the user-facing synthesized answer: settled agreement... | squad-brainstorm | examples · references |
| [`synthesize-report`](./synthesize-report/) | Compose a markdown research report from structured findings, tuned to a named audience, with inline numeric citations and a sources list. Stage 4 of 6 in the squad-researcher workflow: receives `topic`, `research_plan... | squad-researcher | examples · references |

## Security, Governance & Compliance (7)

Security review, policy, sandboxing, spec validation, and governance controls.

| Skill | Purpose | Source family | Extra assets |
|-------|---------|---------------|--------------|
| [`configuring-sandbox-allowlist`](./configuring-sandbox-allowlist/) | Edit the agent-scaffold sandbox hostname allowlist at docker/sandbox-proxy/filter. Adds, reviews, or removes hostname patterns, refusing wildcards, host-internal redirects, IP literals, and metadata-IP patterns. Remin... | existing treasury | references · templates |
| [`devops-infrastructure-hardener`](./devops-infrastructure-hardener/) | Audit an agentic CI/CD or agent-skill codebase for exposed static credentials and emit a remediation plan plus a runtime-secrets architecture that replaces them with short-lived task-scoped OIDC / dynamic secrets (Has... | existing treasury | references · scripts · templates |
| [`governance-policy-generator`](./governance-policy-generator/) | Emit policy-as-code agent-governance artifacts architecturally separated from the agent they govern: an Open Policy Agent (OPA) Rego allowlist with default-deny posture (default allow = false), plus a KILLSWITCH.md de... | existing treasury | references · scripts · templates |
| [`reviewing-software-security`](./reviewing-software-security/) | Defensive software security review for Go/Gin services, Kafka event flows, MySQL/RDS, Kubernetes workloads, Kong/APISIX gateways, and DDD/CQRS event-driven lending and financial systems across SIT/UAT/PRD. Use when th... | squad-delivery | references · templates |
| [`running-sast-security-gate`](./running-sast-security-gate/) | Run a static application security test (SAST) over developer-authored source in a sandbox, detecting insecure logic written in-house — injection, BOLA/IDOR authorization flaws, and hardcoded secrets — then emit a PASS... | workflow-pack | references · schemas |
| [`scanning-appsec-pipeline-gate`](./scanning-appsec-pipeline-gate/) | Scan the built and running application plus its third-party dependencies in a sandbox — DAST attacking the running app from outside, SCA flagging known CVEs in imported dependencies, and a secrets sweep — then emit a... | workflow-pack | references · schemas |
| [`universal-spec-validator`](./universal-spec-validator/) | Validate an agent spec (tool/function JSON Schema, MCP tool manifest, OpenAPI tool def, SKILL.md frontmatter, or data contract) as a CI / pre-commit ENFORCEMENT GATE that fails the build before malformed tool calls or... | existing treasury | references · scripts · templates |

## Observability, Cost & Incident Operations (9)

Telemetry, cost governance, incident handling, readiness, and postmortems.

| Skill | Purpose | Source family | Extra assets |
|-------|---------|---------------|--------------|
| [`analyzing-canary-rollout`](./analyzing-canary-rollout/) | Compare a canary against its baseline across multiple windows with a statistical non-inferiority test and emit a promote, hold, or rollback recommendation from real metric series — never an agent assertion, and never... | workflow-pack | references · schemas |
| [`authoring-workflow-postmortem`](./authoring-workflow-postmortem/) | Author a postmortem for an agent-scaffold workflow within 48h of a trigger event: failed stage with more than $1 spent, rejected implement gate, cap overrun, or sandbox egress test failure. Reads .agent/runlog.jsonl,... | existing treasury | references · templates |
| [`incident-response`](./incident-response/) | Stabilizes a firing production incident — classify severity, identify blast radius, drive a short-term mitigation, communicate, then run a blameless postmortem with prevention follow-ups. Use when the user asks "we ha... | business-analyse | — |
| [`observability-design`](./observability-design/) | Designs how a feature is observed in production — correlation/request/business IDs, logs, metrics, traces, dashboards, and alerts — so production behavior is debuggable and only meaningful problems page anyone. Use wh... | business-analyse | — |
| [`observability-telemetry-instrumenter`](./observability-telemetry-instrumenter/) | Instrument agent code with OpenTelemetry GenAI semantic conventions: an invoke_agent CLIENT root span, execute_tool INTERNAL child spans, the gen_ai.client.token.usage histogram with the exact 14 bucket boundaries, lo... | existing treasury | references · schemas · scripts · templates |
| [`performance-cost-review`](./performance-cost-review/) | Define and review a feature or service's performance budget and cost budget — latency/throughput/resource targets per critical user journey, cost-per-unit, baseline vs target, headroom, and the cheap-levers-first plan... | business-analyse | examples · references |
| [`reviewing-agent-spend`](./reviewing-agent-spend/) | Run the Monday cost-review ritual for an agent-scaffold squad. Wraps just llm-spend, segments by user and model, flags greater-than-2-sigma outliers and any single workflow over $10, and drafts a one-page team-lead re... | existing treasury | references · templates |
| [`running-performance-load-test`](./running-performance-load-test/) | Drive a pre-prod performance/load test against a staging or UAT target and emit a budget-backed PASS, FAIL, or ERROR gate from real runner metrics — p95, p99, error rate, and throughput measured by a load runner, neve... | workflow-pack | references · schemas |
| [`validating-production-slo`](./validating-production-slo/) | Validate a live production release against its declared SLOs by querying live SLIs over a bake window, evaluating multi-window burn-rate, and emitting a promote, hold, or rollback recommendation with a Pass, Marginal,... | workflow-pack | references · schemas |

## Code Analysis, Productivity & Publishing (19)

Code analysis, productivity, translation, publishing, and human-readable report generation.

| Skill | Purpose | Source family | Extra assets |
|-------|---------|---------------|--------------|
| [`business-logic-extractor`](./business-logic-extractor/) | Extract the implemented business logic and rules FROM a codebase — cross-referenced with requirements and agent/execution traces — into a faithful, traceable specification that BOUNDS information loss: salient rules,... | existing treasury | references · scripts · templates |
| [`daily-planner`](./daily-planner/) | Turn a raw list of tasks into a prioritized, trackable daily plan and keep it current across days. Use when the user says "here are my tasks, help me prioritize", "plan my day", "what should I work on first", "make my... | skillify | templates |
| [`drawio`](./drawio/) | Always use when user asks to create, generate, draw, or design a diagram, flowchart, architecture diagram, ER diagram, sequence diagram, class diagram, network diagram, mockup, wireframe, or UI sketch, or mentions dra... | drawio-mcp | — |
| [`drawio-plus`](./drawio-plus/) | Generate clean, standardized, non-overlapping draw.io / diagrams.net diagrams of any kind (architecture, flow, ER, class, network, sequence, mockup) whose boxes never overlap and whose arrows route around boxes, following fixed 80/60/40 spacing and grid stand... | skillify | references · scripts · templates · examples |
| [`extract-anything`](./extract-anything/) | Extract ANY source — a document, a prior workflow stage's output, a spec, code, or notes — into a single LEAN, chainable JSON contract that preserves the source's salient context with bounded, inline-recorded information loss, so workflow stages chain without re-reading the sourc... | skillify | examples · references · schemas · scripts · templates |
| [`generating-pseudocode`](./generating-pseudocode/) | Analyzes requirements or existing code to generate clean, language-agnostic pseudocode that bridges high-level intent and implementation, readable by a Python, Go, or TypeScript developer without translation. Use when... | existing treasury | references |
| [`git-pro`](./git-pro/) | Manage git end-to-end like an expert: branch and stage cleanly, write atomic conventional-commit messages, choose merge vs rebase, resolve conflicts, tidy UNPUBLISHED history (squash, reorder, amend), and undo mistake... | skillify | examples · references · templates |
| [`jira-fix-mr-workflow`](./jira-fix-mr-workflow/) | Run one Jira issue through a gated fix-to-GitLab-MR workflow that advances only on human approval. Use when the user says "fix Jira issue DGL-1234 and open an MR", "run the Jira fix workflow for a given issue key", "f... | jira-flow | — |
| [`milestone-jo-check`](./milestone-jo-check/) | Verify that a vendor's billed milestones stay within the Job Order (JO) limits — that the summed Total Manday across all milestones is within the JO Estimated Manday, and the summed Total Amount is within the JO Estim... | review-milestone | scripts |
| [`organizing-local-files`](./organizing-local-files/) | Plans and applies safe, offline organization of local files and folders. Use when the user says "organize this folder", "clean up my files", "classify my documents", "make folder taxonomy", "generate a safe move plan"... | existing treasury | references · templates |
| [`progressive-bug-hunter`](./progressive-bug-hunter/) | Localize and diagnose a bug by progressively retrieving the MINIMAL sufficient code context: start with cheap agentic grep / structured search, escalate to symbol-graph / call-graph / AST-aware retrieval only when nee... | existing treasury | references · scripts · templates |
| [`publishing-git-review-requests`](./publishing-git-review-requests/) | Publishes one local change to GitHub or GitLab for hosted review: prepare an intended commit, create or attach origin, push the branch, and open a PR or MR. Use when the user asks to "create remote and push", "commit... | existing treasury | references · templates |
| [`rejected-task-replacement`](./rejected-task-replacement/) | For people whose timesheet entries were REJECTED, derive replacement task descriptions sourced from Jira (read-only), written in each person's own style, and save them to a TaskForJO file after the user reviews them.... | review-milestone | references |
| [`rendering-contract-debug-viewer`](./rendering-contract-debug-viewer/) | Render one pipeline JSON contract artifact into a single self-contained, offline HTML viewer styled in the squad-delivery dashboard theme (light plus dark), for local debugging and inspection. Use when the user asks t... | workflow-pack | references · scripts |
| [`rendering-delivery-review-console`](./rendering-delivery-review-console/) | Assemble one run-scoped, offline, byte-deterministic delivery-review.html "Delivery Review Console" from a pipeline run directory: a single standalone document whose left-nav menus are pipeline stages (Epics and Stori... | workflow-pack | references · scripts · templates |
| [`rendering-readable-html`](./rendering-readable-html/) | Render structured content the agent already has - tabular data, a markdown or plain-text report, or findings produced this session - into one clean, self-contained, static HTML file made for a human to read offline. U... | existing treasury | examples · references · scripts · templates |
| [`thai-translator`](./thai-translator/) | Render any source content — pasted text, a file, or a URL — into a structurally identical Thai version written in plain, everyday language a general reader can follow, while preserving the original's meaning, structur... | business-analyse | examples · references |
| [`timesheet-individual-review`](./timesheet-individual-review/) | Review each person's daily timesheet entries, milestone by milestone, against four rules (repeated task over consecutive days, 2-day repeat pattern over the monthly limit, work outside the JO scope, and duplicate work... | review-milestone | references |
| [`transcribing-media-to-minutes`](./transcribing-media-to-minutes/) | Read a meeting recording (audio or video — Thai/English bilingual, usually mostly Thai with English banking and tech terms mixed in) and extract it as faithful Minutes of Meeting (MOM): a Thai-primary, plain-language... | skillify | examples · references · templates |

## Alphabetical Index

| Skill | Purpose Group | Source family |
|-------|---------------|---------------|
| [`agent-context-initializer`](./agent-context-initializer/) | Agent Orchestration & Workflow Infrastructure | existing treasury |
| [`agentic-workflow-design`](./agentic-workflow-design/) | Agent Orchestration & Workflow Infrastructure | business-analyse |
| [`analyzing-banking-requirements`](./analyzing-banking-requirements/) | Banking, BA Delivery & Requirements | business-analyse |
| [`analyzing-canary-rollout`](./analyzing-canary-rollout/) | Observability, Cost & Incident Operations | workflow-pack |
| [`api-contract-design`](./api-contract-design/) | Architecture, Engineering Decisions & Planning | business-analyse |
| [`architecting-fintech-systems`](./architecting-fintech-systems/) | Architecture, Engineering Decisions & Planning | squad-delivery |
| [`architecture-decision`](./architecture-decision/) | Architecture, Engineering Decisions & Planning | business-analyse |
| [`assembling-tl-handoff`](./assembling-tl-handoff/) | Banking, BA Delivery & Requirements | business-analyse |
| [`assessing-ba-feasibility`](./assessing-ba-feasibility/) | Banking, BA Delivery & Requirements | business-analyse |
| [`augment-diagrams`](./augment-diagrams/) | Research, Debate & Knowledge Synthesis | squad-researcher |
| [`authoring-e2e-test-suite`](./authoring-e2e-test-suite/) | Testing, QA & Validation | workflow-pack |
| [`authoring-scaffold-profile`](./authoring-scaffold-profile/) | Agent Orchestration & Workflow Infrastructure | existing treasury |
| [`authoring-workflow-postmortem`](./authoring-workflow-postmortem/) | Observability, Cost & Incident Operations | existing treasury |
| [`befe-contract-design`](./befe-contract-design/) | Architecture, Engineering Decisions & Planning | workflow-pack |
| [`business-logic-extractor`](./business-logic-extractor/) | Code Analysis, Productivity & Publishing | existing treasury |
| [`checking-ba-governance`](./checking-ba-governance/) | Banking, BA Delivery & Requirements | business-analyse |
| [`clean-go-service`](./clean-go-service/) | Implementation, Platform Templates & Code Review | rust-research |
| [`composing-agent-pipelines`](./composing-agent-pipelines/) | Agent Orchestration & Workflow Infrastructure | existing treasury |
| [`configuring-sandbox-allowlist`](./configuring-sandbox-allowlist/) | Security, Governance & Compliance | existing treasury |
| [`contract-testing-pact`](./contract-testing-pact/) | Testing, QA & Validation | workflow-pack |
| [`crafting-backend-code`](./crafting-backend-code/) | Implementation, Platform Templates & Code Review | existing treasury |
| [`crafting-frontend-code`](./crafting-frontend-code/) | Implementation, Platform Templates & Code Review | existing treasury |
| [`crafting-rust-code`](./crafting-rust-code/) | Implementation, Platform Templates & Code Review | existing treasury |
| [`cross-examine`](./cross-examine/) | Research, Debate & Knowledge Synthesis | squad-brainstorm |
| [`daily-planner`](./daily-planner/) | Code Analysis, Productivity & Publishing | skillify |
| [`data-modeling`](./data-modeling/) | Architecture, Engineering Decisions & Planning | business-analyse |
| [`defining-engineering-standards`](./defining-engineering-standards/) | Architecture, Engineering Decisions & Planning | business-analyse |
| [`delegating-to-cli-models`](./delegating-to-cli-models/) | Agent Orchestration & Workflow Infrastructure | skillify |
| [`delivery-planning`](./delivery-planning/) | Architecture, Engineering Decisions & Planning | business-analyse |
| [`designing-tech-lead-handoff`](./designing-tech-lead-handoff/) | Architecture, Engineering Decisions & Planning | squad-delivery |
| [`developing-langgraph-workflows`](./developing-langgraph-workflows/) | Agent Orchestration & Workflow Infrastructure | langgraph-claude-agent |
| [`devops-infrastructure-hardener`](./devops-infrastructure-hardener/) | Security, Governance & Compliance | existing treasury |
| [`discord-harness-extend`](./discord-harness-extend/) | Agent Orchestration & Workflow Infrastructure | discord-llm-harness |
| [`discord-harness-operate`](./discord-harness-operate/) | Agent Orchestration & Workflow Infrastructure | discord-llm-harness |
| [`domain-modeling`](./domain-modeling/) | Architecture, Engineering Decisions & Planning | business-analyse |
| [`drafting-ba-stories`](./drafting-ba-stories/) | Banking, BA Delivery & Requirements | business-analyse |
| [`drafting-stage-prompt`](./drafting-stage-prompt/) | Agent Orchestration & Workflow Infrastructure | existing treasury |
| [`drawio`](./drawio/) | Code Analysis, Productivity & Publishing | drawio-mcp |
| [`drawio-plus`](./drawio-plus/) | Code Analysis, Productivity & Publishing | skillify |
| [`eliciting-banking-brief`](./eliciting-banking-brief/) | Banking, BA Delivery & Requirements | business-analyse |
| [`engineer-growth-planning`](./engineer-growth-planning/) | Architecture, Engineering Decisions & Planning | business-analyse |
| [`engineering-doc-planning`](./engineering-doc-planning/) | Architecture, Engineering Decisions & Planning | business-analyse |
| [`evaluate-banking-compliance`](./evaluate-banking-compliance/) | Banking, BA Delivery & Requirements | craft |
| [`executing-backend-unit-tests`](./executing-backend-unit-tests/) | Testing, QA & Validation | workflow-pack |
| [`executing-frontend-unit-tests`](./executing-frontend-unit-tests/) | Testing, QA & Validation | workflow-pack |
| [`executing-integration-tests`](./executing-integration-tests/) | Testing, QA & Validation | workflow-pack |
| [`executing-qa-test-suite`](./executing-qa-test-suite/) | Testing, QA & Validation | workflow-pack |
| [`extract-anything`](./extract-anything/) | Code Analysis, Productivity & Publishing | skillify |
| [`extract-brief-structure`](./extract-brief-structure/) | Banking, BA Delivery & Requirements | craft |
| [`extract-findings`](./extract-findings/) | Research, Debate & Knowledge Synthesis | squad-researcher |
| [`frame-debate`](./frame-debate/) | Research, Debate & Knowledge Synthesis | squad-brainstorm |
| [`generate-ux-pack`](./generate-ux-pack/) | Architecture, Engineering Decisions & Planning | squad-delivery |
| [`generating-gherkin-acceptance-criteria`](./generating-gherkin-acceptance-criteria/) | Testing, QA & Validation | craft |
| [`generating-pseudocode`](./generating-pseudocode/) | Code Analysis, Productivity & Publishing | existing treasury |
| [`git-pro`](./git-pro/) | Code Analysis, Productivity & Publishing | skillify |
| [`golite`](./golite/) | Implementation, Platform Templates & Code Review | platform-mgmt |
| [`governance-policy-generator`](./governance-policy-generator/) | Security, Governance & Compliance | existing treasury |
| [`handoff-revoke`](./handoff-revoke/) | Agent Orchestration & Workflow Infrastructure | workflow-pack |
| [`handoff-to-deploy`](./handoff-to-deploy/) | Agent Orchestration & Workflow Infrastructure | workflow-pack |
| [`implement-backend-feature`](./implement-backend-feature/) | Implementation, Platform Templates & Code Review | squad-delivery |
| [`implement-frontend-feature`](./implement-frontend-feature/) | Implementation, Platform Templates & Code Review | squad-delivery |
| [`implementing-go-template-requirements`](./implementing-go-template-requirements/) | Implementation, Platform Templates & Code Review | agentic |
| [`incident-response`](./incident-response/) | Observability, Cost & Incident Operations | business-analyse |
| [`integration-design`](./integration-design/) | Architecture, Engineering Decisions & Planning | business-analyse |
| [`jira-fix-mr-workflow`](./jira-fix-mr-workflow/) | Code Analysis, Productivity & Publishing | jira-flow |
| [`langgraph-professional`](./langgraph-professional/) | Implementation, Platform Templates & Code Review | langgraph-claude-agent |
| [`launching-workflow-routines`](./launching-workflow-routines/) | Agent Orchestration & Workflow Infrastructure | skillify |
| [`milestone-jo-check`](./milestone-jo-check/) | Code Analysis, Productivity & Publishing | review-milestone |
| [`model-selection`](./model-selection/) | Architecture, Engineering Decisions & Planning | business-analyse |
| [`multi-agent-handoff-architect`](./multi-agent-handoff-architect/) | Agent Orchestration & Workflow Infrastructure | existing treasury |
| [`observability-design`](./observability-design/) | Observability, Cost & Incident Operations | business-analyse |
| [`observability-telemetry-instrumenter`](./observability-telemetry-instrumenter/) | Observability, Cost & Incident Operations | existing treasury |
| [`opening-debate-panel`](./opening-debate-panel/) | Research, Debate & Knowledge Synthesis | squad-brainstorm |
| [`orchestrate-banking-brief-pipeline`](./orchestrate-banking-brief-pipeline/) | Banking, BA Delivery & Requirements | craft |
| [`orchestrating-agent-scaffold`](./orchestrating-agent-scaffold/) | Agent Orchestration & Workflow Infrastructure | existing treasury |
| [`orchestrating-openclaw-squad`](./orchestrating-openclaw-squad/) | Agent Orchestration & Workflow Infrastructure | existing treasury |
| [`organizing-local-files`](./organizing-local-files/) | Code Analysis, Productivity & Publishing | existing treasury |
| [`panel-open`](./panel-open/) | Research, Debate & Knowledge Synthesis | squad-brainstorm |
| [`performance-cost-review`](./performance-cost-review/) | Observability, Cost & Incident Operations | business-analyse |
| [`plan-research`](./plan-research/) | Research, Debate & Knowledge Synthesis | squad-researcher |
| [`planning-banking-tests`](./planning-banking-tests/) | Testing, QA & Validation | squad-delivery |
| [`platform-common`](./platform-common/) | Implementation, Platform Templates & Code Review | agentic |
| [`platform-go-service`](./platform-go-service/) | Implementation, Platform Templates & Code Review | agentic |
| [`pr-design-review`](./pr-design-review/) | Architecture, Engineering Decisions & Planning | business-analyse |
| [`principal-advisor`](./principal-advisor/) | Architecture, Engineering Decisions & Planning | skillify |
| [`production-readiness`](./production-readiness/) | Architecture, Engineering Decisions & Planning | business-analyse |
| [`progressive-bug-hunter`](./progressive-bug-hunter/) | Code Analysis, Productivity & Publishing | existing treasury |
| [`publishing-git-review-requests`](./publishing-git-review-requests/) | Code Analysis, Productivity & Publishing | existing treasury |
| [`red-teaming-implementation-plan`](./red-teaming-implementation-plan/) | Architecture, Engineering Decisions & Planning | workflow-pack |
| [`refactor-decision`](./refactor-decision/) | Architecture, Engineering Decisions & Planning | business-analyse |
| [`refactoring-go-services`](./refactoring-go-services/) | Implementation, Platform Templates & Code Review | agentic |
| [`rejected-task-replacement`](./rejected-task-replacement/) | Code Analysis, Productivity & Publishing | review-milestone |
| [`rendering-contract-debug-viewer`](./rendering-contract-debug-viewer/) | Code Analysis, Productivity & Publishing | workflow-pack |
| [`rendering-delivery-review-console`](./rendering-delivery-review-console/) | Code Analysis, Productivity & Publishing | workflow-pack |
| [`rendering-readable-html`](./rendering-readable-html/) | Code Analysis, Productivity & Publishing | existing treasury |
| [`report-debate`](./report-debate/) | Research, Debate & Knowledge Synthesis | squad-brainstorm |
| [`reporting-research-run`](./reporting-research-run/) | Research, Debate & Knowledge Synthesis | squad-researcher |
| [`research-vault-librarian`](./research-vault-librarian/) | Research, Debate & Knowledge Synthesis | existing treasury |
| [`researching-ba-problem-space`](./researching-ba-problem-space/) | Banking, BA Delivery & Requirements | workflow-pack |
| [`review-backend-code`](./review-backend-code/) | Implementation, Platform Templates & Code Review | agentic |
| [`review-frontend-code`](./review-frontend-code/) | Implementation, Platform Templates & Code Review | squad-delivery |
| [`review-report`](./review-report/) | Research, Debate & Knowledge Synthesis | squad-researcher |
| [`review-rust-code`](./review-rust-code/) | Implementation, Platform Templates & Code Review | existing treasury |
| [`reviewing-agent-spend`](./reviewing-agent-spend/) | Observability, Cost & Incident Operations | existing treasury |
| [`reviewing-implement-gate`](./reviewing-implement-gate/) | Agent Orchestration & Workflow Infrastructure | existing treasury |
| [`reviewing-software-security`](./reviewing-software-security/) | Security, Governance & Compliance | squad-delivery |
| [`revise-positions`](./revise-positions/) | Research, Debate & Knowledge Synthesis | squad-brainstorm |
| [`risk-estimation`](./risk-estimation/) | Architecture, Engineering Decisions & Planning | business-analyse |
| [`roundtripping-dashboard-data-contract`](./roundtripping-dashboard-data-contract/) | Agent Orchestration & Workflow Infrastructure | workflow-pack |
| [`running-accessibility-tests`](./running-accessibility-tests/) | Testing, QA & Validation | workflow-pack |
| [`running-ba-pipeline`](./running-ba-pipeline/) | Banking, BA Delivery & Requirements | business-analyse |
| [`running-business-analysis-workflow`](./running-business-analysis-workflow/) | Banking, BA Delivery & Requirements | craft |
| [`running-performance-load-test`](./running-performance-load-test/) | Observability, Cost & Incident Operations | workflow-pack |
| [`running-sast-security-gate`](./running-sast-security-gate/) | Security, Governance & Compliance | workflow-pack |
| [`running-smoke-tests`](./running-smoke-tests/) | Testing, QA & Validation | workflow-pack |
| [`scanning-appsec-pipeline-gate`](./scanning-appsec-pipeline-gate/) | Security, Governance & Compliance | workflow-pack |
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
| [`timesheet-individual-review`](./timesheet-individual-review/) | Code Analysis, Productivity & Publishing | review-milestone |
| [`transcribing-media-to-minutes`](./transcribing-media-to-minutes/) | Code Analysis, Productivity & Publishing | skillify |
| [`universal-spec-validator`](./universal-spec-validator/) | Security, Governance & Compliance | existing treasury |
| [`validating-banking-implementation`](./validating-banking-implementation/) | Testing, QA & Validation | squad-delivery |
| [`validating-production-slo`](./validating-production-slo/) | Observability, Cost & Incident Operations | workflow-pack |
