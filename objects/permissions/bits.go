package permissions

//go:generate stringer -type=PermissionBit -trimprefix=PermissionBit -output bits_string.go

import (
	"encoding/json"
	"strconv"
)

type PermissionBit uint64

var _ json.Marshaler = (*PermissionBit)(nil)
var _ json.Unmarshaler = (*PermissionBit)(nil)

func (p *PermissionBit) MarshalJSON() ([]byte, error) {
	return json.Marshal(strconv.FormatUint(uint64(*p), 10))
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

func (p PermissionBit) Has(bits PermissionBit) bool {
	return (p & bits) == bits
}

func (p PermissionBit) HasOrAdmin(bits PermissionBit) bool {
	return (p&bits) == bits || (p&Administrator) == Administrator
}

const (
	CreateInstantInvite PermissionBit = 1 << iota
	KickMembers
	BanMembers
	Administrator
	ManageChannels
	ManageGuild
	AddReactions
	ViewAuditLog
	PrioritySpeaker
	Stream
	ViewChannel
	SendMessages
	SendTTSMessages
	ManageMessages
	EmbedLinks
	AttachFiles
	ReadMessageHistory
	MentionEveryone
	UseExternalEmojis
	ViewGuildInsights
	Connect
	Speak
	MuteMembers
	DeafenMembers
	MoveMembers
	UseVAD
	ChangeNickname
	ManageNicknames
	ManageRoles
	ManageWebhooks
	ManageEmojisAndStickers
	UseApplicationCommands
	RequestToSpeak
	ManageEvents
	ManageThreads
	CreatePublicThreads
	CreatePrivateThreads
	UseExternalStickers
	SendMessagesInThreads
	StartEmbeddedActivities
	ModerateMembers
)
