# Skill activation criteria

Source: ResearchVault run skill-activation-criteria-20260526-164049 (local deep-research pipeline)
Accessed: 2026-05-26
Category: research-vault / activation criteria
Provenance: harvested 2026-07-05 from ResearchVault 05-final_report.md (run executed 2026-05-26); inline bibliography preserved

## Why This Source Matters

Internal deep-research synthesis on when skills should (and should not) activate: positive/negative triggers, over-triggering costs, and description phrasing evidence. Internal-synthesis tier.


## Executive Summary

Skill activation — the decision an agent makes about *when* to bring a given skill to bear on a request — is, in the dominant 2025–2026 architecture, a **model judgment over natural-language metadata**, not a deterministic rule engine. In Anthropic's Agent Skills, the model decides whether to consult a skill purely from each skill's pre-loaded `name` plus `description`, with no regex, keyword matcher, embedding index, or intent classifier sitting between the user and the model's forward pass [1][2]. The description field is therefore the load-bearing activation signal, and the practitioner's lever is its wording: it must state both *what* the skill does and *when* (and when not) to fire, written in third person because it is injected verbatim into the system prompt [3][4]. The same semantic-matching logic governs general LLM tool/function calling, where models match query meaning against description fields rather than literal names [5].

This soft, model-driven design has two characteristic failure directions. Claude tends to **under-trigger** — so guidance pushes authors toward deliberately "pushy" descriptions — yet overly generic descriptions cause **false activation**, where a bare keyword like "Kubernetes" fires an unrelated skill; the canonical remedy is an explicit negative trigger ("Do NOT use for…") plus tighter scoping [6][4]. The other systemic threat is **scale**: tool-selection accuracy collapses from above 90% with a handful of candidates to roughly 13.6% with many, degrading sharply once the pool exceeds about 100 candidates [7]. Because of this, retrieval-based pre-filtering (RAG-MCP and successors), progressive disclosure of metadata, and hierarchical/semantic routing have become near-mandatory scaling mechanisms rather than optional optimizations [7][8][9]. Disambiguation among competing skills is solved primarily through **mutually exclusive descriptions**, namespacing, umbrella/sub-tool hierarchies, and sub-agent scoping that reduce contention in the shared system-prompt space [10][11][12]. Evaluation has matured enough to measure activation directly — the Berkeley Function Calling Leaderboard separates relevance detection from *irrelevance* detection (correct abstention), which is the precise quantification of false-activation avoidance — but substantial selection bias persists across frontier models, so activation reliability is best treated as an engineering discipline, not a solved property [13][14][15].

## Background

"Skill activation criteria" names the gating decision at the front of every skill-equipped agent: given a user turn and a library of installed skills, which skills (if any) should be loaded and applied? This is distinct from how a skill executes once chosen — that is the skill's body and is out of scope here. Activation is where reliability is won or lost, because a skill that never fires is inert and a skill that fires wrongly actively degrades output. The question matters acutely now because Agent Skills shipped as a first-class primitive in late 2025 and were subsequently released as an open standard, so the population of installed skills per agent — and thus the contention among them — is growing fast [16]. The mechanisms below were assembled from vendor documentation and engineering posts (the authoritative substrate for Claude Skills and MCP) cross-referenced against the academic tool-retrieval and tool-selection-bias literature that supplies the empirical thresholds.

## Methodology

The research angle is deliberately mechanism-first: for each activation mechanism we ask what signal it consumes, where the decision physically happens (model forward pass, retrieval layer, or host control logic), and what failure modes it introduces. Primary sources are Anthropic's Agent Skills documentation and engineering writing plus the open SKILL.md spec; these are treated as canonical for the description-as-trigger doctrine and the progressive-disclosure loading model. Empirical claims about thresholds, accuracy degradation, and bias are grounded in peer-reviewed-style arXiv work (RAG-MCP, MCP-Zero, BiasBusters, BFCL) rather than asserted. Scope limits: non-LLM skill systems (voice-assistant intent routing, RPA triggers), post-activation execution and sandboxing, and training procedures are excluded except where a training choice (e.g., name-masking) directly shapes the activation signal.

## Key Findings

- Activation in Claude Agent Skills is a model judgment over always-loaded `name` + `description` metadata; there is no regex, keyword matcher, embedding index, or ML intent classifier in the path [1][2].
- The `description` is the primary trigger signal and must encode both capability and when-to-use, in third person, because it is injected into the system prompt [3][4].
- Auto-activation is the default; `disable-model-invocation: true` converts a skill to explicit `/skill-name` invocation, trading autonomy for determinism [17].
- The dominant failure pair is under-triggering (mitigated by "pushy" descriptions) versus false activation from over-generic descriptions (mitigated by explicit negative triggers and scoping) [6][4].
- Tool/skill selection accuracy collapses with candidate count — >90% to ~13.6% — and degrades sharply past ~100 candidates, making retrieval-based pre-filtering effectively mandatory at scale [7].
- Competing-skill disambiguation is won by mutually exclusive descriptions, namespacing, umbrella/sub-tool hierarchies, and sub-agent scoping [10][11][12].
- Activation is now directly measurable: BFCL separates relevance detection from irrelevance detection (correct abstention), and bias studies quantify positional and provider bias [13][18].

## Mechanisms of Activation: Model Judgment, Retrieval, and Control

The first thing an expert must internalize is that "activation" is not one mechanism but three layers that can stack. The cleanest articulation comes from Claude Agent Skills, where activation is **pure LLM reasoning inside the forward pass**: Claude Code formats every installed skill's `name` and `description` into the system prompt and lets the transformer decide which to consult, with explicitly *no* regex, keyword matching, embeddings, classifiers, or pattern matching in the application code [1]. The decision is therefore not observable as a discrete routing step; it is a token-level judgment conditioned on the metadata block. The metadata is "Level 1" content — always loaded — while the full instruction body (`Level 2`) and any resources (`Level 3`) load only after the skill is selected [2]. This three-tier loading is what makes the model-judgment design tractable: only the cheap metadata competes for attention up front.

```
Level 1: name + description     [always in system prompt]  -> ACTIVATION DECISION
   │  (model judges relevance in the forward pass)
   ▼
Level 2: SKILL.md body          [loaded only when selected]
   │
   ▼
Level 3: scripts / references   [loaded on demand during use]
```

The same semantic logic generalizes beyond Anthropic's skills to LLM function/tool calling at large: models decide tool use by matching the *meaning* of the query against the `description` fields of tool schemas, explicitly in contrast to "simple keyword or name-based matching" [5]. Some function-calling training even masks function and parameter names so the model is forced to rely on the semantic description rather than spurious string matching [5]. This is a deliberate engineering stance: literal name matching is treated as a brittleness to be trained out, not a feature to lean on.

Layered *under* or *around* this model judgment are two other mechanisms. The **retrieval layer** (covered in detail below) pre-filters a large candidate set down to a few before the model ever sees them, so an explicit relevance threshold or ranking can live there even though the model's own decision remains a soft judgment [7][8]. The **host control layer** governs whether the model is even allowed to decide: Claude's default is automatic loading of related skills' descriptions, but `disable-model-invocation: true` removes a skill from auto-consideration and reserves it for explicit `/skill-name` invocation, with `user-invocable` separately controlling slash-command exposure [17]. Expert mental model: activation = (host control gate) → (optional retrieval pre-filter) → (model judgment over surviving metadata). Each layer can independently cause a skill to fire or not fire.

## Description and Trigger Design: Authoring the Activation Signal

Because the model judges from the description, **the description *is* the activation criterion**, and authoring it is the single highest-leverage activation control. Anthropic's documentation is explicit that the description should include both what the skill does and the specific triggers/contexts for when to use it, and that it must be written in third person — first- or second-person phrasing causes discovery problems because the text is injected directly into the system prompt [3]. The skill-creator guidance sharpens this into a design metaphor: write descriptions **like routing rules** that tell the model exactly when to activate and when *not* to, using clear, action-oriented language [4]. An optional `when_to_use` field can supplement the description with finer activation rules — specific user phrases or contexts — and the formal SKILL.md spec caps `name` at 64 characters and `description` at 1024 (non-empty, no XML) [19].

The empirical literature confirms how sensitive activation is to this text. In the BiasBusters study, semantic alignment between the query and tool metadata is the *strongest* driver of which tool gets selected, and small perturbations to descriptions significantly shift the model's choice [18]. That sensitivity cuts both ways: it is why good descriptions work, and why sloppy ones misfire. A separate analysis of MCP tool descriptions found 97.1% contained at least one quality "smell" and 56% failed to state their purpose clearly — meaning the majority of deployed activation signals are defective [20]. Yet richer is not automatically better: augmenting every description component raised task success by 5.85 points but increased execution steps by 67.46% and regressed 16.67% of cases, while shorter, targeted descriptions retained the core semantics at equivalent performance [20]. The expert takeaway is a precision target, not a verbosity target: encode the discriminating triggers and the exclusions, and stop.

## Relevance Thresholds: Soft Judgment vs. Explicit Cutoffs

A recurring expert question is whether activation is *thresholded* — a confidence score crossing a cutoff — or an emergent judgment. The honest answer is that it depends on the layer. At the **model layer**, there is no published numeric threshold; activation is a soft semantic judgment, and the evidence that small description perturbations flip selection shows the decision boundary is continuous and sensitive rather than a stable scored cutoff [18]. At the **retrieval layer**, by contrast, the threshold is explicit and engineered. RAG-MCP offloads tool discovery to semantic retrieval against an external index *before* engaging the LLM, passing only the selected tool descriptions onward [7]. Vector-based MCP discovery makes the ranking concrete: a two-stage matching algorithm first filters candidate servers by platform requirements, then ranks tools within them by semantic similarity, reporting ~89% token reduction with maintained accuracy [8]. MCP-Zero's hierarchical routing similarly matches a requested server domain against descriptions to narrow the space, then ranks individual tools by semantic similarity [9].

The practical consequence is that an architect chooses *where* to put the threshold. Leaving it implicit in the model maximizes flexibility but inherits the model's sensitivity and bias; moving it into a retrieval pre-filter gives an explicit, tunable similarity cutoff and slashes prompt tokens, at the cost of a candidate that scores just below threshold never reaching the model at all. The description-augmentation trade-off reinforces that thresholds interact with description length: padding descriptions to clear a similarity bar can backfire by inflating execution steps and regressing cases [20].

## Keyword and Name Matching vs. Semantic Retrieval

Literal keyword and name matching is the cheapest conceivable activation signal and the most brittle. The brittleness is well documented: function-calling models trained on fixed schema patterns fail on semantically equivalent variants — a model trained only on "city, state" weather queries can fail on "city, country" or coordinates [21]. This is precisely the failure that name-masking during training is designed to prevent, by forcing reliance on semantic descriptions over string matches [5]. Anthropic's skills design takes the same position by construction: there is no keyword matcher in the path at all [1].

Where lexical matching breaks down, **semantic/embedding retrieval** scales. MCP-Zero ranks tools by semantic similarity across nearly 3,000 candidates via hierarchical routing [9], and shared-space approaches embed both tools and their parent agents to improve selection [22]. But pure embedding similarity has its own ceiling: static embedding-based matching "lacks fine-grained discriminability for functionally similar tools," which is exactly why learned, history-aware routers (e.g., ToolACE-MCP) are being trained to disambiguate near-duplicates [23]. The progression — literal name match → semantic embedding retrieval → learned router — is a steady move up the discriminability/ cost curve, each step buying resolution between candidates the previous step confuses.

## Disambiguation Among Competing Skills

When multiple skills could plausibly handle a request, activation becomes a contention problem. The crisp framing is **routing competition**: every skill's description competes in the same system-prompt space, so if three skills all use the word "review," the model is "essentially guessing" [10]. The prescribed fix is **mutually exclusive descriptions** — scope each skill to a distinct slice (one reviews migrations, another reviews API contracts, a third enforces style) so overlapping wording does not create routing ambiguity [10]. This is the description-as-routing-rule doctrine applied across a library rather than to a single skill [4].

```
Competing skills in one system-prompt space
                  ┌─ "review migrations"  ─┐
   user request ──┤  "review API contracts" ├─► model must pick
                  └─ "enforce style"       ─┘
   overlap on "review"  ──► ambiguous routing (guessing)
   disjoint scopes      ──► one clear winner
```

Beyond wording, three structural patterns reduce contention. **Namespacing with hierarchy** groups tools logically and introduces umbrella tools that route to specialized sub-tools, plus example-based selection — directly targeting the confusion that overlapping functionality causes [11]. **Shared-space retrieval** addresses the problem that coarse agent-level description matching obscures fine-grained tool capability; embedding tools and their parent agents together improves selection over matching a single coarse description [22]. **Sub-agent scoping** pushes the problem up a level: narrowly-scoped specialist agents with focused prompts and tool access outperform one overloaded generalist, converting "which skill" into "route to which specialist," a cleaner and more learnable decision [12]. Each pattern trades some flexibility or added routing machinery for sharper activation boundaries; the right choice scales with library size — wording discipline for tens of skills, retrieval and hierarchy for hundreds, sub-agents for thousands.

## False Activation and Failure Modes

The defining tension of model-driven activation is that the two failure directions pull in opposite directions. On one side, Claude tends to **under-trigger** — to not use a skill when it would help — so Anthropic advises making descriptions a little "pushy," firing whenever the user mentions relevant terms even without an explicit request [6]. On the other side, that pushiness produces **false activation** when descriptions are too generic: the documented case is a skill whose description matched any mention of "Kubernetes," firing on general Kubernetes tasks it was not built for and producing poor results [6]. The canonical, low-cost fix is an explicit **negative trigger** in the description — e.g., "Do NOT use for general Kubernetes troubleshooting" — which suppresses the adjacent false positives without losing the intended ones [6][4].

Not all false activation traces to description quality. Models exhibit **positional bias**, over-selecting tools that appear earlier in context; the BiasBusters mitigation (filter to a relevant subset, then sample uniformly) cut positional bias from 0.422 to 0.079 and provider/API bias from 0.338 to 0.108 — large effects independent of any single description's wording [18]. This means an author can write a perfect description and still be mis-activated relative to a competitor merely by ordering, which is an argument for the retrieval-pre-filter layer that normalizes the candidate set before the model sees it.

The mitigation methodology that ties this together is **trigger evaluation**: build roughly 20 eval queries mixing should-trigger and should-not-trigger prompts (including edge cases), save them, and test the skill against them to catch over-triggering before deployment [24]. Finally, false *non*-activation also occurs at the control boundary: a documented Claude Code issue had a skill with `disable-model-invocation: true` refused even on explicit slash-command invocation, because the model read the flag as "cannot use this skill at all"; the workaround pairs `user-invocable: true` with the flag [25]. The lesson for experts is that activation reliability spans description wording, candidate ordering, control-flag semantics, and an eval harness — no single one suffices.

## Scaling Activation: The Tool-Count Problem

The most important quantitative result in this space is that activation accuracy does not hold up as the candidate set grows. RAG-MCP reports tool-selection accuracy dropping from above 90% with few tools to roughly **13.62%** with many, and — critically — that both naive baselines and RAG-MCP itself degrade sharply once the candidate pool exceeds about **100 tools** [7]. The mechanism is attention dilution: a large toolset spreads the model's attention thinly across options, raising the probability of incorrect selection or parameter hallucination through choice overload, inter-tool ambiguity, and "lost in the middle" effects [26]. This is why the field treats retrieval-based selection as mandatory rather than optional beyond roughly 100 candidates [7].

Three scaling mechanisms recur. **Progressive disclosure** keeps activation cheap by loading only name + summary metadata first: with about 50 skills, stage-one metadata is ~1–2% of the context window versus ~75% if all documentation were loaded eagerly [27]. **Retrieval pre-filtering** (RAG-MCP, vector discovery) collapses the candidate set to a handful before the model judges, more than tripling selection accuracy versus the all-tools baseline while cutting prompt tokens by over half [7]. **Active/iterative discovery** lets the agent request relevant tools on demand instead of pre-loading everything — MCP-Zero cut token consumption ~98% on APIBank while staying accurate across nearly 3,000 candidates [9]. All three share a single principle: never make the model judge relevance over a candidate set large enough to dilute its attention; shrink the set first, by loading-tier, by retrieval, or by on-demand request.

```
Selection accuracy vs. candidate count (RAG-MCP)
  few tools         ████████████████████  >90%
  ~100 tools        ██████                 sharp degradation onset
  many tools        ██                     ~13.6%
  + RAG pre-filter  ████████               ~43% (3x baseline)
```

## Evaluation and Thresholds: How Activation Is Measured

Activation reliability is now measurable rather than anecdotal. The Berkeley Function Calling Leaderboard (BFCL) operationalizes the *when-to-fire* decision as distinct metrics: **Relevance Detection (RelAcc)** — the fraction of call-requiring queries for which the model emits at least one correct call — and **Irrelevance Detection (IrrelAcc)** — the fraction of `no_call` cases where the model correctly abstains — alongside precision/recall-based F1 [13]. IrrelAcc is the direct quantification of false-activation avoidance: it scores the model precisely on *not* firing when it should not. BFCL further carves out a dedicated "relevance detection (knowing when to refuse a tool call)" category among six, spanning over 2,000 question-function-answer pairs across multiple languages and REST [14].

The benchmarks deliver a sobering verdict on current reliability. BiasBusters found substantial tool-selection bias persisting across seven evaluated LLMs — including GPT-4.1-mini, Claude 3.5 Sonnet, Gemini 2.5 Flash, and Qwen3 — meaning frontier models still systematically mis-weight candidates by provider and position [15]. Combined with the tool-count collapse, the evaluation picture is that activation is a measurable but unsolved capability: thresholds in the retrieval layer buy back accuracy and bias-mitigation buys back fairness, but neither makes the underlying model judgment robust on its own.

## Synthesis: A Layered Activation-Reliability Stack

Pulling the mechanisms together yields a coherent engineering stance that no single source states in full but that the findings jointly imply. Activation reliability is a *stack*, and an expert tunes it layer by layer. At the **authoring layer**, the description is the activation criterion: make it third-person, action-oriented, mutually exclusive from siblings, "pushy" enough to beat under-triggering, and equipped with explicit negative triggers to beat over-triggering [3][4][6][10]. At the **candidate-set layer**, never let the model judge over more candidates than its attention can sustain — use progressive disclosure, retrieval pre-filtering, or active discovery to keep the live set small, since accuracy collapses past ~100 candidates [7][27][9]. At the **disambiguation layer**, resolve residual overlap with namespacing, hierarchy, and sub-agent scoping so competing descriptions stop fighting in one prompt [11][12]. At the **control layer**, decide per skill whether the model may auto-invoke or whether activation is reserved for explicit invocation, and verify the control flags behave as intended [17][25]. And across all layers, **measure**: a should-trigger/should-not-trigger eval set locally, and RelAcc/IrrelAcc-style metrics for the abstention behavior that constitutes false-activation avoidance [24][13]. Each finding above slots into exactly one of these layers; the bias and threshold results [18][20] explain why the layers are needed rather than adding a sixth. The discipline is to recognize which layer a given activation bug lives in — wording, candidate set, contention, control, or measurement — and fix it there rather than over-tuning the description for a problem that is really about scale or ordering.

## Limitations & Open Questions

Several gaps temper the above. First, the evaluation evidence is dominated by *function-calling/tool* benchmarks (BFCL) rather than benchmarks measuring *skill* activation specifically; skill-level activation evals are still emerging, so the BFCL numbers are an informative proxy, not a direct measurement of SKILL.md activation [14]. Second, the in-model relevance threshold is genuinely unobservable — we infer its softness from perturbation sensitivity rather than reading a published cutoff — so claims about "no numeric threshold at the model layer" describe an architecture, not a measured boundary [18]. Third, much of the strongest authoring guidance comes from vendor and practitioner sources whose claims (e.g., "pushy" descriptions, the Kubernetes false-activation anecdote) are credible and corroborated but not independently benchmarked [6]. Fourth, the control-flag failure (`disable-model-invocation`) is a point-in-time bug report that may already be resolved in current releases [25]. Finally, the field is moving fast — Skills became an open standard and learned routers are actively replacing static embedding matching — so the trajectory findings [16][23] are the least stable in this report.

## Sources

1. Claude Agent Skills: A First Principles Deep Dive — Han-chung Lee — https://leehanchung.github.io/blogs/2025/10/26/claude-skills-deep-dive/
2. Equipping agents for the real world with Agent Skills — Anthropic — https://www.anthropic.com/engineering/equipping-agents-for-the-real-world-with-agent-skills
3. Skill authoring best practices — Claude API Docs — https://platform.claude.com/docs/en/agents-and-tools/agent-skills/best-practices
4. anthropics/skills — skill-creator/SKILL.md — https://github.com/anthropics/skills/blob/main/skills/skill-creator/SKILL.md
5. Function Calling: Structured Tool Use for Large Language Models — Michael Brenndoerfer — https://mbrenndoerfer.com/writing/function-calling-llm-structured-tools
6. Skill Authoring Patterns from Anthropic's Best Practices — https://generativeprogrammer.com/p/skill-authoring-patterns-from-anthropics
7. RAG-MCP: Mitigating Prompt Bloat in LLM Tool Selection via Retrieval-Augmented Generation — https://arxiv.org/abs/2505.03275
8. Semantic Tool Discovery for Large Language Models: A Vector-Based Approach to MCP Tool Selection — https://arxiv.org/html/2603.20313
9. MCP-Zero: Active Tool Discovery for Autonomous LLM Agents — https://arxiv.org/pdf/2506.01056
10. Agent Skills, Stripped of Hype — Steve Kinney — https://stevekinney.com/writing/agent-skills
11. Agent-Skills-for-Context-Engineering — tool-design/SKILL.md — https://github.com/muratcankoylan/Agent-Skills-for-Context-Engineering/blob/main/skills/tool-design/SKILL.md
12. Choosing the Right Multi-Agent Architecture — LangChain — https://blog.langchain.com/choosing-the-right-multi-agent-architecture/
13. Berkeley Function Calling Leaderboard (BFCL) V4 — UC Berkeley (Gorilla) — https://gorilla.cs.berkeley.edu/leaderboard.html
14. Berkeley Function-Calling Benchmark (overview) — https://www.emergentmind.com/topics/berkeley-function-calling-benchmark-bfcl
15. BiasBusters benchmark and model coverage — https://arxiv.org/pdf/2510.00307
16. Agent Skills: Anthropic's Next Bid to Define AI Standards — The New Stack — https://thenewstack.io/agent-skills-anthropics-next-bid-to-define-ai-standards/
17. Extend Claude with skills — Claude Code Docs — https://code.claude.com/docs/en/skills
18. BiasBusters: Uncovering and Mitigating Tool Selection Bias in Large Language Models — https://arxiv.org/pdf/2510.00307
19. SKILL.md Spec: Every Field and Frontmatter Key — Agensi — https://www.agensi.io/learn/skill-md-format-reference
20. MCP Tool Descriptions Are Smelly! Towards Improving AI Agent Efficiency with Augmented MCP Tool Descriptions — https://arxiv.org/html/2602.14878v2
21. Tool Calling Fundamentals — Arun Baby — https://www.arunbaby.com/ai-agents/0004-tool-calling-fundamentals/
22. Tool-to-Agent Retrieval: Bridging Tools and Agents for Scalable LLM Multi-Agent Systems — https://arxiv.org/pdf/2511.01854
23. ToolACE-MCP: Generalizing History-Aware Routing from MCP Tools to the Agent Web — https://arxiv.org/pdf/2601.08276
24. superpowers/writing-skills/anthropic-best-practices.md — https://github.com/obra/superpowers/blob/main/skills/writing-skills/anthropic-best-practices.md
25. Skill with disable-model-invocation: true cannot be invoked via slash command (issue #26251) — GitHub anthropics/claude-code — https://github.com/anthropics/claude-code/issues/26251
26. The MCP Tool Trap — Jentic — https://jentic.com/blog/the-mcp-tool-trap
27. Agent Skills: Progressive Disclosure as a System Design Pattern — SwirlAI — https://www.newsletter.swirlai.com/p/agent-skills-progressive-disclosure
