# Skillify — AI Agent Skill Engineering Toolkit

A research-backed toolkit for creating, validating, and maintaining production-grade `SKILL.md` files for AI coding agents. Built from a systematic review of 55 primary and supporting sources across the Anthropic, OpenAI, GitHub Copilot, open-standard, platform-rule, and academic agent ecosystems.

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
├── literature/              # 55 curated source documents
│   ├── anthropic-claude/    # 13 files — Claude Code, sub-agents, skills docs
│   ├── codex-copilot/       #  9 files — Codex agents, GitHub Copilot skills, VS Code agents
│   ├── openai/              #  4 files — OpenAI skill creator, API tools, cookbook
│   ├── best-practices/      #  4 files — Cross-platform skill authoring guides
│   ├── awesome-lists/       #  2 files — Community skill aggregations
│   ├── open-standard/       #  1 file  — Agent Skills specification
│   ├── other-platforms/     # 13 files — Gemini, Cursor, Cline, OpenCode, Windsurf, and others
│   ├── research-papers/     #  6 files — ReAct, Reflexion, Voyager, Toolformer, MemGPT, SWE-agent
│   └── skill-design-methodology/ # 3 files — synthesized source index, taxonomy, principles
│
└── skillify/                # The skill-creator meta-skill
    ├── SKILL.md             # Core skill definition (8 modes)
    ├── references/          # 9 reference guides
    │   ├── anti-patterns.md
    │   ├── frontmatter-guide.md
    │   ├── lifecycle-and-iteration.md
    │   ├── mode-playbooks.md
    │   ├── platform-compatibility.md
    │   ├── progressive-disclosure.md
    │   ├── security-checklist.md
    │   ├── validation-rubric.md
    │   └── workflow-patterns.md
    ├── templates/           # 4 starter templates
    │   ├── basic-skill-template.md
    │   ├── domain-skill-template.md
    │   ├── mcp-skill-template.md
    │   └── audit-report-template.md
    ├── scripts/             # 3 automation scripts
    │   ├── init_skill.py
    │   ├── quick_validate.py
    │   └── check_links.py
    └── examples/            # 3 worked examples
        ├── create-from-scratch.md
        ├── good-description-examples.md
        └── skill-audit-walkthrough.md
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

## Literature Sources

The `literature/` directory contains the complete research corpus — 55 documents totaling ~758K words — organized by ecosystem:

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
| **Skill Design Methodology** | 3 | Source index, reading taxonomy, synthesized design principles |

## Platform Compatibility

Skillify generates portable `SKILL.md` folders compatible with:

- **Claude Code** - `.claude/skills/` or `~/.claude/skills/`
- **OpenAI Codex** - `.agents/skills/`, `~/.agents/skills/`, or `$CODEX_HOME/skills/`
- **GitHub Copilot** - `.github/skills/`, `.agents/skills/`, `~/.copilot/skills/`, or optional custom instructions
- **Gemini / Antigravity** - `.agents/skills/`, `.gemini/skills/`, `~/.gemini/skills/`, or `~/.gemini/antigravity/skills/`

Use the **Adapt** mode to convert between platform conventions.

## License

This repository is a research and documentation project. The literature files retain their original authors' licensing. The `skillify/` meta-skill and tooling are provided as-is for internal use.
