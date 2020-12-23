package rest

import (
	"encoding/json"
	"fmt"
	"github.com/Postcord/objects"
	"net/http"
)

// Global Commands

func (c *Client) GetCommands(app objects.Snowflake) ([]*objects.ApplicationCommand, error) {
	res, err := c.request(&request{
		method:      http.MethodGet,
		path:        fmt.Sprintf(GlobalApplicationsFmt, app),
		contentType: JsonContentType,
	})
	if err != nil {
		return nil, err
	}

	if err = res.ExpectsStatus(http.StatusOK); err != nil {
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

	res, err := c.request(&request{
		method:      http.MethodPost,
		path:        fmt.Sprintf(GlobalApplicationsFmt, app),
		contentType: JsonContentType,
		body:        data,
	})
	if err != nil {
		return nil, err
	}

	if err = res.ExpectsStatus(http.StatusOK); err != nil {
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

	res, err := c.request(&request{
		method:      http.MethodPatch,
		path:        fmt.Sprintf(GlobalApplicationsUpdateFmt, app, commandID),
		contentType: JsonContentType,
		body:        data,
	})
	if err != nil {
		return nil, err
	}

	if err = res.ExpectsStatus(http.StatusOK); err != nil {
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
	res, err := c.request(&request{
		method:      http.MethodDelete,
		path:        fmt.Sprintf(GlobalApplicationsUpdateFmt, app, commandID),
		contentType: JsonContentType,
		body:        nil,
	})
	if err != nil {
		return err
	}

	if err = res.ExpectsStatus(http.StatusNoContent); err != nil {
		return err
	}

	return nil
}

// Guild Commands

func (c *Client) GetGuildCommand(app, guild objects.Snowflake) ([]*objects.ApplicationCommand, error) {
	res, err := c.request(&request{
		method:      http.MethodGet,
		path:        fmt.Sprintf(GuildApplicationsFmt, app, guild),
		contentType: JsonContentType,
	})
	if err != nil {
		return nil, err
	}

	if err = res.ExpectsStatus(http.StatusOK); err != nil {
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

	res, err := c.request(&request{
		method:      http.MethodPost,
		path:        fmt.Sprintf(GuildApplicationsFmt, app, guild),
		contentType: JsonContentType,
		body:        data,
	})
	if err != nil {
		return nil, err
	}

	if err = res.ExpectsStatus(http.StatusOK); err != nil {
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

	res, err := c.request(&request{
		method:      http.MethodPatch,
		path:        fmt.Sprintf(GuildApplicationsUpdateFmt, app, guild, commandID),
		contentType: JsonContentType,
		body:        data,
	})
	if err != nil {
		return nil, err
	}

	if err = res.ExpectsStatus(http.StatusOK); err != nil {
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
	res, err := c.request(&request{
		method:      http.MethodDelete,
		path:        fmt.Sprintf(GuildApplicationsUpdateFmt, app, guild, commandID),
		contentType: JsonContentType,
		body:        nil,
	})
	if err != nil {
		return err
	}

	if err = res.ExpectsStatus(http.StatusNoContent); err != nil {
		return err
	}
	return nil
}
