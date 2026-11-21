BINARY=bin/omeg4mon-router

.PHONY: all build run clean cross-build

all: build

build:
	@mkdir -p bin
	go build -ldflags="-s -w" -o $(BINARY) ./cmd/omeg4mon-router

run: build
	./$(BINARY)

cross-build:
	@mkdir -p dist
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/omeg4mon-router-linux-amd64 ./cmd/omeg4mon-router
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o dist/omeg4mon-router-linux-arm64 ./cmd/omeg4mon-router
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dist/omeg4mon-router-darwin-amd64 ./cmd/omeg4mon-router
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o dist/omeg4mon-router-darwin-arm64 ./cmd/omeg4mon-router
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o dist/omeg4mon-router-windows-amd64.exe ./cmd/omeg4mon-router

clean:
	rm -rf bin dist
