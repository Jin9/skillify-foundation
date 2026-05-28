---
name: evaluate-banking-compliance
description: >
  Scan structured banking Epics and Stories for PII, regulatory dependencies,
  tipping-off risk, missing Legal or Privacy roles, and TL-handoff governance
  blockers. Use when a user asks to "check this brief for banking compliance",
  "run the PII and legal gap scan", or "evaluate regulatory gaps". Do NOT use
  to restructure stories, extract raw notes, sweep ambiguity, or generate ACs.
---

# Skill: Evaluate Banking Compliance

## Purpose

Evaluate a stage-one banking brief skeleton for governance gaps that must be visible before TL handoff.

## Input

Accept a pipeline state JSON containing a completed extraction result with `epics[]`, `stories[]`, `source_type`, `scope_kind`, and source evidence.

Before scanning, read `references/compliance-rules.md`, then inspect `templates/stage2_compliance.schema.json`.

## Procedure

1. Build a PII inventory for every story. Classify direct, indirect, regulatory-confidential, and financial identifiers. If no PII applies, emit a reasoned `not_applicable` row rather than an empty result.
2. Enumerate stakeholders from the extracted epics and source evidence. Add absent-but-implied rows for Legal, Privacy or DPO, Security Reviewer, SAR Liaison, Model Owner, Migration Owner, Treasury, or card-network representatives when the scope requires them.
3. Determine `legal_status` independently from Compliance presence. Compliance can describe rules; it does not satisfy Legal review.
4. Detect regulatory dependencies and unresolved citations for PII, retention, sanctions, AML, KYC, EDD, SAR, PEP, biometric, payment, PCI, PDPA, GDPR, FATF, OFAC, FinCEN, or named regulators.
5. Run tipping-off scan over every customer-facing status, notification, rejection, email, SMS, push, error, or script string. Use safe wording from `references/compliance-rules.md` and require Legal sign-off for unresolved risky language.
6. Emit P1 governance blockers for Legal absent on regulatory scope, tipping-off violations, missing PII inventory, unresolved T1 citations, missing dual-approval owners, missing compensating action, and unstated retention policy.
7. Return strict JSON matching `templates/stage2_compliance.schema.json`.

## Output Contract

Return only JSON with:

- `stage: "compliance"`
- `stage_status: "complete"` or `"blocked"`
- `pii_inventory[]`, `stakeholders[]`, `legal_status_by_epic[]`
- `regulatory_dependencies[]`, `governance_gaps[]`, `tipping_off_scan`
- `blocks_tl_handoff`

Any P1 unresolved gap must set `blocks_tl_handoff: true`.

## References

- `references/compliance-rules.md` - PII taxonomy, stakeholder rules, tipping-off vocabulary, and blocker rules.
- `templates/stage2_compliance.schema.json` - exact stage output contract.
- `examples/stage2_example.json` - compact example of valid stage output.
