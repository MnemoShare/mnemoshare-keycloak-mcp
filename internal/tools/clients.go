package tools

import (
	"context"
	"fmt"

	"github.com/Nerzal/gocloak/v13"
	"github.com/mnemoshare/mnemoshare-keycloak-mcp/internal/keycloak"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ---------------------------------------------------------------------------
// Argument structs
// ---------------------------------------------------------------------------

type listClientsArgs struct {
	Realm    string `json:"realm,omitempty"    jsonschema:"Keycloak realm (uses default if omitted)"`
	ClientID string `json:"client_id,omitempty" jsonschema:"Filter by clientId"`
}

type getClientArgs struct {
	Realm string `json:"realm,omitempty" jsonschema:"Keycloak realm (uses default if omitted)"`
	ID    string `json:"id"             jsonschema:"Internal UUID of the client"`
}

type createClientArgs struct {
	Realm        string   `json:"realm,omitempty"          jsonschema:"Keycloak realm (uses default if omitted)"`
	ClientID     string   `json:"client_id"                jsonschema:"The clientId for the new client"`
	Name         string   `json:"name,omitempty"           jsonschema:"Display name"`
	Protocol     string   `json:"protocol,omitempty"       jsonschema:"Protocol (default openid-connect)"`
	PublicClient bool     `json:"public_client,omitempty"  jsonschema:"Whether the client is public (no secret)"`
	RootURL      string   `json:"root_url,omitempty"       jsonschema:"Root URL of the client"`
	RedirectURIs []string `json:"redirect_uris,omitempty"  jsonschema:"Valid redirect URIs"`
	WebOrigins   []string `json:"web_origins,omitempty"    jsonschema:"Allowed web origins"`
}

type updateClientArgs struct {
	Realm        string    `json:"realm,omitempty"          jsonschema:"Keycloak realm (uses default if omitted)"`
	ID           string    `json:"id"                       jsonschema:"Internal UUID of the client"`
	Name         *string   `json:"name,omitempty"           jsonschema:"Display name"`
	RootURL      *string   `json:"root_url,omitempty"       jsonschema:"Root URL"`
	RedirectURIs *[]string `json:"redirect_uris,omitempty"  jsonschema:"Valid redirect URIs"`
	WebOrigins   *[]string `json:"web_origins,omitempty"    jsonschema:"Allowed web origins"`
	Enabled      *bool     `json:"enabled,omitempty"        jsonschema:"Whether the client is enabled"`
}

type deleteClientArgs struct {
	Realm string `json:"realm,omitempty" jsonschema:"Keycloak realm (uses default if omitted)"`
	ID    string `json:"id"             jsonschema:"Internal UUID of the client"`
}

type getClientSecretArgs struct {
	Realm string `json:"realm,omitempty" jsonschema:"Keycloak realm (uses default if omitted)"`
	ID    string `json:"id"             jsonschema:"Internal UUID of the client"`
}

type regenerateClientSecretArgs struct {
	Realm string `json:"realm,omitempty" jsonschema:"Keycloak realm (uses default if omitted)"`
	ID    string `json:"id"             jsonschema:"Internal UUID of the client"`
}

type getClientServiceAccountArgs struct {
	Realm string `json:"realm,omitempty" jsonschema:"Keycloak realm (uses default if omitted)"`
	ID    string `json:"id"             jsonschema:"Internal UUID of the client"`
}

type getClientDefaultScopesArgs struct {
	Realm string `json:"realm,omitempty" jsonschema:"Keycloak realm (uses default if omitted)"`
	ID    string `json:"id"             jsonschema:"Internal UUID of the client"`
}

type addClientDefaultScopeArgs struct {
	Realm   string `json:"realm,omitempty" jsonschema:"Keycloak realm (uses default if omitted)"`
	ID      string `json:"id"             jsonschema:"Internal UUID of the client"`
	ScopeID string `json:"scope_id"       jsonschema:"ID of the scope to add"`
}

type removeClientDefaultScopeArgs struct {
	Realm   string `json:"realm,omitempty" jsonschema:"Keycloak realm (uses default if omitted)"`
	ID      string `json:"id"             jsonschema:"Internal UUID of the client"`
	ScopeID string `json:"scope_id"       jsonschema:"ID of the scope to remove"`
}

type getClientOptionalScopesArgs struct {
	Realm string `json:"realm,omitempty" jsonschema:"Keycloak realm (uses default if omitted)"`
	ID    string `json:"id"             jsonschema:"Internal UUID of the client"`
}

type addClientOptionalScopeArgs struct {
	Realm   string `json:"realm,omitempty" jsonschema:"Keycloak realm (uses default if omitted)"`
	ID      string `json:"id"             jsonschema:"Internal UUID of the client"`
	ScopeID string `json:"scope_id"       jsonschema:"ID of the scope to add"`
}

type removeClientOptionalScopeArgs struct {
	Realm   string `json:"realm,omitempty" jsonschema:"Keycloak realm (uses default if omitted)"`
	ID      string `json:"id"             jsonschema:"Internal UUID of the client"`
	ScopeID string `json:"scope_id"       jsonschema:"ID of the scope to remove"`
}

type createClientProtocolMapperArgs struct {
	Realm      string            `json:"realm,omitempty"  jsonschema:"Keycloak realm (uses default if omitted)"`
	ID         string            `json:"id"               jsonschema:"Internal UUID of the client"`
	Name       string            `json:"name"             jsonschema:"Name of the protocol mapper"`
	Protocol   string            `json:"protocol"         jsonschema:"Protocol (e.g. openid-connect)"`
	MapperType string            `json:"mapper_type"      jsonschema:"Protocol mapper type (e.g. oidc-usermodel-attribute-mapper)"`
	Config     map[string]string `json:"config"           jsonschema:"Mapper configuration key-value pairs"`
}

type updateClientProtocolMapperArgs struct {
	Realm    string            `json:"realm,omitempty"  jsonschema:"Keycloak realm (uses default if omitted)"`
	ID       string            `json:"id"               jsonschema:"Internal UUID of the client"`
	MapperID string            `json:"mapper_id"        jsonschema:"ID of the protocol mapper"`
	Name     *string           `json:"name,omitempty"   jsonschema:"New name for the mapper"`
	Config   map[string]string `json:"config,omitempty" jsonschema:"Updated mapper configuration"`
}

type deleteClientProtocolMapperArgs struct {
	Realm    string `json:"realm,omitempty" jsonschema:"Keycloak realm (uses default if omitted)"`
	ID       string `json:"id"             jsonschema:"Internal UUID of the client"`
	MapperID string `json:"mapper_id"      jsonschema:"ID of the protocol mapper"`
}

type getClientSessionsArgs struct {
	Realm string `json:"realm,omitempty" jsonschema:"Keycloak realm (uses default if omitted)"`
	ID    string `json:"id"             jsonschema:"Internal UUID of the client"`
	First *int   `json:"first,omitempty" jsonschema:"Pagination offset"`
	Max   *int   `json:"max,omitempty"   jsonschema:"Maximum number of results"`
}

// ---------------------------------------------------------------------------
// Registration
// ---------------------------------------------------------------------------

func registerClientTools(s *mcp.Server, kc *keycloak.Client) {
	registerListClients(s, kc)
	registerGetClient(s, kc)
	registerCreateClient(s, kc)
	registerUpdateClient(s, kc)
	registerDeleteClient(s, kc)
	registerGetClientSecret(s, kc)
	registerRegenerateClientSecret(s, kc)
	registerGetClientServiceAccount(s, kc)
	registerGetClientDefaultScopes(s, kc)
	registerAddClientDefaultScope(s, kc)
	registerRemoveClientDefaultScope(s, kc)
	registerGetClientOptionalScopes(s, kc)
	registerAddClientOptionalScope(s, kc)
	registerRemoveClientOptionalScope(s, kc)
	registerCreateClientProtocolMapper(s, kc)
	registerUpdateClientProtocolMapper(s, kc)
	registerDeleteClientProtocolMapper(s, kc)
	registerGetClientSessions(s, kc)
}

// ---------------------------------------------------------------------------
// 1. list_clients
// ---------------------------------------------------------------------------

func registerListClients(s *mcp.Server, kc *keycloak.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_clients",
		Description: "List clients in a Keycloak realm, optionally filtered by clientId",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args listClientsArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		params := gocloak.GetClientsParams{}
		if args.ClientID != "" {
			params.ClientID = &args.ClientID
		}

		clients, err := kc.GC.GetClients(ctx, token, realm, params)
		if err != nil {
			return toolError(fmt.Sprintf("failed to list clients: %v", err))
		}
		return toolResult(clients)
	})
}

// ---------------------------------------------------------------------------
// 2. get_client
// ---------------------------------------------------------------------------

func registerGetClient(s *mcp.Server, kc *keycloak.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_client",
		Description: "Get a Keycloak client by its internal UUID",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args getClientArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		client, err := kc.GC.GetClient(ctx, token, realm, args.ID)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get client: %v", err))
		}
		return toolResult(client)
	})
}

// ---------------------------------------------------------------------------
// 3. create_client
// ---------------------------------------------------------------------------

func registerCreateClient(s *mcp.Server, kc *keycloak.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_client",
		Description: "Create a new client in a Keycloak realm",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args createClientArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		protocol := args.Protocol
		if protocol == "" {
			protocol = "openid-connect"
		}

		newClient := gocloak.Client{
			ClientID: gocloak.StringP(args.ClientID),
			Protocol: gocloak.StringP(protocol),
			PublicClient: gocloak.BoolP(args.PublicClient),
		}
		if args.Name != "" {
			newClient.Name = gocloak.StringP(args.Name)
		}
		if args.RootURL != "" {
			newClient.RootURL = gocloak.StringP(args.RootURL)
		}
		if len(args.RedirectURIs) > 0 {
			newClient.RedirectURIs = &args.RedirectURIs
		}
		if len(args.WebOrigins) > 0 {
			newClient.WebOrigins = &args.WebOrigins
		}

		id, err := kc.GC.CreateClient(ctx, token, realm, newClient)
		if err != nil {
			return toolError(fmt.Sprintf("failed to create client: %v", err))
		}
		return toolSuccess(fmt.Sprintf("client created with id: %s", id))
	})
}

// ---------------------------------------------------------------------------
// 4. update_client
// ---------------------------------------------------------------------------

func registerUpdateClient(s *mcp.Server, kc *keycloak.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "update_client",
		Description: "Update an existing Keycloak client",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args updateClientArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		client, err := kc.GC.GetClient(ctx, token, realm, args.ID)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get client for update: %v", err))
		}

		if args.Name != nil {
			client.Name = args.Name
		}
		if args.RootURL != nil {
			client.RootURL = args.RootURL
		}
		if args.RedirectURIs != nil {
			client.RedirectURIs = args.RedirectURIs
		}
		if args.WebOrigins != nil {
			client.WebOrigins = args.WebOrigins
		}
		if args.Enabled != nil {
			client.Enabled = args.Enabled
		}

		err = kc.GC.UpdateClient(ctx, token, realm, *client)
		if err != nil {
			return toolError(fmt.Sprintf("failed to update client: %v", err))
		}
		return toolSuccess("client updated successfully")
	})
}

// ---------------------------------------------------------------------------
// 5. delete_client
// ---------------------------------------------------------------------------

func registerDeleteClient(s *mcp.Server, kc *keycloak.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "delete_client",
		Description: "Delete a client from a Keycloak realm",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args deleteClientArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		err = kc.GC.DeleteClient(ctx, token, realm, args.ID)
		if err != nil {
			return toolError(fmt.Sprintf("failed to delete client: %v", err))
		}
		return toolSuccess("client deleted successfully")
	})
}

// ---------------------------------------------------------------------------
// 6. get_client_secret
// ---------------------------------------------------------------------------

func registerGetClientSecret(s *mcp.Server, kc *keycloak.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_client_secret",
		Description: "Get the secret for a Keycloak client",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args getClientSecretArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		cred, err := kc.GC.GetClientSecret(ctx, token, realm, args.ID)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get client secret: %v", err))
		}
		return toolResult(cred.Value)
	})
}

// ---------------------------------------------------------------------------
// 7. regenerate_client_secret
// ---------------------------------------------------------------------------

func registerRegenerateClientSecret(s *mcp.Server, kc *keycloak.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "regenerate_client_secret",
		Description: "Regenerate the secret for a Keycloak client",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args regenerateClientSecretArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		cred, err := kc.GC.RegenerateClientSecret(ctx, token, realm, args.ID)
		if err != nil {
			return toolError(fmt.Sprintf("failed to regenerate client secret: %v", err))
		}
		return toolResult(cred)
	})
}

// ---------------------------------------------------------------------------
// 8. get_client_service_account
// ---------------------------------------------------------------------------

func registerGetClientServiceAccount(s *mcp.Server, kc *keycloak.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_client_service_account",
		Description: "Get the service account user associated with a Keycloak client",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args getClientServiceAccountArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		user, err := kc.GC.GetClientServiceAccount(ctx, token, realm, args.ID)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get service account: %v", err))
		}
		return toolResult(user)
	})
}

// ---------------------------------------------------------------------------
// 9. get_client_default_scopes
// ---------------------------------------------------------------------------

func registerGetClientDefaultScopes(s *mcp.Server, kc *keycloak.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_client_default_scopes",
		Description: "Get the default scopes assigned to a Keycloak client",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args getClientDefaultScopesArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		scopes, err := kc.GC.GetClientsDefaultScopes(ctx, token, realm, args.ID)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get default scopes: %v", err))
		}
		return toolResult(scopes)
	})
}

// ---------------------------------------------------------------------------
// 10. add_client_default_scope
// ---------------------------------------------------------------------------

func registerAddClientDefaultScope(s *mcp.Server, kc *keycloak.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "add_client_default_scope",
		Description: "Add a default scope to a Keycloak client",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args addClientDefaultScopeArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		err = kc.GC.AddDefaultScopeToClient(ctx, token, realm, args.ID, args.ScopeID)
		if err != nil {
			return toolError(fmt.Sprintf("failed to add default scope: %v", err))
		}
		return toolSuccess("default scope added to client")
	})
}

// ---------------------------------------------------------------------------
// 11. remove_client_default_scope
// ---------------------------------------------------------------------------

func registerRemoveClientDefaultScope(s *mcp.Server, kc *keycloak.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "remove_client_default_scope",
		Description: "Remove a default scope from a Keycloak client",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args removeClientDefaultScopeArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		err = kc.GC.RemoveDefaultScopeFromClient(ctx, token, realm, args.ID, args.ScopeID)
		if err != nil {
			return toolError(fmt.Sprintf("failed to remove default scope: %v", err))
		}
		return toolSuccess("default scope removed from client")
	})
}

// ---------------------------------------------------------------------------
// 12. get_client_optional_scopes
// ---------------------------------------------------------------------------

func registerGetClientOptionalScopes(s *mcp.Server, kc *keycloak.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_client_optional_scopes",
		Description: "Get the optional scopes assigned to a Keycloak client",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args getClientOptionalScopesArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		scopes, err := kc.GC.GetClientsOptionalScopes(ctx, token, realm, args.ID)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get optional scopes: %v", err))
		}
		return toolResult(scopes)
	})
}

// ---------------------------------------------------------------------------
// 13. add_client_optional_scope
// ---------------------------------------------------------------------------

func registerAddClientOptionalScope(s *mcp.Server, kc *keycloak.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "add_client_optional_scope",
		Description: "Add an optional scope to a Keycloak client",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args addClientOptionalScopeArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		err = kc.GC.AddOptionalScopeToClient(ctx, token, realm, args.ID, args.ScopeID)
		if err != nil {
			return toolError(fmt.Sprintf("failed to add optional scope: %v", err))
		}
		return toolSuccess("optional scope added to client")
	})
}

// ---------------------------------------------------------------------------
// 14. remove_client_optional_scope
// ---------------------------------------------------------------------------

func registerRemoveClientOptionalScope(s *mcp.Server, kc *keycloak.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "remove_client_optional_scope",
		Description: "Remove an optional scope from a Keycloak client",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args removeClientOptionalScopeArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		err = kc.GC.RemoveOptionalScopeFromClient(ctx, token, realm, args.ID, args.ScopeID)
		if err != nil {
			return toolError(fmt.Sprintf("failed to remove optional scope: %v", err))
		}
		return toolSuccess("optional scope removed from client")
	})
}

// ---------------------------------------------------------------------------
// 15. create_client_protocol_mapper
// ---------------------------------------------------------------------------

func registerCreateClientProtocolMapper(s *mcp.Server, kc *keycloak.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_client_protocol_mapper",
		Description: "Create a protocol mapper for a Keycloak client",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args createClientProtocolMapperArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		mapper := gocloak.ProtocolMapperRepresentation{
			Name:           gocloak.StringP(args.Name),
			Protocol:       gocloak.StringP(args.Protocol),
			ProtocolMapper: gocloak.StringP(args.MapperType),
			Config:         &args.Config,
		}

		id, err := kc.GC.CreateClientProtocolMapper(ctx, token, realm, args.ID, mapper)
		if err != nil {
			return toolError(fmt.Sprintf("failed to create protocol mapper: %v", err))
		}
		return toolSuccess(fmt.Sprintf("protocol mapper created with id: %s", id))
	})
}

// ---------------------------------------------------------------------------
// 16. update_client_protocol_mapper
// ---------------------------------------------------------------------------

func registerUpdateClientProtocolMapper(s *mcp.Server, kc *keycloak.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "update_client_protocol_mapper",
		Description: "Update a protocol mapper for a Keycloak client",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args updateClientProtocolMapperArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		mapper := gocloak.ProtocolMapperRepresentation{
			ID: gocloak.StringP(args.MapperID),
		}
		if args.Name != nil {
			mapper.Name = args.Name
		}
		if len(args.Config) > 0 {
			mapper.Config = &args.Config
		}

		err = kc.GC.UpdateClientProtocolMapper(ctx, token, realm, args.ID, args.MapperID, mapper)
		if err != nil {
			return toolError(fmt.Sprintf("failed to update protocol mapper: %v", err))
		}
		return toolSuccess("protocol mapper updated successfully")
	})
}

// ---------------------------------------------------------------------------
// 17. delete_client_protocol_mapper
// ---------------------------------------------------------------------------

func registerDeleteClientProtocolMapper(s *mcp.Server, kc *keycloak.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "delete_client_protocol_mapper",
		Description: "Delete a protocol mapper from a Keycloak client",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args deleteClientProtocolMapperArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		err = kc.GC.DeleteClientProtocolMapper(ctx, token, realm, args.ID, args.MapperID)
		if err != nil {
			return toolError(fmt.Sprintf("failed to delete protocol mapper: %v", err))
		}
		return toolSuccess("protocol mapper deleted successfully")
	})
}

// ---------------------------------------------------------------------------
// 18. get_client_sessions
// ---------------------------------------------------------------------------

func registerGetClientSessions(s *mcp.Server, kc *keycloak.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_client_sessions",
		Description: "Get active user sessions for a Keycloak client",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args getClientSessionsArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		params := gocloak.GetClientUserSessionsParams{
			First: args.First,
			Max:   args.Max,
		}

		sessions, err := kc.GC.GetClientUserSessions(ctx, token, realm, args.ID, params)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get client sessions: %v", err))
		}
		return toolResult(sessions)
	})
}
