# skills.sh - AI Agent Skills Package Manager (2026)

**Skills.sh** is a community-driven open directory and package manager designed specifically for **AI agent skills**. Launched by Vercel in early 2026, it aims to standardize how developers discover, install, and manage the procedural knowledge and tools used by AI coding agents.

### What are "Agent Skills"?
Agent Skills are modular packages—typically containing a `SKILL.md` file with instructions and optional scripts—that extend an AI agent's capabilities.
*   **Purpose:** Instead of bloating an AI's context with massive system prompts, skills provide "dormant" manuals and tools that an agent can activate on demand.
*   **Format:** The format uses **progressive disclosure**. An agent reads high-level metadata first to decide if a skill is relevant, then accesses deeper instructions or executable scripts only when necessary.
*   **Use Cases:** Automating specialized tasks like browser testing, commit formatting, React/Next.js best practices, or custom deployment workflows.

### How Skills.sh Works
Skills.sh functions as the "npm for agents," providing a centralized registry and a Command Line Interface (CLI) to manage these capabilities.

*   **CLI Workflow:** Developers use `npx skills` to interact with the ecosystem:
    *   `npx skills find <query>`: Search for available skills.
    *   `npx skills add <package>`: Install a skill directly into an agent's project directory (e.g., `.claude/skills/`).
    *   `npx skills update`: Manage versions and updates for installed skills.
*   **Compatibility:** It is designed to work across various AI platforms, including **Claude Code, Cursor, GitHub Copilot, Aider, Windsurf, Gemini CLI, and others**.

### Key Considerations
*   **Community-Driven:** While the platform offers powerful automation, it is an open directory. The quality of community-submitted skills can vary.
*   **Security:** Because skills can include executable scripts (`/scripts` directory) that may run locally, users are advised to inspect third-party skills before installation to mitigate supply-chain risks and prompt injection.
*   **Evolution:** The ecosystem is still maturing regarding features like robust version pinning and dependency management compared to traditional package managers like npm or pip.
