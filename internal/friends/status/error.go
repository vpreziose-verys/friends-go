package status

import (
	"errors"
	"fmt"
)

var (
	// ErrBad returned if nil or bad status
	ErrBad = errors.New("received bad status")

	//ErrBUID err returned if no buid set on (in/out) Check func
	ErrBUID = errors.New("received bad buid")

	// ErrBadJSONStatus error returned if json contained invalid string enum
	ErrBadJSONStatus = fmt.Errorf("error bad json status received")

	// ErrEncodingJSON used when converting status to json failed
	ErrEncodingJSON = errors.New("error encoding status to json")

	// ErrBadEnum error returned when an enum is past range preventing error
	ErrBadEnum = errors.New("error retreived bad status")

	// ErrMaxExtended err returned when SET presence exceeds max value
	ErrMaxExtended = fmt.Errorf("err too many extended game statuses: %d", MaxExtendedGameStatuses)

	// ErrMaxCustomBytes err returned when SET presence exceeds max value
	ErrMaxCustomBytes = fmt.Errorf("err custom status exceeds max value: %d", MaxCustomDataBytes)

	// ErrMaxGameStatusIDLength err returned when attribute exceeds max value
	ErrMaxGameStatusIDLength = fmt.Errorf("err game status id exceeds max value: %d", MaxGameStatusIDLength)

	// ErrMaxGameStatusDataLength err returned when attribute exceeds max value
	ErrMaxGameStatusDataLength = fmt.Errorf("err custom status exceeds max value: %d", MaxGameStatusDataLength)
)
