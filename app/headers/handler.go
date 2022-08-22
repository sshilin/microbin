package headers

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
)

func Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resp headersResponse

		pod := encodeWorkloadInfo()

		if pod.Namespace != "" || pod.Node != "" || pod.Name != "" {
			resp.Pod = pod
		}

		headers := encodeHeaders(r.Header)

		if authHeader, ok := headers["Authorization"]; ok {
			token, err := encodeToken(authHeader)
			if err != nil {
				// don't fail request, just skip bearer part
				log.Err(err).Msg("error parsing authorization header")
			} else {
				resp.Bearer = token
			}
		}

		resp.Headers = headers

		if err := renderJson(w, resp); err != nil {
			log.Err(err).Msg("error rendering json response")
			http.Error(w, "", http.StatusInternalServerError)
		}
	}
}

func encodeWorkloadInfo() *podResponse {
	return &podResponse{
		Namespace: os.Getenv("K8S_POD_NAMESPACE"),
		Node:      os.Getenv("K8S_NODE_NAME"),
		Name:      os.Getenv("K8S_POD_NAME"),
	}
}

func encodeHeaders(headers http.Header) map[string]string {
	h := make(map[string]string)

	for k, v := range headers {
		h[k] = strings.Join(v, ",")
	}

	return h
}

func encodeToken(authHeader string) (*tokenResponse, error) {
	fields := strings.Fields(authHeader)
	if len(fields) != 2 || fields[0] != "Bearer" {
		return nil, errors.New("unexpected header format")
	}

	parts := strings.Split(fields[1], ".")
	if len(parts) < 2 || len(parts) > 3 {
		return nil, errors.New("unexpected token format")
	}

	decoded := make([][]byte, 0, 3)

	for _, part := range parts {
		data, err := base64.RawURLEncoding.DecodeString(part)
		if err != nil {
			return nil, err
		}

		decoded = append(decoded, data)
	}

	var header json.RawMessage

	err := json.Unmarshal(decoded[0], &header)
	if err != nil {
		return nil, err
	}

	var payload json.RawMessage

	err = json.Unmarshal(decoded[1], &payload)
	if err != nil {
		return nil, err
	}

	var signature string

	if len(decoded) == 3 {
		signature = hex.EncodeToString(decoded[2])
	}

	return &tokenResponse{
		Header:    header,
		Payload:   payload,
		Signature: signature,
	}, nil
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
