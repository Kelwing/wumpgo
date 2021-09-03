package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Postcord/objects"
)

// Global Commands

func (c *Client) GetCommands(app objects.Snowflake) ([]*objects.ApplicationCommand, error) {
	var commands []*objects.ApplicationCommand
	err := NewRequest().
		Method(http.MethodGet).
		Path(fmt.Sprintf(GlobalApplicationsFmt, app)).
		ContentType(JsonContentType).
		Expect(http.StatusOK).
		Bind(&commands).
		Send(c)

	return commands, err
}

func (c *Client) CreateCommand(app objects.Snowflake, command *objects.ApplicationCommand) (*objects.ApplicationCommand, error) {
	data, err := json.Marshal(command)
	if err != nil {
		return nil, err
	}

	cmd := &objects.ApplicationCommand{}

	err = NewRequest().
		Method(http.MethodPost).
		Path(fmt.Sprintf(GlobalApplicationsFmt, app)).
		ContentType(JsonContentType).
		Expect(http.StatusCreated).
		Bind(cmd).
		Body(data).
		Send(c)

	return cmd, err
}

func (c *Client) GetCommand(app, commandID objects.Snowflake) (*objects.ApplicationCommand, error) {
	cmd := &objects.ApplicationCommand{}
	err := NewRequest().
		Method(http.MethodGet).
		Path(fmt.Sprintf(GlobalApplicationsUpdateFmt, app, commandID)).
		ContentType(JsonContentType).
		Expect(http.StatusOK).
		Bind(cmd).
		Send(c)

	return cmd, err
}

func (c *Client) UpdateCommand(app, commandID objects.Snowflake, command *objects.ApplicationCommand) (*objects.ApplicationCommand, error) {
	data, err := json.Marshal(command)
	if err != nil {
		return nil, err
	}

	cmd := &objects.ApplicationCommand{}

	err = NewRequest().
		Method(http.MethodPatch).
		Path(fmt.Sprintf(GlobalApplicationsUpdateFmt, app, commandID)).
		ContentType(JsonContentType).
		Expect(http.StatusOK).
		Bind(cmd).
		Body(data).
		Send(c)

	return cmd, err
}

func (c *Client) DeleteCommand(app, commandID objects.Snowflake) error {
	return NewRequest().
		Method(http.MethodDelete).
		Path(fmt.Sprintf(GlobalApplicationsUpdateFmt, app, commandID)).
		ContentType(JsonContentType).
		Expect(http.StatusNoContent).
		Send(c)

	return nil
}

func (c *Client) BulkOverwriteGlobalCommands(app objects.Snowflake, commands []*objects.ApplicationCommand) ([]*objects.ApplicationCommand, error) {
	data, err := json.Marshal(commands)
	if err != nil {
		return nil, err
	}

	var cmds []*objects.ApplicationCommand

	err = NewRequest().
		Method(http.MethodPut).
		Path(fmt.Sprintf(GlobalApplicationsFmt, app)).
		ContentType(JsonContentType).
		Expect(http.StatusOK).
		Bind(&cmds).
		Body(data).
		Send(c)

	return cmds, err
}

// Guild Commands

func (c *Client) GetGuildCommands(app, guild objects.Snowflake) ([]*objects.ApplicationCommand, error) {
	var commands []*objects.ApplicationCommand
	err := NewRequest().
		Method(http.MethodGet).
		Path(fmt.Sprintf(GuildApplicationsFmt, app, guild)).
		ContentType(JsonContentType).
		Expect(http.StatusOK).
		Bind(&commands).
		Send(c)

	return commands, err
}

func (c *Client) AddGuildCommand(app, guild objects.Snowflake, command *objects.ApplicationCommand) (*objects.ApplicationCommand, error) {
	data, err := json.Marshal(command)
	if err != nil {
		return nil, err
	}

	cmd := &objects.ApplicationCommand{}
	err = NewRequest().
		Method(http.MethodPost).
		Path(fmt.Sprintf(GuildApplicationsFmt, app, guild)).
		ContentType(JsonContentType).
		Expect(http.StatusOK).
		Bind(cmd).
		Body(data).
		Send(c)

	return cmd, err
}

func (c *Client) GetGuildCommand(app, guild, commandID objects.Snowflake) (*objects.ApplicationCommand, error) {
	cmd := &objects.ApplicationCommand{}
	err := NewRequest().
		Method(http.MethodGet).
		Path(fmt.Sprintf(GuildApplicationsUpdateFmt, app, guild, commandID)).
		ContentType(JsonContentType).
		Expect(http.StatusOK).
		Bind(cmd).
		Send(c)

	return cmd, err
}

func (c *Client) UpdateGuildCommand(app, guild, commandID objects.Snowflake, command *objects.ApplicationCommand) (*objects.ApplicationCommand, error) {
	data, err := json.Marshal(command)
	if err != nil {
		return nil, err
	}

	cmd := &objects.ApplicationCommand{}
	err = NewRequest().
		Method(http.MethodPatch).
		Path(fmt.Sprintf(GuildApplicationsUpdateFmt, app, guild, commandID)).
		ContentType(JsonContentType).
		Expect(http.StatusOK).
		Bind(cmd).
		Body(data).
		Send(c)

	return cmd, err
}

func (c *Client) DeleteGuildCommand(app, guild, commandID objects.Snowflake) error {
	return NewRequest().
		Method(http.MethodDelete).
		Path(fmt.Sprintf(GuildApplicationsUpdateFmt, app, guild, commandID)).
		ContentType(JsonContentType).
		Expect(http.StatusNoContent).
		Send(c)
}

func (c *Client) BulkOverwriteGuildCommands(application, guild objects.Snowflake, commands []*objects.ApplicationCommand) ([]*objects.ApplicationCommand, error) {
	data, err := json.Marshal(commands)
	if err != nil {
		return nil, err
	}

	var cmds []*objects.ApplicationCommand

	err = NewRequest().
		Method(http.MethodPut).
		Path(fmt.Sprintf(GuildApplicationsFmt, application, guild)).
		ContentType(JsonContentType).
		Expect(http.StatusOK).
		Body(data).
		Bind(&cmds).
		Send(c)

	return cmds, err
}

func (c *Client) GetGuildApplicationCommandPermissions(app, guild objects.Snowflake) ([]*objects.GuildApplicationCommandPermissions, error) {
	var commands []*objects.GuildApplicationCommandPermissions
	err := NewRequest().
		Method(http.MethodGet).
		Path(fmt.Sprintf(GuildApplicationCommandsPermissionsFmt, app, guild)).
		ContentType(JsonContentType).
		Expect(http.StatusOK).
		Bind(&commands).
		Send(c)

	return commands, err
}

func (c *Client) GetApplicationCommandPermissions(app, guild, cmd objects.Snowflake) (*objects.GuildApplicationCommandPermissions, error) {
	var command objects.GuildApplicationCommandPermissions
	err := NewRequest().
		Method(http.MethodGet).
		Path(fmt.Sprintf(GuildApplicationCommandPermissionsFmt, app, guild, cmd)).
		ContentType(JsonContentType).
		Expect(http.StatusOK).
		Bind(&command).
		Send(c)

	return &command, err
}

func (c *Client) EditApplicationCommandPermissions(app, guild, cmd objects.Snowflake, permissions []*objects.ApplicationCommandPermissions) (*objects.GuildApplicationCommandPermissions, error) {
	data, err := json.Marshal(map[string]interface{}{
		"permissions": permissions,
	})
	if err != nil {
		return nil, err
	}

	var perms objects.GuildApplicationCommandPermissions

	err = NewRequest().
		Method(http.MethodPut).
		Path(fmt.Sprintf(GuildApplicationCommandPermissionsFmt, app, guild, cmd)).
		ContentType(JsonContentType).
		Expect(http.StatusNoContent).
		Body(data).
		Bind(&perms).
		Send(c)

	return &perms, nil
}

func (c *Client) BatchEditApplicationCommandPermissions(app, guild objects.Snowflake, permissions []*objects.GuildApplicationCommandPermissions) ([]*objects.GuildApplicationCommandPermissions, error) {
	data, err := json.Marshal(permissions)
	if err != nil {
		return nil, err
	}

	var perms []*objects.GuildApplicationCommandPermissions

	err = NewRequest().
		Method(http.MethodPut).
		Path(fmt.Sprintf(GuildApplicationCommandsPermissionsFmt, app, guild)).
		ContentType(JsonContentType).
		Expect(http.StatusOK).
		Body(data).
		Bind(&perms).
		Send(c)

	return perms, err
}
