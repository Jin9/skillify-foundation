---
name: agentic-workflow-design
description: >
  Design how an AI-agent pipeline is supervised — its stages and owners, the human-in-the-loop approval gates, the never-do guardrails and command-safety policy, accountability for AI-generated output, and observability — and record it as an agentic-workflow design for human sign-off. Use when the user asks "design our agentic CI/CD or delivery pipeline", "where should the human approval gates go for our AI agents", "set guardrails for our autonomous coding agents", or "how do we keep accountability for AI-generated work". Produces a workflow-design artifact that maps stages to owners, places gates by reversibility and blast radius, defines an allow/confirm/deny command-safety policy, and specifies a per-stage observability and handoff contract. Do NOT use to actually run, deploy, or operate the pipeline (that is your orchestration runtime or CI platform) or for security threat-modeling of the agents (reviewing-software-security).
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
1. **Gate by reversibility and blast radius, never uniformly.** Auto-approve reversible/read-only work; gate irreversible, control-plane, or production-touching actions on **synchronous named-human approval regardless of agent confidence**. Use plan-then-execute or per-stage checkpoints for high-stakes flows; reserve exception-only autonomy for reversible, low-stakes work.
2. **Enforce policy outside the model.** Output guardrails constrain what an agent *says*; an external policy engine with **allow / confirm / deny** tiers at the tool-calling layer constrains what it may *do* — so a hijacked or drifting model cannot override the controls. The agent does not decide what is allowed.
3. **Apply least agency.** Scope each agent's tools, permissions, and credentials to the one task: default read-only, short-lived credentials (not long-lived secrets), allowlisted tool/MCP servers, and read-external isolated from write-sensitive roles. Treat untrusted content (PR titles, issue bodies, tool descriptions) as data, never instructions.
4. **Default to a single agent; justify every handoff boundary.** Add an agent only for genuine parallelism, a context-window limit, or security-domain separation — each unnecessary handoff costs tokens and adds a failure vector. Get decomposition granularity and topology right before adding agents; unstructured "bag of agents" amplifies errors.
5. **Set hard pre-execution limits** (retry cap, spend budget, iteration ceiling, loop detector) that **terminate** the run — not dashboards that alert after the spend.
6. **Model the workflow as a state machine with explicit failure paths.** Classify transient versus fundamental failures; checkpoint at each step to resume (not restart); use idempotency keys for side-effecting actions; make human escalation the final recovery tier.
7. **Keep the human the accountable owner of record.** The pipeline produces artifacts and recommends; it never decides, commits, or deploys. Log **both halves of the loop** (model request + tool effect) under one trace id, tamper-evident, with the model version and a distinct agent identity; prefer signed provenance over anthropomorphic co-authorship; treat reasoning traces as context, not evidence.

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
- Let an agent take an irreversible, control-plane, or production action without a human gate — regardless of its confidence.
- Rely on prompt/output guardrails alone; ship without an external policy engine enforcing allow / confirm / deny.
- Over-provision agents (broad permissions, long-lived credentials, unrestricted tools/MCP) — least agency or nothing.
- Treat untrusted content (PR titles, issue bodies, tool descriptions) as trusted instructions.
- Add agents or handoffs without justification, or forward full history instead of a typed handoff contract.
- Govern cost/loops with alert-only dashboards and no pre-execution caps that terminate the run.
- Ship without logging the tool-effect half of the loop, or with a silently-editable (non-tamper-evident) audit log.
- Let the agent decide architecture alone, approve its own PRs, fix production directly, or modify its own permission config.
- Attribute AI output as an accountable co-author, or leave AI-generated output without a named human owner of record.

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
