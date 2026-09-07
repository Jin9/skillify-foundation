# Workflow Patterns for Skills

Six structural patterns for organizing skill instructions. Choose the pattern that matches your use case; Pattern 0 is the default for judgment work, the numbered patterns are for work where order or a tool contract matters.

See also: `progressive-disclosure.md` for deciding what stays in `SKILL.md` and `mode-playbooks.md` for Skillify's own mode workflows.

## Pattern 0: Goal and Constraints

**Use when**: The agent should judge how to do the work (analysis, review, design, writing) and no single sequence is the only safe one. This is the default degree of freedom for judgment work on current models: a scripted step list here lowers output quality.

**Example structure**:

```markdown
## Approach

Goal: a review of `<target>` that a maintainer can act on without re-reading the code.
Constraints: report only findings you can point to a line for, because unverifiable findings cost the reader more than they save. Do not propose refactors outside the changed files; the request sets the scope.
Done when: every finding has a location, a one-line fix, and a severity; the recap lists what was checked and what was not.
```

**Key techniques**:
- State the goal, the one or two real constraints with their reason, and what done looks like
- Name how the agent verifies its result before reporting it
- Number steps only where order matters; describe everything else as outcomes
- Carry the operating contract (Cross-Cutting Techniques below)

---

## Pattern 1: Sequential Workflow Orchestration

**Use when**: Users need multi-step processes in a specific order.

**Example structure**:

```markdown
## Workflow: Onboard New Customer

### Step 1: Create Account
Call MCP tool: `create_customer`
Parameters: name, email, company

### Step 2: Setup Payment
Call MCP tool: `setup_payment_method`
Wait for: payment method verification

### Step 3: Create Subscription
Call MCP tool: `create_subscription`
Parameters: plan_id, customer_id (from Step 1)

### Step 4: Send Welcome Email
Call MCP tool: `send_email`
Template: welcome_email_template
```

**Key techniques**:
- Numbered steps only where a later step consumes an earlier step's output; describe the rest as outcomes
- Dependencies between steps clearly stated
- Validation at each stage before proceeding
- Rollback instructions for failures

---

## Pattern 2: Multi-MCP Coordination

**Use when**: Workflows span multiple services or tools.

**Example: Design-to-development handoff**

```markdown
### Phase 1: Design Export (Figma MCP)
1. Export design assets from Figma
2. Generate design specifications
3. Create asset manifest

### Phase 2: Asset Storage (Drive MCP)
1. Create project folder in Drive
2. Upload all assets
3. Generate shareable links

### Phase 3: Task Creation (Linear MCP)
1. Create development tasks
2. Attach asset links to tasks
3. Assign to engineering team

### Phase 4: Notification (Slack MCP)
1. Post handoff summary to #engineering
2. Include asset links and task references
```

**Key techniques**:
- Clear phase separation with named phases
- Data passing between MCPs (e.g., asset links from Phase 2 → Phase 3)
- Validation before moving to next phase
- Centralized error handling section

---

## Pattern 3: Iterative Refinement

**Use when**: Output quality improves with iteration and self-correction.

**Example: Report generation**

```markdown
## Iterative Report Creation

### Initial Draft
1. Fetch data via MCP
2. Generate first draft report
3. Save to temporary file

### Quality Check
1. Run validation script: `scripts/check_report.py`
2. Identify issues:
   - Missing sections
   - Inconsistent formatting
   - Data validation errors

### Refinement Loop
1. Address each identified issue
2. Regenerate affected sections
3. Re-validate
4. Repeat until quality threshold met (max 3 iterations)

### Finalization
1. Apply final formatting
2. Generate summary
3. Save final version
```

**Key techniques**:
- Explicit quality criteria defined upfront
- Iterative improvement with validation scripts
- Maximum iteration cap to prevent infinite loops
- Clear finalization gate

---

## Pattern 4: Context-Aware Tool Selection

**Use when**: Same outcome requires different tools depending on context.

**Example: Smart file storage**

```markdown
## Smart File Storage

### Decision Tree
1. Check file type and size
2. Determine best storage location:
   - Large files (>10MB): Use cloud storage MCP
   - Collaborative docs: Use Notion/Docs MCP
   - Code files: Use GitHub MCP
   - Temporary files: Use local storage

### Execute Storage
Based on decision:
- Call appropriate MCP tool
- Apply service-specific metadata
- Generate access link

### Provide Context to User
Explain why that storage was chosen
```

**Key techniques**:
- Clear decision criteria with specific thresholds
- Fallback options for each branch
- Transparency about choices made

---

## Pattern 5: Domain-Specific Intelligence

**Use when**: The skill adds specialized knowledge beyond raw tool access.

**Example: Financial compliance**

```markdown
## Payment Processing with Compliance

### Before Processing (Compliance Check)
1. Fetch transaction details via MCP
2. Apply compliance rules:
   - Check sanctions lists
   - Verify jurisdiction allowances
   - Assess risk level
3. Document compliance decision

### Processing
IF compliance passed:
  - Call payment processing MCP tool
  - Apply appropriate fraud checks
  - Process transaction
ELSE:
  - Flag for review
  - Create compliance case

### Audit Trail
- Log all compliance checks
- Record processing decisions
- Generate audit report
```

**Key techniques**:
- Domain expertise embedded in decision logic
- Compliance/governance checkpoint BEFORE action
- Comprehensive audit trail
- Clear governance with IF/ELSE branching

---

## Cross-Cutting Techniques

These apply across patterns; reach for them when a skill calls tools, will be reused by other skills, or retries.

### Observation–action loop (tool-using skills)

When a skill drives external tools, do not fix the whole plan up front. Alternate: take one concrete action, read the observation, then revise the plan before the next action. This keeps the skill responsive to real tool output instead of a stale plan.

### Composition contract (reusable skills)

A skill that other skills or agents will call should state its contract:

- **Preconditions** — what must be true before it runs.
- **Postconditions** — what is guaranteed true after it succeeds.
- **Failure output** — the structured signal it returns when it cannot succeed.

Explicit contracts let skills compose without the caller reading the body.

### Bounded failure capture (retrying skills)

When a skill retries, cap attempts (default 3), record a one-line lesson per failed attempt, and carry only recent lessons forward. State an escalation path for when the cap is hit instead of looping.

### Operating contract (multi-step skills)

When a skill drives the agent through more than one step, state its operating contract in one block instead of scattering rules: instruction priority (the user's request wins over everything; repo-policy files outrank the skill's defaults), autonomy (act once the inputs exist; at most one question per run), stop conditions (destructive or irreversible actions and genuine scope changes, nothing else), verification (which evidence is checked before a claim), delegation (when sub-agents are worth it), progress (an opening line and a standalone recap), and the default model-cost tier. Current models weight skill text heavily and, when a skill line conflicts with the request, may pause or follow the skill instead of the user; the contract gives them the tie-break. The canonical block with fill rules is `templates/operating-contract.md`; do not paraphrase it into a second version.

### Delegation and parallelism

Say when delegation is desirable, not only that it is allowed. Independent sub-tasks with no shared state go to parallel sub-agents; the final check goes to a fresh-context sub-agent, which beats self-critique; sequential or judgment-heavy work stays in the main thread. Batch independent tool calls in one turn. Prefer asynchronous fan-out with a bounded wait over spawn-and-block, so the lead keeps working while sub-agents run. Messages to sub-agents are self-contained: goal, inputs, expected output shape, and stop condition, with no pointer to context the sub-agent cannot see. Treat sub-agent and external-model output as advisory; the lead reconciles it.

### Model-cost tier and effort hints (per node)

Annotate each phase or step with a model-cost tier (`small` for mechanical extraction and formatting, `mid` for structuring and routine passes, `frontier` for ambiguous design and hard judgment) and, where it differs from the host default, an effort hint (`low`, `medium`, `high`). Two house forms are in use and both are accepted: a sentence or heading suffix, `model-cost tier: frontier (edit) / small (test)` or "Every step is model-cost tier: small.", and an inline tag at the start of a step, `[small/low]` or `[frontier/high]`. Name tiers, never models: model names rot, tiers do not. Tier and effort are separate levers; a frontier tier at low effort often beats a smaller tier at high effort on judgment steps, so raise effort only for hard, verifiable steps.

---

## Choosing Your Pattern

| Pattern | Best For | Degree of Freedom |
|---------|----------|-------------------|
| Goal and Constraints | Review, analysis, design, writing, any judgment work | High |
| Sequential | Onboarding, setup wizards, deployment pipelines | Low |
| Multi-MCP | Cross-service workflows, handoffs | Medium |
| Iterative | Report generation, code review, quality assurance | Medium |
| Context-Aware | File routing, tool selection, adaptive behavior | High |
| Domain-Specific | Compliance, specialized analysis, expert decisions | Low–Medium |

Most skills lean toward one pattern. Complex skills may combine patterns (e.g., sequential + domain-specific for a compliant onboarding flow).
