# Policy-as-code model: three layers, OPA/Rego, default-deny primacy

Deep-dive for Workflow steps 3-4. Do not restate SKILL.md here.

## The three-layer enforcement architecture

Command restrictions in AI coding tools are enforced across three nested
layers, each with different strength and cost characteristics. The Rego
allowlist this skill generates is a **Layer 2** artifact; it is meaningless
without the Layer 3 floor beneath it.

```
Layer 1: Instruction-Level     CLAUDE.md / system prompt rules   weakest / cheapest
Layer 2: Configuration-Level   settings deny rules, PreToolUse
                               hooks, MCP allowlists, OPA/Rego
Layer 3: OS / Sandbox-Level    kernel namespaces, bubblewrap,
                               seccomp, network isolation        strongest / costly
```

- **Layer 1 — instruction-level rules** are embedded in documentation files
  (CLAUDE.md, system prompts, rules files). Cheapest to set up, easiest to
  update, and the weakest: they rely on the model to read and follow the
  instruction, so they are vulnerable to prompt injection and context
  manipulation. Never treat Layer 1 as a security boundary.
- **Layer 2 — configuration-level enforcement** operates through the agent's
  own configuration system, independently of what the model reasons. Deny
  rules take precedence over ask and allow. Enforcement is deterministic and
  configuration-driven, "but it can be circumvented by encoding or path
  manipulation if pattern matching is naïve."
- **Layer 3 — OS/Sandbox-level enforcement** is the strongest tier:
  kernel-level mechanisms prevent the agent process from accessing blocked
  paths or network destinations "regardless of what the model or
  configuration layer says." Anthropic's internal sandboxing for Claude Code
  "reduced permission prompts by 84% while maintaining safety, demonstrating
  that structural boundaries outperform per-command lists at scale." Hard
  boundaries enforced at this layer "cannot be circumvented by any prompt,
  jailbreak, or reasoning trick."

A practical tool gateway combines all three: instructions define intent,
configuration-layer deny rules provide the operational blacklist, and
OS-level sandboxing provides the non-negotiable floor.

## Why OPA / Rego

"Open Policy Agent (OPA) with Rego policies provides a language-agnostic
enforcement layer that can sit across Kubernetes, CI/CD pipelines, API
gateways, and MCP servers simultaneously." For multi-squad organizations,
"adopt OPA or an equivalent policy engine at the tool gateway layer, with
deny-by-default configuration and audit logging for every tool call." A
stateless policy engine of this kind can intercept every agent action before
execution at sub-millisecond latency, and on deviation can soft-block (action
proceeds but trust score drops), hard-block (action denied), or log-only.

## Default-deny is primary; blacklist is supplementary

"The correct architecture for a software-delivery squad is layered:
structural OS-level sandboxing as the floor, a default-deny allowlist as the
primary access model, and a blacklist as a supplementary emergency brake for
high-velocity development environments where allowlist expansion cannot keep
pace." The blacklist's role "is to catch well-understood danger categories
(destructive filesystem operations, outbound network calls, privilege
escalation) in a second layer of defense."

Blacklisting is *threat-centric* (enumerate and block what is dangerous);
allowlisting is *trust-centric* (enumerate and permit only what is safe).
"The 2025 security consensus recommends allowlisting as the baseline, with
blacklisting as a supplementary emergency brake."

## Allowlist vs. blacklist trade-off table (verbatim)

| Dimension | Blacklist | Allowlist |
|---|---|---|
| Default posture | Permit (open) | Deny (closed) |
| Failure mode | New dangerous commands are permitted until listed | New needed commands are blocked until approved |
| Maintenance burden | Grows with threat landscape | Grows with feature surface |
| Onboarding friction | Low (most commands work immediately) | High (must pre-approve every tool) |
| Security posture | Threat-centric; lagging by design | Trust-centric; safer by design |
| Best fit | Supplementary emergency brake; incident response | Primary access model for production agents |

Three trade-offs deserve attention: coverage vs. ease of use (allowlisting
bounds the attack surface but a minimal read-only start blocks legitimate
tasks like `npm install`, `git push`, `docker build`); maintenance discipline
(blacklist expansion is opt-out and "always one step behind adversarial
creativity"; allowlist expansion is opt-in); and **false confidence** — "the
most significant operational risk of a blacklist-only policy is that it
creates a sense of control that the bypass evidence does not support."

## Why model self-policing is not enough

"A model trained to refuse harmful commands may still comply under
adversarial prompting conditions, but a configuration-layer deny rule will
block execution regardless of what the model was told to do." The four
components shaping agent behavior are the model, the harness, the available
tools, and the operating environment; "a well-trained model alone cannot
secure agentic AI if the other three layers are misconfigured."
