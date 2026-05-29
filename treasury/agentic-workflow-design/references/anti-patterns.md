# Anti-patterns (never do)

Supervision-design failures the artifact must not contain. The body of `SKILL.md` points here; this file is the authoritative copy.

- Let an agent take an irreversible, control-plane, or production action without a human gate — regardless of its confidence.
- Rely on prompt/output guardrails alone; ship without an external policy engine enforcing allow / confirm / deny.
- Over-provision agents (broad permissions, long-lived credentials, unrestricted tools/MCP) — least agency or nothing.
- Treat untrusted content (PR titles, issue bodies, tool descriptions) as trusted instructions.
- Add agents or handoffs without justification, or forward full history instead of a typed handoff contract.
- Govern cost/loops with alert-only dashboards and no pre-execution caps that terminate the run.
- Ship without logging the tool-effect half of the loop, or with a silently-editable (non-tamper-evident) audit log.
- Let the agent decide architecture alone, approve its own PRs, fix production directly, or modify its own permission config.
- Attribute AI output as an accountable co-author, or leave AI-generated output without a named human owner of record.
