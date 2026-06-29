# Git Safety Policy

The skill executes git. Every command falls into one of three tiers — classify before running.

## Command tiers

### Allow — run without asking (read-only)
`git status`, `git diff`, `git log`, `git show`, `git branch` (list), `git branch -vv`, `git remote -v`, `git rev-parse`, `git reflog`, `git stash list`, `git blame`, `gh auth status`, `glab auth status`. These never change the working tree, refs, or remote.

### Confirm — state the effect, then wait
Any command that changes the working tree, index, local refs, or remote:
`git add`, `git commit`, `git switch -c` / `git checkout -b`, `git merge`, `git rebase`, `git cherry-pick`, `git stash` / `git stash pop`, `git tag`, `git push`, `git remote add` / `set-url`, `gh pr create`, `glab mr create`.
Confirmation must name: the exact paths/refs, the commit message (if any), the branch, the remote, and the provider and base branch (for a PR/MR).

### Deny by default — refuse unless the user explicitly insists, with recovery stated
Destructive or history-rewriting commands, especially on **published/shared** refs:
`git push --force`, `git reset --hard`, `git clean -fd`, `git branch -D`, `git push origin --delete`, `git rebase` or `git commit --amend` on commits already pushed, `git filter-branch`, `git reflog expire --expire=now`.
When the user insists: restate branch, remote, expected effect, and recovery path; prefer `--force-with-lease` over `--force`; proceed only on an explicit go.

## Recoverable-first principle
Always reach for the operation that can be undone:
- Undo a pushed change with `git revert` (a new commit), not a history rewrite.
- Move HEAD back but keep work with `git reset --soft` / `--mixed`, never `--hard` by default.
- Stash or branch before any risky operation so the prior state is recoverable.
- A local mistake is almost always recoverable via `git reflog` for ~90 days — check it before declaring work lost.

## Published vs unpublished history
- **Unpublished** (never pushed, or pushed only to your own non-shared branch): safe to rewrite — squash, reorder, amend, rebase freely.
- **Published/shared** (pushed to a branch others may have pulled, or any protected/default branch): do NOT rewrite. Add new commits or `revert`. Rewriting shared history forces everyone else to recover.
- Determine state with `git branch -vv` (does it track an upstream?) and `git log @{u}..HEAD` / `git log HEAD..@{u}`.

## Reflog recovery (the safety net)
- `git reflog` lists where HEAD has been; `git reflog show some-branch` for a branch.
- Recover a "lost" commit: `git switch -c rescue <sha>` (non-destructive), or `git reset --hard <sha>` (itself Deny-tier — confirm).
- Recover a deleted branch: find its tip in the reflog and re-create it.

## Never
- Print tokens, credential-helper output, or environment variables.
- Convert a GitHub remote into a GitLab one (or vice versa) silently.
- Run formatters/generators that rewrite broad file sets unless required by the task and approved.
