package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Postcord/objects"
	"github.com/google/go-querystring/query"
)

func (c *Client) GetCurrentUser(ctx context.Context) (*objects.User, error) {
	user := &objects.User{}
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(UsersMeFmt).
		ContentType(JsonContentType).
		Bind(user).
		Send(c)

	return user, err
}

func (c *Client) GetUser(ctx context.Context, user objects.SnowflakeObject) (*objects.User, error) {
	u := &objects.User{}
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(UserFmt, user.GetID())).
		ContentType(JsonContentType).
		Bind(u).
		Send(c)

	return u, err
}

type ModifyCurrentUserParams struct {
	Name   string `json:"name,omitempty"`
	Avatar string `json:"avatar,omitempty"`
	Reason string `json:"-"`
}

func (c *Client) ModifyCurrentUser(ctx context.Context, params *ModifyCurrentUserParams) (*objects.User, error) {
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
		WithContext(ctx).
		Path(UsersMeFmt).
		ContentType(JsonContentType).
		Body(data).
		Reason(reason).
		Bind(u).
		Send(c)

	return u, err
}

type CurrentUserGuildsParams struct {
	Before objects.Snowflake `url:"before,omitempty"`
	After  objects.Snowflake `url:"after,omitempty"`
	Limit  int               `url:"limit,omitempty"`
}

func (c *Client) GetCurrentUserGuilds(ctx context.Context, params *CurrentUserGuildsParams) ([]*objects.Guild, error) {
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
		WithContext(ctx).
		Path(u.String()).
		ContentType(JsonContentType).
		Bind(&guilds).
		Send(c)
	return guilds, err
}

func (c *Client) GetCurrentUserGuildMember(ctx context.Context, guild objects.SnowflakeObject) (*objects.GuildMember, error) {
	member := &objects.GuildMember{}
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(UsersMeGuildMember, guild.GetID())).
		ContentType(JsonContentType).
		Bind(member).
		Send(c)

	return member, err
}

func (c *Client) LeaveGuild(ctx context.Context, guild objects.SnowflakeObject) error {
	return NewRequest().
		Method(http.MethodDelete).
		WithContext(ctx).
		Path(fmt.Sprintf(UsersMeGuild, guild.GetID())).
		ContentType(JsonContentType).
		Send(c)
}

type CreateDMParams struct {
	RecipientID objects.Snowflake `json:"recipient_id"`
	Reason      string            `json:"-"`
}

func (c *Client) CreateDM(ctx context.Context, params *CreateDMParams) (*objects.Channel, error) {
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
		WithContext(ctx).
		Path(UsersMeChannels).
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

func (c *Client) CreateGroupDM(ctx context.Context, params *CreateGroupDMParams) (*objects.Channel, error) {
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
		WithContext(ctx).
		Path(UsersMeChannels).
		ContentType(JsonContentType).
		Body(data).
		Reason(reason).
		Bind(channel).
		Send(c)

	return channel, err
}

func (c *Client) GetUserConnections(ctx context.Context) ([]*objects.Connection, error) {
	connections := []*objects.Connection{}
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(UserConnections).
		ContentType(JsonContentType).
		Bind(&connections).
		Send(c)

	return connections, err
}
