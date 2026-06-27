# Linux Runtime Blackbox

Linux Runtime Blackbox is a CLI-first runtime diagnostics agent for Linux processes.
It inspects a target process through procfs and produces a structured report with
process metadata, file descriptors, loaded mappings, namespaces, cgroups, resource
counters, and operational warnings.

The current v0.1 scope is intentionally small: no daemon, no web API, no database,
and no eBPF. It is a local procfs-based inspector.

## Build and Run

```bash
go test ./...
go run ./cmd/blackbox inspect --pid $$ --pretty
go run ./cmd/blackbox inspect --pid $$ --output /tmp/blackbox-report.json
```

Default output is JSON:

```bash
go run ./cmd/blackbox inspect --pid $$
```

Pretty output is meant for terminal inspection:

```bash
go run ./cmd/blackbox inspect --pid $$ --pretty
```

You can render a saved JSON report as text:

```bash
go run ./cmd/blackbox explain /tmp/blackbox-report.json
```

## What v0.1 Collects

- `/proc/<pid>/comm`, `cmdline`, `exe`, and `cwd`
- process status, UID/GID, thread count, RSS and virtual memory
- `/proc/<pid>/stat` for PPID and process start time ticks
- `/proc/<pid>/io` byte counters when readable
- file descriptors and their symlink targets
- memory mappings from `/proc/<pid>/maps`
- namespaces from `/proc/<pid>/ns`
- cgroup membership from `/proc/<pid>/cgroup`
- simple operational warnings

The tool handles partial procfs read failures by producing the best report it can
and adding warnings. It does not read environment variables or process memory.
