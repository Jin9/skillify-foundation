# Skillify Mode Playbooks

Use this file after `SKILL.md` mode detection. Create mode is defined in `SKILL.md`; this file holds the detailed workflows for every other mode so the main skill stays lean.

## Cross-mode rules

These rules apply to every mode in this file:

- **Target resolution.** If the user pointed to a skill via path, `@`-mention, pasted text, or current working directory, treat that as the target; the clarification budget in `SKILL.md` (Universal preamble, step 3) governs the one question allowed when it is genuinely absent.
- **Preserve the original.** Per `SKILL.md` Universal preamble step 6: for Refactor, Compress, Split, and Merge, write changes as a unified diff or into a sibling directory unless the user authorizes in-place overwrite or supplies an explicit target path or output contract.
- **Iteration cap.** Run at most three improvement passes per mode. If gates still fail, stop and report the remaining blockers.
- **Multi-mode requests.** When a user combines modes (for example, "create and immediately compress"), execute them in this order: Create → Refactor → Compress → Split → Merge → Adapt → Review → Audit. Skip steps the user did not ask for. State the chain in the opening line and proceed; pause first only when the chain includes an in-place overwrite the user has not authorized.
- **Final report.** Each mode finishes with the report shape listed in `SKILL.md` "Completion report".
- **Model-cost tier.** Review, Audit, Refactor, Split, and Merge are judgment work: frontier. Compress and Adapt are mostly mechanical: mid. Gate-1 scripts and `scripts/cruft_scan.py`: small.

## Create

Use `SKILL.md` Core workflow: Create. Return here only for the Cross-mode rules above (chain order, preservation, report shape) when the create request also includes a refactor, audit, compression, split, merge, or platform-adaptation requirement.

## Refactor

1. Read the current `SKILL.md` and inventory all one-level-deep `references/`, `templates/`, `scripts/`, `assets/`, and `examples/`.
2. Identify the user's stated pain point. If the pain point does not map directly to a rubric dimension (for example, "examples are too long"), translate it to the closest dimension (Token Efficiency in that case) and note the translation in the final report. If no pain point is stated, compare the skill against `validation-rubric.md` and `anti-patterns.md`, and run `scripts/cruft_scan.py` for dated patterns. If the pain point is a model change, cruft, or over-prescription, use the Re-baseline sub-flow below.
3. Score each rubric dimension and mark every dimension below 4/5.
4. Apply targeted fixes only to failing or requested areas; preserve passing sections, verification steps, exact commands for fragile operations, and existing useful references. Write changes as a diff or into a sibling folder per the Cross-mode rules.
5. Remove duplicated cross-tier content by keeping the authoritative version in one file and replacing copies with pointers.
6. Run the Validation gate from `SKILL.md` and preserve script output.
7. Produce a before/after summary with changed files, remaining risks, and final rubric scores.

### Refactor sub-flow: Re-baseline for a model generation

Use when the user says "re-baseline", "fit for the new model", "remove cruft", "de-prescribe", or "migrate this skill to the new model", or when Refactor step 2 lands here. Principles: `model-generation-fit.md`. Output is a Refactor output (diff or sibling); a report-only request is Audit.

1. Fix the target generation and the scope without asking: the generation the user named, else the current frontier tier; the skill(s) named. State both at the top of the report.
2. Run `scripts/cruft_scan.py <skill>`; read the keep list and the current-generation profile in `model-generation-fit.md`; inventory `SKILL.md`, references, templates, and examples. Where git history exists, `git blame` the emphatic and prohibitive lines.
3. Classify every scanner finding and every emphatic or numbered line: keep (any item on the keep list, `model-generation-fit.md` section 8), remove, rewrite, move, or add. For each keep-or-remove decision ask which failure, on which generation, the line prevented, and whether it still reproduces.
4. De-prescribe: keep a step number only where a later step consumes an earlier step's output or the operation is fragile or irreversible; collapse judgment choreography to goal, reasoned constraints, verification, and stop condition. Preserve exact commands, scripts, verification text, and model-cost tiers.
5. Add what the generation needs: the operating contract from `templates/operating-contract.md`, reasons beside the one or two real constraints, a learnings location for long-running skills, delegation guidance where work is parallelizable.
6. Produce both deliverables: the report (one entry per finding with location, evidence, pattern, why obsolete, confidence, action; the Cruft Scan section of `templates/audit-report-template.md`) and the diff or sibling with one finding per hunk. Low-confidence items appear in the report only.
7. Verify: gates 1-4 from `SKILL.md`; `scripts/cruft_scan.py <skill> --strict` exits 0, where a High finding retained under the keep list is suppressed with an inline `cruft-scan: allow` comment naming the signal and the reason; skill-local scripts and tests still pass; when an eval or the collected trigger prompts exist, A/B one prompt with the scaffolding removed against the original, otherwise record the removal as untested.
8. Report findings by confidence, hunks applied, lines added versus removed, untested removals, and the generation-profile date used.

## Review

1. Read the submitted `SKILL.md` and only the references needed to understand its claims.
2. Do not edit files unless the user changes the request from review to refactor.
3. Score the skill against `validation-rubric.md`.
4. Identify the top 3 risks, ordered by severity: trigger errors, scope bloat, workflow ambiguity, missing output contract, security, or portability.
5. Borrow the section list from `templates/audit-report-template.md` (rubric scores, top findings, anti-pattern sweep, security sweep, recommendation), but write narratively rather than filling every table form.
6. End with clear next actions: no changes needed, targeted refactor recommended, split recommended, or reject as not worth a skill.

## Audit

1. Run the Validation gate from `SKILL.md` and preserve deterministic script output.
2. Produce a structured report from `templates/audit-report-template.md`, stating the host and model generation the audit assumes; dated scaffolds are cruft relative to a model. Include the `scripts/cruft_scan.py` summary and its High and Medium findings.
3. Do not modify files during Audit mode unless the user explicitly changes the task to Refactor mode.

## Compress

1. Measure the current `SKILL.md`: line count, rough token size, section count, and reference count.
2. Identify content that is not needed on every invocation: long examples, variant-specific guidance, detailed API docs, platform matrices, or repeated checklists.
3. Move that content to one-level-deep reference files with clear ownership.
4. Replace moved content with one-line pointers from `SKILL.md`.
5. Delete duplicate content and restatements of what the agent does unprompted instead of moving them.
6. Run the Validation gate from `SKILL.md` and re-check the `SKILL.md` line count.
7. Stop when `SKILL.md` is under the user-specified budget or, if no budget was given, under 350 lines and 5,000 tokens. Do not compress further than that without confirmation.
8. Report the reduction and any behavior risk introduced by compression.

## Split

1. List every distinct task, trigger phrase, output type, and tool family the current skill handles.
2. Cluster tasks by responsibility; each resulting skill must have one dominant user intent. Aim for two or three resulting skills. If you find more than three clusters, propose merging the thinnest clusters into the closest sibling before splitting.
3. Propose the new skill names, descriptions, and negative triggers in the opening line, then create the sibling folders; renaming a sibling is cheap. Pause first only when the split would overwrite the original in place and the user has not authorized that.
4. Create separate skill folders and allocate references to the one skill that owns them.
5. Add sibling negative triggers where over-triggering is likely.
6. Run the Validation gate from `SKILL.md` on each new folder.
7. Report the mapping from old content to new skills and identify anything intentionally dropped.

## Merge

1. Read all candidate skills and summarize their scope, triggers, outputs, and references.
2. Estimate overlap; reject the merge if the shared responsibility is below roughly 70%.
3. If merging is justified, choose the clearest name and unify the trigger phrases. If neither original name covers the merged scope, create the sibling under the proposed name and say so in the opening line. Updating incoming pointers in other skills is a scope change: list them in the recap and do not edit them without the user's say-so.
4. Deduplicate references and keep each topic in exactly one file.
5. Preserve negative triggers for tasks that were intentionally left out.
6. Run the Validation gate from `SKILL.md` on the merged skill, then validate it against the one-responsibility rule; abort if it becomes a "do everything" skill.
7. Report merged files, removed duplicates, and any rejected source content.

## Adapt

1. Identify the source platform and target platform from the user request. If the source platform is unstated, infer it from the existing folder location (`.claude/skills/`, `skills/`, `.gemini/`, `.agents/skills/`, `.github/skills/`, `~/.copilot/skills/`, `$CODEX_HOME/skills/`) or, if still unclear, from the host-specific frontmatter fields present (`frontmatter-guide.md`), and state the assumed source platform in the opening line.
2. Read `platform-compatibility.md` and compare frontmatter, folder location, root rules, tools, and activation model.
3. Translate platform-specific fields; remove unsupported fields instead of leaving dead metadata; turn a `model:` or `effort:` pin into a body-level model-cost tier and effort hint.
4. Preserve workflow logic and replace vendor/tool-specific language with generic tool capability names.
5. For Codex packaging, document any needed `agents/openai.yaml` UI metadata: `display_name`, `short_description`, and `default_prompt`.
6. Add or update a concise `compatibility:` field when the target platform supports it.
7. Run the Validation gate from `SKILL.md` and report platform-specific caveats.
