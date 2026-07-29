# Changelog

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
