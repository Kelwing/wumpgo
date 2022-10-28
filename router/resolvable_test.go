package router

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
	"wumpgo.dev/wumpgo/objects"
)

func setUnexportedField(field reflect.Value, value any) {
	reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).
		Elem().
		Set(reflect.ValueOf(value))
}

func testSnowflake(t *testing.T, resolvable any) {
	t.Helper()
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
			setUnexportedField(reflect.Indirect(reflect.ValueOf(resolvable)).FieldByName("id"), tt.id)
			r := reflect.ValueOf(resolvable).MethodByName("Snowflake")
			if r.IsZero() {
				t.Fatal("function does not exist")
			}
			if r.Kind() != reflect.Func {
				t.Fatal("not a function")
			}
			res := r.Call([]reflect.Value{})
			if len(res) != 1 {
				t.Fatal("arg count not correct for resolver:", res)
			}
			assert.Equal(t, tt.expects, res[0].Interface())
		})
	}
}

func TestResolvableUser_Snowflake(t *testing.T) {
	testSnowflake(t, &resolvableUser{})
}

func TestResolvableChannel_Snowflake(t *testing.T) {
	testSnowflake(t, &resolvable[objects.Channel]{})
}

func TestResolvableRole_Snowflake(t *testing.T) {
	testSnowflake(t, &resolvable[objects.Role]{})
}

func TestResolvableMessage_Snowflake(t *testing.T) {
	testSnowflake(t, &resolvable[objects.Message]{})
}

func TestResolvableAttachment_Snowflake(t *testing.T) {
	testSnowflake(t, &resolvable[objects.Attachment]{})
}

func testMarshalJSON(t *testing.T, resolvable any) {
	t.Helper()
	setUnexportedField(reflect.Indirect(reflect.ValueOf(resolvable)).FieldByName("id"), "testing")
	r := reflect.ValueOf(resolvable).MethodByName("MarshalJSON")
	if r.IsZero() {
		t.Fatal("function does not exist")
	}
	if r.Kind() != reflect.Func {
		t.Fatal("not a function")
	}
	res := r.Call([]reflect.Value{})
	if len(res) != 2 {
		t.Fatal("arg count not correct for resolver:", res)
	}
	assert.Equal(t, []byte(`"testing"`), res[0].Interface())
	assert.Equal(t, nil, res[1].Interface())
}

func TestResolvableUser_MarshalJSON(t *testing.T) {
	testMarshalJSON(t, &resolvableUser{})
}

func TestResolvableChannel_MarshalJSON(t *testing.T) {
	testMarshalJSON(t, &resolvable[objects.Channel]{})
}

func TestResolvableRole_MarshalJSON(t *testing.T) {
	testMarshalJSON(t, &resolvable[objects.Role]{})
}

func TestResolvableMessage_MarshalJSON(t *testing.T) {
	testMarshalJSON(t, &resolvable[objects.Message]{})
}

func TestResolvableAttachment_MarshalJSON(t *testing.T) {
	testMarshalJSON(t, &resolvable[objects.Attachment]{})
}

func testString(t *testing.T, resolvable any) {
	t.Helper()
	setUnexportedField(reflect.Indirect(reflect.ValueOf(resolvable)).FieldByName("id"), "testing")
	r := reflect.ValueOf(resolvable).MethodByName("String")
	if r.IsZero() {
		t.Fatal("function does not exist")
	}
	if r.Kind() != reflect.Func {
		t.Fatal("not a function")
	}
	res := r.Call([]reflect.Value{})
	if len(res) != 1 {
		t.Fatal("arg count not correct for resolver:", res)
	}
	assert.Equal(t, "testing", res[0].Interface())
}

func TestResolvableUser_String(t *testing.T) {
	testString(t, &resolvableUser{})
}

func TestResolvableChannel_String(t *testing.T) {
	testString(t, &resolvable[objects.Channel]{})
}

func TestResolvableRole_String(t *testing.T) {
	testString(t, &resolvable[objects.Role]{})
}

func TestResolvableMessage_String(t *testing.T) {
	testString(t, &resolvable[objects.Message]{})
}

func TestResolvableAttachment_String(t *testing.T) {
	testString(t, &resolvable[objects.Attachment]{})
}

func TestResolvableUser_Resolve(t *testing.T) {
	tests := []struct {
		name string

		data     *objects.ApplicationCommandInteractionData
		expected *objects.User
	}{
		{
			name:     "nil",
			data:     &objects.ApplicationCommandInteractionData{},
			expected: nil,
		},
		{
			name: "found",
			data: &objects.ApplicationCommandInteractionData{
				Resolved: objects.ApplicationCommandInteractionDataResolved{
					Users: map[objects.Snowflake]objects.User{
						1: {Email: "abc"},
					},
				},
			},
			expected: &objects.User{Email: "abc"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := resolvable[objects.User]{id: "1", data: tt.data}
			assert.Equal(t, tt.expected, r.Resolve())
		})
	}
}

func TestResolvableChannel_Resolve(t *testing.T) {
	tests := []struct {
		name string

		data     *objects.ApplicationCommandInteractionData
		expected *objects.Channel
	}{
		{
			name:     "nil",
			data:     &objects.ApplicationCommandInteractionData{},
			expected: nil,
		},
		{
			name: "found",
			data: &objects.ApplicationCommandInteractionData{
				Resolved: objects.ApplicationCommandInteractionDataResolved{
					Channels: map[objects.Snowflake]objects.Channel{
						1: {Name: "abc"},
					},
				},
			},
			expected: &objects.Channel{Name: "abc"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := resolvable[objects.Channel]{id: "1", data: tt.data}
			assert.Equal(t, tt.expected, r.Resolve())
		})
	}
}

func TestResolvableRole_Resolve(t *testing.T) {
	tests := []struct {
		name string

		data     *objects.ApplicationCommandInteractionData
		expected *objects.Role
	}{
		{
			name:     "nil",
			data:     &objects.ApplicationCommandInteractionData{},
			expected: nil,
		},
		{
			name: "found",
			data: &objects.ApplicationCommandInteractionData{
				Resolved: objects.ApplicationCommandInteractionDataResolved{
					Roles: map[objects.Snowflake]objects.Role{
						1: {Name: "abc"},
					},
				},
			},
			expected: &objects.Role{Name: "abc"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := resolvable[objects.Role]{id: "1", data: tt.data}
			assert.Equal(t, tt.expected, r.Resolve())
		})
	}
}

func TestResolvableMessage_Resolve(t *testing.T) {
	tests := []struct {
		name string

		data     *objects.ApplicationCommandInteractionData
		expected *objects.Message
	}{
		{
			name:     "nil",
			data:     &objects.ApplicationCommandInteractionData{},
			expected: nil,
		},
		{
			name: "found",
			data: &objects.ApplicationCommandInteractionData{
				Resolved: objects.ApplicationCommandInteractionDataResolved{
					Messages: map[objects.Snowflake]objects.Message{
						1: {Content: "abc"},
					},
				},
			},
			expected: &objects.Message{Content: "abc"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := resolvable[objects.Message]{id: "1", data: tt.data}
			assert.Equal(t, tt.expected, r.Resolve())
		})
	}
}

func TestResolvableAttachment_Resolve(t *testing.T) {
	tests := []struct {
		name string

		data     *objects.ApplicationCommandInteractionData
		expected *objects.Attachment
	}{
		{
			name:     "nil",
			data:     &objects.ApplicationCommandInteractionData{},
			expected: nil,
		},
		{
			name: "found",
			data: &objects.ApplicationCommandInteractionData{
				Resolved: objects.ApplicationCommandInteractionDataResolved{
					Attachments: map[objects.Snowflake]objects.Attachment{
						1: {Filename: "abc"},
					},
				},
			},
			expected: &objects.Attachment{Filename: "abc"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := resolvable[objects.Attachment]{id: "1", data: tt.data}
			assert.Equal(t, tt.expected, r.Resolve())
		})
	}
}

func TestResolvableMentionable_Snowflake(t *testing.T) {
	testSnowflake(t, &resolvableMentionable{})
}

func TestResolvableMentionable_MarshalJSON(t *testing.T) {
	testMarshalJSON(t, &resolvableMentionable{})
}

func TestResolvableMentionable_String(t *testing.T) {
	testString(t, &resolvableMentionable{})
}

func TestResolvableMentionable_Resolve(t *testing.T) {
	tests := []struct {
		name string

		data     *objects.ApplicationCommandInteractionData
		expected any
	}{
		{
			name: "nil result",
		},
		{
			name: "channel",
			data: &objects.ApplicationCommandInteractionData{
				Resolved: objects.ApplicationCommandInteractionDataResolved{
					Users:    map[objects.Snowflake]objects.User{1: {ID: 1}},
					Members:  map[objects.Snowflake]objects.GuildMember{1: {Nick: "abc"}},
					Roles:    map[objects.Snowflake]objects.Role{1: {ID: 123}},
					Channels: map[objects.Snowflake]objects.Channel{1: {ID: 123}},
				},
			},
			expected: &objects.Channel{ID: 123},
		},
		{
			name: "role",
			data: &objects.ApplicationCommandInteractionData{
				Resolved: objects.ApplicationCommandInteractionDataResolved{
					Users:    map[objects.Snowflake]objects.User{1: {ID: 1}},
					Members:  map[objects.Snowflake]objects.GuildMember{1: {Nick: "abc"}},
					Roles:    map[objects.Snowflake]objects.Role{1: {ID: 123}},
					Channels: map[objects.Snowflake]objects.Channel{},
				},
			},
			expected: &objects.Role{ID: 123},
		},
		{
			name: "member",
			data: &objects.ApplicationCommandInteractionData{
				Resolved: objects.ApplicationCommandInteractionDataResolved{
					Users:    map[objects.Snowflake]objects.User{1: {ID: 1}},
					Members:  map[objects.Snowflake]objects.GuildMember{1: {Nick: "abc"}},
					Roles:    map[objects.Snowflake]objects.Role{},
					Channels: map[objects.Snowflake]objects.Channel{},
				},
			},
			expected: &objects.GuildMember{Nick: "abc", User: &objects.User{ID: 1}},
		},
		{
			name: "user",
			data: &objects.ApplicationCommandInteractionData{
				Resolved: objects.ApplicationCommandInteractionDataResolved{
					Users:    map[objects.Snowflake]objects.User{1: {ID: 1}},
					Members:  map[objects.Snowflake]objects.GuildMember{},
					Roles:    map[objects.Snowflake]objects.Role{},
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
			r := resolvableMentionable{
				resolvable: resolvable[any]{id: "1", data: tt.data},
			}
			assert.Equal(t, tt.expected, r.Resolve())
		})
	}
}
