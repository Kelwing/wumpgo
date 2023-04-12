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
	"wumpgo.dev/wumpgo/objects/permissions"
)

func (c *Client) GetChannel(ctx context.Context, id objects.Snowflake) (*objects.Channel, error) {
	channel := &objects.Channel{}

	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(ChannelBaseFmt, id)).
		ContentType(JsonContentType).
		Bind(channel).
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

func (c *Client) ModifyChannel(ctx context.Context, id objects.Snowflake, params *ModifyChannelParams) (*objects.Channel, error) {
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
		Path(fmt.Sprintf(ChannelBaseFmt, id)).
		ContentType(JsonContentType).
		Body(data).
		Reason(reason).
		Bind(channel).
		Send(c)

	return channel, err
}

func (c *Client) DeleteChannel(ctx context.Context, id objects.Snowflake, reason string) (*objects.Channel, error) {
	channel := &objects.Channel{}
	err := NewRequest().
		Method(http.MethodDelete).
		WithContext(ctx).
		Path(fmt.Sprintf(ChannelBaseFmt, id)).
		Reason(reason).
		ContentType(JsonContentType).
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

func (c *Client) GetChannelMessages(ctx context.Context, id objects.Snowflake, params *GetChannelMessagesParams) ([]*objects.Message, error) {
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
		WithContext(ctx).
		Path(u.String()).
		ContentType(JsonContentType).
		Bind(&messages).
		Send(c)

	return messages, err
}

func (c *Client) GetChannelMessage(ctx context.Context, channel, message objects.Snowflake) (*objects.Message, error) {
	msg := &objects.Message{}
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(ChannelMessageFmt, channel, message)).
		ContentType(JsonContentType).
		Bind(msg).
		Send(c)

	return msg, err
}

func (c *Client) CrossPostMessage(ctx context.Context, channel, message objects.Snowflake) (*objects.Message, error) {
	msg := &objects.Message{}
	err := NewRequest().
		Method(http.MethodPost).
		WithContext(ctx).
		Path(fmt.Sprintf(CrosspostMessageFmt, channel, message)).
		ContentType(JsonContentType).
		Bind(msg).
		Send(c)

	return msg, err
}

func (c *Client) DeleteMessage(ctx context.Context, channel, message objects.Snowflake) error {
	return NewRequest().
		Method(http.MethodDelete).
		WithContext(ctx).
		Path(fmt.Sprintf(ChannelMessageFmt, channel, message)).
		ContentType(JsonContentType).
		Send(c)
}

type DeleteMessagesParams struct {
	Messages []objects.Snowflake `json:"messages"`
}

func (c *Client) BulkDeleteMessages(ctx context.Context, channel objects.Snowflake, params *DeleteMessagesParams) error {
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}

	return NewRequest().
		Method(http.MethodPost).
		WithContext(ctx).
		Path(fmt.Sprintf(BulkDeleteMessagesFmt, channel)).
		ContentType(JsonContentType).
		Body(data).
		Send(c)
}

type EditChannelParams struct {
	Allow  permissions.PermissionBit `json:"allow"`
	Deny   permissions.PermissionBit `json:"deny"`
	Type   int                       `json:"type"`
	Reason string                    `json:"-"`
}

func (c *Client) EditChannelPermissions(ctx context.Context, channel, overwrite objects.Snowflake, params *EditChannelParams) error {
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
		Path(fmt.Sprintf(ChannelPermissionsFmt, channel, overwrite)).
		ContentType(JsonContentType).
		Body(data).
		Reason(reason).
		Send(c)
}

func (c *Client) DeleteChannelPermission(ctx context.Context, channel, overwrite objects.Snowflake, reason string) error {
	return NewRequest().
		Method(http.MethodDelete).
		WithContext(ctx).
		Path(fmt.Sprintf(ChannelPermissionsFmt, channel, overwrite)).
		Reason(reason).
		ContentType(JsonContentType).
		Send(c)
}

func (c *Client) GetChannelInvites(ctx context.Context, channel objects.Snowflake) ([]*objects.Invite, error) {
	var invites []*objects.Invite
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(ChannelInvitesFmt, channel)).
		ContentType(JsonContentType).
		Bind(&invites).
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

func (c *Client) CreateChannelInvite(ctx context.Context, channel objects.Snowflake, params *CreateInviteParams) (*objects.Invite, error) {
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
		Path(fmt.Sprintf(ChannelInvitesFmt, channel)).
		ContentType(JsonContentType).
		Body(data).
		Reason(reason).
		Bind(invite).
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

func (c *Client) CreateReaction(ctx context.Context, channel, message objects.Snowflake, emoji interface{}) error {
	react, err := c.getEmoji(emoji)
	if err != nil {
		return err
	}

	return NewRequest().
		Method(http.MethodPut).
		WithContext(ctx).
		Path(fmt.Sprintf(ReactionFmt, channel, message, url.QueryEscape(react), "@me")).
		ContentType(JsonContentType).
		Send(c)
}

func (c *Client) DeleteOwnReaction(ctx context.Context, channel, message objects.Snowflake, emoji interface{}) error {
	react, err := c.getEmoji(emoji)
	if err != nil {
		return err
	}

	return NewRequest().
		Method(http.MethodDelete).
		WithContext(ctx).
		Path(fmt.Sprintf(ReactionFmt, channel, message, url.QueryEscape(react), "@me")).
		ContentType(JsonContentType).
		Send(c)
}

func (c *Client) DeleteUserReaction(ctx context.Context, channel, message, user objects.Snowflake, emoji interface{}) error {
	react, err := c.getEmoji(emoji)
	if err != nil {
		return err
	}

	return NewRequest().
		Method(http.MethodDelete).
		WithContext(ctx).
		Path(fmt.Sprintf(ReactionUserFmt, channel, message, url.QueryEscape(react), user)).
		ContentType(JsonContentType).
		Send(c)
}

type GetReactionsParams struct {
	Before objects.Snowflake `url:"before,omitempty"`
	After  objects.Snowflake `url:"after,omitempty"`
	Limit  int               `url:"limit"`
}

func (c *Client) GetReactions(ctx context.Context, channel, message objects.Snowflake, emoji interface{}, params *GetReactionsParams) ([]*objects.User, error) {
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
		WithContext(ctx).
		Path(u.String()).
		ContentType(JsonContentType).
		Bind(&users).
		Send(c)
	return users, err
}

func (c *Client) DeleteAllReactions(ctx context.Context, channel, message objects.Snowflake) error {
	return NewRequest().
		Method(http.MethodDelete).
		WithContext(ctx).
		Path(fmt.Sprintf(ReactionsBaseFmt, channel, message)).
		ContentType(JsonContentType).
		Send(c)
}

func (c *Client) DeleteEmojiReactions(ctx context.Context, channel, message objects.Snowflake, emoji interface{}) error {
	reaction, err := c.getEmoji(emoji)
	if err != nil {
		return err
	}

	return NewRequest().
		Method(http.MethodDelete).
		WithContext(ctx).
		Path(fmt.Sprintf(ReactionsFmt, channel, message, reaction)).
		ContentType(JsonContentType).
		Send(c)
}

func (c *Client) GetPinnedMessages(ctx context.Context, channel objects.Snowflake) ([]*objects.Message, error) {
	var messages []*objects.Message
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(ChannelPinsFmt, channel)).
		ContentType(JsonContentType).
		Bind(&messages).
		Send(c)

	return messages, err
}

func (c *Client) AddPinnedMessage(ctx context.Context, channel, message objects.Snowflake) error {
	return NewRequest().
		Method(http.MethodPut).
		WithContext(ctx).
		Path(fmt.Sprintf(ChannelPinnedFmt, channel, message)).
		ContentType(JsonContentType).
		Send(c)
}

func (c *Client) DeletePinnedMessage(ctx context.Context, channel, message objects.Snowflake) error {
	return NewRequest().
		Method(http.MethodDelete).
		WithContext(ctx).
		Path(fmt.Sprintf(ChannelPinnedFmt, channel, message)).
		ContentType(JsonContentType).
		Send(c)
}

type CreateMessageParams struct {
	Content          string                    `json:"content,omitempty"`
	TTS              bool                      `json:"tts,omitempty"`
	Embeds           []*objects.Embed          `json:"embeds,omitempty"`
	AllowedMentions  *objects.AllowedMentions  `json:"allowed_mentions,omitempty"`
	MessageReference *objects.MessageReference `json:"message_reference,omitempty"`
	Components       []*objects.Component      `json:"components,omitempty"`
	StickerIDs       []objects.Snowflake       `json:"sticker_ids,omitempty"`
	Attachments      []*objects.Attachment     `json:"attachments,omitempty"`
	Flags            objects.MessageFlag       `json:"flags,omitempty"`
	Files            []*objects.DiscordFile    `json:"-"`
}

func (c *Client) CreateMessage(ctx context.Context, channel objects.Snowflake, params *CreateMessageParams) (*objects.Message, error) {
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
		body, err = json.Marshal(params)
		if err != nil {
			return nil, err
		}
	}
	msg := &objects.Message{}
	err := NewRequest().
		Method(http.MethodPost).
		WithContext(ctx).
		Path(fmt.Sprintf(ChannelMessagesFmt, channel)).
		ContentType(contentType).
		Body(body).
		Bind(msg).
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

func (c *Client) EditMessage(ctx context.Context, channel, message objects.Snowflake, params *EditMessageParams) (*objects.Message, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	msg := &objects.Message{}
	err = NewRequest().
		Method(http.MethodPatch).
		WithContext(ctx).
		Path(fmt.Sprintf(ChannelMessageFmt, channel, message)).
		ContentType(JsonContentType).
		Body(body).
		Bind(msg).
		Send(c)
	return msg, err
}

func (c *Client) FollowNewsChannel(ctx context.Context, channel objects.Snowflake) (*objects.FollowedChannel, error) {
	followedChannel := &objects.FollowedChannel{}
	err := NewRequest().
		Method(http.MethodPost).
		WithContext(ctx).
		Path(fmt.Sprintf(ChannelFollowersFmt, channel)).
		ContentType(JsonContentType).
		Bind(followedChannel).
		Send(c)

	return followedChannel, err
}

func (c *Client) StartTyping(ctx context.Context, channel objects.Snowflake) error {
	return NewRequest().
		Method(http.MethodPost).
		WithContext(ctx).
		Path(fmt.Sprintf(ChannelTypingFmt, channel)).
		ContentType(JsonContentType).
		Send(c)
}

type StartThreadParams struct {
	Name                string `json:"name"`
	AutoArchiveDuration int    `json:"auto_archive_duration"`
	Type                int    `json:"type,omitempty"`
	Invitable           bool   `json:"invitable,omitempty"`
}

func (c *Client) StartThreadWithMessage(ctx context.Context, channel, message objects.Snowflake, params *StartThreadParams) (*objects.Channel, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	thread := &objects.Channel{}
	err = NewRequest().
		Method(http.MethodPost).
		WithContext(ctx).
		Path(fmt.Sprintf(ChannelMessageThreadsFmt, channel, message)).
		ContentType(JsonContentType).
		Body(body).
		Bind(thread).
		Send(c)
	return thread, err
}

func (c *Client) StartThread(ctx context.Context, channel objects.Snowflake, params *StartThreadParams) (*objects.Channel, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	thread := &objects.Channel{}
	err = NewRequest().
		Method(http.MethodPost).
		WithContext(ctx).
		Path(fmt.Sprintf(ChannelThreadsFmt, channel)).
		ContentType(JsonContentType).
		Body(body).
		Bind(thread).
		Send(c)
	return thread, err
}

type StartThreadInForumChannelParams struct {
	Name                string                    `json:"name"`
	AutoArchiveDuration uint64                    `json:"auto_archive_duration"`
	RateLimitPerUser    uint64                    `json:"rate_limit_per_user"`
	Message             *ForumThreadMessageParams `json:"message"`
	AppliedTags         []objects.Snowflake       `json:"applied_tags"`
	Reason              string                    `json:"-"`
}

type ForumThreadMessageParams struct {
	Content         string                  `json:"content,omitempty"`
	Embeds          []*objects.Embed        `json:"embeds,omitempty"`
	AllowedMentions objects.AllowedMentions `json:"allowed_mentions,omitempty"`
	Components      []*objects.Component    `json:"components,omitempty"`
	StickerIDs      []objects.Snowflake     `json:"sticker_ids,omitempty"`
	Files           []*objects.DiscordFile  `json:"files,omitempty"`
	Attachments     []*objects.Attachment   `json:"attachments,omitempty"`
	Flags           objects.MessageFlag     `json:"flags,omitempty"`
}

type ForumThreadChannel struct {
	*objects.Channel
	Message *objects.Message
}

func (c *Client) StartThreadInForumChannel(ctx context.Context, channel objects.Snowflake, params *StartThreadInForumChannelParams) (*objects.ForumThreadChannel, error) {
	var contentType string
	var body []byte

	reason := ""
	if params != nil {
		reason = params.Reason
	}

	if len(params.Message.Files) > 0 {
		buffer := new(bytes.Buffer)
		m := multipart.NewWriter(buffer)

		for n, file := range params.Message.Files {
			a, err := file.GenerateAttachment(objects.Snowflake(n+1), m)
			if err != nil {
				continue
			}
			params.Message.Attachments = append(params.Message.Attachments, a)
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
		body, err = json.Marshal(params)
		if err != nil {
			return nil, err
		}
	}
	ch := &objects.ForumThreadChannel{}
	err := NewRequest().
		Method(http.MethodPost).
		WithContext(ctx).
		Path(fmt.Sprintf(ChannelMessagesFmt, channel)).
		ContentType(contentType).
		Body(body).
		Reason(reason).
		Bind(ch).
		Send(c)

	return ch, err
}

func (c *Client) JoinThread(ctx context.Context, thread objects.Snowflake) error {
	return NewRequest().
		Method(http.MethodPost).
		WithContext(ctx).
		Path(fmt.Sprintf(ChannelThreadMembersMeFmt, thread)).
		ContentType(JsonContentType).
		Send(c)
}

func (c *Client) AddThreadMember(ctx context.Context, thread, user objects.Snowflake) error {
	return NewRequest().
		Method(http.MethodPost).
		WithContext(ctx).
		Path(fmt.Sprintf(ChannelThreadMembersUserFmt, thread, user)).
		ContentType(JsonContentType).
		Send(c)
}

func (c *Client) LeaveThread(ctx context.Context, thread objects.Snowflake) error {
	return NewRequest().
		Method(http.MethodDelete).
		WithContext(ctx).
		Path(fmt.Sprintf(ChannelThreadMembersMeFmt, thread)).
		ContentType(JsonContentType).
		Send(c)
}

func (c *Client) RemoveThreadMember(ctx context.Context, thread, user objects.Snowflake) error {
	return NewRequest().
		Method(http.MethodDelete).
		WithContext(ctx).
		Path(fmt.Sprintf(ChannelThreadMembersUserFmt, thread, user)).
		ContentType(JsonContentType).
		Send(c)
}

func (c *Client) GetThreadMember(ctx context.Context, channel objects.Snowflake, user objects.Snowflake) (*objects.ThreadMember, error) {
	member := &objects.ThreadMember{}

	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(ChannelThreadMembersUserFmt, channel, user)).
		ContentType(JsonContentType).
		Bind(member).
		Send(c)

	return member, err
}

func (c *Client) ListThreadMembers(ctx context.Context, thread objects.Snowflake) ([]*objects.ThreadMember, error) {
	members := []*objects.ThreadMember{}
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(ChannelThreadMembersFmt, thread)).
		Bind(&members).
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

func (c *Client) ListPublicArchivedThreads(ctx context.Context, channel objects.Snowflake, params ...*ListThreadsParams) (*ListThreadsResponse, error) {
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
		WithContext(ctx).
		Path(u.String()).
		Bind(threads).
		Send(c)
	return threads, err
}

func (c *Client) ListPrivateArchivedThreads(ctx context.Context, channel objects.Snowflake, params ...*ListThreadsParams) (*ListThreadsResponse, error) {
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
		WithContext(ctx).
		Path(u.String()).
		Bind(threads).
		Send(c)
	return threads, err
}

func (c *Client) ListJoinedPrivateArchivedThreads(ctx context.Context, channel objects.Snowflake, params ...*ListThreadsParams) (*ListThreadsResponse, error) {
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
		WithContext(ctx).
		Path(u.String()).
		Bind(threads).
		Send(c)
	return threads, err
}
