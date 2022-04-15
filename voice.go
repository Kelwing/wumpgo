package rest

import (
	"context"
	"net/http"

	"github.com/Postcord/objects"
)

func (c *Client) GetVoiceRegions(ctx context.Context) ([]*objects.VoiceRegion, error) {
	regions := []*objects.VoiceRegion{}
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(VoiceRegions).
		ContentType(JsonContentType).
		Bind(&regions).
		Send(c)

	return regions, err
}
