---
artifact_type: qa-story-plan
blocking_oqs:
  - OQ-20
coverage_status: partial
epic_id: EPIC-GOVERNANCE-OBS
story_id: EPIC-GOVERNANCE-OBS-3
test_case_count: 6
test_plan_id: TPLAN-EPIC-SHOPPILOT-MVP
---

# Story Test Plan — EPIC-GOVERNANCE-OBS-3

## Test Case Roster

| ID | Scenario | Type | Pyramid | Owner | Risk |
|---|---|---|---|---|---|
| TC-EPIC-GOVERNANCE-OBS-3-001 | Admin views audit row with before/after diff | integration | integration | SDET | medium |
| TC-EPIC-GOVERNANCE-OBS-3-002 | Audit-log payload never contains forbidden fields (banking-grade audit) | security | integration | SDET | critical |
| TC-EPIC-GOVERNANCE-OBS-3-003 | Customer cannot reach audit-log endpoint (authz) | security | integration | SDET | critical |
| TC-EPIC-GOVERNANCE-OBS-3-004 | Cross-cutting: every state-changing event across all epics emits an audit event with 7 fields | integration | integration | SDET | critical |
| TC-EPIC-GOVERNANCE-OBS-3-005 | Cross-cutting: every customer-facing notification logged with notification_logged=true | integration | integration | SDET | high |
| TC-EPIC-GOVERNANCE-OBS-3-006 | Cross-cutting: PII redaction across log volume per pii_inventory masking rules | compliance | integration | SDET | critical |

## Test Case Details

### TC-EPIC-GOVERNANCE-OBS-3-001 — Admin views audit row with before/after diff

- Scenario type: happy
- Test type: integration
- Pyramid level: integration
- Environment: env-gov-integration
- Owner: SDET
- Reviewer role: QA-Lead
- Risk: medium
- Smoke subset: True
- Tags: audit-read, happy
- Compliance tags: CCA-TH-LOG-90D

#### Expected Assertions
- [ui_render] Audit row renders with actor, action='product.price_changed', entity_id, ts, before={price:1290}, after={price:999}, reason.

### TC-EPIC-GOVERNANCE-OBS-3-002 — Audit-log payload never contains forbidden fields (banking-grade audit)

- Scenario type: banking_grade_audit
- Test type: security
- Pyramid level: integration
- Environment: env-gov-integration
- Owner: SDET
- Reviewer role: Compliance-Lead
- Risk: critical
- Smoke subset: True
- Tags: AP-Q12, banking_grade_audit, scrubbing
- Compliance tags: CCA-TH-LOG-90D, PDPA-TH-2019

#### Expected Assertions
- [scrub_invariant] No row in the audit corpus contains a regex match for 'password', 'password_hash', 'token', 'card_number', or PAN-shaped digit sequences in any payload field.
- [fail_closed] If any forbidden match is detected by the self-test scrubber, the export fails closed and an alert is raised.
- [audit_field_completeness] Every sampled audit row enumerates all 7 fields (event, actor, ts, before, after, reason, idem_key) (AP-Q12).

### TC-EPIC-GOVERNANCE-OBS-3-003 — Customer cannot reach audit-log endpoint (authz)

- Scenario type: banking_grade_authz
- Test type: security
- Pyramid level: integration
- Environment: env-gov-integration
- Owner: SDET
- Reviewer role: Compliance-Lead
- Risk: critical
- Smoke subset: True
- Tags: audit-read, banking_grade_authz
- Compliance tags: CCA-TH-LOG-90D

#### Expected Assertions
- [http_status] GET /admin/audit by non-admin returns HTTP 403 or 404.
- [audit_event_emitted] authz_denied audit event emitted with all 7 standard fields.

### TC-EPIC-GOVERNANCE-OBS-3-004 — Cross-cutting: every state-changing event across all epics emits an audit event with 7 fields

- Scenario type: banking_grade_audit
- Test type: integration
- Pyramid level: integration
- Environment: env-gov-integration
- Owner: SDET
- Reviewer role: Compliance-Lead
- Risk: critical
- Smoke subset: False
- Tags: AP-Q12, banking_grade_audit, cross-cutting
- Compliance tags: CCA-TH-LOG-90D

#### Expected Assertions
- [cross_event_completeness] For each enumerated state-changing event class (catalog.*, order.*, inventory.adjusted, review.created, review.hidden, identity.*, payment.*, fulfillment.*), at least one audit row is emitted per occurrence.
- [audit_field_completeness] Every sampled audit row enumerates all 7 fields (event, actor, ts, before, after, reason, idem_key) (AP-Q12).

### TC-EPIC-GOVERNANCE-OBS-3-005 — Cross-cutting: every customer-facing notification logged with notification_logged=true

- Scenario type: banking_grade_notification
- Test type: integration
- Pyramid level: integration
- Environment: env-gov-integration
- Owner: SDET
- Reviewer role: Privacy-Lead
- Risk: high
- Smoke subset: False
- Tags: cross-cutting, notification-log
- Compliance tags: PDPA-TH-2019

#### Expected Assertions
- [notification_log_record] For each customer-facing notification class (order updates, signup confirm, password reset, review-status), a notification-log row is recorded with notification_logged=true.
- [no_pii_leak] Notification-log payload masks PII per pii_inventory[*].masking_rule — email partial-masked, phone partial-masked outside fulfilment context, password/token never appear.

### TC-EPIC-GOVERNANCE-OBS-3-006 — Cross-cutting: PII redaction across log volume per pii_inventory masking rules

- Scenario type: regulatory
- Test type: compliance
- Pyramid level: integration
- Environment: env-gov-integration
- Owner: SDET
- Reviewer role: Privacy-Lead
- Risk: critical
- Smoke subset: False
- Tags: AP-Q7, cross-cutting, pii-redaction
- Compliance tags: PDPA-TH-2019

#### Expected Assertions
- [redaction_invariant] Synthetic fixtures (AP-Q7: no real PII) carry well-known marker tokens; log scan must detect zero plaintext markers for password and session_token (always-masked rule), and the configured partial-mask pattern for phone outside fulfilment surfaces.
- [redaction_invariant] shipping_address line is scrubbed from authz_denied audit payloads per pii_inventory.shipping_address.masking_rule.

