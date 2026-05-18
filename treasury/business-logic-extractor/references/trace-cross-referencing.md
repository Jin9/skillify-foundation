# Agent / execution-trace cross-referencing

Distilled from research report *Agent trace format*. When execution or agent
traces are available, use them to prove which extracted rules actually *fire*
— distinguishing live business logic from coded-but-dead branches — and to
attach runtime provenance alongside `file:line`.

## Trace model (what to read)
A trace = the full execution of one task; a tree of spans (root = outermost
call, children nested). Each span carries `trace_id`, `span_id`,
`parent_span_id`, start, duration, and typed attributes. Two converging
standards, both on the OpenTelemetry substrate:
- **OTel GenAI semantic conventions**: `gen_ai.operation.name`
  (`invoke_agent` / `execute_tool` / `chat`), `gen_ai.agent.id/name`,
  `gen_ai.provider.name`, token usage, opt-in input/output messages.
- **OpenInference**: `openinference.span.kind` ∈ {CHAIN, LLM, EMBEDDING,
  RETRIEVER, RERANKER, TOOL, AGENT, GUARDRAIL, EVALUATOR, PROMPT}; tool calls
  as indexed attribute paths (function name + arguments).

## Six signal categories to mine
model invocation · tool execution · agent path/routing · retrieval context ·
state/memory reads-writes · governance signals (policy/guardrail/PII). For
rule extraction the load-bearing ones are **tool-execution** and
**agent-path/routing** spans (they show which decision/validation actually
ran) and **state/memory** spans (the same code path behaves differently by
accumulated state — a rule conditioned on state is incomplete without it).

## How to cross-reference a rule
For each extracted rule, search the traces for a span whose tool/decision
corresponds to the rule's code site. Annotate the rule with:
`evidence: exercised | not-observed | dead`, plus the `trace_id`/`span_id`
that exercised it (runtime provenance complementing static `file:line`).
- **exercised**: a span shows the branch taken → high-confidence live rule.
- **not-observed**: code exists, no trace covers it → keep the rule, mark it
  unconfirmed at runtime in the ledger (could be rare path or dead code).
- **dead**: reachable only by branches the traces never take and no caller
  supplies the trigger → flag as candidate dead logic, do not silently drop.

## Handoffs and the "why"
In multi-agent traces, handoff records carry the chain of custody (origin
agent → destination, context payload). When a rule's inputs cross an agent
handoff, cite the handoff span — context loss at handoffs is a leading
multi-agent failure class, so a rule depending on handed-off state is a
ledger-worthy risk. Traces answer *why* a path was taken, not just *what*
ran — capture the deciding condition, not only the outcome.

Traces are evidence, not the spec. Do not generate or modify traces; absence
of a trace is a coverage gap recorded in the ledger, never a reason to invent
runtime behavior.
