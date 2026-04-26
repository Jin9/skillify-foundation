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
    mkdir -p "$dest"
    cp "${SKILL_DIR}/SKILL.md" "${dest}/SKILL.md"

    for dir in references templates scripts platforms examples assets; do
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
GEMINI_SKILLS="$HOME/.gemini/skills/${SKILL_NAME}"
COPILOT_SKILLS="$HOME/.copilot/skills/${SKILL_NAME}"
ANTIGRAVITY_SKILLS="$HOME/.gemini/antigravity/skills/${SKILL_NAME}"

echo ""
echo "=== Installing ${SKILL_NAME} skill ==="
echo "Source: ${SKILL_DIR}"
echo ""

install_target "Claude Code" "${CLAUDE_SKILLS}"
install_target "Shared Agent Skills (Codex/Gemini/Copilot)" "${SHARED_SKILLS}"
install_target "Codex compatibility path" "${CODEX_SKILLS}"
install_target "Gemini CLI native path" "${GEMINI_SKILLS}"
install_target "GitHub Copilot native path" "${COPILOT_SKILLS}"
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
echo "  Codex compatibility -> ${CODEX_SKILLS}"
echo "  Gemini CLI         -> ${GEMINI_SKILLS}"
echo "  GitHub Copilot     -> ${COPILOT_SKILLS}"
echo "  Antigravity        -> ${ANTIGRAVITY_SKILLS}"
echo ""
