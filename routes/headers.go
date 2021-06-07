package routes

import (
	"encoding/json"
	"net/http"
	"strings"
)

func RequestHeaders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jsonHeaders := map[string]string{}
		for k, v := range r.Header {
			jsonHeaders[k] = strings.Join(v, ",")
		}
		data, err := json.MarshalIndent(jsonHeaders, "", "  ")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Write(data)
	}
}
