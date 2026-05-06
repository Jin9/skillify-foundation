# Compliance checklists

Apply the checklists relevant to the data classes and geography in the
intake. A regulation is in scope when (a) the operating geography is
covered and (b) the requirement touches data classes the regulation
governs.

## KYC (Know Your Customer)

Applies to onboarding, re-verification, ownership change, account
upgrade, and any flow that promotes a borrower to a higher risk tier.

- [ ] Identity proofing tier specified (basic / enhanced / EDD).
- [ ] Identity-document set named (e.g., national ID + selfie liveness;
      passport + utility bill).
- [ ] Document retention period named (default: 5 years post account
      closure unless local regulation requires longer).
- [ ] Override / manual-review path defined with dual-control.
- [ ] Re-verification triggers listed (large transaction, address
      change, sanctions-list re-screen).
- [ ] Audit-trail requirement: actor + timestamp + before/after on every
      KYC state change, append-only sink.
- [ ] PII access boundaries: who can see plaintext NIK / PAN / DOB.

## AML (Anti-Money Laundering)

Applies to onboarding, transaction monitoring, and any feature changing
how money flows in or out.

- [ ] Risk scoring model named with the inputs it uses.
- [ ] Sanctions / PEP / adverse-media screening provider named with
      refresh cadence.
- [ ] Transaction-monitoring rules listed with thresholds (single tx,
      cumulative-N-day, velocity).
- [ ] Suspicious Transaction Report (STR / SAR) trigger and routing
      defined, with regulator-specific deadline (commonly 24h–7d).
- [ ] Travel rule applied for cross-border transactions ≥ threshold.
- [ ] Look-back window for retroactive screening defined.
- [ ] Investigator workflow: who reviews, what they can see, escalation.

## PCI-DSS (Payment Card Industry)

Applies to any flow that stores, processes, or transmits cardholder data
(PAN, CVV, expiry, magnetic-stripe data).

- [ ] PAN handling tier (none / tokenized / encrypted in scope / not in
      scope) declared and justified.
- [ ] CVV / sensitive auth data retention prohibited post-authorization.
- [ ] Tokenization vendor named with key-management custody.
- [ ] Network segmentation between cardholder data environment (CDE)
      and rest of the platform documented.
- [ ] Logging of PAN-touching events emits masked PAN only.
- [ ] Quarterly ASV scan and annual penetration test on the CDE.
- [ ] Card-on-file consent capture and revocation flow.
- [ ] Out-of-scope reduction strategy documented if PAN is intentionally
      avoided (redirect / iframe / hosted fields).

## PDPA — Personal Data Protection Act (Singapore, Thailand, Malaysia variants)

Applies in SEA jurisdictions. Each country's PDPA has nuances; cite the
relevant national act in the spec.

- [ ] Lawful basis named per data class (consent / contract / legal
      obligation / legitimate interest where allowed).
- [ ] Consent capture UX described (granular, revocable, time-stamped).
- [ ] Cross-border transfer mechanism declared (adequate jurisdiction,
      contractual safeguards, explicit consent).
- [ ] Data subject rights flow: access, rectification, deletion,
      portability — with response SLA (commonly 30 days).
- [ ] Data Protection Officer (DPO) contact surfaced in the spec.
- [ ] Breach notification SLA (commonly 72h to regulator + affected
      subjects).

## GDPR — EU General Data Protection Regulation

Applies when the platform offers services to EU subjects regardless of
where the operator sits.

- [ ] Lawful basis named per processing activity (Art 6).
- [ ] Special category data handling (Art 9) — biometric KYC liveness
      counts.
- [ ] DPIA (Data Protection Impact Assessment) required when the flow
      involves systematic monitoring, profiling, or large-scale
      processing of special-category data.
- [ ] Right-to-erasure flow with conflicts noted (e.g., AML retention
      mandates retention overriding erasure for stipulated period).
- [ ] Data Processor agreements (Art 28) for every vendor.
- [ ] EU Representative named if the operator is outside the EU.
- [ ] 72h breach notification with content per Art 33.

## OJK (Indonesia) / MAS (Singapore) / BOT (Thailand)

Sector-specific regulators. Pull the latest circular for the flow being
specified; the requirements below are stable but version drifts fast.

- [ ] Licensing class declared (e.g., Indonesia OJK fintech P2P, MAS
      digital bank, BOT e-payment service provider).
- [ ] Reporting cadence to the regulator (monthly / quarterly / on-
      event) listed.
- [ ] Local data-residency requirements respected — name the storage
      region for every data class.
- [ ] Outsourcing notification rules (commonly: regulator-notified
      outsourcing for material processors).
- [ ] Customer dispute-resolution SLA per regulator (commonly 5–15
      working days).
- [ ] Mandatory disclosures in customer-facing copy (interest rates,
      fees, late fees, dispute channel).

## Cross-cutting compliance verdict

After running each in-scope checklist, render the verdict:

```
**Compliance verdict:** green | yellow (with conditions) | red (P1 blocker)

Conditions / blockers:
- <regulation>: <unresolved item> — owner: <stakeholder> — by: <date>
```

Verdict rules:

- **green** — every in-scope checklist passes; no open conditions.
- **yellow** — passes with explicit, time-bound conditions (e.g., "STR
  routing TBD; legal-ops to confirm by <date>").
- **red** — at least one P1 violation. The spec is not deliverable until
  the blocker resolves. State the regulation, the offending requirement,
  and the safer alternative.

## What this skill does not do

- Draft the regulator-specific notification text — escalate to legal or
  compliance.
- Sign off on a high-risk change — that requires the named DPO / CCO /
  legal contact, not this skill.
- Validate against regulator drafts that are not yet enacted — name the
  draft and surface the risk; do not bind the platform to draft
  requirements.
