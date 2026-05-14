# Source Document Synthesis: Agent Skill Creation Principles (Expanded)

This document serves as the definitive architectural input for Opus 4.7 Max Reasoning to design the `skill-creator` meta-skill. It synthesizes the exact structural constraints, trigger mechanics, anti-patterns, and lifecycle requirements extracted from 55 source files across the Anthropic, Codex, Copilot, OpenCode, Gemini, Cursor, Cline, Windsurf, open-standard, and academic agent ecosystems.

## 1. Skill Architecture & Progressive Disclosure
The fundamental design pattern for agent skills is **Progressive Disclosure**, ensuring the agent's context window is strictly optimized.

*   **Metadata Layer (YAML Frontmatter):** The most critical layer. Always in context (~100 words). Contains the `name` and `description`.
    *   *Constraint:* Must be strictly valid YAML bounded by `---`.
    *   *Constraint:* The file must be named exactly `SKILL.md` (case-sensitive).
    *   *Constraint:* All trigger information MUST live here. Placing "When to use this skill" inside the Markdown body is an anti-pattern because the body is only loaded *after* the skill triggers.
*   **Body Layer (SKILL.md):** Loaded only *after* triggering (<5k words). Contains declarative instructions, core workflow, and heuristics.
*   **Reference Files:** When a skill supports multiple variations, frameworks, or options, keep only the core workflow in `SKILL.md`. Move variant-specific patterns into separate reference files.
    *   *Crucial Rule:* Avoid deeply nested references. Keep references exactly one level deep from `SKILL.md`.
*   **Bundled Resources (Scripts/Assets):** Executable code or templates. Used for low-freedom, deterministic actions. Scripts bypass context window limits until their output is evaluated.

## 2. Trigger Mechanics & Description Design
A skill is useless if it fails to trigger or triggers hallucinated states.

*   **Bad Descriptions:** "Helps with projects", "Does things...", or generic task outlines.
*   **Good Descriptions:** Highly specific, explicit value propositions that include exact user trigger phrases.
    *   *Example:* `description: "End-to-end customer onboarding workflow. Use specifically when user says 'onboard new customer', 'set up subscription', or 'create PayFlow account'."`
*   **Relevance:** Must mention specific file types (e.g., `.pdf`, `.tsx`) or context states that warrant activation.

## 3. Strict Anti-Patterns (What NOT to do)
To maintain context economy, agents must avoid human-centric documentation.

*   **NO Extraneous Documentation:** A skill must contain *only* essential files supporting agent functionality. Do NOT create `README.md`, `INSTALLATION_GUIDE.md`, `QUICK_REFERENCE.md`, `CHANGELOG.md`, or setup/testing procedures.
*   **NO Vague Validations:** Do not use phrases like "Make sure to validate things properly."
    *   *Instead:* Use explicit checklists: "CRITICAL: Before calling create_project, verify: 1. Project name is non-empty. 2. Start date is not in the past."
*   **NO Language-based Critical Logic:** If a validation is critical and complex, do not rely on LLM interpretation. "Code is deterministic; language interpretation isn't." Bundle a Python/Bash script to perform the check programmatically.

## 4. Degrees of Freedom
Match the instruction format to the task's required rigidity:
1.  **High Freedom (Text instructions):** Best for creative work, heuristic problem solving, or tasks where multiple approaches are valid.
2.  **Medium Freedom (Pseudocode / scripts with parameters):** Best for structured workflows where some variation is acceptable but a preferred pattern exists.
3.  **Low Freedom (Rigid scripts):** Best for deterministic operations, API integrations, and irreversible actions.

## 5. The Skill Lifecycle & Validation
The `skill-creator` must enforce a standardized development lifecycle:

1.  **Understand via Concrete Examples:** Define exact user trigger phrases and expected system outputs before writing instructions.
2.  **Initialize:** Use boilerplate generation scripts (e.g., `init_skill.py`) to create the base directory structure.
3.  **Draft Reusable Content:** Extract rigid logic into `scripts/`, `references/`, and `assets/`. Test these standalone first.
4.  **Validate Structure:** Run quality checklists and validation scripts (`scripts/quick_validate.py <path>`) to catch invalid YAML formatting, unclosed quotes, or missing fields.
5.  **Iterate:** Run the skill on real traces, identify hallucinations or missing context, and compact the instructions.

## 6. Pre-Flight Checklist for `skill-creator` output
Any generated `SKILL.md` must pass this checklist before deployment:
- [ ] Is the frontmatter strictly valid YAML with `---` delimiters?
- [ ] Are all quotes properly closed in the frontmatter?
- [ ] Is the `description` highly specific and populated with literal trigger phrases?
- [ ] Are all auxiliary human-readable docs (`README.md`, etc.) removed?
- [ ] Are complex validations offloaded to deterministic scripts?
- [ ] Are instructions written in actionable, imperative form?

## 7. Surface Selection: Skill vs Rule vs Agent vs Prompt
The expanded platform corpus shows that not all reusable guidance belongs in `SKILL.md`.

* **SKILL.md:** Use for portable, repeatable workflows with optional scripts, references, templates, or assets.
* **AGENTS.md / GEMINI.md / rules:** Use for durable repo, directory, or team conventions that should influence many tasks.
* **Prompt files / custom commands / workflows:** Use for manual, reusable prompt templates without heavy resources.
* **Custom agents / subagents:** Use for persistent role/persona behavior, tool restrictions, model preferences, or handoffs.
* **Memory:** Use for local facts and transient preferences; do not rely on memory for team-shared durable knowledge.

## 8. Agent Research Patterns To Encode
Academic agent literature adds reusable workflow shapes:

* **ReAct:** Interactive skills should alternate planning, action, observation, and plan revision.
* **Reflexion:** Iterative skills should capture failed attempts, feedback signals, concise lessons, and retry strategy.
* **Voyager:** A skill library should favor narrow, composable skills that can be retrieved and reused in new contexts.
* **Toolformer:** Tool guidance must include when to call the tool, how to form arguments, and how to integrate results.
* **MemGPT:** Progressive disclosure is context management; hot context must stay small and cold context should remain external.
* **SWE-agent:** Scripts and command wrappers should be designed as agent-facing interfaces with concise, structured outputs.

## 9. Cross-Platform Portability Rules
- Prefer `.agents/skills/` when the host supports it because it is increasingly scanned by multiple systems.
- Preserve platform-specific paths only when targeting a single host (`.github/skills`, `.windsurf/skills`, `.claude/skills`).
- Keep portable frontmatter valid first; add host-specific fields only when they are optional and harmless elsewhere.
- For hosts without skills, adapt reusable workflows into rules, command files, or imported context documents.
- Never duplicate large guidance across multiple platform files. Create one source file and reference or adapt it.
