# Leakage patterns catalog

Five pathways through which secrets escape an agent workflow. Every step-1
finding is classified into exactly one of these. Detection signatures below
mirror what `scripts/scan_static_creds.py` flags; the script is a first pass,
not a substitute for reading the flagged lines.

## 1. Debug logging and stdout injection — the dominant vector

The largest single source of agent credential leakage is the humble print
statement. Debug logging accounts for 73.5% of agent skill vulnerabilities:
developers add `print(api_key)` or `console.log(secret)` statements during
development, the framework captures stdout and feeds it back into the LLM
context window, and the secret is now visible to the model and any
observability tool that records the context. Seven percent of studied skills
exposed credentials through the LLM context window or output logs.

Detection signatures (example patterns):

- Python: `print(...)`, `logging.*(...)`, `logger.*(...)` whose argument
  references a secret-named variable — e.g. `print(api_key)`,
  `logger.debug(f"token={token}")`.
- JS/TS: `console.log(...)` / `console.error(...)` referencing
  `secret`, `token`, `apiKey`, `api_key`.
- Any log/print line whose interpolated expression name contains
  `secret`, `token`, `api_key`, `apikey`, `password`, `passwd`, `private_key`.

Remediation layer: Layer 3 (output scrubber) plus removal of the statement;
the underlying credential should already be Layer-1 dynamic.

## 2. Environment-variable inheritance

Injecting secrets as environment variables at startup means the agent
inherits all credentials from its environment and may — through logging, tool
outputs, or prompt injection — surface any of them to a potentially hostile
context. Traditional apps handle secrets at build/deploy time with a human
present; an agent calling a third-party API mid-task cannot pause for
approval and holds the credential for the task's duration.

Detection signatures:

- `os.environ[...]` / `os.getenv(...)` / `process.env.*` reads of
  secret-named keys feeding a request, especially with no secrets-manager
  client in the call path.
- Static secret values assigned to env keys in shell, Dockerfile `ENV`,
  compose, or CI config (`API_KEY=...`, `export TOKEN=...`).

Remediation layer: Layer 1 — replace with CLI-injected dynamic credential or
workload identity federation.

## 3. Hard-coded keys

A literal credential committed in source or config.

Detection signatures:

- AWS access key id `AKIA` followed by 16 uppercase/digit chars.
- Private key headers: `-----BEGIN ... PRIVATE KEY-----`.
- Assignment of a high-entropy / structured literal to a name containing
  `secret`, `token`, `api_key`, `apikey`, `password`, `passwd`,
  `client_secret`, `private_key` — e.g. `api_key = "AKIA..."`.

Remediation layer: Layer 1 (eliminate) + Layer 4 (revoke and rotate the
exposed key immediately; assume compromised).

## 4. MCP-config secrets

MCP server configurations frequently store tokens and credentials in
environment variables or config files. 24,008 unique secrets were found in
MCP configuration files alone. A single compromised MCP server can access
credentials for all connected agents; the protocol's lack of built-in
host/server authentication makes tool-redefinition attacks structurally
possible in multi-server deployments.

Detection signatures:

- Secret-named keys with literal values inside MCP config / manifest files
  (e.g. `mcp.json`, server `env` blocks, `.mcp/*`).
- Shared credential reused across multiple MCP server entries.

Remediation layer: Layer 1 (dynamic per-server credential) + Layer 2
(configuration-file protection so the MCP config cannot be silently mutated)
+ unique identity per server (Layer 4).

## 5. Prompt-injection exfiltration

Prompt injection is the OWASP #1 threat for LLM applications. In indirect
injection an attacker embeds instructions in external content the agent
processes (PR comment, web page, document, DB record); the agent acts on
instructions it cannot reliably distinguish from authorized ones (the
confused-deputy problem) and can be made to post or send its credentials.
Demonstrated: a malicious PR title caused a coding agent to post its own API
key as a PR comment. Frameworks feed tool outputs back into the LLM context
without sanitization by default.

Detection signatures (design-level, not a regex):

- Tool outputs / external content fed back into the prompt with no scrubbing
  between the tool boundary and the model.
- No network egress allowlist on the agent runtime.
- Credentials reachable from the agent's context at all (i.e. Layer 1 not in
  place).

Remediation layer: Layer 1 (no static secret to exfiltrate) + Layer 2
(egress allowlist) + Layer 3 (scrub the output boundary).

## Important caveat

No major agent framework activates security controls by default. An empty
scanner result is not a clean bill of health — patterns 4 and 5 are largely
design-level and will not be caught by literal scanning. Always pair the
scan with a manual review of agent/MCP config and the prompt data flow.
