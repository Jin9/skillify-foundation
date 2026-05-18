# KILLSWITCH.md standard: fields, escalation, audit, OVERRIDE, separation

Deep-dive for Workflow step 5. Fill `templates/KILLSWITCH.md` from this.

## The architectural-separation rule (non-negotiable)

"Both the kill switch and the circuit breaker must be implemented in
infrastructure the agent cannot touch — not in the agent's own context or
instruction set." "An agent that can modify its own governance policy
effectively has no constraints." "Kill switches, circuit breakers, and audit
logs must live in infrastructure that agents cannot read or modify ... The
policy engine that evaluates agent actions must be external to the agent's
context — not embedded in its system prompt where a sufficiently capable
model could reason about and circumvent it."

Therefore the generated KILLSWITCH.md and its JSONL audit MUST be placed
outside any path the agent can write, read for self-modification, or include
in its own prompt. The template banner states this; do not remove it.

## The KILLSWITCH.md field set (verbatim source definition)

"The KILLSWITCH.md open standard defines a file-based convention for
embedding emergency shutdown protocols into agent deployments:

- **TRIGGERS** (cost limits, error thresholds),
- **FORBIDDEN actions** (protected files, APIs, system commands),
- a three-level **ESCALATION** protocol (throttle to pause to full stop),
- append-only **JSONL audit** requirements, and
- **OVERRIDE** conditions requiring explicit human approval."

### TRIGGERS

Quantitative conditions that fire the killswitch. Source-grounded examples:
cost limits; error thresholds; circuit-breaker conditions — "if an agent
makes more than N tool calls in M minutes, or repeats the same tool call with
identical parameters three times, the system suspends and alerts." Every
trigger needs a numeric threshold; a trigger with no number is not a trigger.

### FORBIDDEN

Categorical actions never permitted regardless of allowlist state: protected
files, protected APIs, protected system commands. This is the supplementary
emergency brake — it is not the primary access model (that is the
default-deny Rego allowlist).

### ESCALATION — three levels, in order

1. **Throttle** — slow the agent (rate-limit tool calls) on a soft signal.
2. **Pause** — suspend and alert a human; await input.
3. **Full stop** — hard halt of the agent process.

The order is fixed: throttle, then pause, then full stop. Each level names
the threshold that promotes to the next and the human notified.

### Append-only JSONL audit schema

The audit MUST be append-only (no in-place edits or deletes) and live in
separated infrastructure. Decision-level logging shows "timestamps, decision
logic, tool usage, policy evaluations, and identity mappings." One JSON
object per line; suggested keys:

```
{"ts":"<RFC3339>","agent_id":"...","action":"...","decision":"allow|deny|throttle|pause|stop","policy_rule":"...","trigger":null,"blast_radius":"low|medium|high","reversible":true,"identity":"..."}
```

This satisfies record-keeping that "automatically logs events relevant to
risk identification across the system's lifecycle" and the requirement for
audit trails that "capture not just final outputs but intermediate agent
decisions." Integrate the stream into existing SIEM/SOAR so anomalous agent
behavior "is treated as a security event rather than a debugging problem."

### OVERRIDE

Conditions under which a human may resume or relax the killswitch — "OVERRIDE
conditions requiring explicit human approval." Record who can override, the
approval evidence required, and an audit entry for every override. An
override is never self-service by the agent.

## Why the killswitch is not optional

"93% of permission prompts are approved without reading," turning
high-frequency human-in-the-loop into "theatrical oversight rather than real
control." The shift to human-on-the-loop only works if "the telemetry and
alerting infrastructure that makes HOTL interventions timely and effective"
exists — the killswitch plus the JSONL audit is that backbone, and it is only
trustworthy if the agent cannot reach it.
