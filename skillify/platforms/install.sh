#!/usr/bin/env bash
# install.sh — Install skillify to all supported AI platforms
# Run from: skillify/platforms/
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
SKILL_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
SKILL_NAME="skillify"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

log_ok()   { echo -e "${GREEN}✓${NC} $1"; }
log_warn() { echo -e "${YELLOW}⚠${NC} $1"; }
log_err()  { echo -e "${RED}✗${NC} $1"; }

echo ""
echo "=== Installing ${SKILL_NAME} to all platforms ==="
echo "Source: ${SKILL_DIR}"
echo ""

# --- Claude Code ---
CLAUDE_SKILLS="$HOME/.claude/skills/${SKILL_NAME}"
echo "--- Claude Code ---"
mkdir -p "${CLAUDE_SKILLS}/references"
cp "${SKILL_DIR}/SKILL.md" "${CLAUDE_SKILLS}/SKILL.md"
cp "${SKILL_DIR}/references/"*.md "${CLAUDE_SKILLS}/references/"
log_ok "Installed to ${CLAUDE_SKILLS}"

# --- Gemini CLI ---
GEMINI_SKILLS="$HOME/.gemini/skills/${SKILL_NAME}"
echo ""
echo "--- Gemini CLI ---"
mkdir -p "${GEMINI_SKILLS}/references"
cp "${SKILL_DIR}/SKILL.md" "${GEMINI_SKILLS}/SKILL.md"
cp "${SKILL_DIR}/references/"*.md "${GEMINI_SKILLS}/references/"
log_ok "Installed to ${GEMINI_SKILLS}"

# --- Antigravity (shares Gemini skills path) ---
echo ""
echo "--- Antigravity ---"
log_ok "Shares Gemini skills path (${GEMINI_SKILLS})"
echo "  To complete setup, add this user rule in Antigravity settings:"
echo ""
echo "  When the user asks to create, author, write, refactor, review, or critique"
echo "  a SKILL.md or agent skill, activate the 'skillify' skill"
echo "  and follow its authoring workflow, validation loop, and review checklist."
echo ""

# --- Codex (AGENTS.md) ---
echo "--- Codex ---"
if [ -f "${SCRIPT_DIR}/agents.md" ]; then
    log_ok "Generated AGENTS.md available at: ${SCRIPT_DIR}/agents.md"
    echo "  Copy to your project root: cp ${SCRIPT_DIR}/agents.md /path/to/project/AGENTS.md"
else
    log_warn "agents.md not found in platforms/ — run the setup first"
fi

# --- Copilot ---
echo ""
echo "--- Copilot ---"
if [ -f "${SCRIPT_DIR}/copilot-instructions.md" ]; then
    log_ok "Generated copilot-instructions.md available at: ${SCRIPT_DIR}/copilot-instructions.md"
    echo "  Copy to your project: cp ${SCRIPT_DIR}/copilot-instructions.md /path/to/project/.github/copilot-instructions.md"
else
    log_warn "copilot-instructions.md not found in platforms/ — run the setup first"
fi

echo ""
echo "=== Installation complete ==="
echo ""
echo "Summary:"
echo "  Claude Code  → ${CLAUDE_SKILLS}"
echo "  Gemini CLI   → ${GEMINI_SKILLS}"
echo "  Antigravity  → ${GEMINI_SKILLS} (+ user rule)"
echo "  Codex        → Copy platforms/agents.md to project root as AGENTS.md"
echo "  Copilot      → Copy platforms/copilot-instructions.md to .github/"
echo ""
