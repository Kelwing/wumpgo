package objects

var _ SnowflakeObject = (*DiscordBaseObject)(nil)

type SnowflakeObject interface {
	GetID() Snowflake
}

type DiscordBaseObject struct {
	ID Snowflake `json:"id,omitempty"`
}

func (d DiscordBaseObject) GetID() Snowflake {
	return d.ID
}

type Mentionable interface {
	Mention() string
}
