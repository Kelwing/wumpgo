package rest

import (
	"encoding/json"
	"fmt"
	"github.com/Postcord/objects"
	"github.com/google/go-querystring/query"
	"net/http"
	"net/url"
)

func (c *Client) GetChannel(id objects.Snowflake) (*objects.Channel, error) {
	resp, err := c.request(http.MethodGet, fmt.Sprintf(ChannelBaseFmt, id), JsonContentType, nil)
	if err != nil {
		return nil, err
	}

	channel := &objects.Channel{}

	err = resp.JSON(channel)
	if err != nil {
		return nil, err
	}

	return channel, nil
}

type ModifyChannelParams struct {
	Name                 string                        `json:"name"`
	Type                 int64                         `json:"type"`
	Position             int64                         `json:"position,omitempty"`
	Topic                string                        `json:"topic,omitempty"`
	NSFW                 bool                          `json:"nsfw,omitempty"`
	RateLimitPerUser     int64                         `json:"rate_limit_per_user,omitempty"`
	Bitrate              int64                         `json:"bitrate,omitempty"`
	UserLimit            int64                         `json:"user_limit,omitempty"`
	PermissionOverwrites []objects.PermissionOverwrite `json:"permission_overwrites,omitempty"`
	Parent               objects.Snowflake             `json:"parent_id,omitempty"`
}

func (c *Client) ModifyChannel(id objects.Snowflake, params *ModifyChannelParams) (*objects.Channel, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	resp, err := c.request(http.MethodPatch, fmt.Sprintf(ChannelBaseFmt, id), JsonContentType, data)
	if err != nil {
		return nil, err
	}

	channel := &objects.Channel{}

	err = resp.JSON(channel)
	if err != nil {
		return nil, err
	}

	return channel, nil
}

func (c *Client) DeleteChannel(id objects.Snowflake) (*objects.Channel, error) {
	resp, err := c.request(http.MethodDelete, fmt.Sprintf(ChannelBaseFmt, id), JsonContentType, nil)
	if err != nil {
		return nil, err
	}

	channel := &objects.Channel{}

	err = resp.JSON(channel)
	if err != nil {
		return nil, err
	}

	return channel, nil
}

type GetChannelMessagesParams struct {
	Around objects.Snowflake `url:"around,omitempty"`
	Before objects.Snowflake `url:"before,omitempty"`
	After  objects.Snowflake `url:"after,omitempty"`
	Limit  int               `url:"limit,omitempty"`
}

func (c *Client) GetChannelMessages(id objects.Snowflake, params *GetChannelMessagesParams) ([]*objects.Message, error) {
	u, err := url.Parse(fmt.Sprintf(ChannelMessagesFmt, id))
	if err != nil {
		return nil, err
	}
	q, err := query.Values(params)
	if err != nil {
		return nil, err
	}
	u.RawQuery = q.Encode()

	res, err := c.request(http.MethodGet, u.String(), JsonContentType, nil)
	if err != nil {
		return nil, err
	}

	var messages []*objects.Message

	err = res.JSON(&messages)
	if err != nil {
		return nil, err
	}

	return messages, nil
}
