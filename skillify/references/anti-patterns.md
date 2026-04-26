# Skill Creation Anti-Patterns

When building a `SKILL.md` or reviewing an existing one, avoid these common mistakes:

See also: `validation-rubric.md` for scoring and `progressive-disclosure.md` for the architecture behind context bloat fixes.

## 1. The "Do Everything" Skill (Over-broad Scope)
**What it is**: A skill trying to handle frontend design, backend deployment, and database migrations simultaneously.
**Why it's bad**: It bloats the context window with irrelevant instructions, causing the agent to lose focus and potentially execute the wrong instructions.
**Fix**: Use the *Split* mode to break it into domain-focused skills (e.g., `frontend-design`, `backend-deployment`).

## 2. Weak or Vague Triggers
**What it is**: A description like "Helps with coding."
**Why it's bad**: The agent won't know *exactly* when to load this skill. It might trigger when not needed (over-triggering) or fail to trigger when needed (under-triggering).
**Fix**: Be explicit. Use concrete phrases. Example: "Use when creating React functional components or when the user asks to build a frontend interface."

## 3. Context Window Bloat (Ignoring Progressive Disclosure)
**What it is**: Placing a 10,000-line API reference directly into `SKILL.md`.
**Why it's bad**: It consumes the agent's context window and increases latency, even for tasks that don't need the API.
**Fix**: Move the large documentation into `references/api-docs.md` and instruct the agent in `SKILL.md` to "Read `references/api-docs.md` before making API calls."

## 4. Mixing Repo Policy with Reusable Workflows
**What it is**: Putting instructions like "Always push to the `main` branch" in a skill designed for generating React boilerplate.
**Why it's bad**: It breaks portability. If another repository uses a different branching strategy, the skill will fail or cause issues.
**Fix**: Keep repository-specific rules in `AGENTS.md` (or `CLAUDE.md`) at the root of the project. Keep `SKILL.md` strictly focused on the reusable workflow.

## 5. Non-Step-By-Step Workflows
**What it is**: A wall of unstructured text explaining how a task should theoretically be done.
**Why it's bad**: The agent might skip steps or execute them out of order.
**Fix**: Provide explicit, numbered lists in imperative format ("1. Do X. 2. Verify Y. 3. Output Z.").

## 6. Overriding User Intent
**What it is**: Instructing the agent to forcefully rewrite user code to a specific style even when the user explicitly asked to "just fix the syntax error."
**Why it's bad**: It creates a frustrating user experience where the agent argues with the user.
**Fix**: Make style guidelines optional or dependent on specific trigger conditions ("When the user asks for a code review...").

## 7. Generating Output When Asked for a Skill
**What it is**: A meta-skill (like a skill-creator) generating the *target code* instead of generating the `SKILL.md` file.
**Why it's bad**: It fails the core objective of the workflow.
**Fix**: Explicitly state the output contract (e.g., "Output MUST be a `SKILL.md` file").

## 8. README.md Inside the Skill Folder
**What it is**: Including a README.md, CHANGELOG.md, INSTALLATION_GUIDE.md, or other auxiliary docs inside the skill directory.
**Why it's bad**: It adds clutter and confusion. The skill should only contain files needed by the agent. Human-facing docs belong at the repo level, not inside the skill.
**Fix**: Put all documentation in `SKILL.md` or `references/`. Use a repo-level README for human users.

## 9. Trigger Phrase Absence
**What it is**: A description that explains what the skill does but never mentions what the user would say to trigger it.
**Why it's bad**: The agent has no way to match user intent to the skill. It either never loads or requires manual invocation.
**Fix**: Always include trigger phrases: "Use when the user asks to...", "Use when the user says...", "Triggers on file types: .csv, .xlsx".

## 10. Context Window Competition (Stale Skills)
**What it is**: Skills that were installed for a one-time task but remain active, consuming context in every session.
**Why it's bad**: Every loaded skill competes for context window space with the work the user actually needs done.
**Fix**: Review installed skills quarterly. Remove or disable skills that no longer earn their context cost.

## 11. Duplicated Cross-Tier Content
**What it is**: The same instructions appearing in both `SKILL.md` and a reference file, or in both the root `CLAUDE.md` and a skill.
**Why it's bad**: Content drifts apart over time, creating contradictions that confuse the agent.
**Fix**: Information should live in exactly one place. Skills point to references; they don't copy them.

## 12. Hardcoded Platform Assumptions
**What it is**: A skill that assumes it will only run in Claude Code (e.g., referencing Claude-specific frontmatter fields or tool names).
**Why it's bad**: It breaks when used in Codex, Copilot, or other agents that support the open skill standard.
**Fix**: Write platform-agnostic instructions. Use the `compatibility` field to declare supported platforms. See `references/platform-compatibility.md`.
