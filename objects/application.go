package objects

//go:generate stringer -type=ApplicationFlag,TeamMembershipState -output application_string.go

type ApplicationFlag int

const (
	ApplicationFlagGatewayPresence ApplicationFlag = 1 << (iota + 12)
	ApplicationFlagGatewayPresenceLimited
	ApplicationFlagGatewayGuildMembers
	ApplicationFlagGatewayGuildMembersLimited
	ApplicationFlagVerificationPendingGuildLimit
	ApplicationFlagEmbedded
	ApplicationFlagMessageContent
	ApplicationFlagMessageContentLimited
)

// A Discord API Application object.
// https://discord.com/developers/docs/resources/application#application-object-application-structure
type Application struct {
	ID                  Snowflake       `json:"id"`
	Name                string          `json:"name"`
	Icon                string          `json:"icon"`
	Description         string          `json:"description"`
	RPCOrigins          []string        `json:"rpc_origins"`
	BotPublic           bool            `json:"bot_public"`
	BotRequireCodeGrant bool            `json:"bot_require_code_grant"`
	TermsOfServiceURL   string          `json:"terms_of_service_url"`
	PrivacyPolicyURL    string          `json:"privacy_policy_url"`
	Owner               *User           `json:"owner"`
	Summary             string          `json:"summary"`
	VerifyKey           string          `json:"verify_key"`
	Team                *Team           `json:"team"`
	GuildID             Snowflake       `json:"guild_id"`
	PrimarySKUID        Snowflake       `json:"primary_sku_id"`
	Slug                string          `json:"slug"`
	CoverImage          string          `json:"cover_image"`
	Flags               ApplicationFlag `json:"flags"`
	Tags                []string        `json:"tags"`
	InstallParams       InstallParams   `json:"install_params"`
	CustomInstallURL    string          `json:"custom_install_url"`
}

type TeamMembershipState int

const (
	TeamMembershipStateInvited TeamMembershipState = iota + 1
	TeamMembershipStateAccepted
)

type Team struct {
	ID          Snowflake     `json:"id"`
	Icon        string        `json:"icon"`
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

type InstallParams struct {
	Scopes      []string `json:"scopes"`
	Permissions string   `json:"permissions"`
}
