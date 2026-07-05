---
authored_by: chinnawat.w (git author of record; GitHub Jin9)
phase_owner: Codex baseline + Gemini synthesis
pipeline_phase: 0
status: accepted
provenance_dated: 2026-05-19
---

# Research Notes

## Phase 0 Baseline Inventory - Codex/GPT Rerun

Date: 2026-05-09

### Target

- Repository source skill: `skillify/`
- Installed skill target: `~/.codex/skills/skillify/`
- Primary file: `SKILL.md`
- Current source and installed `SKILL.md` are identical.
- Current `SKILL.md` line count: 179.

### Current Skill Boundary

`skillify` is a meta-skill for creating, refactoring, reviewing, auditing, compressing, splitting, merging, and adapting reusable agent skills. It must produce skill folders, `SKILL.md` files, references, templates, scripts, assets, examples, or review/audit reports. It must not generate the downstream artifact a requested skill would later produce. It must not edit repo-policy files unless the user explicitly asks for platform-policy adaptation.

### Existing Support Files

- References: `anti-patterns.md`, `frontmatter-guide.md`, `lifecycle-and-iteration.md`, `mode-playbooks.md`, `platform-compatibility.md`, `progressive-disclosure.md`, `security-checklist.md`, `validation-rubric.md`, `workflow-patterns.md`
- Templates: `audit-report-template.md`, `basic-skill-template.md`, `domain-skill-template.md`, `mcp-skill-template.md`
- Scripts: `check_links.py`, `init_skill.py`, `quick_validate.py`
- Platforms: `agents.md`, `antigravity.md`, `claude.md`, `codex.md`, `copilot-instructions.md`, `copilot.md`, `deployment-guide.md`, `gemini.md`, `install.sh`
- Examples: `create-from-scratch.md`, `good-description-examples.md`, `skill-audit-walkthrough.md`

### Stable Constraints

- Preserve unrelated user edits. Current dirty files before this rerun: `.gitignore`, `skillify/SKILL.md`, `skillify/platforms/deployment-guide.md`, `skillify/platforms/install.sh`, plus untracked `file.txt`.
- Do not create `README.md`, `CHANGELOG.md`, `INSTALLATION_GUIDE.md`, `QUICK_REFERENCE.md`, `CONTRIBUTING.md`, or other human-facing docs inside the skill folder.
- Keep permanent workflow artifacts inside `skill-design-methodology/research/` or the target skill folder only.
- Final target is `~/.codex/skills/skillify/SKILL.md`, with optional focused references only when they reduce load.
- Preserve existing scripts and templates unless a model phase identifies a concrete defect.

### Baseline Assessment

Strengths to preserve:

- Clear meta-skill boundary and explicit negative triggers.
- Concrete mode routing for Create, Refactor, Review, Audit, Compress, Split, Merge, and Adapt.
- Create mode is self-contained in `SKILL.md`; non-create modes point to `references/mode-playbooks.md`.
- Deterministic validation scripts exist for frontmatter/structure and local links.
- Reference files are one-level deep and mostly have clear ownership.

Risks for model phases:

- The skill is broad by mode count, so later phases should challenge whether mode routing is enough to prevent over-triggering.
- Completion-report and validation instructions are inline; later phases should decide whether that is useful immediate context or reference bloat.
- Platform compatibility is broad; later phases should avoid vendor-specific assumptions in the core body.

> 2026-07-05: extended by `research-notes-2026-07.md` (2026-07 rerun — corpus refreshed to 68 clean-markdown sources; see consult-records/ for cross-model provenance).
