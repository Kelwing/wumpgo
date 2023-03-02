package objects

var _ Mentionable = (*GuildMember)(nil)

type GuildMember struct {
	User                       *User           `json:"user,omitempty"`
	Nick                       string          `json:"nick,omitempty"`
	Avatar                     string          `json:"avatar,omitempty"`
	Roles                      []Snowflake     `json:"roles"`
	JoinedAt                   Time            `json:"joined_at"`
	PremiumSince               Time            `json:"premium_since,omitempty"`
	Deaf                       bool            `json:"deaf"`
	Mute                       bool            `json:"mute"`
	Flags                      GuildMemberFlag `json:"flags,omitempty"`
	Pending                    bool            `json:"pending,omitempty"`
	Permissions                string          `json:"permissions,omitempty"`
	CommunicationDisabledUntil *Time           `json:"communication_disabled_until,omitempty"`
}

func (m *GuildMember) Mention() string {
	return m.User.Mention()
}

type GuildMemberFlag uint

const (
	DidRejoin GuildMemberFlag = 1 << iota
	CompletedOnboarding
	BypassesVerification
	StartedOnboarding
)
