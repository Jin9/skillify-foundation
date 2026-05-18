# Output Scrubber Spec (Layer 3) — <scope>

Every agent output boundary runs this regex-based credential scrubber before
results are logged or returned. Debug logging accounts for 73.5% of agent
skill vulnerabilities, so this boundary is the highest-value Layer-3 control.
The scrubber is a containment layer — it assumes Layer 1 has already removed
static secrets; it does not replace Layer 1.

## Credential classes and redaction rules

| credential class | example signature | redaction action |
|------------------|-------------------|------------------|
| AWS access key id | `AKIA` + 16 uppercase/digit chars | replace with `REDACTED_AWS_KEY` |
| Private key block | `-----BEGIN ... PRIVATE KEY-----` ... `-----END ... PRIVATE KEY-----` | drop the entire block, replace with `REDACTED_PRIVATE_KEY` |
| Bearer / OAuth token | `Authorization: Bearer ...`, long opaque token | replace token with `REDACTED_TOKEN` |
| Generic api_key / secret literal | value assigned to a name containing `secret`, `token`, `api_key`, `apikey`, `password`, `passwd`, `client_secret` | replace value with `REDACTED_SECRET` |
| URL-embedded credential | `://user:pass@host` or `?token=...` | strip credential component |
| <project-specific class> | <signature> | <action> |

## Apply points

- Tool-output boundary (before output re-enters the LLM context).
- Log / stdout / stderr boundary (before capture by the framework).
- Trace / observability export (before any span leaves the process).
- Agent response boundary (before returning to a human or another agent).

## Log-sink policy

- **Masking required:** observability tooling MUST mask credential patterns
  and hash or omit sensitive parameters from trace logs.
- **Privacy-first storage:** trace / log stores are on-premises or
  VPC-isolated; no agent context is exported to a third-party sink in clear.
- **Retention:** redacted logs only; raw pre-scrub output is never persisted.
- **Trade-off acknowledged:** the detail that makes logs useful for debugging
  makes them an exfiltration risk if the log store is compromised — default
  to redaction, raise verbosity only behind explicit, time-boxed approval.

## Invariants

- [ ] Scrubber runs at every apply point above, not just final response.
- [ ] Every credential class in the table has a redaction action.
- [ ] Log sink is privacy-first (on-prem / VPC-isolated) with masking on.
- [ ] Raw pre-scrub output is never persisted.
- [ ] Spec states it is containment, not a substitute for Layer 1.
