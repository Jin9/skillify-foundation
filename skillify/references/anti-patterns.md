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

## 5. Wrong Degree of Freedom
**What it is**: Specificity that does not match fragility. Under-specified: a wall of prose for an operation where only one sequence is safe. Over-specified: step-by-step choreography, prohibition walls, or "must" and "never" emphasis for work the agent should judge.
**Why it's bad**: Under-specified fragile work gets skipped or reordered. Over-specified judgment work makes the agent follow the script instead of its better plan, and stacked emphasis over-applies: rigid, hedging output and prohibitions that anchor toward the failure they name. Skills written for earlier model generations are often too prescriptive for current ones and lower output quality.
**Fix**: Match specificity to fragility per `model-generation-fit.md` section 1, and keep every item on its keep list. Keep a prohibition only when it guards a failure seen on the current model or encodes a real policy, and say why. `scripts/cruft_scan.py` flags the common forms.

## 6. Overriding User Intent
**What it is**: A skill instruction that outranks the user: forcing a style rewrite when the user asked to "just fix the syntax error", or any line that makes the agent pause, ask for a confirmation the user did not request, or leave requested work unfinished.
**Why it's bad**: Current models weight skill instructions heavily; when the skill and the user disagree the agent may stop, change direction, or follow the skill's rule instead of the request, and the user experiences an agent that argues or stalls.
**Fix**: State that the user's request sets the scope and wins on conflict (the operating contract in `templates/operating-contract.md`). Reserve stops for destructive or irreversible actions and for genuine changes to the scope the user set. Make style guidance conditional ("When the user asks for a code review..."). For diagnosing and rewriting the offending line, see the "Unrequested pause or divergence" row in `lifecycle-and-iteration.md`.

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

## 12. Hardcoded Platform or Model Assumptions
**What it is**: A skill that assumes one host or one model generation: host-specific frontmatter fields or tool names; pinned model names; workarounds for a retired model's habits; scaffolds the current model has natively ("think step by step", "show your reasoning", "hold all findings until the end", "never use bullets").
**Why it's bad**: Host assumptions break on Codex, Copilot, Gemini, or Antigravity. Model assumptions rot silently at the next release: pinned names degrade, retention crutches waste tokens, "show your reasoning" can trigger a refusal on some models, and narration or formatting suppressors strip output the user wanted.
**Fix**: Write platform-agnostic instructions and declare hosts in `compatibility`. Name model-cost tiers (`small`, `mid`, `frontier`) and effort hints, never model names; a dated "e.g." beside a tier is acceptable only in a data table, never in a rule line (`model-generation-fit.md` section 5). Drop scaffolds for behavior the agent now does unprompted; keep verification steps and exact scripts for fragile operations. Re-audit at every model release (`lifecycle-and-iteration.md`, `model-generation-fit.md`). See `references/platform-compatibility.md`.
