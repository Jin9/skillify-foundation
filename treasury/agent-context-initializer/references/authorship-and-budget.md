# Authorship and the line budget

## Human-curated vs LLM-generated: the evidence

The empirical record is consistent and the direction matters for this
skill's output contract:

- **Human-curated, minimal and precise: roughly +4%.** Human-curated
  AGENTS.md files provide a marginal **+4%** improvement in task success —
  *but only when they are minimal and precise*. A minimal, human-curated
  AGENTS.md of ~100-150 lines can reduce agent runtime by **28.6%**
  (28.64% lower median agent runtime) and token consumption by **16.6%**
  (16.58% lower output token consumption), per the 10-repository,
  124-pull-request study (arXiv 2601.20404).
- **LLM-generated: roughly -3%.** LLM-generated AGENTS.md files *reduce*
  task success rates by approximately **3%** on average (-0.5% on
  SWE-bench Lite, -2% on AGENTBENCH), increase inference cost by over
  20%, and require 2-4 additional reasoning steps (arXiv 2602.11988).
  LLM-generated context files universally triggered performance
  degradation in controlled evaluation.

The operational consequence: this skill's output is a **candidate**, never
a commit-ready file. Auto-generation is acceptable for scaffolding but a
human editorial pass before commit is mandatory — that pass is what
converts the -3% LLM-generated artifact into the +4% human-curated one,
and only if the result stays minimal and precise. `templates/curation-checklist.md`
exists to force that pass.

## The degradation band

Files exceeding roughly **60-200 lines** see degraded instruction-following
as important rules get buried in context; files over ~150 lines begin to
reverse the runtime/token gains entirely. Cited ceilings range from 60 to
200 lines across sources and model versions; no controlled study has
established a precise threshold, so this skill enforces the conservative
operational rule:

- **Hard cap: 200 lines.** `scripts/agents_md_gate.py` FAILs above 200.
- **Target: 100-150 lines.** The gate WARNs above 150.
- If the six sections plus boundaries cannot fit under the cap, the answer
  is never to raise the ceiling — it is to push more rules through the
  deferral filter into the curation checklist (`deferral-meta-principle.md`).

## Multi-agent context-duplication math

In agentic squad architectures, one orchestrator spawns multiple
sub-agents and **each sub-agent inherits a full copy** of the context file.
The cost is multiplicative: a 400-line file spawned across 7 simultaneous
agents is 7 x 400 = 2,800 lines of duplicated context competing directly
with the actual task payload. A 2,000-token file across 30 messages costs
60,000 tokens in context overhead alone. This is the strongest argument for
the line budget: every line saved is saved N times over, once per
sub-agent. Keep the root file operational (commands, constraints,
interfaces); push role-specific detail into the squad-roles block, not
into prose every agent re-loads.

## Progressive-disclosure rationale

AGENTS.md is always-on context — every token loads on every generation
step regardless of relevance. SKILL.md files, by contrast, load on demand
when a task matches their description. The optimization strategy: start
with a single file at the six sections, split into subdirectory AGENTS.md
files when the root exceeds 150-200 lines, and move detailed specialist
knowledge into on-demand skill files rather than inlining it. AGENTS.md
carries what every agent always needs; depth lives elsewhere and is
fetched only when required. This skill applies the same discipline to
itself — its own depth lives in `references/`, not in `SKILL.md`.

## Governance footnote

AGENTS.md was released by OpenAI in August 2025 and **donated to the Linux
Foundation's Agentic AI Foundation in December 2025**, adopted by **60,000+
open-source projects** and supported by 25+ tools. It is stable
infrastructure, not a vendor feature — which is why a minimal, portable,
human-curated file is worth getting right rather than auto-generating and
forgetting.
