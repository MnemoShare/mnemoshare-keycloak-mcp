package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/mnemoshare/mnemoshare-keycloak-mcp/internal/auth"
	"github.com/mnemoshare/mnemoshare-keycloak-mcp/internal/config"
	"github.com/mnemoshare/mnemoshare-keycloak-mcp/internal/keycloak"
	"github.com/mnemoshare/mnemoshare-keycloak-mcp/internal/tools"
)

var version = "dev"

func main() {
	cfg := config.Load()
	initLogger(cfg)

	log.Info().
		Str("transport", cfg.Transport).
		Str("keycloak_url", cfg.KeycloakURL).
		Str("auth_mode", cfg.AuthMode).
		Msg("starting keycloak-mcp server")

	// Token manager + keycloak client
	tm := auth.NewTokenManager(cfg)
	kc := keycloak.NewClient(cfg, tm)

	// MCP server
	s := mcp.NewServer(
		&mcp.Implementation{Name: "keycloak-mcp", Version: version},
		nil,
	)

	tools.RegisterAll(s, kc)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	switch cfg.Transport {
	case "http":
		runHTTP(ctx, cfg, s)
	default:
		runStdio(ctx, s)
	}
}

func runStdio(ctx context.Context, s *mcp.Server) {
	log.Info().Msg("running in stdio mode")
	if err := s.Run(ctx, &mcp.StdioTransport{}); err != nil {
		log.Fatal().Err(err).Msg("stdio server error")
	}
}

func runHTTP(ctx context.Context, cfg *config.Config, s *mcp.Server) {
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Info().Str("addr", addr).Msg("running in HTTP mode")

	httpHandler := mcp.NewStreamableHTTPHandler(
		func(r *http.Request) *mcp.Server { return s },
		nil,
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "ok")
	})
	mux.Handle("/mcp", httpHandler)

	srv := &http.Server{Addr: addr, Handler: mux}

	go func() {
		<-ctx.Done()
		log.Info().Msg("shutting down HTTP server")
		srv.Close()
	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("HTTP server error")
	}
}

func initLogger(cfg *config.Config) {
	level, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	if cfg.LogFormat == "console" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}
