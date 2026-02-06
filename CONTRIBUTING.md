# Contributing to Keycloak MCP Server

Thanks for your interest in contributing! This guide will help you get started.

## Getting Started

### Prerequisites

- Go 1.23 or later
- A running Keycloak instance for testing
- Make

### Setup

```bash
git clone https://github.com/MnemoShare/mnemoshare-keycloak-mcp.git
cd mnemoshare-keycloak-mcp
make build
```

### Running locally

Copy the example environment file and configure it for your Keycloak instance:

```bash
cp .env.example .env
# Edit .env with your Keycloak URL and credentials
```

Run in stdio mode (for use with Claude Code or other MCP clients):

```bash
make run
```

Run in HTTP mode (for testing with curl or other HTTP clients):

```bash
make run-http
```

## How to Contribute

### Reporting Bugs

- Use the [Bug Report](https://github.com/MnemoShare/mnemoshare-keycloak-mcp/issues/new?template=bug_report.yml) template
- Include your keycloak-mcp version and Keycloak version
- Set `LOG_LEVEL=debug` and include relevant log output
- Describe steps to reproduce the issue

### Requesting Features

- Use the [Feature Request](https://github.com/MnemoShare/mnemoshare-keycloak-mcp/issues/new?template=feature_request.yml) template
- Describe the use case and why it would be valuable
- If requesting new Keycloak API coverage, link to the relevant [Keycloak REST API docs](https://www.keycloak.org/docs-api/latest/rest-api/)

### Submitting Code

1. **Fork** the repository
2. **Create a branch** from `main`:
   ```bash
   git checkout -b feature/my-feature
   ```
3. **Make your changes** — keep commits focused and well-described
4. **Test your changes** against a running Keycloak instance
5. **Push** your branch and open a Pull Request

### Pull Request Guidelines

- Keep PRs focused — one feature or fix per PR
- Update the README if you add new tools or change configuration
- Follow existing code patterns and naming conventions
- Include a clear description of what the PR does and why

## Project Structure

```
.
├── cmd/server/          # Entry point
├── internal/
│   ├── config/          # Environment variable configuration
│   ├── keycloak/        # Keycloak client wrapper and token management
│   ├── server/          # MCP server setup and transport handling
│   └── tools/           # Tool implementations (one file per domain)
│       ├── realms.go
│       ├── users.go
│       ├── groups.go
│       ├── clients.go
│       ├── roles.go
│       ├── idp.go
│       ├── auth_flows.go
│       ├── client_scopes.go
│       ├── sessions.go
│       ├── authorization.go
│       ├── components.go
│       ├── attack_detection.go
│       └── server_info.go
├── k8s/                 # Kubernetes deployment manifests
├── Makefile
├── Dockerfile
└── README.md
```

### Adding a New Tool

1. Find the appropriate file in `internal/tools/` (or create a new one for a new domain)
2. Define the tool with its name, description, and input schema
3. Implement the handler function using the gocloak client
4. Register the tool in the server setup

### Naming Conventions

- Tool names: `verb_noun` (e.g. `create_user`, `list_realm_roles`, `get_client_secret`)
- Keep names consistent with existing tools in the same domain
- Use the Keycloak API terminology where possible

## Code of Conduct

Be respectful and constructive. We're all here to make Keycloak administration better.

## Questions?

Open a [discussion](https://github.com/MnemoShare/mnemoshare-keycloak-mcp/issues) or reach out to the maintainers.

## License

By contributing, you agree that your contributions will be licensed under the [MIT License](LICENSE).
