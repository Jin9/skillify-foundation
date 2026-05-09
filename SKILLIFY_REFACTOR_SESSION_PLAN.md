# Skillify Refactor Session Plan

## 1. Objective

Refactor the `skillify` meta-skill so it is clearer, more reliable, and easier to validate without expanding its scope.

Primary target:

- `skillify/SKILL.md`
- `skillify/references/mode-playbooks.md`
- installed copy under `/Users/IF640063/.codex/skills/skillify/`

Supporting targets:

- `skill-design-methodology/research/`
- `skill-design-methodology/evals/`
- `literature/skill-design-methodology/`
- `discovery/security-reviewer-source-request.md`

## 2. Workflow Used

The session followed the cross-model workflow from `skill-design-methodology/plan.md`.

1. **Codex Phase 0:** Baseline inventory of the current Skillify skill, support files, validation scripts, dirty worktree state, and installed copy.
2. **Gemini Phase 1:** Source discovery and source tiering using `gemini-3.1-pro-preview`.
3. **Gemini Phase 2:** Synthesis of operational skill-design rules from the source set.
4. **Claude Phase 3:** Architecture analysis using Claude Opus with max effort.
5. **Claude Phase 4:** Refactor recommendation using Claude Opus with max effort.
6. **Codex Phase 5:** Apply the focused refactor.
7. **Claude Phase 6:** Review the changed skill and return accept/reject verdict.
8. **Codex Phase 7:** Compact, validate, sync installed copy, and record final status.

Important process correction:

- An earlier attempt substituted Codex fallback artifacts when Gemini/Claude ran into CLI issues.
- The workflow was rerun with actual Gemini and Claude command execution.
- Claude plan-mode outputs landed under `.claude/plans/`; their content was then copied into the required research artifacts.

## 3. Refactor Changes Applied

### `skillify/SKILL.md`

Changes:

- Replaced redundant body-level trigger prose with a clearer `Scope Boundary`.
- Compressed non-Create mode routing to point to `references/mode-playbooks.md`.
- Moved Create-specific prompt collection from the Universal preamble into `Core workflow: Create`.
- Wired `scripts/init_skill.py <skill-name>` into Create mode when no target folder exists.
- Shortened Validation gate step 1 to rely on `quick_validate.py` and `check_links.py`.
- Single-sourced the iteration cap in `references/mode-playbooks.md`.
- Removed `platforms/` from the generic generated-skill output tree.
- Clarified that `platforms/` is Skillify-internal install/distribution material.
- Removed the `platforms/deployment-guide.md` row from the general references table.
- Replaced the ad hoc 50-line threshold with a pointer to `references/progressive-disclosure.md`.

### `skillify/references/mode-playbooks.md`

Changes:

- Refactor mode now points to the `SKILL.md` Validation gate rather than naming validator scripts directly.
- Audit mode now points to the `SKILL.md` Validation gate and preserves deterministic script output.
- Iteration cap remains canonical in Cross-mode rules.

## 4. Research Artifacts Created

Files:

- `skill-design-methodology/research/sources.md`
- `skill-design-methodology/research/research-notes.md`
- `skill-design-methodology/research/analysis.md`
- `skill-design-methodology/research/recommendation.md`
- `skill-design-methodology/research/review.md`
- `skill-design-methodology/research/skillify-architecture-design.md`

Purpose:

- Preserve model outputs and decision rationale.
- Make the refactor auditable.
- Separate source discovery, synthesis, architecture analysis, recommendation, and review.

## 5. Evaluation Suite Added

Files:

- `skill-design-methodology/evals/golden_prompts.json`
- `skill-design-methodology/evals/run_skillify_eval.py`
- `skill-design-methodology/evals/latest_report.md`

The eval compares:

- before: committed `HEAD` version of Skillify
- after: current working-tree Skillify

Indicators:

- deterministic validator pass/fail
- banned docs absence
- description length
- trigger count
- mode count
- golden prompt pass rate
- structural indicator pass rate
- key refactor invariants

Latest meaningful pre-commit result during the refactor:

| Indicator | Before | After | Result |
|---|---:|---:|---|
| Golden prompt pass rate | 12/12 | 12/12 | preserved |
| Structural indicator pass rate | 11/18 | 18/18 | improved |
| Approx tokens | 2806 | 2770 | improved |
| Mode count | 8 | 8 | preserved |
| Quoted trigger count | 12 | 12 | preserved |

Improved structural indicators:

- prompt collection moved into Create
- `init_skill.py` wired into Create
- `platforms/` removed from generic output tree
- validation gate script pointer simplified
- iteration cap single-sourced
- 50-line threshold removed
- deployment guide removed from generic reference route map

## 6. Validation Commands

Run these after future Skillify changes:

```sh
python3 skillify/scripts/quick_validate.py skillify
python3 skillify/scripts/check_links.py skillify
python3 skill-design-methodology/evals/run_skillify_eval.py
```

For the installed copy:

```sh
python3 /Users/IF640063/.codex/skills/skillify/scripts/quick_validate.py /Users/IF640063/.codex/skills/skillify
python3 /Users/IF640063/.codex/skills/skillify/scripts/check_links.py /Users/IF640063/.codex/skills/skillify
```

Sync checks:

```sh
cmp -s skillify/SKILL.md /Users/IF640063/.codex/skills/skillify/SKILL.md
cmp -s skillify/references/mode-playbooks.md /Users/IF640063/.codex/skills/skillify/references/mode-playbooks.md
```

## 7. Repository Organization Cleanup

Moved generic literature/source material out of `skill-design-methodology/`:

- `agent_reading_taxonomy.md` -> `literature/skill-design-methodology/agent-reading-task-taxonomy.md`
- `reference_links.md` -> `literature/skill-design-methodology/agent-skill-source-index.md`
- `skill_design_principles.md` -> `literature/skill-design-methodology/agent-skill-design-principles.md`

Renamed Skillify-specific workflow assets:

- `llm_skillification_pipeline.md` -> `skill-design-methodology/cross-model-skillification-pipeline.md`
- `skill_architecture_design.md` -> `skill-design-methodology/research/skillify-architecture-design.md`

Moved unrelated security-reviewer discovery request out of Skillify methodology:

- `source_discovery.md` -> `discovery/security-reviewer-source-request.md`

## 8. Install Path Cleanup

Updated platform install behavior:

- Added `.gemini/` to `.gitignore`.
- Removed Gemini native skill install path from `skillify/platforms/install.sh`.
- Made Codex compatibility install path opt-in via `INSTALL_CODEX_COMPAT=1`.
- Updated `skillify/platforms/deployment-guide.md` to document the duplicate-discovery avoidance.

## 9. Commits Created

Commits pushed to `origin/main`:

- `45504ec Improve Skillify workflow and evals`
- `58dc51b Avoid duplicate Skillify installs`
- `43a2a41 Organize Skillify methodology sources`

## 10. Remaining Follow-Ups

Non-blocking follow-ups from Claude review:

1. Apply the Validation-gate pointer pattern to Compress, Split, Merge, and Adapt in `mode-playbooks.md`.
2. Decide whether `check_links.py` should validate `platforms/` paths or whether platform paths should stay outside validated route maps.
3. Decide whether `init_skill.py` should create standard subdirectories by default or be called with flags.
4. Add a `quick_validate.py` check that each `compatibility:` token has a matching platform file.
5. Reduce overlap between `Scope Boundary` and `Constraints` in `SKILL.md`.

## 11. Definition Of Done

The refactor is considered complete when:

- `quick_validate.py` passes for source and installed Skillify.
- `check_links.py` passes for source and installed Skillify.
- golden prompts pass at 12/12.
- structural indicators do not regress.
- source and installed `SKILL.md` match.
- source and installed `mode-playbooks.md` match.
- no banned human-facing docs are added inside any skill folder.
