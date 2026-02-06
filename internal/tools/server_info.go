package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/mnemoshare/mnemoshare-keycloak-mcp/internal/keycloak"
)

type getServerInfoArgs struct{}

func registerServerInfoTools(s *mcp.Server, kc *keycloak.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_server_info",
		Description: "Get Keycloak server info including system, memory, providers, and themes",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args getServerInfoArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get token: %v", err))
		}

		info, err := kc.GC.GetServerInfo(ctx, token)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get server info: %v", err))
		}
		return toolResult(info)
	})
}
