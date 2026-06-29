# Worked example — Ship a feature branch end to end

A messy feature branch taken from working tree to an opened PR.

## Starting state

```
$ git log --oneline main..HEAD
9c1f2a3 wip
4b8e0d1 fix typo
a23c9f8 more wip
77d1e4b add rate limiter to api client
```

Four commits, three of them WIP — not review-ready. The branch has never been pushed, so its history is unpublished and safe to rewrite.

## 1. Operate — clean the history (unpublished, so safe)

Confirm first, then squash the WIP commits into one coherent change with `git rebase -i main`:

```
pick   77d1e4b add rate limiter to api client
fixup  a23c9f8 more wip
fixup  4b8e0d1 fix typo
fixup  9c1f2a3 wip
```

Result:

```
$ git log --oneline main..HEAD
b5a77c0 feat(client): add token-bucket rate limiter to API client
```

## 2. Describe — generate the simple description

From `git log main..HEAD` plus `git diff --stat main..HEAD`, fill the four sections (intent comes from the commit + diff, nothing invented):

```
# feat(client): add token-bucket rate limiter to API client

## What and why
Adds a token-bucket rate limiter to the API client so we stop tripping the
upstream 429s under burst load.

## Changes
- client: token-bucket limiter (configurable rate/burst), applied per host
- tests: unit tests for refill, burst, and exhaustion paths

## Testing
- go test ./client/... -> pass
- manual: ran the burst script against staging, no 429s

## Risk and rollback
- Risk: low — opt-in via config, defaults to previous behavior
- Rollback: revert this PR
```

## 3. Ship — push and open the PR (confirmed)

Provider inferred as GitHub from the `github.com` remote. On explicit confirmation:

```
$ git push -u origin feature/api-rate-limit
$ gh pr create --base main --head feature/api-rate-limit \
    --title "feat(client): add token-bucket rate limiter to API client" \
    --body-file MR-DESCRIPTION.md
```

Report: PR opened at https://github.com/acme/app/pull/482 · base `main` <- `feature/api-rate-limit` · 1 commit · tests pass.

Note: the WIP squash was safe only because the branch was never pushed. Had it been shared, the history stays as-is and the cleanup step is skipped — new commits or a revert instead.
