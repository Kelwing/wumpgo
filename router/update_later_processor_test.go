package router

import (
	"context"

	"github.com/kelwing/wumpgo/objects"
	"github.com/kelwing/wumpgo/rest"
	"github.com/stretchr/testify/assert"

	"testing"
)

type mockRestEditInteractionParams struct {
	applicationID objects.Snowflake
	token         string
	params        *rest.EditWebhookMessageParams
}

type mockRestEditInteractionResponse struct {
	t *testing.T

	params *mockRestEditInteractionParams
}

func newMockRestEditInteractionResponse(t *testing.T) *mockRestEditInteractionResponse {
	return &mockRestEditInteractionResponse{t: t}
}

func (m *mockRestEditInteractionResponse) EditOriginalInteractionResponse(reqCtx context.Context, applicationID objects.SnowflakeObject, token string, params *rest.EditWebhookMessageParams) (*objects.Message, error) {
	m.t.Helper()
	if m.params != nil {
		m.t.Fatal("edit interaction already called")
		return nil, nil
	}
	m.params = &mockRestEditInteractionParams{
		applicationID: applicationID.GetID(),
		token:         token,
		params:        params,
	}
	return nil, nil
}

func Test_processUpdateLaterResponse(t *testing.T) {
	tests := []struct {
		name string

		applicationID objects.Snowflake
		token         string
		response      *objects.InteractionResponse

		restParams *mockRestEditInteractionParams
	}{
		{
			name:          "deferred message update",
			applicationID: 1,
			token:         "a",
			response:      &objects.InteractionResponse{Type: objects.ResponseDeferredMessageUpdate},
		},
		{
			name:          "deferred channel message with source",
			applicationID: 1,
			token:         "a",
			response:      &objects.InteractionResponse{Type: objects.ResponseDeferredChannelMessageWithSource},
		},
		{
			name:          "update message",
			applicationID: 1,
			token:         "a",
			response: &objects.InteractionResponse{
				Type: objects.ResponseUpdateMessage,
				Data: &objects.InteractionApplicationCommandCallbackData{
					Content: "Hello World",
					Embeds: []*objects.Embed{
						{
							Title: "a embed",
						},
					},
					AllowedMentions: &objects.AllowedMentions{
						Parse:       []string{"1"},
						Roles:       []objects.Snowflake{2},
						Users:       []objects.Snowflake{3},
						RepliedUser: true,
					},
					Components: []*objects.Component{
						{
							Label: "Testing testing 123",
						},
					},
				},
			},
			restParams: &mockRestEditInteractionParams{
				applicationID: 1,
				token:         "a",
				params: &rest.EditWebhookMessageParams{
					Content: "Hello World",
					Embeds: []*objects.Embed{
						{
							Title: "a embed",
						},
					},
					AllowedMentions: &objects.AllowedMentions{
						Parse:       []string{"1"},
						Roles:       []objects.Snowflake{2},
						Users:       []objects.Snowflake{3},
						RepliedUser: true,
					},
					Components: []*objects.Component{
						{
							Label: "Testing testing 123",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			restClient := newMockRestEditInteractionResponse(t)
			processUpdateLaterResponse(context.Background(), restClient, tt.applicationID, tt.token, tt.response)
			assert.Equal(t, tt.restParams, restClient.params)
		})
	}
}
