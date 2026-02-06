package tools

import (
	"context"
	"fmt"

	"github.com/Nerzal/gocloak/v13"
	"github.com/mnemoshare/mnemoshare-keycloak-mcp/internal/keycloak"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ---------------------------------------------------------------------------
// Arg types
// ---------------------------------------------------------------------------

type listRealmRolesArgs struct {
	Realm string `json:"realm,omitempty" jsonschema:"Realm name (uses default if omitted)"`
	First *int   `json:"first,omitempty" jsonschema:"Pagination offset"`
	Max   *int   `json:"max,omitempty"   jsonschema:"Maximum number of results"`
}

type getRealmRoleArgs struct {
	Realm    string `json:"realm,omitempty"  jsonschema:"Realm name (uses default if omitted)"`
	RoleName string `json:"role_name"        jsonschema:"Role name"`
}

type createRealmRoleArgs struct {
	Realm       string `json:"realm,omitempty"       jsonschema:"Realm name (uses default if omitted)"`
	Name        string `json:"name"                  jsonschema:"Role name"`
	Description string `json:"description,omitempty" jsonschema:"Role description"`
}

type updateRealmRoleArgs struct {
	Realm       string  `json:"realm,omitempty"       jsonschema:"Realm name (uses default if omitted)"`
	RoleName    string  `json:"role_name"             jsonschema:"Current role name"`
	Name        *string `json:"name,omitempty"        jsonschema:"New role name"`
	Description *string `json:"description,omitempty" jsonschema:"New role description"`
}

type deleteRealmRoleArgs struct {
	Realm    string `json:"realm,omitempty" jsonschema:"Realm name (uses default if omitted)"`
	RoleName string `json:"role_name"       jsonschema:"Role name"`
}

type getRealmRoleCompositesArgs struct {
	Realm    string `json:"realm,omitempty" jsonschema:"Realm name (uses default if omitted)"`
	RoleName string `json:"role_name"       jsonschema:"Role name"`
}

type addRealmRoleCompositesArgs struct {
	Realm    string   `json:"realm,omitempty" jsonschema:"Realm name (uses default if omitted)"`
	RoleName string   `json:"role_name"       jsonschema:"Role name"`
	Roles    []string `json:"roles"           jsonschema:"List of role names to add as composites"`
}

type removeRealmRoleCompositesArgs struct {
	Realm    string   `json:"realm,omitempty" jsonschema:"Realm name (uses default if omitted)"`
	RoleName string   `json:"role_name"       jsonschema:"Role name"`
	Roles    []string `json:"roles"           jsonschema:"List of role names to remove from composites"`
}

type listClientRolesArgs struct {
	Realm    string `json:"realm,omitempty" jsonschema:"Realm name (uses default if omitted)"`
	ClientID string `json:"client_id"       jsonschema:"Internal client UUID"`
}

type getClientRoleArgs struct {
	Realm    string `json:"realm,omitempty" jsonschema:"Realm name (uses default if omitted)"`
	ClientID string `json:"client_id"       jsonschema:"Internal client UUID"`
	RoleName string `json:"role_name"       jsonschema:"Role name"`
}

type createClientRoleArgs struct {
	Realm       string `json:"realm,omitempty"       jsonschema:"Realm name (uses default if omitted)"`
	ClientID    string `json:"client_id"             jsonschema:"Internal client UUID"`
	Name        string `json:"name"                  jsonschema:"Role name"`
	Description string `json:"description,omitempty" jsonschema:"Role description"`
}

type updateClientRoleArgs struct {
	Realm       string  `json:"realm,omitempty"       jsonschema:"Realm name (uses default if omitted)"`
	ClientID    string  `json:"client_id"             jsonschema:"Internal client UUID"`
	RoleName    string  `json:"role_name"             jsonschema:"Current role name"`
	Name        *string `json:"name,omitempty"        jsonschema:"New role name"`
	Description *string `json:"description,omitempty" jsonschema:"New role description"`
}

type deleteClientRoleArgs struct {
	Realm    string `json:"realm,omitempty" jsonschema:"Realm name (uses default if omitted)"`
	ClientID string `json:"client_id"       jsonschema:"Internal client UUID"`
	RoleName string `json:"role_name"       jsonschema:"Role name"`
}

type getUsersByRealmRoleArgs struct {
	Realm    string `json:"realm,omitempty" jsonschema:"Realm name (uses default if omitted)"`
	RoleName string `json:"role_name"       jsonschema:"Role name"`
}

type getUsersByClientRoleArgs struct {
	Realm    string `json:"realm,omitempty" jsonschema:"Realm name (uses default if omitted)"`
	ClientID string `json:"client_id"       jsonschema:"Internal client UUID"`
	RoleName string `json:"role_name"       jsonschema:"Role name"`
}

type getGroupsByRealmRoleArgs struct {
	Realm    string `json:"realm,omitempty" jsonschema:"Realm name (uses default if omitted)"`
	RoleName string `json:"role_name"       jsonschema:"Role name"`
}

// ---------------------------------------------------------------------------
// Registration
// ---------------------------------------------------------------------------

func registerRoleTools(s *mcp.Server, kc *keycloak.Client) {
	// 1. list_realm_roles
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_realm_roles",
		Description: "List all realm-level roles with optional pagination",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args listRealmRolesArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("auth failed: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)
		params := gocloak.GetRoleParams{
			First: args.First,
			Max:   args.Max,
		}
		roles, err := kc.GC.GetRealmRoles(ctx, token, realm, params)
		if err != nil {
			return toolError(fmt.Sprintf("failed to list realm roles: %v", err))
		}
		return toolResult(roles)
	})

	// 2. get_realm_role
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_realm_role",
		Description: "Get a realm role by name",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args getRealmRoleArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("auth failed: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)
		role, err := kc.GC.GetRealmRole(ctx, token, realm, args.RoleName)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get realm role %q: %v", args.RoleName, err))
		}
		return toolResult(role)
	})

	// 3. create_realm_role
	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_realm_role",
		Description: "Create a new realm-level role",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args createRealmRoleArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("auth failed: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)
		role := gocloak.Role{
			Name:        gocloak.StringP(args.Name),
			Description: gocloak.StringP(args.Description),
		}
		id, err := kc.GC.CreateRealmRole(ctx, token, realm, role)
		if err != nil {
			return toolError(fmt.Sprintf("failed to create realm role %q: %v", args.Name, err))
		}
		return toolSuccess(fmt.Sprintf("Realm role %q created (id=%s)", args.Name, id))
	})

	// 4. update_realm_role
	mcp.AddTool(s, &mcp.Tool{
		Name:        "update_realm_role",
		Description: "Update an existing realm role (name and/or description)",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args updateRealmRoleArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("auth failed: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		existing, err := kc.GC.GetRealmRole(ctx, token, realm, args.RoleName)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get realm role %q: %v", args.RoleName, err))
		}

		if args.Name != nil {
			existing.Name = args.Name
		}
		if args.Description != nil {
			existing.Description = args.Description
		}

		if err := kc.GC.UpdateRealmRole(ctx, token, realm, args.RoleName, *existing); err != nil {
			return toolError(fmt.Sprintf("failed to update realm role %q: %v", args.RoleName, err))
		}
		return toolSuccess(fmt.Sprintf("Realm role %q updated", args.RoleName))
	})

	// 5. delete_realm_role
	mcp.AddTool(s, &mcp.Tool{
		Name:        "delete_realm_role",
		Description: "Delete a realm role by name",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args deleteRealmRoleArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("auth failed: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)
		if err := kc.GC.DeleteRealmRole(ctx, token, realm, args.RoleName); err != nil {
			return toolError(fmt.Sprintf("failed to delete realm role %q: %v", args.RoleName, err))
		}
		return toolSuccess(fmt.Sprintf("Realm role %q deleted", args.RoleName))
	})

	// 6. get_realm_role_composites
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_realm_role_composites",
		Description: "Get composite roles for a realm role",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args getRealmRoleCompositesArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("auth failed: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)
		composites, err := kc.GC.GetCompositeRealmRoles(ctx, token, realm, args.RoleName)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get composites for role %q: %v", args.RoleName, err))
		}
		return toolResult(composites)
	})

	// 7. add_realm_role_composites
	mcp.AddTool(s, &mcp.Tool{
		Name:        "add_realm_role_composites",
		Description: "Add composite roles to a realm role",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args addRealmRoleCompositesArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("auth failed: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		roles := make([]gocloak.Role, 0, len(args.Roles))
		for _, rn := range args.Roles {
			r, err := kc.GC.GetRealmRole(ctx, token, realm, rn)
			if err != nil {
				return toolError(fmt.Sprintf("failed to resolve role %q: %v", rn, err))
			}
			roles = append(roles, *r)
		}

		if err := kc.GC.AddRealmRoleComposite(ctx, token, realm, args.RoleName, roles); err != nil {
			return toolError(fmt.Sprintf("failed to add composites to role %q: %v", args.RoleName, err))
		}
		return toolSuccess(fmt.Sprintf("Added %d composite role(s) to %q", len(roles), args.RoleName))
	})

	// 8. remove_realm_role_composites
	mcp.AddTool(s, &mcp.Tool{
		Name:        "remove_realm_role_composites",
		Description: "Remove composite roles from a realm role",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args removeRealmRoleCompositesArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("auth failed: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		roles := make([]gocloak.Role, 0, len(args.Roles))
		for _, rn := range args.Roles {
			r, err := kc.GC.GetRealmRole(ctx, token, realm, rn)
			if err != nil {
				return toolError(fmt.Sprintf("failed to resolve role %q: %v", rn, err))
			}
			roles = append(roles, *r)
		}

		if err := kc.GC.DeleteRealmRoleComposite(ctx, token, realm, args.RoleName, roles); err != nil {
			return toolError(fmt.Sprintf("failed to remove composites from role %q: %v", args.RoleName, err))
		}
		return toolSuccess(fmt.Sprintf("Removed %d composite role(s) from %q", len(roles), args.RoleName))
	})

	// 9. list_client_roles
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_client_roles",
		Description: "List all roles for a specific client",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args listClientRolesArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("auth failed: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)
		roles, err := kc.GC.GetClientRoles(ctx, token, realm, args.ClientID, gocloak.GetRoleParams{})
		if err != nil {
			return toolError(fmt.Sprintf("failed to list client roles: %v", err))
		}
		return toolResult(roles)
	})

	// 10. get_client_role
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_client_role",
		Description: "Get a client role by name",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args getClientRoleArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("auth failed: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)
		role, err := kc.GC.GetClientRole(ctx, token, realm, args.ClientID, args.RoleName)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get client role %q: %v", args.RoleName, err))
		}
		return toolResult(role)
	})

	// 11. create_client_role
	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_client_role",
		Description: "Create a new role for a specific client",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args createClientRoleArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("auth failed: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)
		role := gocloak.Role{
			Name:        gocloak.StringP(args.Name),
			Description: gocloak.StringP(args.Description),
		}
		id, err := kc.GC.CreateClientRole(ctx, token, realm, args.ClientID, role)
		if err != nil {
			return toolError(fmt.Sprintf("failed to create client role %q: %v", args.Name, err))
		}
		return toolSuccess(fmt.Sprintf("Client role %q created (id=%s)", args.Name, id))
	})

	// 12. update_client_role
	mcp.AddTool(s, &mcp.Tool{
		Name:        "update_client_role",
		Description: "Update an existing client role. Fetches the role first then applies changes using the role ID.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args updateClientRoleArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("auth failed: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		existing, err := kc.GC.GetClientRole(ctx, token, realm, args.ClientID, args.RoleName)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get client role %q: %v", args.RoleName, err))
		}

		if args.Name != nil {
			existing.Name = args.Name
		}
		if args.Description != nil {
			existing.Description = args.Description
		}

		if err := kc.GC.UpdateRole(ctx, token, realm, args.ClientID, *existing); err != nil {
			return toolError(fmt.Sprintf("failed to update client role %q: %v", args.RoleName, err))
		}
		return toolSuccess(fmt.Sprintf("Client role %q updated", args.RoleName))
	})

	// 13. delete_client_role
	mcp.AddTool(s, &mcp.Tool{
		Name:        "delete_client_role",
		Description: "Delete a client role by name",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args deleteClientRoleArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("auth failed: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)
		if err := kc.GC.DeleteClientRole(ctx, token, realm, args.ClientID, args.RoleName); err != nil {
			return toolError(fmt.Sprintf("failed to delete client role %q: %v", args.RoleName, err))
		}
		return toolSuccess(fmt.Sprintf("Client role %q deleted", args.RoleName))
	})

	// 14. get_users_by_realm_role
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_users_by_realm_role",
		Description: "Get all users assigned a specific realm role",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args getUsersByRealmRoleArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("auth failed: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)
		users, err := kc.GC.GetUsersByRoleName(ctx, token, realm, args.RoleName, gocloak.GetUsersByRoleParams{})
		if err != nil {
			return toolError(fmt.Sprintf("failed to get users for role %q: %v", args.RoleName, err))
		}
		return toolResult(users)
	})

	// 15. get_users_by_client_role
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_users_by_client_role",
		Description: "Get all users assigned a specific client role",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args getUsersByClientRoleArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("auth failed: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)
		users, err := kc.GC.GetUsersByClientRoleName(ctx, token, realm, args.ClientID, args.RoleName, gocloak.GetUsersByRoleParams{})
		if err != nil {
			return toolError(fmt.Sprintf("failed to get users for client role %q: %v", args.RoleName, err))
		}
		return toolResult(users)
	})

	// 16. get_groups_by_realm_role
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_groups_by_realm_role",
		Description: "Get all groups assigned a specific realm role",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args getGroupsByRealmRoleArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("auth failed: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)
		groups, err := kc.GC.GetGroupsByRole(ctx, token, realm, args.RoleName)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get groups for role %q: %v", args.RoleName, err))
		}
		return toolResult(groups)
	})
}
