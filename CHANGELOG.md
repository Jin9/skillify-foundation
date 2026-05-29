# Changelog

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
