# MR/PR Description Method

Goal: a description a reviewer can read in under a minute that still carries everything they need. Simple, not sparse.

## Derive, do not invent
1. Collect the branch facts: `git log --oneline some-base..HEAD` (the narrative) and `git diff --stat some-base..HEAD` (the surface area); skim the diff for intent.
2. The **why** comes from the commit bodies and the linked ticket. The **what** comes from the diff. If a section has no basis in either, write the honest minimum ("Testing — not run: local-only change") rather than padding.
3. Never list a change that is not in the diff, and never claim testing that did not happen.

## The four sections
- **What and why** — one or two sentences: the change and the reason it exists. Lead with the user/behavior impact, not the implementation.
- **Changes** — short bullets grouped by area (e.g. `api:`, `store:`, `tests:`). Summarize; the diff has the detail.
- **Testing** — commands run and their result, or "not run: reason". Include manual verification steps when relevant.
- **Risk and rollback** — risk level (low/medium/high) with one line of why, and how to revert (usually "revert this MR/PR").

## Keep-it-simple rules
- One screen or less. If it scrolls, cut.
- No restating the diff line by line; no boilerplate headings with nothing under them.
- Flag breaking changes and migrations explicitly — these are the one thing reviewers must not miss.
- Match the repo's own MR/PR template if one exists (a pull-request or merge-request template under `.github/` or `.gitlab/`); this four-section shape is the default only when none exists.

## GitHub vs GitLab phrasing
- GitHub: "pull request" / "PR"; opened with `gh pr create`.
- GitLab: "merge request" / "MR"; opened with `glab mr create`.
- Use the right term for the provider in both the description and the report. See `references/publishing-flow.md`.

## Relation to publishing-git-review-requests
That skill fills a fixed five-section template (Summary, Changes, Validation, Risks, Reviewer Notes) for a single prepared change. This method generates a leaner four-section description from the whole branch. Prefer the repo's own template when present; otherwise use this default.
