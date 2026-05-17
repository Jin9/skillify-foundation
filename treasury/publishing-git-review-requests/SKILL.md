---
name: publishing-git-review-requests
description: >
  Publishes one local change to GitHub or GitLab for hosted review: prepare an
  intended commit, create or attach origin, push the branch, and open a PR or
  MR. Use when the user asks to "create remote and push", "commit change and
  push it", "open PR", "open MR", "create merge request", or "push this
  branch". Do NOT use for releases, deployments, CI authoring or security
  review, Bitbucket, broad repo administration, or destructive history rewriting.
compatibility: codex, claude-code, copilot, gemini
---

# Publishing Git Review Requests

## Purpose

Move one local change from a checkout to a hosted Git review flow on GitHub or
GitLab, with explicit confirmation before every local or remote mutation.

## When to use this skill

- Use when the user asks to "create remote and push".
- Use when the user asks to "commit change and push it".
- Use when the user asks to "open PR", "open MR", or "create merge request".
- Use when a local branch needs its first upstream push before a PR or MR.
- Do NOT use for release publishing, deployments, CI pipeline authoring, CI
  security review, Bitbucket, repo-policy edits, or rewriting published history.

## Core workflow

1. **Classify the review-publishing request.**
   - Identify the operation: commit-only, remote setup, push-only, review
     request only, or full workflow.
   - Keep the scope to one local change and one target review request. If the
     user asks for unrelated repository administration, split that work out.
   - Read `references/hosted-git-safety.md` before any operation that could
     mutate local history or hosted Git state.

2. **Inspect repository and provider state.**
   - Run read-only checks first: `git status --short`,
     `git branch --show-current`, `git branch -vv`, `git remote -v`, and
     `git rev-parse --show-toplevel`.
   - Infer provider from the existing remote URL, user wording, or available
     CLI: GitHub uses `gh`; GitLab uses `glab`.
   - If provider remains ambiguous, ask the user to choose GitHub or GitLab
     before planning any mutation.

3. **Build the mutation plan.**
   - List the provider, exact files to stage, commit message, branch name,
     remote name, remote URL or hosted repo/project path, push target, and PR/MR
     base branch.
   - Preserve unrelated dirty files. Stage only the paths required for the
     user's request.
   - Default newly created hosted repos/projects to private visibility unless
     the user explicitly asks for public.
   - Use `templates/review-request-description.md` for PR or MR body content.

4. **Request confirmation before mutating.**
   - Show the commands that will change local or hosted state.
   - Wait for explicit user approval before running any of: `git add`,
     `git commit`, `git remote add`, `git remote set-url`, `gh repo create`,
     `glab repo create`, `git push`, `gh pr create`, or `glab mr create`.
   - If approval is only for one step, perform only that step and ask again
     before the next mutation.

5. **Execute in the safest order.**
   - Commit: stage intended files, run the requested tests or checks when
     practical, then commit with the approved message.
   - Remote setup: prefer `origin`; if `origin` exists, verify it points to the
     intended GitHub repository or GitLab project before changing anything.
   - First push: push the current branch with upstream tracking.
   - Review request: create the GitHub PR or GitLab MR against the approved
     base branch and include summary, validation, and risk notes.

6. **Report the result.**
   - Include operation, provider, commit hash, branch, upstream, remote URL,
     PR/MR URL, and tests or checks run.
   - Call out skipped steps, auth failures, existing remote conflicts, protected
     branch issues, ambiguous provider choices, or files left unstaged.

## Output format

Respond with a concise execution report:

```text
Git review publishing result
- Operation: commit, remote setup, push, PR/MR, or full workflow
- Provider: GitHub, GitLab, or unresolved
- Commit: hash and subject, or skipped reason
- Branch: current branch and upstream
- Remote: name and URL
- Review request: PR/MR URL, or skipped reason
- Validation: tests/checks run, or not run reason
- Remaining work: blockers, conflicts, or files intentionally left alone
```

## Constraints and anti-patterns

- DO NOT run mutating git, GitHub, or GitLab commands before explicit
  confirmation.
- DO NOT stage all files blindly; preserve unrelated user changes.
- DO NOT force-push, reset hard, delete branches, or amend published commits
  unless the user explicitly asks for that exact destructive action.
- DO NOT print tokens, environment variables, credential helper output, or CLI
  auth secrets.
- DO NOT create public repositories or projects unless the user explicitly
  confirms public visibility.
- DO NOT bypass repository instructions, branch protection, required review
  gates, or merge request approval rules.

## Validation checklist

Before finalizing, verify:

- [ ] Repository state was inspected before mutation.
- [ ] Provider was explicit or safely inferred from the remote or user wording.
- [ ] User approved every mutating command that was run.
- [ ] Only intended files were staged and committed.
- [ ] Remote points to the intended GitHub repository or GitLab project.
- [ ] Branch push target and PR/MR base branch are explicit.
- [ ] Final report includes provider, commit, branch, remote, PR/MR, and
      validation status.

## References

- For safety rules and failure handling, see `references/hosted-git-safety.md`.
- For PR/MR body structure, see `templates/review-request-description.md`.
