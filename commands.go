package rest

import (
	"encoding/json"
	"fmt"
	"github.com/Postcord/objects"
	"net/http"
)

// Global Commands

func (c *Client) GetCommands(app objects.Snowflake) ([]*objects.ApplicationCommand, error) {
	res, err := c.request(http.MethodGet, fmt.Sprintf(GlobalApplicationsFmt, app), JsonContentType, nil)
	if err != nil {
		return nil, err
	}

	var commands []*objects.ApplicationCommand
	err = res.JSON(&commands)
	if err != nil {
		return nil, err
	}

	return commands, nil
}

func (c *Client) AddCommand(app objects.Snowflake, command *objects.ApplicationCommand) (*objects.ApplicationCommand, error) {
	data, err := json.Marshal(command)
	if err != nil {
		return nil, err
	}

	res, err := c.request(http.MethodPost, fmt.Sprintf(GlobalApplicationsFmt, app), JsonContentType, data)
	if err != nil {
		return nil, err
	}

	cmd := &objects.ApplicationCommand{}
	err = res.JSON(cmd)
	if err != nil {
		return nil, err
	}

	return cmd, nil
}

func (c *Client) UpdateCommand(app, commandID objects.Snowflake, command *objects.ApplicationCommand) (*objects.ApplicationCommand, error) {
	data, err := json.Marshal(command)
	if err != nil {
		return nil, err
	}

	res, err := c.request(http.MethodPatch, fmt.Sprintf(GlobalApplicationsUpdateFmt, app, commandID), JsonContentType, data)
	if err != nil {
		return nil, err
	}

	cmd := &objects.ApplicationCommand{}
	err = res.JSON(cmd)
	if err != nil {
		return nil, err
	}

	return cmd, nil
}

func (c *Client) DeleteCommand(app, commandID objects.Snowflake) error {
	_, err := c.request(http.MethodDelete, fmt.Sprintf(GlobalApplicationsUpdateFmt, app, commandID), JsonContentType, nil)
	if err != nil {
		return err
	}
	return nil
}

// Guild Commands

func (c *Client) GetGuildCommand(app, guild objects.Snowflake) ([]*objects.ApplicationCommand, error) {
	res, err := c.request(http.MethodGet, fmt.Sprintf(GuildApplicationsFmt, app, guild), JsonContentType, nil)
	if err != nil {
		return nil, err
	}

	var commands []*objects.ApplicationCommand
	err = res.JSON(&commands)
	if err != nil {
		return nil, err
	}

	return commands, nil
}

func (c *Client) AddGuildCommand(app, guild objects.Snowflake, command *objects.ApplicationCommand) (*objects.ApplicationCommand, error) {
	data, err := json.Marshal(command)
	if err != nil {
		return nil, err
	}

	res, err := c.request(http.MethodPost, fmt.Sprintf(GuildApplicationsFmt, app, guild), JsonContentType, data)
	if err != nil {
		return nil, err
	}

	cmd := &objects.ApplicationCommand{}
	err = res.JSON(cmd)
	if err != nil {
		return nil, err
	}

	return cmd, nil
}

func (c *Client) UpdateGuildCommand(app, guild, commandID objects.Snowflake, command *objects.ApplicationCommand) (*objects.ApplicationCommand, error) {
	data, err := json.Marshal(command)
	if err != nil {
		return nil, err
	}

	res, err := c.request(http.MethodPatch, fmt.Sprintf(GuildApplicationsUpdateFmt, app, guild, commandID), JsonContentType, data)
	if err != nil {
		return nil, err
	}

	cmd := &objects.ApplicationCommand{}
	err = res.JSON(cmd)
	if err != nil {
		return nil, err
	}

	return cmd, nil
}

func (c *Client) DeleteGuildCommand(app, guild, commandID objects.Snowflake) error {
	_, err := c.request(http.MethodDelete, fmt.Sprintf(GuildApplicationsUpdateFmt, app, guild, commandID), JsonContentType, nil)
	if err != nil {
		return err
	}
	return nil
}
