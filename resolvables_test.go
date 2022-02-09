package router

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/Postcord/objects"
	"github.com/stretchr/testify/assert"
)

func setUnexportedField(field reflect.Value, value interface{}) {
	reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).
		Elem().
		Set(reflect.ValueOf(value))
}

func testSnowflake(t *testing.T, resolvable interface{}) {
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
	testSnowflake(t, &ResolvableUser{})
}

func TestResolvableChannel_Snowflake(t *testing.T) {
	testSnowflake(t, &ResolvableChannel{})
}

func TestResolvableRole_Snowflake(t *testing.T) {
	testSnowflake(t, &ResolvableRole{})
}

func TestResolvableMessage_Snowflake(t *testing.T) {
	testSnowflake(t, &ResolvableMessage{})
}

func TestResolvableAttachment_Snowflake(t *testing.T) {
	testSnowflake(t, &ResolvableAttachment{})
}

func testMarshalJSON(t *testing.T, resolvable interface{}) {
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
	testMarshalJSON(t, &ResolvableUser{})
}

func TestResolvableChannel_MarshalJSON(t *testing.T) {
	testMarshalJSON(t, &ResolvableChannel{})
}

func TestResolvableRole_MarshalJSON(t *testing.T) {
	testMarshalJSON(t, &ResolvableRole{})
}

func TestResolvableMessage_MarshalJSON(t *testing.T) {
	testMarshalJSON(t, &ResolvableMessage{})
}

func TestResolvableAttachment_MarshalJSON(t *testing.T) {
	testMarshalJSON(t, &ResolvableAttachment{})
}

func testString(t *testing.T, resolvable interface{}) {
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
	testString(t, &ResolvableUser{})
}

func TestResolvableChannel_String(t *testing.T) {
	testString(t, &ResolvableChannel{})
}

func TestResolvableRole_String(t *testing.T) {
	testString(t, &ResolvableRole{})
}

func TestResolvableMessage_String(t *testing.T) {
	testString(t, &ResolvableMessage{})
}

func TestResolvableAttachment_String(t *testing.T) {
	testString(t, &ResolvableAttachment{})
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
			r := ResolvableUser{id: "1", data: tt.data}
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
			r := ResolvableChannel{id: "1", data: tt.data}
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
			r := ResolvableRole{id: "1", data: tt.data}
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
			r := ResolvableMessage{id: "1", data: tt.data}
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
			r := ResolvableAttachment{id: "1", data: tt.data}
			assert.Equal(t, tt.expected, r.Resolve())
		})
	}
}
