# Skillify Before/After Eval Report

Before source: `git show HEAD:skillify/...`
After source: working tree `skillify/...`

## Summary Indicators

| Indicator | Before | After | Delta |
|---|---:|---:|---|
| SKILL.md lines | 178 | 180 | +2 |
| Approx tokens | 2806 | 2770 | -36 |
| Description chars | 746 | 746 | same |
| Quoted trigger count | 12 | 12 | same |
| Mode count | 8 | 8 | same |
| Golden prompt pass rate | 12/12 | 12/12 | same |
| Structural indicator pass rate | 11/18 | 18/18 | +7 |

## Structural Indicators

| Indicator | Before | After | Result |
|---|---|---|---|
| `quick_validate` | PASS | PASS | same |
| `check_links` | PASS | PASS | same |
| `banned_docs_absent` | PASS | PASS | same |
| `description_under_1024` | PASS | PASS | same |
| `trigger_count_at_least_10` | PASS | PASS | same |
| `mode_count_8` | PASS | PASS | same |
| `target_artifact_boundary` | PASS | PASS | same |
| `repo_policy_boundary` | PASS | PASS | same |
| `one_off_prompt_boundary` | PASS | PASS | same |
| `review_audit_no_edit` | PASS | PASS | same |
| `non_create_routes_to_playbooks` | PASS | PASS | same |
| `prompt_collection_in_create` | FAIL | PASS | improved |
| `init_skill_wired` | FAIL | PASS | improved |
| `platforms_not_in_generic_tree` | FAIL | PASS | improved |
| `validation_gate_script_pointer` | FAIL | PASS | improved |
| `iteration_cap_single_sourced` | FAIL | PASS | improved |
| `no_50_line_threshold` | FAIL | PASS | improved |
| `no_deployment_guide_reference` | FAIL | PASS | improved |

## Golden Prompt Results

| Prompt ID | Before | After | Expected mode |
|---|---|---|---|
| `create_skill` | PASS | PASS | Create |
| `write_skill_md` | PASS | PASS | Create |
| `refactor_skill` | PASS | PASS | Refactor |
| `review_no_edit` | PASS | PASS | Review |
| `audit_rubric` | PASS | PASS | Audit |
| `compress_skill` | PASS | PASS | Compress |
| `split_skill` | PASS | PASS | Split |
| `merge_skills` | PASS | PASS | Merge |
| `adapt_codex` | PASS | PASS | Adapt |
| `target_artifact_reject` | PASS | PASS | reject |
| `repo_policy_reject` | PASS | PASS | reject |
| `one_off_prompt_reject` | PASS | PASS | reject |

## Validator Output

- Before quick_validate: PASS `ok: /private/var/folders/fk/c29xbmpd00bb2g95dnh2f2_c0000gp/T/skillify-eval-pmfzypw3/skillify`
- Before check_links: PASS `ok: /private/var/folders/fk/c29xbmpd00bb2g95dnh2f2_c0000gp/T/skillify-eval-pmfzypw3/skillify`
- After quick_validate: PASS `ok: /Users/IF640063/Desktop/project/Jin9/skillify-foundation/skillify`
- After check_links: PASS `ok: /Users/IF640063/Desktop/project/Jin9/skillify-foundation/skillify`
