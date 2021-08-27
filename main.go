package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sshilin/microbin/kubernetes"
)

var version = "unknown"

var opts struct {
	Listen string // LISTEN
	Tls    struct {
		Enabled  bool   // TLS_ENABLED
		KeyFile  string // TLS_KEY_FILE
		CertFile string // TLS_CERT_FILE
	}
}

func init() {
	opts.Listen = getenv("LISTEN", "0.0.0.0:8080")
	tlsEnabled, err := strconv.ParseBool(getenv("TLS_ENABLED", "false"))
	if err != nil {
		fmt.Println("Invalid value TLS_ENABLED")
	}
	opts.Tls.Enabled = tlsEnabled
	opts.Tls.KeyFile = getenv("TLS_KEY_FILE", "")
	opts.Tls.CertFile = getenv("TLS_CERT_FILE", "")
}

func main() {
	log.Printf("Starting microbin %s", version)

	if err := run(); err != http.ErrServerClosed {
		log.Fatalln(err)
	}
}

func run() error {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	server := &Server{
		router: router,
		k8s:    kubernetes.NewClient(),
	}
	server.routes()

	if opts.Tls.Enabled {
		return http.ListenAndServeTLS(opts.Listen, opts.Tls.CertFile, opts.Tls.KeyFile, server)
	} else {
		return http.ListenAndServe(opts.Listen, server)
	}
}

type Server struct {
	router *chi.Mux
	k8s    kubernetes.Client
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) routes() {
	s.router.Get("/", s.handleIndex(version))
	s.router.Get("/headers", s.handleHeaders())
}

func (srv *Server) handleIndex(version string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("microbin " + version))
	}
}

func getenv(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}
