package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Postcord/objects"
	"github.com/google/go-querystring/query"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

func (c *Client) GetChannel(id objects.Snowflake) (*objects.Channel, error) {
	resp, err := c.request(&request{
		method:      http.MethodGet,
		path:        fmt.Sprintf(ChannelBaseFmt, id),
		contentType: JsonContentType,
	})
	if err != nil {
		return nil, err
	}
	if err = resp.ExpectsStatus(http.StatusOK); err != nil {
		return nil, err
	}

	channel := &objects.Channel{}

	if err = resp.JSON(channel); err != nil {
		return nil, err
	}

	return channel, nil
}

type ModifyChannelParams struct {
	Name                 string                        `json:"name"`
	Type                 int64                         `json:"type"`
	Position             int64                         `json:"position,omitempty"`
	Topic                string                        `json:"topic,omitempty"`
	NSFW                 bool                          `json:"nsfw,omitempty"`
	RateLimitPerUser     int64                         `json:"rate_limit_per_user,omitempty"`
	Bitrate              int64                         `json:"bitrate,omitempty"`
	UserLimit            int64                         `json:"user_limit,omitempty"`
	PermissionOverwrites []objects.PermissionOverwrite `json:"permission_overwrites,omitempty"`
	Parent               objects.Snowflake             `json:"parent_id,omitempty"`
	Reason               string                        `json:"-"`
}

func (c *Client) ModifyChannel(id objects.Snowflake, params *ModifyChannelParams) (*objects.Channel, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	reason := ""
	if params != nil {
		reason = params.Reason
	}

	resp, err := c.request(&request{
		method:      http.MethodPatch,
		path:        fmt.Sprintf(ChannelBaseFmt, id),
		contentType: JsonContentType,
		body:        data,
		reason:      reason,
	})
	if err != nil {
		return nil, err
	}
	if err = resp.ExpectsStatus(http.StatusOK); err != nil {
		return nil, err
	}

	channel := &objects.Channel{}

	if err = resp.JSON(channel); err != nil {
		return nil, err
	}

	return channel, nil
}

func (c *Client) DeleteChannel(id objects.Snowflake, reason string) (*objects.Channel, error) {
	resp, err := c.request(&request{
		method:      http.MethodDelete,
		path:        fmt.Sprintf(ChannelBaseFmt, id),
		contentType: JsonContentType,
		reason:      reason,
	})
	if err != nil {
		return nil, err
	}

	if err = resp.ExpectsStatus(http.StatusOK); err != nil {
		return nil, err
	}

	channel := &objects.Channel{}

	if err = resp.JSON(channel); err != nil {
		return nil, err
	}

	return channel, nil
}

type GetChannelMessagesParams struct {
	Around objects.Snowflake `url:"around,omitempty"`
	Before objects.Snowflake `url:"before,omitempty"`
	After  objects.Snowflake `url:"after,omitempty"`
	Limit  int               `url:"limit,omitempty"`
}

func (c *Client) GetChannelMessages(id objects.Snowflake, params *GetChannelMessagesParams) ([]*objects.Message, error) {
	u, err := url.Parse(fmt.Sprintf(ChannelMessagesFmt, id))
	if err != nil {
		return nil, err
	}
	q, err := query.Values(params)
	if err != nil {
		return nil, err
	}
	u.RawQuery = q.Encode()

	res, err := c.request(&request{
		method:      http.MethodGet,
		path:        u.String(),
		contentType: JsonContentType,
	})
	if err != nil {
		return nil, err
	}

	if err = res.ExpectsStatus(http.StatusOK); err != nil {
		return nil, err
	}

	var messages []*objects.Message

	if err = res.JSON(&messages); err != nil {
		return nil, err
	}

	return messages, nil
}

func (c *Client) GetChannelMessage(channel, message objects.Snowflake) (*objects.Message, error) {
	res, err := c.request(&request{
		method:      http.MethodGet,
		path:        fmt.Sprintf(ChannelMessageFmt, channel, message),
		contentType: JsonContentType,
	})
	if err != nil {
		return nil, err
	}

	if err = res.ExpectsStatus(http.StatusOK); err != nil {
		return nil, err
	}

	msg := &objects.Message{}

	if err = res.JSON(msg); err != nil {
		return nil, err
	}

	return msg, nil
}

func (c *Client) CrossPostMessage(channel, message objects.Snowflake) (*objects.Message, error) {
	res, err := c.request(&request{
		method:      http.MethodPost,
		path:        fmt.Sprintf(CrosspostMessageFmt, channel, message),
		contentType: JsonContentType,
	})
	if err != nil {
		return nil, err
	}

	if err = res.ExpectsStatus(http.StatusOK); err != nil {
		return nil, err
	}

	msg := &objects.Message{}

	if err = res.JSON(msg); err != nil {
		return nil, err
	}

	return msg, nil
}

func (c *Client) DeleteMessage(channel, message objects.Snowflake) error {
	res, err := c.request(&request{
		method:      http.MethodDelete,
		path:        fmt.Sprintf(ChannelMessageFmt, channel, message),
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

type DeleteMessagesParams struct {
	Messages []objects.Snowflake `json:"messages"`
}

func (c *Client) BulkDeleteMessages(channel objects.Snowflake, params *DeleteMessagesParams) error {
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}

	res, err := c.request(&request{
		method:      http.MethodPost,
		path:        fmt.Sprintf(BulkDeleteMessagesFmt, channel),
		contentType: JsonContentType,
		body:        data,
	})
	if err != nil {
		return err
	}

	if err = res.ExpectsStatus(http.StatusNoContent); err != nil {
		return err
	}

	return nil
}

type EditChannelParams struct {
	Allow  objects.PermissionBit `json:"allow"`
	Deny   objects.PermissionBit `json:"deny"`
	Type   int                   `json:"type"`
	Reason string                `json:"-"`
}

func (c *Client) EditChannelPermissions(channel, overwrite objects.Snowflake, params *EditChannelParams) error {
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}

	reason := ""
	if params != nil {
		reason = params.Reason
	}

	res, err := c.request(&request{
		method:      http.MethodPut,
		path:        fmt.Sprintf(ChannelPermissionsFmt, channel, overwrite),
		contentType: JsonContentType,
		body:        data,
		reason:      reason,
	})
	if err != nil {
		return err
	}

	if err = res.ExpectsStatus(http.StatusNoContent); err != nil {
		return err
	}
	return nil
}

func (c *Client) DeleteChannelPermission(channel, overwrite objects.Snowflake, reason string) error {
	res, err := c.request(&request{
		method:      http.MethodDelete,
		path:        fmt.Sprintf(ChannelPermissionsFmt, channel, overwrite),
		contentType: JsonContentType,
		reason:      reason,
	})
	if err != nil {
		return err
	}

	if err = res.ExpectsStatus(http.StatusNoContent); err != nil {
		return err
	}

	return nil
}

func (c *Client) GetChannelInvites(channel objects.Snowflake) ([]*objects.Invite, error) {
	res, err := c.request(&request{
		method:      http.MethodGet,
		path:        fmt.Sprintf(ChannelInvitesFmt, channel),
		contentType: JsonContentType,
	})
	if err != nil {
		return nil, err
	}

	if err = res.ExpectsStatus(http.StatusOK); err != nil {
		return nil, err
	}

	var invites []*objects.Invite
	if err = res.JSON(&invites); err != nil {
		return nil, err
	}

	return invites, nil
}

type CreateInviteParams struct {
	MaxAge         int64             `json:"max_age,omitempty"`
	MaxUses        int64             `json:"max_uses,omitempty"`
	Temporary      bool              `json:"temporary,,omitempty"`
	Unique         bool              `json:"unique,omitempty"`
	TargetUser     objects.Snowflake `json:"target_user,omitempty"`
	TargetUserType objects.Snowflake `json:"target_user_type,omitempty"`
	Reason         string            `json:"-"`
}

func (c *Client) CreateChannelInvite(channel objects.Snowflake, params *CreateInviteParams) (*objects.Invite, error) {
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
		path:        fmt.Sprintf(ChannelInvitesFmt, channel),
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

	invite := &objects.Invite{}

	if err = res.JSON(invite); err != nil {
		return nil, err
	}
	return invite, nil
}

func (c *Client) getEmoji(emoji interface{}) (string, error) {
	var react string

	switch t := emoji.(type) {
	case objects.Emoji:
		react = fmt.Sprintf("%s:%d", t.Name, t.ID)
	case *objects.Emoji:
		react = fmt.Sprintf("%s:%d", t.Name, t.ID)
	case string:
		react = t
	default:
		return "", errors.New(fmt.Sprintf("invalid emoji type, %T", t))
	}

	return react, nil
}

func (c *Client) CreateReaction(channel, message objects.Snowflake, emoji interface{}) error {
	react, err := c.getEmoji(emoji)
	if err != nil {
		return err
	}

	res, err := c.request(&request{
		method:      http.MethodPut,
		path:        fmt.Sprintf(ReactionFmt, channel, message, url.QueryEscape(react), "@me"),
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

func (c *Client) DeleteOwnReaction(channel, message objects.Snowflake, emoji interface{}) error {
	react, err := c.getEmoji(emoji)
	if err != nil {
		return err
	}

	res, err := c.request(&request{
		method:      http.MethodDelete,
		path:        fmt.Sprintf(ReactionFmt, channel, message, url.QueryEscape(react), "@me"),
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

func (c *Client) DeleteUserReaction(channel, message, user objects.Snowflake, emoji interface{}) error {
	react, err := c.getEmoji(emoji)
	if err != nil {
		return err
	}

	res, err := c.request(&request{
		method:      http.MethodDelete,
		path:        fmt.Sprintf(ReactionUserFmt, channel, message, url.QueryEscape(react), user),
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

type GetReactionsParams struct {
	Before objects.Snowflake `url:"before,omitempty"`
	After  objects.Snowflake `url:"after,omitempty"`
	Limit  int               `url:"limit"`
}

func (c *Client) GetReactions(channel, message objects.Snowflake, emoji interface{}, params *GetReactionsParams) ([]*objects.User, error) {
	react, err := c.getEmoji(emoji)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(fmt.Sprintf(ReactionsFmt, channel, message, url.QueryEscape(react)))
	if err != nil {
		return nil, err
	}

	q, err := query.Values(params)
	if err != nil {
		return nil, err
	}
	u.RawQuery = q.Encode()

	res, err := c.request(&request{
		method:      http.MethodGet,
		path:        u.String(),
		contentType: JsonContentType,
	})
	if err != nil {
		return nil, err
	}

	if err = res.ExpectsStatus(http.StatusOK); err != nil {
		return nil, err
	}

	var users []*objects.User
	if err = res.JSON(&users); err != nil {
		return nil, err
	}
	return users, nil
}

func (c *Client) DeleteAllReactions(channel, message objects.Snowflake) error {
	res, err := c.request(&request{
		method:      http.MethodDelete,
		path:        fmt.Sprintf(ReactionsBaseFmt, channel, message),
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

func (c *Client) DeleteEmojiReactions(channel, message objects.Snowflake, emoji interface{}) error {
	reaction, err := c.getEmoji(emoji)
	if err != nil {
		return err
	}

	res, err := c.request(&request{
		method:      http.MethodDelete,
		path:        fmt.Sprintf(ReactionsFmt, channel, message, reaction),
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

func (c *Client) GetPinnedMessages(channel objects.Snowflake) ([]*objects.Message, error) {
	res, err := c.request(&request{
		method:      http.MethodGet,
		path:        fmt.Sprintf(ChannelPinsFmt, channel),
		contentType: JsonContentType,
	})
	if err != nil {
		return nil, err
	}

	if err = res.ExpectsStatus(http.StatusOK); err != nil {
		return nil, err
	}

	var messages []*objects.Message
	if err = res.JSON(&messages); err != nil {
		return nil, err
	}
	return messages, nil
}

func (c *Client) AddPinnedMessage(channel, message objects.Snowflake) error {
	res, err := c.request(&request{
		method:      http.MethodPut,
		path:        fmt.Sprintf(ChannelPinnedFmt, channel, message),
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

func (c *Client) DeletePinnedMessage(channel, message objects.Snowflake) error {
	res, err := c.request(&request{
		method:      http.MethodDelete,
		path:        fmt.Sprintf(ChannelPinnedFmt, channel, message),
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

type CreateMessageFileParams struct {
	Reader   io.Reader
	Filename string
	Spoiler  bool
}

type CreateMessageParams struct {
	Content          string                     `json:"content,omitempty"`
	Nonce            string                     `json:"nonce,omitempty"`
	TTS              bool                       `json:"tts,omitempty"`
	Files            []*CreateMessageFileParams `json:"-"`
	Embed            *objects.Embed             `json:"embed,omitempty"`
	AllowedMentions  *objects.AllowedMentions   `json:"allowed_mentions,omitempty"`
	MessageReference *objects.MessageReference  `json:"message_reference,omitempty"`
	Components       []*objects.Component       `json:"components,omitempty"`
}

func (c *Client) CreateMessage(channel objects.Snowflake, params *CreateMessageParams) (*objects.Message, error) {
	var contentType string
	var body []byte

	if len(params.Files) > 0 {
		buffer := new(bytes.Buffer)
		m := multipart.NewWriter(buffer)

		b, err := json.Marshal(params)
		if err != nil {
			return nil, err
		}

		if w, err := m.CreateFormField("payload_json"); err != nil {
			return nil, err
		} else {
			if _, err = w.Write(b); err != nil {
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
		body, err = json.Marshal(params)
		if err != nil {
			return nil, err
		}
	}

	res, err := c.request(&request{
		method:      http.MethodPost,
		path:        fmt.Sprintf(ChannelMessagesFmt, channel),
		contentType: contentType,
		body:        body,
	})
	if err != nil {
		return nil, err
	}

	if err = res.ExpectsStatus(http.StatusOK); err != nil {
		return nil, err
	}

	msg := &objects.Message{}
	if err = res.JSON(msg); err != nil {
		return nil, err
	}

	return msg, nil
}

type EditMessageParams struct {
	Content string              `json:"content,omitempty"`
	Embed   *objects.Embed      `json:"embed,omitempty"`
	Flags   objects.MessageFlag `json:"flags,omitempty"`
}

func (c *Client) EditMessage(channel, message objects.Snowflake, params *EditMessageParams) (*objects.Message, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	res, err := c.request(&request{
		method:      http.MethodPatch,
		path:        fmt.Sprintf(ChannelMessageFmt, channel, message),
		contentType: JsonContentType,
		body:        body,
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

func (c *Client) FollowNewsChannel(channel objects.Snowflake) (*objects.FollowedChannel, error) {
	res, err := c.request(&request{
		method:      http.MethodDelete,
		path:        fmt.Sprintf(ChannelFollowersFmt, channel),
		contentType: JsonContentType,
	})
	if err != nil {
		return nil, err
	}

	if err = res.ExpectsStatus(http.StatusOK); err != nil {
		return nil, err
	}

	followedChannel := &objects.FollowedChannel{}
	if err = res.JSON(followedChannel); err != nil {
		return nil, err
	}

	return followedChannel, nil
}

func (c *Client) StartTyping(channel objects.Snowflake) error {
	res, err := c.request(&request{
		method:      http.MethodPost,
		path:        fmt.Sprintf(ChannelTypingFmt, channel),
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
