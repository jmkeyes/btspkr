package btspkr

import "errors"

var (
	ErrTimeout        = errors.New("timeout")
	ErrReaderClosed   = errors.New("channel closed")
	ErrDeviceReset    = errors.New("device reset")
	ErrCommandError   = errors.New("command error")
	ErrCommandUnknown = errors.New("unknown command")
	ErrMemoryAccess   = errors.New("out of bounds")
	ErrMessageInvalid = errors.New("invalid message")
)
