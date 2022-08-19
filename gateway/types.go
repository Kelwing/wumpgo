package gateway

type OpCode int

const (
	OpDispatch OpCode = iota
	OpHeartbeat
	OpIdentify
	OpPresenceUpdate
	OpVoiceStateUpdate
	_
	OpResume
	OpReconnect
	OpRequestGuildMembers
	OpInvalidSession
	OpHello
	OpHeartbeatACK
)

type StatusType string

const (
	StatusOnline       StatusType = "online"
	StatusDoNotDisturb StatusType = "dnd"
	StatusAFK          StatusType = "idle"
	StatusInvisible    StatusType = "invisible"
	StatusOffline      StatusType = "offline"
)

type Intent int

const (
	IntentsGuilds Intent = 1 << iota
	IntentsGuildMembers
	IntentsGuildBans
	IntentsGuildEmojis
	IntentsGuildIntegrations
	IntentsGuildWebhooks
	IntentsGuildInvites
	IntentsGuildVoiceStates
	IntentsGuildPresences
	IntentsGuildMessages
	IntentsGuildMessageReactions
	IntentsGuildMessageTyping
	IntentsDirectMessages
	IntentsDirectMessageReactions
	IntentsDirectMessageTyping
	IntentMessageContent
	IntentGuildScheduledEvents
	_
	_
	_
	IntentAutoModerationConfiguration
	IntentAutoModerationExecution
)

const (
	IntentsAllWithoutPrivileged = IntentsGuilds |
		IntentsGuildBans |
		IntentsGuildEmojis |
		IntentsGuildIntegrations |
		IntentsGuildWebhooks |
		IntentsGuildInvites |
		IntentsGuildVoiceStates |
		IntentsGuildMessages |
		IntentsGuildMessageReactions |
		IntentsGuildMessageTyping |
		IntentsDirectMessages |
		IntentsDirectMessageReactions |
		IntentsDirectMessageTyping |
		IntentGuildScheduledEvents |
		IntentAutoModerationConfiguration |
		IntentAutoModerationExecution
	IntentsAll = IntentsAllWithoutPrivileged |
		IntentsGuildMembers |
		IntentsGuildPresences
	IntentsNone Intent = 0
)
