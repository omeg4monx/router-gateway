export PATH := $(HOME)/.local/go/bin:$(HOME)/.local/bin:$(HOME)/.bun/bin:$(PATH)

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
BUILD_TIME ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildTime=$(BUILD_TIME)

BINARY := bin/omeg4mon-router
UI_DIST := cmd/omeg4mon-router/ui_dist

.PHONY: codegen fmt fmt-check lint check test build-web build smoke dev cross-build clean

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
	cd web && NEXT_OUTPUT_EXPORT=true bun run build
	rm -rf $(UI_DIST)
	cp -R web/out $(UI_DIST)

build: codegen build-web
	mkdir -p bin
	CGO_ENABLED=0 go build -tags embedui -trimpath -ldflags="$(LDFLAGS)" -o $(BINARY) ./cmd/omeg4mon-router

smoke:
	./scripts/smoke-test.sh $(BINARY)

dev:
	@echo "Start omeg4mon-router in two terminals for live reload:"
	@echo "  1) Go API backend  : go run ./cmd/omeg4mon-router"
	@echo "  2) Next.js dev UI  : cd web && bun run dev"
	@echo "Then open http://127.0.0.1:20127 (API proxied to http://127.0.0.1:20128)."

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
	rm -rf bin dist $(UI_DIST) web/out web/.next
