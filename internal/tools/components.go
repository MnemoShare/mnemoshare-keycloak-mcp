package tools

import (
	"context"
	"fmt"

	"github.com/Nerzal/gocloak/v13"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/mnemoshare/mnemoshare-keycloak-mcp/internal/keycloak"
)

// ---------------------------------------------------------------------------
// Arg structs
// ---------------------------------------------------------------------------

type listComponentsArgs struct {
	Realm string `json:"realm,omitempty" jsonschema:"description=Realm name (uses default if omitted)"`
	Name  string `json:"name,omitempty"  jsonschema:"description=Filter by component name"`
	Type  string `json:"type,omitempty"  jsonschema:"description=Filter by provider type"`
}

type getComponentArgs struct {
	Realm       string `json:"realm,omitempty" jsonschema:"description=Realm name (uses default if omitted)"`
	ComponentID string `json:"component_id"    jsonschema:"description=Component ID,required"`
}

type createComponentArgs struct {
	Realm        string              `json:"realm,omitempty"      jsonschema:"description=Realm name (uses default if omitted)"`
	Name         string              `json:"name"                 jsonschema:"description=Component name,required"`
	ProviderType string              `json:"provider_type"        jsonschema:"description=Provider type,required"`
	ProviderID   string              `json:"provider_id"          jsonschema:"description=Provider ID,required"`
	ParentID     string              `json:"parent_id,omitempty"  jsonschema:"description=Parent component ID"`
	Config       map[string][]string `json:"config,omitempty"     jsonschema:"description=Component configuration"`
}

type updateComponentArgs struct {
	Realm       string              `json:"realm,omitempty"      jsonschema:"description=Realm name (uses default if omitted)"`
	ComponentID string              `json:"component_id"         jsonschema:"description=Component ID,required"`
	Name        *string             `json:"name,omitempty"       jsonschema:"description=New component name"`
	Config      map[string][]string `json:"config,omitempty"     jsonschema:"description=Updated component configuration"`
}

type deleteComponentArgs struct {
	Realm       string `json:"realm,omitempty" jsonschema:"description=Realm name (uses default if omitted)"`
	ComponentID string `json:"component_id"    jsonschema:"description=Component ID,required"`
}

// ---------------------------------------------------------------------------
// Registration
// ---------------------------------------------------------------------------

func registerComponentTools(s *mcp.Server, kc *keycloak.Client) {

	// 1. list_components
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_components",
		Description: "List components (user storage, LDAP, etc.) in a realm",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args listComponentsArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		params := gocloak.GetComponentsParams{}
		if args.Name != "" {
			params.Name = gocloak.StringP(args.Name)
		}
		if args.Type != "" {
			params.ProviderType = gocloak.StringP(args.Type)
		}

		components, err := kc.GC.GetComponentsWithParams(ctx, token, realm, params)
		if err != nil {
			return toolError(fmt.Sprintf("failed to list components: %v", err))
		}
		return toolResult(components)
	})

	// 2. get_component
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_component",
		Description: "Get a component by ID",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args getComponentArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		component, err := kc.GC.GetComponent(ctx, token, realm, args.ComponentID)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get component %q: %v", args.ComponentID, err))
		}
		return toolResult(component)
	})

	// 3. create_component
	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_component",
		Description: "Create a component (e.g. user federation provider) in a realm",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args createComponentArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		component := gocloak.Component{
			Name:         gocloak.StringP(args.Name),
			ProviderType: gocloak.StringP(args.ProviderType),
			ProviderID:   gocloak.StringP(args.ProviderID),
		}
		if args.ParentID != "" {
			component.ParentID = gocloak.StringP(args.ParentID)
		}
		if len(args.Config) > 0 {
			component.ComponentConfig = &args.Config
		}

		id, err := kc.GC.CreateComponent(ctx, token, realm, component)
		if err != nil {
			return toolError(fmt.Sprintf("failed to create component: %v", err))
		}
		return toolSuccess(fmt.Sprintf("Component created with ID: %s", id))
	})

	// 4. update_component
	mcp.AddTool(s, &mcp.Tool{
		Name:        "update_component",
		Description: "Update a component in a realm",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args updateComponentArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		existing, err := kc.GC.GetComponent(ctx, token, realm, args.ComponentID)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get component %q for update: %v", args.ComponentID, err))
		}

		if args.Name != nil {
			existing.Name = args.Name
		}
		if len(args.Config) > 0 {
			existing.ComponentConfig = &args.Config
		}

		if err := kc.GC.UpdateComponent(ctx, token, realm, *existing); err != nil {
			return toolError(fmt.Sprintf("failed to update component %q: %v", args.ComponentID, err))
		}
		return toolSuccess(fmt.Sprintf("Component %q updated successfully", args.ComponentID))
	})

	// 5. delete_component
	mcp.AddTool(s, &mcp.Tool{
		Name:        "delete_component",
		Description: "Delete a component from a realm",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args deleteComponentArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		if err := kc.GC.DeleteComponent(ctx, token, realm, args.ComponentID); err != nil {
			return toolError(fmt.Sprintf("failed to delete component %q: %v", args.ComponentID, err))
		}
		return toolSuccess(fmt.Sprintf("Component %q deleted successfully", args.ComponentID))
	})
}
