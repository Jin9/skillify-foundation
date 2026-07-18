---
routine: REPLACE-with-kebab-name-matching-filename
summary: REPLACE - one sentence, max 200 chars, no angle brackets.
version: 0.1.0
default_on_fail: stop
owner: REPLACE-or-delete-this-line
tags: REPLACE-or-delete-this-line
requires: REPLACE-external-preconditions-or-delete-this-line
---

Optional prose: what this routine is for and when to launch it. The parser ignores everything outside the two sections below. Delete every line containing REPLACE, then validate with scripts/validate_routine.py before registering in INDEX.md.

## Inputs

- input_key: what the user must supply (required)
- other_key: what this controls (optional, default VALUE)

## Nodes

### Node: first-step

- purpose: REPLACE - one line, what this node produces and why
- executor: REPLACE-skill-name-or-agent-inline
- tier: small
- inputs: user.input_key
- outputs: 01-first-step.md
- gate: after

### Node: second-step

- purpose: REPLACE
- executor: agent-inline
- executor_mode: REPLACE-or-delete-this-line
- tier: mid
- inputs: 01-first-step.md, user.other_key
- outputs: 02-second-step.md
- gate: none
- on_fail: stop
- notes: REPLACE-caveats-or-delete-this-line
