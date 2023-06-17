package middleware

import (
	"net/http"

	mid "github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func Metrics() func(next http.Handler) http.Handler {
	var httpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Count of all HTTP requests",
	}, []string{"code", "method", "path"})

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := mid.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)
			httpRequestsTotal.WithLabelValues(http.StatusText(ww.Status()), r.Method, r.URL.Path).Inc()
		}

		return http.HandlerFunc(fn)
	}
}
