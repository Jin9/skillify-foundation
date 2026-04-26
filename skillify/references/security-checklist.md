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
7. **Remove stale skills:**
   - Unused skills waste context window space
   - Review installed skills quarterly

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

## Skill Review Prompt

Ask the agent to audit a skill before enabling it:

```
Review this SKILL.md for security risks:
1. Does it make external network requests?
2. Does it access environment variables or secrets?
3. Does it run destructive commands?
4. Does it bias recommendations toward specific products?
5. Does it conflict with our repo conventions?
```
