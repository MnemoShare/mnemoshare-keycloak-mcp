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

type listGroupsArgs struct {
	Realm  string `json:"realm,omitempty"  jsonschema:"description=Realm name (uses default if omitted)"`
	Search string `json:"search,omitempty" jsonschema:"description=Search string for group name"`
	First  *int   `json:"first,omitempty"  jsonschema:"description=Pagination offset"`
	Max    *int   `json:"max,omitempty"    jsonschema:"description=Maximum number of results"`
}

type getGroupArgs struct {
	Realm   string `json:"realm,omitempty"    jsonschema:"description=Realm name (uses default if omitted)"`
	GroupID string `json:"group_id"           jsonschema:"description=Group ID,required"`
}

type createGroupArgs struct {
	Realm string `json:"realm,omitempty" jsonschema:"description=Realm name (uses default if omitted)"`
	Name  string `json:"name"           jsonschema:"description=Group name,required"`
}

type createChildGroupArgs struct {
	Realm         string `json:"realm,omitempty"    jsonschema:"description=Realm name (uses default if omitted)"`
	ParentGroupID string `json:"parent_group_id"   jsonschema:"description=Parent group ID,required"`
	Name          string `json:"name"              jsonschema:"description=Child group name,required"`
}

type updateGroupArgs struct {
	Realm   string `json:"realm,omitempty"    jsonschema:"description=Realm name (uses default if omitted)"`
	GroupID string `json:"group_id"           jsonschema:"description=Group ID,required"`
	Name    string `json:"name"              jsonschema:"description=New group name,required"`
}

type deleteGroupArgs struct {
	Realm   string `json:"realm,omitempty"    jsonschema:"description=Realm name (uses default if omitted)"`
	GroupID string `json:"group_id"           jsonschema:"description=Group ID,required"`
}

type getGroupMembersArgs struct {
	Realm   string `json:"realm,omitempty"    jsonschema:"description=Realm name (uses default if omitted)"`
	GroupID string `json:"group_id"           jsonschema:"description=Group ID,required"`
	First   *int   `json:"first,omitempty"    jsonschema:"description=Pagination offset"`
	Max     *int   `json:"max,omitempty"      jsonschema:"description=Maximum number of results"`
}

type countGroupsArgs struct {
	Realm  string `json:"realm,omitempty"  jsonschema:"description=Realm name (uses default if omitted)"`
	Search string `json:"search,omitempty" jsonschema:"description=Search string to filter groups"`
}

type getGroupRealmRolesArgs struct {
	Realm   string `json:"realm,omitempty"    jsonschema:"description=Realm name (uses default if omitted)"`
	GroupID string `json:"group_id"           jsonschema:"description=Group ID,required"`
}

type addGroupRealmRolesArgs struct {
	Realm   string   `json:"realm,omitempty"    jsonschema:"description=Realm name (uses default if omitted)"`
	GroupID string   `json:"group_id"           jsonschema:"description=Group ID,required"`
	Roles   []string `json:"roles"             jsonschema:"description=List of realm role names to add,required"`
}

type removeGroupRealmRolesArgs struct {
	Realm   string   `json:"realm,omitempty"    jsonschema:"description=Realm name (uses default if omitted)"`
	GroupID string   `json:"group_id"           jsonschema:"description=Group ID,required"`
	Roles   []string `json:"roles"             jsonschema:"description=List of realm role names to remove,required"`
}

type getGroupClientRolesArgs struct {
	Realm    string `json:"realm,omitempty"    jsonschema:"description=Realm name (uses default if omitted)"`
	GroupID  string `json:"group_id"           jsonschema:"description=Group ID,required"`
	ClientID string `json:"client_id"         jsonschema:"description=Client UUID (id of the client),required"`
}

// ---------------------------------------------------------------------------
// Registration
// ---------------------------------------------------------------------------

func registerGroupTools(s *mcp.Server, kc *keycloak.Client) {
	// 1. list_groups
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_groups",
		Description: "List groups in a Keycloak realm with optional search and pagination",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args listGroupsArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		params := gocloak.GetGroupsParams{
			First: args.First,
			Max:   args.Max,
		}
		if args.Search != "" {
			params.Search = gocloak.StringP(args.Search)
		}

		groups, err := kc.GC.GetGroups(ctx, token, realm, params)
		if err != nil {
			return toolError(fmt.Sprintf("failed to list groups: %v", err))
		}
		return toolResult(groups)
	})

	// 2. get_group
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_group",
		Description: "Get a Keycloak group by its ID",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args getGroupArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		group, err := kc.GC.GetGroup(ctx, token, realm, args.GroupID)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get group: %v", err))
		}
		return toolResult(group)
	})

	// 3. create_group
	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_group",
		Description: "Create a new top-level group in a Keycloak realm",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args createGroupArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		groupID, err := kc.GC.CreateGroup(ctx, token, realm, gocloak.Group{
			Name: gocloak.StringP(args.Name),
		})
		if err != nil {
			return toolError(fmt.Sprintf("failed to create group: %v", err))
		}
		return toolSuccess(fmt.Sprintf("Group created with ID: %s", groupID))
	})

	// 4. create_child_group
	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_child_group",
		Description: "Create a child group under an existing parent group",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args createChildGroupArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		childID, err := kc.GC.CreateChildGroup(ctx, token, realm, args.ParentGroupID, gocloak.Group{
			Name: gocloak.StringP(args.Name),
		})
		if err != nil {
			return toolError(fmt.Sprintf("failed to create child group: %v", err))
		}
		return toolSuccess(fmt.Sprintf("Child group created with ID: %s", childID))
	})

	// 5. update_group
	mcp.AddTool(s, &mcp.Tool{
		Name:        "update_group",
		Description: "Update a Keycloak group (rename)",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args updateGroupArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		group, err := kc.GC.GetGroup(ctx, token, realm, args.GroupID)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get group for update: %v", err))
		}

		group.Name = gocloak.StringP(args.Name)

		err = kc.GC.UpdateGroup(ctx, token, realm, *group)
		if err != nil {
			return toolError(fmt.Sprintf("failed to update group: %v", err))
		}
		return toolSuccess(fmt.Sprintf("Group %s updated successfully", args.GroupID))
	})

	// 6. delete_group
	mcp.AddTool(s, &mcp.Tool{
		Name:        "delete_group",
		Description: "Delete a Keycloak group by its ID",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args deleteGroupArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		err = kc.GC.DeleteGroup(ctx, token, realm, args.GroupID)
		if err != nil {
			return toolError(fmt.Sprintf("failed to delete group: %v", err))
		}
		return toolSuccess(fmt.Sprintf("Group %s deleted successfully", args.GroupID))
	})

	// 7. get_group_members
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_group_members",
		Description: "Get the members of a Keycloak group",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args getGroupMembersArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		params := gocloak.GetGroupsParams{
			First: args.First,
			Max:   args.Max,
		}

		members, err := kc.GC.GetGroupMembers(ctx, token, realm, args.GroupID, params)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get group members: %v", err))
		}
		return toolResult(members)
	})

	// 8. count_groups
	mcp.AddTool(s, &mcp.Tool{
		Name:        "count_groups",
		Description: "Count the number of groups in a Keycloak realm",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args countGroupsArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		params := gocloak.GetGroupsParams{}
		if args.Search != "" {
			params.Search = gocloak.StringP(args.Search)
		}

		count, err := kc.GC.GetGroupsCount(ctx, token, realm, params)
		if err != nil {
			return toolError(fmt.Sprintf("failed to count groups: %v", err))
		}
		return toolResult(map[string]int{"count": count})
	})

	// 9. get_group_realm_roles
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_group_realm_roles",
		Description: "Get realm roles assigned to a Keycloak group",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args getGroupRealmRolesArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		roles, err := kc.GC.GetRealmRolesByGroupID(ctx, token, realm, args.GroupID)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get group realm roles: %v", err))
		}
		return toolResult(roles)
	})

	// 10. add_group_realm_roles
	mcp.AddTool(s, &mcp.Tool{
		Name:        "add_group_realm_roles",
		Description: "Add realm roles to a Keycloak group",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args addGroupRealmRolesArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		var rolesToAdd []gocloak.Role
		for _, roleName := range args.Roles {
			role, err := kc.GC.GetRealmRole(ctx, token, realm, roleName)
			if err != nil {
				return toolError(fmt.Sprintf("failed to get realm role %q: %v", roleName, err))
			}
			rolesToAdd = append(rolesToAdd, *role)
		}

		err = kc.GC.AddRealmRoleToGroup(ctx, token, realm, args.GroupID, rolesToAdd)
		if err != nil {
			return toolError(fmt.Sprintf("failed to add realm roles to group: %v", err))
		}
		return toolSuccess(fmt.Sprintf("Added %d realm role(s) to group %s", len(rolesToAdd), args.GroupID))
	})

	// 11. remove_group_realm_roles
	mcp.AddTool(s, &mcp.Tool{
		Name:        "remove_group_realm_roles",
		Description: "Remove realm roles from a Keycloak group",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args removeGroupRealmRolesArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		var rolesToRemove []gocloak.Role
		for _, roleName := range args.Roles {
			role, err := kc.GC.GetRealmRole(ctx, token, realm, roleName)
			if err != nil {
				return toolError(fmt.Sprintf("failed to get realm role %q: %v", roleName, err))
			}
			rolesToRemove = append(rolesToRemove, *role)
		}

		err = kc.GC.DeleteRealmRoleFromGroup(ctx, token, realm, args.GroupID, rolesToRemove)
		if err != nil {
			return toolError(fmt.Sprintf("failed to remove realm roles from group: %v", err))
		}
		return toolSuccess(fmt.Sprintf("Removed %d realm role(s) from group %s", len(rolesToRemove), args.GroupID))
	})

	// 12. get_group_client_roles
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_group_client_roles",
		Description: "Get client roles assigned to a Keycloak group",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args getGroupClientRolesArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get token: %v", err))
		}
		realm := kc.ResolveRealm(args.Realm)

		roles, err := kc.GC.GetClientRolesByGroupID(ctx, token, realm, args.ClientID, args.GroupID)
		if err != nil {
			return toolError(fmt.Sprintf("failed to get group client roles: %v", err))
		}
		return toolResult(roles)
	})
}
