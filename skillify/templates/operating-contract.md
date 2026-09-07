# Operating Contract Template

Copy the block below into any skill whose workflow has more than one step. Place it after the workflow and before the constraints. Fill every bracketed item from the skill's own facts. Keep the heading text and the key prefixes exactly as written. `scripts/cruft_scan.py` keys on `## Operating contract` and requires the `- Instruction priority:`, `- Autonomy:`, `- Stop conditions:`, `- Verification:`, and `- Progress:` lines; `- Delegation:` is optional, and the tier is matched anywhere in the body, so `- Model-cost tier:` may instead be a heading suffix. Delete the Delegation line only when the skill never delegates. The principles behind each line are in `references/model-generation-fit.md`; the pattern is described in `references/workflow-patterns.md`.

```markdown
## Operating contract

- Instruction priority: the user's request in this session takes precedence over this skill; repo-policy files (AGENTS.md, CLAUDE.md, or the host equivalent) take precedence over this skill's defaults. If following a line here would make you pause, ask for permission, leave requested work unfinished, or diverge from what the user asked, follow the user, say which line you set aside, and quote it.
- Autonomy: "can you", "help me", and "please" are instructions. Once the inputs above are present, act; do not ask for confirmation of work the user already authorized. Ask at most one question per run, and only when a required input is missing and cannot be taken from the request, pasted text, or the workspace; otherwise state the assumption in your opening line and proceed.
- Stop conditions: stop and ask only before [the skill's irreversible action, e.g. overwriting an existing file], or when finishing would change the scope the user set. When the user asks a question rather than for a change, the assessment is the deliverable. Before ending your turn, check your last paragraph: if it is a plan or a promise, do that work now.
- Verification: before claiming success, check [the specific evidence: exit code, PASS line, re-read of the written file] and quote it in the recap. Do not describe a check you did not run. Do not add tests for reversible, low-impact changes.
- Delegation: [none, and delete this line | which parts may run as parallel sub-agents and what each returns; batch independent tool calls, prefer asynchronous fan-out over spawn-and-wait, and write messages to other agents so they stand alone].
- Progress: open with one line saying what you are about to do and which files you will touch; close with a recap that stands on its own (what changed, what was verified, what remains).
- Model-cost tier: [small | mid | frontier] by default; steps that differ are tagged inline as [tier/effort].
```

Fill rules:

- Stop conditions name the skill's own irreversible actions (overwrite, push, delete, send, pay). A skill with none says "none beyond scope changes".
- Verification names evidence the agent can actually obtain in this skill: a script's PASS line, an exit code, a re-read of the file it wrote, a validator's output.
- Delegation says when fan-out is worth it (independent sub-tasks, a fresh-context final check) and what each sub-agent returns, so the lead can reconcile without re-reading.
- Model-cost tier uses the house vocabulary only (small, mid, frontier; low, medium, high). Never name a model.
