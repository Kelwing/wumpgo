package rest

import "encoding/json"

type DiscordResponse struct {
	Body   []byte
	Status int
}

func (r *DiscordResponse) JSON(v interface{}) error {
	return json.Unmarshal(r.Body, v)
}
