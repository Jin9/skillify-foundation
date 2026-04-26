---
name: skillify
description: >
  Create, refactor, review, audit, compress, split, merge, or adapt SKILL.md files
  and supporting skill assets for AI agent workflows. Use when the user asks to
  "create a skill", "write a SKILL.md", "design an agent skill", "build a skill
  for X", "refactor a skill", "review my skill", "audit a skill", "reduce token
  size", "split this skill", "merge skills", "adapt this for Codex", or "design
  a reusable agent workflow". Do NOT use for generating the target artifact a
  skill would produce, for editing AGENTS.md / CLAUDE.md /
  .github/copilot-instructions.md repo-policy files, or for one-off slash-command
  or system-prompt files; output MUST be a skill folder containing SKILL.md plus
  optional references, templates, scripts, assets, or examples.
---

# Skillify - Skill Creation Meta-Skill

## Purpose

Create, improve, and validate focused `SKILL.md` files plus the minimum supporting resources an agent needs to perform a reusable workflow.

## When to use this skill

- Use for skill creation, refactoring, review, audit, compression, splitting, merging, or platform adaptation.
- Use when the output should be a skill folder, a `SKILL.md`, skill references, skill templates, skill scripts, or a skill review report.
- Do NOT use when the user wants the target code itself, repository policy files, or a one-off prompt that does not need a reusable skill.

## Modes

| Mode | Use when the user asks to | Output |
|------|---------------------------|--------|
| Create | "create a skill", "write a SKILL.md", "build a skill for X" | New skill folder |
| Refactor | "refactor my skill", "improve this skill", "fix this SKILL.md" | Modified skill folder or patch |
| Review | "review my skill", "is this skill good", "feedback on my SKILL.md" | Findings report, no edits |
| Audit | "audit my skill", "score my skill", "check against the rubric" | Scored rubric and delta list |
| Compress | "shrink this skill", "reduce token cost", "compress my SKILL.md" | Leaner skill with Tier-3 migration |
| Split | "split this skill", "this skill does too much" | Two or more focused skill folders |
| Merge | "merge these skills", "combine X and Y skills" | One merged skill or rejection rationale |
| Adapt | "adapt this for Codex", "make this work in Copilot/Gemini" | Platform-adapted skill |

Run Create mode from this file. For Refactor, Review, Audit, Compress, Split, Merge, and Adapt, read `references/mode-playbooks.md` before changing or reporting on the skill.

## Universal preamble

Run this before every mode.

1. Detect the mode from the user's phrasing and the Modes table.
2. If multiple modes match, pick by the required output rather than the verb:
   - prose feedback only → Review.
   - scored report against the rubric → Audit.
   - edits applied to the existing skill → Refactor (or Compress when the only complaint is token size).
   - two skills out of one → Split.
   - one skill out of two → Merge.
   - same skill on a different platform → Adapt.
3. If no mode matches, or an ambiguity remains after step 2, ask one disambiguating question instead of guessing.
4. Establish the input:
   - Create: collect the intended task, target users, and at least 3 concrete user prompts that should trigger the skill. If the user provided fewer than 3, ask once before continuing; do not invent trigger phrases.
   - All other modes: ask the user for the path to the existing skill folder if not stated, then read its `SKILL.md` and enumerate `references/`, `templates/`, `scripts/`, `assets/`, and `examples/`. Inspect only the files the requested mode needs.
5. State the output contract: list the exact files you will create, modify, or read-only review, and announce the detected mode so the user can correct it.
6. For Refactor, Compress, Split, and Merge, preserve the original by writing changes either as a unified diff or into a sibling directory unless the user explicitly authorizes in-place overwrite. Review and Audit never edit files.

## Core workflow: Create

1. Analyze the target workflow.
   - Identify the specific job the skill must perform and write a one-sentence responsibility statement.
   - Reuse the trigger prompts collected in the universal preamble; do not collect them a second time and do not invent new ones.
   - Classify the skill type: Document/Asset Creation, Workflow Automation, MCP Enhancement, Domain Expertise, Code Review, or Planning.
   - Choose the degree of freedom and apply its consequence in later steps:
     - High → step descriptions in prose with explicit checklists.
     - Medium → numbered steps backed by `templates/` skeletons.
     - Low → numbered steps that delegate fragile or irreversible logic to `scripts/`.
2. Design the skill boundary.
   - Name one responsibility the skill owns.
   - Identify adjacent tasks it must not handle.
   - Split the design if the description needs more than two unrelated "and" clauses.
3. Plan reusable contents.
   - Use `scripts/` for deterministic or frequently repeated operations.
   - Use `references/` for deep guidance that is not always needed.
   - Use `templates/` for reusable skeletons.
   - Use `assets/` only for output resources the agent should not read into context.
4. Draft the frontmatter.
   - Set `name` to lowercase kebab-case, under 64 characters, matching the folder name.
   - Write `description` as: what the skill does, when to use it with literal trigger phrases, and key capabilities.
   - Include at least one negative trigger when adjacent skills or repo-policy files could be confused with this skill.
   - Keep the description under 1024 characters and avoid XML angle brackets.
5. Write the body.
   - Use imperative, numbered steps with entry and exit conditions.
   - State exact output files, paths, formats, and naming rules.
   - Move variant-specific, long, or optional material to one-level-deep reference files.
   - Point to references instead of duplicating their content.
6. Validate and iterate.
   - Run `scripts/quick_validate.py <skill-folder>` when available.
   - Run `scripts/check_links.py <skill-folder>` when available.
   - Apply `references/validation-rubric.md` and produce a filled score table.
   - Sweep `references/anti-patterns.md` and `references/security-checklist.md`.
   - Iterate up to three passes to clear failing gates. If a gate still fails after three passes, stop and surface the blocker rather than guessing further.

## Validation gate

Every mode exits through these gates.

1. Deterministic checks pass: frontmatter parses, `name` matches folder, `description` is under 1024 characters, no XML angle brackets appear in frontmatter, banned human-facing docs are absent, and local links resolve.
2. Rubric score passes: all 10 dimensions in `references/validation-rubric.md` score at least 4/5 and the total is at least 40/50.
3. Anti-pattern sweep passes: every applicable item in `references/anti-patterns.md` is absent or explicitly mitigated.
4. Security sweep passes: `references/security-checklist.md` finds no unreviewed external HTTP, secret access, destructive command, broad permission, or vendor-bias risk.

If gate 1 (deterministic) fails, stop, surface the script error verbatim, and do not write target files. For gates 2-4, behavior depends on the mode: in Create, Refactor, Compress, Split, Merge, and Adapt, iterate up to three passes to clear failures; in Review and Audit, report the failures and do not edit unless the user explicitly upgrades the request.

## Output format

Produce skill files only. Standard structure:

```text
skill-name/
├── SKILL.md
├── references/
├── templates/
├── scripts/
├── platforms/
├── assets/
└── examples/
```

Use optional directories only when they reduce context load or improve reliability. Do not create auxiliary human docs inside the skill folder.

### Completion report

Every mode finishes by reporting to the user:

| Mode | Required report fields |
|------|------------------------|
| Create | Created files, validation gate results, any inferred decisions left for confirmation. |
| Refactor | Changed files (diff or sibling path), failing-then-passing rubric dimensions, remaining risks. |
| Review | Top three risks, rubric scores, recommended next mode. No edits performed. |
| Audit | Filled scoring table, anti-pattern sweep results, security findings, delta list. |
| Compress | Original vs. new line count, files moved to `references/`, any behavior risk introduced. |
| Split | Source-to-target file mapping, new skill names and triggers, content intentionally dropped. |
| Merge | Files removed as duplicates, unified trigger set, retained negative triggers. |
| Adapt | Source platform, target platform, removed or translated frontmatter fields, platform caveats. |

## Constraints

- DO NOT generate the target artifact when the user asked for a skill that would generate it.
- DO NOT edit `AGENTS.md`, `CLAUDE.md`, `.github/copilot-instructions.md`, or repo-policy files unless the user explicitly asks for platform-policy adaptation.
- DO NOT modify or delete files in Review or Audit modes.
- DO NOT invent trigger phrases, target users, or output contracts when the user has not provided them; ask once first.
- DO NOT place long reference material directly in `SKILL.md`; move anything over roughly 50 lines of supporting detail to `references/`.
- DO NOT duplicate guidance between `SKILL.md` and reference files.
- DO NOT create `README.md`, `CHANGELOG.md`, `INSTALLATION_GUIDE.md`, `QUICK_REFERENCE.md`, or `CONTRIBUTING.md` inside a skill folder.
- DO NOT include `claude` or `anthropic` in skill names.
- DO NOT use XML angle brackets in frontmatter.

## Troubleshooting

| Signal | Action |
|--------|--------|
| Under-triggering | Add the missed user phrasing to `description`. |
| Over-triggering | Add a negative trigger or narrow the positive triggers. |
| Context bloat | Run Compress mode and migrate optional detail to Tier 3. |
| Execution drift | Tighten the degree of freedom or replace prose with a script. |
| Stuck after three iteration passes | Stop, surface the remaining blocker, and ask the user before continuing. |

For post-ship maintenance, read `references/lifecycle-and-iteration.md`.

## References

| Need | Reference |
|------|-----------|
| Frontmatter fields and examples | `references/frontmatter-guide.md` |
| Per-mode workflows beyond Create | `references/mode-playbooks.md` |
| Post-ship iteration and retirement | `references/lifecycle-and-iteration.md` |
| Anti-patterns to avoid | `references/anti-patterns.md` |
| Scoring rubric for audits | `references/validation-rubric.md` |
| Progressive disclosure and splitting | `references/progressive-disclosure.md` |
| Workflow structure patterns | `references/workflow-patterns.md` |
| Platform adaptation | `references/platform-compatibility.md` |
| Safety review before enabling a skill | `references/security-checklist.md` |
| Cross-platform deployment | `platforms/deployment-guide.md` |

## Templates and scripts

- `templates/basic-skill-template.md` - Standard single-workflow skill.
- `templates/mcp-skill-template.md` - Multi-tool orchestration skill.
- `templates/domain-skill-template.md` - Domain-expertise skill.
- `templates/audit-report-template.md` - Review and audit report shape.
- `scripts/init_skill.py` - Boilerplate skill folder generator.
- `scripts/quick_validate.py` - Deterministic frontmatter and structure validator.
- `scripts/check_links.py` - Local reference/template link checker.
- `platforms/install.sh` - Cross-platform skill installation script.
