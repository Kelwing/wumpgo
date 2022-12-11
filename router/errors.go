package router

type errCommandNotFound struct{}

func (e *errCommandNotFound) Error() string {
	return "command not found"
}

type errArgumentMismatch struct{}

func (e *errArgumentMismatch) Error() string {
	return "arguments do not match expected"
}

type errInternalCommand struct {
	rec any
}

func (e *errInternalCommand) Error() string {
	return "internal command error"
}

func (e *errInternalCommand) Recover() any {
	return e.rec
}

var (
	ErrCommandNotFound  *errCommandNotFound
	ErrArgumentMismatch *errArgumentMismatch
	ErrInternalCommand  *errInternalCommand
)
