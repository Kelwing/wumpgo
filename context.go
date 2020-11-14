package interactions

import (
	"encoding/json"
	"fmt"
	"github.com/Postcord/objects"
)

type HandlerFunc func(ctx *CommandCtx)

type CommandCtx struct {
	Request  *objects.Interaction
	Response *objects.InteractionResponse
}

func (c *CommandCtx) UnmarshalJSON(data []byte) error {
	c.Request = new(objects.Interaction)
	if err := json.Unmarshal(data, c.Request); err != nil {
		return err
	}
	return nil
}

func (c *CommandCtx) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Response)
}

func (c *CommandCtx) Acknowledge() *CommandCtx {
	c.Response.Type = objects.ResponseAcknowledge
	c.Response.Data = nil
	return c
}

func (c *CommandCtx) AllowedMentions(mentions *objects.AllowedMentions) *CommandCtx {
	c.Response.Data.AllowedMentions = mentions
	return c
}

func (c *CommandCtx) WithSource() *CommandCtx {
	c.Response.Type = objects.ResponseChannelMessageWithSource
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
