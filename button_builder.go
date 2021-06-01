package interactions

import "github.com/Postcord/objects"

// ButtonBuilder makes it easy to generate button components
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

// Style sets the style of the button
func (b *ButtonBuilder) Style(style objects.ButtonStyle) *ButtonBuilder {
	b.button.Style = style
	return b
}

// Label sets the text that appers on the button
func (b *ButtonBuilder) Label(label string) *ButtonBuilder {
	b.button.Label = label
	return b
}

// Emoji adds an emoji to the button
func (b *ButtonBuilder) Emoji(emoji *objects.Emoji) *ButtonBuilder {
	b.button.Emoji = emoji
	return b
}

// CustomID sets the custom ID that's returned with the interaction event for the button.
// WARNING: this value can be manipulated in the client by malicious users, and as such, should not contain sensitive data, nor should the value be trusted.
func (b *ButtonBuilder) CustomID(ID string) *ButtonBuilder {
	b.button.CustomID = ID
	return b
}

// URL sets a URL on the button.
// Style must be objects.ButtonStyleLink
func (b *ButtonBuilder) URL(url string) *ButtonBuilder {
	b.button.URL = url
	return b
}

// Disabled sets the button as a disabled button, can be edited later to enable it.
func (b *ButtonBuilder) Disabled() *ButtonBuilder {
	b.button.Disabled = true
	return b
}

// Returns the complete component
func (b *ButtonBuilder) Build() *objects.Component {
	return b.button
}
