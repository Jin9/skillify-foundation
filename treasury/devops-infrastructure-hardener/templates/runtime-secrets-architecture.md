# Runtime Secrets Architecture — <scope>

Target state: zero static credentials. Every credential class below is
issued short-lived and task-scoped at runtime. The agent never holds a
durable secret value in the repo, agent config, or MCP config.

## Per-credential-class design

| credential class | provider | task-duration TTL | identity principal | injection mechanism | rotation trigger |
|------------------|----------|-------------------|--------------------|--------------------|------------------|
| <e.g. cloud IAM> | workload identity federation: AWS IRSA / GCP Workload Identity / Azure Managed Identity | per-request token (no stored key) | <service account / workload identity> | runtime token exchange (no key mounted) | n/a (no static key) |
| <e.g. DB password> | HashiCorp Vault dynamic secret | <= task duration (e.g. 1h) | <agent workload identity to Vault> | Vault client / `hvac` at task start; lease revoked on completion | lifecycle event: deploy / config / anomaly / scope change |
| <e.g. 3rd-party API key> | AWS Secrets Manager | AWS-managed rotation; fetched per task | <IAM role, scoped to one secret ARN> | CLI-injected at process start; not stored | AWS-managed + scope change |

Pick exactly one provider per credential class.

## Per-class detail (repeat block per class)

### <credential class>

- **Provider:** <Vault dynamic | AWS Secrets Manager | workload identity federation (IRSA / GCP WI / Azure MI)>
- **Why this provider:** <one line>
- **Task-duration TTL:** <minutes/hours; never calendar-based, never unbounded>
- **Identity principal:** <unique non-human identity that authenticates the request>
- **Injection mechanism:** <CLI wrapper / SDK at runtime — value exists only in the wrapped process for the task>
- **Revocation:** <revoked on task completion AND TTL expiry>
- **Rotation trigger:** <deployment change | config change | anomaly alert | scope change — not the calendar>
- **Static value present anywhere?** MUST be "no" (repo / env / agent config / MCP config).

## Invariants

- [ ] No credential class retains a static value in env vars, agent config, or MCP config.
- [ ] Every TTL is task-scoped, not calendar-based, and shorter than the realistic exploitation window.
- [ ] Every request authenticates via a unique non-human identity, not a stored token.
- [ ] CLI/SDK injection keeps the value out of the repo and out of process-persistent config.
