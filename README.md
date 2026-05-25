# Skillify — AI Agent Skill Engineering Toolkit

A research-backed toolkit for creating, validating, and maintaining production-grade `SKILL.md` files for AI coding agents. Built from a systematic review of 56 primary and supporting sources across the Anthropic, OpenAI, GitHub Copilot, open-standard, platform-rule, and academic agent ecosystems.

## What is a Skill?

A skill is a folder containing a `SKILL.md` file — structured instructions in Markdown with YAML frontmatter — that teaches an AI agent how to perform a specific, reusable workflow. Skills are the primary mechanism for extending agent capabilities beyond their base training.

```
my-skill/
├── SKILL.md            # Core instructions (YAML frontmatter + Markdown body)
├── references/         # Deep guidance loaded only when needed
├── templates/          # Reusable output skeletons
├── scripts/            # Deterministic automation (validators, generators)
└── examples/           # Worked examples for agent calibration
```

## Repository Structure

```
.
├── literature/              # 56 curated source documents
│   ├── anthropic-claude/    # 13 files — Claude Code, sub-agents, skills docs
│   ├── codex-copilot/       #  9 files — Codex agents, GitHub Copilot skills, VS Code agents
│   ├── openai/              #  4 files — OpenAI skill creator, API tools, cookbook
│   ├── best-practices/      #  4 files — Cross-platform skill authoring guides
│   ├── awesome-lists/       #  2 files — Community skill aggregations
│   ├── open-standard/       #  1 file  — Agent Skills specification
│   ├── other-platforms/     # 13 files — Gemini, Cursor, Cline, OpenCode, Windsurf, and others
│   ├── research-papers/     #  6 files — ReAct, Reflexion, Voyager, Toolformer, MemGPT, SWE-agent
│   └── skill-design-methodology/ # 4 files — synthesized source index, taxonomy, principles, literature review
│
├── skillify/                # The skill-creator meta-skill
│   ├── SKILL.md             # Core skill definition (8 modes)
│   ├── references/          # 9 reference guides
│   │   ├── anti-patterns.md
│   │   ├── frontmatter-guide.md
│   │   ├── lifecycle-and-iteration.md
│   │   ├── mode-playbooks.md
│   │   ├── platform-compatibility.md
│   │   ├── progressive-disclosure.md
│   │   ├── security-checklist.md
│   │   ├── validation-rubric.md
│   │   └── workflow-patterns.md
│   ├── templates/           # 4 starter templates
│   │   ├── basic-skill-template.md
│   │   ├── domain-skill-template.md
│   │   ├── mcp-skill-template.md
│   │   └── audit-report-template.md
│   ├── scripts/             # 3 automation scripts
│   │   ├── init_skill.py
│   │   ├── quick_validate.py
│   │   └── check_links.py
│   └── examples/            # 3 worked examples
│       ├── create-from-scratch.md
│       ├── good-description-examples.md
│       └── skill-audit-walkthrough.md
│
└── treasury/                # 36 production skills in 7 categories
    ├── README.md            # ← Skill catalog (start here)
    └── <skill>/             # each: SKILL.md + references/ (+ templates, schemas, …)
```

## Core Concepts

### Progressive Disclosure

Skills use a 3-tier loading model to minimize context-window consumption:

| Tier | What | When Loaded | Budget |
|------|------|-------------|--------|
| **1 — Metadata** | `name` + `description` in YAML frontmatter | Always in context | ~100 words |
| **2 — Body** | `SKILL.md` Markdown instructions | When skill triggers | < 5,000 words |
| **3 — References** | `references/`, `templates/`, `scripts/` | On demand by agent | Unlimited |

### Skillify Modes

The `skillify` meta-skill operates in 8 modes:

| Mode | Trigger | Output |
|------|---------|--------|
| **Create** | "create a skill", "write a SKILL.md" | New skill folder |
| **Refactor** | "refactor my skill", "improve this skill" | Modified skill folder |
| **Review** | "review my skill", "is this skill good" | Findings report (no edits) |
| **Audit** | "audit my skill", "score my skill" | Scored rubric (10 dimensions, 50 pts) |
| **Compress** | "shrink this skill", "reduce token cost" | Leaner skill with Tier-3 migration |
| **Split** | "split this skill", "this skill does too much" | Multiple focused skill folders |
| **Merge** | "merge these skills", "combine X and Y" | One merged skill |
| **Adapt** | "adapt this for Codex", "make this work in Copilot" | Platform-adapted skill |

### Validation Rubric

Every skill is scored across 10 dimensions (pass threshold: 40/50):

1. Trigger Quality — Specific user intents in description
2. Scope Focus — Single, clear responsibility
3. Workflow Clarity — Imperative, numbered steps with entry/exit conditions
4. Output Contract — Exact deliverables specified
5. Token Efficiency — Lean body, deep content in references
6. Conflict Risk — No overlap with repo-policy or other skills
7. Reusability — Portable across projects
8. Security & Safety — No destructive commands, secrets, or broad permissions
9. Frontmatter Correctness — Valid YAML, kebab-case name, triggers present
10. Progressive Disclosure — Content lives in exactly one tier

## Quick Start

### Install the meta-skill

Copy the `skillify/` folder into your project's skills directory:

```bash
cp -r skillify/ /path/to/your-project/.skills/skillify/
```

### Create a new skill

```bash
# Generate boilerplate
python skillify/scripts/init_skill.py my-new-skill

# Then ask your agent:
# "Create a skill for [your workflow]"
```

### Validate an existing skill

```bash
python skillify/scripts/quick_validate.py path/to/skill-folder
python skillify/scripts/check_links.py path/to/skill-folder
```

## Treasury — Production Skill Library

The `treasury/` directory holds **36 production-grade skills** built with the `skillify` meta-skill, grouped into 7 functional categories. Every skill is a folder with its own `SKILL.md` and a `references/` tier. Full per-skill detail, asset breakdown, and an alphabetical index live in the catalog: [`treasury/README.md`](treasury/README.md).

### 1. Banking & Fintech Domain (5)

| Skill | Role | What it does |
|-------|------|--------------|
| [`analyzing-banking-requirements`](treasury/analyzing-banking-requirements/) | Business Analyst | BA persona: extract raw banking requests into PRDs/API contracts with strict KYC/AML/PCI-DSS compliance mapping |
| [`architecting-fintech-systems`](treasury/architecting-fintech-systems/) | Senior Fintech Architect | Senior fintech architect for L1–L3 design via DDD/CQRS/event-driven on a Go/AWS/GCP stack |
| [`eliciting-banking-brief`](treasury/eliciting-banking-brief/) | Business Analyst | Turn raw BA input (Jira/Slack/notes) into a structured epic-plus-stories brief; governance gaps as P1 blockers |
| [`planning-banking-tests`](treasury/planning-banking-tests/) | Quality Assurance Engineer | Convert a BA brief into a QA test plan: Gherkin scenarios, NFRs, regulatory deps, compliance, sign-off criteria |
| [`validating-banking-implementation`](treasury/validating-banking-implementation/) | Quality Assurance Engineer | QA persona: adversarial OWASP, transactional-integrity / race / deadlock audit, chaos plan; Approve/Reject verdict |

### 2. Code Implementation & Review (8)

| Skill | Role | What it does |
|-------|------|--------------|
| [`crafting-backend-code`](treasury/crafting-backend-code/) | Senior Backend Engineer | Design, review, and safely implement backend services (Go/Node/Python/Java), pattern-first and evidence-led |
| [`implement-backend-feature`](treasury/implement-backend-feature/) | Backend Software Developer | Generate production-grade Go backend code (HTTP / CQRS / Kafka) for one feature from an approved design |
| [`implementing-go-template-requirements`](treasury/implementing-go-template-requirements/) | Backend Software Developer | Apply one requirement to a go-template service by editing business logic under `app/<domain>/` plus narrow router wiring |
| [`review-backend-code`](treasury/review-backend-code/) | Senior Backend Engineer | Adversarially verify Go code vs the design and 11 banking-grade rules; approve / loop_back / human-queue verdict |
| [`review-rust-code`](treasury/review-rust-code/) | Senior Rust Engineer | Adversarially verify Rust code vs the design and 11 banking-grade rules recast for Rust; machine-readable verdict |
| [`crafting-frontend-code`](treasury/crafting-frontend-code/) | Senior Frontend Engineer | Design, review, and safely implement React/TypeScript frontend with a conservative, repo-first posture |
| [`implement-frontend-feature`](treasury/implement-frontend-feature/) | Frontend Software Developer | Generate production-grade React/TS from an approved UI design: WCAG 2.1 AA, no localStorage auth, PII handling |
| [`review-frontend-code`](treasury/review-frontend-code/) | Senior Frontend Engineer | Adversarially verify React/TS vs the UI design and 12 banking-grade non-negotiables; machine-readable verdict |

### 3. Agent Orchestration & Pipelines (5)

| Skill | Role | What it does |
|-------|------|--------------|
| [`agent-context-initializer`](treasury/agent-context-initializer/) | Agent Systems Architect | Generate a minimal AGENTS.md (100–150 lines, six sections, three-tier boundaries) plus a curation checklist |
| [`composing-agent-pipelines`](treasury/composing-agent-pipelines/) | Agent Orchestrator | Compose a portable multi-agent pipeline (Plan/Gather/Analyze/Review/Validate/Decide/Compact) with versioned artifacts |
| [`multi-agent-handoff-architect`](treasury/multi-agent-handoff-architect/) | Agent Integration Architect | Design inter-agent handoff APIs as versioned JSON Schema contracts; single-writer ownership, autonomy by reversibility |
| [`orchestrating-agent-scaffold`](treasury/orchestrating-agent-scaffold/) | Agent Operations Engineer | Orchestrate a multi-stage research-squad run on agent-scaffold: plan, pick profile/cap, monitor state.json |
| [`orchestrating-openclaw-squad`](treasury/orchestrating-openclaw-squad/) | Technical Program Manager | Orchestrate a BA/Architect/Developer/QA squad with strict human-in-the-loop gates and sandbox isolation |

### 4. Agent-Scaffold Infrastructure (4)

| Skill | Role | What it does |
|-------|------|--------------|
| [`authoring-scaffold-profile`](treasury/authoring-scaffold-profile/) | Agent Infrastructure Engineer | Author/modify agent-scaffold `profiles/NAME.sh` (STAGES, GATED_STAGES, per-stage AGENT/MODEL) with validation |
| [`configuring-sandbox-allowlist`](treasury/configuring-sandbox-allowlist/) | Platform Security Engineer | Edit the sandbox hostname allowlist; refuse wildcard/IP/metadata patterns; remind to rebuild + egress-test |
| [`drafting-stage-prompt`](treasury/drafting-stage-prompt/) | Prompt Engineer | Curate stage prompts into `prompts/library/STAGE/TOPIC.md` with why-it-works, success metric, failure mode |
| [`reviewing-implement-gate`](treasury/reviewing-implement-gate/) | Technical Lead | Walk a researcher through the five-question check before approving the implement gate; approve/reject + command |

### 5. Security, Governance & Compliance (4)

| Skill | Role | What it does |
|-------|------|--------------|
| [`reviewing-software-security`](treasury/reviewing-software-security/) | Security Engineer | Defensive security review for Go/Gin, Kafka, MySQL, K8s, Kong/APISIX, lending flows; mapped to OWASP/CWE/NIST/CIS |
| [`devops-infrastructure-hardener`](treasury/devops-infrastructure-hardener/) | DevOps Engineer | Audit agent CI/CD for static credentials; emit remediation + short-lived OIDC / dynamic-secrets architecture |
| [`universal-spec-validator`](treasury/universal-spec-validator/) | Platform Engineer | CI / pre-commit gate validating agent specs for cross-model drift, unsafe command surface, breaking schema evolution |
| [`governance-policy-generator`](treasury/governance-policy-generator/) | Governance & Risk Officer | Emit policy-as-code: default-deny OPA/Rego allowlist plus KILLSWITCH.md (triggers, escalation, append-only audit) |

### 6. Observability & Cost Governance (3)

| Skill | Role | What it does |
|-------|------|--------------|
| [`observability-telemetry-instrumenter`](treasury/observability-telemetry-instrumenter/) | Observability Engineer | Instrument agent code with OpenTelemetry GenAI conventions: invoke_agent / execute_tool spans, token histogram |
| [`reviewing-agent-spend`](treasury/reviewing-agent-spend/) | FinOps Analyst | Monday cost-review ritual: segment LiteLLM spend by user/model, flag >2σ outliers and workflows over $10 |
| [`authoring-workflow-postmortem`](treasury/authoring-workflow-postmortem/) | Incident Response Manager | Author a postmortem within 48h of a trigger event (failed stage >$1, rejected gate, cap overrun, egress fail) |

### 7. Code Analysis & Productivity (7)

| Skill | Role | What it does |
|-------|------|--------------|
| [`business-logic-extractor`](treasury/business-logic-extractor/) | Requirements Analyst | Extract implemented business logic from code into a traceable spec with file:line provenance + loss ledger |
| [`generating-pseudocode`](treasury/generating-pseudocode/) | Algorithm Designer | Analyze requirements or code into clean language-agnostic pseudocode plus framing and a verification note |
| [`progressive-bug-hunter`](treasury/progressive-bug-hunter/) | Debugging Specialist | Localize/diagnose a bug via minimal progressive retrieval (grep → symbol-graph → AST); ranked diagnosis report |
| [`organizing-local-files`](treasury/organizing-local-files/) | Information Architect | Plan and apply safe offline file organization: inventory, JSON move plan, apply.sh, rollback.sh; never deletes |
| [`research-vault-librarian`](treasury/research-vault-librarian/) | Research Librarian | Read-only librarian for a CLAUDE.md-governed Obsidian research vault: query reports, audit MOC/wikilink/citation drift, emit a ready-to-apply intake patch |
| [`publishing-git-review-requests`](treasury/publishing-git-review-requests/) | Source Control Engineer | Publish one local change to GitHub/GitLab: prepare a commit, attach origin, push the branch, open a PR or MR |
| [`rendering-readable-html`](treasury/rendering-readable-html/) | Technical Writer | Render data, a report, or session findings into one self-contained, static, JavaScript-free HTML file built for a human to read offline |

## Literature Sources

The `literature/` directory contains the complete research corpus — 56 documents totaling ~758K words — organized by ecosystem:

| Category | Files | Key Sources |
|----------|-------|-------------|
| **Anthropic / Claude** | 13 | Claude Code Best Practices, Skills Repo, Sub-Agents, Context Engineering |
| **Codex / Copilot** | 9 | Codex Skills & Agents, GitHub CLI Agent Skills, VS Code Copilot, VS Code Custom Agents |
| **OpenAI** | 4 | Skill Creator (canonical reference), API Tools, Cookbook |
| **Best Practices** | 4 | Cross-platform authoring guides, SKILL.md pattern analysis |
| **Awesome Lists** | 2 | ScienceAIX and VoltAgent community aggregations |
| **Open Standard** | 1 | Open Agent Skills specification |
| **Other Platforms** | 13 | Gemini CLI, Cursor, Cline, OpenCode, Windsurf, Spring AI, Strapi, antfu |
| **Research Papers** | 6 | ReAct, Reflexion, Voyager, Toolformer, MemGPT, SWE-agent |
| **Skill Design Methodology** | 4 | Source index, reading taxonomy, synthesized design principles, literature review |

## Platform Compatibility

Skillify generates portable `SKILL.md` folders compatible with:

- **Claude Code** - `.claude/skills/` or `~/.claude/skills/`
- **OpenAI Codex** - `.agents/skills/`, `~/.agents/skills/`, or `$CODEX_HOME/skills/`
- **GitHub Copilot** - `.github/skills/`, `.agents/skills/`, `~/.copilot/skills/`, or optional custom instructions
- **Gemini / Antigravity** - `.agents/skills/`, `.gemini/skills/`, `~/.gemini/skills/`, or `~/.gemini/antigravity/skills/`

Use the **Adapt** mode to convert between platform conventions.

## License

This repository is a research and documentation project. The literature files retain their original authors' licensing. The `skillify/` meta-skill and tooling are provided as-is for internal use.
