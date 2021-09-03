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

	invite := &objects.Invite{}
	err = NewRequest().
		Method(http.MethodGet).
		Path(u.String()).
		Expect(http.StatusOK).
		ContentType(JsonContentType).
		Bind(invite).
		Send(c)

	return invite, err
}

func (c *Client) DeleteInvite(code, reason string) (*objects.Invite, error) {
	invite := &objects.Invite{}
	err := NewRequest().
		Method(http.MethodDelete).
		Path(fmt.Sprintf(InviteFmt, code)).
		Expect(http.StatusOK).
		ContentType(JsonContentType).
		Bind(invite).
		Send(c)

	return invite, err
}
