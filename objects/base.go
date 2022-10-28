package objects

type SnowflakeObject interface {
	GetID() Snowflake
}

type Mentionable interface {
	Mention() string
}
