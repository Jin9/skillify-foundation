# Changelog

## 2026-05-28

### Added

- Expanded the treasury catalog from 36 committed skills to 93 top-level skills.
- Added new treasury skill folders across banking delivery, architecture, testing, orchestration, research, operations, and domain product-pack workflows.
- Added the ShopPilot `ecom-mvp-2026` domain product pack and supporting artifacts.

### Changed

- Refreshed the root README and `treasury/README.md` to use the 93-skill purpose-group catalog.
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
