# Components

Componets and modals are routed using an http-style router.  We treat custom ID's similar to URL paths.  This allows you to do conditional routing based on parts of the custom ID.

## Path Parameters

Path parameters are specified using curly braces, for example: `/role/{role_id}`[^unsecure]. The values are automatically extracted, and made available as part of the context.

Obtaining the value of a path parameter is fairly simple, for example:
```go
func MyComponent(r router.ComponentResponder, ctx *router.ComponentContext) {
    roleID := ctx.Param("role_id")
    // Ensure the role ID is a valid option in this context
    ...
}
```

## Views

Components can be described in two ways: using views, or by manually building a components list.  Views make things much easier, but are not very customizable.  We will mostly be covering views in this guide, since manually specifying rows is done exactly as it is described in the Discord API docs.

A view can be created using `router.NewView()`.  You can add a component to a view using `View.Add(r router.Renderable)`.  Once created, the view can be attached to any interaction response.
```go
v := router.NewView().Add(
    router.NewButton("/ping").Label("Ping!"),
    router.NewButton("/another").Label("Another button"),
)

ctx.Content("Click a button!").View(v)
```

By default, views will render the components in order, filling out the component rows before adding another, up to 5.  Beyond 25 components, the view will be truncated.  The exception to this is when using a View in a modal, then it will only have 1 component per row, up to 5.  The only valid component type for modals is a Text Input component.

[^unsecure]: Be aware that custom IDs are inherintly not secure.  You cannot trust the contents of a custom ID.  Be sure you validate the data you get as part of a custom ID.  For example, if building a reaction role bot, validate the role is actually a valid choice before assigning it.