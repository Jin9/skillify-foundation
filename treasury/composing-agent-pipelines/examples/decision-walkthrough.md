# Walkthrough: Decision domain

A worked end-to-end run of the `decision` phase shape:
**Plan → Gather → Analyze → Decide → Validate → Compact** (Review skipped).

Note that Validate runs *after* Decide here — to fact-check the recommendation's claims.

## User prompt

> "We need to choose a queue technology for the new event bus: Kafka, NATS JetStream, or RabbitMQ. Help me decide based on what our infra team can support and what the workload needs."

## Step 1 — Universal preamble

- Mode: `Compose`
- Domain: `decision` (matches "choose", "decide", "compare options")
- Phase shape printed: `plan -> gather -> analyze -> decide -> validate -> compact`
- User replies `proceed`.

## Step 2 — Initialize state

```
$ python3 scripts/init_pipeline.py --slug event-bus-choice --domain decision --prompt "Choose between Kafka, NATS JetStream, RabbitMQ..."
created: <cwd>/.agent-pipelines/event-bus-choice-3a8f51
```

## Step 3 — Plan (hard gate)

```
1. q1: What does our existing infra (k8s ops, monitoring, runbooks) already support?
2. q2: What are the throughput, latency, and ordering requirements of the new workload?
3. q3: What's the operational cost / on-call burden of each candidate at our scale?
4. q4: What's the language/SDK fit for our Go services?
```

User confirms `proceed`.

## Step 4 — Gather (4 parallel writer workers)

Workers search the repo (existing docker-compose, Helm values, runbooks) and any docs/. Each writes `02-evidence/qN.md`.

## Step 5 — Analyze

One writer worker produces `03-analysis.md`:
- F1 (P1): Infra team already runs Kafka in 2 other services. Cites `02-evidence/q1.md#F1`.
- F2 (P1): Workload needs partitioned ordering at 10k msg/s. Cites `q2.md#F1`.
- F3 (P2): Go SDK maturity: Kafka (Sarama/segmentio) > NATS > RabbitMQ. Cites `q4.md#F1, F2, F3`.
- F4 (P2): NATS would add a new operational surface (no on-call expertise). Cites `q3.md#F2`.

## Step 6 — Decide (hard gate, no Review pre-pass)

One writer worker reads `03-analysis.md` directly (no `04-review.md` exists in this shape) and produces `06-decision.md`:

| # | Option | Pros | Cons | Risks | Effort | Reversibility |
|---|--------|------|------|-------|--------|---------------|
| 1 | Kafka | Existing infra, ordered, mature Go SDK | Heavier than needed for some topics | Increased coupling to one tech | M | hard |
| 2 | NATS JetStream | Lightweight, good Go SDK | New on-call surface | Operational burden | M | medium |
| 3 | RabbitMQ | Battle-tested | Weaker partitioned-ordering, weaker Go SDK | Requirement misfit | M | medium |

**Recommendation: Option 1 (Kafka)**. Rationale cites F1, F2, F3, F4.

User reviews and types `proceed`.

## Step 7 — Validate

The recommendation rests on F1 (existing Kafka infra) and F2 (10k msg/s ordering). The orchestrator delegates one sequential validation pass to re-verify those claims:

- C1 (F1: Kafka in 2 other services): confirmed — found in `helm/charts/billing/values.yaml` and `helm/charts/notifications/values.yaml`.
- C2 (F2: 10k msg/s ordering needed): unverifiable — found a design doc but the throughput number is from a slack thread not in the repo.

Orchestrator surfaces C2 (unverifiable) to user. User confirms the throughput number from the slack thread, accepts the validation gap, and proceeds.

## Step 8 — Compact (inline)

`07-final.md` (~2000 words) includes the options table, recommendation, and the validation footnote about C2. `summary.md` (~150 words) leads with the recommendation.

## Final state

```
.agent-pipelines/event-bus-choice-3a8f51/
├── manifest.json        # review status=skipped, all others=done
├── 01-plan.md
├── 02-evidence/q1.md ... q4.md
├── 03-analysis.md
├── 06-decision.md
├── 05-validation.md
├── 07-final.md
└── summary.md
```

## Why this shape

- **Review skipped** because Decide's options-comparison already exposes the trade-offs an adversarial reviewer would surface. Adding standalone Review would duplicate that work.
- **Validate kept** to fact-check the recommendation's load-bearing claims before the user acts on them. C2's "unverifiable" verdict is the value Validate adds — without it, the recommendation might rest on an unsourced number.
