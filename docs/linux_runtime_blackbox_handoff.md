# Linux Runtime Blackbox — handoff for Codex

Generated for a fresh project start.
Language: English for tool/agent clarity, with Russian comments where helpful.
Audience: Codex or another coding agent working with the user.

---

## 0. Project identity

Working title: **Linux Runtime Blackbox**

Short idea:

> A Linux runtime diagnostics agent that observes what processes actually do on a machine: process tree, file descriptors, loaded libraries, network sockets, namespaces, cgroups, resource usage, and later eBPF-based runtime events.

This is **not** a web/backend demo. The core value is a system-level Linux agent and CLI. Any API/UI/dashboard is secondary.

The user is a Python/Go systems/backend developer with strong interest in Linux internals, enterprise software, backup/storage agents, packaging, OpenStack/Cinder-like infrastructure, corporate Linux, and low-level debugging. They explicitly do **not** want another boring CRUD/backend dashboard project.

Project vibe:

- practical system software;
- Linux internals first;
- corporate diagnostic tool;
- minimal useless abstractions;
- CLI-first;
- observable, testable, explainable;
- small steps, real working artifacts;
- avoid architecture cosplay.

---

## 1. Primary goal

Build a local Linux agent/CLI that can answer:

```text
What did this process do while it was running?
What files did it open?
What sockets did it use?
What child processes did it create?
Which shared libraries did it load?
Which cgroup/namespace/systemd context did it run in?
How much CPU, memory, and IO did it consume?
What looks suspicious or operationally important?
```

The first usable product is a CLI that produces a structured JSON report for a target process or command.

Example future usage:

```bash
blackbox inspect --pid 1234 --output report.json
blackbox run --duration 30s --output report.json -- /usr/bin/python3 script.py
blackbox explain report.json
blackbox top --sort io
```

---

## 2. Non-goals for the first phase

Do **not** start with:

- web UI;
- Django/FastAPI;
- OpenTelemetry exporter;
- Kubernetes;
- persistent server/daemon;
- complex eBPF tracing;
- machine learning;
- database schema;
- distributed architecture;
- enterprise authorization;
- plugin framework.

These may come later. v0.1 must be boringly useful and locally runnable.

---

## 3. Recommended language and rationale

Recommended implementation language for v0.1: **Go**.

Why Go:

- user knows Go;
- produces easy-to-distribute static-ish binaries;
- good fit for CLI agents;
- good standard library for filesystem/process parsing;
- later integration with `github.com/cilium/ebpf` is natural;
- easier operational deployment than Python for a Linux agent.

Python is allowed for tests, fixtures, helper scripts, and synthetic workload generators.

Rust is interesting, but do not switch the core to Rust unless the user explicitly asks. Rust may be considered later if memory safety and deeper systems work become the main educational goal.

---

## 4. Development environment recommendation

### Preferred environment

Use a dedicated Linux VM, not the daily host.

Reason: the project will inspect `/proc`, `/sys`, namespaces, cgroups, capabilities, and later eBPF. Some commands will need root or elevated capabilities. A VM gives safety and reproducibility.

Suggested VM:

- Ubuntu Server 24.04 LTS or Debian 12/13;
- 2–4 vCPU;
- 4–8 GB RAM;
- 30+ GB disk;
- kernel with BTF available at `/sys/kernel/btf/vmlinux` for future eBPF CO-RE work;
- snapshot before eBPF experiments.

### Avoid as primary environment

- Do not develop the first low-level tests inside an unprivileged container.
- Containers are fine for formatting/linting/build checks, but not as the primary runtime lab for `/proc`, cgroups, namespaces, or eBPF.
- Do not run early experiments on the user's main host unless the command is read-only and does not require root.

---

## 5. Initial repository layout

Create this structure:

```text
linux-runtime-blackbox/
├── README.md
├── Makefile
├── go.mod
├── go.sum
├── cmd/
│   └── blackbox/
│       └── main.go
├── internal/
│   ├── collector/
│   │   ├── collector.go
│   │   └── snapshot.go
│   ├── procfs/
│   │   ├── process.go
│   │   ├── fd.go
│   │   ├── maps.go
│   │   ├── stat.go
│   │   ├── status.go
│   │   ├── cgroup.go
│   │   ├── namespaces.go
│   │   └── sockets.go
│   ├── model/
│   │   ├── process.go
│   │   ├── report.go
│   │   ├── fd.go
│   │   ├── socket.go
│   │   └── warning.go
│   ├── report/
│   │   ├── json.go
│   │   └── explain.go
│   └── run/
│       └── command.go
├── ebpf/
│   └── README.md
├── scripts/
│   ├── workload_open_files.py
│   ├── workload_network.py
│   └── workload_children.sh
├── testdata/
│   └── README.md
└── docs/
    ├── architecture.md
    ├── roadmap.md
    └── linux-notes.md
```

Notes:

- `internal/procfs` is responsible only for parsing Linux procfs/sysfs data.
- `internal/collector` combines raw Linux data into snapshots/reports.
- `internal/model` contains stable domain structures used by output and tests.
- `internal/report` renders JSON and human-readable explanations.
- `internal/run` later runs a command under observation.
- `ebpf/` stays mostly empty in v0.1 except for notes. Do not prematurely implement eBPF.

---

## 6. v0.1 feature scope

v0.1 target: inspect an existing PID and produce JSON.

Command:

```bash
blackbox inspect --pid <PID> --output report.json
```

Minimal fields:

```json
{
  "schema_version": "0.1",
  "collected_at": "2026-06-27T20:00:00Z",
  "target": {
    "pid": 1234,
    "comm": "python3",
    "exe": "/usr/bin/python3.12",
    "cmdline": ["python3", "app.py"],
    "cwd": "/home/user/project"
  },
  "process": {
    "pid": 1234,
    "ppid": 1,
    "uid": 1000,
    "gid": 1000,
    "state": "S",
    "threads": 4,
    "start_time_ticks": 1234567
  },
  "resources": {
    "vmsize_kb": 123456,
    "vmrss_kb": 23456,
    "read_bytes": 1024,
    "write_bytes": 2048
  },
  "fds": [
    {
      "fd": 1,
      "target": "pipe:[12345]",
      "kind": "pipe"
    }
  ],
  "maps": [
    {
      "path": "/usr/lib/x86_64-linux-gnu/libc.so.6",
      "perms": "r-xp"
    }
  ],
  "namespaces": {
    "mnt": "mnt:[4026531841]",
    "pid": "pid:[4026531836]",
    "net": "net:[4026531840]"
  },
  "cgroups": [
    {
      "hierarchy": "0",
      "controllers": "",
      "path": "/user.slice/user-1000.slice/session-2.scope"
    }
  ],
  "warnings": []
}
```

---

## 7. v0.1 data sources

Use these files under `/proc/<pid>`:

```text
/proc/<pid>/comm
/proc/<pid>/cmdline
/proc/<pid>/cwd -> symlink
/proc/<pid>/exe -> symlink
/proc/<pid>/status
/proc/<pid>/stat
/proc/<pid>/io
/proc/<pid>/fd/* -> symlinks
/proc/<pid>/maps
/proc/<pid>/cgroup
/proc/<pid>/ns/* -> symlinks
/proc/<pid>/task/*
```

Optional for network socket resolution:

```text
/proc/net/tcp
/proc/net/tcp6
/proc/net/udp
/proc/net/udp6
/proc/<pid>/net/tcp
/proc/<pid>/net/tcp6
/proc/<pid>/net/udp
/proc/<pid>/net/udp6
```

Important Linux realities:

- processes can disappear while reading `/proc`; handle `ENOENT` gracefully;
- permissions can prevent reading some files; report partial data, do not crash;
- `/proc/<pid>/cmdline` is NUL-separated;
- `/proc/<pid>/stat` parsing is tricky because `comm` is inside parentheses and may contain spaces;
- `fd` symlink targets may be files, sockets, pipes, anon_inode, deleted files;
- `/proc/<pid>/maps` can be long;
- containerized processes may show different namespace/cgroup context;
- never use `ldd` on untrusted binaries in this project path; it can execute code depending on context. For loaded libraries, use `/proc/<pid>/maps` in v0.1.

---

## 8. Warning rules for v0.1

Implement simple heuristic warnings. These are not security claims; they are operational hints.

Examples:

```text
deleted_file_open
  A file descriptor points to a path ending with " (deleted)".

many_open_fds
  Process has more than N open file descriptors.

unexpected_tmp_executable
  exe path is under /tmp, /var/tmp, /dev/shm, or user-writable temp-like directory.

large_rss
  RSS is above configurable threshold.

network_socket_open
  Process has socket file descriptors.

unknown_exe
  /proc/<pid>/exe cannot be resolved.

permission_denied_partial_report
  Some procfs files could not be read.
```

Each warning should include:

```go
type Warning struct {
    ID       string `json:"id"`
    Severity string `json:"severity"` // info, low, medium, high
    Message  string `json:"message"`
    Evidence map[string]string `json:"evidence,omitempty"`
}
```

Do not overclaim. Use wording like "may indicate", "observed", "could mean".

---

## 9. CLI design

Use simple, predictable commands.

Initial commands:

```bash
blackbox inspect --pid 1234
blackbox inspect --pid 1234 --output report.json
blackbox inspect --pid 1234 --pretty
blackbox explain report.json
```

Near-future commands:

```bash
blackbox run --output report.json -- /bin/bash -c 'echo hello'
blackbox watch --pid 1234 --interval 1s --duration 30s
blackbox tree --pid 1234
blackbox version
```

CLI library: keep it simple. Standard `flag` is acceptable for v0.1. Cobra can be added later if command complexity grows.

Exit codes:

```text
0 success
1 generic error
2 invalid CLI usage
3 target process not found
4 insufficient permissions
```

---

## 10. Testing strategy

Do not require root for v0.1 tests.

Test layers:

1. Unit tests for procfs parsers using fixture files.
2. Integration tests that inspect the current test process or a spawned helper process.
3. Script-based workloads under `scripts/` for manual testing.

Examples:

```bash
go test ./...
go run ./cmd/blackbox inspect --pid $$ --pretty
python3 scripts/workload_open_files.py &
blackbox inspect --pid $! --pretty
```

Synthetic workloads:

- process that opens files and sleeps;
- process that opens TCP listener and sleeps;
- process that spawns child processes;
- process that writes to stdout/stderr;
- process that keeps a deleted file open.

Important test cases:

- PID does not exist;
- permission denied;
- process exits during collection;
- weird `/proc/<pid>/stat` comm value;
- fd target is socket/pipe/anon_inode/file/deleted file;
- empty cmdline;
- zombie process.

---

## 11. Architecture principles

Keep collectors small and composable.

Preferred flow:

```text
CLI command
  -> parse target
  -> collector.CollectPID(pid)
  -> procfs readers
  -> model.Report
  -> report JSON/human renderer
```

Rules:

- procfs package should not print to stdout.
- collector should return partial reports plus structured errors/warnings.
- report rendering should not read Linux files.
- CLI should be thin.
- Avoid global state.
- Avoid goroutine complexity until watch mode.
- Every external Linux read should be wrapped with context-rich error handling.

---

## 12. Security and safety requirements

This project observes system state. Be careful.

For Codex:

- Do not add destructive commands.
- Do not write to `/proc` or `/sys` in v0.1.
- Do not change system settings.
- Do not require root unless explicitly needed.
- Do not collect environment variables by default because they may contain secrets.
- Do not dump full memory, credentials, tokens, SSH keys, or private files.
- Do not upload reports anywhere.
- Redact sensitive paths only if the user requests; otherwise keep local output truthful.
- Never claim the tool is a malware detector. It is a runtime diagnostics and observability tool.

---

## 13. Future eBPF plan

Do not implement eBPF in v0.1. Prepare mentally for it.

Future eBPF scope:

```text
trace execve
trace openat/openat2
trace connect/accept
trace process exit
trace file rename/unlink maybe later
trace short-lived processes missed by /proc polling
```

Recommended Go path:

- `github.com/cilium/ebpf`;
- `bpf2go`;
- small C eBPF programs compiled with clang;
- CO-RE if available;
- ring buffer for events.

Potential event model:

```go
type RuntimeEvent struct {
    TimestampUnixNano int64             `json:"timestamp_unix_nano"`
    Type              string            `json:"type"` // exec, open, connect, exit
    PID               int               `json:"pid"`
    TGID              int               `json:"tgid"`
    PPID              int               `json:"ppid,omitempty"`
    Comm              string            `json:"comm"`
    Args              map[string]string `json:"args,omitempty"`
}
```

Important future concern:

- eBPF verifier constraints;
- kernel version differences;
- capabilities/root requirements;
- BTF availability;
- event loss under load;
- ring buffer backpressure;
- short-lived process handling;
- container namespace mapping.

---

## 14. Roadmap

### v0.1 — Static PID snapshot

- initialize repo;
- implement CLI `inspect --pid`;
- collect basic process metadata;
- collect fd list;
- collect memory/io/resource info;
- collect maps-based loaded libraries;
- collect namespaces and cgroups;
- output JSON;
- add human `--pretty` output;
- add unit tests for parsers.

### v0.2 — Command runner

- `blackbox run -- <command>`;
- start command;
- inspect while running;
- wait for exit;
- record exit code and duration;
- save final report.

### v0.3 — Watch mode

- periodic snapshots;
- fd delta;
- memory/io delta;
- process tree delta;
- short text summary.

### v0.4 — Process tree and systemd awareness

- collect descendants;
- map cgroups to likely systemd unit/scope;
- support `--unit name.service` if feasible.

### v0.5 — Network socket enrichment

- map socket inode from fd to `/proc/net/*`;
- decode local/remote addresses;
- TCP states;
- show listening sockets.

### v0.6 — eBPF MVP

- trace `execve` and process exit;
- add event timeline;
- keep procfs snapshot as source of stable metadata.

### v0.7 — eBPF file/network events

- trace `openat`/`openat2`;
- trace `connect`;
- correlate runtime events with snapshots.

### v0.8 — Exporters and integrations

- JSONL event stream;
- optional OpenTelemetry exporter;
- optional local SQLite storage;
- optional HTML report.

---

## 15. First implementation task for Codex

Task:

> Initialize the Go repository for Linux Runtime Blackbox and implement v0.1 skeleton: CLI command `inspect --pid <pid>`, procfs readers for comm/cmdline/exe/cwd/status/io/fd/maps/cgroup/ns, JSON report output, and parser tests where practical.

Acceptance criteria:

```bash
go test ./...
go run ./cmd/blackbox inspect --pid $$ --pretty
go run ./cmd/blackbox inspect --pid $$ --output /tmp/blackbox-report.json
jq . /tmp/blackbox-report.json
```

Expected behavior:

- works without root for processes owned by current user;
- gracefully handles partial read errors;
- produces valid JSON;
- does not write outside requested output path;
- does not require eBPF;
- does not require network;
- includes README with quick local usage.

---

## 16. Suggested README opening

```markdown
# Linux Runtime Blackbox

Linux Runtime Blackbox is a CLI-first runtime diagnostics agent for Linux processes.
It inspects a target process through procfs and produces a structured report with process metadata, file descriptors, loaded mappings, namespaces, cgroups, resource counters, and operational warnings.

The goal is to build a practical Linux systems project that can grow from procfs snapshots into eBPF-based runtime tracing.
```

---

## 17. Style preferences

The user prefers practical, direct, technical work.

When making implementation decisions:

- prefer working code over theoretical architecture;
- explain Linux details when they matter;
- avoid fluffy product language;
- keep commits small;
- keep each step runnable;
- do not introduce large frameworks without need;
- do not start with daemonization or web UI;
- make failures visible and understandable.

---

## 18. Good commit sequence

Recommended first commits:

```text
1. init go module and README
2. add inspect CLI skeleton
3. add procfs process metadata reader
4. add fd reader
5. add status/io parser
6. add maps/cgroup/ns reader
7. add report model and JSON output
8. add pretty output and warnings
9. add parser tests and fixtures
```

---

## 19. Manual smoke commands

```bash
# inspect current shell
./blackbox inspect --pid $$ --pretty

# inspect a sleeping process
sleep 60 &
PID=$!
./blackbox inspect --pid "$PID" --pretty
kill "$PID"

# open-file workload
python3 scripts/workload_open_files.py &
PID=$!
./blackbox inspect --pid "$PID" --pretty
kill "$PID"

# deleted-file workload idea
python3 - <<'PY' &
import os, tempfile, time
f = tempfile.NamedTemporaryFile(delete=False)
path = f.name
f.write(b'hello')
f.flush()
os.unlink(path)
time.sleep(60)
PY
PID=$!
./blackbox inspect --pid "$PID" --pretty
kill "$PID"
```

---

## 20. Definition of done for the first evening

A successful first evening is not a perfect agent.

It is enough if:

- repo exists;
- `go test ./...` runs;
- `blackbox inspect --pid $$ --pretty` prints useful data;
- JSON output can be saved;
- at least fd/cmdline/exe/cwd/status are collected;
- README explains how to run it.

That is already a living project.
