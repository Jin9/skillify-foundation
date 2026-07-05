# SKILL.md design

Source: ResearchVault run skill-md-design-20260522-134139 (local deep-research pipeline)
Accessed: 2026-05-22
Category: research-vault / SKILL.md design synthesis
Provenance: harvested 2026-07-05 from ResearchVault 05-final_report.md (run executed 2026-05-22); inline bibliography preserved

## Why This Source Matters

Internal deep-research synthesis on SKILL.md design: structure, frontmatter, and authoring trade-offs with cited sources. Internal-synthesis tier — below official docs; useful as a cross-checked secondary summary.


## Executive Summary

An Agent Skill is, at its core, a folder containing a `SKILL.md` file: YAML frontmatter (minimally a `name` and a `description`) plus a markdown body of instructions that teaches an agent how to perform a specific task repeatably, optionally bundling scripts, references, and assets [1][2]. The design challenge is not the format — the format is deliberately trivial — but the discipline of authoring under a context-window budget. The architectural idea that makes Skills scale is **progressive disclosure**: only the `name` and `description` of every installed skill (≈50–100 tokens each) are pre-loaded into the system prompt, the full body loads only when a task matches, and bundled files and scripts load or run only when the body calls for them [3][4][5]. Because of this staging, the context a skill can carry is effectively unbounded while its idle cost is near-zero [3].

For an expert author the load-bearing decisions are therefore: write a `description` that triggers at exactly the right time (it is the only always-resident discoverability surface, and Claude uses it to pick among 100+ skills) [6][7]; keep the body lean (Anthropic guidance is roughly under 500 lines, the spec ecosystem recommends under ~5000 tokens, and an empirical study found only ~38.5% of typical skill bodies are actionable) [8][9][10]; push everything else into bundled files that load on demand [8][11]; and choose Skills only when the problem is *procedural-knowledge consistency* rather than external connectivity (MCP), context isolation (subagents), or an explicit terminal entry point (slash commands) [12][13]. Skills launched in preview on October 16, 2025 and the spec was published as an open standard on December 18, 2025, adopted by 32 tools within months — making cross-host portability real but imperfect, since fields like `allowed-tools` are enforced by some hosts and silently ignored by others [14][15][16]. Security is a first-order concern: a large-scale study found 26.1% of surveyed skills carried at least one vulnerability, and natural-language instructions evade conventional scanners [17][18].

## Background

Before Skills, encoding a repeatable workflow into an agent meant one of three unsatisfying options: re-explaining the procedure in every prompt, baking it into an always-loaded instruction file (paying its token cost on every inference), or wiring a bespoke tool. Skills target the gap between *capability* and *procedural knowledge*. A tool or MCP server grants an agent the ability to act on an external system; a Skill encodes the team-specific, repeatable *way* to do a task — a commit-message convention, a document-generation procedure, a compliance checklist [12][1]. Anthropic shipped Skills across Claude.ai, Claude Code, and the API in a preview beginning October 16, 2025, gated to Pro, Max, Team, and Enterprise users with Code Execution as a dependency [19][14]. The same `SKILL.md` is designed to work identically across those surfaces, and — once the format was opened as a standard — across competing hosts [4][15]. This report examines what distinguishes a well-designed `SKILL.md` from a poorly-designed one, treating the format spec, the progressive-disclosure model, the skill-vs-alternatives decision, authoring practice, portability, and security as separable design axes.

## Methodology

The research plan decomposed "SKILL.md design" into ten sub-questions spanning six coverage dimensions: definitions and scope, the format/spec, progressive disclosure and design principles, selection versus alternatives, authoring best practices and pitfalls, and portability and trajectory. Evidence was gathered by web search (2025–2026 window) prioritizing primary Anthropic material — the Agent Skills overview and best-practices docs, the engineering blog, the official `anthropics/skills` repository and skill-development skill, the Skill Creator plugin, and the Complete Guide to Building Skills — cross-checked against the cross-vendor `agentskills.io` specification, independent format references, and security research including two arXiv studies. Findings were extracted evidence-first and grouped by sub-question; contradictions (notably the "only name and description are supported" framing) were preserved rather than flattened. The synthesis weights primary docs above secondary commentary and flags figures that rest on a single source.

## Key Findings

- A `SKILL.md` requires exactly two parts — YAML frontmatter with `name` and `description`, and a markdown instruction body; everything else is optional [20][1]. The allowed frontmatter keys are `name`, `description`, `license`, `allowed-tools`, and `metadata` [20].
- Hard field constraints: `description` ≤ 1024 characters, `name` ≤ 64 characters and must match the parent folder name exactly [21][6].
- Progressive disclosure stages context across advertise (~100 tokens) → load full body (recommended < ~5000 tokens) → read resources → run scripts, so idle skills cost almost nothing [11][3].
- The `description` is the single most consequential design surface: it is always loaded and is what the model matches against to select a skill from many [7][6].
- Skill selection criteria are neutral and orthogonal: MCP for connectivity, Skill for procedural consistency, subagent for isolation/parallelism, slash command for an explicit entry point [12][13][22].
- A skill bundles a root `SKILL.md` plus optional `scripts/`, `references/`, and `assets/` subdirectories; splitting overflow into referenced files is the canonical token-discipline move [11][8].
- The format is filesystem-based and portable across ~32 hosts as of early 2026, but permission semantics like `allowed-tools` are enforced inconsistently [16][21].
- Security is unresolved: 26.1% of surveyed skills carried at least one vulnerability; third-party skills should be treated as untrusted supply-chain components [18][17].

## What a Skill is and where it lives in the Claude ecosystem

A Skill is the smallest unit that packages *procedural knowledge*: a directory whose `SKILL.md` carries metadata and instructions, optionally accompanied by scripts, reference material, templates, and other resources [1][2]. The defining behaviour is dynamic loading — "when a task matches a skill's description, the agent reads the full SKILL.md instructions into context," following them and optionally executing bundled code or reading referenced files [2]. This is what lets an agent keep many skills on hand at a small context footprint.

The same artifact is portable across Anthropic's surfaces: the Complete Guide states that skills "work identically across Claude.ai, Claude Code, and API," so an author writes once and runs everywhere, *provided the host supports any dependencies the skill needs* [4]. The mechanics differ by surface, and an expert should know them. In Claude Code and the Agent SDK, skills are filesystem-based and need no API upload; in the SDK specifically the loader is opt-in — by default "the SDK does not load any filesystem settings," and an author must set `settingSources: ['user', 'project']` for skill metadata to be discovered at startup [5]. The API additionally serves Anthropic-managed skills (e.g. `pptx`, `xlsx`) referenced by `skill_id`, with custom skills uploaded via the `/v1/skills` endpoints, and Code Execution is a hard dependency of the feature [5][19].

```
host surface        skill source            notes
─────────────       ──────────────          ───────────────────────
Claude.ai           managed + uploaded      preview: Pro/Max/Team/Ent
Claude Code         filesystem (custom)     no API upload needed
Agent SDK           filesystem (opt-in)     settingSources required
API                 /v1/skills + managed    skill_id refs; code exec
```

The practical upshot for design: a skill that assumes a particular surface (e.g. an interactive UI, or a tool only present in one host) sacrifices the portability that is one of the format's main selling points.

## The format and frontmatter spec: precise constraints, not loose convention

The format is intentionally minimal but the constraints are exact, and getting them wrong is a common failure. A `SKILL.md` has "two required parts: YAML frontmatter with required name and description fields, and a markdown body with the skill's instructions" [20]. The allowed frontmatter properties are `name`, `description`, `license`, `allowed-tools`, and `metadata`; all of these except `name` and `description` are optional [20]. The hard limits matter for authoring tooling and validation: the `description` field has a hard 1024-character limit, the `name` field a 64-character limit, and the `name` must match the parent folder name exactly [21][6].

There is a genuine dispute worth flagging. The widely-repeated framing that "only name and description are supported" in frontmatter is contested: community tooling (the obra/superpowers `writing-skills` issue) reports that more than the minimal pair is actually supported, and the loose framing can mislead authors into omitting useful optional keys like `allowed-tools` [20][23]. The neutral reading is that `name` and `description` are the only *required* keys, while `license`, `allowed-tools`, and `metadata` are supported optional keys an author should reach for when relevant [20].

The `allowed-tools` field is both a design and a security primitive. It "restricts which tools the skill can use and is a security feature for limiting skill permissions," declaring which tool categories — read, exec, write — the skill may invoke [21][24]. Declaring `allowed-tools` narrowly is a defensive default; omitting it widens the skill's effective permission surface to whatever the host grants.

## Progressive disclosure and the token budget

Progressive disclosure is the principle that organizes everything else. It is "the core design principle that makes Agent Skills flexible and scalable," letting Claude "load information only as needed," so "the amount of context that can be bundled into a skill is effectively unbounded" [3]. The mechanism is staged loading. At startup Claude pre-loads "only the name and description from every installed skill's YAML frontmatter into its system prompt," costing "roughly 50–100 tokens per skill" — just enough for the model to decide relevance [4]. Concretely the ecosystem describes a multi-stage pattern: advertise at ~100 tokens per skill, load the full `SKILL.md` body (recommended under ~5000 tokens) when a task matches, read resources only when required, and run scripts only when needed [11].

```
LEVEL 1  frontmatter (name + description)   ~50-100 tok  always loaded
   │                                                     ↓ on match
LEVEL 2  SKILL.md body                       < ~5000 tok loaded
   │                                                     ↓ on demand
LEVEL 3  references/ scripts/ assets/        0 tok idle  read / run
```

The design implication is a strict token budget. A loaded skill is not free: "a 12,000-token skill file costs 12,000 tokens — every time," and "when multiple skills are active simultaneously, their cumulative cost can dominate the context budget" and degrade performance [25]. This is why discipline at Level 2 is the central authoring skill, and why over-stuffing the body — rather than the format itself — is the dominant design failure. Empirical analysis reinforces the point: in one corpus only ~38.5% of skill-body content was actionable core rules, with over 60% being background, examples, or templates injected regardless of relevance, and compressing skills actually improved functional quality by ~2.8% — a measurable less-is-more effect [10][25]. The lesson is that progressive disclosure is not just a loading convenience; respecting it is a precondition for the model's attention to stay on the task.

## Bundling references, scripts, and assets

Progressive disclosure only pays off if the author actually splits content correctly. A skill directory "typically contains a SKILL.md file with optional subdirectories for resources including scripts/, references/, and assets/," holding extra markdown guidance, executable scripts Claude runs via bash, and resources like database schemas, API docs, or templates [11][1]. The canonical token-discipline rule from Anthropic's own skill-development guidance is: "When the SKILL.md file becomes unwieldy, split its content into separate files and reference them. If certain contexts are mutually exclusive or rarely used together, keeping the paths separate will reduce the token usage" [8]. The benefit is concrete — "if your task only needs the sales schema, Claude loads just that one file. The rest remain on the filesystem consuming zero tokens" [5].

A practical threshold accompanies the rule: keep the `SKILL.md` body under roughly 500 lines for optimal performance, and split overflow into separate files via progressive disclosure [8]. Scripts deserve special mention because they convert work that would otherwise consume both tokens and reasoning into deterministic operations: bundled scripts run via bash and provide deterministic results without injecting their source into context [11]. The mis-bundling anti-pattern is therefore inlining reference material that is only occasionally needed directly in the body — it defeats the unbounded-context promise by forcing every invocation to pay for material most invocations don't use.

## Writing a `description` that triggers correctly

If progressive disclosure is the architecture, the `description` is the API to it. It is the only Level-1 surface, and "Claude uses it to choose the right Skill from potentially 100+ available Skills" [7]. The dominant best practice is to encode both *what* and *when*: "Start with what the skill does in one sentence. Then add explicit trigger phrases with 'Use when'" [7]. The two failure modes are symmetric: "an under-specified description means the skill won't trigger when it should; an over-broad description means it triggers when it shouldn't" [7].

Three concrete techniques sharpen triggering. First, write in third person — "the description is injected into the system prompt, and inconsistent point-of-view can cause discovery problems," so prefer "Generates…", "Use when…" over "I can help you…" [6]. Second, use the concrete tokens users actually type, because "the model matches on concrete terms, not abstractions" — "not 'Python testing' but 'pytest'. Not 'Word documents' but '.docx files'" [26]. Third, suppress false positives with explicit negative scope: for skills prone to misfiring on similar requests, add "Do not use when" clauses (e.g. a test-generation skill that should not trigger on "run existing tests") [11][6]. A rigorous author validates triggering with 5–10 fresh held-out queries — a mix of should-trigger and should-not — run through an eval, since queries used during description tuning give an over-optimistic read [11].

## Skill vs tool, MCP, subagent, or slash command

Choosing the wrong primitive is a recurring design error, and the literature converges on a neutral, criteria-based decision rather than a single prescribed answer. The cleanest distinction is between *access* and *procedure*: "MCP doesn't give Claude procedural knowledge about how to use those tools in your specific context. It gives Claude access to a tool. Skills tell Claude the right way to use it for your workflow"; so use Skills "when the problem is about knowledge consistency, not external connectivity" [12]. The practitioner decision tree generalizes this: external system access → add an MCP server; a workflow you re-explain every prompt → create a Skill; autonomous multi-step action → build an Agent; independent parallel tasks → use Subagents [13].

```
need external system access?      → MCP server
re-explaining a workflow?         → Skill
autonomous multi-step action?     → Agent
independent parallel tasks?       → Subagents
explicit terminal entry point?    → slash command
```

Subagents are not substitutes for skills: they "provide independent task execution with specific tool permissions and context isolation," buying parallelization and keeping exploratory work out of the main context — "each layer does something the others genuinely can't" [12]. Slash commands and skills overlap most and confuse authors most: slash commands "are single-file entries with great terminal discovery/autocomplete," whereas "skills are usually directories with supporting files" that "Claude automatically loads… based on description matching" [22]. The design rule that falls out: reach for a skill when you want auto-invoked, richer, file-backed procedural knowledge; reach for a slash command when you want an explicit, single-file, user-triggered entry point. These layers compose — a mature setup uses MCP/plugins for infrastructure and skills for task-specific procedure [22].

## Authoring best practices and recurring anti-patterns

Anthropic's published authoring guidance is unusually concrete and evaluation-first. The recommended workflow is to "create evaluations before writing extensive documentation," so the skill "solves real problems rather than documenting imagined ones": identify gaps, build three test scenarios, establish a baseline, write *minimal* instructions to close the gaps, then iterate by re-running the evals against the baseline [27][6]. This inverts the naive habit of writing a long body first and testing later.

On the prose itself, two meta-rules recur. Specificity: "explicitly name the specific failure modes to avoid, as generic instructions don't work" [27]. And calibrated control: "Control Tuning [is] the meta-rule for how tightly to constrain Claude, matching instruction freedom to task fragility," paired with "Explain-the-Why patterns [that] state the rule and explain why so Claude can generalize to unanticipated cases" — i.e. for fragile tasks give imperative steps, for robust ones give the rationale and let the model adapt [27]. Naming is a small but high-leverage convention: use gerund form (`processing-pdfs`, `analyzing-spreadsheets`, `managing-databases`) because it "clearly describes the activity or capability the Skill provides" [26]. The anti-patterns are the inverse: vague generic instructions, a bloated body that ignores the ~38.5%-actionable reality, a description that summarizes process instead of triggering conditions, and authoring without evals.

## Portability across hosts and competing formats

The format's portability is real and was made explicit by standardization. Because "the specification is filesystem-based, not API-dependent," a skill authored for one agent "can typically run unchanged in Claude Code, Google Antigravity, Cursor, GitHub Copilot, and dozens of other platforms" [16]. This is the strongest argument for investing in `SKILL.md` over a host-proprietary format. But portability is not total, and an expert should weigh two caveats. First, *semantic* portability lags *syntactic* portability: the `allowed-tools` field "Claude Code and OpenClaw enforce… while other agents silently ignore it," so a skill that relies on `allowed-tools` for safety has different guarantees depending on where it runs [21]. Second, `SKILL.md` is not the only convention in the space: AGENTS.md is "a Markdown file convention for project-level AI instructions" created by OpenAI with Google and Cursor and governed by the Linux Foundation's Agentic AI Foundation, a different artifact (project-level instructions) from SKILL.md's directory-based capability packaging [28][16]. The two are complementary rather than competing, but an author should not conflate them.

## Security, safety, and governance

The happy-path authoring guides understate a first-order risk: skills carry executable instructions and often executable code, and they run with real tool access. A large-scale empirical study ("Agent Skills in the Wild") found that "26.1% of skills contain at least one vulnerability, spanning 14 distinct patterns across four categories — prompt injection, data exfiltration, privilege escalation, and supply chain risks," with 5.2% showing high-severity patterns "strongly suggesting malicious intent" [18]. The structural reason is that "traditional malware scanners fail to analyze natural-language instructions," so an agent "processes these instructions as trusted procedural memory and executes them with full host system permissions" [17]. The governing recommendation is to "always treat third-party AI skills as untrusted supply chain components" [17]. For an author, the defensive design moves follow directly: declare `allowed-tools` as narrowly as the task permits, avoid bundling opaque executables a reviewer cannot read, and treat installed third-party skills with the same scrutiny as any dependency.

## Synthesis

Two cross-cutting principles unify the findings. First, **the format is trivial; the budget is everything.** The required spec is two frontmatter fields and a body, but every design lever that matters — body length, what to inline vs. bundle, how many skills to keep active — is a consequence of progressive disclosure and the per-inference token cost of a loaded skill [3][25][10]. Mastery is token discipline, not YAML knowledge. Second, **the `description` and `allowed-tools` fields are where design, discoverability, and security converge.** The `description` is the only always-resident surface and therefore the single point that determines whether a well-written skill is ever used [7][6]; `allowed-tools` is simultaneously a design constraint and the primary defensive control, yet its enforcement is host-dependent [21][24]. An expert author optimizes the `description` for triggering precision and writes `allowed-tools` for a least-privilege default, accepting that portability of the *permission* guarantee is weaker than portability of the *instructions* [21]. The evaluation-first methodology ties these together: because triggering and behaviour are both empirically testable, the highest-leverage authoring practice is to build small evals before writing the body and iterate against a baseline [27][6].

## Limitations & Open Questions

Several findings rest on single sources or secondary summaries and should be treated as provisional. The "< ~5000 tokens body" and "~500 lines" thresholds are heuristics, not hard validator limits, and the underlying performance claims (e.g. the ~2.8% quality gain from compression, the ~38.5%-actionable figure) come from individual studies on specific corpora and may not generalize [10][9]. The vulnerability prevalence figures (26.1%, 5.2%) likewise derive from one large-scale survey and depend on its sampling [18]. The exact page count and date of the Complete Guide were reported via secondary summaries [29]. This report did not independently verify the live `agentskills.io` spec text against Anthropic's docs field-by-field, nor benchmark cross-host `allowed-tools` enforcement empirically; both are flagged as portability/security caveats rather than measured results [21]. Finally, the field is moving fast — preview-era behaviour, the open-standard adoption count, and tooling like Skill Creator are all dated to 2025–2026 and may shift [14][30].

## Sources

1. GitHub — anthropics/skills: Public repository for Agent Skills — https://github.com/anthropics/skills
2. Agent Skills — Claude API Docs (Overview) — https://platform.claude.com/docs/en/agents-and-tools/agent-skills/overview
3. Equipping agents for the real world with Agent Skills — https://www.anthropic.com/engineering/equipping-agents-for-the-real-world-with-agent-skills
4. The Complete Guide to Building Skills for Claude — https://resources.anthropic.com/hubfs/The-Complete-Guide-to-Building-Skill-for-Claude.pdf
5. Agent Skills in the SDK — Claude Code Docs — https://code.claude.com/docs/en/agent-sdk/skills
6. Skill authoring best practices — Claude API Docs — https://platform.claude.com/docs/en/agents-and-tools/agent-skills/best-practices
7. Optimizing skill descriptions — Agent Skills — https://agentskills.io/skill-creation/optimizing-descriptions
8. claude-code/plugins/plugin-dev/skills/skill-development/SKILL.md — https://github.com/anthropics/claude-code/blob/main/plugins/plugin-dev/skills/skill-development/SKILL.md
9. SKILL.md Format and Asset Bundling | github/awesome-copilot | DeepWiki — https://deepwiki.com/github/awesome-copilot/3.2-skill.md-file-format-and-bundling
10. SkillReducer: Optimizing LLM Agent Skills for Token Efficiency — https://arxiv.org/html/2603.29919v1
11. SKILL.md Format and Asset Bundling (progressive disclosure stages) | DeepWiki — https://deepwiki.com/github/awesome-copilot/3.2-skill.md-file-format-and-bundling
12. Skills explained: How Skills compares to prompts, Projects, MCP, and subagents | Claude — https://claude.com/blog/skills-explained
13. Claude Agents, Subagents, Agent Teams, Skills & MCP: A Developer's Field Guide — https://www.antstack.com/blog/claude-agents-subagents-agent-teams-skills-and-mcp-a-developer-s-field-guide/
14. Claude Skills: Launch Timeline & Technical Overview — Verdent Guides — https://www.verdent.ai/guides/claude-skills-announcement-news
15. Anthropic Opens Agent Skills Standard, Continuing Its Pattern of Building Industry Infrastructure — https://www.unite.ai/anthropic-opens-agent-skills-standard-continuing-its-pattern-of-building-industry-infrastructure/
16. Agent Skills Open Standard Explained | SKILL.md Adopted by Claude Code, Codex, Cursor & 32 Tools — https://www.paperclipped.de/en/blog/agent-skills-open-standard-interoperability/
17. Supply Chain Risks of Agent Skills — Pluto Security — https://pluto.security/blog/supply-chain-risks/
18. Agent Skills in the Wild: An Empirical Study of Security Vulnerabilities at Scale — https://arxiv.org/pdf/2601.10338
19. Claude Introduces Agent Skills for Custom AI Workflows — DevOps.com — https://devops.com/claude-introduces-agent-skills-for-custom-ai-workflows/
20. SKILL.md Format Specification | anthropics/skills | DeepWiki — https://deepwiki.com/anthropics/skills/2.2-skill.md-format-specification
21. SKILL.md Spec: Every Field and Frontmatter Key — https://www.agensi.io/learn/skill-md-format-reference
22. Understanding CLAUDE.md vs Skills vs Slash Commands vs Plugins — https://medium.com/@agustin.ignacio.rossi/understanding-claude-md-vs-skills-vs-slash-commands-vs-plugins-1050c68fa63b
23. superpowers/skills/writing-skills/SKILL.md — https://github.com/obra/superpowers/blob/main/skills/writing-skills/SKILL.md
24. Agent Skills: Explore security threats and controls | Red Hat Developer — https://developers.redhat.com/articles/2026/03/10/agent-skills-explore-security-threats-and-controls
25. What Is Context Rot in Claude Code Skills? How Bloated Skill Files Degrade Agent Performance — https://www.mindstudio.ai/blog/context-rot-claude-code-skills-bloated-files
26. Skill authoring best practices (raw) — https://platform.claude.com/docs/en/agents-and-tools/agent-skills/best-practices.md
27. Skill authoring best practices — Claude Docs — https://docs.claude.com/en/docs/agents-and-tools/agent-skills/best-practices
28. How to use AGENTS.md with Codex, Cursor, and Claude Code — https://benjamincrozat.com/agents-md
29. The Complete Guide to Building Skills for Claude (release summary) — https://resources.anthropic.com/hubfs/The-Complete-Guide-to-Building-Skill-for-Claude.pdf
30. Skill Creator — Claude Plugin | Anthropic — https://claude.com/plugins/skill-creator
