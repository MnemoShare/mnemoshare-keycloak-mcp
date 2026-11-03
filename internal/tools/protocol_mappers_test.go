package tools

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

const existingMapper = `{"id":"m1","name":"orig","protocol":"openid-connect","protocolMapper":"oidc-hardcoded-claim-mapper","consentRequired":false,"config":{"A":"1","B":"2","C":"3"}}`

func mapperHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if strings.HasSuffix(r.URL.Path, "/m1") {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(existingMapper))
			return
		}
		http.NotFound(w, r)
	case http.MethodPost:
		w.Header().Set("Location", r.URL.Path+"/newid")
		w.WriteHeader(http.StatusCreated)
	default:
		w.WriteHeader(http.StatusNoContent)
	}
}

func findReq(t *testing.T, reqs *requestLog, method string) recordedRequest {
	t.Helper()
	for _, r := range reqs.all() {
		if r.method == method {
			return r
		}
	}
	t.Fatalf("no %s request recorded: %+v", method, reqs.all())
	return recordedRequest{}
}

func TestCreateClientScopeProtocolMapperKeepsAllConfigKeys(t *testing.T) {
	cs, reqs := newTestSession(t, mapperHandler)
	cfg := map[string]any{"introspection.token.claim": "true", "lightweight.claim": "false", "claim.name": "x"}
	res := callTool(t, cs, "create_client_scope_protocol_mapper", map[string]any{
		"scope_id": "s1", "name": "n", "protocol": "openid-connect", "mapper_type": "t", "config": cfg,
	})
	if res.IsError {
		t.Fatalf("unexpected error: %s", resultText(res))
	}
	got := findReq(t, reqs, http.MethodPost)
	if got.path != "/admin/realms/r1/client-scopes/s1/protocol-mappers/models" {
		t.Fatalf("path %s", got.path)
	}
	var body struct {
		Config map[string]any `json:"config"`
	}
	if err := json.Unmarshal([]byte(got.body), &body); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(body.Config, cfg) {
		t.Fatalf("config = %v, want %v", body.Config, cfg)
	}
	if !strings.Contains(resultText(res), "newid") {
		t.Fatalf("missing id: %s", resultText(res))
	}
}

func TestCreateClientProtocolMapperKeepsAllConfigKeys(t *testing.T) {
	cs, reqs := newTestSession(t, mapperHandler)
	cfg := map[string]any{"introspection.token.claim": "true", "addOrganizationAttributes": "true"}
	res := callTool(t, cs, "create_client_protocol_mapper", map[string]any{
		"id": "c1", "name": "n", "protocol": "openid-connect", "mapper_type": "t", "config": cfg,
	})
	if res.IsError {
		t.Fatalf("unexpected error: %s", resultText(res))
	}
	var body struct {
		Config map[string]any `json:"config"`
	}
	if err := json.Unmarshal([]byte(findReq(t, reqs, http.MethodPost).body), &body); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(body.Config, cfg) {
		t.Fatalf("config = %v, want %v", body.Config, cfg)
	}
}

func TestUpdateProtocolMapperMergesConfig(t *testing.T) {
	cases := []struct {
		tool string
		args map[string]any
		path string
	}{
		{"update_client_protocol_mapper", map[string]any{"id": "c1", "mapper_id": "m1"}, "/admin/realms/r1/clients/c1/protocol-mappers/models/m1"},
		{"update_client_scope_protocol_mapper", map[string]any{"scope_id": "s1", "mapper_id": "m1"}, "/admin/realms/r1/client-scopes/s1/protocol-mappers/models/m1"},
	}
	for _, tc := range cases {
		t.Run(tc.tool, func(t *testing.T) {
			cs, reqs := newTestSession(t, mapperHandler)
			tc.args["config"] = map[string]any{"B": "changed", "D": "4", "introspection.token.claim": "true"}
			res := callTool(t, cs, tc.tool, tc.args)
			if res.IsError {
				t.Fatalf("unexpected error: %s", resultText(res))
			}
			put := findReq(t, reqs, http.MethodPut)
			if put.path != tc.path {
				t.Fatalf("path %s", put.path)
			}
			var body map[string]any
			if err := json.Unmarshal([]byte(put.body), &body); err != nil {
				t.Fatal(err)
			}
			want := map[string]any{"A": "1", "B": "changed", "C": "3", "D": "4", "introspection.token.claim": "true"}
			if !reflect.DeepEqual(body["config"], want) {
				t.Fatalf("config = %v, want %v", body["config"], want)
			}
			if body["id"] != "m1" || body["protocol"] != "openid-connect" || body["protocolMapper"] != "oidc-hardcoded-claim-mapper" || body["name"] != "orig" {
				t.Fatalf("representation not preserved: %s", put.body)
			}

			cs2, reqs2 := newTestSession(t, mapperHandler)
			tc.args["name"] = "renamed"
			delete(tc.args, "config")
			if res := callTool(t, cs2, tc.tool, tc.args); res.IsError {
				t.Fatalf("unexpected error: %s", resultText(res))
			}
			var body2 map[string]any
			if err := json.Unmarshal([]byte(findReq(t, reqs2, http.MethodPut).body), &body2); err != nil {
				t.Fatal(err)
			}
			if body2["name"] != "renamed" || !reflect.DeepEqual(body2["config"], map[string]any{"A": "1", "B": "2", "C": "3"}) {
				t.Fatalf("name-only update wrong: %v", body2)
			}
		})
	}
}

func TestUpdateProtocolMapperNotFound(t *testing.T) {
	for tool, args := range map[string]map[string]any{
		"update_client_protocol_mapper":       {"id": "c1", "mapper_id": "nope", "config": map[string]any{"x": "y"}},
		"update_client_scope_protocol_mapper": {"scope_id": "s1", "mapper_id": "nope", "config": map[string]any{"x": "y"}},
	} {
		t.Run(tool, func(t *testing.T) {
			cs, reqs := newTestSession(t, mapperHandler)
			res := callTool(t, cs, tool, args)
			if !res.IsError || !strings.Contains(resultText(res), "not found") {
				t.Fatalf("want not-found error, got %q", resultText(res))
			}
			for _, r := range reqs.all() {
				if r.method == http.MethodPut {
					t.Fatalf("unexpected PUT")
				}
			}
		})
	}
}

func TestProtocolMapperEmptyIDsRejected(t *testing.T) {
	cfg := map[string]any{"x": "y"}
	cases := map[string]map[string]any{
		"create_client_scope_protocol_mapper": {"scope_id": "", "name": "n", "protocol": "p", "mapper_type": "t", "config": cfg},
		"update_client_scope_protocol_mapper": {"scope_id": "", "mapper_id": "m1"},
		"update_client_protocol_mapper":       {"id": "c1", "mapper_id": ""},
	}
	for tool, args := range cases {
		t.Run(tool, func(t *testing.T) {
			cs, reqs := newTestSession(t, mapperHandler)
			res := callTool(t, cs, tool, args)
			if !res.IsError {
				t.Fatalf("want error, got %s", resultText(res))
			}
			if n := len(reqs.all()); n != 0 {
				t.Fatalf("%d keycloak calls", n)
			}
		})
	}
}
