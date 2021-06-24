package routes

import (
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockDownwardAPI struct {
	podName      string
	podNamespace string
	nodeName     string
}

func (e mockDownwardAPI) PodName() string {
	return e.podName
}

func (e mockDownwardAPI) PodNamespace() string {
	return e.podNamespace
}

func (e mockDownwardAPI) NodeName() string {
	return e.nodeName
}

func TestHeaders(t *testing.T) {
	testCases := []struct {
		desc        string
		podInfo     mockDownwardAPI
		environment string
	}{
		{
			desc:        "Docker environment",
			environment: "docker",
			podInfo:     mockDownwardAPI{},
		},
		{
			desc:        "Kubernetes environment",
			environment: "kubernetes",
			podInfo: mockDownwardAPI{
				podName:      "pod",
				podNamespace: "namespace",
				nodeName:     "node",
			},
		},
	}
	headers := map[string]string{
		"Header1": "foo",
		"Header2": "bar",
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/headers", nil)
			for k, v := range headers {
				req.Header.Add(k, v)
			}
			rr := httptest.NewRecorder()
			RequestHeaders(tC.podInfo)(rr, req)
			data, err := ioutil.ReadAll(rr.Result().Body)
			assert.NoError(t, err)

			var response Response
			assert.NoError(t, json.Unmarshal(data, &response))

			assert.EqualValues(t, response.Headers, headers)

			switch tC.environment {
			case "docker":
				assert.Nil(t, response.Pod)
			case "kubernetes":
				assert.Equal(t, tC.podInfo.podName, response.Pod.Name)
				assert.Equal(t, tC.podInfo.podNamespace, response.Pod.Namespace)
				assert.Equal(t, tC.podInfo.nodeName, response.Pod.Node)
			}
		})
	}
}
