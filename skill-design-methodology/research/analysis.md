---
authored_by: chinnawat.w (git author of record; GitHub Jin9)
phase_owner: Claude Opus 4.7
pipeline_phase: 3
status: accepted
provenance_dated: 2026-05-19
---

# Phase 3 Architecture Analysis

Claude Opus 4.7 rerun of Phase 3, replacing the prior Codex/GPT fallback. Grounded in the actual `skillify/` files as of 2026-05-09 (`SKILL.md` is 179 lines; references, templates, scripts, platforms, examples folders are populated).

## Strengths To Preserve

- **Frontmatter description is operationally tight.** 752 chars (well under 1024), 12 verbatim trigger phrases, 3 explicit negative triggers, output contract baked into the description, and `compatibility:` declared for five platforms. Hits every Gemini Phase 2 trigger-design rule.
- **Mode routing is explicit and disambiguated.** The 8-row Modes table maps verb -> use-when phrases -> output. Universal preamble step 2 then disambiguates by *required output* rather than verb, which is the right tie-breaker for overlapping requests like "review and fix."
- **Progressive disclosure is correctly tiered.** Create lives inline because it is the most-invoked mode; the other 7 defer to `references/mode-playbooks.md`. References are flat (one level deep), and `progressive-disclosure.md` codifies the 500/5000 rule that the body itself respects.
- **Validation gate is mechanical and bounded.** Four numbered gates, explicit fail-close on gate 1, three-pass iteration cap on gates 2-4, and an explicit "Review and Audit do not edit" carve-out. Addresses Gemini's "iteration caps" synthesis rule.
- **Output contract is enforced both by prose and by script.** `quick_validate.py` blocks frontmatter that is missing trigger markers, contains XML angle brackets, or violates the banned-docs list (`README.md`, `CHANGELOG.md`, `INSTALLATION_GUIDE.md`, `QUICK_REFERENCE.md`, `CONTRIBUTING.md`). The Constraints section lists the same 5 names so a human reviewer hits the rule before the script does.
- **Per-mode Completion report table forces structured handoff.** Eight rows, each one line. Cheap to load, hard to skip.
- **Reference table at the end is a route map, not a glossary.** Each row is "when you need X, read Y" - explicit progressive-disclosure pointers rather than topic descriptions.
- **Constraints encode the meta-skill scope as hard prohibitions.** Rule 1 ("DO NOT generate the target artifact") is the canonical statement of `skillify`'s purpose; rule 2 protects repo-policy files; rules 3-7 cover Review/Audit immutability, trigger invention, body length, cross-tier duplication, and banned docs.

## Defects To Fix

1. **Validation gate step 1 re-enumerates what `quick_validate.py` checks.** Lines 99 say "Frontmatter must parse, `name` must match the folder, `description` must be under 1024 characters, frontmatter must not contain XML angle brackets, banned human-facing docs must be absent, and local links must resolve." This duplicates the script's behavior; if the script changes, `SKILL.md` will silently drift.
   - **Belongs in:** `SKILL.md` (shorten to "Run `scripts/quick_validate.py` and `scripts/check_links.py`. If they exit non-zero, surface the error verbatim and stop.") plus a comment in `scripts/quick_validate.py` that lists the checks. Single source of truth = the script.

2. **Universal preamble step 4 contains a Create-only rule.** "Create: collect the intended task, target users, and at least 3 concrete user prompts that should trigger the skill. If the user provided fewer than 3, ask once; do not invent trigger phrases." This rule does not apply to Refactor/Review/Audit/Compress/Split/Merge/Adapt.
   - **Belongs in:** `SKILL.md` (move to Create workflow step 1; leave universal preamble step 4 as the path-elicitation half only). Net: clearer mode separation, no line-count change.

3. **`platforms/` is leaked into the Output format directory tree.** Lines 110-119 show `platforms/` as part of the standard skill folder shape. `platforms/` is unique to `skillify` itself - generic skills produced by Create mode should not get a `platforms/` directory by default.
   - **Belongs in:** `SKILL.md` (drop `platforms/` from the tree; or label it "skillify-only").

4. **Iteration cap (3) is duplicated between `SKILL.md` gate-2-to-4 prose and `references/mode-playbooks.md` cross-mode rules.** Two homes for the same number. If the cap moves to 2 or 4, both files must update; one will lag.
   - **Belongs in:** `references/mode-playbooks.md` (canonical home, since it owns cross-mode rules). `SKILL.md` should say "iterate per the cap in `mode-playbooks.md`" rather than restating "three passes."

5. **Constraints rule "move optional or mode-specific detail over roughly 50 lines to `references/`" introduces an undocumented threshold.** The 50-line threshold conflicts with the 500-line file-level rule in `progressive-disclosure.md`. The 50-line rule is per-section, not per-file, but the wording does not say so, and no script enforces it.
   - **Belongs in:** `references/progressive-disclosure.md` (clarify the per-section vs per-file thresholds in one place); `SKILL.md` should drop the 50-line number and point to the reference instead.

6. **Modes section trailer re-enumerates 7 mode names.** "Run Create mode from this file. For Refactor, Review, Audit, Compress, Split, Merge, and Adapt, read `references/mode-playbooks.md`..." The Modes table immediately above lists those 7 already.
   - **Belongs in:** `SKILL.md` (compress to "Run Create from this file. For any other mode in the table above, read `references/mode-playbooks.md` first.").

7. **`init_skill.py` is named in the Templates and scripts manifest but never invoked from any workflow step.** Either the Create workflow should call it explicitly, or the manifest should drop it.
   - **Belongs in:** `SKILL.md` Create workflow step 0, e.g. "Run `scripts/init_skill.py <skill-name>` to scaffold the folder if no target exists." (Cheaper than removing the script.)

8. **References table row `platforms/deployment-guide.md` overlaps `references/platform-compatibility.md`.** Two pointers for the platform topic from a single table risks ambiguity at lookup time.
   - **Belongs in:** `SKILL.md` (drop the `platforms/deployment-guide.md` row from the References table; users still find it via the Output format tree if `platforms/` stays).

## Contradictions Or Weak Assumptions

- **The existing Codex `analysis.md` claims "`SKILL.md` has a 'When to use this skill' section that repeats frontmatter trigger logic."** No such H2 exists in the current `SKILL.md`. The Codex finding either references an older revision or is hallucinated. **Reject this defect; it is not a real condition.**
- **Codex `analysis.md` defect "Validation gate wording is dense" is overstated.** The gate is four short numbered steps. The genuine issue is gate-1 enumerating script behavior (defect #1 above), not density.
- **Codex `analysis.md` defect "Completion report table is verbose" is wrong.** The table is 8 one-line rows; removing it would push duplication into 8 mode files. **Keep the table.**
- **Gemini synthesis "Approval Gates for sensitive state changes" applies to skills `skillify` produces, not to `skillify` itself.** `skillify`'s own operations are local file writes inside skill folders - low blast radius. The advice is appropriate for `references/security-checklist.md` (where it already lives implicitly) but does not belong in `SKILL.md`. **Reject as a SKILL.md change.**
- **Gemini synthesis "Standardized Fallbacks for missing external binaries" already passes.** Skillify scripts only need Python and PyYAML. **No skillify change required.**
- **Gemini synthesis "No Deep Nesting", "500/5000 Rule", "On-Demand Loading", and "Deduplicate Cross-Tier Content" are confirmations of existing rules**, not new requirements. They map to `progressive-disclosure.md` (lines 25, 46-50, 107) and `anti-patterns.md` #11. **No change.**
- **Phase 1 Gemini suggested optional license/metadata fields from community catalogs.** The current frontmatter has `name`, `description`, `compatibility`. Adding license/version/author would expand the validation surface (`quick_validate.py`, the rubric, the security checklist) without any official platform requiring them. **Reject; agreed with Codex on this one.**
- **The mode count (8) is broad enough to be a scope risk but not yet a split trigger.** All 8 modes operate on a single artifact (skill folders) and share the same validation gate, anti-pattern sweep, and security checklist. Splitting now would force the user to remember which CLI surface to invoke for which verb. **Hold; revisit if a 9th mode is proposed.**

## Scope Boundaries

Hard rejections that any Phase 4 recommendation must respect:

- `skillify` produces skill artifacts only: `SKILL.md`, `references/`, `templates/`, `scripts/`, `assets/`, `examples/`, and review/audit reports. **It must not generate the downstream artifact a requested skill would later produce** (anti-pattern #7).
- **No edits to `AGENTS.md`, `CLAUDE.md`, `.github/copilot-instructions.md`, `GEMINI.md`, or other repo-policy files** unless the user explicitly asks for platform-policy adaptation. Adapt mode is the only authorized path, and it still does not auto-edit root rules.
- **No one-off slash commands or system prompts.** If a request would not become a reusable skill, decline and route the user to the appropriate non-skillify surface.
- **No multi-model orchestration prescriptions in `SKILL.md`.** The cross-model plan in `skill-design-methodology/plan.md` is implementation scaffolding; the reusable skill must stay agent-agnostic.
- **No human-facing docs inside any skill folder.** Banned: `README.md`, `CHANGELOG.md`, `INSTALLATION_GUIDE.md`, `QUICK_REFERENCE.md`, `CONTRIBUTING.md`. Enforced by `quick_validate.py`.
- **No new mode without a Split.** Adding a 9th mode (e.g., "Test", "Document", "Publish") triggers Split mode against `skillify` itself, not addition.
- **No license/metadata frontmatter expansion** unless an official platform mandates it.

## Compaction Opportunities

Realistic, low-risk reductions to `SKILL.md`. Current 179 lines; target floor about 160 lines without losing operational hooks.

| Change | Source lines | Estimated saving | Risk |
|---|---|---|---|
| Drop the 7-mode-name list in the Modes section trailer; say "any other mode" instead | 43 | about 1 line | None |
| Move Create-only "collect 3 trigger phrases" rule from Universal preamble step 4 to Create workflow step 1 | 59 -> 67 | 0 lines (relocation) | None |
| Replace gate-1 in-line check enumeration with a one-line script invocation; document checks in the script comment | 99 | about 3 lines | Low - `quick_validate.py` becomes the authority |
| Drop `platforms/` from the Output format tree (or label as skillify-only) | 116 | about 1 line | None for downstream skills |
| Drop "roughly 50 lines" threshold from Constraints; point to `progressive-disclosure.md` | 144 | about 1 line | None - reference covers it |
| Drop `platforms/deployment-guide.md` row from References table | 168 | about 1 line | None - `references/platform-compatibility.md` covers it |
| Replace "iterate up to three passes in Create, Refactor, Compress, Split, Merge, and Adapt" with "iterate per the cap in `references/mode-playbooks.md`" | 104 | about 1 line | Low - single source of truth |

Total realistic compaction: 8-10 lines, taking `SKILL.md` to about 169-171. Do **not** pursue aggressive compaction - the file is already lean and most of the prose carries operational weight.

What **not** to compact:
- Per-mode Completion report table (8 rows) - load-bearing for handoff consistency.
- Constraints DO-NOT list - each rule has a distinct enforcement target.
- References table - the route map for progressive disclosure.
- Universal preamble steps 1-3 - disambiguation logic that prevents misrouted invocations.

## Validation Risks

- **`quick_validate.py` does not check the `compatibility:` frontmatter field against the `platforms/` directory contents.** Removing `platforms/codex.md` while leaving `compatibility: codex` in frontmatter would silently leave Adapt mode with broken pointers. Consider adding a cross-check in the script.
- **`check_links.py` only validates relative links from `SKILL.md`.** Links inside reference files (e.g., `references/mode-playbooks.md` referencing `templates/audit-report-template.md`) may not be in scope. Confirm before relying on it for full link integrity.
- **Anti-pattern sweep (gate 3) treats all 12 anti-patterns as creation-time gates.** Anti-pattern #10 (Stale Skills / Context Window Competition) only applies post-deployment and cannot be assessed during Create. Consider tagging anti-patterns with phase applicability (`creation`, `lifecycle`, `both`) or scoping the gate to the creation-applicable subset.
- **Iteration cap (3) is duplicated across `SKILL.md` and `mode-playbooks.md`.** Until defect #4 is fixed, any change to the cap must be applied in both files.
- **The 50-line per-section threshold in Constraints (line 144) is unenforced.** No script measures it. Until clarified, the rule depends on agent judgment, which is a soft validation surface.
- **The banned-docs list (5 named files) is enforced by `quick_validate.py` via filename match.** Adding a 6th banned doc requires updating both the script and the Constraints section. Consider extracting the list into a script constant that the section references.
- **`init_skill.py` is unreferenced from any workflow step**, so it is never executed in the Create path; gate 1 will not exercise it. Either wire it into the Create workflow or accept that it is a manual bootstrap tool.
- **Installed copy at `~/.codex/skills/skillify/SKILL.md` may drift from the repo copy.** The Phase 3 plan instructs Codex to keep the installed target as the primary. Final Phase 7 validation must run against both, or pick one canonical and document the sync method.

> 2026-07-05: extended by `analysis-2026-07.md` (2026-07 rerun — corpus refreshed to 68 clean-markdown sources; see consult-records/ for cross-model provenance).
