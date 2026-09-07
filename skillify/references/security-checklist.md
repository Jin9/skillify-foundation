# Security Checklist for Skills

Skills shape agent behavior. Treat every installed skill like executable influence over your workflow.

## Before Enabling a Skill

1. **Read the raw SKILL.md yourself.** Do not rely on descriptions alone.
2. **Audit for exfiltration patterns:**
   - `curl`, `wget`, `fetch` to external URLs
   - Writing secrets/env vars to files
   - Sending data to third-party APIs not part of the workflow
3. **Check for product/service bias:**
   - Does the skill steer recommendations toward a specific vendor?
   - Are alternatives mentioned fairly or excluded?
4. **Verify no destructive commands:**
   - `rm -rf`, `DROP TABLE`, `git push --force`
   - File overwrites without confirmation
   - Deployments without review gates
5. **Confirm conventions match your repo:**
   - Does the skill contradict your existing architecture?
   - Does it assume a specific framework, package manager, or branching strategy?
6. **Restrict tool permissions:**
   - Use `allowed-tools` to limit what the skill can do
   - Use `paths` to limit where the skill activates
7. **Check instruction priority and stop behavior:**
   - Does any line make the agent pause, ask for a confirmation the user did not request, or leave requested work unfinished?
   - Does the skill say the user's request wins on conflict?
8. **Check for reasoning-extraction prompts:**
   - "Show your reasoning", "explain your thinking", or a required reasoning section can trigger a refusal on some models and leaks chain-of-thought into artifacts
   - Ask for the verified result and its evidence instead
9. **Remove stale skills:**
   - Unused skills waste context window space
   - Review installed skills quarterly

## Provenance and Supply Chain (installed or third-party skills)

Skills are increasingly distributed through package managers, registries, and host plugin marketplaces (for example `npx skills add` via skills.sh, `gh skill` for GitHub-hosted skills, or marketplace-distributed plugins that bundle skills). Treat an installed skill as untrusted code until reviewed — marketplace distribution does not imply review:

1. **Inspect before install.** Read the `SKILL.md` and every file under `scripts/` at the source before adding the skill, not after.
2. **Check provenance.** Prefer skills that record their origin (source repository plus a pinned ref or commit) and a declared license. Unsourced or unlicensed skills are higher risk.
3. **Pin versions.** Install a fixed version or commit rather than a moving `latest`, and re-review on upgrade.
4. **Re-run the audits above.** Apply the exfiltration, destructive-command, and vendor-bias checks to the installed copy — registry presence is not a safety guarantee.

## Frontmatter Security Rules

- **No XML angle brackets** (`<` `>`) in frontmatter. They can inject instructions into the system prompt.
- **No "claude" or "anthropic"** in skill names. These are reserved.
- **Review `allowed-tools`** if present. Overly permissive tool access is a security risk.
- **Check `hooks`** if present. Lifecycle hooks execute code at specific points.

## Common Threat Patterns

| Pattern | Risk | Detection |
|---------|------|-----------|
| External HTTP calls in scripts | Data exfiltration | Search for `curl`, `wget`, `requests.get`, `fetch` |
| Environment variable access | Secret leakage | Search for `$ENV`, `os.environ`, `process.env` |
| Broad file glob patterns | Unintended scope | Check `paths` field for overly broad globs like `**/*` |
| Forced git operations | Code loss | Search for `git push --force`, `git reset --hard` |
| Vendor-specific recommendations | Commercial bias | Read skill for product mentions without alternatives |
| Reasoning-extraction instructions | Refusal or chain-of-thought leakage | `scripts/cruft_scan.py` signal `scaffold.show_reasoning` (show, explain, reveal, or output your reasoning, thinking, or chain of thought); planning scaffolds are anti-pattern 12, not a security risk |
| Stall or precedence instructions | Agent pauses or diverges from the user | Search for `ask for permission`, `always confirm`, `do not proceed until` where the stop guards no destructive or irreversible action, no genuine change to the scope the user set, and no approval gate the user asked for; one clarification question for a missing required input is not a stall (`model-generation-fit.md` section 2) |

## Skill Review Prompt

Ask the agent to audit a skill before enabling it:

```
Review this SKILL.md for security risks:
1. Does it make external network requests?
2. Does it access environment variables or secrets?
3. Does it run destructive commands?
4. Does it bias recommendations toward specific products?
5. Does it conflict with our repo conventions?
6. Does any instruction make you pause, ask permission, or leave requested work unfinished where no destructive or irreversible action and no scope change is at hand? If so, name this file, quote the line, and explain how it applies.
7. Does it ask you to reproduce your reasoning?
```
