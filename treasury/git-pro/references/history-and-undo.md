# History Editing and Undo

Rewrite only **unpublished** commits (see `references/git-safety-policy.md`). For published history, prefer `revert`.

## Interactive rebase — squash, reorder, fixup
`git rebase -i some-base` opens a todo list; change the leading verb per line:
- `pick` — keep the commit as-is.
- `reword` — keep the change, edit the message.
- `squash` / `s` — fold into the previous commit, combine messages.
- `fixup` / `f` — fold into the previous commit, discard this message.
- `drop` — remove the commit.
- reorder by moving lines.
Abort any time with `git rebase --abort`; on conflict, resolve then `git rebase --continue`.
Clean pattern: make small WIP commits while working, then squash them into a few coherent commits before review.

## Amend the last commit
- Fix the message: `git commit --amend`.
- Add a forgotten file: `git add some/path` then `git commit --amend --no-edit`.
- Only amend a commit you have NOT pushed to a shared branch.

## Undo, by intent
| Goal | Command | Notes |
|------|---------|-------|
| Discard staged-but-not-committed | `git restore --staged some/path` | keeps working-tree changes |
| Discard working-tree edits | `git restore some/path` | not recoverable — confirm |
| Move HEAD back, keep changes staged | `git reset --soft HEAD~1` | safest reset |
| Move HEAD back, keep changes unstaged | `git reset --mixed HEAD~1` | default reset |
| Move HEAD back, discard changes | `git reset --hard HEAD~1` | Deny-tier — confirm + state recovery |
| Undo a pushed commit | `git revert <sha>` | new commit; safe on shared history |
| Set work aside | `git stash push -m "note"` then `git stash pop` | inspect with `git stash list` |

## Recover lost work
`git reflog` shows every position HEAD has held. Find the lost commit's sha, then `git switch -c rescue <sha>` (non-destructive) to get it back. Reflog entries persist ~90 days by default.

## Rules
- Confirm before any reset, rebase, or stash drop.
- Never rebase or amend commits already on a shared/upstream branch.
- After rewriting, verify with `git log --oneline` and `git diff some-base..HEAD` before pushing.
