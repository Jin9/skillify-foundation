---
name: agent-context-initializer
description: >
  Generate a minimal, human-curatable AGENTS.md (under 200 lines, target
  100-150) for an agentic-squad repo: exactly six sections (commands,
  testing, project structure, code style, git workflow, boundaries), every
  prohibition paired with an allowed alternative, a three-tier boundaries
  block (always-allowed / requires-human-approval / explicitly-prohibited),
  and an optional squad-roles block. Any rule enforceable in a lint, CI, or
  branch-protection rule is dropped and logged to a curation checklist
  instead. Use when the user asks to "create an AGENTS.md", "scaffold agent
  instructions for this repo", "write a repo context file for the squad",
  "bootstrap AGENTS.md", or "generate an agent constitution". Output: a
  candidate AGENTS.md plus a curation checklist. Do NOT use for
  machine-enforced command / auto-approve policy
  (governance-policy-generator), runtime secrets or identity hardening
  (devops-infrastructure-hardener), inter-agent handoff schemas
  (multi-agent-handoff-architect), or editing this skill's own files.
---

# Agent Context Initializer

## Purpose

Produce the smallest AGENTS.md that still changes agent behavior usefully:
the six highest-performing content areas, every "don't" paired with a "do",
a three-tier boundaries block, and an optional squad-roles contract. The
file is a control surface, not documentation — anything enforceable by
tooling is deliberately excluded and routed to a curation checklist for the
human to wire into lint, CI, or branch protection. This skill *generates* an
AGENTS.md when invoked: that is in scope. Generating an AGENTS.md while
authoring or editing this skill itself is explicitly out of scope — only
`templates/AGENTS.md` (a skeleton) lives in the skill.

## When to use this skill

- Use when: "create an AGENTS.md for this repo" / "scaffold agent instructions for the squad".
- Use when: "bootstrap a repo context file" / "generate an agent constitution / squad operating brief".
- Use when: a repo has agents working in it with no AGENTS.md, or an existing one over 200 lines needs trimming to the six sections.
- Do NOT use when: the ask is machine-enforced command / auto-approve policy (`governance-policy-generator`), runtime secrets or identity hardening (`devops-infrastructure-hardener`), inter-agent handoff schema design (`multi-agent-handoff-architect`), or editing this skill's own files. Hand those to the appropriate skill.

## Inputs

Read access to the target repository: the build/test/lint tooling, the test
layout, the top-level module structure, the lint/format config files, and
the git/PR conventions. Optionally: the squad roster (agent role names and
the paths each owns). None of this is generated here — only read to extract
the six sections. The skill never edits repository files; it emits a
candidate plus a checklist for a human pass.

## Workflow

1. **Scope and read, never write.** Confirm the target repo path and that the deliverable is a *candidate* AGENTS.md plus a curation checklist — not a committed file. Inventory the six sources: exact build/test/lint commands, test layout, top-level structure, lint/format config paths, git/PR conventions, and (if a squad) the role roster. Entry: a repo with agents and no minimal AGENTS.md. Exit: the six raw inputs collected.
2. **Apply the deferral filter first.** For every candidate rule, ask whether it can be expressed in a lint rule, a CI check, or a branch-protection rule. If yes, DROP it from AGENTS.md and log it to the curation checklist's deferral table instead — the tool is the constraint, not prose. Keep only what cannot be enforced programmatically. See `references/deferral-meta-principle.md`.
3. **Draft the six sections, lean.** Fill exactly the six sections (commands, testing, project structure, code style, git workflow, boundaries) from `templates/AGENTS.md`. Commands carry exact flags; project structure is 3-5 bullets; code style is one snippet plus a pointer to the lint config, never restated rules. No seventh section except the optional squad-roles block. See `references/six-section-spec.md`.
4. **Build the three-tier boundaries (and optional squad-roles) block.** Write boundaries as three explicit tiers: always-allowed, requires-human-approval, explicitly-prohibited. For a multi-agent repo, add the squad-roles block (role, owned paths, handoff-target, human-escalation triggers) — point to the handoff schema, do not define it here. See `references/boundaries-and-squad-roles.md`.
5. **Pair every don't with a do; enforce the budget.** Walk every prohibition: each must carry a concrete allowed alternative on the same line or the adjacent bullet ("don't push to main — open a PR"). Cut connective prose until the file is under 200 lines, target 100-150; if it cannot fit, more rules belong in the deferral table, not the file. See `references/authorship-and-budget.md`.
6. **Self-check, then emit for a human pass.** Run `python3 scripts/agents_md_gate.py <candidate>`: it FAILs on over 200 lines, a missing required section, or an unpaired prohibition, and WARNs over 150 lines. Emit the candidate AGENTS.md plus the completed curation checklist. State that a human editorial pass is mandatory before commit — an unreviewed LLM-generated file degrades performance. Do not present the candidate as commit-ready.

Do not let the file accrete: a stale rule actively poisons every request's
context. Treat AGENTS.md as institutional memory the human curates, not a
one-time dump this skill finalizes.

## Output contract

Two markdown artifacts (shaped by the `templates/`):

- **Candidate AGENTS.md** — exactly the six sections plus the three-tier boundaries block and an optional squad-roles block; under 200 lines (target 100-150); every prohibition paired with an allowed alternative; no rule enforceable by lint/CI/branch-protection.
- **Curation checklist** — the human editorial-pass checklist plus a "moved to tooling instead" deferral log table recording every rule dropped under the deferral filter and where it should be enforced.

No repository files are edited, committed, or pushed. No AGENTS.md is
written outside the named candidate artifact. The candidate is explicitly
pre-review, never commit-ready.

## Constraints

- DO NOT emit an AGENTS.md exceeding 200 lines; target 100-150. If rules do not fit, defer more — do not raise the ceiling.
- DO NOT leave any prohibition without a concrete allowed alternative (a "do" for every "don't").
- DO NOT add a section beyond the six (commands, testing, project structure, code style, git workflow, boundaries) except the optional squad-roles block.
- DO NOT restate any rule expressible as a lint rule, CI check, or branch-protection rule; log it to the curation checklist deferral table.
- DO NOT define the inter-agent handoff schema, machine-enforced auto-approve policy, or runtime secret hardening here — point to the sibling skills.
- DO NOT present the candidate as commit-ready; a human editorial pass before commit is mandatory.
- DO NOT edit, commit, or generate an AGENTS.md anywhere except the named candidate artifact; never touch repository files or this skill's own files.
- DO NOT duplicate the methodology here; it lives one level deep in `references/`.

## Validation

- [ ] Frontmatter `name` equals folder, kebab-case, no XML, description under 1024 chars with triggers + negatives and sibling cross-refs.
- [ ] Workflow applies the deferral filter before drafting, pairs every don't with a do, enforces the 200-line cap, and ends with the `agents_md_gate.py` self-check plus a mandatory human pass.
- [ ] Output contract names both artifacts (candidate AGENTS.md, curation checklist) and the no-edit / pre-review boundary.

## References

- What each of the six sections must / must not contain: `references/six-section-spec.md`
- Three-tier boundaries model + squad-roles contract format: `references/boundaries-and-squad-roles.md`
- The lint/CI/branch-protection deferral meta-principle and failure modes: `references/deferral-meta-principle.md`
- Human-curated vs LLM-generated evidence, efficiency figures, line budget: `references/authorship-and-budget.md`
- Skeletons: `templates/AGENTS.md`, `templates/curation-checklist.md`; deterministic gate: `scripts/agents_md_gate.py`
