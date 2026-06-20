"""Tests for the gate state machine: hash chain, schema validation, caps."""

import sys
import unittest
from pathlib import Path

PLUGIN_ROOT = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(PLUGIN_ROOT / "bin"))

import gate_runner  # noqa: E402


class TestHashChain(unittest.TestCase):
    def test_entry_hash_deterministic(self):
        h1 = gate_runner.entry_hash(gate_runner.GENESIS, "start", {"a": 1, "b": 2})
        h2 = gate_runner.entry_hash(gate_runner.GENESIS, "start", {"b": 2, "a": 1})
        self.assertEqual(h1, h2)  # canonical json => key order independent
        self.assertEqual(len(h1), 64)

    def test_chain_links(self):
        h1 = gate_runner.entry_hash(gate_runner.GENESIS, "start", {"k": "DGL-1"})
        h2 = gate_runner.entry_hash(h1, "gate1_submit", {"x": 1})
        self.assertNotEqual(h1, h2)
        # Recomputing the second link from a tampered first link diverges.
        tampered = gate_runner.entry_hash(gate_runner.GENESIS, "start", {"k": "DGL-TAMPER"})
        self.assertNotEqual(gate_runner.entry_hash(tampered, "gate1_submit", {"x": 1}), h2)


class TestSchemaValidation(unittest.TestCase):
    def setUp(self):
        self.schema = gate_runner.load_schema(2)

    def test_valid_payload(self):
        payload = {
            "key": "DGL-1", "branch": "fix/DGL-1-x", "files_changed": ["a.py"],
            "diff_lines": 10, "test_command": "pytest", "test_result": "pass",
            "commit_message": "fix: x (DGL-1)",
        }
        self.assertEqual(gate_runner.validate(payload, self.schema), [])

    def test_missing_required(self):
        errors = gate_runner.validate({"key": "DGL-1"}, self.schema)
        self.assertTrue(any("branch" in e for e in errors))

    def test_wrong_type(self):
        payload = {
            "key": "DGL-1", "branch": "b", "files_changed": "not-a-list",
            "diff_lines": 1, "test_command": "p", "test_result": "r", "commit_message": "m",
        }
        errors = gate_runner.validate(payload, self.schema)
        self.assertTrue(any("files_changed" in e and "array" in e for e in errors))

    def test_bool_is_not_number(self):
        payload = {
            "key": "DGL-1", "branch": "b", "files_changed": [], "diff_lines": True,
            "test_command": "p", "test_result": "r", "commit_message": "m",
        }
        errors = gate_runner.validate(payload, self.schema)
        self.assertTrue(any("diff_lines" in e for e in errors))


class TestCaps(unittest.TestCase):
    def test_diff_lines_breach(self):
        payload = {"files_changed": ["a"], "diff_lines": 99999}
        breaches = gate_runner.check_caps(2, payload)
        self.assertTrue(any("diff_max_lines" in b for b in breaches))

    def test_files_breach(self):
        payload = {"files_changed": [str(i) for i in range(999)], "diff_lines": 1}
        breaches = gate_runner.check_caps(2, payload)
        self.assertTrue(any("diff_max_files" in b for b in breaches))

    def test_within_caps(self):
        self.assertEqual(gate_runner.check_caps(2, {"files_changed": ["a"], "diff_lines": 5}), [])


class TestSelfApproval(unittest.TestCase):
    def test_self_tokens_rejected(self):
        for token in ("agent", "Assistant", "SELF", "ai", "claude"):
            self.assertIn(token.lower(), gate_runner._SELF_APPROVAL)


if __name__ == "__main__":
    unittest.main()
