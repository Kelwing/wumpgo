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

func (c *Client) GetChannel(id objects.Snowflake) (*objects.Channel, error) {
	channel := &objects.Channel{}

	err := NewRequest().
		Method(http.MethodGet).
		Path(fmt.Sprintf(ChannelBaseFmt, id)).
		ContentType(JsonContentType).
		Bind(channel).
		Expect(http.StatusOK).
		Send(c)

	return channel, err
}

type ModifyChannelParams struct {
	Name                 *string                       `json:"name,omitempty"`
	Type                 *objects.ChannelType          `json:"type,omitempty"`
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
	channel := &objects.Channel{}

	err = NewRequest().
		Method(http.MethodPatch).
		Path(fmt.Sprintf(ChannelBaseFmt, id)).
		ContentType(JsonContentType).
		Body(data).
		Reason(reason).
		Expect(http.StatusOK).
		Bind(channel).
		Send(c)

	return channel, err
}

func (c *Client) DeleteChannel(id objects.Snowflake, reason string) (*objects.Channel, error) {
	channel := &objects.Channel{}
	err := NewRequest().
		Method(http.MethodDelete).
		Path(fmt.Sprintf(ChannelBaseFmt, id)).
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

	var messages []*objects.Message
	err = NewRequest().
		Method(http.MethodGet).
		Path(u.String()).
		ContentType(JsonContentType).
		Bind(&messages).
		Expect(http.StatusOK).
		Send(c)

	return messages, err
}

func (c *Client) GetChannelMessage(channel, message objects.Snowflake) (*objects.Message, error) {
	msg := &objects.Message{}
	err := NewRequest().
		Method(http.MethodGet).
		Path(fmt.Sprintf(ChannelMessageFmt, channel, message)).
		ContentType(JsonContentType).
		Bind(msg).
		Expect(http.StatusOK).
		Send(c)

	return msg, err
}

func (c *Client) CrossPostMessage(channel, message objects.Snowflake) (*objects.Message, error) {
	msg := &objects.Message{}
	err := NewRequest().
		Method(http.MethodPost).
		Path(fmt.Sprintf(CrosspostMessageFmt, channel, message)).
		ContentType(JsonContentType).
		Bind(msg).
		Expect(http.StatusOK).
		Send(c)

	return msg, err
}

func (c *Client) DeleteMessage(channel, message objects.Snowflake) error {
	return NewRequest().
		Method(http.MethodDelete).
		Path(fmt.Sprintf(ChannelMessageFmt, channel, message)).
		ContentType(JsonContentType).
		Expect(http.StatusNoContent).
		Send(c)
}

type DeleteMessagesParams struct {
	Messages []objects.Snowflake `json:"messages"`
}

func (c *Client) BulkDeleteMessages(channel objects.Snowflake, params *DeleteMessagesParams) error {
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}

	return NewRequest().
		Method(http.MethodPost).
		Path(fmt.Sprintf(BulkDeleteMessagesFmt, channel)).
		ContentType(JsonContentType).
		Body(data).
		Expect(http.StatusNoContent).
		Send(c)
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

	return NewRequest().
		Method(http.MethodPut).
		Path(fmt.Sprintf(ChannelPermissionsFmt, channel, overwrite)).
		ContentType(JsonContentType).
		Body(data).
		Reason(reason).
		Expect(http.StatusNoContent).
		Send(c)
}

func (c *Client) DeleteChannelPermission(channel, overwrite objects.Snowflake, reason string) error {
	return NewRequest().
		Method(http.MethodDelete).
		Path(fmt.Sprintf(ChannelPermissionsFmt, channel, overwrite)).
		Reason(reason).
		ContentType(JsonContentType).
		Expect(http.StatusNoContent).
		Send(c)
}

func (c *Client) GetChannelInvites(channel objects.Snowflake) ([]*objects.Invite, error) {
	var invites []*objects.Invite
	err := NewRequest().
		Method(http.MethodGet).
		Path(fmt.Sprintf(ChannelInvitesFmt, channel)).
		ContentType(JsonContentType).
		Bind(&invites).
		Expect(http.StatusOK).
		Send(c)

	return invites, err
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

	invite := &objects.Invite{}

	err = NewRequest().
		Method(http.MethodPost).
		Path(fmt.Sprintf(ChannelInvitesFmt, channel)).
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

func (c *Client) CreateReaction(channel, message objects.Snowflake, emoji interface{}) error {
	react, err := c.getEmoji(emoji)
	if err != nil {
		return err
	}

	return NewRequest().
		Method(http.MethodPut).
		Path(fmt.Sprintf(ReactionFmt, channel, message, url.QueryEscape(react), "@me")).
		ContentType(JsonContentType).
		Expect(http.StatusNoContent).
		Send(c)
}

func (c *Client) DeleteOwnReaction(channel, message objects.Snowflake, emoji interface{}) error {
	react, err := c.getEmoji(emoji)
	if err != nil {
		return err
	}

	return NewRequest().
		Method(http.MethodDelete).
		Path(fmt.Sprintf(ReactionFmt, channel, message, url.QueryEscape(react), "@me")).
		ContentType(JsonContentType).
		Expect(http.StatusNoContent).
		Send(c)
}

func (c *Client) DeleteUserReaction(channel, message, user objects.Snowflake, emoji interface{}) error {
	react, err := c.getEmoji(emoji)
	if err != nil {
		return err
	}

	return NewRequest().
		Method(http.MethodDelete).
		Path(fmt.Sprintf(ReactionUserFmt, channel, message, url.QueryEscape(react), user)).
		ContentType(JsonContentType).
		Expect(http.StatusNoContent).
		Send(c)
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

	var users []*objects.User
	err = NewRequest().
		Method(http.MethodGet).
		Path(u.String()).
		ContentType(JsonContentType).
		Bind(&users).
		Expect(http.StatusOK).
		Send(c)
	return users, err
}

func (c *Client) DeleteAllReactions(channel, message objects.Snowflake) error {
	return NewRequest().
		Method(http.MethodDelete).
		Path(fmt.Sprintf(ReactionsBaseFmt, channel, message)).
		ContentType(JsonContentType).
		Expect(http.StatusNoContent).
		Send(c)
}

func (c *Client) DeleteEmojiReactions(channel, message objects.Snowflake, emoji interface{}) error {
	reaction, err := c.getEmoji(emoji)
	if err != nil {
		return err
	}

	return NewRequest().
		Method(http.MethodDelete).
		Path(fmt.Sprintf(ReactionsFmt, channel, message, reaction)).
		ContentType(JsonContentType).
		Expect(http.StatusNoContent).
		Send(c)
}

func (c *Client) GetPinnedMessages(channel objects.Snowflake) ([]*objects.Message, error) {
	var messages []*objects.Message
	err := NewRequest().
		Method(http.MethodGet).
		Path(fmt.Sprintf(ChannelPinsFmt, channel)).
		ContentType(JsonContentType).
		Bind(&messages).
		Expect(http.StatusOK).
		Send(c)

	return messages, err
}

func (c *Client) AddPinnedMessage(channel, message objects.Snowflake) error {
	return NewRequest().
		Method(http.MethodPut).
		Path(fmt.Sprintf(ChannelPinnedFmt, channel, message)).
		ContentType(JsonContentType).
		Expect(http.StatusNoContent).
		Send(c)
}

func (c *Client) DeletePinnedMessage(channel, message objects.Snowflake) error {
	return NewRequest().
		Method(http.MethodDelete).
		Path(fmt.Sprintf(ChannelPinnedFmt, channel, message)).
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
	msg := &objects.Message{}
	err := NewRequest().
		Method(http.MethodPost).
		Path(fmt.Sprintf(ChannelMessagesFmt, channel)).
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

func (c *Client) EditMessage(channel, message objects.Snowflake, params *EditMessageParams) (*objects.Message, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	msg := &objects.Message{}
	err = NewRequest().
		Method(http.MethodPatch).
		Path(fmt.Sprintf(ChannelMessageFmt, channel, message)).
		ContentType(JsonContentType).
		Body(body).
		Bind(msg).
		Expect(http.StatusOK).
		Send(c)
	return msg, err
}

func (c *Client) FollowNewsChannel(channel objects.Snowflake) (*objects.FollowedChannel, error) {
	followedChannel := &objects.FollowedChannel{}
	err := NewRequest().
		Method(http.MethodPost).
		Path(fmt.Sprintf(ChannelFollowersFmt, channel)).
		ContentType(JsonContentType).
		Bind(followedChannel).
		Expect(http.StatusOK).
		Send(c)

	return followedChannel, err
}

func (c *Client) StartTyping(channel objects.Snowflake) error {
	return NewRequest().
		Method(http.MethodPost).
		Path(fmt.Sprintf(ChannelTypingFmt, channel)).
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

func (c *Client) StartThreadWithMessage(channel objects.Snowflake, message objects.Snowflake, params *StartThreadParams) (*objects.Channel, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	thread := &objects.Channel{}
	err = NewRequest().
		Method(http.MethodPost).
		Path(fmt.Sprintf(ChannelMessageThreadsFmt, channel, message)).
		ContentType(JsonContentType).
		Body(body).
		Bind(thread).
		Expect(http.StatusOK, http.StatusCreated).
		Send(c)
	return thread, err
}

func (c *Client) StartThread(channel objects.Snowflake, params *StartThreadParams) (*objects.Channel, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	thread := &objects.Channel{}
	err = NewRequest().
		Method(http.MethodPost).
		Path(fmt.Sprintf(ChannelThreadsFmt, channel)).
		ContentType(JsonContentType).
		Body(body).
		Bind(thread).
		Expect(http.StatusOK, http.StatusCreated).
		Send(c)
	return thread, err
}

func (c *Client) JoinThread(thread objects.Snowflake) error {
	return NewRequest().
		Method(http.MethodPost).
		Path(fmt.Sprintf(ChannelThreadMembersMeFmt, thread)).
		ContentType(JsonContentType).
		Expect(http.StatusNoContent).
		Send(c)
}

func (c *Client) AddThreadMember(thread, user objects.Snowflake) error {
	return NewRequest().
		Method(http.MethodPost).
		Path(fmt.Sprintf(ChannelThreadMembersUserFmt, thread, user)).
		ContentType(JsonContentType).
		Expect(http.StatusNoContent).
		Send(c)
}

func (c *Client) LeaveThread(thread objects.Snowflake) error {
	return NewRequest().
		Method(http.MethodDelete).
		Path(fmt.Sprintf(ChannelThreadMembersMeFmt, thread)).
		ContentType(JsonContentType).
		Expect(http.StatusNoContent).
		Send(c)
}

func (c *Client) RemoveThreadMember(thread, user objects.Snowflake) error {
	return NewRequest().
		Method(http.MethodDelete).
		Path(fmt.Sprintf(ChannelThreadMembersUserFmt, thread, user)).
		ContentType(JsonContentType).
		Expect(http.StatusNoContent).
		Send(c)
}

func (c *Client) ListThreadMembers(thread objects.Snowflake) ([]*objects.ThreadMember, error) {
	members := []*objects.ThreadMember{}
	err := NewRequest().
		Method(http.MethodGet).
		Path(fmt.Sprintf(ChannelThreadMembersFmt, thread)).
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

func (c *Client) ListPublicArchivedThreads(channel objects.Snowflake, params ...*ListThreadsParams) (*ListThreadsResponse, error) {
	u, err := url.Parse(fmt.Sprintf(ChannelThreadsArchivedPublicFmt, channel))
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
		Path(u.String()).
		Bind(threads).
		Expect(http.StatusOK).
		Send(c)
	return threads, err
}

func (c *Client) ListPrivateArchivedThreads(channel objects.Snowflake, params ...*ListThreadsParams) (*ListThreadsResponse, error) {
	u, err := url.Parse(fmt.Sprintf(ChannelThreadsArchivedPrivateFmt, channel))
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
		Path(u.String()).
		Bind(threads).
		Expect(http.StatusOK).
		Send(c)
	return threads, err
}

func (c *Client) ListJoinedPrivateArchivedThreads(channel objects.Snowflake, params ...*ListThreadsParams) (*ListThreadsResponse, error) {
	u, err := url.Parse(fmt.Sprintf(ChannelUsersMeThreadsArchivedFmt, channel))
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
		Path(u.String()).
		Bind(threads).
		Expect(http.StatusOK).
		Send(c)
	return threads, err
}
