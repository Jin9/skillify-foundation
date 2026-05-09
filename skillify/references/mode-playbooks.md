# Skillify Mode Playbooks

Use this file after `SKILL.md` mode detection. Create mode is defined in `SKILL.md`; this file holds the detailed workflows for every other mode so the main skill stays lean.

## Cross-mode rules

These rules apply to every mode in this file:

- **Path elicitation.** If the user pointed to a skill via path, `@`-mention, or current working directory, treat that as the target. Otherwise ask once before reading anything.
- **Preserve the original.** For Refactor, Compress, Split, and Merge, write changes as a unified diff or into a sibling directory by default. Overwrite in place only when the user explicitly authorizes replacement.
- **Iteration cap.** Run at most three improvement passes per mode. If gates still fail, stop and report the remaining blockers.
- **Multi-mode requests.** When a user combines modes (for example, "create and immediately compress"), execute them in this order: Create → Refactor → Compress → Split → Merge → Adapt → Review → Audit. Skip steps the user did not ask for and confirm the chain before starting.
- **Final report.** Each mode finishes with the report shape listed in `SKILL.md` "Completion report".

## Create

Use `SKILL.md` Core workflow: Create. Return here only for the Cross-mode rules above (chain order, preservation, report shape) when the create request also includes a refactor, audit, compression, split, merge, or platform-adaptation requirement.

## Refactor

1. Read the current `SKILL.md` and inventory all one-level-deep `references/`, `templates/`, `scripts/`, `assets/`, and `examples/`.
2. Identify the user's stated pain point. If the pain point does not map directly to a rubric dimension (for example, "examples are too long"), translate it to the closest dimension (Token Efficiency in that case) and note the translation in the final report. If no pain point is stated, compare the skill against `validation-rubric.md` and `anti-patterns.md`.
3. Score each rubric dimension and mark every dimension below 4/5.
4. Apply targeted fixes only to failing or requested areas; preserve passing sections and existing useful references. Write changes as a diff or into a sibling folder per the Cross-mode rules.
5. Remove duplicated cross-tier content by keeping the authoritative version in one file and replacing copies with pointers.
6. Run the Validation gate from `SKILL.md` and preserve script output.
7. Produce a before/after summary with changed files, remaining risks, and final rubric scores.

## Review

1. Read the submitted `SKILL.md` and only the references needed to understand its claims.
2. Do not edit files unless the user changes the request from review to refactor.
3. Score the skill against `validation-rubric.md`.
4. Identify the top 3 risks, ordered by severity: trigger errors, scope bloat, workflow ambiguity, missing output contract, security, or portability.
5. Borrow the section list from `templates/audit-report-template.md` (rubric scores, top findings, anti-pattern sweep, security sweep, recommendation), but write narratively rather than filling every table form.
6. End with clear next actions: no changes needed, targeted refactor recommended, split recommended, or reject as not worth a skill.

## Audit

1. Run the Validation gate from `SKILL.md` and preserve deterministic script output.
2. Produce a structured report from `templates/audit-report-template.md`.
3. Do not modify files during Audit mode unless the user explicitly changes the task to Refactor mode.

## Compress

1. Measure the current `SKILL.md`: line count, rough token size, section count, and reference count.
2. Identify content that is not needed on every invocation: long examples, variant-specific guidance, detailed API docs, platform matrices, or repeated checklists.
3. Move that content to one-level-deep reference files with clear ownership.
4. Replace moved content with one-line pointers from `SKILL.md`.
5. Delete duplicate content instead of keeping both versions.
6. Run `scripts/check_links.py` and re-check the `SKILL.md` line count.
7. Stop when `SKILL.md` is under the user-specified budget or, if no budget was given, under 350 lines and 5,000 tokens. Do not compress further than that without confirmation.
8. Report the reduction and any behavior risk introduced by compression.

## Split

1. List every distinct task, trigger phrase, output type, and tool family the current skill handles.
2. Cluster tasks by responsibility; each resulting skill must have one dominant user intent. Aim for two or three resulting skills. If you find more than three clusters, propose merging the thinnest clusters into the closest sibling before splitting.
3. Propose the new skill names, descriptions, and negative triggers before moving content; wait for user approval on names before creating folders.
4. Create separate skill folders and allocate references to the one skill that owns them.
5. Add sibling negative triggers where over-triggering is likely.
6. Run validation on each new folder.
7. Report the mapping from old content to new skills and identify anything intentionally dropped.

## Merge

1. Read all candidate skills and summarize their scope, triggers, outputs, and references.
2. Estimate overlap; reject the merge if the shared responsibility is below roughly 70%.
3. If merging is justified, choose the clearest name and unify the trigger phrases. If neither original name covers the merged scope, propose a new name and announce that all incoming pointers will need updating; wait for user approval before creating the folder.
4. Deduplicate references and keep each topic in exactly one file.
5. Preserve negative triggers for tasks that were intentionally left out.
6. Validate the merged skill against the one-responsibility rule; abort if it becomes a "do everything" skill.
7. Report merged files, removed duplicates, and any rejected source content.

## Adapt

1. Identify the source platform and target platform from the user request. If the source platform is unstated, infer it from the existing folder location (`.claude/skills/`, `skills/`, `.gemini/`, `.agents/skills/`) or, if still unclear, ask once before translating.
2. Read `platform-compatibility.md` and compare frontmatter, folder location, root rules, tools, and activation model.
3. Translate platform-specific fields; remove unsupported fields instead of leaving dead metadata.
4. Preserve workflow logic and replace vendor/tool-specific language with generic tool capability names.
5. For Codex packaging, document any needed `agents/openai.yaml` UI metadata: `display_name`, `short_description`, and `default_prompt`.
6. Add or update a concise `compatibility:` field when the target platform supports it.
7. Run validation and report platform-specific caveats.
