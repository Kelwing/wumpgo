package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Postcord/objects"
)

func (c *Client) GetTemplate(code string) (*objects.Template, error) {
	template := &objects.Template{}
	err := NewRequest().
		Method(http.MethodGet).
		Path(fmt.Sprintf(TemplateFmt, code)).
		Expect(http.StatusOK).
		ContentType(JsonContentType).
		Bind(template).
		Send(c)

	return template, err
}

func (c *Client) CreateGuildFromTemplate(code string, reason string) (*objects.Guild, error) {
	guild := &objects.Guild{}
	err := NewRequest().
		Method(http.MethodPost).
		Path(fmt.Sprintf(TemplateFmt, code)).
		Expect(http.StatusOK).
		ContentType(JsonContentType).
		Bind(guild).
		Reason(reason).
		Send(c)
	return guild, err
}

func (c *Client) GetGuildTemplates(guild objects.Snowflake) ([]*objects.Template, error) {
	templates := []*objects.Template{}
	err := NewRequest().
		Method(http.MethodGet).
		Path(fmt.Sprintf(GuildTemplateFmt, guild)).
		Expect(http.StatusOK).
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

func (c *Client) CreateGuildTemplate(guild objects.Snowflake, params *CreateGuildTemplateParams) (*objects.Template, error) {
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
		Path(fmt.Sprintf(GuildTemplateFmt, guild)).
		Expect(http.StatusOK).
		ContentType(JsonContentType).
		Bind(template).
		Reason(reason).
		Body(data).
		Send(c)
	return template, err
}

func (c *Client) SyncGuildTemplate(guild objects.Snowflake, code string) (*objects.Template, error) {
	template := &objects.Template{}
	err := NewRequest().
		Method(http.MethodPut).
		Path(fmt.Sprintf(GuildTemplatesFmt, guild, code)).
		Expect(http.StatusOK).
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

func (c *Client) ModifyGuildTemplate(guild objects.Snowflake, code string, params *ModifyGuildTemplateParams) (*objects.Template, error) {
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
		Path(fmt.Sprintf(GuildTemplatesFmt, guild, code)).
		Expect(http.StatusOK).
		ContentType(JsonContentType).
		Reason(reason).
		Body(data).
		Bind(template).
		Send(c)

	return template, err
}

func (c *Client) DeleteGuildTemplate(guild objects.Snowflake, code, reason string) (*objects.Template, error) {
	template := &objects.Template{}
	err := NewRequest().
		Method(http.MethodDelete).
		Path(fmt.Sprintf(GuildTemplatesFmt, guild, code)).
		Expect(http.StatusOK).
		ContentType(JsonContentType).
		Reason(reason).
		Bind(template).
		Send(c)

	return template, err
}
