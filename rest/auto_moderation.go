package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"wumpgo.dev/wumpgo/objects"
)

func (c *Client) GetAutoModerationRules(ctx context.Context, guild objects.Snowflake) ([]*objects.AutoModerationRule, error) {
	rules := []*objects.AutoModerationRule{}

	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(AutoModerationRules, guild)).
		ContentType(JsonContentType).
		Bind(&rules).
		Send(c)

	return rules, err
}

func (c *Client) GetAutoModerationRule(ctx context.Context, guild objects.Snowflake, id objects.Snowflake) (*objects.AutoModerationRule, error) {
	rule := &objects.AutoModerationRule{}

	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(AutoModerationRule, guild, id)).
		ContentType(JsonContentType).
		Bind(rule).
		Send(c)

	return rule, err
}

type CreateAutoModerationRuleParams struct {
	Name            string                                 `json:"name"`
	EventType       objects.AutoModerationEventType        `json:"event_type"`
	TriggerType     objects.AutoModerationTriggerType      `json:"trigger_type"`
	TriggerMetadata *objects.AutoModerationTriggerMetadata `json:"trigger_metadata,omitempty"`
	Actions         []objects.AutoModerationAction         `json:"actions"`
	Enabled         bool                                   `json:"enabled"`
	ExemptRoles     []objects.Snowflake                    `json:"exempt_roles,omitempty"`
	ExemptChannels  []objects.Snowflake                    `json:"exempt_channels,omitempty"`
	Reason          string                                 `json:"-"`
}

func (c *Client) CreateAutoModerationRule(ctx context.Context, guild objects.Snowflake, params *CreateAutoModerationRuleParams) (*objects.AutoModerationRule, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	reason := ""
	if params != nil {
		reason = params.Reason
	}

	rule := &objects.AutoModerationRule{}

	err = NewRequest().
		Method(http.MethodPost).
		WithContext(ctx).
		Path(fmt.Sprintf(AutoModerationRules, guild)).
		ContentType(JsonContentType).
		Body(data).
		Reason(reason).
		Bind(rule).
		Send(c)

	return rule, err
}

type ModifyAutoModerationRuleParams struct {
	Name            string                                 `json:"name"`
	EventType       objects.AutoModerationEventType        `json:"event_type"`
	TriggerMetadata *objects.AutoModerationTriggerMetadata `json:"trigger_metadata,omitempty"`
	Actions         []objects.AutoModerationAction         `json:"actions"`
	Enabled         bool                                   `json:"enabled"`
	ExemptRoles     []objects.Snowflake                    `json:"exempt_roles,omitempty"`
	ExemptChannels  []objects.Snowflake                    `json:"exempt_channels,omitempty"`
	Reason          string                                 `json:"-"`
}

func (c *Client) ModifyAutoModerationRule(ctx context.Context, guild objects.Snowflake, id objects.Snowflake, params *ModifyAutoModerationRuleParams) (*objects.AutoModerationRule, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	reason := ""
	if params != nil {
		reason = params.Reason
	}

	rule := &objects.AutoModerationRule{}

	err = NewRequest().
		Method(http.MethodPatch).
		WithContext(ctx).
		Path(fmt.Sprintf(AutoModerationRule, guild, id)).
		ContentType(JsonContentType).
		Body(data).
		Reason(reason).
		Bind(rule).
		Send(c)

	return rule, err
}

func (c *Client) DeleteAutoModerationRule(ctx context.Context, guild objects.Snowflake, id objects.Snowflake, reason ...string) error {
	realReason := ""
	if len(reason) > 0 {
		realReason = reason[0]
	}

	rule := &objects.AutoModerationRule{}

	return NewRequest().
		Method(http.MethodDelete).
		WithContext(ctx).
		Path(fmt.Sprintf(AutoModerationRule, guild, id)).
		ContentType(JsonContentType).
		Reason(realReason).
		Bind(rule).
		Send(c)
}
