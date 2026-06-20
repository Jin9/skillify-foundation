#!/usr/bin/env python3
"""
run_tests.py — stdlib test runner for the ticket-to-mr enforcement scaffold.

Discovers and runs every test_*.py in this directory. Exit 0 on success, 1 on
any failure. No third-party deps.

    python3 tests/run_tests.py
"""

import sys
import unittest
from pathlib import Path

HERE = Path(__file__).resolve().parent


def main() -> int:
    suite = unittest.defaultTestLoader.discover(str(HERE), pattern="test_*.py")
    result = unittest.TextTestRunner(verbosity=2).run(suite)
    return 0 if result.wasSuccessful() else 1


if __name__ == "__main__":
    sys.exit(main())
