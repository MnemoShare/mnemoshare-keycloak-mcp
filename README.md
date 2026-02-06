# Keycloak MCP Server

A [Model Context Protocol](https://modelcontextprotocol.io) (MCP) server for Keycloak administration, written in Go. Provides 134 tools across 13 domains for managing realms, users, groups, clients, roles, identity providers, authentication flows, client scopes, sessions, authorization, components, attack detection, and server info.

## Quick Start

### Download a release

Grab the latest binary from [Releases](../../releases) for your platform:

| Platform | Binary |
|---|---|
| macOS Apple Silicon | `keycloak-mcp-darwin-arm64` |
| macOS Intel | `keycloak-mcp-darwin-amd64` |
| Linux x86_64 | `keycloak-mcp-linux-amd64` |

```bash
chmod +x keycloak-mcp-*
mv keycloak-mcp-* /usr/local/bin/keycloak-mcp
```

### Build from source

```bash
make build
```

## Configuration

All configuration is via environment variables:

| Variable | Required | Default | Description |
|---|---|---|---|
| `TRANSPORT` | No | `stdio` | Transport mode: `stdio` or `http` |
| `PORT` | No | `8080` | HTTP port (http mode only) |
| `KEYCLOAK_URL` | Yes | — | Keycloak base URL (e.g. `https://id.example.com`) |
| `KEYCLOAK_REALM` | No | `master` | Realm for authentication |
| `KEYCLOAK_AUTH_MODE` | No | `password` | Auth mode: `password` or `client_credentials` |
| `KEYCLOAK_ADMIN_USER` | For password mode | — | Admin username |
| `KEYCLOAK_ADMIN_PASSWORD` | For password mode | — | Admin password |
| `KEYCLOAK_CLIENT_ID` | For client_credentials | — | Service account client ID |
| `KEYCLOAK_CLIENT_SECRET` | For client_credentials | — | Service account client secret |
| `KEYCLOAK_DEFAULT_REALM` | No | — | Default realm for tool operations |
| `LOG_LEVEL` | No | `info` | Log level: debug, info, warn, error |
| `LOG_FORMAT` | No | `json` | Log format: `json` or `console` |

## Usage with Claude Code

Add a `.mcp.json` file to your project root:

```json
{
  "mcpServers": {
    "keycloak": {
      "command": "keycloak-mcp",
      "env": {
        "TRANSPORT": "stdio",
        "KEYCLOAK_URL": "https://id.example.com",
        "KEYCLOAK_AUTH_MODE": "client_credentials",
        "KEYCLOAK_CLIENT_ID": "mcp-admin",
        "KEYCLOAK_CLIENT_SECRET": "your-secret-here",
        "KEYCLOAK_DEFAULT_REALM": "your-realm"
      }
    }
  }
}
```

Or run it from source during development:

```json
{
  "mcpServers": {
    "keycloak": {
      "command": "go",
      "args": ["run", "./cmd/server"],
      "cwd": "/path/to/mnemoshare-keycloak-mcp",
      "env": {
        "TRANSPORT": "stdio",
        "KEYCLOAK_URL": "https://id.example.com",
        "KEYCLOAK_AUTH_MODE": "client_credentials",
        "KEYCLOAK_CLIENT_ID": "mcp-admin",
        "KEYCLOAK_CLIENT_SECRET": "your-secret-here",
        "KEYCLOAK_DEFAULT_REALM": "your-realm"
      }
    }
  }
}
```

## HTTP Mode

Run with `TRANSPORT=http` to expose a Streamable HTTP endpoint:

```bash
TRANSPORT=http KEYCLOAK_URL=https://id.example.com make run-http
```

Endpoints:
- `GET /health` — health check
- `POST /mcp` — MCP Streamable HTTP endpoint

## Keycloak Setup

The server needs a service account with admin privileges. Create one in Keycloak:

1. Go to **master** realm > **Clients** > **Create client**
2. Set Client ID (e.g. `mcp-admin`), enable **Client authentication**
3. Under **Service account roles**, assign the `admin` realm role
4. Copy the client secret from the **Credentials** tab

## Docker

```bash
make docker-build
docker run -e KEYCLOAK_URL=https://id.example.com \
  -e KEYCLOAK_AUTH_MODE=client_credentials \
  -e KEYCLOAK_CLIENT_ID=mcp-admin \
  -e KEYCLOAK_CLIENT_SECRET=secret \
  keycloak-mcp:latest
```

## Tools

134 tools across 13 domains:

- **Realms** (8) — list, get, create, update, delete, clear caches
- **Users** (24) — CRUD, password, credentials, groups, sessions, federated identities, roles
- **Groups** (12) — CRUD, members, count, realm/client role mappings
- **Clients** (18) — CRUD, secrets, service accounts, scopes, protocol mappers, sessions
- **Roles** (16) — realm + client role CRUD, composites, user/group lookups
- **Identity Providers** (8) — CRUD + mappers
- **Authentication Flows** (10) — flows, executions, required actions
- **Client Scopes** (10) — CRUD + protocol mappers + realm defaults
- **Sessions** (5) — logout, events, offline sessions, consent revocation
- **Authorization** (15) — resources, scopes, policies, permissions
- **Components** (5) — CRUD for user federation, LDAP, custom providers
- **Attack Detection** (2) — brute force status + clear
- **Server Info** (1) — Keycloak server info

## License

MIT
