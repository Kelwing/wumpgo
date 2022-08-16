package objects

var _ SnowflakeObject = (*GuildScheduledEvent)(nil)

//go:generate stringer -type=GuildScheduledEventStatus,GuildScheduledEventEntityType -output guild_scheduled_event_string.go
type GuildScheduledEventStatus int

const (
	EventStatusScheduled GuildScheduledEventStatus = iota + 1
	EventStatusActive
	EventStatusCompleted
	EventStatusCanceled
)

type GuildScheduledEventEntityType int

const (
	EntityTypeStageInstance GuildScheduledEventEntityType = iota + 1
	EntityTypeVoice
	EntityTypeExternal
)

type GuildScheduledEventEntityMetadata struct {
	Location string `json:"location"`
}

type GuildScheduledEvent struct {
	DiscordBaseObject
	GuildID            Snowflake                          `json:"guild_id"`
	ChannelID          *Snowflake                         `json:"channel_id,omitempty"`
	Name               string                             `json:"name"`
	Description        *string                            `json:"description,omitempty"`
	ScheduledStartTime Time                               `json:"scheduled_start_time"`
	ScheduledEndTime   Time                               `json:"scheduled_end_time"`
	PrivacyLevel       PrivacyLevel                       `json:"privacy_level"`
	Status             GuildScheduledEventStatus          `json:"status"`
	EntityType         GuildScheduledEventEntityType      `json:"entity_type"`
	EntityID           *Snowflake                         `json:"entity_id,omitempty"`
	EntityMetadata     *GuildScheduledEventEntityMetadata `json:"entity_metadata,omitempty"`
	Creator            *User                              `json:"creator,omitempty"`
	UserCount          int                                `json:"user_count"`
	Image              string                             `json:"image"`
}

type GuildScheduledEventUser struct {
	ScheduledEventID Snowflake    `json:"guild_scheduled_event_id"`
	User             *User        `json:"user"`
	Member           *GuildMember `json:"member"`
}
