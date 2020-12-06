package rest

import (
	"encoding/json"
	"fmt"
	"github.com/Postcord/objects"
	"net/http"
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
