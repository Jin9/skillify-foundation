# Skill Lifecycle and Iteration

Use this runbook after a skill has shipped or after real agent traces reveal trigger, execution, or maintenance problems. It also applies mid-task: when a trace reveals the problem while the skill is in use, propose the smallest change and route it through Refactor's preservation rules (diff or sibling; in place only with the user's authorization), then continue.

## Signal to Action Map

| Signal | Evidence | Action |
|--------|----------|--------|
| Under-triggering | User intent matches the skill, but the agent did not load it | Add the missed intent to `description` as a category plus at most one literal phrase; keep the body unchanged unless workflow behavior also failed. Descriptions that grow one phrase per miss generalize worse. |
| Over-triggering | Skill loads for adjacent or unrelated tasks | Add a negative trigger or narrow the positive trigger phrases. |
| Scope creep | Description or workflow handles unrelated domains | Run Split mode and create focused sibling skills. |
| Context bloat | `SKILL.md` exceeds 500 lines or repeats reference material | Run Compress mode and move optional detail to `references/`. |
| Execution drift | Agents skip steps or produce inconsistent output | Check the trace against tool results first. A skipped fragile step: tighten that step only (checklist to exact command to script). A script followed into a worse result: loosen it to goal, constraints, and verification. |
| Validation drift | A once-passing skill now fails rubric or deterministic checks | Run Refactor mode against the failing dimensions only. |
| Platform drift | A skill assumes one agent runtime and fails elsewhere | Run Adapt mode and update `compatibility:` plus platform notes. |
| Staleness | Skill has not been useful in recent work | Archive or disable it; loaded skills compete for context. |
| Model release | A host ships a new frontier model, or a skill written for an older model is reused | First re-verify the current-generation profile in `model-generation-fit.md` section 9 against the release's prompting guide and update its verified date. Then run `scripts/cruft_scan.py` across skills in active use and the Re-baseline sub-flow of Refactor (`mode-playbooks.md`) on those with High findings or a multi-step workflow with no operating contract. Re-baselining can add text as well as remove it. |
| Unrequested pause or divergence | With the skill loaded, the agent asks for permission, stops early, or changes direction against the user's request, and no destructive, irreversible, or scope-changing action is at hand | Have the agent name the `SKILL.md`, quote the instruction, and explain how it applied; rewrite that line so the user's request sets the scope (anti-pattern 6), then re-score rubric dimension 6. |

## Maintenance Workflow

1. Collect the concrete trace: user prompt, whether the skill loaded, files touched, and observed failure.
2. Classify the trace using the signal table.
3. Apply the smallest targeted change that addresses the signal.
4. Re-run `scripts/quick_validate.py` and `scripts/check_links.py`; run `scripts/cruft_scan.py` and carry High findings into the change.
5. Re-score only the rubric dimensions affected by the change, unless the skill had no recent full audit.
6. Record the changed trigger phrase, workflow step, reference, or script in the user-facing completion note.

## Retirement Criteria

Retire or archive a skill when one of these is true:

1. It handles a one-off task that no longer recurs.
2. It duplicates a stronger skill with clearer triggers.
3. Its instructions are now general model knowledge and no longer earn context cost.
4. It cannot pass safety or conflict checks without becoming too narrow to be useful.
