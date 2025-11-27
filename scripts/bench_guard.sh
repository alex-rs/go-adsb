#!/usr/bin/env bash
set -euo pipefail

# Run a targeted subset of benchmarks and fail if allocation limits are
# exceeded. Throughput varies by runner hardware, so we gate on
# allocations which are far more stable across environments.

tmp="$(mktemp)"
cleanup() { rm -f "$tmp"; }
trap cleanup EXIT

go test -bench=BenchmarkDecoder -benchmem -run=^$ ./beast >"$tmp"
go test -bench=BenchmarkMessage -benchmem -run=^$ ./adsb >>"$tmp"

python - "$tmp" <<'PY'
import sys

thresholds = {
    "BenchmarkDecoderStream": {"allocs": 8},
    "BenchmarkDecoderNoisyStream": {"allocs": 10},
    "BenchmarkMessageICAO": {"allocs": 0},
    "BenchmarkMessageAltitude": {"allocs": 0},
    "BenchmarkMessageCallsign": {"allocs": 0},
    "BenchmarkMessageGlobalPosition": {"allocs": 0},
}

results = {}

with open(sys.argv[1], "r", encoding="utf-8") as fh:
    for line in fh:
        if not line.startswith("Benchmark"):
            continue

        parts = line.split()
        name = parts[0].split("-")[0]

        try:
            alloc_idx = parts.index("allocs/op")
            allocs = float(parts[alloc_idx - 1])
        except ValueError:
            continue

        results[name] = {"allocs": allocs, "raw": line.strip()}

missing = [name for name in thresholds if name not in results]
if missing:
    raise SystemExit(f"missing benchmark results for: {', '.join(missing)}")

failures = []

for name, limits in thresholds.items():
    allocs = results[name]["allocs"]
    if allocs > limits["allocs"]:
        failures.append(f"{name}: allocs/op {allocs} exceeds limit {limits['allocs']}")

if failures:
    msg = "\n".join(failures + ["\nFull output:", *(r['raw'] for r in results.values())])
    raise SystemExit(msg)
PY
