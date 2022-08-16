package objects

import "strings"

var _ Mentionable = (*Emoji)(nil)
var _ SnowflakeObject = (*Emoji)(nil)

type Emoji struct {
	DiscordBaseObject
	Name          string      `json:"name,omitempty"`
	Roles         []Snowflake `json:"roles,omitempty"`
	User          *User       `json:"user,omitempty"`
	RequireColons bool        `json:"require_colons,omitempty"`
	Managed       bool        `json:"managed,omitempty"`
	Animated      bool        `json:"animated,omitempty"`
	Available     bool        `json:"available,omitempty"`
}

func (e *Emoji) Mention() string {
	if e.ID != 0 {
		var b strings.Builder
		b.WriteRune('<')
		if e.Animated {
			b.WriteRune('a')
		}
		b.WriteRune(':')
		b.WriteString(e.Name)
		b.WriteRune(':')
		b.WriteString(e.ID.String())
		b.WriteRune('>')
		return b.String()
	} else {
		return e.Name
	}
}

func (e *Emoji) String() string {
	return e.Mention()
}

type Reaction struct {
	Count int   `json:"count"`
	Me    bool  `json:"me"`
	Emoji Emoji `json:"emoji"`
}
