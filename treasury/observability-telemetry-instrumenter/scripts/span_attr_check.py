#!/usr/bin/env python3
"""observability-telemetry-instrumenter: instrumentation-spec self-check.

Validates the skill's OWN output before it is delivered. PASSES only if all
three hard gates hold:

  1. the token histogram `bucketBoundaries` array EXACTLY equals the 14
     explicit OTel GenAI boundaries
     [1, 4, 16, 64, 256, 1024, 4096, 16384, 65536, 262144, 1048576,
      4194304, 16777216, 67108864];
  2. no metric label set anywhere in the spec contains `user_id` or
     `request_id` (high-cardinality identifiers belong on span attributes and
     histogram exemplars, never on metric labels);
  3. raw-payload capture `optIn` defaults to false (capture is never
     automatic).

No network, no LLM, no code execution — stdlib only. Exit 0 = PASS,
1 = a gate failed, 2 = usage / parse error.

Usage:
  python3 span_attr_check.py instrumentation-spec.json
"""
import argparse
import json
import sys
from pathlib import Path

EXACT_BUCKETS = [
    1, 4, 16, 64, 256, 1024, 4096, 16384, 65536, 262144,
    1048576, 4194304, 16777216, 67108864,
]
FORBIDDEN_LABELS = {"user_id", "request_id"}


def find_bucket_boundaries(obj):
    """Yield every value found under any 'bucketBoundaries' key, recursively."""
    if isinstance(obj, dict):
        for k, v in obj.items():
            if k == "bucketBoundaries":
                yield v
            yield from find_bucket_boundaries(v)
    elif isinstance(obj, list):
        for item in obj:
            yield from find_bucket_boundaries(item)


def find_label_collections(obj):
    """Yield label lists from metricLabels.allowed/forbidden and any 'labels'."""
    if isinstance(obj, dict):
        for k, v in obj.items():
            if k in ("allowed", "forbidden", "labels", "metricLabels"):
                yield k, v
            yield from find_label_collections(v)
    elif isinstance(obj, list):
        for item in obj:
            yield from find_label_collections(item)


def find_opt_in(obj):
    """Yield every value found under any 'optIn' key, recursively."""
    if isinstance(obj, dict):
        for k, v in obj.items():
            if k == "optIn":
                yield v
            yield from find_opt_in(v)
    elif isinstance(obj, list):
        for item in obj:
            yield from find_opt_in(item)


def main():
    ap = argparse.ArgumentParser(
        description="Self-check an OTel GenAI instrumentation spec JSON."
    )
    ap.add_argument("spec", help="path to the instrumentation spec JSON file")
    a = ap.parse_args()

    spec_path = Path(a.spec)
    if not spec_path.is_file():
        print(f"error: spec not found: {a.spec}", file=sys.stderr)
        sys.exit(2)
    try:
        spec = json.loads(spec_path.read_text(encoding="utf-8"))
    except json.JSONDecodeError as e:
        print(f"error: invalid JSON: {e}", file=sys.stderr)
        sys.exit(2)

    failures = []

    # Gate 1: exact bucket boundaries.
    bucket_sets = list(find_bucket_boundaries(spec))
    if not bucket_sets:
        failures.append("no `bucketBoundaries` found — token histogram missing")
    for bb in bucket_sets:
        if bb != EXACT_BUCKETS:
            failures.append(
                f"bucketBoundaries != the exact 14-element array; got {bb!r}"
            )

    # Gate 2: no user_id / request_id on any metric label collection.
    for key, coll in find_label_collections(spec):
        if key == "forbidden":
            # the forbidden list is allowed (and expected) to name them
            continue
        if isinstance(coll, list):
            bad = FORBIDDEN_LABELS.intersection(
                str(x) for x in coll if not isinstance(x, (dict, list))
            )
            if bad:
                failures.append(
                    f"metric label set '{key}' contains forbidden "
                    f"high-cardinality label(s): {sorted(bad)}"
                )

    # Gate 3: raw-payload optIn defaults to false.
    opt_ins = list(find_opt_in(spec))
    if not opt_ins:
        failures.append("no `optIn` found — rawPayloadCapture policy missing")
    for v in opt_ins:
        if v is not False:
            failures.append(
                f"rawPayloadCapture optIn must default to false; got {v!r}"
            )

    if failures:
        print(f"FAIL ({len(failures)}):")
        for f in failures:
            print(f"  - {f}")
        print("FAIL — fix the above before delivering the instrumentation spec.")
        sys.exit(1)

    print("checked: bucketBoundaries, metric labels, rawPayloadCapture.optIn")
    print("PASS — spec uses the exact 14 buckets, no user_id/request_id "
          "metric labels, and opt-in defaults to false.")
    sys.exit(0)


if __name__ == "__main__":
    main()
