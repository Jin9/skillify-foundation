---
artifact_type: qa-story-plan
blocking_oqs:
  - OQ-02
  - OQ-04
  - OQ-11
coverage_status: partial
epic_id: EPIC-IDENTITY
story_id: EPIC-IDENTITY-2
test_case_count: 5
test_plan_id: TPLAN-EPIC-SHOPPILOT-MVP
---

# Story Test Plan — EPIC-IDENTITY-2

## Test Case Roster

| ID | Scenario | Type | Pyramid | Owner | Risk |
|---|---|---|---|---|---|
| TC-EPIC-IDENTITY-2-001 | Valid credentials authenticate the customer | e2e | e2e | SDET | critical |
| TC-EPIC-IDENTITY-2-002 | Wrong password yields generic error (no enumeration) | security | integration | Security-tester | critical |
| TC-EPIC-IDENTITY-2-003 | Unknown email yields the same generic error | security | integration | Security-tester | critical |
| TC-EPIC-IDENTITY-2-004 | Repeated failed logins trigger throttle | security | integration | Security-tester | critical |
| TC-EPIC-IDENTITY-2-005 | Login success audit event (banking_grade_audit derived) | compliance | integration | Compliance-mapper | critical |

## Test Case Details

### TC-EPIC-IDENTITY-2-001 — Valid credentials authenticate the customer

- Scenario type: happy
- Test type: e2e
- Pyramid level: e2e
- Environment: ENV-IDENTITY-e2e
- Owner: SDET
- Reviewer role: qa-lead
- Risk: critical
- Smoke subset: True
- Tags: happy, login, must, session
- Compliance tags: PDPA-TH-2019
- Depends on: TC-EPIC-IDENTITY-1-001

#### Expected Assertions
- [response_header] response sets a session cookie/token bound to customer_id
- [routing] customer routed to account home
- [audit_log_field] audit event 'customer.login_success' emitted with actor=customer_id, ts, no password/hash in payload

### TC-EPIC-IDENTITY-2-002 — Wrong password yields generic error (no enumeration)

- Scenario type: error
- Test type: security
- Pyramid level: integration
- Environment: ENV-IDENTITY-integration
- Owner: Security-tester
- Reviewer role: security-reviewer
- Risk: critical
- Smoke subset: False
- Tags: anti-enumeration, banking_grade_authz, error, login
- Compliance tags: PDPA-TH-2019
- Depends on: TC-EPIC-IDENTITY-1-001

#### Expected Assertions
- [response_body] response returns exact generic message 'อีเมลหรือรหัสผ่านไม่ถูกต้อง'
- [response_body_negative] response does NOT differentiate 'email unknown' from 'password mismatch'
- [audit_log_field] audit event 'customer.login_failure' emitted with reason='credential_mismatch'

### TC-EPIC-IDENTITY-2-003 — Unknown email yields the same generic error

- Scenario type: error
- Test type: security
- Pyramid level: integration
- Environment: ENV-IDENTITY-integration
- Owner: Security-tester
- Reviewer role: security-reviewer
- Risk: critical
- Smoke subset: False
- Tags: anti-enumeration, banking_grade_authz, error, login
- Compliance tags: PDPA-TH-2019

#### Expected Assertions
- [response_body] response returns identical generic message as TC-EPIC-IDENTITY-2-002
- [response_body_negative] response does NOT reveal that the email is unknown
- [response_timing] response timing is statistically indistinguishable from wrong-password case (no enumeration via side-channel)

### TC-EPIC-IDENTITY-2-004 — Repeated failed logins trigger throttle

- Scenario type: banking_grade_authz
- Test type: security
- Pyramid level: integration
- Environment: ENV-IDENTITY-integration
- Owner: Security-tester
- Reviewer role: security-reviewer
- Risk: critical
- Smoke subset: False
- Tags: banking_grade_authz, brute-force, tbd-oq-04, throttle
- Compliance tags: PDPA-TH-2019, banking_grade_authz
- Depends on: TC-EPIC-IDENTITY-2-002

#### Expected Assertions
- [response_status] after threshold failures (threshold TBD pending OQ-04) further attempts on the same email are throttled or require a challenge
- [audit_log_field] audit event 'authz_denied' emitted with reason='throttle' and actor=request_origin
- [config] throttle parameters loaded from configuration, not hard-coded literals

### TC-EPIC-IDENTITY-2-005 — Login success audit event (banking_grade_audit derived)

- Scenario type: banking_grade_audit
- Test type: compliance
- Pyramid level: integration
- Environment: ENV-IDENTITY-integration
- Owner: Compliance-mapper
- Reviewer role: compliance-officer
- Risk: critical
- Smoke subset: False
- Tags: banking_grade_audit, login, pii-scrub
- Compliance tags: PDPA-TH-2019, banking_grade_audit
- Depends on: TC-EPIC-IDENTITY-2-001

#### Expected Assertions
- [audit_log_field] audit event has event='customer.login_success'
- [audit_log_field] audit event has actor=customer_id (surrogate)
- [audit_log_field] audit event has ts ISO-8601
- [audit_log_field] audit event has before=null
- [audit_log_field] audit event has after={session_issued:true}
- [audit_log_field] audit event has reason='credential_match'
- [audit_log_field] audit event has idem_key from request
- [audit_log_negative] audit payload does NOT contain password, password hash, session token value, or email

