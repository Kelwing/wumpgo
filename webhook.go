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

	webhook := &objects.Webhook{}
	err = NewRequest().
		Method(http.MethodPost).
		Path(fmt.Sprintf(ChannelWebhookFmt, channel)).
		ContentType(JsonContentType).
		Body(data).
		Reason(reason).
		Expect(http.StatusOK).
		Bind(webhook).
		Send(c)

	return webhook, err
}

func (c *Client) GetChannelWebhooks(channel objects.Snowflake) ([]*objects.Webhook, error) {
	webhooks := []*objects.Webhook{}
	err := NewRequest().
		Method(http.MethodGet).
		Path(fmt.Sprintf(ChannelWebhookFmt, channel)).
		ContentType(JsonContentType).
		Expect(http.StatusOK).
		Bind(&webhooks).
		Send(c)

	return webhooks, err
}

func (c *Client) GetGuildWebhooks(guild objects.Snowflake) ([]*objects.Webhook, error) {
	webhooks := []*objects.Webhook{}
	err := NewRequest().
		Method(http.MethodGet).
		Path(fmt.Sprintf(GuildWebhookFmt, guild)).
		ContentType(JsonContentType).
		Expect(http.StatusOK).
		Bind(&webhooks).
		Send(c)

	return webhooks, err
}

func (c *Client) GetWebhook(id objects.Snowflake) (*objects.Webhook, error) {
	webhook := &objects.Webhook{}
	err := NewRequest().
		Method(http.MethodGet).
		Path(fmt.Sprintf(WebhookFmt, id)).
		ContentType(JsonContentType).
		Expect(http.StatusOK).
		Bind(webhook).
		Send(c)

	return webhook, err
}

func (c *Client) GetWebhookWithToken(id objects.Snowflake, token string) (*objects.Webhook, error) {
	webhook := &objects.Webhook{}
	err := NewRequest().
		Method(http.MethodGet).
		Path(fmt.Sprintf(WebhookWithTokenFmt, id, token)).
		ContentType(JsonContentType).
		Expect(http.StatusOK).
		Bind(webhook).
		OmitAuth().
		Send(c)

	return webhook, err
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

	webhook := &objects.Webhook{}
	err = NewRequest().
		Method(http.MethodPatch).
		Path(fmt.Sprintf(WebhookFmt, id)).
		ContentType(JsonContentType).
		Body(data).
		Reason(reason).
		Expect(http.StatusOK).
		Bind(webhook).
		Send(c)

	return webhook, err
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

	webhook := &objects.Webhook{}
	err = NewRequest().
		Method(http.MethodPatch).
		Path(fmt.Sprintf(WebhookWithTokenFmt, id, token)).
		ContentType(JsonContentType).
		Body(data).
		Reason(reason).
		Expect(http.StatusOK).
		Bind(webhook).
		OmitAuth().
		Send(c)
	return webhook, err
}

func (c *Client) DeleteWebhook(id objects.Snowflake) error {
	return NewRequest().
		Method(http.MethodDelete).
		Path(fmt.Sprintf(WebhookFmt, id)).
		ContentType(JsonContentType).
		Expect(http.StatusNoContent).
		Send(c)
}

func (c *Client) DeleteWebhookWithToken(id objects.Snowflake, token string) error {
	return NewRequest().
		Method(http.MethodDelete).
		Path(fmt.Sprintf(WebhookWithTokenFmt, id, token)).
		ContentType(JsonContentType).
		Expect(http.StatusNoContent).
		OmitAuth().
		Send(c)
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

	msg := &objects.Message{}
	err = NewRequest().
		Method(http.MethodPost).
		Path(u.String()).
		ContentType(contentType).
		Body(body).
		Expect(http.StatusOK, http.StatusNoContent).
		Bind(msg).
		OmitAuth().Send(c)

	return msg, err
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

	msg := &objects.Message{}
	err = NewRequest().
		Method(http.MethodPatch).
		Path(fmt.Sprintf(WebhookMessageFmt, webhookID, token, messageID)).
		ContentType(JsonContentType).
		Body(data).
		Expect(http.StatusOK).
		Bind(msg).
		OmitAuth().
		Send(c)

	return msg, err
}

func (c *Client) DeleteWebhookMessage(messageID, webhookID objects.Snowflake, token string) error {
	return NewRequest().
		Method(http.MethodDelete).
		Path(fmt.Sprintf(WebhookMessageFmt, webhookID, token, messageID)).
		ContentType(JsonContentType).
		Expect(http.StatusNoContent).
		OmitAuth().
		Send(c)
}
