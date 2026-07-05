# Manage agent skills with GitHub CLI

Source: https://github.blog/changelog/2026-04-16-manage-agent-skills-with-github-cli/
Accessed: 2026-04-26
Category: codex-copilot / gh skill supply chain
Provenance: normalized 2026-07-05 from raw HTML capture; original Accessed date preserved

## Why This Source Matters

GitHub changelog introducing gh skill: preview-before-install, version pinning, and provenance capture. The canonical supply-chain source for skill distribution. Tier 1 official announcement.


[Back to changelog](https://github.blog/changelog/)

Agent skills are reshaping how developers work with AI coding agents. Today we’re launching `gh skill`, a new command in the GitHub CLI that makes it easy to discover, install, manage, and publish agent skills from GitHub repositories.

### What are agent skills?

Agent skills are portable sets of instructions, scripts, and resources that teach AI agents how to perform specific tasks. They follow the open [Agent Skills specification](https://agentskills.io), and work across multiple agent hosts including GitHub Copilot, Claude Code, Cursor, Codex, and Gemini CLI among others.

With the new `gh skill` command, you can now install agent skills in a single command, right from the GitHub CLI.

### Get started

Update the GitHub CLI to version [v2.90.0](https://github.com/cli/cli/releases/tag/v2.90.0) or later.

Then discover and install skills interactively:

```
# Browse skills in a repository and install them interactively
gh skill install github/awesome-copilot

# Or install a specific skill directly
gh skill install github/awesome-copilot documentation-writer

# Install a specific version using @tag
gh skill install github/awesome-copilot documentation-writer@v1.2.0

# Install at a specific commit SHA
gh skill install github/awesome-copilot documentation-writer@abc123def

# Discover skills
gh skill search mcp-apps

```

Skills are automatically installed to the correct directory for your agent host. You can target a specific agent and scope with flags:

```
gh skill install github/awesome-copilot documentation-writer --agent claude-code --scope user

```

### Version pinning and supply chain integrity

Agent skills are executable instructions that shape how AI agents behave. A skill that changes silently between installs is a supply chain risk. `gh skill` brings the same guarantees you expect from package managers to the skills ecosystem, using primitives GitHub already provides.

- **Tags and releases:** Every published release is tied to a git tag. `gh skill publish` offers to enable [immutable releases](https://docs.github.com/repositories/releasing-projects-on-github/about-releases), so release content cannot be altered after publication, even by admins.
- **Content-addressed change detection:** Each installed skill records the git tree SHA of its source directory. `gh skill update` compares local SHAs against the remote to detect real content changes, not just version bumps. By storing this information in skills [front-matter](https://agentskills.io/specification#frontmatter), versioning and pinning are portable too, so you (or your agent) can copy and paste the skill to different projects without losing the ability to track changes and update it.
- **Version pinning:** Lock a skill to a specific tag or commit SHA with `--pin`. Pinned skills are skipped during updates, so you upgrade deliberately, not accidentally.
- **Portable provenance via frontmatter:** When `gh skill` installs a skill, it writes tracking metadata (repository, ref, tree SHA) directly into the `SKILL.md` frontmatter. Because provenance data lives inside the skill file itself, it travels with the skill no matter where it ends up. Skills get moved, copied, and reorganized by users, agents, and scripts.

```
# Pin to a release tag
gh skill install github/awesome-copilot documentation-writer --pin v1.2.0

# Pin to a commit for maximum reproducibility
gh skill install github/awesome-copilot documentation-writer --pin abc123def

```

### Publish your own skills

If you maintain a skills repository, `gh skill publish` validates your skills against the [agentskills.io spec](https://agentskills.io/specification) and checks remote settings like tag protection, secret scanning, and code scanning. These settings are not required, but strongly recommended to improve the supply chain security of your repo.

Enabling immutable releases, for example, means even if someone gets control of your repository they cannot change existing releases, so users installing via tag pinning are fully protected. The `publish` command makes it trivial to enable these features.

```
# Validate all skills
gh skill publish

# Auto-fix metadata issues
gh skill publish --fix

```

## Keep skills up to date

`gh skill update` scans all known agent host directories, reads provenance metadata from each installed skill, and checks for upstream changes:

```
# Check for updates interactively
gh skill update

# Update a specific skill
gh skill update git-commit

# Update everything without prompting
gh skill update --all

```

## Supported agent hosts

| Host | Install command example |
|---|---|

| GitHub Copilot | `gh skill install OWNER/REPOSITORY SKILL` |
| Claude Code | `gh skill install OWNER/REPOSITORY SKILL --agent claude-code` |
| Cursor | `gh skill install OWNER/REPOSITORY SKILL --agent cursor` |
| Codex | `gh skill install OWNER/REPOSITORY SKILL --agent codex` |
| Gemini CLI | `gh skill install OWNER/REPOSITORY SKILL --agent gemini` |
| Antigravity | `gh skill install OWNER/REPOSITORY SKILL --agent antigravity` |

### Learn more

- Check out the [Agent Skills specification](https://agentskills.io).
- Join the discussion in [GitHub Community](https://github.com/orgs/community/discussions).
- Visit [the `gh_skill` documentation](https://cli.github.com/manual/gh_skill).
- Run `gh skill --help` to see all available commands.

`gh skill` is launching in public preview and it’s subject to change without notice.

*Editor’s note (April 17, 2026): Added additional options for learning more.*

Skills are installed at your own discretion. They are not verified by GitHub and may contain prompt injections, hidden instructions, or malicious scripts. We strongly recommend inspecting the content of skills before installation, which can be done via the `gh skill preview` command.

Join the [GitHub Community](https://github.com/orgs/community/discussions/categories/announcements).

## Related Posts

###   Apr.24 Improvement

[Notice about upcoming new format for GitHub App installation tokens](https://github.blog/changelog/2026-04-24-notice-about-upcoming-new-format-for-github-app-installation-tokens)

[actions](https://github.blog/changelog/2026/?label=actions) [application security](https://github.blog/changelog/2026/?label=application-security) [client apps](https://github.blog/changelog/2026/?label=client-apps) [copilot](https://github.blog/changelog/2026/?label=copilot)

###   Apr.24 Release

[GPT-5.5 is generally available for GitHub Copilot](https://github.blog/changelog/2026-04-24-gpt-5-5-is-generally-available-for-github-copilot)

[copilot](https://github.blog/changelog/2026/?label=copilot)

###   Apr.24 Improvement

[Inline agent mode in preview and more in GitHub Copilot for JetBrains IDEs](https://github.blog/changelog/2026-04-24-inline-agent-mode-in-preview-and-more-in-github-copilot-for-jetbrains-ides)

[copilot](https://github.blog/changelog/2026/?label=copilot)

###   Apr.23 Improvement

[Copilot Chat improvements for pull requests](https://github.blog/changelog/2026-04-23-copilot-chat-improvements-for-pull-requests)

[copilot](https://github.blog/changelog/2026/?label=copilot)

###   Apr.23 Release

[Copilot cloud agent fields added to usage metrics](https://github.blog/changelog/2026-04-23-copilot-cloud-agent-fields-added-to-usage-metrics)

[account management](https://github.blog/changelog/2026/?label=account-management) [copilot](https://github.blog/changelog/2026/?label=copilot) [enterprise management tools](https://github.blog/changelog/2026/?label=enterprise-management-tools)

###   Apr.23 Improvement

[Better debugging with GitHub Copilot on the web](https://github.blog/changelog/2026-04-23-better-debugging-with-github-copilot-on-the-web)

[copilot](https://github.blog/changelog/2026/?label=copilot)

###   Apr.23 Release

[View and manage agent sessions from issues and projects](https://github.blog/changelog/2026-04-23-view-and-manage-agent-sessions-from-issues-and-projects)

[copilot](https://github.blog/changelog/2026/?label=copilot) [projects & issues](https://github.blog/changelog/2026/?label=projects-and-issues)

###   Apr.22 Retired

[Pausing new self-serve signups for GitHub Copilot Business](https://github.blog/changelog/2026-04-22-pausing-new-self-serve-signups-for-github-copilot-business)

[copilot](https://github.blog/changelog/2026/?label=copilot)

###   Apr.22 Improvement

[GitHub Copilot for Jira: Our latest enhancements](https://github.blog/changelog/2026-04-22-github-copilot-for-jira-our-latest-enhancements)

[collaboration tools](https://github.blog/changelog/2026/?label=collaboration-tools) [copilot](https://github.blog/changelog/2026/?label=copilot)

## Subscribe to our developer newsletter

 Discover tips, technical guides, and best practices in our biweekly newsletter just for devs.

 Back to top
