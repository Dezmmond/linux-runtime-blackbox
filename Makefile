GO=/snap/bin/go

.PHONY: test run inspect build fmt vet

test:
	$(GO) test ./...

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

run-command:
	$(GO) run ./cmd/blackbox run -- sleep 1

build:
	$(GO) build -o bin/blackbox ./cmd/blackbox

inspect:
	$(GO) run ./cmd/blackbox inspect --pid $$PPID --pretty

run:
	$(GO) run ./cmd/blackbox inspect --pid $$ --pretty
