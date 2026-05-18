# Four-layer defense-in-depth stack

The hardening plan organizes every finding against this four-layer stack.
The layers are sequenced by blast-radius reduction: Layer 1 first.

```
Defense-in-Depth Stack
──────────────────────────────────────────
  Layer 4 │ Lifecycle-event rotation + unique agent identities
  Layer 3 │ Output filtering + privacy-first observability
  Layer 2 │ Least-privilege scoping + sandboxing + egress limits
  Layer 1 │ Runtime credential injection (no static secrets)
──────────────────────────────────────────
```

## Layer 1 — Eliminate static credentials

The most impactful control is architectural: eliminate static credentials
from agent execution environments. This means replacing environment-variable
injection of API keys with runtime credential requests to a secrets manager.
HashiCorp Vault dynamic secrets generate a credential at task start with a
built-in expiration (e.g., 1 hour), then revoke it automatically — a stolen
credential expires before it can be exploited. For cloud workloads, workload
identity federation (AWS IRSA, GCP Workload Identity, Azure Managed Identity)
removes static keys entirely: the agent proves its identity through its
runtime environment, not a credential it possesses.

**What it eliminates:** the static secret itself — there is nothing durable in
env vars, agent config, or MCP config for logging, tool output, or prompt
injection to surface. Detail in `runtime-secrets-architecture.md`.

## Layer 2 — Least-privilege tool scoping

The AWS Well-Architected Generative AI Lens mandates per-tool permission
scoping: read-only vs. write access, specific resource identifiers,
time-bounded grants. Every tool the agent calls should be gated on the
minimum scope required for that task. NVIDIA's sandboxing guidance adds three
mandatory structural controls: network egress allowlists (preventing
credential exfiltration to arbitrary endpoints), workspace write restrictions
(blocking modification of dotfiles and auto-executing configs), and
configuration file protection (blocking changes to MCP server configs and IDE
extension manifests).

**What it eliminates:** the value of a credential that does leak — a tightly
scoped, time-bounded grant plus an egress allowlist caps what an attacker can
do with it and where data can be sent.

## Layer 3 — Output filtering and audit logging

Every agent output boundary should run a regex-based credential-pattern
scrubber before results are logged or returned. Observability tooling must be
configured to mask credential patterns and hash or omit sensitive parameters
from trace logs. This is a non-trivial trade-off: the same level of detail
that makes logs useful for debugging makes them a credential exfiltration
risk if the logging system is compromised. Privacy-first observability
deployments use on-premises or VPC-isolated log stores.

**What it eliminates:** the exfiltration channel — debug logging accounts for
73.5% of agent skill vulnerabilities, so scrubbing the output boundary closes
the dominant leak path. Scrubber rules in `output-scrubber-spec.md` (template).

## Layer 4 — Credential rotation on lifecycle events

Calendar-based 90-day rotation is too slow for AI-era attack speeds:
documented incidents show credentials can be fully exploited within hours of
theft. Rotation should be triggered by deployment changes, configuration
modifications, anomaly detection alerts, or scope changes — not the calendar.
Each AI agent should have a unique service account identity so that a
compromised credential can be revoked without affecting other agents or human
engineers.

**What it eliminates:** persistence and shared blast radius — event-triggered
rotation plus unique per-agent identity means a compromised credential is
short-lived and revocable in isolation.

## Cited statistics (verbatim from the grounding research)

- Debug logging accounts for 73.5% of agent skill vulnerabilities; print/log
  statements pipe secrets directly into the LLM context window. Seven percent
  of studied skills exposed credentials through the LLM context window or
  output logs.
- In 2025, over 1.2 million AI-service credentials were exposed — an 81%
  year-over-year surge (AI-service credentials surged 81% YoY to 1.27 million;
  29 million total secrets leaked to public GitHub in 2025; 24,008 unique
  secrets found in MCP configuration files alone).
- 64% of valid secrets leaked in 2022 remain un-revoked in 2026 (64% of
  leaked secrets from 2022 are still active in 2026) — primarily a governance
  gap, not a tooling gap.
- No major agent framework activates security controls by default; frameworks
  feed tool outputs directly back into the LLM context without sanitization.

## Defense-in-depth ordering rationale

Layer 1 removes the asset; the upper layers contain a leak if Layer 1 is
ever bypassed. Remediate top of the finding queue = Layer 1 items, because
eliminating the static secret is the single largest blast-radius reduction.
Layers 2–4 are necessary but assume a credential can still transit the system
and therefore reduce, not remove, exposure. Never present Layers 2–4 as a
substitute for Layer 1.
