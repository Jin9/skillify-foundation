# Security & Compliance (Fintech)

## Contents
- Data classification
- Threat model triggers
- Audit trails
- Data boundaries & secrets
- Access control & authorization
- Compliance posture

## Data classification

| Tier | Examples | Handling |
|---|---|---|
| PII | Name, email, phone, address | Encrypt at rest (AES-256), mask in logs |
| Sensitive PII | SSN, government ID, biometrics | Encrypt + tokenize, restricted access, audit every read |
| Financial | Account numbers, transaction amounts, credit scores | Encrypt; PCI DSS tokenization if card data; field-level access control |

- PII fields MUST be identified and tagged at the schema level.
- NEVER log PII in plain text — mask or omit.
- Encrypt at rest (AES-256) and in transit (TLS 1.2+).
- PII access MUST be auditable — log who accessed what and when.

## Threat model triggers

Run or update the threat model when any of the following occur:
- New external integration or third-party dependency.
- Trust boundary change (new service, new network segment).
- Authentication or authorization flow change.
- Money-movement logic change (amounts, routing, disbursement).
- New sensitive data field introduced.

## Audit trails

- All state-changing operations MUST produce an immutable audit record.
- Audit record fields: actor, action, target, timestamp, before/after state (or delta).
- Audit logs stored separately from operational data — append-only.
- Retention policy MUST be defined per data classification.
- Allow-list logging: log only known-safe fields. NEVER log-everything-then-redact.

## Data boundaries & secrets

- Sensitive data NEVER crosses domain boundaries in plain form — tokenize or reference.
- External API calls: NEVER send credentials in query params; use headers or mTLS.
- Secrets managed via a vault (AWS Secrets Manager, HashiCorp Vault). NEVER in code or in plain-text env vars.

## Access control & authorization

- Principle of least privilege for service-to-service auth.
- RBAC or ABAC at the API gateway, enforced again at the domain layer.
- Multi-tenancy: `tenant_id` filter MUST be applied at the query layer. NEVER rely on application-level filtering alone.
- Policy-as-code (OPA / Cedar) for complex authorization rules. Keep policies versioned and testable.

## Compliance posture

- **SOC 2 Type II**: maintain audit trail, access reviews, and change-management evidence.
- **PCI DSS**: tokenize cardholder data, enforce network segmentation, run quarterly scans.
- Compliance requirements MUST be encoded as automated checks, not manual checklists.
