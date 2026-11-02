package tools

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/mnemoshare/mnemoshare-keycloak-mcp/internal/auth"
	"github.com/mnemoshare/mnemoshare-keycloak-mcp/internal/config"
	"github.com/mnemoshare/mnemoshare-keycloak-mcp/internal/keycloak"
)

type recordedRequest struct {
	method, path, rawQuery, body string
}

func newTestSession(t *testing.T, handler func(w http.ResponseWriter, r *http.Request)) (*mcp.ClientSession, *[]recordedRequest) {
	t.Helper()
	var reqs []recordedRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/protocol/openid-connect/token") {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"access_token":"tok","expires_in":300}`))
			return
		}
		b, _ := io.ReadAll(r.Body)
		reqs = append(reqs, recordedRequest{r.Method, r.URL.EscapedPath(), r.URL.RawQuery, string(b)})
		handler(w, r)
	}))
	t.Cleanup(srv.Close)

	cfg := &config.Config{
		KeycloakURL:        srv.URL,
		KeycloakRealm:      "master",
		AuthMode:           "password",
		AdminUser:          "admin",
		AdminPassword:      "pw",
		DefaultRealm:       "r1",
		TokenRefreshBuffer: time.Second,
	}
	kc := keycloak.NewClient(cfg, auth.NewTokenManager(cfg))

	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0"}, nil)
	registerGroupTools(server, kc)

	ct, st := mcp.NewInMemoryTransports()
	ctx := context.Background()
	ss, err := server.Connect(ctx, st, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = ss.Close() })
	client := mcp.NewClient(&mcp.Implementation{Name: "c", Version: "0"}, nil)
	cs, err := client.Connect(ctx, ct, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = cs.Close() })
	return cs, &reqs
}

func callTool(t *testing.T, cs *mcp.ClientSession, name string, args map[string]any) *mcp.CallToolResult {
	t.Helper()
	res, err := cs.CallTool(context.Background(), &mcp.CallToolParams{Name: name, Arguments: args})
	if err != nil {
		t.Fatal(err)
	}
	return res
}

func resultText(res *mcp.CallToolResult) string {
	if len(res.Content) == 0 {
		return ""
	}
	if tc, ok := res.Content[0].(*mcp.TextContent); ok {
		return tc.Text
	}
	return ""
}

func TestListGroupChildren(t *testing.T) {
	cs, reqs := newTestSession(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"id":"c1","name":"child"}]`))
	})
	res := callTool(t, cs, "list_group_children", map[string]any{
		"group_id": "g/1", "search": "ch", "first": 5, "max": 10,
	})
	if res.IsError {
		t.Fatalf("unexpected error: %s", resultText(res))
	}
	if !strings.Contains(resultText(res), `"child"`) {
		t.Fatalf("missing child in result: %s", resultText(res))
	}
	got := (*reqs)[0]
	if got.method != "GET" || got.path != "/admin/realms/r1/groups/g%2F1/children" {
		t.Fatalf("unexpected request %s %s", got.method, got.path)
	}
	for _, p := range []string{"search=ch", "first=5", "max=10"} {
		if !strings.Contains(got.rawQuery, p) {
			t.Fatalf("query %q missing %s", got.rawQuery, p)
		}
	}
}

func TestListGroupChildrenHTTPError(t *testing.T) {
	cs, _ := newTestSession(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Could not find group"}`))
	})
	res := callTool(t, cs, "list_group_children", map[string]any{"group_id": "nope"})
	if !res.IsError || !strings.Contains(resultText(res), "404") || !strings.Contains(resultText(res), "Could not find group") {
		t.Fatalf("expected 404 error with body, got %v %s", res.IsError, resultText(res))
	}
}

func TestAddGroupClientRoles(t *testing.T) {
	cs, reqs := newTestSession(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id":"role-id","name":"viewer"}`))
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
	res := callTool(t, cs, "add_group_client_roles", map[string]any{
		"group_id": "g1", "client_id": "c1", "roles": []string{"viewer"},
	})
	if res.IsError {
		t.Fatalf("unexpected error: %s", resultText(res))
	}
	last := (*reqs)[len(*reqs)-1]
	if last.method != "POST" || last.path != "/admin/realms/r1/groups/g1/role-mappings/clients/c1" {
		t.Fatalf("unexpected request %s %s", last.method, last.path)
	}
	var roles []map[string]any
	if err := json.Unmarshal([]byte(last.body), &roles); err != nil || len(roles) != 1 || roles[0]["name"] != "viewer" {
		t.Fatalf("unexpected body %s (%v)", last.body, err)
	}
}

func TestGroupClientRolesRejectEmpty(t *testing.T) {
	cs, reqs := newTestSession(t, func(w http.ResponseWriter, r *http.Request) {})
	for _, name := range []string{"add_group_client_roles", "remove_group_client_roles"} {
		res := callTool(t, cs, name, map[string]any{"group_id": "g1", "client_id": "c1", "roles": []string{}})
		if !res.IsError {
			t.Fatalf("%s: expected error for empty roles", name)
		}
	}
	if len(*reqs) != 0 {
		t.Fatalf("expected no Keycloak calls, got %d", len(*reqs))
	}
}
