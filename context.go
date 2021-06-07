package interactions

import (
	"encoding/json"
	"fmt"

	"github.com/Postcord/objects"
	"github.com/Postcord/rest"
	"github.com/jinzhu/copier"
)

// Ctx represents the base Context for all interaction events
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

// AllowedMentions sets the allowed mentions field on the response
func (c *Ctx) AllowedMentions(mentions *objects.AllowedMentions) *Ctx {
	c.Response.Data.AllowedMentions = mentions
	return c
}

// Ephemeral makes the response ephemeral, only the user running the command can see it
func (c *Ctx) Ephemeral() *Ctx {
	c.Response.Data.Flags = objects.ResponseFlagEphemeral
	return c
}

// AddEmbed adds an embed to the response
func (c *Ctx) AddEmbed(em *objects.Embed) *Ctx {
	if c.Response.Data.Embeds == nil {
		c.Response.Data.Embeds = []*objects.Embed{em}
	} else {
		c.Response.Data.Embeds = append(c.Response.Data.Embeds, em)
	}
	return c
}

// SetEmbed sets the embeds array in the response to the privded single embed
func (c *Ctx) SetEmbed(em *objects.Embed) *Ctx {
	c.Response.Data.Embeds = []*objects.Embed{em}
	return c
}

// EmbedContent sets the embeds array to a single embed containing the content as the description
func (c *Ctx) EmbedContent(content string) *Ctx {
	c.Response.Data.Embeds = []*objects.Embed{
		{
			Description: content,
		},
	}
	return c
}

// SetContent sets the content of the response
func (c *Ctx) SetContent(content string) *Ctx {
	c.Response.Data.Content = content
	return c
}

// TTS enables TTS on the response message
func (c *Ctx) TTS() *Ctx {
	c.Response.Data.TTS = true
	return c
}

// Write writes the byte array to the content of the response
func (c *Ctx) Write(data []byte) (n int, err error) {
	c.Response.Data.Content = fmt.Sprintf("%s%s", c.Response.Data.Content, string(data))
	return len(data), nil
}

// Member retrieves the user who triggered the interaction
func (c *Ctx) Member() *objects.GuildMember {
	return c.Request.Member
}

// Token returns the one time use token for the interaction
func (c *Ctx) Token() string {
	return c.Request.Token
}

// App returns the app that received the interaction
func (c *Ctx) App() *App {
	return c.app
}

// Add component adds a component to the response
func (c *Ctx) AddComponent(component *objects.Component) *Ctx {
	c.Response.Data.Components = append(c.Response.Data.Components, component)
	return c
}

// Edit sends a request to edit the original interaction response with the values in the current context.  Useful for editing a response from a Go routine.
func (c *Ctx) Edit() (*objects.Message, error) {
	return c.app.restClient.EditOriginalInteractionResponse(c.Request.ApplicationID, c.Request.Token, &rest.EditWebhookMessageParams{
		Content:         c.Response.Data.Content,
		Embeds:          c.Response.Data.Embeds,
		AllowedMentions: c.Response.Data.AllowedMentions,
		Components:      c.Response.Data.Components,
	})
}
