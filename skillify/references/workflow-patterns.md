# Workflow Patterns for Skills

Five structural patterns for organizing skill instructions. Choose the pattern that matches your use case.

See also: `progressive-disclosure.md` for deciding what stays in `SKILL.md` and `mode-playbooks.md` for Skillify's own mode workflows.

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
- Explicit step ordering with numbered steps
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

---

## Choosing Your Pattern

| Pattern | Best For | Degree of Freedom |
|---------|----------|-------------------|
| Sequential | Onboarding, setup wizards, deployment pipelines | Low |
| Multi-MCP | Cross-service workflows, handoffs | Medium |
| Iterative | Report generation, code review, quality assurance | Medium |
| Context-Aware | File routing, tool selection, adaptive behavior | High |
| Domain-Specific | Compliance, specialized analysis, expert decisions | Low–Medium |

Most skills lean toward one pattern. Complex skills may combine patterns (e.g., sequential + domain-specific for a compliant onboarding flow).
