package rest

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"wumpgo.dev/wumpgo/objects"
	"github.com/google/go-querystring/query"
)

type GetInviteParams struct {
	WithCounts bool `url:"with_counts"`
}

func (c *Client) GetInvite(ctx context.Context, code string, params *GetInviteParams) (*objects.Invite, error) {
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
		WithContext(ctx).
		Path(u.String()).
		ContentType(JsonContentType).
		Bind(invite).
		Send(c)

	return invite, err
}

func (c *Client) DeleteInvite(ctx context.Context, code, reason string) (*objects.Invite, error) {
	invite := &objects.Invite{}
	err := NewRequest().
		Method(http.MethodDelete).
		WithContext(ctx).
		Path(fmt.Sprintf(InviteFmt, code)).
		ContentType(JsonContentType).
		Bind(invite).
		Send(c)

	return invite, err
}
