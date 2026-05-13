# Stack Mindset Bindings

Concrete stack bindings for the security mindset chain in `SKILL.md`. Use this file when starting a review on a new artifact type so the asset / boundary / entry-point / actor / impact mapping is already grounded.

## Chain (recap)

```
Asset  →  Trust Boundary  →  Entry Point  →  Actor  →  Threat (STRIDE)
       →  Attack Path  →  Impact (CIA + Financial + Regulatory)
       →  Existing Control  →  Gap  →  Recommended Control  →  Residual Risk
```

## Assets

- Loan PII (name, NIK, address, bank account)
- KYC documents and verification artifacts
- Credit-decision artifacts (scorecard outputs, model versions, decision rationale)
- Disbursement instructions and repayment schedules
- MySQL rows touching borrower or loan state
- Kafka topic payloads on `loan.*`, `kyc.*`, `disbursement.*`
- RDS snapshots and backup exports
- JWT signing keys, Kong/APISIX consumer credentials
- K8s secrets and KMS-managed encryption keys

## Trust boundaries

- Internet → Kong/APISIX → Gin service → Kafka → consumer service → MySQL
- SIT ↔ UAT ↔ PRD environment promotion
- Tenant ↔ tenant separation
- Underwriter ↔ borrower role separation
- Service-to-service (mTLS-enforced or not)

## Entry points

- HTTP route (Gin)
- Kafka topic + key + headers
- gRPC method
- Scheduled job / cron
- Admin CLI
- S3 pre-signed URL
- Partner-bank webhook callback

## Actors

- Unauthenticated internet user
- Authenticated borrower
- Authenticated underwriter
- Internal service identity
- Partner-bank API
- Compromised pod (lateral-movement scenario)
- Malicious insider with read-only DB access

## Impact dimensions

- **Confidentiality** — PII / PCI exposure.
- **Integrity** — forged decisions, altered amounts, tampered audit trail.
- **Availability** — DoS during loan-app cutoff, settlement window, or month-end batch.
- **Financial** — direct disbursement loss, fraudulent refund, double-spend.
- **Regulatory** — BOT, PDPA, GDPR, OJK, MAS, PCI-DSS exposure.

## Anti-patterns when applying this chain

- Treating "Asset" as the file in front of you (a handler is not an asset; the data it returns is).
- Skipping Actor enumeration on a route that "looks internal" — the trust boundary may not be where you assumed.
- Equating Impact with severity — impact is dimensional, severity is ranked.
