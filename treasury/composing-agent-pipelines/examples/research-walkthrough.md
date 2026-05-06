# Walkthrough: Research domain

A worked end-to-end run of the `research` phase shape:
**Plan → Gather → Analyze → Review → Validate → Compact** (Decide skipped).

## User prompt

> "Research how the codebase handles auth token expiration. I want a memo I can share with the team — no recommendations needed, just findings."

## Step 1 — Universal preamble

The orchestrator detects:
- Mode: `Compose` (no specific mode hint)
- Domain: `research` (matches "research", "memo", "no recommendations")

It prints:

```
Domain: research
Phase shape: plan -> gather -> analyze -> review -> validate -> compact
Override? Reply with: 'shape: <phase1>,<phase2>,...' or 'proceed'.
```

User replies `proceed`.

## Step 2 — Initialize state

```
$ python3 scripts/init_pipeline.py --slug auth-token-expiry --domain research --prompt "Research how the codebase handles auth token expiration..."
created: <cwd>/.agent-pipelines/auth-token-expiry-7c4e2a
```

## Step 3 — Plan

The orchestrator delegates to one writer worker, or runs inline, with the prompt skeleton from `references/mode-playbooks.md#mode-plan`.

Result written to `01-plan.md`:

```
## Sub-questions
1. q1: Where is token expiration enforced in the auth middleware?
2. q2: How are tokens refreshed when expired?
3. q3: What logging/metrics exist around expiration events?
```

Orchestrator pauses at the hard gate. User reviews and types `proceed`.

## Step 4 — Gather (parallel)

The orchestrator delegates to 3 writer workers in one batch when the host supports parallelism. Each receives the prompt skeleton from `mode-playbooks.md#mode-gather` for its sub-question.

Results: `02-evidence/q1.md`, `q2.md`, `q3.md`. Each cites file paths with line numbers.

Soft gate: orchestrator confirms each file has ≥1 cited finding or "No evidence found", then proceeds.

## Step 5 — Analyze

One writer worker reads all evidence files and synthesizes `03-analysis.md` with findings tagged P1/P2/P3 and themes.

Soft gate. Proceed.

## Step 6 — Review (one pass, hard cap)

One writer worker with the devil's-advocate prompt writes `04-review.md`. Suppose it surfaces 1 P1 issue: "Analysis claims tokens never expire on the WebSocket channel, but only quotes the HTTP middleware — needs WebSocket-side evidence."

Orchestrator surfaces this to the user. User chooses to accept and continue (or to halt and re-run Gather with an additional sub-question — that would be a new pipeline run since Review is capped at 1 pass).

## Step 7 — Validate

The analysis has 4 P1/P2 claims. The orchestrator delegates one sequential validation pass, fact-checking each claim against the cited source.

Result `05-validation.md`:
- C1: confirmed
- C2: confirmed
- C3: unverifiable (cited file no longer exists)
- C4: confirmed

C3 is unverifiable but not refuted, so the pipeline proceeds.

## Step 8 — Compact (inline)

The orchestrator (main agent) reads everything and writes:
- `07-final.md`: 1500-word memo with TL;DR, sub-question answers, findings, validation summary, sources.
- `summary.md`: 180-word TL;DR.

## Final state

```
.agent-pipelines/auth-token-expiry-7c4e2a/
├── manifest.json        # all phases status=done, decide status=skipped
├── 01-plan.md
├── 02-evidence/
│   ├── q1.md
│   ├── q2.md
│   └── q3.md
├── 03-analysis.md
├── 04-review.md
├── 05-validation.md
├── 07-final.md
└── summary.md
```

## What the user sees

A short status table from `pipeline_status.py auth-token-expiry-7c4e2a`, plus the `summary.md` content inline in chat. They can read `07-final.md` for the full memo or `02-evidence/*.md` for raw findings.
