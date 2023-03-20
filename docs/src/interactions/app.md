# App

The interactions package contains a single struct: `App`.  The App struct provides a net/http compatible interface for processing Discord interactions and assigning handlers to deal with those interactions.  It's a very low level package, providing nothing more than authentication, and very basic routing by interaction type.  It provides no capability for routing commands or components beyond calling a single handler if the interaction itself is a command or component.

If you are looking for more advanced routing, like routing a specific slash command or message component to a specific function, take a look at the [router](./router.md).

This package is also only useful when working with the interaction webhook, not for gateway interactions, as those are already authenticated.

## Example

```go
app, err := interactions.New(*key)
if err != nil {
    log.Fatal().Err(err).Msg("failed to create interactions client")
}

app.CommandHandler(myCommandHandler)
app.ComponentHandler(myComponentHandler)

// The App is a net/http compatible handler, so we can just pass it to ListenAndServe
log.Err(http.ListenAndServe(":8080", app)).Msg("failed to listen")
```

## net/http Compatibility

Since the App is compatible with net/http, we can also use it as part of a bigger RESTful service, allowing you to integrate your Discord interactions handler with any net/http router that supports native net/http handlers, like [chi](https://go-chi.io).

This also makes it very easy to use with things like Google Cloud Functions, which use native net/http handlers.