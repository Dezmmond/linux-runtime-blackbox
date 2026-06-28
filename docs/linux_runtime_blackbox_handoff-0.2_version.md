Продолжи проект Linux Runtime Blackbox.

Текущее состояние:
- `go test ./...` проходит.
- Работает команда `blackbox inspect --pid <pid>`.
- Она собирает данные из `/proc`: process metadata, status, io, fd, maps, cgroup, namespaces.
- Есть pretty output и JSON output.

Задача v0.2:
реализовать команду:

    blackbox run -- <command> [args...]

Поведение:
1. Команда запускает указанный процесс.
2. Сохраняет pid, start time, command args.
3. Делает initial snapshot через уже существующий collector.
4. Дожидается завершения процесса.
5. Сохраняет exit code.
6. Если возможно, делает final snapshot. Если `/proc/<pid>` уже исчез, добавляет warning.
7. Печатает pretty report по умолчанию.
8. Поддерживает `--output <path>` для JSON-отчёта.
9. Не добавлять eBPF, daemon mode, REST API, database или web UI.

Также добавить тестовые workload scripts:
- `scripts/workload_open_files.py`
- `scripts/workload_children.sh`
- `scripts/workload_network.py`

Добавить unit-тесты для новой логики там, где это возможно.
После реализации должны проходить:

    go test ./...
    go run ./cmd/blackbox run -- sleep 1
    go run ./cmd/blackbox run -- python3 scripts/workload_open_files.py
