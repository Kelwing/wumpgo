package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Postcord/objects"
)

// Global Commands

func (c *Client) GetCommands(ctx context.Context, app objects.SnowflakeObject) ([]*objects.ApplicationCommand, error) {
	var commands []*objects.ApplicationCommand
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(GlobalApplicationsFmt, app.GetID())).
		ContentType(JsonContentType).
		Expect(http.StatusOK).
		Bind(&commands).
		Send(c)

	return commands, err
}

func (c *Client) CreateCommand(ctx context.Context, app objects.SnowflakeObject, command *objects.ApplicationCommand) (*objects.ApplicationCommand, error) {
	data, err := json.Marshal(command)
	if err != nil {
		return nil, err
	}

	cmd := &objects.ApplicationCommand{}

	err = NewRequest().
		Method(http.MethodPost).
		WithContext(ctx).
		Path(fmt.Sprintf(GlobalApplicationsFmt, app.GetID())).
		ContentType(JsonContentType).
		Expect(http.StatusCreated).
		Bind(cmd).
		Body(data).
		Send(c)

	return cmd, err
}

func (c *Client) GetCommand(ctx context.Context, app, commandID objects.SnowflakeObject) (*objects.ApplicationCommand, error) {
	cmd := &objects.ApplicationCommand{}
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(GlobalApplicationsUpdateFmt, app.GetID(), commandID.GetID())).
		ContentType(JsonContentType).
		Expect(http.StatusOK).
		Bind(cmd).
		Send(c)

	return cmd, err
}

func (c *Client) UpdateCommand(ctx context.Context, app, commandID objects.SnowflakeObject, command *objects.ApplicationCommand) (*objects.ApplicationCommand, error) {
	data, err := json.Marshal(command)
	if err != nil {
		return nil, err
	}

	cmd := &objects.ApplicationCommand{}

	err = NewRequest().
		Method(http.MethodPatch).
		WithContext(ctx).
		Path(fmt.Sprintf(GlobalApplicationsUpdateFmt, app.GetID(), commandID.GetID())).
		ContentType(JsonContentType).
		Expect(http.StatusOK).
		Bind(cmd).
		Body(data).
		Send(c)

	return cmd, err
}

func (c *Client) DeleteCommand(ctx context.Context, app, commandID objects.SnowflakeObject) error {
	return NewRequest().
		Method(http.MethodDelete).
		WithContext(ctx).
		Path(fmt.Sprintf(GlobalApplicationsUpdateFmt, app.GetID(), commandID.GetID())).
		ContentType(JsonContentType).
		Expect(http.StatusNoContent).
		Send(c)
}

func (c *Client) BulkOverwriteGlobalCommands(ctx context.Context, app objects.SnowflakeObject, commands []*objects.ApplicationCommand) ([]*objects.ApplicationCommand, error) {
	data, err := json.Marshal(commands)
	if err != nil {
		return nil, err
	}

	var cmds []*objects.ApplicationCommand

	err = NewRequest().
		Method(http.MethodPut).
		WithContext(ctx).
		Path(fmt.Sprintf(GlobalApplicationsFmt, app.GetID())).
		ContentType(JsonContentType).
		Expect(http.StatusOK).
		Bind(&cmds).
		Body(data).
		Send(c)

	return cmds, err
}

// Guild Commands

func (c *Client) GetGuildCommands(ctx context.Context, app, guild objects.SnowflakeObject) ([]*objects.ApplicationCommand, error) {
	var commands []*objects.ApplicationCommand
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildApplicationsFmt, app.GetID(), guild.GetID())).
		ContentType(JsonContentType).
		Expect(http.StatusOK).
		Bind(&commands).
		Send(c)

	return commands, err
}

func (c *Client) AddGuildCommand(ctx context.Context, app, guild objects.SnowflakeObject, command *objects.ApplicationCommand) (*objects.ApplicationCommand, error) {
	data, err := json.Marshal(command)
	if err != nil {
		return nil, err
	}

	cmd := &objects.ApplicationCommand{}
	err = NewRequest().
		Method(http.MethodPost).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildApplicationsFmt, app.GetID(), guild.GetID())).
		ContentType(JsonContentType).
		Expect(http.StatusOK).
		Bind(cmd).
		Body(data).
		Send(c)

	return cmd, err
}

func (c *Client) GetGuildCommand(ctx context.Context, app, guild, commandID objects.SnowflakeObject) (*objects.ApplicationCommand, error) {
	cmd := &objects.ApplicationCommand{}
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildApplicationsUpdateFmt, app.GetID(), guild.GetID(), commandID.GetID())).
		ContentType(JsonContentType).
		Expect(http.StatusOK).
		Bind(cmd).
		Send(c)

	return cmd, err
}

func (c *Client) UpdateGuildCommand(ctx context.Context, app, guild, commandID objects.SnowflakeObject, command *objects.ApplicationCommand) (*objects.ApplicationCommand, error) {
	data, err := json.Marshal(command)
	if err != nil {
		return nil, err
	}

	cmd := &objects.ApplicationCommand{}
	err = NewRequest().
		Method(http.MethodPatch).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildApplicationsUpdateFmt, app.GetID(), guild.GetID(), commandID.GetID())).
		ContentType(JsonContentType).
		Expect(http.StatusOK).
		Bind(cmd).
		Body(data).
		Send(c)

	return cmd, err
}

func (c *Client) DeleteGuildCommand(ctx context.Context, app, guild, commandID objects.SnowflakeObject) error {
	return NewRequest().
		Method(http.MethodDelete).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildApplicationsUpdateFmt, app.GetID(), guild.GetID(), commandID.GetID())).
		ContentType(JsonContentType).
		Expect(http.StatusNoContent).
		Send(c)
}

func (c *Client) BulkOverwriteGuildCommands(ctx context.Context, application, guild objects.SnowflakeObject, commands []*objects.ApplicationCommand) ([]*objects.ApplicationCommand, error) {
	data, err := json.Marshal(commands)
	if err != nil {
		return nil, err
	}

	var cmds []*objects.ApplicationCommand

	err = NewRequest().
		Method(http.MethodPut).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildApplicationsFmt, application.GetID(), guild.GetID())).
		ContentType(JsonContentType).
		Expect(http.StatusOK).
		Body(data).
		Bind(&cmds).
		Send(c)

	return cmds, err
}

func (c *Client) GetGuildApplicationCommandPermissions(ctx context.Context, app, guild objects.SnowflakeObject) ([]*objects.GuildApplicationCommandPermissions, error) {
	var commands []*objects.GuildApplicationCommandPermissions
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildApplicationCommandsPermissionsFmt, app.GetID(), guild.GetID())).
		ContentType(JsonContentType).
		Expect(http.StatusOK).
		Bind(&commands).
		Send(c)

	return commands, err
}

func (c *Client) GetApplicationCommandPermissions(ctx context.Context, app, guild, cmd objects.SnowflakeObject) (*objects.GuildApplicationCommandPermissions, error) {
	var command objects.GuildApplicationCommandPermissions
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildApplicationCommandPermissionsFmt, app.GetID(), guild.GetID(), cmd.GetID())).
		ContentType(JsonContentType).
		Expect(http.StatusOK).
		Bind(&command).
		Send(c)

	return &command, err
}

func (c *Client) EditApplicationCommandPermissions(ctx context.Context, app, guild, cmd objects.SnowflakeObject, permissions []*objects.ApplicationCommandPermissions) (*objects.GuildApplicationCommandPermissions, error) {
	data, err := json.Marshal(map[string]interface{}{
		"permissions": permissions,
	})
	if err != nil {
		return nil, err
	}

	var perms objects.GuildApplicationCommandPermissions

	err = NewRequest().
		Method(http.MethodPut).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildApplicationCommandPermissionsFmt, app.GetID(), guild.GetID(), cmd.GetID())).
		ContentType(JsonContentType).
		Expect(http.StatusNoContent, http.StatusOK).
		Body(data).
		Bind(&perms).
		Send(c)

	return &perms, err
}

func (c *Client) BatchEditApplicationCommandPermissions(ctx context.Context, app, guild objects.SnowflakeObject, permissions []*objects.GuildApplicationCommandPermissions) ([]*objects.GuildApplicationCommandPermissions, error) {
	data, err := json.Marshal(permissions)
	if err != nil {
		return nil, err
	}

	var perms []*objects.GuildApplicationCommandPermissions

	err = NewRequest().
		Method(http.MethodPut).
		Path(fmt.Sprintf(GuildApplicationCommandsPermissionsFmt, app.GetID(), guild.GetID())).
		ContentType(JsonContentType).
		Expect(http.StatusOK).
		Body(data).
		Bind(&perms).
		Send(c)

	return perms, err
}
