package routes

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/sshilin/microbin/downward"
)

type Pod struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Node      string `json:"node"`
}

type Response struct {
	Pod     *Pod              `json:"pod,omitempty"`
	Headers map[string]string `json:"headers"`
}

func RequestHeaders(dw downward.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := &Response{Headers: formatHeaders(r.Header)}
		response.Pod = getPodInfo(dw)

		data, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Write(data)
	}
}

func getPodInfo(dw downward.API) *Pod {
	podName := dw.PodName()
	podNamespace := dw.PodNamespace()
	nodeName := dw.NodeName()
	if podName == "" && podNamespace == "" && nodeName == "" {
		return nil
	}
	return &Pod{
		Name:      podName,
		Namespace: podNamespace,
		Node:      nodeName,
	}
}

func formatHeaders(headers http.Header) map[string]string {
	h := map[string]string{}
	for k, v := range headers {
		h[k] = strings.Join(v, ",")
	}
	return h
}
