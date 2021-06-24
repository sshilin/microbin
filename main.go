package main

import (
	"log"
	"net/http"
	"os"

	"github.com/sshilin/microbin/routes"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Version can be overridden via ldflags at the build time
var Version = "v.dev"

func main() {
	log.Printf("Starting microbin %s", Version)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", routes.About(Version))
	r.Get("/headers", routes.RequestHeaders())
	log.Fatal(http.ListenAndServe(getenv("IFACE", "")+":"+getenv("PORT", "8080"), r))
}

func getenv(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}
