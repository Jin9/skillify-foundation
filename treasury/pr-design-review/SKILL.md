---
name: pr-design-review
description: >
  Review a PR's design and maintainability — business-logic completeness, test coverage, over-engineering, and template conformance — returning tagged, teachable comments plus a merge-or-iterate verdict. Use when the user asks to "review the design of this PR", "give me must-change / consider / nit feedback", "is this PR ready to merge", or "review the design of this change". Reviews in priority order, tags every comment, and stops at a human approval gate without auto-merging. Do NOT use for diff-correctness, compile, or logic-bug review — that is the built-in /code-review.
---

# pr-design-review

## Purpose
Review a PR's **design and maintainability** — for **business-logic completeness, test coverage, over-engineering, and
template conformance** — and return tagged, teachable comments plus a merge/iterate verdict.

## When to use
- Triggers: *"review the design of this PR"*, *"give me must-change / consider / nit feedback"*, *"is this PR ready to merge"*, *"review the design of this change"*.
- **Not this skill (important negatives):** diff-correctness / compile / logic bugs → built-in **`/code-review`**; security (BOLA, injection, secrets, RBAC) → **`reviewing-software-security`**; localizing a specific bug → **`progressive-bug-hunter`**; opening the PR → **`publishing-git-review-requests`**.

## Input
- A PR diff / changeset **plus** the requirement it implements.

## Output
A **Code Review** artifact + checklist + verdict. Skeleton:

```
# Review — PR <id>
Priority pass:
  1. Business logic complete?   2. Unit tests cover all business logic?
  3. Over-engineered?           4. Matches standard template?
Comments:
  [MUST-CHANGE] <file:line> — <why>
  [CONSIDER]    <file:line> — <why>
  [NIT]         <file:line> — <why>
Verdict:  MERGE  |  ITERATE
```

## Decision rules
1. **Review in priority order**: (1) business-logic completeness, (2) unit-test coverage of all business logic, (3) API/process over-engineering, (4) conformance to the standard template.
2. **Tag every comment**: *must-change* (must fix) · *consider* (fix if time allows) · *nit* (let it pass).
3. **Business-logic coverage = 100%**, and probe edge-case parameters to cut risk.
4. **PR > 400 lines → reject and split** — prefer split by interface change → draft handler → business logic.
5. Handle author disagreement by use-case + a worst-case example; if their argument is stronger, **concede — they may be right**.
6. Don't nitpick naming when there's no time; rarely comment on style alone.

## Checklist
- [ ] Business logic complete for the requirement
- [ ] 100% business-logic test coverage + edge cases
- [ ] Not over-engineered (the "5-minutes-to-explain" test)
- [ ] Matches the standard template
- [ ] Every comment tagged must-change / consider / nit
- [ ] PR size acceptable (else split)
- [ ] Verdict recorded with reasons

## Anti-patterns (never do)
- Approve a PR without a real review.
- Merge code without tests / sign off on missing test coverage.
- Reject a PR without explaining the reasons.
- Skip review of your own PRs.
- Talk down to the author, or fix a junior's bug instead of teaching.

## Example
**Input:** PR implementing wallet top-up. **Output (excerpt):** `[MUST-CHANGE]` missing duplicate-idempotency-key test;
`[CONSIDER]` extract amount validation into the domain; `[NIT]` rename `amt` → `amount`. **Verdict: ITERATE.**

## Human approval gate
**Stop.** The author fixes and re-submits; the reviewer does **not** auto-approve or merge on the team's behalf. The
skill produces review comments + a recommendation only.
