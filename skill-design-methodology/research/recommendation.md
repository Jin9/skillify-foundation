---
authored_by: chinnawat.w (git author of record; GitHub Jin9)
phase_owner: Claude Opus 4.7
pipeline_phase: 4
status: accepted
provenance_dated: 2026-05-19
---

# Recommendation: Skillify Minimal One-Pass Refactor

## Recommended Edits

### 1. Shorten validation gate script duplication

**Where:** `skillify/SKILL.md` "Validation gate" step 1; `references/mode-playbooks.md` Refactor and Audit script lines.
**Change:** Keep the explicit script list in `SKILL.md` only. In `mode-playbooks.md`, replace each script-name line with a pointer to run the Validation gate from `SKILL.md` and preserve script output.
**Why:** Single source of truth for the validation gate; future script additions or renames touch one file.

### 2. Move Create-only prompt collection out of the universal preamble

**Where:** `skillify/SKILL.md` "Universal preamble" step 4.
**Change:** Strip the Create sub-bullet from preamble step 4. Leave only generic input/path elicitation. Insert trigger-prompt collection as the first sub-step inside "Core workflow: Create" step 1.
**Why:** The universal preamble currently special-cases Create. Mode-specific input collection belongs in the mode workflow.

### 3. Remove `platforms/` from the generic output tree

**Where:** `skillify/SKILL.md` "Output format" structure block.
**Change:** Delete the `platforms/` line from the generic tree. Add one sentence that `platforms/` is skillify-internal for cross-host install assets and is not part of generated skill output.
**Why:** The tree teaches downstream skill authors what to produce. Including `platforms/` invites copying skillify deployment assets into unrelated skills.

### 4. Single-source the iteration cap in `mode-playbooks.md`

**Where:** `skillify/SKILL.md` Validation gate trailing paragraph and `references/mode-playbooks.md` Cross-mode rules.
**Change:** Keep the cap statement in `mode-playbooks.md`. In `SKILL.md`, replace the repeated three-pass wording with a pointer to the cap in `references/mode-playbooks.md`.
**Why:** Iteration cap is a cross-mode rule and should live with the other cross-mode rules.

### 5. Remove the 50-line threshold from `SKILL.md`

**Where:** `skillify/SKILL.md` Constraints.
**Change:** Replace the 50-line threshold with a pointer to `references/progressive-disclosure.md` thresholds.
**Why:** `progressive-disclosure.md` already owns line/token thresholds. The 50-line figure conflicts with that reference and is not enforced.

### 6. Compress the non-Create mode trailer

**Where:** paragraph after the Modes table.
**Change:** Collapse mode-name repetition to: "All modes other than Create run from `references/mode-playbooks.md`; combined-mode chains follow the order in that file."
**Why:** Equivalent meaning, fewer always-loaded tokens, and no repeated mode list.

### 7. Wire `scripts/init_skill.py` into Create

**Where:** `skillify/SKILL.md` "Core workflow: Create" reusable-contents step.
**Change:** Add a lead sub-step to run `scripts/init_skill.py <skill-name>` to scaffold the folder when a target folder does not already exist.
**Why:** Listing a script no workflow ever invokes is dead surface area. Wiring it in preserves utility.

### 8. Remove overlapping deployment-guide row from References

**Where:** `skillify/SKILL.md` References table.
**Change:** Delete the `platforms/deployment-guide.md` row.
**Why:** `references/platform-compatibility.md` covers author-time platform guidance; deployment guide is install-time documentation and overlaps in the route map.

## Rejected Edits And Rationale

| Rejected proposal | Rationale |
|---|---|
| Add new metadata fields such as version, owner, or last_reviewed | Out of scope for a minimal patch. Current `name`, `description`, and `compatibility` are the portable contract. |
| Split skillify into multiple skills | No over-triggering or scope-bloat signal meets the Split threshold. Shared validation/rubric/security gates belong together. |
| Add cross-model orchestration to `SKILL.md` | The upgrade pipeline is implementation scaffolding, not reusable Skillify behavior. |
| Remove the per-mode Completion report table | It is the only compact mode-by-mode report contract and is referenced by `mode-playbooks.md`. |

## Expected Behavior After Refactor

- Mode behavior remains unchanged for all eight modes.
- Universal preamble becomes genuinely universal.
- Validation scripts are named from one place.
- Iteration cap is single-sourced in `mode-playbooks.md`.
- Progressive-disclosure thresholds are single-sourced in `progressive-disclosure.md`.
- `init_skill.py` is reachable from Create mode.
- `SKILL.md` shrinks slightly without changing frontmatter or public mode contracts.

## Validation Checklist

1. Run `python3 skillify/scripts/quick_validate.py skillify`.
2. Run `python3 skillify/scripts/check_links.py skillify`.
3. Re-score against `references/validation-rubric.md`; Token Efficiency, Workflow Clarity, and Reusability must not regress.
4. Sweep `references/anti-patterns.md`, especially duplicated cross-tier content.
5. Sweep `references/security-checklist.md`; no new external HTTP, secret access, destructive command, or broad permission may appear.
6. Smoke-test Create routing to ensure trigger-prompt collection and `init_skill.py` wiring are reachable.
7. Smoke-test Refactor and Audit routing to ensure they reach `mode-playbooks.md` and Validation gate pointers resolve.
8. Confirm final source and installed `SKILL.md` match.

## Proposed Patch Notes

- Edit only `skillify/SKILL.md` and `skillify/references/mode-playbooks.md`.
- Do not edit frontmatter.
- Do not edit platform install scripts, deployment docs, validation scripts, or templates.
- Apply the default `init_skill.py` decision: wire it into Create rather than removing it from the script manifest.
- Mirror the final `SKILL.md` and `mode-playbooks.md` to `~/.codex/skills/skillify/`.

> 2026-07-05: extended by `recommendation-2026-07.md` (2026-07 rerun — corpus refreshed to 68 clean-markdown sources; see consult-records/ for cross-model provenance).
