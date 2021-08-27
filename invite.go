package rest

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/Postcord/objects"
	"github.com/google/go-querystring/query"
)

type GetInviteParams struct {
	WithCounts bool `url:"with_counts"`
}

func (c *Client) GetInvite(code string, params *GetInviteParams) (*objects.Invite, error) {
	u, err := url.Parse(fmt.Sprintf(InviteFmt, code))
	if err != nil {
		return nil, err
	}

	v, err := query.Values(params)
	if err != nil {
		return nil, err
	}

	u.RawQuery = v.Encode()

	res, err := c.request(&request{
		method:      http.MethodGet,
		path:        u.String(),
		contentType: JsonContentType,
	})

	if err != nil {
		return nil, err
	}

	if err = res.ExpectsStatus(http.StatusOK); err != nil {
		return nil, err
	}

	inv := &objects.Invite{}
	if err = res.JSON(inv); err != nil {
		return nil, err
	}

	return inv, nil
}

func (c *Client) DeleteInvite(code, reason string) (*objects.Invite, error) {
	res, err := c.request(&request{
		method:      http.MethodDelete,
		path:        fmt.Sprintf(InviteFmt, code),
		contentType: JsonContentType,
		reason:      reason,
	})
	if err != nil {
		return nil, err
	}

	if err = res.ExpectsStatus(http.StatusOK); err != nil {
		return nil, err
	}

	inv := &objects.Invite{}
	if err = res.JSON(inv); err != nil {
		return nil, err
	}

	return inv, nil
}
