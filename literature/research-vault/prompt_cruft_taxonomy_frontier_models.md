# Prompt and Skill Cruft Taxonomy for Frontier-Model Generations (2026-09 synthesis)

Source: Attributed synthesis of Claude Code 2.1.263 bundled `claude-api` skill, files `shared/prompt-audit.md` (Steps 0-7, Groups 1a-1f, 2, 3, 4, keep list) and `shared/agent-design.md` (model parameters, tool surface, scaling, long-running context, caching); local bundled files, not public web pages, so no verbatim copy is redistributed here
Accessed: 2026-09-07
Category: research-vault / dated-pattern taxonomy for skills
Provenance: written 2026-09-07 by the orchestrating agent from the two bundled files above, cross-checked against the public Fable 5.1 and GPT-6 Astra prompting guides captured the same day; section references point at the bundled file, quotations are limited to short phrases

## Why This Source Matters

The bundled audit guide is the only source in the corpus that treats a skill file as a per-model artifact and gives a greppable taxonomy of instructions that helped an older model and now harm a newer one. Its keep list is equally important: it names what an audit must NOT remove. Together they ground skillify's re-baseline sub-flow, its `cruft_scan.py` signals, and the rewritten anti-patterns 5, 6 and 12. Internal-synthesis tier: cite the public vendor guides for any claim that needs an external reference.

**Depth:** Deep | **Audience:** Skill authors and auditors | **Date:** September 2026

## 1. Framing: cruft is relative to a model, not to length

The guide's central move (Step 0, Step 3) is that a prompt or skill instruction is "cruft" only relative to a target model. Text added because an earlier model under-triggered, planned poorly, or rambled is dead weight once the target model does that thing natively, and the leftover text is not merely wasted tokens: over-emphasis produces over-triggering and rigid behaviour, and prohibitions can anchor the model toward the failure they name. The audit therefore looks for specific dated instructions, never for shortness. Its own summary of the frame is that every token must earn its place, which is different from "make it short".

Two classification questions drive everything else. First (Step 3): could the model already know this? Audience, product, environment facts, the quality bar, tool contracts and the reasons behind constraints are context, and context is never cruft. Restatements of trained defaults, behaviour the model already does unprompted, and workarounds for retired failures are the removal candidates. Second: is the line a constraint on behaviour (test it) or context the model cannot get elsewhere (usually keep)? This second check is what stops an audit from degenerating into a length contest.

Provenance matters (Step 2). Where history exists, the audit asks of every emphatic or prohibitive line which failure, on which model, it prevented, and whether that failure still reproduces on the target. Lines whose reason nobody can state are suspect by default. Idiom alone (scratchpad tags, "think step by step", role-context-rules-examples boilerplate) is a low-confidence signal; it earns medium or high confidence only when paired with a reason grounded in the target model's documented behaviour.

## 2. Group 1: dated prompt text

### 1a. Pressure language

Older models needed forcefulness; current models are highly responsive to the system prompt, so the same text over-applies, in both directions: capitalised MUST/NEVER/CRITICAL walls cause over-triggering and rigidity, while leftover hedges ("try to", "if possible") are now read as permission to under-deliver. When several instructions are each marked critical the markers stop carrying information, and the prompt's anxious register becomes the output's register. Emphasis is not banned; it is a scoped, tested fix for one demonstrably under-weighted instruction. Greppable signals: density of MUST/NEVER/ALWAYS/CRITICAL/IMPORTANT in caps, double exclamation marks, emphasis with no adjacent reason, hedges attached to real requirements, "you tend to" trait claims, "don't be too X".

### 1b. Scaffolds replaced by API features

These are swapped, not toned down. "Think step by step" and scratchpad or thinking-tag instructions are replaced by adaptive thinking plus an effort setting. "Plan before acting" is deleted because current models plan without being told and the instruction causes over-planning; if behaviour is still too aggressive, lower effort rather than adding prose. "Show your thinking" or required reasoning sections in the output are replaced by reading thinking blocks through the API, and on the newest models instructing reasoning reproduction can trigger a refusal. Assistant-turn prefills and the JSON-forcing stack around them (stop sequences, regex extraction, retry-on-parse loops, "output ONLY valid JSON") are replaced by structured outputs, and the surrounding code is cruft too. Fixed narration cadences ("summarize every N tool calls") and numeric word caps are deleted and re-baselined, because output caps starve reasoning on hard problems; prefer qualitative length guidance. Inline lookup tables and arithmetic rubrics move to files or code. Forced tool choice becomes a prompt instruction under automatic tool choice, with strict schemas for argument validity.

### 1c. Over-specification

Describe the goal, not the method. Step-by-step choreography for judgment tasks is the headline pattern: skills written for prior models are often too prescriptive for current ones and degrade output, because the model's own plan usually beats a hand-written script. Keep numbered steps only where order truly matters. Prohibition lists are rewritten as positive statements of intent unless the failure reproduces on the target. A single gold example freezes an older model's length, tone and structure into the new one; use several deliberately varied examples labelled illustrative, and keep only examples that pin a genuinely format-sensitive output. Bullet walls for behavioural guidance flatten priority and sever rules from reasons; use prose for behaviour and structure for reference data. Padding (generic virtues, repetition as reinforcement, kitchen-sink edge cases) is applied as actionable signal where it does not fit and inflates thinking spend. Grader vocabulary ("you will be graded on") pushes effort toward being watched; state the requirement instead. Strategy coaching ("it's usually best to") is deleted when removing it would change neither what is legal nor how success is measured.

### 1d. Fossils

Text that outlived its model: version-specific workarounds and "known issue with model X" comments; migration-relative phrasing ("now works differently", "no longer") that describes a diff against a prompt version the model never saw; patch accretion, where many narrow conditionals each trace to one incident and the model navigates a maze instead of a principle; unenforced instructions that nothing checks and nobody misses; identity stubs standing in for real context; update suppressors ("hold all findings for the final response", "don't narrate") written against chatty models, which now make under-narrating models go silent; anti-formatting rules ("never use bullets") written against over-formatting models, which now strip structure the reader wanted; and instruction re-insertion on a cadence, a retention crutch for models that lost instructions over long sessions.

### 1e. Prohibition clusters, judged by provenance

A run of unconditional "never / don't / must not" lines is audited line by line by asking whether each carries a stated reason or encodes a real business or policy constraint, not whether the target model "still needs the guardrail" (that question keeps everything). Refund caps, data rules, compliance language and promises the business must not make stay, ideally with the reason beside them. Banned-phrase lists and tic lists written against an older model's habits are cruft: restate the desired style positively in one line. A cluster of legitimate reasoned prohibitions does not launder the unreasoned ones mixed into it.

### 1f. Output-shaping choreography

Interim-update cadences, numeric output ceilings and cut-the-detail instructions are one pattern and are removed together. A stated operational reason does not convert a numeric clamp into a keeper: re-express the goal as audience or outcome framing and keep genuinely format-sensitive requirements as format instructions, not word counts.

## 3. Group 2: brittle skill files

Skill files inherit everything above plus their own failure modes, and skill size is a tax paid on every trigger. The named patterns: a verbose SKILL.md explaining what the model already knows (apply the deletion rule paragraph by paragraph); wrong degrees of freedom (exact scripts for judgment calls over-constrain, vague prose for fragile operations under-constrains; match specificity to fragility, with "do not modify this command" reserved for narrow bridges); the recency trap, where one session's stumble becomes a permanent rule that later sessions step around for no reason; volatile specifics such as hardcoded paths, flags, version numbers and API claims with no verification date, which rot as code ships; time-sensitive content, option menus and information duplicated across SKILL.md and references, which drift apart; history narratives (past tense, incident IDs, PR numbers, pinned model names) whose authority is the incident rather than the behaviour; and trigger-case enumeration, where a description grows one near-synonymous phrase per missed trigger, taxes every request and generalises worse than intent categories.

## 4. Group 3: tool descriptions

The rubric here is precision and contract accuracy, not brevity, and the most common failure is under-description. A tool description is a man page: what the tool does, when to use it and when not to, what each parameter means, caveats, what it does not return. What changed on current models is which content belongs there: contract and mechanics in; behavioural steering, worked examples, fake dialogue and embedded protocols out (moved to skills and progressive disclosure). Scolding cross-references and behaviour-smuggling do not belong in a description. The system prompt should not name tools. Past a few dozen tools, use tool search and deferred loading rather than always-loading every schema. One deliberate split: text whose job is routing (a skill's frontmatter description, a trigger block) may legitimately carry calibrated urgency because skills currently under-trigger; text whose job is behaviour should explain rather than shout. The two look identical to a grep, so classify by function before flagging.

## 5. Group 4: request configuration and architecture

Reported alongside prompt cruft even though they are not prompt text: API fossils (parameters and headers that error on the target); cache-hostile ordering (volatile content above stable content); budget countdowns rendered into context, which cause premature wrap-up; an LLM executor for a deterministic plan (count the model-call sites and move routing, tallying, normalising and formatting back into code, keeping exactly one call where judgment is real); redundant specialist sub-agents that differ only in a filter or payload field; and the absence of token accounting, which makes every other issue invisible.

## 6. The keep list

The guide is explicit that an audit which only says "delete" hurts the users who follow it most diligently. These stay even when a grep matches:

1. Context is never cruft: audience, product, environment facts, quality bar, constraints and their reasons. Too-short prompts produce generic output because the model fills gaps with safe defaults.
2. Cruft is not length. Never justify a deletion by character count alone.
3. Fragile operations keep exact scripts. Low-freedom prescriptive text is correct where exactly one sequence is safe (destructive commands, auth flows, compliance steps).
4. Tool contract detail stays and often grows.
5. Prohibitions against current, demonstrated failures stay. The discriminator is whether the failure reproduces on the target model in this context.
6. Trigger and routing text may carry calibrated urgency.
7. Format-pinning examples on genuinely format-sensitive outputs stay, labelled illustrative.
8. Working redundancy is not cruft; propose consolidation only when duplicates actually disagree. An audit that finds nothing should change nothing.
9. A one-line role statement is fine; flag identity text only when it substitutes for real context.
10. One deliberate end-of-prompt recap is a reasonable pattern; scattered duplication is the anti-pattern.
11. Re-baselining adds text too: matching a prompt to a new model sometimes means adding guidance for the new model's failure modes. The audit's job is fit, in both directions.

## 7. Method: report, diff, verify

The audit is non-interactive by design and always produces two artifacts: a report (location, quoted evidence, matched pattern, why obsolete for the target model, confidence, proposed action) and a proposed diff with one finding per hunk so effects attribute and hunks can be taken selectively. Confidence is High when documented in current vendor docs or when the pattern errors on the target model, Medium for consistent widely observed behaviour, Low for heuristic or idiom dating (flag, do not edit). A documented-pattern match gets a concrete proposed action; "flag" is reserved for undocumented low-confidence items and out-of-scope items. Rewrites beat bare deletions where the instruction has a live purpose. A removal is complete only when everything referencing it goes too. Removal is a hypothesis: probe behaviour before and after on a scratch copy rather than asking the model whether it needs an instruction, change one thing at a time where stakes are high, re-add in minimal form if a cut regresses, and re-audit at every model release because each new migration section is the trigger.

## 8. Agent-design notes that bear on skills

From the companion agent-design file: adaptive thinking plus effort replaces thinking budgets, and lower effort means fewer and more consolidated tool calls, less preamble and terser confirmations, with medium often a favourable balance and max reserved for correctness over cost. Promote an action from bash to a dedicated tool when it needs gating, staleness checks, rendering or parallel scheduling; start with bash for breadth. Tool search appends schemas rather than swapping them, which preserves the prompt cache, and skills keep task-specific instructions out of the fixed context until they are relevant: both are progressive-disclosure mechanisms. For long runs, context editing prunes stale tool results, compaction summarises near the limit, and memory persists across sessions; many agents use all three. Editing the system prompt or the tool list mid-session invalidates the cache; append a mid-conversation system message or spawn a sub-agent instead.

## 9. Implications recorded for skillify (2026-09-07)

- Anti-pattern 5 becomes two-directional (wrong degree of freedom in either direction), anti-pattern 6 absorbs instruction priority and unrequested stops, anti-pattern 12 absorbs model fossils and dated scaffolds.
- The rubric's workflow dimension rewards specificity matched to fragility rather than numbered steps as such; the token-efficiency dimension names scaffolds and repeated reminders explicitly.
- `scripts/cruft_scan.py` encodes the greppable signals with the keep-list exemptions (trigger sections, table rows, fenced code) and stays advisory, because regex hits are hypotheses.
- The Re-baseline sub-flow of Refactor follows the guide's Step 0 to Step 7 order: target generation and scope first, inventory, provenance, classify, scan, report plus diff, verify.
