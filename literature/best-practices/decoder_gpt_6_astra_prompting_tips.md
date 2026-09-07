# OpenAI shares prompting tips for GPT-6 Astra including a blocklist of slop words

Source: https://the-decoder.com/openai-shares-prompting-tips-for-gpt-6-astra-including-a-blocklist-of-slop-words/
Accessed: 2026-09-07
Category: best-practices / GPT-6 Astra prompting tips (press)
Provenance: new capture 2026-09-07 (frontier-model cohort: Claude Fable 5.1 / GPT-6 Astra generation); normalized from raw HTML capture (Tier 2 secondary source)

## Why This Source Matters

Trade-press summary of OpenAI's GPT-6 Astra prompting guide: initiative, instruction priority over skill files, subagent delegation, testing scope, and the slop-word list. Tier 2 secondary source; useful as an independent reading of what OpenAI emphasised, not as ground truth.


# OpenAI shares prompting tips for GPT-6 Astra including a blocklist of slop words

[https://the-decoder.com/author/matthias-bastian/](https://the-decoder.com/author/matthias-bastian/)

[Matthias Bastian](https://the-decoder.com/author/matthias-bastian/) [**View the LinkedIn Profile of Matthias Bastian](https://www.linkedin.com/in/matthias-bastian-128b71b1/)

 Sep 5, 2026

GPT-Image-2 prompted by THE DECODER

**OpenAI's model documentation spells out where GPT-6 Astra tends toward unwanted behavior and how developers can work around it.**

GPT-6 Astra asks clarifying questions more often than GPT-5.6 Sol instead of making assumptions on its own, according to OpenAI, making it a "more effective collaborator." The trade-off is that the model sometimes stops where users expect it to keep going.

To push it toward more initiative, OpenAI recommends a prompt telling the model to infer the user's "intent" from context and show a "bias towards action." Phrases like "can you...," "I want to...," or "help me..." should be treated as calls to act, not invitations for follow-up questions.

Ad

`You should infer the user’s intent and the scope of the task from the instructions and the conversation context so far. Your task is to demonstrate a tendency to act and to follow through on the user’s intended task until completion. If the user expresses the intention to complete new work or resolve an existing problem, continue working persistently until the user’s intended goal is achieved. Work independently toward the user’s goal (e.g., create isolated work trees/checkouts, resolve merge conflicts, perform read-only actions, create draft PRs, etc.), unless the actions are clearly destructive or irreversible.`

Ad

The model should wait to ask for approval until it has already prepared a concrete, reviewable result. OpenAI prompts it this way: "The user should be approving a concrete, reviewable result." Unsolicited warnings, disclaimers, or safety checklists based on hypothetical risks should be dropped from prompts.

GPT-6 Astra follows longer instructions better than its predecessors but is also more sensitive to context. Unclear or contradictory instructions in skill files like `AGENTS.md` can cause the model to block work or veer off unexpectedly. OpenAI recommends auditing all skill files and context documents the model can access and giving user instructions explicit priority.

Ad

OpenAI also recommends a debugging prompt that forces the model to name the exact skill file and quote the specific instruction that caused it to pause or change direction. This helps developers trace unexpected behavior back to its source.

`If a skill causes you to ask for permission or confirmation, pause, leave requested work unfinished, or diverge from the user's intent, name and link to the exact SKILL.md file you read, quote the relevant instruction, and briefly explain how it applies. Distinguish explicit skill requirements from your interpretation of guidelines.`

Ad

## "Delve," "foster," "leverage" and the rest of the slop words

GPT-6 Astra tends to structure responses with lists, tables, and Markdown formatting, and reuses the same phrases across sessions. OpenAI has specific guidance on shaping the model's writing style. If you want prose, tell the model explicitly to write concise paragraphs using plain language and active voice.

Ad

`By default, use clear, concise paragraphs, each developing a single main idea. Use lists only if the information is truly parallel, sequential, or more easily comparable, and avoid nested lists unless the hierarchy cannot be clearly expressed in prose. Use simple, straightforward language: familiar words, concrete examples, and precise verbs. Favor the active voice and direct statements. Make the main point clear early on, then expand on it with the explanation and details the reader needs. Let each sentence build on the previous one. Develop the points that are important and provide enough evidence to be useful.`

A blocklist of typical AI phrases can help too. OpenAI calls them "slop words." Made-up hyphenated compounds like "exact-head checks" or "editorial-row layouts" should also be avoided. The model should state what it's doing rather than listing what it won't do.

`Avoid using slop words or phrases such as “Conclusion:” in conclusions, “delve into,” “promote,” “use/leverage,” “it’s worth noting,” “what’s important is,” “Question? Answer,” or “This isn’t about X. It’s about Y,” “really/truly,” or compound descriptions and hyphenated adjectives. Do not use summary closing statements such as “In short:...,” “The simplest mental model is:...”. State the intended action directly. Avoid mentioning what you won’t do, what remains unchanged, or how you’ll separate or categorize results. Do not use contrastive phrasing such as “X, not Y” or “X—not Y,” which introduces an unsolicited alternative that the user did not ask for. Avoid made-up compound terms like “exact-head checks” and “editorial-row layouts,” vague qualifiers, and stock transitions; use simple verbs and prepositions to directly express the actual relationship.`

For technical writing, OpenAI recommends keeping jargon to cases where it actually helps: "Use plain language over jargon, and reference technical details only to the degree that it helps illustrate an idea or your work to the user."

## Sub-agents don't delegate enough, tests balloon

The model can hand off work to [sub-agents](https://developers.openai.com/api/docs/guides/responses-multi-agent) running in parallel but does so less often than expected. Developers should spell out when and how much it should delegate, OpenAI says. Messages between agents can also contain grammar or spacing errors.

On coding tasks, GPT-6 Astra runs thorough tests before wrapping up. For small changes, that can mean test suites wildly out of proportion to the actual work. OpenAI recommends telling the model to rerun tests only when new failures or unresolved issues justify it.

More detailed versions of these prompts are available on the [GPT-6 Astra model documentation page](https://developers.openai.com/api/docs/guides/latest-model). Developers who want to switch to GPT-6 Astra can use Codex with the [OpenAI Docs skill](https://github.com/openai/skills/tree/main/skills/.curated/openai-docs) to apply the recommended changes automatically: `$openai-docs migrate this project to GPT-6 Astra`

### AI News Without the Hype – Curated by Humans

 Subscribe to THE DECODER for ad-free reading, a weekly AI newsletter, our exclusive "AI Radar" frontier report six times a year, full archive access, and access to our comment section.

[Subscribe now](https://the-decoder.com/subscription/)

Source: [OpenAI](https://developers.openai.com/api/docs/guides/latest-model)
