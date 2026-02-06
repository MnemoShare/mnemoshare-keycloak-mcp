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

type getResourceServerArgs struct {
	Realm    string `json:"realm,omitempty" jsonschema:"Realm name (uses default if omitted)"`
	ClientID string `json:"client_id"       jsonschema:"Internal client UUID"`
}

type listResourcesArgs struct {
	Realm    string `json:"realm,omitempty"     jsonschema:"Realm name (uses default if omitted)"`
	ClientID string `json:"client_id"           jsonschema:"Internal client UUID"`
	Name     string `json:"name,omitempty"      jsonschema:"Filter by resource name"`
	URI      string `json:"uri,omitempty"       jsonschema:"Filter by resource URI"`
	First    *int   `json:"first,omitempty"     jsonschema:"Pagination offset"`
	Max      *int   `json:"max,omitempty"       jsonschema:"Maximum number of results"`
}

type getResourceArgs struct {
	Realm      string `json:"realm,omitempty" jsonschema:"Realm name (uses default if omitted)"`
	ClientID   string `json:"client_id"       jsonschema:"Internal client UUID"`
	ResourceID string `json:"resource_id"     jsonschema:"Resource ID"`
}

type createResourceArgs struct {
	Realm       string   `json:"realm,omitempty"        jsonschema:"Realm name (uses default if omitted)"`
	ClientID    string   `json:"client_id"              jsonschema:"Internal client UUID"`
	Name        string   `json:"name"                   jsonschema:"Resource name"`
	DisplayName string   `json:"display_name,omitempty" jsonschema:"Human-friendly display name"`
	URIs        []string `json:"uris,omitempty"         jsonschema:"List of URIs protected by this resource"`
	Type        string   `json:"type,omitempty"         jsonschema:"Resource type"`
	Scopes      []string `json:"scopes,omitempty"       jsonschema:"List of scope names to associate with this resource"`
}

type updateResourceArgs struct {
	Realm       string   `json:"realm,omitempty"        jsonschema:"Realm name (uses default if omitted)"`
	ClientID    string   `json:"client_id"              jsonschema:"Internal client UUID"`
	ResourceID  string   `json:"resource_id"            jsonschema:"Resource ID"`
	Name        *string  `json:"name,omitempty"         jsonschema:"New resource name"`
	DisplayName *string  `json:"display_name,omitempty" jsonschema:"New display name"`
	URIs        []string `json:"uris,omitempty"         jsonschema:"New list of URIs"`
}

type deleteResourceArgs struct {
	Realm      string `json:"realm,omitempty" jsonschema:"Realm name (uses default if omitted)"`
	ClientID   string `json:"client_id"       jsonschema:"Internal client UUID"`
	ResourceID string `json:"resource_id"     jsonschema:"Resource ID"`
}

type listAuthScopesArgs struct {
	Realm    string `json:"realm,omitempty" jsonschema:"Realm name (uses default if omitted)"`
	ClientID string `json:"client_id"       jsonschema:"Internal client UUID"`
	Name     string `json:"name,omitempty"  jsonschema:"Filter by scope name"`
	First    *int   `json:"first,omitempty" jsonschema:"Pagination offset"`
	Max      *int   `json:"max,omitempty"   jsonschema:"Maximum number of results"`
}

type createAuthScopeArgs struct {
	Realm       string `json:"realm,omitempty"        jsonschema:"Realm name (uses default if omitted)"`
	ClientID    string `json:"client_id"              jsonschema:"Internal client UUID"`
	Name        string `json:"name"                   jsonschema:"Scope name"`
	DisplayName string `json:"display_name,omitempty" jsonschema:"Human-friendly display name"`
}

type deleteAuthScopeArgs struct {
	Realm    string `json:"realm,omitempty" jsonschema:"Realm name (uses default if omitted)"`
	ClientID string `json:"client_id"       jsonschema:"Internal client UUID"`
	ScopeID  string `json:"scope_id"        jsonschema:"Scope ID"`
}

type listPoliciesArgs struct {
	Realm    string `json:"realm,omitempty" jsonschema:"Realm name (uses default if omitted)"`
	ClientID string `json:"client_id"       jsonschema:"Internal client UUID"`
	Name     string `json:"name,omitempty"  jsonschema:"Filter by policy name"`
	First    *int   `json:"first,omitempty" jsonschema:"Pagination offset"`
	Max      *int   `json:"max,omitempty"   jsonschema:"Maximum number of results"`
}

type getPolicyArgs struct {
	Realm    string `json:"realm,omitempty" jsonschema:"Realm name (uses default if omitted)"`
	ClientID string `json:"client_id"       jsonschema:"Internal client UUID"`
	PolicyID string `json:"policy_id"       jsonschema:"Policy ID"`
}

type createPolicyArgs struct {
	Realm       string            `json:"realm,omitempty"       jsonschema:"Realm name (uses default if omitted)"`
	ClientID    string            `json:"client_id"             jsonschema:"Internal client UUID"`
	Name        string            `json:"name"                  jsonschema:"Policy name"`
	Type        string            `json:"type"                  jsonschema:"Policy type (e.g. role\\, user\\, client\\, js)"`
	Logic       string            `json:"logic,omitempty"       jsonschema:"Policy logic: POSITIVE or NEGATIVE"`
	Description string            `json:"description,omitempty" jsonschema:"Policy description"`
	Config      map[string]string `json:"config,omitempty"      jsonschema:"Policy configuration map"`
}

type deletePolicyArgs struct {
	Realm    string `json:"realm,omitempty" jsonschema:"Realm name (uses default if omitted)"`
	ClientID string `json:"client_id"       jsonschema:"Internal client UUID"`
	PolicyID string `json:"policy_id"       jsonschema:"Policy ID"`
}

type listPermissionsArgs struct {
	Realm    string `json:"realm,omitempty" jsonschema:"Realm name (uses default if omitted)"`
	ClientID string `json:"client_id"       jsonschema:"Internal client UUID"`
	Name     string `json:"name,omitempty"  jsonschema:"Filter by permission name"`
	First    *int   `json:"first,omitempty" jsonschema:"Pagination offset"`
	Max      *int   `json:"max,omitempty"   jsonschema:"Maximum number of results"`
}

type createPermissionArgs struct {
	Realm            string   `json:"realm,omitempty"             jsonschema:"Realm name (uses default if omitted)"`
	ClientID         string   `json:"client_id"                   jsonschema:"Internal client UUID"`
	Name             string   `json:"name"                        jsonschema:"Permission name"`
	Type             string   `json:"type"                        jsonschema:"Permission type: resource or scope"`
	Description      string   `json:"description,omitempty"       jsonschema:"Permission description"`
	Resources        []string `json:"resources,omitempty"         jsonschema:"List of resource IDs"`
	Scopes           []string `json:"scopes,omitempty"            jsonschema:"List of scope IDs"`
	Policies         []string `json:"policies,omitempty"          jsonschema:"List of policy IDs"`
	DecisionStrategy string   `json:"decision_strategy,omitempty" jsonschema:"Decision strategy: UNANIMOUS\\, AFFIRMATIVE\\, or CONSENSUS"`
}

// ---------------------------------------------------------------------------
// Registration
// ---------------------------------------------------------------------------

func registerAuthorizationTools(s *mcp.Server, kc *keycloak.Client) {

	// 1. get_resource_server
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_resource_server",
		Description: "Get the authorization resource server settings for a client",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args getResourceServerArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		rs, err := kc.GC.GetResourceServer(ctx, token, realm, args.ClientID)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get resource server: %v", err))
		}
		return toolResult(rs)
	})

	// 2. list_resources
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_resources",
		Description: "List authorization resources for a client",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args listResourcesArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		params := gocloak.GetResourceParams{
			First: args.First,
			Max:   args.Max,
		}
		if args.Name != "" {
			params.Name = gocloak.StringP(args.Name)
		}
		if args.URI != "" {
			params.URI = gocloak.StringP(args.URI)
		}

		resources, err := kc.GC.GetResources(ctx, token, realm, args.ClientID, params)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to list resources: %v", err))
		}
		return toolResult(resources)
	})

	// 3. get_resource
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_resource",
		Description: "Get an authorization resource by ID",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args getResourceArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		resource, err := kc.GC.GetResource(ctx, token, realm, args.ClientID, args.ResourceID)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get resource %q: %v", args.ResourceID, err))
		}
		return toolResult(resource)
	})

	// 4. create_resource
	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_resource",
		Description: "Create an authorization resource for a client",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args createResourceArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		resource := gocloak.ResourceRepresentation{
			Name: gocloak.StringP(args.Name),
		}
		if args.DisplayName != "" {
			resource.DisplayName = gocloak.StringP(args.DisplayName)
		}
		if len(args.URIs) > 0 {
			resource.URIs = &args.URIs
		}
		if args.Type != "" {
			resource.Type = gocloak.StringP(args.Type)
		}
		if len(args.Scopes) > 0 {
			scopes := make([]gocloak.ScopeRepresentation, len(args.Scopes))
			for i, name := range args.Scopes {
				scopes[i] = gocloak.ScopeRepresentation{
					Name: gocloak.StringP(name),
				}
			}
			resource.Scopes = &scopes
		}

		created, err := kc.GC.CreateResource(ctx, token, realm, args.ClientID, resource)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to create resource: %v", err))
		}
		return toolResult(created)
	})

	// 5. update_resource
	mcp.AddTool(s, &mcp.Tool{
		Name:        "update_resource",
		Description: "Update an authorization resource",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args updateResourceArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		existing, err := kc.GC.GetResource(ctx, token, realm, args.ClientID, args.ResourceID)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get resource %q for update: %v", args.ResourceID, err))
		}

		if args.Name != nil {
			existing.Name = args.Name
		}
		if args.DisplayName != nil {
			existing.DisplayName = args.DisplayName
		}
		if len(args.URIs) > 0 {
			existing.URIs = &args.URIs
		}

		if err := kc.GC.UpdateResource(ctx, token, realm, args.ClientID, *existing); err != nil {
			return toolError(fmt.Sprintf("Error: failed to update resource %q: %v", args.ResourceID, err))
		}
		return toolSuccess(fmt.Sprintf("Resource %q updated successfully", args.ResourceID))
	})

	// 6. delete_resource
	mcp.AddTool(s, &mcp.Tool{
		Name:        "delete_resource",
		Description: "Delete an authorization resource",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args deleteResourceArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		if err := kc.GC.DeleteResource(ctx, token, realm, args.ClientID, args.ResourceID); err != nil {
			return toolError(fmt.Sprintf("Error: failed to delete resource %q: %v", args.ResourceID, err))
		}
		return toolSuccess(fmt.Sprintf("Resource %q deleted successfully", args.ResourceID))
	})

	// 7. list_auth_scopes
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_auth_scopes",
		Description: "List authorization scopes for a client",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args listAuthScopesArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		params := gocloak.GetScopeParams{
			First: args.First,
			Max:   args.Max,
		}
		if args.Name != "" {
			params.Name = gocloak.StringP(args.Name)
		}

		scopes, err := kc.GC.GetScopes(ctx, token, realm, args.ClientID, params)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to list authorization scopes: %v", err))
		}
		return toolResult(scopes)
	})

	// 8. create_auth_scope
	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_auth_scope",
		Description: "Create an authorization scope for a client",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args createAuthScopeArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		scope := gocloak.ScopeRepresentation{
			Name: gocloak.StringP(args.Name),
		}
		if args.DisplayName != "" {
			scope.DisplayName = gocloak.StringP(args.DisplayName)
		}

		created, err := kc.GC.CreateScope(ctx, token, realm, args.ClientID, scope)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to create authorization scope: %v", err))
		}
		return toolResult(created)
	})

	// 9. delete_auth_scope
	mcp.AddTool(s, &mcp.Tool{
		Name:        "delete_auth_scope",
		Description: "Delete an authorization scope",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args deleteAuthScopeArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		if err := kc.GC.DeleteScope(ctx, token, realm, args.ClientID, args.ScopeID); err != nil {
			return toolError(fmt.Sprintf("Error: failed to delete scope %q: %v", args.ScopeID, err))
		}
		return toolSuccess(fmt.Sprintf("Authorization scope %q deleted successfully", args.ScopeID))
	})

	// 10. list_policies
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_policies",
		Description: "List authorization policies for a client",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args listPoliciesArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		params := gocloak.GetPolicyParams{
			First: args.First,
			Max:   args.Max,
		}
		if args.Name != "" {
			params.Name = gocloak.StringP(args.Name)
		}

		policies, err := kc.GC.GetPolicies(ctx, token, realm, args.ClientID, params)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to list policies: %v", err))
		}
		return toolResult(policies)
	})

	// 11. get_policy
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_policy",
		Description: "Get an authorization policy by ID",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args getPolicyArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		policy, err := kc.GC.GetPolicy(ctx, token, realm, args.ClientID, args.PolicyID)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get policy %q: %v", args.PolicyID, err))
		}
		return toolResult(policy)
	})

	// 12. create_policy
	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_policy",
		Description: "Create an authorization policy for a client",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args createPolicyArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		policy := gocloak.PolicyRepresentation{
			Name: gocloak.StringP(args.Name),
			Type: gocloak.StringP(args.Type),
		}
		if args.Description != "" {
			policy.Description = gocloak.StringP(args.Description)
		}
		if args.Logic != "" {
			logic := gocloak.Logic(args.Logic)
			policy.Logic = gocloak.LogicP(logic)
		}
		if len(args.Config) > 0 {
			policy.Config = &args.Config
		}

		created, err := kc.GC.CreatePolicy(ctx, token, realm, args.ClientID, policy)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to create policy: %v", err))
		}
		return toolResult(created)
	})

	// 13. delete_policy
	mcp.AddTool(s, &mcp.Tool{
		Name:        "delete_policy",
		Description: "Delete an authorization policy",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args deletePolicyArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		if err := kc.GC.DeletePolicy(ctx, token, realm, args.ClientID, args.PolicyID); err != nil {
			return toolError(fmt.Sprintf("Error: failed to delete policy %q: %v", args.PolicyID, err))
		}
		return toolSuccess(fmt.Sprintf("Policy %q deleted successfully", args.PolicyID))
	})

	// 14. list_permissions
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_permissions",
		Description: "List authorization permissions for a client",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args listPermissionsArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		params := gocloak.GetPermissionParams{
			First: args.First,
			Max:   args.Max,
		}
		if args.Name != "" {
			params.Name = gocloak.StringP(args.Name)
		}

		permissions, err := kc.GC.GetPermissions(ctx, token, realm, args.ClientID, params)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to list permissions: %v", err))
		}
		return toolResult(permissions)
	})

	// 15. create_permission
	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_permission",
		Description: "Create an authorization permission for a client",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args createPermissionArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		permission := gocloak.PermissionRepresentation{
			Name: gocloak.StringP(args.Name),
			Type: gocloak.StringP(args.Type),
		}
		if args.Description != "" {
			permission.Description = gocloak.StringP(args.Description)
		}
		if len(args.Resources) > 0 {
			permission.Resources = &args.Resources
		}
		if len(args.Scopes) > 0 {
			permission.Scopes = &args.Scopes
		}
		if len(args.Policies) > 0 {
			permission.Policies = &args.Policies
		}
		if args.DecisionStrategy != "" {
			ds := gocloak.DecisionStrategy(args.DecisionStrategy)
			permission.DecisionStrategy = gocloak.DecisionStrategyP(ds)
		}

		created, err := kc.GC.CreatePermission(ctx, token, realm, args.ClientID, permission)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to create permission: %v", err))
		}
		return toolResult(created)
	})
}
