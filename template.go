package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Postcord/objects"
)

func (c *Client) GetTemplate(code string) (*objects.Template, error) {
	res, err := c.request(&request{
		method:      http.MethodGet,
		path:        fmt.Sprintf(TemplateFmt, code),
		contentType: JsonContentType,
	})
	if err != nil {
		return nil, err
	}

	if err = res.ExpectsStatus(http.StatusOK); err != nil {
		return nil, err
	}

	template := &objects.Template{}
	if err = res.JSON(template); err != nil {
		return nil, err
	}

	return template, nil
}

func (c *Client) CreateGuildFromTemplate(code string, reason string) (*objects.Guild, error) {
	res, err := c.request(&request{
		method:      http.MethodPost,
		path:        fmt.Sprintf(TemplateFmt, code),
		contentType: JsonContentType,
		reason:      reason,
	})
	if err != nil {
		return nil, err
	}

	if err = res.ExpectsStatus(http.StatusOK); err != nil {
		return nil, err
	}

	guild := &objects.Guild{}
	if err = res.JSON(guild); err != nil {
		return nil, err
	}
	return guild, nil
}

func (c *Client) GetGuildTemplates(guild objects.Snowflake) ([]*objects.Template, error) {
	res, err := c.request(&request{
		method:      http.MethodGet,
		path:        fmt.Sprintf(GuildTemplateFmt, guild),
		contentType: JsonContentType,
	})
	if err != nil {
		return nil, err
	}

	var templates []*objects.Template
	if err = res.JSON(&templates); err != nil {
		return nil, err
	}

	return templates, nil
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

	res, err := c.request(&request{
		method:      http.MethodPost,
		path:        fmt.Sprintf(GuildTemplateFmt, guild),
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

	template := &objects.Template{}
	if err = res.JSON(template); err != nil {
		return nil, err
	}

	return template, nil
}

func (c *Client) SyncGuildTemplate(guild objects.Snowflake, code string) (*objects.Template, error) {
	res, err := c.request(&request{
		method:      http.MethodPut,
		path:        fmt.Sprintf(GuildTemplatesFmt, guild, code),
		contentType: JsonContentType,
	})
	if err != nil {
		return nil, err
	}

	if err = res.ExpectsStatus(http.StatusOK); err != nil {
		return nil, err
	}

	template := &objects.Template{}
	if err = res.JSON(template); err != nil {
		return nil, err
	}

	return template, nil
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

	res, err := c.request(&request{
		method:      http.MethodPatch,
		path:        fmt.Sprintf(GuildTemplatesFmt, guild, code),
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

	template := &objects.Template{}
	if err = res.JSON(template); err != nil {
		return nil, err
	}
	return template, nil
}

func (c *Client) DeleteGuildTemplate(guild objects.Snowflake, code, reason string) (*objects.Template, error) {
	res, err := c.request(&request{
		method:      http.MethodDelete,
		path:        fmt.Sprintf(GuildTemplatesFmt, guild, code),
		contentType: JsonContentType,
		reason:      reason,
	})
	if err != nil {
		return nil, err
	}

	if err = res.ExpectsStatus(http.StatusOK); err != nil {
		return nil, err
	}

	template := &objects.Template{}
	if err = res.JSON(template); err != nil {
		return nil, err
	}

	return template, nil
}
