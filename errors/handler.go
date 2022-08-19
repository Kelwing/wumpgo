package command_errors

import (
	"github.com/rs/zerolog"
	"wumpgo.dev/wumpgo/objects"
)

func ErrorHandler(logger zerolog.Logger) func(err error) *objects.InteractionResponse {
	return func(err error) *objects.InteractionResponse {
		data := &objects.InteractionApplicationCommandCallbackData{
			Flags: objects.MsgFlagEphemeral,
		}

		switch err.(type) {
		case CommandError:
			data.Content = err.Error()
		default:
			logger.Error().Err(err).Stack().Msg("unhandled error")
			data.Content = "An unknown error occurred"
		}
		return &objects.InteractionResponse{
			Type: objects.ResponseChannelMessageWithSource,
			Data: data,
		}
	}
}
