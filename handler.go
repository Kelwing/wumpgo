package commanderrors

import (
	"runtime/debug"

	"github.com/Postcord/objects"
)

type ErrorLogger interface {
	Errorf(format string, args ...interface{})
}

func ErrorHandler(logger ...ErrorLogger) func(err error) *objects.InteractionResponse {
	return func(err error) *objects.InteractionResponse {
		data := &objects.InteractionApplicationCommandCallbackData{
			Flags: int(objects.MsgFlagEphemeral),
		}

		switch err.(type) {
		case CommandError:
			data.Content = err.Error()
		default:
			if len(logger) > 0 {
				logger[0].Errorf("%s\n", debug.Stack())
				logger[0].Errorf("%+v\n", err)
			}
			data.Content = "An unknown error occurred"
		}
		return &objects.InteractionResponse{
			Type: objects.ResponseChannelMessageWithSource,
			Data: data,
		}
	}
}
