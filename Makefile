SHELL := /usr/bin/env bash
export PATH := $(HOME)/.local/go/bin:$(HOME)/.local/bin:$(HOME)/.bun/bin:$(PATH)

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
BUILD_TIME ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildTime=$(BUILD_TIME)
DEV_LDFLAGS := -X main.version=$(VERSION)-dev -X main.commit=$(COMMIT) -X main.buildTime=$(BUILD_TIME)

BINARY := bin/omeg4mon-router
DEV_BINARY := bin/omeg4mon-router-dev
UI_DIST := cmd/omeg4mon-router/ui_dist
ARGS ?=

.PHONY: codegen fmt fmt-check lint check test build-web build smoke dev dev-bin prod cross-build clean

codegen:
	cd web && bun run gen:api

fmt:
	golangci-lint fmt ./...
	cd web && bun run check:fix

fmt-check:
	@unformatted=$$(gofmt -s -l cmd internal); \
	if [ -n "$$unformatted" ]; then \
		echo "Unformatted Go files:"; echo "$$unformatted"; exit 1; \
	fi
	golangci-lint fmt --diff ./...
	cd web && bun run format:check

lint:
	golangci-lint run ./...
	cd web && bun run check

check: fmt-check lint test

test:
	go test -count=1 -race ./...
	CGO_ENABLED=0 go test -count=1 ./internal/platform/db/...
	cd web && bun test && bun run typecheck
	cd cli && bun test

build-web:
	cd web && bun run build
	rm -rf $(UI_DIST)
	cp -R web/out $(UI_DIST)

build: codegen build-web
	mkdir -p bin
	CGO_ENABLED=0 go build -tags embedui -trimpath -ldflags="$(LDFLAGS)" -o $(BINARY) ./cmd/omeg4mon-router

smoke:
	./scripts/smoke-test.sh $(BINARY)

dev: codegen
	@mkdir -p bin
	@CGO_ENABLED=0 go build -ldflags="$(DEV_LDFLAGS)" -o $(DEV_BINARY) ./cmd/omeg4mon-router
	@set -euo pipefail; \
	HOST="$${OMEG4MON_HOST:-127.0.0.1}"; \
	PORT="$${OMEG4MON_PORT:-20128}"; \
	LOG_LEVEL="$${OMEG4MON_LOG_LEVEL:-debug}"; \
	BACKEND_PID=""; WEB_PID=""; \
	cleanup() { \
		trap - INT TERM EXIT; \
		if [ -n "$$WEB_PID" ]; then \
			pkill -TERM -P "$$WEB_PID" 2>/dev/null || true; \
			kill -TERM "$$WEB_PID" 2>/dev/null || true; \
		fi; \
		if [ -n "$$BACKEND_PID" ]; then \
			kill -TERM "$$BACKEND_PID" 2>/dev/null || true; \
		fi; \
		wait 2>/dev/null || true; \
	}; \
	trap cleanup INT TERM EXIT; \
	echo "Starting omeg4mon-router dev cycle (API: http://$${HOST}:$${PORT}, UI: http://127.0.0.1:20127)..."; \
	OMEG4MON_HOST="$$HOST" OMEG4MON_PORT="$$PORT" OMEG4MON_LOG_LEVEL="$$LOG_LEVEL" ./$(DEV_BINARY) $(ARGS) & \
	BACKEND_PID=$$!; \
	(cd web && exec env OMEG4MON_BACKEND_URL="http://$${HOST}:$${PORT}" bun run dev) & \
	WEB_PID=$$!; \
	wait -n "$$BACKEND_PID" "$$WEB_PID"

dev-bin: codegen build-web
	mkdir -p bin
	CGO_ENABLED=0 go build -tags embedui -ldflags="$(DEV_LDFLAGS)" -o $(DEV_BINARY) ./cmd/omeg4mon-router
	OMEG4MON_LOG_LEVEL="$${OMEG4MON_LOG_LEVEL:-debug}" exec ./$(DEV_BINARY) $(ARGS)

prod: check build smoke
	exec ./$(BINARY) $(ARGS)

cross-build: codegen build-web
	mkdir -p dist
	rm -f dist/omeg4mon-router-* dist/SHA256SUMS.txt
	GOOS=linux   GOARCH=amd64 CGO_ENABLED=0 go build -tags embedui -trimpath -ldflags="$(LDFLAGS)" -o dist/omeg4mon-router-linux-amd64 ./cmd/omeg4mon-router
	GOOS=linux   GOARCH=arm64 CGO_ENABLED=0 go build -tags embedui -trimpath -ldflags="$(LDFLAGS)" -o dist/omeg4mon-router-linux-arm64 ./cmd/omeg4mon-router
	GOOS=darwin  GOARCH=amd64 CGO_ENABLED=0 go build -tags embedui -trimpath -ldflags="$(LDFLAGS)" -o dist/omeg4mon-router-darwin-amd64 ./cmd/omeg4mon-router
	GOOS=darwin  GOARCH=arm64 CGO_ENABLED=0 go build -tags embedui -trimpath -ldflags="$(LDFLAGS)" -o dist/omeg4mon-router-darwin-arm64 ./cmd/omeg4mon-router
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -tags embedui -trimpath -ldflags="$(LDFLAGS)" -o dist/omeg4mon-router-windows-amd64.exe ./cmd/omeg4mon-router
	GOOS=windows GOARCH=arm64 CGO_ENABLED=0 go build -tags embedui -trimpath -ldflags="$(LDFLAGS)" -o dist/omeg4mon-router-windows-arm64.exe ./cmd/omeg4mon-router
	cd dist && sha256sum omeg4mon-router-* > SHA256SUMS.txt

clean:
	rm -rf bin dist $(UI_DIST) web/out
