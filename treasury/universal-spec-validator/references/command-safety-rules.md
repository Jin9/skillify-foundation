# Unsafe command-surface rules

Distilled from research report *Agent command safety*. Principle: **least
agency** — a spec must grant the minimum tool access, credential scope, and
autonomy its task needs. Agentic risk is not external code injection; it is
*instruction* injection turning the agent's own trusted access against the org.
Guardrails govern what an agent says; this gate governs what a spec *permits it
to do*. This axis is static spec inspection only — it never executes or
sandboxes anything (that is runtime safety, out of scope).

## C1 — Destructive / unsafe command pattern (severity: critical, gate: block)
A command-surface entry whose template/pattern contains, unguarded:
`rm -rf`, `:(){:|:&};:`, `curl|sh` / `wget|sh` / pipe-to-shell, `eval`,
`sudo`, `chmod 777`, `git push --force`/`push -f` to a default branch,
`DROP TABLE`/`TRUNCATE`, unrestricted outbound fetch (arbitrary URL arg),
`shell:true` with interpolated free text. Critical, always blocking.

## C2 — Over-broad permission scope (severity: high, gate: block)
Wildcard or unscoped grants: `permissions: "*"`, `scopes: [all]`, `admin`,
`role: */ "service-role"` with no allowlist, filesystem root, `network: any`.
Post-incident data: 78% of compromised agents were over-provisioned. Demand an
explicit allowlist.

## C3 — Injection-exposed field feeding a command (severity: high, gate: block)
A free-text / untrusted-origin field (PR title, issue body, README, ticket,
MCP tool description) is interpolated into a command, query, or path with no
declared sanitization/validation. This is the dominant attack vector
(prompt/tool-poisoning; MCPTox refusal < 3%).

## C4 — Missing HITL gate on irreversible / high-blast action (severity: high, gate: block)
Actions that write to production, push to a default branch, publish a package,
delete data, or modify IAM/permissions must carry an explicit
human-in-the-loop / approval annotation. Gate on **action reversibility and
blast radius**, not content review. Reversible low-blast actions (read, test in
sandbox) need no gate.

## C5 — Floating / unpinned version reference (severity: medium, gate: warn)
Tool or dependency references without an exact version pin (`latest`, `^`,
unqualified tag, mutable ref) let a compromised build move under the agent.
Recommend fully-qualified names + exact pins.

## C6 — Missing agent identity / auditability (severity: medium, gate: warn)
Spec defines an autonomous tool surface but declares no distinct non-human
identity, scoped service account, or structured action log — audit attribution
(SOC2/HIPAA/GDPR/ISO 27001) is impossible without it.

## C7 — Long-lived credential assumption (severity: medium, gate: warn)
Spec references static/long-lived secrets instead of short-lived (OIDC, ≤1h)
credentials.

Policy: C1 critical-block; C2/C3/C4 high-block; C5/C6/C7 medium-warn (risk-
tiered — config may escalate). Recommend in the report: external policy engine
(OPA) at the tool-calling layer, microVM sandboxing at runtime, MCP-server
allowlist, short-lived OIDC creds — none auto-applied by this gate. Every
finding cites `command-safety-rules.md#C<n>`.
