---
artifact_type: qa-story-plan
blocking_oqs:
  - OQ-19
coverage_status: partial
epic_id: EPIC-GOVERNANCE-OBS
story_id: EPIC-GOVERNANCE-OBS-2
test_case_count: 9
test_plan_id: TPLAN-EPIC-SHOPPILOT-MVP
---

# Story Test Plan — EPIC-GOVERNANCE-OBS-2

## Test Case Roster

| ID | Scenario | Type | Pyramid | Owner | Risk |
|---|---|---|---|---|---|
| TC-EPIC-GOVERNANCE-OBS-2-001 | Customer writes a review for a delivered order line | integration | integration | SDET | medium |
| TC-EPIC-GOVERNANCE-OBS-2-002 | Second review on same order line is rejected | functional | integration | SDET | medium |
| TC-EPIC-GOVERNANCE-OBS-2-003 | Review rejected for non-delivered order | functional | integration | SDET | medium |
| TC-EPIC-GOVERNANCE-OBS-2-004 | Rating outside 1-5 is rejected | functional | unit | SDET | low |
| TC-EPIC-GOVERNANCE-OBS-2-005 | Admin soft-hides a review with reason and review is not edited | integration | integration | SDET | high |
| TC-EPIC-GOVERNANCE-OBS-2-006 | Banking-grade authz: only admins can hide reviews; only owning customer can submit | security | integration | SDET | high |
| TC-EPIC-GOVERNANCE-OBS-2-007 | Banking-grade tipping-off: moderation rejection copy does not disclose detection rule | compliance | integration | SDET | high |
| TC-EPIC-GOVERNANCE-OBS-2-008 | Notification log: review-status change to customer is logged with notification_logged=true | integration | integration | SDET | medium |
| TC-EPIC-GOVERNANCE-OBS-2-009 | Idempotency: duplicate review-submit with same idem_key produces single row | integration | integration | SDET | medium |

## Test Case Details

### TC-EPIC-GOVERNANCE-OBS-2-001 — Customer writes a review for a delivered order line

- Scenario type: happy
- Test type: integration
- Pyramid level: integration
- Environment: env-gov-integration
- Owner: SDET
- Reviewer role: QA-Lead
- Risk: medium
- Smoke subset: True
- Tags: audit-emission, happy, review
- Compliance tags: PDPA-TH-2019

#### Expected Assertions
- [state_change] A review row is created tied to {customer:C1, order:O1, line:P1, rating:5, text:'<delivery_speed_remark_th>', is_hidden:false}.
- [derived_metric_update] P1.average_rating and P1.review_count refresh to reflect the new review.
- [audit_event_emitted] Audit event 'review.created' emitted with all 7 fields (event, actor='customer:c1', ts, before=null, after={review_row_id}, reason='customer_initiated', idem_key). (AP-Q12)

### TC-EPIC-GOVERNANCE-OBS-2-002 — Second review on same order line is rejected

- Scenario type: error
- Test type: functional
- Pyramid level: integration
- Environment: env-gov-integration
- Owner: SDET
- Reviewer role: QA-Lead
- Risk: medium
- Smoke subset: False
- Tags: error, review, uniqueness

#### Expected Assertions
- [error_response] Second submit rejected with reason='already_reviewed'.
- [row_count_invariant] Review-row count for {O1,P1} remains 1.

### TC-EPIC-GOVERNANCE-OBS-2-003 — Review rejected for non-delivered order

- Scenario type: error
- Test type: functional
- Pyramid level: integration
- Environment: env-gov-integration
- Owner: SDET
- Reviewer role: QA-Lead
- Risk: medium
- Smoke subset: False
- Tags: error, lifecycle, review

#### Expected Assertions
- [error_response] Request rejected with reason='order_not_delivered'.
- [row_count_invariant] No review row created.

### TC-EPIC-GOVERNANCE-OBS-2-004 — Rating outside 1-5 is rejected

- Scenario type: boundary
- Test type: functional
- Pyramid level: unit
- Environment: env-gov-integration
- Owner: SDET
- Reviewer role: QA-Lead
- Risk: low
- Smoke subset: False
- Tags: boundary, review, validation

#### Expected Assertions
- [error_response] Rating 0 and rating 6 each rejected with reason='rating_out_of_range'.
- [boundary_accepted] Rating 1 and rating 5 each accepted (inclusive bounds).

### TC-EPIC-GOVERNANCE-OBS-2-005 — Admin soft-hides a review with reason and review is not edited

- Scenario type: happy
- Test type: integration
- Pyramid level: integration
- Environment: env-gov-integration
- Owner: SDET
- Reviewer role: Compliance-Lead
- Risk: high
- Smoke subset: True
- Tags: audit-emission, happy, moderation, review
- Compliance tags: CPA-TH-DISCLOSURE, PDPA-TH-2019

#### Expected Assertions
- [state_change] R1.is_hidden becomes true; R1.text unchanged (admin cannot edit).
- [public_view_invariant] Hidden review does not appear on the public PDP.
- [audit_event_emitted] Audit event 'review.hidden' emitted with all 7 fields (event, actor='admin:aoy', ts, before={is_hidden:false}, after={is_hidden:true}, reason='violates_review_policy', idem_key). (AP-Q12)

### TC-EPIC-GOVERNANCE-OBS-2-006 — Banking-grade authz: only admins can hide reviews; only owning customer can submit

- Scenario type: banking_grade_authz
- Test type: security
- Pyramid level: integration
- Environment: env-gov-integration
- Owner: SDET
- Reviewer role: Compliance-Lead
- Risk: high
- Smoke subset: False
- Tags: banking_grade_authz, review-moderation
- Compliance tags: CCA-TH-LOG-90D

#### Expected Assertions
- [http_status] Non-admin attempting to hide a review receives HTTP 403.
- [http_status] Customer attempting to submit a review for an order they do not own receives HTTP 403.
- [audit_event_emitted] authz_denied audit event emitted with all 7 standard fields.

### TC-EPIC-GOVERNANCE-OBS-2-007 — Banking-grade tipping-off: moderation rejection copy does not disclose detection rule

- Scenario type: banking_grade_tipping_off
- Test type: compliance
- Pyramid level: integration
- Environment: env-gov-integration
- Owner: SDET
- Reviewer role: Legal
- Risk: high
- Smoke subset: False
- Tags: AP-Q11, banking_grade_tipping_off, product-copy, review-moderation
- Compliance tags: CPA-TH-DISCLOSURE

#### Expected Assertions
- [copy_lint] Customer-visible moderation-rejection copy contains no detection-rule name, internal code, or threshold value (AP-Q11 applied to product copy).
- [copy_lint] Copy uses generic language only (e.g. 'This review violates our community guidelines') and does not disclose which rule fired.
- [internal_log_invariant] Internal moderation audit event records the specific rule_id, while the customer-facing notification message contains no rule_id.

### TC-EPIC-GOVERNANCE-OBS-2-008 — Notification log: review-status change to customer is logged with notification_logged=true

- Scenario type: banking_grade_notification
- Test type: integration
- Pyramid level: integration
- Environment: env-gov-integration
- Owner: SDET
- Reviewer role: Privacy-Lead
- Risk: medium
- Smoke subset: False
- Tags: notification-log, review
- Compliance tags: PDPA-TH-2019

#### Expected Assertions
- [notification_log_record] Every customer-facing review-status notification emits a notification-log row with notification_logged=true and references the notification template id.
- [no_pii_leak] Notification-log payload contains customer_id (indirect) only, not email/phone/name in plaintext (masking_rule applied per pii_inventory).

### TC-EPIC-GOVERNANCE-OBS-2-009 — Idempotency: duplicate review-submit with same idem_key produces single row

- Scenario type: banking_grade_idempotency
- Test type: integration
- Pyramid level: integration
- Environment: env-gov-integration
- Owner: SDET
- Reviewer role: QA-Lead
- Risk: medium
- Smoke subset: False
- Tags: banking_grade_idempotency, review

#### Expected Assertions
- [row_count_invariant] Replayed POST with same idem_key results in exactly one review row.
- [audit_event_count] Exactly one 'review.created' audit event with that idem_key.

