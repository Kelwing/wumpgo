package objects

type GuildMember struct {
	User         *User       `json:"user,omitempty"`
	Nick         string      `json:"nick,omitempty"`
	Roles        []Snowflake `json:"roles"`
	JoinedAt     Time        `json:"joined_at"`
	PremiumSince Time        `json:"premium_since,omitempty"`
	Deaf         bool        `json:"deaf"`
	Mute         bool        `json:"mute"`
	Pending      bool        `json:"pending,omitempty"`
}
