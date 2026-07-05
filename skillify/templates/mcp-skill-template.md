---
name: your-mcp-skill-name
description: >
  Orchestrates [workflow name] across [Service A] and [Service B].
  Use when the user asks to "[trigger phrase 1]", "[trigger phrase 2]",
  "[trigger phrase 3]", or needs to coordinate actions across [services].
# compatibility: claude-code
# metadata:
#   author: your-name
#   version: 1.0.0
#   mcp-server: your-mcp-server
---

# Your MCP Skill Name

## Purpose

[One sentence describing the multi-service workflow this skill orchestrates.]

## Prerequisites

- MCP server `[server-name]` must be connected (Settings > Extensions)
- Required tools: `[tool_1]`, `[tool_2]`, `[tool_3]`

## Core workflow

### Phase 1: [Data Gathering] (via [Service A] MCP)

1. Call MCP tool: `[tool_name]`
   - Parameters: `[param1]`, `[param2]`
   - Expected output: [describe]
2. Validate response:
   - If success: proceed to Phase 2
   - If error: see Troubleshooting

### Phase 2: [Processing] (local)

1. [Process the data from Phase 1]
2. [Transform or validate as needed]
3. Store intermediate result for Phase 3

### Phase 3: [Action] (via [Service B] MCP)

1. Call MCP tool: `[tool_name]`
   - Parameters: `[param1]`, `[result_from_phase_1]`
   - Expected output: [describe]
2. Verify action completed successfully

### Phase 4: [Notification] (optional)

1. [Notify the user or a channel about the result]
2. Include: [summary of actions taken, links, references]

## Output format

[Describe the final deliverable — created resources, summary report, etc.]

## Error handling

### MCP Connection Failed
If you see "Connection refused":
1. Verify MCP server is running: Check Settings > Extensions
2. Confirm API key is valid
3. Try reconnecting: Settings > Extensions > [Service] > Reconnect

### API Rate Limit
If you see "429 Too Many Requests":
1. Wait 60 seconds before retrying
2. Reduce batch size if processing multiple items
3. Log the partial progress for recovery

### Partial Failure (Phase N)
If a phase fails after earlier phases succeeded:
1. Log what was completed successfully
2. Do NOT retry earlier phases
3. Resume from the failed step with the stored intermediate results

## Constraints

- DO NOT call Phase 3 before Phase 1 data is validated
- DO NOT retry failed operations more than 3 times
- MUST log all MCP calls for audit trail
- MUST handle partial failures gracefully (no orphaned resources)

## Examples

### Example: [Happy path scenario]

**User says**: "[Example prompt]"

**Action**:
1. Phase 1: Fetch [data] from [Service A] → receives [result]
2. Phase 2: Process [result] → produces [intermediate]
3. Phase 3: Create [resource] in [Service B] with [intermediate]
4. Phase 4: Notify user with link to [resource]

**Result**: [Final outcome description]
