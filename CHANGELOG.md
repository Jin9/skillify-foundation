# Changelog

## 2026-09-07

### Added

- Literature corpus refreshed with a 19-source frontier-model cohort for the Claude Fable 5.1 / GPT-6 Astra generation (68 -> 87 sources, ~249K words): the Fable 5.1, Fable 5, and Opus 5 prompting guides, the Fable 5.1 what's-new and migration guide, the cross-model prompting best-practices page, and four gap-fill references (effort, thinking, prompt caching, tool search) from Anthropic; OpenAI's GPT-6 Astra model guidance and launch announcement; the Codex models page for GPT-6 Astra; a 2026-09 re-capture of the Agent Skills specification; four Tier-2 pieces; and an attributed internal synthesis of the dated-prompt-pattern taxonomy (`research-vault/prompt_cruft_taxonomy_frontier_models.md`). No existing capture was re-fetched (scope decision); the 27 April-dated files stay flagged in the source index.
- `skillify/references/model-generation-fit.md`: calibration principles for the current model generation (specificity matched to fragility, instruction priority and stops, verification scope, delegation, effort tiers instead of model names, reporting and style, memory, the keep list, a dated generation profile).
- `skillify/templates/operating-contract.md`: the canonical instruction-priority / autonomy / stop-conditions / verification / delegation / progress / model-cost-tier block that every multi-step skill now carries; embedded (filled) in the three skill templates and the `init_skill.py` scaffold.
- `skillify/scripts/cruft_scan.py`: stdlib advisory scanner for dated patterns (pressure language, thinking and show-your-reasoning scaffolds, narration cadences and numeric caps, unreasoned prohibition walls, narration suppressors, anti-formatting rules, pinned model names, hardcoded paths, missing operating contract or tier) with keep-list exemptions for trigger sections, table rows, fenced code, and quoted mentions; `--strict` is the exit check of the new Re-baseline sub-flow, `--self-test` proves each signal with a hit and a control.
- Refactor sub-flow "Re-baseline for a model generation" in `references/mode-playbooks.md`; lifecycle rows for "Model release" and "Unrequested pause or divergence"; workflow Pattern 0 (goal and constraints) plus Cross-Cutting Techniques for the operating contract, delegation and parallelism, and per-node model-cost tier / effort hints.

### Changed

- `skillify/SKILL.md`: one clarification budget replaces three independent "ask once" rules; the degree-of-freedom ladder chooses by fragility (goal + reasoned constraints for judgment work, numbered steps only where order matters, scripts for fragile operations); gate 1 fixes-and-re-runs within the iteration cap instead of an unsatisfiable "do not write target files"; the eight-bullet "DO NOT" wall became six reasoned constraints; skillify carries its own filled operating contract and per-section model-cost tiers; new trigger "re-baseline this skill for a new model".
- Rubric dimension 3 rewards specificity matched to fragility and an operating contract rather than numbered steps as such; dimension 5 names scaffolds and repeated reminders; dimension 6 covers stalls and user precedence. Anti-patterns stay at twelve: #5 is now "Wrong Degree of Freedom" (two-directional), #6 absorbs instruction priority and unrequested pauses, #12 is "Hardcoded Platform or Model Assumptions" (pinned names, retired-model scaffolds). Audit template gains a Target line, a Cruft Scan section, and two security lines.
- Mode playbooks: "confirm the chain" and the Split/Merge approval pauses are gone (state the plan in the opening line and proceed; pause only for an unauthorised in-place overwrite); Audit states the host and model generation it assumes; Compress drops restatements of what the agent does unprompted; Adapt turns `model:`/`effort:` pins into body-level tiers.
- Templates: MCP template is host-neutral (no "Settings > Extensions" paths) and declares all five hosts; constraints carry their reasons; "Validation checklist" is now "Verification" (name the evidence you will quote).
- Security checklist: instruction-priority / stop-behaviour and reasoning-extraction checks, two threat-table rows, two review-prompt questions.
- Eval harness (`skill-design-methodology/evals/`): the "before" baseline is now the whole skill folder at the git ref; twelve frontier-fit indicators and eight golden prompts (incl. two skill-conflict probes) added; 20/32 indicators and 16/20 prompts at HEAD -> 32/32 and 20/20 after.
- Corpus bookkeeping synced to 87 sources in root `README.md`, `CLAUDE.md`, `AGENTS.md`, the literature review, and the source index; the never-changelogged 2026-07-05 refresh is back-filled below.

## 2026-07-30

### Fixed

- `delegating-to-cli-models`: repaired the agy and codex dispatch contracts, which had stopped working as documented. `run_agy.sh` no longer gates on `agy models` — that call is server-backed, was observed hanging indefinitely, ran before the watchdog was armed, and now emits slugs rather than the display labels the guard whole-line-matched, so it rejected every documented model. The wrapper now accepts either identifier form, writes a per-run `--log-file` so concurrent fan-out dispatches stop racing for the newest shared log, keeps agy's own response timeout under the watchdog (it defaults to 5m, so a longer watchdog could never fire), and exits 3 when the backend model differs from the request. `run_codex.sh` stops passing `-m` unless a model is explicitly pinned — `codex exec` resolves `config.toml` itself, and an unconditional `-m` silently defeated `--profile` and `-c model=` layering — and reads the resolved model back from codex's own run header.

### Changed

- Refreshed `delegating-to-cli-models` model guidance for agy 1.1.8 and codex-cli 0.146.0: `gpt-5.5` → `gpt-5.6-sol`, Flash tier → 3.6, and tier-first routing so concrete model identifiers live in exactly one dated table. Corrected two claims disproved by live dispatch: argument order is *not* what drops `--model` (a controlled 2×2 showed `--model` after `--print` is honored), and agy now has a separate `--effort` flag that hard-conflicts with an effort-encoded model identifier. Verified 2026-07-30 — `gemini-3.1-pro-high` is accepted with exit 0 and no warning while the backend silently runs the default model, so the display-label form is now the documented default and post-run verification is mandatory rather than advisory.

## 2026-07-19

### Added

- Promoted `principal-advisor` into the treasury (130 → 131): a conversational principal/staff-level advisory sparring partner for chat hosts — advise / assess (BUILDABLE / NOT NOW / PHASED) / challenge (steelman-first, evidence-based findings) / brainstorm stances plus on-request in-chat wrap-up templates, a depth dial, and a conditional Thailand-fintech lens (BOT/PDPA/AMLA). Text-only (no scripts), vendor-neutral, packaged for claude.ai zip upload. Verified by both deterministic validators, an independent 49/50 rubric audit, a 20/20 blind trigger eval, and an adjudicated codex (gpt-5.6-sol, max reasoning effort) cross-model consult.

### Changed

- Synced treasury counts to 131 across `treasury/README.md` (summary, group, catalog, and index tables), root `README.md`, `AGENTS.md`, and `CLAUDE.md`.

## 2026-07-18

### Added

- Promoted `launching-workflow-routines` into the treasury (129 → 130): a by-name launcher over a `routines/` registry of designed multi-workflow routines — resolve, script-validate, overview-first plan, gated node-by-node dispatch, script-verified artifacts, run report. Ships `validate_routine.py` / `check_run.py`, the routine-format and run-protocol references, templates, and a demo routine.
- Added the root `routines/` registry with three starter routines (`cross-model-authoring`, `research-squad-chain`, `banking-ba-wrap`) plus `INDEX.md`.

### Changed

- Synced treasury counts to 130 across `treasury/README.md`, root `README.md`, `AGENTS.md` (was stale at 122), and `CLAUDE.md` (was stale at 123).

## 2026-07-05

### Changed

- Literature corpus refreshed to 68 clean-markdown sources (~168K words): every capture normalized to the standard provenance header by `skill-design-methodology/tools/normalize_capture.py`, four former Cloudflare-stub captures recovered, the load-bearing Anthropic and agentskills.io docs re-captured as served markdown, MCP specification (2025-11-25) and Claude Code plugin docs added, and an eight-report `research-vault/` group harvested. `verify_corpus.py` (38 checks) became the corpus acceptance gate. Skillify patched in 15 files for platform-fact drift (Antigravity `antigravity-cli` paths, install-doc link check, three trigger placeholders in scaffolds); rubric 49/50. (Back-filled 2026-09-07: this refresh landed in commits 0b21e45 and fd440ec but was never changelogged.)

## 2026-05-30

### Changed

- Ran the `pass45` audit: brought all 90 treasury skills to ≥45/50 on the validation rubric (every dimension ≥4/5), refactoring 15 skills by moving body bloat into `references/` with no splits (merged PR #16).
- Expanded `implementing-go-template-requirements` with a fuller service-test template, hardened recipe/naming/testing references, and tightened the Go scaffold templates.

### Fixed

- `reporting-research-run` step-1 artifact scan now lists `01-grounding_pack.json`.
- Reconciled the `treasury/README.md` Extra-assets column with the files actually tracked in each skill folder.

## 2026-05-29

### Changed

- Sharpened `skillify`'s scope boundary and added portability/provenance guidance; extended the validation scripts.
- Defaulted the skillify installer to `claude` / `.agents` / `antigravity-cli` targets and made Copilot install opt-in; dropped `platforms/` from the copy set.
- Treasury hygiene pass: merged the redundant Go-service platform pair, compressed the banking-brief skill, and hardened several skills (merged PRs #13–#15).

### Added

- Added new literature source documents to the research corpus.

### Fixed

- Repaired `quick_validate.py` validator false-positives and brought 8 skills to a gate-1 pass.
- Cleared gate-1 / rubric defects across treasury skills.

### Removed

- Retired the duplicate `model-selection-updated` skill.

## 2026-05-28

### Added

- Expanded the treasury catalog from 36 committed skills to 92 top-level skills.
- Added new treasury skill folders across banking delivery, architecture, testing, orchestration, research, and operations workflows.

### Changed

- Refreshed the root README and `treasury/README.md` to use the 92-skill purpose-group catalog.
- Renamed focused treasury skills for clearer action-oriented names:
  - `business-analysis-flow` -> `running-business-analysis-workflow`
  - `requirement-analysis` -> `scoping-technical-requirements`
  - `standardization` -> `defining-engineering-standards`
  - `clean-go-service` -> `refactoring-go-services`
  - `generate-gherkin-ac` -> `generating-gherkin-acceptance-criteria`
  - `langgraph-professional` -> `developing-langgraph-workflows`
  - `panel-open` -> `opening-debate-panel`
  - `report-run` -> `reporting-research-run`
- Updated cross-skill references and catalog links for the renamed folders while preserving workflow stage IDs.
- Normalized selected skill frontmatter metadata into `metadata:` blocks.

### Fixed

- Added valid YAML frontmatter to the Go service refactoring skill.
- Added `.gitignore` coverage for timestamped treasury backup directories.

### Removed

- Removed the ShopPilot `ecom-mvp-2026` domain product pack from the treasury.
