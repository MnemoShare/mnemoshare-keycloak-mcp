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

type listAuthFlowsArgs struct {
	Realm string `json:"realm,omitempty" jsonschema:"Realm name (uses default if omitted)"`
}

type getAuthFlowArgs struct {
	Realm  string `json:"realm,omitempty" jsonschema:"Realm name (uses default if omitted)"`
	FlowID string `json:"flow_id"         jsonschema:"Authentication flow ID"`
}

type createAuthFlowArgs struct {
	Realm       string `json:"realm,omitempty"       jsonschema:"Realm name (uses default if omitted)"`
	Alias       string `json:"alias"                 jsonschema:"Unique alias for the authentication flow"`
	Description string `json:"description,omitempty" jsonschema:"Description of the authentication flow"`
	ProviderID  string `json:"provider_id"           jsonschema:"Provider ID for the flow (e.g. basic-flow)"`
	TopLevel    *bool  `json:"top_level,omitempty"   jsonschema:"Whether this is a top-level flow (default true)"`
	BuiltIn     *bool  `json:"built_in,omitempty"    jsonschema:"Whether this is a built-in flow (default false)"`
}

type deleteAuthFlowArgs struct {
	Realm  string `json:"realm,omitempty" jsonschema:"Realm name (uses default if omitted)"`
	FlowID string `json:"flow_id"         jsonschema:"Authentication flow ID to delete"`
}

type getAuthFlowExecutionsArgs struct {
	Realm     string `json:"realm,omitempty" jsonschema:"Realm name (uses default if omitted)"`
	FlowAlias string `json:"flow_alias"      jsonschema:"Alias of the authentication flow"`
}

type updateAuthFlowExecutionArgs struct {
	Realm       string `json:"realm,omitempty" jsonschema:"Realm name (uses default if omitted)"`
	FlowAlias   string `json:"flow_alias"      jsonschema:"Alias of the parent authentication flow"`
	ExecutionID string `json:"execution_id"    jsonschema:"ID of the execution to update"`
	Requirement string `json:"requirement"     jsonschema:"Requirement level: REQUIRED ALTERNATIVE DISABLED or CONDITIONAL"`
}

type listRequiredActionsArgs struct {
	Realm string `json:"realm,omitempty" jsonschema:"Realm name (uses default if omitted)"`
}

type getRequiredActionArgs struct {
	Realm string `json:"realm,omitempty" jsonschema:"Realm name (uses default if omitted)"`
	Alias string `json:"alias"           jsonschema:"Alias of the required action"`
}

type updateRequiredActionArgs struct {
	Realm         string  `json:"realm,omitempty"          jsonschema:"Realm name (uses default if omitted)"`
	Alias         string  `json:"alias"                    jsonschema:"Alias of the required action to update"`
	Name          *string `json:"name,omitempty"           jsonschema:"Display name of the required action"`
	Enabled       *bool   `json:"enabled,omitempty"        jsonschema:"Whether the required action is enabled"`
	DefaultAction *bool   `json:"default_action,omitempty" jsonschema:"Whether this is a default action for new users"`
}

type deleteRequiredActionArgs struct {
	Realm string `json:"realm,omitempty" jsonschema:"Realm name (uses default if omitted)"`
	Alias string `json:"alias"           jsonschema:"Alias of the required action to delete"`
}

// ---------------------------------------------------------------------------
// Registration
// ---------------------------------------------------------------------------

func registerAuthFlowTools(s *mcp.Server, kc *keycloak.Client) {

	// 1. list_auth_flows
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_auth_flows",
		Description: "List all authentication flows in a realm",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args listAuthFlowsArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get token: %v", err))
		}

		realm := kc.ResolveRealm(args.Realm)

		flows, err := kc.GC.GetAuthenticationFlows(ctx, token, realm)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to list authentication flows: %v", err))
		}

		return toolResult(flows)
	})

	// 2. get_auth_flow
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_auth_flow",
		Description: "Get an authentication flow by ID",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args getAuthFlowArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get token: %v", err))
		}

		realm := kc.ResolveRealm(args.Realm)

		flow, err := kc.GC.GetAuthenticationFlow(ctx, token, realm, args.FlowID)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get authentication flow %q: %v", args.FlowID, err))
		}

		return toolResult(flow)
	})

	// 3. create_auth_flow
	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_auth_flow",
		Description: "Create a new authentication flow in a realm",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args createAuthFlowArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get token: %v", err))
		}

		realm := kc.ResolveRealm(args.Realm)

		topLevel := true
		if args.TopLevel != nil {
			topLevel = *args.TopLevel
		}
		builtIn := false
		if args.BuiltIn != nil {
			builtIn = *args.BuiltIn
		}

		flowRep := gocloak.AuthenticationFlowRepresentation{
			Alias:      gocloak.StringP(args.Alias),
			ProviderID: gocloak.StringP(args.ProviderID),
			TopLevel:   gocloak.BoolP(topLevel),
			BuiltIn:    gocloak.BoolP(builtIn),
		}
		if args.Description != "" {
			flowRep.Description = gocloak.StringP(args.Description)
		}

		err = kc.GC.CreateAuthenticationFlow(ctx, token, realm, flowRep)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to create authentication flow %q: %v", args.Alias, err))
		}

		return toolSuccess(fmt.Sprintf("Authentication flow %q created successfully", args.Alias))
	})

	// 4. delete_auth_flow
	mcp.AddTool(s, &mcp.Tool{
		Name:        "delete_auth_flow",
		Description: "Delete an authentication flow from a realm",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args deleteAuthFlowArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get token: %v", err))
		}

		realm := kc.ResolveRealm(args.Realm)

		if err := kc.GC.DeleteAuthenticationFlow(ctx, token, realm, args.FlowID); err != nil {
			return toolError(fmt.Sprintf("Error: failed to delete authentication flow %q: %v", args.FlowID, err))
		}

		return toolSuccess(fmt.Sprintf("Authentication flow %q deleted successfully", args.FlowID))
	})

	// 5. get_auth_flow_executions
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_auth_flow_executions",
		Description: "Get executions for an authentication flow",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args getAuthFlowExecutionsArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get token: %v", err))
		}

		realm := kc.ResolveRealm(args.Realm)

		executions, err := kc.GC.GetAuthenticationExecutions(ctx, token, realm, args.FlowAlias)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get executions for flow %q: %v", args.FlowAlias, err))
		}

		return toolResult(executions)
	})

	// 6. update_auth_flow_execution
	mcp.AddTool(s, &mcp.Tool{
		Name:        "update_auth_flow_execution",
		Description: "Update an execution within an authentication flow",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args updateAuthFlowExecutionArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get token: %v", err))
		}

		realm := kc.ResolveRealm(args.Realm)

		execution := gocloak.ModifyAuthenticationExecutionRepresentation{
			ID:          gocloak.StringP(args.ExecutionID),
			Requirement: gocloak.StringP(args.Requirement),
		}

		if err := kc.GC.UpdateAuthenticationExecution(ctx, token, realm, args.FlowAlias, execution); err != nil {
			return toolError(fmt.Sprintf("Error: failed to update execution %q in flow %q: %v", args.ExecutionID, args.FlowAlias, err))
		}

		return toolSuccess(fmt.Sprintf("Execution %q updated to %q in flow %q", args.ExecutionID, args.Requirement, args.FlowAlias))
	})

	// 7. list_required_actions
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_required_actions",
		Description: "List all required actions in a realm",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args listRequiredActionsArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get token: %v", err))
		}

		realm := kc.ResolveRealm(args.Realm)

		actions, err := kc.GC.GetRequiredActions(ctx, token, realm)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to list required actions: %v", err))
		}

		return toolResult(actions)
	})

	// 8. get_required_action
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_required_action",
		Description: "Get a required action by alias",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args getRequiredActionArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get token: %v", err))
		}

		realm := kc.ResolveRealm(args.Realm)

		action, err := kc.GC.GetRequiredAction(ctx, token, realm, args.Alias)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get required action %q: %v", args.Alias, err))
		}

		return toolResult(action)
	})

	// 9. update_required_action
	mcp.AddTool(s, &mcp.Tool{
		Name:        "update_required_action",
		Description: "Update a required action in a realm",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args updateRequiredActionArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get token: %v", err))
		}

		realm := kc.ResolveRealm(args.Realm)

		// Fetch the existing required action to apply partial updates.
		action, err := kc.GC.GetRequiredAction(ctx, token, realm, args.Alias)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get required action %q for update: %v", args.Alias, err))
		}

		if args.Name != nil {
			action.Name = args.Name
		}
		if args.Enabled != nil {
			action.Enabled = args.Enabled
		}
		if args.DefaultAction != nil {
			action.DefaultAction = args.DefaultAction
		}

		if err := kc.GC.UpdateRequiredAction(ctx, token, realm, *action); err != nil {
			return toolError(fmt.Sprintf("Error: failed to update required action %q: %v", args.Alias, err))
		}

		return toolSuccess(fmt.Sprintf("Required action %q updated successfully", args.Alias))
	})

	// 10. delete_required_action
	mcp.AddTool(s, &mcp.Tool{
		Name:        "delete_required_action",
		Description: "Delete a required action from a realm",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args deleteRequiredActionArgs) (*mcp.CallToolResult, any, error) {
		token, err := kc.Token(ctx)
		if err != nil {
			return toolError(fmt.Sprintf("Error: failed to get token: %v", err))
		}

		realm := kc.ResolveRealm(args.Realm)

		if err := kc.GC.DeleteRequiredAction(ctx, token, realm, args.Alias); err != nil {
			return toolError(fmt.Sprintf("Error: failed to delete required action %q: %v", args.Alias, err))
		}

		return toolSuccess(fmt.Sprintf("Required action %q deleted successfully", args.Alias))
	})
}
