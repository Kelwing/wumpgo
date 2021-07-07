// Code generated by generate_resolvables.go; DO NOT EDIT.

package router

//go:generate go run generate_resolvables.go

import (
	"encoding/json"
	"github.com/Postcord/objects"
	"strconv"
)

// ResolvableUser is used to define a User in a command option that is potentially resolvable.
type ResolvableUser struct {
	id   string
	data *objects.ApplicationCommandInteractionData
}

// Snowflake is used to return the ID as a snowflake.
func (r ResolvableUser) Snowflake() objects.Snowflake {
	n, _ := strconv.ParseUint(r.id, 10, 64)
	return objects.Snowflake(n)
}

// MarshalJSON implements the json.Marshaler interface.
func (r ResolvableUser) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.id)
}

// String is used to return the ID as a string.
func (r ResolvableUser) String() string {
	return r.id
}

// Resolve is used to attempt to resolve the item to its type. Returns nil if it doesn't exist.
func (r ResolvableUser) Resolve() *objects.User {
	x, _ := r.data.Resolved.Users[r.Snowflake()]
	return &x
}

// ResolvableChannel is used to define a Channel in a command option that is potentially resolvable.
type ResolvableChannel struct {
	id   string
	data *objects.ApplicationCommandInteractionData
}

// Snowflake is used to return the ID as a snowflake.
func (r ResolvableChannel) Snowflake() objects.Snowflake {
	n, _ := strconv.ParseUint(r.id, 10, 64)
	return objects.Snowflake(n)
}

// MarshalJSON implements the json.Marshaler interface.
func (r ResolvableChannel) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.id)
}

// String is used to return the ID as a string.
func (r ResolvableChannel) String() string {
	return r.id
}

// Resolve is used to attempt to resolve the item to its type. Returns nil if it doesn't exist.
func (r ResolvableChannel) Resolve() *objects.Channel {
	x, _ := r.data.Resolved.Channels[r.Snowflake()]
	return &x
}

// ResolvableRole is used to define a Role in a command option that is potentially resolvable.
type ResolvableRole struct {
	id   string
	data *objects.ApplicationCommandInteractionData
}

// Snowflake is used to return the ID as a snowflake.
func (r ResolvableRole) Snowflake() objects.Snowflake {
	n, _ := strconv.ParseUint(r.id, 10, 64)
	return objects.Snowflake(n)
}

// MarshalJSON implements the json.Marshaler interface.
func (r ResolvableRole) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.id)
}

// String is used to return the ID as a string.
func (r ResolvableRole) String() string {
	return r.id
}

// Resolve is used to attempt to resolve the item to its type. Returns nil if it doesn't exist.
func (r ResolvableRole) Resolve() *objects.Role {
	x, _ := r.data.Resolved.Roles[r.Snowflake()]
	return &x
}
