BINARY := bin/bible
PACKAGE := ./cmd/bible
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
BUILD_DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -s -w \
	-X github.com/vmrocha/bible-cli/internal/buildinfo.version=$(VERSION) \
	-X github.com/vmrocha/bible-cli/internal/buildinfo.commit=$(COMMIT) \
	-X github.com/vmrocha/bible-cli/internal/buildinfo.date=$(BUILD_DATE)

.PHONY: build check clean fmt format-check lint test

build:
	go build -trimpath -ldflags "$(LDFLAGS)" -o $(BINARY) $(PACKAGE)

test:
	go test -race ./...

lint:
	go vet ./...

fmt:
	go fmt ./...

format-check:
	@test -z "$$(gofmt -l .)" || (gofmt -l . && exit 1)

check: format-check lint test build

clean:
	go clean
	rm -rf bin
