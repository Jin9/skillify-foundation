# Span hierarchy

## The tree

The entire agent loop is the root span; each LLM call and tool execution is a
child span, and the parent-child links capture the agent's execution
hierarchy. A trace is the full lifecycle of a request, made of spans that nest
into a parent-child tree.

```
agent run (root span)
├── invoke_agent
│   ├── chat  (LLM call: model, tokens)
│   ├── execute_tool  (web_search)
│   │   └── chat  (LLM call: reads result)
│   └── execute_tool  (calculator)
└── invoke_agent  (delegated sub-agent)
    └── chat
```

## Three nested granularities

Conflating these is the most common dashboarding error. The full hierarchy in
production today is: agent-invocation CLIENT span → per-call CLIENT spans →
per-tool-call INTERNAL spans.

- **API call (per-call).** One HTTP request to a provider, one usage block,
  one `gen_ai.*` CLIENT span. OTel names this the *operation*:
  `gen_ai.operation.name` in {chat, embeddings, generate_content,
  invoke_agent, ...}.
- **Turn (per-turn).** One user→agent exchange. "Calling any of the run
  methods can result in one or more agents running (and hence one or more LLM
  calls), but it represents a single logical turn in a chat conversation."
- **Agent invocation (per-agent).** OTel canonicalizes this as
  `gen_ai.operation.name = invoke_agent`, span name
  `"invoke_agent {gen_ai.agent.name}"`, span kind CLIENT, with tool-execution
  spans nested as INTERNAL.

So: `execute_tool` (per-call) nests under the per-turn exchange, which nests
under the per-agent `invoke_agent` root.

## Span kinds and names

- `invoke_agent {gen_ai.agent.name}` — span kind **CLIENT** for the
  agent-invocation root. (A remote agent call uses span kind CLIENT; a local
  one uses INTERNAL — this skill REQUIRES the CLIENT root for the
  agent-invocation boundary.)
- `execute_tool` — span kind **INTERNAL**. GenAI instrumentations SHOULD
  create tool-execution spans (kind INTERNAL) unless a more specific
  instrumentation already covers them (e.g., MCP tool executions are already
  traced by MCP instrumentation and must not be double-instrumented).

The span name follows a fixed pattern so traces are readable and groupable.
Each agent span is **required** to carry `gen_ai.operation.name` and
`gen_ai.provider.name`, and **conditionally** carries `gen_ai.agent.name`,
`gen_ai.agent.id`, and `gen_ai.conversation.id`.

## Conventions reuse

The agent conventions "extend and override" the generic GenAI client-span
conventions (the ones for a plain `chat` or `generate_content` call), so an
agent's LLM-call children reuse the same well-defined model and token schema
rather than inventing their own. Instrumentation maps a natural agent
structure — session, step, think, tool call — onto this span hierarchy,
tagging every span with a session ID so the whole decision process is visible.

## W3C trace-context propagation

OpenTelemetry rides on vendor-neutral open standards such as the W3C Trace
Context specification — which defines how trace identifiers travel between
services — so an agent's trace can span multiple services and tools and still
hang together. Propagate the W3C trace context across every service hop and
every sub-agent delegation so the tree reconstructs end to end.
