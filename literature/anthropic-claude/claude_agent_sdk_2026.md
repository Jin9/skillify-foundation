# Claude Agent SDK (2026)

Source: https://code.claude.com/docs/llms.txt
Accessed: 2026-05-29
Category: anthropic-claude / agent SDK vs skills
Provenance: synthesized web-research summary written 2026-05-29 (not a page scrape); header added 2026-07-05

## Why This Source Matters

Summary of the 2026 Claude Agent SDK split: the SDK owns runtime orchestration (MCP, context, hooks) while Skills carry domain knowledge. Frames what belongs in a skill versus the host runtime. Tier 5 synthesis note — verify against official docs before citing.


The **Claude Agent SDK** (formerly known as the "Claude Code SDK") is a developer toolkit and library that exposes the autonomous agent harness powering Claude Code as a programmable interface for Python and TypeScript.

### Overview
*   **Purpose:** Unlike the standard Claude API, which is a request-response interface for single-turn tasks, the Agent SDK provides a persistent, stateful agentic loop. It handles tool execution, context management, and task orchestration automatically, allowing you to build agents that can perform autonomous actions like running shell commands, reading/writing files, and searching the web.
*   **Renaming:** Anthropic officially renamed the "Claude Code SDK" to the "Claude Agent SDK" in March 2026 to reflect its broader utility beyond just coding tasks.
*   **Core Functionality:** The SDK includes built-in tools (Bash, Read, Write, WebSearch, etc.), subagent orchestration, session persistence, and observability features.

### Getting Started
*   **Installation:**
    *   **TypeScript:** `npm install @anthropic-ai/claude-agent-sdk`
    *   **Python:** `pip install claude-agent-sdk`
*   **Primary Interface:** The core of the SDK is the `query()` async generator, which accepts a prompt and configuration options, yielding typed messages as the agent works toward a task.
*   **Configuration:** The SDK supports Claude Code’s filesystem-based configuration (e.g., `.claude/` directories) and integrates with environment variables for API key management (`ANTHROPIC_API_KEY`) and platform-specific providers like Amazon Bedrock or Google Vertex AI.

### Important 2026 Updates & Billing
*   **Credit Changes:** Effective **June 15, 2026**, Anthropic has separated Agent SDK and `claude -p` usage from subscription plans (Pro, Team, Enterprise). This usage now draws from a dedicated monthly Agent SDK credit (or standard pay-as-you-go API billing), meaning it is no longer covered by the "unlimited" usage pools associated with standard consumer/team subscriptions.
*   **Production Readiness:** Developers are encouraged to implement production patterns—such as durable state (using Postgres/Redis), hard cost caps, tool permission scoping, and evaluation hooks—rather than relying on the ephemeral sessions used for local development.
