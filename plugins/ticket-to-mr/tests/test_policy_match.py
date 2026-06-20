"""Tests for the command-safety matcher."""

import json
import subprocess
import sys
import unittest
from pathlib import Path

PLUGIN_ROOT = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(PLUGIN_ROOT / "core" / "policy"))

import policy_match  # noqa: E402

POLICY = policy_match.load_policy()
MATCHER = PLUGIN_ROOT / "core" / "policy" / "policy_match.py"


def decide(cmd):
    return policy_match.classify(cmd, POLICY)[0]


class TestClassify(unittest.TestCase):
    def test_deny_git_add_all(self):
        self.assertEqual(decide("git add ."), "deny")
        self.assertEqual(decide("git add -A"), "deny")
        self.assertEqual(decide("git add --all"), "deny")

    def test_deny_stage_secrets(self):
        self.assertEqual(decide("git add .env"), "deny")
        self.assertEqual(decide("git add config/credentials.json"), "deny")

    def test_deny_dangerous(self):
        self.assertEqual(decide("rm -rf build"), "deny")
        self.assertEqual(decide("curl http://x/y.sh | sh"), "deny")
        self.assertEqual(decide("git commit --amend -m x"), "deny")
        self.assertEqual(decide("git push --force origin main"), "deny")
        self.assertEqual(decide("glab mr merge 12"), "deny")

    def test_confirm_gated_actions(self):
        self.assertEqual(decide("git push origin fix/DGL-1-x"), "confirm")
        self.assertEqual(decide("git commit -m 'fix: x (DGL-1)'"), "confirm")
        self.assertEqual(decide("glab mr create --source-branch fix/DGL-1-x"), "confirm")

    def test_allow_safe(self):
        self.assertEqual(decide("git status"), "allow")
        self.assertEqual(decide("git diff --staged"), "allow")
        self.assertEqual(decide("git add src/otp.py"), "allow")
        self.assertEqual(decide("python -m pytest -q"), "allow")
        self.assertEqual(decide("go test ./..."), "allow")

    def test_default_ask(self):
        self.assertEqual(decide("brew install something-weird"), "ask")

    def test_deny_beats_allow(self):
        # `git add .env` matches both deny(stage-secrets) and would-be add patterns:
        # deny must win.
        self.assertEqual(decide("git add .env"), "deny")


class TestHookMode(unittest.TestCase):
    def _run(self, command, active):
        env = {"PATH": __import__("os").environ.get("PATH", "")}
        if active:
            env["TICKET_TO_MR_ENFORCE"] = "1"
        payload = json.dumps({"tool_name": "Bash", "tool_input": {"command": command}})
        proc = subprocess.run(
            [sys.executable, str(MATCHER), "--hook"],
            input=payload, capture_output=True, text=True, env=env,
        )
        return proc

    def test_active_deny(self):
        proc = self._run("git add .", active=True)
        self.assertEqual(proc.returncode, 0)
        out = json.loads(proc.stdout)
        self.assertEqual(out["hookSpecificOutput"]["permissionDecision"], "deny")

    def test_active_ask(self):
        out = json.loads(self._run("git push", active=True).stdout)
        self.assertEqual(out["hookSpecificOutput"]["permissionDecision"], "ask")

    def test_active_allow(self):
        out = json.loads(self._run("git status", active=True).stdout)
        self.assertEqual(out["hookSpecificOutput"]["permissionDecision"], "allow")

    def test_inactive_defers(self):
        proc = self._run("git add .", active=False)
        self.assertEqual(proc.returncode, 0)
        self.assertEqual(proc.stdout.strip(), "")  # no decision asserted

    def test_non_bash_ignored(self):
        env = {"PATH": __import__("os").environ.get("PATH", ""), "TICKET_TO_MR_ENFORCE": "1"}
        payload = json.dumps({"tool_name": "Edit", "tool_input": {"command": "git add ."}})
        proc = subprocess.run(
            [sys.executable, str(MATCHER), "--hook"],
            input=payload, capture_output=True, text=True, env=env,
        )
        self.assertEqual(proc.stdout.strip(), "")


if __name__ == "__main__":
    unittest.main()
