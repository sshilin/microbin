package inspect

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
)

type response struct {
	Host    string            `json:"host"`
	Remote  string            `json:"remote"`
	Proto   string            `json:"proto"`
	Method  string            `json:"method"`
	URI     string            `json:"uri"`
	Headers map[string]string `json:"headers,omitempty"`
}

func encode(values map[string][]string) map[string]string {
	h := make(map[string]string)

	for k, v := range values {
		h[k] = strings.Join(v, ",")
	}

	return h
}

func renderJson(w http.ResponseWriter, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	_, err = w.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := &response{
			Remote:  r.RemoteAddr,
			Proto:   r.Proto,
			Method:  r.Method,
			URI:     r.RequestURI,
			Headers: encode(r.Header),
		}

		if host, err := os.Hostname(); err == nil {
			resp.Host = host
		} else {
			log.Err(err).Msg("get hostname;dd")
		}

		if err := renderJson(w, resp); err != nil {
			log.Err(err).Msg("render json response")
			http.Error(w, "", http.StatusInternalServerError)
		}
	}
}
