# Publishing Flow — push and open the PR/MR

The skill opens the review request directly (it does not delegate). Every step here is Confirm-tier (see `references/git-safety-policy.md`).

## Preflight
- `git status --short`, `git branch -vv`, `git remote -v`, `git rev-parse --show-toplevel`.
- Auth: `gh auth status` (GitHub) / `glab auth status` (GitLab) — only to check; never print tokens.
- Confirm there is something to push and a current branch exists; if not, ask.

## Provider inference
- GitHub: `github.com` remotes, `gh` available, "PR"/"pull request" wording.
- GitLab: `gitlab.com` or self-managed remotes, `glab` available, "MR"/"merge request" wording.
- If both are plausible, ask before pushing. Never convert one provider's remote into the other's.

## Push
- Set upstream on first push: `git push -u origin some-branch`.
- Do NOT push to a protected/default branch from a feature flow; open from your branch.
- If the push is rejected (non-fast-forward), inspect with `git status` and `git branch -vv`; do NOT force-push by default. A force-push is Deny-tier — confirm and prefer `--force-with-lease`.

## Open the PR/MR
- Base = the hosted default branch unless the user names another; head = the current upstream branch.
- GitHub: `gh pr create --base some-base --head some-branch --title "..." --body-file MR-DESCRIPTION.md`.
- GitLab: `glab mr create --source-branch some-branch --target-branch some-base --title "..." --description "$(cat MR-DESCRIPTION.md)"`.
- Title: the conventional-commit summary of the branch's primary change.
- Body: the generated four-section description (`references/mr-description-method.md`).

## Defaults and failure handling
- New remote/repo: prefer `origin`; default new repos to private; confirm before `gh repo create` / `glab repo create`.
- Missing or failed auth: surface the `auth status` failure, give the exact manual steps, and stop before pushing.
- Report the final PR/MR URL, provider, head/base branches, and validation status.
