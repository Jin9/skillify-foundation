# Path-Glob Policy Idioms

## Glob syntax

### Prior art at a glance

| Source | Dialect | `**` | `!` negation | Anchoring | Library |
|---|---|---|---|---|---|
| `.gitignore` | gitignore (canonical) | yes (prefix/suffix/infix) | yes (re-include) | leading `/` anchors | n/a (Git built-in) |
| Cargo `package.include` / `package.exclude` | "gitignore-style" (explicit) | yes | yes | leading `/` anchors at package root | gitignore-matcher backed by `globset` family |
| Cargo `workspace.members` / `workspace.exclude` | "typical filename glob" referencing the `glob` crate | `*`, `?` documented; `**` not explicitly listed | no | path-relative | `glob` crate (different from `globset`) |
| `globset` crate | strict superset of gitignore with quirks | yes — but only as prefix `**/`, suffix `/**`, infix `/**/`, or standalone `**` | n/a (matcher only; sets carry their own decision) | controlled by `literal_separator` | `globset` (ripgrep) |
| GitHub CODEOWNERS | "most of the same rules as gitignore" minus `!` and `[…]` | yes | **no** | leading `/` anchors | n/a |
| Semgrep `paths:` | Python `wcmatch` glob | yes | implicit via `exclude:` (no `!`) | matches file and parents | `wcmatch` |

Concrete examples in the wild:

- `cargo-deny` ships `globset = "0.4"` in its `Cargo.toml` and uses it for all path matching in `bans.build` and `sources` checks. The `[[bans.deny]]` table itself has **no** path-scoping field — its "scoped exception" mechanism is `wrappers = ["crate-a", "crate-b"]`, which allow-lists named crates as direct dependents rather than filesystem paths (`{ crate = "ansi_term@0.11.0", wrappers = ["this-crate-directly-depends-on-ansi_term"] }`) ([cargo-deny bans cfg](https://embarkstudios.github.io/cargo-deny/checks/bans/cfg.html), [issue #225](https://github.com/EmbarkStudios/cargo-deny/issues/225)).
- Cargo's manifest reference is explicit: "The patterns should be [gitignore](https://git-scm.com/docs/gitignore)-style patterns" — including `!` negation inside `include` arrays ([Cargo manifest](https://doc.rust-lang.org/cargo/reference/manifest.html#the-exclude-and-include-fields)).
- `globset` is the de-facto Rust glob library (ripgrep, cargo-deny, `ignore`, watchexec). It accepts gitignore patterns but **rejects `**` in non-standard positions** like `a**b` — a real surprise vector ([globset docs](https://docs.rs/globset/latest/globset/)).
- `rustfmt` has `ignore = ["src/generated/**"]` but no per-path override of formatting rules; child `rustfmt.toml` files *replace* rather than merge the `ignore` list ([rustfmt #3881](https://github.com/rust-lang/rustfmt/issues/3881)).

### Decision: what dialect to document for `forbid_under: [string]`

**Recommend: gitignore-style patterns, matched with `globset` semantics, `**` allowed, `!` negation forbidden in the value list.**

Reasoning:

1. Cargo itself promises gitignore-style for `include`/`exclude`, so a Rust audience already has the right reflex.
2. `globset` is what every serious Rust tool ships (cargo-deny, ripgrep, tools using the `ignore` crate). Picking it means zero new dependencies, zero translation layer, and identical error messages to neighbouring tools.
3. Disallow `!` inside `forbid_under` itself — negation belongs to the *schema* (the sibling `allow_under` field), not to the pattern grammar. This avoids the well-known CODEOWNERS footgun where `!` looks supported but isn't ([CODEOWNERS docs](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/customizing-your-repository/about-code-owners) — "`!` is not supported").
4. Document the `globset`-specific `**` placement restriction in a "gotchas" note so users don't write `src**` and get a confusing parse error.

Schema text suggestion: *"Patterns are matched with the `globset` crate using gitignore semantics: `*` matches within a single path segment, `**` matches across segments (only as a full segment: `**/`, `/**`, `/**/`, or `**`), `?` matches one character, `[abc]` and `[a-z]` are character classes, and `{a,b}` alternation is supported. Leading `/` anchors at the repository root. `!` negation is not allowed inside this list — use `allow_under` instead."*

## Precedence model

### Prior art at a glance

| Tool | Shape | Resolution rule |
|---|---|---|
| `.gitignore` | flat list, `!` re-includes | **last-match-wins**, with the carve-out that a re-include cannot resurrect a child of an excluded directory ([gitignore](https://git-scm.com/docs/gitignore)) |
| GitHub CODEOWNERS | flat `path  @owner` list | **last-match-wins**; `!` not supported ([GitHub docs](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/customizing-your-repository/about-code-owners)) |
| GitLab CODEOWNERS | sectioned, `!` supported in 17.10+ | last-match-wins within section and across sections ([GitLab docs](https://docs.gitlab.com/user/project/codeowners/reference/)) |
| ESLint flat config | ordered array of `{files, rules}` objects | **last-match-wins for conflicts; merge for non-conflicting keys**; `ignores` excludes ([ESLint docs](https://eslint.org/docs/latest/use/configure/configuration-files)) |
| Prettier `overrides` | ordered array of `{files, options, excludeFiles}` | array-order precedence (later wins); `excludeFiles` carves out |
| Semgrep `paths:` | `{include: [...], exclude: [...]}` per rule | **exclusion wins over inclusion** when both match ([Semgrep rule syntax](https://semgrep.dev/docs/writing-rules/rule-syntax)) |
| SonarQube file filters | `sonar.inclusions` / `sonar.exclusions` | **exclusions take precedence** when overlapping ([Sonar analysis scope](https://docs.sonarsource.com/sonarqube-server/latest/project-administration/setting-analysis-scope/)) |

### Two competing conventions

There are two coherent traditions, and they disagree:

- **Order-based (gitignore / CODEOWNERS / ESLint / Prettier):** last match wins, regardless of which side it's on. This gives users fine-grained sequencing but requires reading the whole list to predict behaviour.
- **Polarity-based (Semgrep / SonarQube):** the *restrictive* side always wins (`exclude` beats `include`). Order in the file is irrelevant; the safer outcome is structural.

Static-analysis tools — the closest cousins of your code-review skill — uniformly pick polarity-based. This matters because severity-promotion decisions are read by *reviewers under time pressure*, not by config authors carefully arranging order.

### Decision: `forbid_under` wins over `allow_under`

**Keep the current draft.** Document it as: *"When a file path matches both `forbid_under` and `allow_under`, the restrictive policy (`forbid`) wins. This matches Semgrep's `paths.include`/`paths.exclude` and SonarQube's `sonar.inclusions`/`sonar.exclusions` precedence."*

This is preferable to last-match-wins because:

- YAML/JSON arrays carry no meaningful "later in the file" semantics for reviewers — users may sort them alphabetically or with a formatter.
- The failure mode of "we accidentally allowed something we meant to forbid" is strictly worse than the converse.
- Semgrep is the most likely tool a user of a Rust code-review skill has already configured.

Add an escape valve only if needed: a `policy_resolution: "restrictive" | "ordered"` knob, defaulting to `"restrictive"`.

## Severity-zone naming prior art

### How others name the concept

| Tool | Field / concept | Vocabulary |
|---|---|---|
| CodeQL / LGTM | `path_classifiers:` with categories `library`, `test`, `generated`, `documentation`, plus user-defined keys (e.g. `example`) | "classifier", classifies *what kind of code* this is, not its risk level ([lgtm.yml reference](https://help.semmle.com/lgtm-enterprise/archive/1.19/user/help/lgtm.yml-configuration-file.html)) |
| Semgrep | `severity:` on a rule, plus `paths.include` to scope the rule | severity attached to the *rule*, not the path |
| SonarQube | "quality profiles" + "analysis scope" exclusions per path | "scope" / "profile" — risk tier is via the profile assigned to a project |
| GitHub CODEOWNERS | path glob → reviewer set | implicit risk via *who* must approve, not via severity |
| Generic infosec lingo | "blast radius", "trust boundary", "crown jewels" | informal; no tool I found uses these as a config key |

### Decision: rename `money_critical_paths` → `critical_paths` with a `reason` field

**Recommend `critical_paths`** (or `high_risk_paths` as a tied alternative). Rationale:

1. `money_critical_paths` over-fits to one domain. The same field needs to mean "auth boundary" in an SSO repo, "PHI handling" in a healthcare service, and "key material" in a crypto crate. The CodeQL precedent of *generic* category keys (`library`, `test`, `generated`) plus user-defined classifiers is the right template.
2. `critical_paths` reads in plain English and pairs cleanly with the verb you actually want: *escalate findings under critical paths to P1*. Compare: "money-critical paths" forces every reader to translate domain → severity.
3. Avoid `blast_radius` — it's evocative but has no precedent as a config key, and "radius" implies a numeric scope which you're not modelling.
4. Avoid `severity_zones` / `escalation_globs` — they describe the *mechanism* instead of the *meaning*, which makes the YAML harder to skim.

Schema text suggestion:

```yaml
critical_paths:
  - globs: ["src/billing/**", "src/auth/**"]
    reason: "money handling and auth boundary"
    promote_to: P1            # optional; defaults to one level up
```

The `reason` string is borrowed from cargo-deny's `[[bans.deny]] reason = "..."` convention and makes diffs self-documenting. Keeping the field name domain-neutral (`critical_paths`) while letting each repo declare its own reason gives you one field that covers money, auth, credentials, and PII without renaming.
