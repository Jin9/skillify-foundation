---
authored_by: chinnawat.w (git author of record; GitHub Jin9)
phase_owner: Gemini 3.1 Pro
pipeline_phase: 1
status: accepted
provenance_dated: 2026-05-19
---

# Authoritative Source Index: Agent Skill Design

**Location:** `skill-design-methodology/research/sources.md`

## Official Documentation & Primary Providers (Trust Tier 1)

Authoritative frameworks and direct platform guidance.

| Source Name | Local Path / Link | Trust Tier | Relevant Rule | Expected Impact on Skillify | Confidence |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Anthropic Building Effective Agents** | `literature/anthropic-claude/anthropic_building_effective_agents.md` | Tier 1 | Explicit context and modularity in agent instruction. | Structures core decomposition principles in `SKILL.md`. | High |
| **Anthropic Effective Context Engineering** | `literature/anthropic-claude/anthropic_effective_context_engineering.md` | Tier 1 | Precision and token-efficiency in prompting. | Drives concise authoring rules in `frontmatter-guide.md` and templates. | High |
| **Claude Code Best Practices** | `literature/anthropic-claude/claude_code_best_practices.md` | Tier 1 | Single-responsibility, clear examples for agent scaffolding. | Dictates structural layout and required sections of a skill template. | High |
| **OpenAI API Tools & Skills** | `literature/openai/openai_api_tools_skills.md` | Tier 1 | Functional boundaries and strict API parameter definitions. | Hardens tool-calling integrations in `mcp-skill-template.md`. | High |
| **OpenAI Codex Best Practices** | `literature/codex-copilot/codex_best_practices.md` | Tier 1 | Descriptive naming and avoiding logic overlap. | Refines anti-pattern definitions in `anti-patterns.md`. | High |
| **GitHub Copilot Agent Skills** | `literature/codex-copilot/github_copilot_agent_skills.md` | Tier 1 | IDE-integrated context boundaries and progressive disclosure. | Influences `progressive-disclosure.md` and Copilot mappings. | High |
| **Anthropic Agent Skills Overview** | `literature/anthropic-claude/anthropic_agent_skills_overview.md` | Tier 1 | Standardizing agent inputs and capability discovery. | Baselines the methodology for capability exposure. | High |

## Open Standards & Core Repositories (Trust Tier 2)

Widely adopted standards, open formats, and official implementation examples.

| Source Name | Local Path / Link | Trust Tier | Relevant Rule | Expected Impact on Skillify | Confidence |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **AGENTS.md Open Format** | `literature/other-platforms/agents_md.md` | Tier 2 | Standardized repository-level agent instruction files. | Validates open compatibility goals in `skillify/SKILL.md`. | High |
| **OpenAI Skill Creator** | `literature/openai/openai_skill_creator.md` | Tier 2 | Meta-prompting and dynamic generation of skill scaffolds. | Directly informs `init_skill.py` and meta-skill logic. | High |
| **AgentSkills.io** | `literature/other-platforms/agentskills_home.md` | Tier 2 | Interoperability of shared toolsets across platforms. | Guides cross-platform compatibility in `platform-compatibility.md`. | Medium |
| **OpenAI Cookbook Skills** | `literature/openai/openai_cookbook_skills.md` | Tier 2 | Practical recipes for multi-step agent workflows. | Enriches `workflow-patterns.md` with proven architectural recipes. | High |
| **Claude Code Subagents** | `literature/anthropic-claude/claude_sub_agents.md` | Tier 2 | Hierarchical task delegation and state handoff. | Influences `workflow-patterns.md` for orchestrator-based skills. | High |

## Community Guidance & Secondary Sources (Trust Tier 3)

Blogs, community repos, and experimental listicles. Used for validation, not foundational rules.

| Source Name | Local Path / Link | Trust Tier | Relevant Rule | Expected Impact on Skillify | Confidence |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Medium: Cheat Codes for Claude Code** | `literature/anthropic-claude/medium_cheat_codes_claude_code.md` | Tier 3 | Developer experience hacks and formatting tricks. | Minor tweaks to `basic-skill-template.md` readability. | Low |
| **Snyk Top Claude Skills** | `literature/anthropic-claude/snyk_top_claude_skills.md` | Tier 3 | Security-first boundary definitions for agent capabilities. | Cross-checks the `security-checklist.md` coverage. | Medium |
| **MGechev Skills Best Practices** | `literature/best-practices/mgechev_skills_best_practices.md` | Tier 3 | Maintainability and iterative testing of skill definitions. | Adds community weight to `lifecycle-and-iteration.md`. | Medium |
| **Awesome Agent Skills (ScienceAIx)** | `literature/awesome-lists/awesome_agent_skills_scienceaix.md` | Tier 3 | Ecosystem mapping and naming conventions. | Broadens `domain-skill-template.md` categories. | Low |
