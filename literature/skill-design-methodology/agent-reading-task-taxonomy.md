# Paper / Template Reading Task Taxonomy

## 1. Intake / Capture
Purpose: Understand what the material is and why we read it.

Actions:
- Identify document type
- Identify objective
- Identify target audience
- Identify scope
- Identify owner/source
- Identify date/version
- Detect missing context
- Detect outdated content
- Extract table of contents
- Extract key sections
- Extract glossary / terminology

Example command:
"Read this paper and identify its purpose, scope, assumptions, and key sections."

---

## 2. Extract / Parse
Purpose: Pull useful information out of the material.

Actions:
- Extract key points
- Extract definitions
- Extract claims
- Extract evidence
- Extract formulas
- Extract diagrams
- Extract process flow
- Extract entities
- Extract roles
- Extract requirements
- Extract constraints
- Extract risks
- Extract dependencies
- Extract decision points
- Extract open questions
- Extract action items

Example command:
"Extract all requirements, assumptions, constraints, and risks from this document."

---

## 3. Understand / Explain
Purpose: Convert the material into clear understanding.

Actions:
- Explain in simple terms
- Explain in technical terms
- Explain to manager
- Explain to engineer
- Explain to non-technical stakeholder
- Rewrite for clarity
- Translate
- Define terminology
- Map concept to real use case
- Show example
- Show counter-example
- Create analogy
- Create mental model

Example command:
"Explain this paper like I am a tech lead who needs to apply it to a banking system."

---

## 4. Analyze
Purpose: Break down meaning, structure, and implications.

Actions:
- Analyze structure
- Analyze logic
- Analyze business impact
- Analyze technical impact
- Analyze architecture impact
- Analyze data flow
- Analyze process flow
- Analyze edge cases
- Analyze failure cases
- Analyze trade-offs
- Analyze pros and cons
- Analyze cost
- Analyze complexity
- Analyze security impact
- Analyze scalability impact
- Analyze maintainability impact
- Analyze operational impact

Example command:
"Analyze the architecture impact, trade-offs, risks, and failure cases."

---

## 5. Assumption / Hypothesis
Purpose: Identify what is not explicitly stated but required for decision-making.

Actions:
- List assumptions
- Classify assumptions
- Validate assumptions
- Challenge assumptions
- Convert assumptions to questions
- Convert assumptions to risks
- Define hypothesis
- Define testable hypothesis
- Identify weak assumptions
- Identify dangerous assumptions
- Identify hidden dependencies

Example command:
"List all assumptions and separate them into safe, risky, and need-confirmation."

---

## 6. Compare / Benchmark
Purpose: Compare the material with alternatives, existing systems, or best practices.

Actions:
- Compare with current architecture
- Compare with another paper/template
- Compare with industry practice
- Compare with DDD/CQRS/event-driven approach
- Compare old vs new version
- Compare option A vs option B
- Identify gaps
- Identify overlap
- Identify conflict
- Identify reusable parts
- Identify improvement opportunities

Example command:
"Compare this template with our current DDD/CQRS architecture style and show the gaps."

---

## 7. Validate / Verify
Purpose: Check correctness, feasibility, and trustworthiness.

Actions:
- Validate logic
- Validate requirement completeness
- Validate feasibility
- Validate technical correctness
- Validate with source
- Validate with constraints
- Validate with existing system behavior
- Check contradictions
- Check ambiguity
- Check missing acceptance criteria
- Check compliance concern
- Check security concern
- Check operational readiness

Example command:
"Validate whether this design is feasible and identify contradictions or missing details."

---

## 8. Synthesize / Design
Purpose: Convert reading into new design, framework, or working model.

Actions:
- Synthesize key ideas
- Create framework
- Create design principle
- Create architecture pattern
- Create decision model
- Create checklist
- Create workflow
- Create domain model
- Create bounded context
- Create event map
- Create API model
- Create data model
- Create responsibility matrix
- Create reusable template

Example command:
"Use this paper to synthesize a reusable design framework for our lending origination platform."

---

## 9. Plan
Purpose: Turn insight into execution plan.

Actions:
- Create implementation plan
- Create migration plan
- Create refactor plan
- Create rollout plan
- Create test plan
- Create risk mitigation plan
- Create communication plan
- Create timeline
- Create milestone
- Create backlog
- Create epic/story/task
- Define priority
- Define dependency
- Define owner
- Define acceptance criteria

Example command:
"Turn this into an implementation plan with epics, tasks, risks, and acceptance criteria."

---

## 10. Execute / Apply
Purpose: Use the document to produce real output.

Actions:
- Apply template
- Fill template
- Generate code
- Generate config
- Generate SQL
- Generate API spec
- Generate test cases
- Generate checklist
- Generate Confluence page
- Generate decision record
- Generate architecture document
- Generate diagram description
- Generate migration script
- Generate review comments

Example command:
"Apply this template to our loan application flow and generate the first version."

---

## 11. Review / Critique
Purpose: Improve quality before using or sharing.

Actions:
- Review quality
- Review clarity
- Review completeness
- Review consistency
- Review naming
- Review structure
- Review technical risk
- Review maintainability
- Review security
- Review compliance
- Review assumptions
- Review decision quality
- Suggest improvements
- Rewrite improved version

Example command:
"Review this document like a senior architect and give concrete improvement points."

---

## 12. Summarize / Compress
Purpose: Make the content easier to reuse and share.

Actions:
- Summarize short
- Summarize detailed
- Summarize for executive
- Summarize for engineer
- Summarize as bullet points
- Summarize as table
- Summarize as decision memo
- Summarize as action list
- Summarize as risks
- Summarize as architecture insight
- Summarize as reusable prompt
- Summarize into markdown

Example command:
"Summarize this into a markdown note I can share with another model."

---

## 13. Transform / Reformat
Purpose: Change the shape of the material.

Actions:
- Convert to markdown
- Convert to table
- Convert to checklist
- Convert to JSON
- Convert to YAML
- Convert to API spec
- Convert to backlog
- Convert to diagram syntax
- Convert to meeting note
- Convert to prompt
- Convert to skill.md
- Convert to ADR
- Convert to PRD
- Convert to technical design doc

Example command:
"Transform this paper into a reusable skill.md for another coding agent."

---

## 14. Question / Clarify
Purpose: Prepare questions for humans or other teams.

Actions:
- Generate clarification questions
- Generate business questions
- Generate technical questions
- Generate risk questions
- Generate compliance questions
- Generate data questions
- Generate integration questions
- Generate operational questions
- Generate decision questions
- Prioritize questions
- Separate blocker vs non-blocker questions

Example command:
"Generate clarification questions and separate blocker questions from optional questions."

---

## 15. Decide / Recommend
Purpose: Convert analysis into a decision.

Actions:
- Recommend option
- Rank options
- Decide go/no-go
- Decide architecture direction
- Decide implementation approach
- Decide what to ignore
- Decide what to reuse
- Decide what needs proof of concept
- Define decision criteria
- Explain rationale
- Explain trade-offs
- Explain risk acceptance

Example command:
"Recommend the best approach and explain trade-offs, risks, and decision criteria."

---

# Minimum Must-Do Actions When Reading Any Paper / Template

Use this as the default checklist.

1. Identify purpose
2. Identify scope
3. Extract key ideas
4. Extract requirements / rules / constraints
5. Extract assumptions
6. Extract risks
7. Analyze implications
8. Compare with current context
9. Validate feasibility
10. Identify gaps and open questions
11. Convert insight into action plan
12. Summarize into reusable output

---

# Strong Agent Command Pattern

Use this structure when asking another model/agent to read something:

"Read the attached document/template/paper. Then perform:
1. Intake: identify purpose, scope, audience, and version.
2. Extract: list key concepts, requirements, constraints, assumptions, risks, and dependencies.
3. Analyze: explain technical/business impact, trade-offs, edge cases, and failure cases.
4. Validate: check ambiguity, contradiction, feasibility, and missing information.
5. Synthesize: convert useful ideas into reusable framework/checklist/design pattern.
6. Plan: create actionable epics/tasks with acceptance criteria.
7. Summarize: produce a concise markdown output for review by another model or human."

---

# Compact Action Verb List

- Read
- Identify
- Extract
- Explain
- Analyze
- Classify
- Compare
- Validate
- Challenge
- Infer
- Synthesize
- Design
- Plan
- Execute
- Apply
- Review
- Critique
- Improve
- Summarize
- Transform
- Recommend
- Decide
- Document