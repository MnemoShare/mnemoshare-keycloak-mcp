package tools

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/mnemoshare/mnemoshare-keycloak-mcp/internal/keycloak"
)

// RegisterAll wires every tool domain to the MCP server.
func RegisterAll(s *mcp.Server, kc *keycloak.Client) {
	registerRealmTools(s, kc)
	registerUserTools(s, kc)
	registerGroupTools(s, kc)
	registerClientTools(s, kc)
	registerRoleTools(s, kc)
	registerIdentityProviderTools(s, kc)
	registerAuthFlowTools(s, kc)
	registerClientScopeTools(s, kc)
	registerSessionTools(s, kc)
	registerAuthorizationTools(s, kc)
	registerComponentTools(s, kc)
	registerAttackDetectionTools(s, kc)
	registerServerInfoTools(s, kc)
}
