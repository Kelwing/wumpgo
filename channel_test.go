package rest

import (
	"fmt"
	"github.com/Postcord/objects"
	"os"
	"runtime"
	"testing"
)

func TestClient_GetChannel(t *testing.T) {
	client := New(&Config{
		Ratelimiter: NewMemoryRatelimiter(&MemoryConf{
			Authorization: fmt.Sprintf("Bot %s", os.Getenv("TOKEN")),
			MaxRetries:    3,
			UserAgent:     fmt.Sprintf("Postcord/1.0 %s (%s)", runtime.GOOS, runtime.GOARCH),
		}),
	})

	c, err := client.GetChannel(objects.Snowflake(484093378993192971))
	if err != nil {
		t.Error("failed to get channel:", err)
	}

	t.Log(c)
}

func TestClient_GetChannelMessages(t *testing.T) {
	client := New(&Config{
		Ratelimiter: NewMemoryRatelimiter(&MemoryConf{
			Authorization: fmt.Sprintf("Bot %s", os.Getenv("TOKEN")),
			MaxRetries:    3,
			UserAgent:     fmt.Sprintf("Postcord/1.0 %s (%s)", runtime.GOOS, runtime.GOARCH),
		}),
	})

	c, err := client.GetChannelMessages(objects.Snowflake(484093378993192971), nil)
	if err != nil {
		t.Error("failed to get channel messages:", err)
	}

	t.Log(c)
}
