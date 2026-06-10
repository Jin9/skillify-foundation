---
name: agentic-workflow-design
description: >
  Design how an AI-agent pipeline is supervised — its stages and owners, the human-in-the-loop approval gates, the never-do guardrails and command-safety policy, accountability for AI-generated output, and observability — and record it as an agentic-workflow design for human sign-off. Use when the user asks "design our agentic CI/CD or delivery pipeline", "where should the human approval gates go for our AI agents", "set guardrails for our autonomous coding agents", or "how do we keep accountability for AI-generated work". Produces a workflow-design artifact that maps stages to owners, places gates by reversibility and blast radius, defines an allow/confirm/deny command-safety policy, and specifies a per-stage observability and handoff contract. Do NOT use to actually run, deploy, or operate the pipeline (that is your orchestration runtime or CI platform), for security threat-modeling of the agents (reviewing-software-security), or for the handoff payload contract (multi-agent-handoff-architect).
---

# agentic-workflow-design

## Purpose
Design and document how an agentic (AI-driven) workflow is supervised — pipeline stages, human-in-the-loop gates,
never-do guardrails, accountability for AI output, and observability — so the team can run it with a human accountable
at every irreversible step. This is the **supervision-design artifact**: it recommends and records; it never operates
the pipeline.

## When to use
- Triggers: *"design our agentic CI/CD pipeline"*, *"where do the human approval gates go for our AI agents"*, *"set guardrails for our autonomous coding agents"*, *"how do we keep accountability for AI-generated work"*.
- **Not this skill:** running or operating the pipeline itself → your orchestration runtime / CI platform; security threat-modeling of the agents → `reviewing-software-security`.

## Input
- The agent pipeline or initiative: what the agents do, where AI sits in the loop, the tools/credentials/data they touch,
  and which actions are reversible versus irreversible (control-plane / production-touching).

## Output
An **Agentic Workflow Design** artifact + checklist + verdict. Skeleton:

```
# Agentic Workflow Design — <pipeline / initiative>
Goal & autonomy level:  <what it does · L1 assist … L5 autonomous · HITL (approve each) or HOTL (monitor + exception)>
Pipeline stages:        <stage -> agent/tool -> input -> output -> human owner>
Topology:               <single-agent (default) | orchestrator-worker | pipeline | … — why each handoff is justified>
Approval gates:         <action -> reversibility / blast radius -> gate (AUTO | async review | sync named approval) -> who>
Command-safety policy:  ALLOW <reversible, sandboxed>  |  CONFIRM <irreversible / control-plane>  |  DENY <outside allowlist>
Never-do guardrails:    <bounded "agent must NEVER" set — enforced in the harness / policy engine, not the prompt>
Accountability map:     <AI output -> distinct agent identity -> accountable human owner of record>
Observability:          <per stage: model call · tool effect · handoff — one correlation/trace id; tamper-evident; model version>
Handoff contract:       <task id · intent · state · confidence · provenance/trace id · schemaVersion — exactly one owner mutates state>
Failure & recovery:     <state machine: transitions + failure paths · retry/budget/iteration caps · checkpoint · HITL escalation>
Verdict:  DESIGN READY for human sign-off  |  GAPS: <…>
```

## Decision rules
Apply the seven supervision rules: gate by reversibility/blast radius; enforce policy outside the model; least agency; single-agent default + justified handoffs; hard pre-execution caps that terminate; state machine with failure paths; human as accountable owner of record.
Detailed rules: references/supervision-design-rules.md

## Checklist
- [ ] Each stage maps to an agent/tool and a named human owner
- [ ] Topology chosen; every handoff boundary justified (single-agent default)
- [ ] Gates placed by reversibility / blast radius; irreversible / control-plane actions require sync named approval
- [ ] Command-safety policy (allow / confirm / deny) enforced outside the model
- [ ] Least agency: scoped tools, short-lived credentials, allowlist, read/write role isolation
- [ ] Never-do guardrail set defined and harness-enforced
- [ ] Pre-execution caps set (retries, budget, iterations, loop detector)
- [ ] Observability: model call + tool effect + handoff under one trace id, tamper-evident, model version recorded
- [ ] Handoff payload is a versioned typed contract; exactly one owner mutates shared state
- [ ] State machine has explicit failure paths + checkpoint + HITL escalation
- [ ] Accountability map: distinct agent identity → human owner of record

## Anti-patterns (never do)
Reject any design that: ungated irreversible/control-plane actions; prompt-only guardrails; over-provisioned agents; untrusted content as instructions; unjustified handoffs / full-history forwarding; alert-only cost governance; missing tool-effect logging or non-tamper-evident audit; agent self-approval or self-config; AI as accountable co-author or output without a named human owner.
Full list: references/anti-patterns.md

## Example
**Input:** an agentic CI/CD squad — orchestrator coordinating code-writer, test-writer, and reviewer agents into CI/CD.
**Output (excerpt):** single-agent default rejected for *genuine* parallelism (code vs. tests) → orchestrator-worker
with three scoped specialists. *Command-safety:* ALLOW run-tests-in-sandbox · CONFIRM push-to-main / deploy /
schema-change · DENY non-allowlisted MCP. *Gates:* PR review = async human; production deploy = **sync named
approval**. *Caps:* 5 retries / $2 / 50 iterations + loop detector (terminate). *Observability:* one trace id across
handoffs, tool-effect logged at the MCP gateway, hash-chained. *Handoff:* typed schema (task id · intent · state ·
confidence · trace id). *Failure:* transient → backoff; exhausted → dead-letter + HITL. **Verdict: CONDITIONAL** —
DESIGN READY once the deploy gate names an approver and credentials move to 1-hour OIDC tokens.

## Human approval gate
**Stop.** A human owns and signs off the supervision design — and remains the **accountable owner of record** for
everything the pipeline produces. This skill drafts and recommends the workflow; it never runs, deploys, or operates the
pipeline, and never relaxes a gate on the agent's behalf.
