package relic

import "errors"

var (
	// ErrDisabled returned if config is not enabled
	ErrDisabled = errors.New("relic: disabled")

	// ErrBadCreds returned if config is missing name or key
	ErrBadCreds = errors.New("relic: invalid credentials")
)
