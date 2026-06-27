.PHONY: test run build

test:
	go test ./...

run:
	go run ./cmd/blackbox inspect --pid $$$$ --pretty

build:
	go build -o blackbox ./cmd/blackbox
