package router

import (
	"testing"

	"github.com/Postcord/objects"
	"github.com/stretchr/testify/assert"
)

func Test_responseBuilder_ResponseData(t *testing.T) {
	tests := []struct{
		name string

		data *objects.InteractionApplicationCommandCallbackData
	}{
		{
			name: "nil data",
		},
		{
			name: "not nil data",
			data: &objects.InteractionApplicationCommandCallbackData{TTS: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := responseBuilder{dataPtr: tt.data}
			respData := b.ResponseData()
			if tt.data == nil {
				assert.Equal(t, &objects.InteractionApplicationCommandCallbackData{}, respData)
			} else {
				assert.Equal(t, tt.data, respData)
			}
		})
	}
}

func Test_responseBuilder_editEmbed(t *testing.T) {
	tests := []struct{
		name string

		initEmbeds []*objects.Embed
		embed    *objects.Embed
		append   bool
		expected []*objects.Embed
	}{
		{
			name: "all nil",
			expected: ([]*objects.Embed)(nil),
		},
		{
			name:       "non append",
			initEmbeds: []*objects.Embed{{}},
			embed:      &objects.Embed{Description: "a"},
			append:     false,
			expected:   []*objects.Embed{{Description: "a"}},
		},
		{
			name:       "append",
			initEmbeds: []*objects.Embed{{}},
			embed:      &objects.Embed{Description: "a"},
			append:     true,
			expected:   []*objects.Embed{{}, {Description: "a"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := responseBuilder{}
			if tt.initEmbeds != nil {
				b.ResponseData().Embeds = tt.initEmbeds
			}
			b.editEmbed(tt.embed, tt.append)
			assert.Equal(t, tt.expected, b.ResponseData().Embeds)
		})
	}
}

func Test_responseBuilder_editComponent(t *testing.T) {
	tests := []struct{
		name string

		initComponents []*objects.Component
		component    *objects.Component
		append   bool
		expected []*objects.Component
	}{
		{
			name: "all nil",
			expected: ([]*objects.Component)(nil),
		},
		{
			name:       "non append",
			initComponents: []*objects.Component{{}},
			component:  &objects.Component{Placeholder: "test"},
			append:     false,
			expected:   []*objects.Component{{Placeholder: "test"}},
		},
		{
			name:       "append",
			initComponents: []*objects.Component{{}},
			component:     &objects.Component{Placeholder: "a"},
			append:     true,
			expected:   []*objects.Component{{}, {Placeholder: "a"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := responseBuilder{}
			if tt.initComponents != nil {
				b.ResponseData().Components = tt.initComponents
			}
			b.editComponent(tt.component, tt.append)
			assert.Equal(t, tt.expected, b.ResponseData().Components)
		})
	}
}
