package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"wumpgo.dev/wumpgo/objects"
)

// Global Commands

func (c *Client) GetCommands(ctx context.Context, app objects.Snowflake) ([]*objects.ApplicationCommand, error) {
	var commands []*objects.ApplicationCommand
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(GlobalApplicationsFmt, app)).
		ContentType(JsonContentType).
		Bind(&commands).
		Send(c)

	return commands, err
}

func (c *Client) CreateCommand(ctx context.Context, app objects.Snowflake, command *objects.ApplicationCommand) (*objects.ApplicationCommand, error) {
	data, err := json.Marshal(command)
	if err != nil {
		return nil, err
	}

	cmd := &objects.ApplicationCommand{}

	err = NewRequest().
		Method(http.MethodPost).
		WithContext(ctx).
		Path(fmt.Sprintf(GlobalApplicationsFmt, app)).
		ContentType(JsonContentType).
		Bind(cmd).
		Body(data).
		Send(c)

	return cmd, err
}

func (c *Client) GetCommand(ctx context.Context, app, commandID objects.Snowflake) (*objects.ApplicationCommand, error) {
	cmd := &objects.ApplicationCommand{}
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(GlobalApplicationsUpdateFmt, app, commandID)).
		ContentType(JsonContentType).
		Bind(cmd).
		Send(c)

	return cmd, err
}

func (c *Client) UpdateCommand(ctx context.Context, app, commandID objects.Snowflake, command *objects.ApplicationCommand) (*objects.ApplicationCommand, error) {
	data, err := json.Marshal(command)
	if err != nil {
		return nil, err
	}

	cmd := &objects.ApplicationCommand{}

	err = NewRequest().
		Method(http.MethodPatch).
		WithContext(ctx).
		Path(fmt.Sprintf(GlobalApplicationsUpdateFmt, app, commandID)).
		ContentType(JsonContentType).
		Bind(cmd).
		Body(data).
		Send(c)

	return cmd, err
}

func (c *Client) DeleteCommand(ctx context.Context, app, commandID objects.Snowflake) error {
	return NewRequest().
		Method(http.MethodDelete).
		WithContext(ctx).
		Path(fmt.Sprintf(GlobalApplicationsUpdateFmt, app, commandID)).
		ContentType(JsonContentType).
		Send(c)
}

func (c *Client) BulkOverwriteGlobalCommands(ctx context.Context, app objects.Snowflake, commands []*objects.ApplicationCommand) ([]*objects.ApplicationCommand, error) {
	data, err := json.Marshal(commands)
	if err != nil {
		return nil, err
	}

	var cmds []*objects.ApplicationCommand

	err = NewRequest().
		Method(http.MethodPut).
		WithContext(ctx).
		Path(fmt.Sprintf(GlobalApplicationsFmt, app)).
		ContentType(JsonContentType).
		Bind(&cmds).
		Body(data).
		Send(c)

	return cmds, err
}

// Guild Commands

func (c *Client) GetGuildCommands(ctx context.Context, app, guild objects.Snowflake) ([]*objects.ApplicationCommand, error) {
	var commands []*objects.ApplicationCommand
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildApplicationsFmt, app, guild)).
		ContentType(JsonContentType).
		Bind(&commands).
		Send(c)

	return commands, err
}

func (c *Client) AddGuildCommand(ctx context.Context, app, guild objects.Snowflake, command *objects.ApplicationCommand) (*objects.ApplicationCommand, error) {
	data, err := json.Marshal(command)
	if err != nil {
		return nil, err
	}

	cmd := &objects.ApplicationCommand{}
	err = NewRequest().
		Method(http.MethodPost).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildApplicationsFmt, app, guild)).
		ContentType(JsonContentType).
		Bind(cmd).
		Body(data).
		Send(c)

	return cmd, err
}

func (c *Client) GetGuildCommand(ctx context.Context, app, guild, commandID objects.Snowflake) (*objects.ApplicationCommand, error) {
	cmd := &objects.ApplicationCommand{}
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildApplicationsUpdateFmt, app, guild, commandID)).
		ContentType(JsonContentType).
		Bind(cmd).
		Send(c)

	return cmd, err
}

func (c *Client) UpdateGuildCommand(ctx context.Context, app, guild, commandID objects.Snowflake, command *objects.ApplicationCommand) (*objects.ApplicationCommand, error) {
	data, err := json.Marshal(command)
	if err != nil {
		return nil, err
	}

	cmd := &objects.ApplicationCommand{}
	err = NewRequest().
		Method(http.MethodPatch).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildApplicationsUpdateFmt, app, guild, commandID)).
		ContentType(JsonContentType).
		Bind(cmd).
		Body(data).
		Send(c)

	return cmd, err
}

func (c *Client) DeleteGuildCommand(ctx context.Context, app, guild, commandID objects.Snowflake) error {
	return NewRequest().
		Method(http.MethodDelete).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildApplicationsUpdateFmt, app, guild, commandID)).
		ContentType(JsonContentType).
		Send(c)
}

func (c *Client) BulkOverwriteGuildCommands(ctx context.Context, application, guild objects.Snowflake, commands []*objects.ApplicationCommand) ([]*objects.ApplicationCommand, error) {
	data, err := json.Marshal(commands)
	if err != nil {
		return nil, err
	}

	var cmds []*objects.ApplicationCommand

	err = NewRequest().
		Method(http.MethodPut).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildApplicationsFmt, application, guild)).
		ContentType(JsonContentType).
		Body(data).
		Bind(&cmds).
		Send(c)

	return cmds, err
}

func (c *Client) GetGuildApplicationCommandPermissions(ctx context.Context, app, guild objects.Snowflake) ([]*objects.GuildApplicationCommandPermissions, error) {
	var commands []*objects.GuildApplicationCommandPermissions
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildApplicationCommandsPermissionsFmt, app, guild)).
		ContentType(JsonContentType).
		Bind(&commands).
		Send(c)

	return commands, err
}

func (c *Client) GetApplicationCommandPermissions(ctx context.Context, app, guild, cmd objects.Snowflake) (*objects.GuildApplicationCommandPermissions, error) {
	var command objects.GuildApplicationCommandPermissions
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(fmt.Sprintf(GuildApplicationCommandPermissionsFmt, app, guild, cmd)).
		ContentType(JsonContentType).
		Bind(&command).
		Send(c)

	return &command, err
}

func (c *Client) EditApplicationCommandPermissions(ctx context.Context, app, guild, cmd objects.Snowflake, permissions []*objects.ApplicationCommandPermissions) (*objects.GuildApplicationCommandPermissions, error) {
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
		Path(fmt.Sprintf(GuildApplicationCommandPermissionsFmt, app, guild, cmd)).
		ContentType(JsonContentType).
		Body(data).
		Bind(&perms).
		Send(c)

	return &perms, err
}

func (c *Client) BatchEditApplicationCommandPermissions(ctx context.Context, app, guild objects.Snowflake, permissions []*objects.GuildApplicationCommandPermissions) ([]*objects.GuildApplicationCommandPermissions, error) {
	data, err := json.Marshal(permissions)
	if err != nil {
		return nil, err
	}

	var perms []*objects.GuildApplicationCommandPermissions

	err = NewRequest().
		Method(http.MethodPut).
		Path(fmt.Sprintf(GuildApplicationCommandsPermissionsFmt, app, guild)).
		ContentType(JsonContentType).
		Body(data).
		Bind(&perms).
		Send(c)

	return perms, err
}
