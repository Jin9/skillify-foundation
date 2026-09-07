---
name: skillify
description: >
  Create, refactor, review, audit, compress, split, merge, or adapt SKILL.md files
  and supporting skill assets for AI agent workflows. Use when the user asks to
  "create a skill", "write a SKILL.md", "design an agent skill", "build a skill
  for X", "refactor a skill", "review my skill", "audit a skill", "reduce token
  size", "split this skill", "merge skills", "adapt this for Codex", "re-baseline
  this skill for a new model", or "design a reusable agent workflow". Do NOT use for generating the target artifact a
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

- Produce skill folders, `SKILL.md` files, references, templates, scripts, assets, examples, or skill review reports — engineering `SKILL.md` *workflows*, not custom-agent personas, tool permissions, or host model configuration (see `references/platform-compatibility.md`). Per-step model-cost tier and effort hints inside a skill body are in scope (`references/model-generation-fit.md`).
- For what this skill must not do, see Constraints.

## Modes

| Mode | Use when the user asks to | Output |
|------|---------------------------|--------|
| Create | "create a skill", "write a SKILL.md", "build a skill for X" | New skill folder |
| Refactor | "refactor my skill", "improve this skill", "fix this SKILL.md", "re-baseline this skill for a new model", "remove cruft" | Modified skill folder or patch |
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
   - cruft removal or fit for a new model generation: Refactor (Re-baseline sub-flow in `references/mode-playbooks.md`); report-only cruft findings: Audit.
3. Clarification budget: ask at most one question per run, and only when the mode is still ambiguous after step 2 or a required input (the target skill, or the trigger prompts for Create) is absent and cannot be taken from the request, pasted text, or the working directory. Otherwise proceed and state the assumption in the opening line.
4. Establish the input:
   - If the user pointed to a skill by path, mention, pasted text, or current directory, treat that as the target.
   - Then read `SKILL.md`, enumerate one-level `references/`, `templates/`, `scripts/`, `assets/`, and `examples/`, and inspect only the files needed for the mode.
5. State the output contract: list the exact files you will create, modify, or read-only review, and announce the detected mode so the user can correct it.
6. For Refactor, Compress, Split, and Merge, preserve the original by writing a unified diff or sibling directory unless the user authorizes in-place overwrite or supplies an explicit target path/output contract. Review and Audit never edit files.

## Core workflow: Create

1. Analyze the target workflow.
   - Collect the intended task and at least 3 concrete user prompts that should trigger the skill. Phrases the user already wrote count; quote them verbatim. Use the clarification budget only if fewer than 3 can be taken from the request; do not invent prompts, because invented triggers cause over-triggering that is hard to trace. Note target users only when stated.
   - Identify the specific job the skill must perform and write a one-sentence responsibility statement.
   - Classify the skill type: Document/Asset Creation, Workflow Automation, MCP Enhancement, Domain Expertise, Code Review, or Planning.
   - Choose the degree of freedom by fragility, not habit, and apply it in step 5 (`references/model-generation-fit.md` section 1):
     - High (judgment work) → goal, constraints with their reasons, and how to verify; number steps only where order matters.
     - Medium → numbered steps backed by `templates/` skeletons.
     - Low (fragile or irreversible) → numbered steps whose exact commands live in `scripts/`, marked "do not modify".
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
   - Use imperative, numbered steps with entry and exit conditions wherever order or fragility matters. For steps the agent should judge, state the outcome, the one or two real constraints with their reason, and how the agent verifies the result; do not script the method.
   - For skills with more than one step, include the operating contract from `templates/operating-contract.md` with its bracketed lines filled: instruction priority, autonomy, stop conditions, verification, delegation, progress, model-cost tier.
   - Annotate phases or steps with a model-cost tier (`small`, `mid`, `frontier`) and, where it differs from the host default, an effort hint. Name tiers, never models.
   - State exact output files, paths, formats, and naming rules.
   - Move variant-specific, long, or optional material to one-level-deep reference files.
   - Point to references instead of duplicating their content.
6. Run the Validation gate and apply its iteration rules.

## Validation gate

Every mode exits through these gates, in order.

1. Run deterministic checks when scripts are available: `scripts/quick_validate.py` and `scripts/check_links.py` are blocking. If either exits non-zero, surface the error verbatim, fix the cause, and re-run within the iteration cap in `references/mode-playbooks.md`. If it still fails, stop: do not write target files and do not claim the gate passed. Then run `scripts/cruft_scan.py` (advisory): carry High findings into the report as required changes; use `--strict` only as the exit check of the Re-baseline sub-flow.
2. Score `references/validation-rubric.md`: every dimension must be at least 4/5 and total at least 40/50.
3. Sweep `references/anti-patterns.md`: every applicable item must be absent or explicitly mitigated.
4. Sweep `references/security-checklist.md`: no unreviewed external HTTP, secret access, destructive command, broad permission, vendor-bias, unrequested-pause, or reasoning-extraction risk.

Iterate on gates 1-4 per the cap in `references/mode-playbooks.md`. In Review and Audit, report every failure, gate 1 included, and do not edit unless the user upgrades the request. Quote actual validator and scanner output in every report; never describe a check you did not run.

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
| Refactor | Changed files (diff or sibling path), failing-then-passing rubric dimensions, cruft findings before and after, remaining risks. |
| Review | Top three risks, rubric scores, recommended next mode. No edits performed. |
| Audit | Filled scoring table, cruft scan summary, anti-pattern sweep results, security findings, delta list. |
| Compress | Original vs. new line count, files moved to `references/`, any behavior risk introduced. |
| Split | Source-to-target file mapping, new skill names and triggers, content intentionally dropped. |
| Merge | Files removed as duplicates, unified trigger set, retained negative triggers. |
| Adapt | Source platform, target platform, removed or translated frontmatter fields, platform caveats. |

## Operating contract

- Instruction priority: the user's request in this session takes precedence over this skill; repo-policy files (`AGENTS.md`, `CLAUDE.md`, or the host equivalent) take precedence over this skill's defaults. If following a line here would make you pause, ask for permission, leave requested work unfinished, or diverge from what the user asked, follow the user, say which line you set aside, and quote it.
- Autonomy: once the mode and target are established, act; the clarification budget above is the only question. State assumptions in the opening line.
- Stop conditions: stop and ask only before overwriting a skill in place without authorization, before editing a repo-policy file, when the iteration cap is hit with gates still failing, and when finishing would change the scope the user set (for example, updating incoming pointers in other skills during Merge; list them in the recap instead).
- Verification: gates 1-4 above, with the validator and scanner output quoted in the completion report.
- Delegation: none by default. A re-baseline sweep across many skills may fan out one skill per sub-agent; the lead reconciles the reports.
- Progress: open with the detected mode and the files to be created, modified, or read; close with the mode's completion report.
- Model-cost tier: mid by default; Create judgment steps (boundary design, body drafting) are frontier and gate-1 scripts are small. Per-mode tiers for the other modes are in `references/mode-playbooks.md` Cross-mode rules.

## Constraints

- Do not generate the target artifact when the user asked for a skill that would generate it; the skill is the deliverable.
- Do not turn one-off prompts or system-prompt files into skills; a skill must earn its context cost across many sessions.
- Do not edit `AGENTS.md`, `CLAUDE.md`, `.github/copilot-instructions.md`, or other repo-policy files unless the user explicitly asks for platform-policy adaptation; they are always-on policy owned by the repository.
- Review and Audit never edit files; Refactor, Compress, Split, and Merge preserve the original (Universal preamble, step 6).
- Do not invent trigger phrases; quote the user's own (Create step 1), because invented triggers cause over-triggering that is hard to trace. A missing target-user or output-contract detail is stated as an assumption in the opening line under the clarification budget.
- Keep `SKILL.md` lean: long or mode-specific material lives in one reference file per `references/progressive-disclosure.md`, and nothing is stated in both tiers.

For frontmatter rules (kebab-case names, reserved-vendor-name ban, no XML angle brackets, character limits), see `references/frontmatter-guide.md`; for the ban on `README.md`, `CHANGELOG.md`, `INSTALLATION_GUIDE.md`, `QUICK_REFERENCE.md`, or `CONTRIBUTING.md` inside a skill folder, see `references/anti-patterns.md` item 8. `scripts/quick_validate.py` enforces both at gate 1.

## Troubleshooting

- **Stuck after the iteration cap** (see `references/mode-playbooks.md`): stop, surface the remaining blocker, and ask the user before continuing.
- For under-triggering, over-triggering, scope creep, context bloat, execution drift, validation drift, platform drift, staleness, a new model release, or a skill that makes the agent pause or diverge from the user's request, see the Signal-to-Action map in `references/lifecycle-and-iteration.md`.

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
| Fit for the current model generation: calibration, keep list, effort tiers | `references/model-generation-fit.md` |

## Templates and scripts

- `templates/basic-skill-template.md` - Standard single-workflow skill.
- `templates/mcp-skill-template.md` - Multi-tool orchestration skill.
- `templates/domain-skill-template.md` - Domain-expertise skill.
- `templates/audit-report-template.md` - Review and audit report shape.
- `templates/operating-contract.md` - Instruction-priority, autonomy, stop, verification, delegation, progress, and tier block for multi-step skills.
- `scripts/init_skill.py` - Boilerplate skill folder generator.
- `scripts/quick_validate.py` - Deterministic frontmatter and structure validator.
- `scripts/check_links.py` - Local reference/template link checker.
- `scripts/cruft_scan.py` - Advisory dated-pattern scanner; `--strict` for re-baseline exit checks.
