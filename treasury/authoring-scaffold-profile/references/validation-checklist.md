# Profile validation checklist

Run these before declaring the profile done. Failures block the write.

## 1. Bash syntax

```bash
bash -n profiles/<slug>.sh
```

Must exit 0. A syntax error means the dispatcher crashes when sourcing.

## 2. STAGES → runner mapping

```bash
# read STAGES
source profiles/<slug>.sh
for s in $STAGES; do
  if [[ ! -x "scripts/stage-runners/${s}.sh" ]]; then
    echo "missing runner: $s"
  fi
done
```

Every stage must have a runner. Missing runners block the write.

## 3. GATED_STAGES subset of STAGES

```bash
for g in $GATED_STAGES; do
  echo " $STAGES " | grep -q " $g " || echo "gated stage not in STAGES: $g"
done
```

A gated stage that isn't in STAGES does nothing. Either add it to
STAGES or remove it from GATED_STAGES.

## 4. AGENT names match LiteLLM config

The dispatcher passes the AGENT string to LiteLLM as the model alias.
Common names: `gemini`, `codex`, `claude` (mapped in
`docker/litellm/config.yaml`).

```bash
grep -E '^\s*-\s*model_name:' docker/litellm/config.yaml | awk '{print $3}'
```

Cross-reference. If the AGENT name is not present, surface a warning;
the user may accept (LiteLLM may map differently in their deployment).

## 5. Variable name canonicalization

Every variable must match `^[A-Z_]+$` and be one of:

- `<STAGE>_AGENT`
- `<STAGE>_MODEL`
- `<STAGE>_PROMPT_PREFIX`

Where `<STAGE>` is the upper-case form of a stage in `STAGES`. The
dispatcher reads only these names; typos silently no-op.

## 6. PROMPT_PREFIX hygiene

Each non-empty prefix must:

- Be ≤ 200 characters.
- End with a trailing space.
- Contain no real PII / credentials / customer names.
- Contain no `claude` / `anthropic` brand names.
- Not embed full prompt bodies (those go in `prompts/library/`).

## 7. Idempotent re-source

```bash
( source profiles/<slug>.sh; source profiles/<slug>.sh; declare -p STAGES )
```

Sourcing twice should produce the same STAGES list. Any drift indicates
a bug in the file (e.g., `STAGES+="..."` without a guard).

## 8. Final sanity

```bash
WORKFLOW_PROFILE=<slug> just doctor
```

If the doctor command warns or errors, fix before shipping.

## When a check fails

Surface the failure to the user with the command that detected it, the
expected output, and the actual output. Do not write the file until the
user has resolved the failure.
