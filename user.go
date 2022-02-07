package objects

//go:generate stringer -type UserFlags,PremiumType,ActivityType,ActivityFlag -output user_string.go

var _ Mentionable = (*User)(nil)
var _ SnowflakeObject = (*User)(nil)

type UserFlags uint

const UserFlagsNone = UserFlags(0)
const (
	DiscordEmployee UserFlags = 1 << iota
	PartneredServerOwner
	HypesquadEvents
	BugHunterLevel1
	HouseBravery
	HouseBrilliance
	HouseBalance
	EarlySupporter
	TeamUser
	_
	_
	_
	System
	_
	BugHunterLevel2
	_
	VerifiedBot
	EarlyVerifiedBotDeveloper
)

type PremiumType uint

const (
	PremiumTypeNone PremiumType = iota
	NitroClassic
	Nitro
)

type ActivityType uint

const (
	Game ActivityType = iota
	Streaming
	Listening
	Custom
	Competing
)

type ActivityFlag uint

const (
	ActivityInstance ActivityFlag = 1 << iota
	ActivityJoin
	ActivitySpectate
	ActivityJoinRequest
	ActivitySync
	ActivityPlay
)

type User struct {
	DiscordBaseObject
	Username      string      `json:"username"`
	Discriminator string      `json:"discriminator"`
	Avatar        string      `json:"avatar,omitempty"`
	Bot           bool        `json:"bot,omitempty"`
	System        bool        `json:"system,omitempty"`
	MFAEnabled    bool        `json:"mfa_enabled,omitempty"`
	Banner        string      `json:"banner,omitempty"`
	AccentColor   int         `json:"accent_color"`
	Locale        string      `json:"locale,omitempty"`
	Verified      bool        `json:"verified,omitempty"`
	Email         string      `json:"email,omitempty"`
	Flags         UserFlags   `json:"flags"`
	PremiumType   PremiumType `json:"premium_type,omitempty"`
	PublicFlags   UserFlags   `json:"public_flags,omitempty"`
}

func (u *User) Mention() string {
	return "<@" + u.ID.String() + ">"
}

type PresenceUpdate struct {
	User         *User         `json:"user"`
	GuildID      Snowflake     `json:"guild_id"`
	Status       string        `json:"status"`
	Activities   []*Activity   `json:"activities"`
	ClientStatus *ClientStatus `json:"client_status"`
}

type ClientStatus struct {
	Desktop string `json:"desktop,omitempty"`
	Mobile  string `json:"mobile,omitempty"`
	Web     string `json:"web,omitempty"`
}

type Activity struct {
	Name          string           `json:"name,omitempty"`
	Type          ActivityType     `json:"type"`
	URL           string           `json:"url,omitempty"`
	CreatedAt     int              `json:"created_at"`
	Timestamps    Timestamps       `json:"timestamps,omitempty"`
	ApplicationID Snowflake        `json:"application_id,omitempty"`
	Details       string           `json:"details,omitempty"`
	State         string           `json:"state,omitempty"`
	Emoji         *Emoji           `json:"emoji,omitempty"`
	Party         *ActivityParty   `json:"party,omitempty"`
	Assets        *ActivityAssets  `json:"assets,omitempty"`
	Secrets       *ActivitySecrets `json:"secrets,omitempty"`
	Instance      bool             `json:"instance,omitempty"`
	Flags         ActivityFlag     `json:"flags,omitempty"`
}

type Timestamps struct {
	Start int `json:"start,omitempty"`
	End   int `json:"end,omitempty"`
}

type ActivityParty struct {
	ID   string `json:"id,omitempty"`
	Size [2]int `json:"size,omitempty"`
}

type ActivityAssets struct {
	LargeImage string `json:"large_image,omitempty"`
	LargeText  string `json:"large_text,omitempty"`
	SmallImage string `json:"small_image,omitempty"`
	SmallText  string `json:"small_text,omitempty"`
}

type ActivitySecrets struct {
	Join     string `json:"join,omitempty"`
	Spectate string `json:"spectate,omitempty"`
	Match    string `json:"match,omitempty"`
}

type Visibility uint

const (
	VisibilityNone Visibility = iota
	VisibilityEveryone
)

type Connection struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Type         string         `json:"type"`
	Revoked      bool           `json:"revoked,omitempty"`
	Integrations []*Integration `json:"integrations,omitempty"`
	Verified     bool           `json:"verified"`
	FriendSync   bool           `json:"friend_sync"`
	ShowActivity bool           `json:"show_activity"`
	Visibility   Visibility     `json:"visibility"`
}
