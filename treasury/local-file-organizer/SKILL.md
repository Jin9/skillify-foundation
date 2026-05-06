---
name: local-file-organizer
description: Plans and applies safe, offline organization of local files and folders. Use when the user says "organize this folder", "clean up my files", "classify my documents", "make folder taxonomy", "generate a safe move plan", "create rollback.sh", or asks for an offline/local file organizer. Produces an inventory, a JSON move plan, an organization-report.md, an apply.sh, and a rollback.sh; never deletes or overwrites user files, and applies moves only after explicit user approval. Do NOT use for cloud storage sync, code refactoring, repository housekeeping, or content-hash deduplication.
---

# Local File Organizer

Plan and apply safe, offline organization of local files. Always plan first, get explicit approval, then apply with safe move operations only.

## When to use

- Organizing a messy folder by business meaning.
- Producing a Markdown organization report and a rollback script before any move.
- Working in offline / no-internet mode using only local shell tools.

## When NOT to use

- Modifying files inside `.git/`, `node_modules/`, build outputs, or dependency caches.
- Cloud-sync orchestration (Drive, Dropbox, iCloud) — paths must be local.
- Code refactoring or repo housekeeping — those belong in repo-policy or refactor skills.

## Hard safety rules

1. Never delete user files. Never overwrite. Never run destructive commands.
2. Never use `rm` (except on the explicit temp artifacts listed in `references/script-generation.md`), `find -delete`, `trash`, `unlink`, or forced-overwrite commands against user files.
3. Never apply moves without explicit user approval of `move_plan.json`.
4. Always write `organization-report.md` and `rollback.sh` at the target root before applying.
5. Files with confidence `< 0.80` do not enter business folders. `0.60`–`0.79` → `_REVIEW/`. `< 0.60` → leave in place unless the user overrides.
6. On filename conflict, keep both files using a safe `__duplicate-NNN` suffix. Never overwrite.
7. Skip hidden system folders, dependency folders, and secret files (`.env*`, `*.key`, `*.pem`, `id_rsa*`) unless the user explicitly opts in.

## Operating mode

- Assume offline. Use only local shell tools: `pwd`, `ls`, `find`, `stat`, `file`, `du`, `shasum`, `grep`, `awk`, `sed`.
- Do not fetch remote resources or call LLM APIs unless the user explicitly enables them.
- For local-LLM notes, OS-specific inventory commands, and skip lists, see `references/inventory-and-models.md`.

## Workflow

1. **Resolve target root.** Use the current working directory unless the user specifies another path. Confirm with the user.
2. **Build inventory.** Generate `inventory.tsv` with columns `path, filename, extension, size_bytes, modified_at`. Apply skip patterns from `references/inventory-and-models.md`.
3. **Classify by business meaning.** Use filename, extension, path, size, and modified time first. Read file contents only with explicit user approval.
4. **Propose taxonomy.** Default scheme and rules in `references/taxonomy-and-naming.md`. Prefer business meaning over file type.
5. **Emit `move_plan.json`** at the target root following the schema in `templates/move-plan.json` (top-level keys: `folders`, `moves`, `review_required`).
6. **Validate the plan.** Apply confidence routing from `references/taxonomy-and-naming.md`. Reject any plan containing `rm`, overwrite, or hidden-folder modification.
7. **Write `organization-report.md`** at the target root using `templates/organization-report.md`.
8. **Write `rollback.sh`** at the target root following `references/script-generation.md`.
9. **Show a 5–10 line summary** of the report and ask for explicit approval.
10. **On approval, generate and run `apply.sh`** following `references/script-generation.md`. Use `mv -n` only.
11. **On successful apply, run the safe-cleanup command** from `references/script-generation.md` to remove only the known temp artifacts.
12. **Final state at target root:** keep only `organization-report.md` and `rollback.sh`. Do not delete source folders or generated destination folders.

## Output contract

Final files at `<target-root>/` (kept after completion):

- `organization-report.md` — required.
- `rollback.sh` — required, executable, uses `mv -n` only.

Intermediate files (must be removed after successful apply):

- `inventory.tsv`, `move_plan.json`, `apply.sh`, `dry-run.sh`, `organization-plan.json`, `file-classification.json`, `move-log.tsv`.

## Completion criteria

The task is complete only when:

1. Moves performed match the approved `move_plan.json`.
2. No user file was deleted or overwritten.
3. `organization-report.md` and `rollback.sh` exist at the target root.
4. All temp artifacts listed above are removed.
5. Low-confidence files are in `_REVIEW/` or untouched.
6. Filename conflicts were resolved by suffix or `_DUPLICATES/`, never overwrite.

## References

| Need | File |
|------|------|
| Default folder taxonomy, naming, and conflict rules | `references/taxonomy-and-naming.md` |
| `apply.sh` and `rollback.sh` generation rules + safe-cleanup command | `references/script-generation.md` |
| OS-specific inventory commands, skip lists, secret-file list, local-model notes | `references/inventory-and-models.md` |

## Templates

- `templates/organization-report.md` — Markdown report skeleton.
- `templates/move-plan.json` — JSON move-plan skeleton.
