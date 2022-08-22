package headers

import "encoding/json"

type tokenResponse struct {
	Header    json.RawMessage `json:"header"`
	Payload   json.RawMessage `json:"payload"`
	Signature string          `json:"signature,omitempty"`
}

type podResponse struct {
	Namespace string `json:"namespace"`
	Node      string `json:"node"`
	Name      string `json:"name"`
}

type headersResponse struct {
	Pod     *podResponse      `json:"pod,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
	Bearer  *tokenResponse    `json:"bearer,omitempty"`
}
