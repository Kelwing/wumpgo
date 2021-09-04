package router

import (
	"github.com/Postcord/objects"
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestResolvableMentionable_Snowflake(t *testing.T) {
	tests := []struct {
		name string

		id      string
		expects objects.Snowflake
	}{
		{
			name:    "invalid",
			id:      "this_is_not_valid",
			expects: 0,
		},
		{
			name:    "snowflake value",
			id:      "1234",
			expects: 1234,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := ResolvableMentionable{id: tt.id}
			assert.Equal(t, tt.expects, r.Snowflake())
		})
	}
}

func TestResolvableMentionable_MarshalJSON(t *testing.T) {
	r := ResolvableMentionable{id: "testing"}
	b, _ := r.MarshalJSON()
	assert.Equal(t, `"testing"`, string(b))
}

func TestResolvableMentionable_String(t *testing.T) {
	r := ResolvableMentionable{id: "testing"}
	assert.Equal(t, "testing", r.String())
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
					Users: map[objects.Snowflake]objects.User{1: {ID: 1}},
					Members: map[objects.Snowflake]objects.GuildMember{1: {Nick: "abc"}},
					Roles: map[objects.Snowflake]objects.Role{1: {ID: 123}},
					Channels: map[objects.Snowflake]objects.Channel{1: {ID: 123}},
				},
			},
			expected: &objects.Channel{ID: 123},
		},
		{
			name: "role",
			data: &objects.ApplicationCommandInteractionData{
				Resolved: objects.ApplicationCommandInteractionDataResolved{
					Users: map[objects.Snowflake]objects.User{1: {ID: 1}},
					Members: map[objects.Snowflake]objects.GuildMember{1: {Nick: "abc"}},
					Roles: map[objects.Snowflake]objects.Role{1: {ID: 123}},
					Channels: map[objects.Snowflake]objects.Channel{},
				},
			},
			expected: &objects.Role{ID: 123},
		},
		{
			name: "member",
			data: &objects.ApplicationCommandInteractionData{
				Resolved: objects.ApplicationCommandInteractionDataResolved{
					Users: map[objects.Snowflake]objects.User{1: {ID: 1}},
					Members: map[objects.Snowflake]objects.GuildMember{1: {Nick: "abc"}},
					Roles: map[objects.Snowflake]objects.Role{},
					Channels: map[objects.Snowflake]objects.Channel{},
				},
			},
			expected: &objects.GuildMember{Nick: "abc", User: &objects.User{ID: 1}},
		},
		{
			name: "user",
			data: &objects.ApplicationCommandInteractionData{
				Resolved: objects.ApplicationCommandInteractionDataResolved{
					Users: map[objects.Snowflake]objects.User{1: {ID: 1}},
					Members: map[objects.Snowflake]objects.GuildMember{},
					Roles: map[objects.Snowflake]objects.Role{},
					Channels: map[objects.Snowflake]objects.Channel{},
				},
			},
			expected: &objects.User{ID: 1},
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
