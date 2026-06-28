# Linux Runtime Blackbox

Linux Runtime Blackbox - CLI-first агент диагностики Linux-процессов.
Он читает данные из procfs и формирует структурированный отчет о процессе:
метаданные, файловые дескрипторы, memory mappings, namespaces, cgroups,
ресурсные счетчики и простые операционные предупреждения.

Текущий v0.1 намеренно небольшой: без daemon mode, web API, базы данных и eBPF.
Это локальный procfs-based инспектор, который уже можно запускать на обычном
пользовательском процессе без root.

## Быстрый старт

```bash
go test ./...
go run ./cmd/blackbox inspect --pid $$ --pretty
go run ./cmd/blackbox inspect --pid $$ --output /tmp/blackbox-report.json
go run ./cmd/blackbox run -- sleep 1
```

По умолчанию `inspect` печатает JSON в stdout:

```bash
go run ./cmd/blackbox inspect --pid $$
```

Для человекочитаемого отчета используйте `--pretty`:

```bash
go run ./cmd/blackbox inspect --pid $$ --pretty
```

Чтобы сохранить JSON-отчет в файл:

```bash
go run ./cmd/blackbox inspect --pid $$ --output /tmp/blackbox-report.json
```

Сохраненный JSON можно вывести в pretty-формате:

```bash
go run ./cmd/blackbox explain /tmp/blackbox-report.json
```

Запустить команду под наблюдением:

```bash
go run ./cmd/blackbox run -- python3 scripts/workload_open_files.py
go run ./cmd/blackbox run --output /tmp/blackbox-run.json -- sleep 1
```

## Что собирает v0.1

- `/proc/<pid>/comm`, `cmdline`, `exe`, `cwd`
- UID/GID, состояние процесса, количество потоков
- виртуальную память и RSS из `/proc/<pid>/status`
- PPID и `start_time_ticks` из `/proc/<pid>/stat`
- `read_bytes` и `write_bytes` из `/proc/<pid>/io`, если доступны
- файловые дескрипторы из `/proc/<pid>/fd`
- memory mappings из `/proc/<pid>/maps`
- namespaces из `/proc/<pid>/ns`
- cgroup membership из `/proc/<pid>/cgroup`
- простые warnings по эвристикам
- v0.2 run metadata: команда, PID, время старта/завершения, длительность и exit code

Инструмент не читает environment variables и память процесса. Если часть procfs
недоступна из-за прав, zombie-состояния или гонки с завершением процесса,
отчет все равно формируется, а проблема добавляется в `warnings`.

## Пример

```bash
go run ./cmd/blackbox inspect --pid $$ --pretty
```

Примерные секции отчета:

```text
Target
Process
Resources
File descriptors
Mappings
Namespaces
Cgroups
Warnings
```

## Коды выхода

- `0` - успех
- `1` - общая ошибка
- `2` - ошибка CLI usage
- `3` - целевой процесс не найден
- `4` - недостаточно прав

## Ограничения v0.1

- нет eBPF и событий runtime;
- нет daemon/watch mode;
- нет web API и UI;
- нет базы данных;
- network sockets пока определяются только как socket fd, без расшифровки
  адресов из `/proc/net/*`;
- короткоживущие процессы можно пропустить, потому что v0.1 делает статический
  snapshot существующего PID.

## Следующие шаги

Ближайшие логичные расширения:

- `blackbox run -- <command>` для запуска команды под наблюдением;
- watch mode с периодическими snapshot'ами;
- process tree и descendants;
- обогащение socket fd через `/proc/<pid>/net/*`;
- позже - eBPF events для `execve`, `openat`, `connect` и process exit.
