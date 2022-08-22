package headers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWorkloadResponse(t *testing.T) {
	ns, node, pod := "default", "docker-desktop", "microbin-6d4b75cb6d-5vws7"
	os.Setenv("K8S_POD_NAMESPACE", ns)
	os.Setenv("K8S_NODE_NAME", node)
	os.Setenv("K8S_POD_NAME", pod)

	resp := encodeWorkloadInfo()

	assert.Equal(t, ns, resp.Namespace)
	assert.Equal(t, node, resp.Node)
	assert.Equal(t, pod, resp.Name)

	defer func() {
		os.Unsetenv("K8S_POD_NAMESPACE")
		os.Unsetenv("K8S_NODE_NAME")
		os.Unsetenv("K8S_POD_NAME")
	}()
}

func TestNewTokenResponsePositive(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		desc   string
		auth   string
		signed bool
	}{
		{
			desc:   "Good signed token",
			auth:   "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			signed: true,
		},
		{
			desc:   "Good none alg token",
			auth:   "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJub25lIn0.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTY2MDk3NTg1MCwiZXhwIjoxNjYwOTc5NDUwfQ",
			signed: false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			token, err := encodeToken(tC.auth)
			assert.NoError(err)
			assert.NotEmpty(token.Header)
			assert.NotEmpty(token.Payload)
			if tC.signed {
				assert.NotEmpty(token.Signature)
			} else {
				assert.Empty(token.Signature)
			}
		})
	}
}

func TestNewTokenResponseNegative(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		desc string
		auth string
	}{
		{
			desc: "Invalid header format",
			auth: "unexpected",
		},
		{
			desc: "Invalid token format",
			auth: "Bearer unexpected",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			token, err := encodeToken(tC.auth)
			assert.Error(err)
			assert.Nil(token)
		})
	}
}

func TestHandlerHeadersMinimal(t *testing.T) {
	assert := assert.New(t)

	r, _ := http.NewRequest("GET", "/headers", nil)
	w := httptest.NewRecorder()
	Handler().ServeHTTP(w, r)

	assert.Empty(unmarshalToAny(w.Body))
	assert.Equal(w.Result().StatusCode, http.StatusOK)
}

func TestHandlerHeadersFull(t *testing.T) {
	assert := assert.New(t)

	ns, node, pod := "default", "docker-desktop", "microbin-6d4b75cb6d-5vws7"

	os.Setenv("K8S_POD_NAMESPACE", ns)
	os.Setenv("K8S_NODE_NAME", node)
	os.Setenv("K8S_POD_NAME", pod)

	defer func() {
		os.Unsetenv("K8S_POD_NAMESPACE")
		os.Unsetenv("K8S_NODE_NAME")
		os.Unsetenv("K8S_POD_NAME")
	}()

	r, _ := http.NewRequest("GET", "/headers", nil)
	r.Header = http.Header{
		"foo":           {"bar"},
		"Authorization": {"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"},
	}

	w := httptest.NewRecorder()
	Handler().ServeHTTP(w, r)

	resp, err := unmarshalToAny(w.Body)
	assert.NoError(err)

	assert.Contains(resp, "pod")
	assert.Contains(resp, "headers")
	assert.Contains(resp, "bearer")

	assert.Equal(w.Result().StatusCode, http.StatusOK)
}

func unmarshalToAny(buf *bytes.Buffer) (map[string]any, error) {
	data, _ := io.ReadAll(buf)
	var resp map[string]any

	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}
