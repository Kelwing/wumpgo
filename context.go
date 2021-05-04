package interactions

import (
	"encoding/json"
	"fmt"

	"github.com/Postcord/objects"
	"github.com/jinzhu/copier"
)

type HandlerFunc func(ctx *CommandCtx, data *objects.ApplicationCommandInteractionData)
type ButtonHandlerFunc func(ctx *CommandCtx, data *objects.ApplicationComponentInteractionData)

type CommandCtx struct {
	Request  *objects.Interaction
	Response *objects.InteractionResponse
	options  map[string]*CommandOption
	app      *App
}

func (c *CommandCtx) Clone() (*CommandCtx, error) {
	newCtx := &CommandCtx{}
	err := copier.Copy(newCtx, c)
	return newCtx, err
}

func (c *CommandCtx) UnmarshalJSON(data []byte) error {
	c.Request = new(objects.Interaction)
	if err := json.Unmarshal(data, c.Request); err != nil {
		return err
	}
	c.Response = &objects.InteractionResponse{
		Type: objects.ResponseChannelMessageWithSource,
		Data: &objects.InteractionApplicationCommandCallbackData{
			TTS:             false,
			Content:         "",
			Embeds:          nil,
			AllowedMentions: nil,
			Flags:           0,
		},
	}
	return nil
}

func (c *CommandCtx) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Response)
}

func (c *CommandCtx) DeferredMessageUpdate() {
	c.Response.Type = objects.ResponseDeferredMessageUpdate
}

func (c *CommandCtx) UpdateMessage() {
	c.Response.Type = objects.ResponseUpdateMessage
}

func (c *CommandCtx) AllowedMentions(mentions *objects.AllowedMentions) *CommandCtx {
	c.Response.Data.AllowedMentions = mentions
	return c
}

func (c *CommandCtx) Ephemeral() *CommandCtx {
	c.Response.Data.Flags = objects.ResponseFlagEphemeral
	return c
}

func (c *CommandCtx) AddEmbed(em *objects.Embed) *CommandCtx {
	if c.Response.Data.Embeds == nil {
		c.Response.Data.Embeds = []*objects.Embed{em}
	} else {
		c.Response.Data.Embeds = append(c.Response.Data.Embeds, em)
	}
	return c
}

func (c *CommandCtx) SetEmbed(em *objects.Embed) *CommandCtx {
	c.Response.Data.Embeds = []*objects.Embed{em}
	return c
}

func (c *CommandCtx) EmbedContent(content string) *CommandCtx {
	c.Response.Data.Embeds = []*objects.Embed{
		{
			Description: content,
		},
	}
	return c
}

func (c *CommandCtx) SetContent(content string) *CommandCtx {
	c.Response.Data.Content = content
	return c
}

func (c *CommandCtx) TTS() *CommandCtx {
	c.Response.Data.TTS = true
	return c
}

func (c *CommandCtx) Write(data []byte) (n int, err error) {
	c.Response.Data.Content = fmt.Sprintf("%s%s", c.Response.Data.Content, string(data))
	return len(data), nil
}

// Request helper functions
func (c *CommandCtx) Member() *objects.GuildMember {
	return c.Request.Member
}

func (c *CommandCtx) CommandName() string {
	var data objects.ApplicationCommandInteractionData
	err := json.Unmarshal(c.Request.Data, &data)
	if err != nil {
		return ""
	}
	return data.Name
}

func (c *CommandCtx) Options() []objects.ApplicationCommandInteractionDataOption {
	var data objects.ApplicationCommandInteractionData
	err := json.Unmarshal(c.Request.Data, &data)
	if err != nil {
		return nil
	}
	return data.Options
}

func (c *CommandCtx) Token() string {
	return c.Request.Token
}

func (c *CommandCtx) Get(name string) *CommandOption {
	option, ok := c.options[name]
	if !ok {
		return &CommandOption{Value: nil}
	}

	return option
}

func (c *CommandCtx) App() *App {
	return c.app
}

func (c *CommandCtx) AddComponent(component *objects.Component) *CommandCtx {
	c.Response.Data.Components = append(c.Response.Data.Components, component)
	return c
}
