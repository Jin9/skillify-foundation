# Rust Coverage and Unsafe-Tracking Shapes

## Coverage

Rust has three production-grade coverage stacks: `cargo-llvm-cov` (LLVM source-based instrumentation, the dominant choice), `cargo-tarpaulin` (ptrace-based, the older second), and `grcov` (Mozilla; multi-format converter). All can emit **LCOV `.info`** as a lingua franca; `cargo-llvm-cov` additionally emits the full LLVM JSON export.

### Canonical per-function fields the tools converge on

| Concept | cargo-llvm-cov (LLVM JSON) | LCOV `.info` | cargo-tarpaulin (`--out Json`) | grcov |
|---|---|---|---|---|
| Function identifier | `data[].functions[].name` (mangled) | `FN:<line>,<name>` | `traces[].fn_name` (when present on a `Trace`) | inherits source format (`FN/FNDA` in LCOV out) |
| File path | `data[].files[].filename`; `functions[].filenames[]` | `SF:<path>` (one per section) | `files[].path` (Vec<String> components) | per source format |
| Function entry line | (implicit via `regions`) | `FN:<line number of function start>,<name>` | `traces[].line` for the entry trace | per source format |
| Execution count for function | `functions[].count` (uint64) | `FNDA:<count>,<name>` | summed `traces[].stats.Line(n)` for entry | per source format |
| Functions found / hit (file summary) | `summary.functions.count` / `summary.functions.covered` / `.percent` | `FNF:<n>` / `FNH:<n>` | `coverable` / `covered` (file-level totals, not function-segmented) | per source format |
| Per-line execution count | `segments[]` array on file | `DA:<line>,<count>[,<checksum>]` | `traces[].line` + `traces[].stats` (`Line(n)` variant) | per source format |
| Region (sub-line span) | `functions[].regions[]` = `[line_start, col_start, line_end, col_end, exec_count, file_id, expanded_file_id, region_kind]` | not represented | not represented | not represented |
| Branch counts | `functions[].branches[]` (region tuple + `false_execution_count`); `summary.branches.{count,covered,notcovered,percent}` | `BRDA:<line>,[<exception>]<block>,<branch>,<taken>`; `BRF`/`BRH` totals | partial via traces | LCOV pass-through |
| Macro expansions | `files[].expansions[]` and `functions[]` recursive expansions | not represented | not represented | not represented |
| MC/DC | `summary.mcdc`; `mcdc_records[]` | `MCDC:` / `MRF` / `MRH` | not represented | not represented |
| File total coverage % | `summary.lines.percent`, `summary.regions.percent`, `summary.functions.percent` | computed: `LH/LF`, `FNH/FNF`, `BRH/BRF` | `coverage` (top-level float), `covered`/`coverable` | per source format |

Recommended **`per_function_coverage`** declaration keys (verb-noun, tool-agnostic), drawn from the converged set above: `function_name`, `file_path`, `entry_line`, `execution_count`, `lines_covered`, `lines_total`, `regions_covered`, `regions_total`, `branches_covered`, `branches_total`. Each maps cleanly: LLVM JSON gives all of them; LCOV gives all except `regions_*`; tarpaulin gives lines + function presence; grcov forwards whichever input it ingested.

### cargo-llvm-cov / LLVM JSON ([llvm-project/CoverageExporterJson.cpp](https://github.com/llvm/llvm-project/blob/main/llvm/tools/llvm-cov/CoverageExporterJson.cpp), [taiki-e/cargo-llvm-cov](https://github.com/taiki-e/cargo-llvm-cov), [llvm-cov(1)](https://llvm.org/docs/CommandGuide/llvm-cov.html))

```json
{
  "version": "3.1.0",
  "type": "llvm.coverage.json.export",
  "data": [{
    "files": [{
      "filename": "src/lib.rs",
      "segments": [],
      "expansions": [],
      "branches": [],
      "mcdc_records": [],
      "summary": {
        "lines":         {"count": 42, "covered": 38, "percent": 90.5},
        "functions":     {"count": 5,  "covered": 5,  "percent": 100.0},
        "instantiations":{"count": 5,  "covered": 5,  "percent": 100.0},
        "regions":       {"count": 19, "covered": 17, "notcovered": 2, "percent": 89.5},
        "branches":      {"count": 8,  "covered": 7,  "notcovered": 1, "percent": 87.5},
        "mcdc":          {"count": 0,  "covered": 0,  "notcovered": 0, "percent": 0.0}
      }
    }],
    "functions": [{
      "name": "_ZN8my_crate3foo17h...",
      "count": 12,
      "filenames": ["src/lib.rs"],
      "regions":  [[10, 1, 14, 6, 12, 0, 0, 0]],
      "branches": [],
      "mcdc_records": []
    }],
    "totals": { "...same summary shape..." : null }
  }],
  "cargo_llvm_cov": { "version": "0.6.x", "manifest_path": "/abs/Cargo.toml" }
}
```

`--summary-only` omits `segments`, `expansions`, `branches`, `mcdc_records`, `regions`, and the per-file `functions[]`; only the per-file `summary` and totals are retained.

### cargo-tarpaulin `--out Json` ([xd009642/tarpaulin `src/report/json.rs`](https://github.com/xd009642/tarpaulin/blob/develop/src/report/json.rs))

```json
{
  "files": [{
    "path": ["/abs", "src", "lib.rs"],
    "content": "...source text...",
    "traces": [
      {"line": 10, "address": [4198432], "length": 1, "stats": {"Line": 12}, "fn_name": "foo"}
    ],
    "covered": 38,
    "coverable": 42
  }],
  "coverage": 90.476,
  "covered": 38,
  "coverable": 42
}
```

Schema notes: top-level keys are `files`, `coverage` (percent float), `covered`, `coverable`. Each file holds `path` (Vec<String>), `content`, `traces[]`, plus `covered`/`coverable`. Per-line records are `Trace { line, stats: Line(n) | Branch{..} , fn_name? }`. There is no first-class per-region or per-branch totals object; function-level grouping is reconstructed by walking `traces[].fn_name`.

### LCOV `.info` ([geninfo(1)](https://manpages.debian.org/unstable/lcov/geninfo.1.en.html))

```
TN:my_test_name
SF:/abs/src/lib.rs
FN:10,my_crate::foo
FNDA:12,my_crate::foo
FNF:5
FNH:5
BRDA:11,0,0,12
BRDA:11,0,1,0
BRF:2
BRH:1
DA:10,12
DA:11,12
DA:12,0
LF:42
LH:38
end_of_record
```

Record vocabulary used today: `TN`, `SF`, `VER`, `FN`/`FNL`/`FNA`, `FNDA`, `FNF`, `FNH`, `BRDA`, `BRF`, `BRH`, `DA`, `LF`, `LH`, `MCDC`/`MRF`/`MRH`, `end_of_record`.

### grcov ([mozilla/grcov](https://github.com/mozilla/grcov))

grcov ingests `.profraw`, `.gcda`, LCOV, Go, or JaCoCo and re-emits in formats selected with `-t`: `lcov` (default), `coveralls`, `coveralls+` (adds function-level info), `cobertura[-pretty]`, `html`, `markdown`, `covdir` (recursive JSON tree per directory), `ade`, `files`. It does **not** define a per-function JSON schema of its own; for function-level data prefer LCOV out or `coveralls+`.

## Unsafe tracking

Three independent layers exist: (a) static counters across the dep graph (`cargo-geiger`), (b) human attestations of `unsafe` review (`cargo-vet`, `cargo-supply-chain`), (c) dynamic UB checks (`miri`) and source lints (`clippy::undocumented_unsafe_blocks` and friends).

### Canonical per-`unsafe`-block declaration fields

| Concept | Source of truth | Field / record |
|---|---|---|
| File path | clippy / source | absolute or repo-relative path |
| Span (line, col start/end) | clippy diagnostic JSON | `spans[].file_name`, `line_start`, `line_end`, `column_start`, `column_end` |
| Kind of unsafe site | `cargo-geiger-serde::CounterBlock` | one of `functions`, `exprs`, `item_impls`, `item_traits`, `methods` |
| Safety comment present | `clippy::undocumented_unsafe_blocks` (absence -> warn) | boolean derived from "no `// SAFETY:` adjacent" diagnostic |
| Justification text | `// SAFETY:` comment body | free text captured from source |
| Operation count in block | `clippy::multiple_unsafe_ops_per_block` | integer; >1 triggers lint |
| Crate-wide forbid | `cargo-geiger-serde::UnsafeInfo.forbids_unsafe` | bool (`#![forbid(unsafe_code)]`) |
| Used vs unused tally | `UnsafeInfo.used` / `UnsafeInfo.unused` (`CounterBlock`) | each field is a `Count { safe: u64, unsafe_: u64 }` |
| MIRI status | `cargo +nightly miri test` exit code + stdout summary | `ok` / `FAILED`; commit SHA the run pinned to |
| Audit attestation | `supply-chain/audits.toml` | `[[audits.<crate>]]` entry with `who`, `criteria`, and one of `version` / `delta` / `violation` |
| Allow-list reference | source attribute | `#[allow(clippy::undocumented_unsafe_blocks)]` with optional `reason` |

Recommended **`decision_metadata.unsafe_declarations[]`** keys: `file_path`, `line_start`, `line_end`, `column_start`, `column_end`, `kind` (`function|expr|impl|trait|method`), `safety_comment_present` (bool), `safety_justification` (string), `ops_in_block` (int), `miri_status` (`ok|fail|skipped`), `miri_commit` (SHA), `audit_ref` (path + crate + version into `audits.toml`), `allow_lint_attribute` (string or null).

### cargo-geiger ([geiger-rs/cargo-geiger](https://github.com/geiger-rs/cargo-geiger), [cargo-geiger-serde](https://docs.rs/cargo-geiger-serde/latest/cargo_geiger_serde/))

```json
{
  "packages": {
    "my_crate 0.1.0 (path+file:///...)": {
      "package": {
        "id": {"name": "my_crate", "version": "0.1.0", "source": "..."},
        "dependencies": [], "dev_dependencies": [], "build_dependencies": []
      },
      "unsafety": {
        "used":   {"functions": {"safe": 12, "unsafe_": 1},
                   "exprs":     {"safe": 80, "unsafe_": 4},
                   "item_impls":{"safe": 5,  "unsafe_": 0},
                   "item_traits":{"safe": 1, "unsafe_": 0},
                   "methods":   {"safe": 9,  "unsafe_": 0}},
        "unused": {"functions": {"safe": 0,  "unsafe_": 0}, "exprs": {"safe": 0, "unsafe_": 0},
                   "item_impls":{"safe": 0,  "unsafe_": 0}, "item_traits":{"safe": 0,"unsafe_": 0},
                   "methods":   {"safe": 0,  "unsafe_": 0}},
        "forbids_unsafe": false
      }
    }
  },
  "packages_without_metrics": [],
  "used_but_not_scanned_files": []
}
```

Top-level type is `SafetyReport`. `forbids_unsafe: bool` exists on both `UnsafeInfo` (full report) and `QuickReportEntry` (quick scan); it reflects `#![forbid(unsafe_code)]`.

### cargo-vet / cargo-supply-chain ([Recording Audits](https://mozilla.github.io/cargo-vet/recording-audits.html), [Audit Criteria](https://mozilla.github.io/cargo-vet/audit-criteria.html))

```toml
[[audits.serde_json]]
who      = "Alice Foo <alice@example.com>"
criteria = "safe-to-deploy"
version  = "1.0.115"
notes    = "Reviewed all unsafe blocks; only used for SIMD JSON parsing path."

[[audits.serde_json]]
who      = "Bob Bar <bob@example.com>"
criteria = "safe-to-deploy"
delta    = "1.0.115 -> 1.0.116"

[[audits.bad_crate]]
who       = "Charlie Baz <charlie@example.com>"
criteria  = "safe-to-run"
violation = ">=2.0, <2.3"
```

An entry MUST include exactly one of `version` / `delta` / `violation`, plus `who` and `criteria` (default criteria is `safe-to-deploy`, which explicitly requires "review enough to fully reason about the behavior of all unsafe blocks").

### MIRI ([rust-lang/miri](https://github.com/rust-lang/miri))

Invocation: `MIRIFLAGS="-Zmiri-strict-provenance" cargo +nightly miri test`. Output mirrors libtest: a `running N tests` line, dot stream, then `test result: ok. N passed; 0 failed; 0 ignored; 0 measured; 0 filtered out; finished in <t>`. Exit code is `0` on success, nonzero on UB detection. A "MIRI green at commit X" record is faithfully captured as `{miri_status: "ok", miri_commit: "<sha>", miri_flags: "<MIRIFLAGS>", toolchain: "nightly-YYYY-MM-DD"}`.

### Clippy lints ([rust-clippy lint index](https://rust-lang.github.io/rust-clippy/master/), [undocumented_unsafe_blocks.rs](https://github.com/rust-lang/rust-clippy/blob/master/clippy_lints/src/undocumented_unsafe_blocks.rs))

| Lint name | Category | Default | Detects |
|---|---|---|---|
| `clippy::undocumented_unsafe_blocks` | restriction | allow | `unsafe { ... }` / `unsafe impl` lacking adjacent `// SAFETY:` |
| `clippy::unnecessary_safety_comment` | restriction | allow | `// SAFETY:` placed on safe code |
| `clippy::multiple_unsafe_ops_per_block` | restriction | allow | more than one unsafe op in a single `unsafe { ... }` |
| `clippy::missing_safety_doc` | style | warn | `pub unsafe fn` without `# Safety` rustdoc section |

Diagnostic output (rustc/clippy `--message-format=json` shape):

```json
{
  "reason": "compiler-message",
  "message": {
    "code": {"code": "clippy::undocumented_unsafe_blocks"},
    "level": "warning",
    "message": "unsafe block missing a safety comment",
    "spans": [{
      "file_name": "src/lib.rs",
      "line_start": 42, "line_end": 44,
      "column_start": 5, "column_end": 6,
      "is_primary": true
    }]
  }
}
```

This is the per-block record shape the implement skill should ingest directly when emitting `decision_metadata.unsafe_declarations[]`.
