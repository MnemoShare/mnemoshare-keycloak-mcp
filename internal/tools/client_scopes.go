package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Nerzal/gocloak/v13"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/mnemoshare/mnemoshare-keycloak-mcp/internal/keycloak"
)

// ---------------------------------------------------------------------------
// Arg structs
// ---------------------------------------------------------------------------

type listClientScopesArgs struct {
	Realm string `json:"realm,omitempty" jsonschema:"description=The realm name (uses default realm if omitted)"`
}

type getClientScopeArgs struct {
	Realm   string `json:"realm,omitempty" jsonschema:"description=The realm name (uses default realm if omitted)"`
	ScopeID string `json:"scope_id" jsonschema:"description=The client scope ID"`
}

type createClientScopeArgs struct {
	Realm       string `json:"realm,omitempty" jsonschema:"description=The realm name (uses default realm if omitted)"`
	Name        string `json:"name" jsonschema:"description=The name of the client scope"`
	Description string `json:"description,omitempty" jsonschema:"description=A description of the client scope"`
	Protocol    string `json:"protocol,omitempty" jsonschema:"description=The protocol (default openid-connect)"`
	Type        string `json:"type,omitempty" jsonschema:"description=The scope type: default or optional"`
}

type updateClientScopeArgs struct {
	Realm       string  `json:"realm,omitempty" jsonschema:"description=The realm name (uses default realm if omitted)"`
	ScopeID     string  `json:"scope_id" jsonschema:"description=The client scope ID to update"`
	Name        *string `json:"name,omitempty" jsonschema:"description=The new name of the client scope"`
	Description *string `json:"description,omitempty" jsonschema:"description=The new description of the client scope"`
	Protocol    *string `json:"protocol,omitempty" jsonschema:"description=The new protocol of the client scope"`
}

type deleteClientScopeArgs struct {
	Realm   string `json:"realm,omitempty" jsonschema:"description=The realm name (uses default realm if omitted)"`
	ScopeID string `json:"scope_id" jsonschema:"description=The client scope ID to delete"`
}

type listClientScopeProtocolMappersArgs struct {
	Realm   string `json:"realm,omitempty" jsonschema:"description=The realm name (uses default realm if omitted)"`
	ScopeID string `json:"scope_id" jsonschema:"description=The client scope ID"`
}

type createClientScopeProtocolMapperArgs struct {
	Realm      string            `json:"realm,omitempty" jsonschema:"description=The realm name (uses default realm if omitted)"`
	ScopeID    string            `json:"scope_id" jsonschema:"description=The client scope ID"`
	Name       string            `json:"name" jsonschema:"description=The name of the protocol mapper"`
	Protocol   string            `json:"protocol" jsonschema:"description=The protocol (e.g. openid-connect)"`
	MapperType string            `json:"mapper_type" jsonschema:"description=The protocol mapper type (e.g. oidc-usermodel-attribute-mapper)"`
	Config     map[string]string `json:"config" jsonschema:"description=The mapper configuration as key-value pairs"`
}

type updateClientScopeProtocolMapperArgs struct {
	Realm    string             `json:"realm,omitempty" jsonschema:"description=The realm name (uses default realm if omitted)"`
	ScopeID  string             `json:"scope_id" jsonschema:"description=The client scope ID"`
	MapperID string             `json:"mapper_id" jsonschema:"description=The protocol mapper ID to update"`
	Name     *string            `json:"name,omitempty" jsonschema:"description=The new name of the protocol mapper"`
	Config   map[string]string  `json:"config,omitempty" jsonschema:"description=Updated mapper configuration as key-value pairs"`
}

type deleteClientScopeProtocolMapperArgs struct {
	Realm    string `json:"realm,omitempty" jsonschema:"description=The realm name (uses default realm if omitted)"`
	ScopeID  string `json:"scope_id" jsonschema:"description=The client scope ID"`
	MapperID string `json:"mapper_id" jsonschema:"description=The protocol mapper ID to delete"`
}

type getDefaultClientScopesArgs struct {
	Realm string `json:"realm,omitempty" jsonschema:"description=The realm name (uses default realm if omitted)"`
}

// ---------------------------------------------------------------------------
// Registration
// ---------------------------------------------------------------------------

func registerClientScopeTools(s *mcp.Server, kc *keycloak.Client) {
	// 1. list_client_scopes
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_client_scopes",
		Description: "List all client scopes in a Keycloak realm",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args listClientScopesArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get token: %v", err))
		}

		realm := kc.ResolveRealm(args.Realm)
		scopes, err := kc.GC.GetClientScopes(ctx, token, realm)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to list client scopes in realm %q: %v", realm, err))
		}

		return toolResult(scopes)
	})

	// 2. get_client_scope
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_client_scope",
		Description: "Get a client scope by ID",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args getClientScopeArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get token: %v", err))
		}

		realm := kc.ResolveRealm(args.Realm)
		scope, err := kc.GC.GetClientScope(ctx, token, realm, args.ScopeID)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get client scope %q in realm %q: %v", args.ScopeID, realm, err))
		}

		return toolResult(scope)
	})

	// 3. create_client_scope
	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_client_scope",
		Description: "Create a new client scope in a Keycloak realm",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args createClientScopeArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get token: %v", err))
		}

		realm := kc.ResolveRealm(args.Realm)

		protocol := args.Protocol
		if protocol == "" {
			protocol = "openid-connect"
		}

		scopeRep := gocloak.ClientScope{
			Name:     gocloak.StringP(args.Name),
			Protocol: gocloak.StringP(protocol),
		}
		if args.Description != "" {
			scopeRep.Description = gocloak.StringP(args.Description)
		}
		if args.Type != "" {
			scopeRep.ClientScopeAttributes = &gocloak.ClientScopeAttributes{
				DisplayOnConsentScreen: gocloak.StringP("true"),
				IncludeInTokenScope:    gocloak.StringP("true"),
			}
		}

		createdID, err := kc.GC.CreateClientScope(ctx, token, realm, scopeRep)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to create client scope %q in realm %q: %v", args.Name, realm, err))
		}

		return toolSuccess(fmt.Sprintf("Client scope %q created successfully (id: %s)", args.Name, createdID))
	})

	// 4. update_client_scope
	mcp.AddTool(s, &mcp.Tool{
		Name:        "update_client_scope",
		Description: "Update an existing client scope in a Keycloak realm",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args updateClientScopeArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get token: %v", err))
		}

		realm := kc.ResolveRealm(args.Realm)

		existing, err := kc.GC.GetClientScope(ctx, token, realm, args.ScopeID)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get client scope %q for update in realm %q: %v", args.ScopeID, realm, err))
		}

		if args.Name != nil {
			existing.Name = args.Name
		}
		if args.Description != nil {
			existing.Description = args.Description
		}
		if args.Protocol != nil {
			existing.Protocol = args.Protocol
		}

		if err := kc.GC.UpdateClientScope(ctx, token, realm, *existing); err != nil {
			return toolError(fmt.Sprintf("Error: failed to update client scope %q in realm %q: %v", args.ScopeID, realm, err))
		}

		return toolSuccess(fmt.Sprintf("Client scope %q updated successfully", args.ScopeID))
	})

	// 5. delete_client_scope
	mcp.AddTool(s, &mcp.Tool{
		Name:        "delete_client_scope",
		Description: "Delete a client scope from a Keycloak realm",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args deleteClientScopeArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get token: %v", err))
		}

		realm := kc.ResolveRealm(args.Realm)
		if err := kc.GC.DeleteClientScope(ctx, token, realm, args.ScopeID); err != nil {
			return toolError(fmt.Sprintf("Error: failed to delete client scope %q in realm %q: %v", args.ScopeID, realm, err))
		}

		return toolSuccess(fmt.Sprintf("Client scope %q deleted successfully", args.ScopeID))
	})

	// 6. list_client_scope_protocol_mappers
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_client_scope_protocol_mappers",
		Description: "List all protocol mappers for a client scope",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args listClientScopeProtocolMappersArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get token: %v", err))
		}

		realm := kc.ResolveRealm(args.Realm)
		mappers, err := kc.GC.GetClientScopeProtocolMappers(ctx, token, realm, args.ScopeID)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to list protocol mappers for client scope %q in realm %q: %v", args.ScopeID, realm, err))
		}

		return toolResult(mappers)
	})

	// 7. create_client_scope_protocol_mapper
	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_client_scope_protocol_mapper",
		Description: "Create a protocol mapper in a client scope",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args createClientScopeProtocolMapperArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get token: %v", err))
		}

		realm := kc.ResolveRealm(args.Realm)

		mapper := gocloak.ProtocolMappers{
			Name:           gocloak.StringP(args.Name),
			Protocol:       gocloak.StringP(args.Protocol),
			ProtocolMapper: gocloak.StringP(args.MapperType),
		}
		if len(args.Config) > 0 {
			configJSON, _ := json.Marshal(args.Config)
			var pmConfig gocloak.ProtocolMappersConfig
			json.Unmarshal(configJSON, &pmConfig)
			mapper.ProtocolMappersConfig = &pmConfig
		}

		createdID, err := kc.GC.CreateClientScopeProtocolMapper(ctx, token, realm, args.ScopeID, mapper)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to create protocol mapper %q for client scope %q in realm %q: %v", args.Name, args.ScopeID, realm, err))
		}

		return toolSuccess(fmt.Sprintf("Protocol mapper %q created successfully (id: %s)", args.Name, createdID))
	})

	// 8. update_client_scope_protocol_mapper
	mcp.AddTool(s, &mcp.Tool{
		Name:        "update_client_scope_protocol_mapper",
		Description: "Update a protocol mapper in a client scope",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args updateClientScopeProtocolMapperArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get token: %v", err))
		}

		realm := kc.ResolveRealm(args.Realm)

		// Fetch all mappers and find the one to update.
		mappers, err := kc.GC.GetClientScopeProtocolMappers(ctx, token, realm, args.ScopeID)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get protocol mappers for client scope %q in realm %q: %v", args.ScopeID, realm, err))
		}

		var existing *gocloak.ProtocolMappers
		for _, m := range mappers {
			if m.ID != nil && *m.ID == args.MapperID {
				existing = m
				break
			}
		}
		if existing == nil {
			return toolError(fmt.Sprintf("Error: protocol mapper %q not found in client scope %q", args.MapperID, args.ScopeID))
		}

		if args.Name != nil {
			existing.Name = args.Name
		}
		if args.Config != nil {
			configJSON, _ := json.Marshal(args.Config)
			var pmConfig gocloak.ProtocolMappersConfig
			json.Unmarshal(configJSON, &pmConfig)
			existing.ProtocolMappersConfig = &pmConfig
		}

		if err := kc.GC.UpdateClientScopeProtocolMapper(ctx, token, realm, args.ScopeID, *existing); err != nil {
			return toolError(fmt.Sprintf("Error: failed to update protocol mapper %q in client scope %q in realm %q: %v", args.MapperID, args.ScopeID, realm, err))
		}

		return toolSuccess(fmt.Sprintf("Protocol mapper %q updated successfully", args.MapperID))
	})

	// 9. delete_client_scope_protocol_mapper
	mcp.AddTool(s, &mcp.Tool{
		Name:        "delete_client_scope_protocol_mapper",
		Description: "Delete a protocol mapper from a client scope",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args deleteClientScopeProtocolMapperArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get token: %v", err))
		}

		realm := kc.ResolveRealm(args.Realm)
		if err := kc.GC.DeleteClientScopeProtocolMapper(ctx, token, realm, args.ScopeID, args.MapperID); err != nil {
			return toolError(fmt.Sprintf("Error: failed to delete protocol mapper %q from client scope %q in realm %q: %v", args.MapperID, args.ScopeID, realm, err))
		}

		return toolSuccess(fmt.Sprintf("Protocol mapper %q deleted from client scope %q successfully", args.MapperID, args.ScopeID))
	})

	// 10. get_default_client_scopes
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_default_client_scopes",
		Description: "Get the realm's default client scopes",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args getDefaultClientScopesArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get token: %v", err))
		}

		realm := kc.ResolveRealm(args.Realm)
		scopes, err := kc.GC.GetDefaultDefaultClientScopes(ctx, token, realm)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get default client scopes for realm %q: %v", realm, err))
		}

		return toolResult(scopes)
	})
}
