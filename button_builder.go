package interactions

import "github.com/Postcord/objects"

// ButtonBuilder makes it easy to generate new buttons
type ButtonBuilder struct {
	button *objects.Component
}

// NewButtonBuilder creates a new button builder
func NewButtonBuilder() *ButtonBuilder {
	return &ButtonBuilder{
		button: &objects.Component{
			Type: objects.ComponentTypeButton,
		},
	}
}

func (b *ButtonBuilder) Style(style objects.ButtonStyle) *ButtonBuilder {
	b.button.Style = style
	return b
}

func (b *ButtonBuilder) Label(label string) *ButtonBuilder {
	b.button.Label = label
	return b
}

func (b *ButtonBuilder) Emoji(emoji *objects.Emoji) *ButtonBuilder {
	b.button.Emoji = emoji
	return b
}

func (b *ButtonBuilder) CustomID(ID string) *ButtonBuilder {
	b.button.CustomID = ID
	return b
}

func (b *ButtonBuilder) URL(url string) *ButtonBuilder {
	b.button.URL = url
	return b
}

func (b *ButtonBuilder) Disabled() *ButtonBuilder {
	b.button.Disabled = true
	return b
}

func (b *ButtonBuilder) Build() *objects.Component {
	return b.button
}
