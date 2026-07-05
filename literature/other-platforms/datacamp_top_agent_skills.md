# The Top 100+ Agent Skills For OpenClaw, Codex and Claude

Source: https://www.datacamp.com/blog/top-agent-skills
Accessed: 2026-07-05
Category: other-platforms / community catalog
Provenance: re-captured 2026-07-05 via content-extraction proxy (direct fetch is Cloudflare-gated); may be abridged relative to the original page; previous capture was a 139-word Cloudflare challenge stub

## Why This Source Matters

Broad 2026 community catalog of production agent skills across research, coding, DevOps, ML, security, and communication categories, spanning the OpenClaw/ClawHub, Codex, Claude Code, and Cursor ecosystems. Evidences the skill-marketplace/distribution shift (ClawHub downloads, npx installs), documents progressive-disclosure as the standard activation architecture, and carries an explicit supply-chain warning about malicious SKILL.md bundles. Tier 5 catalog source — inspiration and ecosystem signal only; verify before adopting.

---

AI agents are no longer just chat interfaces that generate text. They are execution engines, and the reason they have become this powerful is agent skills.

Agent skills are modular capability bundles defined through structured `SKILL.md` files that give AI agents real-world actions.

Instead of only responding with text, an agent equipped with skills can search the web, run code, query databases, deploy infrastructure, fine-tune models, send emails, automate workflows, and more. You can think of skills as apps for AI agents. Each skill adds a focused, executable capability, documented in its own markdown reference, that turns a language model into a system that can actually do things.

In this article, we have curated a list of 100+ agent skills that can be used with [OpenClaw](https://www.datacamp.com/tutorial/moltbot-clawdbot-tutorial), [Claude Code](https://www.datacamp.com/tutorial/claude-code), [OpenAI Codex](https://www.datacamp.com/blog/gpt-5-3-codex) agents, Cursor CLI agents, and other agent-first environments. These skills span research, coding, infrastructure, machine learning, security, communication, and creative production. The goal is simple: to give you a practical map of what is possible when agents move from conversation to execution.

At the end, we also include a bonus section highlighting the most downloaded skills on ClawHub, so you can see which tools the community is actively using with OpenClaw.

## 1. Top Search and Research Agent Skills

Skills that help agents discover information, query structured sources, and extract reliable insights from technical and scientific data.

* [arxiv-watcher](https://github.com/openclaw/skills/tree/main/skills/rubenfb23/arxiv-watcher/SKILL.md): Search and summarize papers from ArXiv.
* [pubmed-edirect](https://github.com/openclaw/skills/blob/main/skills/killgfat/pubmed-edirect/SKILL.md): Query PubMed for peer-reviewed biomedical and scientific literature.
* [wikipedia](https://github.com/openclaw/skills/blob/main/skills/rachmann-alexander/wikipedia-oc/SKILL.md): Searches, retrieves, and summarizes content from English Wikipedia
* [google-search](https://github.com/openclaw/skills/blob/main/skills/mxfeinberg/google-search/SKILL.md): Search the web using Google Custom Search Engine (PSE).
* [google-search-grounding](https://github.com/openclaw/skills/blob/main/skills/shaharsha/google-search-grounding/SKILL.md): Google web search via Gemini Search Grounding.
* [serper-search](https://github.com/openclaw/skills/blob/main/skills/samoppakiks/serper-search/SKILL.md): Google Search via Serper.dev API.
* [web-scraper-as-a-service](https://github.com/openclaw/skills/blob/main/skills/seanwyngaard/web-scraper-as-a-service/SKILL.md): Build client-ready web scrapers with clean data output.
* [exa-web-search-free](https://github.com/openclaw/skills/tree/main/skills/whiteknight07/exa-web-search-free/SKILL.md): Free AI search via Exa.
* [newsapi-search](https://github.com/openclaw/skills/blob/main/skills/hegghammer/newsapi-search/SKILL.md): Query global news sources for trend tracking and signal analysis.
* [brightdata](https://github.com/openclaw/skills/blob/main/skills/meirkad/bright-data/SKILL.md): Web scraping and search via Bright Data API.

## 2. Top Coding Orchestration and Developer Copilots Skills

Skills for running multi-step coding projects and coordinating "agent teams" end-to-end.

* [buildlog](https://github.com/openclaw/skills/tree/main/skills/espetey/buildlog/SKILL.md): Records and exports coding sessions as shareable "build logs."
* [cc-godmode](https://github.com/openclaw/skills/tree/main/skills/cubetribe/cc-godmode/SKILL.md): Orchestrates multi-agent software work with self-managed coordination.
* [codebuddy-code](https://github.com/openclaw/skills/tree/main/skills/pmwalkercao/codebuddy-code/SKILL.md): Installs/configures a CodeBuddy-style coding CLI assistant.
* [debug-pro](https://github.com/openclaw/skills/tree/main/skills/cmanfre7/debug-pro/SKILL.md): Systematic debugging methodology and language-specific debugging.
* [coder-workspaces](https://github.com/openclaw/skills/tree/main/skills/developmentcats/coder-workspaces/SKILL.md): Manages Coder workspaces for remote/devcontainer-like workflows.
* [cursor-agent](https://github.com/openclaw/skills/tree/main/skills/swiftlysingh/cursor-agent/SKILL.md): A comprehensive skill for using the Cursor CLI agent.
* [ec-task-orchestrator](https://github.com/openclaw/skills/tree/main/skills/henrino3/ec-task-orchestrator/SKILL.md): Autonomous multi-agent task orchestration.
* [codex-orchestration](https://github.com/openclaw/skills/tree/main/skills/shanelindsay/codex-orchestration/SKILL.md): Provides a general orchestration layer for Codex-style agents.
* [codex-quota](https://github.com/openclaw/skills/tree/main/skills/odrobnik/codex-quota/SKILL.md): Checks Codex quota/rate limits so agents don't hit hard stops.
* [coding-agent](https://github.com/openclaw/skills/tree/main/skills/steipete/coding-agent/SKILL.md): Runs common coding agents/CLIs (Codex, Claude Code, etc.) from one skill.

## 3. Top Git, GitHub, PRs, and Repo Intelligence Agent Skills

Skills for version control, PR automation, and understanding what changed in a repo.

* [auto-pr-merger](https://github.com/openclaw/skills/tree/main/skills/autogame-17/auto-pr-merger/SKILL.md): Automates checking and merging PRs when rules are met.
* [backup](https://github.com/openclaw/skills/tree/main/skills/jordanprater/backup/SKILL.md): Backs up and restores an agent's config/skills/settings.
* [bat-cat](https://github.com/openclaw/skills/tree/main/skills/arnarsson/bat-cat/SKILL.md): Provides a "bat"-style file viewer with git-aware output.
* [bitbucket-automation](https://github.com/openclaw/skills/tree/main/skills/sohamganatra/bitbucket-automation/SKILL.md): Automates Bitbucket repo and PR workflows.
* [commit-analyzer](https://github.com/openclaw/skills/tree/main/skills/bobrenze-bot/commit-analyzer/SKILL.md): Analyzes commit patterns to understand change risk and behavior.
* [conventional-commits](https://github.com/openclaw/skills/tree/main/skills/bastos/conventional-commits/SKILL.md): Formats commit messages to the Conventional Commits standard.
* [deepwiki](https://github.com/openclaw/skills/tree/main/skills/arun-8687/deepwiki/SKILL.md): Queries repo documentation/wiki via an MCP-backed "deep wiki" interface.
* [gitclassic](https://github.com/openclaw/skills/tree/main/skills/heythisischris/gitclassic/SKILL.md): Uses a lightweight GitHub browser that works well for agents.
* [gitclaw](https://github.com/openclaw/skills/tree/main/skills/marian2js/gitclaw/SKILL.md): Syncs an agent workspace into a GitHub repo as a backup/mirror.
* [github](https://github.com/openclaw/skills/tree/main/skills/steipete/github/SKILL.md): Operates GitHub via gh for issues, PRs, and repo actions.

## 4. DevOps and Cloud Operations Agent Skills

Skills for provisioning, deploying, monitoring, and securing cloud infrastructure.

* [docker-essentials](https://github.com/openclaw/skills/blob/main/skills/arnarsson/docker-essentials/SKILL.md): Build, tag, and run containers using clean, production-ready Docker workflows.
* [k8-multicluster](https://github.com/openclaw/skills/blob/main/skills/rohitg00/k8-multicluster/SKILL.md): Manage multiple Kubernetes clusters and switch contexts safely.
* [nginx-config-creator](https://github.com/openclaw/skills/blob/main/skills/xieyuanqing/nginx-config-creator/SKILL.md): Generate reverse-proxy configs for Nginx/OpenResty deployments.
* [appdeploy](https://github.com/openclaw/skills/tree/main/skills/avimak/appdeploy/SKILL.md): Deploys web apps, including backend and database pieces.
* [aws-infra](https://github.com/openclaw/skills/tree/main/skills/bmdhodl/aws-infra/SKILL.md): Guides AWS infra work using CLI and best practices.
* [aws-ecs-monitor](https://github.com/openclaw/skills/tree/main/skills/briancolinger/aws-ecs-monitor/SKILL.md): Monitors ECS services and CloudWatch signals for production health.
* [aws-security-scanner](https://github.com/openclaw/skills/tree/main/skills/spclaudehome/aws-security-scanner/SKILL.md): Scans AWS environments for common security issues.
* [azd-deployment](https://github.com/openclaw/skills/tree/main/skills/thegovind/azd-deployment/SKILL.md): Deploys container apps to Azure Container Apps using azd.
* [azure-cli](https://github.com/openclaw/skills/tree/main/skills/ddevaal/azure-cli/SKILL.md): Manages Azure resources through Azure CLI commands and flows.
* [hetzner](https://github.com/openclaw/skills/tree/main/skills/thesethrose/hetzner/SKILL.md): Controls Hetzner Cloud resources using hcloud.

## 5. Top Data Science and Machine Learning Agent Skills

Skills for prepping data, building training workflows, monitoring experiments, and shipping ML models into something usable.

* [peft](https://github.com/openclaw/skills/blob/main/skills/desperado991128/peft/SKILL.md): Fine-tune LLMs with LoRA/QLoRA adapters, merge/swap adapters, and run lightweight post-training workflows.
* [wandb-monitor](https://github.com/openclaw/skills/blob/main/skills/chrisvoncsefalvay/wandb-monitor/SKILL.md): Monitor and compare Weights & Biases runs to spot training issues and performance regressions fast.
* [senior-computer-vision](https://github.com/openclaw/skills/blob/main/skills/alirezarezvani/senior-computer-vision/SKILL.md): Build CV pipelines end-to-end (datasets → training → eval → deployment) with export/optimization guidance (ONNX/TensorRT/CoreML).
* [senior-data-engineer](https://github.com/openclaw/skills/blob/main/skills/alirezarezvani/senior-data-engineer/SKILL.md?): Build scalable ETL/ELT pipelines and modern data infrastructure (Spark/Airflow/dbt/Kafka).
* [hugging-face-model-trainer](https://github.com/huggingface/skills/blob/main/skills/hugging-face-model-trainer/SKILL.md): Train/fine-tune LLMs with TRL methods (SFT/DPO/GRPO) on Hugging Face Jobs and export GGUF.
* [duckdb](https://github.com/openclaw/skills/blob/main/skills/camelsprout/duckdb-cli-ai-skills/SKILL.md): Do fast analytics on CSV/Parquet/JSON with DuckDB CLI.
* [senior-data-scientist](https://github.com/openclaw/skills/tree/main/skills/alirezarezvani/senior-data-scientist/SKILL.md): World-class data science skill.
* [data-analyst](https://github.com/openclaw/skills/blob/main/skills/oyi77/data-analyst/SKILL.md): Analyze data via SQL/spreadsheets, create charts, and generate decision-ready reports.
* [hugging-face-datasets](https://github.com/huggingface/skills/blob/main/skills/hugging-face-datasets/SKILL.md): Create/manage datasets on the Hub, including SQL-based querying/transforms with DuckDB.
* [hugging-face-evaluation](https://github.com/huggingface/skills/blob/main/skills/hugging-face-evaluation/SKILL.md): Add structured eval results to model cards and run/import benchmarks (vLLM/lighteval, etc.).

## 6. Security, Governance, and Safety Agent Skills

Skills that add guardrails, scanning, policy checks, and secure defaults to agents.

* [agentguard](https://github.com/openclaw/skills/tree/main/skills/manas-io-ai/agentguard/SKILL.md): Adds monitoring/guardrails to reduce risky agent behavior.
* [agentmemory](https://github.com/openclaw/skills/tree/main/skills/badaramoni/agentmemory/SKILL.md): Provides encrypted cloud memory for agents across devices.
* [clawscan](https://github.com/openclaw/skills/tree/main/skills/g0head/clawscan/SKILL.md): Scans skill bundles for red flags before installation/use.
* [clawsec-feed](https://github.com/openclaw/skills/tree/main/skills/davida-ps/clawsec-feed/SKILL.md): Pulls security advisories/CVE signals for continuous awareness.
* [clawskillshield](https://github.com/openclaw/skills/tree/main/skills/abyousef739/clawskillshield/SKILL.md): Runs a local-first scanner to detect suspicious skill behavior.
* [config-guardian](https://github.com/openclaw/skills/tree/main/skills/abdhilabs/config-guardian/SKILL.md): Validates config changes to prevent breaking or unsafe updates.
* [prompt-guard](https://github.com/openclaw/skills/tree/main/skills/seojoonkim/prompt-guard/SKILL.md): Defends against prompt injection and unsafe instruction following.
* [skill-flag](https://github.com/openclaw/skills/tree/main/skills/patfire94/skill-flag/SKILL.md): Detects malicious patterns/backdoors in skill code/instructions.
* [skill-scanner](https://github.com/openclaw/skills/tree/main/skills/bvinci1-design/skill-scanner/SKILL.md): Scans skills/MCP servers for spyware-like behavior and risks.
* [skills-audit](https://github.com/openclaw/skills/tree/main/skills/morozred/skill-audit/SKILL.md): Audits installed skills against policy and security checks.

## 7. Communication, Messaging, and Community Agent Skills

Skills that let agents post, reply, and manage conversations across chat platforms.

* [discord-voice](https://github.com/openclaw/skills/blob/main/skills/avatarneil/discord-voice/SKILL.md): Enables real-time voice conversations directly inside Discord voice channels.
* [giphy](https://github.com/openclaw/skills/tree/main/skills/minbang930/giphy/SKILL.md): Finds and sends context-matching GIFs in conversations.
* [mailchannels](https://github.com/openclaw/skills/tree/main/skills/ttulttul/mailchannels/SKILL.md): Sends email through MailChannels and handles signed ingestion flows.
* [google-messages-openclaw-skill](https://github.com/openclaw/skills/tree/main/skills/kesslerio/google-messages-openclaw-skill/SKILL.md): Enables SMS/RCS sending/receiving via Google Messages integration.
* [lark-integration](https://github.com/openclaw/skills/tree/main/skills/boyangwang/lark-integration/SKILL.md): Connects Lark/Feishu messaging into agent workflows via webhooks.
* [clawsignal](https://github.com/openclaw/skills/tree/main/skills/bmcalister/clawsignal/SKILL.md): Adds real-time agent messaging for alerts and coordination.
* [olvid-channel](https://github.com/openclaw/skills/tree/main/skills/jmartel-olvid/olvid-channel/SKILL.md): Integrates the Olvid secure messenger as an agent channel.
* [disclawd](https://github.com/openclaw/skills/tree/main/skills/alexerm/disclawd/SKILL.md): Connects to an agent-first Discord-like environment.
* [agent-mail](https://github.com/openclaw/skills/blob/main/skills/rimelucci/agent-mail/SKILL.md): Email inbox for AI agents.
* [whatsapp-styling-guide](https://github.com/openclaw/skills/tree/main/skills/rubenfb23/whatsapp-styling-guide/SKILL.md): Ensures agent WhatsApp messages follow consistent formatting rules.

## 8. Top Agent Skills For Notes, Knowledge, and Personal Knowledge Management

Skills for storing durable memory, writing notes, and pulling info from your knowledge systems.

* [logseq](https://github.com/openclaw/skills/tree/main/skills/juanirm/logseq/SKILL.md): Lets an agent read/write notes in a local Logseq vault.
* [notesctl-skill-for-openclaw](https://github.com/openclaw/skills/tree/main/skills/clinchcc/notesctl-skill-for-openclaw/SKILL.md): Controls Apple Notes via deterministic CLI-like operations.
* [openclaw-confluence-skill](https://github.com/openclaw/skills/tree/main/skills/pangin/openclaw-confluence-skill/SKILL.md): Uses Confluence Cloud REST APIs for searching and editing pages.
* [openclaw-nextcloud](https://github.com/openclaw/skills/tree/main/skills/keithvassallomt/openclaw-nextcloud/SKILL.md): Connects Nextcloud files/notes/tasks/calendar into the agent.
* [git-notes-memory](https://github.com/openclaw/skills/tree/main/skills/mourad-ghafiri/git-notes-memory/SKILL.md): Stores durable memory using git-notes across sessions.
* [memory-hygiene](https://github.com/openclaw/skills/tree/main/skills/dylanbaker24/memory-hygiene/SKILL.md): Cleans and optimizes vector memory to reduce drift/noise.
* [openclaw-feeds](https://github.com/openclaw/skills/tree/main/skills/nesdeq/openclaw-feeds/SKILL.md): Aggregates RSS feeds for research and daily monitoring.
* [hardcover](https://github.com/openclaw/skills/tree/main/skills/asaphko/hardcover/SKILL.md): Pulls reading lists and book metadata from Hardcover via API.
* [get-tldr](https://github.com/openclaw/skills/tree/main/skills/itobey/get-tldr/SKILL.md): Summarizes long content via a TL;DR summarization API.
* [essence-distiller](https://github.com/openclaw/skills/tree/main/skills/leegitw/essence-distiller/SKILL.md): Extracts the core ideas and "what matters" from messy content.

## 9. Agent Skills For Marketing, Publishing, and Social Growth

Skills to publish content, run ads, create links/QRs, and automate social workflows.

* [wordpress-publishing-skill-for-claude](https://github.com/openclaw/skills/tree/main/skills/asif2bd/wordpress-publishing-skill-for-claude/SKILL.md): Publishes content to WordPress from agent workflows.
* [wp-multi-tool](https://github.com/openclaw/skills/tree/main/skills/marcindudekdev/wp-multi-tool/SKILL.md): Audits WordPress health/performance and suggests fixes.
* [social-scheduler-extended](https://github.com/openclaw/skills/tree/main/skills/coolmanns/social-scheduler-extended/SKILL.md): Schedule and manage social media posts.
* [microsoft-ads-mcp](https://github.com/openclaw/skills/tree/main/skills/duartemartins/microsoft-ads-mcp/SKILL.md): Creates and manages Microsoft Ads campaigns via MCP tooling.
* [go2gg](https://github.com/openclaw/skills/tree/main/skills/rakesh1002/go2gg/SKILL.md): Shortens links, tracks clicks, and generates QR codes.
* [jo4](https://github.com/openclaw/skills/tree/main/skills/anandrathnas/jo4/SKILL.md): Generates short links + QR codes with analytics.
* [aisa-twitter-api](https://github.com/openclaw/skills/tree/main/skills/aisapay/aisa-twitter-api/SKILL.md): Searches X (Twitter) in real time and extracts relevant posts.
* [glasses-to-social](https://github.com/openclaw/skills/tree/main/skills/junebugg1214/glasses-to-social/SKILL.md): Turns smart-glasses photos into ready-to-post social captions.
* [share-usecase](https://github.com/openclaw/skills/tree/main/skills/josephl37/share-usecase/SKILL.md): Publishes your agent use case to a public directory for discovery.
* [openclaw-postsyncer](https://github.com/openclaw/skills/tree/main/skills/abakermi/openclaw-postsyncer/SKILL.md): Automates social posting workflows with a "post sync" routine.

## 10. Media, Images, Video, and Creative Production Agent Skills

Skills for generating images/videos, editing assets, and producing creative artifacts and diagrams.

* [image-router](https://github.com/openclaw/skills/tree/main/skills/dawe35/image-router/SKILL.md): Generates images via an API that can route to multiple models.
* [imagemagick](https://github.com/openclaw/skills/tree/main/skills/kesslerio/imagemagick/SKILL.md): Manipulates images (resize/convert/compose) with ImageMagick tooling.
* [avatar-video-messages](https://github.com/openclaw/skills/tree/main/skills/thewulf7/avatar-video-messages/SKILL.md): Creates avatar-style video messages from text prompts.
* [video-agent](https://github.com/openclaw/skills/tree/main/skills/michaelwang11394/video-agent/SKILL.md): Produces AI avatar videos using HeyGen's Video Agent API.
* [video-cog](https://github.com/openclaw/skills/tree/main/skills/nitishgargiitd/video-cog/SKILL.md): Supports long-form AI video production with multi-step planning.
* [voice-reply](https://github.com/openclaw/skills/tree/main/skills/stolot0mt0m/voice-reply/SKILL.md): Generates local text-to-speech voice replies (offline-friendly).
* [artifacts-builder](https://github.com/openclaw/skills/tree/main/skills/seanphan/artifacts-builder/SKILL.md): Builds multi-part artifacts (docs/assets) from a single agent plan.
* [excalidraw-diagrams](https://github.com/robtaylor/excalidraw-diagrams/blob/main/SKILL.md): Creates Excalidraw-style diagrams and architecture sketches.
* [manim-composer](https://github.com/openclaw/skills/tree/main/skills/inclinedadarsh/manim-composer/SKILL.md): Produces math/animation scenes using Manim-style composition.
* [morfeo-remotion-style](https://github.com/openclaw/skills/tree/main/skills/pauldelavallaz/morfeo-remotion-style/SKILL.md): Applies a consistent Remotion video style/template for output videos.

## Bonus: Top Downloaded ClawHub Skills

These are the most downloaded and widely used skills on [ClawHub](https://clawhub.ai/), based on current marketplace rankings and public download stats (accurate at the time of writing):

* [gog](https://clawhub.ai/steipete/gog): Google Workspace CLI for Gmail, Drive, Docs, and Sheets automation | **~29.4K downloads.**
* [tavily-search](https://clawhub.ai/arun-8687/tavily-search): AI-optimized real-time web search for up-to-date information retrieval | **~23.8K downloads**.
* [summarize](https://clawhub.ai/steipete/summarize): Summarizes URLs, PDFs, documents, and audio into concise outputs | **~22.4K downloads**.
* [github](https://clawhub.ai/steipete/github): Full GitHub CLI integration for issues, PRs, repos, and workflows | **~21.6K downloads**.
* [sonoscli](https://clawhub.ai/steipete/sonoscli): Controls Sonos speakers, including playback, volume, and room management | **~18.6K downloads**.
* [weather](https://clawhub.ai/steipete/weather): Fetches real-time weather conditions and forecasts for any location | **~18.6K downloads**.
* [ontology](https://clawhub.ai/oswalpalash/ontology): Enables structured knowledge graph exploration and semantic reasoning | **~18.1K downloads**.
* [notion](https://clawhub.ai/steipete/notion): Read/write Notion pages and databases directly from your agent | **~11.9K downloads**.
* [api-gateway](https://clawhub.ai/byungkyu/api-gateway): Unified managed API calling layer for secure third-party integrations | **~11.7K downloads**.
* [nano-banana-pro](https://clawhub.ai/steipete/nano-banana-pro): Advanced image generation and creative automation tool | **~11.6K downloads**.

## Final Thoughts

My obsession with agentic skills started when I completely switched my coding workflow from an IDE to OpenAI Codex. Codex is excellent at creating and using skills, but at first, I did not truly understand how to optimize them. I was unsure how `SKILL.md` files should be structured, what made one implementation better than another, and how to design skills that were clean, reusable, and production-ready.

Then I discovered ClawHub, a central marketplace for agent skills that works beyond just the OpenClaw ecosystem. It supports tools like OpenAI Code, OpenCode, Claude Code, and other agent-first environments. Since then, I have been installing skills directly from the hub using simple `npx` commands after reviewing them on the ClawHub web page.

In this article, we reviewed top agent skills across categories and real-world use cases. My advice is simple: start with the bonus section and install the most downloaded ClawHub skills first. They are widely adopted, community-trusted, and a strong foundation for building serious agent systems.

## Agent Skills FAQs

### What is an AI agent skill?

Think of an agent skill as an "npm package" or an app for your AI. Instead of just giving an AI a text prompt, a skill is a self-contained folder that bundles instructions, executable scripts, and reference materials. At its core is a `SKILL.md` file containing metadata and step-by-step procedures. This allows agents (like OpenClaw, Claude Code, or Cursor) to dynamically load domain-specific expertise—like how to deploy to AWS, run a security audit, or query a database—exactly when needed.

### How is an agent skill different from a standard prompt or an MCP server?

* **Prompts** rely entirely on the LLM's general knowledge and can produce inconsistent results. Skills bundle instructions with actual executable code and rules, ensuring the agent performs tasks consistently (idempotently) every time.
* **MCP (Model Context Protocol) Servers** handle the technical integration between AI systems and third-party APIs (like connecting to Slack or Postgres).
* **Skills** sit a layer above. They contain the _know-how_ or business logic, telling the agent _when_ to use those tools, how to sequence the actions, and how to format the output.

### Are agent skills safe to download from public marketplaces?

**You must treat them with caution.** Because skills are executable code, they run with the exact same permissions as your AI agent, meaning they can access your file system, read environment variables, and run shell commands. Recently, security researchers have flagged supply chain vulnerabilities on open marketplaces like ClawHub, where malicious `SKILL.md` files or hidden scripts attempted to exfiltrate API keys or install malware.

Always review the `SKILL.md` and scripts before installing a third-party skill, and consider running agents in sandboxed environments if they are using unverified community tools.

### How does an AI agent know when to use a specific skill?

Agent skills use an architecture called "progressive disclosure" to save on token costs.

* **Discovery:** When the agent starts, it only reads the YAML metadata (the skill's name and a short description) of all installed skills.
* **Semantic Matching:** When you ask a question or assign a task, the agent checks if your request matches any of the skill descriptions.
* **Execution:** If there's a match, the agent loads the full `SKILL.md` instructions and any associated scripts into its context window to execute the task. You can also explicitly trigger them using slash commands (e.g., `/deploy-app`).
