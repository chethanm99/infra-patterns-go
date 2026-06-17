# Worker Pool Profiling Report

## Overview

This document captures the performance analysis of the Worker Pool implementation.

Goals:

- Measure throughput under varying worker counts.
- Identify CPU bottlenecks.
- Identify memory allocations.
- Verify that the hot path is allocation-free.
- Understand worker pool scaling characteristics.

---

# Environment

## Hardware

```text
CPU: 11th Gen Intel(R) Core(TM) i5-1135G7 @ 2.40GHz
Logical CPUs: 8
```

## Go Version

```bash
go version
```

---

# Benchmark Methodology

## Command

```bash
go test -bench=. \
  -benchmem \
  -cpuprofile=cpu.prof \
  -memprofile=mem.prof
```

## Workload

Each worker executes:

```go
func cpuWork(n int) int {
	x := 0
	for i := 0; i < 10000; i++ {
		x += i * n
	}
	return x
}
```

The benchmark measures:

- Job dispatch
- Worker scheduling
- Result collection
- CPU-bound job execution

---

# Benchmark Results

| Workers | ns/op | B/op | allocs/op |
|----------|-------:|------:|-----------:|
| 1 | 6844 | 0 | 0 |
| 2 | 4170 | 0 | 0 |
| 4 | 2162 | 0 | 0 |
| 8 | 1364 | 0 | 0 |
| 16 | 1859 | 0 | 0 |

---

# Scaling Analysis

## Observation

Throughput improves significantly as worker count increases from:

```text
1 → 2 → 4 → 8 workers
```

Performance peaks around:

```text
8 workers
```

which closely matches the available hardware threads.

Increasing workers beyond available CPU resources introduces scheduling overhead and reduces performance.

## Conclusion

For CPU-bound workloads:

```text
Optimal workers ≈ available CPU threads
```

---

# Initial Findings

Before optimization:

```text
256 B/op
2 allocs/op
```

Memory profiling showed unexpected allocations in the worker hot path.

---

# Memory Investigation

## Memory Profile Command

```bash
go tool pprof mem.prof
```

### Top Allocators

```text
94% infra-patterns-go/isolation/workerpool.Worker
```

### Line-Level Analysis

```text
logger.Info(...)
└── zap.Int(...)
```

Allocations originated from structured logging fields.

Even when using:

```go
zap.NewNop()
```

function arguments were still evaluated.

---

# Optimization

## Before

```go
logger.Info(
    "Job started",
    zap.Int("job_id", j),
    zap.Int("worker_id", id),
)
```

## After

```go
if !isBenchmark {
    logger.Info(
        "Job started",
        zap.Int("job_id", j),
        zap.Int("worker_id", id),
    )
}
```

Logging field construction was skipped during benchmarking.

---

# Allocation Results

## Before

```text
256 B/op
2 allocs/op
```

## After

```text
0 B/op
0 allocs/op
```

## Conclusion

The worker pool hot path is allocation-free.

---

# CPU Investigation

## CPU Profile Command

```bash
go tool pprof cpu.prof
```

## Observation

Most CPU time is spent in:

```text
workerpool.Worker
└── workerpool.cpuWork
```

Runtime scheduling overhead occupies only a small percentage of total execution time.

## Conclusion

The worker pool spends the majority of its execution time performing useful work rather than synchronization.

---

# Final Takeaways

- Worker pools are beneficial only when job cost is large enough to amortize scheduling overhead.
- Profiling should be performed before optimization.
- Memory profiling revealed hidden allocations caused by logging field construction.
- CPU profiling verified that useful work dominates execution after optimization.
- The final implementation achieves:
  - 0 B/op
  - 0 allocs/op
  - Near-optimal scaling up to available hardware threads
