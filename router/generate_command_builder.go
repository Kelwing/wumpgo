//go:build ignore
// +build ignore

package main

import (
	"bytes"
	"io/ioutil"
	"strings"
	"text/template"
)

const start = `// Code generated by generate_command_builder.go; DO NOT EDIT.

package router

//go:generate go run generate_command_builder.go

import (
	"fmt"

	"github.com/kelwing/wumpgo/objects"
)

`

const choiceBuilder = `// {{ .TypeName }}AutoCompleteFunc is used to define the auto-complete function for {{ .ImportName }}.
// Note that the command context is a special case in that the response is not used.
type {{ .TypeName }}AutoCompleteFunc = func(*CommandRouterCtx) ([]{{ .ImportName }}, error)

// {{ .TypeName }}ChoiceBuilder is used to choose how this choice is handled.
// This can be nil, or it can pass to one of the functions. The first function adds static choices to the router. The second option adds an autocomplete function.
// Note that you cannot call both functions.
type {{ .TypeName }}ChoiceBuilder = func(addStaticOptions func([]{{ .ImportName }}), addAutocomplete func({{ .TypeName }}AutoCompleteFunc))

// {{ .TypeName }}StaticChoicesBuilder is used to create a shorthand for adding choices.
func {{ .TypeName }}StaticChoicesBuilder(choices []{{ .ImportName }}) {{ .TypeName }}ChoiceBuilder {
	return func(addStaticOptions func([]{{ .ImportName }}), _ func({{ .TypeName }}AutoCompleteFunc)) {
		addStaticOptions(choices)
	}
}

// {{ .TypeName }}AutoCompleteFuncBuilder is used to create a shorthand for adding a auto-complete function.
func {{ .TypeName }}AutoCompleteFuncBuilder(f {{ .TypeName }}AutoCompleteFunc) {{ .TypeName }}ChoiceBuilder {
	return func(_ func([]{{ .ImportName }}), addAutocomplete func({{ .TypeName }}AutoCompleteFunc)) {
		addAutocomplete(f)
	}
}

func (c *commandBuilder[T]) {{ .TypeName }}Option(name, description string, required bool, choiceBuilder {{ .TypeName }}ChoiceBuilder) T {
	var discordifiedChoices []objects.ApplicationCommandOptionChoice
	var f {{ .TypeName }}AutoCompleteFunc
	if choiceBuilder != nil {
		choiceBuilder(func(choices []{{ .TypeName }}Choice) {
			if f != nil {
				panic("cannot set both function and choice slice")
			}
			discordifiedChoices = make([]objects.ApplicationCommandOptionChoice, len(choices))
			for i, v := range choices {
				discordifiedChoices[i] = objects.ApplicationCommandOptionChoice{Name: v.Name, Value: v.Value}
			}
		}, func(autoCompleteFunc {{ .TypeName }}AutoCompleteFunc) {
			if discordifiedChoices != nil {
				panic("cannot set both function and choice slice")
			}
			f = autoCompleteFunc
		})
	}

	c.cmd.Options = append(c.cmd.Options, &objects.ApplicationCommandOption{
		OptionType:   objects.Type{{ .InteractionsTypeName }},
		Name:         name,
		Description:  description,
		Required:     required,
		Choices:      discordifiedChoices,
		Autocomplete: f != nil,
	})
	if f != nil {
		if c.cmd.autocomplete == nil {
			c.cmd.autocomplete = map[string]any{}
		}
		c.cmd.autocomplete[name] = f
	}
	return builderWrapify(c)
}`

const builderShared = `type {{ .Struct }} struct {
	*commandBuilder[{{ .BuilderType }}Builder]
}{{ if not .DoNotHook }}

func (c *commandBuilder[T]) {{ .BuilderType }}() {{ .BuilderType }}Builder {
	c.cmd.commandType = int({{ .CommandType }})
	return {{ .Struct }}{(*commandBuilder[{{ .BuilderType }}Builder])(c)}
}{{ end }}`

const builderWrapify = `func builderWrapify[T any](c *commandBuilder[T]) T {
	var ptr *T
	switch (any)(ptr).(type) {
		case *CommandBuilder:
			return (any)(c).(T){{range $val := .}}{{ if ne $val.BuilderType "CommandBuilder" }}
		case *{{ $val.BuilderType }}Builder:
			return (any)({{ $val.Struct }}{(any)(c).(*commandBuilder[{{ $val.BuilderType }}Builder])}).(T){{ end }}{{ end }}
		default:
			panic(fmt.Errorf("unknown handler: %T", ptr))
	}
}`

var choiceTypes = []struct {
	TypeName             string
	InteractionsTypeName string
	ImportName           string
}{
	{
		TypeName:             "String",
		InteractionsTypeName: "String",
		ImportName:           "StringChoice",
	},
	{
		TypeName:             "Int",
		InteractionsTypeName: "Integer",
		ImportName:           "IntChoice",
	},
	{
		TypeName:             "Double",
		InteractionsTypeName: "Number",
		ImportName:           "DoubleChoice",
	},
}

var builderTypes = []struct {
	Struct      string
	BuilderType string
	CommandType string
	DoNotHook   bool
}{
	{
		Struct:      "textCommandBuilder",
		BuilderType: "TextCommand",
		CommandType: "objects.CommandTypeChatInput",
	},
	{
		Struct:      "subcommandBuilder",
		BuilderType: "SubCommand",
		CommandType: "objects.CommandTypeChatInput",
		DoNotHook:   true,
	},
	{
		Struct:      "messageCommandBuilder",
		BuilderType: "MessageCommand",
		CommandType: "objects.CommandTypeMessage",
	},
	{
		Struct:      "userCommandBuilder",
		BuilderType: "UserCommand",
		CommandType: "objects.CommandTypeUser",
	},
}

func main() {
	file := start
	parts := make([]string, len(choiceTypes)+len(builderTypes)+1)
	t, err := template.New("_").Parse(choiceBuilder)
	if err != nil {
		panic(err)
	}
	for i, v := range choiceTypes {
		buf := &bytes.Buffer{}
		if err := t.Execute(buf, v); err != nil {
			panic(err)
		}
		parts[i] = buf.String()
	}
	t, err = template.New("_").Parse(builderShared)
	if err != nil {
		panic(err)
	}
	for i, v := range builderTypes {
		buf := &bytes.Buffer{}
		if err := t.Execute(buf, v); err != nil {
			panic(err)
		}
		parts[i+len(choiceTypes)] = buf.String()
	}
	t, err = template.New("_").Parse(builderWrapify)
	if err != nil {
		panic(err)
	}
	buf := &bytes.Buffer{}
	if err := t.Execute(buf, builderTypes); err != nil {
		panic(err)
	}
	parts[len(builderTypes)+len(choiceTypes)] = buf.String()
	file += strings.Join(parts, "\n\n") + "\n"
	if err := ioutil.WriteFile("command_builder_gen.go", []byte(file), 0666); err != nil {
		panic(err)
	}
}
