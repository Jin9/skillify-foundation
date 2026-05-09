# Research Notes

## Phase 0 Baseline Inventory - Codex/GPT Rerun

Date: 2026-05-09

### Target

- Repository source skill: `skillify/`
- Installed skill target: `/Users/IF640063/.codex/skills/skillify/`
- Primary file: `SKILL.md`
- Current source and installed `SKILL.md` are identical.
- Current `SKILL.md` line count: 179.

### Current Skill Boundary

`skillify` is a meta-skill for creating, refactoring, reviewing, auditing, compressing, splitting, merging, and adapting reusable agent skills. It must produce skill folders, `SKILL.md` files, references, templates, scripts, assets, examples, or review/audit reports. It must not generate the downstream artifact a requested skill would later produce. It must not edit repo-policy files unless the user explicitly asks for platform-policy adaptation.

### Existing Support Files

- References: `anti-patterns.md`, `frontmatter-guide.md`, `lifecycle-and-iteration.md`, `mode-playbooks.md`, `platform-compatibility.md`, `progressive-disclosure.md`, `security-checklist.md`, `validation-rubric.md`, `workflow-patterns.md`
- Templates: `audit-report-template.md`, `basic-skill-template.md`, `domain-skill-template.md`, `mcp-skill-template.md`
- Scripts: `check_links.py`, `init_skill.py`, `quick_validate.py`
- Platforms: `agents.md`, `antigravity.md`, `claude.md`, `codex.md`, `copilot-instructions.md`, `copilot.md`, `deployment-guide.md`, `gemini.md`, `install.sh`
- Examples: `create-from-scratch.md`, `good-description-examples.md`, `skill-audit-walkthrough.md`

### Stable Constraints

- Preserve unrelated user edits. Current dirty files before this rerun: `.gitignore`, `skillify/SKILL.md`, `skillify/platforms/deployment-guide.md`, `skillify/platforms/install.sh`, plus untracked `file.txt`.
- Do not create `README.md`, `CHANGELOG.md`, `INSTALLATION_GUIDE.md`, `QUICK_REFERENCE.md`, `CONTRIBUTING.md`, or other human-facing docs inside the skill folder.
- Keep permanent workflow artifacts inside `skill-design-methodology/research/` or the target skill folder only.
- Final target is `/Users/IF640063/.codex/skills/skillify/SKILL.md`, with optional focused references only when they reduce load.
- Preserve existing scripts and templates unless a model phase identifies a concrete defect.

### Baseline Assessment

Strengths to preserve:

- Clear meta-skill boundary and explicit negative triggers.
- Concrete mode routing for Create, Refactor, Review, Audit, Compress, Split, Merge, and Adapt.
- Create mode is self-contained in `SKILL.md`; non-create modes point to `references/mode-playbooks.md`.
- Deterministic validation scripts exist for frontmatter/structure and local links.
- Reference files are one-level deep and mostly have clear ownership.

Risks for model phases:

- The skill is broad by mode count, so later phases should challenge whether mode routing is enough to prevent over-triggering.
- Completion-report and validation instructions are inline; later phases should decide whether that is useful immediate context or reference bloat.
- Platform compatibility is broad; later phases should avoid vendor-specific assumptions in the core body.

## Phase 2: Core Skill Design Rules

### 1. Trigger Design

**Source-Backed Rules:**
- **Explicit Constraints:** Avoid vague descriptions ("Helps with coding"). Use concrete user intents and conditions (e.g., "Use when creating React functional components" or "Triggers on file types: .csv, .xlsx"). *(Source: Local `anti-patterns.md`, `validation-rubric.md`)*
- **Verbatim Triggers:** Explicitly document the trigger phrases the user will say (e.g., "Use when the user asks to format CSV data"). *(Source: Local `anti-patterns.md`)*
- **Implicit Activation:** Skills trigger implicitly by matching the `description` field in the frontmatter against the user's prompt. Ensure the description is tightly scoped to a single domain. *(Source: Claude Code Best Practices, Local `platform-compatibility.md`)*

**Synthesis:**
- **Negative Triggers:** Add negative triggers to prevent over-activation (e.g., "Do NOT use for backend business logic").
- **Frontmatter Budget:** Keep the description under 1024 characters to respect context limits while retaining high-signal trigger phrases.

### 2. Progressive Disclosure

**Source-Backed Rules:**
- **The 500/5000 Rule:** Keep the `SKILL.md` body under 500 lines or 5,000 tokens. If it exceeds this, split content. *(Source: Local `progressive-disclosure.md`)*
- **Omit Baseline Knowledge:** Do not explain basic concepts (e.g., what a PDF is). Only include project-specific conventions, edge cases, and non-obvious tools. *(Source: Claude Code Best Practices)*
- **On-Demand Loading:** Move deep reference material (API docs, large templates, exhaustive schemas) into a `references/` or `assets/` directory. Use 1-line pointers in `SKILL.md` specifying exactly *when* the agent should load them (e.g., "Read `references/api-errors.md` if the API returns a non-200 status code"). *(Source: Local `progressive-disclosure.md`, Claude Code Best Practices)*

**Synthesis:**
- **No Deep Nesting:** Keep all reference files exactly one directory level deep from `SKILL.md` to guarantee simple path resolution by the agent.

### 3. Validation

**Source-Backed Rules:**
- **Strict Output Contracts:** Replace vague output requests ("Produce good code") with exact deliverables ("Output a `jest` test file with 100% coverage"). *(Source: Local `validation-rubric.md`)*
- **Validation Loops:** Provide explicit criteria for success. Instruct the agent to run a validator (e.g., a linter, a test suite, or a script), review the errors, fix issues, and repeat until validation passes. *(Source: Claude Code Best Practices)*
- **Plan-Validate-Execute:** For destructive or batch operations, instruct the agent to generate a plan, validate it against a source of truth, and wait for success before executing. *(Source: Claude Code Best Practices, Local `workflow-patterns.md`)*

**Synthesis:**
- **Iteration Caps:** Always include an iteration cap (e.g., "Stop after 3 failed validation attempts and report the blockers") to prevent infinite, token-draining failure loops.

### 4. Safety

**Source-Backed Rules:**
- **Restrict Exfiltration & Destruction:** Explicitly scan for and forbid unauthorized external network requests (`curl`, `wget`), secrets harvesting, and destructive commands (like `git push --force` or unconfirmed file overwrites). *(Source: Local `security-checklist.md`)*
- **Tool Scoping:** Use the `allowed-tools` field to limit tool capabilities and the `paths` field to limit directory access. *(Source: Local `security-checklist.md`)*
- **Frontmatter Hygiene:** Never include XML angle brackets (`<`, `>`) in the YAML frontmatter to prevent prompt injection. Do not use reserved vendor terms (`claude`, `anthropic`) in skill names. *(Source: Local `security-checklist.md`)*

**Synthesis:**
- **Approval Gates:** For sensitive state changes (e.g., deploying, dropping tables), explicitly instruct the agent to halt and request human approval via the CLI prompt before proceeding.

### 5. Portability

**Source-Backed Rules:**
- **Separate Policy from Workflow:** Do not put repository-specific rules ("Always push to the `main` branch") inside a reusable skill. Put repo policies in root `AGENTS.md`, `CLAUDE.md`, or `GEMINI.md` files. *(Source: Local `anti-patterns.md`, `platform-compatibility.md`)*
- **Platform-Agnostic Language:** Avoid hardcoded vendor assumptions (e.g., use "Run the command" instead of "Ask Claude to run"). Declare supported agent platforms using the `compatibility` field. *(Source: Local `anti-patterns.md`, `platform-compatibility.md`)*
- **Skill Bundling:** A skill must be portable as a complete directory. Referenced scripts, templates, and platforms must travel inside the skill folder. *(Source: Local `platform-compatibility.md`)*

**Synthesis:**
- **Standardized Fallbacks:** If a skill requires an external binary (e.g., `jq`, `gh`), provide an alternative path or instruct the agent to use standard language scripts (e.g., Python or Node.js) if the binary is missing.

### 6. Compaction

**Source-Backed Rules:**
- **Favor Procedures over Declarations:** Use numbered, imperative step-by-step instructions ("1. Do X. 2. Verify Y.") instead of long, narrative prose. *(Source: Local `anti-patterns.md`, `validation-rubric.md`)*
- **Provide Defaults, Not Menus:** Do not list exhaustive options. Pick a default tool/approach and state it clearly to save tokens. *(Source: Claude Code Best Practices)*
- **Bundle Logic into Scripts:** If an agent consistently struggles or reinvents logic for a task, extract that logic into a reproducible script in the `scripts/` folder and instruct the agent to execute it. *(Source: Claude Code Best Practices)*
- **Remove Human Docs:** Delete `README.md`, `CHANGELOG.md`, and `CONTRIBUTING.md` from inside the skill folder to prevent context window clutter. *(Source: Local `anti-patterns.md`)*

**Synthesis:**
- **Deduplicate Cross-Tier Content:** Keep instructions in exactly one place. If a rule exists in the root rules file, do not repeat it in the skill; point to it or omit it.
