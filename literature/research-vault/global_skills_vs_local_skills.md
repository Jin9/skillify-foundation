# Global skills vs local skills

Source: ResearchVault run global-skills-vs-local-skills-20260526-162821 (local deep-research pipeline)
Accessed: 2026-05-26
Category: research-vault / global vs project scope
Provenance: harvested 2026-07-05 from ResearchVault 05-final_report.md (run executed 2026-05-26); inline bibliography preserved

## Why This Source Matters

Internal deep-research synthesis on scoping skills globally (~/.claude, ~/.agents) versus per-project — discovery precedence, duplication risk, and when each placement wins. Internal-synthesis tier.


## Executive Summary

Agent skill systems — Claude Code being the reference implementation — scope a skill the same way operating systems and shells have always scoped configuration: by *where the definition lives on disk*. A **global (user-level / personal) skill** lives in the per-user home directory at `~/.claude/skills/` and is visible in every project you open; a **local (project-level) skill** lives in the repository at `.claude/skills/`, is committed to version control, and ships with the codebase so every teammate inherits it [1][2]. A skill itself is just a folder containing a `SKILL.md` file — name plus description metadata and instructions, optionally bundling scripts, references, and templates [3]. The split is therefore not a feature of the skill format at all; it is a *placement decision* about ownership, reproducibility, and reach.

The operationally load-bearing question is precedence: when a global and a local skill collide on name, which one wins? Claude Code's own documentation states an unusual order — **enterprise > personal > project**, with plugin skills lowest — meaning a personal skill *shadows* a same-named project skill [1]. That ordering is the inverse of nearly every adjacent system: GitHub Copilot ranks repository instructions above user instructions [11], Cursor ranks project rules above user rules [13], and even Claude Code's *own* `settings.json` hierarchy ranks project settings above user settings [24]. The open-standard documentation further muddies this, asserting the opposite — that a project skill overrides a personal one [5]. Practitioners have reported the documented override failing outright, with both same-named skills appearing at once [8]. The net expert takeaway: discovery is well-specified, but precedence is **contested, surface-dependent, and version-sensitive** — never assume it; verify it in your installed build.

Beyond precedence, scope is a governance and cost decision. Project skills are the unit of team reproducibility and should be committed; personal skills are individual preference and should not be checked into shared repos [15][17]. Skills are an *executable* surface — organizations must enable code execution before skills run [20] — and inactive global skills carry a real, compounding token cost, with one estimate of ~23,000 wasted tokens per session [31]. The recurring best-practice heuristic is to pick the smallest scope that still covers the need [21].

## Background

The term "skill" collides with the everyday human-capability sense, so precision matters: throughout this report a *skill* is the Agent Skills artifact — a `SKILL.md`-rooted folder of instructions and resources that an agent loads dynamically to improve at a specialized task [3]. Anthropic introduced Agent Skills on October 16, 2025 as a way to teach Claude repeatable workflows [24], and on December 18, 2025 released the format as an open standard alongside a partner directory and central admin controls [17][18]. The format was deliberately minimal — a folder and a markdown file — which is precisely why *scope* (where the folder sits) became the primary axis of system behavior rather than any syntax in the file itself.

For an expert reader the interesting structure is the parallel to decades of configuration-layering practice: dotfiles in `$HOME` versus repo-committed config, user `git config --global` versus repo `.git/config`, shell rc files versus project `.envrc`. Agent skill systems re-instantiate that same global-vs-local dichotomy, and inherit both its ergonomics and its failure modes — shadowing, drift, and discovery surprises.

## Methodology

This report answers a single thesis — how agent skill systems scope skills as global/user-level versus project/local-level, and what governs precedence, overrides, discovery, and the choice between them — decomposed into eight MECE sub-questions spanning definitions, discovery mechanics, precedence, the cross-system landscape, governance/security trade-offs, decision rules, the temporal trajectory, and real-world failure modes. Evidence is drawn from Anthropic's official Claude Code and platform documentation, the open Agent Skills standard, adjacent vendor docs (GitHub Copilot, Cursor), Anthropic GitHub issues documenting deviations, and practitioner write-ups. Where primary sources disagree — notably on override direction — the disagreement is surfaced as a finding rather than silently resolved. The scope deliberately excludes `SKILL.md` authoring syntax except where it bears on discovery or precedence.

## Key Findings

- **Scope = placement.** Global/personal skills live in `~/.claude/skills/`; local/project skills live in the repo's `.claude/skills/` and are committed to git [1][2].
- **Discovery is progressive.** At session start the agent sees only each skill's name and description; the full `SKILL.md` loads only when the skill is judged relevant [6].
- **Documented precedence is enterprise > personal > project > plugin**, an ordering in which personal *outranks* project [1] — but the open standard and field reports contradict it [5][8].
- **The cross-system norm is the inverse.** Copilot and Cursor both rank project/repo above user [11][13]; Claude Code skills are the outlier.
- **Plugin skills are namespaced** (`plugin:skill`) and therefore cannot collide by name with other scopes [7].
- **Scope is governance.** Commit project skills; keep personal skills out of shared repos [15][17]; enterprises can centrally provision skills that users cannot delete [18].
- **Scope has a cost.** Inactive global skills consume listing tokens every session; a `v2.1.129` listing-budget cap can silently truncate large global sets [26][31].

## Defining the global / local split: storage, ownership, lifecycle

The cleanest way to think about the distinction is by the three properties that differ across the boundary: *storage location*, *ownership*, and *lifecycle*. A personal skill is stored in the per-user home tree at `~/.claude/skills/`, is owned by the individual, and lives for as long as that user's machine config does — it follows the *person*, across every project they open [1][2]. A project skill is stored inside the repository at `.claude/skills/` (note the absence of the leading tilde), is owned by the team via the repo, and lives and dies with the codebase — it follows the *project*, and every teammate who clones the repo inherits it identically [2]. This is the entire substance of "global vs local": not a flag, but a directory.

The artifact placed in either location is structurally identical. A skill is a folder containing a `SKILL.md` with at-minimum a name and description plus instructions, optionally accompanied by `scripts/`, `references/`, `assets/`, and arbitrary supporting files [3]. Anthropic's public skills repository and engineering writeup both frame the format as deliberately lightweight precisely so that the same artifact can be dropped into any scope without modification [4]. The intent split is conventional rather than enforced: personal skills are meant for individual preferences (your editor style, your shortcuts), project skills for team standards (code review, testing, architecture) [3].

Because the artifact is portable and the scope is purely positional, "promoting" a personal skill to a project skill — or demoting a project skill to personal — is just a `mv` plus a `git add`. That fluidity is a strength (cheap to relocate) and the root of the dominant failure mode (the same skill ends up in two places at once, discussed below).

## Discovery and resolution: where the agent looks and how it loads

Discovery has two distinct phases that experts should keep separate: *listing* (which skills the agent knows exist) and *loading* (which `SKILL.md` bodies are pulled into context). At session start the agent sees only the name and description of every available skill — not the body. When a task matches a skill's description, the agent reads the full `SKILL.md` into context; referenced files load only if the instructions reach for them. This progressive disclosure is what lets a user keep many skills on hand for a small standing context cost [6].

The filesystem search itself is hierarchical for project scope. Project skills load from `.claude/skills/` in the starting directory and in every parent directory up to the repository root, so launching the agent in a subdirectory still picks up skills defined at the root; additionally, nested `.claude/skills/` directories below the starting point are discovered on demand when you touch files there [1]. Personal skills are read from the single `~/.claude/skills/` root [1]. Plugin skills are discovered through a separate path: an internal routine reads `~/.claude/plugins/installed_plugins.json` to enumerate *installed* marketplace plugins and loads their bundled skills, so that merely-available-but-uninstalled plugins do not leak into the listing [9].

A practically important property: discovery is *live*. Adding, editing, or removing a skill under `~/.claude/skills/`, the project `.claude/skills/`, or an `--add-dir` location takes effect within the current session without a restart [7]. For an expert this means the skill set is mutable mid-session, which is convenient but also means the active skill list is not a stable snapshot — a fact that matters when reasoning about reproducibility and about the precedence conflicts examined next.

## Precedence and overrides: which scope wins on a name collision

This is the contested heart of the topic. Claude Code's documentation states that when skills share a name across levels, **enterprise overrides personal, and personal overrides project**, with plugin skills (namespaced `plugin-name:skill-name`) at the lowest priority [1][7]. Read literally, this means a personal skill in `~/.claude/skills/` *shadows* a same-named project skill in the repo — the individual's version wins over the team's. Higher-priority discovery locations override lower ones on name conflict, and if a skill and a command share a name, the skill takes precedence [9].

That ordering is genuinely surprising, because it inverts the norm both outside and *inside* Claude Code. The open Agent Skills documentation states the **opposite** — that a project-level skill overrides a personal one with the same name, framing it as teams defining defaults that individuals can override [5]. These two primary sources cannot both be literally true of the same build, so the override direction must be read as product-surface- and version-dependent rather than a settled invariant.

```
Claude Code skills (per docs)   vs   adjacent systems
  enterprise   (highest)              Copilot:  repository > workspace > user
    personal                          Cursor:   team > project > user
      project                         CC settings: managed > project > user
        plugin (lowest)               ── project OUTRANKS user everywhere
  ── personal OUTRANKS project ──         except CC *skills* ──
```

The disagreement is not merely academic. A reported Claude Code issue describes a project-level skill and a same-named global skill **both appearing** rather than one overriding the other [8] — i.e., the documented override silently not firing. The defensible expert posture is therefore: treat the documented enterprise > personal > project ordering as the *intended* contract [1], but do not rely on it for correctness; avoid name collisions across scopes entirely, and lean on plugin namespacing (`plugin:skill`) where you need collision-proof identity [7].

## The cross-system landscape: is the pattern general or a Claude Code quirk?

Comparing adjacent agent and tooling ecosystems shows the global-vs-local *pattern* is universal, but Claude Code skills' *precedence direction* is the outlier. GitHub Copilot ranks repository-level configuration (`.github/copilot-instructions.md`, `.github/instructions/`, `.github/skills/`) highest, then workspace settings, then user settings — explicitly so that team-defined instructions apply even when personal settings differ [11]. Copilot's global instructions live entirely outside the repo (e.g., `~/.config/github-copilot/intellij/global-copilot-instructions.md`, or the VS Code user profile), while repo instructions live under `.github/` [12].

Cursor follows the same shape: rules are merged with precedence **Team > Project > User**, where user rules are global preferences in Cursor settings and project rules live in version-controlled `.cursor/rules`; on conflict, the higher (team/project) source wins [13][14]. Nested `AGENTS.md` files combine with parent directories, with the more specific deeper instruction taking precedence — a *locality-of-specificity* rule [14].

Claude Code's own MCP server configuration is a third model living *inside the same product*: three scopes — local (the default, private to you, stored in `~/.claude.json` under the project path), project (`.mcp.json` in the repo root, shareable via git), and user/global (`~/.claude.json`, all projects) — and crucially, project servers *supplement rather than replace* user servers, so both load together [25]. That merge semantics differs again from the override semantics of skills.

```
Three scoping models inside one ecosystem
  Skills :  name-collision OVERRIDE   personal > project (docs) [contested]
  Settings: array-value MERGE         managed > project > user
  MCP     : additive SUPPLEMENT       project + user both load
```

The lesson for experts: "global vs local" is a stable, cross-vendor pattern, but the *resolution rule* (override vs merge vs supplement, and which direction) is not portable — it varies per artifact type even within one tool. The Agent Skills format itself is now an open standard adopted beyond Anthropic (e.g., OpenClaw using `~/.openclaw/skills/` for global skills), which propagates the *pattern* but not necessarily a single precedence semantics [5][16].

## Governance, security, and trade-offs of scope choice

Scope is a governance decision before it is a convenience. The committed-vs-private boundary is sharp: project skills belong in the repository and are committed to git so teammates receive them on clone, while user-level `~/.claude` skills should *not* be committed into project repos because they are personal preferences [15][17]. This makes project scope the unit of team reproducibility — the highest-value team skills encode decisions otherwise scattered across code reviews and Slack: security requirements, error-handling patterns, naming conventions, test-coverage expectations [22].

At the top of the stack, enterprise governance is real and asymmetric. Team and Enterprise admins can centrally provision skills that are enabled by default for the whole organization; only owners can add or remove org-wide skills, and individual users cannot delete provisioned skills (though they may toggle them off) [18]. This is the concrete mechanism behind the documented "enterprise overrides personal" precedence — enterprise scope is privileged in both *enablement* and *immutability*.

The security surface is non-trivial: skills are *executable*, not merely textual. An organization must enable code execution and file creation before skills function at all, because skills depend on code execution [20]. That makes a globally-installed third-party skill a supply-chain consideration on par with a dependency — the mitigating property is transparency: a skill's text is human-readable and can be reviewed for security before being trusted [19]. The trade-off table is therefore: global scope buys reach and personal ergonomics at the cost of cross-project leakage risk and per-session token overhead; local scope buys reproducibility and reviewability at the cost of duplication across repos and the need to commit and maintain it.

## Choosing global vs local: decision rules for practitioners

The recurring, vendor-endorsed heuristic is **pick the smallest layer that still covers the need**: if a rule applies only to one project, make it project-level; if a behavior applies everywhere, make it a personal skill (or part of global instructions) [21]. Operationally this resolves into a simple test on two axes — *generality* (one project vs every project) and *audience* (just me vs the team):

- A behavior you find yourself repeating across conversations and projects is a **personal/global** skill [23].
- A rule that is meaningful only inside one codebase, or that the whole team must share to get reproducible results, is a **project/local** skill that you commit [22][23].
- A capability the whole organization should have by default belongs at the **enterprise** scope, where admins provision it centrally [18].

```
Decision flow
  Does it apply to EVERY project?
    yes → just you?  → yes → personal (~/.claude/skills/)
                      → no  → enterprise (admin-provisioned)
    no  → one repo, shared with team?
                      → yes → project (.claude/skills/, commit it)
                      → no  → personal (private experiment)
```

The reproducibility bias favors committing project skills: a new teammate who clones the repo should get identical agent behavior without manual setup [15][17]. The ergonomics bias favors personal skills for anything that is *about you* rather than *about the work* — your formatting taste, your shortcut vocabulary [3].

## Where the model breaks down: failure modes and mitigations

The idealized decision rule collides with several documented realities. First, **silent shadowing / non-override**: the documented precedence has been reported to not fire, with a project and a same-named global skill both surfacing instead of one winning [8]. Second, **duplication drift**: people copy a skill into two scopes without realizing it, so the agent loads conflicting instructions depending on which path was referenced — and editing one copy leaves the stale other still running, an unpredictable state [28]. The mitigation for both is disciplinary: never let a skill name exist at two scopes, and run a conflict scanner before trusting resolution.

Third, **discovery gaps that are surface-specific**: one report describes a UI that maintains an internal registry and does not scan `~/.claude/skills/` at startup, loading only 3 of 27 personal skills contrary to documentation [29]; another reports third-party marketplace plugin skills failing to load into context entirely [30]. Global-scope discovery, in other words, is only as reliable as the specific product surface (CLI vs Cowork vs plugin host) implementing it. Fourth, **context cost at the global tier**: because every listed skill consumes name+description tokens each session (≈75–150 tokens per skill with XML overhead), a sprawling `~/.claude/skills/` directory imposes a standing tax — one estimate put inactive-skill waste at ~23,000 tokens per session [31]. Claude Code's `v2.1.129` skill-listing budget caps how many descriptions survive into the system prompt and *truncates* the listing when the budget is exceeded — which silently makes some global skills undiscoverable if you over-install [26]. The mitigation is to keep the global directory focused: install only what you actively use, since fewer loaded skills means lower per-session input cost [32].

## Synthesis

Three of the report's findings combine into a single design tension that the per-finding sections only hint at individually. The *placement-defines-scope* property [1], the *progressive, live discovery* model [6][7], and the *contested override direction* [8] together mean that the effective skill that runs for a given prompt is a function of (a) which directories happen to be in scope at that instant, (b) which product surface is doing discovery, and (c) a precedence rule that documentation itself disagrees about [1][5]. In other words, scope resolution in agent skill systems is currently closer to shell `PATH` resolution — order-dependent, environment-dependent, and prone to shadowing — than to a declarative, deterministic config merge. The cross-system comparison sharpens the point: Copilot and Cursor chose project-over-user precedence for exactly the reproducibility reason that Claude Code's governance guidance also endorses [11][13][22], yet Claude Code's *skills* precedence runs the other way [1]. The pragmatic synthesis for an expert is to engineer *around* resolution rather than depend on it — namespace via plugins [7], avoid cross-scope name collisions, commit what must be reproducible [17], and prune the global tier to control both cost and ambiguity [32].

## Limitations & Open Questions

The single most consequential limitation is the unresolved contradiction on override direction: the Claude Code docs say personal overrides project [1] while the open-standard documentation says project overrides personal [5], and a field report shows neither reliably firing [8]. This report flags the disagreement rather than declaring a winner, because the truth is almost certainly version- and surface-specific and was not pinned to a single build here. Several discovery and cost findings rest on practitioner reports and reverse-engineering (e.g., the `installed_plugins.json` mechanism [9], the ~23k-token estimate [31], the `v2.1.129` budget [26]) rather than primary specification, so their precise numbers and internal names should be treated as medium-confidence and re-verified against the installed version. Finally, the cross-system comparison covers Copilot, Cursor, and MCP as the most load-bearing analogues; other agent ecosystems implementing the open standard may have adopted yet different precedence semantics that are out of scope here.

## Sources

1. Extend Claude with skills — Claude Code Docs — https://code.claude.com/docs/en/skills
2. Where Are Claude Skills Stored? Paths for Mac, Windows, Linux — Agensi — https://www.agensi.io/learn/claude-code-skills-folder-location-setup
3. Agent Skills — Claude API Docs — https://platform.claude.com/docs/en/agents-and-tools/agent-skills/overview
4. anthropics/skills: Public repository for Agent Skills — GitHub — https://github.com/anthropics/skills
5. Agent Skills Overview — agentskills.io — https://agentskills.io/home
6. Skills in Claude Aren't About Prompts — They're About Context Design — DEV Community — https://dev.to/akdevcraft/skills-in-claude-arent-about-prompts-theyre-about-context-design-46hf
7. Understanding Claude Skills — CodeSignal Learn — https://codesignal.com/learn/courses/skills-extending-claudes-capabilities/lessons/understanding-claude-skills
8. Project-level skills with same name as global skills show both instead of overriding (Issue #25209) — GitHub — https://github.com/anthropics/claude-code/issues/25209
9. Skill Discovery Locations — opencode-agent-skills — DeepWiki — https://deepwiki.com/joshuadavidthomas/opencode-agent-skills/3.2-skill-discovery-locations
11. Adding repository custom instructions for GitHub Copilot — GitHub Docs — https://docs.github.com/en/copilot/how-tos/custom-instructions/adding-repository-custom-instructions-for-github-copilot
12. Customize GitHub Copilot in JetBrains with Custom Instructions — Microsoft DevBlogs — https://devblogs.microsoft.com/java/customize-github-copilot-in-jetbrains-with-custom-instructions/
13. Rules — Cursor Docs — https://cursor.com/docs/rules
14. Rules Hierarchy in Cursor — Cursor Community Forum — https://forum.cursor.com/t/rules-hierarchy-in-cursor/108589
15. How to Share Claude Code Skills With Your Team (2026) — Agensi — https://www.agensi.io/learn/how-to-share-claude-code-skills-with-team
16. Equipping agents for the real world with Agent Skills — Anthropic — https://www.anthropic.com/engineering/equipping-agents-for-the-real-world-with-agent-skills
17. Which Claude Code Files Should You Commit? — Is Ray, Not Array — https://israynotarray.com/en/ai/2026/03/31/which-claude-code-files-should-you-commit/
18. Provision and manage Skills for your organization — Claude Help Center — https://support.claude.com/en/articles/13119606-provision-and-manage-skills-for-your-organization
19. Best practices for Claude Code — Claude Code Docs — https://code.claude.com/docs/en/best-practices
20. (governance/security — code execution requirement) Provision and manage Skills for your organization — Claude Help Center — https://support.claude.com/en/articles/13119606-provision-and-manage-skills-for-your-organization
21. Best practices for Claude Code — Claude Code Docs — https://code.claude.com/docs/en/best-practices
22. How to Share Claude Code Skills With Your Team (2026) — Agensi — https://www.agensi.io/learn/how-to-share-claude-code-skills-with-team
23. What Are Claude Code Skills (decision rule: smallest scope; repeat-across-projects → personal) — Best practices for Claude Code — Claude Code Docs — https://code.claude.com/docs/en/best-practices
24. Agent Skills: Anthropic's Next Bid to Define AI Standards — The New Stack — https://thenewstack.io/agent-skills-anthropics-next-bid-to-define-ai-standards/
25. Connect Claude Code to tools via MCP — Claude Code Docs — https://code.claude.com/docs/en/mcp
26. Claude Code's Hidden Skill Budget Setting (May 2026) — claudefa.st — https://claudefa.st/blog/guide/mechanics/skill-listing-budget
28. Stop Adding New Claude Skills — Fix the Broken Ones First — Build to Launch — https://buildtolaunch.substack.com/p/claude-skills-not-working-fix
29. [BUG] Cowork only loads 3/27 personal skills from ~/.claude/skills/ (Issue #50669) — GitHub — https://github.com/anthropics/claude-code/issues/50669
30. Skills from third-party marketplace plugins not loaded into context (Issue #27049) — GitHub — https://github.com/anthropics/claude-code/issues/27049
31. Claude Code Skills Waste 23K Tokens Per Session — AI Weekly — https://aiweekly.co/alerts/claude-code-skills-waste-23k-tokens-per-session
32. Claude Code Skills Waste 23K Tokens Per Session — AI Weekly — https://aiweekly.co/alerts/claude-code-skills-waste-23k-tokens-per-session
