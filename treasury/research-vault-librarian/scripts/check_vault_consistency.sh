#!/usr/bin/env bash
# check_vault_consistency.sh — READ-ONLY consistency gate for the ResearchVault.
#
# Operationalizes the consistency one-liners documented in the vault's own
# CLAUDE.md; that CLAUDE.md remains the source of truth for the cross-file
# invariant. This script only reads (grep/find/sed/wc): it never writes, moves,
# deletes, or makes network calls. Exit 0 iff every check passes.
#
# Usage: check_vault_consistency.sh [vault-root]
#   vault-root resolution: $1  ->  $RESEARCH_VAULT_ROOT  ->  default below.

# No portable default: pass a vault root as $1 or set RESEARCH_VAULT_ROOT.
VAULT="${1:-${RESEARCH_VAULT_ROOT:-}}"

if [ -z "$VAULT" ]; then
  printf 'error: no vault root given. Pass it as the first argument or set RESEARCH_VAULT_ROOT.\n' >&2
  exit 2
fi

fail=0
pass_line() { printf '[PASS] %s\n' "$1"; }
fail_line() { printf '[FAIL] %s\n' "$1"; fail=1; }

if [ ! -d "$VAULT" ]; then
  printf 'error: vault root not found: %s\n' "$VAULT" >&2
  exit 2
fi
cd "$VAULT" || { printf 'error: cannot cd into vault: %s\n' "$VAULT" >&2; exit 2; }
if [ ! -f CLAUDE.md ] || [ ! -f index.md ] || [ ! -d maps ] || [ ! -d reports ]; then
  printf 'error: not a ResearchVault layout (need CLAUDE.md, index.md, maps/, reports/): %s\n' "$VAULT" >&2
  exit 2
fi

printf 'Vault: %s\n\n' "$VAULT"

# 1. MOC wikilinks whose target report file is missing.
missing=$(grep -ohE '\[\[reports/[a-z0-9./-]+' maps/*.md | sed 's#\[\[##' \
  | while read -r p; do [ -f "$p.md" ] || echo "  MISSING: $p"; done)
if [ -z "$missing" ]; then pass_line "every MOC wikilink resolves to a report file"
else fail_line "MOC wikilinks with no target file:"; printf '%s\n' "$missing"; fi

# 2. Report files not referenced in any MOC.
unlisted=$(find reports -name '*.md' | while read -r f; do rel="${f%.md}"; \
  grep -q "\[\[$rel|" maps/*.md || echo "  UNLISTED: $rel"; done)
if [ -z "$unlisted" ]; then pass_line "every report is listed in a MOC"
else fail_line "report files in no MOC:"; printf '%s\n' "$unlisted"; fi

# 3. files == MOC links == index.md header total.
files=$(find reports -name '*.md' | wc -l | tr -d '[:space:]')
links=$(grep -hoE '\[\[reports/[a-z0-9./-]+' maps/*.md | wc -l | tr -d '[:space:]')
header=$(grep -E '\*\*Total\*\*' index.md | grep -oE '[0-9]+' | tail -1)
printf 'Reconciliation: files=%s links=%s header-total=%s\n' "$files" "$links" "${header:-?}"
if [ -n "$header" ] && [ "$files" = "$links" ] && [ "$files" = "$header" ]; then
  pass_line "files == links == header total ($files)"
else
  fail_line "count mismatch (files=$files links=$links header=${header:-?})"
fi

# 4. Per-domain folder count == its maps/<folder>.md bullet count.
drift=$(for d in reports/*/; do f=$(basename "$d"); \
  a=$(find "$d" -name '*.md' | wc -l | tr -d '[:space:]'); \
  b=$(grep -c "^- \[\[reports/$f/" "maps/$f.md" 2>/dev/null || echo 0); \
  [ "$a" = "$b" ] || echo "  MOC DRIFT $f: folder=$a moc=$b"; done)
if [ -z "$drift" ]; then pass_line "every domain folder count == its MOC bullet count"
else fail_line "per-domain MOC drift:"; printf '%s\n' "$drift"; fi

# 5. Every report appears in exactly one MOC (its own domain).
not1=$(for x in reports/*/*.md; do fo=$(basename "$(dirname "$x")"); \
  sl=$(basename "$x" .md); \
  n=$(grep -l "\[\[reports/$fo/$sl|" maps/*.md 2>/dev/null | wc -l | tr -d '[:space:]'); \
  [ "$n" = 1 ] || echo "  NOT-1-MOC: $fo/$sl ($n)"; done)
if [ -z "$not1" ]; then pass_line "every report appears in exactly one MOC"
else fail_line "reports not in exactly one MOC:"; printf '%s\n' "$not1"; fi

# 6. index.md links all MOCs.
notlinked=$(ls maps/*.md | sed 's#maps/##;s#\.md##' | while read -r m; do \
  grep -q "\[\[maps/$m|" index.md || echo "  MOC NOT LINKED IN INDEX: $m"; done)
if [ -z "$notlinked" ]; then pass_line "index.md links every MOC"
else fail_line "MOCs not linked in index.md:"; printf '%s\n' "$notlinked"; fi

printf '\n'
if [ "$fail" -eq 0 ]; then printf 'RESULT: ALL CHECKS PASS\n'; exit 0
else printf 'RESULT: FAIL — see [FAIL] lines above\n'; exit 1; fi
