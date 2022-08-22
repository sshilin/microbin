package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	var buf bytes.Buffer

	log := zerolog.New(&buf)

	middleware := Logger(&log)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	f := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	middleware(f).ServeHTTP(w, r)

	fields := make(map[string]any)

	err := json.Unmarshal(buf.Bytes(), &fields)
	assert.NoError(t, err)

	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}

	expected := []string{"path", "status", "message", "level", "proto", "method", "size", "remoteAddr", "latency"}
	assert.ElementsMatch(t, expected, keys)
}
