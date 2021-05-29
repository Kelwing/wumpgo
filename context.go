package interactions

import (
	"encoding/json"
	"fmt"

	"github.com/Postcord/objects"
	"github.com/Postcord/rest"
	"github.com/jinzhu/copier"
)

type Ctx struct {
	Request  *objects.Interaction
	Response *objects.InteractionResponse
	app      *App
}

func (c *Ctx) UnmarshalJSON(data []byte) error {
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

func (c *Ctx) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Response)
}

func (c *Ctx) Clone() (*Ctx, error) {
	newCtx := &Ctx{}
	err := copier.Copy(newCtx, c)
	return newCtx, err
}

func (c *Ctx) AllowedMentions(mentions *objects.AllowedMentions) *Ctx {
	c.Response.Data.AllowedMentions = mentions
	return c
}

func (c *Ctx) Ephemeral() *Ctx {
	c.Response.Data.Flags = objects.ResponseFlagEphemeral
	return c
}

func (c *Ctx) AddEmbed(em *objects.Embed) *Ctx {
	if c.Response.Data.Embeds == nil {
		c.Response.Data.Embeds = []*objects.Embed{em}
	} else {
		c.Response.Data.Embeds = append(c.Response.Data.Embeds, em)
	}
	return c
}

func (c *Ctx) SetEmbed(em *objects.Embed) *Ctx {
	c.Response.Data.Embeds = []*objects.Embed{em}
	return c
}

func (c *Ctx) EmbedContent(content string) *Ctx {
	c.Response.Data.Embeds = []*objects.Embed{
		{
			Description: content,
		},
	}
	return c
}

func (c *Ctx) SetContent(content string) *Ctx {
	c.Response.Data.Content = content
	return c
}

func (c *Ctx) TTS() *Ctx {
	c.Response.Data.TTS = true
	return c
}

func (c *Ctx) Write(data []byte) (n int, err error) {
	c.Response.Data.Content = fmt.Sprintf("%s%s", c.Response.Data.Content, string(data))
	return len(data), nil
}

// Request helper functions
func (c *Ctx) Member() *objects.GuildMember {
	return c.Request.Member
}

func (c *Ctx) Token() string {
	return c.Request.Token
}

func (c *Ctx) App() *App {
	return c.app
}

func (c *Ctx) AddComponent(component *objects.Component) *Ctx {
	c.Response.Data.Components = append(c.Response.Data.Components, component)
	return c
}

func (c *Ctx) Edit() error {
	_, err := c.app.restClient.EditOriginalInteractionResponse(c.Request.ApplicationID, c.Request.Token, &rest.EditWebhookMessageParams{
		Content:         c.Response.Data.Content,
		Embeds:          c.Response.Data.Embeds,
		AllowedMentions: c.Response.Data.AllowedMentions,
	})

	return err
}
