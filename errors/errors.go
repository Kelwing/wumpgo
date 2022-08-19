package command_errors

import (
	"fmt"

	"wumpgo.dev/wumpgo/router"
)

// CommandError is an error that is returned when a command fails
type CommandError interface {
	Error() string
	Command() string
}

// ErrRuntime represents a generic runtime error
type ErrRuntime struct {
	cmd   *router.Command
	Value string
}

func NewErrRuntime(command *router.Command, value string) *ErrRuntime {
	return &ErrRuntime{
		cmd:   command,
		Value: value,
	}
}

func (e *ErrRuntime) Error() string {
	return fmt.Sprintf("Error (%s): %s", e.Command(), e.Value)
}

func (e *ErrRuntime) Command() string {
	return e.cmd.Name
}

// ErrCommandNotImplemented is an error that is returned when a command is not implemented
type ErrCommandNotImplemented struct {
	cmd *router.Command
}

func NewErrCommandNotImplemented(command *router.Command) *ErrCommandNotImplemented {
	return &ErrCommandNotImplemented{
		cmd: command,
	}
}

func (e *ErrCommandNotImplemented) Error() string {
	return fmt.Sprintf("Command %s not implemented", e.Command())
}

func (e *ErrCommandNotImplemented) Command() string {
	return e.cmd.Name
}

// ErrCommandNotAllowed is an error that is returned when the running user is not currently allowed to run a command
type ErrCommandNotAllowed struct {
	cmd *router.Command
}

func NewErrCommandNotAllowed(command *router.Command) *ErrCommandNotAllowed {
	return &ErrCommandNotAllowed{
		cmd: command,
	}
}

func (e *ErrCommandNotAllowed) Error() string {
	return fmt.Sprintf("You are not allowed to run %s at this time.", e.Command())
}

func (e *ErrCommandNotAllowed) Command() string {
	return e.cmd.Name
}
