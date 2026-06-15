# Effort and cost tiers

Match the model and effort to the job. Over-powering a trivial consult wastes tokens and time;
under-powering a hard one produces advice not worth adjudicating.

## Cost tiers

| Tier | Use for | Examples on this setup |
|------|---------|------------------------|
| Small | mechanical extraction, formatting, quick checks | agy `Gemini 3.5 Flash (Low/Medium)`, headless `haiku` |
| Mid | structuring, summarization, routine code passes | agy `Gemini 3.5 Flash (High)`, codex `gpt-5.5`, headless `haiku`/`opus` |
| Frontier | deep reasoning, ambiguous design, hard judgment | agy `Gemini 3.1 Pro (High)`, headless `opus`, codex `gpt-5.5` |

Verify the exact labels with `agy models`; the table lists defaults, not guarantees.

## Per-CLI knobs

- **agy:** effort is encoded in the label itself — `(Low)`, `(Medium)`, `(High)`. Pick the label;
  there is no separate effort flag.
- **codex:** `-m <model>` selects capability; `gpt-5.5` is the working tier on a ChatGPT account.
- **headless host agent:** `--effort low` is mandatory for reliability (see
  `cli-invocation-contracts.md`); raise it only with a longer watchdog and a reason.

## Routing heuristics

- Route routine or cheap work to Small/Mid; reserve Frontier for tasks whose answer the agent would
  otherwise struggle to produce or verify.
- For fan-out, mixing one Mid and one Frontier model often catches more than two Frontier calls at
  higher cost.
- Set the wrapper `TIMEOUT` to the tier: Frontier reasoning over many files needs minutes; a Small
  extraction should finish fast, and a long hang signals a problem, not slow thinking.
