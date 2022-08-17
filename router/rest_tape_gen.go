// Code generated by generate_rest_tape.go; DO NOT EDIT.

package router

//go:generate go run generate_rest_tape.go

import (
	"context"
	"image"

	"github.com/kelwing/wumpgo/objects"
	"github.com/kelwing/wumpgo/rest"
)

type restTape struct {
	tape *tape
	rest rest.RESTClient
}

func (r restTape) AddGuildCommand(a context.Context, b objects.SnowflakeObject, c objects.SnowflakeObject, d *objects.ApplicationCommand) (*objects.ApplicationCommand, error) {
	result := r.tape.write("AddGuildCommand", false, a, b, c, d)
	e, f := r.rest.AddGuildCommand(a, b, c, d)
	result.end(e, f)
	return e, f
}

func (r restTape) AddGuildMember(a context.Context, b objects.SnowflakeObject, c objects.SnowflakeObject, d *rest.AddGuildMemberParams) (*objects.GuildMember, error) {
	result := r.tape.write("AddGuildMember", false, a, b, c, d)
	e, f := r.rest.AddGuildMember(a, b, c, d)
	result.end(e, f)
	return e, f
}

func (r restTape) AddGuildMemberRole(a context.Context, b objects.SnowflakeObject, c objects.SnowflakeObject, d objects.SnowflakeObject, e string) error {
	result := r.tape.write("AddGuildMemberRole", false, a, b, c, d, e)
	x := r.rest.AddGuildMemberRole(a, b, c, d, e)
	result.end(x)
	return x
}

func (r restTape) AddPinnedMessage(a context.Context, b objects.SnowflakeObject, c objects.SnowflakeObject) error {
	result := r.tape.write("AddPinnedMessage", false, a, b, c)
	x := r.rest.AddPinnedMessage(a, b, c)
	result.end(x)
	return x
}

func (r restTape) AddThreadMember(a context.Context, b objects.SnowflakeObject, c objects.SnowflakeObject) error {
	result := r.tape.write("AddThreadMember", false, a, b, c)
	x := r.rest.AddThreadMember(a, b, c)
	result.end(x)
	return x
}

func (r restTape) BatchEditApplicationCommandPermissions(a context.Context, b objects.SnowflakeObject, c objects.SnowflakeObject, d []*objects.GuildApplicationCommandPermissions) ([]*objects.GuildApplicationCommandPermissions, error) {
	result := r.tape.write("BatchEditApplicationCommandPermissions", false, a, b, c, d)
	e, f := r.rest.BatchEditApplicationCommandPermissions(a, b, c, d)
	result.end(e, f)
	return e, f
}

func (r restTape) BeginGuildPrune(a context.Context, b objects.SnowflakeObject, c *rest.BeginGuildPruneParams) (int, error) {
	result := r.tape.write("BeginGuildPrune", false, a, b, c)
	d, e := r.rest.BeginGuildPrune(a, b, c)
	result.end(d, e)
	return d, e
}

func (r restTape) BulkDeleteMessages(a context.Context, b objects.SnowflakeObject, c *rest.DeleteMessagesParams) error {
	result := r.tape.write("BulkDeleteMessages", false, a, b, c)
	x := r.rest.BulkDeleteMessages(a, b, c)
	result.end(x)
	return x
}

func (r restTape) BulkOverwriteGlobalCommands(a context.Context, b objects.SnowflakeObject, c []*objects.ApplicationCommand) ([]*objects.ApplicationCommand, error) {
	result := r.tape.write("BulkOverwriteGlobalCommands", false, a, b, c)
	d, e := r.rest.BulkOverwriteGlobalCommands(a, b, c)
	result.end(d, e)
	return d, e
}

func (r restTape) BulkOverwriteGuildCommands(a context.Context, b objects.SnowflakeObject, c objects.SnowflakeObject, d []*objects.ApplicationCommand) ([]*objects.ApplicationCommand, error) {
	result := r.tape.write("BulkOverwriteGuildCommands", false, a, b, c, d)
	e, f := r.rest.BulkOverwriteGuildCommands(a, b, c, d)
	result.end(e, f)
	return e, f
}

func (r restTape) CreateBan(a context.Context, b objects.SnowflakeObject, c objects.SnowflakeObject, d *rest.CreateGuildBanParams) error {
	result := r.tape.write("CreateBan", false, a, b, c, d)
	x := r.rest.CreateBan(a, b, c, d)
	result.end(x)
	return x
}

func (r restTape) CreateChannelInvite(a context.Context, b objects.SnowflakeObject, c *rest.CreateInviteParams) (*objects.Invite, error) {
	result := r.tape.write("CreateChannelInvite", false, a, b, c)
	d, e := r.rest.CreateChannelInvite(a, b, c)
	result.end(d, e)
	return d, e
}

func (r restTape) CreateCommand(a context.Context, b objects.SnowflakeObject, c *objects.ApplicationCommand) (*objects.ApplicationCommand, error) {
	result := r.tape.write("CreateCommand", false, a, b, c)
	d, e := r.rest.CreateCommand(a, b, c)
	result.end(d, e)
	return d, e
}

func (r restTape) CreateDM(a context.Context, b *rest.CreateDMParams) (*objects.Channel, error) {
	result := r.tape.write("CreateDM", false, a, b)
	c, d := r.rest.CreateDM(a, b)
	result.end(c, d)
	return c, d
}

func (r restTape) CreateFollowupMessage(a context.Context, b objects.SnowflakeObject, c string, d *rest.CreateFollowupMessageParams) (*objects.Message, error) {
	result := r.tape.write("CreateFollowupMessage", false, a, b, c, d)
	e, f := r.rest.CreateFollowupMessage(a, b, c, d)
	result.end(e, f)
	return e, f
}

func (r restTape) CreateGroupDM(a context.Context, b *rest.CreateGroupDMParams) (*objects.Channel, error) {
	result := r.tape.write("CreateGroupDM", false, a, b)
	c, d := r.rest.CreateGroupDM(a, b)
	result.end(c, d)
	return c, d
}

func (r restTape) CreateGuild(a context.Context, b *rest.CreateGuildParams) (*objects.Guild, error) {
	result := r.tape.write("CreateGuild", false, a, b)
	c, d := r.rest.CreateGuild(a, b)
	result.end(c, d)
	return c, d
}

func (r restTape) CreateGuildChannel(a context.Context, b objects.SnowflakeObject, c *rest.ChannelCreateParams) (*objects.Channel, error) {
	result := r.tape.write("CreateGuildChannel", false, a, b, c)
	d, e := r.rest.CreateGuildChannel(a, b, c)
	result.end(d, e)
	return d, e
}

func (r restTape) CreateGuildFromTemplate(a context.Context, b string, c string) (*objects.Guild, error) {
	result := r.tape.write("CreateGuildFromTemplate", false, a, b, c)
	d, e := r.rest.CreateGuildFromTemplate(a, b, c)
	result.end(d, e)
	return d, e
}

func (r restTape) CreateGuildRole(a context.Context, b objects.SnowflakeObject, c *rest.CreateGuildRoleParams) (*objects.Role, error) {
	result := r.tape.write("CreateGuildRole", false, a, b, c)
	d, e := r.rest.CreateGuildRole(a, b, c)
	result.end(d, e)
	return d, e
}

func (r restTape) CreateGuildScheduledEvent(a context.Context, b objects.SnowflakeObject, c *rest.CreateGuildScheduledEventParams) (*objects.GuildScheduledEvent, error) {
	result := r.tape.write("CreateGuildScheduledEvent", false, a, b, c)
	d, e := r.rest.CreateGuildScheduledEvent(a, b, c)
	result.end(d, e)
	return d, e
}

func (r restTape) CreateGuildSticker(a context.Context, b objects.SnowflakeObject, c *rest.CreateGuildStickerParams) (*objects.Sticker, error) {
	result := r.tape.write("CreateGuildSticker", false, a, b, c)
	d, e := r.rest.CreateGuildSticker(a, b, c)
	result.end(d, e)
	return d, e
}

func (r restTape) CreateGuildTemplate(a context.Context, b objects.SnowflakeObject, c *rest.CreateGuildTemplateParams) (*objects.Template, error) {
	result := r.tape.write("CreateGuildTemplate", false, a, b, c)
	d, e := r.rest.CreateGuildTemplate(a, b, c)
	result.end(d, e)
	return d, e
}

func (r restTape) CreateInteractionResponse(a context.Context, b objects.SnowflakeObject, c string, d *objects.InteractionResponse) error {
	result := r.tape.write("CreateInteractionResponse", false, a, b, c, d)
	x := r.rest.CreateInteractionResponse(a, b, c, d)
	result.end(x)
	return x
}

func (r restTape) CreateMessage(a context.Context, b objects.SnowflakeObject, c *rest.CreateMessageParams) (*objects.Message, error) {
	result := r.tape.write("CreateMessage", false, a, b, c)
	d, e := r.rest.CreateMessage(a, b, c)
	result.end(d, e)
	return d, e
}

func (r restTape) CreateReaction(a context.Context, b objects.SnowflakeObject, c objects.SnowflakeObject, d interface {}) error {
	result := r.tape.write("CreateReaction", false, a, b, c, d)
	x := r.rest.CreateReaction(a, b, c, d)
	result.end(x)
	return x
}

func (r restTape) CreateWebhook(a context.Context, b objects.SnowflakeObject, c *rest.CreateWebhookParams) (*objects.Webhook, error) {
	result := r.tape.write("CreateWebhook", false, a, b, c)
	d, e := r.rest.CreateWebhook(a, b, c)
	result.end(d, e)
	return d, e
}

func (r restTape) CrossPostMessage(a context.Context, b objects.SnowflakeObject, c objects.SnowflakeObject) (*objects.Message, error) {
	result := r.tape.write("CrossPostMessage", false, a, b, c)
	d, e := r.rest.CrossPostMessage(a, b, c)
	result.end(d, e)
	return d, e
}

func (r restTape) DeleteAllReactions(a context.Context, b objects.SnowflakeObject, c objects.SnowflakeObject) error {
	result := r.tape.write("DeleteAllReactions", false, a, b, c)
	x := r.rest.DeleteAllReactions(a, b, c)
	result.end(x)
	return x
}

func (r restTape) DeleteChannel(a context.Context, b objects.SnowflakeObject, c string) (*objects.Channel, error) {
	result := r.tape.write("DeleteChannel", false, a, b, c)
	d, e := r.rest.DeleteChannel(a, b, c)
	result.end(d, e)
	return d, e
}

func (r restTape) DeleteChannelPermission(a context.Context, b objects.SnowflakeObject, c objects.SnowflakeObject, d string) error {
	result := r.tape.write("DeleteChannelPermission", false, a, b, c, d)
	x := r.rest.DeleteChannelPermission(a, b, c, d)
	result.end(x)
	return x
}

func (r restTape) DeleteCommand(a context.Context, b objects.SnowflakeObject, c objects.SnowflakeObject) error {
	result := r.tape.write("DeleteCommand", false, a, b, c)
	x := r.rest.DeleteCommand(a, b, c)
	result.end(x)
	return x
}

func (r restTape) DeleteEmojiReactions(a context.Context, b objects.SnowflakeObject, c objects.SnowflakeObject, d interface {}) error {
	result := r.tape.write("DeleteEmojiReactions", false, a, b, c, d)
	x := r.rest.DeleteEmojiReactions(a, b, c, d)
	result.end(x)
	return x
}

func (r restTape) DeleteFollowupMessage(a context.Context, b objects.SnowflakeObject, c string, d objects.SnowflakeObject) error {
	result := r.tape.write("DeleteFollowupMessage", false, a, b, c, d)
	x := r.rest.DeleteFollowupMessage(a, b, c, d)
	result.end(x)
	return x
}

func (r restTape) DeleteGuild(a context.Context, b objects.SnowflakeObject) error {
	result := r.tape.write("DeleteGuild", false, a, b)
	x := r.rest.DeleteGuild(a, b)
	result.end(x)
	return x
}

func (r restTape) DeleteGuildCommand(a context.Context, b objects.SnowflakeObject, c objects.SnowflakeObject, d objects.SnowflakeObject) error {
	result := r.tape.write("DeleteGuildCommand", false, a, b, c, d)
	x := r.rest.DeleteGuildCommand(a, b, c, d)
	result.end(x)
	return x
}

func (r restTape) DeleteGuildIntegration(a context.Context, b objects.SnowflakeObject, c objects.SnowflakeObject, d string) error {
	result := r.tape.write("DeleteGuildIntegration", false, a, b, c, d)
	x := r.rest.DeleteGuildIntegration(a, b, c, d)
	result.end(x)
	return x
}

func (r restTape) DeleteGuildRole(a context.Context, b objects.SnowflakeObject, c objects.SnowflakeObject, d string) error {
	result := r.tape.write("DeleteGuildRole", false, a, b, c, d)
	x := r.rest.DeleteGuildRole(a, b, c, d)
	result.end(x)
	return x
}

func (r restTape) DeleteGuildScheduledEvent(a context.Context, b objects.SnowflakeObject, c objects.SnowflakeObject) error {
	result := r.tape.write("DeleteGuildScheduledEvent", false, a, b, c)
	x := r.rest.DeleteGuildScheduledEvent(a, b, c)
	result.end(x)
	return x
}

func (r restTape) DeleteGuildSticker(a context.Context, b objects.SnowflakeObject, c objects.SnowflakeObject, d ...string) error {
	result := r.tape.write("DeleteGuildSticker", true, a, b, c, d)
	x := r.rest.DeleteGuildSticker(a, b, c, d...)
	result.end(x)
	return x
}

func (r restTape) DeleteGuildTemplate(a context.Context, b objects.SnowflakeObject, c string, d string) (*objects.Template, error) {
	result := r.tape.write("DeleteGuildTemplate", false, a, b, c, d)
	e, f := r.rest.DeleteGuildTemplate(a, b, c, d)
	result.end(e, f)
	return e, f
}

func (r restTape) DeleteInvite(a context.Context, b string, c string) (*objects.Invite, error) {
	result := r.tape.write("DeleteInvite", false, a, b, c)
	d, e := r.rest.DeleteInvite(a, b, c)
	result.end(d, e)
	return d, e
}

func (r restTape) DeleteMessage(a context.Context, b objects.SnowflakeObject, c objects.SnowflakeObject) error {
	result := r.tape.write("DeleteMessage", false, a, b, c)
	x := r.rest.DeleteMessage(a, b, c)
	result.end(x)
	return x
}

func (r restTape) DeleteOriginalInteractionResponse(a context.Context, b objects.SnowflakeObject, c string) error {
	result := r.tape.write("DeleteOriginalInteractionResponse", false, a, b, c)
	x := r.rest.DeleteOriginalInteractionResponse(a, b, c)
	result.end(x)
	return x
}

func (r restTape) DeleteOwnReaction(a context.Context, b objects.SnowflakeObject, c objects.SnowflakeObject, d interface {}) error {
	result := r.tape.write("DeleteOwnReaction", false, a, b, c, d)
	x := r.rest.DeleteOwnReaction(a, b, c, d)
	result.end(x)
	return x
}

func (r restTape) DeletePinnedMessage(a context.Context, b objects.SnowflakeObject, c objects.SnowflakeObject) error {
	result := r.tape.write("DeletePinnedMessage", false, a, b, c)
	x := r.rest.DeletePinnedMessage(a, b, c)
	result.end(x)
	return x
}

func (r restTape) DeleteUserReaction(a context.Context, b objects.SnowflakeObject, c objects.SnowflakeObject, d objects.SnowflakeObject, e interface {}) error {
	result := r.tape.write("DeleteUserReaction", false, a, b, c, d, e)
	x := r.rest.DeleteUserReaction(a, b, c, d, e)
	result.end(x)
	return x
}

func (r restTape) DeleteWebhook(a context.Context, b objects.SnowflakeObject) error {
	result := r.tape.write("DeleteWebhook", false, a, b)
	x := r.rest.DeleteWebhook(a, b)
	result.end(x)
	return x
}

func (r restTape) DeleteWebhookMessage(a context.Context, b objects.SnowflakeObject, c objects.SnowflakeObject, d string) error {
	result := r.tape.write("DeleteWebhookMessage", false, a, b, c, d)
	x := r.rest.DeleteWebhookMessage(a, b, c, d)
	result.end(x)
	return x
}

func (r restTape) DeleteWebhookWithToken(a context.Context, b objects.SnowflakeObject, c string) error {
	result := r.tape.write("DeleteWebhookWithToken", false, a, b, c)
	x := r.rest.DeleteWebhookWithToken(a, b, c)
	result.end(x)
	return x
}

func (r restTape) EditApplicationCommandPermissions(a context.Context, b objects.SnowflakeObject, c objects.SnowflakeObject, d objects.SnowflakeObject, e []*objects.ApplicationCommandPermissions) (*objects.GuildApplicationCommandPermissions, error) {
	result := r.tape.write("EditApplicationCommandPermissions", false, a, b, c, d, e)
	f, g := r.rest.EditApplicationCommandPermissions(a, b, c, d, e)
	result.end(f, g)
	return f, g
}

func (r restTape) EditChannelPermissions(a context.Context, b objects.SnowflakeObject, c objects.SnowflakeObject, d *rest.EditChannelParams) error {
	result := r.tape.write("EditChannelPermissions", false, a, b, c, d)
	x := r.rest.EditChannelPermissions(a, b, c, d)
	result.end(x)
	return x
}

func (r restTape) EditFollowupMessage(a context.Context, b objects.SnowflakeObject, c string, d objects.SnowflakeObject, e *rest.EditWebhookMessageParams) (*objects.Message, error) {
	result := r.tape.write("EditFollowupMessage", false, a, b, c, d, e)
	f, g := r.rest.EditFollowupMessage(a, b, c, d, e)
	result.end(f, g)
	return f, g
}

func (r restTape) EditMessage(a context.Context, b objects.SnowflakeObject, c objects.SnowflakeObject, d *rest.EditMessageParams) (*objects.Message, error) {
	result := r.tape.write("EditMessage", false, a, b, c, d)
	e, f := r.rest.EditMessage(a, b, c, d)
	result.end(e, f)
	return e, f
}

func (r restTape) EditOriginalInteractionResponse(a context.Context, b objects.SnowflakeObject, c string, d *rest.EditWebhookMessageParams) (*objects.Message, error) {
	result := r.tape.write("EditOriginalInteractionResponse", false, a, b, c, d)
	e, f := r.rest.EditOriginalInteractionResponse(a, b, c, d)
	result.end(e, f)
	return e, f
}

func (r restTape) EditWebhookMessage(a context.Context, b objects.SnowflakeObject, c objects.SnowflakeObject, d string, e *rest.EditWebhookMessageParams) (*objects.Message, error) {
	result := r.tape.write("EditWebhookMessage", false, a, b, c, d, e)
	f, g := r.rest.EditWebhookMessage(a, b, c, d, e)
	result.end(f, g)
	return f, g
}

func (r restTape) ExecuteWebhook(a context.Context, b objects.SnowflakeObject, c string, d *rest.ExecuteWebhookParams) (*objects.Message, error) {
	result := r.tape.write("ExecuteWebhook", false, a, b, c, d)
	e, f := r.rest.ExecuteWebhook(a, b, c, d)
	result.end(e, f)
	return e, f
}

func (r restTape) FollowNewsChannel(a context.Context, b objects.SnowflakeObject) (*objects.FollowedChannel, error) {
	result := r.tape.write("FollowNewsChannel", false, a, b)
	c, d := r.rest.FollowNewsChannel(a, b)
	result.end(c, d)
	return c, d
}

func (r restTape) Gateway(a context.Context) (*objects.Gateway, error) {
	result := r.tape.write("Gateway", false, a)
	b, c := r.rest.Gateway(a)
	result.end(b, c)
	return b, c
}

func (r restTape) GatewayBot(a context.Context) (*objects.Gateway, error) {
	result := r.tape.write("GatewayBot", false, a)
	b, c := r.rest.GatewayBot(a)
	result.end(b, c)
	return b, c
}

func (r restTape) GetApplicationCommandPermissions(a context.Context, b objects.SnowflakeObject, c objects.SnowflakeObject, d objects.SnowflakeObject) (*objects.GuildApplicationCommandPermissions, error) {
	result := r.tape.write("GetApplicationCommandPermissions", false, a, b, c, d)
	e, f := r.rest.GetApplicationCommandPermissions(a, b, c, d)
	result.end(e, f)
	return e, f
}

func (r restTape) GetAuditLogs(a context.Context, b objects.SnowflakeObject, c *rest.GetAuditLogParams) (*objects.AuditLog, error) {
	result := r.tape.write("GetAuditLogs", false, a, b, c)
	d, e := r.rest.GetAuditLogs(a, b, c)
	result.end(d, e)
	return d, e
}

func (r restTape) GetChannel(a context.Context, b objects.SnowflakeObject) (*objects.Channel, error) {
	result := r.tape.write("GetChannel", false, a, b)
	c, d := r.rest.GetChannel(a, b)
	result.end(c, d)
	return c, d
}

func (r restTape) GetChannelInvites(a context.Context, b objects.SnowflakeObject) ([]*objects.Invite, error) {
	result := r.tape.write("GetChannelInvites", false, a, b)
	c, d := r.rest.GetChannelInvites(a, b)
	result.end(c, d)
	return c, d
}

func (r restTape) GetChannelMessage(a context.Context, b objects.SnowflakeObject, c objects.SnowflakeObject) (*objects.Message, error) {
	result := r.tape.write("GetChannelMessage", false, a, b, c)
	d, e := r.rest.GetChannelMessage(a, b, c)
	result.end(d, e)
	return d, e
}

func (r restTape) GetChannelMessages(a context.Context, b objects.SnowflakeObject, c *rest.GetChannelMessagesParams) ([]*objects.Message, error) {
	result := r.tape.write("GetChannelMessages", false, a, b, c)
	d, e := r.rest.GetChannelMessages(a, b, c)
	result.end(d, e)
	return d, e
}

func (r restTape) GetChannelWebhooks(a context.Context, b objects.SnowflakeObject) ([]*objects.Webhook, error) {
	result := r.tape.write("GetChannelWebhooks", false, a, b)
	c, d := r.rest.GetChannelWebhooks(a, b)
	result.end(c, d)
	return c, d
}

func (r restTape) GetCommand(a context.Context, b objects.SnowflakeObject, c objects.SnowflakeObject) (*objects.ApplicationCommand, error) {
	result := r.tape.write("GetCommand", false, a, b, c)
	d, e := r.rest.GetCommand(a, b, c)
	result.end(d, e)
	return d, e
}

func (r restTape) GetCommands(a context.Context, b objects.SnowflakeObject) ([]*objects.ApplicationCommand, error) {
	result := r.tape.write("GetCommands", false, a, b)
	c, d := r.rest.GetCommands(a, b)
	result.end(c, d)
	return c, d
}

func (r restTape) GetCurrentUser(a context.Context) (*objects.User, error) {
	result := r.tape.write("GetCurrentUser", false, a)
	b, c := r.rest.GetCurrentUser(a)
	result.end(b, c)
	return b, c
}

func (r restTape) GetCurrentUserGuildMember(a context.Context, b objects.SnowflakeObject) (*objects.GuildMember, error) {
	result := r.tape.write("GetCurrentUserGuildMember", false, a, b)
	c, d := r.rest.GetCurrentUserGuildMember(a, b)
	result.end(c, d)
	return c, d
}

func (r restTape) GetCurrentUserGuilds(a context.Context, b *rest.CurrentUserGuildsParams) ([]*objects.Guild, error) {
	result := r.tape.write("GetCurrentUserGuilds", false, a, b)
	c, d := r.rest.GetCurrentUserGuilds(a, b)
	result.end(c, d)
	return c, d
}

func (r restTape) GetFollowupMessage(a context.Context, b objects.SnowflakeObject, c string, d objects.SnowflakeObject) (*objects.Message, error) {
	result := r.tape.write("GetFollowupMessage", false, a, b, c, d)
	e, f := r.rest.GetFollowupMessage(a, b, c, d)
	result.end(e, f)
	return e, f
}

func (r restTape) GetGuild(a context.Context, b objects.SnowflakeObject) (*objects.Guild, error) {
	result := r.tape.write("GetGuild", false, a, b)
	c, d := r.rest.GetGuild(a, b)
	result.end(c, d)
	return c, d
}

func (r restTape) GetGuildApplicationCommandPermissions(a context.Context, b objects.SnowflakeObject, c objects.SnowflakeObject) ([]*objects.GuildApplicationCommandPermissions, error) {
	result := r.tape.write("GetGuildApplicationCommandPermissions", false, a, b, c)
	d, e := r.rest.GetGuildApplicationCommandPermissions(a, b, c)
	result.end(d, e)
	return d, e
}

func (r restTape) GetGuildBan(a context.Context, b objects.SnowflakeObject, c objects.SnowflakeObject) (*objects.Ban, error) {
	result := r.tape.write("GetGuildBan", false, a, b, c)
	d, e := r.rest.GetGuildBan(a, b, c)
	result.end(d, e)
	return d, e
}

func (r restTape) GetGuildBans(a context.Context, b objects.SnowflakeObject) ([]*objects.Ban, error) {
	result := r.tape.write("GetGuildBans", false, a, b)
	c, d := r.rest.GetGuildBans(a, b)
	result.end(c, d)
	return c, d
}

func (r restTape) GetGuildChannels(a context.Context, b objects.SnowflakeObject) ([]*objects.Channel, error) {
	result := r.tape.write("GetGuildChannels", false, a, b)
	c, d := r.rest.GetGuildChannels(a, b)
	result.end(c, d)
	return c, d
}

func (r restTape) GetGuildCommand(a context.Context, b objects.SnowflakeObject, c objects.SnowflakeObject, d objects.SnowflakeObject) (*objects.ApplicationCommand, error) {
	result := r.tape.write("GetGuildCommand", false, a, b, c, d)
	e, f := r.rest.GetGuildCommand(a, b, c, d)
	result.end(e, f)
	return e, f
}

func (r restTape) GetGuildCommands(a context.Context, b objects.SnowflakeObject, c objects.SnowflakeObject) ([]*objects.ApplicationCommand, error) {
	result := r.tape.write("GetGuildCommands", false, a, b, c)
	d, e := r.rest.GetGuildCommands(a, b, c)
	result.end(d, e)
	return d, e
}

func (r restTape) GetGuildIntegrations(a context.Context, b objects.SnowflakeObject) ([]*objects.Integration, error) {
	result := r.tape.write("GetGuildIntegrations", false, a, b)
	c, d := r.rest.GetGuildIntegrations(a, b)
	result.end(c, d)
	return c, d
}

func (r restTape) GetGuildInvites(a context.Context, b objects.SnowflakeObject) ([]*objects.Invite, error) {
	result := r.tape.write("GetGuildInvites", false, a, b)
	c, d := r.rest.GetGuildInvites(a, b)
	result.end(c, d)
	return c, d
}

func (r restTape) GetGuildMember(a context.Context, b objects.SnowflakeObject, c objects.SnowflakeObject) (*objects.GuildMember, error) {
	result := r.tape.write("GetGuildMember", false, a, b, c)
	d, e := r.rest.GetGuildMember(a, b, c)
	result.end(d, e)
	return d, e
}

func (r restTape) GetGuildPreview(a context.Context, b objects.SnowflakeObject) (*objects.GuildPreview, error) {
	result := r.tape.write("GetGuildPreview", false, a, b)
	c, d := r.rest.GetGuildPreview(a, b)
	result.end(c, d)
	return c, d
}

func (r restTape) GetGuildPruneCount(a context.Context, b objects.SnowflakeObject, c *rest.GetGuildPruneCountParams) (int, error) {
	result := r.tape.write("GetGuildPruneCount", false, a, b, c)
	d, e := r.rest.GetGuildPruneCount(a, b, c)
	result.end(d, e)
	return d, e
}

func (r restTape) GetGuildRoles(a context.Context, b objects.SnowflakeObject) ([]*objects.Role, error) {
	result := r.tape.write("GetGuildRoles", false, a, b)
	c, d := r.rest.GetGuildRoles(a, b)
	result.end(c, d)
	return c, d
}

func (r restTape) GetGuildScheduledEvent(a context.Context, b objects.SnowflakeObject, c objects.SnowflakeObject, d ...*rest.GetGuildScheduledEventParams) (*objects.GuildScheduledEvent, error) {
	result := r.tape.write("GetGuildScheduledEvent", true, a, b, c, d)
	e, f := r.rest.GetGuildScheduledEvent(a, b, c, d...)
	result.end(e, f)
	return e, f
}

func (r restTape) GetGuildScheduledEventUsers(a context.Context, b objects.SnowflakeObject, c objects.SnowflakeObject, d ...*rest.GetGuildScheduledEventUsersParams) ([]*objects.GuildScheduledEventUser, error) {
	result := r.tape.write("GetGuildScheduledEventUsers", true, a, b, c, d)
	e, f := r.rest.GetGuildScheduledEventUsers(a, b, c, d...)
	result.end(e, f)
	return e, f
}

func (r restTape) GetGuildSticker(a context.Context, b objects.SnowflakeObject, c objects.SnowflakeObject) (*objects.Sticker, error) {
	result := r.tape.write("GetGuildSticker", false, a, b, c)
	d, e := r.rest.GetGuildSticker(a, b, c)
	result.end(d, e)
	return d, e
}

func (r restTape) GetGuildTemplates(a context.Context, b objects.SnowflakeObject) ([]*objects.Template, error) {
	result := r.tape.write("GetGuildTemplates", false, a, b)
	c, d := r.rest.GetGuildTemplates(a, b)
	result.end(c, d)
	return c, d
}

func (r restTape) GetGuildVanityURL(a context.Context, b objects.SnowflakeObject) (*objects.Invite, error) {
	result := r.tape.write("GetGuildVanityURL", false, a, b)
	c, d := r.rest.GetGuildVanityURL(a, b)
	result.end(c, d)
	return c, d
}

func (r restTape) GetGuildVoiceRegions(a context.Context, b objects.SnowflakeObject) ([]*objects.VoiceRegion, error) {
	result := r.tape.write("GetGuildVoiceRegions", false, a, b)
	c, d := r.rest.GetGuildVoiceRegions(a, b)
	result.end(c, d)
	return c, d
}

func (r restTape) GetGuildWebhooks(a context.Context, b objects.SnowflakeObject) ([]*objects.Webhook, error) {
	result := r.tape.write("GetGuildWebhooks", false, a, b)
	c, d := r.rest.GetGuildWebhooks(a, b)
	result.end(c, d)
	return c, d
}

func (r restTape) GetGuildWelcomeScreen(a context.Context, b objects.SnowflakeObject) (*objects.MembershipScreening, error) {
	result := r.tape.write("GetGuildWelcomeScreen", false, a, b)
	c, d := r.rest.GetGuildWelcomeScreen(a, b)
	result.end(c, d)
	return c, d
}

func (r restTape) GetGuildWidget(a context.Context, b objects.SnowflakeObject) (*objects.GuildWidgetJSON, error) {
	result := r.tape.write("GetGuildWidget", false, a, b)
	c, d := r.rest.GetGuildWidget(a, b)
	result.end(c, d)
	return c, d
}

func (r restTape) GetGuildWidgetImage(a context.Context, b objects.SnowflakeObject, c *rest.GuildWidgetImageParams) (image.Image, error) {
	result := r.tape.write("GetGuildWidgetImage", false, a, b, c)
	d, e := r.rest.GetGuildWidgetImage(a, b, c)
	result.end(d, e)
	return d, e
}

func (r restTape) GetGuildWidgetSettings(a context.Context, b objects.SnowflakeObject) (*objects.GuildWidget, error) {
	result := r.tape.write("GetGuildWidgetSettings", false, a, b)
	c, d := r.rest.GetGuildWidgetSettings(a, b)
	result.end(c, d)
	return c, d
}

func (r restTape) GetInvite(a context.Context, b string, c *rest.GetInviteParams) (*objects.Invite, error) {
	result := r.tape.write("GetInvite", false, a, b, c)
	d, e := r.rest.GetInvite(a, b, c)
	result.end(d, e)
	return d, e
}

func (r restTape) GetOriginalInteractionResponse(a context.Context, b objects.SnowflakeObject, c string) (*objects.Message, error) {
	result := r.tape.write("GetOriginalInteractionResponse", false, a, b, c)
	d, e := r.rest.GetOriginalInteractionResponse(a, b, c)
	result.end(d, e)
	return d, e
}

func (r restTape) GetPinnedMessages(a context.Context, b objects.SnowflakeObject) ([]*objects.Message, error) {
	result := r.tape.write("GetPinnedMessages", false, a, b)
	c, d := r.rest.GetPinnedMessages(a, b)
	result.end(c, d)
	return c, d
}

func (r restTape) GetReactions(a context.Context, b objects.SnowflakeObject, c objects.SnowflakeObject, d interface {}, e *rest.GetReactionsParams) ([]*objects.User, error) {
	result := r.tape.write("GetReactions", false, a, b, c, d, e)
	f, g := r.rest.GetReactions(a, b, c, d, e)
	result.end(f, g)
	return f, g
}

func (r restTape) GetSticker(a context.Context, b objects.SnowflakeObject) (*objects.Sticker, error) {
	result := r.tape.write("GetSticker", false, a, b)
	c, d := r.rest.GetSticker(a, b)
	result.end(c, d)
	return c, d
}

func (r restTape) GetTemplate(a context.Context, b string) (*objects.Template, error) {
	result := r.tape.write("GetTemplate", false, a, b)
	c, d := r.rest.GetTemplate(a, b)
	result.end(c, d)
	return c, d
}

func (r restTape) GetUser(a context.Context, b objects.SnowflakeObject) (*objects.User, error) {
	result := r.tape.write("GetUser", false, a, b)
	c, d := r.rest.GetUser(a, b)
	result.end(c, d)
	return c, d
}

func (r restTape) GetUserConnections(a context.Context) ([]*objects.Connection, error) {
	result := r.tape.write("GetUserConnections", false, a)
	b, c := r.rest.GetUserConnections(a)
	result.end(b, c)
	return b, c
}

func (r restTape) GetVoiceRegions(a context.Context) ([]*objects.VoiceRegion, error) {
	result := r.tape.write("GetVoiceRegions", false, a)
	b, c := r.rest.GetVoiceRegions(a)
	result.end(b, c)
	return b, c
}

func (r restTape) GetWebhook(a context.Context, b objects.SnowflakeObject) (*objects.Webhook, error) {
	result := r.tape.write("GetWebhook", false, a, b)
	c, d := r.rest.GetWebhook(a, b)
	result.end(c, d)
	return c, d
}

func (r restTape) GetWebhookWithToken(a context.Context, b objects.SnowflakeObject, c string) (*objects.Webhook, error) {
	result := r.tape.write("GetWebhookWithToken", false, a, b, c)
	d, e := r.rest.GetWebhookWithToken(a, b, c)
	result.end(d, e)
	return d, e
}

func (r restTape) JoinThread(a context.Context, b objects.SnowflakeObject) error {
	result := r.tape.write("JoinThread", false, a, b)
	x := r.rest.JoinThread(a, b)
	result.end(x)
	return x
}

func (r restTape) LeaveGuild(a context.Context, b objects.SnowflakeObject) error {
	result := r.tape.write("LeaveGuild", false, a, b)
	x := r.rest.LeaveGuild(a, b)
	result.end(x)
	return x
}

func (r restTape) LeaveThread(a context.Context, b objects.SnowflakeObject) error {
	result := r.tape.write("LeaveThread", false, a, b)
	x := r.rest.LeaveThread(a, b)
	result.end(x)
	return x
}

func (r restTape) ListActiveThreads(a context.Context, b objects.SnowflakeObject) ([]*rest.ListThreadsResponse, error) {
	result := r.tape.write("ListActiveThreads", false, a, b)
	c, d := r.rest.ListActiveThreads(a, b)
	result.end(c, d)
	return c, d
}

func (r restTape) ListGuildMembers(a context.Context, b objects.SnowflakeObject, c *rest.ListGuildMembersParams) ([]*objects.GuildMember, error) {
	result := r.tape.write("ListGuildMembers", false, a, b, c)
	d, e := r.rest.ListGuildMembers(a, b, c)
	result.end(d, e)
	return d, e
}

func (r restTape) ListGuildStickers(a context.Context, b objects.SnowflakeObject) ([]*objects.Sticker, error) {
	result := r.tape.write("ListGuildStickers", false, a, b)
	c, d := r.rest.ListGuildStickers(a, b)
	result.end(c, d)
	return c, d
}

func (r restTape) ListJoinedPrivateArchivedThreads(a context.Context, b objects.SnowflakeObject, c ...*rest.ListThreadsParams) (*rest.ListThreadsResponse, error) {
	result := r.tape.write("ListJoinedPrivateArchivedThreads", true, a, b, c)
	d, e := r.rest.ListJoinedPrivateArchivedThreads(a, b, c...)
	result.end(d, e)
	return d, e
}

func (r restTape) ListNitroStickerPacks(a context.Context) ([]*objects.StickerPack, error) {
	result := r.tape.write("ListNitroStickerPacks", false, a)
	b, c := r.rest.ListNitroStickerPacks(a)
	result.end(b, c)
	return b, c
}

func (r restTape) ListPrivateArchivedThreads(a context.Context, b objects.SnowflakeObject, c ...*rest.ListThreadsParams) (*rest.ListThreadsResponse, error) {
	result := r.tape.write("ListPrivateArchivedThreads", true, a, b, c)
	d, e := r.rest.ListPrivateArchivedThreads(a, b, c...)
	result.end(d, e)
	return d, e
}

func (r restTape) ListPublicArchivedThreads(a context.Context, b objects.SnowflakeObject, c ...*rest.ListThreadsParams) (*rest.ListThreadsResponse, error) {
	result := r.tape.write("ListPublicArchivedThreads", true, a, b, c)
	d, e := r.rest.ListPublicArchivedThreads(a, b, c...)
	result.end(d, e)
	return d, e
}

func (r restTape) ListThreadMembers(a context.Context, b objects.SnowflakeObject) ([]*objects.ThreadMember, error) {
	result := r.tape.write("ListThreadMembers", false, a, b)
	c, d := r.rest.ListThreadMembers(a, b)
	result.end(c, d)
	return c, d
}

func (r restTape) ModifyChannel(a context.Context, b objects.SnowflakeObject, c *rest.ModifyChannelParams) (*objects.Channel, error) {
	result := r.tape.write("ModifyChannel", false, a, b, c)
	d, e := r.rest.ModifyChannel(a, b, c)
	result.end(d, e)
	return d, e
}

func (r restTape) ModifyCurrentUser(a context.Context, b *rest.ModifyCurrentUserParams) (*objects.User, error) {
	result := r.tape.write("ModifyCurrentUser", false, a, b)
	c, d := r.rest.ModifyCurrentUser(a, b)
	result.end(c, d)
	return c, d
}

func (r restTape) ModifyCurrentUserNick(a context.Context, b objects.SnowflakeObject, c *rest.ModifyCurrentUserNickParams) (*rest.ModifyCurrentUserNickParams, error) {
	result := r.tape.write("ModifyCurrentUserNick", false, a, b, c)
	d, e := r.rest.ModifyCurrentUserNick(a, b, c)
	result.end(d, e)
	return d, e
}

func (r restTape) ModifyGuild(a context.Context, b objects.SnowflakeObject, c *rest.ModifyGuildParams) (*objects.Guild, error) {
	result := r.tape.write("ModifyGuild", false, a, b, c)
	d, e := r.rest.ModifyGuild(a, b, c)
	result.end(d, e)
	return d, e
}

func (r restTape) ModifyGuildChannelPositions(a context.Context, b objects.SnowflakeObject, c []*rest.ModifyChannelPositionParams, d string) error {
	result := r.tape.write("ModifyGuildChannelPositions", false, a, b, c, d)
	x := r.rest.ModifyGuildChannelPositions(a, b, c, d)
	result.end(x)
	return x
}

func (r restTape) ModifyGuildMember(a context.Context, b objects.SnowflakeObject, c objects.SnowflakeObject, d *rest.ModifyGuildMemberParams) (*objects.GuildMember, error) {
	result := r.tape.write("ModifyGuildMember", false, a, b, c, d)
	e, f := r.rest.ModifyGuildMember(a, b, c, d)
	result.end(e, f)
	return e, f
}

func (r restTape) ModifyGuildRole(a context.Context, b objects.SnowflakeObject, c objects.SnowflakeObject, d *rest.ModifyGuildRoleParams) (*objects.Role, error) {
	result := r.tape.write("ModifyGuildRole", false, a, b, c, d)
	e, f := r.rest.ModifyGuildRole(a, b, c, d)
	result.end(e, f)
	return e, f
}

func (r restTape) ModifyGuildRolePositions(a context.Context, b objects.SnowflakeObject, c []*rest.ModifyGuildRolePositionsParams) ([]*objects.Role, error) {
	result := r.tape.write("ModifyGuildRolePositions", false, a, b, c)
	d, e := r.rest.ModifyGuildRolePositions(a, b, c)
	result.end(d, e)
	return d, e
}

func (r restTape) ModifyGuildScheduledEvent(a context.Context, b objects.SnowflakeObject, c objects.SnowflakeObject, d *rest.ModifyGuildScheduledEventParams) (*objects.GuildScheduledEvent, error) {
	result := r.tape.write("ModifyGuildScheduledEvent", false, a, b, c, d)
	e, f := r.rest.ModifyGuildScheduledEvent(a, b, c, d)
	result.end(e, f)
	return e, f
}

func (r restTape) ModifyGuildSticker(a context.Context, b objects.SnowflakeObject, c objects.SnowflakeObject, d *rest.BaseStickerParams) (*objects.Sticker, error) {
	result := r.tape.write("ModifyGuildSticker", false, a, b, c, d)
	e, f := r.rest.ModifyGuildSticker(a, b, c, d)
	result.end(e, f)
	return e, f
}

func (r restTape) ModifyGuildTemplate(a context.Context, b objects.SnowflakeObject, c string, d *rest.ModifyGuildTemplateParams) (*objects.Template, error) {
	result := r.tape.write("ModifyGuildTemplate", false, a, b, c, d)
	e, f := r.rest.ModifyGuildTemplate(a, b, c, d)
	result.end(e, f)
	return e, f
}

func (r restTape) ModifyGuildWelcomeScreen(a context.Context, b objects.SnowflakeObject, c *rest.ModifyGuildMembershipScreeningParams) (*objects.MembershipScreening, error) {
	result := r.tape.write("ModifyGuildWelcomeScreen", false, a, b, c)
	d, e := r.rest.ModifyGuildWelcomeScreen(a, b, c)
	result.end(d, e)
	return d, e
}

func (r restTape) ModifyGuildWidget(a context.Context, b objects.SnowflakeObject, c *rest.GuildWidgetParams) (*objects.GuildWidget, error) {
	result := r.tape.write("ModifyGuildWidget", false, a, b, c)
	d, e := r.rest.ModifyGuildWidget(a, b, c)
	result.end(d, e)
	return d, e
}

func (r restTape) ModifyWebhook(a context.Context, b objects.SnowflakeObject, c *rest.ModifyWebhookParams) (*objects.Webhook, error) {
	result := r.tape.write("ModifyWebhook", false, a, b, c)
	d, e := r.rest.ModifyWebhook(a, b, c)
	result.end(d, e)
	return d, e
}

func (r restTape) ModifyWebhookWithToken(a context.Context, b objects.SnowflakeObject, c string, d *rest.ModifyWebhookWithTokenParams) (*objects.Webhook, error) {
	result := r.tape.write("ModifyWebhookWithToken", false, a, b, c, d)
	e, f := r.rest.ModifyWebhookWithToken(a, b, c, d)
	result.end(e, f)
	return e, f
}

func (r restTape) RemoveGuildBan(a context.Context, b objects.SnowflakeObject, c objects.SnowflakeObject, d string) error {
	result := r.tape.write("RemoveGuildBan", false, a, b, c, d)
	x := r.rest.RemoveGuildBan(a, b, c, d)
	result.end(x)
	return x
}

func (r restTape) RemoveGuildMember(a context.Context, b objects.SnowflakeObject, c objects.SnowflakeObject, d string) error {
	result := r.tape.write("RemoveGuildMember", false, a, b, c, d)
	x := r.rest.RemoveGuildMember(a, b, c, d)
	result.end(x)
	return x
}

func (r restTape) RemoveGuildMemberRole(a context.Context, b objects.SnowflakeObject, c objects.SnowflakeObject, d objects.SnowflakeObject, e string) error {
	result := r.tape.write("RemoveGuildMemberRole", false, a, b, c, d, e)
	x := r.rest.RemoveGuildMemberRole(a, b, c, d, e)
	result.end(x)
	return x
}

func (r restTape) RemoveThreadMember(a context.Context, b objects.SnowflakeObject, c objects.SnowflakeObject) error {
	result := r.tape.write("RemoveThreadMember", false, a, b, c)
	x := r.rest.RemoveThreadMember(a, b, c)
	result.end(x)
	return x
}

func (r restTape) StartThread(a context.Context, b objects.SnowflakeObject, c *rest.StartThreadParams) (*objects.Channel, error) {
	result := r.tape.write("StartThread", false, a, b, c)
	d, e := r.rest.StartThread(a, b, c)
	result.end(d, e)
	return d, e
}

func (r restTape) StartThreadWithMessage(a context.Context, b objects.SnowflakeObject, c objects.SnowflakeObject, d *rest.StartThreadParams) (*objects.Channel, error) {
	result := r.tape.write("StartThreadWithMessage", false, a, b, c, d)
	e, f := r.rest.StartThreadWithMessage(a, b, c, d)
	result.end(e, f)
	return e, f
}

func (r restTape) StartTyping(a context.Context, b objects.SnowflakeObject) error {
	result := r.tape.write("StartTyping", false, a, b)
	x := r.rest.StartTyping(a, b)
	result.end(x)
	return x
}

func (r restTape) SyncGuildTemplate(a context.Context, b objects.SnowflakeObject, c string) (*objects.Template, error) {
	result := r.tape.write("SyncGuildTemplate", false, a, b, c)
	d, e := r.rest.SyncGuildTemplate(a, b, c)
	result.end(d, e)
	return d, e
}

func (r restTape) UpdateCommand(a context.Context, b objects.SnowflakeObject, c objects.SnowflakeObject, d *objects.ApplicationCommand) (*objects.ApplicationCommand, error) {
	result := r.tape.write("UpdateCommand", false, a, b, c, d)
	e, f := r.rest.UpdateCommand(a, b, c, d)
	result.end(e, f)
	return e, f
}

func (r restTape) UpdateGuildCommand(a context.Context, b objects.SnowflakeObject, c objects.SnowflakeObject, d objects.SnowflakeObject, e *objects.ApplicationCommand) (*objects.ApplicationCommand, error) {
	result := r.tape.write("UpdateGuildCommand", false, a, b, c, d, e)
	f, g := r.rest.UpdateGuildCommand(a, b, c, d, e)
	result.end(f, g)
	return f, g
}
