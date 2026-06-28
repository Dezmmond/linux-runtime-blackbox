# Release Notes

## v0.2 - Command Runner

Added `blackbox run -- <command> [args...]`.

### Added

- CLI command:
  - `blackbox run -- sleep 1`
  - `blackbox run --output /tmp/blackbox-run.json -- sleep 1`
- Run metadata in JSON reports:
  - command arguments
  - PID
  - start time
  - end time
  - duration in milliseconds
  - exit code
  - signal, when applicable
- Initial procfs snapshot after process start.
- Final procfs snapshot attempt after process exit.
- `final_snapshot_unavailable` warning when `/proc/<pid>` is gone before final collection.
- Pretty output support for run reports.
- JSON output support for run reports through `--output`.
- Unit tests for command runner behavior.
- Workload scripts:
  - `scripts/workload_open_files.py`
  - `scripts/workload_children.sh`
  - `scripts/workload_network.py`
- README and Russian README usage examples for `run`.

### Notes

- v0.2 is still procfs-only.
- No eBPF, daemon mode, REST API, database, or web UI was added.
- The final snapshot is best-effort. For short-lived processes, Linux usually removes `/proc/<pid>` before the collector can read it after `Wait`, so the report records a warning instead of failing.
- If the launched command is a wrapper, shim, or shell script, the observed PID may represent that wrapper. Process tree tracking is planned for a later version.

### Verified

```bash
env GOCACHE=/tmp/go-build /usr/local/go/bin/go test ./...
env GOCACHE=/tmp/go-build /usr/local/go/bin/go run ./cmd/blackbox run -- sleep 1
env GOCACHE=/tmp/go-build /usr/local/go/bin/go run ./cmd/blackbox run -- python3 scripts/workload_open_files.py
env GOCACHE=/tmp/go-build /usr/local/go/bin/go run ./cmd/blackbox run --output /tmp/blackbox-run.json -- sleep 1
jq .schema_version /tmp/blackbox-run.json
```

## v0.1 - Static PID Snapshot

Added `blackbox inspect --pid <pid>`.

### Added

- CLI command:
  - `blackbox inspect --pid <pid>`
  - `blackbox inspect --pid <pid> --pretty`
  - `blackbox inspect --pid <pid> --output /tmp/blackbox-report.json`
- JSON report output.
- Pretty terminal report output.
- Saved report rendering:
  - `blackbox explain report.json`
- Procfs readers for:
  - `/proc/<pid>/comm`
  - `/proc/<pid>/cmdline`
  - `/proc/<pid>/exe`
  - `/proc/<pid>/cwd`
  - `/proc/<pid>/status`
  - `/proc/<pid>/stat`
  - `/proc/<pid>/io`
  - `/proc/<pid>/fd`
  - `/proc/<pid>/maps`
  - `/proc/<pid>/cgroup`
  - `/proc/<pid>/ns`
- Report model for process metadata, resources, file descriptors, mappings, namespaces, cgroups, and warnings.
- Heuristic warnings for partial procfs reads, unknown executable, deleted open files, many open file descriptors, socket file descriptors, temp executables, and large RSS.
- Parser unit tests for procfs status, stat, fd classification, and maps parsing.
- Initial README and Russian README.

### Notes

- v0.1 intentionally avoids environment variable collection and process memory reads.
- Partial procfs failures produce warnings and a best-effort report instead of crashing the CLI.
- Network sockets are classified at the file-descriptor level only; socket address enrichment is planned later.

### Verified

```bash
env GOCACHE=/tmp/go-build /usr/local/go/bin/go test ./...
env GOCACHE=/tmp/go-build /usr/local/go/bin/go run ./cmd/blackbox inspect --pid $$ --pretty
env GOCACHE=/tmp/go-build /usr/local/go/bin/go run ./cmd/blackbox inspect --pid $$ --output /tmp/blackbox-report.json
jq . /tmp/blackbox-report.json
```
