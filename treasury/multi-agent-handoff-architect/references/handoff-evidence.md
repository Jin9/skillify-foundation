# Handoff evidence

## Note on the "17.2x error amplification" framing

The 17.2x figure is carried from the originating brief and is not
located in the cited source reports; the corroborating, source-grounded
evidence is below:

- Multi-agent systems fail at **41–87%** in production when handoffs are
  ad hoc and under-specified.
- A single agent succeeded **28 of 28** times on its benchmark, while
  hierarchical multi-agent organizations failed **36%** of the time and
  self-organized swarms failed **68%** of the time.
- On parallelizable tasks (Finance-Agent benchmark), multi-agent
  coordination produced **+81%** improvement over single-agent; on
  strictly sequential reasoning tasks, multi-agent systems degraded
  performance by **39–70%**.
- Multi-agent systems consumed **4–220x** more tokens than single-agent
  equivalents.
- Chained-agent reliability multiplies: two sequential agents at 95%
  individual reliability give **90.25%**; three give **85.7%**; five
  give **~77%**.
- Structured context objects are typically **200–500 tokens** versus
  **5,000–20,000** for full forwarding.
- Summarized context reduces token count **70–90%** but introduces
  information loss and adds **500ms–1.5s** latency per handoff.
- The orchestrator-worker topology outperformed the same-model
  single-agent baseline by **more than 90%** on Anthropic's research
  evaluation, at roughly **15x** more tokens than a single chat.
- A four-role coding squad (manager, researcher, engineer, reviewer)
  reached **72.2%** on SWE-bench Verified — a **+7.2-point** gain over
  the same-model single-agent baseline attributed to team structure.

This keeps the skill internally honest while preserving the requested
headline. The "17.2x error amplification seen in unstructured agent
loops" phrasing is retained verbatim in `SKILL.md`
(Purpose/description/Constraints framing) as the headline mechanism, per
explicit user decision; this reference is the source of truth for what
is and is not grounded in the cited reports.

## Source reports

- `agent-handoff-protocols.md` — definitions, mechanisms, framework
  landscape, empirical evidence, failure modes, the four principles,
  practical steps.
- `multi-agent-squad-architecture.md` — topologies, the four-step
  decision rule, role-specialization evidence, MAST failure taxonomy,
  token-cost amplification, the single-agent counter-position.
