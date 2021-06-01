package interactions

import "github.com/Postcord/objects"

// SelectMenuBuilder makes it easy to generate select menu components
type SelectMenuBuilder struct {
	selectMenu *objects.Component
}

// NewSelectMenuBuilder returns a new SelectMenuBuilder
func NewSelectMenuBuilder() *SelectMenuBuilder {
	return &SelectMenuBuilder{
		selectMenu: &objects.Component{
			Type: objects.ComponentTypeSelectMenu,
		},
	}
}

// CustomID sets the custom ID that's returned with the interaction event for the button.
// WARNING: this value can be manipulated in the client by malicious users, and as such, should not contain sensitive data, nor should the value be trusted.
func (b *SelectMenuBuilder) CustomID(ID string) *SelectMenuBuilder {
	b.selectMenu.CustomID = ID
	return b
}

// Placeholder 
func (b *SelectMenuBuilder) Placeholder(placeholder string) *SelectMenuBuilder {
	b.selectMenu.Placeholder = placeholder
	return b
}

func (b *SelectMenuBuilder) MinValues(min int) *SelectMenuBuilder {
	b.selectMenu.MinValues = &min
	return b
}

func (b *SelectMenuBuilder) MaxValues(max int) *SelectMenuBuilder {
	b.selectMenu.MaxValues = &max
	return b
}

func (b *SelectMenuBuilder) AddOptions(options ...*objects.SelectOptions) *SelectMenuBuilder {
	if b.selectMenu.Options == nil {
		b.selectMenu.Options = options
	} else {
		b.selectMenu.Options = append(b.selectMenu.Options, options...)
	}
	return b
}

func (b *SelectMenuBuilder) Build() *objects.Component {
	return b.selectMenu
}
