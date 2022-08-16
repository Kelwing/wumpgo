package objects

var _ Mentionable = (*GuildMember)(nil)
var _ SnowflakeObject = (*GuildMember)(nil)

type GuildMember struct {
	User                       *User       `json:"user,omitempty"`
	Nick                       string      `json:"nick,omitempty"`
	Avatar                     string      `json:"avatar,omitempty"`
	Roles                      []Snowflake `json:"roles"`
	JoinedAt                   Time        `json:"joined_at"`
	PremiumSince               Time        `json:"premium_since,omitempty"`
	Deaf                       bool        `json:"deaf"`
	Mute                       bool        `json:"mute"`
	Pending                    bool        `json:"pending,omitempty"`
	Permissions                string      `json:"permissions,omitempty"`
	CommunicationDisabledUntil *Time       `json:"communication_disabled_until,omitempty"`
}

func (m *GuildMember) GetID() Snowflake {
	return m.User.GetID()
}

func (m *GuildMember) Mention() string {
	return m.User.Mention()
}
