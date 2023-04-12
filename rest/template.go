package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"wumpgo.dev/wumpgo/objects"
)

func (c *Client) GetTemplate(ctx context.Context, code string) (*objects.Template, error) {
	template := &objects.Template{}
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(TemplateFmt, code)).
		ContentType(JsonContentType).
		Bind(template).
		Send(c)

	return template, err
}

func (c *Client) CreateGuildFromTemplate(ctx context.Context, code string, reason string) (*objects.Guild, error) {
	guild := &objects.Guild{}
	err := NewRequest().
		Method(http.MethodPost).
		WithContext(ctx).
		Path(fmt.Sprintf(TemplateFmt, code)).
		ContentType(JsonContentType).
		Bind(guild).
		Reason(reason).
		Send(c)
	return guild, err
}

func (c *Client) GetGuildTemplates(ctx context.Context, guild objects.Snowflake) ([]*objects.Template, error) {
	templates := []*objects.Template{}
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildTemplateFmt, guild)).
		ContentType(JsonContentType).
		Bind(&templates).
		Send(c)
	return templates, err
}

type CreateGuildTemplateParams struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Reason      string `json:"-"`
}

func (c *Client) CreateGuildTemplate(ctx context.Context, guild objects.Snowflake, params *CreateGuildTemplateParams) (*objects.Template, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	reason := ""
	if params != nil {
		reason = params.Reason
	}

	template := &objects.Template{}
	err = NewRequest().
		Method(http.MethodPost).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildTemplateFmt, guild)).
		ContentType(JsonContentType).
		Bind(template).
		Reason(reason).
		Body(data).
		Send(c)
	return template, err
}

func (c *Client) SyncGuildTemplate(ctx context.Context, guild objects.Snowflake, code string) (*objects.Template, error) {
	template := &objects.Template{}
	err := NewRequest().
		Method(http.MethodPut).
		Path(fmt.Sprintf(GuildTemplatesFmt, guild, code)).
		ContentType(JsonContentType).
		Bind(template).
		Send(c)
	return template, err
}

type ModifyGuildTemplateParams struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Reason      string `json:"-"`
}

func (c *Client) ModifyGuildTemplate(ctx context.Context, guild objects.Snowflake, code string, params *ModifyGuildTemplateParams) (*objects.Template, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	reason := ""
	if params != nil {
		reason = params.Reason
	}

	template := &objects.Template{}

	err = NewRequest().
		Method(http.MethodPatch).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildTemplatesFmt, guild, code)).
		ContentType(JsonContentType).
		Reason(reason).
		Body(data).
		Bind(template).
		Send(c)

	return template, err
}

func (c *Client) DeleteGuildTemplate(ctx context.Context, guild objects.Snowflake, code, reason string) (*objects.Template, error) {
	template := &objects.Template{}
	err := NewRequest().
		Method(http.MethodDelete).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildTemplatesFmt, guild, code)).
		ContentType(JsonContentType).
		Reason(reason).
		Bind(template).
		Send(c)

	return template, err
}
