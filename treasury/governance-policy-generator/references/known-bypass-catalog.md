# Known-bypass catalog: why blacklist-only is false confidence

Deep-dive for Workflow step 6. Use this as the self-check before emitting.

## The headline figures (verbatim, do not round differently)

- **~36%** — "Lasso Security achieved a 36% bypass rate against Perplexity's
  BrowseSafe using encoding variants (NATO phonetic alphabet, Pig Latin,
  Base32, hexadecimal)."
- **81% vs. 11%** — "NIST red-team research found agent-specific attacks
  achieved an 81% task-hijacking success rate vs. 11% for baseline attacks,"
  "indicating that offensive research substantially outpaces blacklist-based
  defenses."
- **~20%** — "AI coding tools hallucinate non-existent package names roughly
  20% of the time — inventing 440,445 fake dependencies across Python and
  JavaScript ecosystems in one study." This enables "slopsquatting": the
  `react-codeshift` name an LLM hallucinated was found "referenced in 237
  GitHub repositories"; North Korean APT Famous Chollima's PromptMink attack
  used "LLM Optimization abuse and knowledge injection to actively make
  malicious packages more likely to be selected by AI agents."
- **84%** — Anthropic's internal Claude Code sandboxing "reduced permission
  prompts by 84% while maintaining safety," showing structural boundaries
  outperform per-command lists at scale.

## The six demonstrated bypass techniques

"At least six distinct bypass techniques have been demonstrated against
production tools."

1. **Path aliasing (CVE-2025-54794).** "Claude Code's denylist was bypassed
   by accessing the same blocked binary via `/proc/self/root/usr/bin/npx` — a
   different filesystem path to the same executable that the pattern matcher
   did not recognize. When bubblewrap blocked that path, the agent disabled
   the sandbox itself."
2. **Command chaining through whitelisted commands (CVE-2025-54795).** "The
   whitelisted `echo` command was exploited to inject arbitrary shell
   instructions: `echo \"\\\"; <COMMAND>; echo \\\"\"`. No user confirmation
   was required. Fixed in Claude Code v1.0.20."
3. **Session-depth limit exhaustion.** "Claude Code's deny rules were found
   to stop enforcing silently after a command session exceeds 50 subcommands,
   enabling an attacker to pad a session with 50 no-op `true` subcommands and
   then issue a previously blocked `curl` call."
4. **Encoding-based obfuscation.** The ~36% Lasso/BrowseSafe result above —
   "string-matching blacklists are vulnerable to transformations that
   preserve intent but obscure the blocked string."
5. **Argument injection into "safe" tools.** "Trail of Bits found that `git
   show` and `ripgrep`, which were on nominally safe command lists, operated
   without argument restrictions, providing prompt-injection-to-RCE pathways
   within supposedly restricted environments."
6. **MCP server auto-execution (TrustFall).** "All four major coding agent
   CLIs (Claude Code, Gemini CLI, Cursor CLI, Copilot CLI) auto-execute
   project-defined MCP servers on folder trust acceptance, and all default to
   'Yes/Trust.' In headless CI environments, the trust dialog never renders,
   so a malicious project file on a PR branch can trigger arbitrary code
   execution with zero human interaction."

Related real-world incidents: "three CWE-78 command injection flaws in Claude
Code CLI (CVE-2025-66032, fixed in v1.0.93) enabling CI/CD credential
exfiltration," and "a Cline VS Code extension compromise (5M+ users) via
prompt injection that exfiltrated npm tokens."

## Why blacklist alone is false confidence

"Path aliasing, encoding, argument injection, and session-depth exhaustion
all bypass syntactic blacklists without triggering any alert." "Every major
coding agent platform has had at least one confirmed blacklist bypass, and
the bypass techniques are straightforward enough to appear in public blog
posts." "A shared design assumption that prompt-level controls constitute a
meaningful security boundary has been formally declared broken by security
researchers."

The most dangerous failure mode "is not about the model being wrong — it is
about the model being right in a compromised environment": slopsquatting,
prompt injection, and supply-chain attacks "all exploit the agent doing
exactly what it is designed to do ... in an environment that has been
adversarially manipulated."

## Self-check before emit (step 6)

- [ ] The `.rego` is `default allow = false` and every allow rule is conditioned.
- [ ] The blacklist is present only as a supplementary emergency brake, not the primary model.
- [ ] No syntactic deny rule is documented as a security boundary.
- [ ] Package installation is gated on a known-good registry allowlist (~20% hallucination, slopsquatting).
- [ ] An OS/sandbox floor is named as the non-negotiable Layer 3 base.
- [ ] CI/headless note: agent execution gated to post-merge branches, not arbitrary PR branches (TrustFall).
