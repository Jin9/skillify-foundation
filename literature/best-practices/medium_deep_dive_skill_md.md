# Deep Dive SKILL.md (Part 1/2). Lately I have been able to do most of… | by A B Vijay Kumar

Source: https://abvijaykumar.medium.com/deep-dive-skill-md-part-1-2-09fc9a536996
Accessed: 2026-07-05
Category: best-practices / skill folder anatomy
Provenance: re-captured 2026-07-05 and normalized from HTML (previous capture was a 139-word Cloudflare challenge stub)

## Why This Source Matters

Practical walk-through of skill folder anatomy: SKILL.md plus scripts, references, and assets, and how the pieces load. Tier 5 community article; structure examples only.


## Agentic AI Series

# Deep Dive SKILL.md (Part 1/2)

## Lately I have been able to do most of my work with not writing a single line of code, but defining the SKILL.md

A B Vijay Kumar

17 min read

·

Mar 17, 2026

[https://medium.com/m/signin?actionUrl=https%3A%2F%2Fmedium.com%2F_%2Fvote%2Fp%2F09fc9a536996&operation=register&redirect=https%3A%2F%2Fabvijaykumar.medium.com%2Fdeep-dive-skill-md-part-1-2-09fc9a536996&user=A+B+Vijay+Kumar&userId=5be0e5f0ce88&source=---header_actions--09fc9a536996---------------------clap_footer------------------](https://medium.com/m/signin?actionUrl=https%3A%2F%2Fmedium.com%2F_%2Fvote%2Fp%2F09fc9a536996&operation=register&redirect=https%3A%2F%2Fabvijaykumar.medium.com%2Fdeep-dive-skill-md-part-1-2-09fc9a536996&user=A+B+Vijay+Kumar&userId=5be0e5f0ce88&source=---header_actions--09fc9a536996---------------------clap_footer------------------)

--

[https://medium.com/m/signin?actionUrl=https%3A%2F%2Fmedium.com%2F_%2Frepost%2Fp%2F09fc9a536996&operation=register&redirect=https%3A%2F%2Fabvijaykumar.medium.com%2Fdeep-dive-skill-md-part-1-2-09fc9a536996&user=A+B+Vijay+Kumar&userId=5be0e5f0ce88&source=---header_actions--09fc9a536996---------------------repost_header------------------](https://medium.com/m/signin?actionUrl=https%3A%2F%2Fmedium.com%2F_%2Frepost%2Fp%2F09fc9a536996&operation=register&redirect=https%3A%2F%2Fabvijaykumar.medium.com%2Fdeep-dive-skill-md-part-1-2-09fc9a536996&user=A+B+Vijay+Kumar&userId=5be0e5f0ce88&source=---header_actions--09fc9a536996---------------------repost_header------------------)

[https://medium.com/m/signin?actionUrl=https%3A%2F%2Fmedium.com%2F_%2Fbookmark%2Fp%2F09fc9a536996&operation=register&redirect=https%3A%2F%2Fabvijaykumar.medium.com%2Fdeep-dive-skill-md-part-1-2-09fc9a536996&source=---header_actions--09fc9a536996---------------------bookmark_footer------------------](https://medium.com/m/signin?actionUrl=https%3A%2F%2Fmedium.com%2F_%2Fbookmark%2Fp%2F09fc9a536996&operation=register&redirect=https%3A%2F%2Fabvijaykumar.medium.com%2Fdeep-dive-skill-md-part-1-2-09fc9a536996&source=---header_actions--09fc9a536996---------------------bookmark_footer------------------)

[https://medium.com/m/signin?actionUrl=https%3A%2F%2Fmedium.com%2Fplans%3Fdimension%3Dpost_audio_button%26postId%3D09fc9a536996&operation=register&redirect=https%3A%2F%2Fabvijaykumar.medium.com%2Fdeep-dive-skill-md-part-1-2-09fc9a536996&source=---header_actions--09fc9a536996---------------------post_audio_button------------------](https://medium.com/m/signin?actionUrl=https%3A%2F%2Fmedium.com%2Fplans%3Fdimension%3Dpost_audio_button%26postId%3D09fc9a536996&operation=register&redirect=https%3A%2F%2Fabvijaykumar.medium.com%2Fdeep-dive-skill-md-part-1-2-09fc9a536996&source=---header_actions--09fc9a536996---------------------post_audio_button------------------)

Press enter or click to view image in full size

AI-generated (using nano banana 2 Gemini)

This blog will be a 2-part series; in this first part, I am going to walk through the concepts and the architecture of SKILL.md and some internal workings of how agents use SKILL.md to perform complicated tasks. In the second part of this series, I will be walking through the process of building a SKILL.md

Lately, I have realized that there is really no need to program agents to do most of the tasks, I was able to get all my work done with deep agents with SKILL.md. I did cover Deep Agents in the following blog.

## Building Deep Agents

### In this blog, I will share my experience building Autonomous, Memory-Rich, Multi-Step AI Systems

abvijaykumar.medium.com

I am slowly starting to believe in `SKILL.md` will be the most important skill, I need to learn :-D, and I wanted to share my understanding and experience building these `SKILL.md`

The following picture brings everything together.

Press enter or click to view image in full size

Let me walk through this diagram, because it captures how SKILL.md fits into the broader agent architecture.

**Identity Layer: **this is where the agent gets its personality and expertise. SKILL.md files live here alongside constraints and safety rules. Together, they define *what the agent knows* and *what it’s allowed to do*. When you install a skill, you’re extending this layer.

**Orchestration Layer: **is** **the engine. The Dynamic System Prompt is where the skill catalog (all those Level 1 metadata entries) gets embedded. The Large Language Model reads the user’s request, consults the system prompt (which includes the skill catalog), and makes the routing decision we’ll discuss later. The Agent Runtime manages the thought-action loop,the model thinks, decides on an action, the runtime executes it, and the result feeds back to the model for the next step. Task State tracks where the agent is in a multi-step workflow.

**MCP Layer: **connects the agent to the outside world through the Model Context Protocol. MCP Tool Servers expose specific capabilities (file operations, API calls, database queries), and MCP Resources provide read-access to external data. This is how the agent actually *does thing, *skills tell it what to do, and MCP gives it the tools to do it.

**External Systems: **are the real-world endpoint: databases, APIs, and the filesystem where your skill directories actually live.

The key insight from this diagram is that skills operate at the highest level of the stack. They don’t touch the MCP layer directly, they don’t manage tool connections, and they don’t handle state. They influence the model’s behavior through the system prompt, and the model orchestrates everything else.

Ofcourse, I did not cover the semantic layer (RAG, vectors, graphs etc), to keep it simple.

### Agent Skills

A skill is a directory on a filesystem that contains at minimum a `SKILL.md` file, and optionally some subdirectories for scripts, reference docs, and assets. Those files together form a self-contained package of domain expertise that an AI agent can discover, decide is relevant, load up, and use.

Here’s a clean folder-tree representation of a Skill:

```
skill-name/                          ← Skill Directory
│
├── SKILL.md                         ← Required
│   ├── --- YAML Frontmatter ---     ← name + description (the skill's ID card)
│   │       name: skill-name
│   │       description: What it does and when to trigger
│   │   ---
│   └── ## Markdown Body             ← Instructions (the actual playbook)
│         Step-by-step workflows, examples,
│         code patterns, gotchas, etc.
│
├── scripts/                         ← Optional — executable code
│   ├── validate.py                  ← Python scripts
│   ├── setup.sh                     ← Bash scripts
│   └── build.js                     ← JavaScript scripts
│
├── references/                      ← Optional — extra documentation
│   ├── REFERENCE.md                 ← Detailed technical reference
│   ├── domain-guide.md              ← Domain-specific docs
│   └── api-patterns.md              ← API guides
│
└── assets/                          ← Optional — static resources
    ├── template.docx                ← Templates
    ├── config.json                  ← Configuration files
    └── lookup-table.csv             ← Data files
```

So what makes skills architecturally interesting compared to other ways of extending AI agents? A few things stand out.

- **filesystem-based: **They’re just directories and files. You can read them in any text editor, version them with Git, copy them around, and inspect them without any special tools. No databases, no registries, no compiled artifacts. This is intentional—it means the barrier to creating and sharing skills is essentially zero.
- **prompt-based:** When an agent decides to use a skill, what actually happens is pretty simple: the skill’s instructions get loaded into the agent’s context window. It’s basically expanding the conversation with domain-specific guidance. There’s no separate runtime, no plugin architecture, and no function registry behind the scenes. The skill’s power comes entirely from how good its instructions are and how well the agent follows them.
- **composable**. You can have multiple skills active at the same time. They can’t explicitly talk to each other, but the agent can pull in several skills during a single interaction if the task calls for it. This composability just… works, naturally, because of how progressive disclosure is designed.
- **portable**. Thanks to the open standard at agentskills.io, a skill you write for Claude Code also works in VS Code with Copilot, in Cursor, in OpenAI’s Codex, or in any agent that supports the spec (with minimal tweaks).

### Anatomy

The `SKILL.md` file is where all the magic lives. It’s a Markdown file that starts with YAML frontmatter (the stuff between the `---` markers) followed by the actual instructions in regular Markdown. Think of the frontmatter as the skill’s ID card—it tells the system “what this skill is” and “when to use it." The Markdown body is the playbook: it tells the agent “what to actually do."

### The Core Design Principle

Progressive disclosure is the key feature that makes the whole skills system actually work at scale and load only the skills that required and when it is really required. Install ten skills, and your agent’s context window would be stuffed with instructions, leaving barely any room for the actual conversation. With progressive disclosure, you can have dozens or even hundreds of skills installed, and the overhead is negligible.

The concept works on three levels:

- **Level 1 – Just the metadata (~100 tokens per skill).** When the agent starts up, it only loads the `name` and `description` from every installed skill’s frontmatter.
- **Level 2 – The full instructions (under 5,000 tokens, ideally).** When the agent decides a particular skill is relevant to what you’re asking, then it loads the complete SKILL.md body.
- **Level 3 – The deep reference material (only when needed).** Sometimes a task needs even more detail—a specific form-filling guide, a migration risk catalog, or an exhaustive API reference. The agent can selectively pull in files from the skill’s `scripts/`, `references/`, or `assets/` directories. This is the appendix—available if you need it, invisible if you don’t. Because of this third level, the total amount of knowledge you can bundle into a skill is effectively unlimited.

What makes this whole approach so clean is that it degrades gracefully. A simple skill with just a SKILL.md and no extra files? Works perfectly – just two levels instead of three. A skill with a massive reference library? Costs nothing extra until those references are actually needed. And because everything operates through standard filesystem reads, there’s no caching layer to maintain, no separate service to run, and no startup penalty beyond reading a few hundred bytes of YAML per installed skill.

### The Skill Lifecycle

To really understand how skills go from being dormant files on your filesystem to active instructions shaping how the agent works, you need to follow the full lifecycle. It has four phases: installation, discovery, routing, and execution.

Press enter or click to view image in full size

Four phases, left to right. Installation (green) is just putting the folder in the right place. Discovery (blue) is the agent scanning directories and building its catalog of what’s available. Routing (terracotta) is the crucial decision point; the LLM reads the user’s message against all skill descriptions and decides if any match. If yes, the full SKILL.md loads; if no, the agent proceeds without skill help. Execution (burgundy) is the agent following the loaded instructions and pulling in extra resources as needed.

- **Phase 1: Installation** – Installing a skill is literally just putting the folder in the right place. In Claude Code, that means `~/.claude/skills/` for personal cross-project skills, `.claude/skills/` within a repo for project-specific ones, or the plugin marketplace. In Claude.ai, you upload a zip through Settings > Features. Through the API, you reference a `skill_id` in your container config. The key thing to appreciate is that there’s no registration step, no build process, no database entry. You drop a folder and you’re done.
- **Phase 2: Discovery** – When the agent starts a session (or when skill files change mid-session in Claude Code, thanks to live file watching), it scans all the skill directories, reads just the YAML frontmatter from each SKILL.md, and builds a catalog of what’s available. This catalog gets embedded into the agent’s system prompt. The whole scan is lightweight because it only touches the frontmatter – never the full file bodies.
- **Phase 3: Routing** – This distinctive architectural feature eliminates algorithmic routing. Instead, the system formats skill names and descriptions into a text block in the prompt and lets the LLM determine which skill applies to the request. We will be diving deep on this, later in the blog.
- **Phase 4: Execution** – Once a skill is selected, the agent reads the full SKILL.md into its context. The instructions become part of the active conversation, guiding what the agent does next. If those instructions point to additional files; a forms guide, a language mapping reference, whatever – the agent can load those too. There’s no special “skill execution mode”; it’s just the agent operating with a richer context that now includes the skill’s expertise.

### Frontmatter Specification

The YAML frontmatter is the control panel for your skill. Let’s walk through each field and what it actually does.

Press enter or click to view image in full size

**Body Content**

The Markdown body is where your expertise lives. It’s what the agent loads into context when the skill fires, and how good it is directly determines how well the agent performs. The spec doesn’t impose any formatting rules — write whatever helps — but months of real-world usage have surfaced some patterns that work way better than others.

Press enter or click to view image in full size

### **Bundled Resources—Scripts, References, and Assets**

Beyond the SKILL.md file, skills can include three categories of supporting files, each serving a distinct purpose and loaded under different conditions.

Press enter or click to view image in full size

### Skill Triggering

The skill-triggering mechanism is a distinctive aspect of the Agent Skills system, crucial for building reliable skills.

At the code level, there’s no algorithmic routing. Claude Code doesn’t use embeddings, classifiers, regex, keyword matching, or ML-based intent detection to decide which skill to invoke. Instead, it constructs a dynamic tool description by aggregating skill names and descriptions and presents it as part of a meta-tool called “Skill” in the agent’s prompt. The language model’s forward pass through the transformer decides which skill to invoke (if any).

The language model brings its understanding of natural language semantics, context, and nuance to the routing decision. It can understand that “make me a deck about Q3 results” should trigger the `pptx` skill, even though “presentation” isn’t mentioned. It can disambiguate between “extract text from this document” (which might trigger `pdf` or `docx`) and “write a document summarizing this text” (which should trigger `docx` for creation, not extraction).

However, LLM-based routing has characteristics that skill authors must account for. The model tends to “undertrigger,” not invoking a skill even when it’s helpful, especially for requests the model can handle directly. Simple queries like “read this PDF” may not trigger the `pdf` skill because the model believes it can handle the task directly. Complex, multi-step queries and specialized domain tasks are more reliable triggers.

Press enter or click to view image in full size

This diagram traces a specific routing decision through the LLM-based system. When the user asks to fix formatting in a `.docx` file, the language model evaluates this request against the descriptions of all installed skills. The `docx `skill’s description mentions “manipulate Word documents” and “formatting,” producing a strong match. The `pdf`, `xlsx`, and `pptx` skills clearly do not match. With the match identified, the system loads the full `docx/SKILL.md` into the context window, and the agent proceeds to execute the task with the skill’s specialized instructions guiding its behavior. The crucial point is that this entire routing decision—including the semantic understanding that “quarterly report.docx” implies a Word document and “fix the formatting” implies document manipulation—happens within the language model itself, not in any external routing logic.

### Skill Categories

The Agent Skills ecosystem spans a broad range of categories, from Anthropic’s production-ready built-in skills to enterprise-specific custom skills to community-contributed open-source skills.

Press enter or click to view image in full size

Following are some of the best repositories for skill.md

[https://skillsmp.com/](https://skillsmp.com/)

[## GitHub - anthropics/skills: Public repository for Agent Skills

### Public repository for Agent Skills. Contribute to anthropics/skills development by creating an account on GitHub.

github.com](https://github.com/anthropics/skills?source=post_page-----09fc9a536996---------------------------------------)

## Security and Trust Considerations

Since skills inject instructions directly into the agent’s context window, they’re essentially shaping how the agent thinks and acts. When you install a skill from the community marketplace, you’re handing that skill author the ability to influence your agent’s behavior. A well-intentioned skill might tell the agent to always validate inputs before writing files. A malicious one could instruct the agent to exfiltrate data or silently modify code in harmful ways. This is basically the prompt injection problem, but packaged in a neat folder structure.

There are a few guardrails built into the system to mitigate this.

The `allowed-tools` field in the frontmatter (still experimental) lets skill authors declare exactly which tools the skill needs: for example, `Bash(git:*) Bash(jq:*)` Read. This scopes the skill's permissions so it can't reach for tools it shouldn't need. If a skill about formatting Word documents suddenly wants shell access, it's a red flag.

The `license` field serves a dual purpose. Beyond legal clarity, it signals provenance. A skill marked `Proprietary. LICENSE.txt has complete terms` that ships from Anthropic carries very different trust signals than an unlicensed community skill.

For enterprise deployments, treat skills like codebase dependencies. Review, pin versions, and audit changes. Here are some practical guidelines: read a skill’s `SKILL.md` before installing it (it’s a text file). Prefer skills from known repositories like `github.com/anthropics/skills`. Be cautious with skills that include executable scripts in the scripts/ directory. Use project-level skills (`.claude/skills/`) over global ones (`~/.claude/skills/`) to limit the blast radius.

>

The Agent Skills specification is still evolving its security model — things like skill signing, trust levels, and sandboxed execution are active areas of development. For now, the best defense is the simplest one: read what you install.

## Limitations

I truly believe in the power of SKILL.md, but it’s crucial to acknowledge its limitations. Skills aren’t a magical solution, and trying to fit everything into a SKILL.md can actually complicate matters.

- Skill Limitations: Skills lack memory, cannot directly call APIs, and have context window limits.
- Complex Logic Handling: Skills are not ideal for highly dynamic logic with complex branching or real-time feedback.
- Model Dependence: Skill effectiveness varies across different models, impacting cross-platform compatibility.

The ideal use for skills is domain expertise that’s mostly procedural — step-by-step workflows, formatting rules, code patterns, validation checklists, and domain-specific pitfalls.

## Versioning and Skill Updates

One thing that catches people off guard: skills don’t have a built-in versioning mechanism. There’s no `version` field in the required frontmatter spec. However, the `metadata` field — which is an optional catch-all key-value map — is the recommended place to track this.

```
---
name: "name"
description: "descrition"
metadata:
  version: "version number"
  author: "author"
  last-updated: "date"
---
```

Since skills are just folders and files, the most natural versioning strategy is the one you already use for code: Git. Version your skills alongside your project code in `.claude/skills/`, and you get history, diffs, branching, and rollback for free. For global skills in `~/.claude/skills/`, consider maintaining them in a separate Git repository that you clone or symlink.

What happens when a skill is updated mid-session? In Claude Code, file watching detects changes to skill files and updates the agent’s catalog in real time. This means you can edit a SKILL.md, save it, and the agent’s behavior will reflect the changes on the very next interaction — no restart needed. That’s incredibly useful during skill development, but it also means you should be careful about editing production skills while an agent is actively working on a task.

For team environments, I’d recommend treating skills like you’d treat shared configuration: use a pull request process, have someone review changes to SKILL.md files, and tag releases so you can pin to a known-good version if something breaks.

## How Skills Compare to Other Approaches

If you’ve been building with AI agents for a while, you’ve probably used other ways to extend their capabilities. Here’s how Agent Skills stack up against the common alternatives.

Before SKILL.md, every AI coding tool had its own proprietary way of accepting custom instructions — and none of them talked to each other.

- Cursor had `.cursor/rules/` with `.mdc` files using YAML frontmatter and glob patterns.
- GitHub Copilot had `.github/copilot-instructions.md`.
- Claude had `CLAUDE.md`.
- Windsurf had `.windsurfrules`.
- Google's Jules had `JULES.md`.

If you used more than one tool (and most teams do), you were stuck maintaining the same instructions in three or four different formats. Developers literally resorted to symlinking files or writing CLI tools like `rule-porter` just to convert rules between formats. It was a digital Tower of Babel.

OpenAI’s Custom GPTs offered a different approach: configure a role-based persona via a web UI, upload knowledge files, and optionally connect external APIs through OpenAPI schemas (Actions). It was accessible and no-code, but customizations were locked within OpenAI’s ecosystem, preventing script execution, progressive disclosure, and Git version control. GPT Actions were deprecated in early 2024 due to complexity and unpredictable behavior. OpenAI’s replacement, function calling, the Assistants API, and Codex Skills, converged on the SKILL.md pattern pioneered by Anthropic.

Traditional agent frameworks like LangChain, CrewAI, AutoGPT, and the OpenAI Agents SDK define agent capabilities through code (Python classes, tool registrations, function schemas, orchestration logic). These are powerful and flexible but require programming skills, lack portability, and don’t benefit from the simplicity of SKILL.md. Fine-tuning permanently bakes behavior into model weights, making it expensive, slow to iterate, impossible to version, and non-portable across models.

The `AGENTS.md `standard (a vendor-neutral project-level instruction file) provides always-on context about your project, including architecture decisions, coding standards, and build instructions. Several tools read it, but AGENTS.md differs from SKILL.md in scope and behavior. AGENTS.md is always loaded, applies to everything, and doesn’t support progressive disclosure, bundled scripts, or on-demand activation. It tells the agent “here’s how we do things around here” (always-on project context), while SKILL.md tells the agent “here’s how to do this specific task” (on-demand procedural expertise). They’re complementary, not competing.

`SKILL.md` is unique in this landscape due to its four properties: progressive disclosure (load only what’s needed), bundled executable code (scripts that run deterministically alongside LLM-generated output), filesystem-native portability (just folders and files, versionable with Git), and cross-platform standardization (the same skill works in various agents). No other approach delivers all four.

## Standardization & Adoption

The standardization of SKILL.md is one of those rare cases where a single company’s feature became an industry-wide standard in a matter of weeks rather than years. Anthropic published the Agent Skills specification at agentskills.io, releasing it as an open standard under Apache 2.0 (code) and CC-BY-4.0 (docs), with a reference SDK and validation tooling.

The speed of adoption was almost unprecedented for an AI standard. Within days of publication, Microsoft had integrated SKILL.md support into VS Code and GitHub Copilot. OpenAI adopted a structurally identical architecture in Codex CLI and ChatGPT. Cursor, Amp, Goose, OpenCode, Letta, and others followed. Vercel launched [skills.sh](https://skills.sh/). By March 2026, the community [Antigravity Awesome Skills library](https://github.com/sickn33/antigravity-awesome-skills) had cataloged over 1,234 skills compatible with 16+ different agents across the ecosystem. Microsoft even built a full Agent Framework SDK (Python, C#, and TypeScript) with a `SkillsProvider` class that discovers SKILL.md files from filesystem directories, validates their format, and exposes `load_skill`, `read_skill_resource`, and `run_skill_script` tools to agents.

The practical upshot for anyone building skills today is that your investment is protected. A skill you write this afternoon will work across every major AI coding tool on the market — and probably several that haven’t launched yet. That kind of portability guarantee is rare in tech, and it’s the single strongest argument for the SKILL.md format over every proprietary alternative.

## The Open Standard

Portability of skills isn’t just a side effect; it’s an explicit design goal backed by a formal specification. The Agent Skills open standard, hosted at [agentskills.io](https://agentskills.io/home), defines file structure, YAML frontmatter schema, Markdown body conventions, and discovery/routing behavior.

Maintained as an open-source effort (github.com/agentskills/agentskills), the spec isn’t locked to Anthropic or any vendor. Its goal is for skills written with time to work across any compliant agent runtime.

The spec doesn’t prescribe internal routing. It defines the interface (frontmatter metadata) but leaves implementation to agent runtimes. Claude Code uses LLM-based routing; others might use embeddings, classifiers, or hybrids. The spec guarantees that compliant agents can discover and use skills following the structure, but it doesn’t dictate their usage.

>

Skill authors should stick to the spec. Avoid Claude Code-specific behaviors or undocumented features in SKILL.md. Well-written, standard-structured instructions increase the chances of your skill working wherever the spec is supported.

## Comparison

Press enter or click to view image in full size

when to use what-in summary,

**SKILL.md** when you need the agent to follow a specific procedure for a specific type of task—and especially when that procedure involves running scripts, loading reference docs, or producing structured output.

**AGENTS.md** when you want every AI tool touching your repo to know your project conventions—coding style, architecture decisions, and build commands—regardless of which tool it is.

**CLAUDE.md / .cursor/rules / copilot-instructions. md** when you’re committed to one tool and want the tightest integration with that specific platform. Many teams just symlink these to their AGENTS.md.

**Custom GPTs** when you need something quick and no-code for a non-technical audience. But know that you’re trading away portability, version control, scripts, and progressive disclosure.

## Testing and Debugging Skills

Building a SKILL.md is one thing, but knowing if it works is another. Since skill routing is LLM-based and non-deterministic, testing requires a different mindset than traditional code.

The simplest debugging technique is to ask the agent. After sending a request that should trigger your skill, ask: “Which skill did you use?” or “Did you use the deploy-pipeline skill?” The agent will tell you, as it knows its context. This is the fastest way to check if routing works.

For undertriggering issues, start with the description. Often, the problem is that the description lacks enough trigger words or synonyms. If your skill handles Docker deployments but the description only says “containerized deployments,” a user asking about “Docker” might not trigger it. The fix is usually to add more specific trigger terms.

For overtriggering issues, add negative boundaries. The “Do NOT use for…” pattern in the description is your primary tool. Be explicit about what your skill doesn’t handle.

A practical testing workflow is to create a test checklist of 10–15 natural language requests, some that should trigger the skill, some that shouldn’t, and some ambiguous edge cases. Run through them manually and track which ones route correctly. This gives you a triggering accuracy score to improve iteratively.

For debugging the skill body, run a task, evaluate the output, and refine the instructions. Pay attention to where the agent deviates from your intended output. Ambiguous or missing instructions are usually the cause. Adding a “Common Gotchas” or “Critical Rules” section to your SKILL.md often fixes these issues by making the implicit explicit.

In Claude Code, you can view the expanded system prompt to see your skill’s metadata as it appears to the model. This can reveal formatting or truncation issues not visible in the SKILL.md file.

In the next blog, we will build our SKILL.md, and I will walk through the complete process of building and running it.

Hope this blog was useful, and see you soon in my next part…

until then, take care!! Please leave your experiences and feedback

Here are some references, for more reading

- Anthropic. “Introducing Agent Skills.” [https://www.anthropic.com/news/skills](https://www.anthropic.com/news/skills)
- Zhang, B., Lazuka, K., and Murag, M. “Equipping Agents for the Real World with Agent Skills.” Anthropic Engineering Blog, October 16, 2025. [https://www.anthropic.com/engineering/equipping-agents-for-the-real-world-with-agent-skills](https://www.anthropic.com/engineering/equipping-agents-for-the-real-world-with-agent-skills)
- Agent Skills Specification. [agentskills.io](http://agentskills.io/), . [https://agentskills.io/specification](https://agentskills.io/specification)
- Agent Skills GitHub Repository. Anthropic. [https://github.com/anthropics/skills](https://github.com/anthropics/skills)
- Agent Skills Open Standard Repository. [https://github.com/agentskills/agentskills](https://github.com/agentskills/agentskills)
- Anthropic. “Agent Skills — Claude API Documentation.” [https://platform.claude.com/docs/en/agents-and-tools/agent-skills/overview](https://platform.claude.com/docs/en/agents-and-tools/agent-skills/overview)
