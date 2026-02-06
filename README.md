# Keycloak MCP Server

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![MCP](https://img.shields.io/badge/MCP-Compatible-purple)](https://modelcontextprotocol.io)

A comprehensive [Model Context Protocol](https://modelcontextprotocol.io) (MCP) server for Keycloak administration. Provides **134 tools across 13 domains** — manage realms, users, groups, clients, roles, identity providers, authentication flows, and more, all from your AI assistant.

Built and maintained by [MnemoShare](https://mnemoshare.com) — a HIPAA-compliant secure file transfer platform.

## Features

- **134 admin tools** covering the full Keycloak Admin REST API
- **Two transport modes** — stdio (for Claude Code, Cursor, etc.) and HTTP (for remote/container deployments)
- **Two auth modes** — admin password or client credentials (service account)
- **Automatic token refresh** — handles Keycloak token lifecycle transparently
- **Zero configuration files** — everything via environment variables
- **Single binary** — no runtime dependencies

## Quick Start

### Install from release

Download the latest binary from [Releases](https://github.com/MnemoShare/mnemoshare-keycloak-mcp/releases) for your platform:

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
git clone https://github.com/MnemoShare/mnemoshare-keycloak-mcp.git
cd mnemoshare-keycloak-mcp
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
| `LOG_LEVEL` | No | `info` | Log level: `debug`, `info`, `warn`, `error` |
| `LOG_FORMAT` | No | `json` | Log format: `json` or `console` |

## Usage

### Claude Code

Add a `.mcp.json` file to your project root:

```json
{
  "mcpServers": {
    "keycloak": {
      "command": "keycloak-mcp",
      "env": {
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

Or run from source during development:

```json
{
  "mcpServers": {
    "keycloak": {
      "command": "go",
      "args": ["run", "./cmd/server"],
      "cwd": "/path/to/mnemoshare-keycloak-mcp",
      "env": {
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

### HTTP Mode

Run with HTTP transport to expose a Streamable HTTP endpoint — useful for containerized or remote deployments:

```bash
TRANSPORT=http KEYCLOAK_URL=https://id.example.com make run-http
```

Endpoints:
- `GET /health` — health check
- `POST /mcp` — MCP Streamable HTTP endpoint

### Docker

```bash
make docker-build
docker run -e KEYCLOAK_URL=https://id.example.com \
  -e KEYCLOAK_AUTH_MODE=client_credentials \
  -e KEYCLOAK_CLIENT_ID=mcp-admin \
  -e KEYCLOAK_CLIENT_SECRET=your-secret \
  keycloak-mcp:latest
```

## Keycloak Setup

### Option 1: Bootstrap with password mode (fresh install)

On a fresh Keycloak instance you won't have a service account yet. Start with password mode so your AI assistant can connect and create one for you:

```json
{
  "mcpServers": {
    "keycloak": {
      "command": "keycloak-mcp",
      "env": {
        "KEYCLOAK_URL": "https://id.example.com",
        "KEYCLOAK_AUTH_MODE": "password",
        "KEYCLOAK_ADMIN_USER": "admin",
        "KEYCLOAK_ADMIN_PASSWORD": "your-admin-password",
        "KEYCLOAK_DEFAULT_REALM": "your-realm"
      }
    }
  }
}
```

Once connected, ask your AI assistant to:

1. Create a confidential client (e.g. `mcp-admin`) in the **master** realm
2. Assign the `admin` realm role to its service account
3. Retrieve the client secret

Then switch to client credentials mode (Option 2) and remove the admin password from your config.

### Option 2: Client credentials mode (recommended)

For day-to-day use, configure a service account client. To set it up manually in the Keycloak admin console:

1. Go to **master** realm > **Clients** > **Create client**
2. Set Client ID (e.g. `mcp-admin`), enable **Client authentication**
3. Under **Service account roles**, assign the `admin` realm role
4. Copy the client secret from the **Credentials** tab

This is more secure for ongoing use — the service account has a rotatable secret and doesn't expose your admin credentials.

## Tools

134 tools across 13 domains:

| Domain | Tools | Description |
|---|---|---|
| **Realms** | 8 | List, get, create, update, delete, clear caches |
| **Users** | 24 | CRUD, passwords, credentials, groups, sessions, federated identities, roles |
| **Groups** | 12 | CRUD, members, count, realm/client role mappings |
| **Clients** | 18 | CRUD, secrets, service accounts, scopes, protocol mappers, sessions |
| **Roles** | 16 | Realm + client role CRUD, composites, user/group lookups |
| **Identity Providers** | 8 | CRUD + mappers |
| **Authentication Flows** | 10 | Flows, executions, required actions |
| **Client Scopes** | 10 | CRUD + protocol mappers + realm defaults |
| **Sessions** | 5 | Logout, events, offline sessions, consent revocation |
| **Authorization** | 15 | Resources, scopes, policies, permissions |
| **Components** | 5 | CRUD for user federation, LDAP, custom providers |
| **Attack Detection** | 2 | Brute force status + clear |
| **Server Info** | 1 | Keycloak server info |

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/my-feature`)
3. Commit your changes (`git commit -am 'Add my feature'`)
4. Push to the branch (`git push origin feature/my-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License — see the [LICENSE](LICENSE) file for details.

## About MnemoShare

[MnemoShare](https://mnemoshare.com) is a HIPAA-compliant, HITRUST-certified secure file transfer platform. We built this MCP server to manage our own Keycloak infrastructure and are sharing it with the community because we believe AI-powered administration should be accessible to everyone.
