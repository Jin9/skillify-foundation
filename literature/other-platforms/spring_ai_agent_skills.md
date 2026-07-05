# Spring AI Agentic Patterns (Part 1): Agent Skills - Modular, Reusable Capabilities

Source: https://spring.io/blog/2026/01/13/spring-ai-generic-agent-skills
Accessed: 2026-07-05
Category: other-platforms / LLM-agnostic framing
Provenance: re-captured 2026-07-05 via content-extraction proxy (the page renders its article body client-side, so DOM extraction of the saved HTML yields only chrome); previous captures contained no article body

## Why This Source Matters

Spring AI's LLM-agnostic treatment of agent skills as model-decoupled capability registries — evidence that skills should converge on input/output contracts rather than vendor prose, and a concrete tool-based integration (SkillsTool/ShellTools/FileSystemTools) of the agentskills.io standard on the JVM. Also documents current limitations (no sandboxing, no HITL, no versioning) that mirror skillify's security-checklist concerns. Tier 4 reference.

---

_"LLM-Agnostic Skills That Run in Your Environment"_

## What Are Agent Skills?

Agent Skills are modular capabilities packaged as Markdown files with YAML frontmatter. Each skill is a folder containing a `SKILL.md` file with metadata (`name` and `description`, at minimum) and instructions that tell an agent how to perform a specific task. Skills can also bundle scripts, templates, and reference materials. The frontmatter supports both simple string values and complex YAML structures (lists, nested objects) for advanced use cases.

```
my-skill/
├── SKILL.md          # Required: instructions + metadata
├── scripts/          # Optional: executable code
├── references/       # Optional: documentation
└── assets/           # Optional: templates, resources
```

Skills use progressive disclosure to manage context efficiently:

1. **Discovery**: At startup, agents load only the name and description of each available skill, just enough to know when it might be relevant.
2. **Activation**: When a task matches a skill's description, the agent reads the full `SKILL.md` instructions into context.
3. **Execution**: The agent follows the instructions, optionally loading referenced files or executing bundled code as needed.

This approach allows you to register hundreds of skills while keeping your context window lean.

> **💡Tip:** Find out more about Agent Skills in the [official specification](https://agentskills.io/specification).

## Why Use Agent Skills in Spring AI

**Seamless Integration** - Add Agent Skills to your existing Spring AI application by simply registering a few tools—no architectural changes required.

**Portability and Model Agnostic - No Vendor Lock-In** - Unlike implementations tied to specific LLM platforms, this Spring AI implementation works across many LLM providers, letting you switch models without rewriting code or skills.

**Reusable and Composable** - Skills can be shared across projects, version-controlled with your code, combined to create complex workflows, and extended with helper scripts and reference materials. Spring AI Skills seamlessly supports any existing Claude Code Skills.

**Related Spring AI Tools:** Agent Skills work well with other Spring AI tool-based features like [Dynamic Tool Discovery](https://spring.io/blog/2025/12/11/spring-ai-tool-search-tools-tzolov) for efficient tool selection and [Tool Argument Augmentation](https://spring.io/blog/2025/12/23/spring-ai-tool-argument-augmenter-tzolov) for capturing LLM reasoning during skill execution.

## How Spring AI Skills Work

Spring AI uses the [tool-based integration approach](https://agentskills.io/integrate-skills#integration-approaches), implementing tools that allow any LLM to trigger skills and access bundled assets. The implementation closely follows [Claude Code's](https://code.claude.com/docs/en/settings#tools-available-to-claude) tool specifications for `Skills`, `Bash`, and `Read`.

The core set of tools is: [SkillsTool](https://github.com/spring-ai-community/spring-ai-agent-utils/blob/main/spring-ai-agent-utils/docs/SkillsTool.md) (required), [ShellTools](https://github.com/spring-ai-community/spring-ai-agent-utils/blob/main/spring-ai-agent-utils/docs/ShellTools.md) (optional), and [FileSystemTools](https://github.com/spring-ai-community/spring-ai-agent-utils/blob/main/spring-ai-agent-utils/docs/FileSystemTools.md) (optional). SkillsTool provides a `Skill` function that enables AI models to discover and load specified skills on demand, working in conjunction with FileSystemTools (for reading reference files) and ShellTools (for executing helper scripts).

Skills operate through a three-step process:

### 1. Discovery (at startup)

During initialization, SkillsTool scans configured skills directories (such as `.claude/skills/`) and parses the YAML frontmatter from each `SKILL.md` file. It extracts the `name` and `description` fields to build a lightweight skill registry that's embedded directly into the `Skill` tool's description, making it visible to the LLM without consuming conversation context.

### 2. Semantic Matching (during conversation)

When a user makes a request, the LLM examines the skill descriptions embedded in the tool definition. If the LLM determines a user request semantically matches a skill's description, it invokes the `Skill` tool with the skill name as a parameter.

### 3. Execution (on skill invocation)

When the `Skill` tool is called, SkillsTool loads the full `SKILL.md` content from disk and returns it to the LLM along with the skill's base directory path. The LLM then follows the instructions in the skill content. If the skill references additional files or helper scripts, the LLM uses FileSystemTools' `Read` function or ShellTools' `Bash` function to access them on demand.

## Skills in Action

This section demonstrates how skills work in practice with real-world examples.

### Example: Skills with References and Scripts

Step 3's on-demand loading becomes powerful when skills bundle additional resources. Skills can include reference files with supplementary instructions and executable scripts for data processing—all loaded only when needed.

Here's an example from the my-skill skill that includes a YouTube transcript extraction helper and supplementary `research_methodology.md` instructions:

**Skill Directory Structure:**

```
.claude/skills/my-skill/
├── SKILL.md
├── scripts/
│   └── get_youtube_transcript.py
└── research_methodology.md
```

**In SKILL.md:**

```
...
**If concept is unfamiliar or requires research:**
Load `research_methodology.md` for detailed guidance.

**If user provides YouTube video:**
Call `uv run scripts/get_youtube_transcript.py <video_url_or_id>`
for video's transcript.
...
```

When a user asks "Explain the concepts from this video: https://youtube.com/watch?v=abc123. Follow the research methodology", the AI:

1. Invokes the my-skill skill and loads its SKILL.md content
2. Recognizes the need for research methodology and uses `Read` to load `research_methodology.md`
3. Recognizes the YouTube URL and uses `Bash` to execute the helper script via ShellTools
4. Uses the video transcript to explain the concepts following the research methodology instructions

The script code never enters the context window—only the output does, making this approach highly token-efficient.

💡**Demo** Check out the [Skills-Demo](https://github.com/spring-ai-community/spring-ai-agent-utils/tree/main/examples/skills-demo) that implements this workflow.

> ⚠️ **Security Note:** Scripts execute directly on your local machine without sandboxing. You'll need to pre-install any required runtimes (Python, Node.js, etc.). For safer operation, consider running your agentic application in a container.

## Getting Started

Ready to add Agent Skills to your Spring AI project?

### 1. Add the dependency:

```xml
<dependency>
    <groupId>org.springaicommunity</groupId>
    <artifactId>spring-ai-agent-utils</artifactId>
    <version>0.4.2</version>
</dependency>
```

**Note:** For the latest stable release, check the [GitHub releases page](https://github.com/spring-ai-community/spring-ai-agent-utils/releases).

**Note:** You need Spring-AI `2.0.0-M2+`.

### 2. Configure your agent:

```java
@SpringBootApplication
public class Application {

    @Bean
    CommandLineRunner demo(ChatClient.Builder chatClientBuilder) {
        return args -> {
            ChatClient chatClient = chatClientBuilder
                .defaultToolCallbacks(SkillsTool.builder()
                    .addSkillsDirectory(".claude/skills")
                    .build())
                .defaultTools(FileSystemTools.builder().build())
                .defaultTools(ShellTools.builder().build())
                .build();

            String response = chatClient.prompt()
                .user("Your task here")
                .call()
                .content();
        };
    }
}
```

> **💡 Production Tip:** For packaged applications, you can load skills from the classpath using Spring Resources:
>
> ```java
> .defaultToolCallbacks(SkillsTool.builder()
>     .addSkillsResource(resourceLoader.getResource("classpath:.claude/skills"))
>     .build())
> ```
>
> This is particularly useful when distributing skills as part of your JAR/WAR deployment.

### 3. Create your first skill:

```bash
mkdir -p .claude/skills/code-reviewer
cat > .claude/skills/code-reviewer/SKILL.md << 'EOF'
---
name: code-reviewer
description: Reviews Java code for best practices, security issues, and Spring Framework conventions. Use when user asks to review, analyze, or audit code.
---

# Code Reviewer

## Instructions
When reviewing code:
1. Check for security vulnerabilities (SQL injection, XSS, etc.)
2. Verify Spring Boot best practices (proper use of @Service, @Repository, etc.)
3. Look for potential null pointer exceptions
4. Suggest improvements for readability and maintainability
5. Provide specific line-by-line feedback with code examples
EOF
```

### 4. Use your agent with the skill:

```java
String response = chatClient.prompt()
    .user("Review this controller class for best practices: " +
          "src/main/java/com/example/UserController.java")
    .call()
    .content();

System.out.println(response);
```

When you run this, the LLM will:

1. Match "Review this controller" with the code-reviewer skill's description
2. Invoke the `Skill` tool to load the full instructions from `SKILL.md`
3. Use `Read` tool (from FileSystemTools) to access the UserController.java file
4. Follow the review instructions and provide detailed feedback

The skill's instructions guide the LLM's behavior without you hardcoding review logic into your prompts—just update the skill file to change how reviews work.

## Current Limitations

While the Spring AI Agent Skills implementation is powerful and flexible, there are some current limitations to be aware of:

**Script Execution Security** - Scripts executed via ShellTools run directly on your local machine without sandboxing. This means potentially unsafe code could access your filesystem, network, or system resources. Always review skill scripts before use, especially from third-party sources. Consider running your agentic application in a containerized environment (Docker, Kubernetes) to limit exposure.

**No Human-in-the-Loop** - Currently, there's no built-in mechanism to require human approval before executing skills or scripts. The LLM can invoke any registered skill and execute any bundled script automatically. For production environments handling sensitive operations, you may need to implement custom approval workflows using Spring AI's tool callback mechanisms, for example, a `ToolCallback` wrapper.

**Limited Skill Versioning** - There's currently no built-in versioning system for skills. If you update a skill's behavior, all applications using that skill will immediately use the new version. For production deployments, consider implementing your own versioning strategy through directory structure (e.g., `.claude/skills/v1/`, `.claude/skills/v2/`).

## Related: Anthropic's Native Skills API

Spring AI also integrates with Anthropic's native Skills API, which provides a different approach:

* Skills run in Anthropic's sandboxed cloud container (no network access, pre-installed packages only)
* Pre-built document generation: Excel, PowerPoint, Word, PDF
* Skills are uploaded to Anthropic's servers and shared across your workspace
* Requires Claude models (Sonnet 4, Sonnet 4.5, Opus 4)

**Key difference:** Anthropic Skills run in Anthropic's cloud infrastructure; Generic Agent Skills run in your environment.

Use Anthropic's native Skills when you need secure, sandboxed execution or the pre-built document generation capabilities. Use Generic Agent Skills when you need LLM portability, local resource access, or want skills bundled with your application.

**Can you use both?** Yes. You can use Anthropic's native Skills for document generation while using Generic Agent Skills for other portable capabilities in the same application. They serve different purposes and can complement each other.

See the [Anthropic Skills in Spring AI](https://spring.io/blog/2026/01/spring-ai-anthropic-agent-skills) post for details.

## Conclusion

Agent Skills bring modular, reusable capabilities to Spring AI applications without vendor lock-in. By providing domain knowledge on demand, you can update agent behavior without code changes, share skills across projects, and switch between LLM providers seamlessly.

The spring-ai-agent-utils implementation makes this pattern accessible to Java developers with a simple, tool-based approach. Whether building coding assistants, documentation generators, or domain-specific agents, skills provide a scalable foundation for organizing agent knowledge.

**This is just the beginning.** The other posts in this series dive deeper into advanced agentic patterns:

**Series links:**

* [**AskUserQuestionTool**](https://spring.io/blog/2026/01/16/spring-ai-ask-user-question-tool) - Interactive workflows where agents gather user preferences during execution
* [**TodoWriteTool**](https://spring.io/blog/2026/01/20/spring-ai-agentic-patterns-3-todowrite/) - Transparent, trackable agent workflows with multi-step task management
* [**Subagent Orchestration**](https://spring.io/blog/2026/01/27/spring-ai-agentic-patterns-4-task-subagents) - Hierarchical multi-agent architectures with dedicated context windows
* [**A2A Integration**](https://spring.io/blog/2026/01/29/spring-ai-agentic-patterns-a2a-integration) - Building interoperable agents with the Agent2Agent protocol
* **Subagent Extension Framework** (coming soon) - Protocol-agnostic agent orchestration (A2A, MCP, custom)

Start exploring with the [example projects](https://github.com/spring-ai-community/spring-ai-agent-utils/tree/main/examples) or dive into the [Agent Skills specification](https://agentskills.io/specification) to learn more.

## Resources

### Spring AI Agent Utils Toolkit

* **GitHub Repository**: [spring-ai-agent-utils](https://github.com/spring-ai-community/spring-ai-agent-utils)
* **Full Documentation**: [README.md](https://github.com/spring-ai-community/spring-ai-agent-utils/blob/main/spring-ai-agent-utils/README.md)
* **Tool Documentation** - Tools covered in this post: [SkillsTool](https://github.com/spring-ai-community/spring-ai-agent-utils/blob/main/spring-ai-agent-utils/docs/SkillsTool.md), [FileSystemTools](https://github.com/spring-ai-community/spring-ai-agent-utils/blob/main/spring-ai-agent-utils/docs/FileSystemTools.md), [ShellTools](https://github.com/spring-ai-community/spring-ai-agent-utils/blob/main/spring-ai-agent-utils/docs/ShellTools.md)
* **Spring AI Documentation**: [docs.spring.io/spring-ai](https://docs.spring.io/spring-ai/reference/)

### Example Projects

* [skills-demo](https://github.com/spring-ai-community/spring-ai-agent-utils/tree/main/examples/skills-demo) - Focused skills demonstration (this post)
* [code-agent-demo](https://github.com/spring-ai-community/spring-ai-agent-utils/tree/main/examples/code-agent-demo) - Full toolkit integration (Parts 2-3)
* [subagent-demo](https://github.com/spring-ai-community/spring-ai-agent-utils/tree/main/examples/subagent-demo) - Hierarchical agents and A2A integration (Parts 4-5)

### Agent Skills

* **Specification**: [agentskills.io](https://agentskills.io/specification)
* **Claude Code Documentation**: [code.claude.com/docs](https://code.claude.com/docs/en/skills)

### Related Spring AI Blogs

* [Dynamic Tool Discovery](https://spring.io/blog/2025/12/11/spring-ai-tool-search-tools-tzolov) - Achieve 34-64% token savings
* [Tool Argument Augmentation](https://spring.io/blog/2025/12/23/spring-ai-tool-argument-augmenter-tzolov) - Capture LLM reasoning during tool execution
