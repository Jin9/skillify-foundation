**Your push was rejected because the remote branch has commits your local branch does not have.**

1. Someone pushed to `main` after your last `git pull`.
2. Git refuses to overwrite their commits, so it stops with `! [rejected] main -> main (non-fast-forward)`.
3. Nothing is lost. Your commits are still on your machine.

**For you:** run `git pull --rebase origin main`, fix any conflicts, then `git push` again.

**Words:**
- **non-fast-forward**: your branch is missing commits that the remote already has.
- **rebase**: replay your commits on top of the newest remote commits.

Say **more** for the next layer, or tell me which step is unclear.
