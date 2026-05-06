---
name: configuring-sandbox-allowlist
description: >
  Edit the agent-scaffold sandbox hostname allowlist at
  docker/sandbox-proxy/filter. Adds, reviews, or removes hostname patterns,
  refusing wildcards, host-internal redirects, IP literals, and metadata-IP
  patterns. Reminds the user to rebuild the proxy image and run the
  egress-test verification before the next sandboxed implement run. Use
  when the user says "add domain to sandbox allowlist", "sandbox blocks
  HOST", "expand allowlist", "remove HOST from allowlist", "audit the
  sandbox allowlist", or after a postmortem identifies an allowlist gap.
  Does NOT rebuild the proxy itself, modify tinyproxy.conf, or weaken
  default-deny. Do NOT use for routing, firewall, or VPN configuration
  unrelated to the scaffold's tinyproxy filter.
---

# Configuring the sandbox hostname allowlist

## Purpose

The sandboxed implement runner is the scaffold's last line of defense for
credential-adjacent work. The allowlist at `docker/sandbox-proxy/filter`
is what turns a default-deny tinyproxy into a usable model-API egress
path. This skill edits that file with the right discipline: refuse
wildcards, refuse internal redirects, require rebuild + egress-test
verification, leave a paper trail.

## When to use this skill

- The user says "add `<host>` to the sandbox allowlist".
- A sandboxed implement run failed because tinyproxy denied a host the
  user trusts.
- A postmortem flagged an allowlist gap or an over-permissive entry.
- The user wants to audit the allowlist for hosts no longer in use.

Do NOT use this skill to:
- Rebuild the proxy image. Surface the rebuild command for the user.
- Modify `tinyproxy.conf` (Allow rules, FilterDefaultDeny, ConnectPort).
  Those changes weaken the default-deny posture and require a deliberate
  PR with team-lead review, not this skill.
- Configure host-level firewall, route tables, or VPN.
- Author the runner script that consumes the proxy.

## Universal preamble

1. Confirm the working directory is an agent-scaffold checkout
   (`docker/sandbox-proxy/filter`, `docker/sandbox-compose.yml` exist).
2. Read the current filter file to ground the change in current state.
3. Identify the operation:
   - **Add** — user wants `<host>` allowed.
   - **Remove** — user wants `<host>` removed.
   - **Audit** — user wants the existing list reviewed.

## Refusal rules

These refusals are non-negotiable. Stop and explain rather than write:

| Pattern | Refuse because |
|---|---|
| Bare wildcards (`.*`, `^.*$`, `^.*\.com$`) | Defeats default-deny. |
| Multi-label wildcards in subdomain position (`^.*\.example\.com$`) | Allows subdomain takeover to bypass the gate. Require an explicit subdomain list. |
| IP literals in the filter | Filter is hostname-only; IPs bypass DNS-based attribution. Use the scaffold's `tinyproxy.conf` Allow rules for network-level allow, not this filter. |
| Metadata-service hosts (`169.254.169.254`, `metadata.google.internal`, `instance-data.ec2.internal`) | Cloud-credential exfiltration vector. |
| Internal hosts (`*.internal`, `*.corp`, `*.local`, `*.cluster.local`) without team-lead context | Implement runner has no business hitting internal services. |
| Host-internal redirects (anything resolving to `127.0.0.0/8`, `10.0.0.0/8`, `172.16.0.0/12`, `192.168.0.0/16`) where the runner is configured to bypass the proxy | Bypass surface. |
| Hosts that do not exist (typos like `apii.openai.com`) | Surface and ask user to confirm. |

When refusing, name the rule, name the safer alternative, and stop.

## Add operation

1. **Validate the pattern.** It must match POSIX extended regex,
   anchored with `^...$`, escape literal dots (`\.`), and target one
   hostname or one explicit subdomain set joined with `|`.

   Good: `^api\.example\.com$`
   Good: `^(api|cdn)\.example\.com$`
   Bad:  `api.example.com` (unanchored, unescaped)
   Bad:  `^.*\.example\.com$` (multi-label wildcard)

2. **Apply refusal rules** above. If any matches, stop.

3. **Resolve the host** (sanity check it exists):
   ```bash
   getent hosts <host> || dig +short <host>
   ```
   If unresolvable, surface and ask the user to confirm before adding.

4. **Diff the change.** Show the user the before/after of the file with
   the new line in section context (Model APIs / Package managers /
   Public source hosts / Other).

5. **Write the change.** Append to the appropriate section. If no
   section fits, add a new section header (e.g., `# Partner banks`).

6. **Surface rebuild + verification.** Print the exact commands:
   ```bash
   docker compose -f docker/sandbox-compose.yml build proxy --no-cache
   docker compose -f docker/sandbox-compose.yml up -d proxy
   just doctor --full   # exercises the egress-test path
   ```

7. **Suggest a PR shape.**
   ```
   sandbox: allow <host> for <reason>
   ```

## Remove operation

1. **Find the line(s) matching the host or pattern.**
2. **Diff the removal.** Show the section the line lives in.
3. **Apply.** Delete the line.
4. **Surface rebuild + verification** (same commands as Add).
5. **Suggest a PR shape.**
   ```
   sandbox: remove <host> (no longer used by <stage>/<workflow>)
   ```

## Audit operation

1. Read `filter`. Group entries by their section.
2. For each entry, classify:
   - **Active** — referenced in the last 30 days (cross-check with
     proxy logs: `docker compose -f docker/sandbox-compose.yml logs proxy`).
   - **Stale** — no traffic in the last 30 days; candidate for removal.
   - **Suspect** — matches a refusal rule or has an unanchored / overly
     broad pattern.
3. Output an audit table to the user — do not edit during audit.
4. Recommend: which to remove (Remove operation), which to tighten
   (Add operation rewriting the line).

## Output format

Either an updated `docker/sandbox-proxy/filter` (for Add / Remove) or a
markdown audit table (for Audit). The skill never invokes
`docker compose build` — the user must type that.

## Constraints

- DO NOT permit wildcards in the subdomain position.
- DO NOT permit IPs, metadata hosts, or `*.internal`/`*.local` without
  team-lead context.
- DO NOT modify `tinyproxy.conf`, `Dockerfile`, or `sandbox-compose.yml`.
- DO NOT rebuild the proxy from this skill — surface the command.
- DO NOT skip the egress-test reminder. The whole point is verified
  enforcement.
- DO NOT add `claude` or `anthropic` brand names as comment-only entries
  — every entry must be a hostname, period.
- DO NOT batch unrelated changes; one PR per intent.

## Validation gate

Before writing:

1. Pattern is POSIX-anchored, has escaped dots, and matches no refusal
   rule.
2. Host resolves (or user explicitly accepted unresolvable as expected).
3. No duplicate of an existing entry.
4. Section assignment is correct.
5. Diff has been shown to the user and the user said "go".

## References

| Need | Reference |
|---|---|
| Pattern rules and the refusal table | `references/pattern-rules.md` |
| Why default-deny matters and how the proxy verifies it | `references/default-deny.md` |
| Audit cadence and proxy-log reading | `references/audit-cadence.md` |

## Templates

- `templates/allowlist-change.md` — PR description template for filter
  changes.
