# .claude/skills strategy

Source: ResearchVault run claude-skills-strategy-20260522-135029 (local deep-research pipeline)
Accessed: 2026-05-22
Category: research-vault / .claude/skills strategy
Provenance: harvested 2026-07-05 from ResearchVault 05-final_report.md (run executed 2026-05-22); inline bibliography preserved

## Why This Source Matters

Internal deep-research synthesis on organizing a .claude/skills estate: naming, granularity, and curation strategy for a large personal skill library. Internal-synthesis tier.


## Executive Summary

Agent Skills give Claude reusable, on-demand procedural knowledge packaged as a folder with a `SKILL.md` file, and the way a team organizes the `.claude/skills/` directory determines whether those capabilities are discovered, trusted, and maintained over time [1][2]. The most consequential decision is scope. Claude Code resolves skills from two filesystem locations — `~/.claude/skills/` for personal skills available in every project, and `.claude/skills/` committed inside a repository for skills the whole team shares [3][4]. A third, plugin scope, and a fourth, enterprise-managed scope, sit above and beside these, with a documented precedence of enterprise over personal over project and namespaced plugin skills that cannot collide [3][5]. That precedence is the spine of any organizing strategy, though at least one open issue suggests it is not always enforced as documented [6].

Beyond scope, four organizing levers recur across the primary documentation and field reports. Granularity: one skill per coherent workflow, scoped like a single function — narrow enough to activate precisely, broad enough to avoid loading three skills for one task [7][8]. Discovery: skills are selected purely by the model reading the `description` field, so that field is a trigger, not documentation [9][10]. Progressive disclosure: a three-tier loading model (metadata, body, references) keeps the context budget small while a library scales [11][12]. And composition: skills are one of several primitives — alongside subagents, MCP servers, slash commands, and `CLAUDE.md` — each with a distinct job [13][14].

The cross-cutting tension is reliability versus autonomy. Skills are designed to activate automatically, but practitioners report roughly 50% autonomous activation and a separate, harder execution-fidelity problem [15][16]. Combined with an empirical finding that 26.1% of studied skills carried at least one vulnerability [17], this argues for treating `.claude/skills/` not as a scratchpad but as a governed, version-controlled software asset [18]. This report compares these approaches neutrally; where the documentation leaves layout open, it does not prescribe a single mandatory structure.

## Background

Agent Skills are Anthropic's mechanism for extending an agent with portable, composable expertise: a directory containing a `SKILL.md` file (YAML frontmatter plus markdown instructions) and optional bundled scripts and reference files [1][2]. They were introduced and expanded across Anthropic's surfaces during late 2025 and into 2026, and are now positioned as an open cross-platform standard rather than a Claude Code-only feature [1][5]. As soon as a project accumulates more than a handful of skills, the question shifts from "how do I write one skill" to "how do I organize the directory of skills" — which scope to use, how to name and split them, how to share them across a team, and how to keep them from becoming a security liability. This report addresses that organizing problem for an expert audience already fluent in Claude Code.

## Methodology

The research plan decomposed the topic into nine sub-questions spanning six coverage dimensions: scope and precedence, directory layout and packaging, decomposition and naming, discovery and progressive disclosure, lifecycle/versioning/sharing, and governance/composition. Sources were gathered with live web search (2025–2026 material), prioritizing Anthropic primary documentation — the Claude Code skills docs, the Agent Skills engineering post, the API skills guide, and the `anthropics/skills` repository — and triangulating with independent field reports and one empirical security study. Findings were extracted evidence-first and grouped by sub-question; contradictions (notably documented vs. observed precedence behavior) were preserved as separate claims rather than flattened. The report compares approaches neutrally and does not pin a single team layout as canonical where the docs leave it open.

## Key Findings

- Claude Code reads skills from `~/.claude/skills/` (personal) and `.claude/skills/` (project), with project skills committed to the repo for the team [3][4].
- Documented name-conflict precedence is enterprise > personal > project; plugin skills are namespaced `plugin-name:skill-name` and cannot collide; a same-named skill beats a command [3][5].
- A skill is one directory; `SKILL.md` is required, with `name` (≤64 chars, must match the directory) and `description` (≤1024 chars) the only mandatory frontmatter [19][2].
- Granularity guidance is "one skill per coherent workflow," scoped like a function; both over-narrow and over-broad skills fail [7][8].
- Skill selection is prompt-based — the model chooses from `description` text alone; there is no algorithmic router [9].
- Progressive disclosure loads ~100 tokens of metadata per skill at startup, the full body on activation, and references only on demand [11][12].
- Project skills are shared by committing `.claude/skills/` to Git; plugins/marketplaces add versioning and auto-updates at scale [20][21].
- 26.1% of studied skills carried at least one vulnerability; untrusted skills should never be deployed without a full audit [17][13].
- Skills, subagents, MCP, slash commands, and `CLAUDE.md` are distinct primitives with a clean decision rule [13][14].
- The same `SKILL.md` format runs across Claude.ai, Claude Code, the Agent SDK, and the Messages API, with differing deployment modes [5][22].

## Scope and Precedence: Where a Skill Lives Decides Who Sees It

The first organizing decision is scope, and it is largely binary at the filesystem level: a skill is either personal (`~/.claude/skills/`, visible in every project on your machine) or project (`.claude/skills/`, committed to a repository and shipped to every teammate) [3][4]. Project skills resolve hierarchically — from the launch directory up through every parent to the repository root — so a skill defined at the repo root is available even when Claude starts in a deep subdirectory [23]. Two further scopes wrap these: plugin skills, distributed as bundles and namespaced as `plugin-name:skill-name`, and enterprise-managed skills, uploaded through organization settings on Team and Enterprise plans and enabled by default for all users while remaining individually toggleable [20][2].

When two skills share a name, the documented precedence is enterprise over personal over project, with namespaced plugin skills sitting outside the collision space entirely [5][3]. A useful corollary: a skill takes precedence over a like-named slash command [23]. The practical reading is a layered override model — an organization can ship a managed skill that no project can silently shadow, while individuals keep a personal library that overrides project defaults for their own workflow.

```
enterprise-managed   (highest — cannot be user/project-overridden)
   └─ personal  ~/.claude/skills/
        └─ project  .claude/skills/   (lowest of the three)
   plugin  plugin-name:skill-name     (namespaced — no collision)
```

One caveat tempers this clean picture. An open GitHub issue reports that a project skill sharing a name with a personal skill can surface *both* rather than the project one cleanly overriding — a divergence between documented and observed behavior [6]. This is disputed by the official precedence statement [3], so the prudent strategy is to avoid name collisions across scopes by convention rather than relying on override semantics. The enterprise lever is firmer: managed settings are the highest-priority configuration, delivered via admin console, MDM, or on-disk, and cannot be overridden by any user or project setting [24].

## Layout, Packaging, and Naming: What a Skill Physically Is

A skill is a directory, not a file. The required `SKILL.md` opens with YAML frontmatter and a markdown body; alongside it, optional `scripts/` (executable code), `references/` (docs loaded into context on demand), and `assets/` (output templates, icons, fonts) subdirectories hold the rest [2][1]. Frontmatter is tightly validated: `name` must be ≤64 characters, lowercase letters/numbers/hyphens only, free of XML tags and reserved words, and must exactly match the directory name; `description` must be ≤1024 characters, non-empty, and XML-tag-free [19]. Optional fields include `license`, `compatibility`, an arbitrary `metadata` map, and `allowed-tools`, which narrows the tool surface a skill may invoke [19][30].

```
my-skill/                 (dir name == frontmatter name)
├── SKILL.md              required: frontmatter + body
├── scripts/              optional: executable code
├── references/           optional: docs loaded on demand
└── assets/               optional: output templates/fonts
```

Two naming consequences follow directly. First, because the directory name *is* the skill name and the invocation handle, naming is a layout decision, not a cosmetic one — hyphen-case, verb-or-domain-led, and collision-aware across scopes. Second, because `name` and `description` are the only required frontmatter, every other organizing affordance (`allowed-tools` for governance, `compatibility` for portability, `metadata` for cataloguing) is optional but strategically useful for a maintained library.

## Decomposition: When to Split, When to Merge

Granularity is the lever most likely to be set wrong. The documentation frames it as bidirectionally risky: skills scoped too narrowly force multiple skills to load for a single task, multiplying overhead and inviting conflicting instructions, while skills scoped too broadly become hard for the model to activate precisely [7]. The recommended unit is one skill per coherent workflow — "PR review," "commit messages," "run tests" — and the mental model offered is that deciding a skill's scope is like deciding what a single function should do [8][7].

The canonical split signal is mixed responsibility: a skill that queries a database and formats the results is one coherent unit, but a skill that *also* does database administration is doing too much and should be split [7]. For libraries that grow past a handful of skills, a complementary pattern is to avoid loading everything at session start and instead let a small meta-skill act as a router that selects the skill matching the current phase of work [25]. The split-vs-merge decision therefore has two axes — coherence of responsibility (split when responsibilities diverge) and activation precision (merge when fragments can never fire independently).

## Discovery and Progressive Disclosure: Making Skills Found and Cheap

Skill selection is entirely prompt-based. The model reads the textual `description` of every available skill in its system prompt and decides, by reasoning alone, which to invoke; there is no algorithmic selection or code-level intent detector [9]. That single architectural fact makes the `description` field the highest-leverage string in the whole directory: Anthropic's guidance is to state both *what the skill does* and *when to use it*, written in third person to avoid the point-of-view inconsistencies that degrade discovery [12]. Field practitioners extend this with concrete tactics — list five to seven explicit trigger phrases and add negative boundaries ("Do NOT use for…") so a skill fires on the right requests and stays silent on the wrong ones [10].

The efficiency counterpart to discovery is progressive disclosure, a three-tier loading model. Level 1 is metadata — roughly 100 tokens of `name` and `description` per skill, loaded at startup as a menu the model scans. Level 2 is the full `SKILL.md` body (typically under ~5,000 tokens), loaded only once a skill is activated. Level 3 is bundled reference files, pulled in only when the instructions call for them [11]. The payoff is concrete: an agent with ten skills starts each call carrying about 1,000 tokens of metadata instead of ~10,000 in a monolithic prompt [11].

```
L1 metadata  (~100 tok/skill, startup)  → name + description
      │  activate
L2 SKILL.md body  (<~5k tok)            → full instructions
      │  on demand
L3 references/    (as needed)           → specs, examples
```

This tiering shapes authoring strategy. Anthropic recommends keeping `SKILL.md` under 500 lines and, when approaching that limit, adding a hierarchy layer with explicit pointers to reference files (with a table of contents for any reference over 300 lines) [12]. Content placement is itself an organizing decision: load-bearing instructions belong in the body, while bulky specs and examples belong in `references/`. Authors are also advised to avoid time-sensitive conditionals, preferring a "Current method" section with legacy guidance quarantined in an "old patterns" section [26].

Discovery is, however, the directory's softest point. Independent reports describe autonomous activation succeeding only about half the time, because Claude prioritizes completing the task as it understands it over checking whether a relevant skill exists [15]. A related, harder problem is execution fidelity: even an activated skill is not guaranteed to be followed step by step, an invisible failure mode [16]. The organizing implication is that descriptions should be trigger-engineered, and that critical, state-changing workflows should not rely on autonomous activation alone.

## Versioning, Sharing, and Portability Across Surfaces

The default sharing mechanism is version control itself: placing a `.claude/skills/` directory at the repository root and committing it to Git gives every teammate the same standardized skills automatically [20][4]. This is the lowest-friction team strategy and ties skill versions to the repo's own history. Beyond a single repo, plugins distributed through a marketplace add a heavier-weight layer — a marketplace is a git repo carrying one `.claude-plugin/marketplace.json` file, and it provides centralized discovery, version tracking, and automatic updates installed with `/plugin marketplace add owner/repo` [21][22]. For API-deployed skills, versions can be pinned explicitly via a `version` field on the skill reference, giving deterministic version control independent of the filesystem [27].

Portability is a genuine strength. The same `SKILL.md` folder format runs across Claude.ai, Claude Code, the Agent SDK, and the Developer Platform Messages API, and Anthropic frames the format as an open cross-platform standard [5][1]. Deployment modes differ, though: Claude Code loads skills from the filesystem, while the Messages API accepts up to eight skills as request-scoped parameters that run in a fresh, isolated code-execution container per request [27]. The Agent SDK extends the filesystem model, packaging the same directories as portable, versionable units that agents discover and load on demand [28]. The organizing takeaway is that a well-structured `.claude/skills/` directory is largely portable, but a team should not assume identical runtime semantics — filesystem persistence in Claude Code versus ephemeral, count-capped containers in the API.

## Governance: Treating Skills as Software Assets

Because a skill can carry executable scripts and runs with the agent's permissions, the directory is an attack surface, not just a convenience. Anthropic is explicit: never deploy skills from untrusted sources without a full audit, because a malicious skill can direct Claude to execute arbitrary code, access sensitive files, or exfiltrate data [29]. The risk is empirically non-trivial — a study found 26.1% of skills examined carried at least one vulnerability across fourteen patterns (prompt injection, data exfiltration, privilege escalation, supply-chain), with 5.2% showing high-severity patterns suggesting outright malicious intent [17].

The recommended governance posture is to treat skills like any other software dependency: maintain a single source of truth, enforce metadata and versioning, run security scanning, and establish provenance controls [18]. At execution time, a tiered permission model is advised — deny access to critical files, allow read-write within the agent's own workspace, allowlist specific operations, and default-deny everything else — with the `allowed-tools` frontmatter field narrowing each skill's tool surface [30][12]. For organizations, Claude Code's managed settings are the enforcement lever: as the highest-priority configuration delivered via admin console, MDM, or on-disk, they cannot be overridden by any user or project setting, letting administrators set hard guardrails on what skills may do [24].

## Composition: Skill, Subagent, MCP, Command, or CLAUDE.md?

Choosing the wrong primitive is a common organizing mistake, and the field has converged on a reasonably clean decision rule. Skills are for *how to do something* — procedural workflows; MCP is for *access to external systems*; subagents are for *delegating to a specialist with its own isolated context*; and most production setups use all three together [13]. `CLAUDE.md` occupies a different niche: short, always-true conventions (on the order of 200 tokens) that fire on nearly every interaction and must never be skipped, such as commit-message format or house style — knowledge too small and too universal to merit a skill's progressive-disclosure machinery [14].

Slash commands have largely folded into skills. Custom commands in `.claude/commands/` still work, but skills are now the recommended approach, and a same-named skill takes precedence over a command [14][23]. Within a plugin bundle, these primitives co-exist: a plugin can ship `skills/`, `commands/`, `agents/`, `hooks/`, and an `.mcp.json`, though plugin-shipped agents disallow hooks, `mcpServers`, and `permissionMode` for security reasons [17]. The organizing heuristic that emerges: reach for `CLAUDE.md` for tiny universal rules, a skill for a repeatable procedure that should auto-apply, a subagent when context isolation or parallelism matters, MCP when Claude must touch an external system, and a plugin when you need to distribute several of these together with versioning.

## Synthesis

Read together, the findings describe a directory whose organization is governed by a single budget — the model's attention and context — defended at two gates. The discovery gate (the `description` field) decides whether a skill is *found*; the progressive-disclosure gate (the three-tier loading model) decides how *cheaply* it is carried until needed [9][11]. Scope, naming, and granularity all ultimately serve those two gates: scope decides who is exposed to a skill's metadata, naming and granularity decide whether the model can disambiguate one skill's trigger from another's, and the body/reference split decides how much of the budget activation spends [7][12]. Governance and versioning then wrap the whole structure so that the assets feeding those gates are trustworthy and reproducible [18][20]. The recurring failure mode — roughly coin-flip autonomous activation [15] — is best read not as a reason to abandon the model but as the strongest argument for disciplined description engineering and for reserving autonomous skills for non-state-changing work.

## Limitations & Open Questions

This report draws heavily on Anthropic primary documentation triangulated with field reports; some quantitative claims (the ~50% activation rate, the ~100/~1,000/~10,000-token progressive-disclosure figures) originate in community sources and are illustrative rather than benchmarked [15][11]. The precedence behavior is genuinely contested: documentation states personal-over-project override, while an open issue reports both surfacing — the discrepancy was not resolved by the available sources and may be version-dependent [5][6]. The execution-fidelity problem (skills not followed step-by-step once activated) is described qualitatively but not quantified here [16]. Finally, this report deliberately does not prescribe a single canonical directory layout where the documentation leaves it open; the comparative guidance above is meant to inform a team's own choice rather than mandate one.

## Sources

1. Equipping agents for the real world with Agent Skills — https://www.anthropic.com/engineering/equipping-agents-for-the-real-world-with-agent-skills
2. anthropics/skills: Public repository for Agent Skills — https://github.com/anthropics/skills
3. Extend Claude with skills - Claude Code Docs — https://code.claude.com/docs/en/skills
4. Where Are Claude Skills Stored? Paths for Mac, Windows, Linux — https://www.agensi.io/learn/claude-code-skills-folder-location-setup
5. Agent Skills - Claude API Docs — https://platform.claude.com/docs/en/agents-and-tools/agent-skills/overview
6. Project-level skills with same name as global skills show both instead of overriding · Issue #25209 — https://github.com/anthropics/claude-code/issues/25209
7. Best practices for skill creators - Agent Skills — https://agentskills.io/skill-creation/best-practices
8. 3 Principles for Designing Agent Skills | Block Engineering Blog — https://engineering.block.xyz/blog/3-principles-for-designing-agent-skills
9. Claude Agent Skills: A First Principles Deep Dive — https://leehanchung.github.io/blogs/2025/10/26/claude-skills-deep-dive/
10. Claude Agent Skills Not Triggering? The Problem is in Your Description — https://smartscope.blog/en/blog/agent-skills-description-guide/
11. Agent Skills: Progressive Disclosure as a System Design Pattern — https://www.swirlai.com/p/agent-skills-progressive-disclosure
12. Skill authoring best practices - Claude API Docs — https://platform.claude.com/docs/en/agents-and-tools/agent-skills/best-practices
13. Skills, Slash Commands, MCP, Subagents — I Finally Know When to Use Which — https://varunbhanot.substack.com/p/skills-slash-commands-mcp-subagents
14. Claude Code Customization: CLAUDE.md, Slash Commands, Skills, and Subagents — https://alexop.dev/posts/claude-code-customization-guide-claudemd-skills-subagents/
15. Claude Code Skills Don't Auto-Activate (a workaround) — https://scottspence.com/posts/claude-code-skills-dont-auto-activate
16. Claude Skills Have Two Reliability Problems, Not One — https://medium.com/@marc.bara.iniesta/claude-skills-have-two-reliability-problems-not-one-299401842ca8
17. Plugins reference - Claude Code Docs — https://code.claude.com/docs/en/plugins-reference
18. AI agent skills are becoming the next enterprise supply chain risk - here's how to govern them — https://www.techradar.com/pro/ai-agent-skills-are-becoming-the-next-enterprise-supply-chain-risk-heres-how-to-govern-them
19. SKILL.md Format Specification | anthropics/skills | DeepWiki — https://deepwiki.com/anthropics/skills/2.2-skill.md-format-specification
20. Provision and manage Skills for your organization | Claude Help Center — https://support.claude.com/en/articles/13119606-provision-and-manage-skills-for-your-organization
21. Create and distribute a plugin marketplace - Claude Code Docs — https://code.claude.com/docs/en/plugin-marketplaces
22. Claude Code Plugins vs Skills (2026): What Each Does and When to Use Them — https://www.morphllm.com/claude-code-plugins-vs-skills
23. A Mental Model for Claude Code: Skills, Subagents, and Plugins — https://levelup.gitconnected.com/a-mental-model-for-claude-code-skills-subagents-and-plugins-3dea9924bf05
24. Set up Claude Code for your organization - Claude Code Docs — https://code.claude.com/docs/en/admin-setup
25. Agent Skills — AddyOsmani.com — https://addyosmani.com/blog/agent-skills/
26. Skill Authoring Patterns from Anthropic's Best Practices — https://generativeprogrammer.com/p/skill-authoring-patterns-from-anthropics
27. Using Agent Skills with the API - Claude API Docs — https://platform.claude.com/docs/en/build-with-claude/skills-guide
28. Claude Agent SDK Skills: Build Reusable Agent Capabilities — https://www.augmentcode.com/guides/claude-agent-sdk-skills-reusable-agent-capabilities
29. Skills for enterprise - Claude API Docs — https://platform.claude.com/docs/en/agents-and-tools/agent-skills/enterprise
30. Agent Skills: Explore security threats and controls | Red Hat Developer — https://developers.redhat.com/articles/2026/03/10/agent-skills-explore-security-threats-and-controls
