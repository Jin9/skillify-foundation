# Infrastructure as Code (Terraform)

## Contents
- Conventions
- Workflow

## Conventions

- Module boundary = ownership boundary. One team owns one module.
- Remote state with locking: S3 + DynamoDB (AWS) or GCS + Cloud Storage (GCP).
- NEVER hardcode secrets, account IDs, or regions in `.tf` files. Use variables + tfvars.
- Tag every resource: `team`, `environment`, `service`, `cost-center`.

## Workflow

- `terraform plan` output MUST be reviewed in PR before merge.
- Auto-apply only on `main` with guardrails (plan approval, drift detection).
- Use `terraform fmt` and `tflint` in CI. Block merge on failures.
- Pin provider versions explicitly (`~>` for minor, `=` for critical).
