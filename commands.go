package rest

import "github.com/Postcord/objects"

// Global Commands

func (c *Client) GetCommands(app objects.Snowflake) ([]*objects.ApplicationCommand, error) {

}

func (c *Client) AddCommand(app objects.Snowflake, command *objects.ApplicationCommand) (*objects.ApplicationCommand, error) {

}

func (c *Client) UpdateCommand(app, commandID objects.Snowflake, command *objects.ApplicationCommand) (*objects.ApplicationCommand, error) {

}

func (c *Client) DeleteCommand(app, commandID objects.Snowflake) error {

}

// Guild Commands

func (c *Client) GetGuildCommand(app, guild objects.Snowflake) ([]*objects.ApplicationCommand, error) {

}

func (c *Client) AddGuildCommand(app, guild objects.Snowflake, command *objects.ApplicationCommand) (*objects.ApplicationCommand, error) {

}

func (c *Client) UpdateGuildCommand(app, guild, commandID objects.Snowflake, command *objects.ApplicationCommand) (*objects.ApplicationCommand, error) {

}

func (c *Client) DeleteGuildCommand(app, guild, commandID objects.Snowflake) error {

}
