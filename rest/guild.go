package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"net/http"
	"net/url"

	"github.com/kelwing/wumpgo/objects"
	"github.com/kelwing/wumpgo/objects/permissions"
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

func (c *Client) CreateGuild(ctx context.Context, params *CreateGuildParams) (*objects.Guild, error) {
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
		WithContext(ctx).
		Path(GuildCreateFmt).
		ContentType(JsonContentType).
		Body(data).
		Reason(reason).
		Bind(guild).
		Send(c)

	return guild, err
}

func (c *Client) GetGuild(ctx context.Context, id objects.SnowflakeObject) (*objects.Guild, error) {
	guild := &objects.Guild{}
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildBaseFmt, id.GetID())).
		ContentType(JsonContentType).
		Bind(guild).
		Send(c)

	return guild, err
}

func (c *Client) GetGuildPreview(ctx context.Context, id objects.SnowflakeObject) (*objects.GuildPreview, error) {
	preview := &objects.GuildPreview{}
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildPreviewFmt, id.GetID())).
		ContentType(JsonContentType).
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

func (c *Client) ModifyGuild(ctx context.Context, id objects.SnowflakeObject, params *ModifyGuildParams) (*objects.Guild, error) {
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
		WithContext(ctx).
		Path(fmt.Sprintf(GuildBaseFmt, id.GetID())).
		ContentType(JsonContentType).
		Body(data).
		Reason(reason).
		Bind(guild).
		Send(c)

	return guild, err
}

func (c *Client) DeleteGuild(ctx context.Context, id objects.SnowflakeObject) error {
	return NewRequest().
		Method(http.MethodDelete).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildBaseFmt, id.GetID())).
		ContentType(JsonContentType).
		Send(c)
}

func (c *Client) GetGuildChannels(ctx context.Context, id objects.SnowflakeObject) ([]*objects.Channel, error) {
	channels := []*objects.Channel{}
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildChannelsFmt, id.GetID())).
		ContentType(JsonContentType).
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

func (c *Client) CreateGuildChannel(ctx context.Context, id objects.SnowflakeObject, params *ChannelCreateParams) (*objects.Channel, error) {
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
		WithContext(ctx).
		Path(fmt.Sprintf(GuildChannelsFmt, id.GetID())).
		ContentType(JsonContentType).
		Body(data).
		Reason(reason).
		Bind(channel).
		Send(c)

	return channel, err
}

type ModifyChannelPositionParams struct {
	ID       objects.Snowflake `json:"id"`
	Position int               `json:"position"`
}

func (c *Client) ModifyGuildChannelPositions(ctx context.Context, id objects.SnowflakeObject, params []*ModifyChannelPositionParams, reason string) error {
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}

	return NewRequest().
		Method(http.MethodPatch).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildChannelsFmt, id.GetID())).
		ContentType(JsonContentType).
		Body(data).
		Reason(reason).
		Send(c)
}

func (c *Client) ListActiveThreads(ctx context.Context, id objects.SnowflakeObject) ([]*ListThreadsResponse, error) {
	channels := []*ListThreadsResponse{}
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildThreadsFmt, id.GetID())).
		ContentType(JsonContentType).
		Bind(&channels).
		Send(c)

	return channels, err
}

func (c *Client) GetGuildMember(ctx context.Context, guild, user objects.SnowflakeObject) (*objects.GuildMember, error) {
	member := &objects.GuildMember{}
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildMemberFmt, guild.GetID(), user.GetID())).
		ContentType(JsonContentType).
		Bind(member).
		Send(c)

	return member, err
}

type ListGuildMembersParams struct {
	Limit int               `url:"limit,omitempty"`
	After objects.Snowflake `url:"after,omitempty"`
}

func (c *Client) ListGuildMembers(ctx context.Context, guild objects.SnowflakeObject, params *ListGuildMembersParams) ([]*objects.GuildMember, error) {
	u, err := url.Parse(fmt.Sprintf(GuildMembersFmt, guild.GetID()))
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
		WithContext(ctx).
		Path(u.String()).
		ContentType(JsonContentType).
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

func (c *Client) AddGuildMember(ctx context.Context, guild, user objects.SnowflakeObject, params *AddGuildMemberParams) (*objects.GuildMember, error) {
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
		WithContext(ctx).
		Path(fmt.Sprintf(GuildMemberFmt, guild.GetID(), user.GetID())).
		ContentType(JsonContentType).
		Body(data).
		Reason(reason).
		Bind(member).
		Send(c)

	return member, err
}

type ModifyGuildMemberParams struct {
	Nick                       *string              `json:"nick,omitempty"`
	Roles                      *[]objects.Snowflake `json:"roles,omitempty"`
	Mute                       *bool                `json:"mute,omitempty"`
	Deaf                       *bool                `json:"deaf,omitempty"`
	ChannelID                  *objects.Snowflake   `json:"channel_id,omitempty"`
	CommunicationDisabledUntil *objects.Time        `json:"communication_disabled_until,omitempty"`
	Reason                     string               `json:"-"`
}

func (c *Client) ModifyGuildMember(ctx context.Context, guild, member objects.SnowflakeObject, params *ModifyGuildMemberParams) (*objects.GuildMember, error) {
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
		WithContext(ctx).
		Path(fmt.Sprintf(GuildMemberFmt, guild.GetID(), member.GetID())).
		ContentType(JsonContentType).
		Body(data).
		Reason(reason).
		Bind(m).
		Send(c)

	return m, err
}

type ModifyCurrentUserNickParams struct {
	Nick   string `json:"nick"`
	Reason string `json:"-"`
}

func (c *Client) ModifyCurrentUserNick(ctx context.Context, guild objects.SnowflakeObject, params *ModifyCurrentUserNickParams) (*ModifyCurrentUserNickParams, error) {
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
		WithContext(ctx).
		Path(fmt.Sprintf(GuildMemberEditCurrentUserNickFmt, guild.GetID())).
		ContentType(JsonContentType).
		Body(data).
		Reason(reason).
		Bind(newNick).
		Send(c)

	return newNick, err
}

func (c *Client) AddGuildMemberRole(ctx context.Context, guild, user, role objects.SnowflakeObject, reason string) error {
	return NewRequest().
		Method(http.MethodPut).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildMemberRoleFmt, guild.GetID(), user.GetID(), role.GetID())).
		ContentType(JsonContentType).
		Reason(reason).
		Send(c)
}

func (c *Client) RemoveGuildMemberRole(ctx context.Context, guild, user, role objects.SnowflakeObject, reason string) error {
	return NewRequest().
		Method(http.MethodDelete).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildMemberRoleFmt, guild.GetID(), user.GetID(), role.GetID())).
		ContentType(JsonContentType).
		Reason(reason).
		Send(c)
}

func (c *Client) RemoveGuildMember(ctx context.Context, guild, user objects.SnowflakeObject, reason string) error {
	return NewRequest().
		Method(http.MethodDelete).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildMemberFmt, guild.GetID(), user.GetID())).
		ContentType(JsonContentType).
		Reason(reason).
		Send(c)
}

func (c *Client) GetGuildBans(ctx context.Context, guild objects.SnowflakeObject) ([]*objects.Ban, error) {
	bans := []*objects.Ban{}
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildBansFmt, guild.GetID())).
		ContentType(JsonContentType).
		Bind(&bans).
		Send(c)

	return bans, err
}

func (c *Client) GetGuildBan(ctx context.Context, guild, user objects.SnowflakeObject) (*objects.Ban, error) {
	ban := &objects.Ban{}
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildBanUserFmt, guild.GetID(), user.GetID())).
		ContentType(JsonContentType).
		Bind(ban).
		Send(c)

	return ban, err
}

type CreateGuildBanParams struct {
	DeleteMessageDays int    `json:"delete_message_days,omitempty"`
	Reason            string `json:"reason,omitempty"`
}

func (c *Client) CreateBan(ctx context.Context, guild, user objects.SnowflakeObject, params *CreateGuildBanParams) error {
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}

	return NewRequest().
		Method(http.MethodPut).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildBanUserFmt, guild.GetID(), user.GetID())).
		ContentType(JsonContentType).
		Body(data).
		Send(c)
}

func (c *Client) RemoveGuildBan(ctx context.Context, guild, user objects.SnowflakeObject, reason string) error {
	return NewRequest().
		Method(http.MethodDelete).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildBanUserFmt, guild.GetID(), user.GetID())).
		ContentType(JsonContentType).
		Reason(reason).
		Send(c)
}

func (c *Client) GetGuildRoles(ctx context.Context, guild objects.SnowflakeObject) ([]*objects.Role, error) {
	roles := []*objects.Role{}
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildRolesFmt, guild.GetID())).
		ContentType(JsonContentType).
		Bind(&roles).
		Send(c)
	return roles, err
}

type CreateGuildRoleParams struct {
	Name        string                    `json:"name,omitempty"`
	Permissions permissions.PermissionBit `json:"permissions,omitempty"`
	Color       int                       `json:"color,omitempty"`
	Hoist       bool                      `json:"hoist,omitempty"`
	Mentionable bool                      `json:"mentionable,omitempty"`
	Reason      string                    `json:"-"`
}

func (c *Client) CreateGuildRole(ctx context.Context, guild objects.SnowflakeObject, params *CreateGuildRoleParams) (*objects.Role, error) {
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
		WithContext(ctx).
		Path(fmt.Sprintf(GuildRolesFmt, guild.GetID())).
		ContentType(JsonContentType).
		Body(data).
		Reason(reason).
		Bind(role).
		Send(c)
	return role, err
}

type ModifyGuildRolePositionsParams struct {
	ID objects.Snowflake `json:"id"`
}

func (c *Client) ModifyGuildRolePositions(ctx context.Context, guild objects.SnowflakeObject, params []*ModifyGuildRolePositionsParams) ([]*objects.Role, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	roles := []*objects.Role{}
	err = NewRequest().
		Method(http.MethodPatch).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildRolesFmt, guild.GetID())).
		ContentType(JsonContentType).
		Body(data).
		Bind(&roles).
		Send(c)

	return roles, err
}

type ModifyGuildRoleParams struct {
	Name        string                    `json:"name,omitempty"`
	Permissions permissions.PermissionBit `json:"permissions,omitempty"`
	Color       int                       `json:"color,omitempty"`
	Hoist       *bool                     `json:"hoist,omitempty"`
	Mentionable *bool                     `json:"mentionable,omitempty"`
	Reason      string                    `json:"-"`
}

func (c *Client) ModifyGuildRole(ctx context.Context, guild, role objects.SnowflakeObject, params *ModifyGuildRoleParams) (*objects.Role, error) {
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
		WithContext(ctx).
		Path(fmt.Sprintf(GuildRoleFmt, guild.GetID(), role.GetID())).
		ContentType(JsonContentType).
		Body(data).
		Reason(reason).
		Bind(r).
		Send(c)
	return r, err
}

func (c *Client) DeleteGuildRole(ctx context.Context, guild, role objects.SnowflakeObject, reason string) error {
	return NewRequest().
		Method(http.MethodDelete).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildRoleFmt, guild.GetID(), role.GetID())).
		ContentType(JsonContentType).
		Reason(reason).
		Send(c)
}

type GetGuildPruneCountParams struct {
	Days         int                 `url:"days,omitempty"`
	IncludeRoles []objects.Snowflake `url:"include_roles,omitempty"`
}

func (c *Client) GetGuildPruneCount(ctx context.Context, guild objects.SnowflakeObject, params *GetGuildPruneCountParams) (int, error) {
	u, err := url.Parse(fmt.Sprintf(GuildPruneFmt, guild.GetID()))
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
		WithContext(ctx).
		Path(u.String()).
		ContentType(JsonContentType).
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

func (c *Client) BeginGuildPrune(ctx context.Context, guild objects.SnowflakeObject, params *BeginGuildPruneParams) (int, error) {
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
		WithContext(ctx).
		Path(fmt.Sprintf(GuildPruneFmt, guild.GetID())).
		ContentType(JsonContentType).
		Body(body).
		Reason(reason).
		Bind(pruned).
		Send(c)

	return pruned.Pruned, err
}

func (c *Client) GetGuildVoiceRegions(ctx context.Context, guild objects.SnowflakeObject) ([]*objects.VoiceRegion, error) {
	var regions []*objects.VoiceRegion
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildVoiceRegionsFmt, guild.GetID())).
		ContentType(JsonContentType).
		Bind(&regions).
		Send(c)

	return regions, err
}

func (c *Client) GetGuildInvites(ctx context.Context, guild objects.SnowflakeObject) ([]*objects.Invite, error) {
	var invites []*objects.Invite
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildInvitesFmt, guild.GetID())).
		ContentType(JsonContentType).
		Bind(&invites).
		Send(c)

	return invites, err
}

func (c *Client) GetGuildIntegrations(ctx context.Context, guild objects.SnowflakeObject) ([]*objects.Integration, error) {
	var integrations []*objects.Integration
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(IntegrationsBaseFmt, guild.GetID())).
		ContentType(JsonContentType).
		Bind(&integrations).
		Send(c)

	return integrations, err
}

func (c *Client) DeleteGuildIntegration(ctx context.Context, guild, integration objects.SnowflakeObject, reason string) error {
	return NewRequest().
		Method(http.MethodDelete).
		WithContext(ctx).
		Path(fmt.Sprintf(IntegrationBaseFmt, guild.GetID(), integration.GetID())).
		ContentType(JsonContentType).
		Reason(reason).
		Send(c)
}

func (c *Client) GetGuildWidgetSettings(ctx context.Context, guild objects.SnowflakeObject) (*objects.GuildWidget, error) {
	widget := &objects.GuildWidget{}
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildWidgetFmt, guild.GetID())).
		ContentType(JsonContentType).
		Bind(widget).
		Send(c)
	return widget, err
}

type GuildWidgetParams struct {
	Enabled   *bool             `json:"enabled,omitempty"`
	ChannelID objects.Snowflake `json:"channel_id,omitempty"`
	Reason    string            `json:"-"`
}

func (c *Client) ModifyGuildWidget(ctx context.Context, guild objects.SnowflakeObject, params *GuildWidgetParams) (*objects.GuildWidget, error) {
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
		WithContext(ctx).
		Path(fmt.Sprintf(GuildWidgetFmt, guild.GetID())).
		ContentType(JsonContentType).
		Body(data).
		Reason(reason).
		Bind(widget).
		Send(c)

	return widget, err
}

func (c *Client) GetGuildWidget(ctx context.Context, guild objects.SnowflakeObject) (*objects.GuildWidgetJSON, error) {
	widget := &objects.GuildWidgetJSON{}
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildWidgetJSONFmt, guild.GetID())).
		ContentType(JsonContentType).
		Bind(widget).
		Send(c)
	return widget, err
}

func (c *Client) GetGuildVanityURL(ctx context.Context, guild objects.SnowflakeObject) (*objects.Invite, error) {
	invite := &objects.Invite{}
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildVanityURLFmt, guild.GetID())).
		ContentType(JsonContentType).
		Bind(invite).
		Send(c)
	return invite, err
}

type GuildWidgetImageParams struct {
	Style string
}

func (c *Client) GetGuildWidgetImage(ctx context.Context, guild objects.SnowflakeObject, params *GuildWidgetImageParams) (image.Image, error) {
	u, err := url.Parse(fmt.Sprintf(GuildWidgetImageFmt, guild.GetID()))
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
		WithContext(ctx).
		Path(u.String()).
		ContentType(JsonContentType).
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

func (c *Client) GetGuildWelcomeScreen(ctx context.Context, guild objects.SnowflakeObject) (*objects.MembershipScreening, error) {
	screening := &objects.MembershipScreening{}
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildMembershipScreeningFmt, guild.GetID())).
		ContentType(JsonContentType).
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

func (c *Client) ModifyGuildWelcomeScreen(ctx context.Context, guild objects.SnowflakeObject, params *ModifyGuildMembershipScreeningParams) (*objects.MembershipScreening, error) {
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
		WithContext(ctx).
		Path(fmt.Sprintf(GuildMembershipScreeningFmt, guild.GetID())).
		ContentType(JsonContentType).
		Body(data).
		Reason(reason).
		Bind(screening).
		Send(c)

	return screening, err
}
