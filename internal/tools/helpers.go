package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Nerzal/gocloak/v13"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/mnemoshare/mnemoshare-keycloak-mcp/internal/keycloak"
)

// toolError returns an MCP error result for domain-level errors from Keycloak.
func toolError(msg string) (*mcp.CallToolResult, any, error) {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: msg}},
		IsError: true,
	}, nil, nil
}

// toolResult marshals the given data as indented JSON and returns it as a
// successful MCP tool result.
func toolResult(data any) (*mcp.CallToolResult, any, error) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return toolError(fmt.Sprintf("failed to marshal response: %v", err))
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: string(b)}},
	}, nil, nil
}

// toolSuccess returns a plain text success message as an MCP tool result.
func toolSuccess(msg string) (*mcp.CallToolResult, any, error) {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: msg}},
	}, nil, nil
}

// resolveClientRoles looks up each client role by name, preserving order.
func resolveClientRoles(ctx context.Context, kc *keycloak.Client, token, realm, clientID string, names []string) ([]gocloak.Role, error) {
	roles := make([]gocloak.Role, 0, len(names))
	for _, name := range names {
		role, err := kc.GC.GetClientRole(ctx, token, realm, clientID, name)
		if err != nil {
			return nil, fmt.Errorf("failed to get client role %q: %w", name, err)
		}
		roles = append(roles, *role)
	}
	return roles, nil
}
