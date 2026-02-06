package config

import (
	"os"
	"time"
)

type Config struct {
	Transport          string
	Port               string
	KeycloakURL        string
	KeycloakRealm      string
	AuthMode           string // "password" or "client_credentials"
	AdminUser          string
	AdminPassword      string
	ClientID           string
	ClientSecret       string
	DefaultRealm       string
	TokenRefreshBuffer time.Duration
	LogLevel           string
	LogFormat          string
}

func Load() *Config {
	cfg := &Config{
		Transport:          envOr("TRANSPORT", "stdio"),
		Port:               envOr("PORT", "8080"),
		KeycloakURL:        envOr("KEYCLOAK_URL", "http://localhost:8080"),
		KeycloakRealm:      envOr("KEYCLOAK_REALM", "master"),
		AuthMode:           envOr("KEYCLOAK_AUTH_MODE", "password"),
		AdminUser:          envOr("KEYCLOAK_ADMIN_USER", "admin"),
		AdminPassword:      os.Getenv("KEYCLOAK_ADMIN_PASSWORD"),
		ClientID:           envOr("KEYCLOAK_CLIENT_ID", "mcp-admin"),
		ClientSecret:       os.Getenv("KEYCLOAK_CLIENT_SECRET"),
		DefaultRealm:       os.Getenv("KEYCLOAK_DEFAULT_REALM"),
		TokenRefreshBuffer: parseDuration(envOr("KEYCLOAK_TOKEN_REFRESH_BUFFER", "30s")),
		LogLevel:           envOr("LOG_LEVEL", "info"),
		LogFormat:          envOr("LOG_FORMAT", "json"),
	}
	return cfg
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 30 * time.Second
	}
	return d
}
