package rest

import (
	"net/http"

	"github.com/Postcord/objects"
)

func (c *Client) Gateway() (*objects.Gateway, error) {
	var gateway objects.Gateway
	err := NewRequest().
		Method(http.MethodGet).
		Path(GatewayFmt).
		ContentType(JsonContentType).
		Expect(http.StatusOK).
		Bind(&gateway).
		Send(c)

	return &gateway, err
}

func (c *Client) GatewayBot() (*objects.Gateway, error) {
	var gateway objects.Gateway
	err := NewRequest().
		Method(http.MethodGet).
		Path(GatewayBotFmt).
		ContentType(JsonContentType).
		Expect(http.StatusOK).
		Bind(&gateway).
		Send(c)

	return &gateway, err
}
