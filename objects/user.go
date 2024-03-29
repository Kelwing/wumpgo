package objects

//go:generate stringer -type UserFlags,PremiumType,ActivityType,ActivityFlag -output user_string.go

var _ Mentionable = (*User)(nil)

type UserFlags uint

const UserFlagsNone = UserFlags(0)
const (
	DiscordEmployee UserFlags = 1 << iota
	PartneredServerOwner
	HypesquadEvents
	BugHunterLevel1
	_
	_
	HouseBravery
	HouseBrilliance
	HouseBalance
	EarlySupporter
	TeamUser
	_
	_
	_
	BugHunterLevel2
	_
	VerifiedBot
	EarlyVerifiedBotDeveloper
	DiscordCertifiedModerator
	BotHTTPInteractions
	ActiveDeveloper
)

type PremiumType uint

const (
	PremiumTypeNone PremiumType = iota
	NitroClassic
	Nitro
	NitroBasic
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
	ID            Snowflake     `json:"id"`
	Username      string        `json:"username"`
	Discriminator string        `json:"discriminator"`
	Avatar        Asset[Avatar] `json:"avatar,omitempty"`
	Bot           bool          `json:"bot,omitempty"`
	System        bool          `json:"system,omitempty"`
	MFAEnabled    bool          `json:"mfa_enabled,omitempty"`
	Banner        string        `json:"banner,omitempty"`
	AccentColor   int           `json:"accent_color"`
	Locale        string        `json:"locale,omitempty"`
	Verified      bool          `json:"verified,omitempty"`
	Email         string        `json:"email,omitempty"`
	Flags         UserFlags     `json:"flags"`
	PremiumType   PremiumType   `json:"premium_type,omitempty"`
	PublicFlags   UserFlags     `json:"public_flags,omitempty"`
}

func (u *User) Mention() string {
	return "<@" + u.ID.String() + ">"
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

type Service string

const (
	ServiceBattleNet       Service = "battlenet"
	ServiceEbay            Service = "ebay"
	ServiceEpicGames       Service = "epicgames"
	ServiceFacebook        Service = "facebook"
	ServiceGithub          Service = "github"
	ServiceLeagueOfLegends Service = "leagueoflegends"
	ServicePaypal          Service = "paypal"
	ServicePlaystation     Service = "playstation"
	ServiceReddit          Service = "reddit"
	ServiceRiotGames       Service = "riotgames"
	ServiceSpotify         Service = "spotify"
	ServiceSkype           Service = "skype"
	ServiceSteam           Service = "steam"
	ServiceTwitch          Service = "twitch"
	ServiceTwitter         Service = "twitter"
	ServiceXbox            Service = "xbox"
	ServiceYouTube         Service = "youtube"
)

type Connection struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Type         Service        `json:"type"`
	Revoked      bool           `json:"revoked,omitempty"`
	Integrations []*Integration `json:"integrations,omitempty"`
	Verified     bool           `json:"verified"`
	FriendSync   bool           `json:"friend_sync"`
	ShowActivity bool           `json:"show_activity"`
	TwoWayLink   bool           `json:"two_way_link"`
	Visibility   Visibility     `json:"visibility"`
}
