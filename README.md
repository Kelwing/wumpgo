# Postcord Interactions

[![Go Reference](https://pkg.go.dev/badge/github.com/Postcord/interactions.svg)](https://pkg.go.dev/github.com/Postcord/interactions)

Interactions is a simple, batteries included HTTP interactions library for Discord.  It is designed to make is fast and easy to create a new Discord server integration using Discords new interactions system.

## Getting Started

Add Interactions to your project
```
go get -u github.com/Postcord/interactions
```

### Example

```go
package main

import (
    "log"
    "os"

    "github.com/Postcord/interactions"
    "github.com/Postcord/objects"
)

func main() {
	app, err := interactions.New(&interactions.Config{
		PublicKey: os.Getenv("DISCORD_PUBLIC_KEY"),
		Token:     "Bot " + os.Getenv("DISCORD_TOKEN"),
	})
	if err != nil {
		panic("failed to create interactions client")
	}

	app.CommandHandler(testHandler)

	log.Println("test-bot is now running")
	app.Run(1323)
}

func testHandler(ctx *objects.Interaction) *objects.InteractionResponse {
	return &objects.InteractionResponse{
		Type: objects.ResponseChannelMessageWithSource,
		Data: &objects.InteractionApplicationCommandCallbackData{
			Content: "Hello world!",
			Flags:   objects.ResponseFlagEphemeral,
		},
	}
}
```

## Command Routing
If you're looking for a more batteries included solution that includes command and interaction routing.  Check out our [router](https://github.com/Postcord/router) package.

## Documentation

Documentation is still a work in progress.
