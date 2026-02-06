.PHONY: build run run-http clean tidy docker-build

BINARY := keycloak-mcp
IMAGE  := mnemoshare/keycloak-mcp

build:
	go build -o $(BINARY) ./cmd/server

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
