package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/Postcord/objects"
	"github.com/google/go-querystring/query"
)

type CreateWebhookParams struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Reason string `json:"-"`
}

func (c *Client) CreateWebhook(channel objects.Snowflake, params *CreateWebhookParams) (*objects.Webhook, error) {
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
		path:        fmt.Sprintf(ChannelWebhookFmt, channel),
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

	webhook := &objects.Webhook{}
	if err = res.JSON(webhook); err != nil {
		return nil, err
	}

	return webhook, nil
}

func (c *Client) GetChannelWebhooks(channel objects.Snowflake) ([]*objects.Webhook, error) {
	res, err := c.request(&request{
		method:      http.MethodGet,
		path:        fmt.Sprintf(ChannelWebhookFmt, channel),
		contentType: JsonContentType,
	})

	if err != nil {
		return nil, err
	}

	if err = res.ExpectsStatus(http.StatusOK); err != nil {
		return nil, err
	}

	var webhooks []*objects.Webhook
	if err = res.JSON(&webhooks); err != nil {
		return nil, err
	}
	return webhooks, nil
}

func (c *Client) GetGuildWebhooks(guild objects.Snowflake) ([]*objects.Webhook, error) {
	res, err := c.request(&request{
		method:      http.MethodGet,
		path:        fmt.Sprintf(GuildWebhookFmt, guild),
		contentType: JsonContentType,
	})
	if err != nil {
		return nil, err
	}

	if err = res.ExpectsStatus(http.StatusOK); err != nil {
		return nil, err
	}

	var webhooks []*objects.Webhook
	if err = res.JSON(&webhooks); err != nil {
		return nil, err
	}
	return webhooks, nil
}

func (c *Client) GetWebhook(id objects.Snowflake) (*objects.Webhook, error) {
	res, err := c.request(&request{
		method:      http.MethodGet,
		path:        fmt.Sprintf(WebhookFmt, id),
		contentType: JsonContentType,
	})
	if err != nil {
		return nil, err
	}

	if err = res.ExpectsStatus(http.StatusOK); err != nil {
		return nil, err
	}

	webhook := &objects.Webhook{}
	if err = res.JSON(webhook); err != nil {
		return nil, err
	}
	return webhook, nil
}

func (c *Client) GetWebhookWithToken(id objects.Snowflake, token string) (*objects.Webhook, error) {
	res, err := c.request(&request{
		method:      http.MethodGet,
		path:        fmt.Sprintf(WebhookWithTokenFmt, id, token),
		contentType: JsonContentType,
	})
	if err != nil {
		return nil, err
	}

	if err = res.ExpectsStatus(http.StatusOK); err != nil {
		return nil, err
	}

	webhook := &objects.Webhook{}
	if err = res.JSON(webhook); err != nil {
		return nil, err
	}
	return webhook, nil
}

type ModifyWebhookParams struct {
	Name      string            `json:"name,omitempty"`
	Avatar    string            `json:"avatar,omitempty"`
	ChannelID objects.Snowflake `json:"channel_id,omitempty"`
	Reason    string            `json:"-"`
}

func (c *Client) ModifyWebhook(id objects.Snowflake, params *ModifyWebhookParams) (*objects.Webhook, error) {
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
		path:        fmt.Sprintf(WebhookFmt, id),
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

	webhook := &objects.Webhook{}
	if err = res.JSON(webhook); err != nil {
		return nil, err
	}

	return webhook, nil
}

type ModifyWebhookWithTokenParams struct {
	Name   string `json:"name,omitempty"`
	Avatar string `json:"avatar,omitempty"`
	Reason string `json:"-"`
}

func (c *Client) ModifyWebhookWithToken(id objects.Snowflake, token string, params *ModifyWebhookWithTokenParams) (*objects.Webhook, error) {
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
		path:        fmt.Sprintf(WebhookWithTokenFmt, id, token),
		contentType: JsonContentType,
		body:        data,
		reason:      reason,
		omitAuth:    true,
	})
	if err != nil {
		return nil, err
	}

	if err = res.ExpectsStatus(http.StatusNoContent); err != nil {
		return nil, err
	}

	webhook := &objects.Webhook{}
	if err = res.JSON(webhook); err != nil {
		return nil, err
	}
	return webhook, nil
}

func (c *Client) DeleteWebhook(id objects.Snowflake) error {
	res, err := c.request(&request{
		method:      http.MethodDelete,
		path:        fmt.Sprintf(WebhookFmt, id),
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

func (c *Client) DeleteWebhookWithToken(id objects.Snowflake, token string) error {
	res, err := c.request(&request{
		method:      http.MethodDelete,
		path:        fmt.Sprintf(WebhookWithTokenFmt, id, token),
		contentType: JsonContentType,
		omitAuth:    true,
	})
	if err != nil {
		return err
	}
	if err = res.ExpectsStatus(http.StatusNoContent); err != nil {
		return err
	}
	return nil
}

type ExecuteWebhookParams struct {
	Wait bool `json:"-" url:"wait"`

	Content   string                     `json:"content,omitempty" url:"-"`
	Username  string                     `json:"username,omitempty" url:"-"`
	AvatarURL string                     `json:"avatar_url,omitempty" url:"-"`
	TTS       bool                       `json:"tts,omitempty" url:"-"`
	Files     []*CreateMessageFileParams `json:"-" url:"-"`
	Embeds    []*objects.Embed           `json:"embeds,omitempty" url:"-"`

	AllowedMentions *objects.AllowedMentions `json:"allowed_mentions,omitempty" url:"-"`
	Components      []*objects.Component     `json:"components,omitempty"`
}

func (c *Client) ExecuteWebhook(id objects.Snowflake, token string, params *ExecuteWebhookParams) (*objects.Message, error) {
	var contentType string
	var body []byte

	if len(params.Files) > 0 {
		buffer := new(bytes.Buffer)
		m := multipart.NewWriter(buffer)

		b, err := json.Marshal(params)
		if err != nil {
			return nil, err
		}

		if field, err := m.CreateFormField("payload_json"); err != nil {
			return nil, err
		} else {
			if _, err = field.Write(b); err != nil {
				return nil, err
			}
		}

		for n, file := range params.Files {
			if file.Spoiler && !strings.HasPrefix(file.Filename, "SPOILER_") {
				file.Filename = "SPOILER_" + file.Filename
			}

			w, err := m.CreateFormFile(fmt.Sprintf("file%d", n), file.Filename)
			if err != nil {
				return nil, err
			}

			if _, err = io.Copy(w, file.Reader); err != nil {
				return nil, err
			}
		}
		contentType = m.FormDataContentType()
		if err = m.Close(); err != nil {
			return nil, err
		}
		body = buffer.Bytes()
	} else {
		contentType = JsonContentType
		var err error
		if body, err = json.Marshal(params); err != nil {
			return nil, err
		}
	}

	u, err := url.Parse(fmt.Sprintf(WebhookWithTokenFmt, id, token))
	if err != nil {
		return nil, err
	}

	v, err := query.Values(params)
	if err != nil {
		return nil, err
	}

	u.RawQuery = v.Encode()

	res, err := c.request(&request{
		method:      http.MethodPost,
		path:        u.String(),
		contentType: contentType,
		body:        body,
		omitAuth:    true,
	})
	if err != nil {
		return nil, err
	}

	if err = res.ExpectAnyStatus(http.StatusOK, http.StatusNoContent); err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	msg := &objects.Message{}
	if err = res.JSON(msg); err != nil {
		return nil, err
	}
	return msg, nil
}

type EditWebhookMessageParams struct {
	Content         string                   `json:"content"`
	Embeds          []*objects.Embed         `json:"embeds"`
	AllowedMentions *objects.AllowedMentions `json:"allowed_mentions,omitempty"`
	Components      []*objects.Component     `json:"components"`
}

func (c *Client) EditWebhookMessage(messageID, webhookID objects.Snowflake, token string, params *EditWebhookMessageParams) (*objects.Message, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	res, err := c.request(&request{
		method:      http.MethodPatch,
		path:        fmt.Sprintf(WebhookMessageFmt, webhookID, token, messageID),
		contentType: JsonContentType,
		body:        data,
	})
	if err != nil {
		return nil, err
	}

	msg := &objects.Message{}
	if err = res.JSON(msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func (c *Client) DeleteWebhookMessage(messageID, webhookID objects.Snowflake, token string) error {
	res, err := c.request(&request{
		method:      http.MethodDelete,
		path:        fmt.Sprintf(WebhookMessageFmt, webhookID, token, messageID),
		contentType: JsonContentType,
	})
	if err != nil {
		return err
	}

	if err = res.ExpectAnyStatus(http.StatusNoContent); err != nil {
		return err
	}

	return nil
}

func (c *Client) EditOriginalInteractionResponse(applicationID objects.Snowflake, token string, params *EditWebhookMessageParams) (*objects.Message, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	res, err := c.request(&request{
		method:      http.MethodPatch,
		path:        fmt.Sprintf(EditOriginalInteractionResponseFmt, applicationID, token),
		contentType: JsonContentType,
		body:        data,
		omitAuth:    true,
	})
	if err != nil {
		return nil, err
	}

	msg := &objects.Message{}
	if err = res.JSON(msg); err != nil {
		return nil, err
	}
	return msg, nil
}
