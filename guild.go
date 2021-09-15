package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"net/http"
	"net/url"

	"github.com/Postcord/objects"
	"github.com/google/go-querystring/query"
)

type CreateGuildParams struct {
	Name                        string                              `json:"name"`
	Region                      string                              `json:"region,omitempty"`
	Icon                        string                              `json:"icon,omitempty"`
	VerificationLevel           *objects.VerificationLevel          `json:"verification_level,omitempty"`
	DefaultMessageNotifications *objects.MessageNotificationsLevel  `json:"default_message_notifications,omitempty"`
	ExplicitContentFilter       *objects.ExplicitContentFilterLevel `json:"explicit_content_filter,omitempty"`
	Roles                       []*objects.Role                     `json:"roles,omitempty"`
	Channels                    []*objects.Channel                  `json:"channels,omitempty"`
	AFKChannelID                objects.Snowflake                   `json:"afk_channel_id,omitempty"`
	AFKTimeout                  int64                               `json:"afk_timeout,omitempty"`
	SystemChannelID             objects.Snowflake                   `json:"system_channel_id,omitempty"`
	Reason                      string                              `json:"-"`
}

func (c *Client) CreateGuild(params *CreateGuildParams) (*objects.Guild, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	reason := ""
	if params != nil {
		reason = params.Reason
	}

	guild := &objects.Guild{}

	err = NewRequest().
		Method(http.MethodPost).
		Path(GuildCreateFmt).
		ContentType(JsonContentType).
		Body(data).
		Reason(reason).
		Expect(http.StatusOK).
		Bind(guild).
		Send(c)

	return guild, err
}

func (c *Client) GetGuild(id objects.Snowflake) (*objects.Guild, error) {
	guild := &objects.Guild{}
	err := NewRequest().
		Method(http.MethodGet).
		Path(fmt.Sprintf(GuildBaseFmt, id)).
		ContentType(JsonContentType).
		Expect(http.StatusOK).
		Bind(guild).
		Send(c)

	return guild, err
}

func (c *Client) GetGuildPreview(id objects.Snowflake) (*objects.GuildPreview, error) {
	preview := &objects.GuildPreview{}
	err := NewRequest().
		Method(http.MethodGet).
		Path(fmt.Sprintf(GuildPreviewFmt, id)).
		ContentType(JsonContentType).
		Expect(http.StatusOK).
		Bind(preview).
		Send(c)

	return preview, err
}

type ModifyGuildParams struct {
	Name                        string                              `json:"name,omitempty"`
	Region                      string                              `json:"region,omitempty"`
	VerificationLevel           *objects.VerificationLevel          `json:"verification_level,omitempty"`
	DefaultMessageNotifications *objects.MessageNotificationsLevel  `json:"default_message_notifications,omitempty"`
	ExplicitContentFilter       *objects.ExplicitContentFilterLevel `json:"explicit_content_filter,omitempty"`
	AFKChannelID                objects.Snowflake                   `json:"afk_channel_id,omitempty"`
	AFKTimeout                  int64                               `json:"afk_timeout,omitempty"`
	Icon                        string                              `json:"icon,omitempty"`
	OwnerID                     objects.Snowflake                   `json:"owner_id,omitempty"`
	Splash                      string                              `json:"splash,omitempty"`
	Banner                      string                              `json:"banner,omitempty"`
	SystemChannelID             objects.Snowflake                   `json:"system_channel_id,omitempty"`
	RulesChannelID              objects.Snowflake                   `json:"rules_channel_id,omitempty"`
	PublicUpdatesChannelID      objects.Snowflake                   `json:"public_updates_channel_id,omitempty"`
	PreferredLocale             string                              `json:"preferred_locale,omitempty"`
	Reason                      string                              `json:"-"`
}

func (c *Client) ModifyGuild(id objects.Snowflake, params *ModifyGuildParams) (*objects.Guild, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	reason := ""
	if params != nil {
		reason = params.Reason
	}

	guild := &objects.Guild{}
	err = NewRequest().
		Method(http.MethodPatch).
		Path(fmt.Sprintf(GuildBaseFmt, id)).
		ContentType(JsonContentType).
		Body(data).
		Reason(reason).
		Expect(http.StatusOK).
		Bind(guild).
		Send(c)

	return guild, err
}

func (c *Client) DeleteGuild(id objects.Snowflake) error {
	return NewRequest().
		Method(http.MethodDelete).
		Path(fmt.Sprintf(GuildBaseFmt, id)).
		ContentType(JsonContentType).
		Expect(http.StatusNoContent).
		Send(c)
}

func (c *Client) GetGuildChannels(id objects.Snowflake) ([]*objects.Channel, error) {
	channels := []*objects.Channel{}
	err := NewRequest().
		Method(http.MethodGet).
		Path(fmt.Sprintf(GuildChannelsFmt, id)).
		ContentType(JsonContentType).
		Expect(http.StatusOK).
		Bind(&channels).
		Send(c)

	return channels, err
}

type ChannelCreateParams struct {
	Name                 string                         `json:"name"`
	Type                 objects.ChannelType            `json:"type,omitempty"`
	Topic                string                         `json:"topic,omitempty"`
	Bitrate              int                            `json:"bitrate,omitempty"`
	UserLimit            int                            `json:"user_limit,omitempty"`
	RateLimitPerUser     int                            `json:"rate_limit_per_user,omitempty"`
	Position             int                            `json:"position,omitempty"`
	PermissionOverwrites []*objects.PermissionOverwrite `json:"permission_overwrites,omitempty"`
	ParentID             objects.Snowflake              `json:"parent_id,omitempty"`
	NSFW                 bool                           `json:"nsfw,omitempty"`
	Reason               string                         `json:"-"`
}

func (c *Client) CreateGuildChannel(id objects.Snowflake, params *ChannelCreateParams) (*objects.Channel, error) {
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
		Method(http.MethodPost).
		Path(fmt.Sprintf(GuildChannelsFmt, id)).
		ContentType(JsonContentType).
		Body(data).
		Reason(reason).
		Expect(http.StatusOK).
		Bind(channel).
		Send(c)

	return channel, err
}

type ModifyChannelPositionParams struct {
	ID       objects.Snowflake `json:"id"`
	Position int               `json:"position"`
}

func (c *Client) ModifyGuildChannelPositions(id objects.Snowflake, params []*ModifyChannelPositionParams, reason string) error {
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}

	return NewRequest().
		Method(http.MethodPatch).
		Path(fmt.Sprintf(GuildChannelsFmt, id)).
		ContentType(JsonContentType).
		Body(data).
		Reason(reason).
		Expect(http.StatusNoContent).
		Send(c)
}

func (c *Client) ListActiveThreads(id objects.Snowflake) ([]*ListThreadsResponse, error) {
	channels := []*ListThreadsResponse{}
	err := NewRequest().
		Method(http.MethodGet).
		Path(fmt.Sprintf(GuildThreadsFmt, id)).
		ContentType(JsonContentType).
		Expect(http.StatusOK).
		Bind(&channels).
		Send(c)

	return channels, err
}

func (c *Client) GetGuildMember(guild, user objects.Snowflake) (*objects.GuildMember, error) {
	member := &objects.GuildMember{}
	err := NewRequest().
		Method(http.MethodGet).
		Path(fmt.Sprintf(GuildMemberFmt, guild, user)).
		ContentType(JsonContentType).
		Expect(http.StatusOK).
		Bind(member).
		Send(c)

	return member, err
}

type ListGuildMembersParams struct {
	Limit int               `url:"limit,omitempty"`
	After objects.Snowflake `url:"after,omitempty"`
}

func (c *Client) ListGuildMembers(guild objects.Snowflake, params *ListGuildMembersParams) ([]*objects.GuildMember, error) {
	u, err := url.Parse(fmt.Sprintf(GuildMembersFmt, guild))
	if err != nil {
		return nil, err
	}
	v, err := query.Values(params)
	if err != nil {
		return nil, err
	}

	u.RawQuery = v.Encode()

	members := []*objects.GuildMember{}
	err = NewRequest().
		Method(http.MethodGet).
		Path(u.String()).
		ContentType(JsonContentType).
		Expect(http.StatusOK).
		Bind(&members).
		Send(c)

	return members, err
}

type AddGuildMemberParams struct {
	AccessToken string              `json:"access_token"`
	Nick        string              `json:"nick,omitempty"`
	Roles       []objects.Snowflake `json:"roles,omitempty"`
	Mute        bool                `json:"mute"`
	Deaf        bool                `json:"deaf"`
	Reason      string              `json:"-"`
}

func (c *Client) AddGuildMember(guild, user objects.Snowflake, params *AddGuildMemberParams) (*objects.GuildMember, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	reason := ""
	if params != nil {
		reason = params.Reason
	}

	member := &objects.GuildMember{}
	err = NewRequest().
		Method(http.MethodPut).
		Path(fmt.Sprintf(GuildMemberFmt, guild, user)).
		ContentType(JsonContentType).
		Body(data).
		Reason(reason).
		Expect(http.StatusCreated, http.StatusNoContent).
		Bind(member).
		Send(c)

	return member, err
}

type ModifyGuildMemberParams struct {
	Nick      *string              `json:"nick,omitempty"`
	Roles     *[]objects.Snowflake `json:"roles,omitempty"`
	Mute      *bool                `json:"mute,omitempty"`
	Deaf      *bool                `json:"deaf,omitempty"`
	ChannelID *objects.Snowflake   `json:"channel_id,omitempty"`
	Reason    string               `json:"-"`
}

func (c *Client) ModifyGuildMember(guild, member objects.Snowflake, params *ModifyGuildMemberParams) (*objects.GuildMember, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	reason := ""
	if params != nil {
		reason = params.Reason
	}
	m := &objects.GuildMember{}
	err = NewRequest().
		Method(http.MethodPatch).
		Path(fmt.Sprintf(GuildMemberFmt, guild, member)).
		ContentType(JsonContentType).
		Body(data).
		Reason(reason).
		Expect(http.StatusOK).
		Bind(m).
		Send(c)

	return m, err
}

type ModifyCurrentUserNickParams struct {
	Nick   string `json:"nick"`
	Reason string `json:"-"`
}

func (c *Client) ModifyCurrentUserNick(guild objects.Snowflake, params *ModifyCurrentUserNickParams) (*ModifyCurrentUserNickParams, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	reason := ""
	if params != nil {
		reason = params.Reason
	}

	newNick := &ModifyCurrentUserNickParams{}
	err = NewRequest().
		Method(http.MethodPatch).
		Path(fmt.Sprintf(GuildMemberEditCurrentUserNickFmt, guild)).
		ContentType(JsonContentType).
		Body(data).
		Reason(reason).
		Expect(http.StatusOK).
		Bind(newNick).
		Send(c)

	return newNick, err
}

func (c *Client) AddGuildMemberRole(guild, user, role objects.Snowflake, reason string) error {
	return NewRequest().
		Method(http.MethodPut).
		Path(fmt.Sprintf(GuildMemberRoleFmt, guild, user, role)).
		ContentType(JsonContentType).
		Reason(reason).
		Expect(http.StatusNoContent).
		Send(c)
}

func (c *Client) RemoveGuildMemberRole(guild, user, role objects.Snowflake, reason string) error {
	return NewRequest().
		Method(http.MethodDelete).
		Path(fmt.Sprintf(GuildMemberRoleFmt, guild, user, role)).
		ContentType(JsonContentType).
		Reason(reason).
		Expect(http.StatusNoContent).
		Send(c)
}

func (c *Client) RemoveGuildMember(guild, user objects.Snowflake, reason string) error {
	return NewRequest().
		Method(http.MethodDelete).
		Path(fmt.Sprintf(GuildMemberFmt, guild, user)).
		ContentType(JsonContentType).
		Reason(reason).
		Expect(http.StatusNoContent).
		Send(c)
}

func (c *Client) GetGuildBans(guild objects.Snowflake) ([]*objects.Ban, error) {
	bans := []*objects.Ban{}
	err := NewRequest().
		Method(http.MethodGet).
		Path(fmt.Sprintf(GuildBansFmt, guild)).
		ContentType(JsonContentType).
		Expect(http.StatusOK).
		Bind(&bans).
		Send(c)

	return bans, err
}

func (c *Client) GetGuildBan(guild, user objects.Snowflake) (*objects.Ban, error) {
	ban := &objects.Ban{}
	err := NewRequest().
		Method(http.MethodGet).
		Path(fmt.Sprintf(GuildBanUserFmt, guild, user)).
		ContentType(JsonContentType).
		Expect(http.StatusOK).
		Bind(ban).
		Send(c)

	return ban, err
}

type CreateGuildBanParams struct {
	DeleteMessageDays int    `json:"delete_message_days,omitempty"`
	Reason            string `json:"reason,omitempty"`
}

func (c *Client) CreateBan(guild, user objects.Snowflake, params *CreateGuildBanParams) error {
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}

	return NewRequest().
		Method(http.MethodPut).
		Path(fmt.Sprintf(GuildBanUserFmt, guild, user)).
		ContentType(JsonContentType).
		Body(data).
		Expect(http.StatusNoContent).
		Send(c)
}

func (c *Client) RemoveGuildBan(guild, user objects.Snowflake, reason string) error {
	return NewRequest().
		Method(http.MethodDelete).
		Path(fmt.Sprintf(GuildBanUserFmt, guild, user)).
		ContentType(JsonContentType).
		Reason(reason).
		Expect(http.StatusNoContent).
		Send(c)
}

func (c *Client) GetGuildRoles(guild objects.Snowflake) ([]*objects.Role, error) {
	roles := []*objects.Role{}
	err := NewRequest().
		Method(http.MethodGet).
		Path(fmt.Sprintf(GuildRolesFmt, guild)).
		ContentType(JsonContentType).
		Expect(http.StatusOK).
		Bind(&roles).
		Send(c)
	return roles, err
}

type CreateGuildRoleParams struct {
	Name        string                `json:"name,omitempty"`
	Permissions objects.PermissionBit `json:"permissions,omitempty"`
	Color       int                   `json:"color,omitempty"`
	Hoist       bool                  `json:"hoist,omitempty"`
	Mentionable bool                  `json:"mentionable,omitempty"`
	Reason      string                `json:"-"`
}

func (c *Client) CreateGuildRole(guild objects.Snowflake, params *CreateGuildRoleParams) (*objects.Role, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	reason := ""
	if params != nil {
		reason = params.Reason
	}

	role := &objects.Role{}
	err = NewRequest().
		Method(http.MethodPost).
		Path(fmt.Sprintf(GuildRolesFmt, guild)).
		ContentType(JsonContentType).
		Body(data).
		Reason(reason).
		Expect(http.StatusOK).
		Bind(role).
		Send(c)
	return role, err
}

type ModifyGuildRolePositionsParams struct {
	ID objects.Snowflake `json:"id"`
}

func (c *Client) ModifyGuildRolePositions(guild objects.Snowflake, params []*ModifyGuildRolePositionsParams) ([]*objects.Role, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	roles := []*objects.Role{}
	err = NewRequest().
		Method(http.MethodPatch).
		Path(fmt.Sprintf(GuildRolesFmt, guild)).
		ContentType(JsonContentType).
		Body(data).
		Expect(http.StatusOK).
		Bind(&roles).
		Send(c)

	return roles, err
}

type ModifyGuildRoleParams struct {
	Name        string                `json:"name,omitempty"`
	Permissions objects.PermissionBit `json:"permissions,omitempty"`
	Color       int                   `json:"color,omitempty"`
	Hoist       *bool                 `json:"hoist,omitempty"`
	Mentionable *bool                 `json:"mentionable,omitempty"`
	Reason      string                `json:"-"`
}

func (c *Client) ModifyGuildRole(guild, role objects.Snowflake, params *ModifyGuildRoleParams) (*objects.Role, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	reason := ""
	if params != nil {
		reason = params.Reason
	}

	r := &objects.Role{}
	err = NewRequest().
		Method(http.MethodPatch).
		Path(fmt.Sprintf(GuildRoleFmt, guild, role)).
		ContentType(JsonContentType).
		Body(data).
		Reason(reason).
		Expect(http.StatusOK).
		Bind(r).
		Send(c)
	return r, err
}

func (c *Client) DeleteGuildRole(guild, role objects.Snowflake, reason string) error {
	return NewRequest().
		Method(http.MethodDelete).
		Path(fmt.Sprintf(GuildRoleFmt, guild, role)).
		ContentType(JsonContentType).
		Reason(reason).
		Expect(http.StatusNoContent).
		Send(c)
}

type GetGuildPruneCountParams struct {
	Days         int                 `url:"days,omitempty"`
	IncludeRoles []objects.Snowflake `url:"include_roles,omitempty"`
}

func (c *Client) GetGuildPruneCount(guild objects.Snowflake, params *GetGuildPruneCountParams) (int, error) {
	u, err := url.Parse(fmt.Sprintf(GuildPruneFmt, guild))
	if err != nil {
		return 0, err
	}
	v, err := query.Values(params)
	if err != nil {
		return 0, err
	}

	u.RawQuery = v.Encode()

	pruned := &struct {
		Pruned int `json:"pruned"`
	}{}

	err = NewRequest().
		Method(http.MethodGet).
		Path(u.String()).
		ContentType(JsonContentType).
		Expect(http.StatusOK).
		Bind(pruned).
		Send(c)
	return pruned.Pruned, err
}

type BeginGuildPruneParams struct {
	Days              int                 `json:"days,omitempty"`
	ComputePruneCount bool                `json:"compute_prune_count"`
	IncludeRoles      []objects.Snowflake `json:"include_roles,omitempty"`
	Reason            string              `json:"-"`
}

func (c *Client) BeginGuildPrune(guild objects.Snowflake, params *BeginGuildPruneParams) (int, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return 0, err
	}

	reason := ""
	if params != nil {
		reason = params.Reason
	}

	pruned := &struct {
		Pruned int `json:"pruned"`
	}{}

	err = NewRequest().
		Method(http.MethodPost).
		Path(fmt.Sprintf(GuildPruneFmt, guild)).
		ContentType(JsonContentType).
		Body(body).
		Reason(reason).
		Bind(pruned).
		Expect(http.StatusOK).
		Send(c)

	return pruned.Pruned, err
}

func (c *Client) GetGuildVoiceRegions(guild objects.Snowflake) ([]*objects.VoiceRegion, error) {
	var regions []*objects.VoiceRegion
	err := NewRequest().
		Method(http.MethodGet).
		Path(fmt.Sprintf(GuildVoiceRegionsFmt, guild)).
		ContentType(JsonContentType).
		Expect(http.StatusOK).
		Bind(&regions).
		Send(c)

	return regions, err
}

func (c *Client) GetGuildInvites(guild objects.Snowflake) ([]*objects.Invite, error) {
	var invites []*objects.Invite
	err := NewRequest().
		Method(http.MethodGet).
		Path(fmt.Sprintf(GuildInvitesFmt, guild)).
		ContentType(JsonContentType).
		Expect(http.StatusOK).
		Bind(&invites).
		Send(c)

	return invites, err
}

func (c *Client) GetGuildIntegrations(guild objects.Snowflake) ([]*objects.Integration, error) {
	var integrations []*objects.Integration
	err := NewRequest().
		Method(http.MethodGet).
		Path(fmt.Sprintf(IntegrationsBaseFmt, guild)).
		ContentType(JsonContentType).
		Expect(http.StatusOK).
		Bind(&integrations).
		Send(c)

	return integrations, err
}

func (c *Client) DeleteGuildIntegration(guild, integration objects.Snowflake, reason string) error {
	return NewRequest().
		Method(http.MethodDelete).
		Path(fmt.Sprintf(IntegrationBaseFmt, guild, integration)).
		ContentType(JsonContentType).
		Reason(reason).
		Expect(http.StatusNoContent).
		Send(c)
}

func (c *Client) GetGuildWidgetSettings(guild objects.Snowflake) (*objects.GuildWidget, error) {
	widget := &objects.GuildWidget{}
	err := NewRequest().
		Method(http.MethodGet).
		Path(fmt.Sprintf(GuildWidgetFmt, guild)).
		ContentType(JsonContentType).
		Expect(http.StatusOK).
		Bind(widget).
		Send(c)
	return widget, err
}

type GuildWidgetParams struct {
	Enabled   *bool             `json:"enabled,omitempty"`
	ChannelID objects.Snowflake `json:"channel_id,omitempty"`
	Reason    string            `json:"-"`
}

func (c *Client) ModifyGuildWidget(guild objects.Snowflake, params *GuildWidgetParams) (*objects.GuildWidget, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	reason := ""
	if params != nil {
		reason = params.Reason
	}

	widget := &objects.GuildWidget{}
	err = NewRequest().
		Method(http.MethodPatch).
		Path(fmt.Sprintf(GuildWidgetFmt, guild)).
		ContentType(JsonContentType).
		Body(data).
		Reason(reason).
		Bind(widget).
		Expect(http.StatusOK).
		Send(c)

	return widget, err
}

func (c *Client) GetGuildWidget(guild objects.Snowflake) (*objects.GuildWidgetJSON, error) {
	widget := &objects.GuildWidgetJSON{}
	err := NewRequest().
		Method(http.MethodGet).
		Path(fmt.Sprintf(GuildWidgetJSONFmt, guild)).
		ContentType(JsonContentType).
		Expect(http.StatusOK).
		Bind(widget).
		Send(c)
	return widget, err
}

func (c *Client) GetGuildVanityURL(guild objects.Snowflake) (*objects.Invite, error) {
	invite := &objects.Invite{}
	err := NewRequest().
		Method(http.MethodGet).
		Path(fmt.Sprintf(GuildVanityURLFmt, guild)).
		ContentType(JsonContentType).
		Expect(http.StatusOK).
		Bind(invite).
		Send(c)
	return invite, err
}

type GuildWidgetImageParams struct {
	Style string
}

func (c *Client) GetGuildWidgetImage(guild objects.Snowflake, params *GuildWidgetImageParams) (image.Image, error) {
	u, err := url.Parse(fmt.Sprintf(GuildWidgetImageFmt, guild))
	if err != nil {
		return nil, err
	}

	v, err := query.Values(params)
	if err != nil {
		return nil, err
	}

	u.RawQuery = v.Encode()

	res, err := NewRequest().
		Method(http.MethodGet).
		Path(u.String()).
		ContentType(JsonContentType).
		Expect(http.StatusOK).
		SendRaw(c)

	if err != nil {
		return nil, err
	}

	img, err := png.Decode(bytes.NewReader(res.Body))
	if err != nil {
		return nil, err
	}

	return img, nil
}

func (c *Client) GetGuildWelcomeScreen(guild objects.Snowflake) (*objects.MembershipScreening, error) {
	screening := &objects.MembershipScreening{}
	err := NewRequest().
		Method(http.MethodGet).
		Path(fmt.Sprintf(GuildMembershipScreeningFmt, guild)).
		ContentType(JsonContentType).
		Expect(http.StatusOK).
		Bind(screening).
		Send(c)
	return screening, err
}

type ModifyGuildMembershipScreeningParams struct {
	Enabled     *bool  `json:"enabled,omitempty"`
	FormFields  string `json:"form_fields,omitempty"`
	Description string `json:"description,omitempty"`
	Reason      string `json:"-"`
}

func (c *Client) ModifyGuildWelcomeScreen(guild objects.Snowflake, params *ModifyGuildMembershipScreeningParams) (*objects.MembershipScreening, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	reason := ""
	if params != nil {
		reason = params.Reason
	}

	screening := &objects.MembershipScreening{}
	err = NewRequest().
		Method(http.MethodPatch).
		Path(fmt.Sprintf(GuildMembershipScreeningFmt, guild)).
		ContentType(JsonContentType).
		Body(data).
		Reason(reason).
		Bind(screening).
		Expect(http.StatusOK).
		Send(c)

	return screening, err
}
