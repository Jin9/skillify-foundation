# AGENTS.md curation checklist

This checklist accompanies a candidate AGENTS.md. The candidate is
pre-review and NOT commit-ready: an unreviewed LLM-generated AGENTS.md
reduces task success ~3%, while a minimal human-curated one gives ~+4%.
A human MUST complete this pass before the file is committed.

## Human editorial pass (required before commit)

- [ ] File is under 200 lines; ideally 100-150. (`scripts/agents_md_gate.py` run, no FAIL.)
- [ ] Exactly the six sections present: commands, testing, project
      structure, code style, git workflow, boundaries. No seventh section
      except the optional Squad Roles block.
- [ ] Every command is the exact CLI with real flags, copy-pasteable, and
      currently correct against the toolchain.
- [ ] Project structure is 3-5 orientation bullets — not a file tree, not
      architecture history.
- [ ] Code style is one snippet plus a pointer to the lint config; no lint
      rules restated in prose.
- [ ] Boundaries is exactly three tiers (always-allowed /
      requires-human-approval / explicitly-prohibited).
- [ ] Every prohibition is paired with a concrete allowed alternative
      (a "do" for every "don't"). No bare "never X." lines.
- [ ] Deferral filter applied: no rule that a lint rule, CI check, or
      branch-protection rule could enforce remains in the file. Each such
      rule is recorded in the deferral log below.
- [ ] Squad Roles block included only if multi-agent; handoff *schema* is
      pointed to, not defined.
- [ ] All comment/instruction blocks from the skeleton deleted.
- [ ] No secrets, no stale commands, no vague quality statements
      ("write clean code").
- [ ] Reviewer name and date recorded: ______________  ____-__-__

## Moved to tooling instead (deferral log)

Every rule dropped from AGENTS.md under the deferral filter is logged here
with the mechanism that should enforce it. The owning human wires these up.

| Dropped rule (verbatim) | Why deferred | Enforce in | Owner / skill | Wired up? |
|---|---|---|---|---|
| e.g. "use 4-space indents" | formatting | lint config (`pyproject.toml`) | repo owner | [ ] |
| e.g. "tests must pass before merge" | testable gate | CI required check | repo owner | [ ] |
| e.g. "no direct pushes to main" | access control | branch protection | repo admin | [ ] |
| e.g. "auto-approve `npm test` only" | command policy | command policy | governance-policy-generator | [ ] |
| e.g. "rotate the API token at runtime" | secret hardening | runtime config | devops-infrastructure-hardener | [ ] |
| e.g. "handoff payload must include trace id" | handoff schema | handoff contract | multi-agent-handoff-architect | [ ] |
|  |  |  |  | [ ] |

## Maintenance note

AGENTS.md is institutional memory, not a one-time artifact. Stale content
actively poisons every request's context. Update the file when an agent
makes a mistake it could have prevented; trim sections the moment they no
longer apply.
