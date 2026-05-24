# Treasury — Production Skill Library

36 production-grade skills, grouped by function. Each skill name links to its folder. Every skill ships a `references/` tier; the **Extra assets** column lists only what a skill adds on top of that.

| Category | Count |
|----------|-------|
| [1. Banking & Fintech Domain](#1-banking--fintech-domain-5) | 5 |
| [2. Code Implementation & Review](#2-code-implementation--review-8) | 8 |
| [3. Agent Orchestration & Pipelines](#3-agent-orchestration--pipelines-5) | 5 |
| [4. Agent-Scaffold Infrastructure](#4-agent-scaffold-infrastructure-4) | 4 |
| [5. Security, Governance & Compliance](#5-security-governance--compliance-4) | 4 |
| [6. Observability & Cost Governance](#6-observability--cost-governance-3) | 3 |
| [7. Code Analysis & Productivity](#7-code-analysis--productivity-7) | 7 |
| **Total** | **36** |

---

### 1. Banking & Fintech Domain (5)

The regulated BA → Architect → QA spine: requirements, architecture, and verification for KYC/AML/lending workflows.

| Skill | What it does | Extra assets |
|-------|--------------|--------------|
| [`analyzing-banking-requirements`](./analyzing-banking-requirements/) | Business Analyst persona: extract raw banking requests into PRDs/API contracts with strict KYC/AML/PCI-DSS compliance mapping; stops with a P1 blocker on a violation. | — |
| [`architecting-fintech-systems`](./architecting-fintech-systems/) | Senior fintech architect for L1–L3 design (business intent → bounded contexts → strategy) via DDD/CQRS/event-driven on a Go/AWS/GCP stack. | templates |
| [`eliciting-banking-brief`](./eliciting-banking-brief/) | Turn raw BA input (Jira/Slack/notes/email) into a structured epic-plus-stories brief; governance gaps escalated as P1 blockers. | audit · schemas · scripts · tests |
| [`planning-banking-tests`](./planning-banking-tests/) | Convert a BA brief into a structured QA test plan: Gherkin scenarios, NFR targets, regulatory deps, compliance, sign-off criteria. | audit · runs · schemas · scripts · tests |
| [`validating-banking-implementation`](./validating-banking-implementation/) | QA persona: adversarial OWASP testing, transactional-integrity / race / deadlock audit, chaos plan; explicit Approve/Reject verdict. | — |

### 2. Code Implementation & Review (8)

Design, generate, and adversarially review backend (Go/Rust) and frontend (React/TS) code to banking-grade rules.

| Skill | What it does | Extra assets |
|-------|--------------|--------------|
| [`crafting-backend-code`](./crafting-backend-code/) | Design, review, and safely implement backend services (Go/Node/Python/Java) with a pattern-first, evidence-led posture. | platforms |
| [`implement-backend-feature`](./implement-backend-feature/) | Generate production-grade Go backend code (HTTP / CQRS / Kafka) for one feature from an approved design, banking-grade. | schemas · tests |
| [`implementing-go-template-requirements`](./implementing-go-template-requirements/) | Apply one requirement to a go-template service by editing business logic under `app/<domain>/` plus narrow router wiring while locking scaffold files. | examples · templates |
| [`review-backend-code`](./review-backend-code/) | Adversarially verify Go code vs the approved design and 11 banking-grade rules; machine-readable approve / loop_back / human-queue. | schemas · tests |
| [`review-rust-code`](./review-rust-code/) | Adversarially verify Rust code vs the approved design and 11 banking-grade rules recast for Rust idioms; machine-readable verdict. | schemas · tests |
| [`crafting-frontend-code`](./crafting-frontend-code/) | Design, review, and safely implement React/TypeScript frontend with a conservative, repo-first posture. | — |
| [`implement-frontend-feature`](./implement-frontend-feature/) | Generate production-grade React/TS for one feature from an approved UI design: WCAG 2.1 AA, no localStorage auth, PII handling. | schemas · tests |
| [`review-frontend-code`](./review-frontend-code/) | Adversarially verify React/TS vs the approved UI design and 12 banking-grade non-negotiables; machine-readable verdict. | schemas · tests |

### 3. Agent Orchestration & Pipelines (5)

Coordinate multi-agent squads and pipelines, and define how agents hand off to each other.

| Skill | What it does | Extra assets |
|-------|--------------|--------------|
| [`agent-context-initializer`](./agent-context-initializer/) | Generate a minimal AGENTS.md (100–150 lines, six sections, paired prohibitions, three-tier boundaries) plus a curation checklist. | scripts · templates |
| [`composing-agent-pipelines`](./composing-agent-pipelines/) | Compose a portable multi-agent pipeline (Plan/Gather/Analyze/Review/Validate/Decide/Compact) with versioned per-phase artifacts. | examples · scripts · templates |
| [`multi-agent-handoff-architect`](./multi-agent-handoff-architect/) | Design inter-agent handoff APIs as versioned JSON Schema contracts; single-writer ownership, autonomy calibrated to reversibility. | schemas · scripts · templates |
| [`orchestrating-agent-scaffold`](./orchestrating-agent-scaffold/) | Orchestrate a multi-stage research-squad run on agent-scaffold: plan, pick profile/cap, spawn pre-flight, monitor state.json. | scripts · templates |
| [`orchestrating-openclaw-squad`](./orchestrating-openclaw-squad/) | Orchestrate a BA/Architect/Developer/QA squad with strict human-in-the-loop gates and sandbox isolation for regulated work. | templates |

### 4. Agent-Scaffold Infrastructure (4)

Configure and curate the agent-scaffold framework itself: profiles, prompts, allowlists, and gates.

| Skill | What it does | Extra assets |
|-------|--------------|--------------|
| [`authoring-scaffold-profile`](./authoring-scaffold-profile/) | Author/modify agent-scaffold `profiles/NAME.sh` (STAGES, GATED_STAGES, per-stage AGENT/MODEL/PROMPT_PREFIX) with validation. | templates |
| [`configuring-sandbox-allowlist`](./configuring-sandbox-allowlist/) | Edit the sandbox hostname allowlist; refuse wildcard/IP/metadata patterns; remind to rebuild the proxy and run the egress test. | templates |
| [`drafting-stage-prompt`](./drafting-stage-prompt/) | Curate stage prompts into `prompts/library/STAGE/TOPIC.md` with why-it-works, success metric, failure mode, model affinity. | templates |
| [`reviewing-implement-gate`](./reviewing-implement-gate/) | Walk a researcher through the five-question check before approving the implement gate; returns approve/reject + the exact command. | templates |

### 5. Security, Governance & Compliance (4)

Defensive review, secrets hardening, spec gating, and policy-as-code for autonomous agents.

| Skill | What it does | Extra assets |
|-------|--------------|--------------|
| [`reviewing-software-security`](./reviewing-software-security/) | Defensive security review for Go/Gin, Kafka, MySQL, K8s, Kong/APISIX, lending flows; findings mapped to OWASP/CWE/NIST/CIS. | templates |
| [`devops-infrastructure-hardener`](./devops-infrastructure-hardener/) | Audit agent CI/CD for static credentials; emit remediation plus a short-lived OIDC / dynamic-secrets runtime architecture. | scripts · templates |
| [`universal-spec-validator`](./universal-spec-validator/) | CI / pre-commit gate validating agent specs for cross-model drift, unsafe command surface, and breaking schema evolution. | scripts · templates |
| [`governance-policy-generator`](./governance-policy-generator/) | Emit policy-as-code: default-deny OPA/Rego allowlist plus KILLSWITCH.md (triggers, escalation, append-only JSONL audit). | scripts · templates |

### 6. Observability & Cost Governance (3)

Make agent runs observable and accountable: telemetry, spend review, and postmortems.

| Skill | What it does | Extra assets |
|-------|--------------|--------------|
| [`observability-telemetry-instrumenter`](./observability-telemetry-instrumenter/) | Instrument agent code with OpenTelemetry GenAI conventions: invoke_agent / execute_tool spans and the token-usage histogram. | schemas · scripts · templates |
| [`reviewing-agent-spend`](./reviewing-agent-spend/) | Monday cost-review ritual: segment LiteLLM spend by user/model, flag >2σ outliers and any workflow over $10, draft a readout. | templates |
| [`authoring-workflow-postmortem`](./authoring-workflow-postmortem/) | Author a postmortem within 48h of a trigger event (failed stage >$1, rejected gate, cap overrun, sandbox egress failure). | templates |

### 7. Code Analysis & Productivity (7)

Cross-cutting analysis and delivery helpers: extract logic, plan algorithms, hunt bugs, organize files, navigate a research vault, publish reviews, render readable HTML.

| Skill | What it does | Extra assets |
|-------|--------------|--------------|
| [`business-logic-extractor`](./business-logic-extractor/) | Extract implemented business logic from code into a traceable spec with file:line provenance and an information-loss ledger. | scripts · templates |
| [`generating-pseudocode`](./generating-pseudocode/) | Analyze requirements or code into clean language-agnostic pseudocode plus a problem-framing block and a verification note. | — |
| [`progressive-bug-hunter`](./progressive-bug-hunter/) | Localize/diagnose a bug via minimal progressive retrieval (grep → symbol-graph → AST); stops at a ranked diagnosis report. | scripts · templates |
| [`organizing-local-files`](./organizing-local-files/) | Plan and apply safe offline file organization: inventory, JSON move plan, report, apply.sh, rollback.sh; never deletes. | templates |
| [`research-vault-librarian`](./research-vault-librarian/) | Read-only librarian for a CLAUDE.md-governed Obsidian research vault: query reports, audit MOC/wikilink/citation drift, emit a ready-to-apply intake patch. | scripts |
| [`publishing-git-review-requests`](./publishing-git-review-requests/) | Publish one local change to GitHub/GitLab: prepare a commit, attach origin, push the branch, open a PR or MR. | templates |
| [`rendering-readable-html`](./rendering-readable-html/) | Render data, a report, or session findings into one self-contained, static, JavaScript-free HTML file built for a human to read offline. | examples · scripts · templates |

---

### Alphabetical Index

| Skill | Category |
|-------|----------|
| [`agent-context-initializer`](./agent-context-initializer/) | Agent Orchestration & Pipelines |
| [`analyzing-banking-requirements`](./analyzing-banking-requirements/) | Banking & Fintech Domain |
| [`architecting-fintech-systems`](./architecting-fintech-systems/) | Banking & Fintech Domain |
| [`authoring-scaffold-profile`](./authoring-scaffold-profile/) | Agent-Scaffold Infrastructure |
| [`authoring-workflow-postmortem`](./authoring-workflow-postmortem/) | Observability & Cost Governance |
| [`business-logic-extractor`](./business-logic-extractor/) | Code Analysis & Productivity |
| [`composing-agent-pipelines`](./composing-agent-pipelines/) | Agent Orchestration & Pipelines |
| [`configuring-sandbox-allowlist`](./configuring-sandbox-allowlist/) | Agent-Scaffold Infrastructure |
| [`crafting-backend-code`](./crafting-backend-code/) | Code Implementation & Review |
| [`crafting-frontend-code`](./crafting-frontend-code/) | Code Implementation & Review |
| [`devops-infrastructure-hardener`](./devops-infrastructure-hardener/) | Security, Governance & Compliance |
| [`drafting-stage-prompt`](./drafting-stage-prompt/) | Agent-Scaffold Infrastructure |
| [`eliciting-banking-brief`](./eliciting-banking-brief/) | Banking & Fintech Domain |
| [`generating-pseudocode`](./generating-pseudocode/) | Code Analysis & Productivity |
| [`governance-policy-generator`](./governance-policy-generator/) | Security, Governance & Compliance |
| [`implement-backend-feature`](./implement-backend-feature/) | Code Implementation & Review |
| [`implement-frontend-feature`](./implement-frontend-feature/) | Code Implementation & Review |
| [`implementing-go-template-requirements`](./implementing-go-template-requirements/) | Code Implementation & Review |
| [`multi-agent-handoff-architect`](./multi-agent-handoff-architect/) | Agent Orchestration & Pipelines |
| [`observability-telemetry-instrumenter`](./observability-telemetry-instrumenter/) | Observability & Cost Governance |
| [`orchestrating-agent-scaffold`](./orchestrating-agent-scaffold/) | Agent Orchestration & Pipelines |
| [`orchestrating-openclaw-squad`](./orchestrating-openclaw-squad/) | Agent Orchestration & Pipelines |
| [`organizing-local-files`](./organizing-local-files/) | Code Analysis & Productivity |
| [`planning-banking-tests`](./planning-banking-tests/) | Banking & Fintech Domain |
| [`progressive-bug-hunter`](./progressive-bug-hunter/) | Code Analysis & Productivity |
| [`publishing-git-review-requests`](./publishing-git-review-requests/) | Code Analysis & Productivity |
| [`rendering-readable-html`](./rendering-readable-html/) | Code Analysis & Productivity |
| [`research-vault-librarian`](./research-vault-librarian/) | Code Analysis & Productivity |
| [`review-backend-code`](./review-backend-code/) | Code Implementation & Review |
| [`review-frontend-code`](./review-frontend-code/) | Code Implementation & Review |
| [`review-rust-code`](./review-rust-code/) | Code Implementation & Review |
| [`reviewing-agent-spend`](./reviewing-agent-spend/) | Observability & Cost Governance |
| [`reviewing-implement-gate`](./reviewing-implement-gate/) | Agent-Scaffold Infrastructure |
| [`reviewing-software-security`](./reviewing-software-security/) | Security, Governance & Compliance |
| [`universal-spec-validator`](./universal-spec-validator/) | Security, Governance & Compliance |
| [`validating-banking-implementation`](./validating-banking-implementation/) | Banking & Fintech Domain |

**Assets legend** — every skill includes `references/`. Extra assets shown per row: `templates` (output skeletons) · `scripts` (deterministic automation) · `schemas` (structured I/O contracts) · `tests` (skill self-tests/fixtures) · `examples` (worked calibration) · `audit` / `runs` / `platforms` (skill-specific). A `—` means references-only.

> One-line summaries are distilled from each skill's `SKILL.md` `description:` frontmatter — open the linked folder for the full trigger list and the authoritative description.
