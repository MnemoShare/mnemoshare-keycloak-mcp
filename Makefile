.PHONY: build run run-http clean tidy docker-build

BINARY  := keycloak-mcp
IMAGE   := mnemoshare/keycloak-mcp
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)

build:
	go build -ldflags="-s -w -X main.version=$(VERSION)" -o $(BINARY) ./cmd/server

run: build
	TRANSPORT=stdio ./$(BINARY)

run-http: build
	TRANSPORT=http ./$(BINARY)

tidy:
	go mod tidy

clean:
	rm -f $(BINARY)

docker-build:
	docker build -t $(IMAGE):latest .
