# LangGraph v1.x Notes

Use this reference after the skill triggers and before changing LangGraph
orchestration code. Facts are current as of 2026-05-27.

## Source-backed facts

- Official LangGraph docs describe LangGraph as a low-level orchestration
  framework and runtime for long-running, stateful agents, focused on durable
  execution, streaming, human-in-the-loop workflows, and persistence.
  Source: https://docs.langchain.com/oss/python/langgraph/overview
- The Graph API models workflows as state, nodes, and edges. `StateGraph` is
  the main graph class for user-defined state, and graphs must be compiled
  before invocation.
  Source: https://docs.langchain.com/oss/python/langgraph/graph-api
- State schemas are typically `TypedDict`; dataclasses are useful when defaults
  are needed. Reducers control how node updates merge into state, while channels
  without reducers overwrite by default.
  Source: https://docs.langchain.com/oss/python/langgraph/graph-api
- LangGraph v1 is largely backwards compatible, but deprecates
  `langgraph.prebuilt.create_react_agent` in favor of `langchain.agents.create_agent`.
  `MessageGraph` is deprecated in favor of `StateGraph` with a `messages` key.
  Source: https://docs.langchain.com/oss/python/migrate/langgraph-v1
- Persistence writes checkpoints per graph step for a thread. State history,
  replay, and `update_state` operate through checkpoints and reducers.
  Development memory stores are not a production persistence substitute.
  Source: https://docs.langchain.com/oss/python/langgraph/persistence
- PyPI listed `langgraph` 1.2.2 as released on 2026-05-26.
  Source: https://pypi.org/project/langgraph/

## Implementation heuristics

1. Prefer explicit `StateGraph` orchestration when the workflow needs branching,
   multiple nodes, durable state, interrupts, replay, or human review.
2. Use `create_agent` for a standard ReAct loop inside a node or as a standalone
   agent when custom graph topology is not needed.
3. Keep graph state serializable and durable. Store facts, messages, decisions,
   identifiers, and small outputs; do not store live clients, closures, env vars,
   filesystem handles, or rendered prompt strings.
4. Make every node a function boundary with a narrow contract: read state,
   perform one unit of work, and return a partial state update.
5. Use reducers intentionally. For chat/message history use `add_messages`; for
   final answer or status fields rely on overwrite semantics unless accumulation
   is explicitly required.
6. Keep node names stable when checkpointers may already contain persisted
   threads. Renames can affect replay, pending tasks, and state history analysis.
7. Put provider-specific calls behind small functions so graph topology and state
   contracts remain provider-neutral.
8. Keep graph import and compile paths credential-free. Missing API keys or local
   CLI dependencies should produce helpful runtime messages, not import failures.

## Review checklist

- Does each node return a dict or `Command` update instead of mutating `state`?
- Are conditional edges pure and based only on state or explicit runtime context?
- Are recursion limits set when loops or agent calls can recurse?
- Are interrupt and human-in-the-loop paths backed by a real checkpointer if they
  must survive process restarts?
- Are persisted state schema changes backward-compatible, or are migration notes
  included?
- Are tests or smoke checks able to run without live model credentials?
