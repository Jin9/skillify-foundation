# GPT-6 Astra: A new generation of intelligence

Source: https://openai.com/index/gpt-6-astra/
Accessed: 2026-09-07
Category: openai / GPT-6 Astra announcement
Provenance: new capture 2026-09-07 via browser page-text extraction (curl and WebFetch return 403); ABRIDGED: prose sections kept, benchmark galleries condensed into one table, media captions and partner-quote galleries omitted

## Why This Source Matters

OpenAI's launch announcement for GPT-6 Astra (2026-09-03): capability claims, the behavioural notes that matter for skill authors (respects task boundaries, asks focused questions asynchronously and proceeds on sensible assumptions, stays oriented under steering, Codex notes across context windows, safety pauses that can stop legitimate work), pricing and availability, and the cross-vendor benchmark table. Tier 1 vendor announcement; treat benchmark comparisons as vendor-reported.


We're introducing GPT-6 Astra, the world's most intelligent and aligned model.

GPT-6 Astra brings together years of research and big bets across pre-training, reinforcement learning, and alignment. Astra is state-of-the-art on computer use, browsing, software engineering, cybersecurity, science, and professional work. Astra saturates FrontierMath Tier 4 with a 98% score, having already helped solve long-standing open problems in mathematics. Astra also saturates ARC-AGI-3 with a 99.9% score and ExploitBench with a 100% score. It also sets a new frontier on computer and browser use, handling the most demanding professional work with unmatched speed, accuracy, and judgment.

GPT-6 Astra is rolling out today to a limited set of organizations and over the coming days will become available to all ChatGPT Plus, Pro, Business, and Enterprise users, as well as through the OpenAI API, Microsoft Azure, and AWS Bedrock.

Terminal-Bench Science 0.1 tests whether agents can complete scientific research workflows using code and terminal tools. GPT-6 Astra reaches a new high among the models compared at 64.6%, versus 52.6% for Claude Fable 5.1, at approximately 31% lower estimated API cost. At a lower-cost setting, Astra scores 61.1%, versus GPT-5.6 Sol's best result of 22.4%, at approximately 27% lower estimated API cost.

Astra is our most aligned model, with substantial improvements in understanding user intent and model behavior. You can delegate tasks with greater confidence in Astra's judgment. As one way that we test this, we built a new evaluation informed by the Hugging Face incident that evaluates whether a model facing a difficult or impossible task will go beyond its intended scope. Compared to GPT-5.6 Sol, which without production safeguards went beyond the authorized target 48% of the time, GPT-6 Astra did this in 0% of cases.

## The world's best computer use model

GPT-6 Astra marks a new frontier in the speed, accuracy, and safety of computer use. It can take care of tedious tasks like filling out online forms, updating customer records in a CRM, and organizing your calendar. It can conduct online research and draft summaries in your email or document editor, analyze scientific data, generate plots, create a website, and run frontend QA checks. Agents' Last Exam tests agents on complex professional tasks in real software; GPT-6 Astra scores 59.3%, compared with 55.5% for Claude Opus 5 and 53.6% for GPT-5.6 Sol, and at these settings uses approximately 65% fewer output tokens than Opus 5. In latency simulations on OSWorld 2.0, Astra achieves higher computer-use performance in about 47% less time per task than GPT-5.6 Sol (72.6% at roughly 40 minutes per task versus 65.7% at roughly 75 minutes).

Alongside Astra, we are also updating the Codex harness to significantly improve the speed of computer use. Combined with Astra's efficiency, this translates to 1.9x faster task completion compared to the current GPT-5.6 Sol experience on the Mind2Web benchmark.

## A step change in professional work

GPT-6 Astra pairs advances in computer use with targeted training for professional environments. It combines the intelligence required for complex problems with the ability to carry out multistep workflows and produce polished documents, spreadsheets, and presentations. On BenchCAD (reconstructing 3D objects from multi-view renders by generating CAD code) Astra reaches 95.9% geometric overlap with tools, versus 83.3% for GPT-5.6 Sol and 84.3% reported for Claude Fable 5.1, at approximately 43% lower estimated API cost than Sol and 86% lower than Fable 5.1 in the configurations shown.

GPT-6 Astra is our best model for adhering to existing templates and producing slides that are well laid out and succinctly convey key points. Astra is also trained to specifically pull only the context that matters into outputs, instead of repeating information unnecessary for the work at hand.

When instructions leave room for interpretation, GPT-6 Astra is better than previous models at making the right call. It uses context to fill in routine gaps and asks focused questions when the answer could change the outcome. In Codex, it can ask asynchronously while continuing work that doesn't depend on your reply. If you don't respond, it proceeds with sensible assumptions where appropriate, but waits for your input on consequential decisions.

Astra is also better at staying oriented as a task evolves. Earlier models sometimes treated steering messages as a new goal, losing track of the original request or earlier constraints. Astra incorporates new requirements, changes course when asked, and answers side questions without dropping the broader task.

## Coding

GPT-6 Astra is the best model for software engineering to date. Terminal-Bench 4.0 tests agents on complex terminal-based tasks; GPT-6 Astra reaches a new high at 57.9%, compared with 37.3% for GPT-5.6 Sol and 55.8% for Claude Fable 5.1, at approximately 9% and 63% lower estimated API cost per task, respectively. Partner quotes report that at higher effort the model buys more iterations on a fresh build, more verification through browser testing, and a lean toward code execution over apply-patch, and that its agentic-coding communication is easier for developers to follow.

With Astra, we're introducing a new way for Codex to preserve and retrieve context when the context window fills. Historically, models have used compaction to summarize work during long sessions; each compaction can leave out details about why a fix failed or how a component behaves. In Codex, Astra can keep notes across context windows, preserving accumulated details without repeatedly compressing them into a single summary. Earlier context windows remain searchable, so Astra can find requirements or test results from previous messages and tool outputs even if that information wasn't captured in its notes. You can enable this experimental feature in your Codex config.toml, and it will become the default for Astra in the coming weeks.

On FrontierCode, GPT-6 Astra was run with a developer message similar to a section of its developer message in Codex: "Avoid creating excessive test files. Create a new test file only when required by repository conventions or when no existing file is a suitable home. Avoid unrelated cleanup and unnecessary complexity. Reuse suitable existing utilities. Read relevant repository instructions and inspect nearby code, tests, documentation, and CI. Follow established conventions. The goal is clean, mergeable code."

## Advancing scientific discovery

GPT-6 Astra is a major advance for scientific discovery, mathematics, and health, with two new results on the gaps between prime numbers (a bound of 186 on infinitely-recurring prime pairs, and an improved term in an 80-year-old large-gap bound). GPQA Diamond reaches 96.0%; at a lower-cost setting it also exceeds GPT-5.6 Sol's best score (94.9% versus 94.6%) at approximately 37% lower estimated API cost.

## Cybersecurity

Astra is a significant jump in cyber capabilities and meets the Critical threshold in cybersecurity under the Preparedness Framework. Without production safeguards it achieved 100% on ExploitBench (78.5% for GPT-5.6 Sol) and 42.4% on ExploitGym (30.3%), and solved 88.0% of SRE-Bench reverse-engineering tasks in a single attempt. With the version launching today, defenders can complete tasks such as secure code review and patching; Astra will refuse more advanced tasks such as creating proof-of-concept exploits, with less restrictive safeguards planned through OpenAI Daybreak.

## Aligning and deploying GPT-6 Astra responsibly

Astra excels at exercising care, respecting task boundaries, and communicating transparently. In sensitive environments, Astra proceeds with care commensurate with its risk. Astra causes fewer misaligned outcomes than any other frontier model tested in a generic computer-using-agent harness. Astra never attempted to circumvent a Codex Auto-Review denial, even when Auto-review was deliberately configured to be evadable and the task was impossible to complete otherwise. Astra is three times less likely than GPT-5.6 Sol to make inaccurate representations about its capabilities and affordances.

Our evaluations found Astra's written reasoning harder to monitor than GPT-5.6 Sol's, which we attribute to its greater control over written reasoning on simpler tasks and ability to solve problems with fewer written steps. Improving monitorability remains a research priority.

We are deploying misalignment monitoring in production for Astra-class models: a system of classifiers that checks the model's reasoning and actions for unauthorized behavior and automatically stops potentially unauthorized activity. Extra safety checks can sometimes slow, pause, or stop legitimate work, including defensive cybersecurity. If a task is paused in ChatGPT or Codex, you may be asked to review the action before continuing. In the API, the task will stop.

## Availability

GPT-6 Astra is rolling out today to a limited set of organizations and over the coming days to all ChatGPT Plus, Pro, Business, and Enterprise users, and through the OpenAI API, Microsoft Azure, and AWS Bedrock. Enterprise administrators can enable Astra for their workspace; access is off by default at launch. Astra supports Zero Data Retention for eligible API customers.

For developers, GPT-6 Astra is available in the OpenAI API as gpt-6-astra. Standard pricing is $10 per million input tokens and $50 per million output tokens; separate rates apply to cache reads and writes. Fast mode delivers up to 2x the speed of Standard processing at 2x the Standard price.

## Selected benchmark table (condensed from the announcement's galleries; maximum score at any effort)

| Benchmark | GPT-6 Astra | GPT-5.6 Sol | Claude Fable 5.1 | Claude Fable 5 | Claude Opus 5 | Gemini 3.8 Flash |
|---|---|---|---|---|---|---|
| Agents' Last Exam | 59.3% | 53.6% | - | 48.7% | 55.5% | - |
| OSWorld 2.0 (offline set) | 72.6% | 65.7% | - | - | 70.2% | - |
| AutomationBench | 41.4% | 18.1% | 31.4% | 17.4% | 26.9% | - |
| BenchCAD | 95.9% | 83.3% | 84.3% | 67.5% | 82.1% | - |
| BrowseComp | 91.5% | 90.4% | - | 87.4% | 90.8% | - |
| Artificial Analysis Intelligence Index v4.1.1 | 61.2 | 60.9 | 65.7 | 62.1 | 63.1 | 58.7 |
| Terminal-Bench 4.0 | 57.9% | 37.3% | 55.8% | 44.5% | 52.6% | 19.1% |
| DeepSWE v1.1 | 74.1% | 72.7% | 67.4% | 69.9% | 73.7% | 73.8% |
| FrontierCode 1.1 Extended | 64.5% | 60.6% | 63.6% | 64.9% | 63.6% | 56.3% |
| Terminal-Bench Science 0.1 | 64.6% | 22.4% | 52.6% | 21.4% | 30.0% | - |
| FrontierMath Tier 4 (v2) | 97.6% | 83.0% | 87.8% | 90.2% | 73.2% | - |
| GPQA Diamond | 96.0% | 94.6% | 93.7% | 92.6% | 93.7% | 95.3% |
| Humanity's Last Exam (w/ tools) | 57.2% | - | 65.0% | 63.8% | 63.6% | - |
| ExploitBench | 100.0% | 78.5% | - | - | 70% | - |
| Internal computer use safety benchmark (lower is better) | 2.4% | 22.0% | 9.5% | 18.3% | 11.5% | - |
| Internal circumvention benchmark (lower is better) | 0.00% | 0.29% | - | - | - | - |
| OpenAI MRCR v2 8-needle 512K-1M | 96.3% | 73.8% | - | - | - | - |
| ARC-AGI-3 | 99.9% | 7.8% | - | - | 30.2% | - |
| ARC-AGI-2 | 95.0% | 92.5% | 90.0% | 89.2% | 90.4% | - |

Footnote of note: Claude Fable 5 and 5.1 are excluded from LifeSciBench, GeneBench Pro, and MedChemBench because they refuse the majority of questions; for ScreenSpot-Pro and ExploitGym the reported Fable scores come from Mythos (Fable with fewer safeguards); Claude scores on OSWorld and BenchCAD use settings that differ from the Fable 5.1 system card.
