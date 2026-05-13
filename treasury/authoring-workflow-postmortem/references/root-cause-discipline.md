# Root-cause discipline

A postmortem's value is in identifying the **mechanism** that allowed the
failure, not the **agent** that failed. "The model hallucinated" is not a
root cause. "Critique didn't ask for source-level citations, so the
hallucination wasn't surfaced" is.

## The 5-whys filter

Apply at least three `why`s before declaring root cause.

```
Failure: implement stage produced a function that calls a non-existent API.
  Why 1: The plan referenced the API by name without verifying it existed.
  Why 2: The plan stage prompt didn't require source-level citations.
  Why 3: The library prompt for plan was written before the project
         enforced citation discipline.
Root cause: prompts/library/plan/<topic>.md predates the citation rule.
Action:     update the library entry; add a CI check that plan output
            references at least one file:line per claim.
```

## Mechanism categories

When you can't see a clear mechanism, the failure usually fits one of:

| Category | Example mechanism |
|---|---|
| Prompt drift | The prompt asks for X but produces Y under longer goals. |
| Model mismatch | A prompt validated on Codex, run unchanged on Gemini. |
| Gate discipline | An approver skipped Question 2 (critique cross-walk). |
| Cap math | Cap was set assuming $0.30/stage; actual was $0.80. |
| Profile config | `IMPLEMENT_SANDBOXED` defaulted off in the profile. |
| Scaffold bug | Stale-PID recovery didn't fire because pid was reused by an unrelated process. |
| Repo state | Uncommitted changes on the branch produced confusing diffs. |
| Audit-trail gap | Approval log entry was malformed; verify chain failed. |

## What is NOT a root cause

- "The model hallucinated."
- "We didn't have time to review."
- "It was a Friday afternoon."
- "The researcher was new."

These are contributing factors at best. The root cause is whatever the
**system** was missing that allowed those circumstances to produce the
failure.

## Surfacing the root cause

State it in **one sentence** that names a mechanism and a system gap.
Examples:

- "Plan-stage library prompt predates the citation rule, so hallucinated
  APIs survive into implement."
- "Cap-tier $5 default in `profiles/code.sh` doesn't account for the
  doubled token cost on credit-decision goals."
- "Sandbox enforcement check runs before, not after, each implement
  attempt — drift between attempts went undetected."

If you cannot state the root cause in one sentence, you have not yet
found it. Ask another `why`.

## Contributing factors vs. root cause

Contributing factors *enabled* the failure but did not *cause* it.
Examples:

- Researcher approved without re-reading critique on a phone notification.
- `NTFY_TOPIC` was misconfigured, so the researcher noticed late.
- The cap was set at the team-lead-acknowledgment ceiling.

These belong in the contributing-factors section. They become root
causes only if the failure happened *because* of them, not merely *during*
them.
