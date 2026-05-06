# Walkthrough: Code-review domain

A worked end-to-end run of the `code-review` phase shape — every phase runs:
**Plan → Gather → Analyze → Review → Validate → Decide → Compact**.

## User prompt

> "Audit the new payment-retry handler in `internal/payment/retry.go` before we merge. Find security and concurrency issues, and tell me whether it's ready to ship."

## Step 1 — Universal preamble

- Mode: `Compose`
- Domain: `code-review` (matches "audit", "find issues", "ready to ship")
- Phase shape printed; user replies `proceed`.

## Step 2 — Initialize state

```
$ python3 scripts/init_pipeline.py --slug payment-retry-audit --domain code-review --prompt "Audit the new payment-retry handler..."
created: <cwd>/.claude/pipelines/payment-retry-audit-9b1d3e
```

## Step 3 — Plan (hard gate)

1 × `Plan` agent decomposes:

```
1. q1: What does retry.go do, and what state does it touch?
2. q2: Are there race conditions in the retry loop?
3. q3: Is the retry idempotent? What guards against double-charging?
4. q4: Are timeouts and circuit breakers configured?
5. q5: What logging exists for security forensics?
```

User confirms `proceed`.

## Step 4 — Gather (5 parallel `Explore` agents in one message)

Each agent reads `retry.go` and adjacent files (e.g. `payment_test.go`, `internal/idempotency/`).

Results in `02-evidence/q1.md` through `q5.md`. Each finding cites file:line.

## Step 5 — Analyze

1 × `general-purpose` agent produces `03-analysis.md`:

- F1 (P1): Retry uses a per-request mutex but the request-id is not the idempotency key — risk of double-charging on duplicate Kafka delivery. Cites `02-evidence/q3.md#F1`.
- F2 (P1): Circuit breaker timeout (60s) exceeds the upstream gateway timeout (30s) — orphaned charges possible. Cites `02-evidence/q4.md#F2`.
- F3 (P2): No structured logging on retry abandonment — forensic gap.
- F4 (P3): Magic-number constants (`maxAttempts = 5`).

## Step 6 — Review (one pass, hard cap)

1 × `Plan` agent (adversarial) produces `04-review.md`:

- I1 (P1) targets F1: "Analysis assumes Kafka at-least-once delivery — verify the consumer config." → suggested fix: gather config.
- I2 (P2) targets F3: "Forensic gap is real but not unique to this PR — out of scope."

Orchestrator surfaces I1/I2 to user. User accepts I2 as out-of-scope, asks to address I1 in Validate.

## Step 7 — Validate (parallel `Explore` agents)

P1/P2 claims = 3. Orchestrator spawns 3 × `Explore` agents to fact-check, including a new check on the Kafka consumer config (per I1).

Result `05-validation.md`:
- C1 (F1 idempotency): confirmed — Kafka consumer config does NOT dedupe.
- C2 (F2 timeout): confirmed.
- C3 (F3 logging): confirmed.

## Step 8 — Decide (hard gate)

1 × `Plan` agent reads analysis + review + validation, produces `06-decision.md`:

| # | Option | Pros | Cons | Risks | Effort | Reversibility |
|---|--------|------|------|-------|--------|---------------|
| 1 | Block merge — fix F1 + F2 first | Closes both P1s | 1–2 days | Schedule slip | M | easy |
| 2 | Merge with feature flag | Ships traffic-gated | Flag complexity | Misconfig risk | S | easy |
| 3 | Merge as-is | Fastest | Both P1s open | Double-charge incidents | S | hard (revert) |

**Recommendation: Option 1 (block merge)**. Cites F1, F2, C1, C2.

User reviews and types `proceed`.

## Step 9 — Compact (inline)

`07-final.md` (~2200 words) + `summary.md` (~180 words). Final report includes the recommendation and risk table.

## Final state

```
.claude/pipelines/payment-retry-audit-9b1d3e/
├── manifest.json        # all phases status=done
├── 01-plan.md
├── 02-evidence/q1.md ... q5.md
├── 03-analysis.md
├── 04-review.md
├── 05-validation.md
├── 06-decision.md
├── 07-final.md
└── summary.md
```

The user can attach the directory to the PR for reviewer hand-off.
