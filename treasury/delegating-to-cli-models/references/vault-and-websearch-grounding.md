# Vault and websearch grounding

The flagship research-digest mode dispatches agy/"Gemini 3.1 Pro (High)" to digest a research vault
plus current web sources. This file covers reading the vault safely and keeping the digest honest.

## The vault is read via its local mirror, not iCloud

- The Obsidian ResearchVault lives in iCloud. **macOS TCC blocks spawned processes from reading
  iCloud-backed paths** — a CLI the agent launches cannot read the iCloud original even though an
  interactive app could.
- Route all spawned reads to the **local mirror** instead. On this setup that is
  `~/Desktop/project/research/literature/` (deep-research report folders; each report's canonical
  file is its `NN-final_report.md`). Confirm the mirror path at runtime rather than assuming it.
- Pass the CLI **explicit absolute paths** to the report files it should read. Do not ask it to
  "find the vault" — that invites a TCC failure and agentic wandering.

## Scoping the agy read

- `agy --print` is agentic. A broad prompt makes it scan the workspace (shell history, dotfiles,
  the repo). Counter it: list the exact files to read, forbid reading anything else, demand a fixed
  output shape, and run from a neutral working directory (the wrapper's 3rd arg).
- **Lethal trifecta.** An agent with access to private data, exposure to untrusted content, and
  network egress can be steered into exfiltrating secrets. The vault + websearch + an agentic CLI is
  exactly that shape — keep secrets out of the prompt and the working directory.

## Websearch and citation discipline

- The CLI's web grounding produces leads, not proof. **Verify before relying:** a real, resolvable
  URL and a claim you can corroborate, or it stays a lead.
- Be skeptical of bare-domain citations (`medium.com`, `example.ai`) and of repos or sources you
  cannot open — they are commonly confabulated. Down-rank them; do not present them as the skill's
  evidence.
- Prefer the vault's primary reports (which the agent can read directly) as the authoritative layer,
  and use the web pass for currency. Cross-check the digest's claims against the actual report files
  before accepting. See `adjudication-and-provenance.md`.
