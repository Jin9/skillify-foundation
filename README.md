# Skillify — AI Agent Skill Engineering Toolkit

A research-backed toolkit for creating, validating, and maintaining production-grade `SKILL.md` files for AI coding agents. Built from a systematic review of 87 primary and supporting sources across the Anthropic, OpenAI, GitHub Copilot, open-standard, platform-rule, academic, and internal deep-research agent ecosystems.

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
├── literature/              # 87 curated source documents
├── skillify/                # The skill-creator meta-skill and validation tooling
└── treasury/                # 131 top-level skills grouped by purpose
    ├── README.md            # Skill catalog grouped by purpose
    └── <skill>/             # each: SKILL.md plus optional local assets
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

1. Trigger Quality
2. Scope Focus
3. Workflow Clarity
4. Output Contract
5. Token Efficiency
6. Conflict Risk
7. Reusability
8. Security & Safety
9. Frontmatter Correctness
10. Progressive Disclosure

## Quick Start

### Install the meta-skill

Copy the `skillify/` folder into your project's skills directory:

```bash
cp -r skillify/ /path/to/your-project/.skills/skillify/
```

### Create a new skill

```bash
# Generate boilerplate
python3 skillify/scripts/init_skill.py my-new-skill

# Then ask your agent:
# "Create a skill for [your workflow]"
```

### Validate an existing skill

```bash
python3 skillify/scripts/quick_validate.py path/to/skill-folder
python3 skillify/scripts/check_links.py path/to/skill-folder
```

## Treasury — Production Skill Library

The `treasury/` directory holds **131 top-level skills** grouped by purpose. Full per-skill detail, source family, asset breakdown, curated folder names, and an alphabetical index live in the catalog: [`treasury/README.md`](treasury/README.md).

| Purpose Group | Count |
|---------------|-------|
| Banking, BA Delivery & Requirements | 15 |
| Architecture, Engineering Decisions & Planning | 22 |
| Implementation, Platform Templates & Code Review | 15 |
| Testing, QA & Validation | 12 |
| Agent Orchestration & Workflow Infrastructure | 17 |
| Research, Debate & Knowledge Synthesis | 15 |
| Security, Governance & Compliance | 7 |
| Observability, Cost & Incident Operations | 9 |
| Code Analysis, Productivity & Publishing | 19 |
| **Total** | **131** |

## Literature Sources

The `literature/` directory contains the complete research corpus — 87 documents totaling ~249K words (clean markdown; the pre-2026-07 "~763K" figure counted raw-HTML capture noise, normalized away on 2026-07-05; the 2026-09-07 refresh added a 19-source frontier-model cohort for the Claude Fable 5.1 / GPT-6 Astra generation) — organized by ecosystem:

| Category | Files | Key Sources |
|----------|-------|-------------|
| **Anthropic / Claude** | 25 | Claude Code Best Practices, Skills Repo, Sub-Agents, Context Engineering, Plugins & Plugin Marketplaces; Fable 5.1 / Fable 5 / Opus 5 prompting guides, Fable 5.1 what's-new and migration guide, effort, thinking, prompt caching, tool search |
| **Codex / Copilot** | 12 | Codex Skills & Agents, GitHub CLI Agent Skills, VS Code Copilot, VS Code Custom Agents, Codex models (GPT-6 Astra), community Codex CLI Astra notes |
| **OpenAI** | 6 | Skill Creator (canonical reference), API Tools, Cookbook, GPT-6 Astra model guidance and announcement |
| **Best Practices** | 6 | Cross-platform authoring guides, SKILL.md pattern analysis, GPT-6 Astra prompting-tips coverage, Fable 5.1 agent-product-design analysis |
| **Awesome Lists** | 2 | ScienceAIX and VoltAgent community aggregations |
| **Open Standard** | 3 | Open Agent Skills specification, MCP specification (2025-11-25), Agent Skills specification (2026-09 re-capture) |
| **Other Platforms** | 14 | Gemini CLI, Cursor, Cline, OpenCode, Windsurf, Spring AI, Strapi, antfu |
| **Research Papers** | 6 | ReAct, Reflexion, Voyager, Toolformer, MemGPT, SWE-agent |
| **Research Vault** | 9 | Internal deep-research syntheses: SKILL.md design, routing, portability, activation, degradation, scoping, strategy, marketplaces, prompt-cruft taxonomy for frontier models |
| **Skill Design Methodology** | 4 | Source index, reading taxonomy, synthesized design principles, literature review |

## Platform Compatibility

Skillify generates portable `SKILL.md` folders compatible with:

- **Claude Code** - `.claude/skills/` or `~/.claude/skills/`
- **OpenAI Codex** - `.agents/skills/`, `~/.agents/skills/`, or `$CODEX_HOME/skills/`
- **GitHub Copilot** - `.github/skills/`, `.agents/skills/`, `~/.copilot/skills/`, or optional custom instructions
- **Gemini / Antigravity** - `.agents/skills/`, `.gemini/skills/`, `~/.gemini/skills/`, or `~/.gemini/antigravity-cli/skills/`

Use the **Adapt** mode to convert between platform conventions.

## License

This repository is a research and documentation project. The literature files retain their original authors' licensing. The `skillify/` meta-skill and tooling are provided as-is for internal use.
