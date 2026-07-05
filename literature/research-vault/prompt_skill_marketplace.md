# Prompt/skill marketplace

Source: ResearchVault run prompt-skill-marketplace-20260526-050145 (local deep-research pipeline)
Accessed: 2026-05-26
Category: research-vault / marketplace landscape
Provenance: harvested 2026-07-05 from ResearchVault 05-final_report.md (run executed 2026-05-26); inline bibliography preserved

## Why This Source Matters

Internal deep-research synthesis of the prompt/skill marketplace landscape — hubs, monetization, and supply-chain implications. Complements the skills.sh and gh-skill sources. Internal-synthesis tier.


## Executive Summary

A "prompt/skill marketplace" is a platform where AI prompts and, increasingly, executable agent skills are listed, discovered, and monetized. The category leader, PromptBase, proved since 2022 that a well-tuned prompt can be a sellable digital good, taking a 20% commission and paying sellers 80% on items priced roughly $1.99–$9.99 [1][2]. That original thesis — selling strings of text — is now under pressure from two directions at once: prompts are commoditizing as models get better at understanding intent [3], and the buyer base openly questions paying for prompts available free elsewhere [4]. For an executive or founder, the honest read is that the *plain-prompt* business is a thin, defensibility-poor niche, while the *executable-skill* opportunity riding the AI-agent wave is large and fast-growing but increasingly owned by the model platforms themselves.

The market math is favorable in aggregate and treacherous in detail. The narrow prompt-engineering market is small — about $222 million in 2023 growing to roughly $2 billion by 2030 [5] — whereas broader agent-tooling and AI-agent markets are projected at tens of billions by 2030 at 40%+ CAGRs [6][7]. The executable-skill ecosystem is the part actually scaling: over 10,000 active Model Context Protocol (MCP) servers were in production by December 2025, and skills became a cross-vendor open standard when OpenAI adopted Anthropic's SKILL.md format [8][9]. But the most powerful distribution surfaces — OpenAI's ChatGPT app store with ~800 million weekly users and Anthropic's curated partner skills directory — are platform-owned, and their creator monetization remains weak today [10][11].

The decision rule for a founder is therefore narrow but real: a horizontal "sell prompts to everyone" marketplace is a poor bet, while a vetted, vertical, skill-and-workflow marketplace with curation as a moat — built deliberately on the portable open skill standard rather than captive to one model vendor — is the defensible wedge worth pursuing. Trust and security are not a footnote here: independent studies found roughly a quarter to a third of agent skills carry vulnerabilities, and active attacks have already poisoned marketplaces at scale [12][13]. Curation that buyers can trust is simultaneously the largest cost and the clearest moat.

## Background

The question this report answers for an executive or founder is whether there is a defensible, venture-scale business in a marketplace for AI prompts and agent skills, and what determines whether such a marketplace wins, stalls, or gets commoditized. The topic matters now because the underlying artifact is changing shape: the market began with static prompts (sellable text) and is migrating to executable skills (packaged instructions plus scripts and resources that an agent runs). The strategic and economic implications of those two artifact classes differ sharply, and conflating them is the most common error in evaluating the opportunity.

## Methodology

This report synthesizes published industry analysis, primary platform documentation, market-research forecasts, and security research gathered against a nine-part research plan spanning category definition, competitive landscape, market sizing, business model and unit economics, demand and willingness-to-pay, moat and commoditization risk, the prompt-to-skill shift, governance and trust risk, and the founder's decision rule. Where reputable market-sizing sources disagree by an order of magnitude, the report surfaces the disagreement and attributes it to differing category definitions rather than averaging the estimates. Forward-looking forecasts are vendor projections and should be read as directional, not as audited figures.

## Key Findings

- The category is bifurcating into static-prompt marketplaces and executable-skill marketplaces, and the two have very different defensibility profiles [1][14].
- The largest distribution surfaces are platform-owned (GPT Store, ChatGPT app store, Anthropic's partner skills directory), not independent marketplaces [15][10][16].
- Market-sizing estimates diverge by definition: ~$2B by 2030 for narrow prompt engineering versus tens of billions for agent tooling and AI agents [5][6][7].
- Sustainable marketplace take-rates cluster near 20–30% but face downward pressure, and low absolute prompt prices cap gross merchandise value per transaction [2][17][18].
- Demand for AI skills is real and growing, but willingness-to-pay specifically for prompts is contested, and platform-store creator payouts are currently weak [19][4][11].
- Executable skills are scaling fast and have become a cross-vendor open standard, enlarging the market while weakening any single platform's lock-in [8][20][9].
- Trust and security are first-order: roughly a quarter to a third of agent skills carry vulnerabilities, and marketplaces have already been attacked at scale [12][13].

## What the category actually is — and why the boundary matters

A prompt marketplace, in its original form, sells individual prompts as licensed digital goods. PromptBase is the canonical example: creators list text and image prompts, buyers license them per use, and the platform handles payments — a model it has run since 2022 [1]. The defining bet of that era was that "a well-tuned Midjourney or GPT prompt is itself a sellable digital good."

An agent skill is a materially different artifact. Skills are folders of instructions, scripts, and resources — packaged as SKILL.md files — that an AI loads dynamically to improve performance on specialized tasks [14]. Where a prompt is a single string a buyer pastes once, a skill is executable, reusable infrastructure that runs inside an agent. That difference is the whole strategic story: an executable workflow is harder to copy by eyeballing it, can encode proprietary logic and integrations, and is consumed repeatedly rather than once.

```
Static prompt                Executable skill
─────────────                ────────────────
text string         →        SKILL.md folder
  │                            │ instructions
  │ paste once                 │ + scripts
  ▼                            │ + resources
single output                  ▼
                          agent runs it
                          repeatedly, in-workflow
```

The boundary matters commercially because the two artifacts commoditize at different rates and monetize through different surfaces. Tellingly, the category is already converging from the prompt side: PromptBase now lets sellers offer agent skills (SKILL.md files) alongside text and image prompts, signaling that the prompt-only thesis is being abandoned even by its pioneer [1].

## The competitive landscape is platform-owned at the top

The most important structural fact for a new entrant is that the largest stores belong to the model platforms. OpenAI's GPT Store holds roughly 160,000 listed GPTs out of more than 3 million created — a curated subset, but still vastly larger than any independent prompt marketplace [15]. In December 2025 OpenAI opened third-party app submissions inside ChatGPT, giving developers access to roughly 800 million weekly active users through a directory where higher-quality apps are featured more prominently [16]. Anthropic, in parallel, launched a curated partner skills directory with commercial partners including Atlassian, Canva, Cloudflare, Figma, Notion, Ramp, and Sentry [14].

```
        Platform-owned (top of funnel)
        ┌─────────────────────────────────┐
        │ ChatGPT app store (~800M WAU)    │
        │ GPT Store (~160k listed)         │
        │ Anthropic partner skills dir.    │
        └─────────────────────────────────┘
                      ▲
                      │ squeeze
                      │
        ┌─────────────────────────────────┐
        │ Independent marketplaces         │
        │ (PromptBase, et al.)             │
        └─────────────────────────────────┘
```

This is the classic platform-squeeze pattern: the entity that owns the model and the distribution channel can absorb the most valuable parts of any third-party category. An independent marketplace competing head-on for horizontal prompt or skill distribution is competing against a directory attached to hundreds of millions of users. The viable independent plays are therefore the ones the platforms are *least* incentivized to serve well — vetted enterprise-grade curation, vertical depth, and cross-platform portability.

## Market size: large in aggregate, definition-dependent in detail

Market sizing for this category is unusually noisy, and an executive should treat the headline numbers with care because they measure different things. Grand View Research sizes the narrow global prompt-engineering market at about $222 million in 2023, reaching roughly $2.06 billion by 2030 at a 32.8% CAGR [5]. Mordor Intelligence, using a broader "prompt engineering and agent programming tools" definition, sizes the same conceptual space at about $6.95 billion in 2025 growing to $40.87 billion by 2030 at a 42.52% CAGR [6]. These two reputable firms differ by roughly an order of magnitude — a definitional artifact, not a contradiction to be averaged away.

The adjacency that actually bounds the upside is the AI-agents market, projected to grow from roughly $7.8 billion in 2025 to between $48 billion and $53 billion by 2030 at 43–46% CAGRs [7]. A skill marketplace is, in effect, a monetization layer on top of that agent economy. The takeaway for a founder: do not pitch "the prompt-engineering market" as the TAM — it is small and slow-monetizing. Pitch a defensible slice of the agent-skill economy, and be explicit about which definition the sizing rests on.

## Business model and unit economics: thin where it commoditizes

The prevailing marketplace model is a transaction take-rate. PromptBase runs a 20% commission, paying sellers 80%, on prompts typically priced $1.99–$9.99 [2]. That take-rate is reasonable, but the low absolute order values cap gross merchandise value per transaction — a marketplace clearing thousands of $5 prompts generates modest revenue, and customer-acquisition economics are correspondingly tight.

The benchmark ceiling for digital-goods marketplaces is around 30% — the norm across the major app stores — but that ceiling is eroding under competitive and regulatory pressure, with small-business and post-12-month subscription tiers already at 15% [17][18]. Meanwhile the dominant platform store has not yet built digital-goods monetization at all: OpenAI's Apps SDK currently steers developers to external checkout on their own domain rather than capturing an in-store take-rate, with digital-goods monetization still "being explored" [21]. The unit-economics implication is twofold: a take-rate above ~20–30% is hard to defend long-term, and the platforms have not yet locked in the economics — leaving a window, but not a wide one.

## Demand is real for skills, contested for prompts

There is genuine, growing demand for AI capabilities packaged for reuse. Custom-GPT adoption grew quickly — academic-use GPTs up 310% year-over-year, and GPT agents integrated across more than 10,000 SaaS platforms [19]. Usage is not the question.

Willingness to *pay specifically for prompts* is the question, and the evidence is mixed-to-skeptical. Many users openly ask why they should pay for prompts available free from Reddit and thousands of websites, and at least one prominent analysis concluded buyers should "stick with the free stuff" [4]. Worse for the monetization story, even where buyers do engage, the platform-store payout economics are weak: most GPT Store creators earn under $200/month, many earn nothing because they miss minimum-engagement thresholds, and payouts average around $0.03 per conversation — so successful creators bypass the store entirely to capture 95%+ of revenue through their own channels [11]. A marketplace whose best creators route around it has a structural problem.

## The moat is in skills and curation, not in prompts

Two forces are commoditizing plain prompts. First, the models themselves: prompt engineering had roughly an 18-month window as a genuine differentiator before models like GPT-4 Turbo, Claude 3, and Gemini 1.5 began understanding intent without elaborate prompting [3]. Second, free substitutes and disintermediation: PromptBase itself hosts over 2,300 free prompts, and the best platform-store creators bypass the store to keep nearly all their revenue [2][11]. The counter-view is that prompt/skill value is evolving rather than dying, migrating toward domain experts who inject specialized knowledge into AI systems [22] — but note that this counterpoint relocates the value *into vertical depth*, away from generic prompts.

```
Commoditization pressure        Defensibility migrates to
────────────────────────        ─────────────────────────
models understand intent   →    vertical domain depth
free prompts everywhere    →    executable + integrated skills
creators disintermediate   →    trusted curation / vetting
```

This is why the executable-skill shift deepens rather than erodes the opportunity — but only for the right operator. The skill ecosystem is scaling fast: over 10,000 active MCP servers were in production by December 2025 with 97M+ monthly SDK downloads [8], and the MCP registry grew roughly 7.8x year-over-year, from about 1,200 servers in Q1 2025 to 9,400+ by April 2026 [20]. Crucially, skills became a cross-vendor open standard in December 2025 when OpenAI adopted Anthropic's SKILL.md format for Codex CLI and ChatGPT [9]. That portability cuts both ways for a founder: it enlarges the addressable market and reduces lock-in to any single model vendor, but it also means a skill listed on your marketplace can be listed anywhere — so the durable moat is not the catalog, it is the trust, curation, and vertical integration layered on top.

## Governance and trust risk is a first-order liability

Any operator of a skill marketplace inherits a serious security obligation, because executable third-party skills are an attack surface. An empirical study of agent skills found 26.1% contained at least one vulnerability and 36.8% of the ecosystem carried at least one security flaw, spanning prompt injection, data exfiltration, privilege escalation, and supply-chain risks [12]. These are not theoretical: a Snyk study found prompt injection in 36% of skills, and a real campaign saw nearly 1,200 malicious skills infiltrate a major agent marketplace, exfiltrating API keys, cryptocurrency wallets, and browser credentials at scale [13]. The named mitigation is governance — open marketplaces let attackers publish "semantically compliant but logically harmful" skills, which is why security auditing and curation mechanisms are mandatory [13].

For an executive this reframes curation from a cost center into the core product. The platforms' own directories are curated for exactly this reason; an independent marketplace that can offer *deeper, vertical-specific, enterprise-grade* vetting than a horizontal platform store has both a liability shield and a differentiator. Trust is the moat that free prompts and open standards cannot erode.

## Synthesis: the founder's decision rule and viable wedges

Combining the landscape, economics, and risks yields a clear decision rule. A horizontal "sell prompts to everyone" marketplace is a poor venture bet: the artifact is commoditizing [3], willingness-to-pay is contested [4], take-rates are capped and eroding [17], and the largest distribution is platform-owned [16]. A vetted, vertical, skill-and-workflow marketplace is the defensible alternative, and several findings point to the same wedge. Marketplaces follow winner-take-most dynamics, so a sharp early wedge and rapid share capture matter [23]; niche customers prefer a single purpose-designed vendor, so vertical depth is where defensibility lives [24]; the open skill standard means a portable, multi-platform strategy is feasible rather than captive [9]; and trusted curation is both the regulatory-grade requirement and the moat [13].

The viable wedges that follow are: (1) a vertical skill marketplace for a specific industry where curation and compliance matter (the same logic that drives vertical-SaaS defensibility); (2) an enterprise-grade vetting-and-trust layer that sits across the open skill standard and the platform stores; and (3) a tools-and-workflow marketplace riding the MCP/agent economy rather than the prompt-engineering niche. In every case the durable advantage is curation and vertical depth, deliberately built on portable open standards, not a captive prompt catalog.

## Limitations & Open Questions

- All market-sizing figures are vendor forecasts with order-of-magnitude dispersion driven by category definitions; they are directional, not audited [5][6].
- Platform-store monetization is early and changing fast — OpenAI's digital-goods monetization is still being explored, so today's weak creator payouts may not be the steady state [21][11].
- The security-prevalence figures come from specific ecosystem snapshots; the exact rates will move as governance tooling matures, though the direction (real, material risk) is well-supported [12][13].
- This report does not assess any single region-specific marketplace or the deep technical packaging mechanics of MCP/skills beyond what bears on the marketplace opportunity.

## Sources

1. AI Prompts | PromptBase: The #1 Marketplace for AI Prompts — https://promptbase.com/
2. PromptBase Review 2026: Buy, Sell, Make Money AI Prompts — SoftHubTools — https://softhubtools.com/promptbase-review-2026-buy-sell-make-money-ai-prompts/
3. Prompt engineering is dead. Prompt engineering as a skill is dead... — Dmitry Kargaev, Medium — https://deeflect.medium.com/prompt-engineering-is-dead-295169b62a7a
4. Are Premium AI Prompts Worth the Money? — MakeUseOf — https://www.makeuseof.com/should-you-buy-ai-prompts/
5. Prompt Engineering Market Size And Share Report, 2030 — Grand View Research — https://www.grandviewresearch.com/industry-analysis/prompt-engineering-market-report
6. Prompt Engineering And Agent Programming Tools — Market Size, Share & 2030 Growth Trends Report — Mordor Intelligence — https://www.mordorintelligence.com/industry-reports/prompt-engineering-and-agent-programming-tools-market
7. AI Agents Market worth $52.62 billion by 2030 — MarketsandMarkets — https://www.marketsandmarkets.com/PressReleases/ai-agents.asp
8. One Year of MCP: November 2025 Spec Release — Model Context Protocol Blog — https://blog.modelcontextprotocol.io/posts/2025-11-25-first-mcp-anniversary/
9. OpenAI adopts Agent Skills, Anthropic donates MCP, GPT 5.2 and GPT Image 1.5 released — PulseMCP — https://www.pulsemcp.com/posts/openai-agent-skills-anthropic-donates-mcp-gpt-5-2-image-1-5
10. Developers can now submit apps to ChatGPT — OpenAI — https://openai.com/index/developers-can-now-submit-apps-to-chatgpt/
11. OpenAI GPT Store Revenue Sharing Explained: What Creators Actually Earn — The GPT Shop Blog — https://www.thegptshop.online/blog/openai-gpt-store-revenue-sharing
12. Agent Skills in the Wild: An Empirical Study of Security Vulnerabilities at Scale — arXiv — https://arxiv.org/pdf/2601.10338
13. Snyk Finds Prompt Injection in 36%, 1467 Malicious Payloads in a ToxicSkills Study of Agent Skills Supply Chain Compromise — Snyk — https://snyk.io/blog/toxicskills-malicious-ai-agent-skills-clawhub/
14. Agent Skills: Anthropic's Next Bid to Define AI Standards — The New Stack — https://thenewstack.io/agent-skills-anthropics-next-bid-to-define-ai-standards/
15. GPT Store Statistics & Facts: Contains 159,000 of the 3 million created GPTs — SEO.AI — https://seo.ai/blog/gpt-store-statistics-facts
16. Introducing apps in ChatGPT and the new Apps SDK — OpenAI — https://openai.com/index/introducing-apps-in-chatgpt/
17. Apple's App Store and Other Digital Marketplaces: A Comparison of Commission Rates — Analysis Group — https://www.analysisgroup.com/globalassets/insights/publishing/apples_app_store_and_other_digital_marketplaces_a_comparison_of_commission_rates.pdf
18. App Store Small Business Program: Guide for 2026 — Adapty — https://adapty.io/blog/app-store-small-business-program/
19. The Era of Tailored Intelligence: Charting the Growth and Market Impact of Custom GPTs — Originality.AI — https://originality.ai/blog/gpts-statistics
20. MCP Adoption Statistics 2026: Model Context Protocol — Digital Applied — https://www.digitalapplied.com/blog/mcp-adoption-statistics-2026-model-context-protocol
21. Monetization — Apps SDK | OpenAI Developers — https://developers.openai.com/apps-sdk/build/monetization
22. The Death of Prompt Engineering Has Been Greatly Exaggerated — PromptLayer — https://blog.promptlayer.com/the-death-of-prompt-engineering-has-been-greatly-exaggerated/
23. The Next 10 Years Will Be About the AI Agent Economy — NFX — https://www.nfx.com/post/ai-agent-marketplaces
24. Building Vertical AI: An early stage playbook for founders — Bessemer Venture Partners — https://www.bvp.com/atlas/building-vertical-ai-an-early-stage-playbook-for-founders
