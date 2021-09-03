package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Postcord/objects"
	"github.com/google/go-querystring/query"
)

func (c *Client) GetCurrentUser() (*objects.User, error) {
	user := &objects.User{}
	err := NewRequest().
		Method(http.MethodGet).
		Path(UsersMeFmt).
		Expect(http.StatusOK).
		ContentType(JsonContentType).
		Bind(user).
		Send(c)

	return user, err
}

func (c *Client) GetUser(user objects.Snowflake) (*objects.User, error) {
	u := &objects.User{}
	err := NewRequest().
		Method(http.MethodGet).
		Path(fmt.Sprintf(UserFmt, user)).
		Expect(http.StatusOK).
		ContentType(JsonContentType).
		Bind(user).
		Send(c)

	return u, err
}

type ModifyCurrentUserParams struct {
	Name   string `json:"name,omitempty"`
	Avatar string `json:"avatar,omitempty"`
	Reason string `json:"-"`
}

func (c *Client) ModifyCurrentUser(params *ModifyCurrentUserParams) (*objects.User, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	reason := ""
	if params != nil {
		reason = params.Reason
	}

	u := &objects.User{}
	err = NewRequest().
		Method(http.MethodPatch).
		Path(UsersMeFmt).
		Expect(http.StatusOK).
		ContentType(JsonContentType).
		Body(data).
		Reason(reason).
		Bind(u).
		Send(c)

	return u, err
}

type CurrentUserGuildsParams struct {
	Before objects.Snowflake `json:"before,omitempty"`
	After  objects.Snowflake `json:"after,omitempty"`
	Limit  int               `json:"limit,omitempty"`
}

func (c *Client) GetCurrentUserGuilds(params *CurrentUserGuildsParams) ([]*objects.Guild, error) {
	u, err := url.Parse(UsersMeGuilds)
	if err != nil {
		return nil, err
	}

	v, err := query.Values(params)
	if err != nil {
		return nil, err
	}

	u.RawQuery = v.Encode()

	guilds := []*objects.Guild{}
	err = NewRequest().
		Method(http.MethodGet).
		Path(u.String()).
		Expect(http.StatusOK).
		ContentType(JsonContentType).
		Bind(&guilds).
		Send(c)
	return guilds, err
}

func (c *Client) LeaveGuild(guild objects.Snowflake) error {
	return NewRequest().
		Method(http.MethodDelete).
		Path(fmt.Sprintf(UsersMeGuild, guild)).
		ContentType(JsonContentType).
		Expect(http.StatusNoContent).
		Send(c)
}

type CreateDMParams struct {
	RecipientID objects.Snowflake `json:"recipient_id"`
	Reason      string            `json:"-"`
}

func (c *Client) CreateDM(params *CreateDMParams) (*objects.Channel, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	reason := ""
	if params != nil {
		reason = params.Reason
	}

	channel := &objects.Channel{}
	err = NewRequest().
		Method(http.MethodPost).
		Path(UsersMeChannels).
		Expect(http.StatusOK).
		ContentType(JsonContentType).
		Body(data).
		Reason(reason).
		Bind(channel).
		Send(c)

	return channel, err
}

type CreateGroupDMParams struct {
	AccessTokens []string                     `json:"access_tokens"`
	Nicks        map[objects.Snowflake]string `json:"nicks"`
	Reason       string                       `json:"reason,omitempty"`
}

func (c *Client) CreateGroupDM(params *CreateGroupDMParams) (*objects.Channel, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	reason := ""
	if params != nil {
		reason = params.Reason
	}

	channel := &objects.Channel{}
	err = NewRequest().
		Method(http.MethodPost).
		Path(UsersMeChannels).
		Expect(http.StatusOK).
		ContentType(JsonContentType).
		Body(data).
		Reason(reason).
		Bind(channel).
		Send(c)

	return channel, err
}

func (c *Client) GetUserConnections() ([]*objects.Connection, error) {
	connections := []*objects.Connection{}
	err := NewRequest().
		Method(http.MethodGet).
		Path(UserConnections).
		Expect(http.StatusOK).
		ContentType(JsonContentType).
		Bind(&connections).
		Send(c)

	return connections, err
}
