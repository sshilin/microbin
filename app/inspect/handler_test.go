package inspect

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestInfoHandler(t *testing.T) {
	assert := assert.New(t)

	server := httptest.NewServer(Handler())
	defer server.Close()

	testCases := []struct {
		path   string
		method string
		header http.Header
	}{
		{
			path:   "/",
			method: http.MethodPost,
			header: http.Header{},
		},
		{
			path:   "/foo/bar",
			method: http.MethodDelete,
			header: http.Header{},
		},
		{
			path:   "/foo/bar?p1=1&p2=2",
			method: http.MethodGet,
			header: http.Header{},
		},
		{
			path:   "/foo/bar",
			method: http.MethodGet,
			header: http.Header{
				"Foo": []string{"bar"},
			},
		},
	}

	for _, tC := range testCases {
		t.Run("", func(t *testing.T) {
			r, err := http.NewRequest(tC.method, fmt.Sprintf("%s%s", server.URL, tC.path), nil)
			assert.NoError(err)

			r.Header = tC.header

			rr, err := server.Client().Do(r)
			assert.NoError(err)

			data, err := io.ReadAll(rr.Body)
			assert.NoError(err)

			var resp map[string]any

			json.Unmarshal(data, &resp)

			assert.NoError(err)
			assert.Contains(resp, "host")
			assert.Contains(resp, "remote")
			assert.Contains(resp, "proto")
			assert.Equal(tC.method, resp["method"])
			assert.Equal(tC.path, resp["uri"])

			headers, ok := (resp["headers"]).(map[string]any)
			assert.True(ok)

			for h := range tC.header {
				assert.Contains(headers, h)
			}
		})
	}
}
