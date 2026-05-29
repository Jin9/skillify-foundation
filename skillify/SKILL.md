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
compatibility: claude-code, codex, copilot, gemini, antigravity
---

# Skillify - Skill Creation Meta-Skill

## Purpose

Create, improve, and validate focused `SKILL.md` files plus the minimum supporting resources an agent needs for a reusable workflow.

## Scope Boundary

- Produce skill folders, `SKILL.md` files, references, templates, scripts, assets, examples, or skill review reports — engineering `SKILL.md` *workflows*, not custom-agent personas, tool permissions, or model selection (see `references/platform-compatibility.md`).
- For what this skill must not do, see Constraints.

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

All modes other than Create run from `references/mode-playbooks.md`; combined-mode chains follow the order in that file.

## Universal preamble

Run this before every mode.

1. Detect the mode from the user's phrasing and the Modes table.
2. If multiple modes match, pick by the required output rather than the verb:
   - prose feedback only: Review.
   - scored report against the rubric: Audit.
   - edits applied to an existing skill: Refactor, or Compress when the only complaint is token size.
   - two skills out of one: Split.
   - one skill out of two: Merge.
   - same skill on a different platform: Adapt.
3. If no mode matches, or an ambiguity remains after step 2, ask one disambiguating question instead of guessing.
4. Establish the input:
   - If the user pointed to a skill by path, mention, or current directory, treat that as the target.
   - Otherwise ask once for the missing target or creation inputs before reading or writing.
   - Then read `SKILL.md`, enumerate one-level `references/`, `templates/`, `scripts/`, `assets/`, and `examples/`, and inspect only the files needed for the mode.
5. State the output contract: list the exact files you will create, modify, or read-only review, and announce the detected mode so the user can correct it.
6. For Refactor, Compress, Split, and Merge, preserve the original by writing a unified diff or sibling directory unless the user authorizes in-place overwrite or supplies an explicit target path/output contract. Review and Audit never edit files.

## Core workflow: Create

1. Analyze the target workflow.
   - Collect the intended task, target users, and at least 3 concrete user prompts that should trigger the skill. If the user provided fewer than 3, ask once; do not invent trigger phrases.
   - Identify the specific job the skill must perform and write a one-sentence responsibility statement.
   - Reuse the collected trigger prompts; do not collect them a second time and do not invent new ones.
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
   - Run `scripts/init_skill.py <skill-name>` to scaffold the folder when a target folder does not already exist.
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
6. Run the Validation gate and apply its iteration rules.

## Validation gate

Every mode exits through these gates, in order.

1. Run deterministic checks when scripts are available: `scripts/quick_validate.py` and `scripts/check_links.py`. If either exits non-zero, surface the error verbatim and stop.
2. Score `references/validation-rubric.md`: every dimension must be at least 4/5 and total at least 40/50.
3. Sweep `references/anti-patterns.md`: every applicable item must be absent or explicitly mitigated.
4. Sweep `references/security-checklist.md`: no unreviewed external HTTP, secret access, destructive command, broad permission, or vendor-bias risk.

If gate 1 fails, do not write target files. For gates 2-4, iterate per the cap in `references/mode-playbooks.md`. In Review and Audit, report failures and do not edit unless the user upgrades the request.

## Output format

Produce skill files only. Standard structure:

```text
skill-name/
├── SKILL.md
├── references/
├── templates/
├── scripts/
├── assets/
└── examples/
```

Use optional directories only when they reduce context load or improve reliability. `platforms/` is a skillify-internal directory for cross-host install assets and is not part of generated skill output. Do not create auxiliary human docs inside the skill folder.

### Completion report

Finish with the mode-specific report fields:

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
- DO NOT handle one-off prompts or system-prompt files that do not need a reusable skill.
- DO NOT edit `AGENTS.md`, `CLAUDE.md`, `.github/copilot-instructions.md`, or other repo-policy files unless the user explicitly asks for platform-policy adaptation.
- DO NOT modify or delete files in Review or Audit modes.
- DO NOT invent trigger phrases, target users, or output contracts when the user has not provided them; ask once first.
- DO NOT place long or mode-specific reference material directly in `SKILL.md`; move it according to `references/progressive-disclosure.md`.
- DO NOT duplicate guidance between `SKILL.md` and reference files.
- DO NOT create `README.md`, `CHANGELOG.md`, `INSTALLATION_GUIDE.md`, `QUICK_REFERENCE.md`, or `CONTRIBUTING.md` inside a skill folder.

For frontmatter rules (kebab-case names, reserved-vendor-name ban, no XML angle brackets, character limits), see `references/frontmatter-guide.md`. `scripts/quick_validate.py` enforces them at gate 1.

## Troubleshooting

- **Stuck after the iteration cap** (see `references/mode-playbooks.md`): stop, surface the remaining blocker, and ask the user before continuing.
- For under-triggering, over-triggering, scope creep, context bloat, execution drift, validation drift, platform drift, and staleness, see the Signal-to-Action map in `references/lifecycle-and-iteration.md`.

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

## Templates and scripts

- `templates/basic-skill-template.md` - Standard single-workflow skill.
- `templates/mcp-skill-template.md` - Multi-tool orchestration skill.
- `templates/domain-skill-template.md` - Domain-expertise skill.
- `templates/audit-report-template.md` - Review and audit report shape.
- `scripts/init_skill.py` - Boilerplate skill folder generator.
- `scripts/quick_validate.py` - Deterministic frontmatter and structure validator.
- `scripts/check_links.py` - Local reference/template link checker.
- `platforms/install.sh` - Cross-platform skill installation script.
