package middleware

import (
	"net/http"
	"time"

	mid "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

func Logger(logger *zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := mid.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()

			defer func() {
				logger.Info().
					Str("method", r.Method).
					Str("path", r.URL.Path).
					Int("status", ww.Status()).
					Str("remoteAddr", r.RemoteAddr).
					Str("proto", r.Proto).
					Dur("latency", time.Since(t1)).
					Int("size", ww.BytesWritten()).
					Msg("served")
			}()

			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
