# Effort and cost tiers

Match the model and effort to the job. Over-powering a trivial consult wastes tokens and time;
under-powering a hard one produces advice not worth adjudicating.

## Cost tiers

Tiers are the stable abstraction; the identifiers filling them are not. Route by tier, then resolve
the tier to whatever the CLI currently offers.

| Tier | Use for | How to fill it |
|------|---------|----------------|
| Small | mechanical extraction, formatting, quick checks | the cheapest fast model the CLI lists, at its lowest effort variant; headless `haiku` |
| Mid | structuring, summarization, routine code passes | the fast model at its highest effort variant; the account's standard codex model; headless `haiku`/`opus` |
| Frontier | deep reasoning, ambiguous design, hard judgment | the reasoning-class model at high effort; headless `opus` |

**Concrete on this setup, verified 2026-07-30 — an example, not a contract:**

| Tier | agy | codex | headless |
|------|-----|-------|----------|
| Small | `Gemini 3.6 Flash (Low)` / `(Medium)` | — | `haiku` |
| Mid | `Gemini 3.6 Flash (High)` | `gpt-5.6-sol` | `haiku` / `opus` |
| Frontier | `Gemini 3.1 Pro (High)` | `gpt-5.6-sol` | `opus` |

Two things to know before copying those values:

- **`Gemini 3.1 Pro` has no Medium variant** — only High and Low.
- **Use the display-label form for agy.** The equivalent slug `gemini-3.1-pro-high` is accepted
  without error but silently runs the default model instead. See `cli-invocation-contracts.md`.

## Per-CLI knobs

- **agy:** effort is part of the model identifier — `(Low)`, `(Medium)`, `(High)`. A separate
  `--effort` flag exists, but it *conflicts* with an identifier that already encodes effort and is
  rejected. Choose effort by choosing the model.
- **codex:** the account's configured model in `~/.codex/config.toml` sets capability, and
  `model_reasoning_effort` sets depth. Prefer changing the config over pinning `-m` per call.
- **headless host agent:** `--effort low` is mandatory for reliability (see
  `cli-invocation-contracts.md`); raise it only with a longer watchdog and a reason.

## Routing heuristics

- Route routine or cheap work to Small/Mid; reserve Frontier for tasks whose answer the agent would
  otherwise struggle to produce or verify.
- For fan-out, mixing one Mid and one Frontier model often catches more than two Frontier calls at
  higher cost.
- Set the wrapper `TIMEOUT` to the tier: Frontier reasoning over many files needs minutes; a Small
  extraction should finish fast, and a long hang signals a problem, not slow thinking. Remember the
  CLI's own print-timeout bounds the run too — a watchdog alone does not control how long you wait.
- A cheap tier that silently answered from the wrong model is not cheap, it is worthless. Check
  `backend_verified=yes` in the provenance line before you spend adjudication effort on the output.
