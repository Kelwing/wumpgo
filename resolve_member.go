package router

import "github.com/Postcord/objects"

// ResolveMember is used to attempt to resolve the item to a member. Returns nil if not a member.
func (r ResolvableUser) ResolveMember() *objects.GuildMember {
	snowflake := r.Snowflake()
	if x, ok := r.data.Resolved.Members[snowflake]; ok {
		u, _ := r.data.Resolved.Users[snowflake]
		x.User = &u
		return &x
	}
	return nil
}
