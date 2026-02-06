# Keycloak MCP Server

A Model Context Protocol (MCP) server for Keycloak administration, written in Go.

## Build & Run

```bash
make build         # build binary
make run           # run in stdio mode
make run-http      # run in HTTP mode on :8080
make tidy          # go mod tidy
make docker-build  # build Docker image
```

## Configuration

All configuration is via environment variables. See `.env.example` for the full list.

Key variables:
- `TRANSPORT` — `stdio` (default) or `http`
- `KEYCLOAK_URL` — Keycloak base URL
- `KEYCLOAK_AUTH_MODE` — `password` or `client_credentials`
- `KEYCLOAK_DEFAULT_REALM` — optional default realm for tools

## Architecture

- `cmd/server/main.go` — entry point, dual transport (stdio + Streamable HTTP)
- `internal/config/` — env-based configuration
- `internal/auth/` — token manager with auto-refresh
- `internal/keycloak/` — thin gocloak wrapper with token injection
- `internal/tools/` — 134 MCP tools across 13 domains

## Tool Pattern

Each tool file follows the same pattern:
1. Arg struct with `json` + `jsonschema` tags
2. `registerXxxTools(s *mcp.Server, kc *keycloak.Client)` function
3. `mcp.AddTool(s, &mcp.Tool{...}, handler)` with typed generic handler
4. Domain errors → `toolError()`, transport errors → Go `error`

## Dependencies

- `github.com/modelcontextprotocol/go-sdk` — official MCP Go SDK
- `github.com/Nerzal/gocloak/v13` — Keycloak admin client
- `github.com/rs/zerolog` — structured logging
