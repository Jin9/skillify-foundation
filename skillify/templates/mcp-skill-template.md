---
name: your-mcp-skill-name
description: >
  Orchestrates [workflow name] across [Service A] and [Service B].
  Use when the user asks to "[trigger phrase 1]", "[trigger phrase 2]",
  "[trigger phrase 3]", or needs to coordinate actions across [services].
# compatibility: claude-code, codex, copilot, gemini, antigravity
# metadata:
#   author: your-name
#   version: 1.0.0
#   mcp-server: your-mcp-server
---

# Your MCP Skill Name

## Purpose

[One sentence describing the multi-service workflow this skill orchestrates.]

## Prerequisites

- MCP server `[server-name]` connected in the host's MCP configuration
- Required tools: `[tool_1]`, `[tool_2]`, `[tool_3]`

## Core workflow

Phases with no data dependency on each other may run as parallel sub-agents; say which. Phase order below is a data dependency, not a habit.

### Phase 1: [Data Gathering] (via [Service A] MCP) · model-cost tier: [small | mid | frontier]

1. Call MCP tool: `[tool_name]`
   - Parameters: `[param1]`, `[param2]`
   - Expected output: [describe]
2. Validate response:
   - If success: proceed to Phase 2
   - If error: see Troubleshooting

### Phase 2: [Processing] (local) · model-cost tier: [small | mid | frontier]

1. [Process the data from Phase 1]
2. [Transform or validate as needed]
3. Store intermediate result for Phase 3

### Phase 3: [Action] (via [Service B] MCP) · model-cost tier: [small | mid | frontier]

1. Call MCP tool: `[tool_name]`
   - Parameters: `[param1]`, `[result_from_phase_1]`
   - Expected output: [describe]
2. Verify action completed successfully

### Phase 4: [Notification] (optional) · model-cost tier: [small | mid | frontier]

1. [Notify the user or a channel about the result]
2. Include: [summary of actions taken, links, references]

## Output format

[Describe the final deliverable — created resources, summary report, etc.]

## Operating contract

[Filled from the skillify operating-contract template; keep the heading and the key prefixes exact, and delete this note.]

- Instruction priority: the user's request in this session takes precedence over this skill; repo-policy files (AGENTS.md, CLAUDE.md, or the host equivalent) take precedence over this skill's defaults. If following a line here would make you pause, ask for permission, leave requested work unfinished, or diverge from what the user asked, follow the user, say which line you set aside, and quote it.
- Autonomy: "can you", "help me", and "please" are instructions. Once the inputs above are present, act; do not ask for confirmation of work the user already authorized. Ask at most one question per run, and only when a required input is missing and cannot be taken from the request, pasted text, or the workspace; otherwise state the assumption in your opening line and proceed.
- Stop conditions: stop and ask only before [the skill's irreversible action, e.g. overwriting an existing file], or when finishing would change the scope the user set. When the user asks a question rather than for a change, the assessment is the deliverable. Before ending your turn, check your last paragraph: if it is a plan or a promise, do that work now.
- Verification: before claiming success, check [the specific evidence: exit code, PASS line, re-read of the written file] and quote it in the recap. Do not describe a check you did not run. Do not add tests for reversible, low-impact changes.
- Delegation: [none | which parts may run as parallel sub-agents and what each returns]. Batch independent tool calls; prefer asynchronous fan-out over spawn-and-wait; write messages to other agents so they stand alone.
- Progress: open with one line saying what you are about to do and which files you will touch; close with a recap that stands on its own (what changed, what was verified, what remains).
- Model-cost tier: [small | mid | frontier] by default; steps that differ are tagged inline as [tier/effort].

## Error handling

### MCP Connection Failed
If the MCP call cannot connect:
1. Confirm the server is registered and running in the host's MCP configuration
2. Confirm the credential is valid
3. Reconnect through the host and retry once

### API Rate Limit
If you see "429 Too Many Requests":
1. Wait 60 seconds before retrying
2. Reduce batch size if processing multiple items
3. Log the partial progress for recovery

### Partial Failure (Phase N)
If a phase fails after earlier phases succeeded:
1. Log what was completed successfully
2. Resume from the failed step with the stored intermediate results; earlier phases already succeeded, and re-running them creates duplicate resources

## Constraints

- Validate Phase 1 data before Phase 3 runs; Phase 3 writes are expensive to undo.
- Cap retries at 3, record a one-line lesson per failure, then stop and report; a flaky service must not loop.
- Log every MCP call and result; the log is both the audit trail and the recovery point for partial failures.
- Leave no orphaned resources: on a partial failure, report what exists and what was rolled back.

## Examples

### Example: [Happy path scenario]

**User says**: "[Example prompt]"

**Action**:
1. Phase 1: Fetch [data] from [Service A] → receives [result]
2. Phase 2: Process [result] → produces [intermediate]
3. Phase 3: Create [resource] in [Service B] with [intermediate]
4. Phase 4: Notify user with link to [resource]

**Result**: [Final outcome description]
