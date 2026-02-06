package tools

import (
	"context"
	"fmt"

	"github.com/Nerzal/gocloak/v13"
	"github.com/mnemoshare/mnemoshare-keycloak-mcp/internal/keycloak"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ---------------------------------------------------------------------------
// Arg structs
// ---------------------------------------------------------------------------

type logoutUserAllSessionsArgs struct {
	Realm  string `json:"realm,omitempty" jsonschema:"description=Realm name (uses default if omitted)"`
	UserID string `json:"user_id"         jsonschema:"required,description=User ID to logout from all sessions"`
}

type logoutUserSessionArgs struct {
	Realm     string `json:"realm,omitempty" jsonschema:"description=Realm name (uses default if omitted)"`
	SessionID string `json:"session_id"      jsonschema:"required,description=Session ID to logout"`
}

type getEventsArgs struct {
	Realm  string `json:"realm,omitempty"  jsonschema:"description=Realm name (uses default if omitted)"`
	Type   string `json:"type,omitempty"   jsonschema:"description=Event type filter"`
	Client string `json:"client,omitempty" jsonschema:"description=Client filter"`
	User   string `json:"user,omitempty"   jsonschema:"description=User filter"`
	First  *int   `json:"first,omitempty"  jsonschema:"description=Pagination offset"`
	Max    *int   `json:"max,omitempty"    jsonschema:"description=Maximum number of results"`
}

type getClientOfflineSessionsArgs struct {
	Realm    string `json:"realm,omitempty"  jsonschema:"description=Realm name (uses default if omitted)"`
	ClientID string `json:"client_id"        jsonschema:"required,description=Internal UUID of the client"`
	First    *int   `json:"first,omitempty"  jsonschema:"description=Pagination offset"`
	Max      *int   `json:"max,omitempty"    jsonschema:"description=Maximum number of results"`
}

type revokeUserConsentsArgs struct {
	Realm    string `json:"realm,omitempty" jsonschema:"description=Realm name (uses default if omitted)"`
	UserID   string `json:"user_id"         jsonschema:"required,description=User ID"`
	ClientID string `json:"client_id"       jsonschema:"required,description=The clientId string (not UUID)"`
}

// ---------------------------------------------------------------------------
// Registration
// ---------------------------------------------------------------------------

func registerSessionTools(s *mcp.Server, kc *keycloak.Client) {

	// 1. logout_user_all_sessions
	mcp.AddTool(s, &mcp.Tool{
		Name:        "logout_user_all_sessions",
		Description: "Logout a user from all sessions",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args logoutUserAllSessionsArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		err = kc.GC.LogoutAllSessions(ctx, token, realm, args.UserID)
		if err != nil {
			return toolError(fmt.Sprintf("failed to logout user from all sessions: %v", err))
		}
		return toolSuccess("User logged out from all sessions")
	})

	// 2. logout_user_session
	mcp.AddTool(s, &mcp.Tool{
		Name:        "logout_user_session",
		Description: "Logout a specific user session",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args logoutUserSessionArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		err = kc.GC.LogoutUserSession(ctx, token, realm, args.SessionID)
		if err != nil {
			return toolError(fmt.Sprintf("failed to logout session: %v", err))
		}
		return toolSuccess("Session logged out successfully")
	})

	// 3. get_events
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_events",
		Description: "Get events for a realm",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args getEventsArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		params := gocloak.GetEventsParams{}
		if args.First != nil {
			v := int32(*args.First)
			params.First = &v
		}
		if args.Max != nil {
			v := int32(*args.Max)
			params.Max = &v
		}
		if args.Type != "" {
			params.Type = []string{args.Type}
		}
		if args.Client != "" {
			params.Client = &args.Client
		}
		if args.User != "" {
			params.UserID = &args.User
		}

		events, err := kc.GC.GetEvents(ctx, token, realm, params)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get events: %v", err))
		}
		return toolResult(events)
	})

	// 4. get_client_offline_sessions
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_client_offline_sessions",
		Description: "Get offline sessions for a client",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args getClientOfflineSessionsArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		params := gocloak.GetClientUserSessionsParams{
			First: args.First,
			Max:   args.Max,
		}

		sessions, err := kc.GC.GetClientOfflineSessions(ctx, token, realm, args.ClientID, params)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get client offline sessions: %v", err))
		}
		return toolResult(sessions)
	})

	// 5. revoke_user_consents
	mcp.AddTool(s, &mcp.Tool{
		Name:        "revoke_user_consents",
		Description: "Revoke user consents for a client",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args revokeUserConsentsArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		err = kc.GC.RevokeUserConsents(ctx, token, realm, args.UserID, args.ClientID)
		if err != nil {
			return toolError(fmt.Sprintf("failed to revoke user consents: %v", err))
		}
		return toolSuccess("User consents revoked successfully")
	})
}
