package router

import (
	"encoding/json"
	"github.com/Postcord/objects"
	"strconv"
)

// ResolvableMentionable is used to define a Mentionable in a command option that is potentially resolvable.
type ResolvableMentionable struct {
	id   string
	data *objects.ApplicationCommandInteractionData
}

// Snowflake is used to return the ID as a snowflake.
func (r ResolvableMentionable) Snowflake() objects.Snowflake {
	n, _ := strconv.ParseUint(r.id, 10, 64)
	return objects.Snowflake(n)
}

// MarshalJSON implements the json.Marshaler interface.
func (r ResolvableMentionable) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.id)
}

// String is used to return the ID as a string.
func (r ResolvableMentionable) String() string {
	return r.id
}

// Resolve is used to try and resolve the ID into an interface{}. Nil is returned if it wasn't resolved, or a *objects.<type> if it was.
func (r ResolvableMentionable) Resolve() interface{} {
	snowflake := r.Snowflake()
	data := r.data.Resolved
	var peek interface{}
	var ok bool

	// I'm so sorry, this is probably the fastest way to do this but it's not nice looking.
	if peek, ok = data.Channels[snowflake]; !ok {
		if peek, ok = data.Roles[snowflake]; !ok {
			if m, ok := data.Members[snowflake]; ok {
				u := data.Users[snowflake]
				m.User = &u
				peek = m
			} else if peek, ok = data.Users[snowflake]; !ok {
				return nil
			}
		}
	}

	return peek
}
