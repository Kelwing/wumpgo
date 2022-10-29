package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
	"wumpgo.dev/wumpgo/objects"
)

type CreateWebhookParams struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Reason string `json:"-"`
}

func (c *Client) CreateWebhook(ctx context.Context, channel objects.SnowflakeObject, params *CreateWebhookParams) (*objects.Webhook, error) {
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
		WithContext(ctx).
		Path(fmt.Sprintf(ChannelWebhookFmt, channel.GetID())).
		ContentType(JsonContentType).
		Body(data).
		Reason(reason).
		Bind(webhook).
		Send(c)

	return webhook, err
}

func (c *Client) GetChannelWebhooks(ctx context.Context, channel objects.SnowflakeObject) ([]*objects.Webhook, error) {
	webhooks := []*objects.Webhook{}
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(ChannelWebhookFmt, channel.GetID())).
		ContentType(JsonContentType).
		Bind(&webhooks).
		Send(c)

	return webhooks, err
}

func (c *Client) GetGuildWebhooks(ctx context.Context, guild objects.SnowflakeObject) ([]*objects.Webhook, error) {
	webhooks := []*objects.Webhook{}
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildWebhookFmt, guild.GetID())).
		ContentType(JsonContentType).
		Bind(&webhooks).
		Send(c)

	return webhooks, err
}

func (c *Client) GetWebhook(ctx context.Context, id objects.SnowflakeObject) (*objects.Webhook, error) {
	webhook := &objects.Webhook{}
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(WebhookFmt, id.GetID())).
		ContentType(JsonContentType).
		Bind(webhook).
		Send(c)

	return webhook, err
}

func (c *Client) GetWebhookWithToken(ctx context.Context, id objects.SnowflakeObject, token string) (*objects.Webhook, error) {
	webhook := &objects.Webhook{}
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(WebhookWithTokenFmt, id.GetID(), token)).
		ContentType(JsonContentType).
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

func (c *Client) ModifyWebhook(ctx context.Context, id objects.SnowflakeObject, params *ModifyWebhookParams) (*objects.Webhook, error) {
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
		WithContext(ctx).
		Path(fmt.Sprintf(WebhookFmt, id.GetID())).
		ContentType(JsonContentType).
		Body(data).
		Reason(reason).
		Bind(webhook).
		Send(c)

	return webhook, err
}

type ModifyWebhookWithTokenParams struct {
	Name   string `json:"name,omitempty"`
	Avatar string `json:"avatar,omitempty"`
	Reason string `json:"-"`
}

func (c *Client) ModifyWebhookWithToken(ctx context.Context, id objects.SnowflakeObject, token string, params *ModifyWebhookWithTokenParams) (*objects.Webhook, error) {
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
		WithContext(ctx).
		Path(fmt.Sprintf(WebhookWithTokenFmt, id.GetID(), token)).
		ContentType(JsonContentType).
		Body(data).
		Reason(reason).
		Bind(webhook).
		OmitAuth().
		Send(c)
	return webhook, err
}

func (c *Client) DeleteWebhook(ctx context.Context, id objects.SnowflakeObject) error {
	return NewRequest().
		Method(http.MethodDelete).
		WithContext(ctx).
		Path(fmt.Sprintf(WebhookFmt, id.GetID())).
		ContentType(JsonContentType).
		Send(c)
}

func (c *Client) DeleteWebhookWithToken(ctx context.Context, id objects.SnowflakeObject, token string) error {
	return NewRequest().
		Method(http.MethodDelete).
		WithContext(ctx).
		Path(fmt.Sprintf(WebhookWithTokenFmt, id.GetID(), token)).
		ContentType(JsonContentType).
		OmitAuth().
		Send(c)
}

type ExecuteWebhookParams struct {
	Wait     bool              `json:"-" url:"wait"`
	ThreadID objects.Snowflake `json:"-" url:"thread_id,omitempty"`

	Content   string                 `json:"content,omitempty" url:"-"`
	Username  string                 `json:"username,omitempty" url:"-"`
	AvatarURL string                 `json:"avatar_url,omitempty" url:"-"`
	TTS       bool                   `json:"tts,omitempty" url:"-"`
	Files     []*objects.DiscordFile `json:"-" url:"-"`
	Embeds    []*objects.Embed       `json:"embeds,omitempty" url:"-"`

	AllowedMentions *objects.AllowedMentions `json:"allowed_mentions,omitempty" url:"-"`
	Components      []*objects.Component     `json:"components,omitempty" url:"-"`
	Attachments     []*objects.Attachment    `json:"attachments,omitempty" url:"-"`
	Flags           objects.MessageFlag      `json:"flags,omitempty" url:"-"`
}

func (c *Client) ExecuteWebhook(ctx context.Context, id objects.SnowflakeObject, token string, params *ExecuteWebhookParams) (*objects.Message, error) {
	var contentType string
	var body []byte

	if len(params.Files) > 0 {
		buffer := new(bytes.Buffer)
		m := multipart.NewWriter(buffer)

		for n, file := range params.Files {
			a, err := file.GenerateAttachment(objects.Snowflake(n+1), m)
			if err != nil {
				continue
			}
			params.Attachments = append(params.Attachments, a)
		}

		if w, err := m.CreateFormField("payload_json"); err != nil {
			return nil, err
		} else {
			if err := json.NewEncoder(w).Encode(params); err != nil {
				return nil, err
			}
		}
		contentType = m.FormDataContentType()
		if err := m.Close(); err != nil {
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

	u, err := url.Parse(fmt.Sprintf(WebhookWithTokenFmt, id.GetID(), token))
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
		WithContext(ctx).
		Path(u.String()).
		ContentType(contentType).
		Body(body).
		Bind(msg).
		OmitAuth().Send(c)

	return msg, err
}

func (c *Client) GetWebhookMesssage(ctx context.Context, messageID, webhookID objects.SnowflakeObject, token string) (*objects.Message, error) {
	msg := &objects.Message{}
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(WebhookMessageFmt, webhookID.GetID(), token, messageID.GetID())).
		ContentType(JsonContentType).
		Bind(msg).
		OmitAuth().
		Send(c)

	return msg, err
}

type EditWebhookMessageParams struct {
	Content         string                   `json:"content"`
	Embeds          []*objects.Embed         `json:"embeds"`
	AllowedMentions *objects.AllowedMentions `json:"allowed_mentions,omitempty"`
	Components      []*objects.Component     `json:"components"`
}

func (c *Client) EditWebhookMessage(ctx context.Context, messageID, webhookID objects.SnowflakeObject, token string, params *EditWebhookMessageParams) (*objects.Message, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	msg := &objects.Message{}
	err = NewRequest().
		Method(http.MethodPatch).
		WithContext(ctx).
		Path(fmt.Sprintf(WebhookMessageFmt, webhookID.GetID(), token, messageID.GetID())).
		ContentType(JsonContentType).
		Body(data).
		Bind(msg).
		OmitAuth().
		Send(c)

	return msg, err
}

func (c *Client) DeleteWebhookMessage(ctx context.Context, messageID, webhookID objects.SnowflakeObject, token string) error {
	return NewRequest().
		Method(http.MethodDelete).
		WithContext(ctx).
		Path(fmt.Sprintf(WebhookMessageFmt, webhookID.GetID(), token, messageID.GetID())).
		ContentType(JsonContentType).
		OmitAuth().
		Send(c)
}
