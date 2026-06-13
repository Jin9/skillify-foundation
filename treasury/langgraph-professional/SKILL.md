---
name: langgraph-professional
description: >
  Guide professional LangGraph v1.x implementation, refactoring, and review
  workflows. Use when the user asks "Implement a LangGraph v1.x workflow in
  this repo using professional StateGraph patterns", "Refactor this LangGraph
  repo to follow LangGraph v1.x best practices and clean orchestration design",
  or "Review this LangGraph implementation and suggest fixes for state handling,
  node design, agent usage, and checkpointer safety". Do NOT use for generic
  Python cleanup, unrelated agent frameworks, or repo policy docs.
compatibility: claude-code, codex, copilot, gemini, antigravity
---

# LangGraph Professional

## Purpose

Guide agents to implement, refactor, or review Python LangGraph v1.x workflows
with clear state contracts, maintainable graph orchestration, and safe
persistence assumptions.

## When to use this skill

- Use when: implementing a LangGraph v1.x workflow with `StateGraph`.
- Use when: refactoring LangGraph code toward clean orchestration boundaries.
- Use when: reviewing state handling, node design, agent usage, or checkpointer safety.
- Do NOT use when: editing repository policy files, doing generic Python cleanup, or working with non-LangGraph frameworks.

## Core workflow

1. **Ground the repo**: Inspect dependency versions, `langgraph.json`, graph
   entrypoints, state schemas, nodes, tools, runners, tests, and provider
   configuration before proposing or editing behavior.
2. **Choose the abstraction**: Use `StateGraph` for custom orchestration. Use
   `langchain.agents.create_agent` for standard ReAct/tool-calling loops. Avoid
   deprecated LangGraph v0 prebuilt agent helpers.
3. **Design state first**: Define typed state with `TypedDict` or dataclasses.
   Use explicit reducers such as `add_messages` for accumulating channels.
   Keep prompts and presentation formatting out of durable state.
4. **Build nodes as boundaries**: Make each node own one decision or side effect,
   accept current state, and return a partial state update. Do not mutate state
   in place. Keep provider-specific SDK/API logic behind narrow functions.
5. **Wire orchestration explicitly**: Name nodes clearly, use `START` and `END`,
   keep edge conditions deterministic, and document any intentional loops,
   recursion limits, interrupts, or human review points.
6. **Handle persistence deliberately**: Before adding or changing checkpointers,
   read `references/langgraph-v1-notes.md`. Treat node names, state keys, and
   reducers as compatibility-sensitive when persisted threads may exist.
7. **Verify behavior**: Run graph import checks and focused tests. For this repo,
   preserve the local SDK provider path and headless API provider path unless
   the user explicitly requests a provider change.

## Output format

Produce the artifact the user requested, plus a concise completion note that
includes:

- Graph or state behavior changed.
- Commands run.
- Persistence or compatibility risks.
- Any external model/provider requirements that remain.

For review-only requests, lead with prioritized findings and file references
instead of editing files.

## Constraints and anti-patterns

- DO NOT replace custom graph orchestration with a generic agent loop when the
  workflow requires durable state, branching, interrupts, or multiple nodes.
- DO NOT use deprecated `langgraph.prebuilt.create_react_agent` for new v1.x work.
- DO NOT store rendered prompts, secrets, API keys, or transient client objects
  in graph state.
- DO NOT rename persisted nodes or remove state keys without calling out
  checkpoint migration impact.
- MUST keep long LangGraph research notes in `references/`, not in this file.

## Validation checklist

Before finalizing output, verify:

- [ ] The graph compiles and imports without requiring live credentials.
- [ ] Nodes return state updates and have narrow responsibilities.
- [ ] Reducers match append vs overwrite semantics.
- [ ] Provider-specific behavior is isolated and gracefully degrades.
- [ ] Checkpointer or persistence changes include compatibility notes.

## References

- For current LangGraph v1.x guidance: see `references/langgraph-v1-notes.md`.
- For an implementation/review checklist: see `templates/implementation-checklist.md`.
