package rest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseRoute(t *testing.T) {
	cases := []struct {
		route string
		want  string
	}{
		{"/api/v9/gateway/bot", "gateway:bot"},
		{"/api/v9/channels/686053040863969305/messages/872142459847733259", "channels:686053040863969305:messages"},
	}

	for _, c := range cases {
		got := parseRoute(c.route)
		assert.Equal(t, c.want, got)
	}
}
