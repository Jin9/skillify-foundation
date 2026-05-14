# Literature Review and Analysis

A cross-corpus synthesis of the 54 source documents in `literature/`. Companion to `agent-skill-design-principles.md` (prescriptive) and `agent-skill-source-index.md` (catalog); this document is the descriptive review — what the corpus actually says, where it converges, where it disagrees, and what it leaves uncovered.

---

## 1. Executive Summary

The corpus is comprehensive across all major agent ecosystems and converges strongly on the core SKILL.md mechanic: a folder containing a `SKILL.md` with `name` + `description` YAML frontmatter, a ≤500-line body, optional `references/` / `scripts/` / `assets/` directories, and three-tier progressive disclosure (metadata always loaded → body on trigger → references on demand).

Divergence concentrates in five places: (1) description length limits, (2) custom-agent vs. subagent terminology, (3) path conventions, (4) supply-chain and provenance, and (5) ecosystem-specific frontmatter extensions. Six research papers ground the design moves but are currently under-utilized by the synthesized methodology — each gets only a single-sentence treatment.

The current `agent-skill-design-principles.md` captures roughly 80% of what the corpus supports. The remaining 20% — research-grounded design moves, supply-chain rules, the custom-agent surface, skill lifecycle, conditional load hints, and a full frontmatter field inventory — is well-documented in the literature and ready to be folded in.

---

## 2. Corpus Inventory

| Category | Files | Coverage |
|---|---:|---|
| `anthropic-claude/` | 13 | Official Claude Code & Skills docs, Anthropic engineering blog, community Claude Code skill notes, Groff three-tier pattern, Snyk skill catalog |
| `codex-copilot/` | 9 | Codex agents/skills/best-practices, GitHub Copilot agent skills, GitHub CLI `gh skill`, VS Code Copilot, VS Code custom agents |
| `openai/` | 4 | OpenAI skill-creator canonical example, skills repo, API tools/skills guide, cookbook |
| `best-practices/` | 4 | agentskills.io, mgechev best practices, two Medium deep-dives on the SKILL.md pattern |
| `awesome-lists/` | 2 | VoltAgent and ScienceAIX curated catalogs |
| `open-standard/` | 1 | Open Agent Skills specification (openagentskills.dev) |
| `other-platforms/` | 13 | agents.md, Gemini CLI (context + commands), Cursor rules, Cline rules, OpenCode, Windsurf Cascade, Spring AI, Strapi, antfu collection, sohamkamani, DataCamp, agentskills.io home |
| `research-papers/` | 6 | ReAct, Reflexion, Voyager, Toolformer, MemGPT, SWE-agent |
| `skill-design-methodology/` | 3 | Source index, reading-task taxonomy, design principles (already synthesized) |
| **Total** | **54** | ~33,500 lines / ~750k words |

**Extraction caveat:** roughly one in four files is a saved HTML page capture rather than clean markdown (notably most `codex-copilot/` Codex docs, several `vscode_*` files, the OpenAI cookbook, and some `other-platforms/` blog captures). Their content is partially recoverable but biases automated frequency analysis with HTML noise. A normalization pass is the highest-leverage cleanup task.

The README at the project root and `agent-skill-source-index.md` both reference "55 sources" — the discrepancy is `claude_guide.pdf` (a binary file with the markdown companion `claude_guide.md` already counted).

---

## 3. Per-Category Findings

### 3.1 Anthropic / Claude

Defines the immutable core mechanics. The `anthropic_agent_skills_overview.md` and `anthropic_skill_best_practices.md` are the authoritative spec; `claude_guide.md` is the long-form practitioner manual; `groff_implementing_claude_md.md` and `community_claude_skills.md` are the strongest community sources.

**Authoritative rules:**
- Exactly `name` (kebab-case) + `description` required; `description` ≤1024 chars
- File must be named `SKILL.md` (case-sensitive); folder name must equal `name`
- No XML angle brackets (`<>`) anywhere in frontmatter — security restriction (safe YAML parsing only)
- Reserved namespaces: skill names prefixed `claude` or `anthropic` are forbidden
- No `README.md`, `INSTALL.md`, `CHANGELOG.md`, `QUICK_REFERENCE.md` inside the skill folder
- Three-tier loading: metadata (always) → body (on trigger) → linked resources (on demand)

**Community Claude-Code-specific frontmatter (15 fields documented):**
`name`, `description`, `when_to_use` (description + when_to_use combined ≤1536 chars), `allowed-tools`, `model`, `effort` (low/medium/high/xhigh/max), `context` (`fork` creates an isolated subagent), `agent`, `hooks`, `paths` (glob), `shell`, `argument-hint`, `arguments`, `disable-model-invocation`, `user-invocable`, `license`.

**Internal contradictions inside Anthropic sources:**

| Question | Anthropic official | Community Claude Code | Groff (pragmatic) |
|---|---|---|---|
| `description` length | ≤1024 chars | ≤1536 combined with `when_to_use` | (defers to spec) |
| SKILL.md size | ≤300 lines recommended | (silent) | ≤150 lines for most skills |
| Root CLAUDE.md size | (silent) | (silent) | ≤100 lines; HumanLayer cites ≤60 |
| Naming style | mix of gerund and noun in examples | mix | mix; "use whichever is clearest" |

**Most quotable anti-patterns** (from Groff and the official docs):
- Don't auto-generate CLAUDE.md — `/init` produces generic output that wastes the highest-leverage file
- Don't use CLAUDE.md as a linter — use real linters (Ruff, ESLint, Prettier)
- Don't duplicate content across tiers
- "Every line in [the root file] competes for attention. If it does not apply to literally every task, push it down."
- "A single CI workflow file does not need a github-actions skill… Skill count should match actual decision surface, not aspiration."

### 3.2 Codex / Copilot

Adds the **supply-chain dimension** absent from Anthropic. `github_copilot_add_skills.md` is the most extractable file; the OpenAI Codex docs are mostly HTML-wrapped and need re-fetching.

**Key contributions:**
- `gh skill preview <URL>` (inspect before install), `gh skill install ... --host ... --scope ...`, `gh skill publish` (validates against Agent Skills spec)
- Provenance fields captured on install: `source-repo`, `source-ref`, `source-tree-sha`, `pinned: true` (pinned skills skip auto-update)
- Path standardization: project `.github/skills/`, `.claude/skills/`, `.agents/skills/`; personal `~/.copilot/skills/`, `~/.claude/skills/`, `~/.agents/skills/`
- Frontmatter: `name` (lowercase-hyphen), `description`, optional `license`
- **Silent failure** is documented: invalid skill names fail to load without error messages

**The custom-agent surface (VS Code/Copilot, unique to this ecosystem):**

VS Code distinguishes four reusable artifacts where Anthropic has three:

| Artifact | Purpose | Frontmatter | Storage |
|---|---|---|---|
| Custom agent | Persistent persona + tool permissions + model preference | YAML | `.vscode/copilot/agents/` or extension contribution |
| Skill | Repeatable workflow with resources | YAML | `.github/skills/`, `.agents/skills/` |
| Prompt file | Static reusable prompt (no resources) | Minimal | `.github/prompts/` |
| Instructions | Repo conventions | Markdown / no frontmatter | `AGENTS.md` |

The VS Code decision rule: *"Who should the agent be?" → custom agent. "How to execute this repeatable task?" → skill. "How does this repo work?" → instructions.*

### 3.3 OpenAI

`openai_skill_creator.md` is the canonical example. The other three OpenAI files are partly HTML-wrapped.

**Canonical frontmatter (strict):**
```yaml
---
name: skill-creator
description: Guide for creating effective skills. This skill should be used when users want to create a new skill (or update an existing skill) that extends Codex's capabilities with specialized knowledge, workflows, or tool integrations.
metadata:
  short-description: Create or update a skill
---
```

**Canonical folder shape:**
```
skill-name/
├── SKILL.md          # required
├── agents/openai.yaml  # recommended — UI metadata
├── scripts/           # executable code (Python/Bash/...)
├── references/        # documentation for context-on-demand
└── assets/            # templates, icons, fonts, boilerplate
```

**Distinctive contributions:**
- Naming: lowercase letters + digits + hyphens, ≤64 chars, **verb-led** preferred (`plan-mode` not `Planner`)
- Body in imperative / infinitive form
- ≤500 lines body; split into references if approaching
- Three degrees of freedom (high / medium / low) → text instructions / parameterized scripts / rigid scripts
- "Skill as versioned bundle": skills are uploadable, attachable to agent instances via API
- Repository tiers: `.system/` (auto-included), `.curated/` (user-installable), `.experimental/` (opt-in) — a governance pattern worth borrowing

### 3.4 Best-Practices, Awesome-Lists, Open-Standard

**Open Agent Skills specification** is the normative cross-platform baseline:

| Constraint | Value |
|---|---|
| `name` | 1–64 chars, lowercase alphanumeric + hyphens, no leading/trailing/consecutive hyphens, must equal parent directory |
| `description` | 1–1024 chars, combines capability + activation |
| `SKILL.md` body | ≤500 lines (recommended) |
| Optional fields | `license`, `compatibility`, `metadata`, `allowed-tools` (experimental) |
| Reference depth | One level deep from `SKILL.md` |

**Awesome-list patterns** (common categories observed across community catalogs): code quality / review, commit + PR workflows, doc generation, testing + coverage, environment setup, deployment, AI-specific (refactor, write-tests, prompt-engineer).

**Most common anti-patterns observed in community skills:**
- Skills that duplicate linter configs
- Hardcoded shell commands without OS/shell adapter
- Descriptions too terse to trigger auto-load
- ≥200 lines of procedural steps that should be agent guides
- Skills with no error handling or rollback guidance

### 3.5 Other Platforms

Surface taxonomy is the key contribution. **AGENTS.md is the universal portable baseline** — first-class support in Cursor, Cline, OpenCode, Windsurf; fallback in Gemini CLI; community-tracked across 60k+ projects.

**Folder convention table (quoted paths):**

| Platform | Workspace | Global | Notes |
|---|---|---|---|
| Anthropic Claude Code | `.claude/skills/`, `.claude/CLAUDE.md`, `AGENTS.md` | `~/.claude/skills/` | |
| OpenAI Codex | `.agents/skills/`, `AGENTS.md` | `$CODEX_HOME/skills/`, `~/.agents/skills/` | |
| GitHub Copilot | `.github/skills/`, `.agents/skills/` | `~/.copilot/skills/` | Provenance via `gh skill` |
| Gemini CLI | `.gemini/GEMINI.md`, `AGENTS.md`, `.gemini/commands/` | `~/.gemini/GEMINI.md`, `~/.gemini/commands/` | Hierarchical merge; `@file.md` imports |
| Cursor | `.cursor/rules/`, `AGENTS.md` | `~/.cursor/` | `.cursorrules` legacy/deprecated |
| Cline | `.clinerules/`, `AGENTS.md` | `~/Documents/Cline/Rules/` | Auto-detects multiple rule formats; toggleable |
| OpenCode | `AGENTS.md`, `opencode.json` | `~/.config/opencode/AGENTS.md` | Claude fallback supported |
| Windsurf | `.windsurf/skills/`, `.windsurf/rules/`, `AGENTS.md` | `~/.codeium/windsurf/skills/`, `~/.agents/skills/` | Four activation modes |

**Path-conditional activation** (unique to rule-based platforms):

Cline (`paths:` YAML frontmatter):
```yaml
---
paths:
  - '**/*.test.ts'
  - tests/**
---
```

Windsurf rule modes: `always_on` / `model_decision` / `glob` / `manual` — the most explicit activation-mode taxonomy in the corpus.

**LLM-agnostic framings** (Spring AI, Strapi): skills as tool registries with function-calling signatures, decoupled from any specific model. Implies skills should ultimately move toward input/output contracts not prose.

### 3.6 Research Papers

The current methodology gives each paper one sentence. This is materially thin given the empirical claims available. Quantified findings worth citing:

| Paper | Headline finding |
|---|---|
| ReAct | HotpotQA 56% vs. CoT 37%; ALFWorld 71% vs. imitation 34%; WebShop 59% vs. CoT 17% |
| Reflexion | HumanEval 91% vs. baseline 76%; ablation: last-3 reflections > all reflections |
| Voyager | 4.5× more unique items, 2.3× farther travel, 15.3× faster milestones with skill library; transfer across worlds |
| Toolformer | Self-supervised tool use generalizes zero-shot to downstream tasks; structured tool outputs critical |
| MemGPT | Coherent QA over 64k-token documents in a 4k context window via tiered memory |
| SWE-agent | 12% SWE-bench solve rate; 5k-token output budget per action; agent-friendly interfaces beat developer-friendly |

The eight design moves these papers support are listed in §8 below.

---

## 4. Cross-Cutting Consensus

Where the entire corpus agrees:

| Dimension | Consensus |
|---|---|
| Required frontmatter | `name` + `description` only |
| Naming | Lowercase kebab-case, ≤64 chars, folder name = `name` field |
| Description | Encodes BOTH capability AND when-to-use; ≤1024 chars (or ≤1536 with Claude-Code `when_to_use`) |
| Body length | ≤500 lines / ≤5,000 tokens; production sweet spot 150–300 lines |
| Three-tier disclosure | Metadata (always) → body (on trigger) → references (on demand) |
| Reference depth | One semantic level deep from `SKILL.md` |
| Forbidden inside the folder | `README.md`, `INSTALL.md`, `CHANGELOG.md`, `QUICK_REFERENCE.md` |
| Forbidden in frontmatter | XML angle brackets (`<>`) |
| Portable path | `.agents/skills/` for cross-host; platform-specific only for single-host targeting |
| Validation philosophy | Critical + deterministic → script; heuristic → prose; complex non-critical → pseudocode |
| Four-surface model | SKILL.md / subagent (or custom agent) / AGENTS.md / memory — plus prompts/commands as a fifth manual surface |

---

## 5. Contradictions and Resolutions

| # | Contradiction | Resolution |
|---|---|---|
| 1 | Methodology says triggers MUST be in frontmatter; no source defines a `triggers:` field — they live in `description:` | State explicitly: triggers are embedded inside `description`, never in the body. Drop the implication of a separate field. |
| 2 | Description granularity: 2-sentence example with multiple trigger phrases vs. "keep it short" community guidance | Target 1–2 sentences, ≤1024 chars; if longer, move detail to body |
| 3 | Body length: ≤500 (spec) vs. ≤300 (Anthropic) vs. ≤150 (Groff pragmatic) | All three valid — 500 is the spec ceiling, 150 is the production sweet spot. Use 300 as default authoring target. |
| 4 | Path portability: `.agents/skills/` portable vs. platform defaults | Both are legitimate. Default to `.agents/skills/` for portability; add platform-specific symlinks/copies when targeting a host's native discovery. |
| 5 | Custom agent vs. subagent vs. skill — VS Code adds a fourth surface absent from Anthropic vocabulary | Document the four-surface model (skill / subagent or custom agent / AGENTS.md / memory) and the decision rule for each. |
| 6 | Reference subdirectory nesting: methodology says one level; OpenAI examples show topic subfolders | One *semantic* level (`references/<topic>/`), no further. |
| 7 | Validation: methodology mandates scripts; community skills often use prose | Methodology position is correct (deterministic > LLM interpretation for critical logic); note it as aspirational. |
| 8 | Supply chain: GitHub builds provenance/signing/pinning; other ecosystems are silent | Add a dedicated supply-chain section to the methodology covering `gh skill`, version pinning, source-ref tracking, license declaration, pre-install inspection. |

---

## 6. Canonical Frontmatter Field Inventory

The methodology currently lists no specific optional fields. This table consolidates the field inventory from across the corpus, with portability tiers.

| Field | Tier | Sources | Notes |
|---|---|---|---|
| `name` | **Required** | All | Lowercase kebab-case, ≤64 chars, must equal parent directory |
| `description` | **Required** | All | ≤1024 chars; encodes capability + trigger |
| `when_to_use` | Host-specific (Claude Code) | community_claude_skills | Concatenates to description; combined ≤1536 chars |
| `license` | Portable optional | Open Spec, OpenAI, GitHub | SPDX identifier preferred |
| `compatibility` | Portable optional | Open Spec | Host-compat declaration |
| `metadata` | Portable optional | Open Spec, OpenAI | Key-value map; `short-description` common |
| `allowed-tools` | Host-specific (Claude, Copilot) | community_claude_skills, Copilot | Tool allowlist |
| `model` | Host-specific (Claude) | community_claude_skills | Model preference hint |
| `effort` | Host-specific (Claude Code) | community_claude_skills | `low`/`medium`/`high`/`xhigh`/`max` |
| `context` | Host-specific (Claude) | community_claude_skills | `fork` creates isolated subagent |
| `agent` | Host-specific (Claude) | community_claude_skills | Subagent type when `context: fork` |
| `hooks` | Host-specific (Claude Code) | community_claude_skills | Lifecycle hooks |
| `paths` | Host-specific (Claude, Cline) | community_claude_skills, cline_rules | Glob activation |
| `shell` | Host-specific (Claude Code) | community_claude_skills | Shell selection |
| `argument-hint` | Host-specific (Claude, VS Code) | community_claude_skills, vscode_custom_agents | Argument schema hint |
| `arguments` | Host-specific (Claude Code) | community_claude_skills | Argument list |
| `disable-model-invocation` | Host-specific (Claude Code) | community_claude_skills | Manual-only flag |
| `user-invocable` | Host-specific (Claude Code, VS Code) | community_claude_skills, vscode_custom_agents | UI/slash visibility |
| `activation_mode` | Host-specific (Windsurf) | windsurf_cascade | `always_on`/`model_decision`/`glob`/`manual` |
| `glob` | Host-specific (Windsurf, Cline) | windsurf_cascade, cline_rules | Path pattern |
| `version` | Provenance optional (GitHub) | github_copilot_add_skills | Semver; used by `gh skill` pinning |
| `pinned` | Provenance (GitHub-managed) | github_copilot_add_skills | Set on install; skips auto-update |
| `source-repo` | Provenance (GitHub-managed) | github_copilot_add_skills | Captured on install |
| `source-ref` | Provenance (GitHub-managed) | github_copilot_add_skills | Captured on install |
| `source-tree-sha` | Provenance (GitHub-managed) | github_copilot_add_skills | Captured on install |
| `category` | Portable optional | Open Spec, OpenAI examples | Discoverability tag |
| `tags` / `keywords` | Portable optional | Open Spec, antfu | Discoverability tags |

**Rule of thumb:** validators should warn on unknown fields but never fail — frontmatter must be valid YAML and host-specific fields should be harmless when parsed by other platforms.

---

## 7. Numeric Limits Catalog

| Limit | Value | Source(s) | Confidence |
|---|---|---|---|
| `name` length | 1–64 chars | Open Spec | Hard rule |
| `description` length | 1–1024 chars | Open Spec, Anthropic | Hard rule |
| `description` + `when_to_use` | ≤1536 chars combined | Claude Code community notes | Claude-Code-specific |
| `SKILL.md` body | ≤500 lines | Open Spec, agentskills.io | Soft rule |
| `SKILL.md` body | ≤300 lines | Anthropic recommendation | Tighter recommendation |
| `SKILL.md` body | ≤150 lines | Groff pragmatic | Production sweet spot |
| Total body words | ≤5,000 | Methodology | Soft target |
| Frontmatter words | ~100 | Methodology | "Always in context" target |
| Reference depth | 1 semantic level | Methodology, Open Spec | Hard rule |
| Root `CLAUDE.md` | ≤100 lines | Groff | Pragmatic |
| Root `CLAUDE.md` | ≤60 lines | HumanLayer | Tightest |
| Script output | ≤50 lines per call | SWE-agent implied | Research-grounded |
| Reflection memory buffer | Last 3 attempts | Reflexion | Research-grounded |

---

## 8. Research-Grounded Design Moves

Eight specific moves the corpus supports — to be folded into the methodology's §8.

1. **Encode immediate observation loops (ReAct).** For any skill that calls external APIs, executes code, or queries resources: include a section listing the exact observations to collect at each step, and a retry rule when they are missing. *Empirical basis: ReAct 56% vs. CoT 37% on HotpotQA.*

2. **Bundle failure-reflection rules (Reflexion).** For iterative skills: declare a max-attempt count (default 3), a one-line reflection format, and an escalation rule. Cap reflection memory at the last 3 attempts. *Empirical basis: Reflexion 91% vs. 76% on HumanEval; last-3 ablation.*

3. **Design for narrow composability (Voyager).** Each skill solves one decision point. Document explicit preconditions and postconditions so skills can chain. Maintain a retrieval-grade `description` distinct from the body. *Empirical basis: Voyager 4.5× discovery.*

4. **Prescribe tool-use decision rules (Toolformer).** For each tool a skill uses: write a *when* trigger, a *how* argument schema, and an *integration* rule. Vague suggestions don't generalize; sharp rules do. *Empirical basis: Toolformer self-supervised tool-use generalization.*

5. **Specify structured output formats for bundled scripts (Toolformer + SWE-agent).** Scripts emit JSON or named key-value pairs, not prose. Document the success-case schema and the failure-case schema separately. *Empirical basis: SWE-agent interface-design findings.*

6. **Enforce script output-conciseness budgets (SWE-agent).** Default budget: ≤50 lines per script call. Detailed output goes to an artifact file referenced by path. *Empirical basis: SWE-agent 5k-token-per-action budget; verbose output degrades success.*

7. **Structure progressive disclosure with explicit load guards (MemGPT).** Each reference file declares its load condition (e.g., "Load only when initial attempt fails"). Prevents both bloat and under-loading. *Empirical basis: MemGPT selective-loading ablation.*

8. **Include precondition-validation scripts (ReAct + Toolformer + SWE-agent).** Before the main workflow, run a deterministic validator that checks preconditions and returns a structured diagnostic on failure. *Empirical basis: all three papers rely on agents being able to diagnose state.*

---

## 9. Gaps in the Current Methodology

The current `agent-skill-design-principles.md` is solid (82 lines, 9 sections) but the literature now supports filling nine specific gaps:

1. **Frontmatter field inventory** — currently absent; §6 above is the proposed reference table.
2. **Conditional load hints for references** — methodology mentions progressive disclosure but provides no template syntax. MemGPT-grounded.
3. **Retry / reflection policy** — methodology says iterative skills should capture failures but doesn't quantify a max-attempts budget. Reflexion-grounded.
4. **Composability contracts** — methodology calls for narrow skills without requiring documented preconditions/postconditions. Voyager-grounded.
5. **Script output contracts** — methodology says scripts should be agent-facing without mandating structured JSON, error codes, or conciseness budgets. Toolformer + SWE-agent-grounded.
6. **Custom agent surface** — §7 (Surface Selection) lists only SKILL.md / AGENTS.md / prompts / subagents / memory. Missing: VS Code/Copilot's custom-agent surface.
7. **Path-conditional rules** — Cline's `paths:` and Windsurf's `glob` mode are alternative scoping mechanisms absent from the methodology.
8. **Supply chain & provenance** — `gh skill` pinning, source-ref tracking, pre-install inspection, license declaration. Currently uncovered.
9. **Skill lifecycle / pruning** — Voyager and Groff both emphasize active retirement of stale or superseded skills. Methodology covers iteration but not retirement.

---

## 10. Prioritized Recommendations

1. **Resolve the trigger-placement wording** in `agent-skill-design-principles.md` line 11. State explicitly: triggers are embedded inside the `description` field; placing them in the body is the anti-pattern (because the body loads lazily). Avoid implying a separate `triggers:` field exists.

2. **Append four new sections (10–13) to the design-principles document** covering: (10) the frontmatter field inventory of §6, (11) supply-chain and provenance, (12) the four-surface decision model including custom agents, (13) skill lifecycle including retirement.

3. **Replace the six one-liners in §8** with the eight research-grounded design moves from §8 above.

4. **Normalize the HTML-wrapped source files.** Roughly a dozen `.md` files contain saved HTML page captures rather than clean markdown — they extract poorly and inflate substring-based frequency counts. A `scripts/normalize_html.py` pass converting them to clean markdown would materially improve the corpus.

5. **Reconcile the "55 sources" claim** in the project README and `agent-skill-source-index.md` with the actual 54 files (the 55th is `claude_guide.pdf`, a binary companion to `claude_guide.md`). Either count it explicitly or update the references.

6. **Document the description-length contradiction explicitly** (1024 chars baseline, 1536 with `when_to_use` extension on Claude Code) so the validator can enforce the right limit per target host.

7. **Treat all term-frequency analysis as directional only** until the HTML normalization pass — substring matching across HTML noise inflates counts (e.g., the apparent 1,642 mentions of "subagent" is implausible at 54 files).

---

## 11. Appendix: Source Files by Category

### anthropic-claude/
- `anthropic_agent_skills_overview.md`
- `anthropic_skill_best_practices.md`
- `anthropic_skills_repo.md`
- `anthropic_building_effective_agents.md`
- `anthropic_effective_context_engineering.md`
- `claude_code_best_practices.md`
- `claude_sub_agents.md`
- `claude_guide.md` (+ `claude_guide.pdf` binary)
- `community_claude_skills.md`
- `groff_implementing_claude_md.md`
- `medium_cheat_codes_claude_code.md`
- `snyk_top_claude_skills.md`

### codex-copilot/
- `codex_skills.md`, `codex_agents.md`, `codex_best_practices.md`, `codex_customization.md`
- `github_copilot_agent_skills.md`, `github_copilot_add_skills.md`
- `github_cli_agent_skills.md`
- `vscode_copilot_agent_skills.md`, `vscode_custom_agents.md`

### openai/
- `openai_skill_creator.md` (canonical reference)
- `openai_skills_repo.md`, `openai_api_tools_skills.md`, `openai_cookbook_skills.md`

### best-practices/
- `agentskills_best_practices.md`, `mgechev_skills_best_practices.md`
- `medium_skill_md_pattern.md`, `medium_deep_dive_skill_md.md`

### awesome-lists/
- `awesome_agent_skills_voltagent.md`, `awesome_agent_skills_scienceaix.md`

### open-standard/
- `open_agent_skills_specification.md`

### other-platforms/
- `agents_md.md`, `agentskills_home.md`
- `gemini_cli_context_files.md`, `gemini_cli_custom_commands.md`
- `cursor_rules.md`, `cline_rules.md`, `opencode_rules.md`
- `windsurf_cascade_skills_rules_agents.md`
- `spring_ai_agent_skills.md`, `strapi_agent_skills.md`
- `antfu_skills.md`, `sohamkamani_ai_agents.md`, `datacamp_top_agent_skills.md`

### research-papers/
- `react_reasoning_acting.md`, `reflexion_verbal_reinforcement.md`
- `voyager_skill_library.md`, `toolformer_tool_use.md`
- `memgpt_context_management.md`, `swe_agent_computer_interfaces.md`

### skill-design-methodology/
- `agent-skill-source-index.md` (catalog)
- `agent-reading-task-taxonomy.md` (reading framework)
- `agent-skill-design-principles.md` (prescriptive synthesis)
- `literature-review-and-analysis.md` (this document — descriptive review)
