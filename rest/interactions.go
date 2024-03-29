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

func (c *Client) CreateInteractionResponse(ctx context.Context, interactionID objects.Snowflake, token string, response *objects.InteractionResponse) error {
	var contentType string
	var body []byte

	var data *objects.InteractionMessagesCallbackData
	switch d := response.Data.(type) {
	case objects.InteractionMessagesCallbackData:
		data = &d
	case *objects.InteractionMessagesCallbackData:
		data = d
	default:
		data = nil
	}

	if data != nil && len(data.Files) > 0 {
		buffer := new(bytes.Buffer)
		m := multipart.NewWriter(buffer)

		for n, file := range data.Files {
			a, err := file.GenerateAttachment(objects.Snowflake(n+1), m)
			if err != nil {
				continue
			}
			data.Attachments = append(data.Attachments, a)
		}

		if w, err := m.CreateFormField("payload_json"); err != nil {
			return err
		} else {
			if err := json.NewEncoder(w).Encode(response); err != nil {
				return err
			}
		}
		contentType = m.FormDataContentType()
		if err := m.Close(); err != nil {
			return err
		}
		body = buffer.Bytes()
	} else {
		contentType = JsonContentType
		var err error
		if body, err = json.Marshal(response); err != nil {
			return err
		}
	}

	return NewRequest().
		Method(http.MethodPost).
		WithContext(ctx).
		Path(fmt.Sprintf(CreateInteractionResponseFmt, interactionID, token)).
		Body(body).
		ContentType(contentType).
		OmitAuth().
		Send(c)
}

func (c *Client) GetOriginalInteractionResponse(ctx context.Context, applicationID objects.Snowflake, token string) (*objects.Message, error) {
	msg := &objects.Message{}
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(OriginalInteractionResponseFmt, applicationID, token)).
		Bind(msg).
		OmitAuth().
		Send(c)
	return msg, err
}

func (c *Client) EditOriginalInteractionResponse(ctx context.Context, applicationID objects.Snowflake, token string, params *EditWebhookMessageParams) (*objects.Message, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	msg := &objects.Message{}
	err = NewRequest().
		Method(http.MethodPatch).
		WithContext(ctx).
		Path(fmt.Sprintf(OriginalInteractionResponseFmt, applicationID, token)).
		Body(data).
		ContentType(JsonContentType).
		Bind(msg).
		OmitAuth().
		Send(c)
	return msg, err
}

func (c *Client) DeleteOriginalInteractionResponse(ctx context.Context, applicationID objects.Snowflake, token string) error {
	return NewRequest().
		Method(http.MethodDelete).
		WithContext(ctx).
		Path(fmt.Sprintf(OriginalInteractionResponseFmt, applicationID, token)).
		OmitAuth().
		Send(c)
}

type CreateFollowupMessageParams struct {
	Content string                 `json:"content,omitempty" url:"-"`
	TTS     bool                   `json:"tts,omitempty" url:"-"`
	Files   []*objects.DiscordFile `json:"-" url:"-"`
	Embeds  []*objects.Embed       `json:"embeds,omitempty" url:"-"`

	AllowedMentions *objects.AllowedMentions `json:"allowed_mentions,omitempty" url:"-"`
	Components      []*objects.Component     `json:"components,omitempty" url:"-"`
	Attachments     []*objects.Attachment    `json:"attachments,omitempty" url:"-"`
	Flags           objects.MessageFlag      `json:"flags,omitempty" url:"-"`
}

func (c *Client) CreateFollowupMessage(ctx context.Context, applicationID objects.Snowflake, token string, params *CreateFollowupMessageParams) (*objects.Message, error) {
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

	u, err := url.Parse(fmt.Sprintf(WebhookWithTokenFmt, applicationID, token))
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

func (c *Client) GetFollowupMessage(ctx context.Context, applicationID objects.Snowflake, token string, messageID objects.Snowflake) (*objects.Message, error) {
	msg := &objects.Message{}
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(WebhookMessageFmt, applicationID, token, messageID)).
		Bind(msg).
		OmitAuth().
		Send(c)
	return msg, err
}

func (c *Client) EditFollowupMessage(ctx context.Context, applicationID objects.Snowflake, token string, messageID objects.Snowflake, params *EditWebhookMessageParams) (*objects.Message, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	msg := &objects.Message{}
	err = NewRequest().
		Method(http.MethodPatch).
		WithContext(ctx).
		Path(fmt.Sprintf(WebhookMessageFmt, applicationID, token, messageID)).
		Body(data).
		ContentType(JsonContentType).
		Bind(msg).
		OmitAuth().
		Send(c)
	return msg, err
}

func (c *Client) DeleteFollowupMessage(ctx context.Context, applicationID objects.Snowflake, token string, messageID objects.Snowflake) error {
	return NewRequest().
		Method(http.MethodDelete).
		WithContext(ctx).
		Path(fmt.Sprintf(WebhookMessageFmt, applicationID, token, messageID)).
		OmitAuth().
		Send(c)
}
