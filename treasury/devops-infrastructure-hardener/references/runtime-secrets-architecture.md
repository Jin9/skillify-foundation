# Runtime secrets architecture (Layer 1 design)

Layer 1 replaces every static credential with a short-lived, task-scoped
credential issued at runtime. The agent never holds a durable secret value.
A "secret" here is any static credential the agent uses to authenticate with
an external service: API keys, OAuth tokens, cloud IAM access keys, database
passwords. The defining property of an agent secret is that it is needed at
runtime, without a human in the loop — which is exactly why startup env-var
injection creates exposure (the agent inherits all credentials and may surface
any of them via logging, tool output, or prompt injection).

Pick exactly one provider per credential class.

## Option A — HashiCorp Vault dynamic secrets

Vault dynamic secrets generate a credential at task start with a built-in
expiration (e.g., 1 hour), then revoke it automatically — a stolen credential
expires before it can be exploited. The agent authenticates to Vault using its
workload identity (not a Vault token stored in config), requests a lease for
the specific backend (database, cloud, PKI), and uses the leased credential
for the duration of the task only.

- TTL: set to the task duration (minutes to ~1 hour), never open-ended.
- Lease is revoked on task completion or TTL expiry, whichever is first.
- Vault client: `hvac` (Python) or equivalent generates the credential at
  runtime; no credential value lives in the repo or agent config.

## Option B — AWS Secrets Manager

Use when the credential must be a managed AWS-native secret (e.g., an RDS
password or third-party API key under AWS rotation). The agent retrieves the
secret at task start via its IAM role (see Option C for the identity); it does
not store the secret. Pair with AWS-managed rotation so the stored value is
itself short-lived, and scope the `secretsmanager:GetSecretValue` grant to the
single secret ARN the task needs.

## Option C — Workload identity federation (no static keys at all)

The strongest Layer-1 control: there is no credential to leak. The agent
proves its identity through its runtime environment via an OIDC trust
relationship; the cloud issues short-lived, task-scoped credentials on
demand.

- **AWS IRSA** — IAM Roles for Service Accounts: a Kubernetes service account
  is federated to an IAM role; the pod receives short-lived STS credentials.
- **GCP Workload Identity** — a Kubernetes service account is bound to a GCP
  service account; tokens are minted per request.
- **Azure Managed Identity** — the workload's managed identity is exchanged
  for short-lived Azure AD tokens.

No long-lived access key is ever stored, mounted, or printed.

## Task-duration TTL design

- TTL = expected task duration with a small safety margin; never a calendar
  interval and never unbounded.
- Documented incidents show stolen credentials can be fully exploited within
  hours; the TTL must be shorter than the realistic exploitation window.
- The credential is revoked on task completion as well as on TTL expiry.
- Long-running agents re-request a fresh short-lived credential per task
  iteration rather than holding one for the process lifetime.

## CLI injection mechanism

The simplest available mitigation is CLI-based injection: a secrets-manager
CLI wraps the agent process and injects the credential into the process
environment for that invocation only, so the agent never holds the value in
its own configuration. For dynamic secrets, the Vault client or SDK generates
the short-lived credential at runtime, limiting credential lifetime to the
task's duration. The injected value exists only in the wrapped process's
memory for the task and is gone when the process exits — there is no static
secret in the repo, the agent config, or any MCP config.

## Mapping each credential class

For every distinct credential class found in step 1, the architecture doc
records: provider (A/B/C above), task-duration TTL, the identity principal
that authenticates the request, the injection mechanism, and the rotation
trigger. Template: `templates/runtime-secrets-architecture.md`.
