package routes

import (
	"net/http"
)

func About(version string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("microbin " + version))
	}
}
