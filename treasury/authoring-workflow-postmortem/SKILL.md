---
name: authoring-workflow-postmortem
description: >
  Author a postmortem for an agent-scaffold workflow within 48h of a trigger
  event: failed stage with more than $1 spent, rejected implement gate, cap
  overrun, or sandbox egress test failure. Reads .agent/runlog.jsonl,
  .agent/stages/*.md, .agent/approvals.md, and LiteLLM spend logs. Drafts at
  docs/postmortems/YYYY-MM-DD-WORKFLOWID.md and archives evidence in a
  sibling folder. Use when the user says "draft a postmortem", "postmortem
  for the failed workflow", "wasted spend writeup", "we exceeded the cap",
  "the sandbox failed", "the gate was rejected", or after the orchestrator
  flags a postmortem trigger. Does NOT execute remediation actions — those
  are tracked as PRs against prompts/library, profiles, or the playbook.
  Do NOT use for general retrospectives unrelated to a trigger event.
---

# Authoring a workflow postmortem

## Purpose

The PLAYBOOK §4 contract is: every trigger event gets a postmortem within
48 hours. This skill authors the document with the right shape, pulls
evidence from the scaffold's audit trails, archives the inputs alongside
the writeup, and ends with concrete action items that are checkable.

## When to use this skill

The user says "draft a postmortem" or any trigger event has fired:

| Trigger | Signal |
|---|---|
| Failed stage with > $1 spent | `state.json.stages.<x>.status == "failed"` AND `cost_usd > 1` |
| Rejected implement gate | `state.json.stages.implement.status == "rejected"` |
| Cap overrun | dispatcher exited rc=2 with `total_cost_usd >= spend_cap_usd` |
| Sandbox egress failure | `docker compose ... logs proxy` shows non-allowlisted host succeeded |

Do NOT use this skill to:
- Retro a successful workflow ("we should have done X better").
- Document an architecture decision (use ADRs).
- Replace a security incident report when a real exposure occurred —
  escalate first, postmortem second.

## Universal preamble

1. Confirm the working directory is an agent-scaffold checkout.
2. Identify the workflow_id and trigger event. Ask the user if not stated.
3. Confirm a trigger has actually fired by reading `state.json` and the
   relevant logs. If no trigger fired, refuse and explain.
4. Compute the postmortem path:
   `docs/postmortems/YYYY-MM-DD-<workflow_id>.md` using today's local date.
   Archive folder: `docs/postmortems/YYYY-MM-DD-<workflow_id>/`.
   Refuse to overwrite without explicit confirmation.

## Core workflow

### 1 — Collect evidence (read-only)

Pull these into the archive folder:

| Source | Archive path |
|---|---|
| `.agent/runlog.jsonl` | `<archive>/runlog.jsonl` |
| `.agent/runlog.md` (last 100 lines) | `<archive>/runlog.md` |
| `.agent/stages/*.md` | `<archive>/stages/` |
| `.agent/stages/<failed>.log` (redacted tail, 200 lines) | `<archive>/<failed>.log` |
| `.agent/approvals.md` excerpt for this workflow | `<archive>/approvals.md` |
| LiteLLM spend report `just llm-spend 24` | `<archive>/spend.txt` |
| `state.json` snapshot | `<archive>/state.json` |

Use `cp` for static files; redirect `just llm-spend 24` and the redacted
tail. Never edit the originals — archive copies only.

### 2 — Verify the audit chain

Before drafting, run:

```bash
just approvals-verify
```

If the chain fails, **stop**, escalate as a security incident, and do not
proceed with the postmortem until the chain is reconciled.

### 3 — Build the timeline

Generate from `runlog.jsonl`:

```bash
jq -r '"- \(.ts) — \(.msg)"' .agent/runlog.jsonl
```

Trim to events bound to this workflow_id. Keep timestamps in UTC.

### 4 — Identify the root cause

The actual reason — not "the model hallucinated". The right framing is
"what allowed the hallucination to get past critique/review without
being caught?" Use `references/root-cause-discipline.md`.

### 5 — Draft the postmortem

Use `templates/postmortem.md` (a thin extension of the scaffold's
`docs/postmortem-template.md` with this skill's archive-folder
convention). Sections: header, what-happened narrative, timeline, root
cause, contributing factors, what-we'll-change actions with owners and
dates, audit links.

### 6 — Action items

Each action lands as a PR target, not a good intention:

| Action class | Lands as |
|---|---|
| Bad prompt | PR against `prompts/library/<stage>/<topic>.md` (hand off to `drafting-stage-prompt`). |
| Wrong cap default | PR against `profiles/<name>.sh` (hand off to `authoring-scaffold-profile`). |
| Sandbox allowlist gap | PR against `docker/sandbox-proxy/filter` (hand off to `configuring-sandbox-allowlist`). |
| Gate-discipline failure | PR against `docs/PLAYBOOK.md` §10 to clarify the question that was missed. |
| Scaffold bug | PR against the relevant runner or `dispatch.sh`. |

Each action item has an owner (GitHub handle) and a due date within 14
days. Items without an owner do not ship.

### 7 — Surface to the user

Print the path of the drafted postmortem and the archive folder. Suggest
`git add docs/postmortems/<date>-<id>{,.md}` and a PR title shape:

```
postmortem: <workflow_id> — <one-line summary>
```

## Output format

- One markdown file at `docs/postmortems/YYYY-MM-DD-<workflow_id>.md`
  matching `templates/postmortem.md`.
- One archive folder at `docs/postmortems/YYYY-MM-DD-<workflow_id>/`
  with redacted evidence.

## Constraints

- DO NOT include real PII, real credentials, or real customer names.
  Logs are already redacted at write time, but cross-check before
  archiving.
- DO NOT speculate on intent — root cause names mechanism, not motive.
- DO NOT propose actions without owners and due dates.
- DO NOT modify originals under `.agent/`. Archive copies only.
- DO NOT skip the `just approvals-verify` step. A failing chain
  pre-empts the postmortem.
- DO NOT write a postmortem when no trigger fired. Refuse and explain.
- DO NOT execute remediation from this skill — author the actions, hand
  off to the relevant sibling skill.

## Validation gate

Before the postmortem ships:

1. Trigger event verified against `state.json` or logs.
2. Approval chain verified (`just approvals-verify` exit 0).
3. Archive folder exists with all listed files.
4. Every action item has an owner and a due date ≤ 14 days out.
5. Header severity matches the trigger class.
6. No real PII or credentials in the markdown body or archive.

## References

| Need | Reference |
|---|---|
| Trigger conditions and how to verify each | `references/trigger-conditions.md` |
| Evidence collection commands and ordering | `references/evidence-collection.md` |
| Root-cause discipline (mechanism, not motive) | `references/root-cause-discipline.md` |
| Action-item routing to sibling skills | `references/action-routing.md` |

## Templates

- `templates/postmortem.md` — full document skeleton.
