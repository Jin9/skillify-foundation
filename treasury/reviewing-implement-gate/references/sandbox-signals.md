# Sandbox-required signals

If any of the signals below appears in `plan.md`, `critique.md`, repo
metadata (`README.md`, `CLAUDE.md`, `AGENTS.md`), or recent diffs, treat
the workflow as sandbox-required. If the launch was unsandboxed, recommend
reject and ask the user to relaunch with `IMPLEMENT_SANDBOXED=1`.

## File-system signals

- `*.env`, `.env.*`, `.envrc` outside `*.example` patterns
- `secrets/`, `credentials/`, `vault/` directories
- `config/credentials*`, `config/keys*`
- `terraform/`, `.tf` files (touch cloud accounts)
- `database/migrations/`, `db/migrate/`, `migrations/*.sql`
- `kubernetes/`, `k8s/`, `helm/` manifests
- `.github/workflows/` (CI/CD secrets)
- `Dockerfile` lines containing `ARG.*KEY|TOKEN|SECRET|PASSWORD`

## Domain signals

- KYC, AML, PCI, PCI-DSS, GDPR, PDPA, BOT, OJK, MAS
- KMS, Vault, IAM, mTLS, JWT signing key, OAuth client secret
- Partner-bank, payment, disbursement, settlement, ledger
- PII fields: PAN, NIK, account number, customer name, date of birth, SSN
- Loan origination, credit decision, underwriting, repayment, collections

## Repo-tagged

- `README.md` mentions: fintech, banking, lending, neobank, payments
- `CLAUDE.md` or `AGENTS.md` declares the repo as production-bearing
- A nearby skill folder is named `banking-*`, `fintech-*`, `treasury-*`

## Negative signals (sandbox not required)

- Documentation-only changes (`docs/`, `*.md`, `*.mdx`)
- Test fixtures with synthetic data (`testdata/`, `__fixtures__/`)
- Internal CLI tools that do not touch network or filesystem outside
  `~/.config/<tool>`

## When in doubt

Default to **sandbox required**. The cost of a false positive is one
relaunch; the cost of a false negative is exfiltrated credentials.
