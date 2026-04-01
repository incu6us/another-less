VERSION ?= $(shell git describe --tags 2>/dev/null || echo "dev")
LDFLAGS = -ldflags "-X main.version=$(VERSION)"

.PHONY: build test lint smoke clean

build:
	go build $(LDFLAGS) -o bin/aless ./cmd/aless

test:
	go test ./...

lint:
	golangci-lint run ./...

smoke: build
	echo "hello world" | ./bin/aless --version
	@echo "Smoke test passed"

clean:
	rm -rf bin/
