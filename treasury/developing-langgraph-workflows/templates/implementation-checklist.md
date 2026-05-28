# LangGraph Implementation Checklist

Use this checklist for implementation, refactor, and review tasks.

## Discovery

- [ ] Identify graph entrypoint and `langgraph.json` mapping.
- [ ] Inspect state schema, reducers, node functions, edge wiring, tools, and providers.
- [ ] Check dependency versions in `pyproject.toml` and lockfile.
- [ ] Determine whether persisted checkpoints exist or are expected in production.

## Design

- [ ] Choose `StateGraph` for custom orchestration or `create_agent` for a standard agent loop.
- [ ] Define typed state and reducer semantics before editing nodes.
- [ ] Keep state serializable and provider-neutral.
- [ ] Plan node names and edge names as stable compatibility surfaces.

## Implementation

- [ ] Make nodes single-purpose and return partial state updates.
- [ ] Keep provider-specific SDK/API code behind narrow helper functions.
- [ ] Add deterministic routing functions for conditional edges.
- [ ] Avoid deprecated LangGraph v0 APIs in new code.
- [ ] Keep graph import and compile credential-free.

## Verification

- [ ] Run a graph import/compile smoke check.
- [ ] Add or update focused tests for state updates, routing, provider fallback, and tool behavior.
- [ ] Verify missing credentials degrade gracefully.
- [ ] Document checkpoint or persistence migration risk in the final response.
