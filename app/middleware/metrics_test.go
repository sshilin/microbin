package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestMetrics(t *testing.T) {
	expect := `
	# HELP http_requests_total Count of all HTTP requests
	# TYPE http_requests_total counter
	
	http_requests_total{code="",method="GET",path="/"} 2
	`

	middleware := Metrics()

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	f := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	middleware(f).ServeHTTP(w, r)
	middleware(f).ServeHTTP(w, r)

	err := testutil.GatherAndCompare(prometheus.DefaultGatherer, strings.NewReader(expect), "http_requests_total")
	assert.NoError(t, err)
}
