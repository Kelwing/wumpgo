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
	Body   []byte
	Status int
}

func (r *DiscordResponse) JSON(v interface{}) error {
	return json.Unmarshal(r.Body, v)
}

func (r *DiscordResponse) ExpectsStatus(statusCode int) error {
	if r.Status != statusCode {
		return &ErrorREST{
			Message: fmt.Sprintf("expected %d, got %d: %s", statusCode, r.Status, r.Body),
			Status:  r.Status,
		}
	}
	return nil
}

func (r *DiscordResponse) ExpectAnyStatus(statusCodes ...int) error {
	for _, j := range statusCodes {
		if j == r.Status {
			return nil
		}
	}

	return &ErrorREST{
		Message: fmt.Sprintf("expected one of %d, got %d: %s", statusCodes, r.Status, r.Body),
		Status:  r.Status,
	}
}

type request struct {
	method      string
	path        string
	contentType string
	body        []byte
	reason      string

	headers http.Header
}
