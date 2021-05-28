package interactions

import "github.com/Postcord/objects"

type SelectMenuBuilder struct {
	selectMenu *objects.Component
}

func NewSelectMenu() *SelectMenuBuilder {
	return &SelectMenuBuilder{
		selectMenu: &objects.Component{
			Type: objects.ComponentTypeSelectMenu,
		},
	}
}

func (b *SelectMenuBuilder) CustomID(ID string) *SelectMenuBuilder {
	b.selectMenu.CustomID = ID
	return b
}

func (b *SelectMenuBuilder) Placeholder(placeholder string) *SelectMenuBuilder {
	b.selectMenu.Placeholder = placeholder
	return b
}

func (b *SelectMenuBuilder) MinValues(min int) *SelectMenuBuilder {
	b.selectMenu.MinValues = min
	return b
}

func (b *SelectMenuBuilder) MaxValues(max int) *SelectMenuBuilder {
	b.selectMenu.MaxValues = max
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
