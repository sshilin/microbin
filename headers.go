package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (srv *Server) handleHeaders() http.HandlerFunc {
	type podInfo struct {
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
		Node      string `json:"node"`
	}

	type response struct {
		Protocol string            `json:"proto"`
		Headers  map[string]string `json:"headers"`
		PodInfo  *podInfo          `json:"pod,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		response := &response{
			Protocol: r.Proto,
			Headers:  formatHeaders(r.Header),
		}

		pod := podInfo{
			Name:      srv.k8s.PodName(),
			Namespace: srv.k8s.PodNamespace(),
			Node:      srv.k8s.NodeName(),
		}
		if pod != (podInfo{}) {
			response.PodInfo = &pod
		}

		data, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
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
