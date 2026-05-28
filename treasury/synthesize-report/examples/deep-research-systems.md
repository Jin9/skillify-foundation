# Worked example: deep-research-systems

Example invocation of the `synthesize-report` skill at `standard` depth for a `general` audience. Demonstrates: depth-appropriate section set, audience-appropriate plain-English prose, inline `[n]` citations placed adjacent to claims, contradictions named in prose rather than laundered, a v0.2.0 ASCII diagram inside a body section after the introducing paragraph, and the `cited_sources` ↔ `## Sources` bijection.

## Inputs (abbreviated)

```yaml
topic: "What is the current state of LLM-based agentic research systems?"
research_plan:
  depth: standard
  sub_questions:
    - "How do leading systems (Perplexity, OpenAI, GPT Researcher, STORM) structure their reports?"
    - "What are the common failure modes?"
    - "How do they handle citation grounding?"
audience: general
findings:
  - id: f1
    claim: "Perplexity Deep Research produces 5–15 page reports with inline hyperlink citations."
    evidence: "Reports span 5–15 pages, formatted for readability with bolded key facts and inline hyperlinks to originals."
    source_id: s1
    confidence: high
  - id: f2
    claim: "OpenAI Deep Research uses an Executive Summary + Detailed Analysis + Key Findings structure."
    evidence: "Reports follow Executive Summary (1–2 paragraphs of core findings) then Detailed Analysis organized by themes then Key Findings as bullets."
    source_id: s2
    confidence: high
  - id: f3
    claim: "Studies find 3–13% of LLM citation URLs are hallucinated."
    evidence: "Citation URLs are hallucinated at 3–13% rates — no Wayback record, likely never existed."
    source_id: s3
    confidence: medium
  - id: f4
    claim: "STORM achieves +25% organization and +10% breadth vs outline-driven baseline."
    evidence: "Compared to outline-driven retrieval-augmented baseline, STORM articles are deemed organized (25% absolute increase) and broad in coverage (10%)."
    source_id: s4
    confidence: high
```

## Expected output (abbreviated `draft_report`)

````markdown
# What is the current state of LLM-based agentic research systems?

## Executive Summary

LLM-powered "deep research" tools now routinely produce multi-page,
cited research reports — but they share a quiet structural problem:
roughly 3–13% of the source URLs they cite never existed [3]. Major
systems (Perplexity, OpenAI, STORM) have converged on similar report
shapes — executive summary, themed body, citations — yet differ in
how rigorously they ground each claim [1][2][4].

## Background

Agentic research systems take a topic, plan a search, gather sources,
and emit a synthesized report. They are increasingly used in place of
manual research, so understanding their conventions and failure modes
matters even to readers who never run one.

## Key Findings

- **Reports cluster around a common shape.** Perplexity produces 5–15
  page reports with bolded facts and inline hyperlink citations [1].
  OpenAI's Deep Research follows Executive Summary → Detailed Analysis
  → Key Findings [2].
- **Citation hallucination is the dominant failure mode.** 3–13% of
  URLs in unaided LLM reports point to pages that never existed [3].
- **Pre-writing structure matters.** STORM's multi-perspective
  outline approach lifts organization by 25 percentage points over an
  outline-driven baseline [4].

## How leading systems structure their reports

Perplexity, OpenAI, and STORM have converged on report shapes that
look broadly similar — an executive summary at the top, themed body
sections, and a sources list — but the details differ. Perplexity's
reports run 5–15 pages with bolded key facts and inline hyperlinks
to originals [1]. OpenAI's Deep Research uses an Executive Summary
of 1–2 paragraphs, then Detailed Analysis organized by theme, then
Key Findings as bullets [2].

```
Common report shape
┌────────────────────────┐
│ Executive Summary      │  1–2 paragraphs, headline claims
├────────────────────────┤
│ Detailed Analysis      │  themed body sections
├────────────────────────┤
│ Key Findings           │  bullets
├────────────────────────┤
│ Sources                │  inline + appended list
└────────────────────────┘
```

## How they handle citation grounding

This is where the systems most visibly differ in quality. Studies
have found that 3–13% of citation URLs in LLM-generated reports
point to pages that never existed [3]. STORM addresses this in
part by simulating multi-perspective Q&A before drafting, which
both broadens coverage and grounds claims in retrieved sources;
the technique improves organization by 25 percentage points and
coverage breadth by 10 points over a simpler baseline [4].

## Limitations & Open Questions

This report draws from four findings about a small number of
systems. Comparative depth on Elicit, Google Deep Research, and
GPT Researcher's detailed report mode would strengthen the
picture but was not in scope here.

## Sources

1. Perplexity AI — "Introducing Perplexity Deep Research" · https://www.perplexity.ai/hub/blog/introducing-perplexity-deep-research
2. OpenAI — "Introducing Deep Research" · https://openai.com/index/introducing-deep-research/
3. arXiv — "Detecting and Correcting Reference Hallucinations in Commercial LLMs and Deep Research Agents" · https://arxiv.org/html/2604.03173
4. Shao et al. — "Assisting in Writing Wikipedia-like Articles From Scratch with Large Language Models" · https://arxiv.org/abs/2402.14207
````

## Expected output (`cited_sources`)

```yaml
cited_sources:
  - id: 1
    title: "Introducing Perplexity Deep Research"
    url: "https://www.perplexity.ai/hub/blog/introducing-perplexity-deep-research"
    findings_supported: [f1]
  - id: 2
    title: "Introducing Deep Research"
    url: "https://openai.com/index/introducing-deep-research/"
    findings_supported: [f2]
  - id: 3
    title: "Detecting and Correcting Reference Hallucinations in Commercial LLMs and Deep Research Agents"
    url: "https://arxiv.org/html/2604.03173"
    findings_supported: [f3]
  - id: 4
    title: "Assisting in Writing Wikipedia-like Articles From Scratch with Large Language Models"
    url: "https://arxiv.org/abs/2402.14207"
    findings_supported: [f4]
```

Notice: every `[n]` in prose matches a `cited_sources` entry, every entry appears at least once in prose, `## Sources` is the rendering of `cited_sources`, prose tone is `general` (no domain jargon undefined, "what this means" framing, plain English), and the "Common report shape" ASCII diagram sits after the first paragraph of a themed body section with no `[n]` inside the fenced block.
