package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/Postcord/objects"
	"github.com/Postcord/objects/permissions"
	"github.com/google/go-querystring/query"
)

func (c *Client) GetChannel(ctx context.Context, id objects.SnowflakeObject) (*objects.Channel, error) {
	channel := &objects.Channel{}

	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(ChannelBaseFmt, id.GetID())).
		ContentType(JsonContentType).
		Bind(channel).
		Expect(http.StatusOK).
		Send(c)

	return channel, err
}

type ModifyChannelParams struct {
	Name                 *string                        `json:"name,omitempty"`
	Type                 *objects.ChannelType           `json:"type,omitempty"`
	Position             *int64                         `json:"position,omitempty"`
	Topic                *string                        `json:"topic,omitempty"`
	NSFW                 *bool                          `json:"nsfw,omitempty"`
	RateLimitPerUser     *int64                         `json:"rate_limit_per_user,omitempty"`
	Bitrate              *int64                         `json:"bitrate,omitempty"`
	UserLimit            *int64                         `json:"user_limit,omitempty"`
	PermissionOverwrites *[]objects.PermissionOverwrite `json:"permission_overwrites,omitempty"`
	Parent               *objects.Snowflake             `json:"parent_id,omitempty"`
	Reason               string                         `json:"-"`
}

func (c *Client) ModifyChannel(ctx context.Context, id objects.SnowflakeObject, params *ModifyChannelParams) (*objects.Channel, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	reason := ""
	if params != nil {
		reason = params.Reason
	}
	channel := &objects.Channel{}

	err = NewRequest().
		Method(http.MethodPatch).
		WithContext(ctx).
		Path(fmt.Sprintf(ChannelBaseFmt, id.GetID())).
		ContentType(JsonContentType).
		Body(data).
		Reason(reason).
		Expect(http.StatusOK).
		Bind(channel).
		Send(c)

	return channel, err
}

func (c *Client) DeleteChannel(ctx context.Context, id objects.SnowflakeObject, reason string) (*objects.Channel, error) {
	channel := &objects.Channel{}
	err := NewRequest().
		Method(http.MethodDelete).
		WithContext(ctx).
		Path(fmt.Sprintf(ChannelBaseFmt, id.GetID())).
		Reason(reason).
		ContentType(JsonContentType).
		Expect(http.StatusOK).
		Bind(channel).
		Send(c)

	return channel, err
}

type GetChannelMessagesParams struct {
	Around objects.Snowflake `url:"around,omitempty"`
	Before objects.Snowflake `url:"before,omitempty"`
	After  objects.Snowflake `url:"after,omitempty"`
	Limit  int               `url:"limit,omitempty"`
}

func (c *Client) GetChannelMessages(ctx context.Context, id objects.SnowflakeObject, params *GetChannelMessagesParams) ([]*objects.Message, error) {
	u, err := url.Parse(fmt.Sprintf(ChannelMessagesFmt, id.GetID()))
	if err != nil {
		return nil, err
	}
	q, err := query.Values(params)
	if err != nil {
		return nil, err
	}
	u.RawQuery = q.Encode()

	var messages []*objects.Message
	err = NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(u.String()).
		ContentType(JsonContentType).
		Bind(&messages).
		Expect(http.StatusOK).
		Send(c)

	return messages, err
}

func (c *Client) GetChannelMessage(ctx context.Context, channel, message objects.SnowflakeObject) (*objects.Message, error) {
	msg := &objects.Message{}
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(ChannelMessageFmt, channel.GetID(), message.GetID())).
		ContentType(JsonContentType).
		Bind(msg).
		Expect(http.StatusOK).
		Send(c)

	return msg, err
}

func (c *Client) CrossPostMessage(ctx context.Context, channel, message objects.SnowflakeObject) (*objects.Message, error) {
	msg := &objects.Message{}
	err := NewRequest().
		Method(http.MethodPost).
		WithContext(ctx).
		Path(fmt.Sprintf(CrosspostMessageFmt, channel.GetID(), message.GetID())).
		ContentType(JsonContentType).
		Bind(msg).
		Expect(http.StatusOK).
		Send(c)

	return msg, err
}

func (c *Client) DeleteMessage(ctx context.Context, channel, message objects.SnowflakeObject) error {
	return NewRequest().
		Method(http.MethodDelete).
		WithContext(ctx).
		Path(fmt.Sprintf(ChannelMessageFmt, channel.GetID(), message.GetID())).
		ContentType(JsonContentType).
		Expect(http.StatusNoContent).
		Send(c)
}

type DeleteMessagesParams struct {
	Messages []objects.Snowflake `json:"messages"`
}

func (c *Client) BulkDeleteMessages(ctx context.Context, channel objects.SnowflakeObject, params *DeleteMessagesParams) error {
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}

	return NewRequest().
		Method(http.MethodPost).
		WithContext(ctx).
		Path(fmt.Sprintf(BulkDeleteMessagesFmt, channel.GetID())).
		ContentType(JsonContentType).
		Body(data).
		Expect(http.StatusNoContent).
		Send(c)
}

type EditChannelParams struct {
	Allow  permissions.PermissionBit `json:"allow"`
	Deny   permissions.PermissionBit `json:"deny"`
	Type   int                       `json:"type"`
	Reason string                    `json:"-"`
}

func (c *Client) EditChannelPermissions(ctx context.Context, channel, overwrite objects.SnowflakeObject, params *EditChannelParams) error {
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}

	reason := ""
	if params != nil {
		reason = params.Reason
	}

	return NewRequest().
		Method(http.MethodPut).
		WithContext(ctx).
		Path(fmt.Sprintf(ChannelPermissionsFmt, channel.GetID(), overwrite.GetID())).
		ContentType(JsonContentType).
		Body(data).
		Reason(reason).
		Expect(http.StatusNoContent).
		Send(c)
}

func (c *Client) DeleteChannelPermission(ctx context.Context, channel, overwrite objects.SnowflakeObject, reason string) error {
	return NewRequest().
		Method(http.MethodDelete).
		WithContext(ctx).
		Path(fmt.Sprintf(ChannelPermissionsFmt, channel.GetID(), overwrite.GetID())).
		Reason(reason).
		ContentType(JsonContentType).
		Expect(http.StatusNoContent).
		Send(c)
}

func (c *Client) GetChannelInvites(ctx context.Context, channel objects.SnowflakeObject) ([]*objects.Invite, error) {
	var invites []*objects.Invite
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(ChannelInvitesFmt, channel.GetID())).
		ContentType(JsonContentType).
		Bind(&invites).
		Expect(http.StatusOK).
		Send(c)

	return invites, err
}

type CreateInviteParams struct {
	MaxAge            int64                    `json:"max_age,omitempty"`
	MaxUses           int64                    `json:"max_uses,omitempty"`
	Temporary         bool                     `json:"temporary,,omitempty"`
	Unique            bool                     `json:"unique,omitempty"`
	TargetUser        objects.Snowflake        `json:"target_user_id,omitempty"`
	TargetType        objects.InviteTargetType `json:"target_type,omitempty"`
	TargetApplication objects.Snowflake        `json:"target_application_id,omitempty"`
	Reason            string                   `json:"-"`
}

func (c *Client) CreateChannelInvite(ctx context.Context, channel objects.SnowflakeObject, params *CreateInviteParams) (*objects.Invite, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	reason := ""
	if params != nil {
		reason = params.Reason
	}

	invite := &objects.Invite{}

	err = NewRequest().
		Method(http.MethodPost).
		WithContext(ctx).
		Path(fmt.Sprintf(ChannelInvitesFmt, channel.GetID())).
		ContentType(JsonContentType).
		Body(data).
		Reason(reason).
		Bind(invite).
		Expect(http.StatusOK).
		Send(c)

	return invite, err
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
		return "", fmt.Errorf("invalid emoji type, %T", t)
	}

	return react, nil
}

func (c *Client) CreateReaction(ctx context.Context, channel, message objects.SnowflakeObject, emoji interface{}) error {
	react, err := c.getEmoji(emoji)
	if err != nil {
		return err
	}

	return NewRequest().
		Method(http.MethodPut).
		WithContext(ctx).
		Path(fmt.Sprintf(ReactionFmt, channel.GetID(), message.GetID(), url.QueryEscape(react), "@me")).
		ContentType(JsonContentType).
		Expect(http.StatusNoContent).
		Send(c)
}

func (c *Client) DeleteOwnReaction(ctx context.Context, channel, message objects.SnowflakeObject, emoji interface{}) error {
	react, err := c.getEmoji(emoji)
	if err != nil {
		return err
	}

	return NewRequest().
		Method(http.MethodDelete).
		WithContext(ctx).
		Path(fmt.Sprintf(ReactionFmt, channel.GetID(), message.GetID(), url.QueryEscape(react), "@me")).
		ContentType(JsonContentType).
		Expect(http.StatusNoContent).
		Send(c)
}

func (c *Client) DeleteUserReaction(ctx context.Context, channel, message, user objects.SnowflakeObject, emoji interface{}) error {
	react, err := c.getEmoji(emoji)
	if err != nil {
		return err
	}

	return NewRequest().
		Method(http.MethodDelete).
		WithContext(ctx).
		Path(fmt.Sprintf(ReactionUserFmt, channel.GetID(), message.GetID(), url.QueryEscape(react), user.GetID())).
		ContentType(JsonContentType).
		Expect(http.StatusNoContent).
		Send(c)
}

type GetReactionsParams struct {
	Before objects.Snowflake `url:"before,omitempty"`
	After  objects.Snowflake `url:"after,omitempty"`
	Limit  int               `url:"limit"`
}

func (c *Client) GetReactions(ctx context.Context, channel, message objects.SnowflakeObject, emoji interface{}, params *GetReactionsParams) ([]*objects.User, error) {
	react, err := c.getEmoji(emoji)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(fmt.Sprintf(ReactionsFmt, channel.GetID(), message.GetID(), url.QueryEscape(react)))
	if err != nil {
		return nil, err
	}

	q, err := query.Values(params)
	if err != nil {
		return nil, err
	}
	u.RawQuery = q.Encode()

	var users []*objects.User
	err = NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(u.String()).
		ContentType(JsonContentType).
		Bind(&users).
		Expect(http.StatusOK).
		Send(c)
	return users, err
}

func (c *Client) DeleteAllReactions(ctx context.Context, channel, message objects.SnowflakeObject) error {
	return NewRequest().
		Method(http.MethodDelete).
		WithContext(ctx).
		Path(fmt.Sprintf(ReactionsBaseFmt, channel.GetID(), message.GetID())).
		ContentType(JsonContentType).
		Expect(http.StatusNoContent).
		Send(c)
}

func (c *Client) DeleteEmojiReactions(ctx context.Context, channel, message objects.SnowflakeObject, emoji interface{}) error {
	reaction, err := c.getEmoji(emoji)
	if err != nil {
		return err
	}

	return NewRequest().
		Method(http.MethodDelete).
		WithContext(ctx).
		Path(fmt.Sprintf(ReactionsFmt, channel.GetID(), message.GetID(), reaction)).
		ContentType(JsonContentType).
		Expect(http.StatusNoContent).
		Send(c)
}

func (c *Client) GetPinnedMessages(ctx context.Context, channel objects.SnowflakeObject) ([]*objects.Message, error) {
	var messages []*objects.Message
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(ChannelPinsFmt, channel.GetID())).
		ContentType(JsonContentType).
		Bind(&messages).
		Expect(http.StatusOK).
		Send(c)

	return messages, err
}

func (c *Client) AddPinnedMessage(ctx context.Context, channel, message objects.SnowflakeObject) error {
	return NewRequest().
		Method(http.MethodPut).
		WithContext(ctx).
		Path(fmt.Sprintf(ChannelPinnedFmt, channel.GetID(), message.GetID())).
		ContentType(JsonContentType).
		Expect(http.StatusNoContent).
		Send(c)
}

func (c *Client) DeletePinnedMessage(ctx context.Context, channel, message objects.SnowflakeObject) error {
	return NewRequest().
		Method(http.MethodDelete).
		WithContext(ctx).
		Path(fmt.Sprintf(ChannelPinnedFmt, channel.GetID(), message.GetID())).
		ContentType(JsonContentType).
		Expect(http.StatusNoContent).
		Send(c)
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

func (c *Client) CreateMessage(ctx context.Context, channel objects.SnowflakeObject, params *CreateMessageParams) (*objects.Message, error) {
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
	msg := &objects.Message{}
	err := NewRequest().
		Method(http.MethodPost).
		WithContext(ctx).
		Path(fmt.Sprintf(ChannelMessagesFmt, channel.GetID())).
		ContentType(contentType).
		Body(body).
		Bind(msg).
		Expect(http.StatusOK).
		Send(c)

	return msg, err
}

type EditMessageParams struct {
	Content         string                   `json:"content"`
	Embed           *objects.Embed           `json:"embed"`
	Flags           objects.MessageFlag      `json:"flags,omitempty"`
	AllowedMentions *objects.AllowedMentions `json:"allowed_mentions,omitempty"`
	Components      []*objects.Component     `json:"components"`
}

func (c *Client) EditMessage(ctx context.Context, channel, message objects.SnowflakeObject, params *EditMessageParams) (*objects.Message, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	msg := &objects.Message{}
	err = NewRequest().
		Method(http.MethodPatch).
		WithContext(ctx).
		Path(fmt.Sprintf(ChannelMessageFmt, channel.GetID(), message.GetID())).
		ContentType(JsonContentType).
		Body(body).
		Bind(msg).
		Expect(http.StatusOK).
		Send(c)
	return msg, err
}

func (c *Client) FollowNewsChannel(ctx context.Context, channel objects.SnowflakeObject) (*objects.FollowedChannel, error) {
	followedChannel := &objects.FollowedChannel{}
	err := NewRequest().
		Method(http.MethodPost).
		WithContext(ctx).
		Path(fmt.Sprintf(ChannelFollowersFmt, channel.GetID())).
		ContentType(JsonContentType).
		Bind(followedChannel).
		Expect(http.StatusOK).
		Send(c)

	return followedChannel, err
}

func (c *Client) StartTyping(ctx context.Context, channel objects.SnowflakeObject) error {
	return NewRequest().
		Method(http.MethodPost).
		WithContext(ctx).
		Path(fmt.Sprintf(ChannelTypingFmt, channel.GetID())).
		ContentType(JsonContentType).
		Expect(http.StatusNoContent).
		Send(c)
}

type StartThreadParams struct {
	Name                string `json:"name"`
	AutoArchiveDuration int    `json:"auto_archive_duration"`
	Type                int    `json:"type,omitempty"`
	Invitable           bool   `json:"invitable,omitempty"`
}

func (c *Client) StartThreadWithMessage(ctx context.Context, channel, message objects.SnowflakeObject, params *StartThreadParams) (*objects.Channel, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	thread := &objects.Channel{}
	err = NewRequest().
		Method(http.MethodPost).
		WithContext(ctx).
		Path(fmt.Sprintf(ChannelMessageThreadsFmt, channel.GetID(), message.GetID())).
		ContentType(JsonContentType).
		Body(body).
		Bind(thread).
		Expect(http.StatusOK, http.StatusCreated).
		Send(c)
	return thread, err
}

func (c *Client) StartThread(ctx context.Context, channel objects.SnowflakeObject, params *StartThreadParams) (*objects.Channel, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	thread := &objects.Channel{}
	err = NewRequest().
		Method(http.MethodPost).
		WithContext(ctx).
		Path(fmt.Sprintf(ChannelThreadsFmt, channel.GetID())).
		ContentType(JsonContentType).
		Body(body).
		Bind(thread).
		Expect(http.StatusOK, http.StatusCreated).
		Send(c)
	return thread, err
}

func (c *Client) JoinThread(ctx context.Context, thread objects.SnowflakeObject) error {
	return NewRequest().
		Method(http.MethodPost).
		WithContext(ctx).
		Path(fmt.Sprintf(ChannelThreadMembersMeFmt, thread.GetID())).
		ContentType(JsonContentType).
		Expect(http.StatusNoContent).
		Send(c)
}

func (c *Client) AddThreadMember(ctx context.Context, thread, user objects.SnowflakeObject) error {
	return NewRequest().
		Method(http.MethodPost).
		WithContext(ctx).
		Path(fmt.Sprintf(ChannelThreadMembersUserFmt, thread.GetID(), user.GetID())).
		ContentType(JsonContentType).
		Expect(http.StatusNoContent).
		Send(c)
}

func (c *Client) LeaveThread(ctx context.Context, thread objects.SnowflakeObject) error {
	return NewRequest().
		Method(http.MethodDelete).
		WithContext(ctx).
		Path(fmt.Sprintf(ChannelThreadMembersMeFmt, thread.GetID())).
		ContentType(JsonContentType).
		Expect(http.StatusNoContent).
		Send(c)
}

func (c *Client) RemoveThreadMember(ctx context.Context, thread, user objects.SnowflakeObject) error {
	return NewRequest().
		Method(http.MethodDelete).
		WithContext(ctx).
		Path(fmt.Sprintf(ChannelThreadMembersUserFmt, thread.GetID(), user.GetID())).
		ContentType(JsonContentType).
		Expect(http.StatusNoContent).
		Send(c)
}

func (c *Client) ListThreadMembers(ctx context.Context, thread objects.SnowflakeObject) ([]*objects.ThreadMember, error) {
	members := []*objects.ThreadMember{}
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(ChannelThreadMembersFmt, thread.GetID())).
		Bind(&members).
		Expect(http.StatusOK).
		Send(c)
	return members, err
}

type ListThreadsResponse struct {
	Threads []*objects.Channel      `json:"threads"`
	Members []*objects.ThreadMember `json:"members"`
	HasMore bool                    `json:"has_more"`
}

type ListThreadsParams struct {
	Before objects.Time `json:"before,omitempty"`
	Limit  int          `json:"limit,omitempty"`
}

func (c *Client) ListPublicArchivedThreads(ctx context.Context, channel objects.SnowflakeObject, params ...*ListThreadsParams) (*ListThreadsResponse, error) {
	u, err := url.Parse(fmt.Sprintf(ChannelThreadsArchivedPublicFmt, channel.GetID()))
	if err != nil {
		return nil, err
	}
	q, err := query.Values(params)
	if err != nil {
		return nil, err
	}
	u.RawQuery = q.Encode()
	threads := &ListThreadsResponse{}
	err = NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(u.String()).
		Bind(threads).
		Expect(http.StatusOK).
		Send(c)
	return threads, err
}

func (c *Client) ListPrivateArchivedThreads(ctx context.Context, channel objects.SnowflakeObject, params ...*ListThreadsParams) (*ListThreadsResponse, error) {
	u, err := url.Parse(fmt.Sprintf(ChannelThreadsArchivedPrivateFmt, channel.GetID()))
	if err != nil {
		return nil, err
	}
	q, err := query.Values(params)
	if err != nil {
		return nil, err
	}
	u.RawQuery = q.Encode()
	threads := &ListThreadsResponse{}
	err = NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(u.String()).
		Bind(threads).
		Expect(http.StatusOK).
		Send(c)
	return threads, err
}

func (c *Client) ListJoinedPrivateArchivedThreads(ctx context.Context, channel objects.SnowflakeObject, params ...*ListThreadsParams) (*ListThreadsResponse, error) {
	u, err := url.Parse(fmt.Sprintf(ChannelUsersMeThreadsArchivedFmt, channel.GetID()))
	if err != nil {
		return nil, err
	}
	q, err := query.Values(params)
	if err != nil {
		return nil, err
	}
	u.RawQuery = q.Encode()
	threads := &ListThreadsResponse{}
	err = NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(u.String()).
		Bind(threads).
		Expect(http.StatusOK).
		Send(c)
	return threads, err
}
