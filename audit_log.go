package objects

//go:generate stringer -type=AuditLogEvent -trimprefix=AuditLogEvent -output audit_log_string.go

type AuditLogEvent uint

const (
	AuditLogEvtGuildUpdate AuditLogEvent = 1
)
const (
	AuditLogEventChannelCreate AuditLogEvent = iota + 10
	AuditLogEventChannelUpdate
	AuditLogEventChannelDelete
	AuditLogEventOverwriteCreate
	AuditLogEventOverwriteUpdate
	AuditLogEventOverwriteDelete
)

const (
	AuditLogEventMemberKick AuditLogEvent = iota + 20
	AuditLogEventMemberPrune
	AuditLogEventMemberBanAdd
	AuditLogEventMemberBanRemove
	AuditLogEventMemberUpdate
	AuditLogEventMemberRoleUpdate
	AuditLogEventMemberMove
	AuditLogEventMemberDisconnect
	AuditLogEventBotAdd
)

const (
	AuditLogEventRoleCreate AuditLogEvent = iota + 30
	AuditLogEventRoleUpdate
	AuditLogEventRoleDelete
)

const (
	AuditLogEventInviteCreate AuditLogEvent = iota + 40
	AuditLogEventInviteUpdate
	AuditLogEventInviteDelete
)

const (
	AuditLogEventWebhookCreate AuditLogEvent = iota + 50
	AuditLogEventWebhookUpdate
	AuditLogEventWebhookDelete
)

const (
	AuditLogEventEmojiCreate AuditLogEvent = iota + 60
	AuditLogEventEmojiUpdate
	AuditLogEventEmojiDelete
)

const (
	AuditLogEventMessageDelete AuditLogEvent = iota + 72
	AuditLogEventMessageBulkDelete
	AuditLogEventMessagePin
	AuditLogEventMessageUnpin
)

const (
	AuditLogEventIntegrationCreate AuditLogEvent = iota + 80
	AuditLogEventIntegrationUpdate
	AuditLogEventIntegrationDelete
)

type AuditLogChangeKey string

const (
	AuditLogChangeAFKChannelID                AuditLogChangeKey = "afk_channel_id"
	AuditLogChangeAFKTimeout                  AuditLogChangeKey = "afk_timeout"
	AuditLogChangeAdd                         AuditLogChangeKey = "$add"
	AuditLogChangeAllow                       AuditLogChangeKey = "allow"
	AuditLogChangeApplicationID               AuditLogChangeKey = "application_id"
	AuditLogChangeAvatarHash                  AuditLogChangeKey = "avatar_hash"
	AuditLogChangeBannerHash                  AuditLogChangeKey = "banner_hash"
	AuditLogChangeBitrate                     AuditLogChangeKey = "bitrate"
	AuditLogChangeChannelID                   AuditLogChangeKey = "channel_id"
	AuditLogChangeCode                        AuditLogChangeKey = "code"
	AuditLogChangeColor                       AuditLogChangeKey = "color"
	AuditLogChangeDeaf                        AuditLogChangeKey = "deaf"
	AuditLogChangeDefaultMessageNotifications AuditLogChangeKey = "default_message_notifications"
	AuditLogChangeDeny                        AuditLogChangeKey = "deny"
	AuditLogChangeDescription                 AuditLogChangeKey = "description"
	AuditLogChangeDiscoverySplashHash         AuditLogChangeKey = "discovery_splash_hash"
	AuditLogChangeEnableEmoticons             AuditLogChangeKey = "enable_emoticons"
	AuditLogChangeExpireBehaviour             AuditLogChangeKey = "expire_behaviour"
	AuditLogChangeExpireGracePeriod           AuditLogChangeKey = "expire_grace_period"
	AuditLogChangeExplicitContentFilter       AuditLogChangeKey = "explicit_content_filter"
	AuditLogChangeHoist                       AuditLogChangeKey = "hoist"
	AuditLogChangeID                          AuditLogChangeKey = "id"
	AuditLogChangeIconHash                    AuditLogChangeKey = "icon_hash"
	AuditLogChangeInviterID                   AuditLogChangeKey = "inviter_id"
	AuditLogChangeMFALevel                    AuditLogChangeKey = "mfa_level"
	AuditLogChangeMaxAge                      AuditLogChangeKey = "max_age"
	AuditLogChangeMaxUses                     AuditLogChangeKey = "max_uses"
	AuditLogChangeMentionable                 AuditLogChangeKey = "mentionable"
	AuditLogChangeMute                        AuditLogChangeKey = "mute"
	AuditLogChangeNSFW                        AuditLogChangeKey = "nsfw"
	AuditLogChangeName                        AuditLogChangeKey = "name"
	AuditLogChangeNick                        AuditLogChangeKey = "nick"
	AuditLogChangeOwnerID                     AuditLogChangeKey = "owner_id"
	AuditLogChangePermissionOverwrites        AuditLogChangeKey = "permission_overwrites"
	AuditLogChangePermissions                 AuditLogChangeKey = "permissions"
	AuditLogChangePosition                    AuditLogChangeKey = "position"
	AuditLogChangePreferredLocale             AuditLogChangeKey = "preferred_locale"
	AuditLogChangePruneDeleteDays             AuditLogChangeKey = "prune_delete_days"
	AuditLogChangePublicUpdatesChannelID      AuditLogChangeKey = "public_updates_channel_id"
	AuditLogChangeRateLimitPerUser            AuditLogChangeKey = "rate_limit_per_user"
	AuditLogChangeRegion                      AuditLogChangeKey = "region"
	AuditLogChangeRemove                      AuditLogChangeKey = "$remove"
	AuditLogChangeRuleChannelID               AuditLogChangeKey = "rules_channel_id"
	AuditLogChangeSplashHash                  AuditLogChangeKey = "splash_hash"
	AuditLogChangeSystemChannelID             AuditLogChangeKey = "system_channel_id"
	AuditLogChangeTemporary                   AuditLogChangeKey = "temporary"
	AuditLogChangeTopic                       AuditLogChangeKey = "topic"
	AuditLogChangeType                        AuditLogChangeKey = "type"
	AuditLogChangeUserLimit                   AuditLogChangeKey = "user_limit"
	AuditLogChangeUses                        AuditLogChangeKey = "uses"
	AuditLogChangeVanityURLCode               AuditLogChangeKey = "vanity_url_code"
	AuditLogChangeVerificationLevel           AuditLogChangeKey = "verification_level"
	AuditLogChangeWidgetChannelID             AuditLogChangeKey = "widget_channel_id"
	AuditLogChangeWidgetEnabled               AuditLogChangeKey = "widget_enabled"
)

// AuditLogEntry represents a single audit log.
type AuditLogEntry struct {
	TargetID Snowflake         `json:"target_id"`
	Changes  []*AuditLogChange `json:"changes,omitempty"`
	UserID   Snowflake         `json:"user_id"`
	ID       Snowflake         `json:"id"`
	Event    AuditLogEvent     `json:"action_type"`
	Options  *AuditLogOption   `json:"options,omitempty"`
	Reason   string            `json:"reason,omitempty"`
}

// AuditLogChange is the struct representing changes made to the target ID.
// More details can be found at https://discord.com/developers/docs/resources/audit-log#audit-log-change-object-audit-log-change-structure
type AuditLogChange struct {
	NewValue interface{} `json:"new_value,omitempty"`
	OldValue interface{} `json:"old_value,omitempty"`
	Key      string      `json:"key"`
}

// AuditLogOptions is the options for an audit log entry.
// More details can be found at https://discord.com/developers/docs/resources/audit-log#audit-log-entry-object-optional-audit-entry-info
type AuditLogOption struct {

	// number of days after which inactive members were kicked
	// triggered on MEMBER_PRUNE actions
	DeleteMemberDays string `json:"delete_member_days"`

	// number of members removed by the prune
	// triggered on MEMBER_PRUNE actions
	MembersRemoved string `json:"members_removed"`

	// channel in which the entities were targeted
	// triggered on MEMBER_MOVE & MESSAGE_PIN & MESSAGE_UNPIN & MESSAGE_DELETE actions
	ChannelID Snowflake `json:"channel_id"`

	// id of the message that was targeted
	// triggered for MESSAGE_PIN & MESSAGE_UNPIN actions
	MessageID Snowflake `json:"message_id"`

	// number of entities that were targeted
	// triggered on MESSAGE_DELETE & MESSAGE_BULK_DELETE & MEMBER_DISCONNECT & MEMBER_MOVE actions
	Count string `json:"count"`

	// id of the overwritten entity
	// triggered on CHANNEL_OVERWRITE_CREATE & CHANNEL_OVERWRITE_UPDATE & CHANNEL_OVERWRITE_DELETE actions
	ID Snowflake `json:"id"`

	// type of overwritten entity - "0" for "role" or "1" for "member"
	// triggered on CHANNEL_OVERWRITE_CREATE & CHANNEL_OVERWRITE_UPDATE & CHANNEL_OVERWRITE_DELETE actions
	Type string `json:"type"`

	// name of the role if type is "0" (not present if type is "1")
	// triggered on CHANNEL_OVERWRITE_CREATE & CHANNEL_OVERWRITE_UPDATE & CHANNEL_OVERWRITE_DELETE actions
	RoleName string `json:"role_name"`
}

type AuditLog struct {
	Webhooks      []*Webhook       `json:"webhooks"`
	Users         []*User          `json:"users"`
	AuditLogEntry []*AuditLogEntry `json:"audit_log_entry"`
	Integrations  []*Integration
}
