package router

import "github.com/kelwing/wumpgo/objects"

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
