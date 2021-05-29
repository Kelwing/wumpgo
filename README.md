# Postcord Interactions

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
    })
    if err != nil {
        panic("failed to create interactions client")
    }

    app.AddCommand(&objects.ApplicationCommand{
        Name:        "test",
        Description: "Test things",
    }, testHandler)

    log.Println("test-bot is now running")
    app.Run(1323)
}

func testHandler(ctx *interactions.CommandCtx) {
    ctx.SetContent("Hello world!").Ephemeral()
}
```

## Documentation

Documentation is still a work in progress.