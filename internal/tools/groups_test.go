package tools

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
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

type requestLog struct {
	mu   sync.Mutex
	reqs []recordedRequest
}

func (l *requestLog) add(r recordedRequest) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.reqs = append(l.reqs, r)
}

func (l *requestLog) all() []recordedRequest {
	l.mu.Lock()
	defer l.mu.Unlock()
	return append([]recordedRequest(nil), l.reqs...)
}

func newTestSession(t *testing.T, handler func(w http.ResponseWriter, r *http.Request)) (*mcp.ClientSession, *requestLog) {
	t.Helper()
	reqs := &requestLog{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/protocol/openid-connect/token") {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"access_token":"tok","expires_in":300}`))
			return
		}
		b, _ := io.ReadAll(r.Body)
		reqs.add(recordedRequest{r.Method, r.URL.EscapedPath(), r.URL.RawQuery, string(b)})
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
	RegisterAll(server, kc)

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
	return cs, reqs
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
	got := reqs.all()[0]
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
	all := reqs.all()
	if len(all) != 2 {
		t.Fatalf("expected 2 requests, got %d", len(all))
	}
	if all[0].method != "GET" || all[0].path != "/admin/realms/r1/clients/c1/roles/viewer" {
		t.Fatalf("unexpected lookup %s %s", all[0].method, all[0].path)
	}
	last := all[1]
	if last.method != "POST" || last.path != "/admin/realms/r1/groups/g1/role-mappings/clients/c1" {
		t.Fatalf("unexpected request %s %s", last.method, last.path)
	}
	var roles []map[string]any
	if err := json.Unmarshal([]byte(last.body), &roles); err != nil || len(roles) != 1 || roles[0]["name"] != "viewer" || roles[0]["id"] != "role-id" {
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
	if len(reqs.all()) != 0 {
		t.Fatalf("expected no Keycloak calls, got %d", len(reqs.all()))
	}
}

func roleLookupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"role-id","name":"viewer"}`))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func TestRemoveClientRoleMappings(t *testing.T) {
	cases := []struct {
		tool string
		args map[string]any
		path string
	}{
		{"remove_group_client_roles", map[string]any{"group_id": "g1", "client_id": "c1", "roles": []string{"viewer"}}, "/admin/realms/r1/groups/g1/role-mappings/clients/c1"},
		{"remove_user_client_roles", map[string]any{"user_id": "u1", "client_id": "c1", "roles": []string{"viewer"}}, "/admin/realms/r1/users/u1/role-mappings/clients/c1"},
	}
	for _, tc := range cases {
		t.Run(tc.tool, func(t *testing.T) {
			cs, reqs := newTestSession(t, roleLookupHandler)
			res := callTool(t, cs, tc.tool, tc.args)
			if res.IsError {
				t.Fatalf("unexpected error: %s", resultText(res))
			}
			all := reqs.all()
			last := all[len(all)-1]
			if last.method != "DELETE" || last.path != tc.path {
				t.Fatalf("unexpected request %s %s", last.method, last.path)
			}
			var roles []map[string]any
			if err := json.Unmarshal([]byte(last.body), &roles); err != nil || len(roles) != 1 || roles[0]["name"] != "viewer" || roles[0]["id"] != "role-id" {
				t.Fatalf("unexpected body %s (%v)", last.body, err)
			}
		})
	}
}

func TestGetGroupsByClientRole(t *testing.T) {
	cs, reqs := newTestSession(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"id":"g1","name":"grp"}]`))
	})
	res := callTool(t, cs, "get_groups_by_client_role", map[string]any{"client_id": "c1", "role_name": "viewer"})
	if res.IsError || !strings.Contains(resultText(res), `"grp"`) {
		t.Fatalf("unexpected result: %v %s", res.IsError, resultText(res))
	}
	got := reqs.all()[0]
	if got.method != "GET" || got.path != "/admin/realms/r1/clients/c1/roles/viewer/groups" {
		t.Fatalf("unexpected request %s %s", got.method, got.path)
	}
}

func TestToolCount(t *testing.T) {
	cs, _ := newTestSession(t, func(w http.ResponseWriter, r *http.Request) {})
	res, err := cs.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	const want = 139
	if len(res.Tools) != want {
		t.Fatalf("registered %d tools, want %d (update docs when changing)", len(res.Tools), want)
	}
}
