# User story — <short title>   ·   ID: FR-<n>

**As a** <role>
**I want** <capability>
**so that** <benefit / why it matters>

- **Priority (MoSCoW):** Must | Should | Could | Won't
- **Traceability:** supports goal <G-n>
- **Notes / assumptions:** (mark unknowns as open questions, not facts)

## Acceptance criteria
Pick the format that fits (see references/requirements-techniques.md). Aim for 3–6 crisp, testable criteria;
more than ~8 means split the story.

### Option A — Given / When / Then (behavior with branches)
- **Given** <initial context> **When** <action> **Then** <observable outcome>
- **Given** <error/edge context> **When** <action> **Then** <error handling>

### Option B — Rule checklist (validation / UI / non-functional)
- [ ] <observable, testable condition>
- [ ] <edge case condition>
- [ ] <error condition>

## Edge cases & error paths
- <what happens on failure / retry / partial completion>

## INVEST self-check
- [ ] **I**ndependent — can ship without waiting on another story
- [ ] **N**egotiable — describes intent, not a locked solution
- [ ] **V**aluable — delivers value to the named role
- [ ] **E**stimable — enough is known to size it (else split off a Spike)
- [ ] **S**mall — fits comfortably in one iteration
- [ ] **T**estable — the AC above can pass/fail objectively

> If it fails Small/Valuable/Testable, split it (Spike / Path / Interface / Data / Rules).
