# Benchmark Reproduction

GNOSIVELA publishes its benchmark tools so third parties can reproduce results
rather than trusting vendor claims. Two tracks are available.

## 1. GKF-Bench (semantic quality)

Core repository: https://github.com/axisvela/GNOSIVELA

```bash
cd backend
GOWORK=off go test ./bench/... -count=1
```

The suite covers B1–B11 (terminology, identity, temporal, provenance, conflict,
multi-hop, policy, evolution, grounding, federation, composite). The current
baseline and gate values are in core `docs/plans/benchmark.md`.

Report:

- core commit and version;
- Go version and toolchain;
- raw output / JSON;
- any drift from the published baseline.

## 2. Scale benchmark (assertion capacity)

Use a **disposable database**. The tool refuses non-empty tables unless reset is
explicitly enabled.

```bash
createdb gnosivela_scale
export GNOSIVELA_SCALE_DSN='postgres://user:pass@host:5432/gnosivela_scale?sslmode=disable'
GNOSIVELA_PG_DSN="$GNOSIVELA_SCALE_DSN" go run ./backend/cmd/gnosivela -migrate ./backend/migrations
SCALE_COUNT=1000000 ./scripts/scale-bench.sh
```

Record:

- `SCALE_COUNT`, `SCALE_ENTITIES`, `SCALE_SAMPLES`, `SCALE_SEED`;
- PostgreSQL version and key settings (`synchronous_commit`, WAL,
  `maintenance_work_mem`);
- CPU, memory, disk class;
- import throughput, P50/P95/P99, table/index size, projection to target.

Core only supplies the tool and a small calibration. **10B execution, capacity
guarantees, HA and SLA are EE/Managed scope**; a core result is not an SLA claim.

## 3. Submitting a reproduction

Open a **Benchmark reproduction** issue with:

- repository/fork and commit;
- exact commands;
- hardware and database details;
- raw output or a link to a gist/artifact;
- deviations and limitations.

Maintainers will link accepted reproductions from the ecosystem README. Two
independent partner packages plus one external reproduction satisfy the P4-M4
validation gate when combined with the partner-list requirement.
