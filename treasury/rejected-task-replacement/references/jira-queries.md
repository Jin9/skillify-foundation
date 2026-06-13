# Jira Queries (read-only)

Use the cloudId derived from the board URL hostname the user gave at intake. All queries are
read-only. Results often exceed the token budget — save the response to a file and use `jq` to
extract `key`, `summary`, and `status`.

## Dev role — Stories / Tasks by Fix Version

```jql
project = <PROJECT>
AND fixVersion in ("<v1>", "<v2>")
AND issuetype in (Story, Task)
ORDER BY updated DESC
```

- If the Stories returned are not enough, query sub-tasks:
  `parent in (<key1>, <key2>, ...)`.
- If the Stories query returns 0 results, do NOT keep re-querying. Instead derive tasks from the
  **timesheet context of adjacent days** (the day before/after the rejected day).

## Tester role — SIT Defects by date range

`fixVersion` is often incomplete for defects, so query by the milestone's date range:

```jql
project = <PROJECT>
AND issuetype = "SIT Defect"
AND created >= "<YYYY-MM-DD>" AND created <= "<YYYY-MM-DD>"
ORDER BY created DESC
```

Use the milestone's period as the date range.

## Extraction

```bash
jq -r '.issues[] | "\(.key)\t\(.fields.summary)\t\(.fields.status.name)"' jira_response.json
```

Reminder: every ticket reference you place in a derived task must be a real `key` from these
results. If none applies, omit the reference rather than inventing one.
