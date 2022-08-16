package router

import "github.com/Postcord/objects"

type resolvableUser struct {
	resolvable[objects.User]
}

func (r resolvableUser) ResolveMember() *objects.GuildMember {
	snowflake := r.Snowflake()
	data := r.RawData()
	if x, ok := data.Resolved.Members[snowflake]; ok {
		if u, ok := data.Resolved.Users[snowflake]; ok {
			x.User = &u
		}
		return &x
	}
	return nil
}
