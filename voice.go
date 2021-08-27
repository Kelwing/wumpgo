package rest

import (
	"net/http"

	"github.com/Postcord/objects"
)

func (c *Client) GetVoiceRegions() ([]*objects.VoiceRegion, error) {
	res, err := c.request(&request{
		method:      http.MethodGet,
		path:        VoiceRegions,
		contentType: JsonContentType,
	})
	if err != nil {
		return nil, err
	}

	if err = res.ExpectsStatus(http.StatusOK); err != nil {
		return nil, err
	}

	var regions []*objects.VoiceRegion
	if err = res.JSON(&regions); err != nil {
		return nil, err
	}
	return regions, nil
}
