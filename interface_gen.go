// Code generated by generate_interface.go; DO NOT EDIT.

package rest

import (
	"image"

	"github.com/Postcord/objects"
)

// RESTClient is the interface that contains all functions in *rest.Client.
type RESTClient interface {
	AddGuildCommand(objects.Snowflake, objects.Snowflake, *objects.ApplicationCommand) (*objects.ApplicationCommand, error)
	AddGuildMember(objects.Snowflake, objects.Snowflake, *AddGuildMemberParams) (*objects.GuildMember, error)
	AddGuildMemberRole(objects.Snowflake, objects.Snowflake, objects.Snowflake, string) error
	AddPinnedMessage(objects.Snowflake, objects.Snowflake) error
	AddThreadMember(objects.Snowflake, objects.Snowflake) error
	BatchEditApplicationCommandPermissions(objects.Snowflake, objects.Snowflake, []*objects.GuildApplicationCommandPermissions) ([]*objects.GuildApplicationCommandPermissions, error)
	BeginGuildPrune(objects.Snowflake, *BeginGuildPruneParams) (int, error)
	BulkDeleteMessages(objects.Snowflake, *DeleteMessagesParams) error
	BulkOverwriteGlobalCommands(objects.Snowflake, []*objects.ApplicationCommand) ([]*objects.ApplicationCommand, error)
	BulkOverwriteGuildCommands(objects.Snowflake, objects.Snowflake, []*objects.ApplicationCommand) ([]*objects.ApplicationCommand, error)
	CreateBan(objects.Snowflake, objects.Snowflake, *CreateGuildBanParams) error
	CreateChannelInvite(objects.Snowflake, *CreateInviteParams) (*objects.Invite, error)
	CreateCommand(objects.Snowflake, *objects.ApplicationCommand) (*objects.ApplicationCommand, error)
	CreateDM(*CreateDMParams) (*objects.Channel, error)
	CreateFollowupMessage(objects.Snowflake, string, *CreateFollowupMessageParams) (*objects.Message, error)
	CreateGroupDM(*CreateGroupDMParams) (*objects.Channel, error)
	CreateGuild(*CreateGuildParams) (*objects.Guild, error)
	CreateGuildChannel(objects.Snowflake, *ChannelCreateParams) (*objects.Channel, error)
	CreateGuildFromTemplate(string, string) (*objects.Guild, error)
	CreateGuildRole(objects.Snowflake, *CreateGuildRoleParams) (*objects.Role, error)
	CreateGuildTemplate(objects.Snowflake, *CreateGuildTemplateParams) (*objects.Template, error)
	CreateInteractionResponse(objects.Snowflake, string, *objects.InteractionResponse) error
	CreateMessage(objects.Snowflake, *CreateMessageParams) (*objects.Message, error)
	CreateReaction(objects.Snowflake, objects.Snowflake, interface {}) error
	CreateWebhook(objects.Snowflake, *CreateWebhookParams) (*objects.Webhook, error)
	CrossPostMessage(objects.Snowflake, objects.Snowflake) (*objects.Message, error)
	DeleteAllReactions(objects.Snowflake, objects.Snowflake) error
	DeleteChannel(objects.Snowflake, string) (*objects.Channel, error)
	DeleteChannelPermission(objects.Snowflake, objects.Snowflake, string) error
	DeleteCommand(objects.Snowflake, objects.Snowflake) error
	DeleteEmojiReactions(objects.Snowflake, objects.Snowflake, interface {}) error
	DeleteFollowupMessage(objects.Snowflake, string, objects.Snowflake) error
	DeleteGuild(objects.Snowflake) error
	DeleteGuildCommand(objects.Snowflake, objects.Snowflake, objects.Snowflake) error
	DeleteGuildIntegration(objects.Snowflake, objects.Snowflake, string) error
	DeleteGuildRole(objects.Snowflake, objects.Snowflake, string) error
	DeleteGuildTemplate(objects.Snowflake, string, string) (*objects.Template, error)
	DeleteInvite(string, string) (*objects.Invite, error)
	DeleteMessage(objects.Snowflake, objects.Snowflake) error
	DeleteOriginalInteractionResponse(objects.Snowflake, string) error
	DeleteOwnReaction(objects.Snowflake, objects.Snowflake, interface {}) error
	DeletePinnedMessage(objects.Snowflake, objects.Snowflake) error
	DeleteUserReaction(objects.Snowflake, objects.Snowflake, objects.Snowflake, interface {}) error
	DeleteWebhook(objects.Snowflake) error
	DeleteWebhookMessage(objects.Snowflake, objects.Snowflake, string) error
	DeleteWebhookWithToken(objects.Snowflake, string) error
	EditApplicationCommandPermissions(objects.Snowflake, objects.Snowflake, objects.Snowflake, []*objects.ApplicationCommandPermissions) (*objects.GuildApplicationCommandPermissions, error)
	EditChannelPermissions(objects.Snowflake, objects.Snowflake, *EditChannelParams) error
	EditFollowupMessage(objects.Snowflake, string, objects.Snowflake, *EditWebhookMessageParams) (*objects.Message, error)
	EditMessage(objects.Snowflake, objects.Snowflake, *EditMessageParams) (*objects.Message, error)
	EditOriginalInteractionResponse(objects.Snowflake, string, *EditWebhookMessageParams) (*objects.Message, error)
	EditWebhookMessage(objects.Snowflake, objects.Snowflake, string, *EditWebhookMessageParams) (*objects.Message, error)
	ExecuteWebhook(objects.Snowflake, string, *ExecuteWebhookParams) (*objects.Message, error)
	FollowNewsChannel(objects.Snowflake) (*objects.FollowedChannel, error)
	Gateway() (*objects.Gateway, error)
	GatewayBot() (*objects.Gateway, error)
	GetApplicationCommandPermissions(objects.Snowflake, objects.Snowflake, objects.Snowflake) (*objects.GuildApplicationCommandPermissions, error)
	GetAuditLogs(objects.Snowflake, *GetAuditLogParams) (*objects.AuditLog, error)
	GetChannel(objects.Snowflake) (*objects.Channel, error)
	GetChannelInvites(objects.Snowflake) ([]*objects.Invite, error)
	GetChannelMessage(objects.Snowflake, objects.Snowflake) (*objects.Message, error)
	GetChannelMessages(objects.Snowflake, *GetChannelMessagesParams) ([]*objects.Message, error)
	GetChannelWebhooks(objects.Snowflake) ([]*objects.Webhook, error)
	GetCommand(objects.Snowflake, objects.Snowflake) (*objects.ApplicationCommand, error)
	GetCommands(objects.Snowflake) ([]*objects.ApplicationCommand, error)
	GetCurrentUser() (*objects.User, error)
	GetCurrentUserGuilds(*CurrentUserGuildsParams) ([]*objects.Guild, error)
	GetFollowupMessage(objects.Snowflake, string, objects.Snowflake) (*objects.Message, error)
	GetGuild(objects.Snowflake) (*objects.Guild, error)
	GetGuildApplicationCommandPermissions(objects.Snowflake, objects.Snowflake) ([]*objects.GuildApplicationCommandPermissions, error)
	GetGuildBan(objects.Snowflake, objects.Snowflake) (*objects.Ban, error)
	GetGuildBans(objects.Snowflake) ([]*objects.Ban, error)
	GetGuildChannels(objects.Snowflake) ([]*objects.Channel, error)
	GetGuildCommand(objects.Snowflake, objects.Snowflake, objects.Snowflake) (*objects.ApplicationCommand, error)
	GetGuildCommands(objects.Snowflake, objects.Snowflake) ([]*objects.ApplicationCommand, error)
	GetGuildIntegrations(objects.Snowflake) ([]*objects.Integration, error)
	GetGuildInvites(objects.Snowflake) ([]*objects.Invite, error)
	GetGuildMember(objects.Snowflake, objects.Snowflake) (*objects.GuildMember, error)
	GetGuildPreview(objects.Snowflake) (*objects.GuildPreview, error)
	GetGuildPruneCount(objects.Snowflake, *GetGuildPruneCountParams) (int, error)
	GetGuildRoles(objects.Snowflake) ([]*objects.Role, error)
	GetGuildTemplates(objects.Snowflake) ([]*objects.Template, error)
	GetGuildVanityURL(objects.Snowflake) (*objects.Invite, error)
	GetGuildVoiceRegions(objects.Snowflake) ([]*objects.VoiceRegion, error)
	GetGuildWebhooks(objects.Snowflake) ([]*objects.Webhook, error)
	GetGuildWelcomeScreen(objects.Snowflake) (*objects.MembershipScreening, error)
	GetGuildWidget(objects.Snowflake) (*objects.GuildWidgetJSON, error)
	GetGuildWidgetImage(objects.Snowflake, *GuildWidgetImageParams) (image.Image, error)
	GetGuildWidgetSettings(objects.Snowflake) (*objects.GuildWidget, error)
	GetInvite(string, *GetInviteParams) (*objects.Invite, error)
	GetOriginalInteractionResponse(objects.Snowflake, string) (*objects.Message, error)
	GetPinnedMessages(objects.Snowflake) ([]*objects.Message, error)
	GetReactions(objects.Snowflake, objects.Snowflake, interface {}, *GetReactionsParams) ([]*objects.User, error)
	GetTemplate(string) (*objects.Template, error)
	GetUser(objects.Snowflake) (*objects.User, error)
	GetUserConnections() ([]*objects.Connection, error)
	GetVoiceRegions() ([]*objects.VoiceRegion, error)
	GetWebhook(objects.Snowflake) (*objects.Webhook, error)
	GetWebhookWithToken(objects.Snowflake, string) (*objects.Webhook, error)
	JoinThread(objects.Snowflake) error
	LeaveGuild(objects.Snowflake) error
	LeaveThread(objects.Snowflake) error
	ListActiveThreads(objects.Snowflake) ([]*ListThreadsResponse, error)
	ListGuildMembers(objects.Snowflake, *ListGuildMembersParams) ([]*objects.GuildMember, error)
	ListJoinedPrivateArchivedThreads(objects.Snowflake, ...*ListThreadsParams) (*ListThreadsResponse, error)
	ListPrivateArchivedThreads(objects.Snowflake, ...*ListThreadsParams) (*ListThreadsResponse, error)
	ListPublicArchivedThreads(objects.Snowflake, ...*ListThreadsParams) (*ListThreadsResponse, error)
	ListThreadMembers(objects.Snowflake) ([]*objects.ThreadMember, error)
	ModifyChannel(objects.Snowflake, *ModifyChannelParams) (*objects.Channel, error)
	ModifyCurrentUser(*ModifyCurrentUserParams) (*objects.User, error)
	ModifyCurrentUserNick(objects.Snowflake, *ModifyCurrentUserNickParams) (*ModifyCurrentUserNickParams, error)
	ModifyGuild(objects.Snowflake, *ModifyGuildParams) (*objects.Guild, error)
	ModifyGuildChannelPositions(objects.Snowflake, []*ModifyChannelPositionParams, string) error
	ModifyGuildMember(objects.Snowflake, objects.Snowflake, *ModifyGuildMemberParams) (*objects.GuildMember, error)
	ModifyGuildRole(objects.Snowflake, objects.Snowflake, *ModifyGuildRoleParams) (*objects.Role, error)
	ModifyGuildRolePositions(objects.Snowflake, []*ModifyGuildRolePositionsParams) ([]*objects.Role, error)
	ModifyGuildTemplate(objects.Snowflake, string, *ModifyGuildTemplateParams) (*objects.Template, error)
	ModifyGuildWelcomeScreen(objects.Snowflake, *ModifyGuildMembershipScreeningParams) (*objects.MembershipScreening, error)
	ModifyGuildWidget(objects.Snowflake, *GuildWidgetParams) (*objects.GuildWidget, error)
	ModifyWebhook(objects.Snowflake, *ModifyWebhookParams) (*objects.Webhook, error)
	ModifyWebhookWithToken(objects.Snowflake, string, *ModifyWebhookWithTokenParams) (*objects.Webhook, error)
	RemoveGuildBan(objects.Snowflake, objects.Snowflake, string) error
	RemoveGuildMember(objects.Snowflake, objects.Snowflake, string) error
	RemoveGuildMemberRole(objects.Snowflake, objects.Snowflake, objects.Snowflake, string) error
	RemoveThreadMember(objects.Snowflake, objects.Snowflake) error
	StartThread(objects.Snowflake, *StartThreadParams) (*objects.Channel, error)
	StartThreadWithMessage(objects.Snowflake, objects.Snowflake, *StartThreadParams) (*objects.Channel, error)
	StartTyping(objects.Snowflake) error
	SyncGuildTemplate(objects.Snowflake, string) (*objects.Template, error)
	UpdateCommand(objects.Snowflake, objects.Snowflake, *objects.ApplicationCommand) (*objects.ApplicationCommand, error)
	UpdateGuildCommand(objects.Snowflake, objects.Snowflake, objects.Snowflake, *objects.ApplicationCommand) (*objects.ApplicationCommand, error)
}

var _ RESTClient = (*Client)(nil)
