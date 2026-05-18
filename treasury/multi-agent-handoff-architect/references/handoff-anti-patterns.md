# Handoff anti-patterns

These are the failure modes the contract design exists to prevent. The
skill flags any design that exhibits them.

## 1. Free-text narrative notes (the dominant context-loss source)

Passing a prose "here's what I did" note instead of a typed payload.
"Free-text handoff notes are the dominant source of context loss in
practice." A handoff "should be structured data, not a long narrative."

Mitigation (enforced): forbid any `notes`/`narrative`/`freetext`/`prose`
field; require the JSON Schema contract with the six required fields and
`schemaVersion`; `additionalProperties: false`. The schema self-check
script fails the design if a free-text property is present.

## 2. Concurrent writers / state corruption

Two agents mutating the same shared state without explicit ownership
transfer. "Allowing two agents to mutate the same context concurrently
produces state corruption that cascades forward through the pipeline."
State-synchronization issues — agents assuming consistent shared state
without synchronization — cause race conditions and stale reads on
blackboards or shared graph state.

Mitigation: binary ownership — exactly one writer per state element at
any moment; explicit ownership-transfer points; writes single-threaded
even when other agents contribute intelligence.

## 3. Context collapse / lossy compression

The receiving agent's context window is saturated by raw prior history,
so it loses early decisions and contradicts them. Lossy compression
"smooths out edge-case details, after which agents proceed confidently
and build the wrong thing." Unbounded context growth: agents accumulate
context across steps without pruning until capacity is exceeded.

Mitigation: structured typed state objects (200–500 tokens) capturing
key decisions, not raw transcripts; explicit `decisions` field so
already-made choices are not silently re-decided; bounded `summary`.

## 4. Role ambiguity

Overlapping or unclear responsibilities produce duplicated work, skipped
verification, and infinite wait loops when no agent owns a terminal
condition. Practitioner restatement: planners suddenly writing code
instead of outlining, peer suggestions vanishing between turns, agents
withholding context while pursuing divergent plans.

Mitigation: schema-validated communication (typed messages) plus
explicit termination criteria and ownership boundaries; one writer and
one terminal owner per task. (Repo-level role *prose* in AGENTS.md is
the job of `agent-context-initializer`; this skill enforces the
*structural* boundary in the contract.)

## 5. Cascading errors

A bad output from one agent propagates to the next and amplifies.
Chained reliability multiplies: two sequential 95%-reliable agents give
90.25%; three give 85.7%; five give ~77%.

Mitigation: intermediate validation checkpoints at every handoff;
agents flag low-confidence outputs (the `confidence` field) before
passing forward; an independent verifier/judge node that can reject
rather than blend is the highest-reliability option (Step 4 of the
decision rule).

## 6. Prompt-injection propagation through the chain

A compromised agent propagates tainted instructions forward, inheriting
user-level privileges and leaking access tokens or executing malicious
tool calls. Prompt-injection defences at individual agents do not
prevent propagation when the inter-agent channel is implicitly trusted.

Mitigation: privilege isolation — agents that read external content must
not have write access to sensitive systems; agents that write to
sensitive systems must receive input only from internal agents, not raw
external data. The typed contract plus payload signing (e.g. A2A card
signing) reduces the implicitly-trusted-channel risk. (The *policy* of
what each agent is permitted to do is owned by
`governance-policy-generator`.)
