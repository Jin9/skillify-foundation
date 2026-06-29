# Merge, Rebase, and Conflict Resolution

## Merge vs rebase — choosing
| Situation | Prefer | Why |
|-----------|--------|-----|
| Integrating a shared/long-lived branch (e.g. main into your feature) | `merge` | preserves true history; no shared-history rewrite |
| Tidying your own feature branch before review | `rebase` onto the base | linear, readable history (unpublished only) |
| Branch already pushed and shared | `merge` | rebasing would rewrite history others hold |
| Pulling updates on a feature branch | `git pull --rebase` (if unshared) or `merge` | avoids noise merge commits when safe |

Fast-forward: when the base has no new commits, a merge just advances the pointer. Use `--no-ff` only when you want an explicit merge commit to record the integration.

## Resolving conflicts — step by step
1. Start the operation (`git merge some-branch` or `git rebase some-base`). On conflict, git pauses.
2. List conflicts: `git status` (look for "both modified").
3. Open each file; resolve between the conflict markers (the `<<<<<<<` / `=======` / `>>>>>>>` blocks). Understand *both* sides before choosing — do not blindly keep one.
4. For complex cases use `git mergetool`, or `git checkout --ours some/path` / `--theirs some/path` when one side is wholly correct. Know which is which: during a rebase, "ours" is the base you are replaying onto.
5. Stage resolved files: `git add some/path`.
6. Continue: `git merge --continue` or `git rebase --continue`. During a rebase, repeat per conflicting commit.
7. Verify: build and run tests after resolution — a textually clean merge can still be semantically broken.

## Recover / abort
- `git merge --abort` or `git rebase --abort` returns to the pre-operation state.
- If you continued past a bad resolution, `git reflog` holds the pre-merge sha (see `references/git-safety-policy.md`).

## Rules
- Confirm before starting a merge or rebase; state the source and target.
- Never resolve by deleting conflict markers without reading both sides.
- Re-run tests after every non-trivial resolution and report what was run.
