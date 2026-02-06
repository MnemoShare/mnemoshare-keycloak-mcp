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

type listRealmsArgs struct{}

type getRealmArgs struct {
	Realm string `json:"realm" jsonschema:"description=The realm name,required"`
}

type createRealmArgs struct {
	Realm       string `json:"realm" jsonschema:"description=The name of the new realm,required"`
	Enabled     *bool  `json:"enabled,omitempty" jsonschema:"description=Whether the realm is enabled (default true)"`
	DisplayName string `json:"display_name,omitempty" jsonschema:"description=A human-friendly display name for the realm"`
}

type updateRealmArgs struct {
	Realm                 string  `json:"realm" jsonschema:"description=The realm name to update,required"`
	Enabled               *bool   `json:"enabled,omitempty" jsonschema:"description=Whether the realm is enabled"`
	DisplayName           *string `json:"display_name,omitempty" jsonschema:"description=A human-friendly display name for the realm"`
	RegistrationAllowed   *bool   `json:"registration_allowed,omitempty" jsonschema:"description=Whether user self-registration is allowed"`
	ResetPasswordAllowed  *bool   `json:"reset_password_allowed,omitempty" jsonschema:"description=Whether password reset is allowed"`
	RememberMe            *bool   `json:"remember_me,omitempty" jsonschema:"description=Whether the remember-me option is enabled"`
	VerifyEmail           *bool   `json:"verify_email,omitempty" jsonschema:"description=Whether email verification is required"`
	LoginWithEmailAllowed *bool   `json:"login_with_email_allowed,omitempty" jsonschema:"description=Whether login with email is allowed"`
	DuplicateEmailsAllowed *bool  `json:"duplicate_emails_allowed,omitempty" jsonschema:"description=Whether duplicate email addresses are allowed"`
	SSLRequired           *string `json:"ssl_required,omitempty" jsonschema:"description=SSL requirement (none/external/all)"`
}

type deleteRealmArgs struct {
	Realm string `json:"realm" jsonschema:"description=The realm name to delete,required"`
}

type clearRealmCacheArgs struct {
	Realm string `json:"realm,omitempty" jsonschema:"description=The realm name (uses default realm if omitted)"`
}

type clearUserCacheArgs struct {
	Realm string `json:"realm,omitempty" jsonschema:"description=The realm name (uses default realm if omitted)"`
}

type clearKeysCacheArgs struct {
	Realm string `json:"realm,omitempty" jsonschema:"description=The realm name (uses default realm if omitted)"`
}

// ---------------------------------------------------------------------------
// Registration
// ---------------------------------------------------------------------------

func registerRealmTools(s *mcp.Server, kc *keycloak.Client) {
	// 1. list_realms
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_realms",
		Description: "List all Keycloak realms",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args listRealmsArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get token: %v", err))
		}

		realms, err := kc.GC.GetRealms(ctx, token)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to list realms: %v", err))
		}

		return toolResult(realms)
	})

	// 2. get_realm
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_realm",
		Description: "Get a Keycloak realm by name",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args getRealmArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get token: %v", err))
		}

		realm, err := kc.GC.GetRealm(ctx, token, args.Realm)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get realm %q: %v", args.Realm, err))
		}

		return toolResult(realm)
	})

	// 3. create_realm
	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_realm",
		Description: "Create a new Keycloak realm",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args createRealmArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get token: %v", err))
		}

		enabled := true
		if args.Enabled != nil {
			enabled = *args.Enabled
		}

		realmRep := gocloak.RealmRepresentation{
			Realm:   gocloak.StringP(args.Realm),
			Enabled: gocloak.BoolP(enabled),
		}
		if args.DisplayName != "" {
			realmRep.DisplayName = gocloak.StringP(args.DisplayName)
		}

		createdID, err := kc.GC.CreateRealm(ctx, token, realmRep)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to create realm %q: %v", args.Realm, err))
		}

		return toolSuccess(fmt.Sprintf("Realm %q created successfully (id: %s)", args.Realm, createdID))
	})

	// 4. update_realm
	mcp.AddTool(s, &mcp.Tool{
		Name:        "update_realm",
		Description: "Update settings on an existing Keycloak realm",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args updateRealmArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get token: %v", err))
		}

		// Fetch the current representation so we only override supplied fields.
		existing, err := kc.GC.GetRealm(ctx, token, args.Realm)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get realm %q for update: %v", args.Realm, err))
		}

		if args.Enabled != nil {
			existing.Enabled = args.Enabled
		}
		if args.DisplayName != nil {
			existing.DisplayName = args.DisplayName
		}
		if args.RegistrationAllowed != nil {
			existing.RegistrationAllowed = args.RegistrationAllowed
		}
		if args.ResetPasswordAllowed != nil {
			existing.ResetPasswordAllowed = args.ResetPasswordAllowed
		}
		if args.RememberMe != nil {
			existing.RememberMe = args.RememberMe
		}
		if args.VerifyEmail != nil {
			existing.VerifyEmail = args.VerifyEmail
		}
		if args.LoginWithEmailAllowed != nil {
			existing.LoginWithEmailAllowed = args.LoginWithEmailAllowed
		}
		if args.DuplicateEmailsAllowed != nil {
			existing.DuplicateEmailsAllowed = args.DuplicateEmailsAllowed
		}
		if args.SSLRequired != nil {
			existing.SslRequired = args.SSLRequired
		}

		if err := kc.GC.UpdateRealm(ctx, token, *existing); err != nil {
			return toolError(fmt.Sprintf("Error: failed to update realm %q: %v", args.Realm, err))
		}

		return toolSuccess(fmt.Sprintf("Realm %q updated successfully", args.Realm))
	})

	// 5. delete_realm
	mcp.AddTool(s, &mcp.Tool{
		Name:        "delete_realm",
		Description: "Delete a Keycloak realm",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args deleteRealmArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get token: %v", err))
		}

		if err := kc.GC.DeleteRealm(ctx, token, args.Realm); err != nil {
			return toolError(fmt.Sprintf("Error: failed to delete realm %q: %v", args.Realm, err))
		}

		return toolSuccess(fmt.Sprintf("Realm %q deleted successfully", args.Realm))
	})

	// 6. clear_realm_cache
	mcp.AddTool(s, &mcp.Tool{
		Name:        "clear_realm_cache",
		Description: "Clear the realm cache in Keycloak",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args clearRealmCacheArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get token: %v", err))
		}

		realm := kc.ResolveRealm(args.Realm)
		if err := kc.GC.ClearRealmCache(ctx, token, realm); err != nil {
			return toolError(fmt.Sprintf("Error: failed to clear realm cache for %q: %v", realm, err))
		}

		return toolSuccess(fmt.Sprintf("Realm cache cleared for %q", realm))
	})

	// 7. clear_user_cache
	mcp.AddTool(s, &mcp.Tool{
		Name:        "clear_user_cache",
		Description: "Clear the user cache in Keycloak",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args clearUserCacheArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get token: %v", err))
		}

		realm := kc.ResolveRealm(args.Realm)
		if err := kc.GC.ClearUserCache(ctx, token, realm); err != nil {
			return toolError(fmt.Sprintf("Error: failed to clear user cache for %q: %v", realm, err))
		}

		return toolSuccess(fmt.Sprintf("User cache cleared for %q", realm))
	})

	// 8. clear_keys_cache
	mcp.AddTool(s, &mcp.Tool{
		Name:        "clear_keys_cache",
		Description: "Clear the keys cache in Keycloak",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args clearKeysCacheArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get token: %v", err))
		}

		realm := kc.ResolveRealm(args.Realm)
		if err := kc.GC.ClearKeysCache(ctx, token, realm); err != nil {
			return toolError(fmt.Sprintf("Error: failed to clear keys cache for %q: %v", realm, err))
		}

		return toolSuccess(fmt.Sprintf("Keys cache cleared for %q", realm))
	})
}
