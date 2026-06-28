.PHONY: test run inspect build fmt vet

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

build:
	go build -o bin/blackbox ./cmd/blackbox

inspect:
	go run ./cmd/blackbox inspect --pid $$PPID --pretty

run:
	go run ./cmd/blackbox inspect --pid $$ --pretty
