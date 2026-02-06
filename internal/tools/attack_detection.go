package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/mnemoshare/mnemoshare-keycloak-mcp/internal/keycloak"
)

// ---------------------------------------------------------------------------
// Arg structs
// ---------------------------------------------------------------------------

type getBruteForceStatusArgs struct {
	Realm  string `json:"realm,omitempty" jsonschema:"description=Realm name (uses default if omitted)"`
	UserID string `json:"user_id"         jsonschema:"description=User ID,required"`
}

type clearBruteForceStatusArgs struct {
	Realm  string `json:"realm,omitempty" jsonschema:"description=Realm name (uses default if omitted)"`
	UserID string `json:"user_id"         jsonschema:"description=User ID,required"`
}

// ---------------------------------------------------------------------------
// Registration
// ---------------------------------------------------------------------------

func registerAttackDetectionTools(s *mcp.Server, kc *keycloak.Client) {

	// 1. get_brute_force_status
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_brute_force_status",
		Description: "Get brute force detection status for a user",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args getBruteForceStatusArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		status, err := kc.GC.GetUserBruteForceDetectionStatus(ctx, token, realm, args.UserID)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get brute force status: %v", err))
		}
		return toolResult(status)
	})

	// 2. clear_brute_force_status
	mcp.AddTool(s, &mcp.Tool{
		Name:        "clear_brute_force_status",
		Description: "Clear brute force detection status for a user (re-enable login)",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args clearBruteForceStatusArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		// gocloak doesn't expose a ClearBruteForce method, so use raw DELETE.
		resp, err := kc.GC.GetRequestWithBearerAuth(ctx, token).
			Delete(fmt.Sprintf("/admin/realms/%s/attack-detection/brute-force/users/%s", realm, args.UserID))
		if err != nil {
			return toolError(fmt.Sprintf("failed to clear brute force status: %v", err))
		}
		if resp.IsError() {
			return toolError(fmt.Sprintf("failed to clear brute force status: %s", resp.Status()))
		}
		return toolSuccess("Brute force detection status cleared")
	})
}
