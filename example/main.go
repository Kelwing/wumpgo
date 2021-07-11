package main

import (
	"encoding/json"
	"fmt"
	"github.com/Postcord/interactions"
	"github.com/Postcord/objects"
	"github.com/Postcord/router"
	"os"
	"strconv"
)

func main() {
	// Create the response embed and component.
	createResponse := func(amount uint64) (*objects.Embed, *objects.Component) {
		return &objects.Embed{Description: "The value of this is " + strconv.FormatUint(amount, 10)}, &objects.Component{
			Type: objects.ComponentTypeActionRow,
			Components: []*objects.Component{
				{
					Type: objects.ComponentTypeButton,
					Label: "Add One",
					Style: objects.ButtonStylePrimary,
					CustomID: "/set/" + strconv.FormatUint(amount + 1, 10),
				},
			},
		}
	}

	// Creates the component router.
	componentRouter := &router.ComponentRouter{}
	componentRouter.RegisterButton("/set/:number", func(ctx *router.ComponentRouterCtx) error {
		number, err := strconv.ParseUint(ctx.Params["number"], 10, 64)
		if err != nil {
			// This would make the route specified invalid.
			return nil
		}
		embed, component := createResponse(number)
		ctx.Ephemeral().SetEmbed(embed).SetComponent(component)
		return nil
	})

	// Create the command router.
	commandRouter := &router.CommandRouter{}

	// Defines a single command.
	_, err := commandRouter.NewCommandBuilder("add").Description("A demo to show adding numbers.").
		Handler(func(ctx *router.CommandRouterCtx) error {
			embed, component := createResponse(0)
			ctx.Ephemeral().SetEmbed(embed).SetComponent(component)
			return nil
		}).
		Build()
	if err != nil {
		panic(err)
	}

	// Defines a command group.
	group := commandRouter.MustNewCommandGroup("group", "Defines a command group.")
	_, err = group.NewCommandBuilder("command").Description("Defines a description in a group.").
		Handler(func(ctx *router.CommandRouterCtx) error {
			ctx.SetContent("loge!")
			return nil
		}).
		Build()
	if err != nil {
		panic(err)
	}

	// Defines a sub-group.
	subgroups := commandRouter.MustNewCommandGroup("subgroups", "Defines command sub-groups.")
	group1 := subgroups.MustNewCommandGroup("group1", "Defines the first group.")
	group2 := subgroups.MustNewCommandGroup("group2", "Defines the second group.")
	_, err = group1.NewCommandBuilder("command").Description("Defines a description in a group.").
		Handler(func(ctx *router.CommandRouterCtx) error {
			ctx.SetContent("loge!")
			return nil
		}).
		Build()
	if err != nil {
		panic(err)
	}
	_, err = group2.NewCommandBuilder("command").Description("Defines a description in a group.").
		Handler(func(ctx *router.CommandRouterCtx) error {
			ctx.SetContent("loge!")
			return nil
		}).
		Build()
	if err != nil {
		panic(err)
	}

	// Create the interactions app.
	app, err := interactions.New(&interactions.Config{
		PublicKey:  os.Getenv("PUBLIC_KEY"),
		Token:      os.Getenv("TOKEN"),
	})
	if err != nil {
		panic(err)
	}

	if os.Getenv("RUN") == "" {
		j, err := json.Marshal(commandRouter.FormulateDiscordCommands())
		if err != nil {
			panic(err)
		}
		fmt.Println(string(j))
	}

	// Create the router builder.
	router.RouterLoader().ComponentRouter(componentRouter).CommandRouter(commandRouter).Build(app)

	// Create the interactions router.
	if err = app.Run(8000); err != nil {
		panic(err)
	}
}
