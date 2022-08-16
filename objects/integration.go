package objects

//go:generate stringer -type=ExpireBehavior -trimprefix=ExpireBehavior -output integration_string.go

var _ SnowflakeObject = (*Integration)(nil)
var _ SnowflakeObject = (*IntegrationAccount)(nil)
var _ SnowflakeObject = (*IntegrationApplication)(nil)

type ExpireBehavior uint

const (
	ExpireBehaviorRemoveRole ExpireBehavior = iota
	ExpireBehaviorKick
)

type Integration struct {
	DiscordBaseObject
	Name              string                  `json:"name"`
	Type              string                  `json:"type"`
	Enabled           bool                    `json:"enabled"`
	Syncing           bool                    `json:"syncing,omitempty"`
	RoleID            Snowflake               `json:"role_id,omitempty"`
	EnableEmoticons   bool                    `json:"enable_emoticons,omitempty"`
	ExpireBehavior    ExpireBehavior          `json:"expire_behaviour,omitempty"`
	ExpireGracePeriod int64                   `json:"expire_grace_period,omitempty"`
	User              *User                   `json:"user,omitempty"`
	Account           *IntegrationAccount     `json:"account,omitempty"`
	SyncedAt          Time                    `json:"synced_at,omitempty"`
	SubscriberCount   int64                   `json:"subscriber_count,omitempty"`
	Revoked           bool                    `json:"revoked,omitempty"`
	Application       *IntegrationApplication `json:"application,omitempty"`
}

type IntegrationAccount struct {
	DiscordBaseObject
	Name string `json:"name"`
}

type IntegrationApplication struct {
	DiscordBaseObject
	Name        string `json:"name"`
	Icon        string `json:"icon,omitempty"`
	Description string `json:"description"`
	Summary     string `json:"summary"`
	Bot         *User  `json:"bot,omitempty"`
}
