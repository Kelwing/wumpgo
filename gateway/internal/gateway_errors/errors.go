package gatewayerrors

type ErrGenericError struct {
	Message string
}

func (e *ErrGenericError) Error() string {
	return e.Message
}

type ErrReconnect struct{}

func (e *ErrReconnect) Error() string {
	return "reconnect"
}

type ErrInvalidSession struct{}

func (e *ErrInvalidSession) Error() string {
	return "invalid session"
}

type ErrSessionStartLimitReached struct{ ResetAfter int }

func (e *ErrSessionStartLimitReached) Error() string {
	return "session start limit reached"
}
