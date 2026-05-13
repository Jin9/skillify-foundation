# Inventory Commands and Local-Model Notes

## Skip patterns

Always skip unless the user explicitly opts in:

```text
.git/  node_modules/  vendor/  dist/  build/  .next/  .cache/  .DS_Store
venv/  .venv/  target/  coverage/
```

Never move secret files unless the user opts in:

```text
.env  .env.*  *.key  *.pem  *.crt  *.p12  *.jks  id_rsa  id_ed25519
```

When the user opts in, route sensitive files to `_REVIEW/security-sensitive/`.

## Inventory command — Linux

```bash
find . -type f \
  -not -path '*/.git/*' \
  -not -path '*/node_modules/*' \
  -not -path '*/vendor/*' \
  -not -path '*/dist/*' \
  -not -path '*/build/*' \
  -printf '%p\t%f\t%s\t%TY-%Tm-%Td %TH:%TM\n' \
  > inventory.tsv
```

## Inventory command — macOS

```bash
find . -type f \
  -not -path '*/.git/*' \
  -not -path '*/node_modules/*' \
  -not -path '*/vendor/*' \
  -not -path '*/dist/*' \
  -not -path '*/build/*' \
  -exec stat -f '%N\t%z\t%Sm' -t '%Y-%m-%d %H:%M:%S' {} \; \
  > inventory.tsv
```

## Local-model notes (optional)

If the user runs classification through a local LLM (e.g. Ollama), these models have been used successfully — none is required:

- Primary candidate: `qwen3-coder:30b`
- Backups: `devstral-small-2:24b-instruct-2512-q4_K_M`, `mistral-small3.2:24b`
- Embeddings: `bge-m3`

Use whatever the user already has installed; do not prompt them to install anything new.
