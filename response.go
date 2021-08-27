package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorREST struct {
	Message string
	Status  int
}

func (r *ErrorREST) Error() string {
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

func (r *DiscordResponse) ExpectsStatus(statusCode int) error {
	if r.StatusCode != statusCode {
		return &ErrorREST{
			Message: fmt.Sprintf("expected %d, got %d: %s", statusCode, r.StatusCode, r.Body),
			Status:  r.StatusCode,
		}
	}
	return nil
}

func (r *DiscordResponse) ExpectAnyStatus(statusCodes ...int) error {
	for _, j := range statusCodes {
		if j == r.StatusCode {
			return nil
		}
	}

	return &ErrorREST{
		Message: fmt.Sprintf("expected one of %d, got %d: %s", statusCodes, r.StatusCode, r.Body),
		Status:  r.StatusCode,
	}
}

type request struct {
	method      string
	path        string
	contentType string
	body        []byte
	reason      string

	omitAuth bool

	headers http.Header
}
