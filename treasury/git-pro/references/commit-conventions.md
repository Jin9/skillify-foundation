# Commit Conventions

## Atomic commits
One logical change per commit. A commit should build and pass tests on its own and be revertable without dragging in unrelated work.
- Stage by intent, not by "everything dirty": `git add some/path` or `git add -p` to split a file's hunks.
- If a change mixes concerns, split it into separate commits.
- Never bundle a refactor with a behavior change in one commit — reviewers cannot tell them apart.

## Conventional Commits format
```
type(scope): summary

body — why the change is needed and what it does, wrapped at ~72 cols

footer — BREAKING CHANGE: ... / Refs: ABC-123
```
- **type**: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `build`, `ci`, `chore`, `revert`.
- **scope** (optional): the area touched, e.g. `feat(auth):`.
- **summary**: imperative mood, lower case, no trailing period, around 50 chars ("add", not "added"/"adds").
- **body** (optional, preferred for non-trivial changes): the *why*, not a restatement of the diff.
- **BREAKING CHANGE:** in the footer — or a `!` after type/scope, e.g. `feat(api)!:` — for incompatible changes.

## Good vs weak subjects
- Good: `fix(parser): handle empty CSV header row`
- Good: `refactor(store): extract retry logic into helper`
- Weak: `fix bug`, `update code`, `wip`, `changes`.

## Message anatomy checklist
- [ ] Subject is imperative, scoped, and around 50 chars.
- [ ] Body explains *why* (and any non-obvious *how*), wrapped at ~72 cols.
- [ ] Breaking changes flagged in the footer or with `!`.
- [ ] Issue/ticket referenced in the footer when one exists.
- [ ] No secrets, tokens, or absolute local paths in the message.

## House conformance
Before writing a message, read the repo's recent log (`git log --oneline -n 20`) and match its existing convention — type set, scope style, ticket-reference format. Conform to what the repo already does rather than imposing this template.
