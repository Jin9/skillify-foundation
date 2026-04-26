# Skill Lifecycle and Iteration

Use this runbook after a skill has shipped or after real agent traces reveal trigger, execution, or maintenance problems.

## Signal to Action Map

| Signal | Evidence | Action |
|--------|----------|--------|
| Under-triggering | User intent matches the skill, but the agent did not load it | Add the missed phrase to `description`; keep the body unchanged unless workflow behavior also failed. |
| Over-triggering | Skill loads for adjacent or unrelated tasks | Add a negative trigger or narrow the positive trigger phrases. |
| Scope creep | Description or workflow handles unrelated domains | Run Split mode and create focused sibling skills. |
| Context bloat | `SKILL.md` exceeds 500 lines or repeats reference material | Run Compress mode and move optional detail to `references/`. |
| Execution drift | Agents skip steps or produce inconsistent output | Tighten degree of freedom: prose to checklist, checklist to pseudocode, pseudocode to script. |
| Validation drift | A once-passing skill now fails rubric or deterministic checks | Run Refactor mode against the failing dimensions only. |
| Platform drift | A skill assumes one agent runtime and fails elsewhere | Run Adapt mode and update `compatibility:` plus platform notes. |
| Staleness | Skill has not been useful in recent work | Archive or disable it; loaded skills compete for context. |

## Maintenance Workflow

1. Collect the concrete trace: user prompt, whether the skill loaded, files touched, and observed failure.
2. Classify the trace using the signal table.
3. Apply the smallest targeted change that addresses the signal.
4. Re-run `scripts/quick_validate.py` and `scripts/check_links.py`.
5. Re-score only the rubric dimensions affected by the change, unless the skill had no recent full audit.
6. Record the changed trigger phrase, workflow step, reference, or script in the user-facing completion note.

## Retirement Criteria

Retire or archive a skill when one of these is true:

1. It handles a one-off task that no longer recurs.
2. It duplicates a stronger skill with clearer triggers.
3. Its instructions are now general model knowledge and no longer earn context cost.
4. It cannot pass safety or conflict checks without becoming too narrow to be useful.
