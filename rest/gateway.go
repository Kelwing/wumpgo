package rest

import (
	"context"
	"net/http"

	"github.com/Kelwing/wumpgo/objects"
)

func (c *Client) Gateway(ctx context.Context) (*objects.Gateway, error) {
	var gateway objects.Gateway
	err := NewRequest().
		WithContext(ctx).
		Method(http.MethodGet).
		Path(GatewayFmt).
		ContentType(JsonContentType).
		Bind(&gateway).
		Send(c)

	return &gateway, err
}

func (c *Client) GatewayBot(ctx context.Context) (*objects.Gateway, error) {
	var gateway objects.Gateway
	err := NewRequest().
		Method(http.MethodGet).
		WithContext(ctx).
		Path(GatewayBotFmt).
		ContentType(JsonContentType).
		Bind(&gateway).
		Send(c)

	return &gateway, err
}
