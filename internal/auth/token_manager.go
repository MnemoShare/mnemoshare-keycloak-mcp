package auth

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Nerzal/gocloak/v13"
	"github.com/rs/zerolog/log"

	"github.com/mnemoshare/mnemoshare-keycloak-mcp/internal/config"
)

// TokenManager handles Keycloak token acquisition and auto-refresh.
type TokenManager struct {
	gc     *gocloak.GoCloak
	cfg    *config.Config
	mu     sync.RWMutex
	token  *gocloak.JWT
	expiry time.Time
}

func NewTokenManager(cfg *config.Config) *TokenManager {
	gc := gocloak.NewClient(cfg.KeycloakURL)
	return &TokenManager{
		gc:  gc,
		cfg: cfg,
	}
}

// Token returns a valid access token, refreshing if necessary.
func (tm *TokenManager) Token(ctx context.Context) (string, error) {
	tm.mu.RLock()
	if tm.token != nil && time.Now().Before(tm.expiry) {
		t := tm.token.AccessToken
		tm.mu.RUnlock()
		return t, nil
	}
	tm.mu.RUnlock()

	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Double-check after acquiring write lock.
	if tm.token != nil && time.Now().Before(tm.expiry) {
		return tm.token.AccessToken, nil
	}

	jwt, err := tm.authenticate(ctx)
	if err != nil {
		return "", fmt.Errorf("token acquisition failed: %w", err)
	}

	tm.token = jwt
	tm.expiry = time.Now().Add(time.Duration(jwt.ExpiresIn)*time.Second - tm.cfg.TokenRefreshBuffer)
	log.Debug().Time("expiry", tm.expiry).Msg("token acquired")

	return jwt.AccessToken, nil
}

func (tm *TokenManager) authenticate(ctx context.Context) (*gocloak.JWT, error) {
	switch tm.cfg.AuthMode {
	case "client_credentials":
		return tm.gc.LoginClient(ctx, tm.cfg.ClientID, tm.cfg.ClientSecret, tm.cfg.KeycloakRealm)
	default: // "password"
		return tm.gc.LoginAdmin(ctx, tm.cfg.AdminUser, tm.cfg.AdminPassword, tm.cfg.KeycloakRealm)
	}
}

// GoCloak returns the underlying gocloak client (used by the keycloak wrapper).
func (tm *TokenManager) GoCloak() *gocloak.GoCloak {
	return tm.gc
}
