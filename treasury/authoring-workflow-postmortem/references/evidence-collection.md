# Evidence collection

Run these read-only steps in order. Output goes to the archive folder
`docs/postmortems/YYYY-MM-DD-<workflow_id>/`.

## 1. Set variables

```bash
DATE=$(date -u +%F)        # UTC date for the file/folder name
ID="<workflow_id>"         # from state.json
ARCHIVE="docs/postmortems/${DATE}-${ID}"
mkdir -p "$ARCHIVE/stages"
```

## 2. State snapshot

```bash
cp .agent/state.json "$ARCHIVE/state.json"
```

## 3. Runlog (structured + human)

```bash
cp .agent/runlog.jsonl "$ARCHIVE/runlog.jsonl"
tail -n 100 .agent/runlog.md > "$ARCHIVE/runlog.md"
```

## 4. Stage outputs

```bash
cp .agent/stages/*.md "$ARCHIVE/stages/" 2>/dev/null || true
```

For the failed stage(s), capture the redacted log tail too:

```bash
for s in <failed-stage-name>; do
  tail -n 200 ".agent/stages/${s}.log" > "$ARCHIVE/${s}.log"
done
```

(Logs are redacted at write time per the scaffold's secret-redaction
contract; the archive is safe to commit.)

## 5. Approval log

```bash
just approvals-verify > "$ARCHIVE/approvals-verify.txt"
```

If verification fails, **stop**. Treat as a security incident before
proceeding.

```bash
# excerpt approvals for this workflow_id
grep -A 3 "workflow=${ID}" .agent/approvals.md > "$ARCHIVE/approvals.md" || true
```

## 6. Spend snapshot

```bash
just llm-spend 24 > "$ARCHIVE/spend.txt"
```

Compare against `state.json.total_cost_usd`. Drift between LiteLLM and
state.json is itself a postmortem-worthy finding (record it under
contributing factors).

## 7. Sandbox logs (only if sandbox trigger)

```bash
docker compose -f docker/sandbox-compose.yml logs proxy --since 24h \
  > "$ARCHIVE/proxy.log" 2>&1
```

## 8. Cross-check

Before drafting, verify the archive contains:

- [ ] `state.json`
- [ ] `runlog.jsonl` and `runlog.md`
- [ ] `stages/*.md` for every completed stage
- [ ] `<failed>.log` for every failed stage
- [ ] `approvals-verify.txt` showing `OK`
- [ ] `approvals.md` excerpt
- [ ] `spend.txt`
- [ ] `proxy.log` if applicable

Missing items mean the postmortem is incomplete; fill the gap before
drafting.

## Redaction reminder

The scaffold redacts at log-write time, but cross-check the archive for:

- API keys (any string matching `sk-`, `AKIA`, `xoxb-`, JWT `eyJ`)
- Email addresses outside example.com
- Real customer names, account numbers, NIK, PAN

If found, redact in the archive copy before committing.
