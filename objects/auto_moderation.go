package objects

//go:generate stringer -type=AutoModerationEventType,AutoModerationTriggerType,AutoModerationKeywordPresetType,AutoModerationActionType -trimprefix=AutoModeration -output auto_moderation_string.go

type AutoModerationRule struct {
	ID              Snowflake                 `json:"id"`
	GuildID         Snowflake                 `json:"guild_id"`
	Name            string                    `json:"name"`
	CreatorID       Snowflake                 `json:"creator_id"`
	EventType       AutoModerationEventType   `json:"event_type"`
	TriggerType     AutoModerationTriggerType `json:"trigger_type"`
	TriggerMetadata TriggerMetadata           `json:"trigger_metadata"`
	Actions         []*AutoModerationAction   `json:"actions"`
	Enabled         bool                      `json:"enabled"`
	ExemptRoles     []Snowflake               `json:"exempt_roles"`
	ExemptChannel   []Snowflake               `json:"exempt_channels"`
}

type AutoModerationEventType uint64

const (
	AutoModerationEventTypeMessageSend AutoModerationEventType = iota + 1
)

type AutoModerationTriggerType uint64

const (
	AutoModerationTriggerTypeKeyword AutoModerationTriggerType = iota + 1
	AutoModerationTriggerTypeSpam
	AutoModerationTriggerTypeKeywordPreset
	AutoModerationTriggerTypeMentionSpam
)

type TriggerMetadata struct {
	KeywordFilter     []string                          `json:"keyword_filter,omitempty"`
	Presets           []AutoModerationKeywordPresetType `json:"presets,omitempty"`
	AllowList         []string                          `json:"allow_list,omitempty"`
	MentionTotalLimit int                               `json:"mention_total_limit"`
}

type AutoModerationKeywordPresetType uint64

const (
	AutoModerationKeywordPresetTypeProfanity AutoModerationKeywordPresetType = iota + 1
	AutoModerationKeywordPresetTypeSexualContent
	AutoModerationKeywordPresetTypeSlurs
)

type AutoModerationAction struct {
	Type     AutoModerationActionType      `json:"type"`
	Metadata *AutoModerationActionMetadata `json:"metadata,omitempty"`
}

type AutoModerationActionType uint64

const (
	AutoModerationActionTypeBlockMessage AutoModerationActionType = iota + 1
	AutoModerationActionTypeSendAlertMessage
	AutoModerationActionTypeTimeout
)

type AutoModerationActionMetadata struct {
	ChannelID       Snowflake `json:"channel_id,omitempty"`
	DurationSeconds int64     `json:"duration_seconds,omitempty"`
}

type AutoModerationTriggerMetadata struct {
	KeywordFilter     []string                          `json:"keyword_filter,omitempty"`
	Presets           []AutoModerationKeywordPresetType `json:"presets,omitempty"`
	AllowList         []string                          `json:"allow_list,omitempty"`
	MentionTotalLimit int                               `json:"mention_total_limit,omitempty"`
}
