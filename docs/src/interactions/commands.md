# Application Commands

At the very least, application commands are represented as structs that implement the following interface
```go
{{#include ../../../router/command.go:214:216}}
```

## Optional Interfaces

Optionally, you can implement any of the following interfaces for your command.  The purpose of each interface can be easily inferred from the name.

```go
{{#include ../../../router/parser.go:31:73}}
```

## Comments

There are a number of specially formatted comments you can specify on a command struct.  These are only available if you are relying wumpgoctl to generate the interface implementations for your commands.  This can save you a lot of time so you don't have to implement each interface by hand.

| Format | Purpose | Default |
|---------|---------|--------|
| @Name [Name] | The name of the application command | The struct name lowercased |
| @Description [Description] | The description of the application command (only valid, but required for chat input commands) | The name |
| @Name.es-MX [Name for locale] | The name for the given locale | Null |
| @Description.es-MX [Description for locale] | The description for the given locale | Null |
| @Type [ChatInput or Message or User] | The application command type | ChatInput |
| @Option.[OptionName].es-MX.Name | The option name for the given locale | Null |
| @Option.[OptionName].es-MX.Description | The option description for the given locale | Null |
| @DM [true or false] | If the command is available in DMs | false |
| @Permissions [PermissionNames] | Default permissions required to run the command | Null |

## Options

Options are declared as struct fields the types are inferred from the type of the field, and are further customized through struct field tags.

There are two tags used by wumpgo, the first is the `discord` tag, and deals with most customization you might need when declaring the options.  The tag contents should be a comma separated list of customizations, with the first part being the name for the field as it is sent by Discord.  Beyond that the fields can be in any order, and the possible options are as follows.

| Format | Description |
|--------|-------------|
| optional | Marks the field as optional.  wumpgo will ensure optional fields come after required. |
| autocomplete | Marks the field as autocomplete |
| description [description] | Sets the description of the option |
| minLength [min length] | Minimum length for a string argument |
| maxLength [max length] | Maximum length for a string argument |
| minValue [min value] | Minimum value for an int or float |
| maxValue [max value] | Maximum value for an int or float |
| channelTypes [semicolon separated list of channel types] | A semicolon separated list of allowed channel types |

The second tag used is the `choices` tag.  This tag is used to specify the possible choices for a given option.  It is simply a comma separated list of values.

## Example
Let's start with a quick example, then we'll go over what each part does.

```go
{{#include command_example.go}}
```

We start off with a go:generate comment.  This is required on any file that contains one or more slash commands that needs codegen.  wumpgoctl will only perform generation on files that directly call out to gen.

After this, we specify our codegen comments.  You'll see these match the options in the table above.
```go
{{#include command_example.go:11:14}}
```

Next we actually specify our command.  This is a ChatInput command, so `@Type` is not specified.  You'll notice struct annotations specify the option name, as well as whether or not they are optional.
```go
{{#include command_example.go:16:17}}
```

Lastly we supply the actual implementation for the command, this is required.  If this is missing, the command will fail to register at all.

Don't forget to run `go generate` if you're relying on code generation.