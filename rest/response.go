package rest

import (
	"encoding/json"
	"net/http"
)

type ErrorREST struct {
	Message string
	Status  int
	Body    json.RawMessage
}

func (r ErrorREST) Error() string {
	return r.Message
}

type DiscordResponse struct {
	Body       []byte
	StatusCode int
	Header     http.Header
}

func (r *DiscordResponse) JSON(v interface{}) error {
	return json.Unmarshal(r.Body, v)
}
