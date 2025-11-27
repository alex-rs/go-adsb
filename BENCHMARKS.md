# Benchmarks

This project now ships repeatable microbenchmarks to baseline ADS-B
decode performance and track regressions.

## How to run
- Decode hot paths: `go test -bench=. ./...`
- Bench fixtures live in `internal/benchdata` (Beast frames across
  types 1/2/3 and ADS-B DF5/17/20 samples).

## Current baseline (AMD Ryzen 9 7950X, Go toolchain defaults)
- Beast: `BenchmarkDecoderStream` — ~7.3µs/op, ~121 MB/s, 8 allocs/op.
- Beast: `BenchmarkDecoderNoisyStream` — ~0.87µs/op (short run; stops on first error), ~1104 MB/s, 9 allocs/op.
- ADS-B: `BenchmarkMessageICAO` — ~25 ns/op, ~552 MB/s, 0 allocs/op.
- ADS-B: `BenchmarkMessageAltitude` — ~34 ns/op, ~411 MB/s, 0 allocs/op.
- ADS-B: `BenchmarkMessageCallsign` — ~37 ns/op, ~379 MB/s, 0 allocs/op.
- ADS-B: `BenchmarkMessageGlobalPosition` — ~0.76µs/op, ~37 MB/s, 0 allocs/op.

## Targets to hold the line
- Preserve zero allocations for ICAO and altitude reads; prevent new
  allocations in decode path.
- Keep Beast stream throughput above ~120 MB/s and allocs/op at or
  below current 8.
- Callsign/global position paths should stay allocation-free via the
  new buffer-based helpers.
- Treat ±10% drift from the baseline throughput/latency numbers as a
  signal to investigate before merging.

## CI guardrails
- `scripts/bench_guard.sh` runs in CI and fails if allocation ceilings
  are exceeded for the tracked benchmarks. Update thresholds alongside
  intentional perf changes.
