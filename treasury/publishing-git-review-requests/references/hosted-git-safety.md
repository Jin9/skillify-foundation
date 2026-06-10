# Hosted Git Safety Rules

## Preflight checks

Run read-only checks before planning any mutation:

- `git status --short`
- `git branch --show-current`
- `git branch -vv`
- `git remote -v`
- `git rev-parse --show-toplevel`

Use `gh auth status` for GitHub and `glab auth status` for GitLab only to check
authentication. Do not print tokens, credential helper output, or environment
variables.

## Provider selection

- Infer GitHub from `github.com` remotes, `gh` wording, or PR terminology when
  no GitLab signals are present.
- Infer GitLab from `gitlab.com`, self-managed GitLab remotes, `glab` wording,
  MR terminology, or "merge request".
- If both providers are plausible, ask before mutating.
- Do not silently convert a GitHub remote into GitLab or a GitLab remote into
  GitHub.

## Confirmation boundary

Treat these as mutating commands that need explicit approval:

- `git add`
- `git commit`
- `git remote add`
- `git remote set-url`
- `gh repo create`
- `glab repo create`
- `git push`
- `gh pr create`
- `glab mr create`

Approval must identify what will happen: paths to stage, commit message,
provider, target remote, branch name, repo/project visibility, and PR/MR base
branch.

## Preserving user work

- Stage explicit paths, not every dirty file.
- If unrelated files are dirty, leave them unstaged and mention them in the
  final report.
- Do not run formatters or generators that rewrite broad file sets unless they
  are required by the task and approved.
- Do not overwrite an existing remote URL. Surface the mismatch and ask before
  changing it.

## Destructive operations

Refuse or pause for explicit confirmation before:

- `git push --force` or `git push --force-with-lease` (when the user insists
  on a force push, prefer `--force-with-lease` — it still rewrites the remote
  ref but refuses to clobber work pushed by someone else)
- `git reset --hard`
- `git clean`
- deleting local or remote branches
- amending commits that may already be pushed
- rebasing a branch with a shared upstream

When the user explicitly asks for one of these operations, restate the branch,
remote, expected effect, and recovery risk before proceeding.

## Remote creation defaults

- Prefer `origin` as the remote name for a new hosted repository or project.
- Default new GitHub repositories and GitLab projects to private visibility.
- Use the current directory name as the hosted repo/project name only when the
  user did not provide a name and the name is safe kebab-case or snake_case.
- If `gh` or `glab` is unavailable or unauthenticated, provide the exact
  manual hosted-Git steps and stop before pushing.

## PR and MR creation defaults

- Use the tracked upstream branch to infer the head branch.
- Prefer the hosted default branch as base when it can be discovered.
- Do not open a PR or MR from a protected default branch unless the user
  explicitly asked to work directly from that branch.
- Include validation status even when no tests were run.
- Use "PR" for GitHub and "MR" for GitLab in the final report.

## Failure handling

- Existing remote conflict: show current and intended remote URLs; ask before
  changing.
- Push rejected: run read-only inspection such as `git status` and
  `git branch -vv`; do not force-push by default.
- Missing GitHub auth: surface `gh auth status` failure and ask the user to
  authenticate.
- Missing GitLab auth: surface `glab auth status` failure and ask the user to
  authenticate.
- No current branch: ask the user for a branch name before committing or
  pushing.
- No staged or dirty changes: skip commit and continue only if the requested
  operation still makes sense.
