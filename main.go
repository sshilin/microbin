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
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sshilin/microbin/inspect"
	"github.com/sshilin/microbin/middleware"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

// build revision is set using ldflags
var build = "dev"

func httpHandler() http.Handler {
	router := chi.NewRouter()
	router.Use(
		middleware.Logger(&log.Logger),
		middleware.Metrics(),
	)
	router.NotFound(inspect.Handler())
	router.Handle("/metrics", promhttp.Handler())

	return router
}

func gracefullShutdown(sig os.Signal, server *http.Server) error {
	log.Info().Str("sig", sig.String()).Msg("shutdown server")
	if server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			return fmt.Errorf("shutdown server: %w", err)
		}
	}

	return nil
}

func run(cfg Config) error {
	handler := httpHandler()
	if !cfg.TLS.Enabled {
		handler = h2c.NewHandler(handler, &http2.Server{})
	}

	server := &http.Server{
		Addr:    cfg.Listen,
		Handler: handler,
	}

	serverError := make(chan error)

	if cfg.TLS.Enabled {
		go func() { serverError <- server.ListenAndServeTLS(cfg.TLS.CertFile, cfg.TLS.KeyFile) }()
	} else {
		go func() { serverError <- server.ListenAndServe() }()
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

	select {
	case err := <-serverError:
		return err
	case s := <-sig:
		return gracefullShutdown(s, server)
	}
}

func main() {
	cfg := LoadConfig()

	if !cfg.LogFormatJson {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	log.Info().Str("build", build).Interface("cfg", cfg).Send()

	err := run(cfg)
	if err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("http server error")
	}
}
