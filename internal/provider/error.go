package provider

import "errors"

var (
	// ErrConfigNil err ...
	ErrConfigNil = errors.New("error nil config")

	// ErrConfigNotPtr error returned if provider config on open is not a pointer
	ErrConfigNotPtr = errors.New("error config not ptr")

	// ErrConfigInvalid error returned if config type is invalid
	ErrConfigInvalid = errors.New("error invalid config type")

	// ErrNotImplemented error returned if method not implemented
	ErrNotImplemented = errors.New("error not implemented")

	// ErrInvalidBUID error returned if buid is invalid
	ErrInvalidBUID = errors.New("error invalid buid")

	// ErrNilClient error when provider client is nil
	ErrNilClient = errors.New("error nil client")
)
