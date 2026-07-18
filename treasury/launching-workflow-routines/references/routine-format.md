# Routine definition format

One Markdown file per routine. `scripts/validate_routine.py` enforces everything on this page; `scripts/routine_parser.py` is the ground truth if the two ever disagree.

## File and registry conventions

- One routine per file, named `<routine-name>.md`, in a registry directory (default: `routines/` at the project root) beside a thin `INDEX.md` listing every routine (see `templates/registry-index-template.md`).
- `INDEX.md` is never a routine; the launcher and validator skip it when globbing the registry.
- Limits: at most 12 nodes per routine. Larger flows must be split into separate routines launched one after the other — never nested (see the no-recursion rule below).

## Frontmatter — flat scalars only

No nested mappings, no block sequences, no multi-line values. Every key is `key: value` on one line.

| Key | Required | Rule |
|---|---|---|
| `routine` | yes | lowercase kebab-case, ≤64 chars, must equal the filename stem |
| `summary` | yes | ≤200 chars, no angle brackets |
| `version` | yes | non-empty string, e.g. `0.1.0` |
| `default_on_fail` | no | `stop` or `continue`; applies to nodes with no `on_fail` of their own (default `stop`) |
| `owner` | no | free text — who maintains this routine |
| `tags` | no | comma-separated labels |
| `requires` | no | external preconditions in prose, e.g. "web search backend", "advisory CLI configured" |

Unknown frontmatter keys are errors.

## `## Inputs` section

Run-level values the user supplies at launch. One bullet per input:

```
## Inputs

- topic: research question for the run (required)
- depth: quick, standard, or deep (optional, default standard)
```

- Keys are `[a-z_]+`. Mark each description `(required)` or `(optional, default X)`.
- Optional inputs MUST state their default; the launcher passes the default when the user gives no value.
- Nodes reference these as `user.<key>` input tokens.

## `## Nodes` section

Exactly one `## Nodes` section per file. Each node is a `### Node: <id>` subsection holding a fixed bullet list of `- key: value` lines. Nodes execute strictly in file order.

| Key | Required | Rule |
|---|---|---|
| `purpose` | yes | one line — what this node produces and why |
| `executor` | yes | the skill name dispatched as a black box, or `agent-inline` when the launching agent does the work itself |
| `executor_mode` | no | mode/submode passed to the executor skill (e.g. `consult`) |
| `tier` | yes | model-cost class the node should run at: `small`, `mid`, or `frontier` |
| `inputs` | yes | comma-separated input tokens (grammar below) |
| `outputs` | yes | ≥1 comma-separated run-dir-relative artifact paths |
| `gate` | yes | human approval points: `none`, `before`, `after`, or `both` |
| `on_fail` | no | `stop` (default via `default_on_fail`) or `continue` |
| `notes` | no | free text: caveats, executor contract details, deliberate deviations |

Unknown node keys and duplicate node ids are errors.

### Input token grammar

| Token | Meaning |
|---|---|
| `user.<key>` | a run-level input; `<key>` must exist in `## Inputs` |
| `<artifact>` | an output path declared by an EARLIER node in the same routine (forward references are errors) |
| `external:<path>` | a pre-existing path outside the run dir, resolved at launch time |
| `none` | the node takes no inputs |

### Output rules

- Run-dir-relative, no leading `/`, no `..`, globally unique across all nodes in the routine.
- Convention: the first output of node N is `NN-<node-id>.md` (two-digit position). The validator warns — not errors — when a node deviates deliberately (e.g. an executor whose own contract fixes the filename).
- `INDEX.md` and `run-report.md` are reserved for the launcher; nodes must not declare them.

## No-recursion rule

A routine node must never launch another routine. The validator rejects any node whose `executor` is `launching-workflow-routines`, ends in `.md`, or matches the filename stem of a sibling routine in the same registry directory. Executors that are themselves multi-stage runner skills are fine — they are skills, not routines.

## Full example

```
---
routine: demo-inventory-digest
summary: Two-node demo - inventory a directory, then digest the inventory.
version: 0.1.0
default_on_fail: stop
---

Prose before the sections is ignored by the parser.

## Inputs

- target_dir: absolute path of the directory to inventory (required)

## Nodes

### Node: inventory

- purpose: List the files in the target directory with sizes
- executor: agent-inline
- tier: small
- inputs: user.target_dir
- outputs: 01-inventory.md
- gate: after

### Node: digest

- purpose: Compress the inventory into a one-paragraph digest
- executor: agent-inline
- tier: small
- inputs: 01-inventory.md
- outputs: 02-digest.md
- gate: after
```

## Validation

```
python3 scripts/validate_routine.py path/to/routine.md [more.md ...]
```

Exit 0 = every file passed. `error:` lines block launch; `warning:` lines do not.
