package objects

import (
	"encoding/json"
)

type (
	Hello struct {
		HeartbeatInterval int64 `json:"heartbeat_interval"`
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

	Payload struct {
		Op        OpCode          `json:"op"`
		Data      json.RawMessage `json:"d"`
		Sequence  uint64          `json:"s,omitempty"`
		EventName string          `json:"t,omitempty"`
	}

	UpdatePresence struct {
		Since      int64      `json:"since"`
		Activities []Activity `json:"activities"`
		Status     StatusType `json:"status"`
		AFK        bool       `json:"afk"`
	}

	Identify struct {
		Token          string         `json:"token"`
		Properties     Properties     `json:"properties"`
		LargeThreshold int64          `json:"large_threshold"`
		Compress       bool           `json:"compress"`
		Shard          []int          `json:"shard"`
		Presence       UpdatePresence `json:"presence"`
		Intents        Intent         `json:"intents"`
	}

	Properties struct {
		OS      string `json:"$os"`
		Browser string `json:"$browser"`
		Device  string `json:"$device"`
	}

	Resume struct {
		Token     string `json:"token"`
		SessionID string `json:"session_id"`
		Sequence  uint64 `json:"seq"`
	}
)

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
	IntentsGuildModeration
	IntentsGuildEmojisAndStickers
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
	IntentsMessageContent
	IntentsGuildScheduledEvents
	_
	_
	_
	IntentsAutoModerationConfiguration
	IntentsAutoModerationExecution
)

const (
	IntentsAllWithoutPrivileged = IntentsGuilds |
		IntentsGuildModeration |
		IntentsGuildEmojisAndStickers |
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
		IntentsGuildScheduledEvents |
		IntentsAutoModerationConfiguration |
		IntentsAutoModerationExecution
	IntentsAll = IntentsAllWithoutPrivileged |
		IntentsGuildMembers |
		IntentsGuildPresences |
		IntentsMessageContent
	IntentsNone Intent = 0
)
