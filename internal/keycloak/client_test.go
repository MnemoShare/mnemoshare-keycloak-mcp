package keycloak

import (
	"testing"

	"github.com/mnemoshare/mnemoshare-keycloak-mcp/internal/auth"
	"github.com/mnemoshare/mnemoshare-keycloak-mcp/internal/config"
)

func TestAdminRealmURL(t *testing.T) {
	cfg := &config.Config{KeycloakURL: "http://kc.example:8080/"}
	c := NewClient(cfg, auth.NewTokenManager(cfg))

	got := c.AdminRealmURL("my realm/x", "groups", "a/b?c", "children")
	want := "http://kc.example:8080/admin/realms/my%20realm%2Fx/groups/a%2Fb%3Fc/children"
	if got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}
