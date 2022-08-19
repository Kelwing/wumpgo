package gateway

import (
	"encoding/json"

	"wumpgo.dev/wumpgo/objects"
)

type Payload struct {
	Op        OpCode          `json:"op"`
	Data      json.RawMessage `json:"d"`
	Sequence  int64           `json:"s,omitempty"`
	EventName string          `json:"t,omitempty"`
}

type Hello struct {
	HeartbeatInterval int64 `json:"heartbeat_interval"`
}

type UpdatePresence struct {
	Since      int64              `json:"since"`
	Activities []objects.Activity `json:"activities"`
	Status     StatusType         `json:"status"`
	AFK        bool               `json:"afk"`
}

type Identify struct {
	Token          string         `json:"token"`
	Properties     Properties     `json:"properties"`
	LargeThreshold int64          `json:"large_threshold"`
	Compress       bool           `json:"compress"`
	Shard          []int          `json:"shard"`
	Presence       UpdatePresence `json:"presence"`
	Intents        Intent         `json:"intents"`
}

type Properties struct {
	OS      string `json:"$os"`
	Browser string `json:"$browser"`
	Device  string `json:"$device"`
}

type Resume struct {
	Token     string `json:"token"`
	SessionID string `json:"session_id"`
	Sequence  int64  `json:"seq"`
}
