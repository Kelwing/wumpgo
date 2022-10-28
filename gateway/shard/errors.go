package shard

import gatewayerrors "wumpgo.dev/wumpgo/gateway/internal/gateway_errors"

var (
	ErrGeneric                  = errGeneric("")
	ErrReconnect                = errReconnect()
	ErrInvalidSession           = errInvalidSession()
	ErrSessionStartLimitReached = errSessionStartLimitReached(0)
)

func errGeneric(message string) error { return &gatewayerrors.ErrGenericError{Message: message} }
func errReconnect() error             { return &gatewayerrors.ErrReconnect{} }
func errInvalidSession() error        { return &gatewayerrors.ErrInvalidSession{} }
func errSessionStartLimitReached(resetAfter int) error {
	return &gatewayerrors.ErrSessionStartLimitReached{ResetAfter: resetAfter}
}
