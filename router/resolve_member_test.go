package router

import (
	"github.com/stretchr/testify/assert"
	"wumpgo.dev/wumpgo/objects"

	"testing"
)

var fullUser = objects.User{
	ID:            1,
	Username:      "test",
	Discriminator: "1234",
	Avatar:        "a_fgeiogeig",
	Bot:           true,
	System:        true,
	MFAEnabled:    true,
	Locale:        "en",
	Verified:      true,
	Email:         "test@example.com",
	Flags:         1,
	PremiumType:   1,
	PublicFlags:   1,
}

var fullMemberExceptUser = objects.GuildMember{
	Nick:    "testing",
	Roles:   []objects.Snowflake{1},
	Deaf:    true,
	Mute:    true,
	Pending: true,
}

var fullMember *objects.GuildMember

func init() {
	x := fullMemberExceptUser
	userCpy := fullUser
	x.User = &userCpy
	fullMember = &x
}

func TestResolvableUser_ResolveMember(t *testing.T) {
	tests := []struct {
		name string

		data   *objects.ApplicationCommandInteractionData
		member *objects.GuildMember
	}{
		{
			name: "nil member",
			data: &objects.ApplicationCommandInteractionData{},
		},
		{
			name: "nil user",
			data: &objects.ApplicationCommandInteractionData{
				Resolved: objects.ApplicationCommandInteractionDataResolved{
					Users: map[objects.Snowflake]objects.User{},
					Members: map[objects.Snowflake]objects.GuildMember{
						1: {Nick: "a"},
					},
				},
			},
			member: &objects.GuildMember{Nick: "a"},
		},
		{
			name: "full",
			data: &objects.ApplicationCommandInteractionData{
				Resolved: objects.ApplicationCommandInteractionDataResolved{
					Users:   map[objects.Snowflake]objects.User{1: fullUser},
					Members: map[objects.Snowflake]objects.GuildMember{1: fullMemberExceptUser},
				},
			},
			member: fullMember,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolvable := resolvableUser{resolvable: resolvable[objects.User]{
				id: "1", data: tt.data,
			}}
			assert.Equal(t, tt.member, resolvable.ResolveMember())
		})
	}
}
