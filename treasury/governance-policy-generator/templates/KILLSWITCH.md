# KILLSWITCH.md (SKELETON)

> ARCHITECTURAL SEPARATION — MANDATORY
> This file and its JSONL audit MUST live in infrastructure the agent cannot
> read, write, or include in its own prompt/context. "Both the kill switch
> and the circuit breaker must be implemented in infrastructure the agent
> cannot touch — not in the agent's own context or instruction set." An agent
> that can modify its own governance policy effectively has no constraints.
> Do not remove this banner.

## TRIGGERS

Quantitative conditions that fire the killswitch. Every trigger needs a
number; a trigger with no threshold is not a trigger.

- Cost limit: halt if spend exceeds `<COST_LIMIT>` per `<WINDOW>`.
- Error threshold: halt if error rate exceeds `<ERROR_RATE>` over `<N>` actions.
- Circuit breaker — rate: more than `<N>` tool calls in `<M>` minutes.
- Circuit breaker — loop: same tool call with identical parameters `<K>` times (suggest 3).
- `<ADD ENVIRONMENT-SPECIFIC TRIGGERS>`

## FORBIDDEN

Categorical actions never permitted regardless of allowlist state
(supplementary emergency brake, not the primary access model).

- Protected files: `<LIST PATHS>`
- Protected APIs / endpoints: `<LIST>`
- Protected system commands: `<LIST>`
- Self-modification of this file, the Rego policy, or the audit log.

## ESCALATION (three levels, fixed order)

1. THROTTLE — on `<SOFT_SIGNAL>`: rate-limit tool calls to `<RATE>`. Promotes
   to PAUSE if `<THROTTLE_PROMOTE_CONDITION>`. Notify: `<CONTACT>`.
2. PAUSE — suspend the agent and alert a human; await explicit input.
   Promotes to FULL STOP if `<PAUSE_PROMOTE_CONDITION>`. Notify: `<CONTACT>`.
3. FULL STOP — hard halt of the agent process. Requires OVERRIDE to resume.
   Notify: `<CONTACT>` and `<INCIDENT_CHANNEL>`.

## AUDIT (append-only JSONL)

Append-only — no in-place edits or deletes. One JSON object per line, stored
in separated infrastructure, streamed to SIEM/SOAR. Schema:

```
{"ts":"<RFC3339>","agent_id":"<ID>","action":"<ACTION>","decision":"allow|deny|throttle|pause|stop","policy_rule":"<RULE>","trigger":null,"blast_radius":"low|medium|high","reversible":true,"identity":"<WORKLOAD_IDENTITY>"}
```

Captures intermediate agent decisions, not just final outputs: timestamps,
decision logic, tool usage, policy evaluations, and identity mappings.

## OVERRIDE

Conditions under which a human may resume or relax the killswitch — explicit
human approval required; never self-service by the agent.

- Authorized approver(s): `<NAMES/ROLES>`
- Required evidence: `<APPROVAL_EVIDENCE>`
- Every override writes an audit line with the approver and justification.
- Override of FULL STOP additionally requires `<SECOND_APPROVER>`.
