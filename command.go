package interactions

import (
	"encoding/json"

	"github.com/Postcord/objects"
)

type HandlerFunc func(ctx *CommandCtx)

type CommandCtx struct {
	*Ctx
	options map[string]*CommandOption
	Data    *objects.ApplicationCommandInteractionData
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
