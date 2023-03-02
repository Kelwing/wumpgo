package objects

type (
	Ready struct {
		Version          int          `json:"v"`
		User             *User        `json:"user"`
		Guilds           []*Guild     `json:"guilds"`
		SessionID        string       `json:"session_id"`
		ResumeGatewayURL string       `json:"resume_gateway_url"`
		Shard            [2]int       `json:"shard"`
		Application      *Application `json:"application"`
	}

	ApplicationCommandPermissionsUpdate struct {
		*ApplicationCommandPermissions
	}

	AutoModerationRuleCreate struct {
		*AutoModerationRule
	}

	AutoModerationRuleUpdate struct {
		*AutoModerationRule
	}

	AutoModerationRuleDelete struct {
		*AutoModerationRule
	}

	AutoModerationActionExecution struct {
		GuildID              Snowflake                 `json:"guild_id"`
		Action               *AutoModerationAction     `json:"action"`
		RuleID               Snowflake                 `json:"rule_id"`
		RuleTriggerType      AutoModerationTriggerType `json:"rule_trigger_type"`
		UserID               Snowflake                 `json:"user_id"`
		ChannelID            *Snowflake                `json:"channel_id"`
		MessageID            *Snowflake                `json:"message_id"`
		AlertSystemMessageID *Snowflake                `json:"alert_system_message_id"`
		Content              string                    `json:"content"`
		MatchedKeyword       *string                   `json:"matched_keyword"`
		MatchedContent       *string                   `json:"matched_content"`
	}

	ChannelCreate struct {
		*Channel
	}

	ChannelUpdate struct {
		*Channel
	}

	ChannelDelete struct {
		*Channel
	}

	ChannelPinsUpdate struct {
		GuildID          Snowflake `json:"guild_id"`
		ChannelID        Snowflake `json:"channel_id"`
		LastPinTimestamp Time      `json:"last_pin_timestamp"`
	}

	ThreadCreate struct {
		*Channel
		*ThreadMember
	}

	ThreadUpdate struct {
		*Channel
	}

	ThreadDelete struct {
		*Channel
	}

	ThreadListSync struct {
		GuildID    Snowflake       `json:"guild_id"`
		ChannelIDs []Snowflake     `json:"channel_ids"`
		Threads    []*Channel      `json:"threads"`
		Members    []*ThreadMember `json:"members"`
	}

	ThreadMemberUpdate struct {
		*ThreadMember
	}

	ThreadMembersUpdate struct {
		ID               Snowflake       `json:"id"`
		GuildID          Snowflake       `json:"guild_id"`
		MemberCount      int             `json:"member_count"`
		AddedMembers     []*ThreadMember `json:"added_members"`
		RemovedMemberIDs []Snowflake     `json:"removed_member_ids"`
	}

	GuildCreate struct {
		*Guild
	}

	GuildUpdate struct {
		*Guild
	}

	GuildDelete struct {
		*Guild
	}

	GuildAuditLogEntryCreate struct {
		GuildID Snowflake `json:"guild_id"`
		*AuditLogEntry
	}

	GuildBanAdd struct {
		GuildID Snowflake `json:"guild_id"`
		User    *User     `json:"user"`
	}

	GuildBanRemove struct {
		GuildID Snowflake `json:"guild_id"`
		User    *User     `json:"user"`
	}

	GuildEmojisUpdate struct {
		GuildID Snowflake `json:"guild_id"`
		Emojis  []*Emoji  `json:"emojis"`
	}

	GuildStickersUpdate struct {
		GuildID  Snowflake  `json:"guild_id"`
		Stickers []*Sticker `json:"stickers"`
	}

	GuildIntegrationsUpdate struct {
		GuildID Snowflake `json:"guild_id"`
	}

	GuildMemberAdd struct {
		GuildID Snowflake `json:"guild_id"`
		*GuildMember
	}

	GuildMemberRemove struct {
		GuildID Snowflake `json:"guild_id"`
		User    *User     `json:"user"`
	}

	GuildMemberUpdate struct {
		GuildID Snowflake `json:"guild_id"`
		*GuildMember
	}

	GuildMembersChunk struct {
		GuildID    Snowflake      `json:"guild_id"`
		Members    []*GuildMember `json:"members"`
		ChunkIndex int            `json:"chunk_index"`
		ChunkCount int            `json:"chunk_count"`
		NotFound   []Snowflake    `json:"not_found"`
		// Presences []*Presence `json:"presences"`
		Nonce string `json:"nonce"`
	}

	GuildRoleCreate struct {
		GuildID Snowflake `json:"guild_id"`
		Role    *Role     `json:"role"`
	}

	GuildRoleUpdate struct {
		GuildID Snowflake `json:"guild_id"`
		Role    *Role     `json:"role"`
	}

	GuildRoleDelete struct {
		GuildID Snowflake `json:"guild_id"`
		RoleID  Snowflake `json:"role_id"`
	}

	GuildScheduledEventCreate struct {
		*GuildScheduledEvent
	}

	GuildScheduledEventUpdate struct {
		*GuildScheduledEvent
	}

	GuildScheduledEventDelete struct {
		*GuildScheduledEvent
	}

	GuildScheduledEventUserAdd struct {
		EventID Snowflake `json:"guild_scheduled_event_id"`
		UserID  Snowflake `json:"user_id"`
		GuildID Snowflake `json:"guild_id"`
	}

	GuildScheduledEventUserRemove struct {
		EventID Snowflake `json:"guild_scheduled_event_id"`
		UserID  Snowflake `json:"user_id"`
		GuildID Snowflake `json:"guild_id"`
	}

	IntegrationCreate struct {
		GuildID Snowflake `json:"guild_id"`
		*Integration
	}

	IntegrationUpdate struct {
		GuildID Snowflake `json:"guild_id"`
		*Integration
	}

	IntegrationDelete struct {
		ID            Snowflake `json:"id"`
		GuildID       Snowflake `json:"guild_id"`
		ApplicationID Snowflake `json:"application_id"`
	}

	InviteCreate struct {
		ChannelID         Snowflake        `json:"channel_id"`
		Code              string           `json:"code"`
		CreatedAt         Time             `json:"created_at"`
		GuildID           Snowflake        `json:"guild_id"`
		Inviter           *User            `json:"inviter"`
		MaxAge            int              `json:"max_age"`
		MaxUses           int              `json:"max_uses"`
		TargetType        InviteTargetType `json:"target_type"`
		TargetUser        *User            `json:"target_user"`
		TargetApplication *Application     `json:"target_application"`
		Temporary         bool             `json:"temporary"`
		Uses              int              `json:"uses"`
	}

	InviteDelete struct {
		ChannelID Snowflake `json:"channel_id"`
		GuildID   Snowflake `json:"guild_id"`
		Code      string    `json:"code"`
	}

	MessageCreate struct {
		*Message
	}

	MessageUpdate struct {
		*Message
	}

	MessageDelete struct {
		ID        Snowflake `json:"id"`
		ChannelID Snowflake `json:"channel_id"`
		GuildID   Snowflake `json:"guild_id"`
	}

	MessageDeleteBulk struct {
		IDs       []Snowflake `json:"ids"`
		ChannelID Snowflake   `json:"channel_id"`
		GuildID   Snowflake   `json:"guild_id"`
	}

	MessageReactionAdd struct {
		UserID    Snowflake    `json:"user_id"`
		ChannelID Snowflake    `json:"channel_id"`
		MessageID Snowflake    `json:"message_id"`
		GuildID   Snowflake    `json:"guild_id"`
		Member    *GuildMember `json:"member"`
		Emoji     *Emoji       `json:"emoji"`
	}

	MessageReactionRemove struct {
		UserID    Snowflake `json:"user_id"`
		ChannelID Snowflake `json:"channel_id"`
		MessageID Snowflake `json:"message_id"`
		GuildID   Snowflake `json:"guild_id"`
		Emoji     *Emoji    `json:"emoji"`
	}

	MessageReactionRemoveAll struct {
		ChannelID Snowflake `json:"channel_id"`
		MessageID Snowflake `json:"message_id"`
		GuildID   Snowflake `json:"guild_id"`
	}

	MessageReactionRemoveEmoji struct {
		ChannelID Snowflake `json:"channel_id"`
		MessageID Snowflake `json:"message_id"`
		GuildID   Snowflake `json:"guild_id"`
		Emoji     *Emoji    `json:"emoji"`
	}

	PresenceUpdate struct {
		User         *User         `json:"user"`
		GuildID      Snowflake     `json:"guild_id"`
		Status       string        `json:"status"`
		Activities   []*Activity   `json:"activities"`
		ClientStatus *ClientStatus `json:"client_status"`
	}

	TypingStart struct {
		ChannelID Snowflake    `json:"channel_id"`
		GuildID   Snowflake    `json:"guild_id"`
		UserID    Snowflake    `json:"user_id"`
		Timestamp int          `json:"timestamp"`
		Member    *GuildMember `json:"member"`
	}

	UserUpdate struct {
		*User
	}

	VoiceStateUpdate struct {
		*VoiceState
	}

	VoiceServerUpdate struct {
		Token    string    `json:"token"`
		GuildID  Snowflake `json:"guild_id"`
		Endpoint string    `json:"endpoint"`
	}

	WebhooksUpdate struct {
		ChannelID Snowflake `json:"channel_id"`
		GuildID   Snowflake `json:"guild_id"`
	}

	InteractionCreate struct {
		*Interaction
	}

	StageInstanceCreate struct {
		*StageInstance
	}

	StageInstanceUpdate struct {
		*StageInstance
	}

	StageInstanceDelete struct {
		*StageInstance
	}
)
