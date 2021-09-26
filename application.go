package objects

const (
	ApplicationFlag_GATEWAY_PRESENCE = 1 << (iota + 12)
	ApplicationFlag_GATEWAY_PRESENCE_LIMITED
	ApplicationFlag_GATEWAY_GUILD_MEMBERS
	ApplicationFlag_GATEWAY_GUILD_MEMBERS_LIMITED
	ApplicationFlag_VERIFICATION_PENDING_GUILD_LIMIT
	ApplicationFlag_EMBEDDED
)

type Application struct {
	ID                  Snowflake `json:"id"`
	Name                string    `json:"name"`
	Icon                string    `json:"icon"`
	Description         string    `json:"description"`
	RPCOrigins          []string  `json:"rpc_origins"`
	BotPublic           bool      `json:"bot_public"`
	BotRequireCodeGrant bool      `json:"bot_require_code_grant"`
	TermsOfServiceURL   string    `json:"terms_of_service_url"`
	PrivacyPolicyURL    string    `json:"privacy_policy_url"`
	Owner               *User     `json:"owner"`
	Summary             string    `json:"summary"`
	VerifyKey           string    `json:"verify_key"`
	Team                *Team     `json:"team"`
	GuildID             Snowflake `json:"guild_id"`
	PrimarySKUID        Snowflake `json:"primary_sku_id"`
	Slug                string    `json:"slug"`
	CoverImage          string    `json:"cover_image"`
	Flags               int       `json:"flags"`
}

type TeamMembershipState int

const (
	TeamMembershipState_INVITED TeamMembershipState = iota + 1
	TeamMembershipState_ACCEPTED
)

type Team struct {
	Icon        string        `json:"icon"`
	ID          Snowflake     `json:"id"`
	Members     []*TeamMember `json:"members"`
	Name        string        `json:"name"`
	OwnerUserID Snowflake     `json:"owner_user_id"`
}

type TeamMember struct {
	MembershipState TeamMembershipState `json:"membership_state"`
	Permissions     []string            `json:"permissions"`
	TeamID          Snowflake           `json:"team_id"`
	User            *User               `json:"user"`
}
