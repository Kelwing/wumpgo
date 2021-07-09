# router

A high level Postcord interactions router. This code implements a modified version of the binary tree from httprouter (see `tree.go`). This is licensed under the BSD license.

## How does this work?

This package consists of 2 routers (a components router and a commands router), and a function to bootstrap the routers and set any global configurations which affect them both. Each router has been designed to work well for its specific usecase.

## Components Router
Routing components was traditionally quite complicated. You could use a standard map, but then you would lose the ability to track state. You could put the state in the string, but then you would either need to make a cache lookup, or you would need to do complicated and expensive string splitting.

Postcord gets around this by using a fork of the tree that httprouter/fasthttprouter uses for our component routes. This means that we have the ability to track state whilst making the application performant.

To begin with, we will go ahead and define the router somewhere where any of the components that wish to access it can see.

```go
var componentRouter = &router.ComponentRouter{}
```

Note that we can directly initialise this object. There's no new function. From here, we can go ahead and add functionality to this router.

There are 2 functions we can use to add components:
- `RegisterButton(route string, cb ButtonFunc)`: The job of this is to allow you to register button components. The signature for `ButtonFunc` takes `*ComponentRouterCtx` as the first argument and returns an `error`. See [creating responses with the context](#creating-responses-with-the-context) to see how you would make a response with the context handed down for this button.
- `RegisterSelectMenu(route string, cb SelectMenuFunc)`: The job of this is to allow you to register select menu components. The signature for `SelectMenuFunc` takes `*ComponentRouterCtx` as the first argument, and `[]string` as the second (this will contian the choices that the user made). It will then return an error. See [creating responses with the context](#creating-responses-with-the-context) to see how you would make a response with the context handed down for this button.

It is important to note that in both instances, we can use parameters in the path. This is awesome for inputting user data, however it is important to validate the data! Do not trust this input as it can be manipulated by users! As with httprouter, this can be done with `:paramName` in the path. For example, we could go ahead and register the following button:
```go
componentRouter.RegisterButton("/name/:name", func(ctx *router.ComponentRouterCtx) error {
	ctx.SetEmbed(&objects.Embed{Description: "my name " + ctx.Params["name"]})
	return nil
})
```

From here, we could go ahead and import this path in a component with the custom ID `/name/Jeff` and when clicked it would reply with an embed saying `my name Jeff`.

## Commands Router
One thing that is even more difficult to route than components is commands. Routing through sub-commands requires a lot of mind bending iteration, but don't worry, we have your back and have created a high level commands router too!

The same as with interactions, you will firstly want to go ahead and make the commands router:
```go
var commandRouter = &router.CommandRouter{}
```

Now we have a commands router created, we can now add to it.

To register a group, both `CommandRouter` and `CommandGroup` (to allow for one sub-group as per [Discord's nesting rules](https://discord.com/developers/docs/interactions/slash-commands#nested-subcommands-and-groups)) have the `NewCommandGroup` and `MustNewCommandGroup` functions. This will allow you to make a command group, which will then allow you to create commands/groups to be nested inside of it.

To register a command, we can go ahead and call `NewCommandBuilder` on either a router or group with the argument being the name of the command. From here, we can set the following options:

TODO

### Options
TODO

## Creating Responses with the Context
TODO

## Router Loader

So we went ahead and created both our routers. Awesome! But we need to get these into the application somehow. Don't worry, we thought of a very elegant way to do this. The `RouterLoader` function will create a builder which we can use to go ahead to build and execute the loader. We can use the following options in our chain:

- `ComponentRouter(*ComponentRouter) LoaderBuilder`: Adds the component router specified into the loader.
- `CommandRouter(*CommandRouter) LoaderBuilder`: Adds the commands router specified into the loader.
- `ErrorHandler(func(error) *objects.InteractionResponse) LoaderBuilder`: See [error handling](#error-handling).
- `AllowedMentions(*objects.AllowedMentions) LoaderBuilder`: See [allowed mentions](#allowed-mentions).

At the end of this, just call `Build` with your interactions application (you probably want a `*interactions.App` from Postcord/interactions). This will automatically inject the routers into your application and build them with the appropriate allowed mentions configuration.

### Error Handling
So how does error handling work? Error handling is done at a global scope with an error handler that takes a error parameter and returns a `*objects.InteractionResponse`. This can be used to write your own error handling code for actions. Note that there are a few errors that are dispatched by this codebase, and these are documented in the godoc for this project.

### Allowed Mentions
Allowed mention configurations can be set on a command, group, and global scope. Note that it takes affect in that order, so a command level allowed mentions configuration will override a global one.
