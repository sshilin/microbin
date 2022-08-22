package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"github.com/sshilin/microbin/app/headers"
	"github.com/sshilin/microbin/app/middleware"
)

// build revision to be set via ldflags
var build = "develop"

func main() {
	log.Info().Str("build", build).Msg("starting")

	if err := run(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("server error")
	}
}

func run() error {
	cfg := struct {
		Addr string
		TLS  struct {
			Enabled  bool
			KeyFile  string
			CertFile string
		}
	}{
		Addr: getenv("ADDRESS", "0.0.0.0:8080"),
		TLS: struct {
			Enabled  bool
			KeyFile  string
			CertFile string
		}{
			Enabled:  getenv("TLS_ENABLED", "false") == "true",
			KeyFile:  getenv("TLS_KEY_FILE", ""),
			CertFile: getenv("TLS_CERT_FILE", ""),
		},
	}

	log.Info().Msgf("config: %+v", cfg)

	r := chi.NewRouter()
	r.Use(middleware.Logger(&log.Logger))
	r.Use(middleware.Metrics())
	r.Get("/headers", headers.Handler())
	r.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:    cfg.Addr,
		Handler: r,
	}

	serverErrors := make(chan error, 1)

	go func() {
		if cfg.TLS.Enabled {
			serverErrors <- server.ListenAndServeTLS(cfg.TLS.CertFile, cfg.TLS.KeyFile)
		} else {
			serverErrors <- server.ListenAndServe()
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

	select {
	case err := <-serverErrors:
		return err
	case s := <-sig:
		if server != nil {
			log.Info().
				Str("sig", s.String()).
				Msg("shutdown started")
			defer log.Info().
				Str("sig", s.String()).
				Msg("shutdown complete")

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := server.Shutdown(ctx); err != nil {
				return fmt.Errorf("could not stop server gracefully: %w", err)
			}
		}
	}

	return nil
}

func getenv(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return defaultValue
}
