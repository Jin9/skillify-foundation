# CLI invocation contracts

Exact, deterministic invocation rules for each external model CLI this skill drives. The wrappers
in `scripts/` encode these; this file is the rationale and the manual fallback. Verify model
availability at runtime — treat versioned specifics below as observations, not eternal truth. The
agy and codex sections were re-verified by live dispatch on **2026-07-30** (agy 1.1.8,
codex-cli 0.146.0); each dated claim below is something that was actually run, not inferred.

## Shared rules

- **No `timeout` on macOS.** Every dispatch runs under a background-PID watchdog: launch the CLI in
  the background, start a `( sleep N; kill -TERM $pid )` watcher, `wait` on the CLI, then kill the
  watcher. Detect a watchdog kill with a **flag file written by the watcher**, not by testing for
  exit 143 — a CLI that traps SIGTERM, or that times out internally, exits with its own code.
  - *Hardening:* a killed CLI can leave child processes. Where the platform supports it (job
    control via `set -m`, or `setsid`), terminate the whole process group (`kill -TERM -$pgid`)
    instead of the lone PID. The shipped wrappers kill the direct child for portability; add
    group-kill if you observe orphans.
- **Never put a network call in the dispatch path before the watchdog is armed.** A pre-flight that
  hangs wedges the whole consult with no timeout covering it.
- **Redirect stdin from `/dev/null`** for non-interactive calls so a CLI never blocks on input.
- **Print a `PROVENANCE:` line** after each call (CLI, requested model, verified model, timeout,
  exit). Provenance is not optional — see `adjudication-and-provenance.md`.

## codex (Codex CLI — GPT)

- Non-interactive subcommand: `codex exec [OPTIONS] "<prompt>"`. The prompt may also arrive on stdin.
- **Do not pass `-m` by default.** `codex exec` already resolves the model from `~/.codex/config.toml`
  (honoring `CODEX_HOME`), layered with `-p/--profile` and `-c model=…`. An unconditional `-m` in a
  wrapper silently overrides the user's configuration, defeats profile layering, and is a hardcoded
  value guaranteed to go stale. Pass a model only to pin one deliberately for a single call.
- **Ground truth is free here:** `codex exec` prints its resolved configuration as a header before
  the response — `model:`, `provider:`, `approval:`, `sandbox:`, `reasoning effort:`, `session id:`.
  Parse the **first** `^model: ` line, so no model-authored text can impersonate the header.
  `~/.codex/sessions/**/rollout-*.jsonl` carries the same `"model"` plus `cli_version` as a durable
  secondary record.
- codex does not silently downgrade: an unknown model errors rather than falling back.
- Models on this setup: `gpt-5.6-sol` is the configured default (verified 2026-07-30). `gpt-5.5` is
  legacy (last used 2026-06). `gpt-5` is still rejected by the service on this ChatGPT account —
  re-confirmed 2026-07-30, HTTP 400 "The 'gpt-5' model is not supported when using Codex with a
  ChatGPT account" — as is `gpt-5-codex`. `run_codex.sh` blocks both behind
  `CODEX_ALLOW_UNVERIFIED_MODEL=1` rather than treating a dated observation as permanent.
- codex has **no** `models` subcommand — there is no runtime list to enumerate.
- Wrapper: `scripts/run_codex.sh "<prompt>" [model]` (model optional; omit it to use your config).

## agy (Antigravity CLI — Gemini and others)

- Headless form:
  `agy --log-file "<path>" --model "<MODEL>" --print-timeout "<N>s" --print "<prompt>"`.
- **`--model` accepts either form**: a display label (`Gemini 3.1 Pro (High)`) or a slug as listed by
  `agy models` (`gemini-3.6-flash-low`). `agy models` prints slugs; the error message for an unknown
  model prints display labels; the backend log records display labels.
- **The two forms are not equally reliable.** Verified 2026-07-30: `gemini-3.1-pro-high` is accepted
  with **exit 0 and no warning** but the backend runs the `settings.json` default instead —
  a silent fallback on a slug the CLI itself lists. Its display label `Gemini 3.1 Pro (High)` works.
  `gemini-3.1-pro-low`, `gemini-3.6-flash-low/medium`, and `gemini-3.5-flash-high` all worked as
  slugs. **Prefer the display-label form, and treat post-run verification as mandatory, not
  advisory** — this failure is invisible in the output.
- **An unrecognized model is a hard error** (exit 1) that prints every valid display label. So a
  pre-flight model check is unnecessary — and actively harmful: `agy models` is server-backed and
  has been observed to hang indefinitely (killed at 20s and 25s returning nothing, minutes after an
  identical call succeeded). Never gate a dispatch on it.
- **Flag order is not the hazard it was believed to be.** Placing `--model` after `--print` was
  honored in controlled testing; the drops previously attributed to argument order were caused by
  the model *value*. Keep every flag before `--print` as a convention, but do not rely on ordering
  as an explanation for a wrong model — verify the backend instead.
- **`--print-timeout` (default `5m0s`) is a Go duration and must carry a unit** — a bare integer
  aborts at flag-parse. It bounds the run independently of your watchdog, so a watchdog longer than
  5m never fires: agy self-terminates first with exit 1 and `Error: timeout waiting for response`.
  Set it below the watchdog deliberately so the graceful path wins and the log is flushed.
- **`--log-file <path>` gives each run its own log.** This is what makes fan-out safe: scanning a
  shared log directory for the newest `cli-*.log` races against concurrent dispatches and against
  any interactive session, and will silently attribute another run's model to yours.
- **Self-report is worthless.** agy injects a Gemini identity; even a Claude-backed model will say
  "built by Google DeepMind". Ground-truth from the run log, line `Propagating selected model
  override to backend: label="…"`. Compare it to the request after normalizing both (lowercase,
  strip non-alphanumerics) so slug and label forms are comparable; allow a prefix match, since a
  label may carry a qualifier the slug omits (`claude-sonnet-4-6` → `Claude Sonnet 4.6 (Thinking)`).
- **`--effort low|medium|high` exists but is not the routing knob here.** Gemini model identifiers
  already encode effort, and combining them is a hard error
  (`--model gemini-3.6-flash-high --effort low` → "conflicts with --effort=low"). agy also does not
  validate the value at flag-parse time. Choose effort by choosing the model.
- **`agy --print` is agentic** and will scan the workspace (shell history, dotfiles, the repo) on
  vague prompts. Give explicit absolute paths, demand a specific output shape, and run from a
  neutral working directory (3rd arg to the wrapper). `--add-dir` and `--sandbox` are available as
  additional mechanical scoping if prompt discipline proves insufficient. See
  `vault-and-websearch-grounding.md`.
- The default model lives in `~/.gemini/antigravity-cli/settings.json` (`"model"`, display-label
  form) — this is what a silent fallback lands on. agy also keeps a `trustedWorkspaces` list there,
  but OS-level TCC still blocks iCloud reads regardless of trust.
- Wrapper: `scripts/run_agy.sh "<MODEL>" "<prompt>" [neutral_dir]`.

## Headless host-agent child (claude -p)

- Form: `claude -p "<prompt>" --effort low --model <model>`.
- **`--effort low` is mandatory.** A child inherits the parent session's high effort and hangs
  (9+ min, zero output) on substantial prompts. `CLAUDE_EFFORT` / `MAX_THINKING_TOKENS` env vars do
  NOT override it — only the flag does.
- Prefer `opus` or `haiku` for headless driver stages; `sonnet` is slower and less reliable
  headless (the wrapper warns).
- This wrapper does **not** yet ground-truth its backend, unlike the other two. Until it does, treat
  its provenance line as a record of the *request*, not proof of the model that answered.
- Wrapper: `scripts/run_headless_claude.sh "<prompt>" [model]` (default `opus`).

## Trust gates and permission bypasses

- A workspace-trust-gated CLI exits non-zero (e.g. gemini-cli exits 55) until the workspace is
  trusted.
- **Do not self-authorize a bypass.** Named flags that are off-limits to the agent:
  `--skip-trust`, agy's `--dangerously-skip-permissions`, and codex's
  `--dangerously-bypass-approvals-and-sandbox`. A stalled headless run looks like a hang, and these
  flags are exactly what one is tempted to reach for; reaching for them removes the human from a
  decision they never delegated. Many harnesses flag agent self-authorization as a safety
  violation, and it will not clear from inside the agent anyway.
- Resolution: ask the user to run the command themselves (e.g. via the session `!` prefix), then
  consume the output. This is a hard boundary, not a workaround to engineer around.
