package rest

type Config struct {
	Ratelimiter Ratelimiter
}

type Client struct {
	rateLimiter Ratelimiter
}

func New(config *Config) *Client {
	return &Client{
		rateLimiter: config.Ratelimiter,
	}
}

func (c *Client) request(method, path, contentType string, body []byte) (*DiscordResponse, error) {
	return c.rateLimiter.Request(method, path, contentType, body)
}
