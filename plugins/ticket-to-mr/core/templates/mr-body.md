<!--
  ticket-to-mr MR body template. Placeholders {LIKE_THIS} are filled by the
  workflow. The body MUST be secret/path/log-scrubbed before submission: no
  tokens, no absolute local paths, no raw test logs.
-->
## [{KEY}] {SUMMARY}

**Jira:** {JIRA_URL}

### What & why
{PROBLEM_STATEMENT}

### Change
{CHANGE_SUMMARY}

### Files touched
{FILES_CHANGED}

### Tests
- Command: `{TEST_COMMAND}`
- Result: {TEST_RESULT}

### Acceptance criteria
{ACCEPTANCE_CRITERIA}

### Known limitations
{KNOWN_LIMITATIONS}

---
- Source branch: `{SOURCE_BRANCH}` → target `{TARGET_BRANCH}`
- No auto-merge · no remove-source-branch
- Gate-3 approver (MR owner of record): {GATE3_APPROVER}

Assisted-by: {AGENT} ({MODEL})
