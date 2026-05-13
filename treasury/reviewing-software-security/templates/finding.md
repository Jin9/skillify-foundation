### [SEV-N] <short imperative title>

**Severity:** Critical | High | Medium | Low
**Confidence:** High | Medium | Low [needs verification?]
**Category:** <one primary taxonomy area>
**Secondary tags:** <optional adjacent taxonomy areas, or None>
**Standards:** CWE-###, API#:2023, ASVS V#.#.#, NIST SSDF PW.#.# (as applicable)
**Affected area:** path/to/file.go:LL-LL (or yaml: spec.template.spec.containers[0].securityContext)
**Environment scope:** SIT | UAT | PRD | all

**Asset at risk:** <asset and data class>
**Trust boundary crossed:** <boundary>
**Threat (STRIDE + abuse case):** <STRIDE label plus financial-domain abuse case>
**Attack scenario:**
  1. <actor starts at entry point>
  2. <control gap is exercised>
  3. <security impact occurs>
**Risk:** <confidentiality | integrity | availability implications>
**Business impact:** <financial + regulatory + reputational impact>

**Evidence:**
```lang
<minimal cited snippet>
```

**Recommended fix:**
```lang
<corrected snippet, minimal diff>
```
Why this works: <why the control closes the gap>
Trade-offs: <compatibility, operational, or cost trade-off>
Migration: <required rollout or data migration, or "none">

**Safer pattern:** <library / middleware / control name @ version, last-checked CVE date>

**Validation:**
- Test: <test name + assertion that bad behavior is rejected>
- Negative test: <test name + assertion that legitimate behavior still works>
- Static: <tool + rule id>
- Manual: <command or review step + expected result>
- Regression guard: <policy / lint / hook>
- Observability: <log / metric / alert>
- Rollout: <SIT -> UAT -> PRD gate and rollback condition>

**Residual risk:** <one sentence>
