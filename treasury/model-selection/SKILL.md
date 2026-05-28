---
name: model-selection
description: >
  Assign each agent or stage in an agentic workflow a model and reasoning effort by criteria — a privacy/data-class gate, capability-to-role match, cost-per-successful-task, reasoning effort, native structured decoding, and a fallback behind a vendor-neutral abstraction layer — and record it as a model-assignment policy for human sign-off. Use when the user asks "which model should each agent use", "pick the LLM for this stage", "set the model tier and reasoning effort", or "are we overpaying for this agent's model". Produces a per-role assignment matrix (capability tier with dated example models, reasoning effort, cost posture, fallback, rationale) driven by criteria, not a single "best model". Do NOT use to operate a routing gateway or LLM proxy at runtime (use your router, e.g. OpenRouter / RouteLLM / LiteLLM) or for security threat-modeling of models and prompts (use a dedicated security threat-modeling review).
---

# model-selection

## Purpose
Assign each agent or stage in an agentic workflow a **model + reasoning effort** — by a privacy/data-class gate,
capability-to-role match, cost-per-successful-task, native structured decoding, and a fallback behind a vendor-neutral
abstraction layer — and record it as a **model-assignment policy**. This produces a criteria-driven matrix, not a single
"best model"; it recommends and explains, and a human signs off.

## When to use
- Triggers: *"which model should each agent use"*, *"pick the LLM for this stage"*, *"set the model tier and reasoning effort"*, *"are we overpaying for this agent's model"*.
- **Not this skill:** operating a routing gateway / LLM proxy at runtime → your router (e.g. OpenRouter / RouteLLM / LiteLLM); security threat-modeling of models or prompts → a dedicated security threat-modeling review.

## Input
- The workflow's stages / agent roles and their tasks, the **data sensitivity** each role handles, modality requirements
  (vision/audio), and latency/budget constraints (or the skill proposes a baseline from the role shapes).

## Output
A **Model-Assignment Policy** artifact + checklist + verdict. Skeleton:

```
# Model Assignment — <workflow / agent set>
Privacy gate:    <data class -> eligible model class: hosted-standard any · ZDR/BAA frontier-closed OK · strict/no-egress open-weight self-host only>
Per-role matrix: <role/agent -> capability tier (dated example models, e.g. as of 2026-05) -> reasoning effort -> context budget -> native structured decoding -> cost posture -> fallback -> rationale>
Cost-per-task:   <CPST = loaded cost (inference + retries + tool loops + review) ÷ verified successes; agentic loops grow ~O(n²); exhaust caching/batch before any model swap>
Routing:         <static per-role | cascade (cheap-first, escalate) | predictive — by latency tolerance + workload heterogeneity; target 30–50% at held quality>
Fallback:        <per role: alternate model · topology (same-provider/cross-region default) · idempotency + Retry-After · measured quality-loss budget · log the served model>
Re-eval cadence: <reassess quarterly behind a vendor-neutral abstraction layer; capability-per-$ moves fast>
Verdict:  ASSIGNMENT READY for human sign-off  |  GAPS: <…>
```

## Decision rules
1. **Gate on data class first; capability is the last filter.** Classify the data each role handles and eliminate any model that fails a privacy/residency non-negotiable for that class — hosted-standard → any model; ZDR + BAA → frontier-closed OK (contractual); single-tenant VPC → frontier-closed via VPC; strict / no-egress → open-weight self-hostable only. Then pick the cheapest model that clears the role's capability bar. Prefer architectural guarantees (on-device, VPC, self-host) over contractual ones when the threat model includes the provider; open weights ≠ privacy unless you also own the safeguards.
2. **Match capability to the role; don't maximize it.** Frontier-reasoning tier for planning/orchestration, judgment/review, and high-stakes or risk steps; mid-tier for high-volume implementation and structured design; cheap/fast tier for bulk extraction, summarization, and bounded schema-constrained tool calls. The dominant failure is a **planner↔executor capability mismatch** — watch the gap; set reviewer ≥ writer tier and prefer a different model family to suppress self-preference bias.
3. **Set reasoning effort per task, by difficulty — and weigh inference-time compute against base capability.** A mid-tier model at **high** effort often outperforms a frontier model at **zero** effort on logic/code tasks, at lower cost — so size effort and tier together, not separately. Thinking tokens bill as output and inflate latency ~5–30×; raise effort only for hard, verifiable-answer tasks (math/code) and keep it low on simple or mechanical steps, where long traces add cost and can *degrade* accuracy (inverse scaling). Never apply a uniform reasoning budget.
4. **Demand native structured decoding for executors.** For agents that extract data or fill schemas, assign models that natively enforce strict JSON-schema decoding at the decoding layer. Do not burn CPST on LLM self-correction or regex retry loops to coerce malformed output.
5. **Budget in cost-per-successful-task, not token price.** Account for retries, tool loops, and super-linear (~O(n²)) context growth in agentic loops; define and independently audit "success" as the denominator (completion ≠ correctness). Exhaust prompt-caching (with prefix discipline) and batching before recommending any model swap.
6. **Specify routing and fallback explicitly.** Choose static / cascade (cheap-first, escalate) / predictive by latency tolerance and workload heterogeneity; target a defensible 30–50% saving at held quality, not the 85–98% best case, and stage 5–10% of traffic first. For fallback: default same-provider/cross-region, keep retries idempotent, honor `Retry-After`, distinguish terminal vs retryable errors, and budget the secondary's measured quality loss. Log the actually-served model — silent substitution is a contract change.
7. **Model choice is not a safety control, and the matrix decays.** Prompt injection is unmitigated at the model layer, so high-stakes or autonomous actions need architectural guardrails (privilege isolation, output validation, hard limits, human approval) regardless of the model. Cite specific models only as dated examples; commit to tiers and ratios behind a vendor-neutral abstraction layer, validate the shortlist with a private blind eval, and attach a quarterly re-evaluation date.

## Checklist
- [ ] Data class per role classified; privacy/residency gate applied before capability
- [ ] Each role mapped to a capability tier (dated example models) + reasoning effort + rationale
- [ ] Planner↔executor capability gap checked; reviewer ≥ writer tier (different family)
- [ ] Reasoning effort set per task by difficulty, not uniform
- [ ] Inference-time compute vs base-capability tradeoff evaluated (mid+high effort vs frontier+low)
- [ ] Native structured decoding (strict JSON schema mode) assigned for extraction/schema tasks
- [ ] Budgeted in cost-per-successful-task; success criterion defined + independently audited
- [ ] Caching + batching applied before any model swap
- [ ] Routing regime chosen (static/cascade/predictive) with a realistic savings target
- [ ] Fallback per role: alternate model, idempotency, Retry-After, quality-loss budget; served model logged
- [ ] Architectural guardrails (not model choice) cover high-stakes actions
- [ ] Models cited as dated examples; vendor-neutral abstraction layer + quarterly re-eval set

## Anti-patterns (never do)
- Answer "what's the best model?" with a single winner — emit a role × privacy-tier matrix instead.
- Let a capability score override a privacy / data-residency gate.
- Run one model in every role — it overspends on simple tasks and underperforms on hard ones.
- Apply a uniform reasoning budget — high effort on trivial steps wastes cost and can lower accuracy.
- Assume a frontier model at zero reasoning effort beats a mid-tier model at high effort on a pure logic task.
- Burn tokens on regex retry loops instead of using native strict structured decoding for schema-filling agents.
- Choose on token price instead of cost-per-successful-task, or swap models before exhausting caching/batching.
- Treat a stronger model as a safety control for a high-stakes or autonomous action.
- Fall back to a weaker model silently — unmeasured quality regression hides from availability dashboards.
- Hardcode a specific model as a permanent rule, or bind the workflow to one vendor's retrieval schema with no abstraction layer.

## Example
**Input:** a 13-stage agentic delivery squad, payments domain. **Output (excerpt):** privacy gate → regulated data →
frontier-closed via ZDR/BAA or VPC (no public tier). *Matrix:* orchestrator / architect / reviewer / risk →
frontier-reasoning (e.g. as of 2026-05: Opus 4.7 / GPT-5.5 / Qwen 3.7-Max / Gemini 3.1 Pro), **high** effort; coding / contract / test-design → mid-tier
(e.g. Sonnet 4.6 / GPT-5.5 Instant / Grok 4.3 / Gemini 3.5 Flash), **medium** — a mid-tier at high effort beats frontier-at-zero on these logic tasks for less; requirement drafting / extraction → cheap-fast (e.g. Haiku 4.5 /
Gemini 3.5 Flash Lite / DeepSeek V4), **low** effort, with **strict JSON decoding** enabled. Budget in CPST; caching + batch first; cascade routing for the cheap tier with
same-provider fallback + `Retry-After`. **Verdict: ASSIGNMENT READY** once the tiers are blind-eval'd on
representative prompts and a quarterly re-eval date is set.

## Human approval gate
**Stop.** A human signs off the model-assignment policy; one-way-door choices (vendor commitment, deployment boundary)
are decide-then-explain with explicit rationale. The skill recommends the assignment and its rationale — it does not
route, deploy, or operate models.
