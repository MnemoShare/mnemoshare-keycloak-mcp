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

type listUsersArgs struct {
	Realm  string `json:"realm,omitempty"  jsonschema:"description=Realm name (uses default if omitted)"`
	First  *int   `json:"first,omitempty"  jsonschema:"description=Pagination offset"`
	Max    *int   `json:"max,omitempty"    jsonschema:"description=Maximum number of results"`
	Search string `json:"search,omitempty" jsonschema:"description=Search string for users"`
}

type getUserArgs struct {
	Realm  string `json:"realm,omitempty"   jsonschema:"description=Realm name (uses default if omitted)"`
	UserID string `json:"user_id"           jsonschema:"description=User ID,required"`
}

type searchUsersArgs struct {
	Realm     string `json:"realm,omitempty"      jsonschema:"description=Realm name (uses default if omitted)"`
	Username  string `json:"username,omitempty"   jsonschema:"description=Username to search for"`
	Email     string `json:"email,omitempty"      jsonschema:"description=Email to search for"`
	FirstName string `json:"first_name,omitempty" jsonschema:"description=First name to search for"`
	LastName  string `json:"last_name,omitempty"  jsonschema:"description=Last name to search for"`
	Enabled   *bool  `json:"enabled,omitempty"    jsonschema:"description=Filter by enabled status"`
	First     *int   `json:"first,omitempty"      jsonschema:"description=Pagination offset"`
	Max       *int   `json:"max,omitempty"        jsonschema:"description=Maximum number of results"`
}

type createUserArgs struct {
	Realm             string `json:"realm,omitempty"              jsonschema:"description=Realm name (uses default if omitted)"`
	Username          string `json:"username"                     jsonschema:"description=Username for the new user,required"`
	Email             string `json:"email,omitempty"              jsonschema:"description=Email address"`
	FirstName         string `json:"first_name,omitempty"        jsonschema:"description=First name"`
	LastName          string `json:"last_name,omitempty"         jsonschema:"description=Last name"`
	Enabled           *bool  `json:"enabled,omitempty"            jsonschema:"description=Whether the user is enabled (default true)"`
	Password          string `json:"password,omitempty"           jsonschema:"description=Initial password"`
	TemporaryPassword *bool  `json:"temporary_password,omitempty" jsonschema:"description=Whether the password is temporary"`
}

type updateUserArgs struct {
	Realm     string  `json:"realm,omitempty"      jsonschema:"description=Realm name (uses default if omitted)"`
	UserID    string  `json:"user_id"              jsonschema:"description=User ID,required"`
	Email     *string `json:"email,omitempty"      jsonschema:"description=New email address"`
	FirstName *string `json:"first_name,omitempty" jsonschema:"description=New first name"`
	LastName  *string `json:"last_name,omitempty"  jsonschema:"description=New last name"`
	Enabled   *bool   `json:"enabled,omitempty"    jsonschema:"description=Whether the user is enabled"`
}

type deleteUserArgs struct {
	Realm  string `json:"realm,omitempty" jsonschema:"description=Realm name (uses default if omitted)"`
	UserID string `json:"user_id"         jsonschema:"description=User ID,required"`
}

type countUsersArgs struct {
	Realm  string `json:"realm,omitempty"  jsonschema:"description=Realm name (uses default if omitted)"`
	Search string `json:"search,omitempty" jsonschema:"description=Search string to filter count"`
}

type setUserPasswordArgs struct {
	Realm     string `json:"realm,omitempty"     jsonschema:"description=Realm name (uses default if omitted)"`
	UserID    string `json:"user_id"             jsonschema:"description=User ID,required"`
	Password  string `json:"password"            jsonschema:"description=New password,required"`
	Temporary *bool  `json:"temporary,omitempty" jsonschema:"description=Whether the password is temporary"`
}

type getUserCredentialsArgs struct {
	Realm  string `json:"realm,omitempty" jsonschema:"description=Realm name (uses default if omitted)"`
	UserID string `json:"user_id"         jsonschema:"description=User ID,required"`
}

type deleteUserCredentialArgs struct {
	Realm        string `json:"realm,omitempty" jsonschema:"description=Realm name (uses default if omitted)"`
	UserID       string `json:"user_id"         jsonschema:"description=User ID,required"`
	CredentialID string `json:"credential_id"   jsonschema:"description=Credential ID to delete,required"`
}

type sendVerifyEmailArgs struct {
	Realm  string `json:"realm,omitempty" jsonschema:"description=Realm name (uses default if omitted)"`
	UserID string `json:"user_id"         jsonschema:"description=User ID,required"`
}

type executeActionsEmailArgs struct {
	Realm    string   `json:"realm,omitempty"    jsonschema:"description=Realm name (uses default if omitted)"`
	UserID   string   `json:"user_id"            jsonschema:"description=User ID,required"`
	Actions  []string `json:"actions"            jsonschema:"description=List of actions to execute,required"`
	Lifespan *int     `json:"lifespan,omitempty" jsonschema:"description=Lifespan of the action token in seconds"`
}

type getUserGroupsArgs struct {
	Realm  string `json:"realm,omitempty" jsonschema:"description=Realm name (uses default if omitted)"`
	UserID string `json:"user_id"         jsonschema:"description=User ID,required"`
}

type addUserToGroupArgs struct {
	Realm   string `json:"realm,omitempty" jsonschema:"description=Realm name (uses default if omitted)"`
	UserID  string `json:"user_id"         jsonschema:"description=User ID,required"`
	GroupID string `json:"group_id"        jsonschema:"description=Group ID,required"`
}

type removeUserFromGroupArgs struct {
	Realm   string `json:"realm,omitempty" jsonschema:"description=Realm name (uses default if omitted)"`
	UserID  string `json:"user_id"         jsonschema:"description=User ID,required"`
	GroupID string `json:"group_id"        jsonschema:"description=Group ID,required"`
}

type getUserSessionsArgs struct {
	Realm  string `json:"realm,omitempty" jsonschema:"description=Realm name (uses default if omitted)"`
	UserID string `json:"user_id"         jsonschema:"description=User ID,required"`
}

type getUserFederatedIdentitiesArgs struct {
	Realm  string `json:"realm,omitempty" jsonschema:"description=Realm name (uses default if omitted)"`
	UserID string `json:"user_id"         jsonschema:"description=User ID,required"`
}

type createUserFederatedIdentityArgs struct {
	Realm             string `json:"realm,omitempty"      jsonschema:"description=Realm name (uses default if omitted)"`
	UserID            string `json:"user_id"              jsonschema:"description=User ID,required"`
	ProviderID        string `json:"provider_id"          jsonschema:"description=Identity provider alias,required"`
	FederatedUserID   string `json:"federated_user_id"    jsonschema:"description=User ID at the identity provider,required"`
	FederatedUsername string `json:"federated_username"   jsonschema:"description=Username at the identity provider,required"`
}

type deleteUserFederatedIdentityArgs struct {
	Realm      string `json:"realm,omitempty" jsonschema:"description=Realm name (uses default if omitted)"`
	UserID     string `json:"user_id"         jsonschema:"description=User ID,required"`
	ProviderID string `json:"provider_id"     jsonschema:"description=Identity provider alias,required"`
}

type getUserRealmRolesArgs struct {
	Realm  string `json:"realm,omitempty" jsonschema:"description=Realm name (uses default if omitted)"`
	UserID string `json:"user_id"         jsonschema:"description=User ID,required"`
}

type addUserRealmRolesArgs struct {
	Realm  string   `json:"realm,omitempty" jsonschema:"description=Realm name (uses default if omitted)"`
	UserID string   `json:"user_id"         jsonschema:"description=User ID,required"`
	Roles  []string `json:"roles"           jsonschema:"description=List of realm role names to add,required"`
}

type removeUserRealmRolesArgs struct {
	Realm  string   `json:"realm,omitempty" jsonschema:"description=Realm name (uses default if omitted)"`
	UserID string   `json:"user_id"         jsonschema:"description=User ID,required"`
	Roles  []string `json:"roles"           jsonschema:"description=List of realm role names to remove,required"`
}

type getUserClientRolesArgs struct {
	Realm    string `json:"realm,omitempty" jsonschema:"description=Realm name (uses default if omitted)"`
	UserID   string `json:"user_id"         jsonschema:"description=User ID,required"`
	ClientID string `json:"client_id"       jsonschema:"description=Internal client UUID,required"`
}

type addUserClientRolesArgs struct {
	Realm    string   `json:"realm,omitempty" jsonschema:"description=Realm name (uses default if omitted)"`
	UserID   string   `json:"user_id"         jsonschema:"description=User ID,required"`
	ClientID string   `json:"client_id"       jsonschema:"description=Internal client UUID,required"`
	Roles    []string `json:"roles"           jsonschema:"description=List of client role names to add,required"`
}

// ---------------------------------------------------------------------------
// Registration
// ---------------------------------------------------------------------------

func registerUserTools(s *mcp.Server, kc *keycloak.Client) {

	// 1. list_users
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_users",
		Description: "List users in a realm",
	},
		func(ctx context.Context, req *mcp.CallToolRequest, args listUsersArgs) (*mcp.CallToolResult, any, error) {
			token, err := kc.Token(ctx)
			if err != nil {
				return toolError(fmt.Sprintf("token error: %v", err))
			}
			realm := kc.ResolveRealm(args.Realm)

			params := gocloak.GetUsersParams{
				First:  args.First,
				Max:    args.Max,
			}
			if args.Search != "" {
				params.Search = gocloak.StringP(args.Search)
			}

			users, err := kc.GC.GetUsers(ctx, token, realm, params)
			if err != nil {
				return toolError(fmt.Sprintf("failed to list users: %v", err))
			}
			return toolResult(users)
		},
	)

	// 2. get_user
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_user",
		Description: "Get a user by ID",
	},
		func(ctx context.Context, req *mcp.CallToolRequest, args getUserArgs) (*mcp.CallToolResult, any, error) {
			token, err := kc.Token(ctx)
			if err != nil {
				return toolError(fmt.Sprintf("token error: %v", err))
			}
			realm := kc.ResolveRealm(args.Realm)

			user, err := kc.GC.GetUserByID(ctx, token, realm, args.UserID)
			if err != nil {
				return toolError(fmt.Sprintf("failed to get user: %v", err))
			}
			return toolResult(user)
		},
	)

	// 3. search_users
	mcp.AddTool(s, &mcp.Tool{
		Name:        "search_users",
		Description: "Search users with detailed parameters",
	},
		func(ctx context.Context, req *mcp.CallToolRequest, args searchUsersArgs) (*mcp.CallToolResult, any, error) {
			token, err := kc.Token(ctx)
			if err != nil {
				return toolError(fmt.Sprintf("token error: %v", err))
			}
			realm := kc.ResolveRealm(args.Realm)

			params := gocloak.GetUsersParams{
				First:   args.First,
				Max:     args.Max,
				Enabled: args.Enabled,
			}
			if args.Username != "" {
				params.Username = gocloak.StringP(args.Username)
			}
			if args.Email != "" {
				params.Email = gocloak.StringP(args.Email)
			}
			if args.FirstName != "" {
				params.FirstName = gocloak.StringP(args.FirstName)
			}
			if args.LastName != "" {
				params.LastName = gocloak.StringP(args.LastName)
			}

			users, err := kc.GC.GetUsers(ctx, token, realm, params)
			if err != nil {
				return toolError(fmt.Sprintf("failed to search users: %v", err))
			}
			return toolResult(users)
		},
	)

	// 4. create_user
	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_user",
		Description: "Create a new user in a realm",
	},
		func(ctx context.Context, req *mcp.CallToolRequest, args createUserArgs) (*mcp.CallToolResult, any, error) {
			token, err := kc.Token(ctx)
			if err != nil {
				return toolError(fmt.Sprintf("token error: %v", err))
			}
			realm := kc.ResolveRealm(args.Realm)

			enabled := true
			if args.Enabled != nil {
				enabled = *args.Enabled
			}

			user := gocloak.User{
				Username:  gocloak.StringP(args.Username),
				Enabled:   gocloak.BoolP(enabled),
			}
			if args.Email != "" {
				user.Email = gocloak.StringP(args.Email)
			}
			if args.FirstName != "" {
				user.FirstName = gocloak.StringP(args.FirstName)
			}
			if args.LastName != "" {
				user.LastName = gocloak.StringP(args.LastName)
			}

			userID, err := kc.GC.CreateUser(ctx, token, realm, user)
			if err != nil {
				return toolError(fmt.Sprintf("failed to create user: %v", err))
			}

			// Optionally set password
			if args.Password != "" {
				temporary := false
				if args.TemporaryPassword != nil {
					temporary = *args.TemporaryPassword
				}
				if err := kc.GC.SetPassword(ctx, token, userID, realm, args.Password, temporary); err != nil {
					return toolError(fmt.Sprintf("user created (ID: %s) but failed to set password: %v", userID, err))
				}
			}

			return toolSuccess(fmt.Sprintf("User created with ID: %s", userID))
		},
	)

	// 5. update_user
	mcp.AddTool(s, &mcp.Tool{
		Name:        "update_user",
		Description: "Update an existing user",
	},
		func(ctx context.Context, req *mcp.CallToolRequest, args updateUserArgs) (*mcp.CallToolResult, any, error) {
			token, err := kc.Token(ctx)
			if err != nil {
				return toolError(fmt.Sprintf("token error: %v", err))
			}
			realm := kc.ResolveRealm(args.Realm)

			user, err := kc.GC.GetUserByID(ctx, token, realm, args.UserID)
			if err != nil {
				return toolError(fmt.Sprintf("failed to get user: %v", err))
			}

			if args.Email != nil {
				user.Email = args.Email
			}
			if args.FirstName != nil {
				user.FirstName = args.FirstName
			}
			if args.LastName != nil {
				user.LastName = args.LastName
			}
			if args.Enabled != nil {
				user.Enabled = args.Enabled
			}

			if err := kc.GC.UpdateUser(ctx, token, realm, *user); err != nil {
				return toolError(fmt.Sprintf("failed to update user: %v", err))
			}
			return toolSuccess("User updated successfully")
		},
	)

	// 6. delete_user
	mcp.AddTool(s, &mcp.Tool{
		Name:        "delete_user",
		Description: "Delete a user from a realm",
	},
		func(ctx context.Context, req *mcp.CallToolRequest, args deleteUserArgs) (*mcp.CallToolResult, any, error) {
			token, err := kc.Token(ctx)
			if err != nil {
				return toolError(fmt.Sprintf("token error: %v", err))
			}
			realm := kc.ResolveRealm(args.Realm)

			if err := kc.GC.DeleteUser(ctx, token, realm, args.UserID); err != nil {
				return toolError(fmt.Sprintf("failed to delete user: %v", err))
			}
			return toolSuccess("User deleted successfully")
		},
	)

	// 7. count_users
	mcp.AddTool(s, &mcp.Tool{
		Name:        "count_users",
		Description: "Count users in a realm",
	},
		func(ctx context.Context, req *mcp.CallToolRequest, args countUsersArgs) (*mcp.CallToolResult, any, error) {
			token, err := kc.Token(ctx)
			if err != nil {
				return toolError(fmt.Sprintf("token error: %v", err))
			}
			realm := kc.ResolveRealm(args.Realm)

			params := gocloak.GetUsersParams{}
			if args.Search != "" {
				params.Search = gocloak.StringP(args.Search)
			}

			count, err := kc.GC.GetUserCount(ctx, token, realm, params)
			if err != nil {
				return toolError(fmt.Sprintf("failed to count users: %v", err))
			}
			return toolResult(map[string]int{"count": count})
		},
	)

	// 8. set_user_password
	mcp.AddTool(s, &mcp.Tool{
		Name:        "set_user_password",
		Description: "Set a user's password",
	},
		func(ctx context.Context, req *mcp.CallToolRequest, args setUserPasswordArgs) (*mcp.CallToolResult, any, error) {
			token, err := kc.Token(ctx)
			if err != nil {
				return toolError(fmt.Sprintf("token error: %v", err))
			}
			realm := kc.ResolveRealm(args.Realm)

			temporary := false
			if args.Temporary != nil {
				temporary = *args.Temporary
			}

			if err := kc.GC.SetPassword(ctx, token, args.UserID, realm, args.Password, temporary); err != nil {
				return toolError(fmt.Sprintf("failed to set password: %v", err))
			}
			return toolSuccess("Password set successfully")
		},
	)

	// 9. get_user_credentials
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_user_credentials",
		Description: "Get credentials for a user",
	},
		func(ctx context.Context, req *mcp.CallToolRequest, args getUserCredentialsArgs) (*mcp.CallToolResult, any, error) {
			token, err := kc.Token(ctx)
			if err != nil {
				return toolError(fmt.Sprintf("token error: %v", err))
			}
			realm := kc.ResolveRealm(args.Realm)

			creds, err := kc.GC.GetCredentials(ctx, token, realm, args.UserID)
			if err != nil {
				return toolError(fmt.Sprintf("failed to get credentials: %v", err))
			}
			return toolResult(creds)
		},
	)

	// 10. delete_user_credential
	mcp.AddTool(s, &mcp.Tool{
		Name:        "delete_user_credential",
		Description: "Delete a specific credential for a user",
	},
		func(ctx context.Context, req *mcp.CallToolRequest, args deleteUserCredentialArgs) (*mcp.CallToolResult, any, error) {
			token, err := kc.Token(ctx)
			if err != nil {
				return toolError(fmt.Sprintf("token error: %v", err))
			}
			realm := kc.ResolveRealm(args.Realm)

			if err := kc.GC.DeleteCredentials(ctx, token, realm, args.UserID, args.CredentialID); err != nil {
				return toolError(fmt.Sprintf("failed to delete credential: %v", err))
			}
			return toolSuccess("Credential deleted successfully")
		},
	)

	// 11. send_verify_email
	mcp.AddTool(s, &mcp.Tool{
		Name:        "send_verify_email",
		Description: "Send a verification email to a user",
	},
		func(ctx context.Context, req *mcp.CallToolRequest, args sendVerifyEmailArgs) (*mcp.CallToolResult, any, error) {
			token, err := kc.Token(ctx)
			if err != nil {
				return toolError(fmt.Sprintf("token error: %v", err))
			}
			realm := kc.ResolveRealm(args.Realm)

			params := gocloak.SendVerificationMailParams{}
			if err := kc.GC.SendVerifyEmail(ctx, token, args.UserID, realm, params); err != nil {
				return toolError(fmt.Sprintf("failed to send verify email: %v", err))
			}
			return toolSuccess("Verification email sent")
		},
	)

	// 12. execute_actions_email
	mcp.AddTool(s, &mcp.Tool{
		Name:        "execute_actions_email",
		Description: "Send an actions email to a user",
	},
		func(ctx context.Context, req *mcp.CallToolRequest, args executeActionsEmailArgs) (*mcp.CallToolResult, any, error) {
			token, err := kc.Token(ctx)
			if err != nil {
				return toolError(fmt.Sprintf("token error: %v", err))
			}
			realm := kc.ResolveRealm(args.Realm)

			params := gocloak.ExecuteActionsEmail{
				UserID:  gocloak.StringP(args.UserID),
				Actions: &args.Actions,
			}
			if args.Lifespan != nil {
				params.Lifespan = gocloak.IntP(*args.Lifespan)
			}

			if err := kc.GC.ExecuteActionsEmail(ctx, token, realm, params); err != nil {
				return toolError(fmt.Sprintf("failed to send actions email: %v", err))
			}
			return toolSuccess("Actions email sent")
		},
	)

	// 13. get_user_groups
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_user_groups",
		Description: "Get groups for a user",
	},
		func(ctx context.Context, req *mcp.CallToolRequest, args getUserGroupsArgs) (*mcp.CallToolResult, any, error) {
			token, err := kc.Token(ctx)
			if err != nil {
				return toolError(fmt.Sprintf("token error: %v", err))
			}
			realm := kc.ResolveRealm(args.Realm)

			groups, err := kc.GC.GetUserGroups(ctx, token, realm, args.UserID, gocloak.GetGroupsParams{})
			if err != nil {
				return toolError(fmt.Sprintf("failed to get user groups: %v", err))
			}
			return toolResult(groups)
		},
	)

	// 14. add_user_to_group
	mcp.AddTool(s, &mcp.Tool{
		Name:        "add_user_to_group",
		Description: "Add a user to a group",
	},
		func(ctx context.Context, req *mcp.CallToolRequest, args addUserToGroupArgs) (*mcp.CallToolResult, any, error) {
			token, err := kc.Token(ctx)
			if err != nil {
				return toolError(fmt.Sprintf("token error: %v", err))
			}
			realm := kc.ResolveRealm(args.Realm)

			if err := kc.GC.AddUserToGroup(ctx, token, realm, args.UserID, args.GroupID); err != nil {
				return toolError(fmt.Sprintf("failed to add user to group: %v", err))
			}
			return toolSuccess("User added to group")
		},
	)

	// 15. remove_user_from_group
	mcp.AddTool(s, &mcp.Tool{
		Name:        "remove_user_from_group",
		Description: "Remove a user from a group",
	},
		func(ctx context.Context, req *mcp.CallToolRequest, args removeUserFromGroupArgs) (*mcp.CallToolResult, any, error) {
			token, err := kc.Token(ctx)
			if err != nil {
				return toolError(fmt.Sprintf("token error: %v", err))
			}
			realm := kc.ResolveRealm(args.Realm)

			if err := kc.GC.DeleteUserFromGroup(ctx, token, realm, args.UserID, args.GroupID); err != nil {
				return toolError(fmt.Sprintf("failed to remove user from group: %v", err))
			}
			return toolSuccess("User removed from group")
		},
	)

	// 16. get_user_sessions
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_user_sessions",
		Description: "Get active sessions for a user",
	},
		func(ctx context.Context, req *mcp.CallToolRequest, args getUserSessionsArgs) (*mcp.CallToolResult, any, error) {
			token, err := kc.Token(ctx)
			if err != nil {
				return toolError(fmt.Sprintf("token error: %v", err))
			}
			realm := kc.ResolveRealm(args.Realm)

			sessions, err := kc.GC.GetUserSessions(ctx, token, realm, args.UserID)
			if err != nil {
				return toolError(fmt.Sprintf("failed to get user sessions: %v", err))
			}
			return toolResult(sessions)
		},
	)

	// 17. get_user_federated_identities
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_user_federated_identities",
		Description: "Get federated identities for a user",
	},
		func(ctx context.Context, req *mcp.CallToolRequest, args getUserFederatedIdentitiesArgs) (*mcp.CallToolResult, any, error) {
			token, err := kc.Token(ctx)
			if err != nil {
				return toolError(fmt.Sprintf("token error: %v", err))
			}
			realm := kc.ResolveRealm(args.Realm)

			identities, err := kc.GC.GetUserFederatedIdentities(ctx, token, realm, args.UserID)
			if err != nil {
				return toolError(fmt.Sprintf("failed to get federated identities: %v", err))
			}
			return toolResult(identities)
		},
	)

	// 18. create_user_federated_identity
	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_user_federated_identity",
		Description: "Create a federated identity link for a user",
	},
		func(ctx context.Context, req *mcp.CallToolRequest, args createUserFederatedIdentityArgs) (*mcp.CallToolResult, any, error) {
			token, err := kc.Token(ctx)
			if err != nil {
				return toolError(fmt.Sprintf("token error: %v", err))
			}
			realm := kc.ResolveRealm(args.Realm)

			fedIdentity := gocloak.FederatedIdentityRepresentation{
				UserID:   gocloak.StringP(args.FederatedUserID),
				UserName: gocloak.StringP(args.FederatedUsername),
			}

			if err := kc.GC.CreateUserFederatedIdentity(ctx, token, realm, args.UserID, args.ProviderID, fedIdentity); err != nil {
				return toolError(fmt.Sprintf("failed to create federated identity: %v", err))
			}
			return toolSuccess("Federated identity created")
		},
	)

	// 19. delete_user_federated_identity
	mcp.AddTool(s, &mcp.Tool{
		Name:        "delete_user_federated_identity",
		Description: "Delete a federated identity link for a user",
	},
		func(ctx context.Context, req *mcp.CallToolRequest, args deleteUserFederatedIdentityArgs) (*mcp.CallToolResult, any, error) {
			token, err := kc.Token(ctx)
			if err != nil {
				return toolError(fmt.Sprintf("token error: %v", err))
			}
			realm := kc.ResolveRealm(args.Realm)

			if err := kc.GC.DeleteUserFederatedIdentity(ctx, token, realm, args.UserID, args.ProviderID); err != nil {
				return toolError(fmt.Sprintf("failed to delete federated identity: %v", err))
			}
			return toolSuccess("Federated identity deleted")
		},
	)

	// 20. get_user_realm_roles
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_user_realm_roles",
		Description: "Get realm-level roles assigned to a user",
	},
		func(ctx context.Context, req *mcp.CallToolRequest, args getUserRealmRolesArgs) (*mcp.CallToolResult, any, error) {
			token, err := kc.Token(ctx)
			if err != nil {
				return toolError(fmt.Sprintf("token error: %v", err))
			}
			realm := kc.ResolveRealm(args.Realm)

			roles, err := kc.GC.GetRealmRolesByUserID(ctx, token, realm, args.UserID)
			if err != nil {
				return toolError(fmt.Sprintf("failed to get user realm roles: %v", err))
			}
			return toolResult(roles)
		},
	)

	// 21. add_user_realm_roles
	mcp.AddTool(s, &mcp.Tool{
		Name:        "add_user_realm_roles",
		Description: "Add realm-level roles to a user",
	},
		func(ctx context.Context, req *mcp.CallToolRequest, args addUserRealmRolesArgs) (*mcp.CallToolResult, any, error) {
			token, err := kc.Token(ctx)
			if err != nil {
				return toolError(fmt.Sprintf("token error: %v", err))
			}
			realm := kc.ResolveRealm(args.Realm)

			var roles []gocloak.Role
			for _, roleName := range args.Roles {
				role, err := kc.GC.GetRealmRole(ctx, token, realm, roleName)
				if err != nil {
					return toolError(fmt.Sprintf("failed to get realm role %q: %v", roleName, err))
				}
				roles = append(roles, *role)
			}

			if err := kc.GC.AddRealmRoleToUser(ctx, token, realm, args.UserID, roles); err != nil {
				return toolError(fmt.Sprintf("failed to add realm roles to user: %v", err))
			}
			return toolSuccess("Realm roles added to user")
		},
	)

	// 22. remove_user_realm_roles
	mcp.AddTool(s, &mcp.Tool{
		Name:        "remove_user_realm_roles",
		Description: "Remove realm-level roles from a user",
	},
		func(ctx context.Context, req *mcp.CallToolRequest, args removeUserRealmRolesArgs) (*mcp.CallToolResult, any, error) {
			token, err := kc.Token(ctx)
			if err != nil {
				return toolError(fmt.Sprintf("token error: %v", err))
			}
			realm := kc.ResolveRealm(args.Realm)

			var roles []gocloak.Role
			for _, roleName := range args.Roles {
				role, err := kc.GC.GetRealmRole(ctx, token, realm, roleName)
				if err != nil {
					return toolError(fmt.Sprintf("failed to get realm role %q: %v", roleName, err))
				}
				roles = append(roles, *role)
			}

			if err := kc.GC.DeleteRealmRoleFromUser(ctx, token, realm, args.UserID, roles); err != nil {
				return toolError(fmt.Sprintf("failed to remove realm roles from user: %v", err))
			}
			return toolSuccess("Realm roles removed from user")
		},
	)

	// 23. get_user_client_roles
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_user_client_roles",
		Description: "Get client-level roles assigned to a user",
	},
		func(ctx context.Context, req *mcp.CallToolRequest, args getUserClientRolesArgs) (*mcp.CallToolResult, any, error) {
			token, err := kc.Token(ctx)
			if err != nil {
				return toolError(fmt.Sprintf("token error: %v", err))
			}
			realm := kc.ResolveRealm(args.Realm)

			roles, err := kc.GC.GetClientRolesByUserID(ctx, token, realm, args.ClientID, args.UserID)
			if err != nil {
				return toolError(fmt.Sprintf("failed to get user client roles: %v", err))
			}
			return toolResult(roles)
		},
	)

	// 24. add_user_client_roles
	mcp.AddTool(s, &mcp.Tool{
		Name:        "add_user_client_roles",
		Description: "Add client-level roles to a user",
	},
		func(ctx context.Context, req *mcp.CallToolRequest, args addUserClientRolesArgs) (*mcp.CallToolResult, any, error) {
			token, err := kc.Token(ctx)
			if err != nil {
				return toolError(fmt.Sprintf("token error: %v", err))
			}
			realm := kc.ResolveRealm(args.Realm)

			var roles []gocloak.Role
			for _, roleName := range args.Roles {
				role, err := kc.GC.GetClientRole(ctx, token, realm, args.ClientID, roleName)
				if err != nil {
					return toolError(fmt.Sprintf("failed to get client role %q: %v", roleName, err))
				}
				roles = append(roles, *role)
			}

			if err := kc.GC.AddClientRoleToUser(ctx, token, realm, args.ClientID, args.UserID, roles); err != nil {
				return toolError(fmt.Sprintf("failed to add client roles to user: %v", err))
			}
			return toolSuccess("Client roles added to user")
		},
	)
}
