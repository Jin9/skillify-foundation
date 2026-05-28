# Compliance Rules

## PII Taxonomy

- `direct`: NRIC, passport, national ID, biometric, liveness, face match, fingerprint, name, phone, email, address.
- `indirect`: source of funds, PEP classification, adverse media result, sanctions match, DOB, account number, applicant ID.
- `regulatory`: SAR filing, sanctions hold, tipping-off communication, suspicious-activity rationale, EDD outcome.
- `financial`: bank statement, transaction data, wire transfer, card token, chargeback evidence.

Every story gets a PII result. If none applies, emit a `not_applicable` row with a workflow-class justification.

## Required Stakeholders

- Legal: customer-facing language, retention, regulator citation, tipping-off, sanctions, biometric, dual approval.
- Privacy or DPO: any PII inventory row.
- Security Reviewer: biometric, file upload, vendor integration, privileged internal tool, audit logging.
- SAR Liaison: SAR, suspicious activity, AML escalation, sanctions hold.
- Model Owner: score-threshold or risk-routing model.
- Migration Owner: cutover, data migration, backfill, dual-write period.
- Treasury: funds movement, settlement, reversals.
- Card-network representative: card dispute, chargeback, VISA, Mastercard, network rule.

Compliance presence does not satisfy Legal review. Emit `legal_status` independently.

## Tipping-Off Scan

Scan customer-facing strings: email, SMS, push, web status, error, rejection text, call-center script, agent-provided customer explanation.

Forbidden or high-risk terms include: sanctions, AML, flagged, suspicious, regulated, SAR, PEP, adverse media, EDD, investigation, compliance hold, watchlist.

Safe phrases should be neutral and non-specific, such as "The transfer could not be completed. Please contact support." If no safe phrase is available, require Legal sign-off and block TL handoff.

## Governance Blockers

Emit P1 blocker with `blocks_tl_handoff: true` when any of these are unresolved:

- Legal absent on regulatory scope.
- Tipping-off violation in customer-facing copy.
- PII inventory missing or empty without justification.
- T1 regulator citation unresolved.
- Dual-approval owner missing.
- Irreversible funds or data action without compensating action.
- Retention policy unstated for PII, SAR, sanctions, biometric, or financial records.
