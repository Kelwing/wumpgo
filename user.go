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
	res, err := c.request(&request{
		method:      http.MethodGet,
		path:        UsersMeFmt,
		contentType: JsonContentType,
	})
	if err != nil {
		return nil, err
	}

	if err = res.ExpectsStatus(http.StatusOK); err != nil {
		return nil, err
	}

	me := &objects.User{}
	if err = res.JSON(me); err != nil {
		return nil, err
	}

	return me, nil
}

func (c *Client) GetUser(user objects.Snowflake) (*objects.User, error) {
	res, err := c.request(&request{
		method:      http.MethodGet,
		path:        fmt.Sprintf(UserFmt, user),
		contentType: JsonContentType,
	})
	if err != nil {
		return nil, err
	}

	if err = res.ExpectsStatus(http.StatusOK); err != nil {
		return nil, err
	}

	u := &objects.User{}
	if err = res.JSON(u); err != nil {
		return nil, err
	}

	return u, nil
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

	res, err := c.request(&request{
		method:      http.MethodPatch,
		path:        UsersMeFmt,
		contentType: JsonContentType,
		body:        data,
		reason:      reason,
	})
	if err != nil {
		return nil, err
	}

	if err = res.ExpectsStatus(http.StatusOK); err != nil {
		return nil, err
	}

	u := &objects.User{}
	if err = res.JSON(u); err != nil {
		return nil, err
	}

	return u, nil
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

	var guilds []*objects.Guild
	if err = res.JSON(&guilds); err != nil {
		return nil, err
	}
	return guilds, nil
}

func (c *Client) LeaveGuild(guild objects.Snowflake) error {
	res, err := c.request(&request{
		method:      http.MethodDelete,
		path:        fmt.Sprintf(UsersMeGuild, guild),
		contentType: JsonContentType,
	})
	if err != nil {
		return err
	}

	if err = res.ExpectsStatus(http.StatusNoContent); err != nil {
		return err
	}

	return nil
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

	res, err := c.request(&request{
		method:      http.MethodPost,
		path:        UsersMeChannels,
		contentType: JsonContentType,
		body:        data,
		reason:      reason,
	})
	if err != nil {
		return nil, err
	}

	if err = res.ExpectsStatus(http.StatusOK); err != nil {
		return nil, err
	}

	channel := &objects.Channel{}
	if err = res.JSON(channel); err != nil {
		return nil, err
	}

	return channel, nil
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

	res, err := c.request(&request{
		method:      http.MethodPost,
		path:        UsersMeChannels,
		contentType: JsonContentType,
		body:        data,
		reason:      reason,
	})

	if err != nil {
		return nil, err
	}

	if err = res.ExpectsStatus(http.StatusOK); err != nil {
		return nil, err
	}

	channel := &objects.Channel{}
	if err = res.JSON(channel); err != nil {
		return nil, err
	}

	return channel, nil
}

func (c *Client) GetUserConnections() ([]*objects.Connection, error) {
	res, err := c.request(&request{
		method:      http.MethodGet,
		path:        UserConnections,
		contentType: JsonContentType,
	})
	if err != nil {
		return nil, err
	}

	if err = res.ExpectsStatus(http.StatusOK); err != nil {
		return nil, err
	}

	var conn []*objects.Connection
	if err = res.JSON(&conn); err != nil {
		return nil, err
	}

	return conn, nil
}
