package tools

import (
	"context"
	"fmt"

	"github.com/Nerzal/gocloak/v13"
	"github.com/mnemoshare/mnemoshare-keycloak-mcp/internal/keycloak"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ---------------------------------------------------------------------------
// 1. list_identity_providers
// ---------------------------------------------------------------------------

type listIdentityProvidersArgs struct {
	Realm string `json:"realm,omitempty" jsonschema:"target realm, uses default if omitted"`
}

func registerListIdentityProviders(s *mcp.Server, kc *keycloak.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_identity_providers",
		Description: "List all identity providers configured in a realm",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args listIdentityProvidersArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("authentication failed: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		idps, err := kc.GC.GetIdentityProviders(ctx, token, realm)
		if err != nil {
			return toolError(fmt.Sprintf("failed to list identity providers: %v", err))
		}
		return toolResult(idps)
	})
}

// ---------------------------------------------------------------------------
// 2. get_identity_provider
// ---------------------------------------------------------------------------

type getIdentityProviderArgs struct {
	Realm string `json:"realm,omitempty" jsonschema:"target realm, uses default if omitted"`
	Alias string `json:"alias" jsonschema:"alias of the identity provider to retrieve"`
}

func registerGetIdentityProvider(s *mcp.Server, kc *keycloak.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_identity_provider",
		Description: "Get a specific identity provider by alias",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args getIdentityProviderArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("authentication failed: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		idp, err := kc.GC.GetIdentityProvider(ctx, token, realm, args.Alias)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get identity provider %q: %v", args.Alias, err))
		}
		return toolResult(idp)
	})
}

// ---------------------------------------------------------------------------
// 3. create_identity_provider
// ---------------------------------------------------------------------------

type createIdentityProviderArgs struct {
	Realm       string            `json:"realm,omitempty" jsonschema:"target realm, uses default if omitted"`
	Alias       string            `json:"alias" jsonschema:"unique alias for the identity provider"`
	ProviderID  string            `json:"provider_id" jsonschema:"provider type such as google, github, saml, or oidc"`
	DisplayName string            `json:"display_name,omitempty" jsonschema:"human-readable display name"`
	Enabled     *bool             `json:"enabled,omitempty" jsonschema:"whether the identity provider is enabled, defaults to true"`
	Config      map[string]string `json:"config" jsonschema:"provider-specific configuration such as clientId and clientSecret"`
}

func registerCreateIdentityProvider(s *mcp.Server, kc *keycloak.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_identity_provider",
		Description: "Create a new identity provider in a realm",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args createIdentityProviderArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("authentication failed: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		enabled := true
		if args.Enabled != nil {
			enabled = *args.Enabled
		}

		idpRep := gocloak.IdentityProviderRepresentation{
			Alias:       &args.Alias,
			ProviderID:  &args.ProviderID,
			DisplayName: &args.DisplayName,
			Enabled:     &enabled,
			Config:      &args.Config,
		}

		id, err := kc.GC.CreateIdentityProvider(ctx, token, realm, idpRep)
		if err != nil {
			return toolError(fmt.Sprintf("failed to create identity provider %q: %v", args.Alias, err))
		}
		return toolSuccess(fmt.Sprintf("Identity provider %q created successfully (id: %s)", args.Alias, id))
	})
}

// ---------------------------------------------------------------------------
// 4. update_identity_provider
// ---------------------------------------------------------------------------

type updateIdentityProviderArgs struct {
	Realm       string            `json:"realm,omitempty" jsonschema:"target realm, uses default if omitted"`
	Alias       string            `json:"alias" jsonschema:"alias of the identity provider to update"`
	DisplayName *string           `json:"display_name,omitempty" jsonschema:"new display name"`
	Enabled     *bool             `json:"enabled,omitempty" jsonschema:"whether the identity provider should be enabled"`
	Config      map[string]string `json:"config,omitempty" jsonschema:"provider-specific configuration to update"`
}

func registerUpdateIdentityProvider(s *mcp.Server, kc *keycloak.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "update_identity_provider",
		Description: "Update an existing identity provider",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args updateIdentityProviderArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("authentication failed: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		// Fetch the existing identity provider to preserve unmodified fields.
		idp, err := kc.GC.GetIdentityProvider(ctx, token, realm, args.Alias)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get identity provider %q: %v", args.Alias, err))
		}

		if args.DisplayName != nil {
			idp.DisplayName = args.DisplayName
		}
		if args.Enabled != nil {
			idp.Enabled = args.Enabled
		}
		if args.Config != nil {
			idp.Config = &args.Config
		}

		err = kc.GC.UpdateIdentityProvider(ctx, token, realm, args.Alias, *idp)
		if err != nil {
			return toolError(fmt.Sprintf("failed to update identity provider %q: %v", args.Alias, err))
		}
		return toolSuccess(fmt.Sprintf("Identity provider %q updated successfully", args.Alias))
	})
}

// ---------------------------------------------------------------------------
// 5. delete_identity_provider
// ---------------------------------------------------------------------------

type deleteIdentityProviderArgs struct {
	Realm string `json:"realm,omitempty" jsonschema:"target realm, uses default if omitted"`
	Alias string `json:"alias" jsonschema:"alias of the identity provider to delete"`
}

func registerDeleteIdentityProvider(s *mcp.Server, kc *keycloak.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "delete_identity_provider",
		Description: "Delete an identity provider from a realm",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args deleteIdentityProviderArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("authentication failed: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		err = kc.GC.DeleteIdentityProvider(ctx, token, realm, args.Alias)
		if err != nil {
			return toolError(fmt.Sprintf("failed to delete identity provider %q: %v", args.Alias, err))
		}
		return toolSuccess(fmt.Sprintf("Identity provider %q deleted successfully", args.Alias))
	})
}

// ---------------------------------------------------------------------------
// 6. list_identity_provider_mappers
// ---------------------------------------------------------------------------

type listIdentityProviderMappersArgs struct {
	Realm string `json:"realm,omitempty" jsonschema:"target realm, uses default if omitted"`
	Alias string `json:"alias" jsonschema:"alias of the identity provider"`
}

func registerListIdentityProviderMappers(s *mcp.Server, kc *keycloak.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_identity_provider_mappers",
		Description: "List all mappers for an identity provider",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args listIdentityProviderMappersArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("authentication failed: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		mappers, err := kc.GC.GetIdentityProviderMappers(ctx, token, realm, args.Alias)
		if err != nil {
			return toolError(fmt.Sprintf("failed to list mappers for identity provider %q: %v", args.Alias, err))
		}
		return toolResult(mappers)
	})
}

// ---------------------------------------------------------------------------
// 7. create_identity_provider_mapper
// ---------------------------------------------------------------------------

type createIdentityProviderMapperArgs struct {
	Realm      string            `json:"realm,omitempty" jsonschema:"target realm, uses default if omitted"`
	Alias      string            `json:"alias" jsonschema:"alias of the identity provider"`
	Name       string            `json:"name" jsonschema:"name of the mapper"`
	MapperType string            `json:"mapper_type" jsonschema:"the identityProviderMapper type, e.g. hardcoded-user-session-attribute-idp-mapper"`
	Config     map[string]string `json:"config" jsonschema:"mapper configuration key-value pairs"`
}

func registerCreateIdentityProviderMapper(s *mcp.Server, kc *keycloak.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_identity_provider_mapper",
		Description: "Create a mapper for an identity provider",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args createIdentityProviderMapperArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("authentication failed: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		mapper := gocloak.IdentityProviderMapper{
			Name:                   &args.Name,
			IdentityProviderMapper: &args.MapperType,
			IdentityProviderAlias:  &args.Alias,
			Config:                 &args.Config,
		}

		id, err := kc.GC.CreateIdentityProviderMapper(ctx, token, realm, args.Alias, mapper)
		if err != nil {
			return toolError(fmt.Sprintf("failed to create mapper %q for identity provider %q: %v", args.Name, args.Alias, err))
		}
		return toolSuccess(fmt.Sprintf("Mapper %q created successfully for identity provider %q (id: %s)", args.Name, args.Alias, id))
	})
}

// ---------------------------------------------------------------------------
// 8. delete_identity_provider_mapper
// ---------------------------------------------------------------------------

type deleteIdentityProviderMapperArgs struct {
	Realm    string `json:"realm,omitempty" jsonschema:"target realm, uses default if omitted"`
	Alias    string `json:"alias" jsonschema:"alias of the identity provider"`
	MapperID string `json:"mapper_id" jsonschema:"ID of the mapper to delete"`
}

func registerDeleteIdentityProviderMapper(s *mcp.Server, kc *keycloak.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "delete_identity_provider_mapper",
		Description: "Delete a mapper from an identity provider",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args deleteIdentityProviderMapperArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("authentication failed: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		err = kc.GC.DeleteIdentityProviderMapper(ctx, token, realm, args.Alias, args.MapperID)
		if err != nil {
			return toolError(fmt.Sprintf("failed to delete mapper %q from identity provider %q: %v", args.MapperID, args.Alias, err))
		}
		return toolSuccess(fmt.Sprintf("Mapper %q deleted successfully from identity provider %q", args.MapperID, args.Alias))
	})
}

// ---------------------------------------------------------------------------
// Registration entry point
// ---------------------------------------------------------------------------

func registerIdentityProviderTools(s *mcp.Server, kc *keycloak.Client) {
	registerListIdentityProviders(s, kc)
	registerGetIdentityProvider(s, kc)
	registerCreateIdentityProvider(s, kc)
	registerUpdateIdentityProvider(s, kc)
	registerDeleteIdentityProvider(s, kc)
	registerListIdentityProviderMappers(s, kc)
	registerCreateIdentityProviderMapper(s, kc)
	registerDeleteIdentityProviderMapper(s, kc)
}
