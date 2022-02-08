package router

import (
	"github.com/Postcord/objects"
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestResolvableMentionable_Snowflake(t *testing.T) {
	testSnowflake(t, &ResolvableMentionable{})
}

func TestResolvableMentionable_MarshalJSON(t *testing.T) {
	testMarshalJSON(t, &ResolvableMentionable{})
}

func TestResolvableMentionable_String(t *testing.T) {
	testString(t, &ResolvableMentionable{})
}

func TestResolvableMentionable_Resolve(t *testing.T) {
	tests := []struct {
		name string

		data     *objects.ApplicationCommandInteractionData
		expected interface{}
	}{
		{
			name: "nil result",
		},
		{
			name: "channel",
			data: &objects.ApplicationCommandInteractionData{
				Resolved: objects.ApplicationCommandInteractionDataResolved{
					Users:    map[objects.Snowflake]objects.User{1: {DiscordBaseObject: objects.DiscordBaseObject{ID: 1}}},
					Members:  map[objects.Snowflake]objects.GuildMember{1: {Nick: "abc"}},
					Roles:    map[objects.Snowflake]objects.Role{1: {DiscordBaseObject: objects.DiscordBaseObject{ID: 123}}},
					Channels: map[objects.Snowflake]objects.Channel{1: {DiscordBaseObject: objects.DiscordBaseObject{ID: 123}}},
				},
			},
			expected: &objects.Channel{DiscordBaseObject: objects.DiscordBaseObject{ID: 123}},
		},
		{
			name: "role",
			data: &objects.ApplicationCommandInteractionData{
				Resolved: objects.ApplicationCommandInteractionDataResolved{
					Users:    map[objects.Snowflake]objects.User{1: {DiscordBaseObject: objects.DiscordBaseObject{ID: 1}}},
					Members:  map[objects.Snowflake]objects.GuildMember{1: {Nick: "abc"}},
					Roles:    map[objects.Snowflake]objects.Role{1: {DiscordBaseObject: objects.DiscordBaseObject{ID: 123}}},
					Channels: map[objects.Snowflake]objects.Channel{},
				},
			},
			expected: &objects.Role{DiscordBaseObject: objects.DiscordBaseObject{ID: 123}},
		},
		{
			name: "member",
			data: &objects.ApplicationCommandInteractionData{
				Resolved: objects.ApplicationCommandInteractionDataResolved{
					Users:    map[objects.Snowflake]objects.User{1: {DiscordBaseObject: objects.DiscordBaseObject{ID: 1}}},
					Members:  map[objects.Snowflake]objects.GuildMember{1: {Nick: "abc"}},
					Roles:    map[objects.Snowflake]objects.Role{},
					Channels: map[objects.Snowflake]objects.Channel{},
				},
			},
			expected: &objects.GuildMember{Nick: "abc", User: &objects.User{DiscordBaseObject: objects.DiscordBaseObject{ID: 1}}},
		},
		{
			name: "user",
			data: &objects.ApplicationCommandInteractionData{
				Resolved: objects.ApplicationCommandInteractionDataResolved{
					Users:    map[objects.Snowflake]objects.User{1: {DiscordBaseObject: objects.DiscordBaseObject{ID: 1}}},
					Members:  map[objects.Snowflake]objects.GuildMember{},
					Roles:    map[objects.Snowflake]objects.Role{},
					Channels: map[objects.Snowflake]objects.Channel{},
				},
			},
			expected: &objects.User{DiscordBaseObject: objects.DiscordBaseObject{ID: 1}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.data == nil {
				tt.data = &objects.ApplicationCommandInteractionData{}
			}
			r := ResolvableMentionable{id: "1", data: tt.data}
			assert.Equal(t, tt.expected, r.Resolve())
		})
	}
}
