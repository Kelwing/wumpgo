package objects

type (
	Ready struct {
		Version     int          `json:"v"`
		User        *User        `json:"user"`
		Guilds      []*Guild     `json:"guilds"`
		SessionID   string       `json:"session_id"`
		Shard       [2]int       `json:"shard"`
		Application *Application `json:"application"`
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

	ChannelPinsUpdate struct {
		GuildID          Snowflake `json:"guild_id"`
		ChannelID        Snowflake `json:"channel_id"`
		LastPinTimestamp Time      `json:"last_pin_timestamp"`
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
		GuildID      Snowflake   `json:"guild_id"`
		Roles        []Snowflake `json:"roles"`
		User         *User       `json:"user"`
		Nick         string      `json:"nick"`
		JoinedAt     Time        `json:"joined_at"`
		PremiumSince Time        `json:"premium_since"`
		Deaf         bool        `json:"deaf"`
		Mute         bool        `json:"mute"`
		Pending      bool        `json:"pending"`
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

	Gateway struct {
		URL               string `json:"url"`
		Shards            int    `json:"shards"`
		SessionStartLimit struct {
			Total          int `json:"total"`
			Remaining      int `json:"remaining"`
			ResetAfter     int `json:"reset_after"`
			MaxConcurrency int `json:"max_concurrency"`
		} `json:"session_start_limit"`
	}
)
