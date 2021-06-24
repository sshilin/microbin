package routes

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/sshilin/microbin/k8s"
)

type Response struct {
	Pod     *k8s.Pod          `json:"pod,omitempty"`
	Headers map[string]string `json:"headers"`
}

func RequestHeaders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := &Response{Headers: formatHeaders(r.Header)}
		if pod := k8s.GetPodInfo(); pod.Name != "" {
			response.Pod = &pod
		}
		data, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Write(data)
	}
}

func formatHeaders(headers http.Header) map[string]string {
	h := map[string]string{}
	for k, v := range headers {
		h[k] = strings.Join(v, ",")
	}
	return h
}
