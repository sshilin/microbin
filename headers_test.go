package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type k8sClientMock struct {
	podName      string
	podNamespace string
	nodeName     string
}

func (e k8sClientMock) PodName() string {
	return e.podName
}

func (e k8sClientMock) PodNamespace() string {
	return e.podNamespace
}

func (e k8sClientMock) NodeName() string {
	return e.nodeName
}

func TestHeaders(t *testing.T) {
	testCases := []struct {
		desc        string
		podInfo     k8sClientMock
		environment string
	}{
		{
			desc:        "Docker environment",
			environment: "docker",
			podInfo:     k8sClientMock{},
		},
		{
			desc:        "Kubernetes environment",
			environment: "kubernetes",
			podInfo: k8sClientMock{
				podName:      "pod",
				podNamespace: "namespace",
				nodeName:     "node",
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/headers", nil)
			headers := map[string]string{
				"Header1": "foo",
				"Header2": "bar",
			}
			srv := Server{
				router: nil,
				k8s:    &tC.podInfo,
			}
			for k, v := range headers {
				req.Header.Add(k, v)
			}
			rr := httptest.NewRecorder()
			srv.handleHeaders().ServeHTTP(rr, req)
			data, err := ioutil.ReadAll(rr.Result().Body)
			assert.NoError(t, err)

			var response map[string]interface{}
			assert.NoError(t, json.Unmarshal(data, &response))

			switch tC.environment {
			case "docker":
				assert.Nil(t, response["pod"])
				assert.NotNil(t, response["headers"])
			case "kubernetes":
				assert.NotNil(t, response["pod"])
				assert.NotNil(t, response["headers"])
			}
		})
	}
}
