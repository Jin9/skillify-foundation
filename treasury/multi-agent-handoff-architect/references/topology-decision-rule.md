# Single-agent-first and the topology decision rule

## The decision rule the field has converged on

> The decision rule the field has converged on: start with one agent;
> escalate to a squad only when sub-problems are independent, sub-agent
> isolation prevents context pollution, and a hard quality gate (schema +
> verifier) is available.

The skill applies this as a four-step gate before designing any
multi-agent handoff. Every unnecessary handoff is a risk vector and a
token cost.

### Step 1 — Is it a deterministic 90%+ flowchart? Then it is not an agent.

If the task can be expressed as a deterministic flowchart with clear
if/else branches covering 90%+ of cases, do not build an agent at all —
use a workflow with simple classifiers at decision points. Workflows
offer predictability for well-defined tasks; agents trade latency and
cost for flexibility. A non-trivial fraction of "multi-agent"
deployments are actually this category. (Single linear-pipeline stage
composition is the job of `composing-agent-pipelines`, not this skill.)

### Step 2 — Single agent first.

If the task needs an agent, start with one. Single agents work best when
tasks are tightly coupled and mostly sequential, when global context
matters (step one affects step five), when fewer than 20 tools are
needed, and under strict budget/latency constraints. Writes stay
single-threaded; additional agents may contribute intelligence (review,
critique) but not act.

Empirically the single-agent baseline is strong: a single agent
succeeded 28 of 28 times on one benchmark while hierarchical multi-agent
organizations failed 36% of the time and self-organized swarms failed
68%. On strictly sequential reasoning, multi-agent systems degraded
performance by 39–70%. Multi-agent systems also consumed 4–220x more
tokens. See `references/handoff-evidence.md`.

### Step 3 — Escalate only when the task is parallelizable, read-heavy, with independent sub-problems.

Escalate to a squad only when the task has structural properties a squad
addresses: parallelizable read-heavy work with independent sub-problems
(research fan-out, log triage, multi-source enrichment), where an
orchestrator can spawn isolated subagents whose explorations would
otherwise pollute the main context. Parallelize independent searches,
serialize dependent reads; parallel dispatch requires three or more
unrelated tasks with no shared state and clear boundaries. On
parallelizable tasks multi-agent coordination produced +81% over
single-agent; the orchestrator-worker pattern outperformed the
single-agent baseline by more than 90% at roughly 15x tokens.

### Step 4 — Only with a hard quality gate / verifier.

Only use independent multi-agent systems if a hard quality gate is in
place: structured contracts, explicit schemas, and a verifier that can
**reject** outputs rather than **blend** them. Without a verifier the
squad's outputs are an averaged opinion, and averages can be wrong
without being detected.

## Default topology: orchestrator-worker

The dominant production topology in 2026 is orchestrator-worker: a
central LLM analyzes each unique task, dynamically determines subtasks,
and delegates them to specialist worker LLMs running in isolated context
windows, then synthesizes their outputs. This is the skill's default
when a squad is justified. A four-role coding squad (manager,
researcher, engineer, reviewer) reached 72.2% on SWE-bench Verified — a
+7.2-point gain over the same-model single-agent baseline attributed to
team structure. Hierarchical/supervisor is an idiomatic variant; swarm,
peer-to-peer, blackboard, and graph-of-agents exist but are escalations
to justify explicitly, not defaults.

## The three context-transfer strategies and the converged choice

1. **Full context forwarding** — every prior message passed verbatim.
   Token cost scales quadratically with handoff depth; a 50-message
   thread with 4 handoffs means the fifth agent processes ~200 messages.
   Rejected as a default.
2. **Structured context objects** — a typed object (task ID, detected
   intent, extracted entities, resolution status); only relevant fields
   passed. Typical size 200–500 tokens versus 5,000–20,000 for full
   forwarding. **This is the converged choice** — adopted independently
   by Anthropic, OpenAI, AutoGen, and LangChain.
3. **Summarized context** — an LLM compresses prior conversation at each
   handoff. Reduces token count 70–90% but introduces information loss
   and adds 500ms–1.5s latency per handoff. Use only when the structured
   object cannot capture what the receiver needs.

The skill specifies strategy 2 by default and sizes the contract's
`state` object at 200–500 tokens.
