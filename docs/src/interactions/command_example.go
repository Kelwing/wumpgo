//go:build exclude

package example

import (
	"wumpgo.dev/wumpgo/objects"
	"wumpgo.dev/wumpgo/rest"
	"wumpgo.dev/wumpgo/router"
)

//go:generate wumpgoctl gen

// @Name Ban
// @Description Ban a user
// @DM false
// @Permissions BanMembers
type Ban struct {
	User   objects.User `discord:"user"`
	Reason string       `discord:"reason,optional"`
}

func (b Ban) Handle(r router.CommandResponder, ctx *router.CommandContext) {
	err := ctx.Client().CreateBan(ctx.Context(), ctx.GuildID(), b.User.ID, &rest.CreateGuildBanParams{
		DeleteMessageDays: 0,
		Reason:            b.Reason,
	})
	if err != nil {
		r.Ephemeral().Contentf("Failed to ban %s", b.User.Mention())
		return
	}

	r.Ephemeral().Contentf("Banned %s for %s", b.User.Mention(), b.Reason)
}
