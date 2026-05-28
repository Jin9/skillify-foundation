# Failure modes

Detection points and recoveries for the most common failures `review-report`
will encounter mid-run. Each row maps the failure to where in the Procedure
it is detected and the exact recovery action — these are the canonical
recoveries; do not improvise.

| Failure | Detection | Recovery |
|---------|-----------|----------|
| Required input missing or empty | Pre-flight | Emit `final_report = draft_report`, `cited_sources = input cited_sources`, `review_notes` with single entry `[input_incomplete]`. |
| `findings` is empty but `draft_report` has citations | Step 1 | Verify each citation against `cited_sources` alone; mark every claim that has no finding as `citation_unverifiable` flag. Do not delete. |
| Citation marker in draft does not resolve to any `cited_sources.id` | Step 2 | `citation_broken` → swap if a finding supports, else drop. |
| Draft is too long to fully verify in one pass | Step 1 | Verify all citations (this is the load-bearing check) and prioritize step 4 scans by claim density. Note the truncation in review_notes top-risks. |
| Same model "rubber-stamps" with empty notes | Step 6 | Output validator: if `claims_verified == 0` AND citations exist in the draft, the run is incomplete — re-run step 2. |
| Reviewer would rewrite for style | Step 5 | Forbidden — see Step 5 → Hard constraints on revision. Style is the synthesize-report stage's job. |
| Reviewer would add a new claim not in findings | Step 5 | Forbidden by hard constraint (No new claim). |
| Audience-mismatch is pervasive | Step 4 | Flag-only with `audience_mismatch_pervasive`. Do NOT attempt re-targeting. |
