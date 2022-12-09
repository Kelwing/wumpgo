package newrouter

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	"wumpgo.dev/wumpgo/objects"
)

type TestRootCommand struct {
	TestSubCommandGroup
}

type TestSubCommandGroup struct {
	TestCommand
	TestCommandTwo
}

type TestCommand struct {
	Message       string       `discord:"message,minLength:2,maxLength:2000"`
	SomeNum       int64        `discord:"somenum,minValue:5,maxValue:10"`
	DoThing       bool         `discord:"dothing"`
	AnotherNumber float64      `discord:"anothernumber"`
	AnOptional    string       `discord:"anoptional,optional"`
	SomeUser      objects.User `discord:"user"`
}

func (t TestCommand) CommandName() string {
	return "mytestcommand"
}

type TestCommandTwo struct {
	Channel objects.Channel `discord:"channel"`
}

func (t TestCommandTwo) Description() string {
	return "testing"
}

type TestCommandThree struct {
	Channel *objects.Channel `discord:"channel"`
}

type BadCommand struct {
	TestCommandTwo
	Message string `discord:"message"`
}

func TestParseCommand(t *testing.T) {
	p := NewParser()
	x := reflect.ValueOf(TestRootCommand{})
	cmd, err := p.parseCommand(x)
	require.NoError(t, err)
	require.Equal(t, "mytestcommand", cmd.Name)
	require.Equal(t, 1, len(cmd.Options))
	require.Equal(t, 6, len(cmd.Options[0].Options[0].Options))
	require.False(t, cmd.Options[0].Options[0].Options[5].Required)
	require.Equal(t, objects.TypeSubCommandGroup, cmd.Options[0].OptionType)
	require.Equal(t, objects.TypeSubCommand, cmd.Options[0].Options[0].OptionType)
}

func TestParseCommand2(t *testing.T) {
	p := NewParser()
	x := reflect.ValueOf(TestCommandTwo{})

	cmd, err := p.parseCommand(x)
	require.NoError(t, err)
	require.Equal(t, "testcommandtwo", cmd.Name)
	require.Equal(t, 1, len(cmd.Options))
	require.Equal(t, "testing", cmd.Description)
}

func TestParseCommand3(t *testing.T) {
	p := NewParser()
	x := reflect.ValueOf(TestCommandThree{})

	cmd, err := p.parseCommand(x)

	require.NoError(t, err)
	require.Equal(t, "testcommandthree", cmd.Name)
	require.Equal(t, 1, len(cmd.Options))
	require.Equal(t, objects.TypeChannel, cmd.Options[0].OptionType)
}

func TestBadCommand(t *testing.T) {
	p := NewParser()
	x := reflect.ValueOf(BadCommand{})
	_, err := p.parseCommand(x)
	require.Error(t, err)
}
