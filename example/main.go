package main

import (
	"fmt"
	"github.com/Postcord/interactions"
	"github.com/Postcord/objects"
	"github.com/Postcord/router"
	"os"
	"strconv"
)

func builder() (*router.CommandRouter, *interactions.App, router.LoaderBuilder) {
	// Create the interactions app.
	app, err := interactions.New(&interactions.Config{
		PublicKey: os.Getenv("PUBLIC_KEY"),
		Token:     "Bot " + os.Getenv("TOKEN"),
	})
	if err != nil {
		panic(err)
	}

	// Create the response embed and component.
	createResponse := func(amount uint64) (*objects.Embed, []*objects.Component) {
		return &objects.Embed{Description: "The value of this is " + strconv.FormatUint(amount, 10)}, []*objects.Component{
			{
				Type:     objects.ComponentTypeButton,
				Label:    "Add One",
				Style:    objects.ButtonStylePrimary,
				CustomID: "/set/" + strconv.FormatUint(amount+1, 10),
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
		embed, row := createResponse(number)
		ctx.Ephemeral().SetEmbed(embed).AddComponentRow(row)
		return nil
	})

	// Create the command router.
	commandRouter := &router.CommandRouter{}

	// Defines a root level autocomplete.
	commandRouter.NewCommandBuilder("autocomplete").Description("A root level autocomplete.").
		DefaultPermission().
		StringOption("option", "The option that will be autocompleted.", true,
			router.StringAutoCompleteFuncBuilder(func(ctx *router.CommandRouterCtx) ([]router.StringChoice, error) {
				return []router.StringChoice{
					{
						Name:  "d",
						Value: "D",
					},
					{
						Name:  "e",
						Value: "E",
					},
					{
						Name:  "f",
						Value: "F",
					},
				}, nil
			}),
		).
		Handler(func(ctx *router.CommandRouterCtx) error {
			ctx.SetContent(ctx.Options["option"].(string))
			return nil
		}).
		MustBuild()

	// Defines a single command.
	_, err = commandRouter.NewCommandBuilder("add").Description("A demo to show adding numbers.").
		Handler(func(ctx *router.CommandRouterCtx) error {
			embed, row := createResponse(0)
			ctx.Ephemeral().SetEmbed(embed).AddComponentRow(row)
			return nil
		}).
		DefaultPermission().
		Build()
	if err != nil {
		panic(err)
	}

	// Defines a REST command.
	commandRouter.NewCommandBuilder("rest").Description("A command to hit the fuck out of rest.").
		Handler(func(ctx *router.CommandRouterCtx) error {
			_, _ = ctx.RESTClient.GetChannel(1)
			_, err := ctx.RESTClient.GetCurrentUser()
			if err != nil {
				return err
			}
			ctx.SetContent("Hello World!")
			return nil
		}).
		DefaultPermission().
		MustBuild()
	if err != nil {
		panic(err)
	}

	// Defines a command group.
	group := commandRouter.MustNewCommandGroup("group", "Defines a command group.", true)
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
	subgroups := commandRouter.MustNewCommandGroup("subgroups", "Defines command sub-groups.", true)
	group1 := subgroups.MustNewCommandGroup("group1", "Defines the first group.", true)
	group2 := subgroups.MustNewCommandGroup("group2", "Defines the second group.", true)
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

	// Defines a group level autocomplete.
	group2.NewCommandBuilder("autocomplete").Description("A group level autocomplete.").
		StringOption("option", "The option that will be autocompleted.", true,
			router.StringAutoCompleteFuncBuilder(func(ctx *router.CommandRouterCtx) ([]router.StringChoice, error) {
				return []router.StringChoice{
					{
						Name:  "a",
						Value: "A",
					},
					{
						Name:  "b",
						Value: "B",
					},
					{
						Name:  "c",
						Value: "C",
					},
				}, nil
			}),
		).
		Handler(func(ctx *router.CommandRouterCtx) error {
			ctx.SetContent(ctx.Options["option"].(string))
			return nil
		}).
		MustBuild()

	// Add a user target command.
	commandRouter.NewCommandBuilder("user-target").
		UserCommand().
		DefaultPermission().
		Handler(func(ctx *router.CommandRouterCtx, member *objects.GuildMember) error {
			ctx.SetContent("You clicked " + member.User.Username)
			return nil
		}).
		MustBuild()

	// Add a message target command.
	commandRouter.NewCommandBuilder("message-target").
		MessageCommand().
		DefaultPermission().
		Handler(func(ctx *router.CommandRouterCtx, message *objects.Message) error {
			ctx.SetContent("The message was made by " + message.Author.Username)
			return nil
		}).
		MustBuild()

	return commandRouter, app, router.RouterLoader().ComponentRouter(componentRouter).CommandRouter(commandRouter).Build(app)
}

func main() {
	// Dump the Discord commands if specified.
	commandBuilder, app, _ := builder()
	if os.Getenv("DUMP") == "1" {
		commands := commandBuilder.FormulateDiscordCommands()
		me, err := app.Rest().GetCurrentUser()
		if err != nil {
			panic(err)
		}
		if _, err := app.Rest().BulkOverwriteGlobalCommands(me.ID, commands); err != nil {
			panic(err)
		}
		fmt.Println("Commands dumped.")
		return
	}

	// Create the interactions router.
	if err := app.Run(8000); err != nil {
		panic(err)
	}
}
