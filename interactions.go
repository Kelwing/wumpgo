package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/Postcord/objects"
	"github.com/google/go-querystring/query"
)

func (c *Client) CreateInteractionResponse(ctx context.Context, interactionID objects.SnowflakeObject, token string, response *objects.InteractionResponse) error {
	var contentType string
	var body []byte

	if len(response.Data.Files) > 0 {
		buffer := new(bytes.Buffer)
		m := multipart.NewWriter(buffer)

		for n, file := range response.Data.Files {
			a, err := file.GenerateAttachment(objects.Snowflake(n+1), m)
			if err != nil {
				continue
			}
			response.Data.Attachments = append(response.Data.Attachments, a)
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
		Path(fmt.Sprintf(CreateInteractionResponseFmt, interactionID.GetID(), token)).
		Body(body).
		ContentType(contentType).
		OmitAuth().
		Send(c)
}

func (c *Client) GetOriginalInteractionResponse(ctx context.Context, applicationID objects.SnowflakeObject, token string) (*objects.Message, error) {
	msg := &objects.Message{}
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(OriginalInteractionResponseFmt, applicationID.GetID(), token)).
		Bind(msg).
		Expect(http.StatusOK).
		OmitAuth().
		Send(c)
	return msg, err
}

func (c *Client) EditOriginalInteractionResponse(ctx context.Context, applicationID objects.SnowflakeObject, token string, params *EditWebhookMessageParams) (*objects.Message, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	msg := &objects.Message{}
	err = NewRequest().
		Method(http.MethodPatch).
		WithContext(ctx).
		Path(fmt.Sprintf(OriginalInteractionResponseFmt, applicationID.GetID(), token)).
		Body(data).
		ContentType(JsonContentType).
		Bind(msg).
		OmitAuth().
		Send(c)
	return msg, err
}

func (c *Client) DeleteOriginalInteractionResponse(ctx context.Context, applicationID objects.SnowflakeObject, token string) error {
	return NewRequest().
		Method(http.MethodDelete).
		WithContext(ctx).
		Path(fmt.Sprintf(OriginalInteractionResponseFmt, applicationID.GetID(), token)).
		Expect(http.StatusNoContent).
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

func (c *Client) CreateFollowupMessage(ctx context.Context, applicationID objects.SnowflakeObject, token string, params *CreateFollowupMessageParams) (*objects.Message, error) {
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

	u, err := url.Parse(fmt.Sprintf(WebhookWithTokenFmt, applicationID.GetID(), token))
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
		Expect(http.StatusOK, http.StatusNoContent).
		Bind(msg).
		OmitAuth().Send(c)

	return msg, err
}

func (c *Client) GetFollowupMessage(ctx context.Context, applicationID objects.SnowflakeObject, token string, messageID objects.SnowflakeObject) (*objects.Message, error) {
	msg := &objects.Message{}
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(WebhookMessageFmt, applicationID.GetID(), token, messageID.GetID())).
		Bind(msg).
		Expect(http.StatusOK).
		OmitAuth().
		Send(c)
	return msg, err
}

func (c *Client) EditFollowupMessage(ctx context.Context, applicationID objects.SnowflakeObject, token string, messageID objects.SnowflakeObject, params *EditWebhookMessageParams) (*objects.Message, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	msg := &objects.Message{}
	err = NewRequest().
		Method(http.MethodPatch).
		WithContext(ctx).
		Path(fmt.Sprintf(WebhookMessageFmt, applicationID.GetID(), token, messageID.GetID())).
		Body(data).
		Expect(http.StatusOK).
		ContentType(JsonContentType).
		Bind(msg).
		OmitAuth().
		Send(c)
	return msg, err
}

func (c *Client) DeleteFollowupMessage(ctx context.Context, applicationID objects.SnowflakeObject, token string, messageID objects.SnowflakeObject) error {
	return NewRequest().
		Method(http.MethodDelete).
		WithContext(ctx).
		Path(fmt.Sprintf(WebhookMessageFmt, applicationID.GetID(), token, messageID.GetID())).
		Expect(http.StatusNoContent).
		OmitAuth().
		Send(c)
}
