package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorREST struct {
	Message string
	Status  int
	Body    json.RawMessage
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
		j := string(r.Body)
		return &ErrorREST{
			Message: fmt.Sprintf("expected %d, got %d (%s)", statusCode, r.StatusCode, j),
			Status:  r.StatusCode,
			Body:    r.Body,
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

	j := string(r.Body)
	return &ErrorREST{
		Message: fmt.Sprintf("expected %v, got %d (%s)", statusCodes, r.StatusCode, j),
		Status:  r.StatusCode,
		Body:    r.Body,
	}
}
