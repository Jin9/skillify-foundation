# CLI invocation contracts

Exact, deterministic invocation rules for each external model CLI this skill drives. The wrappers
in `scripts/` encode these; this file is the rationale and the manual fallback. Verify model
labels and availability at runtime — treat versioned specifics below as defaults, not eternal truth.

## Shared rules

- **No `timeout` on macOS.** Every dispatch runs under a background-PID watchdog: launch the CLI in
  the background, start a `( sleep N; kill -TERM $pid )` watcher, `wait` on the CLI, then kill the
  watcher. A wrapper maps a watchdog kill (exit 143) to exit 124.
  - *Hardening:* a killed CLI can leave child processes. Where the platform supports it (job
    control via `set -m`, or `setsid`), terminate the whole process group (`kill -TERM -$pgid`)
    instead of the lone PID. The shipped wrappers kill the direct child for portability; add
    group-kill if you observe orphans.
- **Redirect stdin from `/dev/null`** for non-interactive calls so a CLI never blocks on input.
- **Print a `PROVENANCE:` line** after each call (CLI, requested model, verified backend, timeout,
  exit). Provenance is not optional — see `adjudication-and-provenance.md`.

## codex (Codex CLI — GPT)

- Non-interactive subcommand: `codex exec [OPTIONS] "<prompt>"`. The prompt may also arrive on stdin.
- Model flag: `-m, --model <MODEL>`. **On a ChatGPT-account login only `gpt-5.5` works**; `gpt-5`
  and `gpt-5-codex` are rejected by the service. `run_codex.sh` refuses them up front.
- codex does not silently downgrade the model — a bad model errors. The requested model in the
  provenance line is therefore the model used.
- Wrapper: `scripts/run_codex.sh "<prompt>" [model]` (default `gpt-5.5`).

## agy (Antigravity CLI — Gemini and others)

- Headless form: `agy --model "<LABEL>" --print "<prompt>"`.
- **`--model` MUST precede `--print`/the prompt.** Go's flag parser stops at the first positional; a
  trailing `--model` is silently dropped and agy falls back to the `settings.json` default.
- **`<LABEL>` must be an exact display label from `agy models`** (e.g. `Gemini 3.1 Pro (High)`,
  `Claude Opus 4.6 (Thinking)`). Internal IDs (`gemini-3.1-pro`) are rejected and silently fall
  back to the default. `run_agy.sh` whole-line-matches the label against `agy models` and aborts on
  a miss.
- **Self-report is worthless.** agy injects a Gemini identity; even a Claude-backed label will say
  "built by Google DeepMind". Ground-truth the backend from the newest
  `~/.gemini/antigravity-cli/log/cli-*.log`, line `Propagating selected model override to backend:
  label="…"`. `run_agy.sh` parses this into the provenance line.
- **`agy --print` is agentic** and will scan the workspace (shell history, dotfiles, the repo) on
  vague prompts. Give explicit absolute paths, demand a specific output shape, and run from a
  neutral working directory (3rd arg to the wrapper). See `vault-and-websearch-grounding.md`.
- The default model lives in `~/.gemini/antigravity-cli/settings.json` (`"model"`); a wrong label
  falls back to it. agy also keeps a `trustedWorkspaces` list there — but OS-level TCC still blocks
  iCloud reads regardless of trust.
- Wrapper: `scripts/run_agy.sh "<LABEL>" "<prompt>" [neutral_dir]`.

## Headless host-agent child (claude -p)

- Form: `claude -p "<prompt>" --effort low --model <model>`.
- **`--effort low` is mandatory.** A child inherits the parent session's high effort and hangs
  (9+ min, zero output) on substantial prompts. `CLAUDE_EFFORT` / `MAX_THINKING_TOKENS` env vars do
  NOT override it — only the flag does.
- Prefer `opus` or `haiku` for headless driver stages; `sonnet` is slower and less reliable
  headless (the wrapper warns).
- Wrapper: `scripts/run_headless_claude.sh "<prompt>" [model]` (default `opus`).

## Trust-gated CLIs (e.g. gemini-cli proper)

- A workspace-trust-gated CLI exits non-zero (e.g. 55) until the workspace is trusted.
- **Do not self-authorize a bypass** (`--skip-trust` or equivalent). Many harnesses flag agent
  self-authorization as a safety violation, and it will not clear from inside the agent anyway.
- Resolution: ask the user to run the command themselves (e.g. via the session `!` prefix), then
  consume the output. This is a hard boundary, not a workaround to engineer around.
