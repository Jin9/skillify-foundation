# Hardening Plan — <scope>

Audit scope: <repos / agents / pipelines>
Scan command: `python3 scripts/scan_static_creds.py --root <scope>`
Scan result: <N findings | clean (note: empty != clean, see caveat)>
Date: <yyyy-mm-dd>

## Findings (prioritized — Layer 1 first)

| id | leakage pattern | file:line | defense layer | remediation | trust-tier | priority |
|----|-----------------|-----------|---------------|-------------|------------|----------|
| F-001 | hard-coded key | <path:line> | Layer 1 | Replace with <Vault dynamic / AWS Secrets Manager / workload identity federation>; revoke exposed key now | Tier 2 (control-plane? Y/N) | P0 |
| F-002 | debug-logging | <path:line> | Layer 3 + remove statement | Delete log/print; scrub output boundary per output-scrubber-spec | Tier 1 | P0 |
| F-003 | env-var inheritance | <path:line> | Layer 1 | CLI-inject short-lived credential; remove env assignment | Tier 2 | P1 |
| F-004 | MCP-config secret | <path:line> | Layer 1 + Layer 2 | Dynamic per-server credential; config-file protection | Tier 2 | P0 |
| F-005 | prompt-injection exfiltration | <design> | Layer 1 + 2 + 3 | No static secret in context; egress allowlist; scrub | Tier 1 | P1 |

Priority key: P0 = static secret currently exposed / Layer 1; P1 = leak
channel open; P2 = hardening / depth.

## Layer rollup

- **Layer 1 (eliminate static creds):** <count> findings — remediate first.
- **Layer 2 (least-privilege + egress):** <count> findings.
- **Layer 3 (output filtering):** <count> findings.
- **Layer 4 (rotation + unique identity):** <count> findings.

## Control-plane gate

List any remediation whose action touches deployment policy, approval gates,
or rollback thresholds. These remain Tier 1–2, human-gated:

- <action> — Tier <1|2>, approver: <named human>

## Residual risk

State what remains after the plan is fully applied and why it cannot be
eliminated here:

- <residual item> — why it persists, compensating control, owner.
- Design-level patterns (MCP trust, prompt-injection) not caught by literal
  scanning — recommend a human red-team exercise targeting credential
  exfiltration via prompt injection.
- Any credential assumed compromised (was hard-coded) — confirm revocation;
  64% of secrets leaked in 2022 are still active in 2026, so revocation, not
  just rotation, must be verified.

## Self-check (must all be true before emit)

- [ ] Re-ran the scanner against the proposed design: no static secret reintroduced.
- [ ] Every Layer-1 item names a secrets manager / federation with a task-duration TTL.
- [ ] Every agent has a unique non-human identity and a named owner.
- [ ] No control-plane action exceeds Tier 2.
- [ ] Residual-risk section is present and non-empty.
