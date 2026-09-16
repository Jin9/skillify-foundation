# Shapes by target kind

Fillers for the skeleton in `SKILL.md`, one row per kind. Pick the row in workflow step 3; the skeleton, the cap, and the tail do not change.

## First pass

| Kind | Bold line | Steps (2 to 5) | For you | Typical Words |
|---|---|---|---|---|
| Task, plan, or ticket, not started | The goal in one sentence: what will be true when it is done. | What happens, in order, each step naming the real file, service, or ticket it touches. | What you do or decide first. | The one term the ticket leans on. |
| Task in progress ("what am I actually doing here?") | The goal, and where you are in it. | Steps labelled `done:` and `left:`, in order, with real names. | What you must decide or do now; what is blocked and on whom. | Usually none. |
| Code or diff | What the code does, or what changed and why, in one sentence. | Execution order with real function and file names; for a diff, before then after. | What changes for you and what to check. | At most one construct name (`closure`, `goroutine`). |
| Error or stack trace | What failed, quoting the error text verbatim in backticks. | What ran, where it stopped (`file:line`), the real cause, or "most likely" when inferred. | What you can do next, as options, not a choice made for you. | The error's key term (`non-fast-forward`, `nil pointer`). |
| Concept, term, or technology | A one-sentence definition. | What it does, when it is used, what it is not (a real neighbor, not a metaphor). | Why it matters in your current file or ticket. | The one or two words the definition needed. |
| Document or spec | What it decides or asks, in one sentence. | Its 3 to 5 main parts, one line each, with section names. | What it requires from you, and by when if stated. | Any acronym it uses. |

Rules that apply to every row:

- "Most likely" marks an inference; a fact you read is stated as a fact.
- Where a row says options, list them as options; choosing is the user's job, or the principal-advisor skill's.
- If the target is one kind but the question is about another (a ticket, but "what does this error mean?"), follow the kind of the question.

## Layers

| Kind | Layer 2: why and mechanism | Layer 3: details, edge cases, numbers |
|---|---|---|
| Task, plan, or ticket | Why the steps are in this order; what each step depends on. | Deadlines, owners, acceptance criteria, what is out of scope. |
| Task in progress | Why the remaining steps are still open; what changed since the start. | Exact commands or files for the next step; known risks. |
| Code or diff | Why it was written this way; what calls it and what it calls. | Edge cases, error paths, performance or concurrency notes, tests that cover it. |
| Error | The mechanism that produced the error, with the real values involved. | Other triggers of the same error; how to confirm the cause; how to prevent it. |
| Concept | How it works underneath, in the user's stack. | Limits, exceptions, versions, numbers, common mistakes. |
| Document or spec | Why it asks what it asks; what it assumes. | Exact clauses, thresholds, dates, exceptions. |

A layer uses the same skeleton: bold line, 2 to 5 steps, For you, optional Words, tail. Only new sentences; never restate the previous pass. After layer 3, a further "more" narrows to one named part and starts a fresh first pass on it.
