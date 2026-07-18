# Run report: {routine} {version}

- Run dir: {run_dir}
- Launched: {started_utc} UTC — finished: {finished_utc} UTC
- Mode: run
- Verification: {check_run_result_line}
- Verdict: {complete | failed at node X | partial - N of M nodes built}

## Nodes

| # | node | executor | tier | gate | status | honesty | artifacts | bytes | verified |
|---|------|----------|------|------|--------|---------|-----------|-------|----------|
| 1 | {node-id} | {executor} | {tier} | {gate} | {status} | {planned/built/wired} | {paths} | {n} | {ok/MISSING/EMPTY} |

## Gates log

| gate | node | decision | note |
|------|------|----------|------|
| Gate 0 | - | approve | {human note if given} |

## Failures and lessons

- {node-id}: {what failed} — evidence: {script output line} — lesson: {one line}
- (none)

## Next steps

- {follow-up the human should consider, or "none"}
