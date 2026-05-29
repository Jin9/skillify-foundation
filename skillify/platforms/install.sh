#!/usr/bin/env bash
# install.sh - Install skillify to supported AI agent skill paths.
# Run from: skillify/platforms/
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
SKILL_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
SKILL_NAME="skillify"

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_ok() { echo -e "${GREEN}ok${NC} $1"; }
log_warn() { echo -e "${YELLOW}warn${NC} $1"; }

copy_skill() {
    local dest="$1"
    # Clean refresh of just this skill folder so stale tier-3 files do not
    # linger; only the per-skill dest is removed, never sibling skills.
    rm -rf "$dest"
    mkdir -p "$dest"
    cp "${SKILL_DIR}/SKILL.md" "${dest}/SKILL.md"

    # platforms/ is skillify-internal install scaffolding, not skill content;
    # it is intentionally excluded so installed skills stay free of installer files.
    for dir in references templates scripts examples assets; do
        if [ -d "${SKILL_DIR}/${dir}" ]; then
            mkdir -p "${dest}/${dir}"
            cp -R "${SKILL_DIR}/${dir}/." "${dest}/${dir}/"
        fi
    done
}

install_target() {
    local label="$1"
    local dest="$2"
    echo "--- ${label} ---"
    copy_skill "$dest"
    log_ok "Installed to ${dest}"
    echo ""
}

SHARED_SKILLS="${AGENTS_SKILLS_HOME:-$HOME/.agents/skills}/${SKILL_NAME}"
CODEX_SKILLS="${CODEX_HOME:-$HOME/.codex}/skills/${SKILL_NAME}"
CLAUDE_SKILLS="$HOME/.claude/skills/${SKILL_NAME}"
COPILOT_SKILLS="$HOME/.copilot/skills/${SKILL_NAME}"
INSTALL_CODEX_COMPAT="${INSTALL_CODEX_COMPAT:-0}"
INSTALL_COPILOT="${INSTALL_COPILOT:-0}"

# Antigravity moved its skills dir from ~/.gemini/antigravity to the
# ~/.gemini/antigravity-cli home; prefer the active one. Override with
# ANTIGRAVITY_SKILLS_HOME to point at a specific skills directory.
resolve_antigravity_home() {
    if [ -n "${ANTIGRAVITY_SKILLS_HOME:-}" ]; then
        echo "${ANTIGRAVITY_SKILLS_HOME}"
    elif [ -d "$HOME/.gemini/antigravity-cli" ]; then
        echo "$HOME/.gemini/antigravity-cli/skills"
    else
        echo "$HOME/.gemini/antigravity/skills"
    fi
}
ANTIGRAVITY_SKILLS="$(resolve_antigravity_home)/${SKILL_NAME}"

echo ""
echo "=== Installing ${SKILL_NAME} skill ==="
echo "Source: ${SKILL_DIR}"
echo ""

install_target "Claude Code" "${CLAUDE_SKILLS}"
install_target "Shared Agent Skills (Codex/Gemini/Copilot)" "${SHARED_SKILLS}"
if [ "${INSTALL_CODEX_COMPAT}" = "1" ]; then
    install_target "Codex compatibility path" "${CODEX_SKILLS}"
else
    log_warn "Skipping Codex compatibility path to avoid duplicate skill discovery"
    echo "  Set INSTALL_CODEX_COMPAT=1 to also install ${CODEX_SKILLS}"
    echo ""
fi
if [ "${INSTALL_COPILOT}" = "1" ]; then
    install_target "GitHub Copilot native path" "${COPILOT_SKILLS}"
else
    log_warn "Skipping GitHub Copilot native path (not requested by default)"
    echo "  Set INSTALL_COPILOT=1 to also install ${COPILOT_SKILLS}"
    echo ""
fi
install_target "Antigravity global path" "${ANTIGRAVITY_SKILLS}"

echo "Optional always-on wrappers:"
if [ -f "${SCRIPT_DIR}/agents.md" ]; then
    echo "  Codex AGENTS.md wrapper: ${SCRIPT_DIR}/agents.md"
else
    log_warn "Missing optional Codex wrapper: ${SCRIPT_DIR}/agents.md"
fi

if [ -f "${SCRIPT_DIR}/copilot-instructions.md" ]; then
    echo "  Copilot instructions wrapper: ${SCRIPT_DIR}/copilot-instructions.md"
else
    log_warn "Missing optional Copilot wrapper: ${SCRIPT_DIR}/copilot-instructions.md"
fi

echo ""
echo "=== Installation complete ==="
echo ""
echo "Summary:"
echo "  Claude Code        -> ${CLAUDE_SKILLS}"
echo "  Shared agents      -> ${SHARED_SKILLS}"
if [ "${INSTALL_CODEX_COMPAT}" = "1" ]; then
    echo "  Codex compatibility -> ${CODEX_SKILLS}"
else
    echo "  Codex compatibility -> skipped; use INSTALL_CODEX_COMPAT=1 for legacy clients"
fi
if [ "${INSTALL_COPILOT}" = "1" ]; then
    echo "  GitHub Copilot     -> ${COPILOT_SKILLS}"
else
    echo "  GitHub Copilot     -> skipped; use INSTALL_COPILOT=1 to install"
fi
echo "  Antigravity        -> ${ANTIGRAVITY_SKILLS}"
echo ""
