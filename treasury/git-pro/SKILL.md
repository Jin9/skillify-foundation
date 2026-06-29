---
name: git-pro
description: >
  Manage git end-to-end like an expert: branch and stage cleanly, write atomic
  conventional-commit messages, choose merge vs rebase, resolve conflicts, tidy
  UNPUBLISHED history (squash, reorder, amend), and undo mistakes safely with
  revert, reset, or reflog; then derive a simple best-practice MR/PR description,
  push, and open the PR or MR. Confirms before mutating commands and prefers
  recoverable operations. Use when the user asks to manage git, clean up or squash
  or rebase commits, resolve a merge conflict, undo a commit safely, write a simple
  MR or PR description, or take changes all the way to a PR or MR. Do NOT use to
  only push one already-prepared commit with the standard template (use
  publishing-git-review-requests); to run a Jira issue through a fix workflow (use
  jira-fix-mr-workflow); to localize or diagnose a bug (use progressive-bug-hunter);
  to design a CI/CD or release pipeline (use agentic-workflow-design); or to merge
  or approve a PR or MR for someone else.
---

# Git Pro

## Purpose
Act as an expert git operator: run the local git workflow the right way — clean branches, atomic conventional commits, deliberate merge/rebase, careful conflict resolution, honest history, and safe undo — then produce a simple best-practice MR/PR description and open the review request. One skill, from working tree to opened PR/MR.

## The cognitive kernel
Boot this small, fixed core on every task; load everything else on demand.

1. **Safety first.** Confirm before any mutating command; prefer recoverable operations (revert, reset --soft, reflog) over destructive ones; never rewrite published history or force-push a shared branch without explicit approval. → `references/git-safety-policy.md`.
2. **Atomic, conventional commits.** One logical change per commit; a clear `type(scope): summary` subject and a body that says *why*. → `references/commit-conventions.md`.
3. **Honest history.** Tidy only *unpublished* commits; keep the shared record truthful and recoverable. → `references/history-and-undo.md`.
4. **Simple but complete description.** The MR/PR description states what, why, how it was tested, and the rollback — and nothing more. → `references/mr-description-method.md`.

Everything deeper is *userland*: read the matching reference only when a step needs it. The kernel stays small and fixed; modules load when needed.

## Safety and confirmation boundary
Run this before executing anything that changes the repo.

1. **Inspect first (read-only, no confirmation):** `git status --short`, `git branch --show-current`, `git branch -vv`, `git log --oneline -n 10`, `git remote -v`. Establish branch, upstream, provider, and what is already published.
2. **Classify the command** with `references/git-safety-policy.md`:
   - **Allow** — read-only inspection runs without asking.
   - **Confirm** — any mutation (add, commit, branch create/switch, merge, rebase, cherry-pick, stash drop, push, PR/MR create). State exactly what will happen, then wait.
   - **Deny by default** — destructive or history-rewriting commands on *published/shared* refs (`push --force`, `reset --hard`, `clean`, branch deletion, amend/rebase of pushed commits). Proceed only on explicit user confirmation, and prefer `--force-with-lease` over `--force`.
3. **Preserve unrelated work:** stage explicit paths, never blanket `git add -A`; leave unrelated dirty files unstaged and report them.
4. If a destructive operation is unavoidable, restate branch, remote, expected effect, and recovery path, then wait for an explicit go.

## When to use this skill
- Use when: managing git, or doing a git task "the right way" / "like an expert".
- Use when: cleaning up, squashing, reordering, or rebasing commits on your branch.
- Use when: resolving a merge or rebase conflict.
- Use when: undoing a commit, recovering lost work, or fixing a git mistake safely.
- Use when: writing a simple MR/PR description, or taking a change all the way to an opened PR/MR.
- Do NOT use: to ONLY push an already-prepared single commit with the standard five-section template — use `publishing-git-review-requests`. To run a Jira issue through a fix-to-MR workflow — use `jira-fix-mr-workflow`. To localize or diagnose a bug — use `progressive-bug-hunter`. To design a CI/CD or release pipeline — use `agentic-workflow-design`. Never merge or approve a PR/MR on someone else's behalf.

## Modes
| Mode | Use when the user wants to | Output |
|------|----------------------------|--------|
| Operate | perform one git operation (branch, stage, commit, merge, rebase, cherry-pick, resolve conflict, undo) | executed git step(s) + a short summary |
| Describe | turn the current branch into a simple MR/PR description | the filled `templates/mr-description.md` |
| Ship | go end-to-end: clean history → describe → push → open PR/MR | opened PR/MR URL + summary |

Detect the mode from the request. A "ship" request with messy history runs Operate (cleanup), then Describe, then publish.

## Workflow

### Operate
1. Run the Safety preflight; name the branch, upstream, and what is already published.
2. State the single operation and its goal in one sentence.
3. Apply the best-practice steps for that operation, confirming each mutation:
   - commit messages and staging → `references/commit-conventions.md`
   - squash / reorder / amend / undo → `references/history-and-undo.md`
   - merge vs rebase, conflict resolution → `references/merge-and-conflicts.md`
4. Verify (`git status`, `git log --oneline`, `git diff`), confirm nothing unrelated changed, and report.

### Describe
1. Determine the base branch (tracked upstream, else the hosted default) and collect facts: `git log --oneline <base>..HEAD`, `git diff --stat <base>..HEAD`, and the diff itself for intent.
2. Derive what changed and why from the commits and diff — never invent changes that are not in the diff.
3. Fill `templates/mr-description.md` per `references/mr-description-method.md`: keep it simple — What and why, Changes, Testing, Risk and rollback. Reduce an empty section to one honest line.
4. Output the description as `MR-DESCRIPTION.md` or inline, per the user's preference.

### Ship
1. Run Operate cleanup only if history is messy (WIP/fixup commits); otherwise skip.
2. Run Describe to produce the description.
3. Confirm the provider (GitHub PR / GitLab MR) from the remotes (`references/publishing-flow.md`); if both are plausible, ask.
4. On explicit confirmation: push with upstream tracking, then open the PR/MR with the generated description as the body. → `references/publishing-flow.md`.
5. Report the PR/MR URL, head/base branches, commit range, validation status, and anything left unstaged.

## Simple MR/PR description contract
The description is **four sections**, concise, derived from the branch — never padded:

- **What and why** — one or two sentences: the change and its reason.
- **Changes** — the key changes as short bullets, grouped by area.
- **Testing** — what was run and the result, or "not run: <reason>".
- **Risk and rollback** — risk level and how to revert.

This is deliberately leaner than the fixed five-section template in `publishing-git-review-requests`; it favors clarity over ceremony while keeping the best-practice essentials. Full method and provider phrasing: `references/mr-description-method.md`.

## Output contract
- **Operate:** the executed, confirmed git command(s) plus a summary — operation, branch/upstream, resulting `git log --oneline` range, and any unrelated files left untouched.
- **Describe:** a filled MR/PR description in the four-section shape (`MR-DESCRIPTION.md` or inline) containing no invented changes.
- **Ship:** the opened PR/MR URL, provider, head/base branches, commit range, validation status, and remaining work or blockers.

Always report destructive operations that were refused or are pending confirmation.

## Constraints
- Confirm before every mutating command; never run a destructive or history-rewriting command on published/shared refs without explicit approval.
- Prefer recoverable operations; restate the recovery path before anything irreversible.
- Stage explicit paths; never blanket-stage or run broad formatters/generators unprompted.
- Conform to the repo's branch naming, commit style, and provider; do not silently switch a GitHub remote to GitLab or vice versa.
- Put no change in the description that is not in the diff; do not inflate risk or testing claims.
- Never merge, approve, or close a PR/MR on the team's behalf; the human reviewer decides.
- Do not access tokens, secrets, or environment variables; use `gh auth status` / `glab auth status` only to check auth, and never print credentials.

## References
| Need | Reference |
|------|-----------|
| allow/confirm/deny command tiers, recoverable-first, published vs unpublished history, reflog recovery | `references/git-safety-policy.md` |
| conventional commits, atomic-commit discipline, message anatomy | `references/commit-conventions.md` |
| squash/reorder/amend, revert vs reset vs reflog vs stash (unpublished only) | `references/history-and-undo.md` |
| merge vs rebase decision, fast-forward, conflict resolution, abort/recover | `references/merge-and-conflicts.md` |
| derive a simple-but-standard description from commits+diff; GitHub PR vs GitLab MR | `references/mr-description-method.md` |
| push and open PR/MR via gh/glab, provider inference, upstream tracking | `references/publishing-flow.md` |
| worked end-to-end example | `examples/ship-walkthrough.md` |
