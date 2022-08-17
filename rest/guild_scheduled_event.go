package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Kelwing/wumpgo/objects"
	"github.com/google/go-querystring/query"
)

type CreateGuildScheduledEventParams struct {
	ChannelID          *objects.Snowflake                         `json:"channel_id,omitempty"`
	EntityMetadata     *objects.GuildScheduledEventEntityMetadata `json:"entity_metadata,omitempty"`
	Name               string                                     `json:"name"`
	PrivacyLevel       objects.PrivacyLevel                       `json:"privacy_level"`
	ScheduledStartTime objects.Time                               `json:"scheduled_start_time"`
	ScheduledEndTime   *objects.Time                              `json:"scheduled_end_time,omitempty"`
	Description        string                                     `json:"description,omitempty"`
	EntityType         objects.GuildScheduledEventEntityType      `json:"entity_type"`
	Image              []byte                                     `json:"image,omitempty"`
	Reason             string                                     `json:"-"`
}

type GetGuildScheduledEventParams struct {
	WithUserCount bool `url:"with_user_count,omitempty"`
}

type ModifyGuildScheduledEventParams struct {
	CreateGuildScheduledEventParams
	Status objects.GuildScheduledEventStatus `json:"status,omitempty"`
}

type GetGuildScheduledEventUsersParams struct {
	Limit      int                `url:"limit,omitempty"`
	WithMember bool               `url:"with_member,omitempty"`
	Before     *objects.Snowflake `url:"before,omitempty"`
	After      *objects.Snowflake `url:"after,omitempty"`
}

func (c *Client) CreateGuildScheduledEvent(ctx context.Context, guildID objects.SnowflakeObject, params *CreateGuildScheduledEventParams) (*objects.GuildScheduledEvent, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	reason := ""
	if params != nil {
		reason = params.Reason
	}

	event := &objects.GuildScheduledEvent{}

	err = NewRequest().
		Method(http.MethodPost).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildScheduledEventBaseFmt, guildID.GetID())).
		ContentType(JsonContentType).
		Body(data).
		Reason(reason).
		Bind(event).
		Send(c)

	return event, err
}

func (c *Client) GetGuildScheduledEvent(ctx context.Context, guildID objects.SnowflakeObject, id objects.SnowflakeObject, params ...*GetGuildScheduledEventParams) (*objects.GuildScheduledEvent, error) {
	u, err := url.Parse(fmt.Sprintf(GuildScheduledEventFmt, guildID.GetID(), id.GetID()))
	if err != nil {
		return nil, err
	}

	if len(params) > 0 {
		v, err := query.Values(params[0])
		if err != nil {
			return nil, err
		}
		u.RawQuery = v.Encode()
	}

	event := &objects.GuildScheduledEvent{}

	err = NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(u.String()).
		ContentType(JsonContentType).
		Bind(event).
		Send(c)

	return event, err
}

func (c *Client) ModifyGuildScheduledEvent(ctx context.Context, guildID objects.SnowflakeObject, id objects.SnowflakeObject, params *ModifyGuildScheduledEventParams) (*objects.GuildScheduledEvent, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	reason := ""
	if params != nil {
		reason = params.Reason
	}

	event := &objects.GuildScheduledEvent{}

	err = NewRequest().
		Method(http.MethodPatch).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildScheduledEventFmt, guildID.GetID(), id.GetID())).
		ContentType(JsonContentType).
		Body(data).
		Reason(reason).
		Bind(event).
		Send(c)

	return event, err
}

func (c *Client) DeleteGuildScheduledEvent(ctx context.Context, guildID objects.SnowflakeObject, id objects.SnowflakeObject) error {
	return NewRequest().
		Method(http.MethodDelete).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildScheduledEventFmt, guildID.GetID(), id.GetID())).
		ContentType(JsonContentType).
		Send(c)
}

func (c *Client) GetGuildScheduledEventUsers(ctx context.Context, guildID objects.SnowflakeObject, id objects.SnowflakeObject, params ...*GetGuildScheduledEventUsersParams) ([]*objects.GuildScheduledEventUser, error) {
	u, err := url.Parse(fmt.Sprintf(GuildScheduledEventUsersFmt, guildID.GetID(), id.GetID()))
	if err != nil {
		return nil, err
	}

	if len(params) > 0 {
		v, err := query.Values(params[0])
		if err != nil {
			return nil, err
		}
		u.RawQuery = v.Encode()
	}

	users := make([]*objects.GuildScheduledEventUser, 0)

	err = NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(u.String()).
		ContentType(JsonContentType).
		Bind(users).
		Send(c)

	return users, err
}
