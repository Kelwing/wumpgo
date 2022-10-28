package objects

//go:generate stringer -type=InviteTargetType -trimprefix=InviteTargetType -output invite_string.go

// https://discord.com/developers/docs/resources/invite#invite-object-invite-structure
type InviteTargetType int

const (
	InviteTargetTypeStream InviteTargetType = iota + 1
	InviteTargetTypeEmbeddedApplication
)

// Invite represents a code that when used, adds a user to a guild or group DM channel.
// https://discord.com/developers/docs/resources/invite#invite-object-invite-structure
type Invite struct {
	// the invite code (unique ID)
	Code string `json:"code"`

	// the guild this invite is for
	Guild *Guild `json:"guild,omitempty"`

	// the channel this invite is for
	Channel *Channel `json:"channel"`

	// the user who created the invite
	Inviter *User `json:"inviter,omitempty"`

	// The target type for this invite
	TargetType InviteTargetType `json:"target_type,omitempty"`

	// the target user for this invite
	TargetUser *User `json:"target_user,omitempty"`

	// The application the invite is targeting
	TargetApplication *Application `json:"target_application,omitempty"`

	// approximate count of online members (only present when target_user is set)
	ApproximatePresenceCount int `json:"approximate_presence_count,omitempty"`

	// approximate count of total members
	ApproximateMemberCount int `json:"approximate_member_count,omitempty"`

	// the expiration date of this invite, returned from the GET /invites/<code> endpoint when with_expiration is true
	ExpiresAt *Time `json:"expires_at,omitempty"`

	// guild scheduled event data, only included if guild_scheduled_event_id contains a valid guild scheduled event id
	GuildScheduledEvent *GuildScheduledEvent `json:"guild_scheduled_event,omitempty"`

	// number of times this invite has been used
	Uses int `json:"uses"`

	// max number of times this invite can be used
	MaxUses int `json:"max_uses"`

	// duration (in seconds) after which the invite expires
	MaxAge int `json:"max_age"`

	// whether this invite only grants temporary membership
	Temporary bool `json:"temporary"`

	// when this invite was created
	CreatedAt Time `json:"created_at"`
}
