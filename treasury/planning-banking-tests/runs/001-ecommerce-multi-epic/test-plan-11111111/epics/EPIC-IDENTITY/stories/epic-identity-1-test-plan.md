---
artifact_type: qa-story-plan
blocking_oqs:
  - OQ-03
  - OQ-09
coverage_status: partial
epic_id: EPIC-IDENTITY
story_id: EPIC-IDENTITY-1
test_case_count: 4
test_plan_id: TPLAN-EPIC-SHOPPILOT-MVP
---

# Story Test Plan — EPIC-IDENTITY-1

## Test Case Roster

| ID | Scenario | Type | Pyramid | Owner | Risk |
|---|---|---|---|---|---|
| TC-EPIC-IDENTITY-1-001 | New customer signs up successfully | e2e | e2e | SDET | high |
| TC-EPIC-IDENTITY-1-002 | Duplicate email signup is rejected | functional | integration | SDET | high |
| TC-EPIC-IDENTITY-1-003 | Weak password is rejected | security | integration | Security-tester | high |
| TC-EPIC-IDENTITY-1-004 | Signup audit event excludes PII payload | compliance | integration | Compliance-mapper | critical |

## Test Case Details

### TC-EPIC-IDENTITY-1-001 — New customer signs up successfully

- Scenario type: happy
- Test type: e2e
- Pyramid level: e2e
- Environment: ENV-IDENTITY-e2e
- Owner: SDET
- Reviewer role: qa-lead
- Risk: high
- Smoke subset: True
- Tags: happy, must, signup
- Compliance tags: PDPA-TH-2019

#### Expected Assertions
- [db_state] customer row exists with email=fixture.email, name=fixture.name, phone=fixture.phone
- [db_state] password_hash column populated; plaintext password NOT present in any column or log
- [response_body] response body does not echo password or password hash
- [routing] response routes to auto-login OR login page (assertion deferred to OQ-09 resolution)

### TC-EPIC-IDENTITY-1-002 — Duplicate email signup is rejected

- Scenario type: error
- Test type: functional
- Pyramid level: integration
- Environment: ENV-IDENTITY-integration
- Owner: SDET
- Reviewer role: qa-lead
- Risk: high
- Smoke subset: False
- Tags: error, signup, uniqueness
- Compliance tags: PDPA-TH-2019
- Depends on: TC-EPIC-IDENTITY-1-001

#### Expected Assertions
- [response_status] field-level error 'email already in use' returned
- [db_state] no new customer row created (count unchanged)

### TC-EPIC-IDENTITY-1-003 — Weak password is rejected

- Scenario type: error
- Test type: security
- Pyramid level: integration
- Environment: ENV-IDENTITY-integration
- Owner: Security-tester
- Reviewer role: security-reviewer
- Risk: high
- Smoke subset: False
- Tags: error, password-strength, signup, tbd-oq-03
- Compliance tags: PDPA-TH-2019

#### Expected Assertions
- [response_body] field-level error returned describing requirement without revealing internal heuristics (exact predicate TBD pending OQ-03)
- [db_state] no new customer row created

### TC-EPIC-IDENTITY-1-004 — Signup audit event excludes PII payload

- Scenario type: banking_grade_audit
- Test type: compliance
- Pyramid level: integration
- Environment: ENV-IDENTITY-integration
- Owner: Compliance-mapper
- Reviewer role: compliance-officer
- Risk: critical
- Smoke subset: False
- Tags: banking_grade_audit, pii-scrub, signup
- Compliance tags: PDPA-TH-2019, banking_grade_audit
- Depends on: TC-EPIC-IDENTITY-1-001

#### Expected Assertions
- [audit_log_field] audit event has event='customer.signup'
- [audit_log_field] audit event has actor='self'
- [audit_log_field] audit event has ts in ISO-8601
- [audit_log_field] audit event has before=null (signup is creation, no prior state)
- [audit_log_field] audit event has after={subject_id:<customer_id>} (subject_id surrogate only, no PII)
- [audit_log_field] audit event has reason='customer_self_registration'
- [audit_log_field] audit event has idem_key=request idempotency key
- [audit_log_negative] audit payload does NOT contain plaintext password, password hash, email, phone, or name

