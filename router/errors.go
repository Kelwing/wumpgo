package router

import "fmt"

type errCommandNotFound struct{}

func (e *errCommandNotFound) Error() string {
	return "command not found"
}

type errCustomIDNotFound struct{}

func (e *errCustomIDNotFound) Error() string {
	return "this component is not registered"
}

type errArgumentMismatch struct{}

func (e *errArgumentMismatch) Error() string {
	return "arguments do not match expected"
}

type errInternalCommand struct {
	rec any
}

func (e *errInternalCommand) Error() string {
	return fmt.Sprintf("internal command error: %v", e.rec)
}

func (e *errInternalCommand) Recover() any {
	return e.rec
}

var (
	ErrCommandNotFound  *errCommandNotFound
	ErrArgumentMismatch *errArgumentMismatch
	ErrInternalCommand  *errInternalCommand
	ErrCustomIDNotFound *errCustomIDNotFound
)
