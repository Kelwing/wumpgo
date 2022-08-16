package router

import (
	"context"

	"github.com/Postcord/objects"
	"github.com/Postcord/rest"
)

type restEditInteractionResponse interface {
	EditOriginalInteractionResponse(context.Context, objects.SnowflakeObject, string, *rest.EditWebhookMessageParams) (*objects.Message, error)
}

// Process the result and update the webhook.
func processUpdateLaterResponse(reqCtx context.Context, restClient restEditInteractionResponse, applicationID objects.Snowflake, token string, response *objects.InteractionResponse) {
	if response.Type == objects.ResponseDeferredMessageUpdate || response.Type == objects.ResponseDeferredChannelMessageWithSource {
		// We can ignore this! The token will get passed up the chain.
		return
	}
	_, _ = restClient.EditOriginalInteractionResponse(reqCtx, applicationID, token, &rest.EditWebhookMessageParams{
		Content:         response.Data.Content,
		Embeds:          response.Data.Embeds,
		AllowedMentions: response.Data.AllowedMentions,
		Components:      response.Data.Components,
	})
}
