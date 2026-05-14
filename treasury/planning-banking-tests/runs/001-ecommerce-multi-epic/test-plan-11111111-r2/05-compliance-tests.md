---
artifact_type: qa-compliance-tests
test_plan_id: TPLAN-EPIC-SHOPPILOT-MVP
---

# Compliance Tests

## CCA_TH_LOG_90D

### COMP-CCA_TH_LOG_90D-001 — log_retention_enforcement
- Scope: Verify audit-log retention enforces the Computer Crime Act 90-day minimum retention window (rows >= 90 days old retained; deletion before 90d blocked). Retention threshold sourced from CCA-TH-LOG-90D citation, not invented (AP-Q9).
- Scope blocked: true
- Blocking governance gaps: legal_absent_on_regulatory, regulatory_citation_unresolved, retention_policy_unstated

### COMP-CCA_TH_LOG_90D-002 — log_integrity_tamper_evident_hashing
- Scope: Verify audit-log rows carry tamper-evident hashing (e.g. hash-chained per row) and any in-place modification is detectable. Append-only invariant per OQ-20.
- Scope blocked: true
- Blocking governance gaps: legal_absent_on_regulatory, regulatory_citation_unresolved

### COMP-CCA_TH_LOG_90D-003 — access_event_logging
- Scope: Verify every read of the audit-log surface itself is audit-emitted (access audit per 6.6). Reader identity and queried filter captured.
- Scope blocked: true
- Blocking governance gaps: legal_absent_on_regulatory, regulatory_citation_unresolved

### COMP-CCA_TH_LOG_90D-004 — lawful_access_procedure
- Scope: Verify procedure for lawful-access requests (law-enforcement / regulator). Audit-log export workflow exists, is admin-gated, is itself audit-emitted, and uses the scrubbing path that fails-closed on forbidden-field detection.
- Scope blocked: true
- Blocking governance gaps: legal_absent_on_regulatory, regulatory_citation_unresolved, retention_policy_unstated

## CPA

### COMP-CPA-001 — price-transparency
- Scope: Verify final price displayed pre-submit includes VAT + shipping per Consumer Protection Act truthful-pricing rule; no hidden surcharges added post-confirm.
- Scope blocked: true
- Blocking governance gaps: CPA citation_status: pending

### COMP-CPA-001 — refund-disclosure
- Scope: Consumer Protection Act: refund-timeline disclosure on cancellation (manual OOB) - customer-facing copy and notification must state expected refund window; OOB workflow documented.
- Scope blocked: true
- Blocking governance gaps: BA-citation-status-pending-CPA

## PDPA

### COMP-PDPA-001 — pii-minimization-on-order
- Scope: Verify order/audit payloads reference customer_id only and do not duplicate raw PII (full address, phone) beyond the order_line snapshot scope; PDPA minimization principle.

## PDPA_TH_2019

### COMP-PDPA_TH_2019-001 — DSAR-access
- Scope: Customer submits a Data Subject Access Request and receives a structured export of their EPIC-IDENTITY data (email, name, phone, addresses, customer_id, audit-event summary). Test asserts response within statutory window and excludes password hash from the export payload.
- Scope blocked: true
- Blocking governance gaps: pii_inventory_missing, legal_absent_on_regulatory, retention_policy_unstated

### COMP-PDPA_TH_2019-002 — DSAR-rectification
- Scope: Customer submits a rectification request for name / phone / shipping_address; system updates the canonical row, emits banking_grade_audit event with before/after, and propagates correction to active-order snapshots per documented policy.
- Scope blocked: true
- Blocking governance gaps: pii_inventory_missing, legal_absent_on_regulatory, retention_policy_unstated

### COMP-PDPA_TH_2019-003 — DSAR-deletion-and-retention
- Scope: Customer submits a right-to-erasure request. Test asserts: (a) erasure executes on rows whose retention window has expired, (b) rows retained for accounting (e.g., order-attached address snapshots under Thai Revenue Code 5y retention) remain with documented legal basis, (c) erasure audit event emitted.
- Scope blocked: true
- Blocking governance gaps: pii_inventory_missing, legal_absent_on_regulatory, retention_policy_unstated

### COMP-PDPA_TH_2019-004 — breach-notification-72h
- Scope: Simulate a detected PII breach affecting EPIC-IDENTITY data. Test asserts the documented incident-response runbook fires a notification to the PDPC and affected data subjects within 72 hours, and that the breach-notification message contains the statutory required fields without tipping-off vocabulary.
- Scope blocked: true
- Blocking governance gaps: pii_inventory_missing, legal_absent_on_regulatory, retention_policy_unstated

### COMP-PDPA_TH_2019-005 — consent-flow-per-purpose
- Scope: Signup UI presents per-purpose consent toggles (account creation, marketing communications, analytics) at granularity required by PDPA. Test asserts (a) each toggle records lawful basis separately, (b) account creation cannot be conditioned on marketing consent, (c) withdrawing consent is as easy as granting.
- Scope blocked: true
- Blocking governance gaps: pii_inventory_missing, legal_absent_on_regulatory, retention_policy_unstated

### COMP-PDPA_TH_2019-101 — pii_redaction_in_logs
- Scope: Verify PII-redaction rule applied to all log surfaces per pii_inventory[*].masking_rule. Sample-based scan across audit log and notification log; zero leaks for always-masked fields (password, session_token); partial-mask pattern verified for phone/email in non-fulfilment surfaces; address scrubbed from authz_denied payloads.
- Scope blocked: true
- Blocking governance gaps: legal_absent_on_regulatory, pii_inventory_missing, regulatory_citation_unresolved

### COMP-PDPA_TH_2019-102 — cross_border_transfer
- Scope: Verify cross-border-transfer controls when mock-shipping or mock-payment endpoints route through non-Thailand regions. Residency-rule per pii_inventory enforced; transfer events audit-emitted with destination region recorded.
- Scope blocked: true
- Blocking governance gaps: legal_absent_on_regulatory, pii_inventory_missing, regulatory_citation_unresolved

### COMP-PDPA_TH_2019-103 — lawful_basis_record_per_processing_activity
- Scope: Verify each processing activity (signup, order, review-submit, audit-read, notification-send, review-moderation) is annotated with a lawful basis (consent / contract / legitimate-interest / legal-obligation) and that the record-of-processing-activities export contains all activities.
- Scope blocked: true
- Blocking governance gaps: legal_absent_on_regulatory, regulatory_citation_unresolved

## TRC

### COMP-TRC-001 — vat-receipt-display
- Scope: Verify VAT line-item (7%) is displayed on receipts and order confirmation per Thai Revenue Code requirements for retail e-commerce.
- Scope blocked: true
- Blocking governance gaps: TRC citation_status: pending, OQ-CART-tax-invoice-format

### COMP-TRC-001 — record-retention
- Scope: Thai Revenue Code: 5-year accounting record retention for order, payment, audit-log artefacts; verify retention policy on audit rows for order state transitions and admin-cancel reasons.
- Scope blocked: true
- Blocking governance gaps: BA-citation-status-pending-TRC

