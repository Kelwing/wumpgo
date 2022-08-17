package router

import (
	"encoding/json"
	"strconv"

	"github.com/kelwing/wumpgo/objects"
)

// Resolvable is the type which is used for all resolvable items.
type Resolvable[T any] interface {
	// Snowflake is used to return the ID as a snowflake.
	Snowflake() objects.Snowflake

	// MarshalJSON implements the json.Marshaler interface.
	MarshalJSON() ([]byte, error)

	// String is used to return the ID as a string.
	String() string

	// RawData exposes the underlying data.
	RawData() *objects.ApplicationCommandInteractionData

	// Resolve is used to attempt to resolve the item to its type. Returns nil if it doesn't exist.
	Resolve() *T
}

// ResolvableMentionable is a special type that *mostly* implements the Resolvable interface but
// does not return a pointer for resolve.
type ResolvableMentionable interface {
	// Snowflake is used to return the ID as a snowflake.
	Snowflake() objects.Snowflake

	// MarshalJSON implements the json.Marshaler interface.
	MarshalJSON() ([]byte, error)

	// String is used to return the ID as a string.
	String() string

	// RawData exposes the underlying data.
	RawData() *objects.ApplicationCommandInteractionData

	// Resolve is used to try and resolve the ID into an any type. Nil is returned if it wasn't resolved, or a *objects.<type> if it was.
	Resolve() any
}

type resolvable[T any] struct {
	id   string
	data *objects.ApplicationCommandInteractionData
}

func (r resolvable[T]) Snowflake() objects.Snowflake {
	n, _ := strconv.ParseUint(r.id, 10, 64)
	return objects.Snowflake(n)
}

func (r resolvable[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.id)
}

func (r resolvable[T]) Resolve() *T {
	if r.data == nil {
		return nil
	}
	doResolve := func(x any, ok bool) *T {
		if !ok {
			return nil
		}
		genericed := x.(T)
		return &genericed
	}
	var ptr *T
	switch (any)(ptr).(type) {
	case *objects.User:
		a, b := r.data.Resolved.Users[r.Snowflake()]
		return doResolve(a, b)
	case *objects.Channel:
		a, b := r.data.Resolved.Channels[r.Snowflake()]
		return doResolve(a, b)
	case *objects.Role:
		a, b := r.data.Resolved.Roles[r.Snowflake()]
		return doResolve(a, b)
	case *objects.Message:
		a, b := r.data.Resolved.Messages[r.Snowflake()]
		return doResolve(a, b)
	case *objects.Attachment:
		a, b := r.data.Resolved.Attachments[r.Snowflake()]
		return doResolve(a, b)
	default:
		panic("postcord internal error - unknown type")
	}
}

func (r resolvable[T]) String() string {
	return r.id
}

func (r resolvable[T]) RawData() *objects.ApplicationCommandInteractionData {
	return r.data
}

// Used to define a Mentionable in a command option that is potentially resolvable.
type resolvableMentionable struct {
	resolvable[any]
}

// Resolve is used to try and resolve the ID into an any type. Nil is returned if it wasn't resolved, or a *objects.<type> if it was.
func (r resolvableMentionable) Resolve() any {
	snowflake := r.Snowflake()
	data := r.data.Resolved

	if c, ok := data.Channels[snowflake]; ok {
		return &c
	}
	if role, ok := data.Roles[snowflake]; ok {
		return &role
	}
	if m, ok := data.Members[snowflake]; ok {
		if u, ok := data.Users[snowflake]; ok {
			m.User = &u
		}
		return &m
	}
	if u, ok := data.Users[snowflake]; ok {
		return &u
	}
	return nil
}

// ResolvableChannel is used to define a channel in a command option that is potentially resolvable.
type ResolvableChannel Resolvable[objects.Channel]

// ResolvableRole is used to define a role in a command option that is potentially resolvable.
type ResolvableRole Resolvable[objects.Role]

// ResolvableMessage is used to define a message in a command option that is potentially resolvable.
type ResolvableMessage Resolvable[objects.Message]

// ResolvableAttachment is used to define an attachment in a command option that is potentially resolvable.
type ResolvableAttachment Resolvable[objects.Attachment]

// ResolvableUser is used to define a user in a command option that is potentially resolvable.
type ResolvableUser interface {
	Resolvable[objects.User]

	// ResolveMember is used to attempt to resolve the item to a member. Returns nil if not a member.
	ResolveMember() *objects.GuildMember
}
