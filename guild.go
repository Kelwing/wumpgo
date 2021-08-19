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

	res, err := c.request(&request{
		method:      http.MethodPost,
		path:        GuildCreateFmt,
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

	guild := &objects.Guild{}
	if err = res.JSON(guild); err != nil {
		return nil, err
	}

	return guild, nil
}

func (c *Client) GetGuild(id objects.Snowflake) (*objects.Guild, error) {
	res, err := c.request(&request{
		method:      http.MethodGet,
		path:        fmt.Sprintf(GuildBaseFmt, id),
		contentType: JsonContentType,
	})
	if err != nil {
		return nil, err
	}

	if err = res.ExpectsStatus(http.StatusOK); err != nil {
		return nil, err
	}

	guild := &objects.Guild{}
	if err = res.JSON(guild); err != nil {
		return nil, err
	}

	return guild, nil
}

func (c *Client) GetGuildPreview(id objects.Snowflake) (*objects.GuildPreview, error) {
	res, err := c.request(&request{
		method:      http.MethodGet,
		path:        fmt.Sprintf(GuildPreviewFmt, id),
		contentType: JsonContentType,
	})

	if err != nil {
		return nil, err
	}

	if err = res.ExpectsStatus(http.StatusOK); err != nil {
		return nil, err
	}

	preview := &objects.GuildPreview{}
	if err = res.JSON(preview); err != nil {
		return nil, err
	}

	return preview, nil
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

	res, err := c.request(&request{
		method:      http.MethodPatch,
		path:        fmt.Sprintf(GuildBaseFmt, id),
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

	guild := &objects.Guild{}
	if err = res.JSON(guild); err != nil {
		return nil, err
	}

	return guild, nil
}

func (c *Client) DeleteGuild(id objects.Snowflake) error {
	res, err := c.request(&request{
		method:      http.MethodDelete,
		path:        fmt.Sprintf(GuildBaseFmt, id),
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

func (c *Client) GetGuildChannels(id objects.Snowflake) ([]*objects.Channel, error) {
	res, err := c.request(&request{
		method:      http.MethodGet,
		path:        fmt.Sprintf(GuildChannelsFmt, id),
		contentType: JsonContentType,
	})
	if err != nil {
		return nil, err
	}

	if err = res.ExpectsStatus(http.StatusOK); err != nil {
		return nil, err
	}

	var channels []*objects.Channel
	if err = res.JSON(&channels); err != nil {
		return nil, err
	}

	return channels, nil
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

	res, err := c.request(&request{
		method:      http.MethodPost,
		path:        fmt.Sprintf(GuildChannelsFmt, id),
		contentType: JsonContentType,
		body:        data,
		reason:      reason,
	})
	if err != nil {
		return nil, err
	}

	if err = res.ExpectAnyStatus(http.StatusOK, http.StatusCreated); err != nil {
		return nil, err
	}

	channel := &objects.Channel{}
	if err = res.JSON(channel); err != nil {
		return nil, err
	}

	return channel, nil
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

	res, err := c.request(&request{
		method:      http.MethodPatch,
		path:        fmt.Sprintf(GuildChannelsFmt, id),
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

func (c *Client) GetGuildMember(guild, user objects.Snowflake) (*objects.GuildMember, error) {
	res, err := c.request(&request{
		method:      http.MethodGet,
		path:        fmt.Sprintf(GuildMemberFmt, guild, user),
		contentType: JsonContentType,
	})
	if err != nil {
		return nil, err
	}

	if err = res.ExpectsStatus(http.StatusOK); err != nil {
		return nil, err
	}

	member := &objects.GuildMember{}
	if err = res.JSON(member); err != nil {
		return nil, err
	}

	return member, nil
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

	var members []*objects.GuildMember
	if err = res.JSON(&members); err != nil {
		return nil, err
	}

	return members, nil
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

	res, err := c.request(&request{
		method:      http.MethodPut,
		path:        fmt.Sprintf(GuildMemberFmt, guild, user),
		contentType: JsonContentType,
		body:        data,
		reason:      reason,
	})
	if err != nil {
		return nil, err
	}

	if err = res.ExpectAnyStatus(http.StatusNoContent, http.StatusCreated); err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusCreated {
		member := &objects.GuildMember{}
		if err = res.JSON(member); err != nil {
			return nil, err
		}
		return member, nil
	}
	return nil, nil
}

type ModifyGuildMemberParams struct {
	Nick      *string              `json:"nick,omitempty"`
	Roles     *[]objects.Snowflake `json:"roles,omitempty"`
	Mute      *bool                `json:"mute,omitempty"`
	Deaf      *bool                `json:"deaf,omitempty"`
	ChannelID *objects.Snowflake   `json:"channel_id,omitempty"`
	Reason    string               `json:"-"`
}

func (c *Client) ModifyGuildMember(guild, member objects.Snowflake, params *ModifyGuildMemberParams) error {
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}

	reason := ""
	if params != nil {
		reason = params.Reason
	}

	res, err := c.request(&request{
		method:      http.MethodPatch,
		path:        fmt.Sprintf(GuildMemberFmt, guild, member),
		contentType: JsonContentType,
		body:        data,
		reason:      reason,
	})
	if err != nil {
		return nil
	}

	if err = res.ExpectsStatus(http.StatusNoContent); err != nil {
		return err
	}

	return nil
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

	res, err := c.request(&request{
		method:      http.MethodPatch,
		path:        fmt.Sprintf(GuildMemberEditCurrentUserNickFmt, guild),
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

	newNick := &ModifyCurrentUserNickParams{}
	if err = res.JSON(newNick); err != nil {
		return nil, err
	}

	return newNick, nil
}

func (c *Client) AddGuildMemberRole(guild, user, role objects.Snowflake, reason string) error {
	res, err := c.request(&request{
		method:      http.MethodPut,
		path:        fmt.Sprintf(GuildMemberRoleFmt, guild, user, role),
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

func (c *Client) RemoveGuildMemberRole(guild, user, role objects.Snowflake, reason string) error {
	res, err := c.request(&request{
		method:      http.MethodDelete,
		path:        fmt.Sprintf(GuildMemberRoleFmt, guild, user, role),
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

func (c *Client) RemoveGuildMember(guild, user objects.Snowflake, reason string) error {
	res, err := c.request(&request{
		method:      http.MethodDelete,
		path:        fmt.Sprintf(GuildMemberFmt, guild, user),
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

func (c *Client) GetGuildBans(guild objects.Snowflake) ([]*objects.Ban, error) {
	res, err := c.request(&request{
		method:      http.MethodGet,
		path:        fmt.Sprintf(GuildBansFmt, guild),
		contentType: JsonContentType,
	})
	if err != nil {
		return nil, err
	}

	if err = res.ExpectsStatus(http.StatusOK); err != nil {
		return nil, err
	}

	var bans []*objects.Ban
	if err = res.JSON(&bans); err != nil {
		return nil, err
	}

	return bans, nil
}

func (c *Client) GetBan(guild, user objects.Snowflake) (*objects.Ban, error) {
	res, err := c.request(&request{
		method:      http.MethodGet,
		path:        fmt.Sprintf(GuildBanUserFmt, guild, user),
		contentType: JsonContentType,
	})
	if err != nil {
		return nil, err
	}

	if err = res.ExpectsStatus(http.StatusOK); err != nil {
		return nil, err
	}

	ban := &objects.Ban{}
	if err = res.JSON(ban); err != nil {
		return nil, err
	}
	return ban, nil
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

	res, err := c.request(&request{
		method:      http.MethodPut,
		path:        fmt.Sprintf(GuildBanUserFmt, guild, user),
		contentType: JsonContentType,
		body:        data,
	})
	if err != nil {
		return err
	}

	if err = res.ExpectsStatus(http.StatusOK); err != nil {
		return err
	}
	return nil
}

func (c *Client) RemoveGuildBan(guild, user objects.Snowflake, reason string) error {
	res, err := c.request(&request{
		method:      http.MethodDelete,
		path:        fmt.Sprintf(GuildBanUserFmt, guild, user),
		contentType: JsonContentType,
		body:        nil,
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

func (c *Client) GetGuildRoles(guild objects.Snowflake) ([]*objects.Role, error) {
	res, err := c.request(&request{
		method:      http.MethodGet,
		path:        fmt.Sprintf(GuildRolesFmt, guild),
		contentType: JsonContentType,
	})
	if err != nil {
		return nil, err
	}

	if err = res.ExpectsStatus(http.StatusOK); err != nil {
		return nil, err
	}

	var roles []*objects.Role
	if err = res.JSON(&roles); err != nil {
		return nil, err
	}

	return roles, nil
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

	res, err := c.request(&request{
		method:      http.MethodPost,
		path:        fmt.Sprintf(GuildRolesFmt, guild),
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

	role := &objects.Role{}
	if err = res.JSON(role); err != nil {
		return nil, err
	}

	return role, nil
}

type ModifyGuildRolePositionsParams struct {
	ID objects.Snowflake `json:"id"`
}

func (c *Client) ModifyGuildRolePositions(guild objects.Snowflake, params []*ModifyGuildRolePositionsParams) ([]*objects.Role, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	res, err := c.request(&request{
		method:      http.MethodPatch,
		path:        fmt.Sprintf(GuildRolesFmt, guild),
		contentType: JsonContentType,
		body:        data,
	})
	if err != nil {
		return nil, err
	}

	if err = res.ExpectsStatus(http.StatusOK); err != nil {
		return nil, err
	}

	var roles []*objects.Role
	if err = res.JSON(&roles); err != nil {
		return nil, err
	}

	return roles, nil
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

	res, err := c.request(&request{
		method:      http.MethodPatch,
		path:        fmt.Sprintf(GuildRoleFmt, guild, role),
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

	r := &objects.Role{}
	if err = res.JSON(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (c *Client) DeleteGuildRole(guild, role objects.Snowflake, reason string) error {
	res, err := c.request(&request{
		method:      http.MethodDelete,
		path:        fmt.Sprintf(GuildRoleFmt, guild, role),
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

	res, err := c.request(&request{
		method:      http.MethodGet,
		path:        u.String(),
		contentType: JsonContentType,
	})
	if err != nil {
		return 0, err
	}

	if err = res.ExpectsStatus(http.StatusOK); err != nil {
		return 0, err
	}

	pruned := &struct {
		Pruned int `json:"pruned"`
	}{}

	if err = res.JSON(pruned); err != nil {
		return 0, err
	}
	return pruned.Pruned, nil
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

	res, err := c.request(&request{
		method:      http.MethodPost,
		path:        fmt.Sprintf(GuildPruneFmt, guild),
		contentType: JsonContentType,
		body:        body,
		reason:      reason,
	})
	if err != nil {
		return 0, err
	}

	if err = res.ExpectsStatus(http.StatusOK); err != nil {
		return 0, err
	}

	pruned := &struct {
		Pruned int `json:"pruned"`
	}{}

	if err = res.JSON(pruned); err != nil {
		return 0, err
	}

	return pruned.Pruned, nil
}

func (c *Client) GetGuildVoiceRegions(guild objects.Snowflake) ([]*objects.VoiceRegion, error) {
	res, err := c.request(&request{
		method:      http.MethodGet,
		path:        fmt.Sprintf(GuildVoiceRegionsFmt, guild),
		contentType: JsonContentType,
	})
	if err != nil {
		return nil, err
	}

	if err = res.ExpectsStatus(http.StatusOK); err != nil {
		return nil, err
	}

	var regions []*objects.VoiceRegion
	if err = res.JSON(&regions); err != nil {
		return nil, err
	}

	return regions, nil
}

func (c *Client) GetGuildInvites(guild objects.Snowflake) ([]*objects.Invite, error) {
	res, err := c.request(&request{
		method: http.MethodGet,
		path:   fmt.Sprintf(GuildInvitesFmt, guild),
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

func (c *Client) GetGuildIntegrations(guild objects.Snowflake) ([]*objects.Integration, error) {
	res, err := c.request(&request{
		method:      http.MethodGet,
		path:        fmt.Sprintf(IntegrationsBaseFmt, guild),
		contentType: JsonContentType,
	})
	if err != nil {
		return nil, err
	}

	if err = res.ExpectsStatus(http.StatusOK); err != nil {
		return nil, err
	}

	var integrations []*objects.Integration
	if err = res.JSON(&integrations); err != nil {
		return nil, err
	}

	return integrations, nil
}

type CreateGuildIntegrationParams struct {
	Type   string            `json:"type"`
	ID     objects.Snowflake `json:"id"`
	Reason string            `json:"-"`
}

func (c *Client) CreateGuildIntegration(guild objects.Snowflake, params *CreateGuildIntegrationParams) error {
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}

	reason := ""
	if params != nil {
		reason = params.Reason
	}

	res, err := c.request(&request{
		method:      http.MethodPost,
		path:        fmt.Sprintf(IntegrationsBaseFmt, guild),
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

type ModifyGuildIntegrationParams struct {
	ExpireBehaviour   *objects.ExpireBehaviour `json:"expire_behaviour,omitempty"`
	ExpireGracePeriod *int                     `json:"expire_grace_period,omitempty"`
	EnableEmoticons   *bool                    `json:"enable_emoticons,omitempty"`
	Reason            string                   `json:"-"`
}

func (c *Client) ModifyGuildIntegration(guild, integration objects.Snowflake, params *ModifyGuildIntegrationParams) error {
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}

	reason := ""
	if params != nil {
		reason = params.Reason
	}

	res, err := c.request(&request{
		method:      http.MethodPatch,
		path:        fmt.Sprintf(IntegrationBaseFmt, guild, integration),
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

func (c *Client) DeleteGuildIntegration(guild, integration objects.Snowflake, reason string) error {
	res, err := c.request(&request{
		method:      http.MethodDelete,
		path:        fmt.Sprintf(IntegrationBaseFmt, guild, integration),
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

func (c *Client) SyncGuildIntegration(guild, integration objects.Snowflake) error {
	res, err := c.request(&request{
		method:      http.MethodPost,
		path:        fmt.Sprintf(IntegrationSync, guild, integration),
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

func (c *Client) GetGuildWidgetSettings(guild objects.Snowflake) (*objects.GuildWidget, error) {
	res, err := c.request(&request{
		method:      http.MethodGet,
		path:        fmt.Sprintf(GuildWidgetFmt, guild),
		contentType: JsonContentType,
	})
	if err != nil {
		return nil, err
	}

	if err = res.ExpectsStatus(http.StatusOK); err != nil {
		return nil, err
	}

	widget := &objects.GuildWidget{}
	if err = res.JSON(widget); err != nil {
		return nil, err
	}
	return widget, nil
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

	res, err := c.request(&request{
		method:      http.MethodPatch,
		path:        fmt.Sprintf(GuildWidgetFmt, guild),
		body:        data,
		contentType: JsonContentType,
		reason:      reason,
	})
	if err != nil {
		return nil, err
	}

	if err = res.ExpectsStatus(http.StatusOK); err != nil {
		return nil, err
	}

	newWidget := &objects.GuildWidget{}
	if err = res.JSON(newWidget); err != nil {
		return nil, err
	}

	return newWidget, nil
}

func (c *Client) GetGuildWidget(guild objects.Snowflake) (*objects.GuildWidgetJSON, error) {
	res, err := c.request(&request{
		method:      http.MethodGet,
		path:        fmt.Sprintf(GuildWidgetJSONFmt, guild),
		contentType: JsonContentType,
	})
	if err != nil {
		return nil, err
	}

	if err = res.ExpectsStatus(http.StatusOK); err != nil {
		return nil, err
	}

	widget := &objects.GuildWidgetJSON{}
	if err = res.JSON(widget); err != nil {
		return nil, err
	}

	return widget, nil
}

func (c *Client) GetGuildVanityURL(guild objects.Snowflake) (*objects.Invite, error) {
	res, err := c.request(&request{
		method:      http.MethodGet,
		path:        fmt.Sprintf(GuildVanityURLFmt, guild),
		contentType: JsonContentType,
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

	img, err := png.Decode(bytes.NewReader(res.Body))
	if err != nil {
		return nil, err
	}

	return img, nil
}

func (c *Client) GetGuildMembershipScreeningForm(guild objects.Snowflake) (*objects.MembershipScreening, error) {
	res, err := c.request(&request{
		method:      http.MethodGet,
		path:        fmt.Sprintf(GuildMembershipScreeningFmt, guild),
		contentType: JsonContentType,
	})
	if err != nil {
		return nil, err
	}

	if err = res.ExpectsStatus(http.StatusOK); err != nil {
		return nil, err
	}

	form := &objects.MembershipScreening{}
	if err = res.JSON(form); err != nil {
		return nil, err
	}

	return form, err
}

type ModifyGuildMembershipScreeningParams struct {
	Enabled     *bool  `json:"enabled,omitempty"`
	FormFields  string `json:"form_fields,omitempty"`
	Description string `json:"description,omitempty"`
	Reason      string `json:"-"`
}

func (c *Client) ModifyGuildMembershipScreeningForm(guild objects.Snowflake, params *ModifyGuildMembershipScreeningParams) (*objects.MembershipScreening, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	reason := ""
	if params != nil {
		reason = params.Reason
	}

	res, err := c.request(&request{
		method:      http.MethodPatch,
		path:        fmt.Sprintf(GuildMembershipScreeningFmt, guild),
		contentType: JsonContentType,
		reason:      reason,
		body:        data,
	})
	if err != nil {
		return nil, err
	}

	if err = res.ExpectsStatus(http.StatusOK); err != nil {
		return nil, err
	}

	form := &objects.MembershipScreening{}
	if err = res.JSON(form); err != nil {
		return nil, err
	}

	return form, nil
}
