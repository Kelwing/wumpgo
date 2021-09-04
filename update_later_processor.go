package router

import (
	"github.com/Postcord/objects"
	"github.com/Postcord/rest"
)

type restEditInteractionResponse interface {
	EditOriginalInteractionResponse(objects.Snowflake, string, *rest.EditWebhookMessageParams) (*objects.Message, error)
}

// Process the result and update the webhook.
func processUpdateLaterResponse(restClient restEditInteractionResponse, applicationID objects.Snowflake, token string, response *objects.InteractionResponse) {
	if response.Type == objects.ResponseDeferredMessageUpdate || response.Type == objects.ResponseDeferredChannelMessageWithSource {
		// We can ignore this! The token will get passed up the chain.
		return
	}
	if response.Type == objects.ResponseUpdateMessage {
		// In this case, we should do a webhook update.
		_, _ = restClient.EditOriginalInteractionResponse(applicationID, token, &rest.EditWebhookMessageParams{
			Content:         response.Data.Content,
			Embeds:          response.Data.Embeds,
			AllowedMentions: response.Data.AllowedMentions,
			Components:      response.Data.Components,
		})
	}

	// If we get here, the action you are doing is unsupported.
}
