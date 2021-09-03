package rest

import (
	"net/http"

	"github.com/Postcord/objects"
)

func (c *Client) GetVoiceRegions() ([]*objects.VoiceRegion, error) {
	regions := []*objects.VoiceRegion{}
	err := NewRequest().
		Method(http.MethodGet).
		Path(VoiceRegions).
		ContentType(JsonContentType).
		Expect(http.StatusOK).
		Bind(regions).
		Send(c)

	return regions, err
}
