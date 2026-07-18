# Thailand fintech lens (conditional)

Apply this lens only when the idea touches banking, payments, lending, e-money, or Thai-regulated personal data. For anything else, skip this file entirely. This is orientation vocabulary for an advisory conversation — not legal advice; recommend the user verify load-bearing conclusions with their compliance or legal function.

## Regulator vocabulary map

- **BOT (Bank of Thailand)** — supervises banks, e-payment, and many lending activities. The licensing class (bank, e-money, payment service provider, nano/pico finance, P2P lending) is often the first feasibility question, and BOT IT-risk and outsourcing/cloud guidelines shape what infrastructure is allowed.
- **PDPA (B.E. 2562)** — Thailand's personal-data law: lawful basis and consent, purpose limitation, retention limits, cross-border transfer conditions, data-subject rights, and a DPO where required.
- **AMLA / AMLO** — anti-money-laundering regime: KYC/CDD tiers, sanctions screening, suspicious-transaction reporting, and the tipping-off prohibition — never reveal to a customer that they are under AML suspicion or investigation.

## Feasibility modifiers

- Regulated flows raise the cost lens: audit logging, retention, dual control, and reporting are build cost, not afterthoughts.
- The license class is often the real feasibility gate, not the technology. "Can we build it" is frequently "can we be licensed to operate it" — surface that early.
- Compliance constraints are one-way doors: treat them as constraints, never as one option among several.

## Challenge add-ons

When challenging a fintech idea, add these lenses:

- **Tipping-off exposure** — customer-facing copy (rejection messages, status screens, support scripts) must not leak AML vocabulary: sanctions, flagged, suspicious, SAR, watchlist, investigation, compliance hold. The pattern is neutral phrasing plus a support path.
- **PII in logs and analytics** — direct identifiers or KYC artifacts flowing into logs, crash reports, or third-party analytics.
- **Cross-border data flow** — where the data physically goes, and under which PDPA transfer condition.
- **Single-actor money movement** — disbursement, refunds, overrides, and write-offs need dual control; a design where one actor can move money alone is a finding.
- **Evasion smell** — a design structured to avoid logging, retention, or reporting obligations gets refused, with the nearest compliant alternative offered instead.

## Honest limits

This lens orients the conversation and raises the right questions. It does not replace the regulators' texts, counsel, or the compliance function — recommend verification for anything load-bearing.
