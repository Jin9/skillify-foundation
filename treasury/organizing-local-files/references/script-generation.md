# Script Generation Rules

All generated shell scripts must:

- Start with `set -euo pipefail`.
- Quote every path.
- Use `mkdir -p` for target directories.
- Use `mv -n` (no clobber). Never plain `mv`, `cp -f`, or `rsync --delete`.
- Log every action on stdout.
- Never modify `.git/`, `node_modules/`, build outputs, dependency caches, or hidden system folders unless the user explicitly opts in.
- Never use `rm` against user files. The only `rm` allowed is the safe-cleanup command in §Safe cleanup.
- Never use wildcard deletes (`rm -f *.json`, `rm -rf *`, `find . -delete`).

## apply.sh

- Generate only after the user approves `move_plan.json`.
- Place at `<target-root>/apply.sh` temporarily.
- For each move: `mkdir -p` the target dir, then `mv -n "<source>" "<target>"`.
- If `mv -n` refuses (target exists), route the source into `<target-root>/_DUPLICATES/` with a `__duplicate-NNN` suffix.
- Optionally append a row per successful move to `<target-root>/move-log.tsv`.
- Do not remove empty source folders. Do not remove source folders.
- After all moves succeed, run the safe-cleanup command below. Do NOT delete `organization-report.md` or `rollback.sh`.

## rollback.sh

- Place at `<target-root>/rollback.sh`. Keep after apply.
- For each move in `move_plan.json`, generate the reverse `mv -n "<target>" "<source>"`.
- If the reverse target already exists, route the file into `<target-root>/_ROLLBACK_CONFLICTS/` instead.
- Log every rollback action.
- Never delete `organization-report.md` or `rollback.sh`.
- Must be runnable from the target root: `./rollback.sh`.

## Safe cleanup

Only this exact command is permitted, and only after a successful apply:

```bash
rm -f \
  "./inventory.tsv" \
  "./move_plan.json" \
  "./apply.sh" \
  "./dry-run.sh" \
  "./organization-plan.json" \
  "./file-classification.json" \
  "./move-log.tsv"
```

Forbidden cleanup forms:

```bash
rm -f *.json
rm -rf *
find . -delete
```
