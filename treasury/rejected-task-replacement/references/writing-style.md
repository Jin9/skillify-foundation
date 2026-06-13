# Writing Style and Task Derivation

## Five style dimensions (scan PASS days only, at least 5 days)

| Dimension          | What to look at                              | Example contrast                                            |
|--------------------|----------------------------------------------|-------------------------------------------------------------|
| 1. Case            | Title Case / lowercase / mixed               | `Executed SIT Test` vs `Execute test`                       |
| 2. Line structure  | single line / multi-line with `[Detail]`     | one line vs summary + `[Detail]`                            |
| 3. Ticket ref      | specific number / generic series / none      | `(ref DGL-12954)` vs `(DGL-13xxx Series)` vs none           |
| 4. Verb pattern    | the verbs used most and their order          | `Executed/Opened/Retested` vs `Fixed/Implemented/Supported` |
| 5. Scope qualifier | how context is stated at the end of the task | `On NEXT And Paotang Channels` vs `Within LOC Scope`        |

Key rule: people in the same role can have very different styles — always mirror **that
person's** style specifically.

## Role-default fallback (when PASS days are fewer than 5)

| Role               | Case       | Ticket ref               | Verb                          | Qualifier                     |
|--------------------|------------|--------------------------|-------------------------------|-------------------------------|
| Dev FE/BE/Mobile   | Title Case | yes (specific)           | Implemented/Fixed/Enhanced    | `On NEXT And Paotang Channels`|
| Tester             | Title Case | none or generic series   | Executed/Opened/Retested      | scope in parentheses          |
| BA                 | Title Case | none                     | Prepared/Conducted/Analyzed   | none                          |

## [Detail] block

Use only when the person actually uses `[Detail]` in their real timesheet for days with similar
activities (do not force it; mirror their real pattern):

```
- <summary: activity1 , activity2>
	[Detail]
	- <activity1> : <detail>
	- <activity2> : <detail>
```

## Task derivation by role

- Dev FE → UI screens, navigation flow, component, GA tracking, integration test.
- Dev BE → API, database field, patch script, service logic, unit test.
- Dev Mobile → Android/iOS native integration, deeplink, permission, webview fix.
- Tester → Execute test [feature/scenario], Open defect [area], Retest defect [ID], Sync with BA.
  - Spread feature areas across the available defects; do not put everything in one group.
  - Pick defects matching the feature area the person *likely* owns (continuity from the previous
    milestone, or an area teammates did not cover).

## Tester condition-4 recheck and continuity (Step 5.5)

Compare each suggested task against other same-role testers on the same or adjacent day:

- If more than 95% similar, change the feature area or pick a different defect ID.
- Watch for days when several testers worked the same feature at once — give the rejected tester a
  different feature that day.
- Plausible continuity: if the person worked feature X in the previous milestone, the first week
  of the current milestone should continue with X or a feature that builds on X (e.g. Accept flow
  following Apply flow).
