---
name: model-selection
description: >
  Assign each agent or stage in an agentic workflow a model and reasoning effort by criteria — a privacy/data-class gate, capability-to-role match, cost-per-successful-task, reasoning effort, structured decoding, and a fallback behind a vendor-neutral abstraction layer. Use when the user asks "which model should each agent use", "pick the LLM for this stage", "set the model tier and reasoning effort", or "are we overpaying for this agent's model". Produces a per-role assignment matrix driven by criteria, not a single "best model". Do NOT use to operate a routing gateway or LLM proxy at runtime.
---

# model-selection

## Purpose
Assign each agent or stage in an agentic workflow a **model + reasoning effort** — by a privacy/data-class gate,
capability-to-role match, cost-per-successful-task, structured output native support, and a fallback behind a vendor-neutral abstraction layer — and
record it as a **model-assignment policy**. This produces a criteria-driven matrix, not a single "best model"; it
recommends and explains, and a human signs off.

## When to use
- Triggers: *"which model should each agent use"*, *"pick the LLM for this stage"*, *"set the model tier and reasoning effort"*, *"are we overpaying for this agent's model"*.
- **Not this skill:** operating a routing gateway / LLM proxy at runtime → your router (e.g. OpenRouter / RouteLLM / LiteLLM); security threat-modeling of models or prompts → `reviewing-software-security`.

## Input
- The workflow's stages / agent roles and their tasks, the **data sensitivity** each role handles, modality requirements (vision/audio), and latency/budget constraints.

## Output
A **Model-Assignment Policy** artifact + checklist + verdict. Skeleton:

```
# Model Assignment — <workflow / agent set>
Privacy gate:    <data class -> eligible model class: hosted-standard any · ZDR/BAA frontier-closed OK · strict/no-egress open-weight self-host only>
Per-role matrix: <role/agent -> capability tier (dated example models) -> reasoning effort -> context budget -> native structured decoding -> cost posture -> fallback -> rationale>
Cost-per-task:   <CPST = loaded cost (inference + retries + tool loops + review) ÷ verified successes; exhaust prefix-caching before any model swap>
Routing:         <static per-role | cascade (cheap-first, escalate) | predictive — by latency tolerance + workload heterogeneity; target 30–50% at held quality>
Fallback:        <per role: alternate model · topology (same-provider/cross-region default) · idempotency + Retry-After · measured quality-loss budget · log the served model>
Re-eval cadence: <reassess quarterly behind a vendor-neutral abstraction layer; capability-per-$ moves fast>
Verdict:  ASSIGNMENT READY for human sign-off  |  GAPS: <…>
```

## Decision rules
1. **Gate on data class first; capability is the last filter.** Classify the data each role handles and eliminate any model that fails a privacy/residency non-negotiable for that class — hosted-standard → any model; ZDR + BAA → frontier-closed OK (contractual); single-tenant VPC → frontier-closed via VPC; strict / no-egress → open-weight self-hostable only. Then pick the cheapest model that clears the role's capability bar.
2. **Match capability to the role; don't maximize it.** Frontier-reasoning tier for planning/orchestration, judgment/review, and high-stakes or risk steps; mid-tier for high-volume implementation and structured design; cheap/fast tier for bulk extraction, summarization, and bounded schema-constrained tool calls.
3. **Inference-Time Compute Scaling.** Mid-tier models with *High* reasoning effort often outperform Frontier models with *Zero* reasoning effort on logic/coding tasks, but at a lower cost. Set reasoning effort per task by difficulty. Thinking tokens bill as output and inflate latency ~5–30×; raise effort only for hard, verifiable-answer tasks and keep it low on mechanical steps (inverse scaling).
4. **Demand Native Structured Decoding for Executors.** For agents extracting data or filling schemas, always assign models that natively support strict JSON schema enforcement at the decoding layer. Do not burn CPST on LLM self-correction or regex retry loops.
5. **Budget in cost-per-successful-task (CPST), not token price.** Account for retries, tool loops, and super-linear (~O(n²)) context growth in agentic loops. Exhaust prompt-caching (with strict prefix discipline) and batching before recommending any model swap.
6. **Specify routing and fallback explicitly.** Choose static / cascade (cheap-first, escalate) / predictive by latency tolerance and workload heterogeneity. For fallback: default same-provider/cross-region, keep retries idempotent, honor `Retry-After`, distinguish terminal vs retryable errors. Log the actually-served model — silent substitution is a contract change.
7. **Model choice is not a safety control, and the matrix decays.** Prompt injection is unmitigated at the model layer, so high-stakes or autonomous actions need architectural guardrails. Cite specific models only as dated examples; commit to tiers and ratios behind a vendor-neutral abstraction layer, and attach a quarterly re-evaluation date.

## Checklist
- [ ] Data class per role classified; privacy/residency gate applied before capability
- [ ] Each role mapped to a capability tier (dated example models) + reasoning effort + rationale
- [ ] Inference-time compute vs Base-capability tradeoff evaluated
- [ ] Native structured decoding (JSON strict mode) assigned for extraction/schema tasks
- [ ] Planner↔executor capability gap checked; reviewer ≥ writer tier (different family)
- [ ] Reasoning effort set per task by difficulty, not uniform
- [ ] Budgeted in cost-per-successful-task; prefix-caching applied before any model swap
- [ ] Routing regime chosen (static/cascade/predictive) with realistic savings target
- [ ] Fallback per role: alternate model, idempotency, Retry-After; served model logged
- [ ] Models cited as dated examples; vendor-neutral abstraction layer + quarterly re-eval set

## Anti-patterns (never do)
- Answer "what's the best model?" with a single winner — emit a role × privacy-tier matrix instead.
- Burn tokens on regex retry loops instead of utilizing native strict structured decoding.
- Assume a Frontier model with no reasoning effort will beat a Mid-tier model with max reasoning effort on a pure logic task.
- Run one model in every role — it overspends on simple tasks and underperforms on hard ones.
- Choose on token price instead of cost-per-successful-task, or swap models before exhausting caching/batching.
- Treat a stronger model as a safety control for a high-stakes or autonomous action.
- Fall back to a weaker model silently — unmeasured quality regression hides from availability dashboards.

## Example
**Input:** a 13-stage agentic delivery squad, payments domain. **Output (excerpt):** privacy gate → regulated data →
frontier-closed via ZDR/BAA or VPC (no public tier). *Matrix:* orchestrator / architect / reviewer / risk →
frontier-reasoning (e.g. Opus 4.7 / GPT-5.5 / Qwen 3.7-Max / Gemini 3.5 Pro), **high** effort; coding / contract / test-design → mid-tier
(e.g. Sonnet 4.6 / GPT-5.5 Instant / Grok 4.3 / SubQ 1M-Preview), **medium** effort; requirement drafting / extraction → cheap-fast (e.g. Haiku 4.5 / Gemini 3.5 Flash Lite / ZAYA1-8B), **low** effort, with **strict JSON decoding** enabled. Budget in CPST; prefix-caching + batch first; cascade routing for the cheap tier with same-provider fallback + `Retry-After`. **Verdict: ASSIGNMENT READY**.

## Human approval gate
**Stop.** A human signs off the model-assignment policy; one-way-door choices (vendor commitment, deployment boundary)
are decide-then-explain with explicit rationale. The skill recommends the assignment and its rationale — it does not
route, deploy, or operate models.
