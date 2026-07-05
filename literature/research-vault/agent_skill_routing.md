# Agent skill routing

Source: ResearchVault run agent-skill-routing-20260526-162805 (local deep-research pipeline)
Accessed: 2026-05-26
Category: research-vault / skill routing/triggering
Provenance: harvested 2026-07-05 from ResearchVault 05-final_report.md (run executed 2026-05-26); inline bibliography preserved

## Why This Source Matters

Internal deep-research synthesis on how hosts route requests to skills — description matching, trigger phrasing, and failure modes. Directly relevant to skillify's trigger-quality rubric dimension. Internal-synthesis tier.


## Executive Summary

Agent skill routing is the decision layer that selects which skill, tool, or capability an agentic system should invoke for a given task. It is best understood as two coupled problems: **matching** (recall — surfacing the candidate skills that could plausibly serve the request) and **scoring, precedence, and disambiguation** (ranking — choosing one when several candidates compete) [1][8]. Modern agent stacks implement matching through a small set of primitives: embedding/semantic retrieval, LLM function-calling where the model itself reads tool descriptions and decides, hybrid dense-plus-sparse retrieval, and lightweight classifier or rule-based routers [4][5][6]. Each primitive trades cost, latency, and adaptability differently, and production systems increasingly layer them rather than choosing one [1][31].

The hardest part of routing is not recall but ranking. When two tool descriptions overlap, the selection signal becomes ambiguous and the model resolves it arbitrarily; the dominant decision signals are, in order, description text, parameter names, and position in the context window [8][9]. This makes routing systematically biased: large language models favor tools by superficial metadata and by listing position rather than by utility, and they rarely admit uncertainty — instead silently selecting a plausible wrong tool [20][8]. Precedence therefore has to be authored explicitly; absent numbered priority rules, the model will not infer a stable order [11][9].

Routing quality degrades sharply with scale. The Berkeley Function Calling Leaderboard documents accuracy on a calendar-scheduling task collapsing from 43% to 2% as the tool count grew from 4 to 51, while each tool definition imposes fixed per-turn token overhead [18][19]. The remedy that the field has converged on is retrieval-first, hierarchical routing: decouple tool discovery from generation, narrow the candidate set, and only then let the model choose. RAG-MCP reports tool-selection accuracy more than tripling — from 13.62% to 43.13% — when a retrieval stage replaces loading every tool into context [23]. Vendor platforms have implemented the same idea: Anthropic's Tool Search Tool defers tool loading so an agent can address thousands of tools while seeing only the few it needs [15]. The practical guidance that follows is to treat routing as a cost-tiered cascade — fast heuristics, then classifiers, then LLM reasoning only at genuinely ambiguous branch points — and to start with deterministic routing, escalating to model-driven routing only where it earns its inference cost [26][31][28].

## Background

As agentic systems have moved from single-prompt assistants to orchestrated workflows with dozens or hundreds of callable capabilities, the question of *which* capability to invoke for a task has become a first-class engineering concern. The terminology is not fully settled. Practitioners increasingly distinguish a **tool** — a single callable capability such as search, SQL, or a calculator — from a **skill** — a reusable mini-workflow that coordinates multiple tool calls with policy, guardrails, and output structure, placing tools in an execution layer and skills in an orchestration layer [1]. This report uses "skill routing" in the broad sense the topic intends: the act of mapping an incoming task to the right capability, whether that capability is an atomic tool, a composite skill, a sub-agent, or a backend model. The decision mechanism is largely shared across those granularities.

The audience here is assumed to be expert: familiar with embeddings, function-calling, and the basic shape of retrieval-augmented systems. The aim is to be precise about mechanism and to surface measured outcomes and live disagreements rather than to re-explain foundational concepts.

## Methodology

This report synthesizes eight sub-questions spanning the definition and scope of routing, its matching mechanisms, its scoring and precedence logic, the landscape of production implementations, the empirical evidence on routing performance, the failure modes and constraints, the techniques for routing at large scale, and the practical design guidance that follows. Sources span peer-reviewed and preprint research (the Berkeley Function Calling Leaderboard, RouteLLM, RAG-MCP, MCP-Zero, Tool-to-Agent Retrieval, BiasBusters), primary vendor documentation (Anthropic, LangChain), security advisories (OWASP, Invariant Labs), and practitioner engineering analyses. Where a claim is quantitative — an accuracy delta, a cost reduction, a token figure — it is attributed to the source that measured it. Where the field disagrees — most sharply on whether to route deterministically or with an LLM — the disagreement is presented as live, not flattened into a recommendation.

## Key Findings

- Routing decomposes into **matching** (recall) and **scoring/precedence/disambiguation** (ranking); scoping selection to a sub-agent rather than evaluating all skills globally improves precision [3][8].
- The matching primitives are embedding/semantic retrieval, LLM function-calling, hybrid dense+sparse retrieval, and classifier/rule routers; hybrid retrieval tends to beat pure semantic search because tool docs carry exact terminology [4][5][6].
- The model's tool-selection signal runs in order: description text, then parameter names, then context position; overlapping descriptions cause arbitrary disambiguation [8][9].
- Tool selection is systematically biased toward earlier-listed tools and toward superficial metadata, not utility [20].
- Accuracy collapses as the catalog grows (43%→2% from 4→51 tools on one BFCL task), and every tool definition costs tokens on every turn [18][19].
- Retrieval-first, hierarchical routing is the convergent fix: RAG-MCP tripled selection accuracy versus loading all tools; Anthropic, MCP-Zero, and Tool-to-Agent Retrieval all implement two-stage narrowing [23][15][24][25].
- The deterministic-vs-LLM routing choice is a genuine trade-off; the dominant production answer is a cost-tiered hybrid cascade [26][31][28].

## What Skill Routing Is, and What It Is Not

Skill routing sits between intent and execution. Upstream of it, the system has some representation of what the user wants; downstream of it, a concrete capability runs. The routing decision is the choice of *which* capability. It is useful to separate this from three adjacent concepts it is often conflated with. **Intent classification** is a narrower operation — mapping a request to a domain or category — and is frequently used *as* the first stage of routing rather than being identical to it: a lightweight classifier maps the incoming request to a tool category, after which a finer matching stage runs within that category [1]. **Planning** decides the *sequence* of capabilities for a multi-step task; routing decides the single next capability. **Orchestration** is the broader control loop that invokes routing repeatedly, manages state, and handles handoffs between sub-agents.

A second clarifying distinction is the unit being routed to. A tool is atomic; a skill is composite. Anthropic frames Agent Skills as organized folders of instructions, scripts, and resources that an agent can discover and load dynamically, with the model automatically using a skill when it is relevant to the request [2]. The routing decision — discover, match, select, load — is structurally the same whether the target is a tool, a skill, a sub-agent, or a backend model; only the granularity and the metadata available for matching change.

Scope matters for accuracy. Evaluating every skill in a global registry against every request maximizes ambiguity. Scoping intent recognition and skill selection to a specific agent, rather than evaluating all skills globally, significantly reduces ambiguity and improves precision [3]. This is the conceptual seed of the hierarchical routing patterns discussed later: narrow first, then choose.

```
            ┌──────────────────────────────────────────┐
   task ──► │  ROUTING DECISION                          │
            │                                            │
            │  matching (recall)   ─►  scoring (rank)     │
            │  surface candidates      pick + precedence  │
            └───────┬────────────────────────┬───────────┘
                    │                         │
            intent  │                         │  selected
          classify ─┘                         └─► capability ─► execute
         (optional first stage)
```

The diagram makes the load-bearing structure explicit: routing is a two-phase pipeline, optionally fronted by intent classification, that turns a task into a selected capability. Everything that follows is an elaboration of one of these phases.

## Matching Mechanisms: How a Task Meets a Candidate Skill

Four matching primitives dominate. They are not mutually exclusive — production systems combine them — but they are mechanistically distinct.

**Embedding / semantic retrieval.** The request is encoded into a vector and compared against vectors representing each skill (usually built from the skill's name and description). Semantic tool selection uses similarity to surface relevant tools *before the LLM even sees the request*, with semantic routers acting as lightweight proxies that make routing decisions before forwarding to backend model pools [4]. A semantic router encodes prompts into high-dimensional embeddings processed by an intent classifier to choose a route, giving a practical balance of accuracy and efficiency for production systems [7]. The advantage is that matching happens outside the expensive generation call; the cost is that pure semantic similarity can miss exact terminology.

**LLM function-calling (the model as router).** Here the model itself is the router. When processing a tool-calling request, the model reads each tool's description and computes which function best matches the current context, deciding when to call a tool based on the request and the tool's description [5]. This is the most flexible primitive — it can reason about the request — but it is also the most expensive, because every candidate's metadata must sit in the context window, and it is the one most exposed to the biases discussed below.

**Hybrid dense + sparse retrieval.** Because tool and API documentation contains exact terminology — method names, parameter names, domain jargon — hybrid retrieval that combines dense embeddings with sparse BM25/keyword matching consistently outperforms pure semantic search for tool selection [6]. The keyword channel catches the exact-string matches that embeddings blur; the dense channel catches conceptual matches that keywords miss.

**Classifier and rule-based routers.** A dedicated classifier (often BERT-based) or a hand-written rule maps the request to a route. These are fast, cheap, and deterministic, and a dedicated router — ML-based, rule-based, or embedding-based — often yields more robust and efficient routing for larger production systems than a general LLM agent [7].

```
                 ┌─────────────────────────────┐
   request ────► │   matching primitive          │
                 ├─────────────────────────────┤
                 │ semantic   → vector similarity │
                 │ function-  → LLM reads descs   │
                 │   calling                      │
                 │ hybrid     → dense + BM25      │
                 │ classifier → BERT / rules      │
                 └──────────────┬──────────────┘
                                │
                                ▼
                     candidate skill set (top-k)
```

The choice among primitives is the first major routing design decision, and the diagram's right column — cost and exactness rising from classifier up to function-calling — is exactly the axis the practical-guidance section turns into a cascade.

## Scoring, Precedence, and Disambiguation

Recall is the easy half. Once a candidate set exists, the system must *rank* it and break ties — and this is where routing quality is won or lost. The mechanics of how a model discriminates among candidates are now reasonably well characterized. Tool-selection decision signals run in a definite order: description text first, then parameter names, then tool ordering in the context window. When descriptions are distinct and scoped, the selection signal is clean; when they overlap, the model resolves the ambiguity arbitrarily [8]. The ambiguity is not resolved by some hidden principled rule — it is genuinely arbitrary, which means two functionally similar skills with overlapping descriptions will be selected almost by coin-flip [9].

Explicit scoring frameworks borrowed from evaluation transfer directly to ranking candidate skills. The LLM-as-a-Judge paradigm supports pointwise scoring (independent scores per candidate), pairwise comparison (two candidates at a time, ties allowed), and listwise/batch ranking, often driven by explicit structured rubrics [10]. A router can adopt any of these: score each skill independently, compare them head-to-head, or rank the whole shortlist at once. Pointwise is cheapest; listwise is most consistent but most expensive.

Precedence — the rule that decides *who wins* when scores are close — is the part that must be designed rather than assumed. When parallel results conflict, the model needs a defined resolution strategy: either surface the conflict to the user or apply an explicit priority rule [9]. And the priority rule has to be written down. When instructions conflict without explicit priority ordering, the model skips verification steps and rushes to action; the documented fix is to number priorities explicitly, e.g. "Priority 1: Tests pass. Priority 2: Under 5 minutes" [11]. The implication for skill routing is direct: a routing layer that hopes the model will infer the right precedence among overlapping skills is relying on arbitrary behavior. Precedence is authored, not emergent.

```
   candidate set
        │
        ▼
   ┌──────────┐   scores far apart   ┌──────────────┐
   │ score     ├────────────────────►│ select top    │
   │ pointwise │                      └──────────────┘
   │ pairwise  │   scores close (tie)  ┌──────────────┐
   │ listwise  ├────────────────────►│ precedence rule│
   └──────────┘                      │ or escalate    │
                                      └──────────────┘
```

The tie branch is the one most often left undesigned. The diagram's lower path — a tie routed into an explicit precedence rule or an escalation — is the difference between deterministic disambiguation and the arbitrary resolution the failure literature warns about.

## The Landscape of Production Routing Implementations

Concrete systems make the abstractions tangible, and they differ in instructive ways.

**LangGraph (deterministic graph routing).** LangGraph models an agent as a graph; routing happens at conditional edges. A conditional edge attaches a *router function* that reads the current state and returns a string naming the next node — and crucially, that function should not call an LLM, write to state, or produce side effects [12]. This is routing as pure, deterministic state inspection. Notably, recent LangGraph multi-agent implementations have largely switched from conditional edges to a `Command` object: a router node reads state, decides what runs next, and returns a `Command` that both updates state and names the next node, merging the route decision with the state transition [13]. The evolution shows the framework's routing API converging on a single combined "decide-and-move" primitive.

**RouteLLM (learned model routing).** RouteLLM is a framework for serving and evaluating routers that choose between a stronger and a weaker model to save cost without sacrificing quality. It trains four distinct router families on the same task — a similarity-weighted ranking router, a matrix-factorization model, a BERT classifier, and a causal-LLM classifier — and treats the routing decision as a binary classification between a strong and a weak model [14][31]. RouteLLM is "skill routing" in the model-selection sense: the capability being routed to is a model tier rather than a tool, but the matching-and-scoring machinery is identical.

**Anthropic's Tool Search Tool and Agent Skills (retrieval-first platform routing).** Anthropic's 2025 advanced tool-use features implement retrieval-first routing at the platform level. The Tool Search Tool discovers tools on-demand so that the model only sees the tools it needs, letting an agent work with hundreds or thousands of tools without loading all definitions into context, via a `defer_loading: true` flag [15]. Combined with Agent Skills — discoverable, dynamically loaded capability folders [2] — this is a vendor instantiation of the narrow-then-choose pattern, with deferred loading as the narrowing mechanism.

**Semantic routers (embedding routing as infrastructure).** The semantic-router pattern — exemplified by the vLLM Semantic Router and the broader class of router-as-proxy systems — sits in front of model pools and routes by embedding similarity plus a classifier, deciding the route before the request reaches an expensive model [4][29]. This is routing implemented as serving infrastructure rather than as application logic.

These four span the design space: deterministic state routing (LangGraph), learned model routing (RouteLLM), retrieval-first platform routing (Anthropic), and embedding-proxy routing (semantic routers). They share the matching-then-scoring skeleton; they differ in where the routing logic lives and how much of it is learned versus authored.

## What the Evidence Says About Routing Performance

The empirical picture is mixed in an informative way: routing can simultaneously improve accuracy and cut cost when the route choice is well-calibrated, but it degrades badly when the candidate space is large or when abstention is required.

On the positive side, RouteLLM's routers matched baseline performance with up to a 70% cost reduction on MT Bench; its matrix-factorization router reached 95% of GPT-4's quality while making only 26% of the calls to GPT-4 [14]. A semantic router for vLLM improved MMLU-Pro accuracy by 10.2 percentage points while cutting response latency 47.1% and token consumption 48.5% versus direct inference [7]. These are measured deltas, not vendor claims, and they establish that good routing is not merely a cost optimization — it can raise quality by sending each request to the capability best suited to it.

How is routing accuracy actually measured? The Berkeley Function Calling Leaderboard is the most widely cited tool-use benchmark; it grades tool calls by Abstract Syntax Tree comparison of *structure* rather than by executing them, spanning over 2,000 question-function-answer pairs and including categories for multiple-function selection and relevance detection — knowing when to refuse a tool call [16]. Tool-retrieval research complements this with information-retrieval metrics like Recall@K and NDCG@K to quantify whether the correct tool appears in the top-K candidate set, which is precisely the recall half of routing.

The sobering result is about abstention and context. On BFCL, top models ace one-shot tool questions but still stumble when they must remember context, manage long conversations, or decide when *not* to act [17]. Relevance detection — recognizing that no available skill fits and declining to route — is a distinct and harder competency than positive selection, and it is the one current systems are weakest at. A router that always picks *something* will confidently mis-route the requests that should have been refused.

## Failure Modes, Costs, and Constraints

Routing fails in characteristic ways, and the failures compound as the tool catalog grows.

**The scaling collapse.** The headline failure is that selection accuracy degrades predictably with catalog size. The Berkeley Function Calling Leaderboard found accuracy on a calendar-scheduling task dropping from 43% to just 2% when the number of tools expanded from 4 to 51 across multiple domains [18]. A model choosing among five clearly scoped tools substantially outperforms one scanning fifty. This is the single most important empirical constraint on routing design.

**Token overhead.** Even when routing is correct, it is not free. Tool definitions impose fixed per-turn token cost: detailed schemas running ~500 tokens each across ten tools cost ~5,000 tokens before the user asks anything, and that overhead is paid on every request whether or not the tool is used [19]. At scale, naive "load all tools" approaches make routing expensive independent of whether it is accurate.

**Silent wrong-tool selection.** The model rarely says it does not know which tool to use. Instead it picks a plausible-sounding wrong tool, combines multiple tools incorrectly, or calls a tool with parameters that fit a *different* tool's schema [8]. Because the failure is silent and confident, it evades the simplest monitoring; the agent appears to be working while routing to the wrong capability.

**Systematic bias.** Even when a correct tool exists in context, selection is not utility-rational. Evaluated across seven models, LLMs systematically favor providers by superficial metadata — names, descriptions, parameter schemas — and disproportionately favor tools that appear earlier in the context, a positional bias [20]. The mitigation proposed is to filter to a relevant subset and then sample uniformly, reducing bias while preserving coverage — itself an argument for a retrieval-first narrowing stage that decouples selection from raw ordering.

**Security: the routing layer as attack surface.** Because routing reads tool metadata into the model's context to make its decision, that metadata is an injection vector. Tool poisoning is a subset of indirect prompt injection in which malicious instructions are embedded in tool descriptions or metadata that influence agent behavior *even if the poisoned tool is never invoked* [21]. When an agent connects to an MCP server it requests `tools/list`; the server returns names and descriptions that enter the model's context and are used to decide which tools to invoke — and those descriptions can carry hidden instructions [22]. The routing decision is therefore not just a quality problem but a security boundary: the data the router consumes to choose is attacker-controllable.

```
   benign description ──┐
                        ├─► router reads metadata ─► decision
   poisoned description ┘        (context)              │
        (hidden instr.)                                 ▼
                                            behavior altered even if
                                            poisoned tool never called
```

The diagram captures why tool poisoning is insidious: the harm is realized at the moment metadata enters context for the routing decision, not at invocation, so refusing to *call* the malicious tool does not neutralize it.

## Routing at Scale: Hierarchy, Retrieval, and Progressive Disclosure

The scaling collapse and the token-overhead constraint point to the same architectural answer, and the field has converged on it: do not present the whole tool catalog to the model. Treat routing as retrieval over a large capability space, narrow to a small candidate set, and only then let the model choose.

RAG-MCP makes the case quantitatively. By decoupling tool discovery from generation, it lets an LLM scale to hundreds or thousands of MCP servers without prompt bloat or decision fatigue; on benchmark tasks, retrieval-first narrowing more than tripled tool-selection accuracy, from 13.62% with all tools loaded to 43.13% [23]. The mechanism is the same one RAG uses for documents: retrieve only what is relevant rather than overwhelming the model with the full corpus.

Hierarchy adds a second axis. MCP-Zero performs two-stage discovery: first filter candidate *servers* by platform requirements, then match specific *tools* within the selected servers, returning only the top-k tool descriptions to cut context overhead [24]. Tool-to-Agent Retrieval generalizes this by embedding both tools and their parent agents in a shared vector space connected by metadata, enabling retrieval at either the tool level or the agent level from one index [25]. The emerging "Tool RAG" framing makes the analogy explicit: tool requests are processed through hierarchical vector routing to retrieve the most relevant tools, which are then injected into context for immediate use [30]. The common thread is **progressive disclosure**: reveal capabilities to the model in layers, starting coarse (domain or server) and refining to fine (specific tool), so the context window never holds more than the decision needs.

Anthropic's deferred tool loading is the same idea expressed as a platform feature: tools marked `defer_loading: true` are not loaded until the Tool Search Tool surfaces them, so an agent addresses thousands of tools while the context holds only the handful the current task requires [15]. Whether implemented as research framework, vendor feature, or serving proxy, the pattern is identical: routing-as-retrieval with hierarchical narrowing.

## Synthesis: Practical Design Guidance for a Routing Layer

The preceding sections converge on a small set of design rules that combine multiple findings, presented here as cross-cutting synthesis rather than as primary findings.

**Start deterministic; escalate to LLM routing only where it pays.** The deterministic-versus-LLM choice is a real trade-off, not a settled question. Rule-based routers are simple, fast, deterministic, and cheap because they spend no inference tokens; LLM routers adapt to more cases but cost inference and add latency [26][27]. The dominant production answer is a hybrid: use deterministic orchestration at the workflow level, LLM-driven execution at the task level, and narrow LLM-driven routing *only* at the decision points where semantic complexity genuinely justifies it [26]. Best practice is to start simple and introduce agentic routing only when model-driven decisions add value over rules — a simple router can replace a full agent and cut latency by avoiding multi-step loops [28][27]. This rule combines the deterministic/LLM trade-off finding with the start-simple guidance and the scaling evidence: the more candidates and the more semantic the disambiguation, the more an LLM stage earns its cost.

**Build a cost-tiered cascade.** Production intent routing often uses an ensemble: microsecond heuristics for initial filtering, lightweight BERT-based classifiers for domain categorization, and more expensive classification only when needed [7]. Reading this together with the matching-primitive cost axis and the retrieval-first scaling evidence yields a clear cascade: cheap-and-exact methods filter the space first, and expensive-and-flexible LLM reasoning runs last and least. The cascade is simultaneously the answer to cost (cheap stages handle most traffic), to the scaling collapse (the LLM never sees fifty tools), and to bias (a retrieval filter precedes the biased selection step).

**Author precedence and design the tie path.** Because overlapping descriptions are disambiguated arbitrarily and the model will not infer a stable precedence, the routing layer must make descriptions distinct and scoped, and must encode explicit priority rules for genuine ties — or escalate ties to a human or a higher-cost router rather than letting them resolve by position [8][9][11]. Combined with the bias finding, this argues for normalizing tool ordering and filtering before selection so that position cannot decide outcomes.

**Treat routing metadata as an untrusted input.** Since tool descriptions enter the model's context to drive the routing decision and can carry injected instructions, the metadata feeding the router must be validated and treated as a security boundary, not as trusted configuration [21][22]. This is the security corollary of the cost-tiered cascade: a retrieval/validation stage in front of selection is also the natural place to sanitize metadata.

## Limitations & Open Questions

Several gaps remain. Terminology is unsettled: the tool-versus-skill distinction used throughout draws on practitioner conventions rather than a standard, so cross-vendor comparisons should be read with that caveat [1]. The quantitative routing results come from specific benchmarks and tasks — the 43%→2% scaling collapse is one BFCL task, and the RAG-MCP and RouteLLM figures are benchmark-specific — so the magnitudes should be treated as directional rather than universal [18][23][14]. Abstention remains weakly evidenced and weakly solved: BFCL shows that deciding *not* to route is the competency current systems are worst at, and the literature offers diagnosis more than remedy [17]. The deterministic-versus-LLM routing question is genuinely open; the hybrid cascade is the consensus engineering answer but not an empirically settled optimum [26][28]. Finally, the security findings establish that the routing layer is an attack surface but the defenses (metadata validation, sandboxing) are still maturing [21][22].

## Sources

1. LLM Skills vs Tools: The Missing Layer in Agent Design — https://www.abstractalgorithms.dev/llm-skills-vs-tools-in-agent-design
2. Equipping agents for the real world with Agent Skills — https://www.anthropic.com/engineering/equipping-agents-for-the-real-world-with-agent-skills · Anthropic
3. Intent Recognition and Auto-Routing in Multi-Agent Systems — https://gist.github.com/mkbctrl/a35764e99fe0c8e8c00b2358f55cd7fa
4. Semantic Tool Selection: Building Smarter AI Agents with Context-Aware Routing — https://vllm-semantic-router.com/blog/semantic-tool-selection/
5. The Roadmap to Mastering Tool Calling in AI Agents — https://machinelearningmastery.com/the-roadmap-to-mastering-tool-calling-in-ai-agents/
6. The Tool Selection Problem: How Agents Choose What to Call When They Have Dozens of Tools — https://tianpan.co/blog/2026-04-09-tool-selection-problem-agent-tool-routing-at-scale
7. When to Reason: Semantic Router for vLLM — https://arxiv.org/html/2510.08731v1 · arXiv
8. Why AI Agents Call the Wrong Tool — and How to Fix It — https://labs.adaline.ai/p/ai-agent-tool-calling-failures
9. Underlying Factors Behind Inconsistency in LLM Responses with Multi-Tool Calling — https://medium.com/@abhaychougule0907/underlying-factors-behind-inconsistency-in-llm-responses-with-multi-tool-calling-628ce7b4de76
10. LLM-as-a-Judge Scoring — https://www.emergentmind.com/topics/llm-as-a-judge-scoring
11. AGENTS.md Patterns: What Actually Changes Agent Behavior — https://blakecrosley.com/blog/agents-md-patterns
12. LangGraph Conditional Edges Example: Router Pattern Implementation Guide — https://langchain-tutorials.github.io/langgraph-conditional-edges-router-pattern-guide/
13. Graph API overview — Docs by LangChain — https://docs.langchain.com/oss/python/langgraph/graph-api · LangChain
14. RouteLLM: Learning to Route LLMs with Preference Data — https://arxiv.org/html/2406.18665v1 · arXiv (UC Berkeley / Anyscale)
15. Introducing advanced tool use on the Claude Developer Platform — https://www.anthropic.com/engineering/advanced-tool-use · Anthropic
16. The Berkeley Function Calling Leaderboard (BFCL): From Tool Use to Agentic Evaluation of LLMs — https://proceedings.mlr.press/v267/patil25a.html · PMLR (ICML 2025)
17. Berkeley Function Calling Leaderboard (BFCL) V4 — https://gorilla.cs.berkeley.edu/leaderboard.html · UC Berkeley Sky Computing Lab
18. Reduce Agent Errors and Token Costs with Semantic Tool Selection — https://dev.to/aws/reduce-agent-errors-and-token-costs-with-semantic-tool-selection-7mf · AWS / DEV Community
19. Why Bad Tool Calling Makes LLMs Slow and Expensive — https://www.codeant.ai/blogs/poor-tool-calling-llm-cost-latency · CodeAnt AI
20. BiasBusters: Uncovering and Mitigating Tool Selection Bias in Large Language Models — https://arxiv.org/pdf/2510.00307 · arXiv (ICLR)
21. MCP Security Notification: Tool Poisoning Attacks — https://invariantlabs.ai/blog/mcp-security-notification-tool-poisoning-attacks · Invariant Labs
22. MCP Tool Poisoning — OWASP Foundation — https://owasp.org/www-community/attacks/MCP_Tool_Poisoning · OWASP
23. RAG-MCP: Mitigating Prompt Bloat in LLM Tool Selection via Retrieval-Augmented Generation — https://arxiv.org/html/2505.03275v1 · arXiv
24. MCP-Zero: Active Tool Discovery for Autonomous LLM Agents — https://arxiv.org/pdf/2506.01056 · arXiv
25. Tool-to-Agent Retrieval: Bridging Tools and Agents for Scalable LLM Multi-Agent Systems — https://arxiv.org/html/2511.01854v1 · arXiv
26. AI Agent Orchestration: LLM vs Code-Driven Patterns — https://genta.dev/resources/ai-agent-orchestration-patterns-llm-vs-code-driven
27. How simple routing vs. agents can reduce latency in LLM-powered applications — https://medium.com/@tannermcrae/rethinking-ai-agents-why-a-simple-router-may-be-all-you-need-c95031c2d397
28. Best Practices for Building an Agent Router — https://arize.com/blog/best-practices-for-building-an-ai-agent-router/ · Arize AI
29. How to Build an AI Agent With Semantic Router and LLM Tools — https://thenewstack.io/how-to-build-an-ai-agent-with-semantic-router-and-llm-tools/ · The New Stack
30. Tool RAG: The Next Breakthrough in Scalable AI Agents — https://next.redhat.com/2025/11/26/tool-rag-the-next-breakthrough-in-scalable-ai-agents/ · Red Hat Emerging Technologies
31. RouteLLM: A framework for serving and evaluating LLM routers — https://github.com/lm-sys/routellm · GitHub (lm-sys)
