package keycloak

import (
	"context"
	"net/url"
	"strings"

	"github.com/Nerzal/gocloak/v13"

	"github.com/mnemoshare/mnemoshare-keycloak-mcp/internal/auth"
	"github.com/mnemoshare/mnemoshare-keycloak-mcp/internal/config"
)

// Client wraps gocloak with automatic token injection and realm resolution.
type Client struct {
	GC           *gocloak.GoCloak
	tokenManager *auth.TokenManager
	defaultRealm string
	baseURL      string
}

func NewClient(cfg *config.Config, tm *auth.TokenManager) *Client {
	return &Client{
		GC:           tm.GoCloak(),
		tokenManager: tm,
		defaultRealm: cfg.DefaultRealm,
		baseURL:      strings.TrimRight(cfg.KeycloakURL, "/"),
	}
}

// Token returns a valid access token string.
func (c *Client) Token(ctx context.Context) (string, error) {
	return c.tokenManager.Token(ctx)
}

// ResolveRealm returns the provided realm or falls back to the configured default.
func (c *Client) ResolveRealm(realm string) string {
	if realm != "" {
		return realm
	}
	if c.defaultRealm != "" {
		return c.defaultRealm
	}
	return "master"
}

// AdminRealmURL builds {base}/admin/realms/{realm}/{segments...}; every
// variable segment must be passed separately so it is path-escaped.
func (c *Client) AdminRealmURL(realm string, segments ...string) string {
	u := c.baseURL + "/admin/realms/" + url.PathEscape(realm)
	for _, seg := range segments {
		u += "/" + url.PathEscape(seg)
	}
	return u
}
