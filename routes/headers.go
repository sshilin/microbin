package routes

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/sshilin/microbin/k8s"
)

type Response struct {
	Instance k8s.Info
	Headers  map[string]string
}

func RequestHeaders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jsonHeaders := map[string]string{}
		for k, v := range r.Header {
			jsonHeaders[k] = strings.Join(v, ",")
		}
		data, err := json.MarshalIndent(Response{
			Instance: k8s.GetSelfInfo(),
			Headers:  jsonHeaders,
		}, "", "  ")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Write(data)
	}
}
