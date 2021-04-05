package objects

import (
	"encoding/json"
	"strconv"
)

type PermissionBit uint64

var _ json.Marshaler = (*PermissionBit)(nil)
var _ json.Unmarshaler = (*PermissionBit)(nil)

func (p *PermissionBit) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}

func (p *PermissionBit) UnmarshalJSON(bytes []byte) error {
	var perms string
	if err := json.Unmarshal(bytes, &perms); err != nil {
		return err
	}

	v, err := strconv.ParseUint(perms, 10, 64)
	if err != nil {
		return err
	}

	*p = PermissionBit(v)
	return nil
}

func (p PermissionBit) String() string {
	return strconv.FormatUint(uint64(p), 10)
}

func (p PermissionBit) Has(bits PermissionBit) bool {
	return (p & bits) == bits
}

func (p PermissionBit) HasOrAdmin(bits PermissionBit) bool {
	return (p&bits) == bits || (p&ADMINISTRATOR) == ADMINISTRATOR
}

const (
	CREATE_INSTANT_INVITE PermissionBit = 1 << iota
	KICK_MEMBERS
	BAN_MEMBERS
	ADMINISTRATOR
	MANAGE_CHANNELS
	MANAGE_GUILD
	ADD_REACTIONS
	VIEW_AUDIT_LOG
	PRIORITY_SPEAKER
	STREAM
	VIEW_CHANNEL
	SEND_MESSAGES
	SEND_TTS_MESSAGES
	MANAGE_MESSAGES
	EMBED_LINKS
	ATTACH_FILES
	READ_MESSAGE_HISTORY
	MENTION_EVERYONE
	USE_EXTERNAL_EMOJIS
	VIEW_GUILD_INSIGHTS
	CONNECT
	SPEAK
	MUTE_MEMBERS
	DEAFEN_MEMBERS
	MOVE_MEMBERS
	USE_VAD
	CHANGE_NICKNAME
	MANAGE_NICKNAMES
	MANAGE_ROLES
	MANAGE_WEBHOOKS
	MANAGE_EMOJIS
	USE_APPLICATION_COMMANDS
)
