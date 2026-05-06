#!/usr/bin/env bash
# profiles/<slug>.sh
# Domain profile: <one-line description>.
#
# Use: WORKFLOW_PROFILE=<slug> just workflow <id> "<goal>" <cap>

# Stages to run, in order. Every name must match scripts/stage-runners/<name>.sh.
STAGES="<space-separated stage list>"

# Stages requiring human approval before running. Must be a subset of STAGES.
GATED_STAGES="<space-separated subset, or empty>"

# Per-stage agent (LiteLLM model alias) and model identifier.
# Defaults below match the scaffold's `code` profile; adjust as needed.

export RESEARCH_AGENT="${RESEARCH_AGENT:-gemini}"
export RESEARCH_MODEL="${RESEARCH_MODEL:-gemini-research}"
# Optional steering prefix (≤ 200 chars, trailing space). Long prompts go in
# prompts/library/research/<topic>.md instead.
# export RESEARCH_PROMPT_PREFIX='<short prefix> '

export PLAN_AGENT="${PLAN_AGENT:-codex}"
export PLAN_MODEL="${PLAN_MODEL:-codex-plan}"
# export PLAN_PROMPT_PREFIX='<short prefix> '

export CRITIQUE_AGENT="${CRITIQUE_AGENT:-claude}"
export CRITIQUE_MODEL="${CRITIQUE_MODEL:-claude-critique}"
# export CRITIQUE_PROMPT_PREFIX='<short prefix> '

# Drop the implement block when the profile skips implement.
export IMPLEMENT_AGENT="${IMPLEMENT_AGENT:-codex}"
export IMPLEMENT_MODEL="${IMPLEMENT_MODEL:-codex-implement}"
# export IMPLEMENT_PROMPT_PREFIX='<short prefix> '

export REVIEW_AGENT="${REVIEW_AGENT:-claude}"
export REVIEW_MODEL="${REVIEW_MODEL:-claude-review}"
# export REVIEW_PROMPT_PREFIX='<short prefix> '

# `test` stage is local (no model). No AGENT/MODEL/PREFIX here.
